package types

import sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrDuplicate                   = sdkErrors.Register(ModuleName, 41, "duplicate")
	ErrInvalid                     = sdkErrors.Register(ModuleName, 42, "invalid")
	ErrInvalidVote                 = sdkErrors.Register(ModuleName, 43, "invalid vote option")
	ErrResetDelegateKeys           = sdkErrors.Register(ModuleName, 44, "can not set orchestrator addresses more than once")
	ErrEmptyDelegatorAddr          = sdkErrors.Register(ModuleName, 45, "empty delegator address")
	ErrEmptyValidatorAddr          = sdkErrors.Register(ModuleName, 46, "empty validator address")
	ErrInvalidMintingRatio         = sdkErrors.Register(ModuleName, 47, "minting ratio less than 0")
	ErrOrchAddressNotFound         = sdkErrors.Register(ModuleName, 48, "orchestrator address not found")
	ErrInvalidGenesis              = sdkErrors.Register(ModuleName, 49, "invalid genesis state")
	ErrUnknownProposal             = sdkErrors.Register(ModuleName, 50, "unknown proposal")
	ErrInactiveProposal            = sdkErrors.Register(ModuleName, 51, "inactive proposal")
	ErrTxnNotPresentInOutgoingPool = sdkErrors.Register(ModuleName, 52, "txn not present in outgoing pool")
	ErrInvalidStatus               = sdkErrors.Register(ModuleName, 54, "invalid status type")
	ErrTxnDetailsAlreadySent       = sdkErrors.Register(ModuleName, 55, "txn signed details already present")
)
