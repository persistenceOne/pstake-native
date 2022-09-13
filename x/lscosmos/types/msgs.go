package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
)

var (
	_ sdk.Msg = &MsgLiquidStake{}
	_ sdk.Msg = &MsgJuice{}
	_ sdk.Msg = &MsgLiquidUnstake{}
)

// NewMsgLiquidStake returns a new MsgLiquidStake
//
//nolint:interfacer
func NewMsgLiquidStake(amount sdk.Coin, address sdk.AccAddress) *MsgLiquidStake {
	return &MsgLiquidStake{
		DelegatorAddress: address.String(),
		Amount:           amount,
	}
}

// Route should return the name of the module
func (m *MsgLiquidStake) Route() string { return RouterKey }

// Type should return the action
func (m *MsgLiquidStake) Type() string { return MsgTypeLiquidStake }

// ValidateBasic performs stateless checks
func (m *MsgLiquidStake) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.DelegatorAddress); err != nil {
		return sdkErrors.Wrap(sdkErrors.ErrInvalidAddress, m.DelegatorAddress)
	}

	if !m.Amount.IsValid() {
		return sdkErrors.Wrap(sdkErrors.ErrInvalidCoins, m.Amount.String())
	}

	if !m.Amount.IsPositive() {
		return sdkErrors.Wrap(sdkErrors.ErrInvalidCoins, m.Amount.String())
	}
	return ibctransfertypes.ValidateIBCDenom(m.Amount.Denom)
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

// NewMsgJuice returns a new MsgJuice
//
//nolint:interfacer
func NewMsgJuice(amount sdk.Coin, address sdk.AccAddress) *MsgJuice {
	return &MsgJuice{
		RewarderAddress: address.String(),
		Amount:          amount,
	}
}

// Route should return the name of the module
func (m *MsgJuice) Route() string { return RouterKey }

// Type should return the action
func (m *MsgJuice) Type() string { return MsgTypeJuice }

// ValidateBasic performs stateless checks
func (m *MsgJuice) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.RewarderAddress); err != nil {
		return sdkErrors.Wrap(sdkErrors.ErrInvalidAddress, m.RewarderAddress)
	}

	if !m.Amount.IsValid() {
		return sdkErrors.Wrap(sdkErrors.ErrInvalidCoins, m.Amount.String())
	}

	if !m.Amount.IsPositive() {
		return sdkErrors.Wrap(sdkErrors.ErrInvalidCoins, m.Amount.String())
	}
	return ibctransfertypes.ValidateIBCDenom(m.Amount.Denom)
}

// GetSignBytes encodes the message for signing
func (m *MsgJuice) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgJuice) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.RewarderAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

func (m *MsgLiquidUnstake) Route() string { return RouterKey }

func (m *MsgLiquidUnstake) Type() string { return MsgTypeLiquidUnstake }

func (m *MsgLiquidUnstake) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.DelegatorAddress); err != nil {
		return sdkErrors.Wrap(sdkErrors.ErrInvalidAddress, m.DelegatorAddress)
	}

	if !m.Amount.IsValid() {
		return sdkErrors.Wrap(sdkErrors.ErrInvalidCoins, m.Amount.String())
	}

	if !m.Amount.IsPositive() {
		return sdkErrors.Wrap(sdkErrors.ErrInvalidCoins, m.Amount.String())
	}
	return nil
}

func (m *MsgLiquidUnstake) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))

}

func (m *MsgLiquidUnstake) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}
