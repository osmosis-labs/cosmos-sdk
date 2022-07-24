package keeper_test

import (
	gocontext "context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

func (suite *IntegrationTestSuite) TestAddSupplyOffset() {

	testCases := []struct {
		description        string
		initialSupply      sdk.Int
		supplyOffsetAmount sdk.Int
		expPass            bool
		supplyWithOffset   sdk.Int
	}{
		{
			description:        "test positive supply offset",
			initialSupply:      sdk.NewInt(10000),
			supplyOffsetAmount: sdk.NewInt(1000),
			expPass:            true,
			supplyWithOffset:   sdk.NewInt(11000),
		},
		{
			description:        "test negative supply offset",
			initialSupply:      sdk.NewInt(10000),
			supplyOffsetAmount: sdk.NewInt(-1000),
			expPass:            true,
			supplyWithOffset:   sdk.NewInt(9000),
		},
		{
			description:        "make sure that supplyWithOffset can't go negative",
			initialSupply:      sdk.NewInt(10000),
			supplyOffsetAmount: sdk.NewInt(-11000),
			expPass:            false,
			supplyWithOffset:   sdk.NewInt(10000),
		},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(tc.description, func() {
			app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient

			// mint the initialSupply coins
			test1Supply := sdk.NewCoin("test1", tc.initialSupply)
			suite.
				Require().
				NoError(app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, sdk.NewCoins(test1Supply)))

			// make sure that queried supply is correct
			res, err := queryClient.SupplyOf(gocontext.Background(), &types.QuerySupplyOfRequest{Denom: test1Supply.Denom})
			suite.Require().NoError(err)
			suite.Require().NotNil(res)
			suite.Require().Equal(test1Supply, res.Amount)

			// try adding the supply offset
			err = app.BankKeeper.AddSupplyOffset(ctx, "test1", tc.supplyOffsetAmount)
			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}

			// test querying supply results in expected change
			res, err = queryClient.SupplyOf(gocontext.Background(), &types.QuerySupplyOfRequest{Denom: test1Supply.Denom})
			suite.Require().NoError(err)
			suite.Require().NotNil(res)
			suite.Require().Equal(tc.supplyWithOffset, res.Amount.Amount)

			// test querying supply without offset did not change
			res2, err := queryClient.SupplyOfWithoutOffset(gocontext.Background(), &types.QuerySupplyOfWithoutOffsetRequest{Denom: test1Supply.Denom})
			suite.Require().NoError(err)
			suite.Require().NotNil(res2)
			suite.Require().Equal(test1Supply, res2.Amount)
		})
	}
}
