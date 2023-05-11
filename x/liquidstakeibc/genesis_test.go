package liquidstakeibc_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/pstake-native/v2/app/helpers"
	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc"
	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

func TestGenesis(t *testing.T) {
	genesisState := &types.GenesisState{
		Params:     types.DefaultParams(),
		HostChains: make([]*types.HostChain, 0),
		Deposits:   make([]*types.Deposit, 0),
	}

	_, pStakeApp, ctx := helpers.CreateTestApp(t)
	k := pStakeApp.LiquidStakeIBCKeeper
	liquidstakeibc.InitGenesis(ctx, k, genesisState)

	got := liquidstakeibc.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	require.Equal(t, genesisState.Params, got.Params)
	require.Equal(t, genesisState.HostChains, got.HostChains)
	require.Equal(t, genesisState.Deposits, got.Deposits)
}
