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

)

var _ govtypes.Content = &ChangeMultisigProposal{}

func init() {
	govtypes.RegisterProposalType(ProposalTypeChangeMultisig)
	govtypes.RegisterProposalTypeCodec(&ChangeMultisigProposal{}, "persistenceCore/ChangeMultisigProposal")
}

func NewChangeMultisigProposal(title, description string, threshold uint32, orchestratorAddresses []string, accountNumber uint64) *ChangeMultisigProposal {
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
