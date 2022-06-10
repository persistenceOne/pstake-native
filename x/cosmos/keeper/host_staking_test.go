package keeper_test

import (
	"fmt"
	"math"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/persistenceOne/pstake-native/x/cosmos/keeper"
	"github.com/persistenceOne/pstake-native/x/cosmos/types"
)

func TestNegativeCoin(t *testing.T) {
	coinFunc := func() {
		sdk.NewCoin("uatom", sdk.NewInt(-1000))
	}

	assert.Panics(t, coinFunc)
}

func TestMulInt(t *testing.T) {
	w, _ := sdk.NewDecFromStr("0.5")

	a := sdk.NewInt(1000000)

	assert.Equal(t, sdk.NewInt(500000).Int64(), w.Mul(sdk.NewDecFromInt(a)).TruncateInt().Int64())
	assert.Equal(t, sdk.NewInt(0).Int64(), w.Mul(sdk.NewDecFromInt(a.Add(sdk.NewInt(1000000)))).TruncateInt().SubRaw(1000000).Int64())
}

func TestGetIdealCurrentDelegations(t *testing.T) {
	denom := "uatom"
	type testValState struct {
		name   string
		weight string
		amount int64
	}
	testMatrix := []struct {
		amount   int64
		given    []testValState
		expected []testValState
	}{
		{
			amount: 5000000,
			given: []testValState{
				{"cosmosVal1", "0.5", 10000000},
				{"cosmosVal2", "0.5", 5000000},
			},
			expected: []testValState{
				{"cosmosVal1", "", 0},
				{"cosmosVal2", "", 5000000},
			},
		},
		{
			amount: -5000000,
			given: []testValState{
				{"cosmosVal1", "0.5", 10000000},
				{"cosmosVal2", "0.5", 5000000},
			},
			expected: []testValState{
				{"cosmosVal1", "", 5000000},
				{"cosmosVal2", "", 0},
			},
		},
		{
			amount: 5000000,
			given: []testValState{
				{"cosmosVal1", "0.9", 10000000},
				{"cosmosVal2", "0.1", 40000000},
			},
			expected: []testValState{
				{"cosmosVal1", "", 39500000},
				{"cosmosVal2", "", -34500000},
			},
		},
		// Equal distribution
		{
			amount: 0,
			given: []testValState{
				{"cosmosVal1", "0.5", 10000000},
				{"cosmosVal2", "0.5", 10000000},
			},
			expected: []testValState{
				{"cosmosVal1", "", 0},
				{"cosmosVal2", "", 0},
			},
		},
		{
			amount: 30000000,
			given: []testValState{
				{"cosmosVal1", "0.5", 10000000},
				{"cosmosVal2", "0.3", 10000000},
				{"cosmosVal3", "0.2", 10000000},
			},
			expected: []testValState{
				{"cosmosVal1", "", 20000000},
				{"cosmosVal2", "", 8000000},
				{"cosmosVal3", "", 2000000},
			},
		},
		{
			amount: -10000000,
			given: []testValState{
				{"cosmosVal1", "0.5", 10000000},
				{"cosmosVal2", "0.3", 10000000},
				{"cosmosVal3", "0.2", 10000000},
			},
			expected: []testValState{
				{"cosmosVal1", "", 0},
				{"cosmosVal2", "", 4000000},
				{"cosmosVal3", "", 6000000},
			},
		},
		{
			amount: -20000000,
			given: []testValState{
				{"cosmosVal1", "0.5", 10000000},
				{"cosmosVal2", "0.3", 10000000},
				{"cosmosVal3", "0.2", 10000000},
			},
			expected: []testValState{
				{"cosmosVal1", "", 5000000},
				{"cosmosVal2", "", 7000000},
				{"cosmosVal3", "", 8000000},
			},
		},
		{
			amount: 10000000,
			given: []testValState{
				{"cosmosVal1", "0.5", 10000000},
				{"cosmosVal2", "0.5", 10000000},
				{"cosmosVal3", "0", 10000000},
			},
			expected: []testValState{
				{"cosmosVal1", "", 10000000},
				{"cosmosVal2", "", 10000000},
				{"cosmosVal3", "", -10000000},
			},
		},
		{
			amount: 10000000,
			given: []testValState{
				{"cosmosVal1", "0.5", 10000000},
				{"cosmosVal2", "0.4", 10000000},
				{"cosmosVal3", "0", 10000000},
				{"cosmosVal4", "0.1", 10000000},
			},
			expected: []testValState{
				{"cosmosVal1", "0.5", 15000000},
				{"cosmosVal2", "0.4", 10000000},
				{"cosmosVal3", "0", -10000000},
				{"cosmosVal4", "0.1", -5000000},
			},
		},
	}

	for i, test := range testMatrix {
		// Create validator state
		givenState := types.WeightedAddressAmounts{}
		expectedMap := map[string]types.WeightedAddressAmount{}
		for i := 0; i < len(test.given); i++ {
			weight, _ := sdk.NewDecFromStr(test.given[i].weight)
			givenState = append(givenState, types.WeightedAddressAmount{
				Address: test.given[i].name,
				Weight:  weight,
				Denom:   denom,
				Amount:  sdk.NewInt(test.given[i].amount),
			})
			expectedMap[test.expected[i].name] = types.WeightedAddressAmount{
				Address: test.expected[i].name,
				Denom:   denom,
				Amount:  sdk.NewInt(test.expected[i].amount),
			}
		}
		// Call getIdealCurrentDelegations function with params
		state := keeper.GetIdealCurrentDelegations(givenState, sdk.NewInt64Coin(denom, int64(math.Abs(float64(test.amount)))), !(test.amount > 0))

		// Assert state
		for j, s := range state {
			expected, ok := expectedMap[s.Address]
			assert.True(t, ok, "Address not is expected list")
			failMsg := fmt.Sprintf("Amounts should be same. Failed for %d case: %d", i, j)
			assert.Equal(t, expected.Amount.BigInt(), s.Amount.BigInt(), failMsg)
		}
	}
}

