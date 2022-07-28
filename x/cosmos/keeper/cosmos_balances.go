package keeper

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

func (k Keeper) setCosmosBalance(ctx sdk.Context, balance sdk.Coins) {
	store := ctx.KVStore(k.storeKey)

	bz, err := balance.MarshalJSON()
	if err != nil {
		panic(err)
	}

	store.Set(cosmosTypes.KeyCosmosBalances, bz)
}

func (k Keeper) getCosmosBalances(ctx sdk.Context) (balance sdk.Coins) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(cosmosTypes.KeyCosmosBalances)
	err := json.Unmarshal(bz, &balance)
	if err != nil {
		panic(err)
	}

	return balance
}
