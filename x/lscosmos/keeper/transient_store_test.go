package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

func (suite *IntegrationTestSuite) TestTransitionStore() {
	app, ctx := suite.app, suite.ctx

	lscosmosKeeper := app.LSCosmosKeeper

	// set, get and match IBCAmountTransientStore
	ibcAmountTransitionStore := types.IBCAmountTransientStore{
		IBCTransfer: sdk.NewCoins(sdk.NewCoin(lscosmosKeeper.GetHostChainParams(ctx).BaseDenom, sdk.NewInt(10000))),
		ICADelegate: sdk.NewCoin(lscosmosKeeper.GetHostChainParams(ctx).BaseDenom, sdk.NewInt(10000)),
	}
	lscosmosKeeper.SetIBCTransientStore(ctx, ibcAmountTransitionStore)
	suite.Equal(ibcAmountTransitionStore, lscosmosKeeper.GetIBCTransientStore(ctx))

	// add to IBC transfer in transient store
	lscosmosKeeper.AddIBCTransferToTransientStore(ctx, sdk.NewCoin(lscosmosKeeper.GetHostChainParams(ctx).BaseDenom, sdk.NewInt(10000)))
	ibcAmountTransitionStore.IBCTransfer = ibcAmountTransitionStore.IBCTransfer.Add(sdk.NewCoin(lscosmosKeeper.GetHostChainParams(ctx).BaseDenom, sdk.NewInt(10000)))
	suite.Equal(ibcAmountTransitionStore, lscosmosKeeper.GetIBCTransientStore(ctx))

	// remove from IBC transfer in transient store
	lscosmosKeeper.RemoveIBCTransferFromTransientStore(ctx, sdk.NewCoin(lscosmosKeeper.GetHostChainParams(ctx).BaseDenom, sdk.NewInt(10000)))
	ibcAmountTransitionStore.IBCTransfer = sdk.NewCoins(sdk.NewCoin(lscosmosKeeper.GetHostChainParams(ctx).BaseDenom, sdk.NewInt(10000)))
	suite.Equal(ibcAmountTransitionStore, lscosmosKeeper.GetIBCTransientStore(ctx))

	// add to ICA delegate in transient store
	lscosmosKeeper.AddICADelegateToTransientStore(ctx, sdk.NewCoin(lscosmosKeeper.GetHostChainParams(ctx).BaseDenom, sdk.NewInt(10000)))
	ibcAmountTransitionStore.ICADelegate = ibcAmountTransitionStore.ICADelegate.Add(sdk.NewCoin(lscosmosKeeper.GetHostChainParams(ctx).BaseDenom, sdk.NewInt(10000)))
	suite.Equal(ibcAmountTransitionStore, lscosmosKeeper.GetIBCTransientStore(ctx))

	// add to ICA delegate in transient store
	lscosmosKeeper.RemoveICADelegateFromTransientStore(ctx, sdk.NewCoin(lscosmosKeeper.GetHostChainParams(ctx).BaseDenom, sdk.NewInt(10000)))
	ibcAmountTransitionStore.ICADelegate = sdk.NewCoin(lscosmosKeeper.GetHostChainParams(ctx).BaseDenom, sdk.NewInt(10000))
	suite.Equal(ibcAmountTransitionStore, lscosmosKeeper.GetIBCTransientStore(ctx))
}
