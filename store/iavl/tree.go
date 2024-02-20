package iavl

import (
	"fmt"

	dbm "github.com/cometbft/cometbft-db"
	"github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/iavl"
)

var (
	_ Tree = (*immutableTree)(nil)
	_ Tree = (*mutableTree)(nil)
)

type (
	// Tree defines an interface that both mutable and immutable IAVL trees
	// must implement. For mutable IAVL trees, the interface is directly
	// implemented by an iavl.MutableTree. For an immutable IAVL tree, a wrapper
	// must be made.
	Tree interface {
		Has(key []byte) (bool, error)
		Get(key []byte) ([]byte, error)
		Set(key, value []byte) (bool, error)
		Remove(key []byte) ([]byte, bool, error)
		SaveVersion() ([]byte, int64, error)
		DeleteVersionsTo(version int64) error
		Version() int64
		Hash() []byte
		VersionExists(version int64) bool
		GetVersioned(key []byte, version int64) ([]byte, error)
		GetImmutable(version int64) (*iavl.ImmutableTree, error)
		SetInitialVersion(version uint64)
		Iterator(start, end []byte, ascending bool) (types.Iterator, error)
		LoadVersionForOverwriting(targetVersion int64) error
	}

	// immutableTree is a simple wrapper around a reference to an iavl.ImmutableTree
	// that implements the Tree interface. It should only be used for querying
	// and iteration, specifically at previous heights.
	immutableTree struct {
		*iavl.ImmutableTree
	}

	mutableTree struct {
		*iavl.MutableTree
	}
)

func (it *immutableTree) Set(_, _ []byte) (bool, error) {
	panic("cannot call 'Set' on an immutable IAVL tree")
}

func (it *immutableTree) Remove(_ []byte) ([]byte, bool, error) {
	panic("cannot call 'Remove' on an immutable IAVL tree")
}

func (it *immutableTree) SaveVersion() ([]byte, int64, error) {
	panic("cannot call 'SaveVersion' on an immutable IAVL tree")
}

func (it *immutableTree) DeleteVersionsTo(_ int64) error {
	panic("cannot call 'DeleteVersionsTo' on an immutable IAVL tree")
}

func (it *immutableTree) SetInitialVersion(_ uint64) {
	panic("cannot call 'SetInitialVersion' on an immutable IAVL tree")
}

func (it *immutableTree) VersionExists(version int64) bool {
	return it.Version() == version
}

func (it *immutableTree) GetVersioned(key []byte, version int64) ([]byte, error) {
	if it.Version() != version {
		return nil, fmt.Errorf("version mismatch on immutable IAVL tree; got: %d, expected: %d", version, it.Version())
	}

	return it.Get(key)
}

func (it *immutableTree) GetImmutable(version int64) (*iavl.ImmutableTree, error) {
	if it.Version() != version {
		return nil, fmt.Errorf("version mismatch on immutable IAVL tree; got: %d, expected: %d", version, it.Version())
	}

	return it.ImmutableTree, nil
}

func (it *immutableTree) LoadVersionForOverwriting(targetVersion int64) error {
	panic("cannot call 'LoadVersionForOverwriting' on an immutable IAVL tree")
}

func (it *immutableTree) Iterator(start, end []byte, ascending bool) (types.Iterator, error) {
	iterator, err := it.ImmutableTree.Iterator(start, end, ascending)
	return iterator.(dbm.Iterator), err
}

func (mt *mutableTree) Iterator(start, end []byte, ascending bool) (types.Iterator, error) {
	iterator, err := mt.MutableTree.Iterator(start, end, ascending)
	return iterator.(dbm.Iterator), err
}
