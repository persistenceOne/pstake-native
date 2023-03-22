package types

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
)

var (
	_ sdk.Msg = &MsgLiquidStake{}
	_ sdk.Msg = &MsgLiquidUnstake{}
	_ sdk.Msg = &MsgRedeem{}
	_ sdk.Msg = &MsgClaim{}
	_ sdk.Msg = &MsgRecreateICA{}
	_ sdk.Msg = &MsgJumpStart{}
	_ sdk.Msg = &MsgChangeModuleState{}
	_ sdk.Msg = &MsgReportSlashing{}
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
		return errorsmod.Wrap(sdkErrors.ErrInvalidAddress, m.DelegatorAddress)
	}

	if !m.Amount.IsValid() {
		return errorsmod.Wrap(sdkErrors.ErrInvalidCoins, m.Amount.String())
	}

	if !m.Amount.IsPositive() {
		return errorsmod.Wrap(sdkErrors.ErrInvalidCoins, m.Amount.String())
	}
	return ibctransfertypes.ValidateIBCDenom(m.Amount.Denom)
}

// GetSignBytes encodes the message for signing
func (m *MsgLiquidStake) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgLiquidStake) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.DelegatorAddress)
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

// Route should return the name of the module
func (m *MsgLiquidUnstake) Route() string { return RouterKey }

// Type should return the action
func (m *MsgLiquidUnstake) Type() string { return MsgTypeLiquidUnstake }

// ValidateBasic performs stateless checks
func (m *MsgLiquidUnstake) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.DelegatorAddress); err != nil {
		return errorsmod.Wrap(sdkErrors.ErrInvalidAddress, m.DelegatorAddress)
	}

	if !m.Amount.IsValid() {
		return errorsmod.Wrap(sdkErrors.ErrInvalidCoins, m.Amount.String())
	}

	if !m.Amount.IsPositive() {
		return errorsmod.Wrap(sdkErrors.ErrInvalidCoins, m.Amount.String())
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m *MsgLiquidUnstake) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(m))

}

// GetSigners defines whose signature is required
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

// Route should return the name of the module
func (m *MsgRedeem) Route() string { return RouterKey }

// Type should return the action
func (m *MsgRedeem) Type() string { return MsgTypeRedeem }

// ValidateBasic performs stateless checks
func (m *MsgRedeem) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.DelegatorAddress); err != nil {
		return errorsmod.Wrap(sdkErrors.ErrInvalidAddress, m.DelegatorAddress)
	}

	if !m.Amount.IsValid() {
		return errorsmod.Wrap(sdkErrors.ErrInvalidCoins, m.Amount.String())
	}

	if !m.Amount.IsPositive() {
		return errorsmod.Wrap(sdkErrors.ErrInvalidCoins, m.Amount.String())
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m *MsgRedeem) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(m))

}

// GetSigners defines whose signature is required
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
		return errorsmod.Wrap(sdkErrors.ErrInvalidAddress, m.DelegatorAddress)
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (m *MsgClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgClaim) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// NewMsgRecreateICA returns a new MsgRecreateICA
//
//nolint:interfacer
func NewMsgRecreateICA(address sdk.AccAddress) *MsgRecreateICA {
	return &MsgRecreateICA{
		FromAddress: address.String(),
	}
}

// Route should return the name of the module
func (m *MsgRecreateICA) Route() string { return RouterKey }

// Type should return the action
func (m *MsgRecreateICA) Type() string { return MsgTypeRecreateICA }

