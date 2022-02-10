package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	mintTypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramsTypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/persistenceOne/pstake-native/x/cosmos/types"
)

type Keeper struct {
	storeKey   sdk.StoreKey
	paramSpace paramsTypes.Subspace
	bankKeeper types.BankKeeper
	mintKeeper types.MintKeeper
}

func NewKeeper(
	key sdk.StoreKey, paramSpace paramsTypes.Subspace,
	bankKeeper types.BankKeeper, mintKeeper types.MintKeeper,
) Keeper {

	return Keeper{
		storeKey:   key,
		paramSpace: paramSpace.WithKeyTable(types.ParamKeyTable()),
		bankKeeper: bankKeeper,
		mintKeeper: mintKeeper,
	}
}

//______________________________________________________________________

// GetParams returns the total set of parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the total set of parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

//______________________________________________________________________

// GetMintingParams returns the total set of cosmos parameters.
func (k Keeper) GetMintingParams(ctx sdk.Context) (params mintTypes.Params) {
	return k.mintKeeper.GetParams(ctx)
}

// SetMintingParams sets the total set of cosmos parameters.
func (k Keeper) SetMintingParams(ctx sdk.Context, params mintTypes.Params) {
	k.mintKeeper.SetParams(ctx, params)
}
