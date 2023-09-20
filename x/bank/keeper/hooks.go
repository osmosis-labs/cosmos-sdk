package keeper

import (
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
)

// Implements StakingHooks interface
var _ types.BankHooks = BaseSendKeeper{}

// TrackBeforeSend executes the TrackBeforeSend hook if registered.
func (k BaseSendKeeper) TrackBeforeSend(ctx sdk.Context, from, to sdk.AccAddress, amount sdk.Coins, cosmosMsg wasmvmtypes.CosmosMsg) {
	if k.hooks != nil {
		k.hooks.TrackBeforeSend(ctx, from, to, amount, cosmosMsg)
	}
}

// BlockBeforeSend executes the BlockBeforeSend hook if registered.
func (k BaseSendKeeper) BlockBeforeSend(ctx sdk.Context, from, to sdk.AccAddress, amount sdk.Coins, cosmosMsg wasmvmtypes.CosmosMsg) error {
	if k.hooks != nil {
		return k.hooks.BlockBeforeSend(ctx, from, to, amount, cosmosMsg)
	}
	return nil
}
