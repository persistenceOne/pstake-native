package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

var _ types.QueryServer = Keeper{}

// HostChainParams returns the stored host chain params set through proposal.
func (k Keeper) HostChainParams(c context.Context, in *types.QueryHostChainParamsRequest) (*types.QueryHostChainParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	hostChainParams := k.GetHostChainParams(ctx)

	return &types.QueryHostChainParamsResponse{
		HostChainParams: hostChainParams,
	}, nil
}

func (k Keeper) DelegationState(ctx context.Context, request *types.QueryDelegationStateRequest) (*types.QueryDelegationStateResponse, error) {
	sdkctx := sdk.UnwrapSDKContext(ctx)
	delegationState := k.GetDelegationState(sdkctx)

	return &types.QueryDelegationStateResponse{
		DelegationState: delegationState,
	}, nil

}
