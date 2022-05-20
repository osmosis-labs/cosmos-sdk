package types

import (
	"fmt"

	gov "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalTypeSoftwareUpgrade       string = "SoftwareUpgrade"
	ProposalTypeCancelSoftwareUpgrade string = "CancelSoftwareUpgrade"
)

func NewSoftwareUpgradeProposal(title, description string, isExpedited bool, plan Plan) gov.Content {
	return &SoftwareUpgradeProposal{title, description, isExpedited, plan}
}

// Implements Proposal Interface
var _ gov.Content = &SoftwareUpgradeProposal{}

func init() {
	gov.RegisterProposalType(ProposalTypeSoftwareUpgrade)
	gov.RegisterProposalTypeCodec(&SoftwareUpgradeProposal{}, "cosmos-sdk/SoftwareUpgradeProposal")
	gov.RegisterProposalType(ProposalTypeCancelSoftwareUpgrade)
	gov.RegisterProposalTypeCodec(&CancelSoftwareUpgradeProposal{}, "cosmos-sdk/CancelSoftwareUpgradeProposal")
}

func (sup *SoftwareUpgradeProposal) GetTitle() string       { return sup.Title }
func (sup *SoftwareUpgradeProposal) GetDescription() string { return sup.Description }
func (sup *SoftwareUpgradeProposal) GetIsExpedited() bool   { return sup.IsExpedited }
func (sup *SoftwareUpgradeProposal) SetIsExpedited(isExpedited bool) {
	sup.IsExpedited = isExpedited
}
func (sup *SoftwareUpgradeProposal) ProposalRoute() string { return RouterKey }
func (sup *SoftwareUpgradeProposal) ProposalType() string  { return ProposalTypeSoftwareUpgrade }
func (sup *SoftwareUpgradeProposal) ValidateBasic() error {
	if err := sup.Plan.ValidateBasic(); err != nil {
		return err
	}
	return gov.ValidateAbstract(sup)
}

func (sup SoftwareUpgradeProposal) String() string {
	return fmt.Sprintf(`Software Upgrade Proposal:
  Title:       %s
  Description: %s
`, sup.Title, sup.Description)
}

func NewCancelSoftwareUpgradeProposal(title, description string, isExpedited bool) gov.Content {
	return &CancelSoftwareUpgradeProposal{title, description, isExpedited}
}

// Implements Proposal Interface
var _ gov.Content = &CancelSoftwareUpgradeProposal{}

func (csup *CancelSoftwareUpgradeProposal) GetTitle() string       { return csup.Title }
func (csup *CancelSoftwareUpgradeProposal) GetDescription() string { return csup.Description }
func (csup *CancelSoftwareUpgradeProposal) GetIsExpedited() bool   { return csup.IsExpedited }
func (csup *CancelSoftwareUpgradeProposal) SetIsExpedited(isExpedited bool) {
	csup.IsExpedited = isExpedited
}
func (csup *CancelSoftwareUpgradeProposal) ProposalRoute() string { return RouterKey }
func (csup *CancelSoftwareUpgradeProposal) ProposalType() string {
	return ProposalTypeCancelSoftwareUpgrade
}
func (csup *CancelSoftwareUpgradeProposal) ValidateBasic() error {
	return gov.ValidateAbstract(csup)
}

func (csup CancelSoftwareUpgradeProposal) String() string {
	return fmt.Sprintf(`Cancel Software Upgrade Proposal:
  Title:       %s
  Description: %s
`, csup.Title, csup.Description)
}
