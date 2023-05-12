package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/gogo/protobuf/proto"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

type DelegateAmount struct {
	ValAddress string
	ValWeight  sdk.Dec
	Amount     sdk.Dec
}

func (k *Keeper) GenerateDelegateMessages(hc *types.HostChain, depositAmount sdk.Int) ([]proto.Message, error) { //nolint:staticcheck
	return k.generateMessages(hc, depositAmount, false)
	//// calculate the new total delegated amount for the host chain
	//currentDelegation := hc.GetHostChainTotalDelegations()
	//futureDelegation := depositAmount.Add(currentDelegation)
	//
	//delegateAmounts := make([]DelegateAmount, 0)
	//for _, validator := range hc.Validators {
	//	if validator.Weight.Equal(sdk.ZeroDec()) ||
	//		validator.Status != stakingtypes.BondStatusBonded {
	//		continue // skip validators with zero weight or that are not in the active set
	//	}
	//
	//	// calculate the delegated amount difference for the validator:
	//	//     if the difference is positive, new coins have to be delegated
	//	//     if the difference is zero or negative, coins need to be undelegated, but we currently
	//	//     can't un-delegate or re-stake any coins, so don't do anything, it will eventually balance out
	//	newDelegatedAmount := validator.Weight.Mul(sdk.NewDecFromInt(futureDelegation))
	//	newDelegationDifference := newDelegatedAmount.Sub(sdk.NewDecFromInt(validator.DelegatedAmount))
	//	if newDelegationDifference.LTE(sdk.ZeroDec()) {
	//		continue // we can't remove delegation from a validator, and we have limited re-stake operations
	//	}
	//
	//	delegateAmounts = append(delegateAmounts, DelegateAmount{
	//		ValAddress: validator.OperatorAddress,
	//		ValWeight:  validator.Weight,
	//		Amount:     newDelegationDifference,
	//	})
	//}
	//
	//messages := make([]proto.Message, 0)
	//for _, delegationAmount := range delegateAmounts {
	//	message := &stakingtypes.MsgDelegate{
	//		DelegatorAddress: hc.DelegationAccount.Address,
	//		ValidatorAddress: delegationAmount.ValAddress,
	//	}
	//
	//	// return when there is nothing more to delegate
	//	if depositAmount.LTE(delegationAmount.Amount.TruncateInt()) {
	//		message.Amount = sdk.NewCoin(hc.HostDenom, depositAmount)
	//		messages = append(messages, message)
	//		return messages, nil
	//	}
	//
	//	// add the amount to the message and append it
	//	message.Amount = sdk.NewCoin(hc.HostDenom, delegationAmount.Amount.TruncateInt())
	//	messages = append(messages, message)
	//
	//	// subtract the amount to delegate from the total deposited
	//	depositAmount = depositAmount.Sub(delegationAmount.Amount.TruncateInt())
	//}
	//
	//if len(messages) == 0 {
	//	return nil, errorsmod.Wrap(types.ErrInvalidMessages, "no messages to delegate")
	//}
	//
	//return messages, nil
}

func (k *Keeper) GenerateUndelegateMessages(hc *types.HostChain, unbondAmount sdk.Int) ([]proto.Message, error) { //nolint:staticcheck
	return k.generateMessages(hc, unbondAmount, true)
}

func (k *Keeper) generateMessages(
	hc *types.HostChain,
	actionableAmount sdk.Int,
	undelegating bool,
) ([]proto.Message, error) { //nolint:staticcheck
	delegateAmounts := make([]DelegateAmount, 0)
	for _, validator := range hc.Validators {
		// calculate the new total delegated amount for the host chain
		currentDelegation := hc.GetHostChainTotalDelegations()
		futureDelegation := currentDelegation.Add(actionableAmount)
		if undelegating {
			futureDelegation = currentDelegation.Sub(actionableAmount)
		}

		if validator.Weight.Equal(sdk.ZeroDec()) ||
			validator.Status != stakingtypes.BondStatusBonded {
			continue // skip validators with zero weight or that are not in the active set
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
			ValWeight:  validator.Weight,
			Amount:     newDelegationDifference,
		})
	}

	messages := make([]proto.Message, 0)
	for _, delegationAmount := range delegateAmounts {
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

		// return when there is nothing more to delegate/undelegate
		if actionableAmount.LTE(delegationAmount.Amount.TruncateInt()) {
			if !undelegating {
				msgDelegate := message.(*stakingtypes.MsgDelegate)
				msgDelegate.Amount = sdk.NewCoin(hc.HostDenom, actionableAmount)
			} else {
				msgUndelegate := message.(*stakingtypes.MsgUndelegate)
				msgUndelegate.Amount = sdk.NewCoin(hc.HostDenom, actionableAmount)
			}
			messages = append(messages, message)

			return messages, nil
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
