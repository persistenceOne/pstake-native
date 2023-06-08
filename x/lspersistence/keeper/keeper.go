package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/v2/x/lspersistence/types"
)

// Keeper of the liquidstaking store
type Keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.BinaryCodec

	accountKeeper  types.AccountKeeper
	bankKeeper     types.BankKeeper
	stakingKeeper  types.StakingKeeper
	distrKeeper    types.DistrKeeper
	slashingKeeper types.SlashingKeeper

	authority string
}

// NewKeeper returns a liquidstaking keeper. It handles:
// - creating new ModuleAccounts for each pool ReserveAccount
// - sending to and from ModuleAccounts
// - minting, burning PoolCoins
func NewKeeper(cdc codec.BinaryCodec, key storetypes.StoreKey, accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper, stakingKeeper types.StakingKeeper, distrKeeper types.DistrKeeper,
	slashingKeeper types.SlashingKeeper, authority string,
) Keeper {
	// ensure liquidstaking module account is set
	if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	return Keeper{
		storeKey:       key,
		cdc:            cdc,
		accountKeeper:  accountKeeper,
		bankKeeper:     bankKeeper,
		stakingKeeper:  stakingKeeper,
		distrKeeper:    distrKeeper,
		slashingKeeper: slashingKeeper,
		authority:      authority,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// SetParams sets the auth module's parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	if err := params.Validate(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&params)
	store.Set(types.ParamsKey, bz)

	return nil
}

// GetParams gets the auth module's parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return params
	}
	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// GetCodec return codec.Codec object used by the keeper
func (k Keeper) GetCodec() codec.BinaryCodec { return k.cdc }
