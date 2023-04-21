package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k *Keeper) RegisterICAAccount(ctx sdk.Context, connectionId, owner string) error {
	return k.icaControllerKeeper.RegisterInterchainAccount(
		ctx,
		connectionId,
		owner,
		"",
	)
}
