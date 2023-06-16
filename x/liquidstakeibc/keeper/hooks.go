package keeper

import (
	"fmt"
	"time"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/gogoproto/proto"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
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

	// update the c value for each registered host chain
	if epochIdentifier == liquidstakeibctypes.CValueEpoch {
		k.UpdateCValues(ctx)
	}

	return nil
}

func (k *Keeper) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	if epochIdentifier == liquidstakeibctypes.DelegationEpoch {
		k.DepositWorkflow(ctx, epochNumber)
	}

	if epochIdentifier == liquidstakeibctypes.UndelegationEpoch {
		// attempt to fully undelegate any validators that have been more than
		//UnbondingStateEpochLimit epochs in UNBONDING state
		k.ValidatorUndelegationWorkflow(ctx, epochNumber)

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
	k.Logger(ctx).Info(
		"Received incoming IBC transfer.",
		"sequence",
		packet.Sequence,
		"port",
		packet.DestinationPort,
		"channel",
		packet.DestinationChannel,
	)

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

	// the transfer is part of the undelegation process
	if data.GetSender() == hc.DelegationAccount.Address &&
		data.GetReceiver() == k.GetUndelegationModuleAccount(ctx).GetAddress().String() &&
		data.Memo == "" {
		k.Logger(ctx).Info(
			"Received unbonding IBC transfer.",
			"host chain",
			hc.ChainId,
			"sequence",
			packet.Sequence,
			"port",
			packet.DestinationPort,
			"channel",
			packet.DestinationChannel,
		)

		// get all the unbondings for that ibc sequence id
		unbondings := k.FilterUnbondings(
			ctx,
			func(u liquidstakeibctypes.Unbonding) bool {
				return u.ChainId == hc.ChainId && u.State == liquidstakeibctypes.Unbonding_UNBONDING_MATURED
			},
		)

		// update the unbonding states
		for _, unbonding := range unbondings {
			unbonding.IbcSequenceId = ""
			unbonding.State = liquidstakeibctypes.Unbonding_UNBONDING_CLAIMABLE
			k.SetUnbonding(ctx, unbonding)
		}
	}

	// the transfer is part of a total validator unbonding
	if data.GetSender() == hc.DelegationAccount.Address &&
		data.GetReceiver() == k.GetDepositModuleAccount(ctx).GetAddress().String() &&
		data.Memo == "" {
		k.Logger(ctx).Info(
			"Received total validator unbonding IBC transfer.",
			"host chain",
			hc.ChainId,
			"sequence",
			packet.Sequence,
			"port",
			packet.DestinationPort,
			"channel",
			packet.DestinationChannel,
		)

		// add the unbonded amount to the deposit record for that chain/epoch
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

		transferAmount, ok := sdk.NewIntFromString(data.Amount)
		if !ok {
			return errorsmod.Wrapf(
				liquidstakeibctypes.ErrParsingAmount,
				"could not parse transfer amount %s",
				data.Amount,
			)
		}

		deposit.Amount.Amount = deposit.Amount.Amount.Add(transferAmount)
		k.SetDeposit(ctx, deposit)
	}

	// the transfer is part of the autocompounding process
	if data.GetSender() == hc.RewardsAccount.Address &&
		data.GetReceiver() == k.GetDepositModuleAccount(ctx).GetAddress().String() &&
		data.Memo == "" {
		k.Logger(ctx).Info(
			"Received autocompounding IBC transfer.",
			"host chain",
			hc.ChainId,
			"sequence",
			packet.Sequence,
			"port",
			packet.DestinationPort,
			"channel",
			packet.DestinationChannel,
		)

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

	// if the sender is the deposit module account, mark the corresponding deposits as received and update the balance
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

			k.SetHostChain(ctx, hc)

			k.Logger(ctx).Info(
				"Got delegation deposit received ACK.",
				"host chain",
				hc.ChainId,
				"sequence",
				packet.Sequence,
				"port",
				packet.SourceChannel,
				"channel",
				packet.SourceChannel,
			)
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

	k.Logger(ctx).Info(
		"Deposit transfer timed out.",
		"host chain",
		hc.ChainId,
		"sequence",
		packet.Sequence,
		"port",
		packet.SourceChannel,
		"channel",
		packet.SourceChannel,
	)

	return nil
}

// Workflows

func (k *Keeper) DepositWorkflow(ctx sdk.Context, epoch int64) {
	k.Logger(ctx).Info("Running deposit workflow.", "epoch", epoch)

	deposits := k.GetPendingDepositsBeforeEpoch(ctx, epoch)
	for _, deposit := range deposits {
		hc, found := k.GetHostChain(ctx, deposit.ChainId)
		if !found {
			// we can't error out here as all the deposits need to be executed
			continue
		}

		// don't do anything if the chain is not active
		if !hc.Active {
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
			ibctransfertypes.PortID,
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
	k.Logger(ctx).Info("Running undelegation workflow.", "epoch", epoch)

	for _, hc := range k.GetAllHostChains(ctx) {
		// don't do anything if the chain is not active
		if !hc.Active {
			continue
		}

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
			hc.DelegationAccount.Owner,
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

func (k *Keeper) ValidatorUndelegationWorkflow(ctx sdk.Context, epoch int64) {
	k.Logger(ctx).Info("Running validator undelegation workflow.", "epoch", epoch)

	for _, hc := range k.GetAllHostChains(ctx) {
		// don't do anything if the chain is not active
		if !hc.Active {
			continue
		}

		// not an unbonding epoch for the host chain, continue
		if !liquidstakeibctypes.IsUnbondingEpoch(hc.UnbondingFactor, epoch) {
			continue
		}

		for _, validator := range hc.Validators {
			// check if there are validators that need to be unbonded
			if validator.UnbondingEpoch > 0 &&
				validator.UnbondingEpoch+liquidstakeibctypes.UnbondingStateEpochLimit <= epoch {

				// unbond all delegated tokens from the validator
				validatorUnbonding := &liquidstakeibctypes.ValidatorUnbonding{
					ChainId:          hc.ChainId,
					EpochNumber:      epoch,
					MatureTime:       time.Time{},
					ValidatorAddress: validator.OperatorAddress,
					Amount:           sdk.NewCoin(hc.HostDenom, validator.DelegatedAmount),
				}

				// create the MsgUndelegate
				message := &stakingtypes.MsgUndelegate{
					DelegatorAddress: hc.DelegationAccount.Address,
					ValidatorAddress: validatorUnbonding.ValidatorAddress,
					Amount:           validatorUnbonding.Amount,
				}

				// execute the ICA transaction
				sequenceID, err := k.GenerateAndExecuteICATx(
					ctx,
					hc.ConnectionId,
					hc.DelegationAccount.Owner,
					[]proto.Message{message},
				)
				if err != nil {
					k.Logger(ctx).Error(
						"could not send ICA undelegate txs",
						"host_chain",
						hc.ChainId,
					)
					return
				}

				// update the unbonding sequence id
				validatorUnbonding.IbcSequenceId = sequenceID
				k.SetValidatorUnbonding(ctx, validatorUnbonding)

				// redistribute the unbonding validator weight among all the other validators with weight
				k.RedistributeValidatorWeight(ctx, hc, validator)

				k.Logger(ctx).Info(
					"Started total validator unbonding.",
					"host_chain",
					hc.ChainId,
					"validator",
					validatorUnbonding.ValidatorAddress,
					"amount",
					validatorUnbonding.Amount,
					"epoch",
					epoch,
				)
			}
		}
	}
}

func (k *Keeper) RewardsWorkflow(ctx sdk.Context, epoch int64) {
	k.Logger(ctx).Info("Running rewards workflow.", "epoch", epoch)

	for _, hc := range k.GetAllHostChains(ctx) {
		// don't do anything if the chain is not active
		if !hc.Active {
			continue
		}

		// generate the messages
		messages := make([]proto.Message, 0)
		for _, validator := range hc.Validators {
			if validator.DelegatedAmount.GT(sdk.ZeroInt()) {
				message := &distributiontypes.MsgWithdrawDelegatorReward{
					DelegatorAddress: hc.DelegationAccount.Address,
					ValidatorAddress: validator.OperatorAddress,
				}
				messages = append(messages, message)
			}
		}

		if len(messages) > 0 {
			// execute the ICA transactions
			_, err := k.GenerateAndExecuteICATx(
				ctx,
				hc.ConnectionId,
				hc.DelegationAccount.Owner,
				messages,
			)
			if err != nil {
				k.Logger(ctx).Error(
					"Could not send ICA withdraw delegator reward txs",
					"host_chain",
					hc.ChainId,
				)
				continue
			}
		}

		if hc.RewardsAccount != nil &&
			hc.RewardsAccount.ChannelState == liquidstakeibctypes.ICAAccount_ICA_CHANNEL_CREATED {
			if err := k.QueryRewardsHostChainAccountBalance(ctx, hc); err != nil {
				k.Logger(ctx).Error(
					"Could not send rewards account balance ICQ",
					"host_chain",
					hc.ChainId,
				)
				continue
			}
		}
	}
}
