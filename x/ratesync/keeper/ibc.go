package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"encoding/json"
	"fmt"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	liquidstakeibctypes "github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
	"strconv"

	"github.com/persistenceOne/pstake-native/v2/x/ratesync/types"
)

func (k *Keeper) OnChanOpenInit(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID string,
	channelID string,
	channelCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	version string,
) (string, error) {
	return version, nil
}

func (k *Keeper) OnChanOpenAck(
	ctx sdk.Context,
	portID string,
	channelID string,
	counterpartyChannelID string,
	counterpartyVersion string,
) error {
	// get the connection id from the port and channel identifiers
	connID, _, err := k.ibcKeeper.ChannelKeeper.GetChannelConnection(ctx, portID, channelID)
	if err != nil {
		return fmt.Errorf("unable to get connection id using port %s: %w", portID, err)
	}

	// get interchain account address
	address, found := k.icaControllerKeeper.GetInterchainAccountAddress(ctx, connID, portID)
	if !found {
		return fmt.Errorf("couldn't find address for %s/%s", connID, portID)
	}

	// get the port owner from the port id
	portOwner, err := types.OwnerfromPortID(portID)
	if err != nil {
		return fmt.Errorf("unable to parse port id %s, err: %v", portID, err)
	}

	// get the chain id using the connection id
	chainID, err := k.GetChainID(ctx, connID)
	if err != nil {
		return fmt.Errorf("unable to get chain id for connection %s: %w", connID, err)
	}

	id, err := types.IDfromPortID(portID)
	if err != nil {
		return err
	}
	// get host chain
	hc, found := k.GetHostChain(ctx, id)
	if !found {
		return fmt.Errorf("host chain with id %s is not registered", chainID)
	}

	switch {
	case portOwner == hc.IcaAccount.Owner:
		hc.IcaAccount.Address = address
		hc.IcaAccount.ChannelState = liquidstakeibctypes.ICAAccount_ICA_CHANNEL_CREATED
	default:
		k.Logger(ctx).Error("Unrecognised ICA account type for the module", "port-id:", portID, "chain-id", chainID)
		return nil
	}

	// save the changes of the host chain
	k.SetHostChain(ctx, hc)

	k.Logger(ctx).Info(
		"Created new ICA.",
		//"host chain",
		//ChainId,
		"channel",
		channelID,
		"owner",
		portOwner,
		"address",
		address,
	)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventICAChannelCreated,
			sdk.NewAttribute(types.AttributeChainID, hc.ChainId),
			sdk.NewAttribute(types.AttributeICAChannelID, channelID),
			sdk.NewAttribute(types.AttributeICAPortOwner, portOwner),
			sdk.NewAttribute(types.AttributeICAAddress, address),
		),
	)

	return nil
}

func (k *Keeper) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {
	var ack channeltypes.Acknowledgement
	if err := ibctransfertypes.ModuleCdc.UnmarshalJSON(acknowledgement, &ack); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal packet acknowledgement: %v", err)
	}

	var icaPacket icatypes.InterchainAccountPacketData
	if err := icatypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &icaPacket); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-27 packet data: %v", err)
	}

	var icaMemo types.ICAMemo
	err := json.Unmarshal([]byte(icaPacket.Memo), &icaMemo)
	if err != nil {
		return err
	}
	switch resp := ack.Response.(type) {
	case *channeltypes.Acknowledgement_Error:
		err := k.handleUnsuccessfulAck(ctx, icaPacket, packet, icaMemo)
		if err != nil {
			return err
		}
		k.Logger(ctx).Info(fmt.Sprintln("ICS-27 tx failed with ack:", ack.String()))
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypePacket,
				sdk.NewAttribute(types.AttributeKeyAckError, resp.Error),
			),
		)
	case *channeltypes.Acknowledgement_Result:
		err := k.handleSuccessfulAck(ctx, ack, icaPacket, packet, icaMemo)
		if err != nil {
			return err
		}
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypePacket,
				sdk.NewAttribute(types.AttributeKeyAckSuccess, fmt.Sprintln(ack.Success())),
			),
		)
	default:
		// the acknowledgement succeeded on the receiving chain so nothing
		// needs to be executed and no error needs to be returned
		return nil
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypePacket,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyAck, ack.String()),
		),
	)

	return nil
}

