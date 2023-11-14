package pruning_test

import (
	"container/list"
	"errors"
	"fmt"
	"testing"

	"cosmossdk.io/log"
	db "github.com/cosmos/cosmos-db"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/store/pruning"
	"github.com/cosmos/cosmos-sdk/store/pruning/mock"
	"github.com/cosmos/cosmos-sdk/store/pruning/types"
)

const dbErr = "db error"

func TestNewManager(t *testing.T) {
	manager := pruning.NewManager(db.NewMemDB(), log.NewNopLogger())

	require.NotNil(t, manager)
	require.Equal(t, types.PruningNothing, manager.GetOptions().GetPruningStrategy())
}

func TestStrategies(t *testing.T) {
	testcases := map[string]struct {
		strategy         types.PruningOptions
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

	manager := pruning.NewManager(db.NewMemDB(), log.NewNopLogger())

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
				require.Equal(t, uint64(362880), curStrategy.KeepRecent)
				require.Equal(t, uint64(10), curStrategy.Interval)
			case types.PruningNothing:
				require.Equal(t, uint64(0), curStrategy.KeepRecent)
				require.Equal(t, uint64(0), curStrategy.Interval)
			case types.PruningEverything:
				require.Equal(t, uint64(2), curStrategy.KeepRecent)
				require.Equal(t, uint64(10), curStrategy.Interval)
			default:
				//
			}

			manager.SetOptions(curStrategy)
			snHeight := int64(tc.snapshotInterval - 1)
			require.Equal(t, tc.strategy, manager.GetOptions())

			curKeepRecent := curStrategy.KeepRecent
			curInterval := curStrategy.Interval

			for curHeight := int64(0); curHeight < 110000; curHeight++ {
				shouldPruneAtHeightActual := manager.ShouldPruneAtHeight(curHeight)

				pruningHeightActual := manager.GetPruningHeight(curHeight)

				curHeightStr := fmt.Sprintf("height: %d", curHeight)

				switch curStrategy.GetPruningStrategy() {
				case types.PruningNothing:
					require.Equal(t, int64(0), pruningHeightActual, curHeightStr)
				default:
					if curHeight > int64(curKeepRecent) && curHeight%int64(curStrategy.Interval) == 0 {
						pruningHeightExpected := curHeight - int64(curKeepRecent) - 1
						if tc.snapshotInterval > 0 && snHeight < pruningHeightExpected {
							pruningHeightExpected = snHeight
						}
						require.Equal(t, pruningHeightExpected, pruningHeightActual, curHeightStr)
					} else {
						require.Equal(t, int64(0), pruningHeightActual, curHeightStr)
					}
					require.Equal(t, curHeight%int64(curInterval) == 0, shouldPruneAtHeightActual, curHeightStr)
				}
			}
		})
	}
}

func TestPruningHeight_Inputs(t *testing.T) {
	keepRecent := int64(types.NewPruningOptions(types.PruningEverything).KeepRecent)
	interval := int64(types.NewPruningOptions(types.PruningEverything).Interval)

	testcases := map[string]struct {
		height         int64
		expectedResult int64
		strategy       types.PruningStrategy
	}{
		"currentHeight is negative - prune everything - invalid currentHeight": {
			-1,
			0,
			types.PruningEverything,
		},
		"currentHeight is  zero - prune everything - invalid currentHeight": {
			0,
			0,
			types.PruningEverything,
		},
		"currentHeight is positive but within keep recent- prune everything - not kept": {
			keepRecent,
			0,
			types.PruningEverything,
		},
		"currentHeight is positive and equal to keep recent+1 - no kept": {
			keepRecent + 1,
			0,
			types.PruningEverything,
		},
		"currentHeight is positive and greater than keep recent+1 but not multiple of interval - no kept": {
			keepRecent + 2,
			0,
			types.PruningEverything,
		},
		"currentHeight is positive and greater than keep recent+1 and multiple of interval - kept": {
			interval,
			interval - keepRecent - 1,
			types.PruningEverything,
		},
		"pruning nothing, currentHeight is positive and greater than keep recent - not kept": {
			interval,
			0,
			types.PruningNothing,
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			manager := pruning.NewManager(db.NewMemDB(), log.NewNopLogger())
			require.NotNil(t, manager)
			manager.SetOptions(types.NewPruningOptions(tc.strategy))

			pruningHeightActual := manager.GetPruningHeight(tc.height)
			require.Equal(t, tc.expectedResult, pruningHeightActual)
		})
	}
}

func TestHandleHeight_DbErr_Panic(t *testing.T) {
	ctrl := gomock.NewController(t)

	// Setup
	dbMock := mock.NewMockDB(ctrl)

	dbMock.EXPECT().SetSync(gomock.Any(), gomock.Any()).Return(errors.New(dbErr)).Times(1)

	manager := pruning.NewManager(dbMock, log.NewNopLogger())
	manager.SetOptions(types.NewPruningOptions(types.PruningEverything))
	require.NotNil(t, manager)

	defer func() {
		if r := recover(); r == nil {
			t.Fail()
		}
	}()

	manager.HandleHeightSnapshot(10)
}

