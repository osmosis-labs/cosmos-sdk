package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	// ProposalTypeSlashValidator defines the type for a SlashValidatorProposal
	ProposalTypeSlashValidator = "SlashValidator"
)

// String implements the Stringer interface.
func (svp SlashValidatorProposal) String() string {
	return fmt.Sprintf(
		`Slash Validator Proposal:
  Title:       %s
  Description: %s
  Validator Address: %s
  Slash Factor: %s
`, svp.Title, svp.Description, svp.ValidatorAddress, svp.SlashFactor.String())
}

// GetTitle returns the title of a parameter change proposal.
func (svp SlashValidatorProposal) GetTitle() string { return svp.Title }

// GetDescription returns the description of a parameter change proposal.
func (svp SlashValidatorProposal) GetDescription() string { return svp.Description }

// ProposalRoute returns the routing key of a parameter change proposal.
func (svp SlashValidatorProposal) ProposalRoute() string { return ProposalTypeSlashValidator }

// ProposalType returns the type of a parameter change proposal.
func (svp SlashValidatorProposal) ProposalType() string { return ProposalTypeSlashValidator }

// ValidateBasic validates the parameter change proposal
func (svp SlashValidatorProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(svp)
	if err != nil {
		return err
	}
	if svp.SlashFactor.LTE(sdk.ZeroDec()) || svp.SlashFactor.GT(sdk.OneDec()) {
		return ErrBadSlashingFactor
	}
	_, valAddrErr := sdk.ValAddressFromBech32(svp.ValidatorAddress)
	if valAddrErr != nil {
		return ErrBadValidatorAddr
	}
	return nil
}
