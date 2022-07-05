package keeper

import (
	"math"

	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkTx "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/x/authz"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

// generateUnbondingOutgoingTxn Generate unbonding outgoing transaction and set in outoging pool with given txID
func (k Keeper) generateUnbondingOutgoingTxn(ctx sdk.Context, listOfValidatorsAndUnbondingAmount []ValAddressAmount,
	epochNumber int64, cValue sdk.Dec) {
	params := k.GetParams(ctx)

	chunkMsgs := ChunkDelegateAndUndelegateSlice(listOfValidatorsAndUnbondingAmount, params.ChunkSize)

	for _, chunk := range chunkMsgs {
		nextID := k.autoIncrementID(ctx, []byte(cosmosTypes.KeyLastTXPoolID))

		var undelegateMsgsAny []*codecTypes.Any
		var undelegategMsgs []stakingTypes.MsgUndelegate
		for _, element := range chunk {
			msg := stakingTypes.MsgUndelegate{
				DelegatorAddress: params.CustodialAddress,
				ValidatorAddress: element.Validator,
				Amount:           element.Amount,
			}
			anyMsg, err := codecTypes.NewAnyWithValue(&msg)
			if err != nil {
				panic(err)
			}
			undelegateMsgsAny = append(undelegateMsgsAny, anyMsg)
			undelegategMsgs = append(undelegategMsgs, msg)
		}

		cosmosAddrr, err := cosmosTypes.Bech32ifyAddressBytes(cosmosTypes.Bech32PrefixAccAddr, k.GetCurrentAddress(ctx))
		if err != nil {
			panic(err)
		}
		execMsg := authz.MsgExec{
			Grantee: cosmosAddrr,
			Msgs:    undelegateMsgsAny,
		}

		execMsgAny, err := codecTypes.NewAnyWithValue(&execMsg)
		if err != nil {
			panic(err)
		}

		tx := cosmosTypes.CosmosTx{
			Tx: sdkTx.Tx{
				Body: &sdkTx.TxBody{
					Messages:      []*codecTypes.Any{execMsgAny},
					Memo:          "",
					TimeoutHeight: 0,
				},
				AuthInfo: &sdkTx.AuthInfo{
					SignerInfos: nil,
					Fee: &sdkTx.Fee{
						Amount:   nil,
						GasLimit: 400000,
						Payer:    "",
					},
				},
				Signatures: nil,
			},
			EventEmitted:      false,
			Status:            "",
			TxHash:            "",
			ActiveBlockHeight: ctx.BlockHeight() + cosmosTypes.StorageWindow,
			SignerAddress:     cosmosAddrr,
		}

		err = k.setIDInEpochPoolForWithdrawals(ctx, nextID, undelegategMsgs, epochNumber, cValue)
		if err != nil {
			panic(err)
		}
		//Once event is emitted, store it in KV store for orchestrators to query transactions and sign them
		k.SetNewTxnInOutgoingPool(ctx, nextID, tx)

		k.setNewInTransactionQueue(ctx, nextID)
	}
}

// setIDInEpochPoolForWithdrawals Sets ID in epoch pool for withdrawals for the given aaray of undelegate messages
func (k Keeper) setIDInEpochPoolForWithdrawals(ctx sdk.Context, txID uint64,
	undelegateMsgs []stakingTypes.MsgUndelegate, epochNumber int64, cValue sdk.Dec) error {
	unbondingEpochStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyOutgoingUnbondStore)
	key := cosmosTypes.UInt64Bytes(txID)
	value := cosmosTypes.NewValueOutgoingUnbondStore(undelegateMsgs, epochNumber, cValue)
	bz, err := k.cdc.Marshal(&value)
	if err != nil {
		return err
	}
	unbondingEpochStore.Set(key, bz)
	return nil
}

// getIDInEpochPoolForWithdrawals Gets the details corresponding to the given txID in the epoch pool for withdrawals
func (k Keeper) getIDInEpochPoolForWithdrawals(ctx sdk.Context, txID uint64) (value cosmosTypes.ValueOutgoingUnbondStore) {
	unbondingEpochStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyOutgoingUnbondStore)
	key := cosmosTypes.UInt64Bytes(txID)
	k.cdc.MustUnmarshal(unbondingEpochStore.Get(key), &value)
	return value
}

// deleteIDInEpochPoolForWithdrawals Removes the details corresponding to the given ID in the epoch pool for withdrawals
func (k Keeper) deleteIDInEpochPoolForWithdrawals(ctx sdk.Context, txID uint64) {
	unbondingEpochStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyOutgoingUnbondStore)
	key := cosmosTypes.UInt64Bytes(txID)
	unbondingEpochStore.Delete(key)
}

//_____________________________________________________________________________________

// setEpochWithdrawSuccessStore Set given epoch number with status "false" in epoch withdraw success store
func (k Keeper) setEpochWithdrawSuccessStore(ctx sdk.Context, epochNumber int64, cValue sdk.Dec) {
	epochWithdrawStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyEpochStoreForWithdrawSuccess)
	key := cosmosTypes.Int64Bytes(epochNumber)
	if !epochWithdrawStore.Has(key) {
		value := cosmosTypes.NewEpochWithdrawSuccessStoreValue(cValue)
		epochWithdrawStore.Set(key, k.cdc.MustMarshal(&value))
	}
}

