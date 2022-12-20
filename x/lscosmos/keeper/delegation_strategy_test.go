package keeper_test

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/stretchr/testify/assert"

	"github.com/persistenceOne/pstake-native/x/lscosmos/keeper"
	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

var HostStakingDenom = "uatom"

func (suite *IntegrationTestSuite) SetupAllowListedValSetAndDelegationState(ws types.WeightedAddressAmounts) {
	app, ctx := suite.app, suite.ctx
	allList := make([]types.AllowListedValidator, ws.Len())
	allowListedVal := types.AllowListedValidators{AllowListedValidators: allList}

	delList := make([]types.HostAccountDelegation, ws.Len())
	delegationState := types.DelegationState{HostAccountDelegations: delList}
	delegationState.HostAccountDelegations = delList
	delegationState.HostChainDelegationAddress = "cosmosdelegationAddr1"
	delegationState.HostDelegationAccountBalance = sdk.NewCoins(sdk.NewCoin(HostStakingDenom, sdk.NewInt(0)))

	for i, w := range ws {
		allowListedVal.AllowListedValidators[i].ValidatorAddress = w.Address
		allowListedVal.AllowListedValidators[i].TargetWeight = w.Weight
		delegationState.HostAccountDelegations[i].ValidatorAddress = w.Address
		delegationState.HostAccountDelegations[i].Amount = w.Coin()

	}
	app.LSCosmosKeeper.SetModuleState(ctx, true)
	app.LSCosmosKeeper.SetAllowListedValidators(ctx, allowListedVal)
	app.LSCosmosKeeper.SetDelegationState(ctx, delegationState)
}

func (suite *IntegrationTestSuite) TestDivideAmountIntoValidatorSet() {
	_, ctx := suite.app, suite.ctx

	denom := HostStakingDenom
	state := testStateData(denom)
	suite.SetupAllowListedValSetAndDelegationState(state)

	// Test data
	testMatrix := []struct {
		given    int64
		expected map[string]int64
	}{
		{
			given: 1000,
			expected: map[string]int64{
				"cosmosvalidatorAddr3": 1000,
			},
		},
		{
			given: 10000000,
			expected: map[string]int64{
				"cosmosvalidatorAddr3": 8500000,
				"cosmosvalidatorAddr4": 1500000,
			},
		},
		{
			given: 20000000,
			expected: map[string]int64{
				"cosmosvalidatorAddr3": 11500000,
				"cosmosvalidatorAddr1": 7000000,
				"cosmosvalidatorAddr4": 1500000,
			},
		},
		{
			given: 30000000,
			expected: map[string]int64{
				"cosmosvalidatorAddr1": 11000000,
				"cosmosvalidatorAddr3": 14500000,
				"cosmosvalidatorAddr4": 4500000,
			},
		},
		{
			given: 50000000,
			expected: map[string]int64{
				"cosmosvalidatorAddr1": 19000000,
				"cosmosvalidatorAddr2": 2000000,
				"cosmosvalidatorAddr3": 20500000,
				"cosmosvalidatorAddr4": 8500000,
			},
		},
	}

	for _, test := range testMatrix {
		givenCoin := sdk.NewInt64Coin(HostStakingDenom, test.given)
		expectedMap := map[string]int64{}

		for k, v := range test.expected {
			valAddress, _ := Bech32ifyValAddressBytes(types.CosmosValOperPrefix, sdk.ValAddress(k))
			expectedMap[valAddress] = v
		}

		allowlistedVals := suite.app.LSCosmosKeeper.GetAllowListedValidators(ctx)
		delegationState := suite.app.LSCosmosKeeper.GetDelegationState(ctx)

		// Run getIdealCurrentDelegations function with params
		valAmounts, err := keeper.FetchValidatorsToDelegate(allowlistedVals, delegationState, givenCoin)
		suite.Nil(err, "Error is not nil for validator to delegate")

		// Check outputs
		actualMap := map[string]int64{}
		for _, va := range valAmounts {
			actualMap[va.ValidatorAddr] = va.Amount.Amount.Int64()
		}
		suite.Equal(expectedMap, actualMap, "Matching val distribution")

	}
}

func (suite *IntegrationTestSuite) TestDivideAmountIntoStateValidatorSet() {
	_, ctx := suite.app, suite.ctx

	// Test data
	testMatrix := []struct {
		state    map[string][]string
		given    int64
		expected map[string]int64
	}{
		{
			state: map[string][]string{
				"cosmosvalidatorAddr1": {"4000000", "0.1"},
				"cosmosvalidatorAddr2": {"8000000", "0.2"},
				"cosmosvalidatorAddr3": {"8000000", "0.2"},
				"cosmosvalidatorAddr4": {"20000000", "0.5"},
			},
			given: 13028679724,
			expected: map[string]int64{
				"cosmosvalidatorAddr1": 1302867972,
				"cosmosvalidatorAddr2": 2605735944,
				"cosmosvalidatorAddr3": 2605735944,
				"cosmosvalidatorAddr4": 6514339864,
			},
		},
	}

	// Create input parameters
	for _, test := range testMatrix {
		// Set validator set weighted amount
		state := createStateFromMap(test.state, HostStakingDenom)
		suite.SetupAllowListedValSetAndDelegationState(state)

		// Create state
		givenCoin := sdk.NewInt64Coin(HostStakingDenom, test.given)
		expectedMap := map[string]int64{}
		for k, v := range test.expected {
			valAddress, _ := Bech32ifyValAddressBytes(types.CosmosValOperPrefix, sdk.ValAddress(k))
			expectedMap[valAddress] = v
		}
		allowlistedVals := suite.app.LSCosmosKeeper.GetAllowListedValidators(ctx)
		delegationState := suite.app.LSCosmosKeeper.GetDelegationState(ctx)

		// Run getIdealCurrentDelegations function with params

		valAmounts, err := keeper.FetchValidatorsToDelegate(allowlistedVals, delegationState, givenCoin)
		suite.Nil(err, "Error is not nil for validator to delegate")
		// Check outputs
		actualMap := map[string]int64{}
		for _, va := range valAmounts {
			actualMap[va.ValidatorAddr] = va.Amount.Amount.Int64()
		}
		suite.Equal(expectedMap, actualMap, "Matching val distribution")
	}
}

func (suite *IntegrationTestSuite) TestUndelegateDivideAmountIntoValidatorSet() {
	// Test setup
	app, ctx := suite.app, suite.ctx

	// Set Params
	denom := HostStakingDenom
	state := testStateData(denom)
	suite.SetupAllowListedValSetAndDelegationState(state)

	// Test data
	// Test data
	testMatrix := []struct {
		given    int64
		expected map[string]int64
	}{
		{
			given: 1000,
			expected: map[string]int64{
				"cosmosvalidatorAddr5": 1000,
			},
		},
		{
			given: 10000000,
			expected: map[string]int64{
				"cosmosvalidatorAddr5": 5000000,
				"cosmosvalidatorAddr1": 5000000,
			},
		},
		{
			given:    0,
			expected: map[string]int64{},
		},
		{
			given: 20000000,
			expected: map[string]int64{
				"cosmosvalidatorAddr1": 9000000,
				"cosmosvalidatorAddr2": 6000000,
				"cosmosvalidatorAddr5": 5000000,
			},
		},
		{
			given: 30000000,
			expected: map[string]int64{
				"cosmosvalidatorAddr1": 13000000,
				"cosmosvalidatorAddr2": 9000000,
				"cosmosvalidatorAddr3": 3000000,
				"cosmosvalidatorAddr5": 5000000,
			},
		},
		{
			given: 35000000,
			expected: map[string]int64{
				"cosmosvalidatorAddr1": 15000000,
				"cosmosvalidatorAddr2": 10000000,
				"cosmosvalidatorAddr3": 5000000,
				"cosmosvalidatorAddr5": 5000000,
			},
		},
	}

	// Create input parameters

	for _, test := range testMatrix {
		givenCoin := sdk.NewInt64Coin(HostStakingDenom, test.given)
		expectedMap := map[string]int64{}

		for k, v := range test.expected {
			valAddress, _ := Bech32ifyValAddressBytes(types.CosmosValOperPrefix, sdk.ValAddress(k))
			expectedMap[valAddress] = v
		}

		allowlistedVals := app.LSCosmosKeeper.GetAllowListedValidators(ctx)
		delegationState := app.LSCosmosKeeper.GetDelegationState(ctx)

		// Run getIdealCurrentDelegations function with params
		valAmounts, err := keeper.FetchValidatorsToUndelegate(allowlistedVals, delegationState, givenCoin)
		suite.Nil(err, "Error is not nil for validator to delegate")

		// Check outputs
		actualMap := map[string]int64{}
		for _, va := range valAmounts {
			actualMap[va.ValidatorAddr] = va.Amount.Amount.Int64()
		}
		suite.Equal(expectedMap, actualMap, "Matching val distribution")

	}

}

func (suite *IntegrationTestSuite) TestDivideAmountIntoValidatorSetEqual() {
	_, ctx := suite.app, suite.ctx

	denom := HostStakingDenom
	state := testStateDataEqual(denom)
	suite.SetupAllowListedValSetAndDelegationState(state)

	// Test data
	testMatrix := []struct {
		given    int64
		expected map[string]int64
	}{
		{
			given: 50000000,
			expected: map[string]int64{
				"cosmosvalidatorAddr1": 10000000,
				"cosmosvalidatorAddr2": 10000000,
				"cosmosvalidatorAddr3": 10000000,
				"cosmosvalidatorAddr4": 10000000,
				"cosmosvalidatorAddr5": 10000000,
			},
		},
	}

	for _, test := range testMatrix {
		givenCoin := sdk.NewInt64Coin(HostStakingDenom, test.given)
		expectedMap := map[string]int64{}

		for k, v := range test.expected {
			valAddress, _ := Bech32ifyValAddressBytes(types.CosmosValOperPrefix, sdk.ValAddress(k))
			expectedMap[valAddress] = v
		}

		allowlistedVals := suite.app.LSCosmosKeeper.GetAllowListedValidators(ctx)
		delegationState := suite.app.LSCosmosKeeper.GetDelegationState(ctx)

		// Run getIdealCurrentDelegations function with params
		valAmounts, err := keeper.FetchValidatorsToDelegate(allowlistedVals, delegationState, givenCoin)
		suite.Nil(err, "Error is not nil for validator to delegate")

		// sort the val address amount based on address to avoid generating different lists
		// by all validators
		sort.Sort(valAmounts)

		for i := 1; i < len(valAmounts); i++ {
			if valAmounts[i-1].ValidatorAddr > valAmounts[i].ValidatorAddr {
				panic("not sorted by string")
			}
		}

		// Check outputs
		actualMap := map[string]int64{}
		for _, va := range valAmounts {
			actualMap[va.ValidatorAddr] = va.Amount.Amount.Int64()
		}
		suite.Equal(expectedMap, actualMap, "Matching val distribution")

	}
}

func TestGetIdealCurrentDelegations(t *testing.T) {
	denom := HostStakingDenom

	type testValState struct {
		name   string
		weight string
		amount int64
	}

	testMatrix := []struct {
		amount               int64
		givenValset          types.AllowListedValidators
		givenDelegationState types.DelegationState
		expected             []testValState
	}{
		{
			amount: 5000000,
			givenValset: types.AllowListedValidators{
				AllowListedValidators: []types.AllowListedValidator{
					{"cosmosvalidatorAddr1", sdk.MustNewDecFromStr("0.5")},
					{"cosmosvalidatorAddr2", sdk.MustNewDecFromStr("0.5")}},
			},
			givenDelegationState: types.DelegationState{
				HostDelegationAccountBalance: sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(0))),
				HostChainDelegationAddress:   "cosmosdelegationAddr1",
				HostAccountDelegations: []types.HostAccountDelegation{
					{"cosmosvalidatorAddr1", sdk.NewCoin(denom, sdk.NewInt(10000000))}, {"cosmosvalidatorAddr2",
						sdk.NewCoin(denom, sdk.NewInt(5000000))}},
			},
			expected: []testValState{
				{"cosmosvalidatorAddr1", "", 0},
				{"cosmosvalidatorAddr2", "", 5000000},
			},
		},
		{
			amount: -5000000,
			givenValset: types.AllowListedValidators{
				AllowListedValidators: []types.AllowListedValidator{
					{"cosmosvalidatorAddr1", sdk.MustNewDecFromStr("0.5")},
					{"cosmosvalidatorAddr2", sdk.MustNewDecFromStr("0.5")}},
			},
			givenDelegationState: types.DelegationState{
				HostDelegationAccountBalance: sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(0))),
				HostChainDelegationAddress:   "cosmosdelegationAddr1",
				HostAccountDelegations: []types.HostAccountDelegation{
					{"cosmosvalidatorAddr1", sdk.NewCoin(denom, sdk.NewInt(10000000))},
					{"cosmosvalidatorAddr2", sdk.NewCoin(denom, sdk.NewInt(5000000))}},
			},
			expected: []testValState{
				{"cosmosvalidatorAddr1", "", 5000000},
				{"cosmosvalidatorAddr2", "", 0},
			},
		},
		{
			amount: 5000000,
			givenValset: types.AllowListedValidators{
				AllowListedValidators: []types.AllowListedValidator{
					{"cosmosvalidatorAddr1", sdk.MustNewDecFromStr("0.9")},
					{"cosmosvalidatorAddr2", sdk.MustNewDecFromStr("0.1")}},
			},
			givenDelegationState: types.DelegationState{
				HostDelegationAccountBalance: sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(0))),
				HostChainDelegationAddress:   "cosmosdelegationAddr1",
				HostAccountDelegations:       []types.HostAccountDelegation{{"cosmosvalidatorAddr1", sdk.NewCoin(denom, sdk.NewInt(10000000))}, {"cosmosvalidatorAddr2", sdk.NewCoin(denom, sdk.NewInt(40000000))}},
			},
			expected: []testValState{
				{"cosmosvalidatorAddr1", "", 39500000},
				{"cosmosvalidatorAddr2", "", -34500000},
			},
		},
		{
			amount: 0,
			givenValset: types.AllowListedValidators{
				AllowListedValidators: []types.AllowListedValidator{
					{"cosmosvalidatorAddr1", sdk.MustNewDecFromStr("0.5")},
					{"cosmosvalidatorAddr2", sdk.MustNewDecFromStr("0.5")}},
			},
			givenDelegationState: types.DelegationState{
				HostDelegationAccountBalance: sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(0))),
				HostChainDelegationAddress:   "cosmosdelegationAddr1",
				HostAccountDelegations:       []types.HostAccountDelegation{{"cosmosvalidatorAddr1", sdk.NewCoin(denom, sdk.NewInt(10000000))}, {"cosmosvalidatorAddr2", sdk.NewCoin(denom, sdk.NewInt(10000000))}},
			},
			expected: []testValState{
				{"cosmosvalidatorAddr1", "", 0},
				{"cosmosvalidatorAddr2", "", 0},
			},
		},
		{
			amount: 30000000,
			givenValset: types.AllowListedValidators{
				AllowListedValidators: []types.AllowListedValidator{
					{"cosmosvalidatorAddr1", sdk.MustNewDecFromStr("0.5")},
					{"cosmosvalidatorAddr2", sdk.MustNewDecFromStr("0.3")},
					{"cosmosvalidatorAddr3", sdk.MustNewDecFromStr("0.2")}},
			},
			givenDelegationState: types.DelegationState{
				HostDelegationAccountBalance: sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(0))),
				HostChainDelegationAddress:   "cosmosdelegationAddr1",
				HostAccountDelegations:       []types.HostAccountDelegation{{"cosmosvalidatorAddr1", sdk.NewCoin(denom, sdk.NewInt(10000000))}, {"cosmosvalidatorAddr2", sdk.NewCoin(denom, sdk.NewInt(10000000))}, {"cosmosvalidatorAddr3", sdk.NewCoin(denom, sdk.NewInt(10000000))}},
			},
			expected: []testValState{
				{"cosmosvalidatorAddr1", "", 20000000},
				{"cosmosvalidatorAddr2", "", 8000000},
				{"cosmosvalidatorAddr3", "", 2000000},
			},
		},
		{
			amount: -10000000,
			givenValset: types.AllowListedValidators{
				AllowListedValidators: []types.AllowListedValidator{
					{"cosmosvalidatorAddr1", sdk.MustNewDecFromStr("0.5")},
					{"cosmosvalidatorAddr2", sdk.MustNewDecFromStr("0.3")},
					{"cosmosvalidatorAddr3", sdk.MustNewDecFromStr("0.2")}},
			},
			givenDelegationState: types.DelegationState{
				HostDelegationAccountBalance: sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(0))),
				HostChainDelegationAddress:   "cosmosdelegationAddr1",
				HostAccountDelegations:       []types.HostAccountDelegation{{"cosmosvalidatorAddr1", sdk.NewCoin(denom, sdk.NewInt(10000000))}, {"cosmosvalidatorAddr2", sdk.NewCoin(denom, sdk.NewInt(10000000))}, {"cosmosvalidatorAddr3", sdk.NewCoin(denom, sdk.NewInt(10000000))}},
			},
			expected: []testValState{
				{"cosmosvalidatorAddr1", "", 0},
				{"cosmosvalidatorAddr2", "", 4000000},
				{"cosmosvalidatorAddr3", "", 6000000},
			},
		},
		{
			amount: -20000000,
			givenValset: types.AllowListedValidators{
				AllowListedValidators: []types.AllowListedValidator{
					{"cosmosvalidatorAddr1", sdk.MustNewDecFromStr("0.5")},
					{"cosmosvalidatorAddr2", sdk.MustNewDecFromStr("0.3")},
					{"cosmosvalidatorAddr3", sdk.MustNewDecFromStr("0.2")}},
			},
			givenDelegationState: types.DelegationState{
				HostDelegationAccountBalance: sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(0))),
				HostChainDelegationAddress:   "cosmosdelegationAddr1",
				HostAccountDelegations:       []types.HostAccountDelegation{{"cosmosvalidatorAddr1", sdk.NewCoin(denom, sdk.NewInt(10000000))}, {"cosmosvalidatorAddr2", sdk.NewCoin(denom, sdk.NewInt(10000000))}, {"cosmosvalidatorAddr3", sdk.NewCoin(denom, sdk.NewInt(10000000))}},
			},
			expected: []testValState{
				{"cosmosvalidatorAddr1", "", 5000000},
				{"cosmosvalidatorAddr2", "", 7000000},
				{"cosmosvalidatorAddr3", "", 8000000},
			},
		},
		{
			amount: 10000000,
			givenValset: types.AllowListedValidators{
				AllowListedValidators: []types.AllowListedValidator{
					{"cosmosvalidatorAddr1", sdk.MustNewDecFromStr("0.5")},
					{"cosmosvalidatorAddr2", sdk.MustNewDecFromStr("0.5")},
					{"cosmosvalidatorAddr3", sdk.MustNewDecFromStr("0")}},
			},
			givenDelegationState: types.DelegationState{
				HostDelegationAccountBalance: sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(0))),
				HostChainDelegationAddress:   "cosmosdelegationAddr1",
				HostAccountDelegations:       []types.HostAccountDelegation{{"cosmosvalidatorAddr1", sdk.NewCoin(denom, sdk.NewInt(10000000))}, {"cosmosvalidatorAddr2", sdk.NewCoin(denom, sdk.NewInt(10000000))}, {"cosmosvalidatorAddr3", sdk.NewCoin(denom, sdk.NewInt(10000000))}},
			},
			expected: []testValState{
				{"cosmosvalidatorAddr1", "", 10000000},
				{"cosmosvalidatorAddr2", "", 10000000},
				{"cosmosvalidatorAddr3", "", -10000000},
			},
		},
		{
			amount: 10000000,
			givenValset: types.AllowListedValidators{
				AllowListedValidators: []types.AllowListedValidator{
					{"cosmosvalidatorAddr1", sdk.MustNewDecFromStr("0.5")},
					{"cosmosvalidatorAddr2", sdk.MustNewDecFromStr("0.4")},
					{"cosmosvalidatorAddr3", sdk.MustNewDecFromStr("0")},
					{"cosmosvalidatorAddr4", sdk.MustNewDecFromStr("0.1")}},
			},
			givenDelegationState: types.DelegationState{
				HostDelegationAccountBalance: sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(0))),
				HostChainDelegationAddress:   "cosmosdelegationAddr1",
				HostAccountDelegations:       []types.HostAccountDelegation{{"cosmosvalidatorAddr1", sdk.NewCoin(denom, sdk.NewInt(10000000))}, {"cosmosvalidatorAddr2", sdk.NewCoin(denom, sdk.NewInt(10000000))}, {"cosmosvalidatorAddr3", sdk.NewCoin(denom, sdk.NewInt(10000000))}, {"cosmosvalidatorAddr4", sdk.NewCoin(denom, sdk.NewInt(10000000))}},
			},
			expected: []testValState{
				{"cosmosvalidatorAddr1", "", 15000000},
				{"cosmosvalidatorAddr2", "", 10000000},
				{"cosmosvalidatorAddr3", "", -10000000},
				{"cosmosvalidatorAddr4", "", -5000000},
			},
		},
	}
	for i, test := range testMatrix {
		// Create validator state
		givenState := types.WeightedAddressAmounts{}
		delegationMap := types.GetHostAccountDelegationMap(test.givenDelegationState.HostAccountDelegations)
		expectedMap := map[string]types.WeightedAddressAmount{}
		for i := 0; i < len(test.givenValset.AllowListedValidators); i++ {
			givenState = append(givenState, types.WeightedAddressAmount{
				Address: test.givenValset.AllowListedValidators[i].ValidatorAddress,
				Weight:  test.givenValset.AllowListedValidators[i].TargetWeight,
				Denom:   denom,
				Amount:  delegationMap[test.givenValset.AllowListedValidators[i].ValidatorAddress].Amount,
			})
			expectedMap[test.expected[i].name] = types.WeightedAddressAmount{
				Address: test.expected[i].name,
				Denom:   denom,
				Amount:  sdk.NewInt(test.expected[i].amount),
			}
		}
		// Call getIdealCurrentDelegations function with params
		state, err := keeper.GetIdealCurrentDelegations(test.givenValset, test.givenDelegationState, sdk.NewInt64Coin(denom, int64(math.Abs(float64(test.amount)))), !(test.amount > 0))
		assert.NoError(t, err)

		// Assert state
		for j, s := range state {
			expected, ok := expectedMap[s.Address]
			assert.True(t, ok, "Address not is expected list")
			failMsg := fmt.Sprintf("Amounts should be same. Failed for %d case: %d", i, j)
			assert.Equal(t, expected.Amount.BigInt(), s.Amount.BigInt(), failMsg)
		}
	}
}

