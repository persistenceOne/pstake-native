package keeper_test

import (
	"github.com/persistenceOne/pstake-native/x/cosmos/keeper"
	"testing"
)

func Keeper_test(t *testing.T) {
	app := newTestApp()
	app.CosmosKeeper = keeper.NewKeeper()
}
