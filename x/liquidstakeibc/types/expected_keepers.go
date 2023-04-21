package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
)

// AccountKeeper defines the account contract that must be fulfilled when
// creating a x/liquidstakeibc keeper.
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
}

// BankKeeper defines the bankkeeper contract that must be fulfilled when
// creating a x/liquidstakeibc keeper.
type BankKeeper interface {
	MintCoins(ctx sdk.Context, name string, amt sdk.Coins) error
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}

type ScopedKeeper interface {
	ClaimCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) error
}

type ICAControllerKeeper interface {
	RegisterInterchainAccount(ctx sdk.Context, connectionID, owner string, version string) error
	GetInterchainAccountAddress(ctx sdk.Context, connectionID, portID string) (string, bool)
}
