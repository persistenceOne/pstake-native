package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v6/modules/core/exported"
	"github.com/persistenceOne/persistence-sdk/v2/utils"
	epochstypes "github.com/persistenceOne/persistence-sdk/v2/x/epochs/types"
	ibchookertypes "github.com/persistenceOne/persistence-sdk/v2/x/ibchooker/types"

	liquidstakeibctypes "github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

type EpochsHooks struct {
	k Keeper
}

var _ epochstypes.EpochHooks = EpochsHooks{}

func (h EpochsHooks) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	return h.k.BeforeEpochStart(ctx, epochIdentifier, epochNumber)
}

func (h EpochsHooks) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	return h.k.AfterEpochEnd(ctx, epochIdentifier, epochNumber)
}

type IBCTransferHooks struct {
	k Keeper
}

var _ ibchookertypes.IBCHandshakeHooks = IBCTransferHooks{}

func (k *Keeper) NewIBCTransferHooks() IBCTransferHooks {
	return IBCTransferHooks{*k}
}

func (i IBCTransferHooks) OnRecvPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
	transferAck ibcexported.Acknowledgement,
) error {
	return i.k.OnRecvIBCTransferPacket(ctx, packet, relayer, transferAck)
}

func (i IBCTransferHooks) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
	transferAckErr error,
) error {
	return i.k.OnAcknowledgementIBCTransferPacket(ctx, packet, acknowledgement, relayer, transferAckErr)
}

func (i IBCTransferHooks) OnTimeoutPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
	transferTimeoutErr error,
) error {
	return i.k.OnTimeoutIBCTransferPacket(ctx, packet, relayer, transferTimeoutErr)
}

// Module hooks

func (k *Keeper) NewEpochHooks() EpochsHooks {
	return EpochsHooks{*k}
}

func (k *Keeper) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	// create a batch of user deposits for the new deposit epoch
	if epochIdentifier == liquidstakeibctypes.DelegationEpoch {
		k.CreateDeposits(ctx, epochNumber)
	}

	return nil
}

func (k *Keeper) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	if epochIdentifier == liquidstakeibctypes.DelegationEpoch {
		workflow := func(ctx sdk.Context) error {
			return k.DepositWorkflow(ctx, epochNumber)
		}
		err := utils.ApplyFuncIfNoError(ctx, workflow)
		if err != nil {
			k.Logger(ctx).Error(
				"failed delegation workflow",
				"epoch_identifier",
				epochIdentifier,
				"epoch_number",
				epochNumber,
				"error",
				err,
			)
		}
	}

	if epochIdentifier == liquidstakeibctypes.UndelegationEpoch {
		workflow := func(ctx sdk.Context) error {
			return k.UndelegationWorkflow(ctx, epochNumber)
		}
		err := utils.ApplyFuncIfNoError(ctx, workflow)
		if err != nil { //nolint:staticcheck
			// handle this case
		}
	}

	return nil
}

// IBC transfer hooks

func (k *Keeper) OnRecvIBCTransferPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
	transferAck ibcexported.Acknowledgement,
) error {
	return nil
}

func (k *Keeper) OnAcknowledgementIBCTransferPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
	transferAckErr error,
) error {
	if transferAckErr != nil {
		return transferAckErr
	}

	// validate the ack
	var ack channeltypes.Acknowledgement
	if err := ibctransfertypes.ModuleCdc.UnmarshalJSON(acknowledgement, &ack); err != nil {
		return err
	}
	if !ack.Success() {
		return channeltypes.ErrInvalidAcknowledgement
	}

	var data ibctransfertypes.FungibleTokenPacketData
	if err := ibctransfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return err
	}

	transferAmount, ok := sdk.NewIntFromString(data.Amount)
	if !ok {
		return fmt.Errorf("could not parse ibc transfer amount %s", data.Amount)
	}

	// if the sender is the deposit module account, mark the corresponding deposits as received and send an
	// ICQ query to get the new host delegator account balance
	if data.GetSender() == authtypes.NewModuleAddress(liquidstakeibctypes.DepositModuleAccount).String() {
		deposits := k.GetDepositsWithSequenceID(ctx, k.GetTransactionSequenceID(packet.SourceChannel, packet.Sequence))
		for _, deposit := range deposits {
			// update the deposit state
			deposit.IbcSequenceId = ""
			deposit.State = liquidstakeibctypes.Deposit_DEPOSIT_RECEIVED
			k.SetDeposit(ctx, deposit)

			hc, found := k.GetHostChain(ctx, deposit.ChainId)
			if !found {
				return fmt.Errorf("host chain with id %s is not registered", deposit.ChainId)
			}

			hc.DelegationAccount.Balance = hc.DelegationAccount.Balance.Add(
				sdk.Coin{
					Denom:  hc.DelegationAccount.Balance.Denom,
					Amount: transferAmount,
				},
			)

			hc.CValue = k.GetHostChainCValue(ctx, hc)
			k.SetHostChain(ctx, hc)
		}
	}

	return nil
}

