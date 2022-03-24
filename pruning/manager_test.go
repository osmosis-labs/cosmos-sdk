package pruning_test

import (
	"container/list"
	"fmt"

	"sync"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/pruning"
	"github.com/cosmos/cosmos-sdk/pruning/types"
	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"
)

func Test_NewManager(t *testing.T) {
	manager := pruning.NewManager(log.NewNopLogger(), db.NewMemDB())

	require.NotNil(t, manager)
	require.NotNil(t, manager.GetPruningHeights())
	require.Equal(t, types.PruningNothing, manager.GetOptions().GetPruningStrategy())
}

func Test_Strategies(t *testing.T) {
	testcases := map[string]struct {
		strategy         *types.PruningOptions
		snapshotInterval uint64
		strategyToAssert types.PruningStrategy
		isValid          bool
	}{
		"prune nothing - no snapshot": {
			strategy:         types.NewPruningOptions(types.PruningNothing),
			strategyToAssert: types.PruningNothing,
		},
		"prune nothing - snapshot": {
			strategy:         types.NewPruningOptions(types.PruningNothing),
			strategyToAssert: types.PruningNothing,
			snapshotInterval: 100,
		},
		"prune default - no snapshot": {
			strategy:         types.NewPruningOptions(types.PruningDefault),
			strategyToAssert: types.PruningDefault,
		},
		"prune default - snapshot": {
			strategy:         types.NewPruningOptions(types.PruningDefault),
			strategyToAssert: types.PruningDefault,
			snapshotInterval: 100,
		},
		"prune everything - no snapshot": {
			strategy:         types.NewPruningOptions(types.PruningEverything),
			strategyToAssert: types.PruningEverything,
		},
		"prune everything - snapshot": {
			strategy:         types.NewPruningOptions(types.PruningEverything),
			strategyToAssert: types.PruningEverything,
			snapshotInterval: 100,
		},
		"custom 100-10-15": {
			strategy:         types.NewCustomPruningOptions(100, 15),
			snapshotInterval: 10,
			strategyToAssert: types.PruningCustom,
		},
		"custom 10-10-15": {
			strategy:         types.NewCustomPruningOptions(10, 15),
			snapshotInterval: 10,
			strategyToAssert: types.PruningCustom,
		},
		"custom 100-0-15": {
			strategy:         types.NewCustomPruningOptions(100, 15),
			snapshotInterval: 0,
			strategyToAssert: types.PruningCustom,
		},
	}

	manager := pruning.NewManager(log.NewNopLogger(), db.NewMemDB())

	require.NotNil(t, manager)

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			curStrategy := tc.strategy
			manager.SetSnapshotInterval(tc.snapshotInterval)

			pruneStrategy := curStrategy.GetPruningStrategy()
			require.Equal(t, tc.strategyToAssert, pruneStrategy)

			// Validate strategy parameters
			switch pruneStrategy {
			case types.PruningDefault:
				require.Equal(t, uint64(100000), curStrategy.KeepRecent)
				require.Equal(t, uint64(100), curStrategy.Interval)
			case types.PruningNothing:
				require.Equal(t, uint64(0), curStrategy.KeepRecent)
				require.Equal(t, uint64(0), curStrategy.Interval)
			case types.PruningEverything:
				require.Equal(t, uint64(10), curStrategy.KeepRecent)
				require.Equal(t, uint64(10), curStrategy.Interval)
			default:
				//
			}

			manager.SetOptions(curStrategy)
			require.Equal(t, tc.strategy, manager.GetOptions())

			curKeepRecent := curStrategy.KeepRecent
			curInterval := curStrategy.Interval

			for curHeight := int64(0); curHeight < 110000; curHeight++ {
				handleHeightActual := manager.HandleHeight(curHeight)
				shouldPruneAtHeightActual := manager.ShouldPruneAtHeight(curHeight)

				curPruningHeihts := manager.GetPruningHeights()

				curHeightStr := fmt.Sprintf("height: %d", curHeight)

				switch curStrategy.GetPruningStrategy() {
				case types.PruningNothing:
					require.Equal(t, int64(0), handleHeightActual, curHeightStr)
					require.False(t, shouldPruneAtHeightActual, curHeightStr)

					require.Equal(t, 0, len(manager.GetPruningHeights()))
				default:
					if curHeight > int64(curKeepRecent) && (tc.snapshotInterval != 0 && (curHeight-int64(curKeepRecent))%int64(tc.snapshotInterval) != 0 || tc.snapshotInterval == 0) {
						expectedHeight := curHeight - int64(curKeepRecent)
						require.Equal(t, curHeight-int64(curKeepRecent), handleHeightActual, curHeightStr)

						require.Contains(t, curPruningHeihts, expectedHeight, curHeightStr)
					} else {
						require.Equal(t, int64(0), handleHeightActual, curHeightStr)

						require.Equal(t, 0, len(manager.GetPruningHeights()))
					}
					require.Equal(t, curHeight%int64(curInterval) == 0, shouldPruneAtHeightActual, curHeightStr)
				}
				manager.ResetPruningHeights()
				require.Equal(t, 0, len(manager.GetPruningHeights()))
			}
		})
	}
}

