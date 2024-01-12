package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
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

// Route should return the name of the module
func (m *MsgLiquidStake) Route() string { return RouterKey }

// Type should return the action
func (m *MsgLiquidStake) Type() string { return MsgTypeLiquidStake }

// ValidateBasic performs stateless checks
func (m *MsgLiquidStake) ValidateBasic() error {
	return nil
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

// Route should return the name of the module
func (m *MsgLiquidUnstake) Route() string { return RouterKey }

// Type should return the action
func (m *MsgLiquidUnstake) Type() string { return MsgTypeLiquidUnstake }

// ValidateBasic performs stateless checks
func (m *MsgLiquidUnstake) ValidateBasic() error {
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

// Route should return the name of the module
func (m *MsgRedeem) Route() string { return RouterKey }

// Type should return the action
func (m *MsgRedeem) Type() string { return MsgTypeRedeem }

// ValidateBasic performs stateless checks
func (m *MsgRedeem) ValidateBasic() error {
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

// Route should return the name of the module
func (m *MsgClaim) Route() string { return RouterKey }

// Type should return the action
func (m *MsgClaim) Type() string { return MsgTypeClaim }

// ValidateBasic performs stateless checks
func (m *MsgClaim) ValidateBasic() error {
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

// Route should return the name of the module
func (m *MsgRecreateICA) Route() string { return RouterKey }

// Type should return the action
func (m *MsgRecreateICA) Type() string { return MsgTypeRecreateICA }

// ValidateBasic performs stateless checks
func (m *MsgRecreateICA) ValidateBasic() error {
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

// Route should return the name of the module
func (m *MsgJumpStart) Route() string { return RouterKey }

// Type should return the action
func (m *MsgJumpStart) Type() string { return MsgTypeJumpStart }

// ValidateBasic performs stateless checks
func (m *MsgJumpStart) ValidateBasic() error {
	return nil
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

// Route should return the name of the module
func (m *MsgChangeModuleState) Route() string { return RouterKey }

// Type should return the action
func (m *MsgChangeModuleState) Type() string { return MsgTypeChangeModuleState }

// ValidateBasic performs stateless checks
func (m *MsgChangeModuleState) ValidateBasic() error {
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

// Route should return the name of the module
func (m *MsgReportSlashing) Route() string { return RouterKey }

// Type should return the action
func (m *MsgReportSlashing) Type() string { return MsgTypeReportSlashing }

// ValidateBasic performs stateless checks
func (m *MsgReportSlashing) ValidateBasic() error {
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
