package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgDummy{}

func NewMsgDummy(fromAddress sdk.AccAddress) *MsgDummy {
	return &MsgDummy{FromAddress: fromAddress.String()}
}

// Route Implements Msg.
func (m MsgDummy) Route() string { return sdk.MsgTypeURL(&m) }

// Type Implements Msg.
func (m MsgDummy) Type() string { return sdk.MsgTypeURL(&m) }

// GetSignBytes Implements Msg.
func (m MsgDummy) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

// GetSigners returns the expected signers for a MsgDummy.
func (m MsgDummy) GetSigners() []sdk.AccAddress {
	fromAddress := sdk.MustAccAddressFromBech32(m.FromAddress)

	return []sdk.AccAddress{fromAddress}
}

// ValidateBasic does a sanity check on the provided data
func (m MsgDummy) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.FromAddress)
	if err != nil {
		return sdkerrors.Wrapf(err, "invalid from address: %s", m.FromAddress)
	}
	return nil
}
