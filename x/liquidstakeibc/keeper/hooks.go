package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v6/modules/core/exported"
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
		k.DepositWorkflow(ctx, epochNumber)
	}

	if epochIdentifier == liquidstakeibctypes.UndelegationEpoch {
		k.UndelegationWorkflow(ctx, epochNumber)
	}

	if epochIdentifier == liquidstakeibctypes.RewardsEpochIdentifier {
		k.RewardsWorkflow(ctx, epochNumber)
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
	if !transferAck.Success() {
		return nil
	}

	var data ibctransfertypes.FungibleTokenPacketData
	if err := ibctransfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return err
	}

	// if the transfer isn't from any of the registered host chains, return
	denom := data.GetDenom()
	hc, found := k.GetHostChainFromHostDenom(ctx, denom)
	if !found {
		return nil
	}

	// the transfer goes delegationAddress -> undelegationAccount, update corresponding unbondings
	if data.GetSender() == hc.DelegationAccount.Address &&
		data.GetReceiver() == k.GetUndelegationModuleAccount(ctx).GetAddress().String() {
		// get all the unbondings for that ibc sequence id
		unbondings := k.FilterUnbondings(
			ctx,
			func(u liquidstakeibctypes.Unbonding) bool {
				return u.UnbondAmount.Denom == hc.HostDenom && u.State == liquidstakeibctypes.Unbonding_UNBONDING_MATURED
			},
		)

		// update the unbonding states
		for _, unbonding := range unbondings {
			unbonding.IbcSequenceId = ""
			unbonding.State = liquidstakeibctypes.Unbonding_UNBONDING_CLAIMABLE
			k.SetUnbonding(ctx, unbonding)
		}
	}

	if data.GetSender() == hc.RewardsAccount.Address &&
		data.GetReceiver() == k.GetDepositModuleAccount(ctx).GetAddress().String() {
		// parse the transfer amount
		transferAmount, ok := sdk.NewIntFromString(data.Amount)
		if !ok {
			return errorsmod.Wrapf(
				liquidstakeibctypes.ErrParsingAmount,
				"could not parse transfer amount %s",
				data.Amount,
			)
		}

		// calculate protocol fee
		feeAmount := hc.Params.RestakeFee.MulInt(transferAmount)
		fee, _ := sdk.NewDecCoinFromDec(hc.IBCDenom(), feeAmount).TruncateDecimal()

		// send the protocol fee
		err := k.SendProtocolFee(
			ctx,
			sdk.NewCoins(fee),
			liquidstakeibctypes.DepositModuleAccount,
			k.GetParams(ctx).FeeAddress,
		)
		if err != nil {
			return errorsmod.Wrapf(
				liquidstakeibctypes.ErrFailedDeposit,
				"failed to send restake fee to module fee address %s: %s",
				k.GetParams(ctx).FeeAddress,
				err.Error(),
			)
		}

		// add the deposit amount to the deposit record for that chain/epoch
		currentEpoch := k.GetEpochNumber(ctx, liquidstakeibctypes.DelegationEpoch)
		deposit, found := k.GetDepositForChainAndEpoch(ctx, hc.ChainId, currentEpoch)
		if !found {
			return errorsmod.Wrapf(
				liquidstakeibctypes.ErrDepositNotFound,
				"deposit not found for chain %s and epoch %v",
				hc.ChainId,
				currentEpoch,
			)
		}

		// update the deposit
		deposit.Amount.Amount = deposit.Amount.Amount.Add(transferAmount.Sub(feeAmount.TruncateInt()))
		k.SetDeposit(ctx, deposit)
	}

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

	// revert the state of the deposits that timed out
	k.RevertDepositsState(
		ctx,
		k.GetDepositsWithSequenceID(ctx, k.GetTransactionSequenceID(packet.SourceChannel, packet.Sequence)),
	)

	return nil
}

// Workflows

func (k *Keeper) DepositWorkflow(ctx sdk.Context, epoch int64) {
	deposits := k.GetPendingDepositsBeforeEpoch(ctx, epoch)
	for _, deposit := range deposits {
		hc, found := k.GetHostChain(ctx, deposit.ChainId)
		if !found {
			// we can't error out here as all the deposits need to be executed
			continue
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
			// we can't error out here as all the deposits need to be executed
			continue
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
			// we can't error out here as all the deposits need to be executed
			continue
		}
		ctx.EventManager().EmitEvents(res.GetEvents())

		var msgTransferResponse ibctransfertypes.MsgTransferResponse
		if err = k.cdc.Unmarshal(res.MsgResponses[0].Value, &msgTransferResponse); err != nil {
			// we can't error out here as all the deposits need to be executed
			continue
		}

		deposit.State = liquidstakeibctypes.Deposit_DEPOSIT_SENT
		deposit.IbcSequenceId = k.GetTransactionSequenceID(hc.ChannelId, msgTransferResponse.Sequence)
		k.SetDeposit(ctx, deposit)
	}
}

func (k *Keeper) UndelegationWorkflow(ctx sdk.Context, epoch int64) {
	for _, hc := range k.GetAllHostChains(ctx) {
		// not an unbonding epoch for the host chain, continue
		if !liquidstakeibctypes.IsUnbondingEpoch(hc.UnbondingFactor, epoch) {
			continue
		}

		// retrieve the unbonding for the current epoch
		unbonding, found := k.GetUnbonding(
			ctx,
			hc.ChainId,
			liquidstakeibctypes.CurrentUnbondingEpoch(hc.UnbondingFactor, epoch),
		)
		if !found {
			// nothing to unbond for this epoch
			continue
		}

		// check if there is anything to unbond
		if !unbonding.UnbondAmount.Amount.GT(sdk.ZeroInt()) {
			k.Logger(ctx).Info(
				"No tokens to unbond.",
				"host_chain",
				hc.ChainId,
				"epoch",
				epoch,
			)
			continue
		}

		// generate the undelegation messages based on the total unbonding amount for the epoch
		messages, err := k.GenerateUndelegateMessages(hc, unbonding.UnbondAmount.Amount)
		if err != nil {
			k.Logger(ctx).Error(
				"could not generate undelegate messages",
				"host_chain",
				hc.ChainId,
			)
			return
		}

		// execute the ICA transactions
		sequenceID, err := k.GenerateAndExecuteICATx(
			ctx,
			hc.ConnectionId,
			k.DelegateAccountPortOwner(hc.ChainId),
			messages,
		)
		if err != nil {
			k.Logger(ctx).Error(
				"could not send ICA undelegate txs",
				"host_chain",
				hc.ChainId,
			)
			return
		}

		// update the unbonding ibc sequence id and state
		unbonding.IbcSequenceId = sequenceID
		unbonding.State = liquidstakeibctypes.Unbonding_UNBONDING_INITIATED
		k.SetUnbonding(ctx, unbonding)
	}
}

func (k *Keeper) RewardsWorkflow(ctx sdk.Context, epoch int64) {
	for _, hc := range k.GetAllHostChains(ctx) {
		if hc.RewardsAccount != nil &&
			hc.RewardsAccount.ChannelState == liquidstakeibctypes.ICAAccount_ICA_CHANNEL_CREATED {
			if err := k.QueryHostChainAccountBalance(ctx, hc, hc.RewardsAccount.Address); err != nil {
				k.Logger(ctx).Info(
					"Could not send rewards account balance ICQ.",
					"host_chain",
					hc.ChainId,
					"epoch",
					epoch,
				)
			}
		}
	}
}
