package pruning

import (
	"encoding/binary"
	"fmt"

	"github.com/cosmos/cosmos-sdk/pruning/types"

	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

type Manager struct {
	logger         log.Logger
	opts *types.PruningOptions
	pruneHeights   []int64
}

const(
	pruneHeightsKey  = "s/pruneheights"
)

func NewManager(logger         log.Logger) *Manager {
	return &Manager{
		logger: logger,
		opts: types.PruneNothing,
		pruneHeights: []int64{},
	}
}

// SetOptions sets the pruning strategy on the manager.
func (m *Manager) SetOptions(opts *types.PruningOptions) {
	m.opts = opts
}

// GetOptions fetches the pruning strategy from the manager.
func (m *Manager) GetOptions() *types.PruningOptions {
	return m.opts
}

// GetPruningHeights returns all heights to be pruned during the next call to Prune().
func (m *Manager) GetPruningHeights() []int64 {
	return m.pruneHeights
}

// ResetPruningHeights resets the heights to be pruned.
func (m *Manager) ResetPruningHeights() {
	m.pruneHeights = make([]int64, 0)
}

// HandleHeight determines if pruneHeight height needs to be kept for pruning at the right interval prescribed by
// the pruning strategy. Returns true if the given height was kept to be pruned at the next call to Prune(), false otherwise
func (m *Manager) HandleHeight(height int64) int64{
	if m.opts == types.PruneNothing {
		return 0
	}
	if m.opts == types.PruneEverything {
		return height
	}

	if int64(m.opts.KeepRecent) < height {
		pruneHeight := height - int64(m.opts.KeepRecent)
		// We consider this height to be pruned iff:
		//
		// - KeepEvery is zero as that means that all heights should be pruned.
		// - KeepEvery % (height - KeepRecent) != 0 as that means the height is not
		// a 'snapshot' height.
		if m.opts.KeepEvery == 0 || pruneHeight%int64(m.opts.KeepEvery) != 0 {
			m.pruneHeights = append(m.pruneHeights, pruneHeight)
			return pruneHeight
		}
	}
	return 0
}

// ShouldPruneAtHeight return true if the given height should be pruned, false otherwise
func (m *Manager) ShouldPruneAtHeight(height int64) bool {
	return m.opts != types.PruneNothing && m.opts.Interval > 0 && height%int64(m.opts.Interval) == 0
}

// FlushPruningHeights flushes the pruning heights to the database for crash recovery.
func (m *Manager) FlushPruningHeights(batch dbm.Batch) {
	bz := make([]byte, 0)
	for _, ph := range m.pruneHeights {
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, uint64(ph))
		bz = append(bz, buf...)
	}

	batch.Set([]byte(pruneHeightsKey), bz)
}

// LoadPruningHeights loads the pruning heights from the database as a crash recovery.
func (m *Manager) LoadPruningHeights(db dbm.DB) error {
	bz, err := db.Get([]byte(pruneHeightsKey))
	if err != nil {
		return fmt.Errorf("failed to get pruned heights: %w", err)
	}
	if len(bz) == 0 {
		return nil
	}

	prunedHeights := make([]int64, len(bz)/8)
	i, offset := 0, 0
	for offset < len(bz) {
		prunedHeights[i] = int64(binary.BigEndian.Uint64(bz[offset : offset+8]))
		i++
		offset += 8
	}

	if len(prunedHeights) > 0 {
		m.pruneHeights = prunedHeights
	}

	return nil
}