func (k *Keeper) OnTimeoutPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) error {
	var icaPacket icatypes.InterchainAccountPacketData
	if err := icatypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &icaPacket); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-27 tx message data: %v", err)
	}

	var icaMemo types.ICAMemo
	err := json.Unmarshal([]byte(icaPacket.Memo), &icaMemo)
	if err != nil {
		return err
	}

	if err := k.handleUnsuccessfulAck(ctx, icaPacket, packet, icaMemo); err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTimeout,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		),
	)

	k.Logger(ctx).Info(
		"ICA transaction timed out.",
		"sequence",
		packet.Sequence,
		"channel",
		packet.SourceChannel,
		"port",
		packet.SourcePort,
	)

	return nil
}

func (k *Keeper) handleUnsuccessfulAck(
	ctx sdk.Context,
	icaPacket icatypes.InterchainAccountPacketData,
	packet channeltypes.Packet, icaMemo types.ICAMemo,
) error {
	messages, err := icatypes.DeserializeCosmosTx(k.cdc, icaPacket.GetData())
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "cannot deserialize ica packet data: %v", err)
	}
	hc, found := k.GetHostChain(ctx, icaMemo.HostChainId)
	if !found {
		return errorsmod.Wrapf(sdkerrors.ErrNotFound, "hostchain not found for id %v", icaMemo.HostChainId)
	}
	for _, msg := range messages {
		switch sdk.MsgTypeURL(msg) {
		case sdk.MsgTypeURL(&wasmtypes.MsgInstantiateContract{}):
			// parse the MsgInstantiateContract  to emit the instantiate error event
			parsedMsg, ok := msg.(*wasmtypes.MsgInstantiateContract)
			if !ok {
				k.Logger(ctx).Error(
					"Could not parse MsgInstantiateContract while handling unsuccessful ack.",
					"channel", packet.SourceChannel, "sequence", packet.Sequence,
				)
				continue
			}
			//reset instantiation state so can be retried.
			switch icaMemo.FeatureType {
			case types.FeatureType_LIQUID_STAKE_IBC:
				hc.Features.LiquidStakeIBC.Instantiation = types.InstantiationState_INSTANTIATION_NOT_INITIATED
			case types.FeatureType_LIQUID_STAKE:
				hc.Features.LiquidStake.Instantiation = types.InstantiationState_INSTANTIATION_NOT_INITIATED
			}

			// emit an event for the instantiate confirmation
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeUnsuccessfulInstantiateContract,
					sdk.NewAttribute(types.AttributeChainID, hc.ChainId),
					sdk.NewAttribute(types.AttributeSender, parsedMsg.Sender),
					sdk.NewAttribute(channeltypes.AttributeKeySequence, strconv.FormatUint(packet.Sequence, 10)),
					sdk.NewAttribute(channeltypes.AttributeKeySrcChannel, packet.SourceChannel),
					sdk.NewAttribute(channeltypes.AttributeKeyPortID, packet.SourcePort),
				),
			)
		case sdk.MsgTypeURL(&wasmtypes.MsgExecuteContract{}):

			// parse the MsgExecuteContract to emit the execute error event
			parsedMsg, ok := msg.(*wasmtypes.MsgExecuteContract)
			if !ok {
				k.Logger(ctx).Error(
					"Could not parse MsgExecuteContract while handling unsuccessful ack.",
					"channel", packet.SourceChannel, "sequence", packet.Sequence,
				)
				continue
			}
			//Do nothing, relay next epoch

			// emit an event for the execution confirmation
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeUnsuccessfulExecuteContract,
					sdk.NewAttribute(types.AttributeChainID, hc.ChainId),
					sdk.NewAttribute(types.AttributeSender, parsedMsg.Sender),
					sdk.NewAttribute(channeltypes.AttributeKeySequence, strconv.FormatUint(packet.Sequence, 10)),
					sdk.NewAttribute(channeltypes.AttributeKeySrcChannel, packet.SourceChannel),
					sdk.NewAttribute(channeltypes.AttributeKeyPortID, packet.SourcePort),
				),
			)

		}
	}
	k.SetHostChain(ctx, hc)

	return nil
}

