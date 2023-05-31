package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/gogoproto/proto"
	"github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"

	lscosmostypes "github.com/persistenceOne/pstake-native/v2/x/lscosmos/types"
)

// GenerateAndExecuteICATx does ica transactions with messages,
// optimistic bool does not check for channel to be open. only use to do icatxns when channel is getting created.
func (k Keeper) GenerateAndExecuteICATx(ctx sdk.Context, connectionID string, ownerID string, msgs []proto.Message) error {

	msgData, err := icatypes.SerializeCosmosTx(k.cdc, msgs)
	if err != nil {
		k.Logger(ctx).Error(fmt.Sprintf("could not serialize cosmostx err %v", err))
		return err
	}

	icaPacketData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: msgData,
	}

	msg := &types.MsgSendTx{
		Owner:           ownerID,
		ConnectionId:    connectionID,
		PacketData:      icaPacketData,
		RelativeTimeout: uint64(lscosmostypes.ICATimeoutTimestamp.Nanoseconds()),
	}
	handler := k.msgRouter.Handler(msg)

	res, err := handler(ctx, msg)
	if err != nil {
		k.Logger(ctx).Error(fmt.Sprintf("send ica txn of msgs: %s failed with err: %v", msgs, err))
		return errorsmod.Wrapf(lscosmostypes.ErrICATxFailure, "Failed to send ica msgs with err: %v", err)
	}
	ctx.EventManager().EmitEvents(res.GetEvents())

	for _, msgResponse := range res.MsgResponses {
		var parsedMsgResponse types.MsgSendTxResponse
		if err := k.cdc.Unmarshal(msgResponse.Value, &parsedMsgResponse); err != nil {
			return errorsmod.Wrapf(sdkerrors.ErrJSONUnmarshal, "cannot unmarshal ica sendtx response message: %s", err.Error())
		}
		k.Logger(ctx).Info(fmt.Sprintf("sent ICA transactions with seq: %v,  connectionID: %s, ownerID: %s, msgs: %s", parsedMsgResponse.Sequence, connectionID, ownerID, msgs))
	}

	return nil
}

// CheckPendingICATxs checks if there are any ongoing ica transaction which are stuck
func (k Keeper) CheckPendingICATxs(ctx sdk.Context) (bool, error) {
	hostChainParams := k.GetHostChainParams(ctx)
	hostAccounts := k.GetHostAccounts(ctx)
	delegationChannelID, ok := k.icaControllerKeeper.GetOpenActiveChannel(ctx, hostChainParams.ConnectionID, hostAccounts.DelegatorAccountPortID())
	if !ok {
		return true, errorsmod.Wrapf(channeltypes.ErrChannelNotFound, "PortID: %s, connectionID: %s", hostAccounts.DelegatorAccountPortID(), hostChainParams.ConnectionID)
	}
	delegationNextSendSeq, ok := k.channelKeeper.GetNextSequenceSend(ctx, hostAccounts.DelegatorAccountPortID(), delegationChannelID)
	if !ok {
		return true, errorsmod.Wrapf(channeltypes.ErrSequenceSendNotFound, "PortID: %s, channelID: %s", hostAccounts.DelegatorAccountPortID(), delegationChannelID)
	}
	delegationNextAckSeq, ok := k.channelKeeper.GetNextSequenceAck(ctx, hostAccounts.DelegatorAccountPortID(), delegationChannelID)
	if !ok {
		return true, errorsmod.Wrapf(channeltypes.ErrSequenceAckNotFound, "PortID: %s, channelID: %s", hostAccounts.DelegatorAccountPortID(), delegationChannelID)
	}
	if delegationNextSendSeq != delegationNextAckSeq {
		return true, errorsmod.Wrapf(channeltypes.ErrPacketSequenceOutOfOrder, "PortID: %s, channelID: %s, NextSendSequence: %v, NextAckSequence: %v", hostAccounts.DelegatorAccountPortID(), delegationChannelID, delegationNextSendSeq, delegationNextAckSeq)
	}
	rewardsChannelID, ok := k.icaControllerKeeper.GetOpenActiveChannel(ctx, hostChainParams.ConnectionID, hostAccounts.RewardsAccountPortID())
	if !ok {
		return true, errorsmod.Wrapf(channeltypes.ErrChannelNotFound, "PortID: %s, connectionID: %s", hostAccounts.RewardsAccountPortID(), hostChainParams.ConnectionID)
	}
	rewardsNextSendSeq, ok := k.channelKeeper.GetNextSequenceSend(ctx, hostAccounts.RewardsAccountPortID(), rewardsChannelID)
	if !ok {
		return true, errorsmod.Wrapf(channeltypes.ErrSequenceSendNotFound, "PortID: %s, channelID: %s", hostAccounts.RewardsAccountPortID(), rewardsChannelID)
	}
	rewardsNextAckSeq, ok := k.channelKeeper.GetNextSequenceAck(ctx, hostAccounts.RewardsAccountPortID(), rewardsChannelID)
	if !ok {
		return true, errorsmod.Wrapf(channeltypes.ErrSequenceAckNotFound, "PortID: %s, channelID: %s", hostAccounts.RewardsAccountPortID(), rewardsChannelID)
	}
	if rewardsNextSendSeq != rewardsNextAckSeq {
		return true, errorsmod.Wrapf(channeltypes.ErrPacketSequenceOutOfOrder, "PortID: %s, channelID: %s, NextSendSequence: %v, NextAckSequence: %v", hostAccounts.RewardsAccountPortID(), rewardsChannelID, rewardsNextSendSeq, rewardsNextAckSeq)
	}
	return false, nil
}
