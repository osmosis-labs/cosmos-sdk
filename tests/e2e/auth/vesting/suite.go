package testutil

import (
	"fmt"
	"log"
	"os"

	"github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting/client/cli"
)

type E2ETestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network
}

func NewE2ETestSuite(cfg network.Config) *E2ETestSuite {
	return &E2ETestSuite{cfg: cfg}
}

func (s *E2ETestSuite) SetupSuite() {
	s.T().Log("setting up e2e test suite")

	var err error
	s.network, err = network.New(s.T(), s.T().TempDir(), s.cfg)
	s.Require().NoError(err)
	s.Require().NoError(s.network.WaitForNextBlock())
}

func (s *E2ETestSuite) TearDownSuite() {
	s.T().Log("tearing down e2e test suite")
	s.network.Cleanup()
}

func (s *E2ETestSuite) TestNewMsgCreateVestingAccountCmd() {
	val := s.network.Validators[0]

	testCases := map[string]struct {
		args         []string
		expectErr    bool
		expectedCode uint32
		respType     proto.Message
	}{
		"create a continuous vesting account": {
			args: []string{
				sdk.AccAddress("addr2_______________").String(),
				sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String(),
				"4070908800",
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			expectErr:    false,
			expectedCode: 0,
			respType:     &sdk.TxResponse{},
		},
		"create a delayed vesting account": {
			args: []string{
				sdk.AccAddress("addr3_______________").String(),
				sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String(),
				"4070908800",
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
				fmt.Sprintf("--%s=true", cli.FlagDelayed),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			expectErr:    false,
			expectedCode: 0,
			respType:     &sdk.TxResponse{},
		},
		"invalid address": {
			args: []string{
				"addr4",
				sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String(),
				"4070908800",
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
			},
			expectErr:    true,
			expectedCode: 0,
			respType:     &sdk.TxResponse{},
		},
		"invalid coins": {
			args: []string{
				sdk.AccAddress("addr4_______________").String(),
				"fooo",
				"4070908800",
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
			},
			expectErr:    true,
			expectedCode: 0,
			respType:     &sdk.TxResponse{},
		},
		"invalid end time": {
			args: []string{
				sdk.AccAddress("addr4_______________").String(),
				sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String(),
				"-4070908800",
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
			},
			expectErr:    true,
			expectedCode: 0,
			respType:     &sdk.TxResponse{},
		},
	}

	// Synchronize height between test runs, to ensure sequence numbers are
	// properly updated.
	height, err := s.network.LatestHeight()
	if err != nil {
		s.T().Fatalf("Getting initial latest height: %v", err)
	}
	s.T().Logf("Initial latest height: %d", height)
	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			clientCtx := val.ClientCtx

			bw, err := clitestutil.ExecTestCLICmd(clientCtx, cli.NewMsgCreateVestingAccountCmd(), tc.args)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(bw.Bytes(), tc.respType), bw.String())

				txResp := tc.respType.(*sdk.TxResponse)
				s.Require().NoError(clitestutil.CheckTxCode(s.network, clientCtx, txResp.TxHash, tc.expectedCode))
			}
		})

		next, err := s.network.WaitForHeight(height + 1)
		if err != nil {
			s.T().Fatalf("Waiting for height %d: %v", height+1, err)
		}
		height = next
		s.T().Logf("Height now: %d", height)
	}
}

