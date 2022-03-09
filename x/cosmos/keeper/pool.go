package keeper

import (
	"encoding/binary"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
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

func (k Keeper) setMintAddressAndAmount(ctx sdk.Context, chainID string, blockHeight int64, txHash string, destinationAddress sdk.AccAddress, amount sdk.Coins) {
	store := ctx.KVStore(k.storeKey)
	mintAddressAndAmountStore := prefix.NewStore(store, []byte(cosmosTypes.AddressAndAmountKey))

	chainIDHeightAndTxHash := cosmosTypes.NewChainIDHeightAndTxHash(chainID, blockHeight, txHash)
	key, err := chainIDHeightAndTxHash.Marshal()
	if err != nil {
		panic("error in marshaling chainID, height and txHash")
	}

	addressAndAmount := cosmosTypes.NewAddressAndAmount(destinationAddress, amount)
	bz, err := addressAndAmount.Marshal()
	if err != nil {
		panic("error in marshaling address and amount")
	}
	mintAddressAndAmountStore.Set(key, bz)
}

type KeyAndValueForMinting struct {
	Key   cosmosTypes.ChainIDHeightAndTxHash
	Value cosmosTypes.AddressAndAmount
	Ratio float64
}

func (k Keeper) GetAllMintAddressAndAmount(ctx sdk.Context, list []KeyAndValueForMinting) ([]KeyAndValueForMinting, error) {
	store := ctx.KVStore(k.storeKey)
	mintAddressAndAmountStore := prefix.NewStore(store, []byte(cosmosTypes.AddressAndAmountKey))

	iterator := mintAddressAndAmountStore.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var chainIDHeightAndTxHash cosmosTypes.ChainIDHeightAndTxHash

		err := chainIDHeightAndTxHash.Unmarshal(iterator.Key())
		if err != nil {
			return nil, err
		}

		var addressAndAmount cosmosTypes.AddressAndAmount

		err = addressAndAmount.Unmarshal(iterator.Value())
		if err != nil {
			return nil, err
		}

		a := KeyAndValueForMinting{
			Key:   chainIDHeightAndTxHash,
			Value: addressAndAmount,
		}

		list = append(list, a)
	}
	return list, nil
}

func (k Keeper) deleteMintedAddressAndAmountKeys(ctx sdk.Context, keyHash cosmosTypes.ChainIDHeightAndTxHash) {
	store := ctx.KVStore(k.storeKey)
	mintAddressAndAmountStore := prefix.NewStore(store, []byte(cosmosTypes.AddressAndAmountKey))

	chainIDHeightAndTxHash := cosmosTypes.NewChainIDHeightAndTxHash(keyHash.ChainID, keyHash.BlockHeight, keyHash.TxHash)
	key, err := chainIDHeightAndTxHash.Marshal()
	if err != nil {
		panic("error in marshaling chainID, height and txHash")
	}

	mintAddressAndAmountStore.Delete(key)
}
