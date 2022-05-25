package keeper

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/persistenceOne/pstake-native/x/cosmos/types"
)

func TestNegativeCoin(t *testing.T) {
	coinFunc := func() {
		sdk.NewCoin("uatom", sdk.NewInt(-1000))
	}
	
	assert.Panics(t, coinFunc)
}

func TestValidatorAddr(t *testing.T) {
	addr := "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt"

	valAddr, err := sdk.ValAddressFromBech32(addr)
	assert.NoError(t, err, "Error from conversion of text to val address")
	assert.Equal(t, addr, valAddr.String(), "Val addr from bech32 should be same")
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
					name:   "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt",
					weight: "0.5",
					amount: 10000000,
				},
				{
					name:   "cosmosvaloper1lcck2cxh7dzgkrfk53kysg9ktdrsjj6jfwlnm2",
					weight: "0.5",
					amount: 5000000,
				},
			},
			expected: []testValState{
				{
					name:   "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt",
					amount: -2500000,
				},
				{
					name:   "cosmosvaloper1lcck2cxh7dzgkrfk53kysg9ktdrsjj6jfwlnm2",
					amount: 2500000,
				},
			},
		},
		// Equal distribution
		{
			given: []testValState{
				{
					name:   "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt",
					weight: "0.5",
					amount: 10000000,
				},
				{
					name:   "cosmosvaloper1lcck2cxh7dzgkrfk53kysg9ktdrsjj6jfwlnm2",
					weight: "0.5",
					amount: 10000000,
				},
			},
			expected: []testValState{
				{
					name:   "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt",
					amount: 0,
				},
				{
					name:   "cosmosvaloper1lcck2cxh7dzgkrfk53kysg9ktdrsjj6jfwlnm2",
					amount: 0,
				},
			},
		},
		{
			given: []testValState{
				{
					name:   "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt",
					weight: "0.5",
					amount: 10000000,
				},
				{
					name:   "cosmosvaloper1lcck2cxh7dzgkrfk53kysg9ktdrsjj6jfwlnm2",
					weight: "0.3",
					amount: 10000000,
				},
				{
					name:   "cosmosvaloper10vcqjzphfdlumas0vp64f0hruhrqxv0cd7wdy2",
					weight: "0.2",
					amount: 10000000,
				},
			},
			expected: []testValState{
				{
					name:   "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt",
					amount: 5000000,
				},
				{
					name:   "cosmosvaloper1lcck2cxh7dzgkrfk53kysg9ktdrsjj6jfwlnm2",
					amount: -1000000,
				},
				{
					name:   "cosmosvaloper10vcqjzphfdlumas0vp64f0hruhrqxv0cd7wdy2",
					amount: -4000000,
				},
			},
		},
		{
			given: []testValState{
				{
					name:   "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt",
					weight: "0.5",
					amount: 10000000,
				},
				{
					name:   "cosmosvaloper1lcck2cxh7dzgkrfk53kysg9ktdrsjj6jfwlnm2",
					weight: "0.5",
					amount: 10000000,
				},
				{
					name:   "cosmosvaloper10vcqjzphfdlumas0vp64f0hruhrqxv0cd7wdy2",
					weight: "0",
					amount: 10000000,
				},
			},
			expected: []testValState{
				{
					name:   "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt",
					amount: 5000000,
				},
				{
					name:   "cosmosvaloper1lcck2cxh7dzgkrfk53kysg9ktdrsjj6jfwlnm2",
					amount: 5000000,
				},
				{
					name:   "cosmosvaloper10vcqjzphfdlumas0vp64f0hruhrqxv0cd7wdy2",
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
		state := getIdealCurrentDelegations(givenState, denom)

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
		state := normalizedWeightedAddressAmounts(givenState)

		// Assert state
		for _, s := range state {
			expected, ok := expectedMap[s.Address]
			assert.True(t, ok, "Address not is expected list")
			
			assert.Equal(t, expected.Amount.BigInt(), s.Amount.BigInt(), "Amounts should be same")
		}
	}
}


func TestDivideAmountIntoValidatorSet(t *testing.T) {
	denom := "uatom"
	testState := []struct {
		name string
		weight string
		amount int64
	}{
		{
			name:   "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt",
			weight: "0.5",
			amount: 15000000,
		},
		{
			name:   "cosmosvaloper1lcck2cxh7dzgkrfk53kysg9ktdrsjj6jfwlnm2",
			weight: "0.2",
			amount: 10000000,
		},
		{
			name:   "cosmosvaloper10vcqjzphfdlumas0vp64f0hruhrqxv0cd7wdy2",
			weight: "0.3",
			amount: 5000000,
		},
		{
			name:   "cosmosvaloper10khgeppewe4rgfrcy809r9h00aquwxxxgwgwa5",
			weight: "0.1",
			amount: 0,
		},
		{
			name:   "cosmosvaloper1ey69r37gfxvxg62sh4r0ktpuc46pzjrm873ae8",
			weight: "0",
			amount: 5000000,
		},
	}
	testMatrix := []struct {
		given int64
		expected map[string]int64
	}{
		{
			given: 1000,
			expected: map[string]int64{
				"cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt": 1000,
			},
		},
		{
			given: 10000000,
			expected: map[string]int64{
				"cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt": 10000000,
			},
		},
		{
			given: 0,
			expected: map[string]int64{},
		},
		{
			given: 20000000,
			expected: map[string]int64{
				"cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt": 15000000,
				"cosmosvaloper1lcck2cxh7dzgkrfk53kysg9ktdrsjj6jfwlnm2": 5000000,
			},
		},
		{
			given: 30000000,
			expected: map[string]int64{
				"cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt": 15000000,
				"cosmosvaloper1lcck2cxh7dzgkrfk53kysg9ktdrsjj6jfwlnm2": 10000000,
				"cosmosvaloper10vcqjzphfdlumas0vp64f0hruhrqxv0cd7wdy2": 5000000,
			},
		},
		{
			given: 50000000,
			expected: map[string]int64{
				"cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt": 25000000,
				"cosmosvaloper1lcck2cxh7dzgkrfk53kysg9ktdrsjj6jfwlnm2": 14000000,
				"cosmosvaloper10vcqjzphfdlumas0vp64f0hruhrqxv0cd7wdy2": 11000000,
				"cosmosvaloper10khgeppewe4rgfrcy809r9h00aquwxxxgwgwa5": 2000000,
			},
		},
	}

	// Create state
	state := types.WeightedAddressAmounts{}
	for _, ts := range testState{
		weight, _ :=  sdk.NewDecFromStr(ts.weight)
		state = append(state, types.WeightedAddressAmount{
			Weight: weight,
			Amount: sdk.NewInt(ts.amount),
			Address: ts.name,
			Denom: denom,
		})
	}
	// Create input parameters
	for _, test := range testMatrix {
		// Create state
		givenCoin := sdk.NewInt64Coin(denom, test.given)
		expectedMap := map[string]int64{}
		for k, v := range test.expected {
			addr, _ := sdk.ValAddressFromBech32(k)
			expectedMap[addr.String()] = v
		}

		// Run getIdealCurrentDelegations function with params
		valAmounts, _ := divideAmountIntoValidatorSet(state, givenCoin)
		
		// Assert outputs
		actualMap := map[string]int64{}
		for _, va := range valAmounts {
			actualMap[va.Validator.String()] = va.Amount.Amount.Int64()
		}

		assert.Equal(t, expectedMap, actualMap, "Matching val distribution")
	}
}