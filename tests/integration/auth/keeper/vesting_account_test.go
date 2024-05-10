package keeper_test

import (
	"testing"
	"time"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtime "github.com/cometbft/cometbft/types/time"
	"github.com/stretchr/testify/require"

	"cosmossdk.io/math"
	"cosmossdk.io/simapp"

	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	vesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	banktestutil "github.com/cosmos/cosmos-sdk/x/bank/testutil"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

var (
	stakeDenom = "stake"
	feeDenom   = "fee"
)

func TestAddGrantClawbackVestingAcc(t *testing.T) {
	c := sdk.NewCoins
	fee := func(amt int64) sdk.Coin { return sdk.NewInt64Coin(feeDenom, amt) }
	now := tmtime.Now()

	// set up simapp
	app := simapp.Setup(t, false)
	ctx := app.BaseApp.NewContextLegacy(false, tmproto.Header{}).WithBlockTime((now))
	bondDenom, err := app.StakingKeeper.BondDenom(ctx)
	require.NoError(t, err)
	require.Equal(t, "stake", bondDenom)

	// create an account with an initial grant
	_, _, funder := testdata.KeyTestPubAddr()
	lockupPeriods := vesting.Periods{{Length: int64(12 * 3600), Amount: c(fee(1000))}} // noon
	vestingPeriods := vesting.Periods{
		{Length: int64(8 * 3600), Amount: c(fee(200))}, // 8am
		{Length: int64(1 * 3600), Amount: c(fee(200))}, // 9am
		{Length: int64(6 * 3600), Amount: c(fee(200))}, // 3pm
		{Length: int64(2 * 3600), Amount: c(fee(200))}, // 5pm
		{Length: int64(1 * 3600), Amount: c(fee(200))}, // 6pm
	}
	bacc, origCoins := initBaseAccount()
	va := vesting.NewClawbackVestingAccount(bacc, funder, origCoins, now.Unix(), lockupPeriods, vestingPeriods)
	addr := va.GetAddress()

	ctx = ctx.WithBlockTime(now.Add(11 * time.Hour))
	require.Equal(t, int64(1000), va.GetVestingCoins(ctx.BlockTime()).AmountOf(feeDenom).Int64())

	// Add a new grant(1000fee, 100stake) while all slashing is covered by unvested tokens
	grantAction := vesting.NewClawbackGrantAction(funder.String(), ctx.BlockTime().Unix(),
		lockupPeriods, vestingPeriods, origCoins)
	err = va.AddGrant(ctx, grantAction)
	require.NoError(t, err)

	// locked coin is expected to be 2000feetoken(1000fee + 1000fee)
	require.Equal(t, int64(2000), va.GetVestingCoins(ctx.BlockTime()).AmountOf(feeDenom).Int64())
	require.Equal(t, int64(0), va.DelegatedVesting.AmountOf(feeDenom).Int64())
	require.Equal(t, int64(0), va.DelegatedFree.AmountOf(feeDenom).Int64())

	ctx = ctx.WithBlockTime(now.Add(13 * time.Hour))
	require.Equal(t, int64(1600), va.GetVestingCoins(ctx.BlockTime()).AmountOf(feeDenom).Int64())

	ctx = ctx.WithBlockTime(now.Add(17 * time.Hour))
	require.Equal(t, int64(1200), va.GetVestingCoins(ctx.BlockTime()).AmountOf(feeDenom).Int64())

	ctx = ctx.WithBlockTime(now.Add(20 * time.Hour))
	require.Equal(t, int64(1000), va.GetVestingCoins(ctx.BlockTime()).AmountOf(feeDenom).Int64())

	ctx = ctx.WithBlockTime(now.Add(22 * time.Hour))
	require.Equal(t, int64(1000), va.GetVestingCoins(ctx.BlockTime()).AmountOf(feeDenom).Int64())

	// fund the vesting account with new grant (old has vested and transferred out)
	err = banktestutil.FundAccount(ctx, app.BankKeeper, addr, origCoins)
	require.NoError(t, err)
	require.Equal(t, int64(100), app.BankKeeper.GetBalance(ctx, addr, stakeDenom).Amount.Int64())

	feeAmt := app.BankKeeper.GetBalance(ctx, addr, feeDenom).Amount
	require.Equal(t, int64(1000), feeAmt.Int64())
}

func TestClawback(t *testing.T) {
	c := sdk.NewCoins
	fee := func(x int64) sdk.Coin { return sdk.NewInt64Coin(feeDenom, x) }
	stake := func(x int64) sdk.Coin { return sdk.NewInt64Coin(stakeDenom, x) }
	now := tmtime.Now()

	lockupPeriods := vesting.Periods{
		{Length: int64(12 * 3600), Amount: c(fee(1000), stake(100))}, // noon
	}
	vestingPeriods := vesting.Periods{
		{Length: int64(8 * 3600), Amount: c(fee(200))},            // 8am
		{Length: int64(1 * 3600), Amount: c(fee(200), stake(50))}, // 9am
		{Length: int64(6 * 3600), Amount: c(fee(200), stake(50))}, // 3pm
		{Length: int64(2 * 3600), Amount: c(fee(200))},            // 5pm
		{Length: int64(1 * 3600), Amount: c(fee(200))},            // 6pm
	}
	// each test creates a new clawback vesting account, with the lockup and vesting periods defined above.
	// the clawback is executed at the test case's provided time, and expects that post clawback,
	// the address has a total of `vestingAccBalance` coins, but only `spendableCoins` are spendable.
	// It expects the clawback acct funder to have `funderBalance` (aka that amt clawed back)
	testCases := []struct {
		name              string
		ctxTime           time.Time
		vestingAccBalance sdk.Coins
		spendableCoins    sdk.Coins
		funderBalance     sdk.Coins
	}{
		{
			"clawback before all vesting periods, before cliff ended",
			now.Add(7 * time.Hour),
			// vesting account should not have funds after clawback
			sdk.NewCoins(),
			sdk.Coins{},
			// all funds should be returned to funder account
			sdk.NewCoins(sdk.NewCoin(feeDenom, math.NewInt(1000)), sdk.NewCoin(stakeDenom, math.NewInt(100))),
		},
		{
			"clawback after two vesting periods, before cliff ended",
			now.Add(10 * time.Hour),
			sdk.NewCoins(fee(400), stake(50)),
			sdk.Coins{},
			// everything but first two vesting periods of fund should be returned to sender
			sdk.NewCoins(sdk.NewCoin(feeDenom, math.NewInt(600)), sdk.NewCoin(stakeDenom, math.NewInt(50))),
		},
		{
			"clawback right after cliff has finsihed",
			now.Add(13 * time.Hour),
			sdk.NewCoins(sdk.NewCoin(feeDenom, math.NewInt(400)), sdk.NewCoin(stakeDenom, math.NewInt(50))),
			sdk.NewCoins(sdk.NewCoin(feeDenom, math.NewInt(400)), sdk.NewCoin(stakeDenom, math.NewInt(50))),
			sdk.NewCoins(sdk.NewCoin(feeDenom, math.NewInt(600)), sdk.NewCoin(stakeDenom, math.NewInt(50))),
		},
		{
			"clawback after cliff has finished, 3 vesting periods have finished",
			now.Add(16 * time.Hour),
			sdk.NewCoins(sdk.NewCoin(feeDenom, math.NewInt(600)), sdk.NewCoin(stakeDenom, math.NewInt(100))),
			sdk.NewCoins(sdk.NewCoin(feeDenom, math.NewInt(600)), sdk.NewCoin(stakeDenom, math.NewInt(100))),
			sdk.NewCoins(sdk.NewCoin(feeDenom, math.NewInt(400))),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			// set up simapp and validators
			app := simapp.Setup(t, false)
			ctx := app.BaseApp.NewContextLegacy(false, tmproto.Header{}).WithBlockTime((now))
			valAddr, val := createValidator(t, ctx, app, 100)
			bondDenom, err := app.StakingKeeper.BondDenom(ctx)
			require.NoError(t, err)
			require.Equal(t, "stake", bondDenom)

			// Set up funder
			origCoins := sdk.Coins{sdk.NewInt64Coin(feeDenom, 1000), sdk.NewInt64Coin(stakeDenom, 100)}
			_, _, to := testdata.KeyTestPubAddr()
			bacc := types.NewBaseAccountWithAddress(to)
			addrDel := sdk.AccAddress([]byte("addr"))
			acc := app.AccountKeeper.NewAccountWithAddress(ctx, addrDel)
			app.AccountKeeper.SetAccount(ctx, acc)
			funder := acc.GetAddress()

			// create a clawback vesting account
			va := vesting.NewClawbackVestingAccount(bacc, funder, origCoins, now.Unix(), lockupPeriods, vestingPeriods)
			addr := va.GetAddress()
			app.AccountKeeper.NewAccount(ctx, va)
			app.AccountKeeper.SetAccount(ctx, va)

			// fund the vesting account
			err = banktestutil.FundAccount(ctx, app.BankKeeper, addr, c(fee(1000), stake(100)))
			require.NoError(t, err)
			require.Equal(t, int64(1000), app.BankKeeper.GetBalance(ctx, addr, feeDenom).Amount.Int64())
			require.Equal(t, int64(100), app.BankKeeper.GetBalance(ctx, addr, stakeDenom).Amount.Int64())

			// try delegating, clawback vesting account not allowed to delegate
			_, err = app.StakingKeeper.Delegate(ctx, addr, math.NewInt(65), stakingtypes.Unbonded, val, true)
			require.Error(t, err)

			// undelegation should emit an error(delegator does not contain delegation)
			_, _, err = app.StakingKeeper.Undelegate(ctx, addr, valAddr, math.LegacyNewDec(5))
			require.Error(t, err)

			ctx = ctx.WithBlockTime(tc.ctxTime)
			va = app.AccountKeeper.GetAccount(ctx, addr).(*vesting.ClawbackVestingAccount)
			clawbackAction := vesting.NewClawbackAction(funder, funder, app.AccountKeeper, app.BankKeeper)
			err = va.Clawback(ctx, clawbackAction)
			require.NoError(t, err)
			app.AccountKeeper.SetAccount(ctx, va)

			vestingAccBalance := app.BankKeeper.GetAllBalances(ctx, addr)
			require.Equal(t, tc.vestingAccBalance, vestingAccBalance, "vesting account balance test")

			funderBalance := app.BankKeeper.GetAllBalances(ctx, funder)
			require.Equal(t, tc.funderBalance, funderBalance, "funder account balance test")

			spendableCoins := app.BankKeeper.SpendableCoins(ctx, addr)
			require.Equal(t, tc.spendableCoins, spendableCoins, "vesting account spendable test")
		})
	}
}

