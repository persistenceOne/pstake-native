package ratesync

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/v2/x/ratesync/keeper"
	"github.com/persistenceOne/pstake-native/v2/x/ratesync/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the chain
	for _, elem := range genState.HostChains {
		k.SetHostChain(ctx, elem)
	}
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.HostChains = k.GetAllHostChain(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
