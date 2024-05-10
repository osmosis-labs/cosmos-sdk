package bank_test

import (
	"context"
	"fmt"
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/stretchr/testify/require"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktestutil "github.com/cosmos/cosmos-sdk/x/bank/testutil"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

var _ types.BankHooks = &MockBankHooksReceiver{}

// BankHooks event hooks for bank (noalias)
type MockBankHooksReceiver struct{}

// Mock BlockBeforeSend bank hook that doesn't allow the sending of exactly 100 coins of any denom.
func (h *MockBankHooksReceiver) BlockBeforeSend(ctx context.Context, from, to sdk.AccAddress, amount sdk.Coins) error {
	for _, coin := range amount {
		if coin.Amount.Equal(math.NewInt(100)) {
			return fmt.Errorf("not allowed; expected %v, got: %v", 100, coin.Amount)
		}
	}
	return nil
}

// variable for counting `TrackBeforeSend`
var (
	countTrackBeforeSend = 0
	expNextCount         = 1
)

// Mock TrackBeforeSend bank hook that simply tracks the sending of exactly 50 coins of any denom.
func (h *MockBankHooksReceiver) TrackBeforeSend(ctx context.Context, from, to sdk.AccAddress, amount sdk.Coins) {
	for _, coin := range amount {
		if coin.Amount.Equal(math.NewInt(50)) {
			countTrackBeforeSend++
		}
	}
}

func TestHooks(t *testing.T) {
	acc := &authtypes.BaseAccount{
		Address: addr1.String(),
	}

	genAccs := []authtypes.GenesisAccount{acc}
	app := createTestSuite(t, genAccs)
	baseApp := app.App.BaseApp
	ctx := baseApp.NewContextLegacy(false, tmproto.Header{})

	require.NoError(t, banktestutil.FundAccount(ctx, app.BankKeeper, addr1, sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 10000))))
	require.NoError(t, banktestutil.FundAccount(ctx, app.BankKeeper, addr2, sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 10000))))
	banktestutil.FundModuleAccount(ctx, app.BankKeeper, stakingtypes.BondedPoolName, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1000))))

	// create a valid send amount which is 1 coin, and an invalidSendAmount which is 100 coins
	validSendAmount := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1)))
	triggerTrackSendAmount := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(50)))
	invalidBlockSendAmount := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(100)))

	// setup our mock bank hooks receiver that prevents the send of 100 coins
	bankHooksReceiver := MockBankHooksReceiver{}
	baseBankKeeper, ok := app.BankKeeper.(keeper.BaseKeeper)
	require.True(t, ok)
	baseBankKeeper.SetHooks(
		types.NewMultiBankHooks(&bankHooksReceiver),
	)
	app.BankKeeper = baseBankKeeper

	// try sending a validSendAmount and it should work
	err := app.BankKeeper.SendCoins(ctx, addr1, addr2, validSendAmount)
	require.NoError(t, err)

	// try sending an trigger track send amount and it should work
	err = app.BankKeeper.SendCoins(ctx, addr1, addr2, triggerTrackSendAmount)
	require.NoError(t, err)

	require.Equal(t, countTrackBeforeSend, expNextCount)
	expNextCount++

	// try sending an invalidSendAmount and it should not work
	err = app.BankKeeper.SendCoins(ctx, addr1, addr2, invalidBlockSendAmount)
	require.Error(t, err)

	// try doing SendManyCoins and make sure if even a single subsend is invalid, the entire function fails
	err = app.BankKeeper.SendManyCoins(ctx, addr1, []sdk.AccAddress{addr1, addr2}, []sdk.Coins{invalidBlockSendAmount, validSendAmount})
	require.Error(t, err)

	err = app.BankKeeper.SendManyCoins(ctx, addr1, []sdk.AccAddress{addr1, addr2}, []sdk.Coins{triggerTrackSendAmount, validSendAmount})
	require.NoError(t, err)
	require.Equal(t, countTrackBeforeSend, expNextCount)
	expNextCount++

	// make sure that account to module doesn't bypass hook
	err = app.BankKeeper.SendCoinsFromAccountToModule(ctx, addr1, stakingtypes.BondedPoolName, validSendAmount)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromAccountToModule(ctx, addr1, stakingtypes.BondedPoolName, invalidBlockSendAmount)
	require.Error(t, err)
	err = app.BankKeeper.SendCoinsFromAccountToModule(ctx, addr1, stakingtypes.BondedPoolName, triggerTrackSendAmount)
	require.NoError(t, err)
	require.Equal(t, countTrackBeforeSend, expNextCount)
	expNextCount++

	// make sure that module to account doesn't bypass hook
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, stakingtypes.BondedPoolName, addr1, validSendAmount)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, stakingtypes.BondedPoolName, addr1, invalidBlockSendAmount)
	require.Error(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, stakingtypes.BondedPoolName, addr1, triggerTrackSendAmount)
	require.NoError(t, err)
	require.Equal(t, countTrackBeforeSend, expNextCount)
	expNextCount++

	// make sure that module to module doesn't bypass hook
	err = app.BankKeeper.SendCoinsFromModuleToModule(ctx, stakingtypes.BondedPoolName, stakingtypes.NotBondedPoolName, validSendAmount)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToModule(ctx, stakingtypes.BondedPoolName, stakingtypes.NotBondedPoolName, invalidBlockSendAmount)
	// there should be no error since module to module does not call block before send hooks
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToModule(ctx, stakingtypes.BondedPoolName, stakingtypes.NotBondedPoolName, triggerTrackSendAmount)
	require.NoError(t, err)
	require.Equal(t, countTrackBeforeSend, expNextCount)
	expNextCount++

	// make sure that module to many accounts doesn't bypass hook
	err = app.BankKeeper.SendCoinsFromModuleToManyAccounts(ctx, stakingtypes.BondedPoolName, []sdk.AccAddress{addr1, addr2}, []sdk.Coins{validSendAmount, validSendAmount})
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToManyAccounts(ctx, stakingtypes.BondedPoolName, []sdk.AccAddress{addr1, addr2}, []sdk.Coins{validSendAmount, invalidBlockSendAmount})
	require.Error(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToManyAccounts(ctx, stakingtypes.BondedPoolName, []sdk.AccAddress{addr1, addr2}, []sdk.Coins{validSendAmount, triggerTrackSendAmount})
	require.NoError(t, err)
	require.Equal(t, countTrackBeforeSend, expNextCount)
	expNextCount++

	// make sure that DelegateCoins doesn't bypass the hook
	err = app.BankKeeper.DelegateCoins(ctx, addr1, app.AccountKeeper.GetModuleAddress(stakingtypes.BondedPoolName), validSendAmount)
	require.NoError(t, err)
	err = app.BankKeeper.DelegateCoins(ctx, addr1, app.AccountKeeper.GetModuleAddress(stakingtypes.BondedPoolName), invalidBlockSendAmount)
	require.Error(t, err)
	err = app.BankKeeper.DelegateCoins(ctx, addr1, app.AccountKeeper.GetModuleAddress(stakingtypes.BondedPoolName), triggerTrackSendAmount)
	require.NoError(t, err)
	require.Equal(t, countTrackBeforeSend, expNextCount)
	expNextCount++

	// make sure that UndelegateCoins doesn't bypass the hook
	err = app.BankKeeper.UndelegateCoins(ctx, app.AccountKeeper.GetModuleAddress(stakingtypes.BondedPoolName), addr1, validSendAmount)
	require.NoError(t, err)
	err = app.BankKeeper.UndelegateCoins(ctx, app.AccountKeeper.GetModuleAddress(stakingtypes.BondedPoolName), addr1, invalidBlockSendAmount)
	require.Error(t, err)

	err = app.BankKeeper.UndelegateCoins(ctx, app.AccountKeeper.GetModuleAddress(stakingtypes.BondedPoolName), addr1, triggerTrackSendAmount)
	require.NoError(t, err)
	require.Equal(t, countTrackBeforeSend, expNextCount)
	expNextCount++
}
