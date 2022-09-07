package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

func (suite *IntegrationTestSuite) TestHostChainParamsQuery() {
	app, ctx := suite.app, suite.ctx

	c := sdk.WrapSDKContext(ctx)

	qrysrv := types.QueryServer(app.LSCosmosKeeper)

	res, err := qrysrv.HostChainParams(c, &types.QueryHostChainParamsRequest{})
	suite.NoError(err)
	suite.Equal(&types.QueryHostChainParamsResponse{HostChainParams: app.LSCosmosKeeper.GetHostChainParams(ctx)}, res)
}

func (suite *IntegrationTestSuite) TestQueryDelegationState() {
	app, ctx := suite.app, suite.ctx

	baseDenom := app.LSCosmosKeeper.GetHostChainParams(ctx).BaseDenom
	delegationState := types.DelegationState{
		HostDelegationAccountBalance: sdk.NewCoins(sdk.NewInt64Coin(baseDenom, 100)),
		HostChainDelegationAddress:   "address________________",
		HostAccountDelegations: []types.HostAccountDelegation{
			{
				ValidatorAddress: "address________________",
				Amount:           sdk.NewInt64Coin(baseDenom, 25),
			},
			{
				ValidatorAddress: "address________________",
				Amount:           sdk.NewInt64Coin(baseDenom, 75),
			},
		},
	}
	app.LSCosmosKeeper.SetDelegationState(ctx, delegationState)

	c := sdk.WrapSDKContext(ctx)

	qrysrv := types.QueryServer(app.LSCosmosKeeper)

	res, err := qrysrv.DelegationState(c, &types.QueryDelegationStateRequest{})
	suite.NoError(err)
	suite.Equal(&types.QueryDelegationStateResponse{DelegationState: delegationState}, res)
}

func (suite *IntegrationTestSuite) TestQueryAllowListedValidators() {
	app, ctx := suite.app, suite.ctx

	c := sdk.WrapSDKContext(ctx)

	allowListedValidators := types.AllowListedValidators{
		AllowListedValidators: []types.AllowListedValidator{
			{
				ValidatorAddress: "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt",
				TargetWeight:     sdk.NewDecWithPrec(33, 2),
			},
			{
				ValidatorAddress: "cosmosvaloper1lcck2cxh7dzgkrfk53kysg9ktdrsjj6jfwlnm2",
				TargetWeight:     sdk.NewDecWithPrec(33, 2),
			},
			{
				ValidatorAddress: "cosmosvaloper10khgeppewe4rgfrcy809r9h00aquwxxxgwgwa5",
				TargetWeight:     sdk.NewDecWithPrec(34, 2),
			},
		},
	}
	app.LSCosmosKeeper.SetAllowListedValidators(ctx, allowListedValidators)

	qrysrv := types.QueryServer(app.LSCosmosKeeper)

	res, err := qrysrv.AllowListedValidators(c, &types.QueryAllowListedValidatorsRequest{})
	suite.NoError(err)
	suite.Equal(&types.QueryAllowListedValidatorsResponse{AllowListedValidators: allowListedValidators}, res)
}

func (suite *IntegrationTestSuite) TestQueryCValue() {
	app, ctx := suite.app, suite.ctx

	c := sdk.WrapSDKContext(ctx)

	qrysrv := types.QueryServer(app.LSCosmosKeeper)

	res, err := qrysrv.CValue(c, &types.QueryCValueRequest{})
	suite.NoError(err)
	suite.Equal(&types.QueryCValueResponse{CValue: sdk.NewDec(1)}, res)
}

func (suite *IntegrationTestSuite) TestQueryModuleState() {
	app, ctx := suite.app, suite.ctx

	c := sdk.WrapSDKContext(ctx)

	qrysrv := types.QueryServer(app.LSCosmosKeeper)

	res, err := qrysrv.ModuleState(c, &types.QueryModuleStateRequest{})
	suite.NoError(err)
	suite.Equal(&types.QueryModuleStateResponse{ModuleState: false}, res)

	app.LSCosmosKeeper.SetModuleState(ctx, true)

	res, err = qrysrv.ModuleState(c, &types.QueryModuleStateRequest{})
	suite.NoError(err)
	suite.Equal(&types.QueryModuleStateResponse{ModuleState: true}, res)
}