// getLeastEpochNumberWithWithdrawStatusFalse Gets the least epoch number with withdraw status false from withdraw success store
func (k Keeper) getLeastEpochNumberWithWithdrawStatusFalse(ctx sdk.Context) (epochNumber int64, cValue sdk.Dec) {
	epochWithdrawStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyEpochStoreForWithdrawSuccess)
	iterator := epochWithdrawStore.Iterator(nil, nil)
	defer iterator.Close()
	min := int64(math.MaxInt64)
	for ; iterator.Valid(); iterator.Next() {
		var storeValue cosmosTypes.EpochWithdrawSuccessStoreValue
		k.cdc.MustUnmarshal(iterator.Value(), &storeValue)
		if !storeValue.Status {
			epochNumber = cosmosTypes.Int64FromBytes(iterator.Key())
			if min > epochNumber {
				min = epochNumber
				cValue = storeValue.CValue
			}
		}
	}
	return min, cValue
}

// deleteEpochWithdrawSuccessStore Removes the given epoch number entry from the epoch withdraw success store
func (k Keeper) deleteEpochWithdrawSuccessStore(ctx sdk.Context, epochNumber int64) {
	epochWithdrawStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyEpochStoreForWithdrawSuccess)
	key := cosmosTypes.Int64Bytes(epochNumber)
	epochWithdrawStore.Delete(key)
}

//___________________________________________________________________________________________

// setEpochNumberAndUndelegateDetailsOfValidators Set epoch number and undelegate details of given validators in the epoch number store
func (k Keeper) setEpochNumberAndUndelegateDetailsOfValidators(ctx sdk.Context, details cosmosTypes.ValueOutgoingUnbondStore) {
	epochNumberStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.GetEpochStoreForUndelegationKey(details.EpochNumber))
	for _, element := range details.UndelegateMessages {
		a := append([]byte(element.ValidatorAddress), []byte(element.Amount.String())...)
		epochNumberStore.Set(a, []byte("false"))
	}
}

// setEpochNumberAndUndelegateDetailsOfIndividualValidator Set epoch number and undelegate details of validator in the epoch number store
func (k Keeper) setEpochNumberAndUndelegateDetailsOfIndividualValidator(ctx sdk.Context, validatorAddress string, epochNumber int64, amount sdk.Coin) {
	epochNumberStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.GetEpochStoreForUndelegationKey(epochNumber))
	a := append([]byte(validatorAddress), []byte(amount.String())...)
	epochNumberStore.Set(a, []byte("true"))
}

// getEpochNumberAndUndelegateDetailsOfValidators Gets the undelegate details of validator mapped to the given epoch number in the epoch number store
func (k Keeper) getEpochNumberAndUndelegateDetailsOfValidators(ctx sdk.Context, epochNumber int64) bool {
	epochNumberStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.GetEpochStoreForUndelegationKey(epochNumber))
	iterator := epochNumberStore.Iterator(nil, nil)
	defer iterator.Close()
	counter := 0
	for ; iterator.Valid(); iterator.Next() {
		counter++
		if string(iterator.Value()) == "false" {
			return false
		}
	}
	return true
}

// deleteEpochNumberAndUndelegateDetailsOfValidators deletes epoch number and unbonding details from store
func (k Keeper) deleteEpochNumberAndUndelegateDetailsOfValidators(ctx sdk.Context, epochNumber int64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(cosmosTypes.GetEpochStoreForUndelegationKey(epochNumber))

}

// setEpochAndValidatorDetailsForAllUndelegations set the validator details for all MsgUndelegate success entries
func (k Keeper) setEpochAndValidatorDetailsForAllUndelegations(ctx sdk.Context, txID uint64) {
	details := k.getIDInEpochPoolForWithdrawals(ctx, txID)
	k.setEpochNumberAndUndelegateDetailsOfValidators(ctx, details) //sets undelegations txns for future verifications
	k.deleteIDInEpochPoolForWithdrawals(ctx, txID)
	k.setEpochWithdrawSuccessStore(ctx, details.EpochNumber, details.CValue) //sets withdraw batch success as false
}

// ChunkDelegateAndUndelegateSlice divides 1D slice of ValAddressAmount into chunks of given size and
// returns it by putting it in a 2D slice
func ChunkDelegateAndUndelegateSlice(slice []ValAddressAmount, chunkSize int64) (chunks [][]ValAddressAmount) {
	for {
		if len(slice) == 0 {
			break
		}

		// necessary check to avoid slicing beyond
		// slice capacity
		if int64(len(slice)) < chunkSize {
			chunkSize = int64(len(slice))
		}

		chunks = append(chunks, slice[0:chunkSize])
		slice = slice[chunkSize:]
	}

	return chunks
}
