package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/gogo/protobuf/proto"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

func (k *Keeper) GenerateDelegateMessages(hc *types.HostChain, newDelegation sdk.Int) ([]proto.Message, error) {
	// calculate the new total delegated amount for the host chain
	currentDelegation := hc.GetHostChainTotalDelegations()
	futureDelegation := newDelegation.Add(currentDelegation)

	var messages []proto.Message
	for _, validator := range hc.Validators {
		if validator.Weight.Equal(sdk.ZeroDec()) ||
			validator.Status != stakingtypes.BondStatusBonded {
			continue // skip validators with zero weight or that are not in the active set
		}

		// calculate the delegated amount difference for the validator:
		//     if the difference is positive, new coins have to be delegated
		//     if the difference is zero or negative, coins need to be undelegated, but we currently
		//     can't un-delegate or re-stake any coins, so don't do anything, it will eventually balance out
		newDelegatedAmount := validator.Weight.Mul(sdk.NewDecFromInt(futureDelegation))
		newDelegationDifference := newDelegatedAmount.Sub(sdk.NewDecFromInt(validator.DelegatedAmount))
		if newDelegationDifference.LTE(sdk.ZeroDec()) {
			continue // we can't remove delegation from a validator, and we have limited re-stake operations
		}

		// create the delegate message
		messages = append(
			messages,
			&stakingtypes.MsgDelegate{
				DelegatorAddress: hc.DelegationAccount.Address,
				ValidatorAddress: validator.OperatorAddress,
				Amount:           sdk.NewCoin(hc.HostDenom, newDelegationDifference.TruncateInt()),
			},
		)
	}

	if len(messages) == 0 {
		return nil, errorsmod.Wrap(types.ErrInvalidMessages, "no messages to delegate")
	}

	return messages, nil
}
