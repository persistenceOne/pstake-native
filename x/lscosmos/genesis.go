package lscosmos

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	host "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	"github.com/persistenceOne/pstake-native/x/lscosmos/keeper"
	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init

	k.SetParams(ctx, genState.Params)
	k.SetModuleState(ctx, genState.ModuleEnabled)
	k.SetCosmosIBCParams(ctx, genState.CosmosIBCParams)
	if !genState.CosmosIBCParams.IsEmpty() {
		err := k.NewCapability(ctx, host.ChannelCapabilityPath(genState.CosmosIBCParams.TokenTransferPort, genState.CosmosIBCParams.TokenTransferChannel))
		if err != nil {
			panic(err)
		}
	}
	k.SetAllowListedValidators(ctx, genState.AllowListedValidators)
	k.GetDepositAccount(ctx)
	k.GetDelegationAccount(ctx)
	k.GetRewardAccount(ctx)
	k.GetUndelegationAccount(ctx)

}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.ModuleEnabled = k.GetModuleState(ctx)
	genesis.CosmosIBCParams = k.GetCosmosIBCParams(ctx)
	genesis.AllowListedValidators = k.GetAllowListedValidators(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
