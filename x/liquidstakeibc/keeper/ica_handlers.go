package keeper

import (
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"

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

	k.SetHostChain(ctx, hc)

	// emit an event for the delegation confirmation
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventSuccessfulDelegation,
			sdk.NewAttribute(types.AttributeChainID, hc.ChainId),
			sdk.NewAttribute(types.AttributeDelegatorAddress, parsedMsg.DelegatorAddress),
			sdk.NewAttribute(types.AttributeValidatorAddress, parsedMsg.ValidatorAddress),
			sdk.NewAttribute(types.AttributeDelegatedAmount, sdk.NewCoin(hc.HostDenom, parsedMsg.Amount.Amount).String()),
			sdk.NewAttribute(types.AttributeIBCSequenceID, k.GetTransactionSequenceID(channel, sequence)),
		),
	)

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

		// emit an event for the burned coins
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventBurn,
				sdk.NewAttribute(types.AttributeChainID, hc.ChainId),
				sdk.NewAttribute(types.AttributeTotalEpochBurnAmount, sdk.NewCoin(hc.MintDenom(), unbonding.BurnAmount.Amount).String()),
			),
		)

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

	// emit an event for the undelegation confirmation
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventSuccessfulUndelegation,
			sdk.NewAttribute(types.AttributeChainID, hc.ChainId),
			sdk.NewAttribute(types.AttributeDelegatorAddress, parsedMsg.DelegatorAddress),
			sdk.NewAttribute(types.AttributeValidatorAddress, parsedMsg.ValidatorAddress),
			sdk.NewAttribute(types.AttributeUndelegatedAmount, sdk.NewCoin(hc.HostDenom, parsedMsg.Amount.Amount).String()),
			sdk.NewAttribute(types.AttributeIBCSequenceID, k.GetTransactionSequenceID(channel, sequence)),
		),
	)

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
			"unable to cast msg of type %s to MsgTransfer",
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

		// emit the transfer ack event
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventSuccessfulUndelegationTransfer,
				sdk.NewAttribute(types.AttributeChainID, hc.ChainId),
				sdk.NewAttribute(types.AttributeIBCSequenceID, k.GetTransactionSequenceID(hc.ChannelId, resp.Sequence)),
			),
		)
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

		// emit the transfer ack event
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventSuccessfulValidatorUndelegationTransfer,
				sdk.NewAttribute(types.AttributeChainID, hc.ChainId),
				sdk.NewAttribute(types.AttributeIBCSequenceID, k.GetTransactionSequenceID(hc.ChannelId, resp.Sequence)),
			),
		)
	}

	return nil
}

func (k *Keeper) HandleMsgRedeemTokensForShares(
	ctx sdk.Context,
	msg sdk.Msg,
	resp stakingtypes.MsgRedeemTokensForSharesResponse,
	channel string,
	sequence uint64,
) error {
	parsedMsg, ok := msg.(*stakingtypes.MsgRedeemTokensForShares)
	if !ok {
		return errorsmod.Wrapf(
			sdkerrors.ErrInvalidType,
			"unable to cast msg of type %s to MsgRedeemTokensForShares",
			sdk.MsgTypeURL(msg),
		)
	}

	// remove LSM deposits for this sequence (if any)
	deposits := k.GetLSMDepositsFromIbcSequenceID(ctx, k.GetTransactionSequenceID(channel, sequence))
	for _, deposit := range deposits {
		k.DeleteLSMDeposit(ctx, deposit)
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

	// parse the validator address from the LSM token denom
	operatorAddress, _, found := strings.Cut(parsedMsg.Amount.Denom, "/")
	if !found {
		return errorsmod.Wrapf(
			types.ErrInvalidLSMDenom,
			"could not parse validator address from LSM token %s",
			operatorAddress,
		)
	}

	// get the validator that the delegation was performed to
	validator, found := hc.GetValidator(operatorAddress)
	if !found {
		return errorsmod.Wrapf(
			types.ErrValidatorNotFound,
			"validator with operator address %s not found",
			operatorAddress,
		)
	}

	// update the validator delegated amount
	validator.DelegatedAmount = validator.DelegatedAmount.Add(resp.Amount.Amount)
	k.SetHostChainValidator(ctx, hc, validator)

	k.SetHostChain(ctx, hc)

	// emit an event for the redeem confirmation
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventSuccessfulLSMRedeem,
			sdk.NewAttribute(types.AttributeChainID, hc.ChainId),
			sdk.NewAttribute(types.AttributeDelegatorAddress, parsedMsg.DelegatorAddress),
			sdk.NewAttribute(types.AttributeRedeemedAmount, sdk.NewCoin(hc.HostDenom, parsedMsg.Amount.Amount).String()),
			sdk.NewAttribute(types.AttributeIBCSequenceID, k.GetTransactionSequenceID(channel, sequence)),
		),
	)

	k.Logger(ctx).Info(
		"Received lsm token redeem acknowledgement",
		"delegator",
		parsedMsg.DelegatorAddress,
		"validator",
		operatorAddress,
		"amount",
		resp.Amount.String(),
	)

	return nil
}

