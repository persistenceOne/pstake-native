package keeper_test

import (
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/gogoproto/proto"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"

	"github.com/persistenceOne/pstake-native/v2/x/ratesync/keeper"
	"github.com/persistenceOne/pstake-native/v2/x/ratesync/types"
)

func (suite *IntegrationTestSuite) TestOnChanOpenAck() {
	k, ctx := suite.app.RatesyncKeeper, suite.ctx
	_ = createNChain(k, ctx, 2)
	hc, _ := k.GetHostChain(ctx, 1)
	suite.Require().NoError(k.OnChanOpenAck(ctx, types.MustICAPortIDFromOwner(hc.ICAAccount.Owner),
		suite.ratesyncPathAB.EndpointA.ChannelID, "", ""))
}

func (suite *IntegrationTestSuite) TestOnAcknowledgementPacket() {
	k, ctx := suite.app.RatesyncKeeper, suite.ctx
	_ = createNChain(k, ctx, 2)
	hc, _ := k.GetHostChain(ctx, 1)
	// case 1, instantiate msg.
	{
		msg, memo, err := keeper.GenerateInstantiateLiquidStakeContractMsg(hc.ICAAccount, hc.Features.LiquidStakeIBC, hc.ID)
		suite.Require().NoError(err)
		msgData, err := icatypes.SerializeCosmosTx(suite.app.AppCodec(), []proto.Message{msg})
		suite.Require().NoError(err)
		data := &icatypes.InterchainAccountPacketData{
			Type: 0,
			Data: msgData,
			Memo: string(memo),
		}
		databz, err := suite.app.AppCodec().MarshalJSON(data)
		suite.Require().NoError(err)
		packet := channeltypes.Packet{
			SourcePort: types.MustICAPortIDFromOwner(hc.ICAAccount.Owner),
			Data:       databz,
		}
		wasmResponse := &wasmtypes.MsgInstantiateContractResponse{
			Address: authtypes.NewModuleAddress("Contract").String(),
			Data:    nil,
		}

		msgResult, _ := codectypes.NewAnyWithValue(wasmResponse)

		result := &sdk.TxMsgData{MsgResponses: []*codectypes.Any{msgResult}}
		resultbz, err := suite.app.AppCodec().Marshal(result)
		suite.Require().NoError(err)
		ack := channeltypes.NewResultAcknowledgement(resultbz)
		ackbz, err := suite.app.AppCodec().MarshalJSON(&ack)
		suite.Require().NoError(k.OnAcknowledgementPacket(ctx, packet,
			ackbz, authtypes.NewModuleAddress("test")))
	}
	// case 2 instantiate failure
	{
		msg, memo, err := keeper.GenerateInstantiateLiquidStakeContractMsg(hc.ICAAccount, hc.Features.LiquidStakeIBC, hc.ID)
		suite.Require().NoError(err)
		msgData, err := icatypes.SerializeCosmosTx(suite.app.AppCodec(), []proto.Message{msg})
		suite.Require().NoError(err)
		data := &icatypes.InterchainAccountPacketData{
			Type: 0,
			Data: msgData,
			Memo: string(memo),
		}
		databz, err := suite.app.AppCodec().MarshalJSON(data)
		suite.Require().NoError(err)
		packet := channeltypes.Packet{
			SourcePort: types.MustICAPortIDFromOwner(hc.ICAAccount.Owner),
			Data:       databz,
		}

		ack := channeltypes.NewErrorAcknowledgement(types.ErrICATxFailure)
		ackbz, err := suite.app.AppCodec().MarshalJSON(&ack)
		suite.Require().NoError(k.OnAcknowledgementPacket(ctx, packet,
			ackbz, authtypes.NewModuleAddress("test")))
	}
	// case 3, execute msg.
	{
		msg, memo, err := keeper.GenerateExecuteLiquidStakeRateTxMsg(ctx.BlockTime().Unix(), hc.Features.LiquidStake,
			"stk/uatom", "uatom", sdk.OneDec(), hc.ID, hc.ICAAccount)
		suite.Require().NoError(err)
		msgData, err := icatypes.SerializeCosmosTx(suite.app.AppCodec(), []proto.Message{msg})
		suite.Require().NoError(err)
		data := &icatypes.InterchainAccountPacketData{
			Type: 0,
			Data: msgData,
			Memo: string(memo),
		}
		databz, err := suite.app.AppCodec().MarshalJSON(data)
		suite.Require().NoError(err)
		packet := channeltypes.Packet{
			SourcePort: types.MustICAPortIDFromOwner(hc.ICAAccount.Owner),
			Data:       databz,
		}
		wasmResponse := &wasmtypes.MsgExecuteContractResponse{Data: []byte{}}

		msgResult, _ := codectypes.NewAnyWithValue(wasmResponse)

		result := &sdk.TxMsgData{MsgResponses: []*codectypes.Any{msgResult}}
		resultbz, err := suite.app.AppCodec().Marshal(result)
		suite.Require().NoError(err)
		ack := channeltypes.NewResultAcknowledgement(resultbz)
		ackbz, err := suite.app.AppCodec().MarshalJSON(&ack)
		suite.Require().NoError(k.OnAcknowledgementPacket(ctx, packet,
			ackbz, authtypes.NewModuleAddress("test")))
	}
	// case 4 execute failure
	{
		msg, memo, err := keeper.GenerateExecuteLiquidStakeRateTxMsg(ctx.BlockTime().Unix(), hc.Features.LiquidStake,
			"stk/uatom", "uatom", sdk.OneDec(), hc.ID, hc.ICAAccount)
		suite.Require().NoError(err)
		msgData, err := icatypes.SerializeCosmosTx(suite.app.AppCodec(), []proto.Message{msg})
		suite.Require().NoError(err)
		data := &icatypes.InterchainAccountPacketData{
			Type: 0,
			Data: msgData,
			Memo: string(memo),
		}
		databz, err := suite.app.AppCodec().MarshalJSON(data)
		suite.Require().NoError(err)
		packet := channeltypes.Packet{
			SourcePort: types.MustICAPortIDFromOwner(hc.ICAAccount.Owner),
			Data:       databz,
		}

		ack := channeltypes.NewErrorAcknowledgement(types.ErrICATxFailure)
		ackbz, err := suite.app.AppCodec().MarshalJSON(&ack)
		suite.Require().NoError(k.OnAcknowledgementPacket(ctx, packet,
			ackbz, authtypes.NewModuleAddress("test")))
	}
}

