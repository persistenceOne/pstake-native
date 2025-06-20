package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	_ sdk.Msg = &MsgUpdateParams{}
)

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
	return nil
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
	return nil
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
	return nil
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
	return nil
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
	return nil
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

func (m *MsgUpdateParams) ValidateBasic() error { return nil }
