package keeper

import (
	"fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	icatypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"

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

	// get the ica account type from the ownership string
	_, icaAccountType, found := strings.Cut(portOwner, ".")
	if !found {
		return fmt.Errorf("unable to parse port owner %s", portOwner)
	}

	// create the ica account
	icaAccount := &types.ICAAccount{
		Address:      address,
		Balance:      sdk.Coin{Amount: sdk.ZeroInt(), Denom: hc.HostDenom},
		Owner:        portOwner,
		ChannelState: types.ICAAccount_ICA_CHANNEL_CREATED,
	}

	switch icaAccountType {
	case types.DelegateICAType:
		hc.DelegationAccount = icaAccount
	case types.RewardsICAType:
		hc.RewardsAccount = icaAccount
	}

	if hc.DelegationAccount != nil && hc.RewardsAccount != nil {
		err := k.SetWithdrawAddress(ctx, hc)
		if err != nil {
			k.Logger(ctx).Error("Could not set withdraw address.", "chain_id", hc.ChainId)
		}
	}

	// save the changes of the host chain
	k.SetHostChain(ctx, hc)

	// revert the state for all the deposits that were being delegated on that host chain
	k.RevertDepositsState(ctx, k.GetDelegatingDepositsForChain(ctx, hc.ChainId))

	// send an ICQ query to get the delegator account balance
	if err := k.QueryHostChainAccountBalance(ctx, hc, hc.DelegationAccount.Address); err != nil {
		return fmt.Errorf(
			"error querying host chain %s for delegation account balances: %v",
			hc.ChainId,
			err,
		)
	}

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
		err := k.handleUnsuccessfulAck(ctx, packet.SourceChannel, packet.Sequence)
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
		err := k.handleSuccessfulAck(ctx, icaPacket, packet.SourceChannel, packet.Sequence)
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

	messages, err := icatypes.DeserializeCosmosTx(k.cdc, icaPacket.GetData())
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "cannot deserialize ica packet data: %v", err)
	}

	for _, msg := range messages {
		switch sdk.MsgTypeURL(msg) { //nolint:gocritic
		case sdk.MsgTypeURL(&stakingtypes.MsgDelegate{}):
			// revert all the deposits for that sequence to its previous state
			k.RevertDepositsState(
				ctx,
				k.GetDepositsWithSequenceID(ctx, k.GetTransactionSequenceID(packet.SourceChannel, packet.Sequence)),
			)
		}
	}

	k.Logger(ctx).Info(
		fmt.Sprintf(
			"ICA packet timed out with seq: %v, channel: %s, port: %s, msgs: %s",
			packet.Sequence,
			packet.SourceChannel,
			packet.SourcePort,
			messages,
		),
	)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTimeout,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		),
	)
	return nil
}

func (k *Keeper) handleUnsuccessfulAck(
	ctx sdk.Context,
	channel string,
	sequence uint64,
) error {
	// revert all the deposits for that sequence back to the previous state
	k.RevertDepositsState(ctx, k.GetDepositsWithSequenceID(ctx, k.GetTransactionSequenceID(channel, sequence)))

	return nil
}

func (k *Keeper) handleSuccessfulAck(
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
		switch sdk.MsgTypeURL(msg) { //nolint:gocritic
		case sdk.MsgTypeURL(&stakingtypes.MsgDelegate{}):
			if err = k.HandleDelegateResponse(ctx, msg, channel, sequence); err != nil {
				return err
			}
		}
	}

	return nil
}
