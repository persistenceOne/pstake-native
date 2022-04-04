package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

// InitGenesis new cosmos genesis
func InitGenesis(ctx sdk.Context, keeper Keeper, data *cosmosTypes.GenesisState) {
	keeper.SetParams(ctx, data.Params)
	keeper.SetProposalID(ctx, 1)
	keeper.setID(ctx, 0, []byte(cosmosTypes.KeyLastTXPoolID))
	keeper.setCosmosValidatorParams(ctx, nil)
	keeper.setTotalDelegatedAmountTillDate(ctx, sdk.Coin{})
	//keeper.SetVotingParams(ctx, data.Params.CosmosProposalParams)
	//TODO add remaining
}
