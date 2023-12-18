package types

import (
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"
)

var ValidHostChainInMsg = HostChain{
	Id:           1,
	ChainId:      "test-1",
	ConnectionId: ibcexported.LocalhostConnectionID,
	IcaAccount:   types.ICAAccount{},
	Features: Feature{LiquidStakeIBC: LiquidStake{
		FeatureType:     0,
		CodeID:          0,
		Instantiation:   0,
		ContractAddress: "",
		Denoms:          []string{},
		Enabled:         false,
	}},
}

func TestMsgUpdateParams_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateParams
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdateParams{
				Authority: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUpdateParams{
				Authority: authtypes.NewModuleAddress("addr1").String(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMsgCreateHostChain_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCreateHostChain
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgCreateHostChain{
				Authority: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgCreateHostChain{
				Authority: authtypes.NewModuleAddress("addr1").String(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMsgUpdateHostChain_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateHostChain
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdateHostChain{
				Authority: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgUpdateHostChain{
				Authority: authtypes.NewModuleAddress("addr1").String(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMsgDeleteHostChain_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgDeleteHostChain
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgDeleteHostChain{
				Authority: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgDeleteHostChain{
				Authority: authtypes.NewModuleAddress("addr1").String(),
				Id:        1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
