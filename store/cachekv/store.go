package cachekv

import (
	"bytes"
	"io"
	"reflect"
	"sync"
	"time"
	"unsafe"

  "github.com/google/btree"

	"github.com/cosmos/cosmos-sdk/store/tracekv"
	"github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
)

// If value is nil but deleted is false, it means the parent doesn't have the
// key.  (No need to delete upon Write())
type cValue struct {
  key []byte
	value   []byte
	deleted bool

  infinity bool
}

var _ btree.Item = &cValue{}

func newKey(key []byte) btree.Item {
  return &cValue{key: key}
}

func newInfinityKey() btree.Item {
  return &cValue{infinity: true}
}

func (value *cValue) Less(i btree.Item) bool {
  value2, ok := i.(*cValue)
  if !ok {
    panic("invalid access")
  }
  if value.infinity {
    return false
  }
  if value2.infinity {
    return true
  }
  return bytes.Compare(value.key, value2.key) == -1
}

// Store wraps an in-memory cache around an underlying types.KVStore.
type Store struct {
	mtx     sync.Mutex
	cache   *btree.BTree
  parent  types.KVStore
}

var _ types.CacheKVStore = (*Store)(nil)

// NewStore creates a new Store object
func NewStore(parent types.KVStore) *Store {
	return &Store{
    cache: btree.New(2),
    parent:        parent,
	}
}

// GetStoreType implements Store.
func (store *Store) GetStoreType() types.StoreType {
	return store.parent.GetStoreType()
}

// Get implements types.KVStore.
func (store *Store) Get(key []byte) (value []byte) {
	store.mtx.Lock()
	defer store.mtx.Unlock()
	defer telemetry.MeasureSince(time.Now(), "store", "cachekv", "get")

	types.AssertValidKey(key)

  cacheValue, ok := store.cache.Get(newKey(key)).(*cValue)
	if !ok {
		value = store.parent.Get(key)
		store.setCacheValue(key, value, false)
	} else {
		value = cacheValue.value
	}

	return value
}

// Set implements types.KVStore.
func (store *Store) Set(key []byte, value []byte) {
	store.mtx.Lock()
	defer store.mtx.Unlock()
	defer telemetry.MeasureSince(time.Now(), "store", "cachekv", "set")

	types.AssertValidKey(key)
	types.AssertValidValue(value)

	store.setCacheValue(key, value, false)
}

// Has implements types.KVStore.
func (store *Store) Has(key []byte) bool {
	value := store.Get(key)
	return value != nil
}

// Delete implements types.KVStore.
func (store *Store) Delete(key []byte) {
	store.mtx.Lock()
	defer store.mtx.Unlock()
	defer telemetry.MeasureSince(time.Now(), "store", "cachekv", "delete")

	types.AssertValidKey(key)
	store.setCacheValue(key, nil, true)
}

// Implements Cachetypes.KVStore.
func (store *Store) Write() {
	store.mtx.Lock()
	defer store.mtx.Unlock()
	defer telemetry.MeasureSince(time.Now(), "store", "cachekv", "write")

	// We need a copy of all of the keys.
	// Not the best, but probably not a bottleneck depending.
  values := make([]*cValue, 0, store.cache.Len())
  store.cache.Ascend(func(i btree.Item) bool {
    value, ok := i.(*cValue)
    if !ok {
      panic("invalid access")
    }
    values = append(values, value)
    return true
  })

	// TODO: Consider allowing usage of Batch, which would allow the write to
	// at least happen atomically.
	for _, cacheValue := range values {
		switch {
		case cacheValue.deleted:
			store.parent.Delete([]byte(cacheValue.key))
		case cacheValue.value == nil:
			// Skip, it already doesn't exist in parent.
		default:
			store.parent.Set([]byte(cacheValue.key), cacheValue.value)
		}
	}

	// Clear the cache
  store.cache.Clear(false)
}

// CacheWrap implements CacheWrapper.
func (store *Store) CacheWrap() types.CacheWrap {
	return NewStore(store)
}

// CacheWrapWithTrace implements the CacheWrapper interface.
func (store *Store) CacheWrapWithTrace(w io.Writer, tc types.TraceContext) types.CacheWrap {
	return NewStore(tracekv.NewStore(store, w, tc))
}

//----------------------------------------
// Iteration

// Iterator implements types.KVStore.
func (store *Store) Iterator(start, end []byte) types.Iterator {
	return store.iterator(start, end, true)
}

// ReverseIterator implements types.KVStore.
func (store *Store) ReverseIterator(start, end []byte) types.Iterator {
	return store.iterator(start, end, false)
}

func (store *Store) iterator(start, end []byte, ascending bool) types.Iterator {
	store.mtx.Lock()
	defer store.mtx.Unlock()

	var parent, cache types.Iterator

	if ascending {
		parent = store.parent.Iterator(start, end)
	} else {
		parent = store.parent.ReverseIterator(start, end)
	}

	cache = newMemIterator(start, end, store.cache.Clone(), ascending)

	return newCacheMergeIterator(parent, cache, ascending)
}

// strToByte is meant to make a zero allocation conversion
// from string -> []byte to speed up operations, it is not meant
// to be used generally, but for a specific pattern to check for available
// keys within a domain.
func strToByte(s string) []byte {
	var b []byte
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	hdr.Cap = len(s)
	hdr.Len = len(s)
	hdr.Data = (*reflect.StringHeader)(unsafe.Pointer(&s)).Data
	return b
}

// byteSliceToStr is meant to make a zero allocation conversion
// from []byte -> string to speed up operations, it is not meant
// to be used generally, but for a specific pattern to delete keys
// from a map.
func byteSliceToStr(b []byte) string {
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&b))
	return *(*string)(unsafe.Pointer(hdr))
}

//----------------------------------------
// etc

// Only entrypoint to mutate store.cache.
func (store *Store) setCacheValue(key, value []byte, deleted bool) {
	store.cache.ReplaceOrInsert(&cValue{
    key: key,
		value:   value,
		deleted: deleted,
	})
}
