package keeper

import (
	"fmt"

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
	deposits := k.GetDepositsWithSequenceId(ctx, k.GetDepositSequenceId(channel, sequence))
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

	// send an ICQ query to get the delegator account balance
	if err := k.QueryHostChainAccountBalance(ctx, hc, hc.DelegationAccount.Address); err != nil {
		return fmt.Errorf(
			"error querying host chain %s for delegation account balances: %v",
			hc.ChainId,
			err,
		)
	}

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
