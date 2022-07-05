package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalTypeChangeCosmosValidatorWeights = "ChangeCosmosValidatorWeights"
	ProposalTypeChangeOracleValidatorWeights = "ChangeOracleValidatorWeights"
	ProposalTypeChangeMultisig               = "ChangeMultisig" // unrelated to Changing cosmos validator weights.
	ProposalTypeEnableModule                 = "EnableModule"
)

var _, _, _, _ govTypes.Content = &ChangeMultisigProposal{}, &EnableModuleProposal{}, &ChangeCosmosValidatorWeightsProposal{}, &ChangeOracleValidatorWeightsProposal{}

func init() {
	govTypes.RegisterProposalType(ProposalTypeChangeMultisig)
	govTypes.RegisterProposalType(ProposalTypeEnableModule)
	govTypes.RegisterProposalType(ProposalTypeChangeCosmosValidatorWeights)
	govTypes.RegisterProposalType(ProposalTypeChangeOracleValidatorWeights)
	govTypes.RegisterProposalTypeCodec(&ChangeMultisigProposal{}, "persistenceCore/ChangeMultisigProposal")
	govTypes.RegisterProposalTypeCodec(&EnableModuleProposal{}, "persistenceCore/EnableModuleProposal")
	govTypes.RegisterProposalTypeCodec(&ChangeCosmosValidatorWeightsProposal{}, "persistenceCore/ChangeCosmosValidatorWeightsProposal")
	govTypes.RegisterProposalTypeCodec(&ChangeOracleValidatorWeightsProposal{}, "persistenceCore/ChangeOracleValidatorWeightsProposal")
}

// NewChangeMultisigProposal creates a new multisig change proposal.
func NewChangeMultisigProposal(title, description string, threshold uint64, orchestratorAddresses []string, accountNumber uint64) *ChangeMultisigProposal {
	return &ChangeMultisigProposal{
		Title:                title,
		Description:          description,
		Threshold:            threshold,
		OrcastratorAddresses: orchestratorAddresses,
		AccountNumber:        accountNumber,
	}
}

// GetTitle returns the title of the multisig change proposal.
func (m *ChangeMultisigProposal) GetTitle() string {
	return m.Title
}

// GetDescription returns the description of multisig change proposal.
func (m *ChangeMultisigProposal) GetDescription() string {
	return m.Description
}

// ProposalRoute returns the proposal route of multisig change proposal.
func (m *ChangeMultisigProposal) ProposalRoute() string {
	return RouterKey
}

// ProposalType returns the proposal type of multisig change proposal.
func (m *ChangeMultisigProposal) ProposalType() string {
	return ProposalTypeChangeMultisig
}

// ValidateBasic runs basic stateless validity checks
func (m *ChangeMultisigProposal) ValidateBasic() error {
	err := govTypes.ValidateAbstract(m)
	if err != nil {
		return err
	}

	if m.Threshold > uint64(len(m.OrcastratorAddresses)) {
		return fmt.Errorf("threshold cannot be greated than the number of addresses")
	}

	return nil
}

// String returns the string of proposal details
func (m *ChangeMultisigProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Update Pool Incentives Proposal:
  Title:                 %s
  Description:           %s
  Threshold:             %d
  OrcastratorAddresses:  %s
`, m.Title, m.Description, m.Threshold, m.OrcastratorAddresses))
	return b.String()
}

// NewEnableModuleProposal returns a new module enable proposal
func NewEnableModuleProposal(title, description, custodialAddress, chainID string, threshold, accountNumber,
	sequenceNumber uint64, orchestratorAddresses []string) *EnableModuleProposal {
	return &EnableModuleProposal{
		Title:                 title,
		Description:           description,
		Threshold:             threshold,
		AccountNumber:         accountNumber,
		OrchestratorAddresses: orchestratorAddresses,
		SequenceNumber:        sequenceNumber,
		CustodialAddress:      custodialAddress,
		ChainID:               chainID,
	}
}

// GetTitle returns the title of module enable proposal
func (m *EnableModuleProposal) GetTitle() string {
	return m.Title
}

// GetDescription returns the description of module enable proposal
func (m *EnableModuleProposal) GetDescription() string {
	return m.Description
}

// ProposalRoute returns the proposal route for the module enable proposal
func (m *EnableModuleProposal) ProposalRoute() string {
	return RouterKey
}

// ProposalType returns the proposal type for the module enable proposal
func (m *EnableModuleProposal) ProposalType() string {
	return ProposalTypeEnableModule
}

// ValidateBasic runs basic stateless validity checks
func (m *EnableModuleProposal) ValidateBasic() error {
	err := govTypes.ValidateAbstract(m)
	if err != nil {
		return err
	}

	for _, acc := range m.OrchestratorAddresses {
		_, err = sdk.AccAddressFromBech32(acc)
		if err != nil {
			return err
		}
	}

	_, err = AccAddressFromBech32(m.CustodialAddress, Bech32PrefixAccAddr)
	if err != nil {
		return err
	}

	if m.ChainID == "" {
		return fmt.Errorf("chain ID can not be empty")
	}

	if m.Threshold > uint64(len(m.OrchestratorAddresses)) {
		return fmt.Errorf("threshold cannot be greater than the number of addresses")
	}

	return nil
}

// String returns the string of proposal details
func (m *EnableModuleProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Update Pool Incentives Proposal:
  Title:                 %s
  Description:           %s
`, m.Title, m.Description))
	return b.String()
}