func (suite *IntegrationTestSuite) TestOnTimeoutPacket() {
	k, ctx := suite.app.RatesyncKeeper, suite.ctx
	_ = createNChain(k, ctx, 2)
	hc, _ := k.GetHostChain(ctx, 1)
	// case 1, instantiate msg.
	{
		msg, memo, err := keeper.GenerateInstantiateLiquidStakeContractMsg(hc.ICAAccount, hc.Features.LiquidStakeIBC, hc.ID)
		suite.Require().NoError(err)
		msgData, err := icatypes.SerializeCosmosTx(suite.app.AppCodec(), []proto.Message{msg})
		suite.Require().NoError(err)
		data := &icatypes.InterchainAccountPacketData{
			Type: 0,
			Data: msgData,
			Memo: string(memo),
		}
		databz, err := suite.app.AppCodec().MarshalJSON(data)
		suite.Require().NoError(err)
		packet := channeltypes.Packet{
			SourcePort: types.MustICAPortIDFromOwner(hc.ICAAccount.Owner),
			Data:       databz,
		}
		suite.Require().NoError(k.OnTimeoutPacket(ctx, packet, authtypes.NewModuleAddress("test")))
	}
	// case 2, execute msg.
	{
		msg, memo, err := keeper.GenerateExecuteLiquidStakeRateTxMsg(ctx.BlockTime().Unix(), hc.Features.LiquidStake,
			"stk/uatom", "uatom", sdk.OneDec(), hc.ID, hc.ICAAccount)
		suite.Require().NoError(err)
		msgData, err := icatypes.SerializeCosmosTx(suite.app.AppCodec(), []proto.Message{msg})
		suite.Require().NoError(err)
		data := &icatypes.InterchainAccountPacketData{
			Type: 0,
			Data: msgData,
			Memo: string(memo),
		}
		databz, err := suite.app.AppCodec().MarshalJSON(data)
		suite.Require().NoError(err)
		packet := channeltypes.Packet{
			SourcePort: types.MustICAPortIDFromOwner(hc.ICAAccount.Owner),
			Data:       databz,
		}
		suite.Require().NoError(k.OnTimeoutPacket(ctx, packet, authtypes.NewModuleAddress("test")))
	}
}
