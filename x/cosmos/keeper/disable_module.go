package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

// disableModule disables module by setting param to true
func (k Keeper) disableModule(ctx sdk.Context) {
	k.paramSpace.Set(ctx, cosmosTypes.KeyModuleEnabled, false)
}

// enableModule enables module by setting param to true
func (k Keeper) enableModule(ctx sdk.Context) {
	k.paramSpace.Set(ctx, cosmosTypes.KeyModuleEnabled, true)
}
