package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/testutil"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	connectiontypes "github.com/cosmos/ibc-go/v6/modules/core/03-connection/types"
	ibctmtypes "github.com/cosmos/ibc-go/v6/modules/light-clients/07-tendermint/types"
	ibctesting "github.com/cosmos/ibc-go/v6/testing"
	"github.com/stretchr/testify/suite"

	"github.com/persistenceOne/pstake-native/v2/app"
	"github.com/persistenceOne/pstake-native/v2/app/helpers"
	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

var (
	ChainID          = "cosmoshub-4"
	ConnectionID     = "connection-0"
	TransferChannel  = "channel-0"
	TransferPort     = "transfer"
	HostDenom        = "uatom"
	MintDenom        = "stk/uatom"
	MinDeposit       = sdk.NewInt(5)
	PstakeFeeAddress = "persistence1xruvjju28j0a5ud5325rfdak8f5a04h0s30mld"
)

func init() {
	ibctesting.DefaultTestingAppInit = helpers.SetupTestingApp
}

type IntegrationTestSuite struct {
	suite.Suite

	app        *app.PstakeApp
	ctx        sdk.Context
	govHandler govtypes.Handler

	coordinator *ibctesting.Coordinator
	chainA      *ibctesting.TestChain
	chainB      *ibctesting.TestChain
	path        *ibctesting.Path
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (suite *IntegrationTestSuite) SetupTest() {
	_, pstakeApp, ctx := helpers.CreateTestApp(suite.T())

	keeper := pstakeApp.LiquidStakeIBCKeeper

	params := types.DefaultParams()
	keeper.SetParams(ctx, params)

	suite.app = &pstakeApp
	suite.ctx = ctx

	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2)
	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(1))
	suite.chainB = suite.coordinator.GetChain(ibctesting.GetChainID(2))

	suite.path = ibctesting.NewPath(suite.chainA, suite.chainB)
	suite.path.EndpointA.ChannelConfig.PortID = ibctesting.TransferPort
	suite.path.EndpointB.ChannelConfig.PortID = ibctesting.TransferPort
	suite.coordinator.SetupConnections(suite.path)

	// set host chain params
	depositFee, err := sdk.NewDecFromStr("0.01")
	suite.NoError(err)

	restakeFee, err := sdk.NewDecFromStr("0.02")
	suite.NoError(err)

	unstakeFee, err := sdk.NewDecFromStr("0.03")
	suite.NoError(err)

	redemptionFee, err := sdk.NewDecFromStr("0.03")
	suite.NoError(err)

	hostChainLSParams := &types.HostChainLSParams{
		DepositFee:    depositFee,
		RestakeFee:    restakeFee,
		UnstakeFee:    unstakeFee,
		RedemptionFee: redemptionFee,
	}

	validators := make([]*types.Validator, 0)
	for _, validator := range suite.chainB.Vals.Validators {
		validators = append(validators, &types.Validator{
			OperatorAddress: validator.Address.String(),
			Status:          stakingtypes.Bonded.String(),
			Weight:          sdk.ZeroDec(),
			DelegatedAmount: sdk.ZeroInt(),
		})
	}

	hc := &types.HostChain{
		ChainId:      suite.path.EndpointB.Chain.ChainID,
		ConnectionId: suite.path.EndpointA.ConnectionID,
		Params:       hostChainLSParams,
		HostDenom:    HostDenom,
		ChannelId:    "channel-0", //suite.path.EndpointA.ChannelID,
		PortId:       suite.path.EndpointA.ChannelConfig.PortID,
		DelegationAccount: &types.ICAAccount{
			Address:      "cosmos1mykw6u6dq4z7qhw9aztpk5yp8j8y5n0c6usg9faqepw83y2u4nzq2qxaxc",
			Balance:      sdk.Coin{Denom: HostDenom, Amount: sdk.ZeroInt()},
			Owner:        suite.chainB.ChainID + "." + types.DelegateICAType,
			ChannelState: types.ICAAccount_ICA_CHANNEL_CREATED,
		},
		RewardsAccount: &types.ICAAccount{
			Address:      "cosmos19dade3sxq2wqvy6fenytxmn0y3njw8r2p88cn27pj4naxcyzzs8qgxrun3",
			Balance:      sdk.Coin{Denom: HostDenom, Amount: sdk.ZeroInt()},
			Owner:        suite.chainB.ChainID + "." + types.RewardsICAType,
			ChannelState: types.ICAAccount_ICA_CHANNEL_CREATED,
		},
		Validators:     validators,
		MinimumDeposit: MinDeposit,
		CValue:         sdk.OneDec(),
		NextValsetHash: nil,
	}

	suite.app.LiquidStakeIBCKeeper.SetHostChain(suite.ctx, hc)

	pstakeApp.IBCKeeper.ClientKeeper.SetClientState(ctx, "07-tendermint-0", &ibctmtypes.ClientState{ChainId: suite.chainB.ChainID, TrustingPeriod: time.Hour, LatestHeight: clienttypes.Height{RevisionNumber: 1, RevisionHeight: 100}})
	pstakeApp.IBCKeeper.ClientKeeper.SetClientConsensusState(ctx, "07-tendermint-0", clienttypes.Height{RevisionNumber: 1, RevisionHeight: 100}, &ibctmtypes.ConsensusState{Timestamp: ctx.BlockTime()})
	pstakeApp.IBCKeeper.ConnectionKeeper.SetConnection(ctx, suite.path.EndpointA.ConnectionID, connectiontypes.ConnectionEnd{ClientId: "07-tendermint-0"})
}

