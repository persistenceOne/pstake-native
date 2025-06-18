package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/v3/x/liquidstakeibc/types"
)

func (k *Keeper) SetRedelegationTx(ctx sdk.Context, redelegationTx *types.RedelegateTx) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.RedelegationTxKey)
	bytes := k.cdc.MustMarshal(redelegationTx)
	store.Set(types.GetRedelegationTxStoreKey(redelegationTx.ChainId, redelegationTx.IbcSequenceId), bytes)
}

func (k *Keeper) GetRedelegationTx(ctx sdk.Context, chainID, ibcSequenceID string) (*types.RedelegateTx, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.RedelegationTxKey)
	bz := store.Get(types.GetRedelegationTxStoreKey(chainID, ibcSequenceID))
	if bz == nil {
		return nil, false
	}

	var redelegationtx types.RedelegateTx
	k.cdc.MustUnmarshal(bz, &redelegationtx)
	return &redelegationtx, true
}

func (k *Keeper) GetAllRedelegationTx(ctx sdk.Context) []*types.RedelegateTx {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.RedelegationTxKey)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	txs := make([]*types.RedelegateTx, 0)
	for ; iterator.Valid(); iterator.Next() {
		redelegateTx := types.RedelegateTx{}
		k.cdc.MustUnmarshal(iterator.Value(), &redelegateTx)
		txs = append(txs, &redelegateTx)
	}

	return txs
}

func (k *Keeper) FilterRedelegationTx(
	ctx sdk.Context,
	filter func(d types.RedelegateTx) bool,
) []*types.RedelegateTx {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.RedelegationTxKey)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	redelegationTxs := make([]*types.RedelegateTx, 0)
	for ; iterator.Valid(); iterator.Next() {
		deposit := types.RedelegateTx{}
		k.cdc.MustUnmarshal(iterator.Value(), &deposit)
		if filter(deposit) {
			redelegationTxs = append(redelegationTxs, &deposit)
		}
	}

	return redelegationTxs
}

func (k *Keeper) DeleteRedelegationTx(ctx sdk.Context, chainID, ibcSequenceID string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.RedelegationTxKey)
	store.Delete(types.GetRedelegationTxStoreKey(chainID, ibcSequenceID))
}
