package types

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMsgLiquidStakeRoute(t *testing.T) {
	delegatorAddr := sdk.AccAddress([]byte("delegatorAddress"))
	depositToken := sdk.NewInt64Coin("atom", 10)
	var msg = NewMsgLiquidStake(depositToken, delegatorAddr)

	require.Equal(t, msg.Route(), RouterKey)
	require.Equal(t, msg.Type(), "msg_liquid_stake")
}

func TestMsgLiquidStakeValidation(t *testing.T) {
	addr := sdk.AccAddress([]byte("delegatorAdd______________________"))
	addrEmpty := sdk.AccAddress([]byte(""))
	addrLong := sdk.AccAddress([]byte("Purposefully long address"))

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
		msg         *MsgLiquidStake
	}{
		{"", NewMsgLiquidStake(atom123, addr)},                                                                                // valid send
		{"", NewMsgLiquidStake(atom123, addrLong)},                                                                            // valid send with long addr sender
		{"0atom: invalid coins", NewMsgLiquidStake(atom0, addr)},                                                              // Zero Coin
		{": invalid address", NewMsgLiquidStake(atom123, addrEmpty)},                                                          // Nil address
		{"-1atom: invalid coins", NewMsgLiquidStake(atomNegative, addr)},                                                      // Negative coin
		{"invalid denom trace hash A: encoding/hex: odd length hex string", NewMsgLiquidStake(InvalidIBCDenom, addr)},         // Invalid IBC hash
		{"invalid denom trace hash AE: expected size to be 32 bytes, got 1 bytes", NewMsgLiquidStake(InvalidIBCDenom2, addr)}, // Negative IBC hash len

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
	var msg = NewMsgLiquidStake(coin, addr)
	res := msg.GetSignBytes()
	excepted := `{"amount":{"amount":"10","denom":"atom"},"delegator_address":"cosmos1d9h8qat57ljhcm"}`
	require.Equal(t, excepted, string(res))

}

func TestMsgLiquidStakeGetSigners(t *testing.T) {
	var msg = NewMsgLiquidStake(sdk.NewCoin("atom", sdk.NewInt(10)), sdk.AccAddress([]byte("input111111111111111")))
	res := msg.GetSigners()
	require.Equal(t, fmt.Sprintf("%v", res), "[696E707574313131313131313131313131313131]")
}
