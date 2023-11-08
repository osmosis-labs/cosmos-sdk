package keeper_test

import (
	"context"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/comet"
	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"cosmossdk.io/simapp"
	"cosmossdk.io/x/evidence"

	storetypes "cosmossdk.io/store/types"
	"cosmossdk.io/x/evidence/exported"
	evidencekeeper "cosmossdk.io/x/evidence/keeper"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/slashing/testutil"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtestutil "github.com/cosmos/cosmos-sdk/x/staking/testutil"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"

	"github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil/integration"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	vesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktestutil "github.com/cosmos/cosmos-sdk/x/bank/testutil"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	tmtime "github.com/cometbft/cometbft/types/time"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"

	evidencetypes "cosmossdk.io/x/evidence/types"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	stakeDenom = "stake"
	feeDenom   = "fee"

	pubkeys = []cryptotypes.PubKey{
		newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB50"),
		newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB51"),
		newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB52"),
	}

	valAddresses = []sdk.ValAddress{
		sdk.ValAddress(pubkeys[0].Address()),
		sdk.ValAddress(pubkeys[1].Address()),
		sdk.ValAddress(pubkeys[2].Address()),
	}

	// The default power validators are initialized to have within tests
	initAmt   = sdk.TokensFromConsensusPower(200, sdk.DefaultPowerReduction)
	initCoins = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, initAmt))
)

type fixture struct {
	app *integration.App

	ctx sdk.Context
	cdc codec.Codec

	bankKeeper     bankkeeper.Keeper
	accountKeeper  *keeper.AccountKeeper
	evidenceKeeper *evidencekeeper.Keeper
	slashingKeeper slashingkeeper.Keeper
	stakingKeeper  *stakingkeeper.Keeper
}

func initFixture(t testing.TB) *fixture {
	keys := storetypes.NewKVStoreKeys(
		authtypes.StoreKey, banktypes.StoreKey, paramtypes.StoreKey, consensusparamtypes.StoreKey, evidencetypes.StoreKey, stakingtypes.StoreKey, slashingtypes.StoreKey,
	)
	cdc := moduletestutil.MakeTestEncodingConfig(auth.AppModuleBasic{}, evidence.AppModuleBasic{})

	logger := log.NewTestLogger(t)
	cms := integration.CreateMultiStore(keys, logger)

	newCtx := sdk.NewContext(cms, cmtproto.Header{}, true, logger)

	authority := authtypes.NewModuleAddress("gov")

	maccPerms := map[string][]string{
		minttypes.ModuleName:           {authtypes.Minter},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		distrtypes.ModuleName:          {authtypes.Minter, authtypes.Burner},
	}

	accountKeeper := authkeeper.NewAccountKeeper(
		cdc.Codec,
		runtime.NewKVStoreService(keys[authtypes.StoreKey]),
		authtypes.ProtoBaseAccount,
		maccPerms,
		addresscodec.NewBech32Codec(sdk.Bech32MainPrefix),
		sdk.Bech32MainPrefix,
		authority.String(),
	)

	vestingtypes.RegisterInterfaces(cdc.InterfaceRegistry)

	blockedAddresses := map[string]bool{
		accountKeeper.GetAuthority(): false,
	}
	bankKeeper := bankkeeper.NewBaseKeeper(
		cdc.Codec,
		runtime.NewKVStoreService(keys[banktypes.StoreKey]),
		accountKeeper,
		blockedAddresses,
		authority.String(),
		log.NewNopLogger(),
	)

	stakingKeeper := stakingkeeper.NewKeeper(cdc.Codec, runtime.NewKVStoreService(keys[stakingtypes.StoreKey]), accountKeeper, bankKeeper, authority.String(), addresscodec.NewBech32Codec(sdk.Bech32PrefixValAddr), addresscodec.NewBech32Codec(sdk.Bech32PrefixConsAddr))

	slashingKeeper := slashingkeeper.NewKeeper(cdc.Codec, codec.NewLegacyAmino(), runtime.NewKVStoreService(keys[slashingtypes.StoreKey]), stakingKeeper, authority.String())

	evidenceKeeper := evidencekeeper.NewKeeper(cdc.Codec, runtime.NewKVStoreService(keys[evidencetypes.StoreKey]), stakingKeeper, slashingKeeper, addresscodec.NewBech32Codec("cosmos"), runtime.ProvideCometInfoService())
	router := evidencetypes.NewRouter()
	router = router.AddRoute(evidencetypes.RouteEquivocation, testEquivocationHandler(evidenceKeeper))
	evidenceKeeper.SetRouter(router)

	authModule := auth.NewAppModule(cdc.Codec, accountKeeper, authsims.RandomGenesisAccounts, nil)
	bankModule := bank.NewAppModule(cdc.Codec, bankKeeper, accountKeeper, nil)
	stakingModule := staking.NewAppModule(cdc.Codec, stakingKeeper, accountKeeper, bankKeeper, nil)
	slashingModule := slashing.NewAppModule(cdc.Codec, slashingKeeper, accountKeeper, bankKeeper, stakingKeeper, nil, cdc.Codec.InterfaceRegistry())
	evidenceModule := evidence.NewAppModule(*evidenceKeeper)

	integrationApp := integration.NewIntegrationApp(newCtx, logger, keys, cdc.Codec, map[string]appmodule.AppModule{
		authtypes.ModuleName:     authModule,
		banktypes.ModuleName:     bankModule,
		stakingtypes.ModuleName:  stakingModule,
		slashingtypes.ModuleName: slashingModule,
		evidencetypes.ModuleName: evidenceModule,
	})

	sdkCtx := sdk.UnwrapSDKContext(integrationApp.Context())

	// Register MsgServer and QueryServer
	types.RegisterMsgServer(integrationApp.MsgServiceRouter(), keeper.NewMsgServerImpl(accountKeeper))
	types.RegisterQueryServer(integrationApp.QueryHelper(), keeper.NewQueryServer(accountKeeper))

	assert.NilError(t, slashingKeeper.SetParams(sdkCtx, testutil.TestParams()))

	// set default staking params
	assert.NilError(t, stakingKeeper.SetParams(sdkCtx, stakingtypes.DefaultParams()))

	return &fixture{
		app:            integrationApp,
		ctx:            sdkCtx,
		cdc:            cdc.Codec,
		bankKeeper:     bankKeeper,
		evidenceKeeper: evidenceKeeper,
		accountKeeper:  &accountKeeper,
		slashingKeeper: slashingKeeper,
		stakingKeeper:  stakingKeeper,
	}
}

