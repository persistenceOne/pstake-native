package types

import (
	"fmt"
	"strings"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalTypeRegisterHostChain = "RegisterHostChain"
)

var _ govtypes.Content = &RegisterHostChainProposal{}

func init() {
	govtypes.RegisterProposalType(ProposalTypeRegisterHostChain)
	govtypes.RegisterProposalTypeCodec(&RegisterHostChainProposal{}, "persistenceCore/RegisterHostChain")
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

	if !m.AllowListedValidators.Valid() {
		return sdkerrors.Wrapf(ErrInValidAllowListedValidators, "allow listed validators is not valid")
	}

	_, err = sdktypes.AccAddressFromBech32(m.PstakeFeeAddress)
	if err != nil {
		return err
	}

	if m.PstakeDepositFee.IsNegative() {
		return sdkerrors.Wrapf(ErrInvalidFee, "pstake deposit fee must be positive")
	}

	if m.PstakeRestakeFee.IsNegative() {
		return sdkerrors.Wrapf(ErrInvalidFee, "pstake restake fee must be positive")
	}

	if m.PstakeUnstakeFee.IsNegative() {
		return sdkerrors.Wrapf(ErrInvalidFee, "pstake unstake fee must be positive")
	}

	if m.MinDeposit.IsNegative() || m.MinDeposit.IsZero() {
		return sdkerrors.Wrapf(ErrInvalidDeposit, "min deposit must be positive")
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
