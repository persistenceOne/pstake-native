package keeper

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
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

type ValAddressAndAmountForStakingAndUndelegating struct {
	validator sdk.ValAddress
	amount    sdk.Coin
}

// gives a list of all validators having weighted amount for few and 1uatom for rest in order to auto claim all rewards accumulated in current epoch
func (k Keeper) fetchValidatorsToDelegate(ctx sdk.Context, amount sdk.Coin) []ValAddressAndAmountForStakingAndUndelegating {
	//TODO : Add pseudo code for filtering out validators to delegate
	return nil
}

// gives a list of validators having weighted amount for few validators
func (k Keeper) fetchValidatorsToUndelegate(ctx sdk.Context, amount sdk.Coin) []ValAddressAndAmountForStakingAndUndelegating {
	//TODO : Implement opposite of fetchValidatorsToDelegate
	return nil
}
