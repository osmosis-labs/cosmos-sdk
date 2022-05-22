package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

// SetProtocolVersion sets the protocol version to version. Used for testing the upgrade keeper.
func (k *Keeper) SetProtocolVersion(ctx sdk.Context, version uint64) error {
	return k.versionManager.SetProtocolVersion(ctx, version)
}
