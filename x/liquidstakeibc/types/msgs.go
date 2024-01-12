package types

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	connectiontypes "github.com/cosmos/ibc-go/v7/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
)

const (
	MsgTypeRegisterHostChain string = "msg_register_host_chain"
	MsgTypeUpdateHostChain   string = "msg_update_host_chain"
	MsgTypeLiquidStake       string = "msg_liquid_stake"
	MsgTypeLiquidStakeLSM    string = "msg_liquid_stake_lsm"
	MsgTypeLiquidUnstake     string = "msg_liquid_unstake"
	MsgTypeRedeem            string = "msg_redeem"
	MsgTypeUpdateParams      string = "msg_update_params"
)

var (
	_ sdk.Msg = &MsgRegisterHostChain{}
	_ sdk.Msg = &MsgUpdateHostChain{}
	_ sdk.Msg = &MsgLiquidStake{}
	_ sdk.Msg = &MsgLiquidUnstake{}
	_ sdk.Msg = &MsgRedeem{}
	_ sdk.Msg = &MsgLiquidStakeLSM{}
)

func NewMsgRegisterHostChain(
	connectionID string,
	channelID string,
	portID string,
	depositFee string,
	restakeFee string,
	unstakeFee string,
	redemptionFee string,
	hostDenom string,
	minimumDeposit math.Int,
	unbondingFactor int64,
	authority string,
	autocompoundFactor int64,
) *MsgRegisterHostChain {
	depositFeeDec, _ := sdk.NewDecFromStr(depositFee)
	restakeFeeDec, _ := sdk.NewDecFromStr(restakeFee)
	unstakeFeeDec, _ := sdk.NewDecFromStr(unstakeFee)
	redemptionFeeDec, _ := sdk.NewDecFromStr(redemptionFee)

	return &MsgRegisterHostChain{
		ConnectionId:       connectionID,
		HostDenom:          hostDenom,
		ChannelId:          channelID,
		PortId:             portID,
		MinimumDeposit:     minimumDeposit,
		UnbondingFactor:    unbondingFactor,
		DepositFee:         depositFeeDec,
		RestakeFee:         restakeFeeDec,
		UnstakeFee:         unstakeFeeDec,
		RedemptionFee:      redemptionFeeDec,
		Authority:          authority,
		AutoCompoundFactor: autocompoundFactor,
	}
}

func (m *MsgRegisterHostChain) Route() string {
	return RouterKey
}

func (m *MsgRegisterHostChain) Type() string {
	return MsgTypeRegisterHostChain
}

func (m *MsgRegisterHostChain) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgRegisterHostChain) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (m *MsgRegisterHostChain) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address %q: %v", m.Authority, err)
	}

	// connection id cannot be empty and must begin with "connection"
	if m.ConnectionId == "" {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "connection id cannot be empty")
	}
	if !strings.HasPrefix(m.ConnectionId, connectiontypes.ConnectionPrefix) {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, fmt.Sprintf("invalid connection id: %s, must begin with '%s'", m.ConnectionId, connectiontypes.ConnectionPrefix))
	}

	// validate host denom
	if err := sdk.ValidateDenom(m.HostDenom); err != nil {
		return errorsmod.Wrapf(
			sdkerrors.ErrInvalidRequest,
			fmt.Sprintf("invalid host denom: %s", err.Error()),
		)
	}

	// validate channel id
	if valid := strings.HasPrefix(m.ChannelId, channeltypes.ChannelPrefix); !valid {
		return errorsmod.Wrapf(
			sdkerrors.ErrInvalidRequest,
			fmt.Sprintf("invalid channel id: %s, must begin with '%s'", m.ChannelId, channeltypes.ChannelPrefix),
		)
	}

	// deposit fee must be positive or zero
	if m.DepositFee.LT(sdk.ZeroDec()) || m.DepositFee.GT(sdk.OneDec()) {
		return errorsmod.Wrapf(
			sdkerrors.ErrInvalidRequest,
			"deposit fee quantity must be greater or equal than zero and less than equal one",
		)
	}

	// restake fee must be positive or zero
	if m.RestakeFee.LT(sdk.ZeroDec()) || m.RestakeFee.GT(sdk.OneDec()) {
		return errorsmod.Wrapf(
			sdkerrors.ErrInvalidRequest,
			"restake fee quantity must be greater or equal than zero and less than equal one",
		)
	}

	// unstake fee must be positive or zero
	if m.UnstakeFee.LT(sdk.ZeroDec()) || m.UnstakeFee.GT(sdk.OneDec()) {
		return errorsmod.Wrapf(
			sdkerrors.ErrInvalidRequest,
			"unstake fee quantity must be greater or equal than zero and less than equal one",
		)
	}

	// redemption deposit must be positive or zero
	if m.RedemptionFee.LT(sdk.ZeroDec()) || m.RedemptionFee.GT(sdk.OneDec()) {
		return errorsmod.Wrapf(
			sdkerrors.ErrInvalidRequest,
			"redemption fee quantity must be greater or equal than zero and less than equal one",
		)
	}

	// minimum deposit must be at least one
	if m.MinimumDeposit.LTE(sdk.ZeroInt()) {
		return sdkerrors.ErrInvalidRequest.Wrapf(
			"minimum deposit should be greater than zero",
		)
	}

	// unbonding factor must be greater than zero
	if m.UnbondingFactor <= 0 {
		return sdkerrors.ErrInvalidRequest.Wrapf(
			"unbonding factor should be greater than zero",
		)
	}

	// autocompound factor must be greater than zero
	if m.AutoCompoundFactor <= 0 {
		return sdkerrors.ErrInvalidRequest.Wrapf(
			"autocompound factor should be greater than zero",
		)
	}

	return nil
}

