package types

import sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrDuplicate          = sdkErrors.Register(ModuleName, 1, "duplicate")
	ErrInvalid            = sdkErrors.Register(ModuleName, 2, "invalid")
	ErrInvalidVote        = sdkErrors.Register(ModuleName, 3, "invalid vote option")
	ErrResetDelegateKeys  = sdkErrors.Register(ModuleName, 4, "can not set orchestrator addresses more than once")
	ErrEmptyDelegatorAddr = sdkErrors.Register(ModuleName, 5, "empty delegator address")
	ErrEmptyValidatorAddr = sdkErrors.Register(ModuleName, 6, "empty validator address")
)
