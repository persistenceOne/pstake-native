package types

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	epochsTypes "github.com/persistenceOne/pstake-native/x/epochs/types"
)

// BankKeeper defines the contract needed to be fulfilled for banking and supply
// dependencies.
type BankKeeper interface {
	MintCoins(ctx sdkTypes.Context, name string, amt sdkTypes.Coins) error
	BurnCoins(ctx sdkTypes.Context, name string, amt sdkTypes.Coins) error
	GetBalance(ctx sdkTypes.Context, addr sdkTypes.AccAddress, denom string) sdkTypes.Coin
	SendCoinsFromModuleToAccount(ctx sdkTypes.Context, senderModule string, recipientAddr sdkTypes.AccAddress, amt sdkTypes.Coins) error
	SendCoinsFromModuleToModule(ctx sdkTypes.Context, senderModule, recipientModule string, amt sdkTypes.Coins) error
	SendCoinsFromAccountToModule(ctx sdkTypes.Context, senderAddr sdkTypes.AccAddress, recipientModule string, amt sdkTypes.Coins) error
}

// DBHelper defines the contract needed to be fulfilled for structs with validator addresses array and ratios
type DBHelper interface {
	Find(address string) bool
	UpdateValues(address string, validatorCount int64)
}

// GovHooks event hooks for governance proposal object (noalias)
type GovHooks interface {
	AfterProposalSubmission(ctx sdkTypes.Context, proposalID uint64)                          // Must be called after proposal is submitted
	AfterProposalVote(ctx sdkTypes.Context, proposalID uint64, voterAddr sdkTypes.AccAddress) // Must be called after a vote on a proposal is cast
	AfterProposalVotingPeriodEnded(ctx sdkTypes.Context, proposalID uint64)                   // Must be called when proposal's finishes it's voting period
}

// EpochKeeper defines the contract needed to be fulfilled for epochs keeper
type EpochKeeper interface {
	GetEpochInfo(ctx sdkTypes.Context, identifier string) epochsTypes.EpochInfo
}
