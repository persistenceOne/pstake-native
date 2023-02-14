package types

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	govv1beta1types "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"strings"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	ProposalTypeRegisterHostChain         = "RegisterHostChain"
	ProposalTypeMinDepositAndFeeChange    = "MinDepositAndFeeChange"
	ProposalPstakeFeeAddressChange        = "PstakeFeeAddressChange"
	ProposalAllowListedValidatorSetChange = "AllowListedValidatorSetChange"
)

var (
	_ govv1beta1types.Content = &RegisterHostChainProposal{}
	_ govv1beta1types.Content = &MinDepositAndFeeChangeProposal{}
	_ govv1beta1types.Content = &PstakeFeeAddressChangeProposal{}
	_ govv1beta1types.Content = &AllowListedValidatorSetChangeProposal{}
)

func init() {
	amino := codec.NewLegacyAmino()
	govv1beta1types.RegisterProposalType(ProposalTypeRegisterHostChain)
	amino.RegisterConcrete(&RegisterHostChainProposal{}, "pstake/RegisterHostChain", nil)
	govv1beta1types.RegisterProposalType(ProposalTypeMinDepositAndFeeChange)
	amino.RegisterConcrete(&MinDepositAndFeeChangeProposal{}, "pstake/MinDepositAndFeeChange", nil)
	govv1beta1types.RegisterProposalType(ProposalPstakeFeeAddressChange)
	amino.RegisterConcrete(&PstakeFeeAddressChangeProposal{}, "pstake/PstakeFeeAddressChange", nil)
	govv1beta1types.RegisterProposalType(ProposalAllowListedValidatorSetChange)
	amino.RegisterConcrete(&AllowListedValidatorSetChangeProposal{}, "pstake/AllowListedValidatorSetChange", nil)
}

// NewRegisterHostChainProposal creates a new host chain register proposal.
func NewRegisterHostChainProposal(title, description string, moduleEnabled bool, chainID, connectionID, transferChannel,
	transferPort, baseDenom, mintDenom, pstakeFeeAddress string, minDeposit sdktypes.Int, allowListedValidators AllowListedValidators,
	pstakeDepositFee, pstakeRestakeFee, pstakeUnstakeFee, pstakeRedemptionFee sdktypes.Dec) *RegisterHostChainProposal {

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
		PstakeParams: PstakeParams{
			PstakeDepositFee:    pstakeDepositFee,
			PstakeRestakeFee:    pstakeRestakeFee,
			PstakeUnstakeFee:    pstakeUnstakeFee,
			PstakeRedemptionFee: pstakeRedemptionFee,
			PstakeFeeAddress:    pstakeFeeAddress,
		},
	}
}

// GetTitle returns the title of the host chain register proposal.
func (m *RegisterHostChainProposal) GetTitle() string {
	return m.Title
}

// GetDescription returns the description of host chain register proposal.
func (m *RegisterHostChainProposal) GetDescription() string {
	return m.Description
}

// ProposalRoute returns the proposal route of host chain register proposal.
func (m *RegisterHostChainProposal) ProposalRoute() string {
	return RouterKey
}

// ProposalType returns the proposal type of host chain register proposal.
func (m *RegisterHostChainProposal) ProposalType() string {
	return ProposalTypeRegisterHostChain
}

// ValidateBasic runs basic stateless validity checks
func (m *RegisterHostChainProposal) ValidateBasic() error {
	err := govv1beta1types.ValidateAbstract(m)
	if err != nil {
		return err
	}

	if !m.AllowListedValidators.Valid() {
		return sdkerrors.Wrapf(ErrInValidAllowListedValidators, "allow listed validators is not valid")
	}

	if ConvertBaseDenomToMintDenom(m.BaseDenom) != m.MintDenom {
		return ErrInvalidMintDenom
	}

	err = m.PstakeParams.Validate()
	if err != nil {
		return err
	}

	if m.MinDeposit.LTE(sdktypes.ZeroInt()) {
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
PstakeRedemptionFee:   %s

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
		m.PstakeParams.PstakeDepositFee,
		m.PstakeParams.PstakeRestakeFee,
		m.PstakeParams.PstakeUnstakeFee,
		m.PstakeParams.PstakeRedemptionFee),
	)
	return b.String()
}

