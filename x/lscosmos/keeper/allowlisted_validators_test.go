package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"sort"

	"github.com/persistenceOne/pstake-native/x/lscosmos/keeper"
	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

func (suite *IntegrationTestSuite) TestAllowListedValidators() {
	app, ctx := suite.app, suite.ctx

	// set empty allow listed validators
	app.LSCosmosKeeper.SetAllowListedValidators(ctx, types.AllowListedValidators{})

	resAllowListedValidators := app.LSCosmosKeeper.GetAllowListedValidators(ctx)
	suite.Nil(resAllowListedValidators.AllowListedValidators)

	// set the filled allow listed validators
	app.LSCosmosKeeper.SetAllowListedValidators(ctx, allowListedValidators)

	resAllowListedValidators = app.LSCosmosKeeper.GetAllowListedValidators(ctx)
	suite.Equal(allowListedValidators, resAllowListedValidators)
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
	updateValList, hostAccountDelegations := k.GetAllValidatorsState(ctx)

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
	updateValList, hostAccountDelegations = k.GetAllValidatorsState(ctx)

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
