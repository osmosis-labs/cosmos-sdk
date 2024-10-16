package keeper

import (
	"context"

	abci "github.com/cometbft/cometbft/abci/types"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

// BeginBlocker will persist the current header and validator set as a historical entry
// and prune the oldest entry based on the HistoricalEntries parameter
func (k *Keeper) BeginBlocker(ctx context.Context) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, telemetry.Now(), telemetry.MetricKeyBeginBlocker)
	return k.TrackHistoricalInfo(ctx)
}

// EndBlocker called at every block, update validator set
func (k *Keeper) EndBlocker(ctx context.Context) ([]abci.ValidatorUpdate, error) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, telemetry.Now(), telemetry.MetricKeyEndBlocker)

	// TODO: Remove migration code and panic catch in the next upgrade
	// Wrap the migration call in a function that can recover from panics
	func() {
		defer func() {
			if r := recover(); r != nil {
				k.Logger(sdk.UnwrapSDKContext(ctx)).Error("Panic in MigrateDelegationsByValidatorIndex", "recover", r)
			}
		}()

		// Only migrate 10000 items per block to make the migration as fast as possible
		k.MigrateDelegationsByValidatorIndex(sdk.UnwrapSDKContext(ctx), 10000)
	}()

	return k.BlockValidatorUpdates(ctx)
}
