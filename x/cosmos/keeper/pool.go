package keeper

import (
	"encoding/binary"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
	"strconv"
)

func (k Keeper) AddToOutgoingPool(ctx sdk.Context, sender sdk.AccAddress, msg *types.Any) (uint64, error) {
	nextID := k.autoIncrementID(ctx, []byte(cosmosTypes.KeyLastTXPoolID))
	outgoing := &cosmosTypes.OutgoingTransferTx{
		Id:      nextID,
		Message: msg,
	}

	err := k.addUnbatchedTX(ctx, outgoing)
	if err != nil {
		panic(err)
	}

	poolEvent := sdk.NewEvent(
		cosmosTypes.EventTypeAddToOutgoingPool,
		sdk.NewAttribute(sdk.AttributeKeyModule, cosmosTypes.ModuleName),
		sdk.NewAttribute(cosmosTypes.AttributeSender, sender.String()),
		sdk.NewAttribute(cosmosTypes.AttributeKeyOutgoingTXID, strconv.Itoa(int(nextID))),
		sdk.NewAttribute(cosmosTypes.AttributeKeyNonce, fmt.Sprint(nextID)),
	)
	ctx.EventManager().EmitEvent(poolEvent)

	return nextID, nil
}

// a specialized function used for iterating store counters, handling
// returning, initializing and incrementing all at once. This is particularly
// used for the transaction pool and batch pool where each batch or transaction is
// assigned a unique ID.
func (k Keeper) autoIncrementID(ctx sdk.Context, idKey []byte) uint64 {
	id := k.getID(ctx, idKey)
	id += 1
	k.setID(ctx, id, idKey)
	return id
}

// gets a generic uint64 counter from the store, initializing to 1 if no value exists
func (k Keeper) getID(ctx sdk.Context, idKey []byte) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(idKey)
	id := binary.BigEndian.Uint64(bz)
	return id
}

// sets a generic uint64 counter in the store
func (k Keeper) setID(ctx sdk.Context, id uint64, idKey []byte) {
	store := ctx.KVStore(k.storeKey)
	bz := sdk.Uint64ToBigEndian(id)
	store.Set(idKey, bz)
}
