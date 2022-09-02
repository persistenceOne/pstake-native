package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/pstake-native/app/helpers"
	"github.com/persistenceOne/pstake-native/x/lscosmos/keeper"
	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

func (suite *IntegrationTestSuite) setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	_, app, ctx := helpers.CreateTestApp()
	k := app.LSCosmosKeeper
	return keeper.NewMsgServerImpl(k), sdk.WrapSDKContext(ctx)
}
