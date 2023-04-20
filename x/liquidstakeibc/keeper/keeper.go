package keeper

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

type Keeper struct {
	cdc codec.BinaryCodec

	AccountKeeper types.AccountKeeper
	BankKeeper    types.BankKeeper

	storeKey   storetypes.StoreKey
	paramSpace paramtypes.Subspace

	msgRouter *baseapp.MsgServiceRouter
}

func NewKeeper(cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	accountKeeper types.AccountKeeper, bankKeeper types.BankKeeper,
	paramSpace paramtypes.Subspace, msgRouter *baseapp.MsgServiceRouter,
) Keeper {
	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}
	return Keeper{
		cdc:           cdc,
		AccountKeeper: accountKeeper,
		BankKeeper:    bankKeeper,
		storeKey:      storeKey,
		paramSpace:    paramSpace,
		msgRouter:     msgRouter,
	}
}

// GetParams gets the total set of liquidstakeibc parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the total set of liquidstakeibc parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}
