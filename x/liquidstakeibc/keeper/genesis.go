package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

// InitGenesis initializes the liquidstakeibc module's state from a given genesis state.
func (k Keeper) InitGenesis(ctx sdk.Context, genState *types.GenesisState) {
	k.SetParams(ctx, genState.Params)

}

// ExportGenesis returns the liquidstakeibc module's genesis state.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {

	return types.NewGenesisState(
		k.GetParams(ctx),
	)
}
