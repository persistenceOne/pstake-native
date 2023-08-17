package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	liquidstakeibctypes "github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

func (k *Keeper) Migrate(ctx sdk.Context) error {

	for _, hc := range k.GetAllHostChains(ctx) {
		// set the default LSM flag value to avoid null pointer references
		hc.Flags = &liquidstakeibctypes.HostChainFlags{Lsm: false}
		k.SetHostChain(ctx, hc)
	}

	return nil
}
