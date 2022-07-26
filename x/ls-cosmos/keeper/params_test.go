package keeper_test

import (
	"testing"

	testkeeper "github.com/persistenceOne/pstake-native/testutil/keeper"
	"github.com/persistenceOne/pstake-native/x/ls-cosmos/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.LscosmosKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
