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

func (k *Keeper) RevertDepositsState(ctx sdk.Context, deposits []*liquidstakeibctypes.Deposit) {
	for _, deposit := range deposits {
		deposit.IbcSequenceId = ""

		if deposit.State != liquidstakeibctypes.Deposit_DEPOSIT_PENDING {
			deposit.State--
		}

		k.SetDeposit(ctx, deposit)
	}
}

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

func (k *Keeper) GetTransactionSequenceID(channelID string, sequence uint64) string {
	sequenceStr := strconv.FormatUint(sequence, 10)
	return channelID + "-sequence-" + sequenceStr
}

func (k *Keeper) AdjustDepositsForRedemption(
	ctx sdk.Context,
	hc *liquidstakeibctypes.HostChain,
	redeemableAmount sdk.Coin,
) error {
	redeemableDeposits, depositsAmount := k.GetRedeemableDepositsForHostChain(ctx, hc)
	if depositsAmount.LT(redeemableAmount.Amount) {
		return nil
	}

	for _, deposit := range redeemableDeposits {
		// there is enough tokens in this deposit to fulfill the redeem request
		if deposit.Amount.Amount.GT(redeemableAmount.Amount) || redeemableAmount.IsZero() {
			deposit.Amount = deposit.Amount.Sub(redeemableAmount)
			k.SetDeposit(ctx, deposit)
			return nil
		}

		// the deposit is not enough to fulfill the redeem request, use it and remove it
		redeemableAmount = redeemableAmount.Sub(deposit.Amount)
		k.DeleteDeposit(ctx, deposit)
	}

	return nil
}

// TODO: There is many repeated code, have just 1 iterative method and pass in a condition.

func (k *Keeper) GetDepositForChainAndEpoch(
	ctx sdk.Context,
	chainID string,
	epoch int64,
) (*liquidstakeibctypes.Deposit, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), liquidstakeibctypes.DepositKey)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		deposit := &liquidstakeibctypes.Deposit{}
		k.cdc.MustUnmarshal(iterator.Value(), deposit)

		if deposit.Epoch.Int64() == epoch &&
			deposit.ChainId == chainID {
			return deposit, true
		}
	}

	return nil, false
}

func (k *Keeper) GetDepositsWithSequenceID(ctx sdk.Context, sequenceID string) []*liquidstakeibctypes.Deposit {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), liquidstakeibctypes.DepositKey)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	deposits := make([]*liquidstakeibctypes.Deposit, 0)
	for ; iterator.Valid(); iterator.Next() {
		deposit := &liquidstakeibctypes.Deposit{}
		k.cdc.MustUnmarshal(iterator.Value(), deposit)

		if deposit.IbcSequenceId == sequenceID {
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

func (k *Keeper) GetRedeemableDepositsForHostChain(
	ctx sdk.Context,
	hc *liquidstakeibctypes.HostChain,
) ([]*liquidstakeibctypes.Deposit, sdk.Int) { //nolint:staticcheck
	store := prefix.NewStore(ctx.KVStore(k.storeKey), liquidstakeibctypes.DepositKey)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	deposits := make([]*liquidstakeibctypes.Deposit, 0)
	redeemableAmount := sdk.ZeroInt()
	for ; iterator.Valid(); iterator.Next() {
		deposit := &liquidstakeibctypes.Deposit{}
		k.cdc.MustUnmarshal(iterator.Value(), deposit)

		if deposit.ChainId == hc.ChainId &&
			deposit.State == liquidstakeibctypes.Deposit_DEPOSIT_PENDING &&
			!deposit.Amount.IsZero() {
			redeemableAmount = redeemableAmount.Add(deposit.Amount.Amount)
			deposits = append(deposits, deposit)
		}
	}

	return deposits, redeemableAmount
}

func (k *Keeper) GetDelegableDepositsForChain(ctx sdk.Context, chainID string) []*liquidstakeibctypes.Deposit {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), liquidstakeibctypes.DepositKey)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	deposits := make([]*liquidstakeibctypes.Deposit, 0)
	for ; iterator.Valid(); iterator.Next() {
		deposit := &liquidstakeibctypes.Deposit{}
		k.cdc.MustUnmarshal(iterator.Value(), deposit)

		if deposit.ChainId == chainID &&
			deposit.State == liquidstakeibctypes.Deposit_DEPOSIT_RECEIVED {
			deposits = append(deposits, deposit)
		}
	}

	return deposits
}

func (k *Keeper) GetDelegatingDepositsForChain(ctx sdk.Context, chainID string) []*liquidstakeibctypes.Deposit {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), liquidstakeibctypes.DepositKey)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	deposits := make([]*liquidstakeibctypes.Deposit, 0)
	for ; iterator.Valid(); iterator.Next() {
		deposit := &liquidstakeibctypes.Deposit{}
		k.cdc.MustUnmarshal(iterator.Value(), deposit)

		if deposit.ChainId == chainID &&
			deposit.State == liquidstakeibctypes.Deposit_DEPOSIT_DELEGATING {
			deposits = append(deposits, deposit)
		}
	}

	return deposits
}
