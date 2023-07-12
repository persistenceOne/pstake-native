package keeper_test

import (
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	connectiontypes "github.com/cosmos/ibc-go/v7/modules/core/03-connection/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	commitmenttypes "github.com/cosmos/ibc-go/v7/modules/core/23-commitment/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	solomachine "github.com/cosmos/ibc-go/v7/modules/light-clients/06-solomachine"
	"github.com/stretchr/testify/require"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/testutil"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	ibctesting "github.com/cosmos/ibc-go/v7/testing"
	"github.com/stretchr/testify/suite"

	"github.com/persistenceOne/pstake-native/v2/app"
	"github.com/persistenceOne/pstake-native/v2/app/helpers"
	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

var (
	//ChainID          = "cosmoshub-4"
	//ConnectionID     = "connection-0"
	//TransferChannel  = "channel-0"
	//TransferPort     = "transfer"
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
	chainA      *ibctesting.TestChain
	chainB      *ibctesting.TestChain
	path        *ibctesting.Path
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (suite *IntegrationTestSuite) SetupTest() {

	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2)
	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(1))
	suite.chainB = suite.coordinator.GetChain(ibctesting.GetChainID(2))

	suite.path = NewTransferPath(suite.chainA, suite.chainB)
	suite.coordinator.Setup(suite.path)

	suite.app = suite.chainA.App.(*app.PstakeApp)
	suite.ctx = suite.chainA.GetContext()

	keeper := suite.app.LiquidStakeIBCKeeper

	params := types.DefaultParams()
	keeper.SetParams(suite.ctx, params)

	suite.SetupHostChain()

}
func (suite *IntegrationTestSuite) SetupHostChain() {
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
		Validators:      validators,
		MinimumDeposit:  MinDeposit,
		CValue:          sdk.OneDec(),
		UnbondingFactor: 4,
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
func (suite *IntegrationTestSuite) SetupICAChannels() *ibctesting.Path {
	icapath := NewICAPath(suite.chainA, suite.chainB)
	icapath.EndpointA.ClientID = suite.path.EndpointA.ClientID
	icapath.EndpointB.ClientID = suite.path.EndpointB.ClientID
	icapath.EndpointA.ConnectionID = suite.path.EndpointA.ConnectionID
	icapath.EndpointB.ConnectionID = suite.path.EndpointB.ConnectionID
	icapath.EndpointA.ClientConfig = suite.path.EndpointA.ClientConfig
	icapath.EndpointB.ClientConfig = suite.path.EndpointB.ClientConfig
	icapath.EndpointA.ConnectionConfig = suite.path.EndpointA.ConnectionConfig
	icapath.EndpointB.ConnectionConfig = suite.path.EndpointB.ConnectionConfig

	err := suite.SetupICAPath(icapath, types.DefaultDelegateAccountPortOwner(suite.chainB.ChainID))
	suite.Require().NoError(err)

	return icapath
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
func (suite *IntegrationTestSuite) TestGetSetParams() {
	tc := []struct {
		name     string
		params   types.Params
		expected types.Params
	}{
		{
			name: "normal params",
			params: types.Params{
				AdminAddress:     "persistence10khgeppewe4rgfrcy809r9h00aquwxxxrk6glr",
				FeeAddress:       "persistence1xruvjju28j0a5ud5325rfdak8f5a04h0s30mld",
				UpperCValueLimit: decFromStr("1.1"),
				LowerCValueLimit: decFromStr("0.9"),
			},
			expected: types.Params{
				AdminAddress:     "persistence10khgeppewe4rgfrcy809r9h00aquwxxxrk6glr",
				FeeAddress:       "persistence1xruvjju28j0a5ud5325rfdak8f5a04h0s30mld",
				UpperCValueLimit: decFromStr("1.1"),
				LowerCValueLimit: decFromStr("0.9"),
			},
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
	hc, found := pstakeApp.LiquidStakeIBCKeeper.GetHostChain(ctx, suite.chainB.ChainID)
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
	hc, found := pstakeApp.LiquidStakeIBCKeeper.GetHostChain(ctx, suite.chainB.ChainID)
	suite.Require().Equal(found, true)

	suite.Require().Equal(
		hc.DelegationAccount.Owner,
		hc.ChainId+"."+types.DelegateICAType,
	)
}

func (suite *IntegrationTestSuite) TestRewardsAccountPortOwner() {
	pstakeApp, ctx := suite.app, suite.ctx
	hc, found := pstakeApp.LiquidStakeIBCKeeper.GetHostChain(ctx, suite.chainB.ChainID)
	suite.Require().Equal(found, true)

	suite.Require().Equal(
		hc.RewardsAccount.Owner,
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

func (suite *IntegrationTestSuite) TestGetClientState() {
	pstakeApp, ctx := suite.app, suite.ctx

	// check client state
	state, err := pstakeApp.LiquidStakeIBCKeeper.GetClientState(ctx, suite.path.EndpointA.ConnectionID)
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), ibcexported.Tendermint, state.ClientType())

	// check localhost client exists
	state, err = pstakeApp.LiquidStakeIBCKeeper.GetClientState(ctx, ibcexported.LocalhostConnectionID)
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), ibcexported.Localhost, state.ClientType())

	//no connection found
	_, err = pstakeApp.LiquidStakeIBCKeeper.GetClientState(ctx, "connection-1")
	require.Error(suite.T(), err)

	// set connection without an active client-id
	pstakeApp.IBCKeeper.ConnectionKeeper.SetConnection(ctx, "connection-1", connectiontypes.ConnectionEnd{ClientId: "client-1"})
	_, err = pstakeApp.LiquidStakeIBCKeeper.GetClientState(ctx, "connection-1")
	require.Error(suite.T(), err)
}

func (suite *IntegrationTestSuite) TestGetChainID() {
	pstakeApp, ctx := suite.app, suite.ctx

	chainID, err := pstakeApp.LiquidStakeIBCKeeper.GetChainID(ctx, suite.path.EndpointA.ConnectionID)
	suite.Require().NoError(err)
	suite.Require().Equal(suite.chainB.ChainID, chainID)

	chainID, err = pstakeApp.LiquidStakeIBCKeeper.GetChainID(ctx, ibcexported.LocalhostConnectionID)
	suite.Require().NoError(err)
	suite.Require().Equal(suite.chainA.ChainID, chainID)

	// random type of client not supported
	solomachine.RegisterInterfaces(pstakeApp.InterfaceRegistry())
	pstakeApp.IBCKeeper.ClientKeeper.SetClientState(ctx, "client-1", &solomachine.ClientState{ConsensusState: &solomachine.ConsensusState{}})
	pstakeApp.IBCKeeper.ConnectionKeeper.SetConnection(ctx, "connection-1", connectiontypes.NewConnectionEnd(connectiontypes.OPEN, "client-1", connectiontypes.NewCounterparty("--", "--", commitmenttypes.NewMerklePrefix([]byte("New"))), nil, 1))
	_, err = pstakeApp.LiquidStakeIBCKeeper.GetChainID(ctx, "connection-1")
	suite.Require().Error(err)

	//connection not found
	_, err = pstakeApp.LiquidStakeIBCKeeper.GetChainID(ctx, "connection-2")
	suite.Require().Error(err)

}

func (suite *IntegrationTestSuite) TestGetPortID() {
	portID := suite.app.LiquidStakeIBCKeeper.GetPortID("owner")
	suite.Require().Equal(icatypes.ControllerPortPrefix+"owner", portID)
}

func (suite *IntegrationTestSuite) TestRegisterICAAccount() {
	pstakeApp, ctx := suite.app, suite.ctx
	err := pstakeApp.LiquidStakeIBCKeeper.RegisterICAAccount(ctx, suite.path.EndpointA.ConnectionID, types.DefaultDelegateAccountPortOwner(suite.chainB.ChainID))
	suite.Require().NoError(err)
}

func (suite *IntegrationTestSuite) TestSetWithdrawAddress() {
	pstakeApp, ctx := suite.app, suite.ctx
	hc, found := pstakeApp.LiquidStakeIBCKeeper.GetHostChain(ctx, suite.chainB.ChainID)
	require.Equal(suite.T(), true, found)
	require.NotNil(suite.T(), hc)

	_ = suite.SetupICAChannels()
	err := pstakeApp.LiquidStakeIBCKeeper.SetWithdrawAddress(ctx, hc)
	require.NoError(suite.T(), err)
}
