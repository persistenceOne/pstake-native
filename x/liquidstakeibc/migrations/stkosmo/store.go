package stkosmo

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

const (
	OsmosisTestnetChainID   = "osmo-test-5"
	OsmosisTestnetMintDenom = "stk/uosmo"
)

// MigrateStore performs in-place store migrations from v2.3.0 to remove stkOSMO from test-core-2.
// The migration includes:
//
// - Migrate stkOSMO host chain to remove it from the store.
// - Burn all minted stkOSMO.
func MigrateStore(ctx sdk.Context, storeKey storetypes.StoreKey, k types.BankKeeper) error {

	// Burn all the coins present in the module account.
	if err := k.BurnCoins(
		ctx,
		types.ModuleName,
		sdk.NewCoins(k.GetBalance(ctx, authtypes.NewModuleAddress(types.ModuleName), OsmosisTestnetMintDenom)),
	); err != nil {
		return err
	}

	// Get the specific coin bank store.
	// This store maps denomination to account balance for that denomination.
	denomPrefixStore := prefix.NewStore(
		ctx.KVStore(sdk.NewKVStoreKey(banktypes.StoreKey)),
		banktypes.CreateDenomAddressPrefix(OsmosisTestnetMintDenom),
	)

	// Loop through all the entries in the store.
	_, err := query.FilteredPaginate(
		denomPrefixStore,
		nil,
		func(key []byte, _ []byte, accumulate bool) (bool, error) {
			// Get the address for each entry.
			address, _, err := banktypes.AddressAndDenomFromBalancesStore(key)
			if err != nil {
				return false, err
			}

			// Get the balance for that address.
			balance := sdk.NewCoins(k.GetBalance(ctx, address, OsmosisTestnetMintDenom))

			// Send the whole address balance to the module account.
			if err = k.SendCoinsFromAccountToModule(ctx, address, types.ModuleName, balance); err != nil {
				return false, err
			}

			// Burn the coins.
			if err = k.BurnCoins(ctx, types.ModuleName, balance); err != nil {
				return false, err
			}

			return true, nil
		},
	)
	if err != nil {
		return err
	}

	// Remove the host chain from the store
	store := prefix.NewStore(ctx.KVStore(storeKey), types.HostChainKey)
	store.Delete([]byte(OsmosisTestnetChainID))

	return nil
}
