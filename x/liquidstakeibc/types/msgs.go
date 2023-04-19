package types

import (
	"fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
)

const (
	MsgTypeRegisterHostChain string = "msg_register_host_chain"
	MsgTypeUpdateHostChain   string = "msg_update_host_chain"
	MsgTypeLiquidStake       string = "msg_liquid_stake"
)

var (
	_ sdk.Msg = &MsgRegisterHostChain{}
	_ sdk.Msg = &MsgUpdateHostChain{}
	_ sdk.Msg = &MsgLiquidStake{}
)

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

func NewMsgUpdateHostChain(chainId string, updates []*KVUpdate) *MsgUpdateHostChain {
	return &MsgUpdateHostChain{
		ChainId: chainId,
		Updates: updates,
	}
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

func NewMsgLiquidStake(amount sdk.Coin, address sdk.AccAddress) *MsgLiquidStake {
	return &MsgLiquidStake{
		DelegatorAddress: address.String(),
		Amount:           amount,
	}
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
	if _, err := sdk.AccAddressFromBech32(m.DelegatorAddress); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, m.DelegatorAddress)
	}

	if !m.Amount.IsValid() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, m.Amount.String())
	}

	if !m.Amount.IsPositive() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, m.Amount.String())
	}

	return ibctransfertypes.ValidateIBCDenom(m.Amount.Denom)
}
