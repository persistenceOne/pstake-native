package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgSetOrchestrator{}
	_ sdk.Msg = &MsgSendWithFees{}
	_ sdk.Msg = &MsgVoteWithFees{}
	_ sdk.Msg = &MsgDelegateWithFees{}
	_ sdk.Msg = &MsgUndelegateWithFees{}
	_ sdk.Msg = &MsgMintTokensForAccount{}
)

// NewMsgSetOrchestrator returns a new MsgSetOrchestrator
func NewMsgSetOrchestrator(val sdk.ValAddress, operator sdk.AccAddress) *MsgSetOrchestrator {
	return &MsgSetOrchestrator{
		Validator:    val.String(),
		Orchestrator: operator.String(),
	}
}

// Route should return the name of the module
func (m *MsgSetOrchestrator) Route() string { return RouterKey }

// Type should return the action
func (m *MsgSetOrchestrator) Type() string { return "msg_set_orchestrator" }

// ValidateBasic performs stateless checks
func (m *MsgSetOrchestrator) ValidateBasic() error {
	if _, err := sdk.ValAddressFromBech32(m.Validator); err != nil {
		return sdkErrors.Wrap(sdkErrors.ErrInvalidAddress, m.Validator)
	}
	if _, err := sdk.AccAddressFromBech32(m.Orchestrator); err != nil {
		return sdkErrors.Wrap(sdkErrors.ErrInvalidAddress, m.Orchestrator)
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m *MsgSetOrchestrator) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgSetOrchestrator) GetSigners() []sdk.AccAddress {
	acc, err := sdk.ValAddressFromBech32(m.Validator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sdk.AccAddress(acc)}
}

// NewMsgSendWithFees returns a new MsgSendWithFees
func NewMsgSendWithFees(from string, to string, amount sdk.Coins, fees sdk.Coin) *MsgSendWithFees {
	return &MsgSendWithFees{
		MessageSend: &MsgSend{
			FromAddress: from,
			ToAddress:   to,
			Amount:      amount,
		},
		Fees: fees,
	}
}

// Route should return the name of the module
func (m *MsgSendWithFees) Route() string { return RouterKey }

// Type should return the action
func (m *MsgSendWithFees) Type() string { return "msg_send_with_fees" }