func TestNegativeCoin(t *testing.T) {
	coinFunc := func() {
		sdk.NewCoin(HostStakingDenom, sdk.NewInt(-1000))
	}

	assert.Panics(t, coinFunc)
}

func TestMulInt(t *testing.T) {
	w, _ := sdk.NewDecFromStr("0.5")

	a := sdk.NewInt(1000000)

	assert.Equal(t, sdk.NewInt(500000).Int64(), w.Mul(sdk.NewDecFromInt(a)).TruncateInt().Int64())
	assert.Equal(t, sdk.NewInt(0).Int64(), w.Mul(sdk.NewDecFromInt(a.Add(sdk.NewInt(1000000)))).TruncateInt().SubRaw(1000000).Int64())
}

func testStateData(denom string) types.WeightedAddressAmounts {
	testStruct := []struct {
		name   string
		weight string
		amount int64
	}{
		{
			name:   "cosmosvalidatorAddr1",
			weight: "0.4",
			amount: 15000000, // ideal: 14000000
		},
		{
			name:   "cosmosvalidatorAddr2",
			weight: "0.2",
			amount: 10000000, // ideal: 7000000
		},
		{
			name:   "cosmosvalidatorAddr3",
			weight: "0.3",
			amount: 5000000, // ideal: 10500000
		},
		{
			name:   "cosmosvalidatorAddr4",
			weight: "0.1",
			amount: 0, // ideal: 3500000
		},
		{
			name:   "cosmosvalidatorAddr5",
			weight: "0",
			amount: 5000000, // ideal: 0
		},
	}
	// Create state
	state := types.WeightedAddressAmounts{}
	for _, ts := range testStruct {
		weight, _ := sdk.NewDecFromStr(ts.weight)
		valAddress, _ := Bech32ifyValAddressBytes(types.CosmosValOperPrefix, sdk.ValAddress(ts.name))
		state = append(state, types.WeightedAddressAmount{
			Weight:  weight,
			Amount:  sdk.NewInt(ts.amount),
			Address: valAddress,
			Denom:   denom,
		})
	}
	return state
}

func Bech32ifyValAddressBytes(prefix string, address sdk.ValAddress) (string, error) {
	if address.Empty() {
		return "", nil
	}
	if len(address.Bytes()) == 0 {
		return "", nil
	}
	if len(prefix) == 0 {
		return "", errors.New("prefix cannot be empty")
	}
	return bech32.ConvertAndEncode(prefix, address.Bytes())
}

func createStateFromMap(stateMap map[string][]string, denom string) types.WeightedAddressAmounts {
	// Create state
	state := types.WeightedAddressAmounts{}
	for addr, wa := range stateMap {
		amt, _ := sdk.NewIntFromString(wa[0])
		weight, _ := sdk.NewDecFromStr(wa[1])
		valAddress, _ := Bech32ifyValAddressBytes(types.CosmosValOperPrefix, sdk.ValAddress(addr))
		state = append(state, types.WeightedAddressAmount{
			Weight:  weight,
			Amount:  amt,
			Address: valAddress,
			Denom:   denom,
		})
	}
	return state
}

func testStateDataEqual(denom string) types.WeightedAddressAmounts {
	testStruct := []struct {
		name   string
		weight string
		amount int64
	}{
		{
			name:   "cosmosvalidatorAddr1",
			weight: "0.2",
			amount: 15000000, // ideal: 14000000
		},
		{
			name:   "cosmosvalidatorAddr2",
			weight: "0.2",
			amount: 15000000, // ideal: 7000000
		},
		{
			name:   "cosmosvalidatorAddr3",
			weight: "0.2",
			amount: 15000000, // ideal: 10500000
		},
		{
			name:   "cosmosvalidatorAddr4",
			weight: "0.2",
			amount: 15000000, // ideal: 3500000
		},
		{
			name:   "cosmosvalidatorAddr5",
			weight: "0.2",
			amount: 15000000, // ideal: 0
		},
	}
	// Create state
	state := types.WeightedAddressAmounts{}
	for _, ts := range testStruct {
		weight, _ := sdk.NewDecFromStr(ts.weight)
		valAddress, _ := Bech32ifyValAddressBytes(types.CosmosValOperPrefix, sdk.ValAddress(ts.name))
		state = append(state, types.WeightedAddressAmount{
			Weight:  weight,
			Amount:  sdk.NewInt(ts.amount),
			Address: valAddress,
			Denom:   denom,
		})
	}
	return state
}

func (suite *IntegrationTestSuite) TestGetAllValidatorsState() {
	app, ctx := suite.app, suite.ctx

	k := app.LSCosmosKeeper

	hostChainParams := k.GetHostChainParams(ctx)

	allowListedValidatorsSet := types.AllowListedValidators{
		AllowListedValidators: []types.AllowListedValidator{
			{
				ValidatorAddress: "cosmosvaloper10vcqjzphfdlumas0vp64f0hruhrqxv0cd7wdy2",
				TargetWeight:     sdk.NewDecWithPrec(3, 1),
			},
			{
				ValidatorAddress: "cosmosvaloper10khgeppewe4rgfrcy809r9h00aquwxxxgwgwa5",
				TargetWeight:     sdk.NewDecWithPrec(7, 1),
			},
		},
	}

	k.SetAllowListedValidators(ctx, allowListedValidatorsSet)

	delegationState := types.DelegationState{
		HostAccountDelegations: []types.HostAccountDelegation{
			{
				ValidatorAddress: "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt",
				Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 200),
			},
			{
				ValidatorAddress: "cosmosvaloper1lcck2cxh7dzgkrfk53kysg9ktdrsjj6jfwlnm2",
				Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 100),
			},
			{
				ValidatorAddress: "cosmosvaloper10khgeppewe4rgfrcy809r9h00aquwxxxgwgwa5",
				Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 400),
			},
		},
	}
	k.SetDelegationState(ctx, delegationState)

	// fetch a combined updated val set list and delegation state
	updateValList, hostAccountDelegations := k.GetAllValidatorsState(ctx, hostChainParams.BaseDenom)

	// sort both updatedValList and hostAccountDelegations
	sort.Sort(updateValList)
	sort.Sort(hostAccountDelegations)

	// get the current delegation state and
	// assign the updated validator delegation state to the current delegation state
	delegationStateS := k.GetDelegationState(ctx)
	delegationStateS.HostAccountDelegations = hostAccountDelegations

	allowListerValidators := types.AllowListedValidators{AllowListedValidators: updateValList}

	list, err := keeper.FetchValidatorsToUndelegate(allowListerValidators, delegationStateS, sdk.NewInt64Coin(hostChainParams.BaseDenom, 600))
	suite.NoError(err)
	suite.Equal(sdk.NewInt64Coin(hostChainParams.BaseDenom, 200), list[0].Amount)
	suite.Equal(sdk.NewInt64Coin(hostChainParams.BaseDenom, 100), list[1].Amount)
	suite.Equal(sdk.NewInt64Coin(hostChainParams.BaseDenom, 300), list[2].Amount)

	list, err = keeper.FetchValidatorsToDelegate(allowListerValidators, delegationStateS, sdk.NewInt64Coin(hostChainParams.BaseDenom, 2000))
	suite.NoError(err)
	suite.Equal(sdk.NewInt64Coin(hostChainParams.BaseDenom, 1490), list[0].Amount)
	suite.Equal(sdk.NewInt64Coin(hostChainParams.BaseDenom, 510), list[1].Amount)

	delegationState = types.DelegationState{
		HostAccountDelegations: []types.HostAccountDelegation{
			{
				ValidatorAddress: "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt",
				Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 0),
			},
			{
				ValidatorAddress: "cosmosvaloper1lcck2cxh7dzgkrfk53kysg9ktdrsjj6jfwlnm2",
				Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 0),
			},
			{
				ValidatorAddress: "cosmosvaloper10khgeppewe4rgfrcy809r9h00aquwxxxgwgwa5",
				Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1890),
			},
			{
				ValidatorAddress: "cosmosvaloper10vcqjzphfdlumas0vp64f0hruhrqxv0cd7wdy2",
				Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 510),
			},
		},
	}
	k.SetDelegationState(ctx, delegationState)

	// fetch a combined updated val set list and delegation state
	updateValList, hostAccountDelegations = k.GetAllValidatorsState(ctx, hostChainParams.BaseDenom)

	// sort both updatedValList and hostAccountDelegations
	sort.Sort(updateValList)
	sort.Sort(hostAccountDelegations)

	// get the current delegation state and
	// assign the updated validator delegation state to the current delegation state
	delegationStateS = k.GetDelegationState(ctx)
	delegationStateS.HostAccountDelegations = hostAccountDelegations

	allowListerValidators = types.AllowListedValidators{AllowListedValidators: updateValList}

	list, err = keeper.FetchValidatorsToDelegate(allowListerValidators, delegationStateS, sdk.NewInt64Coin(hostChainParams.BaseDenom, 0))
	suite.NoError(err)
	suite.Equal(0, len(list))
}

