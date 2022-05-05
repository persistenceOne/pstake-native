package keeper

import (
	"fmt"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkTx "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/x/authz"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

// Generate an event for delegating on cosmos chain once staking epoch is called
func (k Keeper) generateDelegateOutgoingEvent(ctx sdk.Context, validatorSet []ValAddressAndAmountForStakingAndUndelegating, epochNumber int64) error {
	params := k.GetParams(ctx)

	//create chunks for delegation on cosmos chain
	chunkSlice := ChunkStakeAndUnStakeSlice(validatorSet, params.ChunkSize)

	for _, chunk := range chunkSlice {
		nextID := k.autoIncrementID(ctx, []byte(cosmosTypes.KeyLastTXPoolID))

		var delegateMsgsAny []*codecTypes.Any
		for _, element := range chunk {
			msg := stakingTypes.MsgDelegate{
				DelegatorAddress: params.CustodialAddress,
				ValidatorAddress: element.validator.String(),
				Amount:           element.amount,
			}
			anyMsg, err := codecTypes.NewAnyWithValue(&msg)
			if err != nil {
				return err
			}
			delegateMsgsAny = append(delegateMsgsAny, anyMsg)
		}

		execMsg := authz.MsgExec{
			Grantee: params.CustodialAddress,
			Msgs:    delegateMsgsAny,
		}

		execMsgAny, err := codecTypes.NewAnyWithValue(&execMsg)
		if err != nil {
			return err
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

		// set acknowledgment flag true for future reference (not any yet)

		//ctx.EventManager().EmitEvent(
		//	sdk.NewEvent(
		//		cosmosTypes.EventTypeOutgoing,
		//		sdk.NewAttribute(cosmosTypes.AttributeKeyOutgoingTXID, fmt.Sprint(nextID)),
		//	),
		//)

		err = k.setInEpochPoolForMinting(ctx, epochNumber, nextID, false)
		if err != nil {
			return err
		}
		//Once event is emitted, store it in KV store for orchestrators to query transactions and sign them
		k.setNewTxnInOutgoingPool(ctx, nextID, tx)

		k.setNewInTransactionQueue(ctx, nextID)
	}

	return nil
}

//______________________________________________________________________________________________________________________

type EpochNumberAndDetailsForMinting struct {
	epochNumber       int64
	mintingEpochValue cosmosTypes.MintingEpochValue
}

func (k Keeper) setInEpochPoolForMinting(ctx sdk.Context, epochNumber int64, nextID uint64, status bool) error {
	mintingEpochStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyMintingEpochStore)
	key := cosmosTypes.Int64Bytes(epochNumber)
	if mintingEpochStore.Has(key) {
		var mintingEpochStoreValue cosmosTypes.MintingEpochValue
		err := k.cdc.Unmarshal(mintingEpochStore.Get(key), &mintingEpochStoreValue)
		if err != nil {
			return err
		}
		mintingEpochStoreValue.TxIDAndStatus = append(mintingEpochStoreValue.TxIDAndStatus, cosmosTypes.MintingEpochValueMember{TxID: nextID, Status: status})
		bz, err := k.cdc.Marshal(&mintingEpochStoreValue)
		if err != nil {
			return err
		}
		mintingEpochStore.Set(key, bz)
		return nil
	}
	mintingEpochStoreValue := cosmosTypes.NewMintingEpochValue(cosmosTypes.MintingEpochValueMember{TxID: nextID, Status: status})
	bz, err := k.cdc.Marshal(&mintingEpochStoreValue)
	if err != nil {
		return err
	}
	mintingEpochStore.Set(key, bz)
	return nil
}

func (k Keeper) setListInEpochPoolForMinting(ctx sdk.Context, epochNumber int64, mintingEpochStoreValue cosmosTypes.MintingEpochValue) error {
	mintingEpochStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyMintingEpochStore)
	key := cosmosTypes.Int64Bytes(epochNumber)
	bz, err := k.cdc.Marshal(&mintingEpochStoreValue)
	if err != nil {
		return err
	}
	mintingEpochStore.Set(key, bz)
	return nil
}

func (k Keeper) fetchInEpochPoolForMinting(ctx sdk.Context, epochNumber int64) (cosmosTypes.MintingEpochValue, error) {
	mintingEpochStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyMintingEpochStore)
	key := cosmosTypes.Int64Bytes(epochNumber)
	if mintingEpochStore.Has(key) {
		var mintingEpochStoreValue cosmosTypes.MintingEpochValue
		err := k.cdc.Unmarshal(mintingEpochStore.Get(key), &mintingEpochStoreValue)
		if err != nil {
			return cosmosTypes.MintingEpochValue{}, err
		}
		return mintingEpochStoreValue, nil
	}
	return cosmosTypes.MintingEpochValue{}, fmt.Errorf("could not locate %d in pool", epochNumber)
}

func (k Keeper) fetchAllInEpochPoolForMinting(ctx sdk.Context) (list []EpochNumberAndDetailsForMinting, err error) {
	mintingEpochStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyMintingEpochStore)
	iterator := mintingEpochStore.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		epochNumber := cosmosTypes.Int64FromBytes(iterator.Key())
		var mintingEpochStoreValue cosmosTypes.MintingEpochValue
		if err = k.cdc.Unmarshal(iterator.Value(), &mintingEpochStoreValue); err != nil {
			return list, err
		}
		list = append(list, EpochNumberAndDetailsForMinting{epochNumber: epochNumber, mintingEpochValue: mintingEpochStoreValue})
	}
	return list, nil
}

func (k Keeper) deleteInEpochPoolForMinting(ctx sdk.Context, epochNumber int64) {
	mintingEpochStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.KeyMintingEpochStore)
	key := cosmosTypes.Int64Bytes(epochNumber)
	mintingEpochStore.Delete(key)
}

//______________________________________________________________________________________________________________________
func (k Keeper) setTotalDelegatedAmountTillDate(ctx sdk.Context, addToTotal sdk.Coin) {
	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.Marshal(&addToTotal)
	if err != nil {
		panic(err)
	}
	store.Set([]byte(cosmosTypes.KeyTotalDelegationTillDate), bz)
}

func (k Keeper) getTotalDelegatedAmountTillDate(ctx sdk.Context) sdk.Coin {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(cosmosTypes.KeyTotalDelegationTillDate))
	var amount sdk.Coin
	err := k.cdc.Unmarshal(bz, &amount)
	if err != nil {
		panic(err)
	}
	return amount
}

//______________________________________________________________________________________________________________________
func (k Keeper) processStakingSuccessTxns(ctx sdk.Context, txID uint64) error {
	epochNumberAndTxIDStatusList, err := k.fetchAllInEpochPoolForMinting(ctx)
	if err != nil {
		return err
	}
	for _, en := range epochNumberAndTxIDStatusList {
		for i := range en.mintingEpochValue.TxIDAndStatus {
			if en.mintingEpochValue.TxIDAndStatus[i].TxID == txID {
				en.mintingEpochValue.TxIDAndStatus[i].Status = true
				err = k.setListInEpochPoolForMinting(ctx, en.epochNumber, en.mintingEpochValue)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (k Keeper) emitStakingTxnForClaimedRewards(ctx sdk.Context, msgs []sdk.Msg) {
	//totalAmountInClaimMsgs := sdk.NewInt64Coin(k.GetParams(ctx).BondDenom, 0)
	//TODO : Ask which impl to go forwards with txn response for claimRewards and minting rewards for devs and validators
}
