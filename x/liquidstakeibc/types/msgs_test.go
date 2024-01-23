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

func TestMsgLiquidStakeLSM(t *testing.T) {
	msgLiquidStakeLSM := &types.MsgLiquidStakeLSM{
		DelegatorAddress: addr1.String(),
		Delegations:      sdk.NewCoins(sdk.NewCoin(ibcDenom, sdk.NewInt(10000))),
	}
	newMsgLiquidStake := types.NewMsgLiquidStakeLSM(sdk.NewCoins(sdk.NewCoin(ibcDenom, sdk.NewInt(10000))), addr1)
	require.Equal(t, msgLiquidStakeLSM, newMsgLiquidStake)
	require.Equal(t, types.ModuleName, msgLiquidStakeLSM.Route())
	require.Equal(t, types.MsgTypeLiquidStakeLSM, msgLiquidStakeLSM.Type())
	require.Equal(t, addr1, msgLiquidStakeLSM.GetSigners()[0])
	require.NotPanics(t, func() { msgLiquidStakeLSM.GetSignBytes() })

	require.Equal(t, nil, msgLiquidStakeLSM.ValidateBasic())

	invalidDelegations := sdk.Coins{sdk.Coin{Denom: ibcDenom, Amount: sdk.NewInt(-10000)}}
	invalidCoinMsg := types.NewMsgLiquidStakeLSM(invalidDelegations, addr1)
	require.Error(t, invalidCoinMsg.ValidateBasic())

	zeroCoinMsg := types.NewMsgLiquidStakeLSM(sdk.Coins{sdk.Coin{Denom: ibcDenom, Amount: sdk.NewInt(0)}}, addr1)
	require.Error(t, zeroCoinMsg.ValidateBasic())

	zeroCoinMsg2 := types.NewMsgLiquidStakeLSM(sdk.NewCoins(sdk.NewCoin(ibcDenom, sdk.ZeroInt())), addr1)
	require.Error(t, zeroCoinMsg2.ValidateBasic())

	invalidAddrMsg := types.NewMsgLiquidStakeLSM(sdk.NewCoins(sdk.NewCoin(ibcDenom, sdk.NewInt(10000))), sdk.AccAddress("test"))
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

	invalidMsg = *msgRegisterHostChain
	invalidMsg.MinimumDeposit = sdk.ZeroInt()
	require.Error(t, invalidMsg.ValidateBasic())
}

