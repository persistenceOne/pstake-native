package types

import sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrInvalid                     = sdkErrors.Register(ModuleName, 42, "invalid")
	ErrInvalidVote                 = sdkErrors.Register(ModuleName, 43, "invalid vote option")
	ErrInvalidProposal             = sdkErrors.Register(ModuleName, 44, "invalid proposal sender")
	ErrOrchAddressNotFound         = sdkErrors.Register(ModuleName, 48, "orchestrator address not found")
	ErrInvalidGenesis              = sdkErrors.Register(ModuleName, 49, "invalid genesis state")
	ErrUnknownProposal             = sdkErrors.Register(ModuleName, 50, "unknown proposal")
	ErrInactiveProposal            = sdkErrors.Register(ModuleName, 51, "inactive proposal")
	ErrTxnNotPresentInOutgoingPool = sdkErrors.Register(ModuleName, 52, "txn not present in outgoing pool")
	ErrInvalidStatus               = sdkErrors.Register(ModuleName, 54, "invalid status type")
	ErrTxnDetailsAlreadySent       = sdkErrors.Register(ModuleName, 55, "txn signed details already present")
	ErrModuleNotEnabled            = sdkErrors.Register(ModuleName, 56, "module not enabled")
	ErrInvalidWithdrawDenom        = sdkErrors.Register(ModuleName, 57, "invalid withdraw denom")
	ErrInvalidBondDenom            = sdkErrors.Register(ModuleName, 58, "invalid bond denom")
)
