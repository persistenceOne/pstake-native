package keeper

import (
	"fmt"

	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkTx "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/cosmos/cosmos-sdk/x/bank/types"

	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

// addToWithdrawPool adds details to withdraw pool for ubonding epoch
func (k Keeper) addToWithdrawPool(ctx sdk.Context, asset cosmosTypes.MsgWithdrawStkAsset) error {
	withdrawStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyWithdrawStore)

	//get module params
	moduleParams := k.GetParams(ctx)

	//get both epochs info
	undelegateEpochInfo := k.epochsKeeper.GetEpochInfo(ctx, moduleParams.UndelegateEpochIdentifier)
	stakingEpochInfo := k.epochsKeeper.GetEpochInfo(ctx, moduleParams.StakingEpochIdentifier)

	//calculate time just 2*stakingEpochDuration before current epoch ends
	diffTime := undelegateEpochInfo.CurrentEpochStartTime.Add(undelegateEpochInfo.Duration - 2*stakingEpochInfo.Duration)

	//move withdraw transaction to next undelegating epoch if current block time is after the diffTime
	currentEpoch := undelegateEpochInfo.CurrentEpoch
	if ctx.BlockTime().Before(diffTime) {
		currentEpoch++
	}

	//if store has the key then append new withdrawals to the existing array, else make a new key value pair
	key := cosmosTypes.Int64Bytes(currentEpoch)
	if withdrawStore.Has(key) {
		bz := withdrawStore.Get(key)
		if bz == nil {
			return fmt.Errorf("withdraw store has key but nothing assigned to it")
		}
		var withdrawStoreValue cosmosTypes.WithdrawStoreValue
		err := k.cdc.Unmarshal(bz, &withdrawStoreValue)
		if err != nil {
			return err
		}
		withdrawStoreValue.WithdrawDetails = append(withdrawStoreValue.WithdrawDetails, asset)
		withdrawStoreValue.UnbondEmitFlag = append(withdrawStoreValue.UnbondEmitFlag, false)

		bz1, err := k.cdc.Marshal(&withdrawStoreValue)
		if err != nil {
			return err
		}
		withdrawStore.Set(key, bz1)
	} else {
		withdrawDetails := cosmosTypes.NewWithdrawStoreValue(asset)
		bz, err := k.cdc.Marshal(&withdrawDetails)
		if err != nil {
			return err
		}
		withdrawStore.Set(key, bz)
	}
	return nil
}

// fetchWithdrawTxnsWithCurrentEpochInfo Gets withdraw transaction mapped to current epoch number
func (k Keeper) fetchWithdrawTxnsWithCurrentEpochInfo(ctx sdk.Context, currentEpoch int64) (withdrawStoreValue cosmosTypes.WithdrawStoreValue) {
	withdrawStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyWithdrawStore)
	if !withdrawStore.Has(cosmosTypes.Int64Bytes(currentEpoch)) {
		return cosmosTypes.WithdrawStoreValue{WithdrawDetails: []cosmosTypes.MsgWithdrawStkAsset{{Amount: sdk.NewInt64Coin(k.GetParams(ctx).MintDenom, 0)}}}
	}
	k.cdc.MustUnmarshal(withdrawStore.Get(cosmosTypes.Int64Bytes(currentEpoch)), &withdrawStoreValue)
	return withdrawStoreValue
}

// deleteWithdrawTxnWithCurrentEpochInfo Remove the details mapped to the current epoch number
func (k Keeper) deleteWithdrawTxnWithCurrentEpochInfo(ctx sdk.Context, currentEpoch int64) {
	withdrawStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyWithdrawStore)
	withdrawStore.Delete(cosmosTypes.Int64Bytes(currentEpoch))
}

// totalAmountToBeUnbonded Get the total amount that is to be unbonded
func (k Keeper) totalAmountToBeUnbonded(value cosmosTypes.WithdrawStoreValue, denom string) sdk.Coin {
	amount := sdk.NewInt64Coin(denom, 0)
	for _, element := range value.WithdrawDetails {
		amount = amount.Add(sdk.NewCoin(denom, element.Amount.Amount))
	}
	return amount
}

// generateSendTransactionForAllWithdrawals Generates send transaction for the withdrawals and add it to the outgoing pool with the given txID
func (k Keeper) generateSendTransactionForAllWithdrawals(ctx sdk.Context, epochNumber int64, cValue sdk.Dec) error {
	withdrawStoreValue := k.fetchWithdrawTxnsWithCurrentEpochInfo(ctx, epochNumber)
	params := k.GetParams(ctx)
	bondDenom, err := params.GetBondDenomOf(cosmosTypes.DefaultStakingDenom)
	if err != nil {
		return err
	}
	chunkSlice := ChunkWithdrawSlice(withdrawStoreValue.WithdrawDetails, params.ChunkSize)
	for _, chunk := range chunkSlice {
		nextID := k.autoIncrementID(ctx, []byte(cosmosTypes.KeyLastTXPoolID))
		var sendMsgsAny []*codecTypes.Any
		for _, element := range chunk {
			withdrawalAmount, _ := sdk.NewDecCoinFromDec(bondDenom, element.Amount.Amount.ToDec().Mul(sdk.NewDec(1).Quo(cValue))).TruncateDecimal()
			msg := types.MsgSend{
				FromAddress: params.CustodialAddress,
				ToAddress:   element.ToAddress,
				Amount:      sdk.NewCoins(withdrawalAmount),
			}
			anyMsg, err := codecTypes.NewAnyWithValue(&msg)
			if err != nil {
				panic(err)
			}
			sendMsgsAny = append(sendMsgsAny, anyMsg)
		}

		cosmosAddrr, err := cosmosTypes.Bech32ifyAddressBytes(cosmosTypes.Bech32PrefixAccAddr, k.GetCurrentAddress(ctx))
		if err != nil {
			return err
		}
		execMsg := authz.MsgExec{
			Grantee: cosmosAddrr,
			Msgs:    sendMsgsAny,
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
						GasLimit: cosmosTypes.MinGasFee,
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

		//Once event is emitted, store it in KV store for orchestrators to query transactions and sign them
		k.SetNewTxnInOutgoingPool(ctx, nextID, tx)

		k.SetNewInTransactionQueue(ctx, nextID)
	}

	// delete from all the stores
	k.deleteEpochWithdrawSuccessStore(ctx, epochNumber)
	k.deleteWithdrawTxnWithCurrentEpochInfo(ctx, epochNumber)
	k.deleteEpochNumberAndUndelegateDetailsOfValidators(ctx, epochNumber)
	return nil
}

// ChunkWithdrawSlice divides 1D slice of MsgWithdrawStkAsset into chunks of given size and
// returns it by putting it in a 2D slice
func ChunkWithdrawSlice(slice []cosmosTypes.MsgWithdrawStkAsset, chunkSize int64) (chunks [][]cosmosTypes.MsgWithdrawStkAsset) {
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
