package keeper

import (
	"encoding/json"
	"fmt"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/persistenceOne/pstake-native/v2/x/ratesync/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/gogoproto/proto"
	icacontrollertypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"

	liquidstakeibctypes "github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
	//"github.com/persistenceOne/pstake-native/v2/x/ratesync/types"
)

func (k *Keeper) GenerateAndExecuteICATx(
	ctx sdk.Context,
	connectionID string,
	ownerID string,
	messages []proto.Message,
	memo string,
) (icacontrollertypes.MsgSendTxResponse, error) {

	msgData, err := icatypes.SerializeCosmosTx(k.cdc, messages)
	if err != nil {
		k.Logger(ctx).Error(fmt.Sprintf("could not serialize tx data: %v", err))
		return icacontrollertypes.MsgSendTxResponse{}, err
	}

	icaPacketData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: msgData,
		Memo: memo,
	}

	msgSendTx := &icacontrollertypes.MsgSendTx{
		Owner:           ownerID,
		ConnectionId:    connectionID,
		PacketData:      icaPacketData,
		RelativeTimeout: uint64(liquidstakeibctypes.ICATimeoutTimestamp.Nanoseconds()),
	}

	handler := k.msgRouter.Handler(msgSendTx)
	res, err := handler(ctx, msgSendTx)
	if err != nil {
		k.Logger(ctx).Error(fmt.Sprintf("sending ica tx with msg: %s failed with err: %v", msgData, err))
		return icacontrollertypes.MsgSendTxResponse{}, errorsmod.Wrapf(liquidstakeibctypes.ErrICATxFailure, "failed to send ica msg with err: %v", err)
	}
	ctx.EventManager().EmitEvents(res.GetEvents())

	portID, err := icatypes.NewControllerPortID(ownerID)
	if err != nil {
		return icacontrollertypes.MsgSendTxResponse{}, errorsmod.Wrapf(
			liquidstakeibctypes.ErrICATxFailure,
			"failed to create portID from ownerID: %v",
			err,
		)
	}
	_, found := k.icaControllerKeeper.GetOpenActiveChannel(ctx, connectionID, portID)
	if !found {
		return icacontrollertypes.MsgSendTxResponse{}, errorsmod.Wrapf(
			liquidstakeibctypes.ErrICATxFailure,
			"failed to get ica active channel: %v",
			err,
		)
	}

	// responses length should always be 1 since we are just sending one MsgSendTx at a time
	if len(res.MsgResponses) != 1 {
		return icacontrollertypes.MsgSendTxResponse{}, errorsmod.Wrapf(
			liquidstakeibctypes.ErrInvalidResponses,
			"not enough message responses for ica tx: %v",
			err,
		)
	}

	var msgSendTxResponse icacontrollertypes.MsgSendTxResponse
	if err = k.cdc.Unmarshal(res.MsgResponses[0].Value, &msgSendTxResponse); err != nil {
		return icacontrollertypes.MsgSendTxResponse{}, errorsmod.Wrapf(
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

	return msgSendTxResponse, nil
}

func (k *Keeper) ExecuteLiquidStakeRateTx(ctx sdk.Context, feature types.LiquidStake,
	mintDenom, hostDenom string, cValue sdk.Dec, hostchainId uint64,
	connectionID string, icaAccount liquidstakeibctypes.ICAAccount) error {
	if feature.AllowsDenom(mintDenom) {
		contractMsg := types.ExecuteLiquidStakeRate{
			LiquidStakeRate: types.LiquidStakeRate{
				DefaultBondDenom:    hostDenom,
				StkDenom:            mintDenom,
				CValue:              cValue,
				ControllerChainTime: ctx.BlockTime(),
			},
		}
		contractBz, err := json.Marshal(contractMsg)
		if err != nil {
			return err
		}
		msg := &wasmtypes.MsgExecuteContract{
			Sender:   icaAccount.Address,
			Contract: feature.ContractAddress,
			Msg:      contractBz,
			Funds:    nil,
		}
		memo := types.ICAMemo{
			FeatureType: feature.FeatureType,
			HostChainId: hostchainId,
		}
		memoBz, err := json.Marshal(memo)
		if err != nil {
			return err

		}
		_, err = k.GenerateAndExecuteICATx(ctx, connectionID, icaAccount.Owner, []proto.Message{msg}, string(memoBz))
		if err != nil {
			return err
		}
	}
	return nil
}

func (k *Keeper) InstantiateLiquidStakeContract(ctx sdk.Context, icaAccount liquidstakeibctypes.ICAAccount,
	feature types.LiquidStake, id uint64, connectionID string) error {
	// generate contract msg{msg}
	contractMsg := types.InstantiateLiquidStakeRateContract{
		Admin: icaAccount.Address,
	}
	contractMsgBz, err := json.Marshal(contractMsg)
	if err != nil {
		return errorsmod.Wrapf(err, "unable to marshal InstantiateLiquidStakeRateContract")
	}

	msg := &wasmtypes.MsgInstantiateContract{
		Sender: icaAccount.Address,
		Admin:  icaAccount.Address,
		CodeID: feature.CodeID,
		Label:  fmt.Sprintf("PSTAKE ratesync, ID-%v", id),
		Msg:    contractMsgBz,
		Funds:  sdk.Coins{},
	}
	memo := types.ICAMemo{
		FeatureType: feature.FeatureType,
		HostChainId: id,
	}
	memobz, err := json.Marshal(memo)
	if err != nil {
		return err
	}
	_, err = k.GenerateAndExecuteICATx(ctx, connectionID, icaAccount.Owner, []proto.Message{msg}, string(memobz))
	if err != nil {
		return err
	}
	return nil
}