func TestHandleHeightSnapshot_FlushLoadFromDisk(t *testing.T) {
	loadedHeightsMirror := []int64{}

	// Setup
	db := db.NewMemDB()
	manager := pruning.NewManager(db, log.NewNopLogger())
	require.NotNil(t, manager)

	manager.SetOptions(types.NewPruningOptions(types.PruningEverything))

	for snapshotHeight := int64(-1); snapshotHeight < 100; snapshotHeight++ {
		// Test flush
		manager.HandleHeightSnapshot(snapshotHeight)

		// Post test
		if snapshotHeight > 0 {
			loadedHeightsMirror = append(loadedHeightsMirror, snapshotHeight)
		}

		loadedSnapshotHeights, err := pruning.LoadPruningSnapshotHeights(db)
		require.NoError(t, err)
		require.Equal(t, len(loadedHeightsMirror), len(loadedSnapshotHeights))

		// Test load back
		err = manager.LoadPruningSnapshotHeights(db)
		require.NoError(t, err)

		loadedSnapshotHeights, err = pruning.LoadPruningSnapshotHeights(db)
		require.NoError(t, err)
		require.Equal(t, len(loadedHeightsMirror), len(loadedSnapshotHeights))
	}
}

func TestHandleHeightSnapshot_DbErr_Panic(t *testing.T) {
	ctrl := gomock.NewController(t)

	// Setup
	dbMock := mock.NewMockDB(ctrl)

	dbMock.EXPECT().SetSync(gomock.Any(), gomock.Any()).Return(errors.New(dbErr)).Times(1)

	manager := pruning.NewManager(dbMock, log.NewNopLogger())
	manager.SetOptions(types.NewPruningOptions(types.PruningEverything))
	require.NotNil(t, manager)

	defer func() {
		if r := recover(); r == nil {
			t.Fail()
		}
	}()

	manager.HandleHeightSnapshot(10)
}

func TestLoadSnapshotPruningHeights(t *testing.T) {
	var (
		manager = pruning.NewManager(db.NewMemDB(), log.NewNopLogger())
		err     error
	)
	require.NotNil(t, manager)

	// must not be PruningNothing
	manager.SetOptions(types.NewPruningOptions(types.PruningDefault))

	testcases := map[string]struct {
		flushedPruningHeights            []int64
		getFlushedPruningSnapshotHeights func() *list.List
		expectedResult                   error
	}{
		"negative pruningHeight - error": {
			flushedPruningHeights: []int64{10, 0, -1},
			expectedResult:        &pruning.NegativeHeightsError{Height: -1},
		},
		"negative snapshotPruningHeight - error": {
			getFlushedPruningSnapshotHeights: func() *list.List {
				l := list.New()
				l.PushBack(int64(5))
				l.PushBack(int64(-2))
				l.PushBack(int64(3))
				return l
			},
			expectedResult: &pruning.NegativeHeightsError{Height: -2},
		},
		"both have negative - pruningHeight error": {
			flushedPruningHeights: []int64{10, 0, -1},
			getFlushedPruningSnapshotHeights: func() *list.List {
				l := list.New()
				l.PushBack(int64(5))
				l.PushBack(int64(-2))
				l.PushBack(int64(3))
				return l
			},
			expectedResult: &pruning.NegativeHeightsError{Height: -1},
		},
		"both non-negative - success": {
			flushedPruningHeights: []int64{10, 0, 3},
			getFlushedPruningSnapshotHeights: func() *list.List {
				l := list.New()
				l.PushBack(int64(5))
				l.PushBack(int64(0))
				l.PushBack(int64(3))
				return l
			},
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			db := db.NewMemDB()
			if tc.flushedPruningHeights != nil {
				err = db.Set(pruning.PruneHeightsKey, pruning.Int64SliceToBytes(tc.flushedPruningHeights))
				require.NoError(t, err)
			}

			if tc.getFlushedPruningSnapshotHeights != nil {
				err = db.Set(pruning.PruneSnapshotHeightsKey, pruning.ListToBytes(tc.getFlushedPruningSnapshotHeights()))
				require.NoError(t, err)
			}

			err = manager.LoadPruningSnapshotHeights(db)
			require.Equal(t, tc.expectedResult, err)
		})
	}
}

func TestLoadSnapshotPruningHeights_PruneNothing(t *testing.T) {
	manager := pruning.NewManager(db.NewMemDB(), log.NewNopLogger())
	require.NotNil(t, manager)

	manager.SetOptions(types.NewPruningOptions(types.PruningNothing))

	require.Nil(t, manager.LoadPruningSnapshotHeights(db.NewMemDB()))
}
