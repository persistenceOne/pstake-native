package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

var _ types.QueryServer = Keeper{}

// CosmosIBCParams returns the stored cosoms IBC params set through proposal.
func (k Keeper) CosmosIBCParams(c context.Context, in *types.QueryCosmosIBCParamsRequest) (*types.QueryCosmosIBCParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	ibcParams := k.GetCosmosIBCParams(ctx)

	return &types.QueryCosmosIBCParamsResponse{
		CosmosIBCParams: ibcParams,
	}, nil
}

func (k Keeper) DelegationState(ctx context.Context, request *types.QueryDelegationStateRequest) (*types.QueryDelegationStateResponse, error) {
	sdkctx := sdk.UnwrapSDKContext(ctx)
	delegationState := k.GetDelegationState(sdkctx)

	return &types.QueryDelegationStateResponse{
		DelegationState: delegationState,
	}, nil

}
