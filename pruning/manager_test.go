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
	manager := pruning.NewManager(log.NewNopLogger())

	require.NotNil(t, manager)
	require.NotNil(t, manager.GetPruningHeights())
	require.Equal(t, types.Nothing, manager.GetOptions().GetType())
}

func Test_Strategies(t *testing.T) {
	testcases := map[string]struct {
		strategy *types.PruningOptions
		snapshotInterval uint64
		typeToAssert types.Type
		isValid  bool
	}{
		"prune nothing - no snapshot": {
			strategy: types.NewPruningOptions(types.Nothing),
			typeToAssert: types.Nothing,
		},
		"prune nothing - snapshot": {
			strategy: types.NewPruningOptions(types.Nothing),
			typeToAssert: types.Nothing,
			snapshotInterval: 100,
		},
		"prune default - no snapshot": {
			strategy: types.NewPruningOptions(types.Default),
			typeToAssert: types.Default,
		},
		"prune default - snapshot": {
			strategy: types.NewPruningOptions(types.Default),
			typeToAssert: types.Default,
			snapshotInterval: 100,
		},
		"prune everything - no snapshot": {
			strategy: types.NewPruningOptions(types.Everything),
			typeToAssert: types.Everything,
		},
		"prune everything - snapshot": {
			strategy: types.NewPruningOptions(types.Everything),
			typeToAssert: types.Everything,
			snapshotInterval: 100,
		},
		"custom 100-10-15": {
			strategy: types.NewCustomPruningOptions(100, 15),
			snapshotInterval: 10,
			typeToAssert: types.Custom,
		},
		"custom 10-10-15": {
			strategy: types.NewCustomPruningOptions(10, 15),
			snapshotInterval: 10,
			typeToAssert: types.Custom,
		},
		"custom 100-0-15": {
			strategy: types.NewCustomPruningOptions(100, 15),
			snapshotInterval: 0,
			typeToAssert: types.Custom,
		},
	}

	manager := pruning.NewManager(log.NewNopLogger())

	require.NotNil(t, manager)

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			curStrategy := tc.strategy 
			manager.SetSnapshotInterval(tc.snapshotInterval)
			
			pruneType := curStrategy.GetType()
			require.Equal(t, tc.typeToAssert, pruneType)

			// Validate strategy parameters
			switch pruneType {
			case types.Default:
				require.Equal(t, uint64(100000), curStrategy.KeepRecent)
				require.Equal(t, uint64(100), curStrategy.Interval)
			case types.Nothing:
				require.Equal(t, uint64(0), curStrategy.KeepRecent)
				require.Equal(t, uint64(0), curStrategy.Interval)
			case types.Everything:
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

				switch curStrategy.GetType() {
				case types.Nothing:
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

func Test_FlushLoad(t *testing.T) {
	manager := pruning.NewManager(log.NewNopLogger())
	require.NotNil(t, manager)

	db := db.NewMemDB()

	curStrategy := types.NewCustomPruningOptions(100, 15)

	snapshotInterval := uint64(10)
	manager.SetSnapshotInterval(snapshotInterval)

	manager.SetOptions(curStrategy)
	require.Equal(t, curStrategy, manager.GetOptions())

	keepRecent := curStrategy.KeepRecent

	heightsToPruneMirror := make([]int64, 0)

	for curHeight := int64(0); curHeight < 1000; curHeight++ {
		handleHeightActual := manager.HandleHeight(curHeight)

		curHeightStr := fmt.Sprintf("height: %d", curHeight)

		if curHeight > int64(keepRecent) && (snapshotInterval != 0 && (curHeight-int64(keepRecent))%int64(snapshotInterval) != 0 || snapshotInterval == 0) {
			expectedHandleHeight := curHeight - int64(keepRecent)
			require.Equal(t, expectedHandleHeight, handleHeightActual, curHeightStr)
			heightsToPruneMirror = append(heightsToPruneMirror, expectedHandleHeight)
		} else {
			require.Equal(t, int64(0), handleHeightActual, curHeightStr)
		}

		if manager.ShouldPruneAtHeight(curHeight) {
			manager.ResetPruningHeights()
			heightsToPruneMirror = make([]int64, 0)
		}

		// N.B.: There is no reason behind the choice of 3.
		if curHeight%3 == 0 {
			require.Equal(t, heightsToPruneMirror, manager.GetPruningHeights(), curHeightStr)
			batch := db.NewBatch()
			manager.FlushPruningHeights(batch)
			require.NoError(t, batch.Write())
			require.NoError(t, batch.Close())

			manager.ResetPruningHeights()
			require.Equal(t, make([]int64, 0), manager.GetPruningHeights(), curHeightStr)

			err := manager.LoadPruningHeights(db)
			require.NoError(t, err)
			require.Equal(t, heightsToPruneMirror, manager.GetPruningHeights(), curHeightStr)
		}
	}
}

func Test_WithSnapshot(t *testing.T) {
	manager := pruning.NewManager(log.NewNopLogger())
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

		snapshotHeightsToPruneMirrorNew := list.New() // TODO: optimize
		for e := snapshotHeightsToPruneMirror.Front(); e != nil; e = e.Next() {
			snapshotHeight := e.Value.(int64)
			if snapshotHeight < curHeight-int64(keepRecent) {
				heightsToPruneMirror = append(heightsToPruneMirror, snapshotHeight)
			} else {
				snapshotHeightsToPruneMirrorNew.PushBack(snapshotHeight)
			}
		}
		snapshotHeightsToPruneMirror = snapshotHeightsToPruneMirrorNew

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
