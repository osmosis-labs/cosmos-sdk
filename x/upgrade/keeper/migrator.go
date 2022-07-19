package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Migrator defines a wrapper around the x/upgrade keeper for performing consensus
// state migrations.
type Migrator struct {
	Keeper
}

func NewMigrator(k Keeper) Migrator {
	return Migrator{Keeper: k}
}

// Migrate1to2 migrates from consensus version 1 to version 2.
func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	return v2.MigrateStore(ctx, m.storeKey)
}