func TestMsgUpdateHostChain(t *testing.T) {
	validKVUpdates := []*types.KVUpdate{
		{
			Key:   types.KeySetWithdrawAddress,
			Value: "",
		},
		{
			Key:   types.KeyAddValidator,
			Value: "{\"operator_address\":\"cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt\",\"status\":\"BOND_STATUS_UNSPECIFIED\",\"weight\":\"0\",\"delegated_amount\":\"0\",\"exchange_rate\":\"1\"}",
		},
		{
			Key:   types.KeyRemoveValidator,
			Value: "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt",
		},
		{
			Key:   types.KeyValidatorUpdate,
			Value: "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt",
		},
		{
			Key:   types.KeyValidatorWeight,
			Value: "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt,1",
		},
		{
			Key:   types.KeyRedemptionFee,
			Value: "0",
		},
		{
			Key:   types.KeyDepositFee,
			Value: "0",
		},
		{
			Key:   types.KeyRestakeFee,
			Value: "0",
		},
		{
			Key:   types.KeyUnstakeFee,
			Value: "0",
		},
		{
			Key:   types.KeyUpperCValueLimit,
			Value: "1.1",
		},
		{
			Key:   types.KeyLowerCValueLimit,
			Value: "0.9",
		},
		{
			Key:   types.KeyMinimumDeposit,
			Value: "1",
		},
		{
			Key:   types.KeyLSMValidatorCap,
			Value: "0",
		},
		{
			Key:   types.KeyLSMValidatorCap,
			Value: "1",
		},
		{
			Key:   types.KeyLSMValidatorCap,
			Value: "0.5",
		},
		{
			Key:   types.KeyLSMBondFactor,
			Value: "-1",
		},
		{
			Key:   types.KeyLSMBondFactor,
			Value: "250",
		},
		{
			Key:   types.KeyLSMBondFactor,
			Value: "0",
		},
		{
			Key:   types.KeyActive,
			Value: "true",
		},
		{
			Key:   types.KeyAutocompoundFactor,
			Value: "2",
		},
	}
	msgUpdateHostChain := &types.MsgUpdateHostChain{
		Authority: addr1.String(),
		ChainId:   "chain-1",
		Updates:   validKVUpdates,
	}
	newMsgUpdateHostChain := types.NewMsgUpdateHostChain("chain-1", addr1.String(), validKVUpdates)
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

	invalidKVUpdates := []*types.KVUpdate{
		{
			Key:   types.KeyAddValidator,
			Value: "InvalidJson",
		}, {
			Key:   types.KeyAddValidator,
			Value: "{\"operator_address\":\"cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt\"}",
		}, {
			Key:   types.KeyRemoveValidator,
			Value: "testval",
		}, {
			Key:   types.KeyValidatorUpdate,
			Value: "testval",
		}, {
			Key:   types.KeyValidatorWeight,
			Value: "testval",
		}, {
			Key:   types.KeyValidatorWeight,
			Value: "testval,1",
		}, {
			Key:   types.KeyValidatorWeight,
			Value: "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt,2",
		}, {
			Key:   types.KeyValidatorWeight,
			Value: "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt,invalidDec",
		}, {
			Key:   types.KeyDepositFee,
			Value: "2",
		}, {
			Key:   types.KeyDepositFee,
			Value: "InvalidDec",
		}, {
			Key:   types.KeyRestakeFee,
			Value: "2",
		}, {
			Key:   types.KeyRestakeFee,
			Value: "InvalidDec",
		}, {
			Key:   types.KeyUnstakeFee,
			Value: "2",
		}, {
			Key:   types.KeyUnstakeFee,
			Value: "invalidDec",
		}, {
			Key:   types.KeyRedemptionFee,
			Value: "2",
		}, {
			Key:   types.KeyRedemptionFee,
			Value: "invalidDec",
		}, {
			Key:   types.KeyUpperCValueLimit,
			Value: "-1",
		}, {
			Key:   types.KeyLowerCValueLimit,
			Value: "-1",
		}, {
			Key:   types.KeyLSMValidatorCap,
			Value: "-0.5",
		}, {
			Key:   types.KeyLSMValidatorCap,
			Value: "2",
		}, {
			Key:   types.KeyLSMBondFactor,
			Value: "-0.5",
		}, {
			Key:   types.KeyLSMBondFactor,
			Value: "-1.5",
		}, {
			Key:   types.KeyMinimumDeposit,
			Value: "0",
		}, {
			Key:   types.KeyMinimumDeposit,
			Value: "InvalidInt",
		}, {
			Key:   types.KeyActive,
			Value: "not bool",
		}, {
			Key:   types.KeySetWithdrawAddress,
			Value: "SomeStrHere",
		}, {
			Key:   types.KeyAutocompoundFactor,
			Value: "0",
		}, {
			Key:   types.KeyAutocompoundFactor,
			Value: "InvalidDec",
		}, {
			Key:   "InvalidKey",
			Value: "InvalidKey",
		},
	}
	for _, update := range invalidKVUpdates {
		invalidMsg := types.NewMsgUpdateHostChain("chain-1", addr1.String(), []*types.KVUpdate{update})
		require.Error(t, invalidMsg.ValidateBasic())
	}
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
		AdminAddress: addr1.String(),
		FeeAddress:   addr1.String(),
	})
	require.Error(t, invalidAddrMsg.ValidateBasic())
	require.Panics(t, func() { invalidAddrMsg.GetSigners() })

	invalidParamsMsg := *msgUpdateParams
	invalidParamsMsg.Params.AdminAddress = "test"
	require.Error(t, invalidParamsMsg.ValidateBasic())
}
