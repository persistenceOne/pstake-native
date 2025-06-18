package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	liquidstakeibctypes "github.com/persistenceOne/pstake-native/v3/x/liquidstakeibc/types"
)

const TypeMsgUpdateParams = "msg_update_params"
const (
	TypeMsgCreateHostChain = "create_host_chain"
	TypeMsgUpdateHostChain = "update_host_chain"
	TypeMsgDeleteHostChain = "delete_host_chain"
)

var _ sdk.Msg = &MsgUpdateParams{}

func NewMsgUpdateParams(authority string, params Params) *MsgUpdateParams {
	return &MsgUpdateParams{
		Authority: authority,
		Params:    params,
	}
}

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
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	return msg.Params.Validate()
}

var _ sdk.Msg = &MsgCreateHostChain{}

func NewMsgCreateHostChain(
	authority string,
	hc HostChain,
) *MsgCreateHostChain {
	return &MsgCreateHostChain{
		Authority: authority,
		HostChain: hc,
	}
}

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
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	err = msg.HostChain.ValidateBasic()
	if err != nil {
		return err
	}

	if msg.HostChain.ID != 0 {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "hostchain ID for create msg should be 0")
	}
	if msg.HostChain.ICAAccount.Owner != "" {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "owner should not be specified as app uses default")
	}
	if msg.HostChain.ICAAccount.ChannelState != liquidstakeibctypes.ICAAccount_ICA_CHANNEL_CREATING {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "channel state should be creating")
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateHostChain{}

func NewMsgUpdateHostChain(
	creator string,
	hc HostChain,
) *MsgUpdateHostChain {
	return &MsgUpdateHostChain{
		Authority: creator,
		HostChain: hc,
	}
}

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
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	err = msg.HostChain.ValidateBasic()
	if err != nil {
		return err
	}

	if msg.HostChain.ID == 0 {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "hostchain ID for update msg should not be 0")
	}

	return nil
}

var _ sdk.Msg = &MsgDeleteHostChain{}

func NewMsgDeleteHostChain(
	creator string,
	id uint64,
) *MsgDeleteHostChain {
	return &MsgDeleteHostChain{
		Authority: creator,
		ID:        id,
	}
}

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
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.ID == 0 {
		return errors.Wrapf(sdkerrors.ErrInvalidRequest, "hostchain ID for delete msg should not be 0")
	}

	return nil
}
