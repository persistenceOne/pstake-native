package main

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	abci "github.com/cometbft/cometbft/abci/types"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	types2 "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	commitmenttypes "github.com/cosmos/ibc-go/v7/modules/core/23-commitment/types"
	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	tmclient "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	"github.com/persistenceOne/pstake-native/relay/configs"
	lensclient "github.com/strangelove-ventures/lens/client"
	lensquery "github.com/strangelove-ventures/lens/client/query"
	"go.uber.org/zap"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	home := os.ExpandEnv("data")
	keyhome := fmt.Sprintf("%s/%s", home, "keys")
	ctx := context.Background()
	persistenceClient, err := configs.GetPersistenceClient(home, keyhome)
	if err != nil {
		panic(err)
	}
	if persistenceClient.KeyExists(persistenceClient.Config.Key) {
		err = persistenceClient.DeleteKey(persistenceClient.Config.Key)
		if err != nil {
			panic(err)
		}
	}
	signerkey, err := persistenceClient.KeyAddOrRestore(persistenceClient.Config.Key, 118, "apology bargain author wood window mosquito air peanut skin visual intact swamp urge capital dawn squeeze focus medal oil skin wedding trip never key")
	//persistence1r9c7gml7njgz4r34jp4yyssdww3el9aueh22uj
	if err != nil {
		panic(err)
	}

	channelID := "channel-66"
	portID := "icacontroller-lscosmos_pstake_delegation_account"
	cosmoschannelID := "channel-490"
	cosmosPortID := "icahost"
	seq := uint64(325)
	q := sendPacketQuery(channelID, portID, seq)
	page := 1
	perpage := 1000
	packet := channeltypes.Packet{}
	txFoundheight := int64(0)
	blockSearch, err := persistenceClient.RPCClient.BlockSearch(ctx, q, &page, &perpage, "")
	if err != nil {
		panic(err)
	}
	fmt.Println(len(blockSearch.Blocks))
	br0 := blockSearch.Blocks[0]
	txFoundheight = br0.Block.Height
	br, err := persistenceClient.RPCClient.BlockResults(ctx, &txFoundheight)
	if err != nil {
		panic(err)
	}
	packet, err = ParsePacketFromEvents(br.BeginBlockEvents, seq)
	if err != nil {
		packet, err = ParsePacketFromEvents(br.EndBlockEvents, seq)
		if err != nil {
			panic(err)
		}
	}

	//txSearch, err := persistenceClient.RPCClient.TxSearch(ctx, q, true, &page, &perpage, "")
	//if err != nil {
	//	panic(err)
	//}
	//if len(txSearch.Txs) > 0 {
	//	tx0 := txSearch.Txs[0]
	//  txFoundheight = tx0.txFoundheight
	//	packet, err = ParsePacketFromEvents(tx0.TxResult.Events)
	//	if err != nil {
	//		panic(err)
	//	}
	//}else {
	//	panic("NO TXS, try block search")
	//}

	srctmClientID := "07-tendermint-36"
	//destClientID := "07-tendermint-391"
	cosmosClient, err := configs.GetCosmosClient(home, keyhome)
	if err != nil {
		panic(err)
	}
	cosmosLatestBlock, err := cosmosClient.RPCClient.Block(ctx, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(cosmosLatestBlock.Block.Height)

	//header, err := getHeader(ctx, cosmosClient, persistenceClient, srctmClientID, txFoundheight, true)
	header, err := getHeader(ctx, cosmosClient, persistenceClient, srctmClientID, cosmosLatestBlock.Block.Height-10, true)
	if err != nil {
		panic(err)
	}
	anyHeader, err := clienttypes.PackClientMessage(header)
	if err != nil {
		panic(err)
	}
	updateClientMsg := &clienttypes.MsgUpdateClient{
		ClientId:      srctmClientID,
		ClientMessage: anyHeader,
		Signer:        signerkey.Address,
	}
	fmt.Println(updateClientMsg)

	persistenceLatestBlock, err := persistenceClient.RPCClient.Block(ctx, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(persistenceLatestBlock.Block.Height)
	pffHeight := header.TrustedHeight.RevisionHeight
	//pffHeight := persistenceLatestBlock.Block.Height
	key := host.NextSequenceRecvKey(cosmosPortID, cosmoschannelID)
	value, proofBz, proofHeight, err := QueryTendermintProof(cosmosClient, ctx, int64(pffHeight), key)
	if err != nil {
		panic(err)
	}
	//check if next sequence receive exists
	if len(value) == 0 {
		panic("ERROR: len(value) = 0")
	}
	sequence := binary.BigEndian.Uint64(value)

	timeoutMsg := &types2.MsgTimeout{
		Packet:           packet,
		ProofUnreceived:  proofBz,
		ProofHeight:      proofHeight,
		NextSequenceRecv: sequence,
		Signer:           signerkey.Address,
	}
	fmt.Println(timeoutMsg)
	//channelkey := host.ChannelKey(portID, channelID)
	channelkey := host.ChannelKey(cosmosPortID, cosmoschannelID)
	_, proofCloseBz, _, err := QueryTendermintProof(cosmosClient, ctx, int64(pffHeight), channelkey)
	if err != nil {
		panic(err)
	}
	//check if next sequence receive exists
	if len(value) == 0 {
		panic("ERROR: len(value) = 0")
	}
	timeoutOnCloseMsg := &types2.MsgTimeoutOnClose{
		Packet:           packet,
		ProofUnreceived:  proofBz,
		ProofClose:       proofCloseBz,
		ProofHeight:      proofHeight,
		NextSequenceRecv: sequence,
		Signer:           signerkey.Address,
	}
	fmt.Println(timeoutOnCloseMsg)
	////header, err := getHeader(ctx, cosmosClient, persistenceClient, srctmClientID, txFoundheight, true)
	//header2, err := getHeader(ctx, persistenceClient, cosmosClient, destClientID, txFoundheight, true)
	//if err != nil {
	//	panic(err)
	//}
	//anyHeader2, err := clienttypes.PackClientMessage(header2)
	//if err != nil {
	//	panic(err)
	//}
	//msg2 := &clienttypes.MsgUpdateClient{
	//	ClientId:      destClientID,
	//	ClientMessage: anyHeader2,
	//	Signer:        signerkey.Address,
	//}
	//fmt.Println(msg2)
	txresp, err := persistenceClient.SendMsgs(ctx, []sdk.Msg{updateClientMsg, timeoutOnCloseMsg}, "lets see")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(txresp)
	//txresp, err = persistenceClient.SendMsgs(ctx, []sdk.Msg{timeoutMsg}, "lets see")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(txresp)

}
func sendPacketQuery(channelID string, portID string, seq uint64) string {
	spTag := "send_packet"
	x := []string{
		fmt.Sprintf("%s.packet_src_channel='%s'", spTag, channelID),
		fmt.Sprintf("%s.packet_sequence='%d'", spTag, seq),
	}
	return strings.Join(x, " AND ")
}
func abciToSdkEvents(events []abci.Event) []sdk.StringEvent {
	var newEvts []sdk.StringEvent
	for _, evt := range events {
		newEvts = append(newEvts, parseBase64Event(nil, evt))
	}
	return newEvts
}
func parseBase64Event(log *zap.Logger, event abci.Event) sdk.StringEvent {
	evt := sdk.StringEvent{Type: event.Type}
	for _, attr := range event.Attributes {
		key, err := base64.StdEncoding.DecodeString(attr.Key)
		if err != nil {
			log.Error("Failed to decode legacy key as base64", zap.String("base64", attr.Key), zap.Error(err))
			continue
		}
		value, err := base64.StdEncoding.DecodeString(attr.Value)
		if err != nil {
			log.Error("Failed to decode legacy value as base64", zap.String("base64", attr.Value), zap.Error(err))
			continue
		}
		evt.Attributes = append(evt.Attributes, sdk.Attribute{
			Key:   string(key),
			Value: string(value),
		})
	}
	return evt
}

// ParsePacketFromEvents parses events emitted from a MsgRecvPacket and returns the
// acknowledgement.
func ParsePacketFromEvents(events []abci.Event, matchSeq uint64) (channeltypes.Packet, error) {
Main:
	for _, ev := range events {
		if ev.Type == channeltypes.EventTypeSendPacket {
			ev2 := parseBase64Event(nil, ev)
			packet := channeltypes.Packet{}
			for _, attr := range ev2.Attributes {
				switch attr.Key {
				case channeltypes.AttributeKeyData: //nolint:staticcheck // DEPRECATED
					packet.Data = []byte(attr.Value)

				case channeltypes.AttributeKeySequence:
					seq, err := strconv.ParseUint(attr.Value, 10, 64)
					if err != nil {
						return channeltypes.Packet{}, err
					}
					if matchSeq != seq {
						continue Main
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
						return channeltypes.Packet{}, err
					}

					packet.TimeoutHeight = height

				case channeltypes.AttributeKeyTimeoutTimestamp:
					timestamp, err := strconv.ParseUint(attr.Value, 10, 64)
					if err != nil {
						return channeltypes.Packet{}, err
					}

					packet.TimeoutTimestamp = timestamp

				default:
					continue
				}
			}

			return packet, nil
		}
	}
	return channeltypes.Packet{}, fmt.Errorf("acknowledgement event attribute not found")
}
func QueryClientStateABCI(c *lensclient.ChainClient,
	ctx context.Context, height int64, clientID string,
) (*clienttypes.QueryClientStateResponse, error) {
	key := host.FullClientStateKey(clientID)

	value, proofBz, proofHeight, err := QueryTendermintProof(c, ctx, height, key)
	if err != nil {
		return nil, err
	}

	// check if client exists
	if len(value) == 0 {
		return nil, sdkerrors.Wrap(clienttypes.ErrClientNotFound, clientID)
	}

	cdc := codec.NewProtoCodec(c.Codec.InterfaceRegistry)

	clientState, err := clienttypes.UnmarshalClientState(cdc, value)
	if err != nil {
		return nil, err
	}

	anyClientState, err := clienttypes.PackClientState(clientState)
	if err != nil {
		return nil, err
	}

	clientStateRes := clienttypes.NewQueryClientStateResponse(anyClientState, proofBz, proofHeight)
	return clientStateRes, nil
}

func QueryTendermintProof(c *lensclient.ChainClient, ctx context.Context, height int64, key []byte) ([]byte, []byte, clienttypes.Height, error) {

	// ABCI queries at heights 1, 2 or less than or equal to 0 are not supported.
	// Base app does not support queries for height less than or equal to 1.
	// Therefore, a query at height 2 would be equivalent to a query at height 3.
	// A height of 0 will query with the lastest state.
	if height != 0 && height <= 2 {
		return nil, nil, clienttypes.Height{}, fmt.Errorf("proof queries at height <= 2 are not supported")
	}

	// Use the IAVL height if a valid tendermint height is passed in.
	// A height of 0 will query with the latest state.
	if height != 0 {
		height--
	}

	req := abci.RequestQuery{
		Path:   fmt.Sprintf("store/%s/key", ibcexported.StoreKey),
		Height: height,
		Data:   key,
		Prove:  true,
	}

	res, err := c.QueryABCI(ctx, req)
	if err != nil {
		return nil, nil, clienttypes.Height{}, err
	}

	merkleProof, err := commitmenttypes.ConvertProofs(res.ProofOps)
	if err != nil {
		return nil, nil, clienttypes.Height{}, err
	}

	cdc := codec.NewProtoCodec(c.Codec.InterfaceRegistry)

	proofBz, err := cdc.Marshal(&merkleProof)
	if err != nil {
		return nil, nil, clienttypes.Height{}, err
	}

	revision := clienttypes.ParseChainID(c.Config.ChainID)
	return res.Value, proofBz, clienttypes.NewHeight(revision, uint64(res.Height)+1), nil
}

func getHeader(ctx context.Context, client, submitClient *lensclient.ChainClient, clientId string, requestHeight int64, historicOk bool) (*tmclient.Header, error) {
	submitQuerier := lensquery.Query{Client: submitClient, Options: lensquery.DefaultOptions()}
	state, err := submitQuerier.Ibc_ClientState(clientId) // pass in from request
	if err != nil {
		return nil, fmt.Errorf("error: Could not get state from chain: %q ", err.Error())
	}
	unpackedState, err := clienttypes.UnpackClientState(state.ClientState)
	if err != nil {
		return nil, fmt.Errorf("error: Could not unpack state from chain: %q ", err.Error())
	}

	trustedHeight := unpackedState.GetLatestHeight()
	clientHeight, ok := trustedHeight.(clienttypes.Height)
	if !ok {
		return nil, fmt.Errorf("error: Could coerce trusted height")

	}

	if !historicOk && clientHeight.RevisionHeight >= uint64(requestHeight+1) {
		return nil, fmt.Errorf("trusted height >= request height")
	}

	newBlock, err := retryLightblock(ctx, client, requestHeight+1)
	if err != nil {
		panic(fmt.Sprintf("Error: Could not fetch updated LC from chain - bailing: %v", err))
	}

	trustedBlock, err := retryLightblock(ctx, client, int64(clientHeight.RevisionHeight)+1)
	if err != nil {
		panic(fmt.Sprintf("Error: Could not fetch updated LC from chain - bailing (2): %v", err))
	}

	valSet := tmtypes.NewValidatorSet(newBlock.ValidatorSet.Validators)
	trustedValSet := tmtypes.NewValidatorSet(trustedBlock.ValidatorSet.Validators)
	protoVal, err := valSet.ToProto()
	if err != nil {
		panic(fmt.Sprintf("Error: Could not get valset from chain: %v", err))
	}
	trustedProtoVal, err := trustedValSet.ToProto()
	if err != nil {
		panic(fmt.Sprintf("Error: Could not get trusted valset from chain: %v", err))
	}

	header := &tmclient.Header{
		SignedHeader:      newBlock.SignedHeader.ToProto(),
		ValidatorSet:      protoVal,
		TrustedHeight:     clientHeight,
		TrustedValidators: trustedProtoVal,
	}

	return header, nil
}
func retryLightblock(ctx context.Context, client *lensclient.ChainClient, height int64) (*tmtypes.LightBlock, error) {
	lightBlock, err := client.LightProvider.LightBlock(ctx, height)
	if err != nil {
		for i := 0; i < 4; i++ {
			time.Sleep(time.Duration(6) * time.Second)
			lightBlock, err = client.LightProvider.LightBlock(ctx, height)
			if err == nil {
				break
			}
		}
	}
	return lightBlock, err
}
