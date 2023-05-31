package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/gogoproto/proto"
	"github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"

	liquidstakeibctypes "github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

func (k *Keeper) GenerateAndExecuteICATx(
	ctx sdk.Context,
	connectionID string,
	ownerID string,
	messages []proto.Message,
) (string, error) {

	msgData, err := icatypes.SerializeCosmosTx(k.cdc, messages)
	if err != nil {
		k.Logger(ctx).Error(fmt.Sprintf("could not serialize tx data: %v", err))
		return "", err
	}

	icaPacketData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: msgData,
	}

	msgSendTx := &types.MsgSendTx{
		Owner:           ownerID,
		ConnectionId:    connectionID,
		PacketData:      icaPacketData,
		RelativeTimeout: uint64(liquidstakeibctypes.ICATimeoutTimestamp.Nanoseconds()),
	}

	handler := k.msgRouter.Handler(msgSendTx)
	res, err := handler(ctx, msgSendTx)
	if err != nil {
		k.Logger(ctx).Error(fmt.Sprintf("sending ica tx with msg: %s failed with err: %v", msgData, err))
		return "", errorsmod.Wrapf(liquidstakeibctypes.ErrICATxFailure, "failed to send ica msg with err: %v", err)
	}
	ctx.EventManager().EmitEvents(res.GetEvents())

	channelID, found := k.icaControllerKeeper.GetOpenActiveChannel(ctx, connectionID, k.GetPortID(ownerID))
	if !found {
		return "", errorsmod.Wrapf(
			liquidstakeibctypes.ErrICATxFailure,
			"failed to get ica active channel: %v",
			err,
		)
	}

	// responses length should always be 1 since we are just sending one MsgSendTx at a time
	if len(res.MsgResponses) != 1 {
		return "", errorsmod.Wrapf(
			liquidstakeibctypes.ErrInvalidResponses,
			"not enough message responses for ica tx: %v",
			err,
		)
	}

	var msgSendTxResponse types.MsgSendTxResponse
	if err = k.cdc.Unmarshal(res.MsgResponses[0].Value, &msgSendTxResponse); err != nil {
		return "", errorsmod.Wrapf(
			sdkerrors.ErrJSONUnmarshal,
			"cannot unmarshal ica send tx response message: %v",
			err,
		)
	}
	k.Logger(ctx).Info(
		fmt.Sprintf(
			"sent ICA transactions with seq: %v, connectionID: %s, ownerID: %s, msgs: %s",
			msgSendTxResponse.Sequence,
			connectionID,
			ownerID,
			messages,
		),
	)

	return k.GetTransactionSequenceID(channelID, msgSendTxResponse.Sequence), nil
}
