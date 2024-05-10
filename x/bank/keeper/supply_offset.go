package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"

	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
)

// GetSupplyOffset retrieves the SupplyOffset from store for a specific denom
func (k BaseViewKeeper) GetSupplyOffset(ctx context.Context, denom string) math.Int {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	supplyOffsetStore := prefix.NewStore(store, types.SupplyOffsetKey)

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

// setSupplyOffset sets the supply offset for the given denom
func (k BaseKeeper) setSupplyOffset(ctx context.Context, denom string, offsetAmount math.Int) {
	intBytes, err := offsetAmount.Marshal()
	if err != nil {
		panic(fmt.Errorf("unable to marshal amount value %v", err))
	}

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	supplyOffsetStore := prefix.NewStore(store, types.SupplyOffsetKey)

	// Bank invariants and IBC requires to remove zero coins.
	if offsetAmount.IsZero() {
		supplyOffsetStore.Delete([]byte(denom))
	} else {
		supplyOffsetStore.Set([]byte(denom), intBytes)
	}
}

// AddSupplyOffset adjusts the current supply offset of a denom by the inputted offsetAmount
func (k BaseKeeper) AddSupplyOffset(ctx context.Context, denom string, offsetAmount math.Int) {
	k.setSupplyOffset(ctx, denom, k.GetSupplyOffset(ctx, denom).Add(offsetAmount))
}

// GetSupplyWithOffset retrieves the Supply of a denom and offsets it by SupplyOffset
// If SupplyWithOffset is negative, it returns 0.  This is because sdk.Coin is not valid
// with a negative amount
func (k BaseKeeper) GetSupplyWithOffset(ctx context.Context, denom string) sdk.Coin {
	supply := k.GetSupply(ctx, denom)
	supply.Amount = supply.Amount.Add(k.GetSupplyOffset(ctx, denom))

	if supply.Amount.IsNegative() {
		supply.Amount = math.ZeroInt()
	}

	return supply
}

// GetPaginatedTotalSupplyWithOffsets queries for the supply with offset, ignoring 0 coins, with a given pagination
func (k BaseKeeper) GetPaginatedTotalSupplyWithOffsets(ctx context.Context, pagination *query.PageRequest) (sdk.Coins, *query.PageResponse, error) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	supplyStore := prefix.NewStore(store, types.SupplyKey)

	supply := sdk.NewCoins()

	pageRes, err := query.Paginate(supplyStore, pagination, func(key, value []byte) error {
		denom := string(key)

		var amount math.Int
		err := amount.Unmarshal(value)
		if err != nil {
			return fmt.Errorf("unable to convert amount string to Int %v", err)
		}

		amount = amount.Add(k.GetSupplyOffset(ctx, denom))

		// `Add` omits the 0 coins addition to the `supply`.
		supply = supply.Add(sdk.NewCoin(denom, amount))
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return supply, pageRes, nil
}

// IterateTotalSupplyWithOffsets iterates over the total supply with offsets calling the given cb (callback) function
// with the balance of each coin.
// The iteration stops if the callback returns true.
func (k BaseViewKeeper) IterateTotalSupplyWithOffsets(ctx context.Context, cb func(sdk.Coin) bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	supplyStore := prefix.NewStore(store, types.SupplyKey)

	iterator := supplyStore.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var amount math.Int
		err := amount.Unmarshal(iterator.Value())
		if err != nil {
			panic(fmt.Errorf("unable to unmarshal supply value %v", err))
		}

		balance := sdk.Coin{
			Denom:  string(iterator.Key()),
			Amount: amount.Add(k.GetSupplyOffset(ctx, string(iterator.Key()))),
		}

		if cb(balance) {
			break
		}
	}
}