func testStateData(denom string) types.WeightedAddressAmounts {
	testStruct := []struct {
		name   string
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
	for _, ts := range testStruct {
		weight, _ := sdk.NewDecFromStr(ts.weight)
		state = append(state, types.WeightedAddressAmount{
			Weight:  weight,
			Amount:  sdk.NewInt(ts.amount),
			Address: sdk.ValAddress(ts.name).String(),
			Denom:   denom,
		})
	}
	return state
}

func (suite *IntegrationTestSuite) TestDivideAmountIntoValidatorSet() {
	// Test setup
	app, ctx := suite.app, suite.ctx

	// Set validator set weighted amount
	params := app.CosmosKeeper.GetParams(ctx)
	state := testStateData(params.StakingDenom)
	suite.SetupValWeightedAmounts(state)

	// Test data
	testMatrix := []struct {
		given    int64
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
				"cosmosVal3": 8500000,
				"cosmosVal4": 1500000,
			},
		},
		{
			given: 20000000,
			expected: map[string]int64{
				"cosmosVal3": 11500000,
				"cosmosVal1": 7000000,
				"cosmosVal4": 1500000,
			},
		},
		{
			given: 30000000,
			expected: map[string]int64{
				"cosmosVal1": 11000000,
				"cosmosVal3": 14500000,
				"cosmosVal4": 4500000,
			},
		},
		{
			given: 50000000,
			expected: map[string]int64{
				"cosmosVal1": 19000000,
				"cosmosVal2": 2000000,
				"cosmosVal3": 20500000,
				"cosmosVal4": 8500000,
			},
		},
	}

	// Create input parameters
	for _, test := range testMatrix {
		// Create state
		givenCoin := sdk.NewInt64Coin(params.StakingDenom, test.given)
		expectedMap := map[string]int64{}
		for k, v := range test.expected {
			expectedMap[sdk.ValAddress(k).String()] = v
		}

		// Run getIdealCurrentDelegations function with params
		valAmounts, err := app.CosmosKeeper.FetchValidatorsToDelegate(ctx, givenCoin)
		suite.Nil(err, "Error is not nil for validator to delegate")
		// Check outputs
		actualMap := map[string]int64{}
		for _, va := range valAmounts {
			actualMap[va.Validator.String()] = va.Amount.Amount.Int64()
		}
		suite.Equal(expectedMap, actualMap, "Matching val distribution")
	}
}

func (suite *IntegrationTestSuite) TestUndelegateDivideAmountIntoValidatorSet() {
	// Test setup
	app, ctx := suite.app, suite.ctx

	// Set validator set weighted amount
	params := app.CosmosKeeper.GetParams(ctx)
	state := testStateData(params.StakingDenom)
	suite.SetupValWeightedAmounts(state)

	// Test data
	testMatrix := []struct {
		given    int64
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
				"cosmosVal2": 5000000,
			},
		},
		{
			given:    0,
			expected: map[string]int64{},
		},
		{
			given: 20000000,
			expected: map[string]int64{
				"cosmosVal1": 9000000,
				"cosmosVal2": 6000000,
				"cosmosVal5": 5000000,
			},
		},
		{
			given: 30000000,
			expected: map[string]int64{
				"cosmosVal1": 13000000,
				"cosmosVal2": 9000000,
				"cosmosVal3": 3000000,
				"cosmosVal5": 5000000,
			},
		},
		{
			given: 35000000,
			expected: map[string]int64{
				"cosmosVal1": 15000000,
				"cosmosVal2": 10000000,
				"cosmosVal3": 5000000,
				"cosmosVal5": 5000000,
			},
		},
	}

	// Create input parameters
	for _, test := range testMatrix {
		// Create state
		givenCoin := sdk.NewInt64Coin(params.StakingDenom, test.given)
		expectedMap := map[string]int64{}
		for k, v := range test.expected {
			expectedMap[sdk.ValAddress(k).String()] = v
		}

		// Run getIdealCurrentDelegations function with params
		valAmounts, err := app.CosmosKeeper.FetchValidatorsToUndelegate(ctx, givenCoin)
		suite.Nil(err, "Error is not nil for validator to delegate")

		// Check outputs
		actualMap := map[string]int64{}
		for _, va := range valAmounts {
			actualMap[va.Validator.String()] = va.Amount.Amount.Int64()
		}
		suite.Equal(expectedMap, actualMap, "Matching val distribution")
	}
}
