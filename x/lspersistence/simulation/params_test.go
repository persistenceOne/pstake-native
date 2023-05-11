package simulation_test

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/pstake-native/v2/x/lspersistence/simulation"
)

func TestParamChanges(t *testing.T) {
	s := rand.NewSource(1)
	r := rand.New(s)

	expected := []struct {
		composedKey string
		key         string
		simValue    string
		subspace    string
	}{
		{"lspersistence/WhitelistedValidators", "WhitelistedValidators", "[]", "lspersistence"},
		{"lspersistence/LiquidBondDenom", "LiquidBondDenom", "\"stk/uxprt\"", "lspersistence"},
		{"lspersistence/UnstakeFeeRate", "UnstakeFeeRate", "\"0.010000000000000000\"", "lspersistence"},
		{"lspersistence/MinLiquidStakingAmount", "MinLiquidStakingAmount", "\"9727887\"", "lspersistence"},
	}

	paramChanges := simulation.ParamChanges(r)
	require.Len(t, paramChanges, 4)

	for i, p := range paramChanges {
		require.Equal(t, expected[i].composedKey, p.ComposedKey())
		require.Equal(t, expected[i].key, p.Key())
		require.Equal(t, expected[i].simValue, p.SimValue()(r))
		require.Equal(t, expected[i].subspace, p.Subspace())
	}
}
