package keeper

import (
	"fmt"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkTx "github.com/cosmos/cosmos-sdk/types/tx"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
	"math"
)

func (k Keeper) generateUnbondingOutgoingEvent(ctx sdk.Context, listOfValidatorsAndUnbondingAmount []ValAddressAndAmountForStakingAndUnstaking, epochNumber int64) {
	params := k.GetParams(ctx)

	chunkMsgs := ChunkUndelegationSlice(listOfValidatorsAndUnbondingAmount, params.ChunkSize)

	for _, chunk := range chunkMsgs {
		nextID := k.autoIncrementID(ctx, []byte(cosmosTypes.KeyLastTXPoolID))

		var undelegateMsgsAny []*codecTypes.Any
		var undelegategMsgs []stakingTypes.MsgUndelegate
		for _, element := range chunk {
			msg := stakingTypes.MsgUndelegate{
				DelegatorAddress: params.CustodialAddress,
				ValidatorAddress: element.validator.String(),
				Amount:           element.amount,
			}
			anyMsg, err := codecTypes.NewAnyWithValue(&msg)
			if err != nil {
				panic(err)
			}
			undelegateMsgsAny = append(undelegateMsgsAny, anyMsg)
			undelegategMsgs = append(undelegategMsgs, msg)
		}

		tx := cosmosTypes.CosmosTx{
			Tx: sdkTx.Tx{
				Body: &sdkTx.TxBody{
					Messages:      undelegateMsgsAny,
					Memo:          "",
					TimeoutHeight: 0,
				},
				AuthInfo: &sdkTx.AuthInfo{
					SignerInfos: nil,
					Fee: &sdkTx.Fee{
						Amount:   nil,
						GasLimit: 200000,
						Payer:    "",
					},
				},
				Signatures: nil,
			},
			EventEmitted:      true,
			Status:            "",
			TxHash:            "",
			NativeBlockHeight: ctx.BlockHeight(),
			ActiveBlockHeight: ctx.BlockHeight() + cosmosTypes.StorageWindow,
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				cosmosTypes.EventTypeOutgoing,
				sdk.NewAttribute(cosmosTypes.AttributeKeyOutgoingTXID, fmt.Sprint(nextID)),
			),
		)

		err := k.setIDInEpochPoolForWithdrawals(ctx, nextID, undelegategMsgs, params.CustodialAddress, epochNumber)
		if err != nil {
			panic(err)
		}
		//Once event is emitted, store it in KV store for orchestrators to query transactions and sign them
		k.setNewTxnInOutgoingPool(ctx, nextID, tx)
	}
}

func (k Keeper) setIDInEpochPoolForWithdrawals(ctx sdk.Context, txID uint64, undelegateMsgs []stakingTypes.MsgUndelegate, custodialAddress string, epochNumber int64) error {
	unbondingEpochStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyOutgoingUnbondStore)
	key := cosmosTypes.UInt64Bytes(txID)
	value := cosmosTypes.NewValueOutgoingUnbondStore(undelegateMsgs, epochNumber)
	bz, err := value.Marshal()
	if err != nil {
		return err
	}
	unbondingEpochStore.Set(key, bz)
	return nil
}

func (k Keeper) getIDInEpochPoolForWithdrawals(ctx sdk.Context, txID uint64) (value cosmosTypes.ValueOutgoingUnbondStore, err error) {
	unbondingEpochStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyOutgoingUnbondStore)
	key := cosmosTypes.UInt64Bytes(txID)
	err = value.Unmarshal(unbondingEpochStore.Get(key))
	if err != nil {
		return cosmosTypes.ValueOutgoingUnbondStore{}, err
	}
	return value, nil
}

func (k Keeper) deleteIDInEpochPoolForWithdrawals(ctx sdk.Context, txID uint64) {
	unbondingEpochStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyOutgoingUnbondStore)
	key := cosmosTypes.UInt64Bytes(txID)
	unbondingEpochStore.Delete(key)
}

//_____________________________________________________________________________________

func (k Keeper) setEpochWithdrawSuccessStore(ctx sdk.Context, epochNumber int64) {
	epochWithdrawStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyEpochStoreForWithdrawSuccess)
	key := cosmosTypes.Int64Bytes(epochNumber)
	if !epochWithdrawStore.Has(key) {
		epochWithdrawStore.Set(key, []byte("false"))
	}
}

func (k Keeper) getEpochWithdrawSuccessStore(ctx sdk.Context, epochNumber int64) string {
	epochWithdrawStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyEpochStoreForWithdrawSuccess)
	key := cosmosTypes.Int64Bytes(epochNumber)
	bz := epochWithdrawStore.Get(key)
	return string(bz)
}

func (k Keeper) getLeastEpochNumberWithWithdrawStatusFalse(ctx sdk.Context) int64 {
	epochWithdrawStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyEpochStoreForWithdrawSuccess)
	iterator := epochWithdrawStore.Iterator(nil, nil)
	min := int64(math.MaxInt64)
	for ; iterator.Valid(); iterator.Next() {
		if string(iterator.Value()) == "false" {
			epochNumber := cosmosTypes.Int64FromBytes(iterator.Key())
			if min > epochNumber {
				min = epochNumber
			}
		}
	}
	return min
}

func (k Keeper) deleteEpochWithdrawSuccessStore(ctx sdk.Context, epochNumber int64) {
	epochWithdrawStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyEpochStoreForWithdrawSuccess)
	key := cosmosTypes.Int64Bytes(epochNumber)
	epochWithdrawStore.Delete(key)
}

//___________________________________________________________________________________________

func (k Keeper) setEpochNumberAndUndelegateDetailsOfValidators(ctx sdk.Context, details cosmosTypes.ValueOutgoingUnbondStore) {
	epochNumberStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.GetEpochStoreForUndelegationKey(details.EpochNumber))
	for _, element := range details.UndelegateMessages {
		a := append([]byte(element.ValidatorAddress), []byte(element.Amount.String())...)
		epochNumberStore.Set(a, []byte("false"))
	}
}

func (k Keeper) setEpochNumberAndUndelegateDetailsOfIndividualValidator(ctx sdk.Context, validatorAddress string, epochNumber int64, amount sdk.Coin) {
	epochNumberStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.GetEpochStoreForUndelegationKey(epochNumber))
	a := append([]byte(validatorAddress), []byte(amount.String())...)
	epochNumberStore.Set(a, []byte("true"))
}

func (k Keeper) getEpochNumberAndUndelegateDetailsOfValidators(ctx sdk.Context, epochNumber int64) bool {
	epochNumberStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.GetEpochStoreForUndelegationKey(epochNumber))
	iterator := epochNumberStore.Iterator(nil, nil)
	counter := 0
	for ; iterator.Valid(); iterator.Next() {
		counter++
		if string(iterator.Value()) == "false" {
			return false
		}
	}
	if counter > 0 {
		return true
	}
	return false
}

func (k Keeper) setEpochAndValidatorDetailsForAllUndelegations(ctx sdk.Context, txID uint64) error {
	details, err := k.getIDInEpochPoolForWithdrawals(ctx, txID)
	if err != nil {
		return err
	}
	k.setEpochNumberAndUndelegateDetailsOfValidators(ctx, details) //sets undelegations txns for future verifications
	k.deleteIDInEpochPoolForWithdrawals(ctx, txID)
	k.setEpochWithdrawSuccessStore(ctx, details.EpochNumber) //sets withdraw batch success as false
	return nil
}

func ChunkUndelegationSlice(slice []ValAddressAndAmountForStakingAndUnstaking, chunkSize int64) (chunks [][]ValAddressAndAmountForStakingAndUnstaking) {
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
