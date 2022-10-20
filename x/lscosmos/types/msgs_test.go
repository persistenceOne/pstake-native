package types_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

func TestMsgLiquidStakeRoute(t *testing.T) {
	delegatorAddr := sdk.AccAddress([]byte("delegatorAddress"))
	depositToken := sdk.NewInt64Coin("atom", 10)
	var msg = types.NewMsgLiquidStake(depositToken, delegatorAddr)

	require.Equal(t, msg.Route(), types.RouterKey)
	require.Equal(t, msg.Type(), "msg_liquid_stake")
}

func TestMsgLiquidStakeValidation(t *testing.T) {
	addr := sdk.AccAddress("addr________________")
	addrEmpty := sdk.AccAddress("")
	addrLong := sdk.AccAddress("Purposefully long address")

	atom123 := sdk.NewInt64Coin("atom", 123)
	atom0 := sdk.NewInt64Coin("atom", 0)
	InvalidIBCDenom := sdk.NewInt64Coin("ibc/A", 1)
	InvalidIBCDenom2 := sdk.NewInt64Coin("ibc/AE", 1)
	atomNegative := sdk.Coin{
		Denom:  "atom",
		Amount: sdk.NewInt(-1),
	}
	cases := []struct {
		expectedErr string // empty means no error expected
		msg         *types.MsgLiquidStake
	}{
		{"", types.NewMsgLiquidStake(atom123, addr)}, // valid send
		{"persistence12p6hyur0wdjkvatvd3ujqmr0denjqctyv3ex2umn4nhuy6: invalid address", types.NewMsgLiquidStake(atom123, addrLong)}, // invalid send with long addr sender
		{"0atom: invalid coins", types.NewMsgLiquidStake(atom0, addr)},                                                              // Zero Coin
		{": invalid address", types.NewMsgLiquidStake(atom123, addrEmpty)},                                                          // Nil address
		{"-1atom: invalid coins", types.NewMsgLiquidStake(atomNegative, addr)},                                                      // Negative coin
		{"invalid denom trace hash A: encoding/hex: odd length hex string", types.NewMsgLiquidStake(InvalidIBCDenom, addr)},         // Invalid IBC hash
		{"invalid denom trace hash AE: expected size to be 32 bytes, got 1 bytes", types.NewMsgLiquidStake(InvalidIBCDenom2, addr)}, // Negative IBC hash len

	}

	for _, tc := range cases {
		err := tc.msg.ValidateBasic()
		if tc.expectedErr == "" {
			require.Nil(t, err)
		} else {
			require.EqualError(t, err, tc.expectedErr)
		}
	}
}

func TestNewMsgLiquidStakeGetSignBytes(t *testing.T) {
	addr := sdk.AccAddress([]byte("input"))
	coin := sdk.NewInt64Coin("atom", 10)
	var msg = types.NewMsgLiquidStake(coin, addr)
	res := msg.GetSignBytes()
	excepted := `{"type":"cosmos/MsgLiquidStake","value":{"amount":{"amount":"10","denom":"atom"},"delegator_address":"persistence1d9h8qat5et0urd"}}`
	require.Equal(t, excepted, string(res))

}

func TestMsgLiquidStakeGetSigners(t *testing.T) {
	var msg = types.NewMsgLiquidStake(sdk.NewCoin("atom", sdk.NewInt(10)), sdk.AccAddress([]byte("input111111111111111")))
	res := msg.GetSigners()
	require.Equal(t, fmt.Sprintf("%v", res), "[696E707574313131313131313131313131313131]")
}
