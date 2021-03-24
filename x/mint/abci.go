package mint

import (
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mint/keeper"
	"github.com/cosmos/cosmos-sdk/x/mint/types"
)

// BeginBlocker mints new tokens for the previous block.
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	// Check if we are at an epoch boundary. If not, exit early
	epochDuration := k.GetParams(ctx).EpochDuration
	nextEpochTimeEst := k.GetLastEpochTime(ctx).Add(epochDuration)
	if ctx.BlockTime().Before(nextEpochTimeEst) {
		return
	}

	k.SetLastEpochTime(ctx, ctx.BlockTime())
	k.SetEpochNum(ctx, k.GetEpochNum(ctx)+1)

	// fetch stored minter & params
	minter := k.GetMinter(ctx)
	params := k.GetParams(ctx)

	// Check if we have hit an epoch where we update the inflation parameter.
	// Since epochs only update based on BFT time data, it is safe to store the "halvening period time"
	// in terms of the number of epochs that have transpired.
	if k.GetEpochNum(ctx) >= k.GetParams(ctx).HalvenPeriodInEpoch+k.GetLastHalvenEpochNum(ctx) {
		// Halven the reward per halven period
		minter.AnnualProvisions = minter.NextAnnualProvisions(params)
		k.SetMinter(ctx, minter)
		k.SetLastHalvenEpochNum(ctx, k.GetEpochNum(ctx))
	}

	// mint coins, update supply
	mintedCoin := minter.EpochProvision(params)
	mintedCoins := sdk.NewCoins(mintedCoin)

	err := k.MintCoins(ctx, mintedCoins)
	if err != nil {
		panic(err)
	}

	// send the minted coins to the fee collector account
	err = k.AddCollectedFees(ctx, mintedCoins)
	if err != nil {
		panic(err)
	}

	if mintedCoin.Amount.IsInt64() {
		defer telemetry.ModuleSetGauge(types.ModuleName, float32(mintedCoin.Amount.Int64()), "minted_tokens")
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMint,
			sdk.NewAttribute(types.AttributeKeyAnnualProvisions, minter.AnnualProvisions.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, mintedCoin.Amount.String()),
		),
	)
}