func (k *Keeper) OnTimeoutIBCTransferPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
	transferTimeoutErr error,
) error {
	if transferTimeoutErr != nil {
		return transferTimeoutErr
	}

	var data ibctransfertypes.FungibleTokenPacketData
	if err := ibctransfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return err
	}

	// if the transfer doesn't belong to any of the registered host chains, return
	ibcDenom := ibctransfertypes.ParseDenomTrace(data.GetDenom()).IBCDenom()
	hc, found := k.GetHostChainFromIbcDenom(ctx, ibcDenom)
	if !found {
		return nil
	}

	// if the transfer is not from deposit module account -> delegation host account, return
	if data.GetSender() != authtypes.NewModuleAddress(liquidstakeibctypes.DepositModuleAccount).String() ||
		data.GetReceiver() != hc.DelegationAccount.Address ||
		data.GetDenom() != ibctransfertypes.GetPrefixedDenom(hc.PortId, hc.ChannelId, hc.HostDenom) {
		return nil
	}

	// revert all the deposits for that sequence to its previous state
	k.RevertDepositsState(
		ctx,
		k.GetDepositsWithSequenceID(ctx, k.GetTransactionSequenceID(packet.SourceChannel, packet.Sequence)),
	)

	return nil
}

// Workflows

func (k *Keeper) DepositWorkflow(ctx sdk.Context, epoch int64) error {
	deposits := k.GetPendingDepositsBeforeEpoch(ctx, epoch)
	for _, deposit := range deposits {
		hc, found := k.GetHostChain(ctx, deposit.ChainId)
		if !found {
			return fmt.Errorf("host chain with id %s is not registered", deposit.ChainId)
		}

		// check if the deposit amount is larger than 0
		if deposit.Amount.Amount.LTE(sdk.NewInt(0)) {
			// delete empty deposits to save on storage
			if deposit.Epoch.Int64() < epoch {
				k.DeleteDeposit(ctx, deposit)
			}

			continue
		}

		clientState, err := k.GetClientState(ctx, hc.ConnectionId)
		if err != nil {
			return fmt.Errorf("client state not found for connection \"%s\": \"%s\"", hc.ConnectionId, err.Error())
		}

		timeoutHeight := clienttypes.NewHeight(
			clientState.GetLatestHeight().GetRevisionNumber(),
			clientState.GetLatestHeight().GetRevisionHeight()+liquidstakeibctypes.IBCTimeoutHeightIncrement,
		)

		msg := ibctransfertypes.NewMsgTransfer(
			ibctransfertypes.TypeMsgTransfer,
			hc.ChannelId,
			deposit.Amount,
			authtypes.NewModuleAddress(liquidstakeibctypes.DepositModuleAccount).String(),
			hc.DelegationAccount.Address,
			timeoutHeight,
			0,
			"",
		)

		handler := k.msgRouter.Handler(msg)
		res, err := handler(ctx, msg)
		if err != nil {
			k.Logger(ctx).Error(fmt.Sprintf("could not send transfer msg via MsgServiceRouter, error: %s", err))
			return err
		}
		ctx.EventManager().EmitEvents(res.GetEvents())

		var msgTransferResponse ibctransfertypes.MsgTransferResponse
		if err = k.cdc.Unmarshal(res.MsgResponses[0].Value, &msgTransferResponse); err != nil {
			return errorsmod.Wrapf(
				sdkerrors.ErrJSONUnmarshal,
				"cannot unmarshal ibc transfer tx response message: %v",
				err,
			)
		}

		deposit.State = liquidstakeibctypes.Deposit_DEPOSIT_SENT
		deposit.IbcSequenceId = k.GetTransactionSequenceID(hc.ChannelId, msgTransferResponse.Sequence)
		k.SetDeposit(ctx, deposit)
	}

	return nil
}

func (k *Keeper) UndelegationWorkflow(ctx sdk.Context, epoch int64) error {
	return nil
}
