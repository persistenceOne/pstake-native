package types

import (
	"fmt"
	"strings"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalTypeRegisterHostChain      = "RegisterHostChain"
	ProposalTypeMinDepositAndFeeChange = "MinDepositAndFeeChange"
	ProposalPstakeFeeAddressChange     = "PstakeFeeAddressChange"
)

var (
	_ govtypes.Content = &RegisterHostChainProposal{}
	_ govtypes.Content = &MinDepositAndFeeChangeProposal{}
	_ govtypes.Content = &PstakeFeeAddressChangeProposal{}
)

func init() {
	govtypes.RegisterProposalType(ProposalTypeRegisterHostChain)
	govtypes.RegisterProposalTypeCodec(&RegisterHostChainProposal{}, "persistenceCore/RegisterHostChain")
	govtypes.RegisterProposalType(ProposalTypeMinDepositAndFeeChange)
	govtypes.RegisterProposalTypeCodec(&MinDepositAndFeeChangeProposal{}, "persistenceCore/MinDepositAndFeeChange")
	govtypes.RegisterProposalType(ProposalPstakeFeeAddressChange)
	govtypes.RegisterProposalTypeCodec(&PstakeFeeAddressChangeProposal{}, "persistenceCore/PstakeFeeAddressChange")
}

// NewRegisterHostChainProposal creates a new multisig change proposal.
func NewRegisterHostChainProposal(title, description string, moduleEnabled bool, chainID, connectionID, transferChannel,
	transferPort, baseDenom, mintDenom, pstakeFeeAddress string, minDeposit sdktypes.Int, allowListedValidators AllowListedValidators,
	pstakeDepositFee, pstakeRestakeFee, pstakeUnstakeFee sdktypes.Dec) *RegisterHostChainProposal {

	return &RegisterHostChainProposal{
		Title:                 title,
		Description:           description,
		ModuleEnabled:         moduleEnabled,
		ChainID:               chainID,
		ConnectionID:          connectionID,
		TransferChannel:       transferChannel,
		TransferPort:          transferPort,
		BaseDenom:             baseDenom,
		MintDenom:             mintDenom,
		MinDeposit:            minDeposit,
		AllowListedValidators: allowListedValidators,
		PstakeDepositFee:      pstakeDepositFee,
		PstakeRestakeFee:      pstakeRestakeFee,
		PstakeUnstakeFee:      pstakeUnstakeFee,
		PstakeFeeAddress:      pstakeFeeAddress,
	}
}

// GetTitle returns the title of the multisig change proposal.
func (m *RegisterHostChainProposal) GetTitle() string {
	return m.Title
}

// GetDescription returns the description of multisig change proposal.
func (m *RegisterHostChainProposal) GetDescription() string {
	return m.Description
}

// ProposalRoute returns the proposal route of multisig change proposal.
func (m *RegisterHostChainProposal) ProposalRoute() string {
	return RouterKey
}

// ProposalType returns the proposal type of multisig change proposal.
func (m *RegisterHostChainProposal) ProposalType() string {
	return ProposalTypeRegisterHostChain
}

// ValidateBasic runs basic stateless validity checks
func (m *RegisterHostChainProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(m)
	if err != nil {
		return err
	}

	return nil
}

