package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/persistenceOne/pstake-native/x/lscosmos/types"
    "github.com/persistenceOne/pstake-native/x/lscosmos/keeper"
    keepertest "github.com/persistenceOne/pstake-native/testutil/keeper"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.LscosmosKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
