package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = (*MsgLiquidStake)(nil)
	_ sdk.Msg = (*MsgLiquidUnstake)(nil)
)

// Message types for the liquidstaking module
const (
	TypeMsgLiquidStake   = "liquid_stake"
	TypeMsgLiquidUnstake = "liquid_unstake"
)

// NewMsgLiquidStake creates a new MsgLiquidStake.
func NewMsgLiquidStake(
	liquidStaker sdk.AccAddress, //nolint: interfacer
	amount sdk.Coin,
) *MsgLiquidStake {
	return &MsgLiquidStake{
		DelegatorAddress: liquidStaker.String(),
		Amount:           amount,
	}
}

func (msg MsgLiquidStake) Route() string { return RouterKey }

func (msg MsgLiquidStake) Type() string { return TypeMsgLiquidStake }

func (msg MsgLiquidStake) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.DelegatorAddress); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid delegator address %q: %v", msg.DelegatorAddress, err)
	}
	if ok := msg.Amount.IsZero(); ok {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "staking amount must not be zero")
	}
	if err := msg.Amount.Validate(); err != nil {
		return err
	}
	return nil
}

func (msg MsgLiquidStake) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgLiquidStake) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgLiquidStake) GetDelegator() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return addr
}

// NewMsgLiquidUnstake creates a new MsgLiquidUnstake.
func NewMsgLiquidUnstake(
	liquidStaker sdk.AccAddress, //nolint: interfacer
	amount sdk.Coin,
) *MsgLiquidUnstake {
	return &MsgLiquidUnstake{
		DelegatorAddress: liquidStaker.String(),
		Amount:           amount,
	}
}

func (msg MsgLiquidUnstake) Route() string { return RouterKey }

func (msg MsgLiquidUnstake) Type() string { return TypeMsgLiquidUnstake }

func (msg MsgLiquidUnstake) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.DelegatorAddress); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid delegator address %q: %v", msg.DelegatorAddress, err)
	}
	if ok := msg.Amount.IsZero(); ok {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "unstaking amount must not be zero")
	}
	if err := msg.Amount.Validate(); err != nil {
		return err
	}
	return nil
}

func (msg MsgLiquidUnstake) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgLiquidUnstake) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgLiquidUnstake) GetDelegator() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return addr
}
