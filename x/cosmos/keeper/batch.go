package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/tx"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

const OutgoingTxBatchSize = 100

func (k Keeper) BuildOutgoingTxBatch(ctx sdk.Context, batch_size int) (tx.TxBody, error) {
	//TODO
	return tx.TxBody{}, nil
}

// addUnbatchedTx creates a new transaction in the pool
// WARNING: Do not make this function public
func (k Keeper) addUnbatchedTX(ctx sdk.Context, val *cosmosTypes.OutgoingTransferTx) error {
	store := ctx.KVStore(k.storeKey)
	idxKey := []byte(cosmosTypes.GetOutgoingTxPoolKey(val.Fees, val.Id))
	if store.Has(idxKey) {
		return sdkerrors.Wrap(cosmosTypes.ErrDuplicate, "transaction already in pool")
	}

	bz, err := k.cdc.Marshal(val)
	if err != nil {
		return err
	}

	store.Set(idxKey, bz)
	return err
}

func (k Keeper) addToMintingPoolTx(ctx sdk.Context, destinationAddress sdk.AccAddress, orchestratorAddress sdk.AccAddress, amount sdk.Coins) error {
	store := ctx.KVStore(k.storeKey)
	mintingPoolStore := prefix.NewStore(store, []byte(cosmosTypes.MintingPoolStoreKey))
	key := []byte(cosmosTypes.GetDestinationAddressAndAmountKey(destinationAddress, amount))
	if mintingPoolStore.Has(key) {
		var txnDetails cosmosTypes.IncomingMintTx
		bz := mintingPoolStore.Get(key)
		err := txnDetails.Unmarshal(bz)

		found := txnDetails.Find(orchestratorAddress.String())
		if !found {
			txnDetails.AddAndIncrement(orchestratorAddress.String())
		}

		bz, err = txnDetails.Marshal()
		if err != nil {
			return err
		}
		mintingPoolStore.Set(key, bz)
	} else {
		txnDetails := cosmosTypes.NewIncomingMintTx(orchestratorAddress, 1)
		bz, _ := txnDetails.Marshal()
		mintingPoolStore.Set(key, bz)
	}
	return nil
}

func (k Keeper) FetchFromMintPoolTx(ctx sdk.Context, keyAndValueForMinting []KeyAndValueForMinting) []KeyAndValueForMinting {
	store := ctx.KVStore(k.storeKey)
	mintingPoolStore := prefix.NewStore(store, []byte(cosmosTypes.MintingPoolStoreKey))
	totalCount := float64(k.getTotalValidatorOrchestratorCount(ctx))
	for i, _ := range keyAndValueForMinting {
		destinationAddress, err := sdk.AccAddressFromBech32(keyAndValueForMinting[i].Value.DestinationAddress)
		if err != nil {
			panic("Error in parsing destination address")
		}

		key := []byte(cosmosTypes.GetDestinationAddressAndAmountKey(destinationAddress, keyAndValueForMinting[i].Value.Amount))
		bz := mintingPoolStore.Get(key)

		var txnDetails cosmosTypes.IncomingMintTx
		err = txnDetails.Unmarshal(bz)
		if err != nil {
			panic("Error in unmarshalling txn Details")
		}

		sizeOfOrchAddress := float64(len(txnDetails.OrchAddresses))
		ratio := sizeOfOrchAddress / totalCount

		if ratio > 0.66 {
			keyAndValueForMinting[i].Ratio = ratio
		}
	}
	return keyAndValueForMinting
}

func (k Keeper) deleteFromMintPoolTx(ctx sdk.Context, destinationAddress sdk.AccAddress, amount sdk.Coins) {
	store := ctx.KVStore(k.storeKey)
	mintingPoolStore := prefix.NewStore(store, []byte(cosmosTypes.MintingPoolStoreKey))
	mintingPoolStore.Delete([]byte(cosmosTypes.GetDestinationAddressAndAmountKey(destinationAddress, amount)))
}
