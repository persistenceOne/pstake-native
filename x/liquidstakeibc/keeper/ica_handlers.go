package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

func (k *Keeper) HandleDelegateResponse(ctx sdk.Context, msg sdk.Msg, channel string, sequence uint64) error {
	parsedMsg, ok := msg.(*stakingtypes.MsgDelegate)
	if !ok {
		return errorsmod.Wrapf(
			sdkerrors.ErrInvalidType,
			"unable to cast msg of type %s to MsgDelegate",
			sdk.MsgTypeURL(msg),
		)
	}

	// remove delegated deposits for this sequence (if any)
	deposits := k.GetDepositsWithSequenceID(ctx, k.GetTransactionSequenceID(channel, sequence))
	for _, deposit := range deposits {
		k.DeleteDeposit(ctx, deposit)
	}

	// get the host chain of the delegation using its delegator address
	hc, found := k.GetHostChainFromDelegatorAddress(ctx, parsedMsg.DelegatorAddress)
	if !found {
		return errorsmod.Wrapf(
			types.ErrInvalidHostChain,
			"host chain with delegator address %s not registered, or account not associated",
			parsedMsg.DelegatorAddress,
		)
	}

	// update delegation account balance
	hc.DelegationAccount.Balance = hc.DelegationAccount.Balance.Sub(parsedMsg.Amount)

	// get the validator that the delegation was performed to
	validator, found := hc.GetValidator(parsedMsg.ValidatorAddress)
	if !found {
		return errorsmod.Wrapf(
			types.ErrValidatorNotFound,
			"validator with operator address %s not found",
			parsedMsg.ValidatorAddress,
		)
	}

	// update the validator delegated amount
	validator.DelegatedAmount = validator.DelegatedAmount.Add(parsedMsg.Amount.Amount)
	k.SetHostChainValidator(ctx, hc, validator)

	// update host the host chain c value
	hc.CValue = k.GetHostChainCValue(ctx, hc)
	k.SetHostChain(ctx, hc)

	k.Logger(ctx).Info(
		"Received delegation acknowledgement",
		"delegator",
		parsedMsg.DelegatorAddress,
		"validator",
		parsedMsg.ValidatorAddress,
		"amount",
		parsedMsg.Amount.String(),
	)

	return nil
}

func (k *Keeper) HandleUndelegateResponse(
	ctx sdk.Context,
	msg sdk.Msg,
	resp stakingtypes.MsgUndelegateResponse,
	channel string,
	sequence uint64,
) error {
	parsedMsg, ok := msg.(*stakingtypes.MsgUndelegate)
	if !ok {
		return errorsmod.Wrapf(
			sdkerrors.ErrInvalidType,
			"unable to cast msg of type %s to MsgUndelegate",
			sdk.MsgTypeURL(msg),
		)
	}

	// get the host chain of the delegation using its delegator address
	hc, found := k.GetHostChainFromDelegatorAddress(ctx, parsedMsg.DelegatorAddress)
	if !found {
		return errorsmod.Wrapf(
			types.ErrInvalidHostChain,
			"host chain with delegator address %s not registered, or account not associated",
			parsedMsg.DelegatorAddress,
		)
	}

	// get the validator that the delegation was performed to
	validator, found := hc.GetValidator(parsedMsg.ValidatorAddress)
	if !found {
		return errorsmod.Wrapf(
			types.ErrValidatorNotFound,
			"validator with operator address %s not found",
			parsedMsg.ValidatorAddress,
		)
	}

	// update the validator delegated amount
	validator.DelegatedAmount = validator.DelegatedAmount.Sub(parsedMsg.Amount.Amount)
	k.SetHostChainValidator(ctx, hc, validator)

	// update host the host chain c value
	hc.CValue = k.GetHostChainCValue(ctx, hc)
	k.SetHostChain(ctx, hc)

	// update the state of all the unbondings associated with the undelegation
	unbondings := k.FilterUnbondings(
		ctx,
		func(u types.Unbonding) bool { return u.IbcSequenceId == k.GetTransactionSequenceID(channel, sequence) },
	)

	for _, unbonding := range unbondings {
		// burn the undelegated stk tokens
		err := k.bankKeeper.BurnCoins(
			ctx,
			types.UndelegationModuleAccount,
			sdk.NewCoins(unbonding.BurnAmount),
		)
		if err != nil {
			return err
		}

		// update the mature time and the state for the undelegation
		unbonding.IbcSequenceId = ""
		unbonding.MatureTime = resp.CompletionTime
		unbonding.State = types.Unbonding_UNBONDING_MATURING
		k.SetUnbonding(ctx, unbonding)

		k.Logger(ctx).Info(
			"Received unbonding acknowledgement",
			"delegator",
			parsedMsg.DelegatorAddress,
			"validator",
			parsedMsg.ValidatorAddress,
			"amount",
			parsedMsg.Amount.String(),
		)
	}

	// update the state of all the validator unbondings associated with the undelegation
	validatorUnbondings := k.FilterValidatorUnbondings(
		ctx,
		func(u types.ValidatorUnbonding) bool {
			return u.IbcSequenceId == k.GetTransactionSequenceID(channel, sequence)
		},
	)

	for _, validatorUnbonding := range validatorUnbondings {
		// update the mature time and the state for the validator undelegation
		validatorUnbonding.IbcSequenceId = ""
		validatorUnbonding.MatureTime = resp.CompletionTime
		k.SetValidatorUnbonding(ctx, validatorUnbonding)

		k.Logger(ctx).Info(
			"Received validator unbonding acknowledgement",
			"delegator",
			parsedMsg.DelegatorAddress,
			"validator",
			parsedMsg.ValidatorAddress,
			"amount",
			parsedMsg.Amount.String(),
		)
	}

	return nil
}