// Here, we focus on testing the correctness of flushing pruneHeights to disk and loading back.
func Test_FlushLoad_PruneHeights(t *testing.T) {
	db := db.NewMemDB()

	manager := pruning.NewManager(log.NewNopLogger(), db)
	require.NotNil(t, manager)

	const pruningKeepRecent = 100

	curStrategy := types.NewCustomPruningOptions(uint64(pruningKeepRecent), 15)

	snapshotInterval := uint64(10)
	manager.SetSnapshotInterval(snapshotInterval)

	manager.SetOptions(curStrategy)
	require.Equal(t, curStrategy, manager.GetOptions())

	heightsToPruneMirror := make([]int64, 0)

	for curHeight := int64(1); curHeight < 1000; curHeight++ {
		handleHeightActual := manager.HandleHeight(curHeight)

		curHeightStr := fmt.Sprintf("height: %d", curHeight)

		// flush should happen every time a height is persisted
		if handleHeightActual > 0 {
			heightsToPruneMirror = append(heightsToPruneMirror, curHeight-pruningKeepRecent)

			actualPruningHeightsBeforeLoad := manager.GetPruningHeights()
			require.Equal(t, heightsToPruneMirror, actualPruningHeightsBeforeLoad)

			manager.ResetPruningHeights()

			err := manager.LoadPruningHeights(db)
			require.NoError(t, err)

			actualPruningHeightsAfterLoad := manager.GetPruningHeights()
			require.Equal(t, actualPruningHeightsBeforeLoad, actualPruningHeightsAfterLoad, curHeightStr)
		}
	}
}

