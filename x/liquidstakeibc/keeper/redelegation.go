package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

func (k *Keeper) SetRedelegations(ctx sdk.Context, chainID string, redelegations []*stakingtypes.Redelegation) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.RedelegationsKey)
	bytes := k.cdc.MustMarshal(&types.Redelegations{
		ChainID:       chainID,
		Redelegations: redelegations,
	})
	store.Set(types.GetRedelegationsStoreKey(chainID), bytes)
}

func (k *Keeper) GetRedelegations(ctx sdk.Context, chainID string) (*types.Redelegations, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.RedelegationsKey)
	bz := store.Get(types.GetRedelegationsStoreKey(chainID))
	if bz == nil {
		return nil, false
	}

	var redelegations types.Redelegations
	k.cdc.MustUnmarshal(bz, &redelegations)
	return &redelegations, true
}

func (k *Keeper) DeleteRedelegations(ctx sdk.Context, chainID string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.RedelegationsKey)
	store.Delete(types.GetRedelegationsStoreKey(chainID))
}

func (k *Keeper) AddRedelegationEntry(ctx sdk.Context, chainID string, redelegationMsg stakingtypes.MsgBeginRedelegate, response stakingtypes.MsgBeginRedelegateResponse) {
	redelegations, ok := k.GetRedelegations(ctx, chainID)
	if !ok {
		redelegations = &types.Redelegations{
			ChainID: chainID,
			Redelegations: []*stakingtypes.Redelegation{{
				DelegatorAddress:    redelegationMsg.DelegatorAddress,
				ValidatorSrcAddress: redelegationMsg.ValidatorSrcAddress,
				ValidatorDstAddress: redelegationMsg.ValidatorDstAddress,
				Entries: []stakingtypes.RedelegationEntry{{
					CompletionTime: response.CompletionTime,
				}},
			}},
		}
	} else {
		found := false
		for i, redelegation := range redelegations.Redelegations {
			if redelegation.DelegatorAddress == redelegationMsg.DelegatorAddress &&
				redelegation.ValidatorSrcAddress == redelegationMsg.ValidatorSrcAddress &&
				redelegation.ValidatorDstAddress == redelegationMsg.ValidatorDstAddress {
				found = true
				redelegations.Redelegations[i].Entries = append(redelegations.Redelegations[i].Entries, stakingtypes.RedelegationEntry{
					CompletionTime: response.CompletionTime,
				})
			}
		}
		if !found {
			redelegations.Redelegations = append(redelegations.Redelegations, &stakingtypes.Redelegation{
				DelegatorAddress:    redelegationMsg.DelegatorAddress,
				ValidatorSrcAddress: redelegationMsg.ValidatorSrcAddress,
				ValidatorDstAddress: redelegationMsg.ValidatorDstAddress,
				Entries: []stakingtypes.RedelegationEntry{{
					CompletionTime: response.CompletionTime,
				}},
			})
		}
	}
	k.SetRedelegations(ctx, chainID, redelegations.Redelegations)
}