// NewHostChainParams returns HostChainParams with the input provided
func NewHostChainParams(chainID, connectionID, channel, port, baseDenom, mintDenom, pstakefeeAddress string, minDeposit sdktypes.Int, pstakeDepositFee, pstakeRestakeFee, pstakeUnstakeFee, pstakeRedemptionFee sdktypes.Dec) HostChainParams {
	return HostChainParams{
		ChainID:         chainID,
		ConnectionID:    connectionID,
		TransferChannel: channel,
		TransferPort:    port,
		BaseDenom:       baseDenom,
		MintDenom:       mintDenom,
		MinDeposit:      minDeposit,
		PstakeParams: PstakeParams{
			PstakeDepositFee:    pstakeDepositFee,
			PstakeRestakeFee:    pstakeRestakeFee,
			PstakeUnstakeFee:    pstakeUnstakeFee,
			PstakeRedemptionFee: pstakeRedemptionFee,
			PstakeFeeAddress:    pstakefeeAddress,
		},
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
		c.PstakeParams.PstakeFeeAddress == "" {
		return true
	}
	// can add more, but this should be good enough

	return false
}

// NewMinDepositAndFeeChangeProposal creates a protocol fee and min deposit change proposal.
func NewMinDepositAndFeeChangeProposal(title, description string, minDeposit sdktypes.Int, pstakeDepositFee,
	pstakeRestakeFee, pstakeUnstakeFee, pstakeRedemptionFee sdktypes.Dec) *MinDepositAndFeeChangeProposal {

	return &MinDepositAndFeeChangeProposal{
		Title:               title,
		Description:         description,
		MinDeposit:          minDeposit,
		PstakeDepositFee:    pstakeDepositFee,
		PstakeRestakeFee:    pstakeRestakeFee,
		PstakeUnstakeFee:    pstakeUnstakeFee,
		PstakeRedemptionFee: pstakeRedemptionFee,
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
	err := govv1beta1types.ValidateAbstract(m)
	if err != nil {
		return err
	}

	if m.PstakeDepositFee.IsNegative() || m.PstakeDepositFee.GTE(MaxPstakeDepositFee) {
		return sdkerrors.Wrapf(ErrInvalidFee, "pstake deposit fee must be between %s and %s", sdktypes.ZeroDec(), MaxPstakeDepositFee)
	}

	if m.PstakeRestakeFee.IsNegative() || m.PstakeRestakeFee.GTE(MaxPstakeRestakeFee) {
		return sdkerrors.Wrapf(ErrInvalidFee, "pstake restake fee must be between %s and %s", sdktypes.ZeroDec(), MaxPstakeRestakeFee)
	}

	if m.PstakeUnstakeFee.IsNegative() || m.PstakeUnstakeFee.GTE(MaxPstakeUnstakeFee) {
		return sdkerrors.Wrapf(ErrInvalidFee, "pstake unstake fee must be between %s and %s", sdktypes.ZeroDec(), MaxPstakeUnstakeFee)
	}

	if m.PstakeRedemptionFee.IsNegative() || m.PstakeRedemptionFee.GTE(MaxPstakeRedemptionFee) {
		return sdkerrors.Wrapf(ErrInvalidFee, "pstake redemption fee must be between %s and %s", sdktypes.ZeroDec(), MaxPstakeRedemptionFee)
	}

	if m.MinDeposit.LTE(sdktypes.ZeroInt()) {
		return sdkerrors.Wrapf(ErrInvalidDeposit, "min deposit must be positive")
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
PstakeRedemptionFee:   %s

`,
		m.Title,
		m.Description,
		m.MinDeposit,
		m.PstakeDepositFee,
		m.PstakeRestakeFee,
		m.PstakeUnstakeFee,
		m.PstakeRedemptionFee),
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
	err := govv1beta1types.ValidateAbstract(m)
	if err != nil {
		return err
	}

	_, err = sdktypes.AccAddressFromBech32(m.PstakeFeeAddress)
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

// NewAllowListedValidatorSetChangeProposal creates a allowListed validator set change proposal.
func NewAllowListedValidatorSetChangeProposal(title, description string, allowListedValidators AllowListedValidators) *AllowListedValidatorSetChangeProposal {
	return &AllowListedValidatorSetChangeProposal{
		Title:                 title,
		Description:           description,
		AllowListedValidators: allowListedValidators,
	}
}

// GetTitle returns the title of allowListed validator set change proposal.
func (m *AllowListedValidatorSetChangeProposal) GetTitle() string {
	return m.Title
}

// GetDescription returns the description of allowListed validator set change proposal.
func (m *AllowListedValidatorSetChangeProposal) GetDescription() string {
	return m.Description
}

// ProposalRoute returns the proposal-route of allowListed validator set change proposal.
func (m *AllowListedValidatorSetChangeProposal) ProposalRoute() string {
	return RouterKey
}

// ProposalType returns the proposal-type of allowListed validator set change proposal.
func (m *AllowListedValidatorSetChangeProposal) ProposalType() string {
	return ProposalAllowListedValidatorSetChange
}

// ValidateBasic runs basic stateless validity checks
func (m *AllowListedValidatorSetChangeProposal) ValidateBasic() error {
	err := govv1beta1types.ValidateAbstract(m)
	if err != nil {
		return err
	}

	if !m.AllowListedValidators.Valid() {
		return sdkerrors.Wrapf(ErrInValidAllowListedValidators, "allow listed validators is not valid")
	}

	return nil
}

// String returns the string of proposal details
func (m *AllowListedValidatorSetChangeProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`AllowListedValidatorSetChange:
Title:                 %s
Description:           %s
AllowListedValidators: 	   %s

`,
		m.Title,
		m.Description,
		m.AllowListedValidators,
	),
	)
	return b.String()
}
