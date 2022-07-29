package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v3/modules/core/exported"
)

// ls-cosmosChannelKeeper is a stub of cosmosibckeeper.ChannelKeeper.
type lscosmosChannelKeeper struct{}

func (lscosmosChannelKeeper) GetChannel(ctx sdk.Context, srcPort, srcChan string) (channel channeltypes.Channel, found bool) {
	return channeltypes.Channel{}, false
}
func (lscosmosChannelKeeper) GetNextSequenceSend(ctx sdk.Context, portID, channelID string) (uint64, bool) {
	return 0, false
}
func (lscosmosChannelKeeper) SendPacket(ctx sdk.Context, channelCap *capabilitytypes.Capability, packet ibcexported.PacketI) error {
	return nil
}
func (lscosmosChannelKeeper) ChanCloseInit(ctx sdk.Context, portID, channelID string, chanCap *capabilitytypes.Capability) error {
	return nil
}

// lscosmosPortKeeper is a stub of cosmosibckeeper.PortKeeper
type lscosmosPortKeeper struct{}

func (lscosmosPortKeeper) BindPort(ctx sdk.Context, portID string) *capabilitytypes.Capability {
	return &capabilitytypes.Capability{}
}

// LscosmosKeeper
// todo : either create test app or fix this (genesis test and msg server test depends on it)
// Can totally get rid of this and create test suites for individual folder is ls-cosmos
//func LscosmosKeeper(t testing.TB) (*keeper.Keeper, sdk.Context) {
//	logger := log.NewNopLogger()
//
//	storeKey := sdk.NewKVStoreKey(types.StoreKey)
//	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)
//
//	db := tmdb.NewMemDB()
//	stateStore := store.NewCommitMultiStore(db)
//	stateStore.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
//	stateStore.MountStoreWithDB(memStoreKey, sdk.StoreTypeMemory, nil)
//	require.NoError(t, stateStore.LoadLatestVersion())
//
//	registry := codectypes.NewInterfaceRegistry()
//	appCodec := codec.NewProtoCodec(registry)
//	capabilityKeeper := capabilitykeeper.NewKeeper(appCodec, storeKey, memStoreKey)
//
//	paramsSubspace := typesparams.NewSubspace(appCodec,
//		types.Amino,
//		storeKey,
//		memStoreKey,
//		"LscosmosParams",
//	)
//
//	ibcKeeper := ibckeeper.NewKeeper(appCodec,
//		storeKey,
//		paramsSubspace,
//		stakingkeeper.Keeper{},
//		upgradekeeper.Keeper{},
//		capabilityKeeper.ScopeToModule(types.ModuleName))
//
//	k := keeper.NewKeeper(
//		appCodec,
//		storeKey,
//		memStoreKey,
//		paramsSubspace,
//		*ibcKeeper,
//		capabilityKeeper.ScopeToModule(types.ModuleName),
//	)
//
//	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, logger)
//
//	// Initialize params
//	k.SetParams(ctx, types.DefaultParams())
//
//	return &k, ctx
//}
