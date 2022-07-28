package keeper

import (
	"math"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

// UndelegateSuccessKeyAndValue :
type UndelegateSuccessKeyAndValue struct {
	ChainIDHeightAndTxHashKey   cosmosTypes.ChainIDHeightAndTxHashKey
	ValueUndelegateSuccessStore cosmosTypes.ValueUndelegateSuccessStore
}

/*
setUndelegateSuccessDetails Adds the undelegate success message entry to the undelegate success store with the given validator address.
Performs the following actions :
  1. Checks if store has the key or not. If not then create new entry
  2. Checks if store has it and matches all the details present in the message. If not then create a new entry.
  3. Finally, if all the details match then append the validator address to keep track.
*/
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

// getAllUndelegateSuccessDetails Gets all the undelegate success details present in the undelegate success store
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

// deleteUndelegateSuccessDetails Removes the given key from the undelegate success store
func (k Keeper) deleteUndelegateSuccessDetails(ctx sdk.Context, key cosmosTypes.ChainIDHeightAndTxHashKey) {
	undelegateSuccessStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyUndelegateSuccessStore)
	storeKey := k.cdc.MustMarshal(&key)
	undelegateSuccessStore.Delete(storeKey)
}

/*
ProcessAllUndelegateSuccess processes all the undelegate success requests
This function is called every EndBlocker to perform the defined set of actions as mentioned below :
   1. Get the list of all undelegate success requests and the last epoch with withdraw status false.
   2. Checks if the majority of the validator orchestrator have sent the minting request.
   3. If majority is reached, set the epch number and undelegate details if individual details of validator.
   4. Another check is present for setting send transaction in outgoing pool.
*/
func (k Keeper) ProcessAllUndelegateSuccess(ctx sdk.Context) {
	list := k.getAllUndelegateSuccessDetails(ctx)
	epochNumber, cValue := k.getLeastEpochNumberWithWithdrawStatusFalse(ctx)
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
		err := k.generateSendTransactionForAllWithdrawals(ctx, epochNumber, cValue)
		if err != nil {
			panic(any(err))
		}
	}
}

// StoreValueEqualOrNotUndelegateSuccess Helper function for undelegate success store to check if the relevant details in the message matches or not.
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
