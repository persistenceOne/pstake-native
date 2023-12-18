package ratesync_test

import (
	"github.com/persistenceOne/pstake-native/v2/app/helpers"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/pstake-native/v2/x/ratesync"
	"github.com/persistenceOne/pstake-native/v2/x/ratesync/types"
)

func TestGenesis(t *testing.T) {
	_, pStakeApp, ctx := helpers.CreateTestApp(t)

	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		HostChains: []types.HostChain{
			{
				ID: 0,
			},
			{
				ID: 1,
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k := pStakeApp.RatesyncKeeper
	ratesync.InitGenesis(ctx, *k, genesisState)
	got := ratesync.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	require.Equal(t, genesisState.Params, got.Params)

	require.ElementsMatch(t, genesisState.HostChains, got.HostChains)
	// this line is used by starport scaffolding # genesis/test/assert
}
