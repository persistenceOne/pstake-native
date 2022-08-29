package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	lscosmostypes "github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

// DelegateMsgs // Replace this function, does not consider delegation strategy with weights
// CONTRACT: allowlistedValidators should never have 0 elements, amount > len(validators).
func DelegateMsgs(delegatorAddr string, allowlistedValidators lscosmostypes.AllowListedValidators, amount sdk.Int, denom string) []sdk.Msg {
	equalDelegation := amount.QuoRaw(int64(len(allowlistedValidators.AllowListedValidators)))
	change := amount.ModRaw(int64(len(allowlistedValidators.AllowListedValidators)))
	msgs := make([]sdk.Msg, len(allowlistedValidators.AllowListedValidators))
	for i, val := range allowlistedValidators.AllowListedValidators {
		delegationAmount := equalDelegation
		if i == 0 {
			delegationAmount = delegationAmount.Add(change)
		}
		msg := &stakingtypes.MsgDelegate{
			DelegatorAddress: delegatorAddr,
			ValidatorAddress: val.ValidatorAddress,
			Amount:           sdk.NewCoin(denom, delegationAmount),
		}
		msgs[i] = msg
	}
	return msgs
}
