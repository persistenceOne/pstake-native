package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/gogo/protobuf/proto"
	icqtypes "github.com/persistenceOne/persistence-sdk/v2/x/interchainquery/types"

	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

func (suite *IntegrationTestSuite) TestHandleRewardsAccountCallbacks() {
	app, ctx := suite.app, suite.ctx

	lscosmosKeeper := app.LSCosmosKeeper
	hostChainParams := lscosmosKeeper.GetHostChainParams(ctx)

	balance := sdk.NewInt64Coin(hostChainParams.BaseDenom, 0)
	queryBalancesResponse, err := proto.Marshal(&banktypes.QueryBalanceResponse{Balance: &balance})
	suite.NoError(err)

	err = lscosmosKeeper.HandleRewardsAccountBalanceCallback(ctx, queryBalancesResponse, icqtypes.Query{})
	suite.NoError(err)

	delegationState := types.DelegationState{
		HostDelegationAccountBalance: sdk.NewCoins(sdk.NewInt64Coin(hostChainParams.BaseDenom, 100)),
		HostChainDelegationAddress:   "address_________________",
		HostAccountDelegations: []types.HostAccountDelegation{
			{
				ValidatorAddress: "address_______________1",
				Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 25),
			},
			{
				ValidatorAddress: "address_______________2",
				Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 75),
			},
		},
	}
	app.LSCosmosKeeper.SetDelegationState(ctx, delegationState)

	balance = sdk.NewInt64Coin(hostChainParams.BaseDenom, 150)
	queryBalancesResponse, err = proto.Marshal(&banktypes.QueryBalanceResponse{Balance: &balance})
	suite.NoError(err)

	err = lscosmosKeeper.HandleRewardsAccountBalanceCallback(ctx, queryBalancesResponse, icqtypes.Query{})
	suite.Error(err)
}
