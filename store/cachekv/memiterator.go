package cachekv

import (
	"bytes"
	"errors"
  "sync"

  "github.com/google/btree"

)

// Iterates over iterKVCache items.
// if key is nil, means it was deleted.
// Implements Iterator.
type memIterator struct {
  mtx *sync.Mutex
	start, end []byte
  itemch chan *cValue
  item *cValue
	ascending  bool
}

func newMemIterator(start, end []byte, cache *btree.BTree, ascending bool) *memIterator {
  mi := &memIterator {
    mtx: new(sync.Mutex),
    start: start,
    end: end,
    itemch: make(chan *cValue),
    item: nil,
    ascending: ascending,
  }

  startKey := newKey(start)
  endKey := newKey(end)
  if end == nil {
    endKey = newInfinityKey()
  }

  if ascending {
    go func() {
      cache.AscendGreaterOrEqual(
        startKey,
        func(i btree.Item) bool {
          value, ok := i.(*cValue)
          if !ok {
            panic("invalid access")
          }
          if end != nil {
            if bytes.Compare(value.key, end) != -1 {
              return false
            }
          }
          mi.itemch <- value
          return true
        },
      )
      mi.itemch <- nil
    }()
    mi.item = <- mi.itemch
  } else {
    go func() {
      cache.DescendLessOrEqual(
        endKey,
        func(i btree.Item) bool {
          value, ok := i.(*cValue)
          if !ok {
            panic("invalid access")
          }
          if bytes.Compare(value.key, start) == -1 {
            return false
          }
          mi.itemch <- value
          return true
        },
      )
      mi.itemch <- nil
    }()
    mi.item = <- mi.itemch
    if bytes.Equal(mi.item.key, end) {
      mi.Next()
    }
  }

  return mi
}

func (mi *memIterator) Domain() ([]byte, []byte) {
	return mi.start, mi.end
}

func (mi *memIterator) Valid() bool {
  return mi.item != nil
}

func (mi *memIterator) assertValid() {
	if err := mi.Error(); err != nil {
		panic(err)
	}
}

func (mi *memIterator) Next() {
  mi.mtx.Lock()
  defer mi.mtx.Unlock()

  mi.assertValid()

  mi.item = <-mi.itemch
  if mi.item == nil {
    mi.Close()
  }
}

func (mi *memIterator) Key() []byte {
  mi.assertValid()

  return mi.item.key
}

func (mi *memIterator) Value() []byte {
  mi.assertValid()

  return mi.item.value
}

func (mi *memIterator) Close() error {
	mi.start = nil
	mi.end = nil
  mi.item = nil
	mi.itemch = nil

	return nil
}

// Error returns an error if the memIterator is invalid defined by the Valid
// method.
func (mi *memIterator) Error() error {
	if !mi.Valid() {
		return errors.New("invalid memIterator")
	}

	return nil
}
