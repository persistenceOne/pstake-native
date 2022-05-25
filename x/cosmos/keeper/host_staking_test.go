package keeper_test

import (
	"fmt"
	"sort"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/persistenceOne/pstake-native/x/cosmos/types"
	"github.com/persistenceOne/pstake-native/x/cosmos/keeper"
)

func TestNegativeCoin(t *testing.T) {
	coinFunc := func() {
		sdk.NewCoin("uatom", sdk.NewInt(-1000))
	}
	
	assert.Panics(t, coinFunc)
}

func TestGetIdealCurrentDelegations(t *testing.T) {
	denom := "uatom"
	type testValState struct {
		name   string
		weight string
		amount int64
	}
	testMatrix := []struct {
		given  []testValState
		expected []testValState
	}{
		{
			given: []testValState{
				{
					name:   "cosmosVal1",
					weight: "0.5",
					amount: 10000000,
				},
				{
					name:   "cosmosVal2",
					weight: "0.5",
					amount: 5000000,
				},
			},
			expected: []testValState{
				{
					name:   "cosmosVal1",
					amount: -2500000,
				},
				{
					name:   "cosmosVal2",
					amount: 2500000,
				},
			},
		},
		// Equal distribution
		{
			given: []testValState{
				{
					name:   "cosmosVal1",
					weight: "0.5",
					amount: 10000000,
				},
				{
					name:   "cosmosVal2",
					weight: "0.5",
					amount: 10000000,
				},
			},
			expected: []testValState{
				{
					name:   "cosmosVal1",
					amount: 0,
				},
				{
					name:   "cosmosVal2",
					amount: 0,
				},
			},
		},
		{
			given: []testValState{
				{
					name:   "cosmosVal1",
					weight: "0.5",
					amount: 10000000,
				},
				{
					name:   "cosmosVal2",
					weight: "0.3",
					amount: 10000000,
				},
				{
					name:   "cosmosVal3",
					weight: "0.2",
					amount: 10000000,
				},
			},
			expected: []testValState{
				{
					name:   "cosmosVal1",
					amount: 5000000,
				},
				{
					name:   "cosmosVal2",
					amount: -1000000,
				},
				{
					name:   "cosmosVal3",
					amount: -4000000,
				},
			},
		},
		{
			given: []testValState{
				{
					name:   "cosmosVal1",
					weight: "0.5",
					amount: 10000000,
				},
				{
					name:   "cosmosVal2",
					weight: "0.5",
					amount: 10000000,
				},
				{
					name:   "cosmosVal3",
					weight: "0",
					amount: 10000000,
				},
			},
			expected: []testValState{
				{
					name:   "cosmosVal1",
					amount: 5000000,
				},
				{
					name:   "cosmosVal2",
					amount: 5000000,
				},
				{
					name:   "cosmosVal3",
					amount: -10000000,
				},
			},
		},
	}

	for _, test := range testMatrix {
		// Create validator state
		givenState := types.WeightedAddressAmounts{}
		expectedMap := map[string]types.WeightedAddressAmount{}
		for i := 0; i < len(test.given); i++ {
			weight, _ :=  sdk.NewDecFromStr(test.given[i].weight)
			givenState = append(givenState, types.WeightedAddressAmount{
				Address: test.given[i].name,
				Weight: weight,
				Denom: denom,
				Amount: sdk.NewInt(test.given[i].amount),
			})
			expectedMap[test.expected[i].name] = types.WeightedAddressAmount{
				Address: test.expected[i].name,
				Denom: denom,
				Amount: sdk.NewInt(test.expected[i].amount),
			}
		}
		// Call getIdealCurrentDelegations function with params
		state := keeper.GetIdealCurrentDelegations(givenState, denom)

		// Assert state
		for _, s := range state {
			expected, ok := expectedMap[s.Address]
			assert.True(t, ok, "Address not is expected list")
			
			assert.Equal(t, expected.Amount.BigInt(), s.Amount.BigInt(), "Amounts should be same")
		}
	}
}

func TestNormalizedWeightedAddressAmounts(t *testing.T) {
	denom := "uatom"
	testMatrix := []struct {
		given  []int64
		expected []int64
	}{
		{
			given: []int64{10000000, -10000000, 5000000, 10000000},
			expected: []int64{20000000, 0, 15000000, 20000000},
		},
		{
			given: []int64{10000000, 0, 5000000, 10000000},
			expected: []int64{10000000, 0, 5000000, 10000000},
		},
		{
			given: []int64{10000000, -10000000, -50000000, 10000000},
			expected: []int64{60000000, 40000000, 0, 60000000},
		},
	}

	for _, test := range testMatrix {
		// Create state
		givenState := types.WeightedAddressAmounts{}
		expectedMap := map[string]types.WeightedAddressAmount{}
		for i := 0; i < len(test.given); i++ {
			name := fmt.Sprintf("test%d", i)
			givenState = append(givenState, types.WeightedAddressAmount{
				Address: name,
				Denom: denom,
				Amount: sdk.NewInt(test.given[i]),
			})
			expectedMap[name] = types.WeightedAddressAmount{
				Address: name,
				Denom: denom,
				Amount: sdk.NewInt(test.expected[i]),
			}
		}
		// Call getIdealCurrentDelegations function with params
		state := keeper.NormalizedWeightedAddressAmounts(givenState)

		// Assert state
		for _, s := range state {
			expected, ok := expectedMap[s.Address]
			assert.True(t, ok, "Address not is expected list")
			
			assert.Equal(t, expected.Amount.BigInt(), s.Amount.BigInt(), "Amounts should be same")
		}
	}
}