func (s *E2ETestSuite) TestNewMsgCreatePermanentLockedAccountCmd() {
	val := s.network.Validators[0]

	testCases := map[string]struct {
		args         []string
		expectErr    bool
		expectedCode uint32
		respType     proto.Message
	}{
		"create a permanent locked account": {
			args: []string{
				sdk.AccAddress("addr4_______________").String(),
				sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			expectErr:    false,
			expectedCode: 0,
			respType:     &sdk.TxResponse{},
		},
		"invalid address": {
			args: []string{
				"addr4",
				sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String(),
				"4070908800",
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
			},
			expectErr:    true,
			expectedCode: 0,
			respType:     &sdk.TxResponse{},
		},
		"invalid coins": {
			args: []string{
				sdk.AccAddress("addr4_______________").String(),
				"fooo",
				"4070908800",
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
			},
			expectErr:    true,
			expectedCode: 0,
			respType:     &sdk.TxResponse{},
		},
	}

	// Synchronize height between test runs, to ensure sequence numbers are
	// properly updated.
	height, err := s.network.LatestHeight()
	s.Require().NoError(err, "Getting initial latest height")
	s.T().Logf("Initial latest height: %d", height)
	for name, tc := range testCases {
		tc := tc

		s.Run(name, func() {
			clientCtx := val.ClientCtx

			bw, err := clitestutil.ExecTestCLICmd(clientCtx, cli.NewMsgCreatePermanentLockedAccountCmd(), tc.args)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(bw.Bytes(), tc.respType), bw.String())

				txResp := tc.respType.(*sdk.TxResponse)
				s.Require().NoError(clitestutil.CheckTxCode(s.network, clientCtx, txResp.TxHash, tc.expectedCode))
			}
		})
		next, err := s.network.WaitForHeight(height + 1)
		s.Require().NoError(err, "Waiting for height...")
		height = next
		s.T().Logf("Height now: %d", height)
	}
}

