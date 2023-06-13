package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/pstake-native/v2/x/lscosmos/types"
)

func (k Keeper) GetGenesisState(ctx sdk.Context) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.ModuleEnabled = k.GetModuleState(ctx)
	genesis.HostChainParams = k.GetHostChainParams(ctx)
	genesis.AllowListedValidators = k.GetAllowListedValidators(ctx)
	genesis.DelegationState = k.GetDelegationState(ctx)
	genesis.HostChainRewardAddress = k.GetHostChainRewardAddress(ctx)
	genesis.IBCAmountTransientStore = k.GetIBCTransientStore(ctx)
	genesis.UnbondingEpochCValues = k.IterateAllUnbondingEpochCValues(ctx)
	genesis.DelegatorUnbondingEpochEntries = k.IterateAllDelegatorUnbondingEpochEntry(ctx)
	genesis.HostAccounts = k.GetHostAccounts(ctx)
	return genesis
}