func (k *Keeper) HandleMsgTransfer(
	ctx sdk.Context,
	msg sdk.Msg,
	resp ibctransfertypes.MsgTransferResponse,
	channel string,
	sequence uint64,
) error {
	parsedMsg, ok := msg.(*ibctransfertypes.MsgTransfer)
	if !ok {
		return errorsmod.Wrapf(
			sdkerrors.ErrInvalidType,
			"unable to cast msg of type %s to MsgUndelegate",
			sdk.MsgTypeURL(msg),
		)
	}

	// get the host chain of the transfer using its host denom
	hc, found := k.GetHostChainFromHostDenom(ctx, parsedMsg.Token.Denom)
	if !found {
		return errorsmod.Wrapf(
			types.ErrInvalidHostChain,
			"host chain with host denom %s not registered",
			parsedMsg.Token.Denom,
		)
	}

	// the transfer is part of the undelegation process
	if parsedMsg.Sender == hc.DelegationAccount.Address &&
		parsedMsg.Receiver == k.GetUndelegationModuleAccount(ctx).GetAddress().String() {
		// get all the unbondings for that ibc sequence id
		unbondings := k.FilterUnbondings(
			ctx,
			func(u types.Unbonding) bool {
				return u.IbcSequenceId == k.GetTransactionSequenceID(channel, sequence)
			},
		)

		// update the unbonding ibc sequence id to the transfer id
		for _, unbonding := range unbondings {
			unbonding.IbcSequenceId = k.GetTransactionSequenceID(hc.ChannelId, resp.Sequence)
			k.SetUnbonding(ctx, unbonding)
		}
	}

	if parsedMsg.Sender == hc.DelegationAccount.Address &&
		parsedMsg.Receiver == k.GetDepositModuleAccount(ctx).GetAddress().String() {
		validatorUnbondings := k.FilterValidatorUnbondings(
			ctx,
			func(u types.ValidatorUnbonding) bool {
				return u.ChainId == hc.ChainId && u.IbcSequenceId == k.GetTransactionSequenceID(channel, sequence)
			},
		)

		// remove the unbonding entries as the transfer has succeeded on our part
		for _, validatorUnbonding := range validatorUnbondings {
			k.DeleteValidatorUnbonding(ctx, validatorUnbonding)
		}
	}

	return nil
}
