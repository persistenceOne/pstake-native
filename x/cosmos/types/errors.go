package types

import sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrDuplicate           = sdkErrors.Register(ModuleName, 1, "duplicate")
	ErrInvalid             = sdkErrors.Register(ModuleName, 2, "invalid")
	ErrInvalidVote         = sdkErrors.Register(ModuleName, 3, "invalid vote option")
	ErrResetDelegateKeys   = sdkErrors.Register(ModuleName, 4, "can not set orchestrator addresses more than once")
	ErrEmptyDelegatorAddr  = sdkErrors.Register(ModuleName, 5, "empty delegator address")
	ErrEmptyValidatorAddr  = sdkErrors.Register(ModuleName, 6, "empty validator address")
	ErrInvalidMintingRatio = sdkErrors.Register(ModuleName, 7, "minting ratio less than 0")
	ErrOrchAddressNotFound = sdkErrors.Register(ModuleName, 8, "orchestrator address not found")
	ErrInvalidGenesis      = sdkErrors.Register(ModuleName, 9, "invalid genesis state")
	ErrUnknownProposal     = sdkErrors.Register(ModuleName, 10, "unknown proposal")
	ErrInactiveProposal    = sdkErrors.Register(ModuleName, 11, "inactive proposal")
)
