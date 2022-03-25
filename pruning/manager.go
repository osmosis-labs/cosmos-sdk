package pruning

import (
	"container/list"
	"encoding/binary"
	"fmt"
	"sync"

	"github.com/cosmos/cosmos-sdk/pruning/types"

	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

type Manager struct {
	logger               log.Logger
	db dbm.DB
	opts                 *types.PruningOptions
	snapshotInterval     uint64
	pruneHeights         []int64
	pruneSnapshotHeights *list.List
	mx                   sync.Mutex
}

const (
	pruneHeightsKey         = "s/pruneheights"
	pruneSnapshotHeightsKey = "s/pruneSnheights"

	uint64Size = 8
)

func NewManager(logger log.Logger, db dbm.DB) *Manager {
	return &Manager{
		logger:       logger,
		db: db,
		opts:         types.NewPruningOptions(types.PruningNothing),
		pruneHeights: []int64{},
		// These are the heights that are multiples of snapshotInterval and kept for state sync snapshots.
		// The heights are added to this list to be pruned when a snapshot is complete.
		pruneSnapshotHeights: list.New(),
		mx:                   sync.Mutex{},
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
func (m *Manager) HandleHeight(previousHeight int64) int64 {
	if m.opts.GetPruningStrategy() == types.PruningNothing {
		return 0
	}

	// Flag indicating whether should flush to disk.
	// It is set to true when an update to one of
	// pruneHeights or pruneSnapshotHeights is made
	shouldFlush := false

	defer func() {
		// handle persisted snapshot heights
		m.mx.Lock()
		defer m.mx.Unlock()
		defer func() {
			if shouldFlush {
				// Must be the unlocked implementation since we are under mutex
				m.flushAllPruningHeightsUnlocked()
			}
		}()

		var next *list.Element
		for e := m.pruneSnapshotHeights.Front(); e != nil; e = next {
			snHeight := e.Value.(int64)
			if snHeight < previousHeight-int64(m.opts.KeepRecent) {
				m.pruneHeights = append(m.pruneHeights, snHeight)

				// We must get next before removing to be able to continue iterating.
				next = e.Next()
				m.pruneSnapshotHeights.Remove(e)

				shouldFlush = true
			} else {
				next = e.Next()
			}
		}
	}()

	if int64(m.opts.KeepRecent) < previousHeight {
		pruneHeight := previousHeight - int64(m.opts.KeepRecent)
		// We consider this height to be pruned iff:
		//
		// - snapshotInterval is zero as that means that all heights should be pruned.
		// - snapshotInterval % (height - KeepRecent) != 0 as that means the height is not
		// a 'snapshot' height.
		if m.snapshotInterval == 0 || pruneHeight%int64(m.snapshotInterval) != 0 {
			m.pruneHeights = append(m.pruneHeights, pruneHeight)
			shouldFlush = true
			return pruneHeight
		}
	}
	return 0
}

func (m *Manager) HandleHeightSnapshot(height int64) {
	if m.opts.GetPruningStrategy() == types.PruningNothing {
		return
	}
	m.mx.Lock()
	defer m.mx.Unlock()
	m.logger.Debug("HandleHeightSnapshot", "height", height) // TODO: change log level to Debug
	m.pruneSnapshotHeights.PushBack(height)
	// m.flushPruningSnapshotHeightsUnlocked()
}

// SetSnapshotInterval sets the interval at which the snapshots are taken.
func (m *Manager) SetSnapshotInterval(snapshotInterval uint64) {
	m.snapshotInterval = snapshotInterval
}

// ShouldPruneAtHeight return true if the given height should be pruned, false otherwise
func (m *Manager) ShouldPruneAtHeight(height int64) bool {
	return m.opts.GetPruningStrategy() != types.PruningNothing && m.opts.Interval > 0 && height%int64(m.opts.Interval) == 0
}

// LoadPruningHeights loads the pruning heights from the database as a crash recovery.
func (m *Manager) LoadPruningHeights(db dbm.DB) error {
	if m.opts.GetPruningStrategy() == types.PruningNothing {
		return nil
	}
	if err := m.loadPruningHeights(db); err != nil {
		return err
	}
	if err := m.loadPruningSnapshotHeights(db); err != nil {
		return err
	}
	return nil
}

func (m *Manager) loadPruningHeights(db dbm.DB) error {
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

func (m *Manager) loadPruningSnapshotHeights(db dbm.DB) error {
	bz, err := db.Get([]byte(pruneSnapshotHeightsKey))
	if err != nil {
		return fmt.Errorf("failed to get post-snapshot pruned heights: %w", err)
	}
	if len(bz) == 0 {
		return nil
	}

	pruneSnapshotHeights := list.New()
	i, offset := 0, 0
	for offset < len(bz) {
		pruneSnapshotHeights.PushBack(int64(binary.BigEndian.Uint64(bz[offset : offset+8])))
		i++
		offset += 8
	}

	if pruneSnapshotHeights.Len() > 0 {
		m.mx.Lock()
		defer m.mx.Unlock()
		m.pruneSnapshotHeights = pruneSnapshotHeights
	}

	return nil
}

// FlushAllPruningHeights flushes the pruning heights to the database for crash recovery.
// "All" refers to regular pruning heights and snapshot heights, if any.
func (m *Manager) FlushAllPruningHeights() {
	if m.opts.GetPruningStrategy() == types.PruningNothing {
		return
	}
	batch := m.db.NewBatch()
	defer batch.Close()
	m.flushPruningHeightsBatch(batch)
	m.flushPruningSnapshotHeightsBatch(batch)

	if err := batch.Write(); err != nil {
		panic(fmt.Errorf("error on batch write %w", err))
	}
}

// flushAllPruningHeightsUnlocked flushes the pruning heights to the database for crash recovery.
// "All" refers to regular pruning heights and snapshot heights, if any.
// It serves the same function as exported FlushPruningHeights. However, it assummes that
// mutex was acquired prior to calling this method.
func (m *Manager) flushAllPruningHeightsUnlocked() {
	if m.opts.GetPruningStrategy() == types.PruningNothing {
		return
	}
	batch := m.db.NewBatch()
	defer batch.Close()
	m.flushPruningHeightsBatch(batch)
	m.flushPruningSnapshotHeightsUnlockedBatch(batch)

	if err := batch.Write(); err != nil {
		panic(fmt.Errorf("error on batch write %w", err))
	}
}

func (m *Manager) flushPruningHeightsBatch(batch dbm.Batch) {
	bz := make([]byte, 0, uint64Size * len(m.pruneHeights))
	for _, ph := range m.pruneHeights {
		buf := make([]byte, uint64Size)
		binary.BigEndian.PutUint64(buf, uint64(ph))
		bz = append(bz, buf...)
	}

	if err := batch.Set([]byte(pruneHeightsKey), bz); err != nil {
		panic(err)
	}
}

func (m *Manager) flushPruningSnapshotHeightsBatch(batch dbm.Batch) {
	m.mx.Lock()
	defer m.mx.Unlock()
	m.flushPruningSnapshotHeightsUnlockedBatch(batch)
}

func (m *Manager) flushPruningSnapshotHeightsUnlockedBatch(batch dbm.Batch) {
	bz := convertInt64ListToBytes(m.pruneSnapshotHeights)
	if err := batch.Set([]byte(pruneSnapshotHeightsKey), bz); err != nil {
		panic(err)
	}
}

func (m *Manager) flushPruningSnapshotHeightsUnlocked() {
	bz := convertInt64ListToBytes(m.pruneSnapshotHeights)
	if err := m.db.Set([]byte(pruneSnapshotHeightsKey), bz); err != nil {
		panic(err)
	}
}

// INVARIANT: list contains ints
func convertInt64ListToBytes(list *list.List) []byte {
	bz := make([]byte, 0, uint64Size * list.Len())
	for e := list.Front(); e != nil; e = e.Next() {
		buf := make([]byte, uint64Size)
		binary.BigEndian.PutUint64(buf, uint64(e.Value.(int64)))
		bz = append(bz, buf...)
	}
	return bz
}
