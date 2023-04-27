package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	liquidstakeibctypes "github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

func (k *Keeper) SetDeposit(ctx sdk.Context, deposit *liquidstakeibctypes.Deposit) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), liquidstakeibctypes.DepositKey)
	bytes := k.cdc.MustMarshal(deposit)
	store.Set([]byte(deposit.ChainId+deposit.Epoch.String()), bytes)
}

func (k *Keeper) DeleteDeposit(ctx sdk.Context, deposit *liquidstakeibctypes.Deposit) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), liquidstakeibctypes.DepositKey)
	store.Delete([]byte(deposit.ChainId + deposit.Epoch.String()))
}

func (k *Keeper) CreateDeposits(ctx sdk.Context, epoch int64) {
	hostChains := k.GetAllHostChains(ctx)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), liquidstakeibctypes.DepositKey)
	for _, hc := range hostChains {
		deposit := &liquidstakeibctypes.Deposit{
			ChainId:     hc.ChainId,
			Amount:      sdk.NewCoin(hc.GetIBCDenom(), sdk.NewInt(0)),
			Epoch:       sdk.NewInt(epoch),
			State:       0,
			IbcSequence: sdk.NewInt(0),
		}
		bytes := k.cdc.MustMarshal(deposit)
		store.Set([]byte(deposit.ChainId+deposit.Epoch.String()), bytes)
	}
}

func (k *Keeper) GetDepositForChainAndEpoch(
	ctx sdk.Context,
	chainId string,
	epoch int64,
) (*liquidstakeibctypes.Deposit, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), liquidstakeibctypes.DepositKey)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		deposit := &liquidstakeibctypes.Deposit{}
		k.cdc.MustUnmarshal(iterator.Value(), deposit)

		if deposit.Epoch.Int64() == epoch &&
			deposit.ChainId == chainId {
			return deposit, true
		}
	}

	return nil, false
}

func (k *Keeper) GetPendingDepositsBeforeEpoch(ctx sdk.Context, epoch int64) []*liquidstakeibctypes.Deposit {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), liquidstakeibctypes.DepositKey)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	userDeposits := make([]*liquidstakeibctypes.Deposit, 0)
	for ; iterator.Valid(); iterator.Next() {
		deposit := &liquidstakeibctypes.Deposit{}
		k.cdc.MustUnmarshal(iterator.Value(), deposit)

		if deposit.Epoch.Int64() <= epoch &&
			deposit.State == liquidstakeibctypes.Deposit_DEPOSIT_PENDING {
			userDeposits = append(userDeposits, deposit)
		}
	}

	return userDeposits
}

func (k *Keeper) GetDepositFromSequence(ctx sdk.Context, sequence uint64) (*liquidstakeibctypes.Deposit, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), liquidstakeibctypes.DepositKey)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		deposit := &liquidstakeibctypes.Deposit{}
		k.cdc.MustUnmarshal(iterator.Value(), deposit)

		if deposit.IbcSequence.Uint64() == sequence {
			return deposit, true
		}
	}

	return nil, false
}
