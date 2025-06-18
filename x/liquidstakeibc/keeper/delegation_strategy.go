package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/gogoproto/proto"

	"github.com/persistenceOne/pstake-native/v3/x/liquidstakeibc/types"
)

type DelegateAmount struct {
	ValAddress string
	Amount     sdk.Dec
}

// GenerateDelegateMessages produces the same result regardless the LSM flag on the host chain.
func (k *Keeper) GenerateDelegateMessages(hc *types.HostChain, depositAmount math.Int) ([]proto.Message, error) {
	// filter out validators which are non-delegable (which reached any LSM cap)
	delegableValidators := make([]*types.Validator, 0)
	nonDelegableWeight := sdk.ZeroDec()
	nonDelegableDelegations := sdk.ZeroInt()
	for _, validator := range hc.Validators {
		if validator.Delegable {
			delegableValidators = append(delegableValidators, validator)
		} else {
			nonDelegableWeight = nonDelegableWeight.Add(validator.Weight)
			nonDelegableDelegations = nonDelegableDelegations.Add(validator.DelegatedAmount)
		}
	}

	// if there are no delegable validators, do nothing
	if len(delegableValidators) == 0 {
		return nil, errorsmod.Wrap(types.ErrInvalidMessages, "no delegable validators")
	}

	// the weight of the un-delegable validators is distributed evenly among the others
	if nonDelegableWeight.GT(sdk.ZeroDec()) {
		weightDelta := nonDelegableWeight.Quo(sdk.NewDec(int64(len(delegableValidators))))
		for _, validator := range delegableValidators {
			validator.Weight = validator.Weight.Add(weightDelta)
		}
	}

	// subtract the delegations from non-delegable validators to get the effective total delegated amount
	effectiveTotalDelegatedAmount := hc.GetHostChainTotalDelegations().Sub(nonDelegableDelegations)

	return k.generateMessages(hc, delegableValidators, effectiveTotalDelegatedAmount, depositAmount, false)
}

func (k *Keeper) GenerateUndelegateMessages(hc *types.HostChain, unbondAmount math.Int) ([]proto.Message, error) {
	return k.generateMessages(hc, hc.Validators, hc.GetHostChainTotalDelegations(), unbondAmount, true)
}

func (k *Keeper) generateMessages(
	hc *types.HostChain,
	validators []*types.Validator,
	totalDelegatedAmount math.Int,
	actionableAmount math.Int,
	undelegating bool,
) ([]proto.Message, error) {
	delegateAmounts := make([]DelegateAmount, 0)
	for _, validator := range validators {
		// calculate the new total delegated amount for the host chain
		futureDelegation := totalDelegatedAmount.Add(actionableAmount)
		if undelegating {
			futureDelegation = totalDelegatedAmount.Sub(actionableAmount)
		}

		// calculate the delegated/undelegated amount difference for the validator:
		//     if the difference is positive, new coins have to be delegated/undelegated
		//     if the difference is zero or negative, don't do anything, it will eventually balance out
		newDelegatedAmount := validator.Weight.Mul(sdk.NewDecFromInt(futureDelegation))
		newDelegationDifference := newDelegatedAmount.Sub(sdk.NewDecFromInt(validator.DelegatedAmount))
		if undelegating {
			newDelegationDifference = sdk.NewDecFromInt(validator.DelegatedAmount).Sub(newDelegatedAmount)
		}
		if newDelegationDifference.LTE(sdk.ZeroDec()) {
			continue // we can't remove delegation from a validator, and we have limited re-stake operations
		}

		delegateAmounts = append(delegateAmounts, DelegateAmount{
			ValAddress: validator.OperatorAddress,
			Amount:     newDelegationDifference,
		})
	}

	messages := make([]proto.Message, 0)
	for i, delegationAmount := range delegateAmounts {
		// create the basic structure of the delegate / undelegate message
		// containing both the delegator and validator addresses
		var message proto.Message
		if !undelegating {
			message = &stakingtypes.MsgDelegate{
				DelegatorAddress: hc.DelegationAccount.Address,
				ValidatorAddress: delegationAmount.ValAddress,
			}
		} else {
			message = &stakingtypes.MsgUndelegate{
				DelegatorAddress: hc.DelegationAccount.Address,
				ValidatorAddress: delegationAmount.ValAddress,
			}
		}

		// if what's left to delegate is less than what needs to be delegated OR we are in the last validator just delegate everything that is left
		// this will also remove any remainder tokens that can be left because of precision issues
		if actionableAmount.LTE(delegationAmount.Amount.TruncateInt()) || i == len(delegateAmounts)-1 {
			if !undelegating {
				msgDelegate := message.(*stakingtypes.MsgDelegate)
				msgDelegate.Amount = sdk.NewCoin(hc.HostDenom, actionableAmount)
			} else {
				msgUndelegate := message.(*stakingtypes.MsgUndelegate)
				msgUndelegate.Amount = sdk.NewCoin(hc.HostDenom, actionableAmount)
			}
			messages = append(messages, message)

			break
		}

		// add the amount to the message and append it
		if !undelegating {
			msgDelegate := message.(*stakingtypes.MsgDelegate)
			msgDelegate.Amount = sdk.NewCoin(hc.HostDenom, delegationAmount.Amount.TruncateInt())
		} else {
			msgUndelegate := message.(*stakingtypes.MsgUndelegate)
			msgUndelegate.Amount = sdk.NewCoin(hc.HostDenom, delegationAmount.Amount.TruncateInt())
		}
		messages = append(messages, message)

		// subtract the amount to delegate/undelegate from the actionable total
		actionableAmount = actionableAmount.Sub(delegationAmount.Amount.TruncateInt())
	}

	if len(messages) == 0 {
		err := errorsmod.Wrap(types.ErrInvalidMessages, "no messages to delegate")
		if undelegating {
			err = errorsmod.Wrap(types.ErrInvalidMessages, "no messages to undelegate")
		}
		return nil, err
	}

	return messages, nil
}
