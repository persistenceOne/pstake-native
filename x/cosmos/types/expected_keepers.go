package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	mintTypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

type BankKeeper interface {
	MintCoins(ctx sdk.Context, name string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, name string, amt sdk.Coins) error
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
}

type MintKeeper interface {
	GetParams(ctx sdk.Context) (params mintTypes.Params)
	SetParams(ctx sdk.Context, params mintTypes.Params)
}

type DBHelper interface {
	Find(address string) bool
	AddAndIncrement(address string)
}

// GovHooks event hooks for governance proposal object (noalias)
type GovHooks interface {
	AfterProposalSubmission(ctx sdk.Context, proposalID uint64)                     // Must be called after proposal is submitted
	AfterProposalVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress) // Must be called after a vote on a proposal is cast
	AfterProposalVotingPeriodEnded(ctx sdk.Context, proposalID uint64)              // Must be called when proposal's finishes it's voting period
}
