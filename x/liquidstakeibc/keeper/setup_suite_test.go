package keeper_test

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/gogoproto/proto"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
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

	delegationPathAB *ibctesting.Path // chainA - chain B delegation ica path
	rewardsPathAB    *ibctesting.Path // chainA - chainB rewards ica path
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

	suite.SetupHostChainAB()
	suite.SetupICAChannelsAB()

	suite.Transfer(suite.transferPathAB, sdk.NewCoin("uatom", sdk.NewInt(1000000000000)))
	suite.Transfer(suite.transferPathAC, sdk.NewCoin("uosmo", sdk.NewInt(1000000000000)))

	suite.SetupLSM()

	suite.CleanupSetup()
	suite.ctx = suite.chainA.GetContext()
}

func (suite *IntegrationTestSuite) CleanupSetup() {
	pstakeApp := suite.app

	params := pstakeApp.LiquidStakeIBCKeeper.GetParams(suite.chainA.GetContext())
	params.AdminAddress = suite.chainA.SenderAccount.GetAddress().String()
	suite.app.LiquidStakeIBCKeeper.SetParams(suite.chainA.GetContext(), params)

	epoch := suite.app.EpochsKeeper.GetEpochInfo(suite.chainA.GetContext(), types.DelegationEpoch).CurrentEpoch
	pstakeApp.LiquidStakeIBCKeeper.DepositWorkflow(suite.chainA.GetContext(), epoch)
}
func (suite *IntegrationTestSuite) ResetEpochs() {
	ctx := suite.chainA.GetContext()

	//ctxCheck := app.BaseApp.NewContext(true, tmproto.Header{})
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

	lsmValidatorCap, err := sdk.NewDecFromStr("0.5")
	suite.NoError(err)

	lsmBondFactor, err := sdk.NewDecFromStr("250")
	suite.NoError(err)

	hostChainLSParams := &types.HostChainLSParams{
		DepositFee:      depositFee,
		RestakeFee:      restakeFee,
		UnstakeFee:      unstakeFee,
		RedemptionFee:   redemptionFee,
		LsmValidatorCap: lsmValidatorCap,
		LsmBondFactor:   lsmBondFactor,
	}

	validators := make([]*types.Validator, 0)
	equalWeight := sdk.OneDec().Quo(sdk.NewDecFromInt(sdk.NewInt(int64(len(suite.chainB.Vals.Validators)))))
	for _, validator := range suite.chainB.Vals.Validators {
		validators = append(validators, &types.Validator{
			OperatorAddress: sdk.MustBech32ifyAddressBytes(app.Bech32PrefixValAddr, validator.Address),
			Status:          stakingtypes.Bonded.String(),
			Weight:          equalWeight,
			DelegatedAmount: sdk.ZeroInt(),
			ExchangeRate:    sdk.OneDec(),
			Delegable:       true,
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
		Validators:         validators,
		MinimumDeposit:     MinDeposit,
		CValue:             sdk.OneDec(),
		UnbondingFactor:    4,
		AutoCompoundFactor: suite.app.LiquidStakeIBCKeeper.CalculateAutocompoundLimit(sdk.NewDec(20)),
		Active:             true,
		Flags:              &types.HostChainFlags{Lsm: true},
	}

	suite.app.LiquidStakeIBCKeeper.SetHostChain(suite.chainA.GetContext(), hc)
}

func (suite *IntegrationTestSuite) SetupLSM() {
	// delegate from an address on Cosmos
	msgDelegate := &stakingtypes.MsgDelegate{
		DelegatorAddress: suite.chainB.SenderAccount.GetAddress().String(),
		ValidatorAddress: sdk.MustBech32ifyAddressBytes(app.Bech32PrefixValAddr, suite.chainB.Vals.Validators[0].Address),
		Amount:           sdk.NewCoin(HostDenom, sdk.NewInt(1000000000)),
	}
	res, err := suite.chainB.SendMsgs(msgDelegate)
	suite.Require().NotNil(res)
	suite.NoError(err)

	// tokenize the whole delegation
	msgTokenizeShares := &stakingtypes.MsgTokenizeShares{
		ValidatorAddress:    sdk.MustBech32ifyAddressBytes(app.Bech32PrefixValAddr, suite.chainB.Vals.Validators[0].Address),
		DelegatorAddress:    suite.chainB.SenderAccount.GetAddress().String(),
		Amount:              sdk.NewCoin(HostDenom, sdk.NewInt(1000000000)),
		TokenizedShareOwner: suite.chainB.SenderAccount.GetAddress().String(),
	}
	res, err = suite.chainB.SendMsgs(msgTokenizeShares)
	suite.Require().NotNil(res)
	suite.NoError(err)

	// send it via IBC to persistence
	suite.Transfer(
		suite.transferPathAB,
		sdk.NewCoin(
			sdk.MustBech32ifyAddressBytes(app.Bech32PrefixValAddr, suite.chainB.Vals.Validators[0].Address)+"/1",
			sdk.NewInt(1000),
		),
	)
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
	suite.delegationPathAB = icapath

	icapath2 := NewICAPath(suite.chainA, suite.chainB)
	icapath2.EndpointA.ClientID = suite.transferPathAB.EndpointA.ClientID
	icapath2.EndpointB.ClientID = suite.transferPathAB.EndpointB.ClientID
	icapath2.EndpointA.ConnectionID = suite.transferPathAB.EndpointA.ConnectionID
	icapath2.EndpointB.ConnectionID = suite.transferPathAB.EndpointB.ConnectionID
	icapath2.EndpointA.ClientConfig = suite.transferPathAB.EndpointA.ClientConfig
	icapath2.EndpointB.ClientConfig = suite.transferPathAB.EndpointB.ClientConfig
	icapath2.EndpointA.ConnectionConfig = suite.transferPathAB.EndpointA.ConnectionConfig
	icapath2.EndpointB.ConnectionConfig = suite.transferPathAB.EndpointB.ConnectionConfig
	suite.rewardsPathAB = icapath2

	err := suite.SetupICAPath(suite.delegationPathAB, types.DefaultDelegateAccountPortOwner(suite.chainB.ChainID))
	suite.Require().NoError(err)

	err = suite.SetupICAPath(suite.rewardsPathAB, types.DefaultRewardsAccountPortOwner(suite.chainB.ChainID))
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

func (suite *IntegrationTestSuite) TestOneFullFlow() {
	pstakeApp := suite.app
	k := pstakeApp.LiquidStakeIBCKeeper
	hc, ok := k.GetHostChain(suite.chainA.GetContext(), suite.chainB.ChainID)
	suite.True(ok)

	epoch := pstakeApp.EpochsKeeper.GetEpochInfo(suite.chainA.GetContext(), types.DelegationEpoch)
	suite.NotNil(epoch)
	err := k.BeforeEpochStart(suite.chainA.GetContext(), epoch.Identifier, epoch.CurrentEpoch)
	suite.Require().NoError(err)

	senderAcc := suite.chainA.SenderAccount
	//user liquidstakes
	msgLiquidStake := types.NewMsgLiquidStake(sdk.NewInt64Coin(hc.IBCDenom(), 1000000), senderAcc.GetAddress())
	result, err := suite.app.MsgServiceRouter().Handler(msgLiquidStake)(suite.chainA.GetContext(), msgLiquidStake)
	suite.NotNil(result)
	suite.NoError(err)

	// user redeems
	msgRedeem := types.NewMsgRedeem(sdk.NewInt64Coin(hc.MintDenom(), 100000), senderAcc.GetAddress())
	result, err = suite.app.MsgServiceRouter().Handler(msgRedeem)(suite.chainA.GetContext(), msgRedeem)
	suite.NotNil(result)
	suite.NoError(err)

	// Do ica staking
	epoch = pstakeApp.EpochsKeeper.GetEpochInfo(suite.chainA.GetContext(), types.DelegationEpoch)
	suite.NotNil(epoch)

	deposit, found := k.GetDepositForChainAndEpoch(suite.chainA.GetContext(), hc.ChainId, epoch.CurrentEpoch)
	suite.Require().True(found)
	suite.Require().Equal(types.Deposit_DEPOSIT_PENDING, deposit.State)

	ctx := suite.chainA.GetContext() //use separate context so we can fetch events out of it
	err = k.AfterEpochEnd(ctx, epoch.Identifier, epoch.CurrentEpoch)
	suite.NoError(err)
	packets, err := ParsePacketsFromEvents(ctx.EventManager().Events())
	suite.NoError(err)

	deposit, found = k.GetDepositForChainAndEpoch(suite.chainA.GetContext(), hc.ChainId, epoch.CurrentEpoch)
	suite.Require().True(found)
	suite.Require().Equal(types.Deposit_DEPOSIT_SENT, deposit.State)

	suite.chainA.NextBlock() //commit the packets and their commitments so its available in context

	suite.RelayAllPacketsAB(packets) //also calls for Next Block causing Deposit_DEPOSIT_RECEIVED to just pass

	deposit, found = k.GetDepositForChainAndEpoch(suite.chainA.GetContext(), hc.ChainId, epoch.CurrentEpoch)
	suite.Require().True(found)
	suite.Require().Equal(types.Deposit_DEPOSIT_DELEGATING, deposit.State)

	timeoutTimestamp := uint64(suite.chainA.GetContext().BlockTime().UnixNano()) + uint64(types.ICATimeoutTimestamp.Nanoseconds()) - uint64(time.Second*5) //sub one b
	data, err := suite.CreateICAData(deposit.Amount.Amount, hc, 0)
	suite.NoError(err)

	packet, err := CreateICADelegatePacketHardcoded(data,
		"1", "0-0", fmt.Sprintf("%d", timeoutTimestamp))
	suite.NoError(err)
	suite.chainA.NextBlock() //commit the packets and their commitments so its available in context
	suite.RelayAllPacketsAB([]channeltypes.Packet{packet})
	deposit, found = k.GetDepositForChainAndEpoch(suite.chainA.GetContext(), hc.ChainId, epoch.CurrentEpoch)
	suite.Require().False(found)
	suite.Require().Nil(deposit)
	// ^ Fin staking

	// Do unstake
	undelegateAmount := int64(100000)
	msgUnstake := types.NewMsgLiquidUnstake(sdk.NewInt64Coin(hc.MintDenom(), undelegateAmount), senderAcc.GetAddress())
	result, err = suite.app.MsgServiceRouter().Handler(msgUnstake)(suite.chainA.GetContext(), msgUnstake)
	suite.NotNil(result)
	suite.NoError(err)

	epoch = pstakeApp.EpochsKeeper.GetEpochInfo(suite.chainA.GetContext(), types.DelegationEpoch)
	suite.NotNil(epoch)

	userUnbonding, found := k.GetUserUnbonding(suite.chainA.GetContext(), hc.ChainId, senderAcc.GetAddress().String(), types.CurrentUnbondingEpoch(hc.UnbondingFactor, epoch.CurrentEpoch))
	suite.Require().True(found)
	suite.Require().NotNil(userUnbonding)

	unbonding, found := k.GetUnbonding(suite.chainA.GetContext(), hc.ChainId, types.CurrentUnbondingEpoch(hc.UnbondingFactor, epoch.CurrentEpoch))
	suite.Require().True(found)
	suite.Require().Equal(types.Unbonding_UNBONDING_PENDING, unbonding.State)
	//Force undelegation by manipulating epoch number
	ctx = suite.chainA.GetContext()
	err = k.AfterEpochEnd(ctx, epoch.Identifier, types.CurrentUnbondingEpoch(hc.UnbondingFactor, epoch.CurrentEpoch))
	packets, err = ParsePacketsFromEvents(ctx.EventManager().Events())
	suite.NoError(err)
	unbonding, found = k.GetUnbonding(suite.chainA.GetContext(), hc.ChainId, types.CurrentUnbondingEpoch(hc.UnbondingFactor, epoch.CurrentEpoch))
	suite.Require().True(found)
	suite.Require().Equal(types.Unbonding_UNBONDING_INITIATED, unbonding.State)

	suite.chainA.NextBlock() //commit the packets and their commitments so its available in context

	suite.RelayAllPacketsAB(packets)
	unbonding, found = k.GetUnbonding(suite.chainA.GetContext(), hc.ChainId, types.CurrentUnbondingEpoch(hc.UnbondingFactor, epoch.CurrentEpoch))
	suite.Require().True(found)
	suite.Require().Equal(types.Unbonding_UNBONDING_MATURING, unbonding.State)
}

func (suite *IntegrationTestSuite) RelayAllPacketsAB(packets []channeltypes.Packet) {
	suite.Require().NotEqual(0, len(packets), "No packets to relay")
	hc, _ := suite.app.LiquidStakeIBCKeeper.GetHostChain(suite.chainA.GetContext(), suite.chainB.ChainID)
	for _, packet := range packets {
		if packet.SourcePort == hc.PortId {
			err := suite.transferPathAB.RelayPacket(packet)
			suite.Require().NoError(err)
		} else if packet.SourcePort == suite.app.LiquidStakeIBCKeeper.GetPortID(hc.DelegationAccount.Owner) {
			err := suite.delegationPathAB.RelayPacket(packet)
			suite.Require().NoError(err)
		} else if packet.SourcePort == suite.app.LiquidStakeIBCKeeper.GetPortID(hc.RewardsAccount.Owner) {
			err := suite.rewardsPathAB.RelayPacket(packet)
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

func (suite *IntegrationTestSuite) CreateICAData(amount math.Int, hc *types.HostChain, txtype int) (string, error) {

	messages := make([]proto.Message, 0)
	for _, vals := range hc.Validators {
		var message proto.Message
		switch txtype {
		case 0:
			message = &stakingtypes.MsgDelegate{
				DelegatorAddress: hc.DelegationAccount.Address,
				ValidatorAddress: vals.OperatorAddress,
				Amount:           sdk.NewCoin(hc.HostDenom, vals.Weight.MulInt(amount).TruncateInt()),
			}
		case 1:
		case 2:
			message = &distributiontypes.MsgWithdrawDelegatorReward{
				DelegatorAddress: hc.DelegationAccount.Address,
				ValidatorAddress: vals.OperatorAddress,
			}

		}

		messages = append(messages, message)
	}
	msgData, err := icatypes.SerializeCosmosTx(suite.app.AppCodec(), messages)
	if err != nil {
		return "", err
	}

	icaPacketData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: msgData,
	}
	return string(icaPacketData.GetBytes()), nil
}

func CreateICADelegatePacketHardcoded(data, sequence, timeoutHeight, timeoutTimestamp string) (channeltypes.Packet, error) {
	seq, err := strconv.ParseUint(sequence, 10, 64)
	if err != nil {
		return channeltypes.Packet{}, err
	}
	height, err := clienttypes.ParseHeight(timeoutHeight)
	if err != nil {
		return channeltypes.Packet{}, err
	}

	timestamp, err := strconv.ParseUint(timeoutTimestamp, 10, 64)
	if err != nil {
		return channeltypes.Packet{}, err
	}
	packet := channeltypes.Packet{
		Sequence:           seq,
		SourcePort:         "icacontroller-testchain2-1.delegate",
		SourceChannel:      "channel-2",
		DestinationPort:    "icahost",
		DestinationChannel: "channel-1",
		Data:               []byte(data),
		TimeoutHeight:      height,
		TimeoutTimestamp:   timestamp,
	}
	return packet, nil
}
func CreateICARewardsPacketHardcoded(data, sequence, timeoutHeight, timeoutTimestamp string) (channeltypes.Packet, error) {
	seq, err := strconv.ParseUint(sequence, 10, 64)
	if err != nil {
		return channeltypes.Packet{}, err
	}
	height, err := clienttypes.ParseHeight(timeoutHeight)
	if err != nil {
		return channeltypes.Packet{}, err
	}

	timestamp, err := strconv.ParseUint(timeoutTimestamp, 10, 64)
	if err != nil {
		return channeltypes.Packet{}, err
	}
	packet := channeltypes.Packet{
		Sequence:           seq,
		SourcePort:         "icacontroller-testchain2-1.rewards",
		SourceChannel:      "channel-3",
		DestinationPort:    "icahost",
		DestinationChannel: "channel-2",
		Data:               []byte(data),
		TimeoutHeight:      height,
		TimeoutTimestamp:   timestamp,
	}
	return packet, nil
}

func (suite *IntegrationTestSuite) UpdateChainActive(active bool, hc *types.HostChain) {
	hc.Active = active
	suite.app.LiquidStakeIBCKeeper.SetHostChain(suite.ctx, hc)
}

func (suite *IntegrationTestSuite) UpdateChainLSMActive(active bool, hc *types.HostChain) {
	hc.Flags.Lsm = active
	suite.app.LiquidStakeIBCKeeper.SetHostChain(suite.ctx, hc)
}