func testStateData(denom string) types.WeightedAddressAmounts {
	testStruct := []struct {
		name string
		weight string
		amount int64
	}{
		{
			name:   "cosmosVal1",
			weight: "0.4",
			amount: 15000000, // ideal: 14000000
		},
		{
			name:   "cosmosVal2",
			weight: "0.2",
			amount: 10000000, // ideal: 7000000
		},
		{
			name:   "cosmosVal3",
			weight: "0.3",
			amount: 5000000, // ideal: 10500000
		},
		{
			name:   "cosmosVal4",
			weight: "0.1",
			amount: 0, // ideal: 3500000
		},
		{
			name:   "cosmosVal5",
			weight: "0",
			amount: 5000000, // ideal: 0
		},
	}
	// Create state
	state := types.WeightedAddressAmounts{}
	for _, ts := range testStruct{
		weight, _ :=  sdk.NewDecFromStr(ts.weight)
		state = append(state, types.WeightedAddressAmount{
			Weight: weight,
			Amount: sdk.NewInt(ts.amount),
			Address: sdk.ValAddress(ts.name).String(),
			Denom: denom,
		})
	}
	return state
}

func TestDivideAmountIntoValidatorSet(t *testing.T) {
	denom := "uatom"
	state := testStateData(denom)
	testMatrix := []struct {
		given int64
		expected map[string]int64
	}{
		{
			given: 1000,
			expected: map[string]int64{
				"cosmosVal3": 1000,
			},
		},
		{
			given: 10000000,
			expected: map[string]int64{
				"cosmosVal3": 10000000,
			},
		},
		{
			given: 0,
			expected: map[string]int64{},
		},
		{
			given: 20000000,
			expected: map[string]int64{
				"cosmosVal3": 10500000,
				"cosmosVal4": 8500000,
				"cosmosVal1": 1000000,
			},
		},
		{
			given: 30000000,
			expected: map[string]int64{
				"cosmosVal1": 6000000,
				"cosmosVal2": 3000000,
				"cosmosVal3": 12000000,
				"cosmosVal4": 9000000,
			},
		},
		{
			given: 50000000,
			expected: map[string]int64{
				"cosmosVal1": 14000000,
				"cosmosVal2": 7000000,
				"cosmosVal3": 18000000,
				"cosmosVal4": 11000000,
			},
		},
	}

	idealCurDis := keeper.GetIdealCurrentDelegations(state, denom)
	idealCurDis = keeper.NormalizedWeightedAddressAmounts(idealCurDis)

	sort.Sort(sort.Reverse(idealCurDis))

	// Create input parameters
	for _, test := range testMatrix {
		// Create state
		givenCoin := sdk.NewInt64Coin(denom, test.given)
		expectedMap := map[string]int64{}
		for k, v := range test.expected {
			expectedMap[sdk.ValAddress(k).String()] = v
		}

		// Run getIdealCurrentDelegations function with params
		valAmounts, err := keeper.DivideAmountIntoValidatorSet(idealCurDis, givenCoin)
		assert.Nil(t, err, "Error is not nil for divideAmountIntoValidatorSet")
		
		// Assert outputs
		actualMap := map[string]int64{}
		for _, va := range valAmounts {
			actualMap[va.Validator.String()] = va.Amount.Amount.Int64()
		}

		assert.Equal(t, expectedMap, actualMap, "Matching val distribution")
	}
}

func TestUndelegateDivideAmountIntoValidatorSet(t *testing.T) {
	denom := "uatom"
	state := testStateData(denom)
	testMatrix := []struct {
		given int64
		expected map[string]int64
	}{
		{
			given: 1000,
			expected: map[string]int64{
				"cosmosVal5": 1000,
			},
		},
		{
			given: 10000000,
			expected: map[string]int64{
				"cosmosVal5": 5000000, 
				"cosmosVal4": 5000000,
			},
		},
		{
			given: 0,
			expected: map[string]int64{},
		},
		{
			given: 20000000,
			expected: map[string]int64{
				"cosmosVal1": 15000000,
				"cosmosVal2": 5000000,
			},
		},
	}

	// Create input parameters
	for _, test := range testMatrix {
		// Create state
		givenCoin := sdk.NewInt64Coin(denom, test.given)
		expectedMap := map[string]int64{}
		for k, v := range test.expected {
			expectedMap[sdk.ValAddress(k).String()] = v
		}

		// Run getIdealCurrentDelegations function with params
		valAmounts, _ := keeper.DivideUndelegateAmountIntoValidatorSet(state, givenCoin)
		
		// Assert outputs
		actualMap := map[string]int64{}
		for _, va := range valAmounts {
			actualMap[va.Validator.String()] = va.Amount.Amount.Int64()
		}

		assert.Equal(t, expectedMap, actualMap, "Matching val distribution")
	}
}