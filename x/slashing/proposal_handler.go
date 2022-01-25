package slashing

import (
	"github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	"github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// NewSlashProposalsHandler  creates a new governance Handler for a slashing proposals
func NewSlashProposalsHandler(slashingKeeper keeper.Keeper, stakingKeeper stakingkeeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.SlashValidatorProposal:
			return handleSlashValidatorProposal(ctx, c, slashingKeeper, stakingKeeper)
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized param proposal content type: %T", c)
		}
	}
}

func handleSlashValidatorProposal(
	ctx sdk.Context, p *types.SlashValidatorProposal,
	slashingKeeper keeper.Keeper, stakingKeeper stakingkeeper.Keeper,
) error {
	valAddr, valAddrErr := sdk.ValAddressFromBech32(p.ValidatorAddress)
	if valAddrErr != nil {
		return valAddrErr
	}
	validator, found := stakingKeeper.GetValidator(ctx, valAddr)
	if !found {
		return stakingtypes.ErrNoValidatorFound
	}

	if validator.IsUnbonded() {
		return types.ErrSlashUnbondedValidator
	}

	valConsAddr, valConsdAddrErr := validator.GetConsAddr()
	if valConsdAddrErr != nil {
		return valConsdAddrErr
	}

	valPower := stakingKeeper.GetLastValidatorPower(ctx, valAddr)
	slashingKeeper.Slash(ctx, valConsAddr, p.SlashFactor, valPower, ctx.BlockHeight())
	return nil
}
