package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

var (
	addr1             = authtypes.NewModuleAddressOrBech32Address("test1")
	ibcDenom          = "ibc/C4CFF46FD6DE35CA4CF4CE031E643C8FDC9BA4B99AE598E9B0ED98FE3A2319F9"
	amount1           = sdk.NewInt64Coin(ibcDenom, 10000)
	liquidstakedDenom = "stk/uatom"
	stkAmount1        = sdk.NewInt64Coin(liquidstakedDenom, 10000)
)

func TestMsgLiquidStake(t *testing.T) {
	msgLiquidStake := &types.MsgLiquidStake{
		DelegatorAddress: addr1.String(),
		Amount:           amount1,
	}
	newMsgLiquidStake := types.NewMsgLiquidStake(amount1, addr1)
	require.Equal(t, msgLiquidStake, newMsgLiquidStake)
	require.Equal(t, types.ModuleName, msgLiquidStake.Route())
	require.Equal(t, types.MsgTypeLiquidStake, msgLiquidStake.Type())
	require.Equal(t, addr1, msgLiquidStake.GetSigners()[0])
	require.NotPanics(t, func() { msgLiquidStake.GetSignBytes() })

	require.Equal(t, nil, msgLiquidStake.ValidateBasic())

	invalidCoin := amount1
	invalidCoin.Amount = sdk.NewInt(-1)
	invalidCoinMsg := types.NewMsgLiquidStake(invalidCoin, addr1)
	require.Error(t, invalidCoinMsg.ValidateBasic())

	zeroCoinMsg := types.NewMsgLiquidStake(sdk.NewCoin(ibcDenom, sdk.ZeroInt()), addr1)
	require.Error(t, zeroCoinMsg.ValidateBasic())

	invalidAddrMsg := types.NewMsgLiquidStake(amount1, sdk.AccAddress("test"))
	require.Error(t, invalidAddrMsg.ValidateBasic())
	require.Panics(t, func() { invalidAddrMsg.GetSigners() })
}
func TestMsgLiquidUnstake(t *testing.T) {
	msgLiquidUnstake := &types.MsgLiquidUnstake{
		DelegatorAddress: addr1.String(),
		Amount:           stkAmount1,
	}
	newMsgLiquidUnstake := types.NewMsgLiquidUnstake(stkAmount1, addr1)
	require.Equal(t, msgLiquidUnstake, newMsgLiquidUnstake)
	require.Equal(t, types.ModuleName, msgLiquidUnstake.Route())
	require.Equal(t, types.MsgTypeLiquidUnstake, msgLiquidUnstake.Type())
	require.Equal(t, addr1, msgLiquidUnstake.GetSigners()[0])
	require.NotPanics(t, func() { msgLiquidUnstake.GetSignBytes() })

	require.Equal(t, nil, msgLiquidUnstake.ValidateBasic())

	invalidCoin := stkAmount1
	invalidCoin.Amount = sdk.NewInt(-10)
	invalidCoinMsg := types.NewMsgLiquidUnstake(invalidCoin, addr1)
	require.Error(t, invalidCoinMsg.ValidateBasic())

	zeroCoinMsg := types.NewMsgLiquidUnstake(sdk.NewCoin(liquidstakedDenom, sdk.ZeroInt()), addr1)
	require.Error(t, zeroCoinMsg.ValidateBasic())

	invalidDenomMsg := types.NewMsgLiquidUnstake(amount1, addr1)
	require.Error(t, invalidDenomMsg.ValidateBasic())

	invalidAddrMsg := types.NewMsgLiquidUnstake(stkAmount1, sdk.AccAddress("test"))
	require.Error(t, invalidAddrMsg.ValidateBasic())
	require.Panics(t, func() { invalidAddrMsg.GetSigners() })
}
func TestMsgRedeem(t *testing.T) {
	msgRedeem := &types.MsgRedeem{
		DelegatorAddress: addr1.String(),
		Amount:           stkAmount1,
	}
	newMsgRedeem := types.NewMsgRedeem(stkAmount1, addr1)
	require.Equal(t, msgRedeem, newMsgRedeem)
	require.Equal(t, types.ModuleName, msgRedeem.Route())
	require.Equal(t, types.MsgTypeRedeem, msgRedeem.Type())
	require.Equal(t, addr1, msgRedeem.GetSigners()[0])
	require.NotPanics(t, func() { msgRedeem.GetSignBytes() })

	require.Equal(t, nil, msgRedeem.ValidateBasic())

	invalidCoin := stkAmount1
	invalidCoin.Amount = sdk.NewInt(-10)
	invalidCoinMsg := types.NewMsgRedeem(invalidCoin, addr1)
	require.Error(t, invalidCoinMsg.ValidateBasic())

	zeroCoinMsg := types.NewMsgRedeem(sdk.NewCoin(liquidstakedDenom, sdk.ZeroInt()), addr1)
	require.Error(t, zeroCoinMsg.ValidateBasic())

	invalidDenomMsg := types.NewMsgRedeem(amount1, addr1)
	require.Error(t, invalidDenomMsg.ValidateBasic())

	invalidAddrMsg := types.NewMsgRedeem(stkAmount1, sdk.AccAddress("test"))
	require.Error(t, invalidAddrMsg.ValidateBasic())
	require.Panics(t, func() { invalidAddrMsg.GetSigners() })
}

