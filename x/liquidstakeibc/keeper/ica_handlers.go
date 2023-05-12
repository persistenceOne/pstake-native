package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	liquidstakeibctypes "github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
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
			liquidstakeibctypes.ErrInvalidHostChain,
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
			liquidstakeibctypes.ErrValidatorNotFound,
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
	unbondings := k.GetAllUnbondingsForSequenceID(ctx, k.GetTransactionSequenceID(channel, sequence))

	for _, unbonding := range unbondings {
		// burn the undelegated stk tokens
		err := k.bankKeeper.BurnCoins(
			ctx,
			liquidstakeibctypes.UndelegationModuleAccount,
			sdk.NewCoins(unbonding.BurnAmount),
		)
		if err != nil {
			return err
		}

		// set the mature time for the undelegation
		unbonding.MatureTime = resp.CompletionTime
		k.SetUnbonding(ctx, unbonding)

		parsedMsg, ok := msg.(*stakingtypes.MsgUndelegate)
		if !ok {
			return errorsmod.Wrapf(
				sdkerrors.ErrInvalidType,
				"unable to cast msg of type %s to MsgUndelegate",
				sdk.MsgTypeURL(msg),
			)
		}

		// TODO: Remove delegation from validators
		// TODO: Update host chain C value

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
	return nil
}

func (k *Keeper) HandleSetWithdrawAddressResponse(ctx sdk.Context, msg sdk.Msg) error {
	return nil
}
