package keeper

import (
	"strconv"

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
			ChainId:       hc.ChainId,
			Amount:        sdk.NewCoin(hc.IBCDenom(), sdk.NewInt(0)),
			Epoch:         sdk.NewInt(epoch),
			State:         liquidstakeibctypes.Deposit_DEPOSIT_PENDING,
			IbcSequenceId: "",
		}
		bytes := k.cdc.MustMarshal(deposit)
		store.Set([]byte(deposit.ChainId+deposit.Epoch.String()), bytes)
	}
}

func (k *Keeper) RevertDepositsWithSequenceId(
	ctx sdk.Context,
	sequenceId string,
) {
	deposits := k.GetDepositsWithSequenceId(ctx, sequenceId)
	for _, deposit := range deposits {
		if deposit.State == liquidstakeibctypes.Deposit_DEPOSIT_PENDING {
			continue
		}

		deposit.IbcSequenceId = ""
		deposit.State = deposit.State - 1
		k.SetDeposit(ctx, deposit)
	}
}

// GetAllDeposits retrieves all deposits
func (k *Keeper) GetAllDeposits(ctx sdk.Context) []*liquidstakeibctypes.Deposit {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), liquidstakeibctypes.DepositKey)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	deposits := make([]*liquidstakeibctypes.Deposit, 0)
	for ; iterator.Valid(); iterator.Next() {
		deposit := liquidstakeibctypes.Deposit{}
		k.cdc.MustUnmarshal(iterator.Value(), &deposit)
		deposits = append(deposits, &deposit)
	}

	return deposits
}

func (k *Keeper) GetDepositSequenceId(channelId string, sequence uint64) string {
	sequenceStr := strconv.FormatUint(sequence, 10)
	return channelId + "-sequence-" + sequenceStr
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

func (k *Keeper) GetDepositsWithSequenceId(ctx sdk.Context, sequenceId string) []*liquidstakeibctypes.Deposit {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), liquidstakeibctypes.DepositKey)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	deposits := make([]*liquidstakeibctypes.Deposit, 0)
	for ; iterator.Valid(); iterator.Next() {
		deposit := &liquidstakeibctypes.Deposit{}
		k.cdc.MustUnmarshal(iterator.Value(), deposit)

		if deposit.IbcSequenceId == sequenceId {
			deposits = append(deposits, deposit)
		}
	}

	return deposits
}

func (k *Keeper) GetPendingDepositsBeforeEpoch(ctx sdk.Context, epoch int64) []*liquidstakeibctypes.Deposit {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), liquidstakeibctypes.DepositKey)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	deposits := make([]*liquidstakeibctypes.Deposit, 0)
	for ; iterator.Valid(); iterator.Next() {
		deposit := &liquidstakeibctypes.Deposit{}
		k.cdc.MustUnmarshal(iterator.Value(), deposit)

		if deposit.Epoch.Int64() <= epoch &&
			deposit.State == liquidstakeibctypes.Deposit_DEPOSIT_PENDING {
			deposits = append(deposits, deposit)
		}
	}

	return deposits
}

func (k *Keeper) GetDelegableDepositsForChain(ctx sdk.Context, chainId string) []*liquidstakeibctypes.Deposit {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), liquidstakeibctypes.DepositKey)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	deposits := make([]*liquidstakeibctypes.Deposit, 0)
	for ; iterator.Valid(); iterator.Next() {
		deposit := &liquidstakeibctypes.Deposit{}
		k.cdc.MustUnmarshal(iterator.Value(), deposit)

		if deposit.ChainId == chainId &&
			deposit.State == liquidstakeibctypes.Deposit_DEPOSIT_RECEIVED {
			deposits = append(deposits, deposit)
		}
	}

	return deposits
}

// TODO: There is many repeated code, have just 1 iterative method and pass in a condition.