func (k *Keeper) HandleMsgBeginRedelegate(
	ctx sdk.Context,
	msg sdk.Msg,
	resp stakingtypes.MsgBeginRedelegateResponse,
	channel string,
	sequence uint64,
) error {
	parsedMsg, ok := msg.(*stakingtypes.MsgBeginRedelegate)
	if !ok {
		return errorsmod.Wrapf(
			sdkerrors.ErrInvalidType,
			"unable to cast msg of type %s to MsgRedeemTokensForShares",
			sdk.MsgTypeURL(msg),
		)
	}
	hc, found := k.GetHostChainFromHostDenom(ctx, parsedMsg.Amount.Denom)
	if !found {
		return errorsmod.Wrapf(
			types.ErrInvalidHostChain,
			"host chain with host denom %s not registered",
			parsedMsg.Amount.Denom,
		)
	}
	// remove redebelgation tx for this sequence (if any)
	tx, ok := k.GetRedelegationTx(ctx, hc.ChainId, k.GetTransactionSequenceID(channel, sequence))
	if !ok {
		k.Logger(ctx).Error("unidentified ica tx acked")
		return nil
	}
	tx.State = types.RedelegateTx_REDELEGATE_ACKED
	k.SetRedelegationTx(ctx, tx)

	// add dst validator tokens
	toValidator, found := hc.GetValidator(parsedMsg.ValidatorDstAddress)
	if !found {
		return errorsmod.Wrapf(
			types.ErrValidatorNotFound,
			"validator with operator address %s not found",
			parsedMsg.ValidatorDstAddress,
		)
	}

	toValidator.DelegatedAmount = toValidator.DelegatedAmount.Add(parsedMsg.Amount.Amount)
	k.SetHostChainValidator(ctx, hc, toValidator)

	// remove src validator tokens
	fromValidator, found := hc.GetValidator(parsedMsg.ValidatorSrcAddress)
	if !found {
		return errorsmod.Wrapf(
			types.ErrValidatorNotFound,
			"validator with operator address %s not found",
			parsedMsg.ValidatorSrcAddress,
		)
	}

	fromValidator.DelegatedAmount = fromValidator.DelegatedAmount.Sub(parsedMsg.Amount.Amount)
	k.SetHostChainValidator(ctx, hc, fromValidator)

	// add redelegation entry.
	k.AddRedelegationEntry(ctx, hc.ChainId, *parsedMsg, resp)

	// emit an event for the redelegation confirmation
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventSuccessfulRedelegation,
			sdk.NewAttribute(types.AttributeChainID, hc.ChainId),
			sdk.NewAttribute(types.AttributeDelegatorAddress, parsedMsg.DelegatorAddress),
			sdk.NewAttribute(types.AttributeValidatorSrcAddress, parsedMsg.ValidatorSrcAddress),
			sdk.NewAttribute(types.AttributeValidatorDstAddress, parsedMsg.ValidatorDstAddress),
			sdk.NewAttribute(types.AttributeRedelegatedAmount, sdk.NewCoin(hc.HostDenom, parsedMsg.Amount.Amount).String()),
			sdk.NewAttribute(types.AttributeIBCSequenceID, k.GetTransactionSequenceID(channel, sequence)),
		),
	)
	k.Logger(ctx).Info(
		"Received redelegate tx acknowledgement",
		"delegator",
		parsedMsg.DelegatorAddress,
		"from-validator",
		parsedMsg.ValidatorSrcAddress,
		"to-validator",
		parsedMsg.ValidatorDstAddress,
		"amount",
		parsedMsg.Amount.String(),
	)

	return nil
}
