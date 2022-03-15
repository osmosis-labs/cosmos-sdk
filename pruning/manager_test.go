package pruning_test

import (
	"fmt"
	"testing"

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
	require.Equal(t, types.PruneNothing, manager.GetOptions())
}

func Test_Strategies(t *testing.T) {
	testcases := map[string]struct {
		strategy *types.PruningOptions
		isValid  bool
	}{
		"prune nothing": {
			strategy: types.PruneNothing,
			isValid:  true,
		},
		"prune default": {
			strategy: types.PruneDefault,
			isValid:  true,
		},
		"prune everything": {
			strategy: types.PruneEverything,
			isValid:  true,
		},
		"custom 100-10-15": {
			strategy: types.NewPruningOptions(100, 10, 15),
			isValid:  true,
		},
		"custom 0-10-15": {
			strategy: types.NewPruningOptions(0, 10, 15),
			isValid:  true,
		},
		"custom 100-0-15": {
			strategy: types.NewPruningOptions(0, 10, 15),
			isValid:  true,
		},
		"custom 100-0-0": { // does not make sense
			strategy: types.NewPruningOptions(100, 0, 0),
			isValid:  false,
		},
		"custom 0-10-0": { // does not make sense
			strategy: types.NewPruningOptions(0, 10, 0),
			isValid:  false,
		},
		"custom 0-0-1": { // does not make sense
			strategy: types.NewPruningOptions(0, 0, 1),
			isValid:  true,
		},
	}

	manager := pruning.NewManager(log.NewNopLogger())

	require.NotNil(t, manager)

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			curStrategy := tc.strategy

			require.Equal(t, tc.isValid, curStrategy.Validate() == nil)
			if !tc.isValid {
				return
			}

			manager.SetOptions(curStrategy)
			require.Equal(t, tc.strategy, manager.GetOptions())

			curKeepRecent := curStrategy.KeepRecent
			curKeepEvery := curStrategy.KeepEvery
			curInterval := curStrategy.Interval

			for h := int64(0); h < 110000; h++ {
				handleHeightActual := manager.HandleHeight(h)
				shouldPruneAtHeightActual := manager.ShouldPruneAtHeight(h)

				curHeightStr := fmt.Sprintf("height: %d", h)

				switch curStrategy {
				case types.PruneNothing:
					require.Equal(t, int64(0), handleHeightActual, curHeightStr)
					require.False(t, shouldPruneAtHeightActual, curHeightStr)
				case types.PruneEverything:
					require.Equal(t, h, handleHeightActual, fmt.Sprintf("height: %d", h))
					require.Equal(t, h%int64(types.PruneEverything.Interval) == 0, shouldPruneAtHeightActual, curHeightStr)
				case types.PruneDefault:
					if h > int64(types.PruneDefault.KeepRecent) {
						require.Equal(t, h-int64(types.PruneDefault.KeepRecent), handleHeightActual, curHeightStr)
					} else {
						require.Equal(t, int64(0), handleHeightActual, curHeightStr)
					}
					require.Equal(t, h%int64(types.PruneDefault.Interval) == 0, shouldPruneAtHeightActual, curHeightStr)
				default:
					if h > int64(curKeepRecent) && (curKeepEvery != 0 && h%int64(curKeepEvery) != 0 || curKeepEvery == 0) {
						require.Equal(t, h-int64(curKeepRecent), handleHeightActual, curHeightStr)
					} else {
						require.Equal(t, int64(0), handleHeightActual, curHeightStr)
					}
					require.Equal(t, h%int64(curInterval) == 0, shouldPruneAtHeightActual, curHeightStr)
				}
			}
		})
	}
}

func Test_FlushLoad(t *testing.T) {
	manager := pruning.NewManager(log.NewNopLogger())
	require.NotNil(t, manager)

	db := db.NewMemDB()

	curStrategy := types.NewPruningOptions(100, 10, 15)

	manager.SetOptions(curStrategy)
	require.Equal(t, curStrategy, manager.GetOptions())

	curKeepRecent := curStrategy.KeepRecent
	curKeepEvery := curStrategy.KeepEvery

	heightsToPruneMirror := make([]int64, 0)

	for h := int64(0); h < 1000; h++ {
		handleHeightActual := manager.HandleHeight(h)

		curHeightStr := fmt.Sprintf("height: %d", h)

		if h > int64(curKeepRecent) && (curKeepEvery != 0 && h%int64(curKeepEvery) != 0 || curKeepEvery == 0) {
			expectedHandleHeight := h - int64(curKeepRecent)
			require.Equal(t, expectedHandleHeight, handleHeightActual, curHeightStr)
			heightsToPruneMirror = append(heightsToPruneMirror, expectedHandleHeight)
		} else {
			require.Equal(t, int64(0), handleHeightActual, curHeightStr)
		}

		if manager.ShouldPruneAtHeight(h) {
			manager.ResetPruningHeights()
			heightsToPruneMirror = make([]int64, 0)
		}

		if h%3 == 0 {
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
