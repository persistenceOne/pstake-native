package keeper

import (
	"fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
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
	_, portOwner, found := strings.Cut(portID, icatypes.ControllerPortPrefix)
	if !found {
		return fmt.Errorf("unable to parse port id %s", portID)
	}

	// get the chain id using the connection id
	chainID, err := k.GetChainID(ctx, connID)
	if err != nil {
		return fmt.Errorf("unable to get chain id for connection %s: %w", connID, err)
	}

	// get host chain
	hc, found := k.GetHostChain(ctx, chainID)
	if !found {
		return fmt.Errorf("host chain with id %s is not registered", chainID)
	}

	switch {
	case portOwner == hc.DelegationAccount.Owner:
		hc.DelegationAccount.Address = address
		hc.DelegationAccount.Owner = portOwner
		hc.DelegationAccount.ChannelState = types.ICAAccount_ICA_CHANNEL_CREATED
	case portOwner == hc.RewardsAccount.Owner:
		hc.RewardsAccount.Address = address
		hc.RewardsAccount.Owner = portOwner
		hc.RewardsAccount.ChannelState = types.ICAAccount_ICA_CHANNEL_CREATED
	default:
		k.Logger(ctx).Error("Unrecognised ICA account type for the module", "port-id:", portID, "chain-id", chainID)
		return nil
	}

	// save the changes of the host chain
	k.SetHostChain(ctx, hc)

	// send an ICQ query to get the delegator account balance
	if hc.DelegationAccount != nil && hc.DelegationAccount.ChannelState == types.ICAAccount_ICA_CHANNEL_CREATED {
		if err := k.QueryDelegationHostChainAccountBalance(ctx, hc); err != nil {
			return fmt.Errorf(
				"error querying host chain %s for delegation account balances: %v",
				hc.ChainId,
				err,
			)
		}
	}

	k.Logger(ctx).Info(
		"Created new ICA.",
		"host chain",
		hc.ChainId,
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

	switch resp := ack.Response.(type) {
	case *channeltypes.Acknowledgement_Error:
		err := k.handleUnsuccessfulAck(ctx, icaPacket, packet.SourceChannel, packet.Sequence)
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
		err := k.handleSuccessfulAck(ctx, ack, icaPacket, packet.SourceChannel, packet.Sequence)
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

	if err := k.handleUnsuccessfulAck(ctx, icaPacket, packet.SourceChannel, packet.Sequence); err != nil {
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
	channel string,
	sequence uint64,
) error {
	messages, err := icatypes.DeserializeCosmosTx(k.cdc, icaPacket.GetData())
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "cannot deserialize ica packet data: %v", err)
	}

	for _, msg := range messages {
		switch sdk.MsgTypeURL(msg) {
		case sdk.MsgTypeURL(&stakingtypes.MsgDelegate{}):
			// revert all the deposits for that sequence back to the previous state
			k.RevertDepositsState(ctx, k.GetDepositsWithSequenceID(ctx, k.GetTransactionSequenceID(channel, sequence)))
		case sdk.MsgTypeURL(&stakingtypes.MsgUndelegate{}):
			// mark all the unbondings for the previous epoch as failed
			k.FailAllUnbondingsForSequenceID(ctx, k.GetTransactionSequenceID(channel, sequence))
			// delete all validator unbondings so they can be picked up again
			k.DeleteValidatorUnbondingsForSequenceID(ctx, k.GetTransactionSequenceID(channel, sequence))
		case sdk.MsgTypeURL(&ibctransfertypes.MsgTransfer{}):
			unbondings := k.FilterUnbondings(
				ctx,
				func(u types.Unbonding) bool {
					return u.IbcSequenceId == k.GetTransactionSequenceID(channel, sequence)
				},
			)
			// revert unbonding state so it can be picked up again
			// this won't conflict with failed rewards transfers since the transaction sequence id won't match
			k.RevertUnbondingsState(ctx, unbondings)

			validatorUnbondings := k.FilterValidatorUnbondings(
				ctx,
				func(u types.ValidatorUnbonding) bool {
					return u.IbcSequenceId == k.GetTransactionSequenceID(channel, sequence)
				},
			)

			// empty the ibc sequence id, so they will be picked up again while processing mature delegations
			for _, validatorUnbonding := range validatorUnbondings {
				validatorUnbonding.IbcSequenceId = ""
				k.SetValidatorUnbonding(ctx, validatorUnbonding)
			}
		case sdk.MsgTypeURL(&stakingtypes.MsgRedeemTokensForShares{}):
			deposits := k.FilterLSMDeposits(
				ctx,
				func(d types.LSMDeposit) bool {
					return d.IbcSequenceId == k.GetTransactionSequenceID(channel, sequence)
				},
			)

			// revert the state of the deposit, so it will be retried
			k.RevertLSMDepositsState(ctx, deposits)
		}
	}

	return nil
}

func (k *Keeper) handleSuccessfulAck(
	ctx sdk.Context,
	ack channeltypes.Acknowledgement,
	icaPacket icatypes.InterchainAccountPacketData,
	channel string,
	sequence uint64,
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
		case sdk.MsgTypeURL(&stakingtypes.MsgDelegate{}):
			if err = k.HandleDelegateResponse(ctx, msg, channel, sequence); err != nil {
				return err
			}
		case sdk.MsgTypeURL(&stakingtypes.MsgUndelegate{}):
			var data []byte
			if len(txMsgData.Data) == 0 {
				data = txMsgData.GetMsgResponses()[i].Value
			} else {
				data = txMsgData.Data[i].Data
			}

			var msgResponse stakingtypes.MsgUndelegateResponse
			if err := k.cdc.Unmarshal(data, &msgResponse); err != nil {
				return errorsmod.Wrapf(
					sdkerrors.ErrJSONUnmarshal, "cannot unmarshal undelegate response message: %s",
					err.Error(),
				)
			}

			if err = k.HandleUndelegateResponse(ctx, msg, msgResponse, channel, sequence); err != nil {
				return err
			}
		case sdk.MsgTypeURL(&ibctransfertypes.MsgTransfer{}):
			var data []byte
			if len(txMsgData.Data) == 0 {
				data = txMsgData.GetMsgResponses()[i].Value
			} else {
				data = txMsgData.Data[i].Data
			}

			var msgResponse ibctransfertypes.MsgTransferResponse
			if err := k.cdc.Unmarshal(data, &msgResponse); err != nil {
				return errorsmod.Wrapf(
					sdkerrors.ErrJSONUnmarshal, "cannot unmarshal undelegate response message: %s",
					err.Error(),
				)
			}

			if err = k.HandleMsgTransfer(ctx, msg, msgResponse, channel, sequence); err != nil {
				return err
			}
		case sdk.MsgTypeURL(&stakingtypes.MsgRedeemTokensForShares{}):
			var data []byte
			if len(txMsgData.Data) == 0 {
				data = txMsgData.GetMsgResponses()[i].Value
			} else {
				data = txMsgData.Data[i].Data
			}

			var msgResponse stakingtypes.MsgRedeemTokensForSharesResponse
			if err := k.cdc.Unmarshal(data, &msgResponse); err != nil {
				return errorsmod.Wrapf(
					sdkerrors.ErrJSONUnmarshal, "cannot unmarshal redeem tokens response message: %s",
					err.Error(),
				)
			}

			if err = k.HandleMsgRedeemTokensForShares(ctx, msg, msgResponse, channel, sequence); err != nil {
				return err
			}
		}
	}

	k.Logger(ctx).Info(
		"ICA transaction ACK success.",
		"sequence",
		sequence,
		"channel",
		channel,
		"messages",
		messages,
	)

	return nil
}