// createValidator creates a validator in the given SimApp.
func createValidator(t *testing.T, ctx sdk.Context, app *simapp.SimApp, powers int64) (sdk.ValAddress, stakingtypes.Validator) {
	valTokens := sdk.TokensFromConsensusPower(powers, sdk.DefaultPowerReduction)
	addrs := simapp.AddTestAddrsIncremental(app, ctx, 1, valTokens)
	valAddrs := simtestutil.ConvertAddrsToValAddrs(addrs)
	pks := simtestutil.CreateTestPubKeys(1)
	cdc := app.AppCodec()

	authority := types.NewModuleAddress("gov")
	app.StakingKeeper = stakingkeeper.NewKeeper(cdc, runtime.NewKVStoreService(app.GetKey(stakingtypes.StoreKey)), app.AccountKeeper, app.BankKeeper, authority.String(), addresscodec.NewBech32Codec(sdk.Bech32PrefixValAddr), addresscodec.NewBech32Codec(sdk.Bech32PrefixConsAddr))

	val, err := stakingtypes.NewValidator(valAddrs[0].String(), pks[0], stakingtypes.Description{})
	require.NoError(t, err)

	app.StakingKeeper.SetValidator(ctx, val)
	require.NoError(t, app.StakingKeeper.SetValidatorByConsAddr(ctx, val))
	app.StakingKeeper.SetNewValidatorByPowerIndex(ctx, val)

	_, err = app.StakingKeeper.Delegate(ctx, addrs[0], valTokens, stakingtypes.Unbonded, val, true)
	require.NoError(t, err)

	_, err = app.StakingKeeper.EndBlocker(ctx)
	require.NoError(t, err)

	return valAddrs[0], val
}

func initBaseAccount() (*types.BaseAccount, sdk.Coins) {
	_, _, addr := testdata.KeyTestPubAddr()
	origCoins := sdk.Coins{sdk.NewInt64Coin(feeDenom, 1000), sdk.NewInt64Coin(stakeDenom, 100)}
	bacc := types.NewBaseAccountWithAddress(addr)

	return bacc, origCoins
}
