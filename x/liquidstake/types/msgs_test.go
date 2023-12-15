package types_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstake/types"
)

func TestMsgLiquidStake(t *testing.T) {
	delegatorAddr := sdk.AccAddress(crypto.AddressHash([]byte("delegatorAddr")))
	stakingCoin := sdk.NewCoin("uxprt", math.NewInt(1))

	testCases := []struct {
		expectedErr string
		msg         *types.MsgLiquidStake
	}{
		{
			"", // empty means no error expected
			types.NewMsgLiquidStake(delegatorAddr, stakingCoin),
		},
		{
			"invalid delegator address \"\": empty address string is not allowed: invalid address",
			types.NewMsgLiquidStake(sdk.AccAddress{}, stakingCoin),
		},
		{
			"staking amount must not be zero: invalid request",
			types.NewMsgLiquidStake(delegatorAddr, sdk.NewCoin("token", math.NewInt(0))),
		},
	}

	for _, tc := range testCases {
		require.IsType(t, &types.MsgLiquidStake{}, tc.msg)
		require.Equal(t, types.MsgTypeLiquidStake, tc.msg.Type())
		require.Equal(t, types.RouterKey, tc.msg.Route())
		require.Equal(t, sdk.MustSortJSON(types.ModuleCdc.MustMarshalJSON(tc.msg)), tc.msg.GetSignBytes())

		err := tc.msg.ValidateBasic()
		if tc.expectedErr == "" {
			require.Nil(t, err)
			signers := tc.msg.GetSigners()
			require.Len(t, signers, 1)
			require.Equal(t, tc.msg.GetDelegator(), signers[0])
		} else {
			require.EqualError(t, err, tc.expectedErr)
		}
	}
}

func TestMsgLiquidUnstake(t *testing.T) {
	delegatorAddr := sdk.AccAddress(crypto.AddressHash([]byte("delegatorAddr")))
	stakingCoin := sdk.NewCoin("stk/uxprt", math.NewInt(1))

	testCases := []struct {
		expectedErr string
		msg         *types.MsgLiquidUnstake
	}{
		{
			"", // empty means no error expected
			types.NewMsgLiquidUnstake(delegatorAddr, stakingCoin),
		},
		{
			"invalid delegator address \"\": empty address string is not allowed: invalid address",
			types.NewMsgLiquidUnstake(sdk.AccAddress{}, stakingCoin),
		},
		{
			"unstaking amount must not be zero: invalid request",
			types.NewMsgLiquidUnstake(delegatorAddr, sdk.NewCoin("btoken", math.NewInt(0))),
		},
	}

	for _, tc := range testCases {
		require.IsType(t, &types.MsgLiquidUnstake{}, tc.msg)
		require.Equal(t, types.MsgTypeLiquidUnstake, tc.msg.Type())
		require.Equal(t, types.RouterKey, tc.msg.Route())
		require.Equal(t, sdk.MustSortJSON(types.ModuleCdc.MustMarshalJSON(tc.msg)), tc.msg.GetSignBytes())

		err := tc.msg.ValidateBasic()
		if tc.expectedErr == "" {
			require.Nil(t, err)
			signers := tc.msg.GetSigners()
			require.Len(t, signers, 1)
			require.Equal(t, tc.msg.GetDelegator(), signers[0])
		} else {
			require.EqualError(t, err, tc.expectedErr)
		}
	}
}
