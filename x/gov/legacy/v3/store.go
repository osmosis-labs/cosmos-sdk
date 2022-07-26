package v3

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const minInitialDepositPercent uint32 = 25

// MigrateStore performs in-place store migrations for consensus version 3
// in the gov module.
// Please note that this is the first version that switches from using
// SDK versioning (v043 etc) for package names to consensus versioning
// of the gov module.
// The migration includes:
//
// - Setting the Min param in the paramstore
func MigrateStore(ctx sdk.Context, paramstore paramtypes.Subspace) error {
	migrateParamsStore(ctx, paramstore)
	return nil
}

func migrateParamsStore(ctx sdk.Context, paramstore paramtypes.Subspace) {
	var depositParams types.DepositParams
	paramstore.Get(ctx, types.ParamStoreKeyDepositParams, &depositParams)
	depositParams.MinInitialDepositPercent = minInitialDepositPercent
	paramstore.Set(ctx, types.ParamStoreKeyDepositParams, depositParams)
}
