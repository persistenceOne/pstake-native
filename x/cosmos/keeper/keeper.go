package keeper

import (
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	authKeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	mintKeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	mintTypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramsTypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingKeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	tmLog "github.com/tendermint/tendermint/libs/log"

	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
	epochsTypes "github.com/persistenceOne/pstake-native/x/epochs/types"
)

// Keeper of the cosmos store
type Keeper struct {
	cdc           codec.BinaryCodec
	storeKey      sdkTypes.StoreKey
	paramSpace    paramsTypes.Subspace
	authKeeper    *authKeeper.AccountKeeper
	bankKeeper    *bankKeeper.BaseKeeper
	mintKeeper    *mintKeeper.Keeper
	stakingKeeper *stakingKeeper.Keeper
	hooks         cosmosTypes.GovHooks
	epochsKeeper  cosmosTypes.EpochKeeper
	distrKeeper   *distrkeeper.Keeper
}

// NewKeeper creates a new cosmos Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec, key sdkTypes.StoreKey, paramSpace paramsTypes.Subspace, authKeeper *authKeeper.AccountKeeper,
	bankKeeper *bankKeeper.BaseKeeper, mintKeeper *mintKeeper.Keeper, stakingKeeper *stakingKeeper.Keeper,
	epochKeeper cosmosTypes.EpochKeeper, distrKeeper *distrkeeper.Keeper,
) Keeper {

	return Keeper{
		cdc:           cdc,
		storeKey:      key,
		paramSpace:    paramSpace.WithKeyTable(cosmosTypes.ParamKeyTable()),
		authKeeper:    authKeeper,
		bankKeeper:    bankKeeper,
		mintKeeper:    mintKeeper,
		stakingKeeper: stakingKeeper,
		epochsKeeper:  epochKeeper,
		distrKeeper:   distrKeeper,
	}
}

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

//______________________________________________________________________

// mints tokens in the given account after multiplying with the cValue
func (k Keeper) mintTokens(ctx sdkTypes.Context, mintStoreValue cosmosTypes.MsgMintTokensForAccount) error {
	destinationAddress, err := sdkTypes.AccAddressFromBech32(mintStoreValue.AddressFromMemo)
	if err != nil {
		return err
	}

	// convert the amount to minting amount and multiply by C value
	toBeMinted := mintStoreValue.Amount.Amount.ToDec().Mul(k.GetCValue(ctx))

	// remainder is not considered as they cannot be minted
	mintingAmount, _ := sdkTypes.NewDecCoinFromDec(k.GetParams(ctx).MintDenom, toBeMinted).TruncateDecimal()

	if mintingAmount.Amount.GT(k.GetParams(ctx).MinMintingAmount.Amount) && mintingAmount.Amount.LT(k.GetParams(ctx).MaxMintingAmount.Amount) {
		amnt := sdkTypes.NewCoins(mintingAmount)
		err = k.bankKeeper.MintCoins(ctx, cosmosTypes.ModuleName, amnt)
		if err != nil {
			return err
		}
		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, cosmosTypes.ModuleName, destinationAddress, amnt)
		if err != nil {
			return err
		}

		// add the minted amount
		k.AddToMinted(ctx, mintingAmount)

		// add to virtually staked amount
		k.AddToVirtuallyStaked(ctx, mintStoreValue.Amount)
	}
	return nil
}

// mint tokens for the given rewards receiver account with the given amount
func (k Keeper) mintTokensForRewardReceivers(ctx sdkTypes.Context, address sdkTypes.AccAddress, amount sdkTypes.Coin) error {
	toBeMinted := sdkTypes.NewCoins(amount)

	err := k.bankKeeper.MintCoins(ctx, cosmosTypes.ModuleName, toBeMinted)
	if err != nil {
		return err
	}

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, cosmosTypes.ModuleName, address, toBeMinted)
	if err != nil {
		return err
	}

	// add the minted amount
	k.AddToMinted(ctx, amount)

	return nil
}

// InsertActiveProposalQueue inserts a ProposalID into the active proposal queue at endTime
func (k Keeper) InsertActiveProposalQueue(ctx sdkTypes.Context, proposalID uint64, endTime time.Time) {
	store := ctx.KVStore(k.storeKey)
	bz := cosmosTypes.GetProposalIDBytes(proposalID)
	store.Set(cosmosTypes.ActiveProposalQueueKey(proposalID, endTime), bz)
}

// RemoveFromActiveProposalQueue removes a proposalID from the Active Proposal Queue
func (k Keeper) RemoveFromActiveProposalQueue(ctx sdkTypes.Context, proposalID uint64, endTime time.Time) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(cosmosTypes.ActiveProposalQueueKey(proposalID, endTime))
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdkTypes.Context) tmLog.Logger {
	return ctx.Logger().With("module", "x/"+cosmosTypes.ModuleName)
}
