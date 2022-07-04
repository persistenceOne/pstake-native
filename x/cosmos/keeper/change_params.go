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

// setCustodialAddress sets custodial address in params
func (k Keeper) setCustodialAddress(ctx sdk.Context, address string) {
	k.paramSpace.Set(ctx, cosmosTypes.KeyCustodialAddress, address)
}

// setCosmosChainID sets cosmos chain ID in params
func (k Keeper) setCosmosChainID(ctx sdk.Context, chainID string) {
	k.paramSpace.Set(
		ctx, cosmosTypes.KeyCosmosProposalParams,
		cosmosTypes.CosmosChainProposalParams{
			ChainID:              chainID,
			ReduceVotingPeriodBy: cosmosTypes.DefaultPeriod,
		},
	)
}