// String returns the string of proposal details
func (m *RegisterHostChainProposal) String() string {
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

func NewHostChainParams(chainID, connectionID, channel, port, baseDenom, mintDenom, pstakefeeAddress string, minDeposit sdktypes.Int, pstakeDepositFee, pstakeRestakeFee, pstakeUnstakeFee sdktypes.Dec) HostChainParams {
	return HostChainParams{
		ChainID:          chainID,
		ConnectionID:     connectionID,
		TransferChannel:  channel,
		TransferPort:     port,
		BaseDenom:        baseDenom,
		MintDenom:        mintDenom,
		MinDeposit:       minDeposit,
		PstakeDepositFee: pstakeDepositFee,
		PstakeRestakeFee: pstakeRestakeFee,
		PstakeUnstakeFee: pstakeUnstakeFee,
		PstakeFeeAddress: pstakefeeAddress,
	}
}

// IsEmpty Checks if HostChainParams were initialised
func (c *HostChainParams) IsEmpty() bool {
	if c.TransferChannel == "" ||
		c.TransferPort == "" ||
		c.ConnectionID == "" ||
		c.ChainID == "" ||
		c.BaseDenom == "" ||
		c.MintDenom == "" ||
		c.PstakeFeeAddress == "" {
		return true
	}
	// can add more, but this should be good enough

	return false
}

// NewMinDepositAndFeeChangeProposal creates a protocol fee and min deposit change proposal.
func NewMinDepositAndFeeChangeProposal(title, description string, minDeposit sdktypes.Int, pstakeDepositFee,
	pstakeRestakeFee, pstakeUnstakeFee sdktypes.Dec) *MinDepositAndFeeChangeProposal {

	return &MinDepositAndFeeChangeProposal{
		Title:            title,
		Description:      description,
		MinDeposit:       minDeposit,
		PstakeDepositFee: pstakeDepositFee,
		PstakeRestakeFee: pstakeRestakeFee,
		PstakeUnstakeFee: pstakeUnstakeFee,
	}
}

// GetTitle returns the title of the min-deposit and fee change proposal.
func (m *MinDepositAndFeeChangeProposal) GetTitle() string {
	return m.Title
}

// GetDescription returns the description of the min-deposit and fee change proposal.
func (m *MinDepositAndFeeChangeProposal) GetDescription() string {
	return m.Description
}

// ProposalRoute returns the proposal-route of the min-deposit and fee change proposal.
func (m *MinDepositAndFeeChangeProposal) ProposalRoute() string {
	return RouterKey
}

// ProposalType returns the proposal-type of the min-deposit and fee change proposal.
func (m *MinDepositAndFeeChangeProposal) ProposalType() string {
	return ProposalTypeMinDepositAndFeeChange
}

// ValidateBasic runs basic stateless validity checks
func (m *MinDepositAndFeeChangeProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(m)
	if err != nil {
		return err
	}

	return nil
}

// String returns the string of proposal details
func (m *MinDepositAndFeeChangeProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`MinDepositAndFeeChange:
Title:                 %s
Description:           %s
MinDeposit:             %s
PstakeDepositFee:	   %s
PstakeRestakeFee: 	   %s
PstakeUnstakeFee: 	   %s

`,
		m.Title,
		m.Description,
		m.MinDeposit,
		m.PstakeDepositFee,
		m.PstakeRestakeFee,
		m.PstakeUnstakeFee),
	)
	return b.String()
}

// NewPstakeFeeAddressChangeProposal creates a pstake fee  address change proposal.
func NewPstakeFeeAddressChangeProposal(title, description,
	pstakeFeeAddress string) *PstakeFeeAddressChangeProposal {
	return &PstakeFeeAddressChangeProposal{
		Title:            title,
		Description:      description,
		PstakeFeeAddress: pstakeFeeAddress,
	}
}

// GetTitle returns the title of fee collector pstake fee address change proposal.
func (m *PstakeFeeAddressChangeProposal) GetTitle() string {
	return m.Title
}

// GetDescription returns the description of the pstake fee address proposal.
func (m *PstakeFeeAddressChangeProposal) GetDescription() string {
	return m.Description
}

// ProposalRoute returns the proposal-route of pstake fee address proposal.
func (m *PstakeFeeAddressChangeProposal) ProposalRoute() string {
	return RouterKey
}

// ProposalType returns the proposal-type of pstake fee address change proposal.
func (m *PstakeFeeAddressChangeProposal) ProposalType() string {
	return ProposalPstakeFeeAddressChange
}

// ValidateBasic runs basic stateless validity checks
func (m *PstakeFeeAddressChangeProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(m)
	if err != nil {
		return err
	}

	return nil
}

// String returns the string of proposal details
func (m *PstakeFeeAddressChangeProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`PstakeFeeAddressChange:
Title:                 %s
Description:           %s
PstakeFeeAddress: 	   %s

`,
		m.Title,
		m.Description,
		m.PstakeFeeAddress,
	),
	)
	return b.String()
}
