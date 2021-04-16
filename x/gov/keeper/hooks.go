package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/types"
)

// Implements GovHooks interface
var _ types.GovHooks = Keeper{}

// AfterProposalSubmission - call hook if registered
func (k Keeper) AfterProposalSubmission(ctx sdk.Context, proposalID uint64) {
	if k.hooks != nil {
		k.hooks.AfterProposalSubmission(ctx, proposalID)
	}
}

// AfterProposalDeposit - call hook if registered
func (k Keeper) AfterProposalDeposit(ctx sdk.Context, proposalID uint64, depositAmount sdk.Coins) {
	if k.hooks != nil {
		k.hooks.AfterProposalDeposit(ctx, proposalID, depositAmount)
	}
}

// AfterProposalVote - call hook if registered
func (k Keeper) AfterProposalVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress) {
	if k.hooks != nil {
		k.hooks.AfterProposalVote(ctx, proposalID, voterAddr)
	}
}

// AfterProposalInactive - call hook if registered
func (k Keeper) AfterProposalInactive(ctx sdk.Context, proposalID uint64) {
	if k.hooks != nil {
		k.hooks.AfterProposalInactive(ctx, proposalID)
	}
}

// AfterProposalActive - call hook if registered
func (k Keeper) AfterProposalActive(ctx sdk.Context, proposalID uint64) {
	if k.hooks != nil {
		k.hooks.AfterProposalActive(ctx, proposalID)
	}
}