// ValidateBasic performs stateless checks
func (m *MsgRecreateICA) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.FromAddress); err != nil {
		return errorsmod.Wrap(sdkErrors.ErrInvalidAddress, m.FromAddress)
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (m *MsgRecreateICA) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgRecreateICA) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// NewMsgJumpStart returns a new MsgJumpStart
//
//nolint:interfacer
func NewMsgJumpStart(address sdk.AccAddress, chainID, connectionID, transferChannel, transferPort, baseDenom, mintDenom string,
	minDeposit math.Int, allowListedValidators AllowListedValidators, pstakeParams PstakeParams, hostAccounts HostAccounts) *MsgJumpStart {
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
		return errorsmod.Wrap(sdkErrors.ErrInvalidAddress, m.PstakeAddress)
	}
	if m.ChainID == "" ||
		m.ConnectionID == "" ||
		m.TransferChannel == "" ||
		m.TransferPort == "" ||
		m.BaseDenom == "" ||
		m.MintDenom == "" {
		return errorsmod.Wrap(sdkErrors.ErrInvalidRequest, "params cannot be empty")
	}

	if !m.AllowListedValidators.Valid() {
		return ErrInValidAllowListedValidators
	}
	if _, err := sdk.AccAddressFromBech32(m.PstakeParams.PstakeFeeAddress); err != nil {
		return errorsmod.Wrap(sdkErrors.ErrInvalidAddress, m.PstakeParams.PstakeFeeAddress)
	}
	if m.PstakeParams.PstakeFeeAddress != m.PstakeAddress {
		return errorsmod.Wrap(sdkErrors.ErrInvalidAddress, fmt.Sprintf("pstakeAddress should be equal to PstakeParams.PstakeFeeAddress, got %s, %s", m.PstakeParams.PstakeFeeAddress, m.PstakeAddress))
	}
	if m.MinDeposit.LTE(sdk.ZeroInt()) {
		return errorsmod.Wrapf(ErrInvalidDeposit, "min deposit must be positive")
	}
	if ConvertBaseDenomToMintDenom(m.BaseDenom) != m.MintDenom {
		return ErrInvalidMintDenom
	}
	return m.HostAccounts.Validate()
}

// GetSignBytes encodes the message for signing
func (m *MsgJumpStart) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgJumpStart) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.PstakeAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// NewMsgChangeModuleState returns a new MsgChangeModuleState
//
//nolint:interfacer
func NewMsgChangeModuleState(address sdk.AccAddress, moduleState bool) *MsgChangeModuleState {
	return &MsgChangeModuleState{
		PstakeAddress: address.String(),
		ModuleState:   moduleState,
	}
}

// Route should return the name of the module
func (m *MsgChangeModuleState) Route() string { return RouterKey }

// Type should return the action
func (m *MsgChangeModuleState) Type() string { return MsgTypeChangeModuleState }

// ValidateBasic performs stateless checks
func (m *MsgChangeModuleState) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.PstakeAddress); err != nil {
		return errorsmod.Wrap(sdkErrors.ErrInvalidAddress, m.PstakeAddress)
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (m *MsgChangeModuleState) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgChangeModuleState) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.PstakeAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

// NewMsgReportSlashing returns a new MsgReportSlashing
//
//nolint:interfacer
func NewMsgReportSlashing(address sdk.AccAddress, validatorAddress sdk.ValAddress) *MsgReportSlashing {
	valAddr, err := Bech32FromValAddress(validatorAddress, CosmosValOperPrefix)
	if err != nil {
		panic(err)
	}
	return &MsgReportSlashing{
		PstakeAddress:    address.String(),
		ValidatorAddress: valAddr,
	}
}

// Route should return the name of the module
func (m *MsgReportSlashing) Route() string { return RouterKey }

// Type should return the action
func (m *MsgReportSlashing) Type() string { return MsgTypeReportSlashing }

// ValidateBasic performs stateless checks
func (m *MsgReportSlashing) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.PstakeAddress); err != nil {
		return errorsmod.Wrap(sdkErrors.ErrInvalidAddress, m.PstakeAddress)
	}
	if _, err := ValAddressFromBech32(m.ValidatorAddress, CosmosValOperPrefix); err != nil {
		return errorsmod.Wrap(sdkErrors.ErrInvalidAddress, m.ValidatorAddress)
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (m *MsgReportSlashing) GetSignBytes() []byte {
	return sdk.MustSortJSON(AminoCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgReportSlashing) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(m.PstakeAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}
