package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
	"math"
)

type UndelegateSuccessKeyAndValue struct {
	ChainIDHeightAndTxHashKey   cosmosTypes.ChainIDHeightAndTxHashKey
	ValueUndelegateSuccessStore cosmosTypes.ValueUndelegateSuccessStore
}

func (k Keeper) setUndelegateSuccessDetails(ctx sdk.Context, validatorAddress sdk.ValAddress, orchestratorAddress sdk.AccAddress, amount sdk.Coin, txHash string, chainID string, blockHeight int64) error {
	undelegateSuccessStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyUndelegateSuccessStore)
	chainIDHeightAndTxHash := cosmosTypes.NewChainIDHeightAndTxHash(chainID, blockHeight, txHash)
	key, err := k.cdc.Marshal(&chainIDHeightAndTxHash)
	if err != nil {
		return err
	}
	if undelegateSuccessStore.Has(key) {
		var valueUndelegateSuccessStore cosmosTypes.ValueUndelegateSuccessStore
		err = k.cdc.Unmarshal(undelegateSuccessStore.Get(key), &valueUndelegateSuccessStore)
		if err != nil {
			panic("error in unmarshalling valueUndelegateSuccessStore")
		}
		if !valueUndelegateSuccessStore.Find(orchestratorAddress.String()) {
			valueUndelegateSuccessStore.OrchestratorAddresses = append(valueUndelegateSuccessStore.OrchestratorAddresses, orchestratorAddress.String())
			valueUndelegateSuccessStore.Counter++
			valueUndelegateSuccessStore.Ratio = float32(valueUndelegateSuccessStore.Counter) / float32(k.getTotalValidatorOrchestratorCount(ctx))
			bz, err := k.cdc.Marshal(&valueUndelegateSuccessStore)
			if err != nil {
				panic("error in marshaling txHashValue")
			}
			undelegateSuccessStore.Set(key, bz)
		}
	} else {
		ratio := float32(1) / float32(k.getTotalValidatorOrchestratorCount(ctx))
		newValue := cosmosTypes.NewValueUndelegateSuccessStore(validatorAddress, orchestratorAddress, ratio, amount, ctx.BlockHeight(), ctx.BlockHeight()+cosmosTypes.StorageWindow)
		bz, err := k.cdc.Marshal(&newValue)
		if err != nil {
			panic("error in marshaling valueUndelegateSuccessStore")
		}
		undelegateSuccessStore.Set(key, bz)
	}
	return nil
}

func (k Keeper) getAllUndelegateSuccessDetails(ctx sdk.Context) (list []UndelegateSuccessKeyAndValue, err error) {
	undelegateSuccessStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyUndelegateSuccessStore)
	iterator := undelegateSuccessStore.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {
		var chainIDHeightAndTxHashKey cosmosTypes.ChainIDHeightAndTxHashKey
		err = k.cdc.Unmarshal(iterator.Key(), &chainIDHeightAndTxHashKey)
		if err != nil {
			return nil, err
		}

		var valueUndelegateSuccessStore cosmosTypes.ValueUndelegateSuccessStore
		err = k.cdc.Unmarshal(iterator.Value(), &valueUndelegateSuccessStore)
		if err != nil {
			return nil, err
		}
		list = append(list, UndelegateSuccessKeyAndValue{ChainIDHeightAndTxHashKey: chainIDHeightAndTxHashKey, ValueUndelegateSuccessStore: valueUndelegateSuccessStore})
	}
	return list, nil
}

func (k Keeper) deleteUndelegateSuccessDetails(ctx sdk.Context, key cosmosTypes.ChainIDHeightAndTxHashKey) {
	undelegateSuccessStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyUndelegateSuccessStore)
	storeKey := k.cdc.MustMarshal(&key)
	undelegateSuccessStore.Delete(storeKey)
}

func (k Keeper) ProcessAllUndelegateSuccess(ctx sdk.Context) error {
	list, err := k.getAllUndelegateSuccessDetails(ctx)
	if err != nil {
		return err
	}
	epochNumber := k.getLeastEpochNumberWithWithdrawStatusFalse(ctx)
	if epochNumber == int64(math.MaxInt64) {
		return cosmosTypes.ErrInvalidEpochNumber
	}
	for _, element := range list {
		if element.ValueUndelegateSuccessStore.Ratio > cosmosTypes.MinimumRatioForMajority {
			k.setEpochNumberAndUndelegateDetailsOfIndividualValidator(ctx, element.ValueUndelegateSuccessStore.ValidatorAddress, epochNumber, element.ValueUndelegateSuccessStore.Amount)
		}

		if element.ValueUndelegateSuccessStore.ActiveBlockHeight <= ctx.BlockHeight() {
			k.deleteUndelegateSuccessDetails(ctx, element.ChainIDHeightAndTxHashKey)
		}
	}

	flagForWithdrawSuccess := k.getEpochNumberAndUndelegateDetailsOfValidators(ctx, epochNumber)
	if flagForWithdrawSuccess {
		err = k.emitSendTransactionForAllWithdrawals(ctx, epochNumber)
		if err != nil {
			return err
		}
	}
	return nil
}
