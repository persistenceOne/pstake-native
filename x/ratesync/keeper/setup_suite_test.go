package keeper_test

import (
	"fmt"
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	ibctesting "github.com/cosmos/ibc-go/v7/testing"
	"github.com/stretchr/testify/suite"

	"github.com/persistenceOne/pstake-native/v2/app"
	"github.com/persistenceOne/pstake-native/v2/app/helpers"
	"github.com/persistenceOne/pstake-native/v2/x/ratesync/types"
)

var (
	HostDenom        = "uatom"
	MintDenom        = "stk/uatom"
	MinDeposit       = sdk.NewInt(5)
	PstakeFeeAddress = "persistence1xruvjju28j0a5ud5325rfdak8f5a04h0s30mld"
	GovAddress       = authtypes.NewModuleAddress("gov")
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
	chainA      *ibctesting.TestChain // pstake chain
	chainB      *ibctesting.TestChain // host chain, run tests of active chains
	chainC      *ibctesting.TestChain // host chain 2, run tests of to activate chains

	transferPathAB *ibctesting.Path // chainA - chainB transfer path
	transferPathAC *ibctesting.Path // chainA - chainC transfer path

	ratesyncPathAB *ibctesting.Path // chainA - chain B ratesync ica path
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (suite *IntegrationTestSuite) SetupTest() {
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 0)

	ibctesting.DefaultTestingAppInit = helpers.SetupTestingApp
	sdk.DefaultBondDenom = "uxprt"
	suite.chainA = ibctesting.NewTestChain(suite.T(), suite.coordinator, ibctesting.GetChainID(1))
	suite.ResetEpochs()
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

	// suite.SetupHostChainAB()
	suite.SetupICAChannelsAB()

	suite.Transfer(suite.transferPathAB, sdk.NewCoin("uatom", sdk.NewInt(1000000000000)))
	suite.Transfer(suite.transferPathAC, sdk.NewCoin("uosmo", sdk.NewInt(1000000000000)))

	// suite.SetupLSM()

	suite.CleanupSetup()
	suite.ctx = suite.chainA.GetContext()
}

func (suite *IntegrationTestSuite) CleanupSetup() {
}

func (suite *IntegrationTestSuite) ResetEpochs() {
	ctx := suite.chainA.GetContext()

	// ctxCheck := app.BaseApp.NewContext(true, tmproto.Header{})
	epochsKeeper := suite.chainA.App.(*app.PstakeApp).EpochsKeeper

	for _, epoch := range epochsKeeper.AllEpochInfos(ctx) {
		epoch.StartTime = ctx.BlockTime()
		epoch.CurrentEpoch = int64(1)
		epoch.CurrentEpochStartTime = ctx.BlockTime()
		epoch.CurrentEpochStartHeight = ctx.BlockHeight()
		epochsKeeper.DeleteEpochInfo(ctx, epoch.Identifier)
		err := epochsKeeper.AddEpochInfo(ctx, epoch)
		if err != nil {
			panic(err)
		}
	}
}

func NewTransferPath(chainA, chainB *ibctesting.TestChain) *ibctesting.Path {
	path := ibctesting.NewPath(chainA, chainB)
	path.EndpointA.ChannelConfig.PortID = ibctesting.TransferPort
	path.EndpointB.ChannelConfig.PortID = ibctesting.TransferPort
	path.EndpointA.ChannelConfig.Version = ibctransfertypes.Version
	path.EndpointB.ChannelConfig.Version = ibctransfertypes.Version

	return path
}

func (suite *IntegrationTestSuite) Transfer(path *ibctesting.Path, coin sdk.Coin) {
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
	suite.ratesyncPathAB = icapath

	err := suite.SetupICAPath(suite.ratesyncPathAB, types.DefaultPortOwner(1))
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

func (suite *IntegrationTestSuite) RelayAllPacketsAB(packets []channeltypes.Packet) {
	suite.Require().NotEqual(0, len(packets), "No packets to relay")
	hc, _ := suite.app.LiquidStakeIBCKeeper.GetHostChain(suite.chainA.GetContext(), suite.chainB.ChainID)
	for _, packet := range packets {
		if packet.SourcePort == hc.PortId {
			err := suite.transferPathAB.RelayPacket(packet)
			suite.Require().NoError(err)
		} else if packet.SourcePort == suite.app.LiquidStakeIBCKeeper.GetPortID(hc.DelegationAccount.Owner) {
			err := suite.ratesyncPathAB.RelayPacket(packet)
			suite.Require().NoError(err)
		}
	}
}

// ParsePacketsFromEvents parses events emitted from a MsgRecvPacket and returns the
// acknowledgement.
func ParsePacketsFromEvents(events sdk.Events) ([]channeltypes.Packet, error) {
	packets := make([]channeltypes.Packet, 0)
	for _, ev := range events {
		if ev.Type == channeltypes.EventTypeSendPacket {
			packet := channeltypes.Packet{}
			for _, attr := range ev.Attributes {
				switch attr.Key {
				case channeltypes.AttributeKeyData: //nolint:staticcheck // DEPRECATED
					packet.Data = []byte(attr.Value)

				case channeltypes.AttributeKeySequence:
					seq, err := strconv.ParseUint(attr.Value, 10, 64)
					if err != nil {
						return []channeltypes.Packet{}, err
					}

					packet.Sequence = seq

				case channeltypes.AttributeKeySrcPort:
					packet.SourcePort = attr.Value

				case channeltypes.AttributeKeySrcChannel:
					packet.SourceChannel = attr.Value

				case channeltypes.AttributeKeyDstPort:
					packet.DestinationPort = attr.Value

				case channeltypes.AttributeKeyDstChannel:
					packet.DestinationChannel = attr.Value

				case channeltypes.AttributeKeyTimeoutHeight:
					height, err := clienttypes.ParseHeight(attr.Value)
					if err != nil {
						return []channeltypes.Packet{}, err
					}

					packet.TimeoutHeight = height

				case channeltypes.AttributeKeyTimeoutTimestamp:
					timestamp, err := strconv.ParseUint(attr.Value, 10, 64)
					if err != nil {
						return []channeltypes.Packet{}, err
					}

					packet.TimeoutTimestamp = timestamp

				default:
					continue
				}
			}
			packets = append(packets, packet)
		}
	}
	if len(packets) == 0 {
		return []channeltypes.Packet{}, fmt.Errorf("acknowledgement event attribute not found")
	} else {
		return packets, nil
	}
}
