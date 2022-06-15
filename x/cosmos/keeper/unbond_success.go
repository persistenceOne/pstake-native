package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"math"

	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

type UndelegateSuccessKeyAndValue struct {
	ChainIDHeightAndTxHashKey   cosmosTypes.ChainIDHeightAndTxHashKey
	ValueUndelegateSuccessStore cosmosTypes.ValueUndelegateSuccessStore
}

func (k Keeper) setUndelegateSuccessDetails(ctx sdk.Context, msg cosmosTypes.MsgUndelegateSuccess, validatorAddress sdk.ValAddress) {
	undelegateSuccessStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyUndelegateSuccessStore)
	chainIDHeightAndTxHash := cosmosTypes.NewChainIDHeightAndTxHash(msg.ChainID, msg.BlockHeight, msg.TxHash)
	key := k.cdc.MustMarshal(&chainIDHeightAndTxHash)
	totalValidatorCount := k.GetTotalValidatorOrchestratorCount(ctx)

	// check if key present or not
	if !undelegateSuccessStore.Has(key) {
		ratio := sdk.NewDec(1).Quo(sdk.NewDec(totalValidatorCount))
		newValue := cosmosTypes.NewValueUndelegateSuccessStore(msg, validatorAddress, ratio, ctx.BlockHeight()+cosmosTypes.StorageWindow)
		undelegateSuccessStore.Set(key, k.cdc.MustMarshal(&newValue))
		return
	}

	var valueUndelegateSuccessStore cosmosTypes.ValueUndelegateSuccessStore
	k.cdc.MustUnmarshal(undelegateSuccessStore.Get(key), &valueUndelegateSuccessStore)

	// Match if the message value and stored value are same
	// if not equal then initialize by new value in store
	if !StoreValueEqualOrNotUndelegateSuccess(valueUndelegateSuccessStore, msg) {
		ratio := sdk.NewDec(1).Quo(sdk.NewDec(totalValidatorCount))
		newValue := cosmosTypes.NewValueUndelegateSuccessStore(msg, validatorAddress, ratio, ctx.BlockHeight()+cosmosTypes.StorageWindow)
		undelegateSuccessStore.Set(key, k.cdc.MustMarshal(&newValue))
		return
	}

	if !valueUndelegateSuccessStore.Find(validatorAddress.String()) {
		valueUndelegateSuccessStore.UpdateValues(validatorAddress.String(), totalValidatorCount)
		undelegateSuccessStore.Set(key, k.cdc.MustMarshal(&valueUndelegateSuccessStore))
	}
}

func (k Keeper) getAllUndelegateSuccessDetails(ctx sdk.Context) (list []UndelegateSuccessKeyAndValue) {
	undelegateSuccessStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyUndelegateSuccessStore)
	iterator := undelegateSuccessStore.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var chainIDHeightAndTxHashKey cosmosTypes.ChainIDHeightAndTxHashKey
		k.cdc.MustUnmarshal(iterator.Key(), &chainIDHeightAndTxHashKey)

		var valueUndelegateSuccessStore cosmosTypes.ValueUndelegateSuccessStore
		k.cdc.MustUnmarshal(iterator.Value(), &valueUndelegateSuccessStore)

		list = append(list, UndelegateSuccessKeyAndValue{ChainIDHeightAndTxHashKey: chainIDHeightAndTxHashKey, ValueUndelegateSuccessStore: valueUndelegateSuccessStore})
	}
	return list
}

func (k Keeper) deleteUndelegateSuccessDetails(ctx sdk.Context, key cosmosTypes.ChainIDHeightAndTxHashKey) {
	undelegateSuccessStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyUndelegateSuccessStore)
	storeKey := k.cdc.MustMarshal(&key)
	undelegateSuccessStore.Delete(storeKey)
}

func (k Keeper) ProcessAllUndelegateSuccess(ctx sdk.Context) {
	list := k.getAllUndelegateSuccessDetails(ctx)
	epochNumber := k.getLeastEpochNumberWithWithdrawStatusFalse(ctx)
	if epochNumber == int64(math.MaxInt64) {
		return
	}
	for _, element := range list {
		if element.ValueUndelegateSuccessStore.Ratio.GT(cosmosTypes.MinimumRatioForMajority) {
			k.setEpochNumberAndUndelegateDetailsOfIndividualValidator(
				ctx, element.ValueUndelegateSuccessStore.UndelegateSuccess.ValidatorAddress,
				epochNumber, element.ValueUndelegateSuccessStore.UndelegateSuccess.Amount,
			)
		}

		if element.ValueUndelegateSuccessStore.ActiveBlockHeight <= ctx.BlockHeight() {
			k.deleteUndelegateSuccessDetails(ctx, element.ChainIDHeightAndTxHashKey)
		}
	}

	flagForWithdrawSuccess := k.getEpochNumberAndUndelegateDetailsOfValidators(ctx, epochNumber)
	if flagForWithdrawSuccess {
		err := k.emitSendTransactionForAllWithdrawals(ctx, epochNumber)
		if err != nil {
			panic(err)
		}
	}
}

func StoreValueEqualOrNotUndelegateSuccess(storeValue cosmosTypes.ValueUndelegateSuccessStore,
	msgValue cosmosTypes.MsgUndelegateSuccess) bool {
	if storeValue.UndelegateSuccess.DelegatorAddress != msgValue.DelegatorAddress {
		return false
	}
	if storeValue.UndelegateSuccess.ValidatorAddress != msgValue.ValidatorAddress {
		return false
	}
	if !storeValue.UndelegateSuccess.Amount.IsEqual(msgValue.Amount) {
		return false
	}
	if storeValue.UndelegateSuccess.TxHash != msgValue.TxHash {
		return false
	}
	if storeValue.UndelegateSuccess.ChainID != msgValue.ChainID {
		return false
	}
	if storeValue.UndelegateSuccess.BlockHeight != msgValue.BlockHeight {
		return false
	}
	return true
}
