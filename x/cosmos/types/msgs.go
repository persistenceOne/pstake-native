package types

import (
	"fmt"
	"github.com/ghodss/yaml"

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
	_ sdk.Msg = &MsgMakeProposal{}
	_ sdk.Msg = &MsgVote{}
	_ sdk.Msg = &MsgVoteWeighted{}
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

// NewMsgMakeProposal returns a new MsgMakeProposal
func NewMsgMakeProposal(title string, description string, orchestratorAddress sdk.AccAddress, chainID string, blockHeight int64, proposalID int64) *MsgMakeProposal {
	return &MsgMakeProposal{
		Title:               title,
		Description:         description,
		OrchestratorAddress: orchestratorAddress.String(),
		ChainID:             chainID,
		BlockHeight:         blockHeight,
		ProposalID:          proposalID,
	}
}

// Route should return the name of the module
func (m *MsgMakeProposal) Route() string { return RouterKey }

// Type should return the action
func (m *MsgMakeProposal) Type() string { return "msg_make_proposal" }

// ValidateBasic performs stateless checks
func (m *MsgMakeProposal) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.OrchestratorAddress); err != nil {
		return sdkErrors.Wrap(sdkErrors.ErrInvalidAddress, m.OrchestratorAddress)
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m *MsgMakeProposal) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgMakeProposal) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.OrchestratorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

//nolint:interfacer
func NewMsgVote(voter sdk.AccAddress, proposalID uint64, option VoteOption) *MsgVote {
	return &MsgVote{proposalID, voter.String(), option}
}

// Route implements Msg
func (m *MsgVote) Route() string { return RouterKey }

// Type implements Msg
func (m *MsgVote) Type() string { return "vote" }

// ValidateBasic implements Msg
func (m *MsgVote) ValidateBasic() error {
	if m.Voter == "" {
		return sdkErrors.Wrap(sdkErrors.ErrInvalidAddress, m.Voter)
	}

	if !ValidVoteOption(m.Option) {
		return sdkErrors.Wrap(ErrInvalidVote, m.Option.String())
	}

	return nil
}

// String implements the Stringer interface
func (m *MsgVote) String() string {
	out, _ := yaml.Marshal(m)
	return string(out)
}

// GetSignBytes implements Msg
func (m *MsgVote) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements Msg
func (m *MsgVote) GetSigners() []sdk.AccAddress {
	voter, _ := sdk.AccAddressFromBech32(m.Voter)
	return []sdk.AccAddress{voter}
}

// NewMsgVoteWeighted creates a message to cast a vote on an active proposal
//nolint:interfacer
func NewMsgVoteWeighted(voter sdk.AccAddress, proposalID uint64, options WeightedVoteOptions) *MsgVoteWeighted {
	return &MsgVoteWeighted{proposalID, voter.String(), options}
}

// Route implements Msg
func (m *MsgVoteWeighted) Route() string { return RouterKey }

// Type implements Msg
func (m *MsgVoteWeighted) Type() string { return "weighted_vote" }

// ValidateBasic implements Msg
func (m *MsgVoteWeighted) ValidateBasic() error {
	if m.Voter == "" {
		return sdkErrors.Wrap(sdkErrors.ErrInvalidAddress, m.Voter)
	}

	if len(m.Options) == 0 {
		return sdkErrors.Wrap(sdkErrors.ErrInvalidRequest, WeightedVoteOptions(m.Options).String())
	}

	totalWeight := sdk.NewDec(0)
	usedOptions := make(map[VoteOption]bool)
	for _, option := range m.Options {
		if !ValidWeightedVoteOption(option) {
			return sdkErrors.Wrap(ErrInvalidVote, option.String())
		}
		totalWeight = totalWeight.Add(option.Weight)
		if usedOptions[option.Option] {
			return sdkErrors.Wrap(ErrInvalidVote, "Duplicated vote option")
		}
		usedOptions[option.Option] = true
	}

	if totalWeight.GT(sdk.NewDec(1)) {
		return sdkErrors.Wrap(ErrInvalidVote, "Total weight overflow 1.00")
	}

	if totalWeight.LT(sdk.NewDec(1)) {
		return sdkErrors.Wrap(ErrInvalidVote, "Total weight lower than 1.00")
	}

	return nil
}

// String implements the Stringer interface
func (m *MsgVoteWeighted) String() string {
	out, _ := yaml.Marshal(m)
	return string(out)
}

// GetSignBytes implements Msg
func (m *MsgVoteWeighted) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements Msg
func (m *MsgVoteWeighted) GetSigners() []sdk.AccAddress {
	voter, _ := sdk.AccAddressFromBech32(m.Voter)
	return []sdk.AccAddress{voter}
}

var _ DBHelper = &IncomingMintTx{}
var _ DBHelper = &ProposalValue{}

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

func NewProposalValue(title string, description string, orchAddress string, ratio float32) ProposalValue {
	return ProposalValue{
		Title:                 title,
		Description:           description,
		OrchestratorAddresses: []string{orchAddress},
		Ratio:                 ratio,
		Counter:               1,
		ProposalPosted:        false,
	}
}
func (m *ProposalValue) Find(orchAddress string) bool {
	for _, address := range m.OrchestratorAddresses {
		if address == orchAddress {
			return true
		}
	}
	return false
}

func (m *ProposalValue) AddAndIncrement(orchAddress string) {
	m.OrchestratorAddresses = append(m.OrchestratorAddresses, orchAddress)
	m.Counter++
}
