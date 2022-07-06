package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

// InitGenesis new genesis for the cosmos module
func InitGenesis(ctx sdk.Context, keeper Keeper, data *cosmosTypes.GenesisState) {
	keeper.SetParams(ctx, data.Params)
	keeper.SetProposalID(ctx, 1)
	keeper.setID(ctx, 0, []byte(cosmosTypes.KeyLastTXPoolID))
}
