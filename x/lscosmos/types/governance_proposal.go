package types

import (
	"fmt"
	"strings"

	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalTypeRegisterCosmosChain = "RegisterCosmosChain"
)

var _ govTypes.Content = &RegisterCosmosChainProposal{}

func init() {
	govTypes.RegisterProposalType(ProposalTypeRegisterCosmosChain)
	govTypes.RegisterProposalTypeCodec(&RegisterCosmosChainProposal{}, "persistenceCore/RegisterCosmosChain")
}

// NewRegisterCosmosChainProposal creates a new multisig change proposal.
func NewRegisterCosmosChainProposal(title, description, ibcConnection, tokenTransferChannel, tokenTransferPort, baseDenom, mintDenom string) *RegisterCosmosChainProposal {
	return &RegisterCosmosChainProposal{
		Title:                title,
		Description:          description,
		IBCConnection:        ibcConnection,
		TokenTransferChannel: tokenTransferChannel,
		TokenTransferPort:    tokenTransferPort,
		BaseDenom:            baseDenom,
		MintDenom:            mintDenom,
	}
}

// GetTitle returns the title of the multisig change proposal.
func (m *RegisterCosmosChainProposal) GetTitle() string {
	return m.Title
}

// GetDescription returns the description of multisig change proposal.
func (m *RegisterCosmosChainProposal) GetDescription() string {
	return m.Description
}

// ProposalRoute returns the proposal route of multisig change proposal.
func (m *RegisterCosmosChainProposal) ProposalRoute() string {
	return RouterKey
}

// ProposalType returns the proposal type of multisig change proposal.
func (m *RegisterCosmosChainProposal) ProposalType() string {
	return ProposalTypeRegisterCosmosChain
}

// ValidateBasic runs basic stateless validity checks
func (m *RegisterCosmosChainProposal) ValidateBasic() error {
	err := govTypes.ValidateAbstract(m)
	if err != nil {
		return err
	}

	return nil
}

// String returns the string of proposal details
func (m *RegisterCosmosChainProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Update Pool Incentives Proposal:
Title:                 %s
Description:           %s
IBCConnection:         %s
TokenTransferChannel:  %s
TokenTransferPort:     %s
BaseDenom: 			 %s
MintDenom: 			 %s
`,
		m.Title,
		m.Description,
		m.IBCConnection,
		m.TokenTransferChannel,
		m.TokenTransferPort,
		m.BaseDenom,
		m.MintDenom),
	)
	return b.String()
}

func NewCosmosIBCParams(ibcConnection, channel, port, baseDenom, mintDenom string) CosmosIBCParams {
	return CosmosIBCParams{
		IBCConnection:        ibcConnection,
		TokenTransferChannel: channel,
		TokenTransferPort:    port,
		BaseDenom:            baseDenom,
		MintDenom:            mintDenom,
	}
}
