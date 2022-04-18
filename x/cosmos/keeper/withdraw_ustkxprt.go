package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

// adds details to withdraw pool for ubonding epoch
func (k Keeper) addToWithdrawPool(ctx sdk.Context, asset cosmosTypes.MsgWithdrawStkAsset) error {
	withdrawStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyWithdrawStore)
	currentEpoch := k.epochsKeeper.GetEpochInfo(ctx, k.GetParams(ctx).UndelegateEpochIdentifier).CurrentEpoch
	key := cosmosTypes.Int64Bytes(currentEpoch)
	if withdrawStore.Has(key) {
		bz := withdrawStore.Get(key)
		if bz == nil {
			return fmt.Errorf("withdraw store has key but nothing assigned to it")
		}
		var withdrawStoreValue cosmosTypes.WithdrawStoreValue
		err := withdrawStoreValue.Unmarshal(bz)
		if err != nil {
			return err
		}
		withdrawStoreValue.WithdrawDetails = append(withdrawStoreValue.WithdrawDetails, asset)
		withdrawStoreValue.UnbondEmitFlag = append(withdrawStoreValue.UnbondEmitFlag, false)

		bz1, err := withdrawStoreValue.Marshal()
		if err != nil {
			return err
		}
		withdrawStore.Set(key, bz1)
	} else {
		withdrawDetails := cosmosTypes.NewWithdrawStoreValue(asset)
		bz, err := withdrawDetails.Marshal()
		if err != nil {
			return err
		}
		withdrawStore.Set(key, bz)
	}
	return nil
}

func (k Keeper) fetchWithdrawTxnsWithCurrentEpochInfo(ctx sdk.Context, currentEpoch int64) (withdrawStoreValue cosmosTypes.WithdrawStoreValue, err error) {
	withdrawStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyWithdrawStore)
	bz := withdrawStore.Get(cosmosTypes.Int64Bytes(currentEpoch))
	err = withdrawStoreValue.Unmarshal(bz)
	if err != nil {
		return cosmosTypes.WithdrawStoreValue{}, err
	}
	return withdrawStoreValue, nil
}

func (k Keeper) totalAmountToBeUnbonded(value cosmosTypes.WithdrawStoreValue, denom string) sdk.Coin {
	amount := sdk.NewInt64Coin(denom, 0)
	for _, element := range value.WithdrawDetails {
		amount = amount.Add(sdk.NewCoin(denom, element.Amount.Amount))
	}
	return amount
}
