package keeper_test

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
)

func (suite *IntegrationTestSuite) TestQuerier_QueryBalance() {
	app, ctx := suite.app, suite.ctx
	legacyAmino := app.LegacyAmino()
	_, _, addr := testdata.KeyTestPubAddr()
	req := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", types.ModuleName, types.QueryBalance),
		Data: []byte{},
	}

	querier := keeper.NewQuerier(app.BankKeeper, legacyAmino)

	res, err := querier(ctx, []string{types.QueryBalance}, req)
	suite.Require().NotNil(err)
	suite.Require().Nil(res)

	req.Data = legacyAmino.MustMarshalJSON(types.NewQueryBalanceRequest(addr, fooDenom))
	res, err = querier(ctx, []string{types.QueryBalance}, req)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)

	var balance sdk.Coin
	suite.Require().NoError(legacyAmino.UnmarshalJSON(res, &balance))
	suite.True(balance.IsZero())

	origCoins := sdk.NewCoins(newFooCoin(50), newBarCoin(30))
	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr)

	app.AccountKeeper.SetAccount(ctx, acc)
	suite.Require().NoError(app.BankKeeper.SetBalances(ctx, acc.GetAddress(), origCoins))

	res, err = querier(ctx, []string{types.QueryBalance}, req)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().NoError(legacyAmino.UnmarshalJSON(res, &balance))
	suite.True(balance.IsEqual(newFooCoin(50)))
}

func (suite *IntegrationTestSuite) TestQuerier_QueryAllBalances() {
	app, ctx := suite.app, suite.ctx
	legacyAmino := app.LegacyAmino()
	_, _, addr := testdata.KeyTestPubAddr()
	req := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", types.ModuleName, types.QueryAllBalances),
		Data: []byte{},
	}

	querier := keeper.NewQuerier(app.BankKeeper, legacyAmino)

	res, err := querier(ctx, []string{types.QueryAllBalances}, req)
	suite.Require().NotNil(err)
	suite.Require().Nil(res)

	req.Data = legacyAmino.MustMarshalJSON(types.NewQueryAllBalancesRequest(addr, nil))
	res, err = querier(ctx, []string{types.QueryAllBalances}, req)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)

	var balances sdk.Coins
	suite.Require().NoError(legacyAmino.UnmarshalJSON(res, &balances))
	suite.True(balances.IsZero())

	origCoins := sdk.NewCoins(newFooCoin(50), newBarCoin(30))
	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr)

	app.AccountKeeper.SetAccount(ctx, acc)
	suite.Require().NoError(app.BankKeeper.SetBalances(ctx, acc.GetAddress(), origCoins))

	res, err = querier(ctx, []string{types.QueryAllBalances}, req)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
	suite.Require().NoError(legacyAmino.UnmarshalJSON(res, &balances))
	suite.True(balances.IsEqual(origCoins))
}

func (suite *IntegrationTestSuite) TestQuerier_QueryTotalSupply() {
	SetOsmosisAddressPrefixes()
	defer SetCosmosAddressPrefixes()
	app, ctx := suite.app, suite.ctx
	legacyAmino := app.LegacyAmino()
	uosmoSupply := sdk.NewInt64Coin(keeper.OsmoBondDenom, 100000000)
	uosmoDevRewards := sdk.NewInt64Coin(keeper.OsmoBondDenom, 100000)
	devRewardAddr, err := sdk.AccAddressFromBech32(keeper.DevRewardsAddr)
	suite.Require().NoError(err)
	totalSupply := types.NewSupply(sdk.NewCoins(sdk.NewInt64Coin("test", 400000000), uosmoSupply))
	app.BankKeeper.SetSupply(ctx, totalSupply)
	app.BankKeeper.SetBalance(ctx, devRewardAddr, uosmoDevRewards)

	req := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", types.ModuleName, types.QueryTotalSupply),
		Data: []byte{},
	}

	querier := keeper.NewQuerier(app.BankKeeper, legacyAmino)

	res, err := querier(ctx, []string{types.QueryTotalSupply}, req)
	suite.Require().NotNil(err)
	suite.Require().Nil(res)

	req.Data = legacyAmino.MustMarshalJSON(types.NewQueryTotalSupplyParams(1, 100))
	res, err = querier(ctx, []string{types.QueryTotalSupply}, req)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)

	var resp sdk.Coins
	suite.Require().NoError(legacyAmino.UnmarshalJSON(res, &resp))
	suite.Require().Equal(totalSupply.Total.Sub(sdk.Coins{uosmoDevRewards}), resp)
}

func (suite *IntegrationTestSuite) TestQuerier_QueryTotalSupplyOf() {
	SetOsmosisAddressPrefixes()
	defer SetCosmosAddressPrefixes()
	app, ctx := suite.app, suite.ctx
	legacyAmino := app.LegacyAmino()
	test1Supply := sdk.NewInt64Coin("test1", 4000000)
	test2Supply := sdk.NewInt64Coin("test2", 700000000)
	uosmoSupply := sdk.NewInt64Coin(keeper.OsmoBondDenom, 100000000)
	uosmoDevRewards := sdk.NewInt64Coin(keeper.OsmoBondDenom, 100000)
	devRewardAddr, err := sdk.AccAddressFromBech32(keeper.DevRewardsAddr)
	suite.Require().NoError(err)
	totalSupply := types.NewSupply(sdk.NewCoins(test1Supply, test2Supply, uosmoSupply))
	app.BankKeeper.SetSupply(ctx, totalSupply)
	app.BankKeeper.SetBalance(ctx, devRewardAddr, uosmoDevRewards)

	req := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", types.ModuleName, types.QuerySupplyOf),
		Data: []byte{},
	}

	querier := keeper.NewQuerier(app.BankKeeper, legacyAmino)

	res, err := querier(ctx, []string{types.QuerySupplyOf}, req)
	suite.Require().NotNil(err)
	suite.Require().Nil(res)

	req.Data = legacyAmino.MustMarshalJSON(types.NewQuerySupplyOfParams(test1Supply.Denom))
	res, err = querier(ctx, []string{types.QuerySupplyOf}, req)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)

	var resp sdk.Coin
	suite.Require().NoError(legacyAmino.UnmarshalJSON(res, &resp))
	suite.Require().Equal(test1Supply, resp)

	req.Data = legacyAmino.MustMarshalJSON(types.NewQuerySupplyOfParams(keeper.OsmoBondDenom))
	res, err = querier(ctx, []string{types.QuerySupplyOf}, req)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)

	suite.Require().NoError(legacyAmino.UnmarshalJSON(res, &resp))
	suite.Require().Equal(uosmoSupply.Sub(uosmoDevRewards), resp)
}

func (suite *IntegrationTestSuite) TestQuerierRouteNotFound() {
	app, ctx := suite.app, suite.ctx
	legacyAmino := app.LegacyAmino()
	req := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/invalid", types.ModuleName),
		Data: []byte{},
	}

	querier := keeper.NewQuerier(app.BankKeeper, legacyAmino)
	_, err := querier(ctx, []string{"invalid"}, req)
	suite.Error(err)
}
