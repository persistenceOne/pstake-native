package types_test

import (
	"github.com/persistenceOne/pstake-native/x/cosmos/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTypes(t *testing.T) {
	num := uint64(1)
	bz := types.UInt64Bytes(num)
	num1 := types.UInt64FromBytes(bz)
	require.Equal(t, num, num1)

	num2 := int64(1)
	bz1 := types.Int64Bytes(num2)
	num3 := types.Int64FromBytes(bz1)
	require.Equal(t, num2, num3)
}
