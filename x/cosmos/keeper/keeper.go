package keeper

import (
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	epochsTypes "github.com/persistenceOne/pstake-native/x/epochs/types"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	mintKeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	mintTypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramsTypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingKeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
	tmLog "github.com/tendermint/tendermint/libs/log"
)

type Keeper struct {
	cdc           codec.BinaryCodec
	storeKey      sdkTypes.StoreKey
	paramSpace    paramsTypes.Subspace
	authKeeper    *authkeeper.AccountKeeper
	bankKeeper    *bankKeeper.BaseKeeper
	mintKeeper    *mintKeeper.Keeper
	stakingKeeper *stakingKeeper.Keeper
	hooks         cosmosTypes.GovHooks
	epochsKeeper  cosmosTypes.EpochKeeper
}

func NewKeeper(
	key sdkTypes.StoreKey, paramSpace paramsTypes.Subspace, authKeeper *authkeeper.AccountKeeper,
	bankKeeper *bankKeeper.BaseKeeper, mintKeeper *mintKeeper.Keeper, stakingKeeper *stakingKeeper.Keeper,
	epochKeeper cosmosTypes.EpochKeeper,
) Keeper {

	return Keeper{
		storeKey:      key,
		paramSpace:    paramSpace.WithKeyTable(cosmosTypes.ParamKeyTable()),
		authKeeper:    authKeeper,
		bankKeeper:    bankKeeper,
		mintKeeper:    mintKeeper,
		stakingKeeper: stakingKeeper,
		epochsKeeper:  epochKeeper,
	}
}

//TODO : Add hooks in app.go
//TODO : Add epoch hooks

// SetHooks sets the hooks for governance
func (k *Keeper) SetHooks(gh cosmosTypes.GovHooks, eh epochsTypes.EpochHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set governance hooks twice")
	}

	k.hooks = gh

	return k
}

//______________________________________________________________________

// GetParams returns the total set of parameters.
func (k Keeper) GetParams(ctx sdkTypes.Context) (params cosmosTypes.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the total set of parameters.
func (k Keeper) SetParams(ctx sdkTypes.Context, params cosmosTypes.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

//______________________________________________________________________

// GetMintingParams returns the total set of cosmos parameters.
func (k Keeper) GetMintingParams(ctx sdkTypes.Context) (params mintTypes.Params) {
	return k.mintKeeper.GetParams(ctx)
}

// SetMintingParams sets the total set of cosmos parameters.
func (k Keeper) SetMintingParams(ctx sdkTypes.Context, params mintTypes.Params) {
	k.mintKeeper.SetParams(ctx, params)
}

func prefixRange(prefix []byte) ([]byte, []byte) {
	if prefix == nil {
		panic("nil key not allowed")
	}
	// special case: no prefix is whole range
	if len(prefix) == 0 {
		return nil, nil
	}

	// copy the prefix and update last byte
	end := make([]byte, len(prefix))
	copy(end, prefix)
	l := len(end) - 1
	end[l]++

	// wait, what if that overflowed?....
	for end[l] == 0 && l > 0 {
		l--
		end[l]++
	}

	// okay, funny guy, you gave us FFF, no end to this range...
	if l == 0 && end[0] == 0 {
		end = nil
	}
	return prefix, end
}

func (k Keeper) mintTokensOnMajority(ctx sdkTypes.Context, key cosmosTypes.ChainIDHeightAndTxHashKey, value cosmosTypes.AddressAndAmountKey) error {
	//TODO incorporate minting_ratio
	if value.Amount.Amount.GT(k.GetParams(ctx).MinMintingAmount.Amount) && value.Amount.Amount.LT(k.GetParams(ctx).MaxMintingAmount.Amount) {
		destinationAddress, err := sdkTypes.AccAddressFromBech32(value.DestinationAddress)
		if err != nil {
			return err
		}
		amnt := sdkTypes.NewCoins(value.Amount)
		err = k.bankKeeper.MintCoins(ctx, cosmosTypes.ModuleName, amnt)
		if err != nil {
			return err
		}
		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, cosmosTypes.ModuleName, destinationAddress, amnt)
		if err != nil {
			return err
		}
	}
	k.setMintedFlagTrue(ctx, key)
	return nil
}

func (k Keeper) mintTokensForRewardReceivers(ctx sdkTypes.Context, address string, amount sdkTypes.Coins) error {
	//TODO : incorporate minting_ratio

	accAddress, err := sdkTypes.AccAddressFromBech32(address)
	if err != nil {
		return err
	}

	err = k.bankKeeper.MintCoins(ctx, cosmosTypes.ModuleName, amount)
	if err != nil {
		return err
	}

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, cosmosTypes.ModuleName, accAddress, amount)
	if err != nil {
		return err
	}

	return nil
}

// InsertActiveProposalQueue inserts a ProposalID into the active proposal queue at endTime
func (keeper Keeper) InsertActiveProposalQueue(ctx sdkTypes.Context, proposalID uint64, endTime time.Time) {
	store := ctx.KVStore(keeper.storeKey)
	bz := cosmosTypes.GetProposalIDBytes(proposalID)
	store.Set(cosmosTypes.ActiveProposalQueueKey(proposalID, endTime), bz)
}

// RemoveFromActiveProposalQueue removes a proposalID from the Active Proposal Queue
func (keeper Keeper) RemoveFromActiveProposalQueue(ctx sdkTypes.Context, proposalID uint64, endTime time.Time) {
	store := ctx.KVStore(keeper.storeKey)
	store.Delete(cosmosTypes.ActiveProposalQueueKey(proposalID, endTime))
}

// Logger returns a module-specific logger.
func (keeper Keeper) Logger(ctx sdkTypes.Context) tmLog.Logger {
	return ctx.Logger().With("module", "x/"+cosmosTypes.ModuleName)
}
