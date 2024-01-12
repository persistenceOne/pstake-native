package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	persistencetypes "github.com/persistenceOne/persistence-sdk/v2/x/epochs/types"

	liquidstaketypes "github.com/persistenceOne/pstake-native/v2/x/liquidstake/types"
	liquidstakeibctypes "github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	GetModuleAccount(ctx sdk.Context, moduleName string) types.ModuleAccountI
}

type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	MintCoins(ctx sdk.Context, name string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, name string, amt sdk.Coins) error
	GetSupply(ctx sdk.Context, denom string) sdk.Coin
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	SendCoins(ctx sdk.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}

type ICAControllerKeeper interface {
	RegisterInterchainAccount(ctx sdk.Context, connectionID, owner, version string) error
	GetInterchainAccountAddress(ctx sdk.Context, connectionID, portID string) (string, bool)
	GetOpenActiveChannel(ctx sdk.Context, connectionID, portID string) (string, bool)
}

type EpochsKeeper interface {
	GetEpochInfo(ctx sdk.Context, identifier string) persistencetypes.EpochInfo
}

type LiquidStakeIBCKeeper interface {
	GetHostChain(ctx sdk.Context, chainID string) (*liquidstakeibctypes.HostChain, bool)
}

type LiquidStakeKeeper interface {
	// add for stkxprt
	GetNetAmountState(ctx sdk.Context) liquidstaketypes.NetAmountState
	LiquidBondDenom(ctx sdk.Context) string
}
