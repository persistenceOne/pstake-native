package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	"github.com/gogo/protobuf/proto"
	icqtypes "github.com/persistenceOne/persistence-sdk/v2/x/interchainquery/types"

	"github.com/persistenceOne/pstake-native/v2/x/lscosmos/types"
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

func (suite *IntegrationTestSuite) TestHandleDelegationCallbacks() {
	ctx := suite.ctx
	app := suite.app
	lscosmosKeeper := app.LSCosmosKeeper
	hostChainParams := lscosmosKeeper.GetHostChainParams(ctx)
	hostAccounts := lscosmosKeeper.GetHostAccounts(ctx)

	valAddr1 := sdk.ValAddress("valAddr1")
	valAddrStr1, err := types.Bech32FromValAddress(valAddr1, types.CosmosValOperPrefix)
	suite.NoError(err)
	valAddr2 := sdk.ValAddress("valAddr2")
	valAddrStr2, err := types.Bech32FromValAddress(valAddr2, types.CosmosValOperPrefix)
	suite.NoError(err)

	delegationState := types.DelegationState{
		HostDelegationAccountBalance: sdk.NewCoins(sdk.NewInt64Coin(hostChainParams.BaseDenom, 100)),
		HostChainDelegationAddress:   "address_________________",
		HostAccountDelegations: []types.HostAccountDelegation{
			{
				ValidatorAddress: valAddrStr1,
				Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 25),
			},
			{
				ValidatorAddress: valAddrStr2,
				Amount:           sdk.NewInt64Coin(hostChainParams.BaseDenom, 75),
			},
		},
	}
	app.LSCosmosKeeper.SetDelegationState(ctx, delegationState)

	delegationResponse, err := proto.Marshal(&stakingtypes.QueryDelegationResponse{DelegationResponse: &stakingtypes.DelegationResponse{
		Delegation: stakingtypes.Delegation{
			DelegatorAddress: delegationState.HostChainDelegationAddress,
			ValidatorAddress: valAddrStr1,
			Shares:           sdk.MustNewDecFromStr("1000000"),
		},
		Balance: sdk.NewInt64Coin(hostChainParams.BaseDenom, 100),
	}})
	err = lscosmosKeeper.HandleDelegationCallback(ctx, delegationResponse, icqtypes.Query{})
	// no set ibc states
	suite.Error(err)
	// Old delegation remains
	suite.Equal(lscosmosKeeper.GetHostAccountDelegation(ctx, valAddrStr1).Amount, sdk.NewInt64Coin(hostChainParams.BaseDenom, 25))

	//setIBCStates
	app.ICAControllerKeeper.SetActiveChannelID(ctx, hostChainParams.ConnectionID, hostAccounts.DelegatorAccountPortID(), "channel-1")
	app.IBCKeeper.ChannelKeeper.SetChannel(ctx, hostAccounts.DelegatorAccountPortID(), "channel-1", channeltypes.Channel{
		State:          channeltypes.OPEN,
		Ordering:       2,
		Counterparty:   channeltypes.Counterparty{},
		ConnectionHops: []string{hostChainParams.ConnectionID},
		Version:        "",
	})
	app.IBCKeeper.ChannelKeeper.SetNextSequenceSend(ctx, hostAccounts.DelegatorAccountPortID(), "channel-1", 1)
	app.IBCKeeper.ChannelKeeper.SetNextSequenceAck(ctx, hostAccounts.DelegatorAccountPortID(), "channel-1", 1)
	app.ICAControllerKeeper.SetActiveChannelID(ctx, hostChainParams.ConnectionID, hostAccounts.RewardsAccountPortID(), "channel-2")
	app.IBCKeeper.ChannelKeeper.SetChannel(ctx, hostAccounts.RewardsAccountPortID(), "channel-2", channeltypes.Channel{
		State:          channeltypes.OPEN,
		Ordering:       2,
		Counterparty:   channeltypes.Counterparty{},
		ConnectionHops: []string{hostChainParams.ConnectionID},
		Version:        "",
	})
	app.IBCKeeper.ChannelKeeper.SetNextSequenceSend(ctx, hostAccounts.RewardsAccountPortID(), "channel-2", 1)
	app.IBCKeeper.ChannelKeeper.SetNextSequenceAck(ctx, hostAccounts.RewardsAccountPortID(), "channel-2", 1)

	delegationResponse, err = proto.Marshal(&stakingtypes.QueryDelegationResponse{DelegationResponse: &stakingtypes.DelegationResponse{
		Delegation: stakingtypes.Delegation{
			DelegatorAddress: delegationState.HostChainDelegationAddress,
			ValidatorAddress: valAddrStr1,
			Shares:           sdk.MustNewDecFromStr("1000000"),
		},
		Balance: sdk.NewInt64Coin(hostChainParams.BaseDenom, 24),
	}})
	err = lscosmosKeeper.HandleDelegationCallback(ctx, delegationResponse, icqtypes.Query{})
	suite.NoError(err)
	//slashed
	suite.Equal(lscosmosKeeper.GetHostAccountDelegation(ctx, valAddrStr1).Amount, sdk.NewInt64Coin(hostChainParams.BaseDenom, 24))
}
