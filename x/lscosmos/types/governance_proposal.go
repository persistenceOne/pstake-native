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
func NewRegisterCosmosChainProposal(title, description string, moduleEnabled bool, connectionID, transferChannel,
	TransferPort, baseDenom, mintDenom string, minDeposit sdktypes.Int, allowListedValidators AllowListedValidators,
	pstakeDepositFee, pstakeRestakeFee, pstakeUnstakeFee sdktypes.Dec) *RegisterCosmosChainProposal {

	return &RegisterCosmosChainProposal{
		Title:                 title,
		Description:           description,
		ModuleEnabled:         moduleEnabled,
		ConnectionID:          connectionID,
		TransferChannel:       transferChannel,
		TransferPort:          TransferPort,
		BaseDenom:             baseDenom,
		MintDenom:             mintDenom,
		MinDeposit:            minDeposit,
		AllowListedValidators: allowListedValidators,
		PstakeDepositFee:      pstakeDepositFee,
		PstakeRestakeFee:      pstakeRestakeFee,
		PstakeUnstakeFee:      pstakeUnstakeFee,
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
ConnectionID:         %s
TransferChannel:  %s
TransferPort:     %s
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
		m.ConnectionID,
		m.TransferChannel,
		m.TransferPort,
		m.BaseDenom,
		m.MintDenom,
		m.AllowListedValidators,
		m.PstakeDepositFee,
		m.PstakeRestakeFee,
		m.PstakeUnstakeFee),
	)
	return b.String()
}

func NewCosmosParams(connectionID, channel, port, baseDenom, mintDenom string, minDeposit sdktypes.Int, pstakeDepositFee, pstakeRestakeFee, pstakeUnstakeFee sdktypes.Dec) CosmosParams {
	return CosmosParams{
		ConnectionID:     connectionID,
		TransferChannel:  channel,
		TransferPort:     port,
		BaseDenom:        baseDenom,
		MintDenom:        mintDenom,
		MinDeposit:       minDeposit,
		PstakeDepositFee: pstakeDepositFee,
		PstakeRestakeFee: pstakeRestakeFee,
		PstakeUnstakeFee: pstakeUnstakeFee,
	}
}

// IsEmpty Checks if CosmosParams were initialised
func (c *CosmosParams) IsEmpty() bool {
	if c.TransferChannel == "" ||
		c.TransferPort == "" ||
		c.ConnectionID == "" ||
		c.BaseDenom == "" ||
		c.MintDenom == "" {
		return true
	}
	// can add more, but this should be good enough

	return false
}