// ValidateBasic performs stateless checks
func (m *MsgSendWithFees) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.MessageSend.FromAddress)
	if err != nil {
		return sdkErrors.Wrapf(sdkErrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(m.MessageSend.ToAddress)
	if err != nil {
		return sdkErrors.Wrapf(sdkErrors.ErrInvalidAddress, "Invalid recipient address (%s)", err)
	}

	if !m.MessageSend.Amount.IsValid() {
		return sdkErrors.Wrap(sdkErrors.ErrInvalidCoins, m.MessageSend.Amount.String())
	}

	if !m.MessageSend.Amount.IsAllPositive() {
		return sdkErrors.Wrap(sdkErrors.ErrInvalidCoins, m.MessageSend.Amount.String())
	}

	if !m.Fees.Amount.IsPositive() {
		return fmt.Errorf("fees %s amount is not positive", m.Fees.Denom)
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (m *MsgSendWithFees) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgSendWithFees) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.MessageSend.FromAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// NewMsgVoteWithFees returns a new MsgVoteWithFees
func NewMsgVoteWithFees(id uint64, voter string, option VoteOption, fees sdk.Coin) *MsgVoteWithFees {
	return &MsgVoteWithFees{
		MessageVote: &MsgVote{
			ProposalId: id,
			Voter:      voter,
			Option:     option,
		},
		Fees: fees,
	}
}

// Route should return the name of the module
func (m *MsgVoteWithFees) Route() string { return RouterKey }

// Type should return the action
func (m *MsgVoteWithFees) Type() string { return "msg_vote_with_fees" }

// ValidateBasic performs stateless checks
func (m *MsgVoteWithFees) ValidateBasic() error {
	if m.MessageVote.Voter == "" {
		return sdkErrors.Wrap(sdkErrors.ErrInvalidAddress, m.MessageVote.Voter)
	}

	if !ValidVoteOption(m.MessageVote.Option) {
		return sdkErrors.Wrap(ErrInvalidVote, m.MessageVote.Option.String())
	}

	if !m.Fees.Amount.IsPositive() {
		return fmt.Errorf("fees %s amount is not positive", m.Fees.Denom)
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (m *MsgVoteWithFees) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgVoteWithFees) GetSigners() []sdk.AccAddress {
	voter, _ := sdk.AccAddressFromBech32(m.MessageVote.Voter)
	return []sdk.AccAddress{voter}
}

// NewMsgDelegateWithFees returns a new MsgDelegateWithFees
func NewMsgDelegateWithFees(delegator string, validator string, amount sdk.Coin, fees sdk.Coin) *MsgDelegateWithFees {
	return &MsgDelegateWithFees{
		MessageDelegate: &MsgDelegate{
			DelegatorAddress: delegator,
			ValidatorAddress: validator,
			Amount:           amount,
		},
		Fees: fees,
	}
}

// Route should return the name of the module
func (m *MsgDelegateWithFees) Route() string { return RouterKey }

// Type should return the action
func (m *MsgDelegateWithFees) Type() string { return "msg_delegate_with_fees" }

// ValidateBasic performs stateless checks
func (m *MsgDelegateWithFees) ValidateBasic() error {
	if m.MessageDelegate.DelegatorAddress == "" {
		return ErrEmptyDelegatorAddr
	}

	if m.MessageDelegate.ValidatorAddress == "" {
		return ErrEmptyValidatorAddr
	}

	if !m.MessageDelegate.Amount.IsValid() || !m.MessageDelegate.Amount.Amount.IsPositive() {
		return sdkErrors.Wrap(
			sdkErrors.ErrInvalidRequest,
			"invalid delegation amount",
		)
	}

	if !m.Fees.Amount.IsPositive() {
		return fmt.Errorf("fees %s amount is not positive", m.Fees.Denom)
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (m *MsgDelegateWithFees) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgDelegateWithFees) GetSigners() []sdk.AccAddress {
	delAddr, err := sdk.AccAddressFromBech32(m.MessageDelegate.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{delAddr}
}

// NewMsgUndelegateWithFees returns a new MsgUndelegateWithFees
func NewMsgUndelegateWithFees() *MsgUndelegateWithFees {
	return &MsgUndelegateWithFees{}
}

// Route should return the name of the module
func (m *MsgUndelegateWithFees) Route() string { return RouterKey }

// Type should return the action
func (m *MsgUndelegateWithFees) Type() string { return "msg_undelegate_with_fees" }

// ValidateBasic performs stateless checks
func (m *MsgUndelegateWithFees) ValidateBasic() error {
	if m.MessageUndelegate.DelegatorAddress == "" {
		return ErrEmptyDelegatorAddr
	}

	if m.MessageUndelegate.ValidatorAddress == "" {
		return ErrEmptyValidatorAddr
	}

	if !m.MessageUndelegate.Amount.IsValid() || !m.MessageUndelegate.Amount.Amount.IsPositive() {
		return sdkErrors.Wrap(
			sdkErrors.ErrInvalidRequest,
			"invalid shares amount",
		)
	}

	if !m.Fees.Amount.IsPositive() {
		return fmt.Errorf("fees %s amount is not positive", m.Fees.Denom)
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (m *MsgUndelegateWithFees) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgUndelegateWithFees) GetSigners() []sdk.AccAddress {
	delAddr, err := sdk.AccAddressFromBech32(m.MessageUndelegate.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{delAddr}
}

// NewMsgMintTokensForAccount returns a new MsgMintTokensForAccount
func NewMsgMintTokensForAccount(address sdk.AccAddress, orchAddress sdk.AccAddress, amount sdk.Coins, txHash string, chainID string, blockHeight int64) *MsgMintTokensForAccount {
	return &MsgMintTokensForAccount{
		AddressFromMemo:     address.String(),
		OrchestratorAddress: orchAddress.String(),
		Amount:              amount,
		TxHash:              txHash,
		ChainID:             chainID,
		BlockHeight:         blockHeight,
	}
}

// Route should return the name of the module
func (m *MsgMintTokensForAccount) Route() string { return RouterKey }

// Type should return the action
func (m *MsgMintTokensForAccount) Type() string { return "msg_mint_coins" }

// ValidateBasic performs stateless checks
func (m *MsgMintTokensForAccount) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.AddressFromMemo); err != nil {
		return sdkErrors.Wrap(sdkErrors.ErrInvalidAddress, m.AddressFromMemo)
	}
	if _, err := sdk.AccAddressFromBech32(m.OrchestratorAddress); err != nil {
		return sdkErrors.Wrap(sdkErrors.ErrInvalidAddress, m.OrchestratorAddress)
	}
	if !m.Amount.IsValid() {
		return sdkErrors.Wrap(sdkErrors.ErrInvalidCoins, m.Amount.String())
	}

	if !m.Amount.IsAllPositive() {
		return sdkErrors.Wrap(sdkErrors.ErrInvalidCoins, m.Amount.String())
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (m *MsgMintTokensForAccount) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgMintTokensForAccount) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.OrchestratorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

var _ MintTokensForAccountInterface = &IncomingMintTx{}

func (m *IncomingMintTx) Find(orchAddress string) bool {
	for _, address := range m.OrchAddresses {
		if address == orchAddress {
			return true
		}
	}
	return false
}

func (m *IncomingMintTx) AddAndIncrement(orchAddress string) {
	m.OrchAddresses = append(m.OrchAddresses, orchAddress)
	m.Counter++
}