// NewChangeCosmosValidatorWeightsProposal returns a new cosmos validator weights change proposal
func NewChangeCosmosValidatorWeightsProposal(title, description string, weightedAddresses []WeightedAddressAmount) *ChangeCosmosValidatorWeightsProposal {
	return &ChangeCosmosValidatorWeightsProposal{
		Title:             title,
		Description:       description,
		WeightedAddresses: weightedAddresses,
	}
}

// GetTitle returns the title of cosmos validator weights change proposal
func (m *ChangeCosmosValidatorWeightsProposal) GetTitle() string {
	return m.Title
}

// GetDescription returns the description of cosmos validator weights change proposal
func (m *ChangeCosmosValidatorWeightsProposal) GetDescription() string {
	return m.Description
}

// ProposalRoute returns the proposal route for the cosmos validator weights change proposal
func (m *ChangeCosmosValidatorWeightsProposal) ProposalRoute() string {
	return RouterKey
}

// ProposalType returns the proposal type for the cosmos validator weights change proposal
func (m *ChangeCosmosValidatorWeightsProposal) ProposalType() string {
	return ProposalTypeChangeCosmosValidatorWeights
}

// ValidateBasic runs basic stateless validity checks
func (m *ChangeCosmosValidatorWeightsProposal) ValidateBasic() error {
	err := govTypes.ValidateAbstract(m)
	if err != nil {
		return err
	}

	if len(m.WeightedAddresses) == 0 {
		return fmt.Errorf("address should be more than zero")
	}

	weightSum := sdk.NewDec(0)
	for i, w := range m.WeightedAddresses {
		if w.Address != "" {
			_, err := ValAddressFromBech32(w.Address, Bech32PrefixValAddr)
			if err != nil {
				return fmt.Errorf("invalid address at %dth", i)
			}
		}
		if !w.Weight.IsPositive() {
			return fmt.Errorf("non-positive weight at %dth", i)
		}
		if w.Weight.GT(sdk.NewDec(1)) {
			return fmt.Errorf("more than 1 weight at %dth", i)
		}
		weightSum = weightSum.Add(w.Weight)
	}

	if !weightSum.Equal(sdk.NewDec(1)) {
		return fmt.Errorf("invalid weight sum: %s", weightSum.String())
	}

	return nil
}

// String returns the string of proposal details
func (m *ChangeCosmosValidatorWeightsProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Update Pool Incentives Proposal:
  Title:                 %s
  Description:           %s
`, m.Title, m.Description))
	return b.String()
}

// NewChangeOracleValidatorWeightsProposal returns a new oracle validator weights change proposal
func NewChangeOracleValidatorWeightsProposal(title, description string, weightedAddresses []WeightedAddress) *ChangeOracleValidatorWeightsProposal {
	return &ChangeOracleValidatorWeightsProposal{
		Title:             title,
		Description:       description,
		WeightedAddresses: weightedAddresses,
	}
}

// GetTitle returns the title of oracle validator weights change proposal
func (m *ChangeOracleValidatorWeightsProposal) GetTitle() string {
	return m.Title
}

// GetDescription returns the description of oracle validator weights change proposal
func (m *ChangeOracleValidatorWeightsProposal) GetDescription() string {
	return m.Description
}

// ProposalRoute returns the proposal route for the oracle validator weights change proposal
func (m *ChangeOracleValidatorWeightsProposal) ProposalRoute() string {
	return RouterKey
}

// ProposalType returns the proposal type for the oracle validator weights change proposal
func (m *ChangeOracleValidatorWeightsProposal) ProposalType() string {
	return ProposalTypeChangeOracleValidatorWeights
}

// ValidateBasic runs basic stateless validity checks
func (m *ChangeOracleValidatorWeightsProposal) ValidateBasic() error {
	err := govTypes.ValidateAbstract(m)
	if err != nil {
		return err
	}

	if len(m.WeightedAddresses) == 0 {
		return fmt.Errorf("address should be more than zero")
	}

	weightSum := sdk.NewDec(0)
	for i, w := range m.WeightedAddresses {
		if w.Address != "" {
			_, err := sdk.ValAddressFromBech32(w.Address)
			if err != nil {
				return fmt.Errorf("invalid address at %dth", i)
			}
		}
		if !w.Weight.IsPositive() {
			return fmt.Errorf("non-positive weight at %dth", i)
		}
		if w.Weight.GT(sdk.NewDec(1)) {
			return fmt.Errorf("more than 1 weight at %dth", i)
		}
		weightSum = weightSum.Add(w.Weight)
	}

	if !weightSum.Equal(sdk.NewDec(1)) {
		return fmt.Errorf("invalid weight sum: %s", weightSum.String())
	}

	return nil
}

// String returns the string of proposal details
func (m *ChangeOracleValidatorWeightsProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Update Pool Incentives Proposal:
  Title:                 %s
  Description:           %s
`, m.Title, m.Description))
	return b.String()
}
