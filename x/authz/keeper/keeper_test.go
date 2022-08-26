package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

var (
	bankSendAuthMsgType = banktypes.SendAuthorization{}.MsgTypeURL()
	coins10             = sdk.NewCoins(sdk.NewInt64Coin("stake", 10))
)

type TestSuite struct {
	suite.Suite

	app         *simapp.SimApp
	ctx         sdk.Context
	addrs       []sdk.AccAddress
	queryClient authz.QueryClient
}

func (s *TestSuite) SetupTest() {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	now := tmtime.Now()
	ctx = ctx.WithBlockHeader(tmproto.Header{Time: now})
	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	authz.RegisterQueryServer(queryHelper, app.AuthzKeeper)
	queryClient := authz.NewQueryClient(queryHelper)
	s.queryClient = queryClient

	s.app = app
	s.ctx = ctx
	s.queryClient = queryClient
	s.addrs = simapp.AddTestAddrsIncremental(app, ctx, 7, sdk.NewInt(30000000))
}

func (s *TestSuite) TestKeeper() {
	app, ctx, addrs := s.app, s.ctx, s.addrs

	granterAddr := addrs[0]
	granteeAddr := addrs[1]
	recipientAddr := addrs[2]

	s.T().Log("verify that no authorization returns nil")
	authorization, expiration := app.AuthzKeeper.GetCleanAuthorization(ctx, granteeAddr, granterAddr, bankSendAuthMsgType)
	s.Require().Nil(authorization)
	s.Require().Equal(expiration, time.Time{})
	now := s.ctx.BlockHeader().Time

	newCoins := sdk.NewCoins(sdk.NewInt64Coin("steak", 100))
	s.T().Log("verify if expired authorization is rejected")
	x := &banktypes.SendAuthorization{SpendLimit: newCoins}
	err := app.AuthzKeeper.SaveGrant(ctx, granterAddr, granteeAddr, x, now.Add(-1*time.Hour))
	s.Require().Error(err)
	authorization, _ = app.AuthzKeeper.GetCleanAuthorization(ctx, granteeAddr, granterAddr, bankSendAuthMsgType)
	s.Require().Nil(authorization)

	s.T().Log("verify if authorization is accepted")
	x = &banktypes.SendAuthorization{SpendLimit: newCoins}
	err = app.AuthzKeeper.SaveGrant(ctx, granteeAddr, granterAddr, x, now.Add(time.Hour))
	s.Require().NoError(err)
	authorization, _ = app.AuthzKeeper.GetCleanAuthorization(ctx, granteeAddr, granterAddr, bankSendAuthMsgType)
	s.Require().NotNil(authorization)
	s.Require().Equal(authorization.MsgTypeURL(), bankSendAuthMsgType)

	s.T().Log("verify fetching authorization with wrong msg type fails")
	authorization, _ = app.AuthzKeeper.GetCleanAuthorization(ctx, granteeAddr, granterAddr, sdk.MsgTypeURL(&govtypes.MsgDeposit{}))
	s.Require().Nil(authorization)

	s.T().Log("verify fetching authorization with wrong grantee fails")
	authorization, _ = app.AuthzKeeper.GetCleanAuthorization(ctx, recipientAddr, granterAddr, bankSendAuthMsgType)
	s.Require().Nil(authorization)

	s.T().Log("verify revoke fails with wrong information")
	err = app.AuthzKeeper.DeleteGrant(ctx, recipientAddr, granterAddr, bankSendAuthMsgType)
	s.Require().Error(err)
	authorization, _ = app.AuthzKeeper.GetCleanAuthorization(ctx, recipientAddr, granterAddr, bankSendAuthMsgType)
	s.Require().Nil(authorization)

	s.T().Log("verify revoke executes with correct information")
	err = app.AuthzKeeper.DeleteGrant(ctx, granteeAddr, granterAddr, bankSendAuthMsgType)
	s.Require().NoError(err)
	authorization, _ = app.AuthzKeeper.GetCleanAuthorization(ctx, granteeAddr, granterAddr, bankSendAuthMsgType)
	s.Require().Nil(authorization)

}