func NewMsgUpdateHostChain(chainID, authority string, updates []*KVUpdate) *MsgUpdateHostChain {
	return &MsgUpdateHostChain{
		ChainId:   chainID,
		Authority: authority,
		Updates:   updates,
	}
}

func (m *MsgUpdateHostChain) Route() string {
	return RouterKey
}

func (m *MsgUpdateHostChain) Type() string {
	return MsgTypeUpdateHostChain
}

func (m *MsgUpdateHostChain) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgUpdateHostChain) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (m *MsgUpdateHostChain) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address %q: %v", m.Authority, err)
	}
	for _, update := range m.Updates {
		switch update.Key {
		case KeyAddValidator:
			var validator Validator
			err := json.Unmarshal([]byte(update.Value), &validator)
			if err != nil {
				return fmt.Errorf("unable to unmarshal validator update string")
			}
			err = validator.Validate()
			if err != nil {
				return err
			}

		case KeyRemoveValidator:
			_, _, err := bech32.DecodeAndConvert(update.Value)
			if err != nil {
				return err
			}
		case KeyValidatorUpdate:
			_, _, err := bech32.DecodeAndConvert(update.Value)
			if err != nil {
				return err
			}
		case KeyValidatorWeight:
			validator, weight, valid := strings.Cut(update.Value, ",")
			if !valid {
				return fmt.Errorf("unable to parse validator update string")
			}
			_, _, err := bech32.DecodeAndConvert(validator)
			if err != nil {
				return err
			}
			decWt, err := sdk.NewDecFromStr(weight)
			if err != nil {
				return err
			}
			if decWt.GT(sdk.OneDec()) || decWt.LT(sdk.ZeroDec()) {
				return fmt.Errorf("weight should be, 0 <= weight <= 1")
			}

		case KeyDepositFee:
			fee, err := sdk.NewDecFromStr(update.Value)
			if err != nil {
				return fmt.Errorf("unable to parse string to sdk.Dec: %w", err)
			}

			if fee.LT(sdk.ZeroDec()) || fee.GT(sdk.OneDec()) {
				return sdkerrors.ErrInvalidRequest.Wrapf("invalid deposit fee value should be 0 <= fee <= 1")
			}
		case KeyRestakeFee:
			fee, err := sdk.NewDecFromStr(update.Value)
			if err != nil {
				return fmt.Errorf("unable to parse string to sdk.Dec: %w", err)
			}

			if fee.LT(sdk.ZeroDec()) || fee.GT(sdk.OneDec()) {
				return sdkerrors.ErrInvalidRequest.Wrapf("invalid restake fee value should be 0 <= fee <= 1")
			}
		case KeyRedemptionFee:
			fee, err := sdk.NewDecFromStr(update.Value)
			if err != nil {
				return fmt.Errorf("unable to parse string to sdk.Dec: %w", err)
			}

			if fee.LT(sdk.ZeroDec()) || fee.GT(sdk.OneDec()) {
				return sdkerrors.ErrInvalidRequest.Wrapf("invalid redemption fee value should be 0 <= fee <= 1")
			}
		case KeyUnstakeFee:
			fee, err := sdk.NewDecFromStr(update.Value)
			if err != nil {
				return fmt.Errorf("unable to parse string to sdk.Dec: %w", err)
			}

			if fee.LT(sdk.ZeroDec()) || fee.GT(sdk.OneDec()) {
				return sdkerrors.ErrInvalidRequest.Wrapf("invalid unstake fee value should be 0 <= fee <= 1")
			}
		case KeyLSMValidatorCap:
			validatorCap, err := sdk.NewDecFromStr(update.Value)
			if err != nil {
				return fmt.Errorf("unable to parse string to sdk.Dec: %w", err)
			}

			if validatorCap.LT(sdk.ZeroDec()) || validatorCap.GT(sdk.OneDec()) {
				return sdkerrors.ErrInvalidRequest.Wrapf("invalid validator cap value should be 0 <= cap <= 1")
			}
		case KeyLSMBondFactor:
			bondFactor, err := sdk.NewDecFromStr(update.Value)
			if err != nil {
				return fmt.Errorf("unable to parse string to sdk.Dec: %w", err)
			}

			// -1 is the default bond factor value
			if bondFactor.LT(sdk.ZeroDec()) && !bondFactor.Equal(sdk.NewDec(-1)) {
				return sdkerrors.ErrInvalidRequest.Wrapf("invalid validator bond factor value should be bond_factor == -1 || bond_factor >= 0")
			}
		case KeyMaxEntries:
			entries, err := strconv.ParseUint(update.Value, 10, 32)
			if err != nil {
				return err
			}
			if entries <= 0 {
				return fmt.Errorf("max entries undelegation/redelegation cannot be zero or lesser, found %v", entries)
			}
		case KeyRedelegationAcceptableDelta:
			redelegationAcceptableDelta, ok := sdk.NewIntFromString(update.Value)
			if !ok {
				return fmt.Errorf("unable to parse redeleagtion acceptable delta string %v to sdk.Int", update.Value)
			}
			if redelegationAcceptableDelta.LTE(math.ZeroInt()) {
				return fmt.Errorf("acceptable skew in validator delegations cannot be less that equal to zero, found %v", redelegationAcceptableDelta.String())
			}
		case KeyMinimumDeposit:
			minimumDeposit, ok := sdk.NewIntFromString(update.Value)
			if !ok {
				return sdkerrors.ErrInvalidRequest.Wrapf("")
			}

			if minimumDeposit.LTE(sdk.ZeroInt()) {
				return fmt.Errorf("invalid minimum deposit value less or equal than zero")
			}
		case KeyActive:
			_, err := strconv.ParseBool(update.Value)
			if err != nil {
				return fmt.Errorf("unable to parse string to bool")
			}
		case KeySetWithdrawAddress:
			if update.Value != "" {
				return fmt.Errorf("expected value for key:SetWithdrawAddress is empty")
			}
		case KeyAutocompoundFactor:
			autocompoundFactor, err := sdk.NewDecFromStr(update.Value)
			if err != nil {
				return fmt.Errorf("unable to parse string to sdk.Dec")
			}

			if autocompoundFactor.LTE(sdk.ZeroDec()) {
				return fmt.Errorf("invalid autocompound factor value less or equal than zero")
			}
		case KeyFlags:
			var flags HostChainFlags
			err := json.Unmarshal([]byte(update.Value), &flags)
			if err != nil {
				return fmt.Errorf("unable to unmarshal flags update string")
			}
		case KeyRewardParams:
			var params RewardParams
			err := json.Unmarshal([]byte(update.Value), &params)
			if err != nil {
				return fmt.Errorf("unable to unmarshal reward params update string")
			}

			if err := sdk.ValidateDenom(params.Denom); err != nil {
				return fmt.Errorf("invalid rewards denom: %s", err.Error())
			}
		default:
			return fmt.Errorf("invalid or unexpected update key: %s", update.Key)
		}
	}

	return nil
}

