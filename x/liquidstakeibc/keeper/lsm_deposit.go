package keeper

import (
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	liquidstakeibctypes "github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

func (k *Keeper) SetLSMDeposit(ctx sdk.Context, deposit *liquidstakeibctypes.LSMDeposit) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), liquidstakeibctypes.LSMDepositKey)
	bytes := k.cdc.MustMarshal(deposit)
	store.Set(liquidstakeibctypes.GetLSMDepositStoreKey(deposit.ChainId, deposit.DelegatorAddress, deposit.Denom), bytes)
}

func (k *Keeper) GetLSMDeposit(ctx sdk.Context, chainID, delegator, denom string) (*liquidstakeibctypes.LSMDeposit, bool) {
	hc := liquidstakeibctypes.LSMDeposit{}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), liquidstakeibctypes.LSMDepositKey)
	bytes := store.Get(liquidstakeibctypes.GetLSMDepositStoreKey(chainID, delegator, denom))
	if len(bytes) == 0 {
		return &hc, false
	}

	k.cdc.MustUnmarshal(bytes, &hc)
	return &hc, true
}

func (k *Keeper) GetLSMDepositsFromIbcDenom(ctx sdk.Context, ibcDenom string) []*liquidstakeibctypes.LSMDeposit {
	return k.FilterLSMDeposits(
		ctx,
		func(d liquidstakeibctypes.LSMDeposit) bool {
			return d.IbcDenom == ibcDenom
		},
	)
}

func (k *Keeper) GetLSMDepositsFromIbcSequenceID(ctx sdk.Context, ibcSequenceID string) []*liquidstakeibctypes.LSMDeposit {
	return k.FilterLSMDeposits(
		ctx,
		func(d liquidstakeibctypes.LSMDeposit) bool {
			return d.IbcSequenceId == ibcSequenceID
		},
	)
}

func (k *Keeper) GetTransferableLSMDeposits(ctx sdk.Context, chainID string) []*liquidstakeibctypes.LSMDeposit {
	return k.FilterLSMDeposits(
		ctx,
		func(d liquidstakeibctypes.LSMDeposit) bool {
			return d.ChainId == chainID && d.State == liquidstakeibctypes.LSMDeposit_DEPOSIT_PENDING
		},
	)
}

func (k *Keeper) GetRedeemableLSMDeposits(ctx sdk.Context, chainID string) []*liquidstakeibctypes.LSMDeposit {
	return k.FilterLSMDeposits(
		ctx,
		func(d liquidstakeibctypes.LSMDeposit) bool {
			return d.ChainId == chainID && d.State == liquidstakeibctypes.LSMDeposit_DEPOSIT_RECEIVED
		},
	)
}

func (k *Keeper) DeleteLSMDeposit(ctx sdk.Context, deposit *liquidstakeibctypes.LSMDeposit) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), liquidstakeibctypes.LSMDepositKey)
	store.Delete(liquidstakeibctypes.GetLSMDepositStoreKey(deposit.ChainId, deposit.DelegatorAddress, deposit.Denom))
}

func (k *Keeper) RevertLSMDepositsState(ctx sdk.Context, deposits []*liquidstakeibctypes.LSMDeposit) {
	for _, deposit := range deposits {
		deposit.IbcSequenceId = ""

		if deposit.State != liquidstakeibctypes.LSMDeposit_DEPOSIT_PENDING {
			deposit.State--
		}

		k.SetLSMDeposit(ctx, deposit)
	}
}

func (k *Keeper) UpdateLSMDepositsStateAndSequence(
	ctx sdk.Context,
	deposits []*liquidstakeibctypes.LSMDeposit,
	state liquidstakeibctypes.LSMDeposit_LSMDepositState,
	ibcSequence string,
) {
	for _, deposit := range deposits {
		deposit.IbcSequenceId = ibcSequence
		deposit.State = state
		k.SetLSMDeposit(ctx, deposit)
	}
}

func (k *Keeper) FilterLSMDeposits(
	ctx sdk.Context,
	filter func(d liquidstakeibctypes.LSMDeposit) bool,
) []*liquidstakeibctypes.LSMDeposit {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), liquidstakeibctypes.LSMDepositKey)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	deposits := make([]*liquidstakeibctypes.LSMDeposit, 0)
	for ; iterator.Valid(); iterator.Next() {
		deposit := liquidstakeibctypes.LSMDeposit{}
		k.cdc.MustUnmarshal(iterator.Value(), &deposit)
		if filter(deposit) {
			deposits = append(deposits, &deposit)
		}
	}

	return deposits
}

func (k *Keeper) GetLSMDepositAmountUntokenized(ctx sdk.Context, chainID string) math.Int {
	amount := sdk.ZeroInt()

	deposits := k.FilterLSMDeposits(
		ctx,
		func(d liquidstakeibctypes.LSMDeposit) bool {
			return d.ChainId == chainID
		},
	)

	for _, deposit := range deposits {
		amount = amount.Add(deposit.Amount)
	}

	return amount
}
