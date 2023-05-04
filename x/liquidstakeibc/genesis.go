package liquidstakeibc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/keeper"
	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

// InitGenesis initializes the liquidstakeibc module's state from a given genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState *types.GenesisState) {
	k.SetParams(ctx, genState.Params)

	for _, hc := range genState.HostChains {
		k.SetHostChain(ctx, hc)
	}

	for _, deposit := range genState.Deposits {
		k.SetDeposit(ctx, deposit)
	}

	k.GetDepositModuleAccount(ctx)
}

// ExportGenesis returns the liquidstakeibc module's genesis state.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {

	return &types.GenesisState{
		Params:     k.GetParams(ctx),
		HostChains: k.GetAllHostChains(ctx),
		Deposits:   k.GetAllDeposits(ctx),
	}
}