func (s *TestSuite) TestKeeperIter() {
	app, ctx, addrs := s.app, s.ctx, s.addrs

	granterAddr := addrs[0]
	granteeAddr := addrs[1]

	s.T().Log("verify that no authorization returns nil")
	authorization, expiration := app.AuthzKeeper.GetCleanAuthorization(ctx, granteeAddr, granterAddr, "Abcd")
	s.Require().Nil(authorization)
	s.Require().Equal(time.Time{}, expiration)
	now := s.ctx.BlockHeader().Time.Add(time.Second)

	newCoins := sdk.NewCoins(sdk.NewInt64Coin("steak", 100))
	s.T().Log("verify if expired authorization is rejected")
	x := &banktypes.SendAuthorization{SpendLimit: newCoins}
	err := app.AuthzKeeper.SaveGrant(ctx, granteeAddr, granterAddr, x, now.Add(-1*time.Hour))
	s.Require().Error(err)
	authorization, _ = app.AuthzKeeper.GetCleanAuthorization(ctx, granteeAddr, granterAddr, "abcd")
	s.Require().Nil(authorization)

	app.AuthzKeeper.IterateGrants(ctx, func(granter, grantee sdk.AccAddress, grant authz.Grant) bool {
		s.Require().Equal(granter, granterAddr)
		s.Require().Equal(grantee, granteeAddr)
		return true
	})

}

func (s *TestSuite) TestKeeperFees() {
	app, addrs := s.app, s.addrs

	granterAddr := addrs[0]
	granteeAddr := addrs[1]
	recipientAddr := addrs[2]
	s.Require().NoError(simapp.FundAccount(app.BankKeeper, s.ctx, granterAddr, sdk.NewCoins(sdk.NewInt64Coin("steak", 10000))))
	expiration := s.ctx.BlockHeader().Time.Add(1 * time.Second)

	smallCoin := sdk.NewCoins(sdk.NewInt64Coin("steak", 20))
	someCoin := sdk.NewCoins(sdk.NewInt64Coin("steak", 123))

	msgs := authz.NewMsgExec(granteeAddr, []sdk.Msg{
		&banktypes.MsgSend{
			Amount:      sdk.NewCoins(sdk.NewInt64Coin("steak", 2)),
			FromAddress: granterAddr.String(),
			ToAddress:   recipientAddr.String(),
		},
	})

	s.Require().NoError(msgs.UnpackInterfaces(app.AppCodec()))

	s.T().Log("verify dispatch fails with invalid authorization")
	executeMsgs, err := msgs.GetMessages()
	s.Require().NoError(err)
	result, err := app.AuthzKeeper.DispatchActions(s.ctx, granteeAddr, executeMsgs)

	s.Require().Nil(result)
	s.Require().NotNil(err)

	s.T().Log("verify dispatch executes with correct information")
	// grant authorization
	err = app.AuthzKeeper.SaveGrant(s.ctx, granteeAddr, granterAddr, &banktypes.SendAuthorization{SpendLimit: smallCoin}, expiration)
	s.Require().NoError(err)
	authorization, _ := app.AuthzKeeper.GetCleanAuthorization(s.ctx, granteeAddr, granterAddr, bankSendAuthMsgType)
	s.Require().NotNil(authorization)

	s.Require().Equal(authorization.MsgTypeURL(), bankSendAuthMsgType)

	executeMsgs, err = msgs.GetMessages()
	s.Require().NoError(err)

	result, err = app.AuthzKeeper.DispatchActions(s.ctx, granteeAddr, executeMsgs)
	s.Require().NoError(err)
	s.Require().NotNil(result)

	authorization, _ = app.AuthzKeeper.GetCleanAuthorization(s.ctx, granteeAddr, granterAddr, bankSendAuthMsgType)
	s.Require().NotNil(authorization)

	s.T().Log("verify dispatch fails with overlimit")
	// grant authorization

	msgs = authz.NewMsgExec(granteeAddr, []sdk.Msg{
		&banktypes.MsgSend{
			Amount:      someCoin,
			FromAddress: granterAddr.String(),
			ToAddress:   recipientAddr.String(),
		},
	})

	s.Require().NoError(msgs.UnpackInterfaces(app.AppCodec()))
	executeMsgs, err = msgs.GetMessages()
	s.Require().NoError(err)

	result, err = app.AuthzKeeper.DispatchActions(s.ctx, granteeAddr, executeMsgs)
	s.Require().Nil(result)
	s.Require().NotNil(err)

	authorization, _ = app.AuthzKeeper.GetCleanAuthorization(s.ctx, granteeAddr, granterAddr, bankSendAuthMsgType)
	s.Require().NotNil(authorization)
}

