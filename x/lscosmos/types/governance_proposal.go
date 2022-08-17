package types

import (
	"fmt"
	"strings"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalTypeRegisterCosmosChain = "RegisterCosmosChain"
)

var _ govtypes.Content = &RegisterCosmosChainProposal{}

func init() {
	govtypes.RegisterProposalType(ProposalTypeRegisterCosmosChain)
	govtypes.RegisterProposalTypeCodec(&RegisterCosmosChainProposal{}, "persistenceCore/RegisterCosmosChain")
}

// NewRegisterCosmosChainProposal creates a new multisig change proposal.
func NewRegisterCosmosChainProposal(title, description string, moduleEnabled bool, ibcConnection, tokenTransferChannel,
	tokenTransferPort, baseDenom, mintDenom string, minDeposit sdktypes.Int, allowListedValidators AllowListedValidators,
	pStakeDepositFee, pstakeRestakeFee, pStakeUnstakeFee sdktypes.Dec) *RegisterCosmosChainProposal {

	return &RegisterCosmosChainProposal{
		Title:                 title,
		Description:           description,
		ModuleEnabled:         moduleEnabled,
		IBCConnection:         ibcConnection,
		TokenTransferChannel:  tokenTransferChannel,
		TokenTransferPort:     tokenTransferPort,
		BaseDenom:             baseDenom,
		MintDenom:             mintDenom,
		MinDeposit:            minDeposit,
		AllowListedValidators: allowListedValidators,
		PStakeDepositFee:      pStakeDepositFee,
		PStakeRestakeFee:      pstakeRestakeFee,
		PStakeUnstakeFee:      pStakeUnstakeFee,
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
	err := govtypes.ValidateAbstract(m)
	if err != nil {
		return err
	}

	return nil
}

// String returns the string of proposal details
func (m *RegisterCosmosChainProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Register host chain:
Title:                 %s
Description:           %s
ModuleEnabled:		   %v
IBCConnection:         %s
TokenTransferChannel:  %s
TokenTransferPort:     %s
BaseDenom: 			   %s
MintDenom: 			   %s
AllowlistedValidators: %s
PstakeDepositFee:	   %s
PstakeRestakeFee: 	   %s
PstakeUnstakeFee: 	   %s

`,
		m.Title,
		m.Description,
		m.ModuleEnabled,
		m.IBCConnection,
		m.TokenTransferChannel,
		m.TokenTransferPort,
		m.BaseDenom,
		m.MintDenom,
		m.AllowListedValidators,
		m.PStakeDepositFee,
		m.PStakeRestakeFee,
		m.PStakeUnstakeFee),
	)
	return b.String()
}

func NewCosmosIBCParams(ibcConnection, channel, port, baseDenom, mintDenom string, minDeposit sdktypes.Int, pStakeDepositFee, pstakeRestakeFee, pstakeUnstakeFee sdktypes.Dec) CosmosIBCParams {
	return CosmosIBCParams{
		IBCConnection:        ibcConnection,
		TokenTransferChannel: channel,
		TokenTransferPort:    port,
		BaseDenom:            baseDenom,
		MintDenom:            mintDenom,
		MinDeposit:           minDeposit,
		PStakeDepositFee:     pStakeDepositFee,
		PStakeRestakeFee:     pstakeRestakeFee,
		PStakeUnstakeFee:     pstakeUnstakeFee,
	}
}

// Checks if cosmosIBC params were initialised
func (c *CosmosIBCParams) IsEmpty() bool {
	if c.TokenTransferChannel == "" ||
		c.TokenTransferPort == "" ||
		c.IBCConnection == "" ||
		c.BaseDenom == "" ||
		c.MintDenom == "" {
		return true
	}
	// can add more, but this should be good enough

	return false
}
