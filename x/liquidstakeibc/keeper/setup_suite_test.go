package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	ibctesting "github.com/cosmos/ibc-go/v7/testing"
	"github.com/stretchr/testify/suite"

	"github.com/persistenceOne/pstake-native/v2/app"
	"github.com/persistenceOne/pstake-native/v2/app/helpers"
	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

var (
	HostDenom        = "uatom"
	MintDenom        = "stk/uatom"
	MinDeposit       = sdk.NewInt(5)
	PstakeFeeAddress = "persistence1xruvjju28j0a5ud5325rfdak8f5a04h0s30mld"
	// TestVersion defines a reusable interchainaccounts version string for testing purposes
	TestVersion = string(icatypes.ModuleCdc.MustMarshalJSON(&icatypes.Metadata{
		Version:                icatypes.Version,
		ControllerConnectionId: ibctesting.FirstConnectionID,
		HostConnectionId:       ibctesting.FirstConnectionID,
		Encoding:               icatypes.EncodingProtobuf,
		TxType:                 icatypes.TxTypeSDKMultiMsg,
	}))
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
	chainA      *ibctesting.TestChain //pstake chain
	chainB      *ibctesting.TestChain //host chain, run tests of active chains
	chainC      *ibctesting.TestChain //host chain 2, run tests of to activate chains

	transferPathAB *ibctesting.Path // chainA - chainB transfer path
	transferPathAC *ibctesting.Path // chainA - chainC transfer path
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (suite *IntegrationTestSuite) SetupTest() {

	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 0)

	ibctesting.DefaultTestingAppInit = helpers.SetupTestingApp
	sdk.DefaultBondDenom = "uxprt"
	suite.chainA = ibctesting.NewTestChain(suite.T(), suite.coordinator, ibctesting.GetChainID(1))

	ibctesting.DefaultTestingAppInit = ibctesting.SetupTestingApp
	sdk.DefaultBondDenom = HostDenom
	suite.chainB = ibctesting.NewTestChain(suite.T(), suite.coordinator, ibctesting.GetChainID(2))
	sdk.DefaultBondDenom = "uosmo"
	suite.chainC = ibctesting.NewTestChain(suite.T(), suite.coordinator, ibctesting.GetChainID(3))

	suite.coordinator.Chains = map[string]*ibctesting.TestChain{
		ibctesting.GetChainID(1): suite.chainA,
		ibctesting.GetChainID(2): suite.chainB,
		ibctesting.GetChainID(3): suite.chainC,
	}

	suite.transferPathAB = NewTransferPath(suite.chainA, suite.chainB)
	suite.coordinator.Setup(suite.transferPathAB)

	suite.transferPathAC = NewTransferPath(suite.chainA, suite.chainC)
	suite.coordinator.Setup(suite.transferPathAC)

	suite.app = suite.chainA.App.(*app.PstakeApp)
	suite.ctx = suite.chainA.GetContext()

	suite.SetupHostChainAB()
	suite.SetupICAChannelsAB()

	suite.Transfer(suite.transferPathAB, "uatom")
	suite.Transfer(suite.transferPathAC, "uosmo")

	suite.CleanupSetup()
}

func (suite *IntegrationTestSuite) CleanupSetup() {
	pstakeApp, ctx := suite.app, suite.ctx
	params := pstakeApp.LiquidStakeIBCKeeper.GetParams(ctx)
	params.AdminAddress = suite.chainA.SenderAccount.GetAddress().String()
	suite.app.LiquidStakeIBCKeeper.SetParams(ctx, params)

	epoch := suite.app.EpochsKeeper.GetEpochInfo(ctx, types.DelegationEpoch).CurrentEpoch
	pstakeApp.LiquidStakeIBCKeeper.DepositWorkflow(ctx, epoch)
}
func (suite *IntegrationTestSuite) TestOneFullFlow() {
	pstakeApp, ctx := suite.app, suite.ctx
	k := pstakeApp.LiquidStakeIBCKeeper
	hc, ok := k.GetHostChain(ctx, suite.chainB.ChainID)
	suite.True(ok)
	senderAcc := suite.chainA.SenderAccount
	//user liquidstakes
	msgLiquidStake := types.NewMsgLiquidStake(sdk.NewInt64Coin(hc.IBCDenom(), 1000000), senderAcc.GetAddress())
	result, err := suite.chainA.SendMsgs(msgLiquidStake)
	suite.NotNil(result)
	suite.NoError(err)

	msgRedeem := types.NewMsgRedeem(sdk.NewInt64Coin(hc.MintDenom(), 100000), senderAcc.GetAddress())
	result, err = suite.chainA.SendMsgs(msgRedeem)
	suite.NotNil(result)
	suite.NoError(err)
}

