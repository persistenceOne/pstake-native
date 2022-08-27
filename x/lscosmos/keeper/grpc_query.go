package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

var _ types.QueryServer = Keeper{}

// CosmosParams returns the stored cosoms IBC params set through proposal.
func (k Keeper) CosmosParams(c context.Context, in *types.QueryCosmosParamsRequest) (*types.QueryCosmosParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	cosmosParams := k.GetCosmosParams(ctx)

	return &types.QueryCosmosParamsResponse{
		CosmosParams: cosmosParams,
	}, nil
}

func (k Keeper) DelegationState(ctx context.Context, request *types.QueryDelegationStateRequest) (*types.QueryDelegationStateResponse, error) {
	sdkctx := sdk.UnwrapSDKContext(ctx)
	delegationState := k.GetDelegationState(sdkctx)

	return &types.QueryDelegationStateResponse{
		DelegationState: delegationState,
	}, nil

}
