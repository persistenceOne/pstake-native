package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
)

var (
	_ sdk.Msg = &MsgLiquidStake{}
	_ sdk.Msg = &MsgJuice{}
	_ sdk.Msg = &MsgLiquidUnstake{}
	_ sdk.Msg = &MsgRedeem{}
	_ sdk.Msg = &MsgClaim{}
	_ sdk.Msg = &MsgJumpStart{}
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

// NewMsgLiquidUnstake returns a new MsgLiquidUnstake
//
//nolint:interfacer
func NewMsgLiquidUnstake(address sdk.AccAddress, amount sdk.Coin) *MsgLiquidUnstake {
	return &MsgLiquidUnstake{
		DelegatorAddress: address.String(),
		Amount:           amount,
	}
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

// NewMsgRedeem returns a new MsgRedeem
//
//nolint:interfacer
func NewMsgRedeem(address sdk.AccAddress, amount sdk.Coin) *MsgRedeem {
	return &MsgRedeem{
		DelegatorAddress: address.String(),
		Amount:           amount,
	}
}

func (m *MsgRedeem) Route() string { return RouterKey }

func (m *MsgRedeem) Type() string { return MsgTypeRedeem }

func (m *MsgRedeem) ValidateBasic() error {
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

func (m *MsgRedeem) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))

}

func (m *MsgRedeem) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// NewMsgClaim returns a new MsgClaim
//
//nolint:interfacer
func NewMsgClaim(address sdk.AccAddress) *MsgClaim {
	return &MsgClaim{
		DelegatorAddress: address.String(),
	}
}

// Route should return the name of the module
func (m *MsgClaim) Route() string { return RouterKey }

// Type should return the action
func (m *MsgClaim) Type() string { return MsgTypeClaim }

// ValidateBasic performs stateless checks
func (m *MsgClaim) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.DelegatorAddress); err != nil {
		return sdkErrors.Wrap(sdkErrors.ErrInvalidAddress, m.DelegatorAddress)
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (m *MsgClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgClaim) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// NewMsgLiquidStake returns a new MsgLiquidStake
//
//nolint:interfacer
func NewMsgJumpStart(address sdk.AccAddress, chainID, connectionID, transferChannel, transferPort, baseDenom, mintDenom string,
	minDeposit sdk.Int, allowListedValidators AllowListedValidators, pstakeParams PstakeParams, hostAccounts HostAccounts) *MsgJumpStart {
	return &MsgJumpStart{
		PstakeAddress:         address.String(),
		ChainID:               chainID,
		ConnectionID:          connectionID,
		TransferChannel:       transferChannel,
		TransferPort:          transferPort,
		BaseDenom:             baseDenom,
		MintDenom:             mintDenom,
		MinDeposit:            minDeposit,
		AllowListedValidators: allowListedValidators,
		PstakeParams:          pstakeParams,
		HostAccounts:          hostAccounts,
	}
}

// Route should return the name of the module
func (m *MsgJumpStart) Route() string { return RouterKey }

// Type should return the action
func (m *MsgJumpStart) Type() string { return MsgTypeJumpStart }

// ValidateBasic performs stateless checks
func (m *MsgJumpStart) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.PstakeAddress); err != nil {
		return sdkErrors.Wrap(sdkErrors.ErrInvalidAddress, m.PstakeAddress)
	}
	if m.ChainID == "" ||
		m.ConnectionID == "" ||
		m.TransferChannel == "" ||
		m.TransferPort == "" ||
		m.BaseDenom == "" ||
		m.MintDenom == "" {
		return sdkErrors.Wrap(sdkErrors.ErrInvalidRequest, "params cannot be empty")
	}

	if !m.AllowListedValidators.Valid() {
		return ErrInValidAllowListedValidators
	}
	if _, err := sdk.AccAddressFromBech32(m.PstakeParams.PstakeFeeAddress); err != nil {
		return sdkErrors.Wrap(sdkErrors.ErrInvalidAddress, m.PstakeParams.PstakeFeeAddress)
	}
	if m.PstakeParams.PstakeFeeAddress != m.PstakeAddress {
		return sdkErrors.Wrap(sdkErrors.ErrInvalidAddress, fmt.Sprintf("pstakeAddress should be equal to PstakeParams.PstakeFeeAddress, got %s, %s", m.PstakeParams.PstakeFeeAddress, m.PstakeAddress))
	}
	return m.HostAccounts.Validate()
}

// GetSignBytes encodes the message for signing
func (m *MsgJumpStart) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgJumpStart) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.PstakeAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}