// In this test, we focus on the correctness of flushing pruneSnapshotHeight to disk and loading back.
// To do so, we assert that, eventually, all heights are moved from pruningSnapshotHeights
// to pruningHeighs. Additionally, we claim all heights are preserved even in the conditions
// when they are periodically flushed to disk, reset and then loaded back from disk into memory.
func Test_FlushLoad_PruneSnapshotHeights(t *testing.T) {
	const (
		pruningKeepRecent = 30
		totalHeights      = int64(1000)
		snapshotInterval  = uint64(10)
	)

	var (
		db          = db.NewMemDB()
		manager     = pruning.NewManager(log.NewNopLogger(), db)
		curStrategy = types.NewCustomPruningOptions(uint64(pruningKeepRecent), 15)
		wg          = &sync.WaitGroup{}
	)
	require.NotNil(t, manager)

	manager.SetSnapshotInterval(snapshotInterval)
	manager.SetOptions(curStrategy)
	require.Equal(t, curStrategy, manager.GetOptions())

	for curHeight := int64(1); curHeight < totalHeights; curHeight++ {
		// We always use previous height for handling pruning
		// but current height for snapshots
		previousHeight := curHeight - 1

		handleHeightActual := manager.HandleHeight(previousHeight)

		// No flush / update should happen at snapshot heights
		if previousHeight%int64(snapshotInterval) == 0 {
			require.True(t, handleHeightActual == 0)
		}

		// flush should happen every time a height is persisted
		if handleHeightActual > 0 {

			manager.ResetPruningHeights()

			err := manager.LoadPruningHeights(db)
			require.NoError(t, err)
		}

		if curHeight%int64(snapshotInterval) == 0 {
			wg.Add(1)
			// Mimic taking snapshot in a separate goroutine
			go func(h int64) {

				// Mimic delay to take snapshots
				time.Sleep(30 * time.Millisecond)
				manager.HandleHeightSnapshot(h)
				wg.Done()
			}(curHeight)
		}
	}

	wg.Wait()

	// We call HandleHeight(lastHeight) after wg.Wait() because the pruningSnapshotHeights are moved to pruningHeights
	// only when HandleHeight(...) is called.
	// Doing the call to HandleHeight(...) after all snapshots are complete (this is ensured by the wait group above),
	// allows us to assert with confidence on what pruneHeights should contain.
	lastHeight := totalHeights + 1
	handleHeightActual := manager.HandleHeight(lastHeight)
	require.Equal(t, int64(lastHeight-pruningKeepRecent), handleHeightActual)

	pruningHeights := manager.GetPruningHeights()
	require.Equal(t, int(totalHeights+1-1-pruningKeepRecent), len(pruningHeights))
}

func Test_WithSnapshot(t *testing.T) {
	manager := pruning.NewManager(log.NewNopLogger(), db.NewMemDB())
	require.NotNil(t, manager)

	curStrategy := types.NewCustomPruningOptions(10, 10)

	snapshotInterval := uint64(15)
	manager.SetSnapshotInterval(snapshotInterval)

	manager.SetOptions(curStrategy)
	require.Equal(t, curStrategy, manager.GetOptions())

	keepRecent := curStrategy.KeepRecent

	heightsToPruneMirror := make([]int64, 0)

	mx := sync.Mutex{}
	snapshotHeightsToPruneMirror := list.New()

	wg := sync.WaitGroup{}

	for curHeight := int64(1); curHeight < 100000; curHeight++ {
		mx.Lock()
		handleHeightActual := manager.HandleHeight(curHeight)

		curHeightStr := fmt.Sprintf("height: %d", curHeight)

		if curHeight > int64(keepRecent) && (curHeight-int64(keepRecent))%int64(snapshotInterval) != 0 {
			expectedHandleHeight := curHeight - int64(keepRecent)
			require.Equal(t, expectedHandleHeight, handleHeightActual, curHeightStr)
			heightsToPruneMirror = append(heightsToPruneMirror, expectedHandleHeight)
		} else {
			require.Equal(t, int64(0), handleHeightActual, curHeightStr)
		}

		actualHeightsToPrune := manager.GetPruningHeights()

		var next *list.Element
		for e := snapshotHeightsToPruneMirror.Front(); e != nil; e = next {
			snapshotHeight := e.Value.(int64)
			if snapshotHeight < curHeight-int64(keepRecent) {
				heightsToPruneMirror = append(heightsToPruneMirror, snapshotHeight)

				// We must get next before removing to be able to continue iterating.
				next = e.Next()
				snapshotHeightsToPruneMirror.Remove(e)
			} else {
				next = e.Next()
			}
		}

		require.Equal(t, heightsToPruneMirror, actualHeightsToPrune, curHeightStr)
		mx.Unlock()

		if manager.ShouldPruneAtHeight(curHeight) {
			manager.ResetPruningHeights()
			heightsToPruneMirror = make([]int64, 0)
		}

		// Mimic taking snapshots in the background
		if curHeight%int64(snapshotInterval) == 0 {
			wg.Add(1)
			go func(curHeightCp int64) {
				time.Sleep(time.Nanosecond * 500)

				mx.Lock()
				manager.HandleHeightSnapshot(curHeightCp)
				snapshotHeightsToPruneMirror.PushBack(curHeightCp)
				mx.Unlock()
				wg.Done()
			}(curHeight)
		}
	}

	wg.Wait()
}