func NewMsgLiquidStake(amount sdk.Coin, address sdk.AccAddress) *MsgLiquidStake {
	return &MsgLiquidStake{
		DelegatorAddress: address.String(),
		Amount:           amount,
	}
}

func (m *MsgLiquidStake) Route() string {
	return RouterKey
}

// Type should return the action
func (m *MsgLiquidStake) Type() string {
	return MsgTypeLiquidStake
}

// GetSignBytes encodes the message for signing
func (m *MsgLiquidStake) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgLiquidStake) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// ValidateBasic performs stateless checks
func (m *MsgLiquidStake) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.DelegatorAddress); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, m.DelegatorAddress)
	}

	if !m.Amount.IsValid() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, m.Amount.String())
	}

	if !m.Amount.IsPositive() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, m.Amount.String())
	}

	return ibctransfertypes.ValidateIBCDenom(m.Amount.Denom)
}

func NewMsgLiquidStakeLSM(delegations sdk.Coins, address sdk.AccAddress) *MsgLiquidStakeLSM {
	return &MsgLiquidStakeLSM{
		DelegatorAddress: address.String(),
		Delegations:      delegations,
	}
}

func (m *MsgLiquidStakeLSM) Route() string {
	return RouterKey
}

// Type should return the action
func (m *MsgLiquidStakeLSM) Type() string {
	return MsgTypeLiquidStakeLSM
}