func testEquivocationHandler(_ interface{}) evidencetypes.Handler {
	return func(ctx context.Context, e exported.Evidence) error {
		if err := e.ValidateBasic(); err != nil {
			return err
		}

		ee, ok := e.(*evidencetypes.Equivocation)
		if !ok {
			return fmt.Errorf("unexpected evidence type: %T", e)
		}
		if ee.Height%2 == 0 {
			return fmt.Errorf("unexpected even evidence height: %d", ee.Height)
		}

		return nil
	}
}

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
			t.Parallel()
			f := initFixture(t)

			ctx := f.ctx.WithIsCheckTx(false).WithBlockHeight(1)
			populateValidators(t, f)

			power := int64(100)
			stakingParams, err := f.stakingKeeper.GetParams(ctx)
			assert.NilError(t, err)
			operatorAddr, valpubkey := valAddresses[0], pubkeys[0]
			tstaking := stakingtestutil.NewHelper(t, ctx, f.stakingKeeper)

			selfDelegation := tstaking.CreateValidatorWithValPower(operatorAddr, valpubkey, power, true)

			// execute end-blocker and verify validator attributes
			_, err = f.stakingKeeper.EndBlocker(ctx)
			assert.NilError(t, err)
			assert.DeepEqual(t,
				f.bankKeeper.GetAllBalances(ctx, sdk.AccAddress(operatorAddr)).String(),
				sdk.NewCoins(sdk.NewCoin(stakingParams.BondDenom, initAmt.Sub(selfDelegation))).String(),
			)
			valI, err := f.stakingKeeper.Validator(ctx, operatorAddr)
			assert.NilError(t, err)
			assert.DeepEqual(t, selfDelegation, valI.GetBondedTokens())

			assert.NilError(t, f.slashingKeeper.AddPubkey(ctx, valpubkey))

			info := slashingtypes.NewValidatorSigningInfo(sdk.ConsAddress(valpubkey.Address()), ctx.BlockHeight(), int64(0), time.Unix(0, 0), false, int64(0))
			f.slashingKeeper.SetValidatorSigningInfo(ctx, sdk.ConsAddress(valpubkey.Address()), info)

			// handle a signature to set signing info
			f.slashingKeeper.HandleValidatorSignature(ctx, valpubkey.Address(), selfDelegation.Int64(), comet.BlockIDFlagCommit)

			val, found := f.stakingKeeper.GetValidator(ctx, operatorAddr)
			assert.Assert(t, found)

			bondDenom, err := f.stakingKeeper.BondDenom(ctx)
			require.NoError(t, err)
			require.Equal(t, "stake", bondDenom)

			bacc, origCoins := initBaseAccount()
			_, _, funder := testdata.KeyTestPubAddr()
			va := vesting.NewClawbackVestingAccount(bacc, funder, origCoins, now.Unix(), lockupPeriods, vestingPeriods)
			acc := f.accountKeeper.NewAccount(ctx, va)
			f.accountKeeper.SetAccount(ctx, acc)
			addr := acc.GetAddress()

			// fund the vesting account
			err = banktestutil.FundAccount(ctx, f.bankKeeper, addr, c(fee(1000), stake(100)))
			require.NoError(t, err)
			require.Equal(t, int64(1000), f.bankKeeper.GetBalance(ctx, addr, feeDenom).Amount.Int64())
			require.Equal(t, int64(100), f.bankKeeper.GetBalance(ctx, addr, stakeDenom).Amount.Int64())

			// try delegating, clawback vesting account not allowed to delegate
			_, err = f.stakingKeeper.Delegate(ctx, addr, math.NewInt(65), stakingtypes.Unbonded, val, true)
			require.Error(t, err)

			// undelegation should emit an error(delegator does not contain delegation)
			_, _, err = f.stakingKeeper.Undelegate(ctx, addr, operatorAddr, math.LegacyNewDec(5))
			require.Error(t, err)

			ctx = ctx.WithBlockTime(tc.ctxTime)
			va = f.accountKeeper.GetAccount(ctx, addr).(*vesting.ClawbackVestingAccount)
			clawbackAction := vesting.NewClawbackAction(funder, funder, f.accountKeeper, f.bankKeeper)
			err = va.Clawback(ctx, clawbackAction)
			require.NoError(t, err)
			f.accountKeeper.SetAccount(ctx, va)

			vestingAccBalance := f.bankKeeper.GetAllBalances(ctx, addr)
			require.Equal(t, tc.vestingAccBalance, vestingAccBalance, "vesting account balance test")

			funderBalance := f.bankKeeper.GetAllBalances(ctx, funder)
			require.Equal(t, tc.funderBalance, funderBalance, "funder account balance test")

			spendableCoins := f.bankKeeper.SpendableCoins(ctx, addr)
			require.Equal(t, tc.spendableCoins, spendableCoins, "vesting account spendable test")
		})
	}
}

