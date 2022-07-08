package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewWeightedAddressAmounts(t *testing.T) {
	weightedAddresses := WeightedAddressAmounts{
		NewWeightedAddressAmount("A", sdk.NewDec(1), sdk.NewInt64Coin("test", 100), sdk.NewInt64Coin("test", 100)),
		NewWeightedAddressAmount("B", sdk.NewDec(1), sdk.NewInt64Coin("test", 1000), sdk.NewInt64Coin("test", 100)),
		NewWeightedAddressAmount("C", sdk.NewDec(1), sdk.NewInt64Coin("test", 10000), sdk.NewInt64Coin("test", 100)),
		NewWeightedAddressAmount("D", sdk.NewDec(1), sdk.NewInt64Coin("test", 102), sdk.NewInt64Coin("test", 100)),
		NewWeightedAddressAmount("E", sdk.NewDec(1), sdk.NewInt64Coin("test", 109), sdk.NewInt64Coin("test", 100)),
	}

	weightedAddresses1 := WeightedAddressAmounts{
		NewWeightedAddressAmount("A", sdk.NewDec(1), sdk.NewInt64Coin("test", 100), sdk.NewInt64Coin("test", 100)),
		NewWeightedAddressAmount("B", sdk.NewDec(1), sdk.NewInt64Coin("test", 1000), sdk.NewInt64Coin("test", 100)),
		NewWeightedAddressAmount("D", sdk.NewDec(1), sdk.NewInt64Coin("test", 102), sdk.NewInt64Coin("test", 100)),
		NewWeightedAddressAmount("E", sdk.NewDec(1), sdk.NewInt64Coin("test", 109), sdk.NewInt64Coin("test", 100)),
		NewWeightedAddressAmount("C", sdk.NewDec(1), sdk.NewInt64Coin("test", 10000), sdk.NewInt64Coin("test", 100)),
	}

	require.NotEqual(t, weightedAddresses, weightedAddresses1)

	require.Equal(t, weightedAddresses.Sort(), weightedAddresses1.Sort())

}
