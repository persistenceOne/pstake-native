package types

import (
	"fmt"
	"strings"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

const (
	ProposalTypeMinDepositAndFeeChange    = "MinDepositAndFeeChange"
	ProposalPstakeFeeAddressChange        = "PstakeFeeAddressChange"
	ProposalAllowListedValidatorSetChange = "AllowListedValidatorSetChange"
)

var (
	_ govtypes.Content = &MinDepositAndFeeChangeProposal{}
	_ govtypes.Content = &PstakeFeeAddressChangeProposal{}
	_ govtypes.Content = &AllowListedValidatorSetChangeProposal{}
)

func init() {
	govtypes.RegisterProposalType(ProposalTypeMinDepositAndFeeChange)
	govtypes.RegisterProposalType(ProposalPstakeFeeAddressChange)
	govtypes.RegisterProposalType(ProposalAllowListedValidatorSetChange)
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
