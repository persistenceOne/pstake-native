package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	liquidstakeibctypes "github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
	"github.com/persistenceOne/pstake-native/v2/x/ratesync/types"
)

func (k *Keeper) BeginBlock(ctx sdk.Context) {

	// perform BeginBlocker tasks for each chain
	for _, hc := range k.GetAllHostChain(ctx) {
		if !hc.IsActive() {
			// don't do anything on inactive chains
			continue
		}
		// attempt to recreate closed ICA channels
		k.DoRecreateICA(ctx, hc)

		// reset hc before going into next function, as it might have changed in earlier function
		// as we do not want to re-write and omit the last write.
	}

}

func (k *Keeper) DoRecreateICA(ctx sdk.Context, hc types.HostChain) {
	// return early if any of the accounts is currently being recreated
	if hc.ICAAccount.ChannelState == liquidstakeibctypes.ICAAccount_ICA_CHANNEL_CREATING {
		return
	}

	// if the channel is closed, and it is not being recreated, recreate it
	portID := types.MustICAPortIDFromOwner(hc.ICAAccount.Owner)
	_, isActive := k.icaControllerKeeper.GetOpenActiveChannel(ctx, hc.ConnectionID, portID)
	if !isActive {
		if err := k.icaControllerKeeper.RegisterInterchainAccount(ctx, hc.ConnectionID, portID, ""); err != nil {
			k.Logger(ctx).Error("error recreating %s ratesync ica: %w", hc.ChainID, err)
		} else {
			k.Logger(ctx).Info("Recreating ratesync ICA.", "chain", hc.ChainID)

			hc.ICAAccount.ChannelState = liquidstakeibctypes.ICAAccount_ICA_CHANNEL_CREATING
			k.SetHostChain(ctx, hc)
		}
	}
}
