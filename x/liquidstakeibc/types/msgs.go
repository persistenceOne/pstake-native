package types

import (
	"fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	connectiontypes "github.com/cosmos/ibc-go/v7/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
)

const (
	MsgTypeRegisterHostChain string = "msg_register_host_chain"
	MsgTypeUpdateHostChain   string = "msg_update_host_chain"
	MsgTypeLiquidStake       string = "msg_liquid_stake"
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
) *MsgRegisterHostChain {
	depositFeeDec, _ := sdk.NewDecFromStr(depositFee)
	restakeFeeDec, _ := sdk.NewDecFromStr(restakeFee)
	unstakeFeeDec, _ := sdk.NewDecFromStr(unstakeFee)
	redemptionFeeDec, _ := sdk.NewDecFromStr(redemptionFee)

	return &MsgRegisterHostChain{
		ConnectionId:    connectionID,
		HostDenom:       hostDenom,
		ChannelId:       channelID,
		PortId:          portID,
		MinimumDeposit:  minimumDeposit,
		UnbondingFactor: unbondingFactor,
		DepositFee:      depositFeeDec,
		RestakeFee:      restakeFeeDec,
		UnstakeFee:      unstakeFeeDec,
		RedemptionFee:   redemptionFeeDec,
		Authority:       authority,
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
	if m.DepositFee.LT(sdk.NewDec(0)) {
		return errorsmod.Wrapf(
			sdkerrors.ErrInvalidRequest,
			"deposit fee quantity must be greater or equal than zero",
		)
	}

	// restake fee must be positive or zero
	if m.RestakeFee.LT(sdk.NewDec(0)) {
		return errorsmod.Wrapf(
			sdkerrors.ErrInvalidRequest,
			"restake fee quantity must be greater or equal than zero",
		)
	}

	// unstake fee must be positive or zero
	if m.UnstakeFee.LT(sdk.NewDec(0)) {
		return errorsmod.Wrapf(
			sdkerrors.ErrInvalidRequest,
			"unstake fee quantity must be greater or equal than zero",
		)
	}

	// redemption deposit must be positive or zero
	if m.RedemptionFee.LT(sdk.NewDec(0)) {
		return errorsmod.Wrapf(
			sdkerrors.ErrInvalidRequest,
			"redemption fee quantity must be greater or equal than zero",
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
	return nil
}

//nolint:interfacer
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

//nolint:interfacer
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

	return nil
}

//nolint:interfacer
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

	return nil
}

//nolint:interfacer
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
