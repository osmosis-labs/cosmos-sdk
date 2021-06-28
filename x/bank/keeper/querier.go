package keeper

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
)

var OsmoBondDenom = "uosmo"
var DevUnvestedRewardsAddr = "osmo1vqy8rqqlydj9wkcyvct9zxl3hc4eqgu3d7hd9k"

// NewQuerier returns a new sdk.Keeper instance.
func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryBalance:
			return queryBalance(ctx, req, k, legacyQuerierCdc)

		case types.QueryAllBalances:
			return queryAllBalance(ctx, req, k, legacyQuerierCdc)

		case types.QueryTotalSupply:
			return queryTotalSupply(ctx, req, k, legacyQuerierCdc)

		case types.QuerySupplyOf:
			return querySupplyOf(ctx, req, k, legacyQuerierCdc)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query endpoint: %s", types.ModuleName, path[0])
		}
	}
}

func queryBalance(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryBalanceRequest

	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	address, err := sdk.AccAddressFromBech32(params.Address)
	if err != nil {
		return nil, err
	}

	balance := k.GetBalance(ctx, address, params.Denom)

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, balance)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryAllBalance(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryAllBalancesRequest

	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	address, err := sdk.AccAddressFromBech32(params.Address)
	if err != nil {
		return nil, err
	}

	balances := k.GetAllBalances(ctx, address)

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, balances)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryTotalSupply(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryTotalSupplyParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	totalSupply := k.GetSupply(ctx).GetTotal()

	start, end := client.Paginate(len(totalSupply), params.Page, params.Limit, 100)
	if start < 0 || end < 0 {
		totalSupply = sdk.Coins{}
	} else {
		totalSupply = totalSupply[start:end]
	}

	if totalSupply.AmountOf(OsmoBondDenom).IsPositive() {
		devUnvestedRewardsAddr, err := sdk.AccAddressFromBech32(DevUnvestedRewardsAddr)
		if err != nil {
			panic(fmt.Errorf("developer rewards address is not parsable: %s", devUnvestedRewardsAddr))
		}
		totalSupply = totalSupply.Sub(sdk.Coins{k.GetBalance(ctx, devUnvestedRewardsAddr, OsmoBondDenom)})
	}

	res, err := legacyQuerierCdc.MarshalJSON(totalSupply)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

func querySupplyOf(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QuerySupplyOfParams

	err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	amount := k.GetSupply(ctx).GetTotal().AmountOf(params.Denom)
	supply := sdk.NewCoin(params.Denom, amount)
	if params.Denom == OsmoBondDenom {
		devUnvestedRewardsAddr, err := sdk.AccAddressFromBech32(DevUnvestedRewardsAddr)
		if err != nil {
			panic(fmt.Errorf("developer rewards address is not parsable: %s", devUnvestedRewardsAddr))
		}
		supply = supply.Sub(k.GetBalance(ctx, devUnvestedRewardsAddr, OsmoBondDenom))
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, supply)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
