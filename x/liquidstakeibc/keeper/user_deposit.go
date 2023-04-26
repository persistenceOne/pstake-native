package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	liquidstakeibctypes "github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

func (k *Keeper) SetUserDeposit(ctx sdk.Context, userDeposit *liquidstakeibctypes.UserDeposit) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), liquidstakeibctypes.UserDepositKey)
	bytes := k.cdc.MustMarshal(userDeposit)
	store.Set([]byte(userDeposit.ChainId+userDeposit.Epoch.String()), bytes)
}

func (k *Keeper) DeleteUserDeposit(ctx sdk.Context, userDeposit *liquidstakeibctypes.UserDeposit) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), liquidstakeibctypes.UserDepositKey)
	store.Delete([]byte(userDeposit.ChainId + userDeposit.Epoch.String()))
}

func (k *Keeper) CreateUserDeposits(ctx sdk.Context, epoch int64) {
	hostChains := k.GetAllHostChains(ctx)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), liquidstakeibctypes.UserDepositKey)
	for _, hc := range hostChains {
		userDeposit := &liquidstakeibctypes.UserDeposit{
			ChainId: hc.ChainId,
			Amount:  sdk.NewCoin(hc.GetIBCDenom(), sdk.NewInt(0)),
			Epoch:   sdk.NewInt(epoch),
			State:   0,
		}
		bytes := k.cdc.MustMarshal(userDeposit)
		store.Set([]byte(userDeposit.ChainId+userDeposit.Epoch.String()), bytes)
	}
}

func (k *Keeper) GetUserDepositForChainAndEpoch(
	ctx sdk.Context,
	chainId string,
	epoch int64,
) *liquidstakeibctypes.UserDeposit {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), liquidstakeibctypes.UserDepositKey)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		userDeposit := &liquidstakeibctypes.UserDeposit{}
		k.cdc.MustUnmarshal(iterator.Value(), userDeposit)

		if userDeposit.Epoch.Int64() == epoch &&
			userDeposit.ChainId == chainId {
			return userDeposit
		}
	}

	return nil
}

func (k *Keeper) GetPendingUserDepositsBeforeEpoch(ctx sdk.Context, epoch int64) []*liquidstakeibctypes.UserDeposit {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), liquidstakeibctypes.UserDepositKey)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	userDeposits := make([]*liquidstakeibctypes.UserDeposit, 0)
	for ; iterator.Valid(); iterator.Next() {
		userDeposit := &liquidstakeibctypes.UserDeposit{}
		k.cdc.MustUnmarshal(iterator.Value(), userDeposit)

		if userDeposit.Epoch.Int64() <= epoch &&
			userDeposit.State == liquidstakeibctypes.UserDeposit_DEPOSIT_PENDING {
			userDeposits = append(userDeposits, userDeposit)
		}
	}

	return userDeposits
}
