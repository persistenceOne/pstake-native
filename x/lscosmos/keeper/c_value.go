package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

func (k Keeper) GetCValue(ctx sdk.Context) sdk.Dec {
	//TODO: C-value logic to be implemented
	return sdk.NewDec(1)
}
