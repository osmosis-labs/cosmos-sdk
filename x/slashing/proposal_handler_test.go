package slashing_test

import (
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/staking/teststaking"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"testing"
)

func TestSlashValidatorProposal(t *testing.T) {
	// initial setup
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	proposalHandler := slashing.NewSlashProposalsHandler(app.SlashingKeeper, app.StakingKeeper)

	addrDels := simapp.AddTestAddrsIncremental(app, ctx, 2, app.StakingKeeper.TokensFromConsensusPower(ctx, 200))
	valAddrs := simapp.ConvertAddrsToValAddrs(addrDels)
	pks := simapp.CreateTestPubKeys(2)

	params := app.StakingKeeper.GetParams(ctx)
	params.MaxValidators = uint32(1)
	app.StakingKeeper.SetParams(ctx, params)
	tstaking := teststaking.NewHelper(t, ctx, app.StakingKeeper)

	tstaking.CreateValidatorWithValPower(valAddrs[0], pks[0], 100, true)
	tstaking.CreateValidatorWithValPower(valAddrs[1], pks[1], 1, true)

	staking.EndBlocker(ctx, app.StakingKeeper)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	// slash unbonded validator
	err := proposalHandler(ctx, &types.SlashValidatorProposal{
		Title:            "test proposal",
		Description:      "test proposal descr",
		ValidatorAddress: valAddrs[1].String(),
		SlashFactor:      sdk.MustNewDecFromStr("0.1"),
	})
	require.Error(t, err)

	// slash active validator
	activeValidator, _ := app.StakingKeeper.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(pks[0]))
	err = proposalHandler(ctx, &types.SlashValidatorProposal{
		Title:            "test proposal",
		Description:      "test proposal descr",
		ValidatorAddress: valAddrs[0].String(),
		SlashFactor:      sdk.MustNewDecFromStr("0.1"),
	})
	require.NoError(t, err)

	updates, err := app.StakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	require.Equal(t, 1, len(updates))

	slashedValidator, _ := app.StakingKeeper.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(pks[0]))
	if slashedValidator.Tokens.GTE(activeValidator.Tokens) {
		t.Fail()
	}
}