func (suite *IntegrationTestSuite) TestGetSetParams() {
	tc := []struct {
		name     string
		params   types.Params
		expected types.Params
	}{
		{
			name:     "normal params",
			params:   types.Params{FeeAddress: "persistence1xruvjju28j0a5ud5325rfdak8f5a04h0s30mld"},
			expected: types.Params{FeeAddress: "persistence1xruvjju28j0a5ud5325rfdak8f5a04h0s30mld"},
		},
	}

	for _, t := range tc {
		suite.Run(t.name, func() {
			pstakeApp, ctx := suite.app, suite.ctx

			pstakeApp.LiquidStakeIBCKeeper.SetParams(ctx, t.params)
			params := pstakeApp.LiquidStakeIBCKeeper.GetParams(ctx)
			suite.Require().Equal(params, t.expected)
		})
	}
}

func (suite *IntegrationTestSuite) TestSendProtocolFee() {
	tc := []struct {
		name       string
		fee        sdk.Coins
		module     string
		feeAddress string
		success    bool
	}{
		{
			name:       "successful case",
			fee:        sdk.Coins{sdk.Coin{Denom: MintDenom, Amount: sdk.NewInt(100)}},
			module:     types.ModuleName,
			feeAddress: PstakeFeeAddress,
			success:    true,
		},
		{
			name:       "invalid fee address",
			fee:        sdk.Coins{sdk.Coin{Denom: MintDenom, Amount: sdk.NewInt(100)}},
			module:     types.ModuleName,
			feeAddress: "1234",
			success:    false,
		},
		{
			name:       "not enough tokens",
			fee:        sdk.Coins{sdk.Coin{Denom: MintDenom, Amount: sdk.NewInt(1000)}},
			module:     types.ModuleName,
			feeAddress: PstakeFeeAddress,
			success:    false,
		},
	}

	pstakeApp, ctx := suite.app, suite.ctx
	hc, found := pstakeApp.LiquidStakeIBCKeeper.GetHostChain(ctx, suite.path.EndpointB.Chain.ChainID)
	suite.Require().Equal(found, true)

	baseFee := sdk.NewInt64Coin(hc.MintDenom(), 100)
	suite.Require().NoError(testutil.FundModuleAccount(pstakeApp.BankKeeper, ctx, types.ModuleName, sdk.NewCoins(baseFee)))

	for _, t := range tc {
		suite.Run(t.name, func() {
			if t.success {
				suite.Require().NoError(
					pstakeApp.LiquidStakeIBCKeeper.SendProtocolFee(
						ctx,
						t.fee,
						t.module,
						t.feeAddress,
					),
				)

				feeAddress, _ := sdk.AccAddressFromBech32(t.feeAddress)
				currBalance := pstakeApp.BankKeeper.GetBalance(ctx, feeAddress, hc.MintDenom())
				suite.Require().Equal(baseFee, currBalance)
			} else {
				suite.Require().Error(
					pstakeApp.LiquidStakeIBCKeeper.SendProtocolFee(
						ctx,
						t.fee,
						t.module,
						t.feeAddress,
					),
				)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestDelegateAccountPortOwner() {
	pstakeApp, ctx := suite.app, suite.ctx
	hc, found := pstakeApp.LiquidStakeIBCKeeper.GetHostChain(ctx, suite.path.EndpointB.Chain.ChainID)
	suite.Require().Equal(found, true)

	suite.Require().Equal(
		pstakeApp.LiquidStakeIBCKeeper.DelegateAccountPortOwner(hc.ChainId),
		hc.ChainId+"."+types.DelegateICAType,
	)
}

func (suite *IntegrationTestSuite) TestRewardsAccountPortOwner() {
	pstakeApp, ctx := suite.app, suite.ctx
	hc, found := pstakeApp.LiquidStakeIBCKeeper.GetHostChain(ctx, suite.path.EndpointB.Chain.ChainID)
	suite.Require().Equal(found, true)

	suite.Require().Equal(
		pstakeApp.LiquidStakeIBCKeeper.RewardsAccountPortOwner(hc.ChainId),
		hc.ChainId+"."+types.RewardsICAType,
	)
}

func (suite *IntegrationTestSuite) TestGetEpochNumber() {
	pstakeApp, ctx := suite.app, suite.ctx

	suite.Require().Equal(
		pstakeApp.LiquidStakeIBCKeeper.GetEpochNumber(ctx, types.DelegationEpoch),
		pstakeApp.EpochsKeeper.GetEpochInfo(ctx, types.DelegationEpoch).CurrentEpoch,
	)
}
