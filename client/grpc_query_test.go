package client_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

type IntegrationTestSuite struct {
	suite.Suite

	network *network.Network
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	s.network = network.New(s.T(), network.DefaultConfig())
	s.Require().NotNil(s.network)

	_, err := s.network.WaitForHeight(2)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func (s *IntegrationTestSuite) TestGRPCQuery() {
	val0 := s.network.Validators[0]

	// gRPC query to test service should work
	testClient := testdata.NewQueryClient(val0.ClientCtx)
	testRes, err := testClient.Echo(context.Background(), &testdata.EchoRequest{Message: "hello"})
	s.Require().NoError(err)
	s.Require().Equal("hello", testRes.Message)

	// gRPC query to bank service should work
	denom := fmt.Sprintf("%stoken", val0.Moniker)
	bankClient := banktypes.NewQueryClient(val0.ClientCtx)
	var header metadata.MD
	bankRes, err := bankClient.Balance(
		context.Background(),
		&banktypes.QueryBalanceRequest{Address: val0.Address.String(), Denom: denom},
		grpc.Header(&header), // Also fetch grpc header
	)
	s.Require().NoError(err)
	s.Require().Equal(
		sdk.NewCoin(denom, s.network.Config.AccountTokens),
		*bankRes.GetBalance(),
	)
	blockHeight := header.Get(grpctypes.GRPCBlockHeightHeader)
	s.Require().NotEmpty(blockHeight[0]) // Should contain the block height

	// Request metadata should work
	val0.ClientCtx = val0.ClientCtx.WithHeight(1) // We set clientCtx to height 1
	bankClient = banktypes.NewQueryClient(val0.ClientCtx)
	bankRes, err = bankClient.Balance(
		context.Background(),
		&banktypes.QueryBalanceRequest{Address: val0.Address.String(), Denom: denom},
		grpc.Header(&header),
	)
	blockHeight = header.Get(grpctypes.GRPCBlockHeightHeader)
	s.Require().Equal([]string{"1"}, blockHeight)
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func TestSelectHeight(t *testing.T) {
	// if height is set to this, the tests assume that it is not set
	const heightNotSetFlag = int64(-1)

	testcases := map[string]struct {
		clientContextHeight int64
		grpcHeight          int64
		expectedHeight      int64
	}{
		"clientContextHeight 1; grpcHeight not set - clientContextHeight selected": {
			clientContextHeight: 1,
			grpcHeight:          heightNotSetFlag,
			expectedHeight:      1,
		},
		"clientContextHeight not set; grpcHeight is 2 - grpcHeight is chosen": {
			clientContextHeight: heightNotSetFlag,
			grpcHeight:          2,
			expectedHeight:      2,
		},
		"both not set - 0 returned": {
			clientContextHeight: heightNotSetFlag,
			grpcHeight:          heightNotSetFlag,
			expectedHeight:      0,
		},
		"clientContextHeight 3; grpcHeight is 0 - grpcHeight is chosen": {
			clientContextHeight: 3,
			grpcHeight:          0,
			expectedHeight:      0,
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			clientCtx := client.Context{}
			if tc.clientContextHeight != heightNotSetFlag {
				clientCtx = clientCtx.WithHeight(tc.clientContextHeight)
			}

			grpcContxt := context.Background()
			if tc.grpcHeight != heightNotSetFlag {
				header := metadata.Pairs(grpctypes.GRPCBlockHeightHeader, fmt.Sprintf("%d", tc.grpcHeight))
				grpcContxt = metadata.NewOutgoingContext(grpcContxt, header)
			}

			height, err := client.SelectHeight(clientCtx, grpcContxt)
			require.NoError(t, err)
			require.Equal(t, tc.expectedHeight, height)
		})
	}
}
