package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/pstake-native/x/cosmos/types"
)

var _ types.QueryServer = Keeper{}

// QueryParams queries all the params in genesis
func (k Keeper) QueryParams(context context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(context)
	params := k.GetParams(ctx)
	return &types.QueryParamsResponse{Params: params}, nil
}

// QueryTxByID Query txns by ID for orchestrators to sign
func (k Keeper) QueryTxByID(context context.Context, req *types.QueryOutgoingTxByIDRequest) (*types.QueryOutgoingTxByIDResponse, error) {
	ctx := sdk.UnwrapSDKContext(context)
	cosmosTxDetails, err := k.getTxnFromOutgoingPoolByID(ctx, req.TxID)
	if err != nil {
		return nil, err
	}
	return &cosmosTxDetails, nil
}
