package cachekv

import (
	"time"

	"github.com/tidwall/btree"

	"github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
)

// Store wraps an in-memory cache around an underlying types.KVStore.
type WriteBuffer struct {
	// this is a map from key -> value, string -> []byte
	// if value is nil, means its deleted.
	buffer *btree.Map[string, []byte]
}

// NewStore creates a new Store object
func NewWriteBuffer() *WriteBuffer {
	return &WriteBuffer{
		buffer: btree.NewMap[string, []byte](16),
	}
}

func (wb *WriteBuffer) GetString(key string) (value []byte, found bool) {
	value, found = wb.bufferedWrites.Get(key)
	return value, found
}

func (wb *WriteBuffer) SetString(key string, value []byte) {
	wb.buffer.Set(key, value)
}

func (wb *WriteBuffer) HasString(key string) bool {
	_, found := wb.buffer.Get(key)
	return found
}

func (wb *WriteBuffer) DeleteString(key string) {
	wb.buffer.Delete(key)
}

// Implements Cachetypes.KVStore.
func (wb *WriteBuffer) WriteTo(parent types.KVStore) {
	defer telemetry.MeasureSince(time.Now(), "store", "WriteBuffer", "write")

	// TODO: Consider allowing usage of Batch, which would allow the write to
	// at least happen atomically.
	wb.buffer.Ascend("", func(key string, value []byte) bool {
		if value == nil {
			// We use []byte(key) instead of conv.UnsafeStrToBytes because we cannot
			// be sure if the underlying store might do a save with the byteslice or
			// not. Once we get confirmation that .Delete is guaranteed not to
			// save the byteslice, then we can assume only a read-only copy is sufficient.
			parent.Delete([]byte(key))
			return false
		}
		// It already exists in the parent, hence delete it.
		parent.Set([]byte(key), value)
		return false
	})

	// TODO: Investigate getting better memory re-use in this.
	// Its a waste that we have to go re-allocate everything later :(
	wb.buffer.Clear()
}

//----------------------------------------
// Iteration

// Iterator implements types.KVStore.
func (wb *WriteBuffer) Iterator(start, end []byte) types.Iterator {
	var a types.Iterator
	return a
}

// ReverseIterator implements types.KVStore.
func (wb *WriteBuffer) ReverseIterator(start, end []byte) types.Iterator {
	var a types.Iterator
	return a
}
