package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	icatypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v3/modules/core/05-port/types"
	host "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	ibcexported "github.com/cosmos/ibc-go/v3/modules/core/exported"

	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

// OnChanOpenInit implements the IBCModule interface
func (k Keeper) OnChanOpenInit(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID string,
	channelID string,
	chanCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	version string,
) error {

	// Require portID is the portID module is bound to
	if portID != types.DelegationAccountPortID &&
		portID != types.RewardAccountPortID {
		return sdkerrors.Wrapf(porttypes.ErrInvalidPort, "invalid port: %s, expected either of %s or %s",
			portID, types.DelegationAccountPortID, types.RewardAccountPortID)
	}
	var versionData icatypes.Metadata
	if err := icatypes.ModuleCdc.UnmarshalJSON([]byte(version), &versionData); err != nil {
		return err
	}
	if versionData.Version != icatypes.Version {
		return sdkerrors.Wrapf(types.ErrInvalidVersion, "got %s, expected %s", versionData.Version, icatypes.Version)
	}

	// Claim channel capability passed back by IBC module
	if err := k.ClaimCapability(ctx, chanCap, host.ChannelCapabilityPath(portID, channelID)); err != nil {
		return err
	}

	return nil
}

// OnChanOpenTry implements the IBCModule interface
func (k Keeper) OnChanOpenTry(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID,
	channelID string,
	chanCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	counterpartyVersion string,
) (string, error) {
	// Controller Auth Module does not do OnChanOpenTry
	return "", nil
}

// OnChanOpenAck implements the IBCModule interface
func (k Keeper) OnChanOpenAck(
	ctx sdk.Context,
	portID,
	channelID string,
	counterpartyChannelID string,
	counterpartyVersion string,
) error {
	var counterpartyVersionData icatypes.Metadata
	if err := icatypes.ModuleCdc.UnmarshalJSON([]byte(counterpartyVersion), &counterpartyVersionData); err != nil {
		return err
	}

	if counterpartyVersionData.Version != icatypes.Version {
		return sdkerrors.Wrapf(types.ErrInvalidVersion, "invalid counterparty version: %s, expected %s", counterpartyVersion, types.Version)
	}
	//TODO more checks

	hostchainparams := k.GetCosmosIBCParams(ctx)

	if portID == types.DelegationAccountPortID {
		address, found := k.icaControllerKeeper.GetInterchainAccountAddress(ctx, hostchainparams.IBCConnection, portID)
		if !found {
			ctx.Logger().Error(fmt.Sprintf("expected to find an address for %s/%s", hostchainparams.IBCConnection, portID))
			return icatypes.ErrInterchainAccountNotFound
		}
		if err := k.SetHostChainDelegationAddress(ctx, address); err != nil {
			return err
		}
	}

	// On Ack we enable module and it's transactions
	// TODO add checks if both delegation and rewards ica exists only then enable module
	k.SetModuleState(ctx, true)

	return nil
}

// OnChanOpenConfirm implements the IBCModule interface
func (k Keeper) OnChanOpenConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return nil
}

// OnChanCloseInit implements the IBCModule interface
func (k Keeper) OnChanCloseInit(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	// Disallow user-initiated channel closing for channels
	return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "user cannot close channel")
}

// OnChanCloseConfirm implements the IBCModule interface
func (k Keeper) OnChanCloseConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return nil
}

// OnRecvPacket implements the IBCModule interface
func (k Keeper) OnRecvPacket(
	ctx sdk.Context,
	modulePacket channeltypes.Packet,
	relayer sdk.AccAddress,
) ibcexported.Acknowledgement {
	// Controller Auth Module does not do OnRecvPacket
	return nil
}