func (k *Keeper) handleSuccessfulAck(
	ctx sdk.Context,
	ack channeltypes.Acknowledgement,
	icaPacket icatypes.InterchainAccountPacketData,
	packet channeltypes.Packet, icaMemo types.ICAMemo,
) error {
	txMsgData := &sdk.TxMsgData{}
	if err := k.cdc.Unmarshal(ack.GetResult(), txMsgData); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ics-27 tx ack data: %v", err)
	}

	messages, err := icatypes.DeserializeCosmosTx(k.cdc, icaPacket.GetData())
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "cannot deserialize ica packet data: %v", err)
	}

	for i, msg := range messages {
		switch sdk.MsgTypeURL(msg) {
		case sdk.MsgTypeURL(&wasmtypes.MsgInstantiateContract{}):
			var data []byte
			if len(txMsgData.Data) == 0 {
				data = txMsgData.GetMsgResponses()[i].Value
			} else {
				data = txMsgData.Data[i].Data
			}

			var msgResponse wasmtypes.MsgInstantiateContractResponse
			if err := k.cdc.Unmarshal(data, &msgResponse); err != nil {
				return errorsmod.Wrapf(
					sdkerrors.ErrJSONUnmarshal, "cannot unmarshal MsgInstantiateContract response message: %s",
					err.Error(),
				)
			}
			parsedMsg, ok := msg.(*wasmtypes.MsgInstantiateContract)
			if !ok {
				return errorsmod.Wrapf(
					sdkerrors.ErrInvalidType,
					"unable to cast msg of type %s to MsgInstantiateContract",
					sdk.MsgTypeURL(msg),
				)
			}
			if err = k.HandleInstantiateContractResponse(ctx, parsedMsg, msgResponse, icaMemo); err != nil {
				return err
			}
		case sdk.MsgTypeURL(&wasmtypes.MsgExecuteContract{}):
			var data []byte
			if len(txMsgData.Data) == 0 {
				data = txMsgData.GetMsgResponses()[i].Value
			} else {
				data = txMsgData.Data[i].Data
			}

			var msgResponse wasmtypes.MsgExecuteContractResponse
			if err := k.cdc.Unmarshal(data, &msgResponse); err != nil {
				return errorsmod.Wrapf(
					sdkerrors.ErrJSONUnmarshal, "cannot unmarshal MsgExecuteContract response message: %s",
					err.Error(),
				)
			}
			parsedMsg, ok := msg.(*wasmtypes.MsgExecuteContract)
			if !ok {
				return errorsmod.Wrapf(
					sdkerrors.ErrInvalidType,
					"unable to cast msg of type %s to MsgExecuteContract",
					sdk.MsgTypeURL(msg),
				)
			}

			if err = k.HandleExecuteContractResponse(ctx, parsedMsg, msgResponse); err != nil {
				return err
			}
		}
	}

	k.Logger(ctx).Info(
		"ICA transaction ACK success.",
		"sequence",
		packet.Sequence,
		"channel",
		packet.SourceChannel,
		"messages",
		messages,
	)

	return nil
}

func (k Keeper) HandleInstantiateContractResponse(ctx sdk.Context,
	msg *wasmtypes.MsgInstantiateContract,
	resp wasmtypes.MsgInstantiateContractResponse,
	icaMemo types.ICAMemo,
) error {
	hc, found := k.GetHostChain(ctx, icaMemo.HostChainId)
	if !found {
		return errorsmod.Wrapf(sdkerrors.ErrNotFound, "hostchain not found for id %v", icaMemo.HostChainId)
	}
	switch icaMemo.FeatureType {
	case types.FeatureType_LIQUID_STAKE_IBC:
		hc.Features.LiquidStakeIBC.Instantiation = types.InstantiationState_INSTANTIATION_COMPLETED
		hc.Features.LiquidStakeIBC.ContractAddress = resp.Address
		hc.Features.LiquidStakeIBC.Enabled = true
	case types.FeatureType_LIQUID_STAKE:
		hc.Features.LiquidStake.Instantiation = types.InstantiationState_INSTANTIATION_COMPLETED
		hc.Features.LiquidStake.ContractAddress = resp.Address
		hc.Features.LiquidStake.Enabled = true
	}
	k.SetHostChain(ctx, hc)
	return nil
}

func (k Keeper) HandleExecuteContractResponse(ctx sdk.Context,
	msg *wasmtypes.MsgExecuteContract,
	resp wasmtypes.MsgExecuteContractResponse,
) error {
	// cool do nothing
	return nil
}