func populateValidators(t assert.TestingT, f *fixture) {
	// add accounts and set total supply
	totalSupplyAmt := initAmt.MulRaw(int64(len(valAddresses)))
	totalSupply := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, totalSupplyAmt))
	assert.NilError(t, f.bankKeeper.MintCoins(f.ctx, minttypes.ModuleName, totalSupply))

	for _, addr := range valAddresses {
		assert.NilError(t, f.bankKeeper.SendCoinsFromModuleToAccount(f.ctx, minttypes.ModuleName, (sdk.AccAddress)(addr), initCoins))
	}
}

func newPubKey(pk string) (res cryptotypes.PubKey) {
	pkBytes, err := hex.DecodeString(pk)
	if err != nil {
		panic(err)
	}

	pubkey := &ed25519.PubKey{Key: pkBytes}

	return pubkey
}

// createValidator creates a validator in the given SimApp.
func createValidator(t *testing.T, ctx sdk.Context, app *simapp.SimApp, powers int64) (sdk.ValAddress, stakingtypes.Validator) {
	valTokens := sdk.TokensFromConsensusPower(powers, sdk.DefaultPowerReduction)
	addrs := simapp.AddTestAddrsIncremental(app, ctx, 1, valTokens)
	valAddrs := simtestutil.ConvertAddrsToValAddrs(addrs)
	pks := simtestutil.CreateTestPubKeys(1)
	cdc := app.AppCodec() //simapp.MakeTestEncodingConfig().Marshaler
	key := storetypes.NewKVStoreKey(stakingtypes.StoreKey)
	storeService := runtime.NewKVStoreService(key)

	app.StakingKeeper = stakingkeeper.NewKeeper(
		cdc,
		storeService,
		app.AccountKeeper,
		app.BankKeeper,
		types.NewModuleAddress(types.ModuleName).String(),
		address.NewBech32Codec("cosmosvaloper"),
		address.NewBech32Codec("cosmosvalcons"),
	)

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
