package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	mintKeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	mintTypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramsTypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingKeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/persistenceOne/pstake-native/x/cosmos/types"
	"time"
)

type Keeper struct {
	cdc           codec.BinaryCodec
	storeKey      sdk.StoreKey
	paramSpace    paramsTypes.Subspace
	bankKeeper    *bankKeeper.BaseKeeper
	mintKeeper    *mintKeeper.Keeper
	stakingKeeper *stakingKeeper.Keeper
	hooks         types.GovHooks
}

func NewKeeper(
	key sdk.StoreKey, paramSpace paramsTypes.Subspace,
	bankKeeper *bankKeeper.BaseKeeper, mintKeeper *mintKeeper.Keeper, stakingKeeper *stakingKeeper.Keeper,
) Keeper {

	return Keeper{
		storeKey:      key,
		paramSpace:    paramSpace.WithKeyTable(types.ParamKeyTable()),
		bankKeeper:    bankKeeper,
		mintKeeper:    mintKeeper,
		stakingKeeper: stakingKeeper,
	}
}

// SetHooks sets the hooks for governance
func (k *Keeper) SetHooks(gh types.GovHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set governance hooks twice")
	}

	k.hooks = gh

	return k
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

func (k Keeper) GetDelegateKeys(ctx sdk.Context) []types.MsgSetOrchestrator {
	store := ctx.KVStore(k.storeKey)
	prefix := []byte(types.KeyOrchestratorAddress)
	iter := store.Iterator(prefixRange(prefix))
	defer iter.Close()

	orchAddresses := make(map[string]string)

	for ; iter.Valid(); iter.Next() {
		key := iter.Key()[len(types.KeyOrchestratorAddress):]
		value := iter.Value()
		orchAddress := sdk.AccAddress(key)
		if err := sdk.VerifyAddressFormat(orchAddress); err != nil {
			panic(sdkErrors.Wrapf(err, "invalid orchAddress in key %v", orchAddresses))
		}
		valAddress := sdk.ValAddress(value)
		if err := sdk.VerifyAddressFormat(valAddress); err != nil {
			panic(sdkErrors.Wrapf(err, "invalid val address stored for orchestrator %s", valAddress.String()))
		}

		orchAddresses[valAddress.String()] = orchAddress.String()
	}

	var result []types.MsgSetOrchestrator

	for valAddr := range orchAddresses {
		orch, ok := orchAddresses[valAddr]
		if !ok {
			// this should never happen unless the store
			// is somehow inconsistent
			panic("Can't find address")
		}
		result = append(result, types.MsgSetOrchestrator{
			Orchestrator: orch,
			Validator:    valAddr,
		})

	}

	return result
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

func (k Keeper) MintTokensOnMajority(ctx sdk.Context, key types.ChainIDHeightAndTxHashKey, value types.AddressAndAmountKey) error {
	//TODO incorporate minting_ratio
	destinationAddress, err := sdk.AccAddressFromBech32(value.DestinationAddress)
	if err != nil {
		return err
	}
	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, value.Amount)
	if err != nil {
		return err
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, destinationAddress, value.Amount)
	if err != nil {
		return err
	}

	k.setMintedFlagTrue(ctx, key)
	//k.deleteFromMintPoolTx(ctx, destinationAddress, value.Amount)
	return nil
}

// InsertActiveProposalQueue inserts a ProposalID into the active proposal queue at endTime
func (keeper Keeper) InsertActiveProposalQueue(ctx sdk.Context, proposalID uint64, endTime time.Time) {
	store := ctx.KVStore(keeper.storeKey)
	bz := types.GetProposalIDBytes(proposalID)
	store.Set(types.ActiveProposalQueueKey(proposalID, endTime), bz)
}

// RemoveFromActiveProposalQueue removes a proposalID from the Active Proposal Queue
func (keeper Keeper) RemoveFromActiveProposalQueue(ctx sdk.Context, proposalID uint64, endTime time.Time) {
	store := ctx.KVStore(keeper.storeKey)
	store.Delete(types.ActiveProposalQueueKey(proposalID, endTime))
}
