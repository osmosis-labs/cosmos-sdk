package mempool

import (
	"github.com/cosmos/cosmos-sdk/types"
)

// Tx we define an app-side mempool transaction interface that is as
// minimal as possible, only requiring applications to define the size of the
// transaction to be used when reaping and getting the transaction itself.
// Interface type casting can be used in the actual app-side mempool implementation.
type Tx interface {
	types.Tx

	// Size returns the size of the transaction in bytes.
	Size() int

	// Hash returns the hash of the transaction.
	Hash() [32]byte
}

type Mempool interface {
	// Insert attempts to insert a Tx into the app-side mempool returning
	// an error upon failure.
	Insert(types.Context, Tx) error

	// Select returns the next set of available transactions from the app-side
	// mempool, up to maxBytes or until the mempool is empty. The application can
	// decide to return transactions from its own mempool, from the incoming
	// txs, or some combination of both.
	Select(ctx types.Context, txs [][]byte, maxBytes int) ([]Tx, error)

	// CountTx returns the number of transactions currently in the mempool.
	CountTx() int

	// Remove attempts to remove a transaction from the mempool, returning an error
	// upon failure.
	Remove(types.Context, Tx) error
}
