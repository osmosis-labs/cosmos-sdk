package cachekv

import (
	"time"

	"github.com/tidwall/btree"

	"github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
)

// Store wraps an in-memory cache around an underlying types.KVStore.
type Store struct {
	// this is a map from key -> value, string -> []byte
	// if value is nil, means its deleted.
	bufferedWrites *btree.Map[string, []byte]
	parent         types.KVStore
}

// NewStore creates a new Store object
func NewStore(parent types.KVStore) *Store {
	return &Store{
		bufferedWrites: btree.NewMap[string, []byte](16),
		parent:         parent,
	}
}

// GetStoreType implements Store.
func (store *Store) GetStoreType() types.StoreType {
	return store.parent.GetStoreType()
}

func (store *Store) GetString(key string) (value []byte) {
	value, _ = store.bufferedWrites.Get(key)
	return value
}

func (store *Store) SetString(key string, value []byte) {
	store.bufferedWrites.Set(key, value)
}

func (store *Store) HasString(key string) bool {
	_, found := store.bufferedWrites.Get(key)
	return found
}

func (store *Store) DeleteString(key string) {
	store.bufferedWrites.Delete(key)
}

// Implements Cachetypes.KVStore.
func (store *Store) Write() {
	defer telemetry.MeasureSince(time.Now(), "store", "cachekv", "write")

	// TODO: Consider allowing usage of Batch, which would allow the write to
	// at least happen atomically.
	store.bufferedWrites.Ascend("", func(key string, value []byte) bool {
		if value == nil {
			// We use []byte(key) instead of conv.UnsafeStrToBytes because we cannot
			// be sure if the underlying store might do a save with the byteslice or
			// not. Once we get confirmation that .Delete is guaranteed not to
			// save the byteslice, then we can assume only a read-only copy is sufficient.
			store.parent.Delete([]byte(key))
			return false
		}
		// It already exists in the parent, hence delete it.
		store.parent.Set([]byte(key), value)
		return false
	})

	// TODO: Investigate getting better memory re-use in this.
	// Its a waste that we have to go re-allocate everything later :(
	store.bufferedWrites.Clear()
}

//----------------------------------------
// Iteration

// Iterator implements types.KVStore.
func (store *Store) Iterator(start, end []byte) types.Iterator {
	var a types.Iterator
	return a
}

// ReverseIterator implements types.KVStore.
func (store *Store) ReverseIterator(start, end []byte) types.Iterator {
	var a types.Iterator
	return a
}

// func (store *Store) iterator(start, end []byte, ascending bool) types.Iterator {
// 	expectedId := nextIteratorId
// 	fmt.Printf("calling lock for iterator creation, cur next id %v\n", expectedId)
// 	fmt.Printf("current store mtx address %p\n", &store.mtx)
// 	fmt.Printf("current memdb address %p\n", store.sortedCache)
// 	store.mtx.Lock()
// 	defer store.mtx.Unlock()
// 	fmt.Printf("got past iterator lock with expected id %v, next iterator id %v\n", expectedId, nextIteratorId)

// 	var parent, cache types.Iterator

// 	fmt.Println("Creating a parent iterator")
// 	if ascending {
// 		parent = store.parent.Iterator(start, end)
// 	} else {
// 		parent = store.parent.ReverseIterator(start, end)
// 	}
// 	fmt.Println("Finished creating a parent iterator")

// 	store.dirtyItems(start, end)

// 	fmt.Println("Creating a mem iterator")
// 	cache = newMemIterator(start, end, store.sortedCache, store.deleted, ascending)
// 	fmt.Println("Finished creating a mem iterator")

// 	iter := newCacheMergeIterator(parent, cache, ascending)
// 	return iter
// }
