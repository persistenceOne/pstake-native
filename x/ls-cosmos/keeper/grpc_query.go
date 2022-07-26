package keeper

import (
	"github.com/persistenceOne/pstake-native/x/ls-cosmos/types"
)

var _ types.QueryServer = Keeper{}
