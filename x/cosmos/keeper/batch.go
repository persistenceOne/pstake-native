package keeper

import (
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