func (suite *IntegrationTestSuite) TestDelegateAndUndelegate() {
	app, ctx := suite.app, suite.ctx

	k := app.LSCosmosKeeper

	hostChainParams := k.GetHostChainParams(ctx)

	testCases := []struct {
		types.AllowListedValidators
		types.DelegationState
		DelegateAmounts                        []sdk.Coin
		UndelegateAmounts                      []sdk.Coin
		ExpectedListWithDelegateDistribution   []types.ValAddressAmounts
		ExpectedListWithUndelegateDistribution []types.ValAddressAmounts
	}{
		{
			AllowListedValidators: types.AllowListedValidators{
				AllowListedValidators: []types.AllowListedValidator{
					{
						ValidatorAddress: "cosmosvaloper10vcqjzphfdlumas0vp64f0hruhrqxv0cd7wdy2",
						TargetWeight:     sdk.NewDecWithPrec(3, 1),
					},
					{
						ValidatorAddress: "cosmosvaloper10khgeppewe4rgfrcy809r9h00aquwxxxgwgwa5",
						TargetWeight:     sdk.NewDecWithPrec(7, 1),
					},
				},
			},
			DelegationState: types.DelegationState{
				HostAccountDelegations: []types.HostAccountDelegation{
					{
						ValidatorAddress: "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 200),
					},
					{
						ValidatorAddress: "cosmosvaloper1lcck2cxh7dzgkrfk53kysg9ktdrsjj6jfwlnm2",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 100),
					},
					{
						ValidatorAddress: "cosmosvaloper10khgeppewe4rgfrcy809r9h00aquwxxxgwgwa5",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 400),
					},
				},
			},
			DelegateAmounts: []sdk.Coin{
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 1),
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 2),
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 3),
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 40),
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 222),
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 223),
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 4000),
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 4001),
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 4002),
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 4003),
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 4004),
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 4005),
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 4006),
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 4007),
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 4008),
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 4009),
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 4010),
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 4011),
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 1000000000),
			},
			UndelegateAmounts: []sdk.Coin{
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 1),
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 2),
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 3),
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 40),
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 222),
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 223),
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 400),
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 500),
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 700),
			},
			ExpectedListWithDelegateDistribution: []types.ValAddressAmounts{
				{types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10vcqjzphfdlumas0vp64f0hruhrqxv0cd7wdy2", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))}},
				{types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10vcqjzphfdlumas0vp64f0hruhrqxv0cd7wdy2", Amount: sdk.NewCoin("uatom", sdk.NewInt(2))}},
				{types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10vcqjzphfdlumas0vp64f0hruhrqxv0cd7wdy2", Amount: sdk.NewCoin("uatom", sdk.NewInt(3))}},
				{types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10vcqjzphfdlumas0vp64f0hruhrqxv0cd7wdy2", Amount: sdk.NewCoin("uatom", sdk.NewInt(40))}},
				{types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10vcqjzphfdlumas0vp64f0hruhrqxv0cd7wdy2", Amount: sdk.NewCoin("uatom", sdk.NewInt(222))}},
				{types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10vcqjzphfdlumas0vp64f0hruhrqxv0cd7wdy2", Amount: sdk.NewCoin("uatom", sdk.NewInt(223))}},
				{
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10khgeppewe4rgfrcy809r9h00aquwxxxgwgwa5", Amount: sdk.NewCoin("uatom", sdk.NewInt(2890))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10vcqjzphfdlumas0vp64f0hruhrqxv0cd7wdy2", Amount: sdk.NewCoin("uatom", sdk.NewInt(1110))},
				},
				{
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10khgeppewe4rgfrcy809r9h00aquwxxxgwgwa5", Amount: sdk.NewCoin("uatom", sdk.NewInt(2890))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10vcqjzphfdlumas0vp64f0hruhrqxv0cd7wdy2", Amount: sdk.NewCoin("uatom", sdk.NewInt(1111))},
				},
				{
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10khgeppewe4rgfrcy809r9h00aquwxxxgwgwa5", Amount: sdk.NewCoin("uatom", sdk.NewInt(2891))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10vcqjzphfdlumas0vp64f0hruhrqxv0cd7wdy2", Amount: sdk.NewCoin("uatom", sdk.NewInt(1111))},
				},
				{
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10khgeppewe4rgfrcy809r9h00aquwxxxgwgwa5", Amount: sdk.NewCoin("uatom", sdk.NewInt(2892))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10vcqjzphfdlumas0vp64f0hruhrqxv0cd7wdy2", Amount: sdk.NewCoin("uatom", sdk.NewInt(1111))},
				},
				{
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10khgeppewe4rgfrcy809r9h00aquwxxxgwgwa5", Amount: sdk.NewCoin("uatom", sdk.NewInt(2892))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10vcqjzphfdlumas0vp64f0hruhrqxv0cd7wdy2", Amount: sdk.NewCoin("uatom", sdk.NewInt(1112))},
				},
				{
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10khgeppewe4rgfrcy809r9h00aquwxxxgwgwa5", Amount: sdk.NewCoin("uatom", sdk.NewInt(2893))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10vcqjzphfdlumas0vp64f0hruhrqxv0cd7wdy2", Amount: sdk.NewCoin("uatom", sdk.NewInt(1112))},
				},
				{
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10khgeppewe4rgfrcy809r9h00aquwxxxgwgwa5", Amount: sdk.NewCoin("uatom", sdk.NewInt(2894))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10vcqjzphfdlumas0vp64f0hruhrqxv0cd7wdy2", Amount: sdk.NewCoin("uatom", sdk.NewInt(1112))},
				},
				{
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10khgeppewe4rgfrcy809r9h00aquwxxxgwgwa5", Amount: sdk.NewCoin("uatom", sdk.NewInt(2894))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10vcqjzphfdlumas0vp64f0hruhrqxv0cd7wdy2", Amount: sdk.NewCoin("uatom", sdk.NewInt(1113))},
				},
				{
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10khgeppewe4rgfrcy809r9h00aquwxxxgwgwa5", Amount: sdk.NewCoin("uatom", sdk.NewInt(2895))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10vcqjzphfdlumas0vp64f0hruhrqxv0cd7wdy2", Amount: sdk.NewCoin("uatom", sdk.NewInt(1113))},
				},
				{
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10khgeppewe4rgfrcy809r9h00aquwxxxgwgwa5", Amount: sdk.NewCoin("uatom", sdk.NewInt(2896))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10vcqjzphfdlumas0vp64f0hruhrqxv0cd7wdy2", Amount: sdk.NewCoin("uatom", sdk.NewInt(1113))},
				},
				{
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10khgeppewe4rgfrcy809r9h00aquwxxxgwgwa5", Amount: sdk.NewCoin("uatom", sdk.NewInt(2897))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10vcqjzphfdlumas0vp64f0hruhrqxv0cd7wdy2", Amount: sdk.NewCoin("uatom", sdk.NewInt(1113))},
				},
				{
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10khgeppewe4rgfrcy809r9h00aquwxxxgwgwa5", Amount: sdk.NewCoin("uatom", sdk.NewInt(2897))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10vcqjzphfdlumas0vp64f0hruhrqxv0cd7wdy2", Amount: sdk.NewCoin("uatom", sdk.NewInt(1114))},
				},
				{
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10khgeppewe4rgfrcy809r9h00aquwxxxgwgwa5", Amount: sdk.NewCoin("uatom", sdk.NewInt(700000090))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10vcqjzphfdlumas0vp64f0hruhrqxv0cd7wdy2", Amount: sdk.NewCoin("uatom", sdk.NewInt(299999910))},
				},
			},
			ExpectedListWithUndelegateDistribution: []types.ValAddressAmounts{
				{types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))}},
				{types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt", Amount: sdk.NewCoin("uatom", sdk.NewInt(2))}},
				{types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt", Amount: sdk.NewCoin("uatom", sdk.NewInt(3))}},
				{types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt", Amount: sdk.NewCoin("uatom", sdk.NewInt(40))}},
				{types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt", Amount: sdk.NewCoin("uatom", sdk.NewInt(200))}, types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1lcck2cxh7dzgkrfk53kysg9ktdrsjj6jfwlnm2", Amount: sdk.NewCoin("uatom", sdk.NewInt(22))}},
				{types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt", Amount: sdk.NewCoin("uatom", sdk.NewInt(200))}, types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1lcck2cxh7dzgkrfk53kysg9ktdrsjj6jfwlnm2", Amount: sdk.NewCoin("uatom", sdk.NewInt(23))}},
				{types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt", Amount: sdk.NewCoin("uatom", sdk.NewInt(200))}, types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1lcck2cxh7dzgkrfk53kysg9ktdrsjj6jfwlnm2", Amount: sdk.NewCoin("uatom", sdk.NewInt(100))}, types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10khgeppewe4rgfrcy809r9h00aquwxxxgwgwa5", Amount: sdk.NewCoin("uatom", sdk.NewInt(100))}},
				{types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt", Amount: sdk.NewCoin("uatom", sdk.NewInt(200))}, types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1lcck2cxh7dzgkrfk53kysg9ktdrsjj6jfwlnm2", Amount: sdk.NewCoin("uatom", sdk.NewInt(100))}, types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10khgeppewe4rgfrcy809r9h00aquwxxxgwgwa5", Amount: sdk.NewCoin("uatom", sdk.NewInt(200))}},
				{types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt", Amount: sdk.NewCoin("uatom", sdk.NewInt(200))}, types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1lcck2cxh7dzgkrfk53kysg9ktdrsjj6jfwlnm2", Amount: sdk.NewCoin("uatom", sdk.NewInt(100))}, types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10khgeppewe4rgfrcy809r9h00aquwxxxgwgwa5", Amount: sdk.NewCoin("uatom", sdk.NewInt(400))}},
			},
		},
		{
			AllowListedValidators: types.AllowListedValidators{
				AllowListedValidators: []types.AllowListedValidator{
					{
						ValidatorAddress: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0",
						TargetWeight:     sdk.NewDecWithPrec(16131000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper196ax4vc0lwpxndu9dyhvca7jhxp70rmcvrj90c",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper1v5y0tg0jllvxf5c3afml8s3awue0ymju89frut",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper1qaa9zej9a0ge3ugpx3pxyx602lxh3ztqgfnp42",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper1tflk30mq5vgqjdly92kkhhq3raev2hnz6eete3",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper1hjct6q7npsspsg3dgvzk3sdf89spmlpfdn6m9d",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper1clpqr4nrk4khgkxj78fcwwh6dl3uw4epsluffn",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper132juzk0gdmwuxvx4phug7m3ymyatxlh9734g4w",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper15urq2dtp9qce4fyc85m6upwm9xul3049e02707",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper1vf44d85es37hwl9f4h9gv0e064m0lla60j9luj",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper1lzhlnpahvznwfv4jmay2tgaha5kmz5qxerarrl",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper1g48268mu5vfp4wk7dk89r0wdrakm9p5xk0q50k",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper16k579jk6yt2cwmqx9dz5xvq9fug2tekvlu9qdv",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper1gpx52r9h3zeul45amvcy2pysgvcwddxrgx6cnv",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper1grgelyng2v6v3t8z87wu3sxgt9m5s03xfytvz7",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper14kn0kk33szpwus9nh8n87fjel8djx0y070ymmj",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper1vvwtk805lxehwle9l4yudmq6mn0g32px9xtkhc",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper1rpgtz9pskr5geavkjz02caqmeep7cwwpv73axj",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper1y0us8xvsvfvqkk9c6nt5cfyu5au5tww2ztve7q",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper130mdu9a0etmeuw52qfxk73pn0ga6gawkxsrlwf",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper10e4vsut6suau8tk9m6dnrm0slgd6npe3jx5xpv",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper1kgddca7qj96z0qcxr2c45z73cfl0c75p7f3s2e",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper199mlc7fr6ll5t54w7tts7f4s0cvnqgc59nmuxf",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper1k2d9ed9vgfuk2m58a2d80q9u6qljkh4vfaqjfq",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper1uutuwrwt3z2a5z8z3uasml3rftlpmu25aga5c6",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper1ssm0d433seakyak8kcf93yefhknjleeds4y3em",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper1ehkfl7palwrh6w2hhr2yfrgrq8jetgucudztfe",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper1n3mhyp9fvcmuu8l0q8qvjy07x0rql8q46fe2xk",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper1gxju9ky3hwxvqqagrl3dxtl49kjpxq6wlqe6m5",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper157v7tczs40axfgejp2m43kwuzqe0wsy0rv8puv",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper14l0fp639yudfl46zauvv8rkzjgd4u0zk2aseys",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper1lktjhnzkpkz3ehrg8psvmwhafg56kfss3q3t8m",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper1udpsgkgyutgsglauk9vk9rs03a3skc62gup9ny",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper124maqmcqv8tquy764ktz7cu0gxnzfw54n3vww8",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper1uhnsxv6m83jj3328mhrql7yax3nge5svrv6t6c",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper1083svrca4t350mphfv9x45wq9asrs60cdmrflj",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper17mggn4znyeyg25wd7498qxl7r2jhgue8u4qjcq",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper1pjmngrwcsatsuyy8m3qrunaun67sr9x7z5r2qs",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper1gjtvly9lel6zskvwtvlg5vhwpu9c9waw7sxzwx",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper1ptyzewnns2kn37ewtmv6ppsvhdnmeapvtfc9y5",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper1gf4wlkutql95j7wwsxz490s6fahlvk2s9xpwax",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper1m73mgwn3cm2e8x9a9axa0kw8nqz8a492ms63vn",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper1hdrlqvyjfy5sdrseecjrutyws9khtxxaux62l7",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper17h2x3j7u44qkrq0sk8ul0r2qr440rwgjkfg0gh",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper1gp957czryfgyvxwn3tfnyy2f0t9g2p4pqeemx8",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper10nzaaeh2kq28t3nqsh5m8kmyv90vx7ym5mpakx",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper1ddle9tczl87gsvmeva3c48nenyng4n56nghmjk",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper1fhr7e04ct0zslmkzqt9smakg3sxrdve6ulclj2",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper18extdhzzl5c8tr6453e5hzaj3exrdlea90fj3y",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper1e859xaue4k2jzqw20cv6l7p3tmc378pc3k8g2u",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper1jlr62guqwrwkdt4m3y00zh2rrsamhjf9num5xr",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper16yupepagywvlk7uhpfchtwa0stu5f8cyhh54f2",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper1485u80fdxjan4sd3esrvyw6cyurpvddvzuh48y",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper1fqzqejwkk898fcslw4z4eeqjzesynvrdfr5hte",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper1e0plfg475phrsvrlzw8gwppeva0zk5yg9fgg8c",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper1xwazl8ftks4gn00y5x3c47auquc62ssuqlj02r",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper1jmykcq8gylmy5tgqtel4xj4q62fdt49sl584xd",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper14qazscc80zgzx3m0m0aa30ths0p9hg8vdglqrc",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper13x77yexvf6qexfjg9czp6jhpv7vpjdwwkyhe4p",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper1jxv0u20scum4trha72c7ltfgfqef6nsch7q6cu",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
					{
						ValidatorAddress: "cosmosvaloper1ualhu3fjgg77g485gmyswkq3w0dp7gys6qzwrv",
						TargetWeight:     sdk.NewDecWithPrec(16129000000000000, 18),
					},
				},
			},
			DelegationState: types.DelegationState{
				HostAccountDelegations: []types.HostAccountDelegation{
					{
						ValidatorAddress: "cosmosvaloper1083svrca4t350mphfv9x45wq9asrs60cdmrflj",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper10e4vsut6suau8tk9m6dnrm0slgd6npe3jx5xpv",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339554),
					},
					{
						ValidatorAddress: "cosmosvaloper10nzaaeh2kq28t3nqsh5m8kmyv90vx7ym5mpakx",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper124maqmcqv8tquy764ktz7cu0gxnzfw54n3vww8",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper130mdu9a0etmeuw52qfxk73pn0ga6gawkxsrlwf",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper132juzk0gdmwuxvx4phug7m3ymyatxlh9734g4w",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper13x77yexvf6qexfjg9czp6jhpv7vpjdwwkyhe4p",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper1485u80fdxjan4sd3esrvyw6cyurpvddvzuh48y",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper14kn0kk33szpwus9nh8n87fjel8djx0y070ymmj",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper14l0fp639yudfl46zauvv8rkzjgd4u0zk2aseys",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper14qazscc80zgzx3m0m0aa30ths0p9hg8vdglqrc",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper157v7tczs40axfgejp2m43kwuzqe0wsy0rv8puv",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper15urq2dtp9qce4fyc85m6upwm9xul3049e02707",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper16k579jk6yt2cwmqx9dz5xvq9fug2tekvlu9qdv",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper16yupepagywvlk7uhpfchtwa0stu5f8cyhh54f2",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper17h2x3j7u44qkrq0sk8ul0r2qr440rwgjkfg0gh",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper17mggn4znyeyg25wd7498qxl7r2jhgue8u4qjcq",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper18extdhzzl5c8tr6453e5hzaj3exrdlea90fj3y",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper196ax4vc0lwpxndu9dyhvca7jhxp70rmcvrj90c",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper199mlc7fr6ll5t54w7tts7f4s0cvnqgc59nmuxf",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper1clpqr4nrk4khgkxj78fcwwh6dl3uw4epsluffn",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper1ddle9tczl87gsvmeva3c48nenyng4n56nghmjk",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper1e0plfg475phrsvrlzw8gwppeva0zk5yg9fgg8c",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper1e859xaue4k2jzqw20cv6l7p3tmc378pc3k8g2u",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper1ehkfl7palwrh6w2hhr2yfrgrq8jetgucudztfe",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper1fhr7e04ct0zslmkzqt9smakg3sxrdve6ulclj2",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper1fqzqejwkk898fcslw4z4eeqjzesynvrdfr5hte",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper1g48268mu5vfp4wk7dk89r0wdrakm9p5xk0q50k",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper1gf4wlkutql95j7wwsxz490s6fahlvk2s9xpwax",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper1gjtvly9lel6zskvwtvlg5vhwpu9c9waw7sxzwx",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper1gp957czryfgyvxwn3tfnyy2f0t9g2p4pqeemx8",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper1gpx52r9h3zeul45amvcy2pysgvcwddxrgx6cnv",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper1grgelyng2v6v3t8z87wu3sxgt9m5s03xfytvz7",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper1gxju9ky3hwxvqqagrl3dxtl49kjpxq6wlqe6m5",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper1hdrlqvyjfy5sdrseecjrutyws9khtxxaux62l7",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper1hjct6q7npsspsg3dgvzk3sdf89spmlpfdn6m9d",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper1jlr62guqwrwkdt4m3y00zh2rrsamhjf9num5xr",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper1jmykcq8gylmy5tgqtel4xj4q62fdt49sl584xd",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper1jxv0u20scum4trha72c7ltfgfqef6nsch7q6cu",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper1k2d9ed9vgfuk2m58a2d80q9u6qljkh4vfaqjfq",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper1kgddca7qj96z0qcxr2c45z73cfl0c75p7f3s2e",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper1lktjhnzkpkz3ehrg8psvmwhafg56kfss3q3t8m",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper1lzhlnpahvznwfv4jmay2tgaha5kmz5qxerarrl",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper1m73mgwn3cm2e8x9a9axa0kw8nqz8a492ms63vn",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper1n3mhyp9fvcmuu8l0q8qvjy07x0rql8q46fe2xk",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper1pjmngrwcsatsuyy8m3qrunaun67sr9x7z5r2qs",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper1ptyzewnns2kn37ewtmv6ppsvhdnmeapvtfc9y5",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper1qaa9zej9a0ge3ugpx3pxyx602lxh3ztqgfnp42",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper1rpgtz9pskr5geavkjz02caqmeep7cwwpv73axj",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339690),
					},
					{
						ValidatorAddress: "cosmosvaloper1ssm0d433seakyak8kcf93yefhknjleeds4y3em",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper1tflk30mq5vgqjdly92kkhhq3raev2hnz6eete3",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper1ualhu3fjgg77g485gmyswkq3w0dp7gys6qzwrv",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper1udpsgkgyutgsglauk9vk9rs03a3skc62gup9ny",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper1uhnsxv6m83jj3328mhrql7yax3nge5svrv6t6c",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper1uutuwrwt3z2a5z8z3uasml3rftlpmu25aga5c6",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper1v5y0tg0jllvxf5c3afml8s3awue0ymju89frut",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper1vf44d85es37hwl9f4h9gv0e064m0lla60j9luj",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper1vvwtk805lxehwle9l4yudmq6mn0g32px9xtkhc",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper1xwazl8ftks4gn00y5x3c47auquc62ssuqlj02r",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
					{
						ValidatorAddress: "cosmosvaloper1y0us8xvsvfvqkk9c6nt5cfyu5au5tww2ztve7q",
						Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 1339524),
					},
				},
			},
			DelegateAmounts: []sdk.Coin{
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 1),
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 2),
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 10),
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 40),
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 70),
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 222),
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 4000),
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 1000000000),
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 10000000000001),
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 1000000000100000),
			},
			UndelegateAmounts: []sdk.Coin{
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 1),
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 2),
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 10),
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 40),
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 70),
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 222),
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 4000),
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 83050660),
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 83050680),
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 83050683),
				sdk.NewInt64Coin(hostChainParams.BaseDenom, 83050684),
			},
			ExpectedListWithDelegateDistribution: []types.ValAddressAmounts{
				{
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1083svrca4t350mphfv9x45wq9asrs60cdmrflj", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10e4vsut6suau8tk9m6dnrm0slgd6npe3jx5xpv", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10nzaaeh2kq28t3nqsh5m8kmyv90vx7ym5mpakx", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper124maqmcqv8tquy764ktz7cu0gxnzfw54n3vww8", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper130mdu9a0etmeuw52qfxk73pn0ga6gawkxsrlwf", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper132juzk0gdmwuxvx4phug7m3ymyatxlh9734g4w", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper13x77yexvf6qexfjg9czp6jhpv7vpjdwwkyhe4p", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1485u80fdxjan4sd3esrvyw6cyurpvddvzuh48y", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14kn0kk33szpwus9nh8n87fjel8djx0y070ymmj", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14l0fp639yudfl46zauvv8rkzjgd4u0zk2aseys", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14qazscc80zgzx3m0m0aa30ths0p9hg8vdglqrc", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper157v7tczs40axfgejp2m43kwuzqe0wsy0rv8puv", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper15urq2dtp9qce4fyc85m6upwm9xul3049e02707", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper16k579jk6yt2cwmqx9dz5xvq9fug2tekvlu9qdv", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper16yupepagywvlk7uhpfchtwa0stu5f8cyhh54f2", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper17h2x3j7u44qkrq0sk8ul0r2qr440rwgjkfg0gh", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper17mggn4znyeyg25wd7498qxl7r2jhgue8u4qjcq", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper18extdhzzl5c8tr6453e5hzaj3exrdlea90fj3y", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper196ax4vc0lwpxndu9dyhvca7jhxp70rmcvrj90c", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper199mlc7fr6ll5t54w7tts7f4s0cvnqgc59nmuxf", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1clpqr4nrk4khgkxj78fcwwh6dl3uw4epsluffn", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ddle9tczl87gsvmeva3c48nenyng4n56nghmjk", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1e0plfg475phrsvrlzw8gwppeva0zk5yg9fgg8c", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1e859xaue4k2jzqw20cv6l7p3tmc378pc3k8g2u", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ehkfl7palwrh6w2hhr2yfrgrq8jetgucudztfe", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1fhr7e04ct0zslmkzqt9smakg3sxrdve6ulclj2", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1fqzqejwkk898fcslw4z4eeqjzesynvrdfr5hte", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1g48268mu5vfp4wk7dk89r0wdrakm9p5xk0q50k", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gf4wlkutql95j7wwsxz490s6fahlvk2s9xpwax", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gjtvly9lel6zskvwtvlg5vhwpu9c9waw7sxzwx", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gp957czryfgyvxwn3tfnyy2f0t9g2p4pqeemx8", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gpx52r9h3zeul45amvcy2pysgvcwddxrgx6cnv", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1grgelyng2v6v3t8z87wu3sxgt9m5s03xfytvz7", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gxju9ky3hwxvqqagrl3dxtl49kjpxq6wlqe6m5", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1hdrlqvyjfy5sdrseecjrutyws9khtxxaux62l7", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1hjct6q7npsspsg3dgvzk3sdf89spmlpfdn6m9d", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jlr62guqwrwkdt4m3y00zh2rrsamhjf9num5xr", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jmykcq8gylmy5tgqtel4xj4q62fdt49sl584xd", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jxv0u20scum4trha72c7ltfgfqef6nsch7q6cu", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1k2d9ed9vgfuk2m58a2d80q9u6qljkh4vfaqjfq", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1kgddca7qj96z0qcxr2c45z73cfl0c75p7f3s2e", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1lktjhnzkpkz3ehrg8psvmwhafg56kfss3q3t8m", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1lzhlnpahvznwfv4jmay2tgaha5kmz5qxerarrl", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1m73mgwn3cm2e8x9a9axa0kw8nqz8a492ms63vn", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1n3mhyp9fvcmuu8l0q8qvjy07x0rql8q46fe2xk", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1pjmngrwcsatsuyy8m3qrunaun67sr9x7z5r2qs", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ptyzewnns2kn37ewtmv6ppsvhdnmeapvtfc9y5", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1qaa9zej9a0ge3ugpx3pxyx602lxh3ztqgfnp42", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1rpgtz9pskr5geavkjz02caqmeep7cwwpv73axj", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ssm0d433seakyak8kcf93yefhknjleeds4y3em", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1tflk30mq5vgqjdly92kkhhq3raev2hnz6eete3", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ualhu3fjgg77g485gmyswkq3w0dp7gys6qzwrv", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1udpsgkgyutgsglauk9vk9rs03a3skc62gup9ny", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1uhnsxv6m83jj3328mhrql7yax3nge5svrv6t6c", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1uutuwrwt3z2a5z8z3uasml3rftlpmu25aga5c6", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1v5y0tg0jllvxf5c3afml8s3awue0ymju89frut", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1vf44d85es37hwl9f4h9gv0e064m0lla60j9luj", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1vvwtk805lxehwle9l4yudmq6mn0g32px9xtkhc", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1xwazl8ftks4gn00y5x3c47auquc62ssuqlj02r", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1y0us8xvsvfvqkk9c6nt5cfyu5au5tww2ztve7q", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
				},
				{
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1083svrca4t350mphfv9x45wq9asrs60cdmrflj", Amount: sdk.NewCoin("uatom", sdk.NewInt(2))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10e4vsut6suau8tk9m6dnrm0slgd6npe3jx5xpv", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10nzaaeh2kq28t3nqsh5m8kmyv90vx7ym5mpakx", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper124maqmcqv8tquy764ktz7cu0gxnzfw54n3vww8", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper130mdu9a0etmeuw52qfxk73pn0ga6gawkxsrlwf", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper132juzk0gdmwuxvx4phug7m3ymyatxlh9734g4w", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper13x77yexvf6qexfjg9czp6jhpv7vpjdwwkyhe4p", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1485u80fdxjan4sd3esrvyw6cyurpvddvzuh48y", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14kn0kk33szpwus9nh8n87fjel8djx0y070ymmj", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14l0fp639yudfl46zauvv8rkzjgd4u0zk2aseys", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14qazscc80zgzx3m0m0aa30ths0p9hg8vdglqrc", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper157v7tczs40axfgejp2m43kwuzqe0wsy0rv8puv", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper15urq2dtp9qce4fyc85m6upwm9xul3049e02707", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper16k579jk6yt2cwmqx9dz5xvq9fug2tekvlu9qdv", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper16yupepagywvlk7uhpfchtwa0stu5f8cyhh54f2", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper17h2x3j7u44qkrq0sk8ul0r2qr440rwgjkfg0gh", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper17mggn4znyeyg25wd7498qxl7r2jhgue8u4qjcq", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper18extdhzzl5c8tr6453e5hzaj3exrdlea90fj3y", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper196ax4vc0lwpxndu9dyhvca7jhxp70rmcvrj90c", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper199mlc7fr6ll5t54w7tts7f4s0cvnqgc59nmuxf", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1clpqr4nrk4khgkxj78fcwwh6dl3uw4epsluffn", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ddle9tczl87gsvmeva3c48nenyng4n56nghmjk", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1e0plfg475phrsvrlzw8gwppeva0zk5yg9fgg8c", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1e859xaue4k2jzqw20cv6l7p3tmc378pc3k8g2u", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ehkfl7palwrh6w2hhr2yfrgrq8jetgucudztfe", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1fhr7e04ct0zslmkzqt9smakg3sxrdve6ulclj2", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1fqzqejwkk898fcslw4z4eeqjzesynvrdfr5hte", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1g48268mu5vfp4wk7dk89r0wdrakm9p5xk0q50k", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gf4wlkutql95j7wwsxz490s6fahlvk2s9xpwax", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gjtvly9lel6zskvwtvlg5vhwpu9c9waw7sxzwx", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gp957czryfgyvxwn3tfnyy2f0t9g2p4pqeemx8", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gpx52r9h3zeul45amvcy2pysgvcwddxrgx6cnv", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1grgelyng2v6v3t8z87wu3sxgt9m5s03xfytvz7", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gxju9ky3hwxvqqagrl3dxtl49kjpxq6wlqe6m5", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1hdrlqvyjfy5sdrseecjrutyws9khtxxaux62l7", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1hjct6q7npsspsg3dgvzk3sdf89spmlpfdn6m9d", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jlr62guqwrwkdt4m3y00zh2rrsamhjf9num5xr", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jmykcq8gylmy5tgqtel4xj4q62fdt49sl584xd", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jxv0u20scum4trha72c7ltfgfqef6nsch7q6cu", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1k2d9ed9vgfuk2m58a2d80q9u6qljkh4vfaqjfq", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1kgddca7qj96z0qcxr2c45z73cfl0c75p7f3s2e", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1lktjhnzkpkz3ehrg8psvmwhafg56kfss3q3t8m", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1lzhlnpahvznwfv4jmay2tgaha5kmz5qxerarrl", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1m73mgwn3cm2e8x9a9axa0kw8nqz8a492ms63vn", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1n3mhyp9fvcmuu8l0q8qvjy07x0rql8q46fe2xk", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1pjmngrwcsatsuyy8m3qrunaun67sr9x7z5r2qs", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ptyzewnns2kn37ewtmv6ppsvhdnmeapvtfc9y5", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1qaa9zej9a0ge3ugpx3pxyx602lxh3ztqgfnp42", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1rpgtz9pskr5geavkjz02caqmeep7cwwpv73axj", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ssm0d433seakyak8kcf93yefhknjleeds4y3em", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1tflk30mq5vgqjdly92kkhhq3raev2hnz6eete3", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ualhu3fjgg77g485gmyswkq3w0dp7gys6qzwrv", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1udpsgkgyutgsglauk9vk9rs03a3skc62gup9ny", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1uhnsxv6m83jj3328mhrql7yax3nge5svrv6t6c", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1uutuwrwt3z2a5z8z3uasml3rftlpmu25aga5c6", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1v5y0tg0jllvxf5c3afml8s3awue0ymju89frut", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1vf44d85es37hwl9f4h9gv0e064m0lla60j9luj", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1vvwtk805lxehwle9l4yudmq6mn0g32px9xtkhc", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1xwazl8ftks4gn00y5x3c47auquc62ssuqlj02r", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1y0us8xvsvfvqkk9c6nt5cfyu5au5tww2ztve7q", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
				},
				{
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1083svrca4t350mphfv9x45wq9asrs60cdmrflj", Amount: sdk.NewCoin("uatom", sdk.NewInt(10))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10e4vsut6suau8tk9m6dnrm0slgd6npe3jx5xpv", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10nzaaeh2kq28t3nqsh5m8kmyv90vx7ym5mpakx", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper124maqmcqv8tquy764ktz7cu0gxnzfw54n3vww8", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper130mdu9a0etmeuw52qfxk73pn0ga6gawkxsrlwf", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper132juzk0gdmwuxvx4phug7m3ymyatxlh9734g4w", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper13x77yexvf6qexfjg9czp6jhpv7vpjdwwkyhe4p", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1485u80fdxjan4sd3esrvyw6cyurpvddvzuh48y", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14kn0kk33szpwus9nh8n87fjel8djx0y070ymmj", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14l0fp639yudfl46zauvv8rkzjgd4u0zk2aseys", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14qazscc80zgzx3m0m0aa30ths0p9hg8vdglqrc", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper157v7tczs40axfgejp2m43kwuzqe0wsy0rv8puv", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper15urq2dtp9qce4fyc85m6upwm9xul3049e02707", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper16k579jk6yt2cwmqx9dz5xvq9fug2tekvlu9qdv", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper16yupepagywvlk7uhpfchtwa0stu5f8cyhh54f2", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper17h2x3j7u44qkrq0sk8ul0r2qr440rwgjkfg0gh", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper17mggn4znyeyg25wd7498qxl7r2jhgue8u4qjcq", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper18extdhzzl5c8tr6453e5hzaj3exrdlea90fj3y", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper196ax4vc0lwpxndu9dyhvca7jhxp70rmcvrj90c", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper199mlc7fr6ll5t54w7tts7f4s0cvnqgc59nmuxf", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1clpqr4nrk4khgkxj78fcwwh6dl3uw4epsluffn", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ddle9tczl87gsvmeva3c48nenyng4n56nghmjk", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1e0plfg475phrsvrlzw8gwppeva0zk5yg9fgg8c", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1e859xaue4k2jzqw20cv6l7p3tmc378pc3k8g2u", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ehkfl7palwrh6w2hhr2yfrgrq8jetgucudztfe", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1fhr7e04ct0zslmkzqt9smakg3sxrdve6ulclj2", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1fqzqejwkk898fcslw4z4eeqjzesynvrdfr5hte", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1g48268mu5vfp4wk7dk89r0wdrakm9p5xk0q50k", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gf4wlkutql95j7wwsxz490s6fahlvk2s9xpwax", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gjtvly9lel6zskvwtvlg5vhwpu9c9waw7sxzwx", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gp957czryfgyvxwn3tfnyy2f0t9g2p4pqeemx8", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gpx52r9h3zeul45amvcy2pysgvcwddxrgx6cnv", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1grgelyng2v6v3t8z87wu3sxgt9m5s03xfytvz7", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gxju9ky3hwxvqqagrl3dxtl49kjpxq6wlqe6m5", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1hdrlqvyjfy5sdrseecjrutyws9khtxxaux62l7", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1hjct6q7npsspsg3dgvzk3sdf89spmlpfdn6m9d", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jlr62guqwrwkdt4m3y00zh2rrsamhjf9num5xr", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jmykcq8gylmy5tgqtel4xj4q62fdt49sl584xd", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jxv0u20scum4trha72c7ltfgfqef6nsch7q6cu", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1k2d9ed9vgfuk2m58a2d80q9u6qljkh4vfaqjfq", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1kgddca7qj96z0qcxr2c45z73cfl0c75p7f3s2e", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1lktjhnzkpkz3ehrg8psvmwhafg56kfss3q3t8m", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1lzhlnpahvznwfv4jmay2tgaha5kmz5qxerarrl", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1m73mgwn3cm2e8x9a9axa0kw8nqz8a492ms63vn", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1n3mhyp9fvcmuu8l0q8qvjy07x0rql8q46fe2xk", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1pjmngrwcsatsuyy8m3qrunaun67sr9x7z5r2qs", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ptyzewnns2kn37ewtmv6ppsvhdnmeapvtfc9y5", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1qaa9zej9a0ge3ugpx3pxyx602lxh3ztqgfnp42", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1rpgtz9pskr5geavkjz02caqmeep7cwwpv73axj", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ssm0d433seakyak8kcf93yefhknjleeds4y3em", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1tflk30mq5vgqjdly92kkhhq3raev2hnz6eete3", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ualhu3fjgg77g485gmyswkq3w0dp7gys6qzwrv", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1udpsgkgyutgsglauk9vk9rs03a3skc62gup9ny", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1uhnsxv6m83jj3328mhrql7yax3nge5svrv6t6c", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1uutuwrwt3z2a5z8z3uasml3rftlpmu25aga5c6", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1v5y0tg0jllvxf5c3afml8s3awue0ymju89frut", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1vf44d85es37hwl9f4h9gv0e064m0lla60j9luj", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1vvwtk805lxehwle9l4yudmq6mn0g32px9xtkhc", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1xwazl8ftks4gn00y5x3c47auquc62ssuqlj02r", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1y0us8xvsvfvqkk9c6nt5cfyu5au5tww2ztve7q", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
				},
				{
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1083svrca4t350mphfv9x45wq9asrs60cdmrflj", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10nzaaeh2kq28t3nqsh5m8kmyv90vx7ym5mpakx", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper124maqmcqv8tquy764ktz7cu0gxnzfw54n3vww8", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper130mdu9a0etmeuw52qfxk73pn0ga6gawkxsrlwf", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper132juzk0gdmwuxvx4phug7m3ymyatxlh9734g4w", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper13x77yexvf6qexfjg9czp6jhpv7vpjdwwkyhe4p", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1485u80fdxjan4sd3esrvyw6cyurpvddvzuh48y", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14kn0kk33szpwus9nh8n87fjel8djx0y070ymmj", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14l0fp639yudfl46zauvv8rkzjgd4u0zk2aseys", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14qazscc80zgzx3m0m0aa30ths0p9hg8vdglqrc", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper157v7tczs40axfgejp2m43kwuzqe0wsy0rv8puv", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper15urq2dtp9qce4fyc85m6upwm9xul3049e02707", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper16k579jk6yt2cwmqx9dz5xvq9fug2tekvlu9qdv", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper16yupepagywvlk7uhpfchtwa0stu5f8cyhh54f2", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper17h2x3j7u44qkrq0sk8ul0r2qr440rwgjkfg0gh", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper17mggn4znyeyg25wd7498qxl7r2jhgue8u4qjcq", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper18extdhzzl5c8tr6453e5hzaj3exrdlea90fj3y", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper196ax4vc0lwpxndu9dyhvca7jhxp70rmcvrj90c", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper199mlc7fr6ll5t54w7tts7f4s0cvnqgc59nmuxf", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1clpqr4nrk4khgkxj78fcwwh6dl3uw4epsluffn", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ddle9tczl87gsvmeva3c48nenyng4n56nghmjk", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1e0plfg475phrsvrlzw8gwppeva0zk5yg9fgg8c", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1e859xaue4k2jzqw20cv6l7p3tmc378pc3k8g2u", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ehkfl7palwrh6w2hhr2yfrgrq8jetgucudztfe", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1fhr7e04ct0zslmkzqt9smakg3sxrdve6ulclj2", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1fqzqejwkk898fcslw4z4eeqjzesynvrdfr5hte", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1g48268mu5vfp4wk7dk89r0wdrakm9p5xk0q50k", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gf4wlkutql95j7wwsxz490s6fahlvk2s9xpwax", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gjtvly9lel6zskvwtvlg5vhwpu9c9waw7sxzwx", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gp957czryfgyvxwn3tfnyy2f0t9g2p4pqeemx8", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gpx52r9h3zeul45amvcy2pysgvcwddxrgx6cnv", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1grgelyng2v6v3t8z87wu3sxgt9m5s03xfytvz7", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gxju9ky3hwxvqqagrl3dxtl49kjpxq6wlqe6m5", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1hdrlqvyjfy5sdrseecjrutyws9khtxxaux62l7", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1hjct6q7npsspsg3dgvzk3sdf89spmlpfdn6m9d", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jlr62guqwrwkdt4m3y00zh2rrsamhjf9num5xr", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jmykcq8gylmy5tgqtel4xj4q62fdt49sl584xd", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jxv0u20scum4trha72c7ltfgfqef6nsch7q6cu", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1k2d9ed9vgfuk2m58a2d80q9u6qljkh4vfaqjfq", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1kgddca7qj96z0qcxr2c45z73cfl0c75p7f3s2e", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
				},
				{
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1083svrca4t350mphfv9x45wq9asrs60cdmrflj", Amount: sdk.NewCoin("uatom", sdk.NewInt(10))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10e4vsut6suau8tk9m6dnrm0slgd6npe3jx5xpv", Amount: sdk.NewCoin("uatom", sdk.NewInt(0))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10nzaaeh2kq28t3nqsh5m8kmyv90vx7ym5mpakx", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper124maqmcqv8tquy764ktz7cu0gxnzfw54n3vww8", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper130mdu9a0etmeuw52qfxk73pn0ga6gawkxsrlwf", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper132juzk0gdmwuxvx4phug7m3ymyatxlh9734g4w", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper13x77yexvf6qexfjg9czp6jhpv7vpjdwwkyhe4p", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1485u80fdxjan4sd3esrvyw6cyurpvddvzuh48y", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14kn0kk33szpwus9nh8n87fjel8djx0y070ymmj", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14l0fp639yudfl46zauvv8rkzjgd4u0zk2aseys", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14qazscc80zgzx3m0m0aa30ths0p9hg8vdglqrc", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper157v7tczs40axfgejp2m43kwuzqe0wsy0rv8puv", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper15urq2dtp9qce4fyc85m6upwm9xul3049e02707", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper16k579jk6yt2cwmqx9dz5xvq9fug2tekvlu9qdv", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper16yupepagywvlk7uhpfchtwa0stu5f8cyhh54f2", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper17h2x3j7u44qkrq0sk8ul0r2qr440rwgjkfg0gh", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper17mggn4znyeyg25wd7498qxl7r2jhgue8u4qjcq", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper18extdhzzl5c8tr6453e5hzaj3exrdlea90fj3y", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper196ax4vc0lwpxndu9dyhvca7jhxp70rmcvrj90c", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper199mlc7fr6ll5t54w7tts7f4s0cvnqgc59nmuxf", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1clpqr4nrk4khgkxj78fcwwh6dl3uw4epsluffn", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ddle9tczl87gsvmeva3c48nenyng4n56nghmjk", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1e0plfg475phrsvrlzw8gwppeva0zk5yg9fgg8c", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1e859xaue4k2jzqw20cv6l7p3tmc378pc3k8g2u", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ehkfl7palwrh6w2hhr2yfrgrq8jetgucudztfe", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1fhr7e04ct0zslmkzqt9smakg3sxrdve6ulclj2", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1fqzqejwkk898fcslw4z4eeqjzesynvrdfr5hte", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1g48268mu5vfp4wk7dk89r0wdrakm9p5xk0q50k", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gf4wlkutql95j7wwsxz490s6fahlvk2s9xpwax", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gjtvly9lel6zskvwtvlg5vhwpu9c9waw7sxzwx", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gp957czryfgyvxwn3tfnyy2f0t9g2p4pqeemx8", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gpx52r9h3zeul45amvcy2pysgvcwddxrgx6cnv", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1grgelyng2v6v3t8z87wu3sxgt9m5s03xfytvz7", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gxju9ky3hwxvqqagrl3dxtl49kjpxq6wlqe6m5", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1hdrlqvyjfy5sdrseecjrutyws9khtxxaux62l7", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1hjct6q7npsspsg3dgvzk3sdf89spmlpfdn6m9d", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jlr62guqwrwkdt4m3y00zh2rrsamhjf9num5xr", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jmykcq8gylmy5tgqtel4xj4q62fdt49sl584xd", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jxv0u20scum4trha72c7ltfgfqef6nsch7q6cu", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1k2d9ed9vgfuk2m58a2d80q9u6qljkh4vfaqjfq", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1kgddca7qj96z0qcxr2c45z73cfl0c75p7f3s2e", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1lktjhnzkpkz3ehrg8psvmwhafg56kfss3q3t8m", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1lzhlnpahvznwfv4jmay2tgaha5kmz5qxerarrl", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1m73mgwn3cm2e8x9a9axa0kw8nqz8a492ms63vn", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1n3mhyp9fvcmuu8l0q8qvjy07x0rql8q46fe2xk", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1pjmngrwcsatsuyy8m3qrunaun67sr9x7z5r2qs", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ptyzewnns2kn37ewtmv6ppsvhdnmeapvtfc9y5", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1qaa9zej9a0ge3ugpx3pxyx602lxh3ztqgfnp42", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1rpgtz9pskr5geavkjz02caqmeep7cwwpv73axj", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ssm0d433seakyak8kcf93yefhknjleeds4y3em", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1tflk30mq5vgqjdly92kkhhq3raev2hnz6eete3", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ualhu3fjgg77g485gmyswkq3w0dp7gys6qzwrv", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1udpsgkgyutgsglauk9vk9rs03a3skc62gup9ny", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1uhnsxv6m83jj3328mhrql7yax3nge5svrv6t6c", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1uutuwrwt3z2a5z8z3uasml3rftlpmu25aga5c6", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1v5y0tg0jllvxf5c3afml8s3awue0ymju89frut", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1vf44d85es37hwl9f4h9gv0e064m0lla60j9luj", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1vvwtk805lxehwle9l4yudmq6mn0g32px9xtkhc", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1xwazl8ftks4gn00y5x3c47auquc62ssuqlj02r", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1y0us8xvsvfvqkk9c6nt5cfyu5au5tww2ztve7q", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
				},
				{
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1083svrca4t350mphfv9x45wq9asrs60cdmrflj", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10nzaaeh2kq28t3nqsh5m8kmyv90vx7ym5mpakx", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper124maqmcqv8tquy764ktz7cu0gxnzfw54n3vww8", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper130mdu9a0etmeuw52qfxk73pn0ga6gawkxsrlwf", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper132juzk0gdmwuxvx4phug7m3ymyatxlh9734g4w", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper13x77yexvf6qexfjg9czp6jhpv7vpjdwwkyhe4p", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1485u80fdxjan4sd3esrvyw6cyurpvddvzuh48y", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14kn0kk33szpwus9nh8n87fjel8djx0y070ymmj", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14l0fp639yudfl46zauvv8rkzjgd4u0zk2aseys", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14qazscc80zgzx3m0m0aa30ths0p9hg8vdglqrc", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper157v7tczs40axfgejp2m43kwuzqe0wsy0rv8puv", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper15urq2dtp9qce4fyc85m6upwm9xul3049e02707", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper16k579jk6yt2cwmqx9dz5xvq9fug2tekvlu9qdv", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper16yupepagywvlk7uhpfchtwa0stu5f8cyhh54f2", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper17h2x3j7u44qkrq0sk8ul0r2qr440rwgjkfg0gh", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper17mggn4znyeyg25wd7498qxl7r2jhgue8u4qjcq", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper18extdhzzl5c8tr6453e5hzaj3exrdlea90fj3y", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper196ax4vc0lwpxndu9dyhvca7jhxp70rmcvrj90c", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper199mlc7fr6ll5t54w7tts7f4s0cvnqgc59nmuxf", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1clpqr4nrk4khgkxj78fcwwh6dl3uw4epsluffn", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ddle9tczl87gsvmeva3c48nenyng4n56nghmjk", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1e0plfg475phrsvrlzw8gwppeva0zk5yg9fgg8c", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1e859xaue4k2jzqw20cv6l7p3tmc378pc3k8g2u", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ehkfl7palwrh6w2hhr2yfrgrq8jetgucudztfe", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1fhr7e04ct0zslmkzqt9smakg3sxrdve6ulclj2", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1fqzqejwkk898fcslw4z4eeqjzesynvrdfr5hte", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1g48268mu5vfp4wk7dk89r0wdrakm9p5xk0q50k", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gf4wlkutql95j7wwsxz490s6fahlvk2s9xpwax", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gjtvly9lel6zskvwtvlg5vhwpu9c9waw7sxzwx", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gp957czryfgyvxwn3tfnyy2f0t9g2p4pqeemx8", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gpx52r9h3zeul45amvcy2pysgvcwddxrgx6cnv", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1grgelyng2v6v3t8z87wu3sxgt9m5s03xfytvz7", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gxju9ky3hwxvqqagrl3dxtl49kjpxq6wlqe6m5", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1hdrlqvyjfy5sdrseecjrutyws9khtxxaux62l7", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1hjct6q7npsspsg3dgvzk3sdf89spmlpfdn6m9d", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jlr62guqwrwkdt4m3y00zh2rrsamhjf9num5xr", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jmykcq8gylmy5tgqtel4xj4q62fdt49sl584xd", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jxv0u20scum4trha72c7ltfgfqef6nsch7q6cu", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1k2d9ed9vgfuk2m58a2d80q9u6qljkh4vfaqjfq", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1kgddca7qj96z0qcxr2c45z73cfl0c75p7f3s2e", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1lktjhnzkpkz3ehrg8psvmwhafg56kfss3q3t8m", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1lzhlnpahvznwfv4jmay2tgaha5kmz5qxerarrl", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1m73mgwn3cm2e8x9a9axa0kw8nqz8a492ms63vn", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1n3mhyp9fvcmuu8l0q8qvjy07x0rql8q46fe2xk", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1pjmngrwcsatsuyy8m3qrunaun67sr9x7z5r2qs", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ptyzewnns2kn37ewtmv6ppsvhdnmeapvtfc9y5", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1qaa9zej9a0ge3ugpx3pxyx602lxh3ztqgfnp42", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1rpgtz9pskr5geavkjz02caqmeep7cwwpv73axj", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ssm0d433seakyak8kcf93yefhknjleeds4y3em", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1tflk30mq5vgqjdly92kkhhq3raev2hnz6eete3", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ualhu3fjgg77g485gmyswkq3w0dp7gys6qzwrv", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1udpsgkgyutgsglauk9vk9rs03a3skc62gup9ny", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1uhnsxv6m83jj3328mhrql7yax3nge5svrv6t6c", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1uutuwrwt3z2a5z8z3uasml3rftlpmu25aga5c6", Amount: sdk.NewCoin("uatom", sdk.NewInt(2))},
				},
				{
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1083svrca4t350mphfv9x45wq9asrs60cdmrflj", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10e4vsut6suau8tk9m6dnrm0slgd6npe3jx5xpv", Amount: sdk.NewCoin("uatom", sdk.NewInt(34))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10nzaaeh2kq28t3nqsh5m8kmyv90vx7ym5mpakx", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper124maqmcqv8tquy764ktz7cu0gxnzfw54n3vww8", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper130mdu9a0etmeuw52qfxk73pn0ga6gawkxsrlwf", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper132juzk0gdmwuxvx4phug7m3ymyatxlh9734g4w", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper13x77yexvf6qexfjg9czp6jhpv7vpjdwwkyhe4p", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1485u80fdxjan4sd3esrvyw6cyurpvddvzuh48y", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14kn0kk33szpwus9nh8n87fjel8djx0y070ymmj", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14l0fp639yudfl46zauvv8rkzjgd4u0zk2aseys", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14qazscc80zgzx3m0m0aa30ths0p9hg8vdglqrc", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper157v7tczs40axfgejp2m43kwuzqe0wsy0rv8puv", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper15urq2dtp9qce4fyc85m6upwm9xul3049e02707", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper16k579jk6yt2cwmqx9dz5xvq9fug2tekvlu9qdv", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper16yupepagywvlk7uhpfchtwa0stu5f8cyhh54f2", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper17h2x3j7u44qkrq0sk8ul0r2qr440rwgjkfg0gh", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper17mggn4znyeyg25wd7498qxl7r2jhgue8u4qjcq", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper18extdhzzl5c8tr6453e5hzaj3exrdlea90fj3y", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper196ax4vc0lwpxndu9dyhvca7jhxp70rmcvrj90c", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper199mlc7fr6ll5t54w7tts7f4s0cvnqgc59nmuxf", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1clpqr4nrk4khgkxj78fcwwh6dl3uw4epsluffn", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ddle9tczl87gsvmeva3c48nenyng4n56nghmjk", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1e0plfg475phrsvrlzw8gwppeva0zk5yg9fgg8c", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1e859xaue4k2jzqw20cv6l7p3tmc378pc3k8g2u", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ehkfl7palwrh6w2hhr2yfrgrq8jetgucudztfe", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1fhr7e04ct0zslmkzqt9smakg3sxrdve6ulclj2", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1fqzqejwkk898fcslw4z4eeqjzesynvrdfr5hte", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1g48268mu5vfp4wk7dk89r0wdrakm9p5xk0q50k", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gf4wlkutql95j7wwsxz490s6fahlvk2s9xpwax", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gjtvly9lel6zskvwtvlg5vhwpu9c9waw7sxzwx", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gp957czryfgyvxwn3tfnyy2f0t9g2p4pqeemx8", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gpx52r9h3zeul45amvcy2pysgvcwddxrgx6cnv", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1grgelyng2v6v3t8z87wu3sxgt9m5s03xfytvz7", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gxju9ky3hwxvqqagrl3dxtl49kjpxq6wlqe6m5", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1hdrlqvyjfy5sdrseecjrutyws9khtxxaux62l7", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1hjct6q7npsspsg3dgvzk3sdf89spmlpfdn6m9d", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jlr62guqwrwkdt4m3y00zh2rrsamhjf9num5xr", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jmykcq8gylmy5tgqtel4xj4q62fdt49sl584xd", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jxv0u20scum4trha72c7ltfgfqef6nsch7q6cu", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1k2d9ed9vgfuk2m58a2d80q9u6qljkh4vfaqjfq", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1kgddca7qj96z0qcxr2c45z73cfl0c75p7f3s2e", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1lktjhnzkpkz3ehrg8psvmwhafg56kfss3q3t8m", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1lzhlnpahvznwfv4jmay2tgaha5kmz5qxerarrl", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1m73mgwn3cm2e8x9a9axa0kw8nqz8a492ms63vn", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1n3mhyp9fvcmuu8l0q8qvjy07x0rql8q46fe2xk", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1pjmngrwcsatsuyy8m3qrunaun67sr9x7z5r2qs", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ptyzewnns2kn37ewtmv6ppsvhdnmeapvtfc9y5", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1qaa9zej9a0ge3ugpx3pxyx602lxh3ztqgfnp42", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1rpgtz9pskr5geavkjz02caqmeep7cwwpv73axj", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", Amount: sdk.NewCoin("uatom", sdk.NewInt(126))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ssm0d433seakyak8kcf93yefhknjleeds4y3em", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1tflk30mq5vgqjdly92kkhhq3raev2hnz6eete3", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ualhu3fjgg77g485gmyswkq3w0dp7gys6qzwrv", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1udpsgkgyutgsglauk9vk9rs03a3skc62gup9ny", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1uhnsxv6m83jj3328mhrql7yax3nge5svrv6t6c", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1uutuwrwt3z2a5z8z3uasml3rftlpmu25aga5c6", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1v5y0tg0jllvxf5c3afml8s3awue0ymju89frut", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1vf44d85es37hwl9f4h9gv0e064m0lla60j9luj", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1vvwtk805lxehwle9l4yudmq6mn0g32px9xtkhc", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1xwazl8ftks4gn00y5x3c47auquc62ssuqlj02r", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1y0us8xvsvfvqkk9c6nt5cfyu5au5tww2ztve7q", Amount: sdk.NewCoin("uatom", sdk.NewInt(64))},
				},
				{
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1083svrca4t350mphfv9x45wq9asrs60cdmrflj", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10e4vsut6suau8tk9m6dnrm0slgd6npe3jx5xpv", Amount: sdk.NewCoin("uatom", sdk.NewInt(16128970))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10nzaaeh2kq28t3nqsh5m8kmyv90vx7ym5mpakx", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper124maqmcqv8tquy764ktz7cu0gxnzfw54n3vww8", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper130mdu9a0etmeuw52qfxk73pn0ga6gawkxsrlwf", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper132juzk0gdmwuxvx4phug7m3ymyatxlh9734g4w", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper13x77yexvf6qexfjg9czp6jhpv7vpjdwwkyhe4p", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1485u80fdxjan4sd3esrvyw6cyurpvddvzuh48y", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14kn0kk33szpwus9nh8n87fjel8djx0y070ymmj", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14l0fp639yudfl46zauvv8rkzjgd4u0zk2aseys", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14qazscc80zgzx3m0m0aa30ths0p9hg8vdglqrc", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper157v7tczs40axfgejp2m43kwuzqe0wsy0rv8puv", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper15urq2dtp9qce4fyc85m6upwm9xul3049e02707", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper16k579jk6yt2cwmqx9dz5xvq9fug2tekvlu9qdv", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper16yupepagywvlk7uhpfchtwa0stu5f8cyhh54f2", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper17h2x3j7u44qkrq0sk8ul0r2qr440rwgjkfg0gh", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper17mggn4znyeyg25wd7498qxl7r2jhgue8u4qjcq", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper18extdhzzl5c8tr6453e5hzaj3exrdlea90fj3y", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper196ax4vc0lwpxndu9dyhvca7jhxp70rmcvrj90c", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper199mlc7fr6ll5t54w7tts7f4s0cvnqgc59nmuxf", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1clpqr4nrk4khgkxj78fcwwh6dl3uw4epsluffn", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ddle9tczl87gsvmeva3c48nenyng4n56nghmjk", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1e0plfg475phrsvrlzw8gwppeva0zk5yg9fgg8c", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1e859xaue4k2jzqw20cv6l7p3tmc378pc3k8g2u", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ehkfl7palwrh6w2hhr2yfrgrq8jetgucudztfe", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1fhr7e04ct0zslmkzqt9smakg3sxrdve6ulclj2", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1fqzqejwkk898fcslw4z4eeqjzesynvrdfr5hte", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1g48268mu5vfp4wk7dk89r0wdrakm9p5xk0q50k", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gf4wlkutql95j7wwsxz490s6fahlvk2s9xpwax", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gjtvly9lel6zskvwtvlg5vhwpu9c9waw7sxzwx", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gp957czryfgyvxwn3tfnyy2f0t9g2p4pqeemx8", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gpx52r9h3zeul45amvcy2pysgvcwddxrgx6cnv", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1grgelyng2v6v3t8z87wu3sxgt9m5s03xfytvz7", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gxju9ky3hwxvqqagrl3dxtl49kjpxq6wlqe6m5", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1hdrlqvyjfy5sdrseecjrutyws9khtxxaux62l7", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1hjct6q7npsspsg3dgvzk3sdf89spmlpfdn6m9d", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jlr62guqwrwkdt4m3y00zh2rrsamhjf9num5xr", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jmykcq8gylmy5tgqtel4xj4q62fdt49sl584xd", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jxv0u20scum4trha72c7ltfgfqef6nsch7q6cu", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1k2d9ed9vgfuk2m58a2d80q9u6qljkh4vfaqjfq", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1kgddca7qj96z0qcxr2c45z73cfl0c75p7f3s2e", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1lktjhnzkpkz3ehrg8psvmwhafg56kfss3q3t8m", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1lzhlnpahvznwfv4jmay2tgaha5kmz5qxerarrl", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1m73mgwn3cm2e8x9a9axa0kw8nqz8a492ms63vn", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1n3mhyp9fvcmuu8l0q8qvjy07x0rql8q46fe2xk", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1pjmngrwcsatsuyy8m3qrunaun67sr9x7z5r2qs", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ptyzewnns2kn37ewtmv6ppsvhdnmeapvtfc9y5", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1qaa9zej9a0ge3ugpx3pxyx602lxh3ztqgfnp42", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1rpgtz9pskr5geavkjz02caqmeep7cwwpv73axj", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", Amount: sdk.NewCoin("uatom", sdk.NewInt(16131030))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ssm0d433seakyak8kcf93yefhknjleeds4y3em", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1tflk30mq5vgqjdly92kkhhq3raev2hnz6eete3", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ualhu3fjgg77g485gmyswkq3w0dp7gys6qzwrv", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1udpsgkgyutgsglauk9vk9rs03a3skc62gup9ny", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1uhnsxv6m83jj3328mhrql7yax3nge5svrv6t6c", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1uutuwrwt3z2a5z8z3uasml3rftlpmu25aga5c6", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1v5y0tg0jllvxf5c3afml8s3awue0ymju89frut", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1vf44d85es37hwl9f4h9gv0e064m0lla60j9luj", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1vvwtk805lxehwle9l4yudmq6mn0g32px9xtkhc", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1xwazl8ftks4gn00y5x3c47auquc62ssuqlj02r", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1y0us8xvsvfvqkk9c6nt5cfyu5au5tww2ztve7q", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000))},
				},
				{
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1083svrca4t350mphfv9x45wq9asrs60cdmrflj", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10e4vsut6suau8tk9m6dnrm0slgd6npe3jx5xpv", Amount: sdk.NewCoin("uatom", sdk.NewInt(161289999970))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10nzaaeh2kq28t3nqsh5m8kmyv90vx7ym5mpakx", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper124maqmcqv8tquy764ktz7cu0gxnzfw54n3vww8", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper130mdu9a0etmeuw52qfxk73pn0ga6gawkxsrlwf", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper132juzk0gdmwuxvx4phug7m3ymyatxlh9734g4w", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper13x77yexvf6qexfjg9czp6jhpv7vpjdwwkyhe4p", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1485u80fdxjan4sd3esrvyw6cyurpvddvzuh48y", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14kn0kk33szpwus9nh8n87fjel8djx0y070ymmj", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14l0fp639yudfl46zauvv8rkzjgd4u0zk2aseys", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14qazscc80zgzx3m0m0aa30ths0p9hg8vdglqrc", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper157v7tczs40axfgejp2m43kwuzqe0wsy0rv8puv", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper15urq2dtp9qce4fyc85m6upwm9xul3049e02707", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper16k579jk6yt2cwmqx9dz5xvq9fug2tekvlu9qdv", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper16yupepagywvlk7uhpfchtwa0stu5f8cyhh54f2", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper17h2x3j7u44qkrq0sk8ul0r2qr440rwgjkfg0gh", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper17mggn4znyeyg25wd7498qxl7r2jhgue8u4qjcq", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper18extdhzzl5c8tr6453e5hzaj3exrdlea90fj3y", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper196ax4vc0lwpxndu9dyhvca7jhxp70rmcvrj90c", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper199mlc7fr6ll5t54w7tts7f4s0cvnqgc59nmuxf", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1clpqr4nrk4khgkxj78fcwwh6dl3uw4epsluffn", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ddle9tczl87gsvmeva3c48nenyng4n56nghmjk", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1e0plfg475phrsvrlzw8gwppeva0zk5yg9fgg8c", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1e859xaue4k2jzqw20cv6l7p3tmc378pc3k8g2u", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ehkfl7palwrh6w2hhr2yfrgrq8jetgucudztfe", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1fhr7e04ct0zslmkzqt9smakg3sxrdve6ulclj2", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1fqzqejwkk898fcslw4z4eeqjzesynvrdfr5hte", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1g48268mu5vfp4wk7dk89r0wdrakm9p5xk0q50k", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gf4wlkutql95j7wwsxz490s6fahlvk2s9xpwax", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gjtvly9lel6zskvwtvlg5vhwpu9c9waw7sxzwx", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gp957czryfgyvxwn3tfnyy2f0t9g2p4pqeemx8", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gpx52r9h3zeul45amvcy2pysgvcwddxrgx6cnv", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1grgelyng2v6v3t8z87wu3sxgt9m5s03xfytvz7", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gxju9ky3hwxvqqagrl3dxtl49kjpxq6wlqe6m5", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1hdrlqvyjfy5sdrseecjrutyws9khtxxaux62l7", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1hjct6q7npsspsg3dgvzk3sdf89spmlpfdn6m9d", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jlr62guqwrwkdt4m3y00zh2rrsamhjf9num5xr", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jmykcq8gylmy5tgqtel4xj4q62fdt49sl584xd", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jxv0u20scum4trha72c7ltfgfqef6nsch7q6cu", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1k2d9ed9vgfuk2m58a2d80q9u6qljkh4vfaqjfq", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1kgddca7qj96z0qcxr2c45z73cfl0c75p7f3s2e", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1lktjhnzkpkz3ehrg8psvmwhafg56kfss3q3t8m", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1lzhlnpahvznwfv4jmay2tgaha5kmz5qxerarrl", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1m73mgwn3cm2e8x9a9axa0kw8nqz8a492ms63vn", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1n3mhyp9fvcmuu8l0q8qvjy07x0rql8q46fe2xk", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1pjmngrwcsatsuyy8m3qrunaun67sr9x7z5r2qs", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ptyzewnns2kn37ewtmv6ppsvhdnmeapvtfc9y5", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1qaa9zej9a0ge3ugpx3pxyx602lxh3ztqgfnp42", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1rpgtz9pskr5geavkjz02caqmeep7cwwpv73axj", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", Amount: sdk.NewCoin("uatom", sdk.NewInt(161310000031))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ssm0d433seakyak8kcf93yefhknjleeds4y3em", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1tflk30mq5vgqjdly92kkhhq3raev2hnz6eete3", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ualhu3fjgg77g485gmyswkq3w0dp7gys6qzwrv", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1udpsgkgyutgsglauk9vk9rs03a3skc62gup9ny", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1uhnsxv6m83jj3328mhrql7yax3nge5svrv6t6c", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1uutuwrwt3z2a5z8z3uasml3rftlpmu25aga5c6", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1v5y0tg0jllvxf5c3afml8s3awue0ymju89frut", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1vf44d85es37hwl9f4h9gv0e064m0lla60j9luj", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1vvwtk805lxehwle9l4yudmq6mn0g32px9xtkhc", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1xwazl8ftks4gn00y5x3c47auquc62ssuqlj02r", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1y0us8xvsvfvqkk9c6nt5cfyu5au5tww2ztve7q", Amount: sdk.NewCoin("uatom", sdk.NewInt(161290000000))},
				},
				{
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1083svrca4t350mphfv9x45wq9asrs60cdmrflj", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10e4vsut6suau8tk9m6dnrm0slgd6npe3jx5xpv", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001583))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10nzaaeh2kq28t3nqsh5m8kmyv90vx7ym5mpakx", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper124maqmcqv8tquy764ktz7cu0gxnzfw54n3vww8", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper130mdu9a0etmeuw52qfxk73pn0ga6gawkxsrlwf", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper132juzk0gdmwuxvx4phug7m3ymyatxlh9734g4w", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper13x77yexvf6qexfjg9czp6jhpv7vpjdwwkyhe4p", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1485u80fdxjan4sd3esrvyw6cyurpvddvzuh48y", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14kn0kk33szpwus9nh8n87fjel8djx0y070ymmj", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14l0fp639yudfl46zauvv8rkzjgd4u0zk2aseys", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14qazscc80zgzx3m0m0aa30ths0p9hg8vdglqrc", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper157v7tczs40axfgejp2m43kwuzqe0wsy0rv8puv", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper15urq2dtp9qce4fyc85m6upwm9xul3049e02707", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper16k579jk6yt2cwmqx9dz5xvq9fug2tekvlu9qdv", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper16yupepagywvlk7uhpfchtwa0stu5f8cyhh54f2", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper17h2x3j7u44qkrq0sk8ul0r2qr440rwgjkfg0gh", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper17mggn4znyeyg25wd7498qxl7r2jhgue8u4qjcq", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper18extdhzzl5c8tr6453e5hzaj3exrdlea90fj3y", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper196ax4vc0lwpxndu9dyhvca7jhxp70rmcvrj90c", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper199mlc7fr6ll5t54w7tts7f4s0cvnqgc59nmuxf", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1clpqr4nrk4khgkxj78fcwwh6dl3uw4epsluffn", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ddle9tczl87gsvmeva3c48nenyng4n56nghmjk", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1e0plfg475phrsvrlzw8gwppeva0zk5yg9fgg8c", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1e859xaue4k2jzqw20cv6l7p3tmc378pc3k8g2u", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ehkfl7palwrh6w2hhr2yfrgrq8jetgucudztfe", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1fhr7e04ct0zslmkzqt9smakg3sxrdve6ulclj2", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1fqzqejwkk898fcslw4z4eeqjzesynvrdfr5hte", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1g48268mu5vfp4wk7dk89r0wdrakm9p5xk0q50k", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gf4wlkutql95j7wwsxz490s6fahlvk2s9xpwax", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gjtvly9lel6zskvwtvlg5vhwpu9c9waw7sxzwx", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gp957czryfgyvxwn3tfnyy2f0t9g2p4pqeemx8", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gpx52r9h3zeul45amvcy2pysgvcwddxrgx6cnv", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1grgelyng2v6v3t8z87wu3sxgt9m5s03xfytvz7", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gxju9ky3hwxvqqagrl3dxtl49kjpxq6wlqe6m5", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1hdrlqvyjfy5sdrseecjrutyws9khtxxaux62l7", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1hjct6q7npsspsg3dgvzk3sdf89spmlpfdn6m9d", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jlr62guqwrwkdt4m3y00zh2rrsamhjf9num5xr", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jmykcq8gylmy5tgqtel4xj4q62fdt49sl584xd", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jxv0u20scum4trha72c7ltfgfqef6nsch7q6cu", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1k2d9ed9vgfuk2m58a2d80q9u6qljkh4vfaqjfq", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1kgddca7qj96z0qcxr2c45z73cfl0c75p7f3s2e", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1lktjhnzkpkz3ehrg8psvmwhafg56kfss3q3t8m", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1lzhlnpahvznwfv4jmay2tgaha5kmz5qxerarrl", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1m73mgwn3cm2e8x9a9axa0kw8nqz8a492ms63vn", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1n3mhyp9fvcmuu8l0q8qvjy07x0rql8q46fe2xk", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1pjmngrwcsatsuyy8m3qrunaun67sr9x7z5r2qs", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ptyzewnns2kn37ewtmv6ppsvhdnmeapvtfc9y5", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1qaa9zej9a0ge3ugpx3pxyx602lxh3ztqgfnp42", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1rpgtz9pskr5geavkjz02caqmeep7cwwpv73axj", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", Amount: sdk.NewCoin("uatom", sdk.NewInt(16131000001637))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ssm0d433seakyak8kcf93yefhknjleeds4y3em", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1tflk30mq5vgqjdly92kkhhq3raev2hnz6eete3", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ualhu3fjgg77g485gmyswkq3w0dp7gys6qzwrv", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1udpsgkgyutgsglauk9vk9rs03a3skc62gup9ny", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1uhnsxv6m83jj3328mhrql7yax3nge5svrv6t6c", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1uutuwrwt3z2a5z8z3uasml3rftlpmu25aga5c6", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1v5y0tg0jllvxf5c3afml8s3awue0ymju89frut", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1vf44d85es37hwl9f4h9gv0e064m0lla60j9luj", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1vvwtk805lxehwle9l4yudmq6mn0g32px9xtkhc", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1xwazl8ftks4gn00y5x3c47auquc62ssuqlj02r", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1y0us8xvsvfvqkk9c6nt5cfyu5au5tww2ztve7q", Amount: sdk.NewCoin("uatom", sdk.NewInt(16129000001613))},
				},
			},
			ExpectedListWithUndelegateDistribution: []types.ValAddressAmounts{
				{
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10e4vsut6suau8tk9m6dnrm0slgd6npe3jx5xpv", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
				},
				{
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10e4vsut6suau8tk9m6dnrm0slgd6npe3jx5xpv", Amount: sdk.NewCoin("uatom", sdk.NewInt(2))},
				},
				{
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10e4vsut6suau8tk9m6dnrm0slgd6npe3jx5xpv", Amount: sdk.NewCoin("uatom", sdk.NewInt(10))},
				},
				{
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10e4vsut6suau8tk9m6dnrm0slgd6npe3jx5xpv", Amount: sdk.NewCoin("uatom", sdk.NewInt(31))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1083svrca4t350mphfv9x45wq9asrs60cdmrflj", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10nzaaeh2kq28t3nqsh5m8kmyv90vx7ym5mpakx", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper124maqmcqv8tquy764ktz7cu0gxnzfw54n3vww8", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper130mdu9a0etmeuw52qfxk73pn0ga6gawkxsrlwf", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper132juzk0gdmwuxvx4phug7m3ymyatxlh9734g4w", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper13x77yexvf6qexfjg9czp6jhpv7vpjdwwkyhe4p", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1485u80fdxjan4sd3esrvyw6cyurpvddvzuh48y", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14kn0kk33szpwus9nh8n87fjel8djx0y070ymmj", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14l0fp639yudfl46zauvv8rkzjgd4u0zk2aseys", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
				},
				{
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10e4vsut6suau8tk9m6dnrm0slgd6npe3jx5xpv", Amount: sdk.NewCoin("uatom", sdk.NewInt(31))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1083svrca4t350mphfv9x45wq9asrs60cdmrflj", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10nzaaeh2kq28t3nqsh5m8kmyv90vx7ym5mpakx", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper124maqmcqv8tquy764ktz7cu0gxnzfw54n3vww8", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper130mdu9a0etmeuw52qfxk73pn0ga6gawkxsrlwf", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper132juzk0gdmwuxvx4phug7m3ymyatxlh9734g4w", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper13x77yexvf6qexfjg9czp6jhpv7vpjdwwkyhe4p", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1485u80fdxjan4sd3esrvyw6cyurpvddvzuh48y", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14kn0kk33szpwus9nh8n87fjel8djx0y070ymmj", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14l0fp639yudfl46zauvv8rkzjgd4u0zk2aseys", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14qazscc80zgzx3m0m0aa30ths0p9hg8vdglqrc", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper157v7tczs40axfgejp2m43kwuzqe0wsy0rv8puv", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper15urq2dtp9qce4fyc85m6upwm9xul3049e02707", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper16k579jk6yt2cwmqx9dz5xvq9fug2tekvlu9qdv", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper16yupepagywvlk7uhpfchtwa0stu5f8cyhh54f2", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper17h2x3j7u44qkrq0sk8ul0r2qr440rwgjkfg0gh", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper17mggn4znyeyg25wd7498qxl7r2jhgue8u4qjcq", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper18extdhzzl5c8tr6453e5hzaj3exrdlea90fj3y", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper196ax4vc0lwpxndu9dyhvca7jhxp70rmcvrj90c", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper199mlc7fr6ll5t54w7tts7f4s0cvnqgc59nmuxf", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1clpqr4nrk4khgkxj78fcwwh6dl3uw4epsluffn", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ddle9tczl87gsvmeva3c48nenyng4n56nghmjk", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1e0plfg475phrsvrlzw8gwppeva0zk5yg9fgg8c", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1e859xaue4k2jzqw20cv6l7p3tmc378pc3k8g2u", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ehkfl7palwrh6w2hhr2yfrgrq8jetgucudztfe", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1fhr7e04ct0zslmkzqt9smakg3sxrdve6ulclj2", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1fqzqejwkk898fcslw4z4eeqjzesynvrdfr5hte", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1g48268mu5vfp4wk7dk89r0wdrakm9p5xk0q50k", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gf4wlkutql95j7wwsxz490s6fahlvk2s9xpwax", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gjtvly9lel6zskvwtvlg5vhwpu9c9waw7sxzwx", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gp957czryfgyvxwn3tfnyy2f0t9g2p4pqeemx8", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gpx52r9h3zeul45amvcy2pysgvcwddxrgx6cnv", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1grgelyng2v6v3t8z87wu3sxgt9m5s03xfytvz7", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gxju9ky3hwxvqqagrl3dxtl49kjpxq6wlqe6m5", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1hdrlqvyjfy5sdrseecjrutyws9khtxxaux62l7", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1hjct6q7npsspsg3dgvzk3sdf89spmlpfdn6m9d", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jlr62guqwrwkdt4m3y00zh2rrsamhjf9num5xr", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jmykcq8gylmy5tgqtel4xj4q62fdt49sl584xd", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jxv0u20scum4trha72c7ltfgfqef6nsch7q6cu", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1k2d9ed9vgfuk2m58a2d80q9u6qljkh4vfaqjfq", Amount: sdk.NewCoin("uatom", sdk.NewInt(1))},
				},
				{
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10e4vsut6suau8tk9m6dnrm0slgd6npe3jx5xpv", Amount: sdk.NewCoin("uatom", sdk.NewInt(34))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1083svrca4t350mphfv9x45wq9asrs60cdmrflj", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10nzaaeh2kq28t3nqsh5m8kmyv90vx7ym5mpakx", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper124maqmcqv8tquy764ktz7cu0gxnzfw54n3vww8", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper130mdu9a0etmeuw52qfxk73pn0ga6gawkxsrlwf", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper132juzk0gdmwuxvx4phug7m3ymyatxlh9734g4w", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper13x77yexvf6qexfjg9czp6jhpv7vpjdwwkyhe4p", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1485u80fdxjan4sd3esrvyw6cyurpvddvzuh48y", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14kn0kk33szpwus9nh8n87fjel8djx0y070ymmj", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14l0fp639yudfl46zauvv8rkzjgd4u0zk2aseys", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14qazscc80zgzx3m0m0aa30ths0p9hg8vdglqrc", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper157v7tczs40axfgejp2m43kwuzqe0wsy0rv8puv", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper15urq2dtp9qce4fyc85m6upwm9xul3049e02707", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper16k579jk6yt2cwmqx9dz5xvq9fug2tekvlu9qdv", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper16yupepagywvlk7uhpfchtwa0stu5f8cyhh54f2", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper17h2x3j7u44qkrq0sk8ul0r2qr440rwgjkfg0gh", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper17mggn4znyeyg25wd7498qxl7r2jhgue8u4qjcq", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper18extdhzzl5c8tr6453e5hzaj3exrdlea90fj3y", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper196ax4vc0lwpxndu9dyhvca7jhxp70rmcvrj90c", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper199mlc7fr6ll5t54w7tts7f4s0cvnqgc59nmuxf", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1clpqr4nrk4khgkxj78fcwwh6dl3uw4epsluffn", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ddle9tczl87gsvmeva3c48nenyng4n56nghmjk", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1e0plfg475phrsvrlzw8gwppeva0zk5yg9fgg8c", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1e859xaue4k2jzqw20cv6l7p3tmc378pc3k8g2u", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ehkfl7palwrh6w2hhr2yfrgrq8jetgucudztfe", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1fhr7e04ct0zslmkzqt9smakg3sxrdve6ulclj2", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1fqzqejwkk898fcslw4z4eeqjzesynvrdfr5hte", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1g48268mu5vfp4wk7dk89r0wdrakm9p5xk0q50k", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gf4wlkutql95j7wwsxz490s6fahlvk2s9xpwax", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gjtvly9lel6zskvwtvlg5vhwpu9c9waw7sxzwx", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gp957czryfgyvxwn3tfnyy2f0t9g2p4pqeemx8", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gpx52r9h3zeul45amvcy2pysgvcwddxrgx6cnv", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1grgelyng2v6v3t8z87wu3sxgt9m5s03xfytvz7", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gxju9ky3hwxvqqagrl3dxtl49kjpxq6wlqe6m5", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1hdrlqvyjfy5sdrseecjrutyws9khtxxaux62l7", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1hjct6q7npsspsg3dgvzk3sdf89spmlpfdn6m9d", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jlr62guqwrwkdt4m3y00zh2rrsamhjf9num5xr", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jmykcq8gylmy5tgqtel4xj4q62fdt49sl584xd", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jxv0u20scum4trha72c7ltfgfqef6nsch7q6cu", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1k2d9ed9vgfuk2m58a2d80q9u6qljkh4vfaqjfq", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1kgddca7qj96z0qcxr2c45z73cfl0c75p7f3s2e", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1lktjhnzkpkz3ehrg8psvmwhafg56kfss3q3t8m", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1lzhlnpahvznwfv4jmay2tgaha5kmz5qxerarrl", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1m73mgwn3cm2e8x9a9axa0kw8nqz8a492ms63vn", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1n3mhyp9fvcmuu8l0q8qvjy07x0rql8q46fe2xk", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1pjmngrwcsatsuyy8m3qrunaun67sr9x7z5r2qs", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ptyzewnns2kn37ewtmv6ppsvhdnmeapvtfc9y5", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1qaa9zej9a0ge3ugpx3pxyx602lxh3ztqgfnp42", Amount: sdk.NewCoin("uatom", sdk.NewInt(4))},
				},
				{
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10e4vsut6suau8tk9m6dnrm0slgd6npe3jx5xpv", Amount: sdk.NewCoin("uatom", sdk.NewInt(95))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1083svrca4t350mphfv9x45wq9asrs60cdmrflj", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10nzaaeh2kq28t3nqsh5m8kmyv90vx7ym5mpakx", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper124maqmcqv8tquy764ktz7cu0gxnzfw54n3vww8", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper130mdu9a0etmeuw52qfxk73pn0ga6gawkxsrlwf", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper132juzk0gdmwuxvx4phug7m3ymyatxlh9734g4w", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper13x77yexvf6qexfjg9czp6jhpv7vpjdwwkyhe4p", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1485u80fdxjan4sd3esrvyw6cyurpvddvzuh48y", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14kn0kk33szpwus9nh8n87fjel8djx0y070ymmj", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14l0fp639yudfl46zauvv8rkzjgd4u0zk2aseys", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14qazscc80zgzx3m0m0aa30ths0p9hg8vdglqrc", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper157v7tczs40axfgejp2m43kwuzqe0wsy0rv8puv", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper15urq2dtp9qce4fyc85m6upwm9xul3049e02707", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper16k579jk6yt2cwmqx9dz5xvq9fug2tekvlu9qdv", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper16yupepagywvlk7uhpfchtwa0stu5f8cyhh54f2", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper17h2x3j7u44qkrq0sk8ul0r2qr440rwgjkfg0gh", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper17mggn4znyeyg25wd7498qxl7r2jhgue8u4qjcq", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper18extdhzzl5c8tr6453e5hzaj3exrdlea90fj3y", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper196ax4vc0lwpxndu9dyhvca7jhxp70rmcvrj90c", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper199mlc7fr6ll5t54w7tts7f4s0cvnqgc59nmuxf", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1clpqr4nrk4khgkxj78fcwwh6dl3uw4epsluffn", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ddle9tczl87gsvmeva3c48nenyng4n56nghmjk", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1e0plfg475phrsvrlzw8gwppeva0zk5yg9fgg8c", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1e859xaue4k2jzqw20cv6l7p3tmc378pc3k8g2u", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ehkfl7palwrh6w2hhr2yfrgrq8jetgucudztfe", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1fhr7e04ct0zslmkzqt9smakg3sxrdve6ulclj2", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1fqzqejwkk898fcslw4z4eeqjzesynvrdfr5hte", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1g48268mu5vfp4wk7dk89r0wdrakm9p5xk0q50k", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gf4wlkutql95j7wwsxz490s6fahlvk2s9xpwax", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gjtvly9lel6zskvwtvlg5vhwpu9c9waw7sxzwx", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gp957czryfgyvxwn3tfnyy2f0t9g2p4pqeemx8", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gpx52r9h3zeul45amvcy2pysgvcwddxrgx6cnv", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1grgelyng2v6v3t8z87wu3sxgt9m5s03xfytvz7", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gxju9ky3hwxvqqagrl3dxtl49kjpxq6wlqe6m5", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1hdrlqvyjfy5sdrseecjrutyws9khtxxaux62l7", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1hjct6q7npsspsg3dgvzk3sdf89spmlpfdn6m9d", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jlr62guqwrwkdt4m3y00zh2rrsamhjf9num5xr", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jmykcq8gylmy5tgqtel4xj4q62fdt49sl584xd", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jxv0u20scum4trha72c7ltfgfqef6nsch7q6cu", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1k2d9ed9vgfuk2m58a2d80q9u6qljkh4vfaqjfq", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1kgddca7qj96z0qcxr2c45z73cfl0c75p7f3s2e", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1lktjhnzkpkz3ehrg8psvmwhafg56kfss3q3t8m", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1lzhlnpahvznwfv4jmay2tgaha5kmz5qxerarrl", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1m73mgwn3cm2e8x9a9axa0kw8nqz8a492ms63vn", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1n3mhyp9fvcmuu8l0q8qvjy07x0rql8q46fe2xk", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1pjmngrwcsatsuyy8m3qrunaun67sr9x7z5r2qs", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ptyzewnns2kn37ewtmv6ppsvhdnmeapvtfc9y5", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1qaa9zej9a0ge3ugpx3pxyx602lxh3ztqgfnp42", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1rpgtz9pskr5geavkjz02caqmeep7cwwpv73axj", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ssm0d433seakyak8kcf93yefhknjleeds4y3em", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1tflk30mq5vgqjdly92kkhhq3raev2hnz6eete3", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ualhu3fjgg77g485gmyswkq3w0dp7gys6qzwrv", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1udpsgkgyutgsglauk9vk9rs03a3skc62gup9ny", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1uhnsxv6m83jj3328mhrql7yax3nge5svrv6t6c", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1uutuwrwt3z2a5z8z3uasml3rftlpmu25aga5c6", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1v5y0tg0jllvxf5c3afml8s3awue0ymju89frut", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1vf44d85es37hwl9f4h9gv0e064m0lla60j9luj", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1vvwtk805lxehwle9l4yudmq6mn0g32px9xtkhc", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1xwazl8ftks4gn00y5x3c47auquc62ssuqlj02r", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1y0us8xvsvfvqkk9c6nt5cfyu5au5tww2ztve7q", Amount: sdk.NewCoin("uatom", sdk.NewInt(65))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", Amount: sdk.NewCoin("uatom", sdk.NewInt(5))},
				},
				{
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339690))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10e4vsut6suau8tk9m6dnrm0slgd6npe3jx5xpv", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339554))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1083svrca4t350mphfv9x45wq9asrs60cdmrflj", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10nzaaeh2kq28t3nqsh5m8kmyv90vx7ym5mpakx", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper124maqmcqv8tquy764ktz7cu0gxnzfw54n3vww8", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper130mdu9a0etmeuw52qfxk73pn0ga6gawkxsrlwf", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper132juzk0gdmwuxvx4phug7m3ymyatxlh9734g4w", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper13x77yexvf6qexfjg9czp6jhpv7vpjdwwkyhe4p", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1485u80fdxjan4sd3esrvyw6cyurpvddvzuh48y", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14kn0kk33szpwus9nh8n87fjel8djx0y070ymmj", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14l0fp639yudfl46zauvv8rkzjgd4u0zk2aseys", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14qazscc80zgzx3m0m0aa30ths0p9hg8vdglqrc", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper157v7tczs40axfgejp2m43kwuzqe0wsy0rv8puv", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper15urq2dtp9qce4fyc85m6upwm9xul3049e02707", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper16k579jk6yt2cwmqx9dz5xvq9fug2tekvlu9qdv", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper16yupepagywvlk7uhpfchtwa0stu5f8cyhh54f2", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper17h2x3j7u44qkrq0sk8ul0r2qr440rwgjkfg0gh", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper17mggn4znyeyg25wd7498qxl7r2jhgue8u4qjcq", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper18extdhzzl5c8tr6453e5hzaj3exrdlea90fj3y", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper196ax4vc0lwpxndu9dyhvca7jhxp70rmcvrj90c", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper199mlc7fr6ll5t54w7tts7f4s0cvnqgc59nmuxf", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1clpqr4nrk4khgkxj78fcwwh6dl3uw4epsluffn", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ddle9tczl87gsvmeva3c48nenyng4n56nghmjk", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1e0plfg475phrsvrlzw8gwppeva0zk5yg9fgg8c", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1e859xaue4k2jzqw20cv6l7p3tmc378pc3k8g2u", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ehkfl7palwrh6w2hhr2yfrgrq8jetgucudztfe", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1fhr7e04ct0zslmkzqt9smakg3sxrdve6ulclj2", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1fqzqejwkk898fcslw4z4eeqjzesynvrdfr5hte", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1g48268mu5vfp4wk7dk89r0wdrakm9p5xk0q50k", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gf4wlkutql95j7wwsxz490s6fahlvk2s9xpwax", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gjtvly9lel6zskvwtvlg5vhwpu9c9waw7sxzwx", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gp957czryfgyvxwn3tfnyy2f0t9g2p4pqeemx8", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gpx52r9h3zeul45amvcy2pysgvcwddxrgx6cnv", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1grgelyng2v6v3t8z87wu3sxgt9m5s03xfytvz7", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gxju9ky3hwxvqqagrl3dxtl49kjpxq6wlqe6m5", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1hdrlqvyjfy5sdrseecjrutyws9khtxxaux62l7", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1hjct6q7npsspsg3dgvzk3sdf89spmlpfdn6m9d", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jlr62guqwrwkdt4m3y00zh2rrsamhjf9num5xr", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jmykcq8gylmy5tgqtel4xj4q62fdt49sl584xd", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jxv0u20scum4trha72c7ltfgfqef6nsch7q6cu", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1k2d9ed9vgfuk2m58a2d80q9u6qljkh4vfaqjfq", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1kgddca7qj96z0qcxr2c45z73cfl0c75p7f3s2e", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1lktjhnzkpkz3ehrg8psvmwhafg56kfss3q3t8m", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1lzhlnpahvznwfv4jmay2tgaha5kmz5qxerarrl", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1m73mgwn3cm2e8x9a9axa0kw8nqz8a492ms63vn", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1n3mhyp9fvcmuu8l0q8qvjy07x0rql8q46fe2xk", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1pjmngrwcsatsuyy8m3qrunaun67sr9x7z5r2qs", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ptyzewnns2kn37ewtmv6ppsvhdnmeapvtfc9y5", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1qaa9zej9a0ge3ugpx3pxyx602lxh3ztqgfnp42", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1rpgtz9pskr5geavkjz02caqmeep7cwwpv73axj", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ssm0d433seakyak8kcf93yefhknjleeds4y3em", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1tflk30mq5vgqjdly92kkhhq3raev2hnz6eete3", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ualhu3fjgg77g485gmyswkq3w0dp7gys6qzwrv", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1udpsgkgyutgsglauk9vk9rs03a3skc62gup9ny", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1uhnsxv6m83jj3328mhrql7yax3nge5svrv6t6c", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1uutuwrwt3z2a5z8z3uasml3rftlpmu25aga5c6", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1v5y0tg0jllvxf5c3afml8s3awue0ymju89frut", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1vf44d85es37hwl9f4h9gv0e064m0lla60j9luj", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1vvwtk805lxehwle9l4yudmq6mn0g32px9xtkhc", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1xwazl8ftks4gn00y5x3c47auquc62ssuqlj02r", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1y0us8xvsvfvqkk9c6nt5cfyu5au5tww2ztve7q", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339500))},
				},
				{
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339690))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10e4vsut6suau8tk9m6dnrm0slgd6npe3jx5xpv", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339554))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1083svrca4t350mphfv9x45wq9asrs60cdmrflj", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10nzaaeh2kq28t3nqsh5m8kmyv90vx7ym5mpakx", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper124maqmcqv8tquy764ktz7cu0gxnzfw54n3vww8", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper130mdu9a0etmeuw52qfxk73pn0ga6gawkxsrlwf", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper132juzk0gdmwuxvx4phug7m3ymyatxlh9734g4w", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper13x77yexvf6qexfjg9czp6jhpv7vpjdwwkyhe4p", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1485u80fdxjan4sd3esrvyw6cyurpvddvzuh48y", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14kn0kk33szpwus9nh8n87fjel8djx0y070ymmj", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14l0fp639yudfl46zauvv8rkzjgd4u0zk2aseys", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14qazscc80zgzx3m0m0aa30ths0p9hg8vdglqrc", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper157v7tczs40axfgejp2m43kwuzqe0wsy0rv8puv", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper15urq2dtp9qce4fyc85m6upwm9xul3049e02707", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper16k579jk6yt2cwmqx9dz5xvq9fug2tekvlu9qdv", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper16yupepagywvlk7uhpfchtwa0stu5f8cyhh54f2", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper17h2x3j7u44qkrq0sk8ul0r2qr440rwgjkfg0gh", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper17mggn4znyeyg25wd7498qxl7r2jhgue8u4qjcq", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper18extdhzzl5c8tr6453e5hzaj3exrdlea90fj3y", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper196ax4vc0lwpxndu9dyhvca7jhxp70rmcvrj90c", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper199mlc7fr6ll5t54w7tts7f4s0cvnqgc59nmuxf", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1clpqr4nrk4khgkxj78fcwwh6dl3uw4epsluffn", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ddle9tczl87gsvmeva3c48nenyng4n56nghmjk", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1e0plfg475phrsvrlzw8gwppeva0zk5yg9fgg8c", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1e859xaue4k2jzqw20cv6l7p3tmc378pc3k8g2u", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ehkfl7palwrh6w2hhr2yfrgrq8jetgucudztfe", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1fhr7e04ct0zslmkzqt9smakg3sxrdve6ulclj2", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1fqzqejwkk898fcslw4z4eeqjzesynvrdfr5hte", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1g48268mu5vfp4wk7dk89r0wdrakm9p5xk0q50k", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gf4wlkutql95j7wwsxz490s6fahlvk2s9xpwax", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gjtvly9lel6zskvwtvlg5vhwpu9c9waw7sxzwx", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gp957czryfgyvxwn3tfnyy2f0t9g2p4pqeemx8", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gpx52r9h3zeul45amvcy2pysgvcwddxrgx6cnv", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1grgelyng2v6v3t8z87wu3sxgt9m5s03xfytvz7", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gxju9ky3hwxvqqagrl3dxtl49kjpxq6wlqe6m5", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1hdrlqvyjfy5sdrseecjrutyws9khtxxaux62l7", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1hjct6q7npsspsg3dgvzk3sdf89spmlpfdn6m9d", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jlr62guqwrwkdt4m3y00zh2rrsamhjf9num5xr", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jmykcq8gylmy5tgqtel4xj4q62fdt49sl584xd", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jxv0u20scum4trha72c7ltfgfqef6nsch7q6cu", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1k2d9ed9vgfuk2m58a2d80q9u6qljkh4vfaqjfq", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1kgddca7qj96z0qcxr2c45z73cfl0c75p7f3s2e", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1lktjhnzkpkz3ehrg8psvmwhafg56kfss3q3t8m", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1lzhlnpahvznwfv4jmay2tgaha5kmz5qxerarrl", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1m73mgwn3cm2e8x9a9axa0kw8nqz8a492ms63vn", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1n3mhyp9fvcmuu8l0q8qvjy07x0rql8q46fe2xk", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1pjmngrwcsatsuyy8m3qrunaun67sr9x7z5r2qs", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ptyzewnns2kn37ewtmv6ppsvhdnmeapvtfc9y5", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1qaa9zej9a0ge3ugpx3pxyx602lxh3ztqgfnp42", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1rpgtz9pskr5geavkjz02caqmeep7cwwpv73axj", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ssm0d433seakyak8kcf93yefhknjleeds4y3em", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1tflk30mq5vgqjdly92kkhhq3raev2hnz6eete3", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ualhu3fjgg77g485gmyswkq3w0dp7gys6qzwrv", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1udpsgkgyutgsglauk9vk9rs03a3skc62gup9ny", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1uhnsxv6m83jj3328mhrql7yax3nge5svrv6t6c", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1uutuwrwt3z2a5z8z3uasml3rftlpmu25aga5c6", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1v5y0tg0jllvxf5c3afml8s3awue0ymju89frut", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1vf44d85es37hwl9f4h9gv0e064m0lla60j9luj", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1vvwtk805lxehwle9l4yudmq6mn0g32px9xtkhc", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1xwazl8ftks4gn00y5x3c47auquc62ssuqlj02r", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1y0us8xvsvfvqkk9c6nt5cfyu5au5tww2ztve7q", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339520))},
				},
				{
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339690))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10e4vsut6suau8tk9m6dnrm0slgd6npe3jx5xpv", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339554))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1083svrca4t350mphfv9x45wq9asrs60cdmrflj", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10nzaaeh2kq28t3nqsh5m8kmyv90vx7ym5mpakx", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper124maqmcqv8tquy764ktz7cu0gxnzfw54n3vww8", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper130mdu9a0etmeuw52qfxk73pn0ga6gawkxsrlwf", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper132juzk0gdmwuxvx4phug7m3ymyatxlh9734g4w", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper13x77yexvf6qexfjg9czp6jhpv7vpjdwwkyhe4p", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1485u80fdxjan4sd3esrvyw6cyurpvddvzuh48y", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14kn0kk33szpwus9nh8n87fjel8djx0y070ymmj", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14l0fp639yudfl46zauvv8rkzjgd4u0zk2aseys", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14qazscc80zgzx3m0m0aa30ths0p9hg8vdglqrc", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper157v7tczs40axfgejp2m43kwuzqe0wsy0rv8puv", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper15urq2dtp9qce4fyc85m6upwm9xul3049e02707", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper16k579jk6yt2cwmqx9dz5xvq9fug2tekvlu9qdv", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper16yupepagywvlk7uhpfchtwa0stu5f8cyhh54f2", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper17h2x3j7u44qkrq0sk8ul0r2qr440rwgjkfg0gh", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper17mggn4znyeyg25wd7498qxl7r2jhgue8u4qjcq", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper18extdhzzl5c8tr6453e5hzaj3exrdlea90fj3y", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper196ax4vc0lwpxndu9dyhvca7jhxp70rmcvrj90c", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper199mlc7fr6ll5t54w7tts7f4s0cvnqgc59nmuxf", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1clpqr4nrk4khgkxj78fcwwh6dl3uw4epsluffn", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ddle9tczl87gsvmeva3c48nenyng4n56nghmjk", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1e0plfg475phrsvrlzw8gwppeva0zk5yg9fgg8c", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1e859xaue4k2jzqw20cv6l7p3tmc378pc3k8g2u", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ehkfl7palwrh6w2hhr2yfrgrq8jetgucudztfe", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1fhr7e04ct0zslmkzqt9smakg3sxrdve6ulclj2", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1fqzqejwkk898fcslw4z4eeqjzesynvrdfr5hte", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1g48268mu5vfp4wk7dk89r0wdrakm9p5xk0q50k", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gf4wlkutql95j7wwsxz490s6fahlvk2s9xpwax", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gjtvly9lel6zskvwtvlg5vhwpu9c9waw7sxzwx", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gp957czryfgyvxwn3tfnyy2f0t9g2p4pqeemx8", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gpx52r9h3zeul45amvcy2pysgvcwddxrgx6cnv", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1grgelyng2v6v3t8z87wu3sxgt9m5s03xfytvz7", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gxju9ky3hwxvqqagrl3dxtl49kjpxq6wlqe6m5", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1hdrlqvyjfy5sdrseecjrutyws9khtxxaux62l7", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1hjct6q7npsspsg3dgvzk3sdf89spmlpfdn6m9d", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jlr62guqwrwkdt4m3y00zh2rrsamhjf9num5xr", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jmykcq8gylmy5tgqtel4xj4q62fdt49sl584xd", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jxv0u20scum4trha72c7ltfgfqef6nsch7q6cu", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1k2d9ed9vgfuk2m58a2d80q9u6qljkh4vfaqjfq", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1kgddca7qj96z0qcxr2c45z73cfl0c75p7f3s2e", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1lktjhnzkpkz3ehrg8psvmwhafg56kfss3q3t8m", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1lzhlnpahvznwfv4jmay2tgaha5kmz5qxerarrl", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1m73mgwn3cm2e8x9a9axa0kw8nqz8a492ms63vn", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1n3mhyp9fvcmuu8l0q8qvjy07x0rql8q46fe2xk", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1pjmngrwcsatsuyy8m3qrunaun67sr9x7z5r2qs", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ptyzewnns2kn37ewtmv6ppsvhdnmeapvtfc9y5", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1qaa9zej9a0ge3ugpx3pxyx602lxh3ztqgfnp42", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1rpgtz9pskr5geavkjz02caqmeep7cwwpv73axj", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ssm0d433seakyak8kcf93yefhknjleeds4y3em", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1tflk30mq5vgqjdly92kkhhq3raev2hnz6eete3", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ualhu3fjgg77g485gmyswkq3w0dp7gys6qzwrv", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1udpsgkgyutgsglauk9vk9rs03a3skc62gup9ny", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1uhnsxv6m83jj3328mhrql7yax3nge5svrv6t6c", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1uutuwrwt3z2a5z8z3uasml3rftlpmu25aga5c6", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1v5y0tg0jllvxf5c3afml8s3awue0ymju89frut", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1vf44d85es37hwl9f4h9gv0e064m0lla60j9luj", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1vvwtk805lxehwle9l4yudmq6mn0g32px9xtkhc", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1xwazl8ftks4gn00y5x3c47auquc62ssuqlj02r", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1y0us8xvsvfvqkk9c6nt5cfyu5au5tww2ztve7q", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339523))},
				},
				{
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339690))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10e4vsut6suau8tk9m6dnrm0slgd6npe3jx5xpv", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339554))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1083svrca4t350mphfv9x45wq9asrs60cdmrflj", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper10nzaaeh2kq28t3nqsh5m8kmyv90vx7ym5mpakx", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper124maqmcqv8tquy764ktz7cu0gxnzfw54n3vww8", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper130mdu9a0etmeuw52qfxk73pn0ga6gawkxsrlwf", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper132juzk0gdmwuxvx4phug7m3ymyatxlh9734g4w", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper13x77yexvf6qexfjg9czp6jhpv7vpjdwwkyhe4p", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1485u80fdxjan4sd3esrvyw6cyurpvddvzuh48y", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14kn0kk33szpwus9nh8n87fjel8djx0y070ymmj", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14l0fp639yudfl46zauvv8rkzjgd4u0zk2aseys", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper14qazscc80zgzx3m0m0aa30ths0p9hg8vdglqrc", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper157v7tczs40axfgejp2m43kwuzqe0wsy0rv8puv", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper15urq2dtp9qce4fyc85m6upwm9xul3049e02707", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper16k579jk6yt2cwmqx9dz5xvq9fug2tekvlu9qdv", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper16yupepagywvlk7uhpfchtwa0stu5f8cyhh54f2", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper17h2x3j7u44qkrq0sk8ul0r2qr440rwgjkfg0gh", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper17mggn4znyeyg25wd7498qxl7r2jhgue8u4qjcq", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper18extdhzzl5c8tr6453e5hzaj3exrdlea90fj3y", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper196ax4vc0lwpxndu9dyhvca7jhxp70rmcvrj90c", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper199mlc7fr6ll5t54w7tts7f4s0cvnqgc59nmuxf", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1clpqr4nrk4khgkxj78fcwwh6dl3uw4epsluffn", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ddle9tczl87gsvmeva3c48nenyng4n56nghmjk", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1e0plfg475phrsvrlzw8gwppeva0zk5yg9fgg8c", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1e859xaue4k2jzqw20cv6l7p3tmc378pc3k8g2u", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ehkfl7palwrh6w2hhr2yfrgrq8jetgucudztfe", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1fhr7e04ct0zslmkzqt9smakg3sxrdve6ulclj2", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1fqzqejwkk898fcslw4z4eeqjzesynvrdfr5hte", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1g48268mu5vfp4wk7dk89r0wdrakm9p5xk0q50k", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gf4wlkutql95j7wwsxz490s6fahlvk2s9xpwax", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gjtvly9lel6zskvwtvlg5vhwpu9c9waw7sxzwx", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gp957czryfgyvxwn3tfnyy2f0t9g2p4pqeemx8", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gpx52r9h3zeul45amvcy2pysgvcwddxrgx6cnv", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1grgelyng2v6v3t8z87wu3sxgt9m5s03xfytvz7", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1gxju9ky3hwxvqqagrl3dxtl49kjpxq6wlqe6m5", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1hdrlqvyjfy5sdrseecjrutyws9khtxxaux62l7", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1hjct6q7npsspsg3dgvzk3sdf89spmlpfdn6m9d", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jlr62guqwrwkdt4m3y00zh2rrsamhjf9num5xr", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jmykcq8gylmy5tgqtel4xj4q62fdt49sl584xd", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1jxv0u20scum4trha72c7ltfgfqef6nsch7q6cu", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1k2d9ed9vgfuk2m58a2d80q9u6qljkh4vfaqjfq", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1kgddca7qj96z0qcxr2c45z73cfl0c75p7f3s2e", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1lktjhnzkpkz3ehrg8psvmwhafg56kfss3q3t8m", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1lzhlnpahvznwfv4jmay2tgaha5kmz5qxerarrl", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1m73mgwn3cm2e8x9a9axa0kw8nqz8a492ms63vn", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1n3mhyp9fvcmuu8l0q8qvjy07x0rql8q46fe2xk", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1pjmngrwcsatsuyy8m3qrunaun67sr9x7z5r2qs", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ptyzewnns2kn37ewtmv6ppsvhdnmeapvtfc9y5", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1qaa9zej9a0ge3ugpx3pxyx602lxh3ztqgfnp42", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1rpgtz9pskr5geavkjz02caqmeep7cwwpv73axj", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ssm0d433seakyak8kcf93yefhknjleeds4y3em", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1tflk30mq5vgqjdly92kkhhq3raev2hnz6eete3", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1ualhu3fjgg77g485gmyswkq3w0dp7gys6qzwrv", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1udpsgkgyutgsglauk9vk9rs03a3skc62gup9ny", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1uhnsxv6m83jj3328mhrql7yax3nge5svrv6t6c", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1uutuwrwt3z2a5z8z3uasml3rftlpmu25aga5c6", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1v5y0tg0jllvxf5c3afml8s3awue0ymju89frut", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1vf44d85es37hwl9f4h9gv0e064m0lla60j9luj", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1vvwtk805lxehwle9l4yudmq6mn0g32px9xtkhc", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1xwazl8ftks4gn00y5x3c47auquc62ssuqlj02r", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
					types.ValAddressAmount{ValidatorAddr: "cosmosvaloper1y0us8xvsvfvqkk9c6nt5cfyu5au5tww2ztve7q", Amount: sdk.NewCoin("uatom", sdk.NewInt(1339524))},
				},
			},
		},
	}

	for _, tc := range testCases {
		k.SetDelegationState(ctx, tc.DelegationState)
		k.SetAllowListedValidators(ctx, tc.AllowListedValidators)
		for i, amount := range tc.DelegateAmounts {
			// fetch a combined updated val set list and delegation state
			updateValList, hostAccountDelegations := k.GetAllValidatorsState(ctx, hostChainParams.BaseDenom)

			// sort both updatedValList and hostAccountDelegations
			sort.Sort(updateValList)
			sort.Sort(hostAccountDelegations)

			// get the current delegation state and
			// assign the updated validator delegation state to the current delegation state
			delegationStateS := k.GetDelegationState(ctx)
			delegationStateS.HostAccountDelegations = hostAccountDelegations

			allowListerValidators := types.AllowListedValidators{AllowListedValidators: updateValList}

			// get list of validator with respective amounts to delegate
			list, err := keeper.FetchValidatorsToDelegate(allowListerValidators, delegationStateS, amount)
			suite.NoError(err)
			suite.Equal(len(list), len(tc.ExpectedListWithDelegateDistribution[i]))
			for j := range list {
				suite.Equal(list[j].Amount.Amount.ToDec(), tc.ExpectedListWithDelegateDistribution[i][j].Amount.Amount.ToDec())
				suite.Equal(list[j].ValidatorAddr, tc.ExpectedListWithDelegateDistribution[i][j].ValidatorAddr)
			}
		}

		for i, amount := range tc.UndelegateAmounts {
			// fetch a combined updated val set list and delegation state
			updateValList, hostAccountDelegations := k.GetAllValidatorsState(ctx, hostChainParams.BaseDenom)

			// sort both updatedValList and hostAccountDelegations
			sort.Sort(updateValList)
			sort.Sort(hostAccountDelegations)

			// get the current delegation state and
			// assign the updated validator delegation state to the current delegation state
			delegationStateS := k.GetDelegationState(ctx)
			delegationStateS.HostAccountDelegations = hostAccountDelegations

			allowListerValidators := types.AllowListedValidators{AllowListedValidators: updateValList}

			list, err := keeper.FetchValidatorsToUndelegate(allowListerValidators, delegationStateS, amount)
			suite.NoError(err)
			for j := range list {
				suite.Equal(list[j].Amount.Amount.ToDec(), tc.ExpectedListWithUndelegateDistribution[i][j].Amount.Amount.ToDec())
				suite.Equal(list[j].ValidatorAddr, tc.ExpectedListWithUndelegateDistribution[i][j].ValidatorAddr)
			}
		}
	}
}
