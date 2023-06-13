package keeper

import (
	"context"

	"github.com/persistenceOne/pstake-native/v2/x/lscosmos/types"
)

// Params queries the module params
func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	return nil, types.ErrDeprecated
}
