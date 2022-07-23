package keeper_test

import (
	gocontext "context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

func (suite *IntegrationTestSuite) TestAddSupplyOffset() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient

	test1Supply := sdk.NewInt64Coin("test1", 4000000)
	expectedTotalSupply := sdk.NewCoins(test1Supply)
	suite.
		Require().
		NoError(app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, expectedTotalSupply))

	_, err := queryClient.SupplyOf(gocontext.Background(), &types.QuerySupplyOfRequest{})
	suite.Require().Error(err)

	res, err := queryClient.SupplyOf(gocontext.Background(), &types.QuerySupplyOfRequest{Denom: test1Supply.Denom})
	suite.Require().NoError(err)
	suite.Require().NotNil(res)

	suite.Require().Equal(test1Supply, res.Amount)

	// test total supply of query with supply offset
	err = app.BankKeeper.AddSupplyOffset(ctx, "test1", sdk.NewInt(-1000000))
	suite.Require().NoError(err)
	res, err = queryClient.SupplyOf(gocontext.Background(), &types.QuerySupplyOfRequest{Denom: test1Supply.Denom})
	suite.Require().NoError(err)
	suite.Require().NotNil(res)

	suite.Require().Equal(test1Supply.Sub(sdk.NewInt64Coin("test1", 1000000)), res.Amount)

	// make sure query without offsets hasn't changed
	res2, err := queryClient.SupplyOfWithoutOffset(gocontext.Background(), &types.QuerySupplyOfWithoutOffsetRequest{Denom: test1Supply.Denom})
	suite.Require().NoError(err)
	suite.Require().NotNil(res2)

	suite.Require().Equal(test1Supply, res2.Amount)

	// make sure that the supplyWithoutOffset can't go negative
	err = app.BankKeeper.AddSupplyOffset(ctx, "test1", sdk.NewInt(-4000000))
	suite.Require().Error(err)

	// make sure supply didn't change
	res2, err = queryClient.SupplyOfWithoutOffset(gocontext.Background(), &types.QuerySupplyOfWithoutOffsetRequest{Denom: test1Supply.Denom})
	suite.Require().NoError(err)
	suite.Require().NotNil(res2)
	suite.Require().Equal(test1Supply, res2.Amount)
}