func (suite *IntegrationTestSuite) SetupHostChainAB() {
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
		ChainId:      suite.chainB.ChainID,
		ConnectionId: suite.transferPathAB.EndpointA.ConnectionID,
		Params:       hostChainLSParams,
		HostDenom:    HostDenom,
		ChannelId:    "channel-0", //suite.transferPathAB.EndpointA.ChannelID,
		PortId:       suite.transferPathAB.EndpointA.ChannelConfig.PortID,
		DelegationAccount: &types.ICAAccount{
			Address:      "cosmos1mykw6u6dq4z7qhw9aztpk5yp8j8y5n0c6usg9faqepw83y2u4nzq2qxaxc", //gets replaced
			Balance:      sdk.Coin{Denom: HostDenom, Amount: sdk.ZeroInt()},
			Owner:        types.DefaultDelegateAccountPortOwner(suite.chainB.ChainID),
			ChannelState: types.ICAAccount_ICA_CHANNEL_CREATED,
		},
		RewardsAccount: &types.ICAAccount{
			Address:      "cosmos19dade3sxq2wqvy6fenytxmn0y3njw8r2p88cn27pj4naxcyzzs8qgxrun3", //gets replaced
			Balance:      sdk.Coin{Denom: HostDenom, Amount: sdk.ZeroInt()},
			Owner:        types.DefaultRewardsAccountPortOwner(suite.chainB.ChainID),
			ChannelState: types.ICAAccount_ICA_CHANNEL_CREATED,
		},
		Validators:      validators,
		MinimumDeposit:  MinDeposit,
		CValue:          sdk.OneDec(),
		UnbondingFactor: 4,
		Active:          true,
	}

	suite.app.LiquidStakeIBCKeeper.SetHostChain(suite.ctx, hc)
}

