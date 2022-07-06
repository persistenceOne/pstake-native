package types_test

import (
	"github.com/persistenceOne/pstake-native/app"
	"github.com/persistenceOne/pstake-native/x/cosmos/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParams_Validate(t *testing.T) {
	app.SetAddressPrefixes()
	params := types.DefaultParams()

	// default params have no error
	require.NoError(t, params.Validate())
}