func TestMsgRegisterHostChain(t *testing.T) {
	msgRegisterHostChain := &types.MsgRegisterHostChain{
		Authority:          addr1.String(),
		ConnectionId:       "connection-localhost",
		DepositFee:         sdk.ZeroDec(),
		RestakeFee:         sdk.ZeroDec(),
		UnstakeFee:         sdk.ZeroDec(),
		RedemptionFee:      sdk.ZeroDec(),
		ChannelId:          "channel-1",
		PortId:             "transfer",
		HostDenom:          "uatom",
		MinimumDeposit:     sdk.OneInt(),
		UnbondingFactor:    4,
		AutoCompoundFactor: 2,
	}
	newMsgRegisterHostChain := types.NewMsgRegisterHostChain("connection-localhost", "channel-1", "transfer",
		"0", "0", "0", "0", "uatom", sdk.OneInt(), 4,
		addr1.String(), 2)
	require.Equal(t, msgRegisterHostChain, newMsgRegisterHostChain)
	require.Equal(t, types.ModuleName, msgRegisterHostChain.Route())
	require.Equal(t, types.MsgTypeRegisterHostChain, msgRegisterHostChain.Type())
	require.Equal(t, addr1, msgRegisterHostChain.GetSigners()[0])
	require.NotPanics(t, func() { msgRegisterHostChain.GetSignBytes() })

	require.Equal(t, nil, msgRegisterHostChain.ValidateBasic())

	invalidAddrMsg := *msgRegisterHostChain
	invalidAddrMsg.Authority = "test"
	require.Error(t, invalidAddrMsg.ValidateBasic())
	require.Panics(t, func() { invalidAddrMsg.GetSigners() })

	invalidMsg := *msgRegisterHostChain
	invalidMsg.ConnectionId = ""
	require.Error(t, invalidMsg.ValidateBasic())

	invalidMsg = *msgRegisterHostChain
	invalidMsg.ConnectionId = "notconnection-0"
	require.Error(t, invalidMsg.ValidateBasic())

	invalidMsg = *msgRegisterHostChain
	invalidMsg.HostDenom = "s" // small denom invalid
	require.Error(t, invalidMsg.ValidateBasic())

	invalidMsg = *msgRegisterHostChain
	invalidMsg.ChannelId = "notchannel-0"
	require.Error(t, invalidMsg.ValidateBasic())

	invalidMsg = *msgRegisterHostChain
	invalidMsg.RestakeFee = sdk.MustNewDecFromStr("-1")
	require.Error(t, invalidMsg.ValidateBasic())

	invalidMsg = *msgRegisterHostChain
	invalidMsg.DepositFee = sdk.MustNewDecFromStr("-1")
	require.Error(t, invalidMsg.ValidateBasic())

	invalidMsg = *msgRegisterHostChain
	invalidMsg.UnstakeFee = sdk.MustNewDecFromStr("-1")
	require.Error(t, invalidMsg.ValidateBasic())

	invalidMsg = *msgRegisterHostChain
	invalidMsg.RedemptionFee = sdk.MustNewDecFromStr("-1")
	require.Error(t, invalidMsg.ValidateBasic())
}

func TestMsgUpdateHostChain(t *testing.T) {
	msgUpdateHostChain := &types.MsgUpdateHostChain{
		Authority: addr1.String(),
		ChainId:   "chain-1",
		Updates: []*types.KVUpdate{
			{
				Key:   "add_val",
				Value: "cosmos1someval",
			},
		},
	}
	newMsgUpdateHostChain := types.NewMsgUpdateHostChain("chain-1", addr1.String(), []*types.KVUpdate{{
		Key:   "add_val",
		Value: "cosmos1someval",
	}})
	require.Equal(t, msgUpdateHostChain, newMsgUpdateHostChain)
	require.Equal(t, types.ModuleName, msgUpdateHostChain.Route())
	require.Equal(t, types.MsgTypeUpdateHostChain, msgUpdateHostChain.Type())
	require.Equal(t, addr1, msgUpdateHostChain.GetSigners()[0])
	require.NotPanics(t, func() { msgUpdateHostChain.GetSignBytes() })

	require.Equal(t, nil, msgUpdateHostChain.ValidateBasic())

	invalidAddrMsg := *msgUpdateHostChain
	invalidAddrMsg.Authority = "test"
	require.Error(t, invalidAddrMsg.ValidateBasic())
	require.Panics(t, func() { invalidAddrMsg.GetSigners() })
}

func TestMsgUpdateParams(t *testing.T) {
	msgUpdateParams := &types.MsgUpdateParams{
		Authority: addr1.String(),
		Params:    types.DefaultParams(),
	}
	newMsgUpdateParams := types.NewMsgUpdateParams(addr1, types.DefaultParams())
	require.Equal(t, msgUpdateParams, newMsgUpdateParams)
	require.Equal(t, types.ModuleName, msgUpdateParams.Route())
	require.Equal(t, types.MsgTypeUpdateParams, msgUpdateParams.Type())
	require.Equal(t, addr1, msgUpdateParams.GetSigners()[0])
	require.NotPanics(t, func() { msgUpdateParams.GetSignBytes() })

	require.Equal(t, nil, msgUpdateParams.ValidateBasic())

	invalidAddrMsg := types.NewMsgUpdateParams(sdk.AccAddress("test"), types.Params{
		AdminAddress:     addr1.String(),
		FeeAddress:       addr1.String(),
		UpperCValueLimit: sdk.OneDec(),
		LowerCValueLimit: sdk.ZeroDec(),
	})
	require.Error(t, invalidAddrMsg.ValidateBasic())
	require.Panics(t, func() { invalidAddrMsg.GetSigners() })

	invalidParamsMsg := *msgUpdateParams
	invalidParamsMsg.Params.AdminAddress = "test"
	require.Error(t, invalidParamsMsg.ValidateBasic())

}