func (s *E2ETestSuite) TestNewMsgCreateClawbackVestingAccountCmd() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pwd)
	val := s.network.Validators[0]
	for _, tc := range []struct {
		name         string
		args         []string
		expectErr    bool
		expectedCode uint32
		respType     proto.Message
	}{
		{
			name: "basic",
			args: []string{
				sdk.AccAddress("addr10______________").String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
				fmt.Sprintf("--%s=%s", cli.FlagLockup, fmt.Sprintf(pwd+"/testdata/periods1.json")),
				fmt.Sprintf("--%s=%s", cli.FlagVesting, fmt.Sprintf(pwd+"/testdata/periods1.json")),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			expectErr:    false,
			expectedCode: 0,
			respType:     &sdk.TxResponse{},
		},
		{
			name: "defaultLockup",
			args: []string{
				sdk.AccAddress("addr11______________").String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
				fmt.Sprintf("--%s=%s", cli.FlagVesting, fmt.Sprintf(pwd+"/testdata/periods1.json")),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			expectErr:    false,
			expectedCode: 0,
			respType:     &sdk.TxResponse{},
		},
		{
			name: "defaultVesting",
			args: []string{
				sdk.AccAddress("addr12______________").String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
				fmt.Sprintf("--%s=%s", cli.FlagLockup, fmt.Sprintf(pwd+"/testdata/periods1.json")),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			expectErr:    false,
			expectedCode: 0,
			respType:     &sdk.TxResponse{},
		},
		{
			name: "merge",
			args: []string{
				sdk.AccAddress("addr10______________").String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
				fmt.Sprintf("--%s=%s", cli.FlagLockup, fmt.Sprintf(pwd+"/testdata/periods1.json")),
				fmt.Sprintf("--%s=%s", cli.FlagVesting, fmt.Sprintf(pwd+"/testdata/periods1.json")),
				fmt.Sprintf("--%s=%s", cli.FlagMerge, "true"),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			expectErr:    false,
			expectedCode: 0,
			respType:     &sdk.TxResponse{},
		},
		{
			name: "bad vesting addr",
			args: []string{
				"foo",
			},
			expectErr: true,
		},
		{
			name: "no files",
			args: []string{
				sdk.AccAddress("addr13______________").String(),
			},
			expectErr: true,
		},
		{
			name: "bad lockup filename",
			args: []string{
				sdk.AccAddress("addr13______________").String(),
				fmt.Sprintf("--%s=%s", cli.FlagLockup, fmt.Sprintf(pwd+"testdata/noexist")),
			},
			expectErr: true,
		},
		{
			name: "bad lockup json",
			args: []string{
				sdk.AccAddress("addr13______________").String(),
				fmt.Sprintf("--%s=%s", cli.FlagLockup, fmt.Sprintf(pwd+"testdata/badjson")),
			},
			expectErr: true,
		},
		{
			name: "bad lockup periods",
			args: []string{
				sdk.AccAddress("addr13______________").String(),
				fmt.Sprintf("--%s=%s", cli.FlagLockup, fmt.Sprintf(pwd+"testdata/badperiod.json")),
			},
			expectErr: true,
		},
		{
			name: "bad vesting filename",
			args: []string{
				sdk.AccAddress("addr13______________").String(),
				fmt.Sprintf("--%s=%s", cli.FlagVesting, fmt.Sprintf(pwd+"testdata/noexist")),
			},
			expectErr: true,
		},
		{
			name: "bad vesting json",
			args: []string{
				sdk.AccAddress("addr13______________").String(),
				fmt.Sprintf("--%s=%s", cli.FlagVesting, fmt.Sprintf(pwd+"testdata/badjson")),
			},
			expectErr: true,
		},
		{
			name: "bad vesting periods",
			args: []string{
				sdk.AccAddress("addr13______________").String(),
				fmt.Sprintf("--%s=%s", cli.FlagVesting, fmt.Sprintf(pwd+"testdata/badperiod.json")),
			},
			expectErr: true,
		},
	} {
		s.Run(tc.name, func() {
			clientCtx := val.ClientCtx

			bw, err := clitestutil.ExecTestCLICmd(clientCtx, cli.NewMsgCreateClawbackVestingAccountCmd(), tc.args)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(bw.Bytes(), tc.respType), bw.String())

				txResp := tc.respType.(*sdk.TxResponse)
				s.Require().Equal(tc.expectedCode, txResp.Code)
			}
		})
	}
}

func (s *E2ETestSuite) TestNewMsgClawbackCmd() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	val := s.network.Validators[0]
	addr := sdk.AccAddress("addr30______________")

	_, err = clitestutil.ExecTestCLICmd(val.ClientCtx, cli.NewMsgCreateClawbackVestingAccountCmd(), []string{
		addr.String(),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
		fmt.Sprintf("--%s=%s", cli.FlagLockup, fmt.Sprintf(pwd+"/testdata/periods1.json")),
		fmt.Sprintf("--%s=%s", cli.FlagVesting, fmt.Sprintf(pwd+"/testdata/periods1.json")),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	})
	s.Require().NoError(err)

	for _, tc := range []struct {
		name         string
		args         []string
		expectErr    bool
		expectedCode uint32
		respType     proto.Message
	}{
		{
			name: "basic",
			args: []string{
				addr.String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
				fmt.Sprintf("--%s=%s", cli.FlagDest, sdk.AccAddress("addr32______________").String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			expectErr:    false,
			expectedCode: 0,
			respType:     &sdk.TxResponse{},
		},
		{
			name: "bad vesting addr",
			args: []string{
				"foo",
			},
			expectErr: true,
		},
		{
			name: "bad dest addr",
			args: []string{
				addr.String(),
				fmt.Sprintf("--%s=%s", cli.FlagDest, "bar"),
			},
			expectErr: true,
		},
		{
			name: "default dest",
			args: []string{
				addr.String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
			},
			expectErr:    false,
			expectedCode: 0,
			respType:     &sdk.TxResponse{},
		},
	} {
		s.Run(tc.name, func() {
			clientCtx := val.ClientCtx

			bw, err := clitestutil.ExecTestCLICmd(clientCtx, cli.NewMsgClawbackCmd(), tc.args)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(bw.Bytes(), tc.respType), bw.String())

				txResp := tc.respType.(*sdk.TxResponse)
				s.Require().Equal(tc.expectedCode, txResp.Code)
			}
		})
	}
}
