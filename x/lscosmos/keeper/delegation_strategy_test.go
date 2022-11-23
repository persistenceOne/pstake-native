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
		state := keeper.GetIdealCurrentDelegations(test.givenValset, test.givenDelegationState, sdk.NewInt64Coin(denom, int64(math.Abs(float64(test.amount)))), !(test.amount > 0))

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