// GetSignBytes encodes the message for signing
func (m *MsgLiquidStakeLSM) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgLiquidStakeLSM) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// ValidateBasic performs stateless checks
func (m *MsgLiquidStakeLSM) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.DelegatorAddress); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, m.DelegatorAddress)
	}

	if !m.Delegations.IsValid() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, m.Delegations.String())
	}
	if !m.Delegations.IsAllPositive() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, m.Delegations.String())
	}

	for _, delegation := range m.Delegations {
		if err := ibctransfertypes.ValidateIBCDenom(delegation.Denom); err != nil {
			return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, delegation.Amount.String())
		}
	}

	return nil
}

func NewMsgLiquidUnstake(amount sdk.Coin, address sdk.AccAddress) *MsgLiquidUnstake {
	return &MsgLiquidUnstake{
		DelegatorAddress: address.String(),
		Amount:           amount,
	}
}

func (m *MsgLiquidUnstake) Route() string {
	return RouterKey
}

// Type should return the action
func (m *MsgLiquidUnstake) Type() string {
	return MsgTypeLiquidUnstake
}

// GetSignBytes encodes the message for signing
func (m *MsgLiquidUnstake) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgLiquidUnstake) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// ValidateBasic performs stateless checks
func (m *MsgLiquidUnstake) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.DelegatorAddress); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, m.DelegatorAddress)
	}

	if !m.Amount.IsValid() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, m.Amount.String())
	}

	if !m.Amount.IsPositive() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, m.Amount.String())
	}

	if !IsLiquidStakingDenom(m.Amount.Denom) {
		return sdkerrors.ErrInvalidCoins.Wrapf("invalid denom, required stk/{host-denom} got %s", m.Amount.Denom)
	}

	return nil
}

func NewMsgRedeem(amount sdk.Coin, address sdk.AccAddress) *MsgRedeem {
	return &MsgRedeem{
		DelegatorAddress: address.String(),
		Amount:           amount,
	}
}

func (m *MsgRedeem) Route() string {
	return RouterKey
}

// Type should return the action
func (m *MsgRedeem) Type() string {
	return MsgTypeRedeem
}

// GetSignBytes encodes the message for signing
func (m *MsgRedeem) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgRedeem) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// ValidateBasic performs stateless checks
func (m *MsgRedeem) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.DelegatorAddress); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, m.DelegatorAddress)
	}

	if !m.Amount.IsValid() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, m.Amount.String())
	}

	if !m.Amount.IsPositive() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, m.Amount.String())
	}

	if !IsLiquidStakingDenom(m.Amount.Denom) {
		return sdkerrors.ErrInvalidCoins.Wrapf("invalid denom, required stk/{host-denom} got %s", m.Amount.Denom)
	}
	return nil
}

func NewMsgUpdateParams(authority sdk.AccAddress, amount Params) *MsgUpdateParams {
	return &MsgUpdateParams{
		Authority: authority.String(),
		Params:    amount,
	}
}

func (m *MsgUpdateParams) Route() string {
	return RouterKey
}

// Type should return the action
func (m *MsgUpdateParams) Type() string {
	return MsgTypeUpdateParams
}

// GetSignBytes encodes the message for signing
func (m *MsgUpdateParams) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (m *MsgUpdateParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address %q: %v", m.Authority, err)
	}

	err := m.Params.Validate()
	if err != nil {
		return err
	}
	return nil
}
