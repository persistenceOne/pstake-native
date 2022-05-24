package types

import (
	"fmt"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"strings"
)

const (
	ProposalTypeChangeCosmosValidatorWeights = "ChangeCosmosValidatorWeights"
	ProposalTypeChangeOracleValidatorWeights = "ChangeOracleValidatorWeights"
	ProposalTypeChangeMultisig               = "ChangeMultisig" // unrelated to Changing cosmos validator weights.
	ProposalTypeEnableModule                 = "EnableModule"
)

var _, _, _, _ govtypes.Content = &ChangeMultisigProposal{}, &EnableModuleProposal{}, &ChangeCosmosValidatorWeightsProposal{}, &ChangeOracleValidatorWeightsProposal{}

func init() {
	govtypes.RegisterProposalType(ProposalTypeChangeMultisig)
	govtypes.RegisterProposalType(ProposalTypeEnableModule)
	govtypes.RegisterProposalType(ProposalTypeChangeCosmosValidatorWeights)
	govtypes.RegisterProposalType(ProposalTypeChangeOracleValidatorWeights)
	govtypes.RegisterProposalTypeCodec(&ChangeMultisigProposal{}, "persistenceCore/ChangeMultisigProposal")
	govtypes.RegisterProposalTypeCodec(&EnableModuleProposal{}, "persistenceCore/EnableModuleProposal")
	govtypes.RegisterProposalTypeCodec(&ChangeCosmosValidatorWeightsProposal{}, "persistenceCore/ChangeCosmosValidatorWeightsProposal")
	govtypes.RegisterProposalTypeCodec(&ChangeOracleValidatorWeightsProposal{}, "persistenceCore/ChangeOracleValidatorWeightsProposal")
}

func NewChangeMultisigProposal(title, description string, threshold uint64, orchestratorAddresses []string, accountNumber uint64) *ChangeMultisigProposal {
	return &ChangeMultisigProposal{
		Title:                title,
		Description:          description,
		Threshold:            threshold,
		OrcastratorAddresses: orchestratorAddresses,
		AccountNumber:        accountNumber,
	}
}

func (m *ChangeMultisigProposal) GetTitle() string {
	return m.Title
}

func (m *ChangeMultisigProposal) GetDescription() string {
	return m.Description
}

func (m *ChangeMultisigProposal) ProposalRoute() string {
	return RouterKey
}

func (m *ChangeMultisigProposal) ProposalType() string {
	return ProposalTypeChangeMultisig
}

func (m *ChangeMultisigProposal) ValidateBasic() error {
	//TODO add validations
	return nil
}

func (m *ChangeMultisigProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Update Pool Incentives Proposal:
  Title:                 %s
  Description:           %s
  Threshold:             %s
  OrcastratorAddresses:  %s
`, m.Title, m.Description, m.Threshold, m.OrcastratorAddresses))
	return b.String()
}

func NewEnableModuleProposal(title, description string, threshold uint64, accountNumber uint64) *EnableModuleProposal {
	return &EnableModuleProposal{
		Title:         title,
		Description:   description,
		Threshold:     threshold,
		AccountNumber: accountNumber,
	}
}

func (m *EnableModuleProposal) GetTitle() string {
	return m.Title
}

func (m *EnableModuleProposal) GetDescription() string {
	return m.Description
}

func (m *EnableModuleProposal) ProposalRoute() string {
	return RouterKey
}

func (m *EnableModuleProposal) ProposalType() string {
	return ProposalTypeEnableModule
}

func (m *EnableModuleProposal) ValidateBasic() error {
	//TODO add validations
	return nil
}

func (m *EnableModuleProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Update Pool Incentives Proposal:
  Title:                 %s
  Description:           %s
`, m.Title, m.Description))
	return b.String()
}

func NewChangeCosmosValidatorWeightsProposal(title, description string, weightedAddresses []WeightedAddressAmount) *ChangeCosmosValidatorWeightsProposal {
	return &ChangeCosmosValidatorWeightsProposal{
		Title:             title,
		Description:       description,
		WeightedAddresses: weightedAddresses,
	}
}

func (m *ChangeCosmosValidatorWeightsProposal) GetTitle() string {
	return m.Title
}

func (m *ChangeCosmosValidatorWeightsProposal) GetDescription() string {
	return m.Description
}

func (m *ChangeCosmosValidatorWeightsProposal) ProposalRoute() string {
	return RouterKey
}

func (m *ChangeCosmosValidatorWeightsProposal) ProposalType() string {
	return ProposalTypeChangeCosmosValidatorWeights
}

func (m *ChangeCosmosValidatorWeightsProposal) ValidateBasic() error {
	//TODO add validations
	return nil
}

func (m *ChangeCosmosValidatorWeightsProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Update Pool Incentives Proposal:
  Title:                 %s
  Description:           %s
`, m.Title, m.Description))
	return b.String()
}

func NewChangeOracleValidatorWeightsProposal(title, description string, weightedAddresses []WeightedAddress) *ChangeOracleValidatorWeightsProposal {
	return &ChangeOracleValidatorWeightsProposal{
		Title:             title,
		Description:       description,
		WeightedAddresses: weightedAddresses,
	}
}

func (m *ChangeOracleValidatorWeightsProposal) GetTitle() string {
	return m.Title
}

func (m *ChangeOracleValidatorWeightsProposal) GetDescription() string {
	return m.Description
}

func (m *ChangeOracleValidatorWeightsProposal) ProposalRoute() string {
	return RouterKey
}

func (m *ChangeOracleValidatorWeightsProposal) ProposalType() string {
	return ProposalTypeChangeOracleValidatorWeights
}

func (m *ChangeOracleValidatorWeightsProposal) ValidateBasic() error {
	//TODO add validations
	return nil
}

func (m *ChangeOracleValidatorWeightsProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Update Pool Incentives Proposal:
  Title:                 %s
  Description:           %s
`, m.Title, m.Description))
	return b.String()
}
