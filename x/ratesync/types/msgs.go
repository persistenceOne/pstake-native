package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgUpdateParams = "msg_update_params"
const (
	TypeMsgCreateHostChain = "create_host_chain"
	TypeMsgUpdateHostChain = "update_host_chain"
	TypeMsgDeleteHostChain = "delete_host_chain"
)

var _ sdk.Msg = &MsgUpdateParams{}

func (msg *MsgUpdateParams) Route() string {
	return RouterKey
}

func (msg *MsgUpdateParams) Type() string {
	return TypeMsgUpdateParams
}

func (msg *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgUpdateParams) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateParams) ValidateBasic() error {
	return nil
}

var _ sdk.Msg = &MsgCreateHostChain{}

func (msg *MsgCreateHostChain) Route() string {
	return RouterKey
}

func (msg *MsgCreateHostChain) Type() string {
	return TypeMsgCreateHostChain
}

func (msg *MsgCreateHostChain) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateHostChain) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateHostChain) ValidateBasic() error {
	return nil
}

var _ sdk.Msg = &MsgUpdateHostChain{}

func (msg *MsgUpdateHostChain) Route() string {
	return RouterKey
}

func (msg *MsgUpdateHostChain) Type() string {
	return TypeMsgUpdateHostChain
}

func (msg *MsgUpdateHostChain) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateHostChain) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateHostChain) ValidateBasic() error {
	return nil
}

var _ sdk.Msg = &MsgDeleteHostChain{}

func (msg *MsgDeleteHostChain) Route() string {
	return RouterKey
}

func (msg *MsgDeleteHostChain) Type() string {
	return TypeMsgDeleteHostChain
}

func (msg *MsgDeleteHostChain) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeleteHostChain) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteHostChain) ValidateBasic() error {
	return nil
}
