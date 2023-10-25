package stkosmo

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

const (
	OsmosisTestnetChainID   = "osmo-test-5"
	OsmosisTestnetMintDenom = "stk/uosmo"
)

// MigrateStore performs in-place store migrations from v2.3.0 to remove stkOSMO from test-core-2.
// The migration includes:
//
// - Burn all minted stkOSMO.
// - Remove all the related store objects.
// - Migrate stkOSMO host chain to remove it from the store.
func MigrateStore(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec, k types.BankKeeper) error {

	// Burn all the coins present in the module account.
	if err := k.BurnCoins(
		ctx,
		types.ModuleName,
		sdk.NewCoins(k.GetBalance(ctx, authtypes.NewModuleAddress(types.ModuleName), OsmosisTestnetMintDenom)),
	); err != nil {
		return err
	}

	// Burn all coins on other addresses.
	k.IterateAllBalances(
		ctx,
		func(address sdk.AccAddress, coin sdk.Coin) (stop bool) {

			if coin.Denom == OsmosisTestnetMintDenom {
				// Send the whole address balance to the module account.
				if err := k.SendCoinsFromAccountToModule(ctx, address, types.ModuleName, sdk.NewCoins(coin)); err != nil {
					return false
				}

				// Burn the coins.
				if err := k.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(coin)); err != nil {
					return false
				}
			}

			return false
		},
	)

	// Remove all the object stores related to the chain.
	deleteDeposits(ctx, storeKey, cdc)
	deleteLSMDeposits(ctx, storeKey, cdc)
	deleteUnbondings(ctx, storeKey, cdc)
	deleteUserUnbondings(ctx, storeKey, cdc)
	deleteValidatorUnbondings(ctx, storeKey, cdc)

	// Remove the host chain from the store.
	store := prefix.NewStore(ctx.KVStore(storeKey), types.HostChainKey)
	store.Delete([]byte(OsmosisTestnetChainID))

	return nil
}

func deleteDeposits(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) {
	store := prefix.NewStore(ctx.KVStore(storeKey), types.DepositKey)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		deposit := types.Deposit{}
		cdc.MustUnmarshal(iterator.Value(), &deposit)

		if deposit.ChainId == OsmosisTestnetChainID {
			store.Delete(types.GetDepositStoreKey(deposit.ChainId, deposit.Epoch))
		}
	}
}

func deleteLSMDeposits(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) {
	store := prefix.NewStore(ctx.KVStore(storeKey), types.LSMDepositKey)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		deposit := types.LSMDeposit{}
		cdc.MustUnmarshal(iterator.Value(), &deposit)

		if deposit.ChainId == OsmosisTestnetChainID {
			store.Delete(types.GetLSMDepositStoreKey(deposit.ChainId, deposit.DelegatorAddress, deposit.Denom))
		}
	}
}

func deleteUnbondings(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) {
	store := prefix.NewStore(ctx.KVStore(storeKey), types.UnbondingKey)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		ub := types.Unbonding{}
		cdc.MustUnmarshal(iterator.Value(), &ub)

		if ub.ChainId == OsmosisTestnetChainID {
			store.Delete(types.GetUnbondingStoreKey(ub.ChainId, ub.EpochNumber))
		}
	}
}

func deleteUserUnbondings(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) {
	store := prefix.NewStore(ctx.KVStore(storeKey), types.UserUnbondingKey)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		ub := types.UserUnbonding{}
		cdc.MustUnmarshal(iterator.Value(), &ub)

		if ub.ChainId == OsmosisTestnetChainID {
			store.Delete(types.GetUserUnbondingStoreKey(ub.ChainId, ub.Address, ub.EpochNumber))
		}
	}
}

func deleteValidatorUnbondings(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) {
	store := prefix.NewStore(ctx.KVStore(storeKey), types.ValidatorUnbondingKey)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		ub := types.ValidatorUnbonding{}
		cdc.MustUnmarshal(iterator.Value(), &ub)

		if ub.ChainId == OsmosisTestnetChainID {
			store.Delete(types.GetValidatorUnbondingStoreKey(ub.ChainId, ub.ValidatorAddress, ub.EpochNumber))
		}
	}
}