// OnAcknowledgementPacket implements the IBCModule interface
func (k Keeper) OnAcknowledgementPacket(
	ctx sdk.Context,
	modulePacket channeltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {
	var ack channeltypes.Acknowledgement
	if err := types.ModuleCdc.UnmarshalJSON(acknowledgement, &ack); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal packet acknowledgement: %v", err)
	}
	if !ack.Success() {
		return sdkerrors.Wrapf(channeltypes.ErrInvalidAcknowledgement, "acknowledgement failed")
	}
	// this line is used by starport scaffolding # oracle/packet/module/ack

	txMsgData := &sdk.TxMsgData{}
	if err := k.cdc.Unmarshal(ack.GetResult(), txMsgData); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-27 tx message data: %v", err)
	}

	icaPacket := &icatypes.InterchainAccountPacketData{}
	if err := icatypes.ModuleCdc.UnmarshalJSON(modulePacket.GetData(), icaPacket); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-27 tx message data: %v", err)
	}
	msgs, err := icatypes.DeserializeCosmosTx(k.cdc, icaPacket.GetData())
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot Deserialise icapacket data: %v", err)
	}
	var eventType string

	// Dispatch packet
	switch len(txMsgData.Data) {
	case 0:
		// TODO: handle for sdk 0.46.x
		return nil
	default:
		for i, msgData := range txMsgData.Data {
			response, err := k.handleAckMsgData(ctx, msgData, msgs[i])
			if err != nil {
				return err
			}

			k.Logger(ctx).Info("message response in ICS-27 packet response", "response", response)
		}
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			eventType,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyAck, fmt.Sprintf("%v", ack)),
		),
	)

	switch resp := ack.Response.(type) {
	case *channeltypes.Acknowledgement_Result:
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				eventType,
				sdk.NewAttribute(types.AttributeKeyAckSuccess, string(resp.Result)),
			),
		)
	case *channeltypes.Acknowledgement_Error:
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				eventType,
				sdk.NewAttribute(types.AttributeKeyAckError, resp.Error),
			),
		)
	}

	return nil
}

// OnTimeoutPacket implements the IBCModule interface
func (k Keeper) OnTimeoutPacket(
	ctx sdk.Context,
	modulePacket channeltypes.Packet,
	relayer sdk.AccAddress,
) error {
	// this line is used by starport scaffolding # oracle/packet/module/ack

	icaPacket := &icatypes.InterchainAccountPacketData{}
	if err := icatypes.ModuleCdc.UnmarshalJSON(modulePacket.GetData(), icaPacket); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-27 tx message data: %v", err)
	}
	msgs, err := icatypes.DeserializeCosmosTx(k.cdc, icaPacket.GetData())
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot Deserialise icapacket data: %v", err)
	}
	// Dispatch packet
	switch len(icaPacket.Data) {
	case 0:
		// TODO: handle for sdk 0.46.x
		return nil
	default:
		for _, msg := range msgs {
			response, err := k.handleTimeoutMsgData(ctx, msg)
			if err != nil {
				return err
			}

			k.Logger(ctx).Info("message response in ICS-27 packet response", "response", response)
		}
	}

	return nil
}

func (k Keeper) handleAckMsgData(ctx sdk.Context, msgData *sdk.MsgData, msg sdk.Msg) (string, error) {
	switch msgData.MsgType {
	case sdk.MsgTypeURL(&stakingtypes.MsgDelegate{}):
		parsedMsg := msg.(*stakingtypes.MsgDelegate)
		var msgResponse stakingtypes.MsgDelegateResponse
		if err := k.cdc.Unmarshal(msgData.Data, &msgResponse); err != nil {
			return "", sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, "cannot unmarshal send response message: %s", err.Error())
		}
		// remove from host-balance
		k.RemoveBalanceToDelegationState(ctx, sdk.NewCoins(parsedMsg.Amount))
		// Add delegation state

		return msgResponse.String(), nil

	// TODO: handle other messages

	default:
		return "", nil
	}
}

func (k Keeper) handleTimeoutMsgData(_ sdk.Context, msg sdk.Msg) (string, error) {
	switch sdk.MsgTypeURL(msg) {
	case sdk.MsgTypeURL(&stakingtypes.MsgDelegate{}):
		return msg.String(), nil
	case sdk.MsgTypeURL(&stakingtypes.MsgBeginRedelegate{}):
		return msg.String(), sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "Not implemented, unexpected msg %s", msg.String())
	default:
		return "", nil
	}
}
