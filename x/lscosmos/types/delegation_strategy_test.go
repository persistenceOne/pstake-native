package types_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
	"github.com/stretchr/testify/require"
	"sort"
	"testing"
)

func TestWeightedAddressAmounts(t *testing.T) {
	wa1 := types.NewWeightedAddressAmount("addr1", sdk.NewDecWithPrec(0, 1), sdk.NewCoin("uatom", sdk.NewInt(0)), sdk.NewCoin("uatom", sdk.NewInt(100)))
	wa2 := types.NewWeightedAddressAmount("addr2", sdk.NewDecWithPrec(5, 1), sdk.NewCoin("uatom", sdk.NewInt(1001)), sdk.NewCoin("uatom", sdk.NewInt(100)))
	wa3 := types.NewWeightedAddressAmount("addr3", sdk.NewDecWithPrec(2, 1), sdk.NewCoin("uatom", sdk.NewInt(1010)), sdk.NewCoin("uatom", sdk.NewInt(100)))
	wa4 := types.NewWeightedAddressAmount("addr4", sdk.NewDecWithPrec(3, 1), sdk.NewCoin("uatom", sdk.NewInt(1003)), sdk.NewCoin("uatom", sdk.NewInt(100)))

	// create new weighted address amounts using the above 4 weightedAddressAmount
	weightedAddressAmounts := types.NewWeightedAddressAmounts([]types.WeightedAddressAmount{wa1, wa2, wa3, wa4})

	// get zero weightedAddressAmounts
	zeroWeighted := weightedAddressAmounts.GetZeroWeighted()
	require.Equal(t, types.WeightedAddressAmounts{wa1}, zeroWeighted)

	// get zero valued weightedAddressAmounts
	zeroValued := weightedAddressAmounts.GetZeroValued()
	require.Equal(t, types.WeightedAddressAmounts{wa1}, zeroValued)

	// calculate the total amount
	totalAmount := weightedAddressAmounts.TotalAmount("uatom")
	require.Equal(t, sdk.NewCoin("uatom", sdk.NewInt(3014)), totalAmount)

	// get weighted address map
	weightedAddressMap := types.GetWeightedAddressMap(weightedAddressAmounts)
	newWeightedAddressMap := map[string]sdk.Dec{}
	newWeightedAddressMap[wa1.Address] = wa1.Weight
	newWeightedAddressMap[wa2.Address] = wa2.Weight
	newWeightedAddressMap[wa3.Address] = wa3.Weight
	newWeightedAddressMap[wa4.Address] = wa4.Weight
	require.Equal(t, newWeightedAddressMap, weightedAddressMap)

	// get zero and non-zero weightedAddressAmounts
	zeroWeighted, nonZeroWeighted := types.GetZeroNonZeroWightedAddrAmts(weightedAddressAmounts)
	require.Equal(t, types.WeightedAddressAmounts{wa1}, zeroWeighted)
	require.Equal(t, types.WeightedAddressAmounts{wa2, wa3, wa4}, nonZeroWeighted)

	// sort weightedAddressAmounts
	sort.Sort(weightedAddressAmounts)
	require.Equal(t, wa1, weightedAddressAmounts[0])
	require.Equal(t, wa2, weightedAddressAmounts[1])
	require.Equal(t, wa4, weightedAddressAmounts[2])
	require.Equal(t, wa3, weightedAddressAmounts[3])
}