func NewTransferPath(chainA, chainB *ibctesting.TestChain) *ibctesting.Path {
	path := ibctesting.NewPath(chainA, chainB)
	path.EndpointA.ChannelConfig.PortID = ibctesting.TransferPort
	path.EndpointB.ChannelConfig.PortID = ibctesting.TransferPort
	path.EndpointA.ChannelConfig.Version = ibctransfertypes.Version
	path.EndpointB.ChannelConfig.Version = ibctransfertypes.Version

	return path
}
func (suite *IntegrationTestSuite) Transfer(path *ibctesting.Path, denom string) {
	coin := sdk.NewCoin(denom, sdk.NewInt(1000000000000))

	transferMsg := ibctransfertypes.NewMsgTransfer(path.EndpointB.ChannelConfig.PortID,
		path.EndpointB.ChannelID, coin, path.EndpointB.Chain.SenderAccount.GetAddress().String(),
		path.EndpointA.Chain.SenderAccount.GetAddress().String(), path.EndpointA.Chain.GetTimeoutHeight(),
		0, "")
	result, err := path.EndpointB.Chain.SendMsgs(transferMsg)
	suite.Require().NoError(err) // message committed

	packet, err := ibctesting.ParsePacketFromEvents(result.GetEvents())
	suite.Require().NoError(err)

	err = path.RelayPacket(packet)
	suite.Require().NoError(err)
}
func (suite *IntegrationTestSuite) SetupICAChannelsAB() {
	icapath := NewICAPath(suite.chainA, suite.chainB)
	icapath.EndpointA.ClientID = suite.transferPathAB.EndpointA.ClientID
	icapath.EndpointB.ClientID = suite.transferPathAB.EndpointB.ClientID
	icapath.EndpointA.ConnectionID = suite.transferPathAB.EndpointA.ConnectionID
	icapath.EndpointB.ConnectionID = suite.transferPathAB.EndpointB.ConnectionID
	icapath.EndpointA.ClientConfig = suite.transferPathAB.EndpointA.ClientConfig
	icapath.EndpointB.ClientConfig = suite.transferPathAB.EndpointB.ClientConfig
	icapath.EndpointA.ConnectionConfig = suite.transferPathAB.EndpointA.ConnectionConfig
	icapath.EndpointB.ConnectionConfig = suite.transferPathAB.EndpointB.ConnectionConfig

	icapath2 := NewICAPath(suite.chainA, suite.chainB)
	icapath2.EndpointA.ClientID = suite.transferPathAB.EndpointA.ClientID
	icapath2.EndpointB.ClientID = suite.transferPathAB.EndpointB.ClientID
	icapath2.EndpointA.ConnectionID = suite.transferPathAB.EndpointA.ConnectionID
	icapath2.EndpointB.ConnectionID = suite.transferPathAB.EndpointB.ConnectionID
	icapath2.EndpointA.ClientConfig = suite.transferPathAB.EndpointA.ClientConfig
	icapath2.EndpointB.ClientConfig = suite.transferPathAB.EndpointB.ClientConfig
	icapath2.EndpointA.ConnectionConfig = suite.transferPathAB.EndpointA.ConnectionConfig
	icapath2.EndpointB.ConnectionConfig = suite.transferPathAB.EndpointB.ConnectionConfig

	err := suite.SetupICAPath(icapath, types.DefaultDelegateAccountPortOwner(suite.chainB.ChainID))
	suite.Require().NoError(err)

	err = suite.SetupICAPath(icapath2, types.DefaultRewardsAccountPortOwner(suite.chainB.ChainID))
	suite.Require().NoError(err)
}
func NewICAPath(chainA, chainB *ibctesting.TestChain) *ibctesting.Path {
	path := ibctesting.NewPath(chainA, chainB)
	path.EndpointA.ChannelConfig.PortID = icatypes.HostPortID
	path.EndpointB.ChannelConfig.PortID = icatypes.HostPortID
	path.EndpointA.ChannelConfig.Order = channeltypes.ORDERED
	path.EndpointB.ChannelConfig.Order = channeltypes.ORDERED
	path.EndpointA.ChannelConfig.Version = TestVersion
	path.EndpointB.ChannelConfig.Version = TestVersion

	return path
}

func (suite *IntegrationTestSuite) RegisterInterchainAccount(endpoint *ibctesting.Endpoint, owner string) error {
	portID, err := icatypes.NewControllerPortID(owner)
	if err != nil {
		return err
	}

	channelSequence := suite.app.GetIBCKeeper().ChannelKeeper.GetNextChannelSequence(endpoint.Chain.GetContext())

	if err := suite.app.ICAControllerKeeper.RegisterInterchainAccount(endpoint.Chain.GetContext(), endpoint.ConnectionID, owner, TestVersion); err != nil {
		return err
	}

	// commit state changes for proof verification
	endpoint.Chain.NextBlock()

	// update port/channel ids
	endpoint.ChannelID = channeltypes.FormatChannelIdentifier(channelSequence)
	endpoint.ChannelConfig.PortID = portID
	endpoint.ChannelConfig.Version = TestVersion

	return nil
}

// SetupICAPath invokes the InterchainAccounts entrypoint and subsequent channel handshake handlers
func (suite *IntegrationTestSuite) SetupICAPath(path *ibctesting.Path, owner string) error {
	if err := suite.RegisterInterchainAccount(path.EndpointA, owner); err != nil {
		return err
	}

	if err := path.EndpointB.ChanOpenTry(); err != nil {
		return err
	}

	if err := path.EndpointA.ChanOpenAck(); err != nil {
		return err
	}

	if err := path.EndpointB.ChanOpenConfirm(); err != nil {
		return err
	}

	return nil
}
