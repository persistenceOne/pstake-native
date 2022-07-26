package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/persistenceOne/pstake-native/testutil/keeper"
	"github.com/persistenceOne/pstake-native/x/icacosmos/keeper"
	"github.com/persistenceOne/pstake-native/x/icacosmos/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.IcacosmosKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
