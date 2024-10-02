package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"

	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
)

// NOTE: All these functions should only be used in the v27 migration
// this file should be removed completely after the migration

// GetSupplyOffset retrieves the SupplyOffset from store for a specific denom using the pre v26 (old) key
// TODO: Remove after v27 migration
func (k BaseViewKeeper) GetSupplyOffsetOld(ctx context.Context, denom string) math.Int {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	supplyOffsetStore := prefix.NewStore(store, types.SupplyOffsetKeyOld)

	bz := supplyOffsetStore.Get([]byte(denom))
	if bz == nil {
		return math.NewInt(0)
	}

	var amount math.Int
	err := amount.Unmarshal(bz)
	if err != nil {
		panic(fmt.Errorf("unable to unmarshal supply offset value %v", err))
	}

	return amount
}

// RemoveOldSupplyOffset removes the old supply offset key
// TODO: Remove after v27 migration
func (k BaseViewKeeper) RemoveOldSupplyOffset(ctx context.Context, denom string) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	supplyOffsetStore := prefix.NewStore(store, types.SupplyOffsetKeyOld)

	supplyOffsetStore.Delete([]byte(denom))
}

// setSupplyOffsetOld sets the supply offset for the given denom using the pre v26 (old) key
// TODO: Remove after v27 migration
func (k BaseKeeper) setSupplyOffsetOld(ctx context.Context, denom string, offsetAmount math.Int) {
	intBytes, err := offsetAmount.Marshal()
	if err != nil {
		panic(fmt.Errorf("unable to marshal amount value %v", err))
	}

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	supplyOffsetStore := prefix.NewStore(store, types.SupplyOffsetKeyOld)

	// Bank invariants and IBC requires to remove zero coins.
	if offsetAmount.IsZero() {
		supplyOffsetStore.Delete([]byte(denom))
	} else {
		supplyOffsetStore.Set([]byte(denom), intBytes)
	}
}

// AddSupplyOffset adjusts the current supply offset of a denom by the inputted offsetAmount using the pre v26 (old) key
// TODO: Remove after v27 migration
func (k BaseKeeper) AddSupplyOffsetOld(ctx context.Context, denom string, offsetAmount math.Int) {
	k.setSupplyOffsetOld(ctx, denom, k.GetSupplyOffsetOld(ctx, denom).Add(offsetAmount))
}
