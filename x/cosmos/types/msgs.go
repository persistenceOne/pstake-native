package types

import (
	"github.com/ghodss/yaml"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	sdkTx "github.com/cosmos/cosmos-sdk/types/tx"
)

var (
	_ sdk.Msg = &MsgSetOrchestrator{}
	_ sdk.Msg = &MsgWithdrawStkAsset{}
	_ sdk.Msg = &MsgMintTokensForAccount{}
	_ sdk.Msg = &MsgMakeProposal{}
	_ sdk.Msg = &MsgVote{}
	_ sdk.Msg = &MsgVoteWeighted{}
	_ sdk.Msg = &MsgSignedTx{}
	_ sdk.Msg = &MsgTxStatus{}
	_ sdk.Msg = &MsgUndelegateSuccess{}
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

// NewMsgWithdrawStkAsset returns a new MsgWithdrawStkAsset
func NewMsgWithdrawStkAsset(from sdk.AccAddress, to sdk.Address, amount sdk.Coin) *MsgWithdrawStkAsset {
	return &MsgWithdrawStkAsset{
		FromAddress: from.String(),
		ToAddress:   to.String(),
		Amount:      amount,
	}
}

// Route should return the name of the module
func (m *MsgWithdrawStkAsset) Route() string { return RouterKey }

// Type should return the action
func (m *MsgWithdrawStkAsset) Type() string { return "msg_withdraw_stk_asset" }

// ValidateBasic performs stateless checks
func (m *MsgWithdrawStkAsset) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return sdkErrors.Wrapf(sdkErrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(m.ToAddress)
	if err != nil {
		return sdkErrors.Wrapf(sdkErrors.ErrInvalidAddress, "Invalid recipient address (%s)", err)
	}

	if !m.Amount.IsValid() {
		return sdkErrors.Wrap(sdkErrors.ErrInvalidCoins, m.Amount.String())
	}

	if !m.Amount.IsPositive() {
		return sdkErrors.Wrap(sdkErrors.ErrInvalidCoins, m.Amount.String())
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (m *MsgWithdrawStkAsset) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgWithdrawStkAsset) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
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
	if m.BlockHeight <= 0 {
		return sdkErrors.Wrap(sdkErrors.ErrInvalidHeight, "BlockHeight should be greater than zero")
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
func NewMsgMakeProposal(title string, description string, orchestratorAddress sdk.AccAddress, chainID string,
	blockHeight int64, proposalID uint64, votingStartTime time.Time, votingEndTime time.Time) *MsgMakeProposal {
	return &MsgMakeProposal{
		Title:               title,
		Description:         description,
		OrchestratorAddress: orchestratorAddress.String(),
		ChainID:             chainID,
		BlockHeight:         blockHeight,
		ProposalID:          proposalID,
		VotingStartTime:     votingStartTime,
		VotingEndTime:       votingEndTime,
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
	// TODO add more checks
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

	if _, err := sdk.AccAddressFromBech32(m.Voter); err != nil {
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
	if _, err := sdk.AccAddressFromBech32(m.Voter); err != nil {
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

// NewMsgSignedTx returns a new MsgSignedTx
func NewMsgSignedTx(txID uint64, tx sdkTx.Tx, orchAddress sdk.AccAddress) *MsgSignedTx {
	return &MsgSignedTx{
		TxID:                txID,
		Tx:                  tx,
		OrchestratorAddress: orchAddress.String(),
	}
}

// Route should return the name of the module
func (m *MsgSignedTx) Route() string { return RouterKey }

// Type should return the action
func (m *MsgSignedTx) Type() string { return "signed_tx" }

// ValidateBasic performs stateless checks
func (m *MsgSignedTx) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.OrchestratorAddress); err != nil {
		return sdkErrors.Wrap(sdkErrors.ErrInvalidAddress, m.OrchestratorAddress)
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m *MsgSignedTx) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgSignedTx) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.OrchestratorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// NewMsgTxStatus returns a new MsgTxStatus
func NewMsgTxStatus(orchAddress sdk.AccAddress, status string, txHash string) *MsgTxStatus {
	return &MsgTxStatus{
		OrchestratorAddress: orchAddress.String(),
		TxHash:              txHash,
		Status:              status,
	}
}

// Route should return the name of the module
func (m *MsgTxStatus) Route() string { return RouterKey }

// Type should return the action
func (m *MsgTxStatus) Type() string { return "msg_tx_status" }

// ValidateBasic performs stateless checks
func (m *MsgTxStatus) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.OrchestratorAddress); err != nil {
		return sdkErrors.Wrap(sdkErrors.ErrInvalidAddress, m.OrchestratorAddress)
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m *MsgTxStatus) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgTxStatus) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.OrchestratorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// NewMsgUndelegateSuccess returns a new MsgUndelegateSuccess
func NewMsgUndelegateSuccess(val sdk.ValAddress, delegatorAddress sdk.AccAddress, amount sdk.Coin, orchAddress sdk.AccAddress) *MsgUndelegateSuccess {
	return &MsgUndelegateSuccess{
		ValidatorAddress:    val.String(),
		DelegatorAddress:    delegatorAddress.String(),
		Amount:              amount,
		OrchestratorAddress: orchAddress.String(),
	}
}

// Route should return the name of the module
func (m *MsgUndelegateSuccess) Route() string { return RouterKey }

// Type should return the action
func (m *MsgUndelegateSuccess) Type() string { return "msg_undelegation_success" }

// ValidateBasic performs stateless checks
func (m *MsgUndelegateSuccess) ValidateBasic() error {
	if _, err := sdk.ValAddressFromBech32(m.ValidatorAddress); err != nil {
		return sdkErrors.Wrap(sdkErrors.ErrInvalidAddress, m.ValidatorAddress)
	}
	if _, err := sdk.AccAddressFromBech32(m.OrchestratorAddress); err != nil {
		return sdkErrors.Wrap(sdkErrors.ErrInvalidAddress, m.OrchestratorAddress)
	}
	if _, err := sdk.AccAddressFromBech32(m.DelegatorAddress); err != nil {
		return sdkErrors.Wrap(sdkErrors.ErrInvalidAddress, m.DelegatorAddress)
	}
	if !m.Amount.IsValid() || !m.Amount.Amount.IsPositive() {
		return sdkErrors.Wrap(
			sdkErrors.ErrInvalidRequest,
			"invalid delegation amount",
		)
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m *MsgUndelegateSuccess) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgUndelegateSuccess) GetSigners() []sdk.AccAddress {
	acc, err := sdk.ValAddressFromBech32(m.OrchestratorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sdk.AccAddress(acc)}
}

// TODO find a better place for this.
var _ DBHelper = &IncomingMintTx{}
var _ DBHelper = &ProposalValue{}
var _ DBHelper = &TxHashValue{}
var _ DBHelper = &ValueUndelegateSuccessStore{}

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

func (m *TxHashValue) Find(orchAddress string) bool {
	for _, address := range m.OrchestratorAddresses {
		if address == orchAddress {
			return true
		}
	}
	return false
}

func (m *TxHashValue) AddAndIncrement(orchAddress string) {
	m.OrchestratorAddresses = append(m.OrchestratorAddresses, orchAddress)
	m.Counter++
}

func (m *ValueUndelegateSuccessStore) Find(orchAddress string) bool {
	for _, address := range m.OrchestratorAddresses {
		if address == orchAddress {
			return true
		}
	}
	return false
}

func (m *ValueUndelegateSuccessStore) AddAndIncrement(orchAddress string) {
	m.OrchestratorAddresses = append(m.OrchestratorAddresses, orchAddress)
	m.Counter++
}
