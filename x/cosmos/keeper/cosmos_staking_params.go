package keeper

import (
	"encoding/json"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

type InternalWeightedAddressCosmos []cosmosTypes.WeightedAddressCosmos

var _ sort.Interface = InternalWeightedAddressCosmos{}

func ConvertToInternalWeightedAddressCosmos(weightedAddressCosmos []cosmosTypes.WeightedAddressCosmos) (internalWeightedAddressCosmos InternalWeightedAddressCosmos) {
	for _, element := range weightedAddressCosmos {
		internalWeightedAddressCosmos = append(internalWeightedAddressCosmos, element)
	}
	return internalWeightedAddressCosmos
}

func (w InternalWeightedAddressCosmos) Len() int {
	return len(w)
}

func (w InternalWeightedAddressCosmos) Less(i, j int) bool {
	// TODO refactor
	//return w[i].Difference.Amount.Uint64() < w[j].Difference.Amount.Uint64()
	return false
}

func (w InternalWeightedAddressCosmos) Swap(i, j int) {
	w[i], w[j] = w[j], w[i]
}

func (w InternalWeightedAddressCosmos) Sort() InternalWeightedAddressCosmos {
	sort.Sort(w)
	return w
}

func (w InternalWeightedAddressCosmos) Marshal() ([]byte, error) {
	if w == nil {
		return json.Marshal(InternalWeightedAddressCosmos{})
	}
	return json.Marshal(w)
}

func (w InternalWeightedAddressCosmos) Unmarshal(bz []byte) error {
	err := json.Unmarshal(bz, &w)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) setCosmosValidatorParams(ctx sdk.Context, details InternalWeightedAddressCosmos) {
	store := ctx.KVStore(k.storeKey)
	key := []byte(cosmosTypes.KeyCosmosValidatorSet)
	if store.Has(key) {
		bz, err := details.Sort().Marshal()
		if err != nil {
			panic(err)
		}
		store.Set(key, bz)
	} else {
		newWeightedAddress := ConvertToInternalWeightedAddressCosmos(k.GetParams(ctx).ValidatorSetCosmosChain)
		bz, err := newWeightedAddress.Sort().Marshal()
		if err != nil {
			panic(err)
		}
		store.Set(key, bz)
	}
}

func (k Keeper) getCosmosValidatorParams(ctx sdk.Context) (internalWeightedAddressCosmos InternalWeightedAddressCosmos) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(cosmosTypes.KeyCosmosValidatorSet))
	err := internalWeightedAddressCosmos.Unmarshal(bz)
	if err != nil {
		panic(err)
	}
	return internalWeightedAddressCosmos
}

func (k Keeper) updateCosmosValidatorStakingParams(ctx sdk.Context, msgs []sdk.Msg) error {
	uatomDenom, err := k.GetParams(ctx).GetBondDenomOf("uatom")
	if err != nil {
		return err
	}
	totalAmountInDelegateMsgs := sdk.NewInt64Coin(uatomDenom, 0)
	msgsMap := make(map[string]stakingTypes.MsgDelegate, len(msgs))
	for _, msg := range msgs {
		delegateMsg := msg.(*stakingTypes.MsgDelegate)
		totalAmountInDelegateMsgs = totalAmountInDelegateMsgs.Add(delegateMsg.Amount)
		msgsMap[delegateMsg.ValidatorAddress] = *delegateMsg
	}

	k.setTotalDelegatedAmountTillDate(ctx, totalAmountInDelegateMsgs)

	internalWeightedAddressCosmos := k.getCosmosValidatorParams(ctx)
	for _, element := range internalWeightedAddressCosmos {
		if val, ok := msgsMap[element.Address]; ok {
			element.CurrentDelegatedAmount.Add(val.Amount)
			// TODO refactor this, difference and ideal delegated amount was deleted.
			//element.IdealDelegatedAmount = sdk.NewCoin(element.IdealDelegatedAmount.Denom,
			//	k.getTotalDelegatedAmountTillDate(ctx).Amount.ToDec().Mul(element.Weight).TruncateInt(),
			//)
			//element.Difference = element.IdealDelegatedAmount.Sub(element.CurrentDelegatedAmount)
		}
	}
	k.setCosmosValidatorParams(ctx, internalWeightedAddressCosmos)
	return nil
	//TODO : Update c token ratio
}

type ValAddressAndAmountForStakingAndUnstaking struct {
	validator sdk.ValAddress
	amount    sdk.Coin
}

func (k Keeper) fetchValidatorsToDelegate(ctx sdk.Context, amount sdk.Coin) []ValAddressAndAmountForStakingAndUnstaking {
	//internalWeightedAddressCosmos := k.getCosmosValidatorParams(ctx)
	//uatomAmount := amount.AmountOf(cosmosTypes.StakeDenom)
	//for _, element := range internalWeightedAddressCosmos {
	//	delegationThreshold := k.GetParams(ctx).DelegationThreshold
	//	//process element
	//	//TODO : Add pseudo code for filtering out validators to delegate
	//}
	return nil
}

func (k Keeper) fetchValidatorsToUndelegate(ctx sdk.Context, amount sdk.Coin) []ValAddressAndAmountForStakingAndUnstaking {
	//TODO : Implement opposite of fetchValidatorsToDelegate
	return nil
}
