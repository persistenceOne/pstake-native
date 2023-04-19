package types

import (
	"fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgRegisterHostChain{}

func NewMsgRegisterHostChain(
	connectionId string,
	hostDenom string,
	localDenom string,
	minimumDeposit math.Int,
) *MsgRegisterHostChain {

	return &MsgRegisterHostChain{
		ConnectionId:   connectionId,
		HostDenom:      hostDenom,
		LocalDenom:     localDenom,
		MinimumDeposit: minimumDeposit,
	}
}

func (m *MsgRegisterHostChain) Route() string {
	return sdk.MsgTypeURL(m)
}

func (m *MsgRegisterHostChain) Type() string {
	return sdk.MsgTypeURL(m)
}

func (m *MsgRegisterHostChain) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgRegisterHostChain) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Authority)
	return []sdk.AccAddress{addr}
}

func (m *MsgRegisterHostChain) ValidateBasic() error {
	// connection id cannot be empty and must begin with "connection"
	if m.ConnectionId == "" {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "connection id cannot be empty")
	}
	if !strings.HasPrefix(m.ConnectionId, "connection") {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "connection id must begin with 'connection'")
	}

	// validate host denom
	if err := sdk.ValidateDenom(m.HostDenom); err != nil {
		return errorsmod.Wrapf(
			sdkerrors.ErrInvalidRequest,
			fmt.Sprintf("invalid host denom: %s", err.Error()),
		)
	}

	// validate local denom
	if err := sdk.ValidateDenom(m.LocalDenom); err != nil {
		return errorsmod.Wrapf(
			sdkerrors.ErrInvalidRequest,
			fmt.Sprintf("invalid local denom: %s", err.Error()),
		)
	}

	// minimum deposit must be positive or zero
	if m.MinimumDeposit.LT(sdk.NewInt(0)) {
		return errorsmod.Wrapf(
			sdkerrors.ErrInvalidRequest,
			"minimum deposit quantity must be greater or equal than zero",
		)
	}

	return nil
}

func NewMsgUpdateHostChain(
	chainId string,
	updates []*KVUpdate,
) *MsgUpdateHostChain {

	return &MsgUpdateHostChain{
		ChainId: chainId,
		Updates: updates,
	}
}

func (m *MsgUpdateHostChain) Route() string {
	return sdk.MsgTypeURL(m)
}

func (m *MsgUpdateHostChain) Type() string {
	return sdk.MsgTypeURL(m)
}

func (m *MsgUpdateHostChain) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgUpdateHostChain) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Authority)
	return []sdk.AccAddress{addr}
}

func (m *MsgUpdateHostChain) ValidateBasic() error {
	return nil
}