// Tests that all msg events included in an authz MsgExec tx
// Ref: https://github.com/cosmos/cosmos-sdk/issues/9501
func (s *TestSuite) TestDispatchedEvents() {
	require := s.Require()
	app, addrs := s.app, s.addrs
	granterAddr := addrs[0]
	granteeAddr := addrs[1]
	recipientAddr := addrs[2]
	require.NoError(simapp.FundAccount(app.BankKeeper, s.ctx, granterAddr, sdk.NewCoins(sdk.NewInt64Coin("steak", 10000))))
	expiration := s.ctx.BlockHeader().Time.Add(1 * time.Second)

	smallCoin := sdk.NewCoins(sdk.NewInt64Coin("steak", 20))
	msgs := authz.NewMsgExec(granteeAddr, []sdk.Msg{
		&banktypes.MsgSend{
			Amount:      sdk.NewCoins(sdk.NewInt64Coin("steak", 2)),
			FromAddress: granterAddr.String(),
			ToAddress:   recipientAddr.String(),
		},
	})

	// grant authorization
	err := app.AuthzKeeper.SaveGrant(s.ctx, granteeAddr, granterAddr, &banktypes.SendAuthorization{SpendLimit: smallCoin}, expiration)
	require.NoError(err)
	authorization, _ := app.AuthzKeeper.GetCleanAuthorization(s.ctx, granteeAddr, granterAddr, bankSendAuthMsgType)
	require.NotNil(authorization)
	require.Equal(authorization.MsgTypeURL(), bankSendAuthMsgType)

	executeMsgs, err := msgs.GetMessages()
	require.NoError(err)

	result, err := app.AuthzKeeper.DispatchActions(s.ctx, granteeAddr, executeMsgs)
	require.NoError(err)
	require.NotNil(result)
	events := s.ctx.EventManager().Events()
	// get last 5 events (events that occur *after* the grant)
	events = events[len(events)-5:]
	requiredEvents := map[string]bool{
		"coin_spent":    false,
		"coin_received": false,
		"transfer":      false,
		"message":       false,
	}
	for _, e := range events {
		requiredEvents[e.Type] = true
	}
	for _, v := range requiredEvents {
		require.True(v)
	}
}

func (s *TestSuite) TestGetAuthorization() {
	addr1 := s.addrs[3]
	addr2 := s.addrs[4]
	addr3 := s.addrs[5]
	addr4 := s.addrs[6]

	genAuthMulti := authz.NewGenericAuthorization(sdk.MsgTypeURL(&banktypes.MsgMultiSend{}))
	genAuthSend := authz.NewGenericAuthorization(sdk.MsgTypeURL(&banktypes.MsgSend{}))
	sendAuth := banktypes.NewSendAuthorization(coins10)

	start := s.ctx.BlockHeader().Time
	expired := start.Add(time.Duration(1) * time.Second)
	notExpired := start.Add(time.Duration(5) * time.Hour)

	s.Require().NoError(s.app.AuthzKeeper.SaveGrant(s.ctx, addr1, addr2, genAuthMulti, nil), "creating grant 1->2")
	s.Require().NoError(s.app.AuthzKeeper.SaveGrant(s.ctx, addr1, addr3, genAuthSend, &expired), "creating grant 1->3")
	s.Require().NoError(s.app.AuthzKeeper.SaveGrant(s.ctx, addr1, addr4, sendAuth, &notExpired), "creating grant 1->4")
	// Without access to private keeper methods, I don't know how to save a grant with an invalid authorization.
	newCtx := s.ctx.WithBlockTime(start.Add(time.Duration(1) * time.Minute))

	tests := []struct {
		name    string
		grantee sdk.AccAddress
		granter sdk.AccAddress
		msgType string
		expAuth authz.Authorization
		expExp  *time.Time
	}{
		{
			name:    "grant has nil exp and is returned",
			grantee: addr1,
			granter: addr2,
			msgType: genAuthMulti.MsgTypeURL(),
			expAuth: genAuthMulti,
			expExp:  nil,
		},
		{
			name:    "grant is expired not returned",
			grantee: addr1,
			granter: addr3,
			msgType: genAuthSend.MsgTypeURL(),
			expAuth: nil,
			expExp:  nil,
		},
		{
			name:    "grant is not expired and is returned",
			grantee: addr1,
			granter: addr4,
			msgType: sendAuth.MsgTypeURL(),
			expAuth: sendAuth,
			expExp:  &notExpired,
		},
		{
			name:    "grant is not expired but wrong msg type returns nil",
			grantee: addr1,
			granter: addr4,
			msgType: genAuthMulti.MsgTypeURL(),
			expAuth: nil,
			expExp:  nil,
		},
		{
			name:    "no grant exists between the two",
			grantee: addr2,
			granter: addr3,
			msgType: genAuthSend.MsgTypeURL(),
			expAuth: nil,
			expExp:  nil,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			actAuth, actExp := s.app.AuthzKeeper.GetAuthorization(newCtx, tc.grantee, tc.granter, tc.msgType)
			s.Assert().Equal(tc.expAuth, actAuth, "authorization")
			s.Assert().Equal(tc.expExp, actExp, "expiration")
		})
	}
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
