package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/lscosmos module sentinel errors
var (
	ErrSample                                = sdkerrors.Register(ModuleName, 1100, "sample error")
	ErrInvalidPacketTimeout                  = sdkerrors.Register(ModuleName, 1500, "invalid packet timeout")
	ErrInvalidVersion                        = sdkerrors.Register(ModuleName, 1501, "invalid version")
	ErrInvalidMessage                        = sdkerrors.Register(ModuleName, 64, "invalid message")
	ErrInvalidDenom                          = sdkerrors.Register(ModuleName, 65, "denom not whitelisted/ invalid denom")
	ErrInvalidArgs                           = sdkerrors.Register(ModuleName, 66, "invalid arguments")
	ErrFailedDeposit                         = sdkerrors.Register(ModuleName, 67, "deposit failed")
	ErrMintFailed                            = sdkerrors.Register(ModuleName, 68, "minting failed")
	ErrMinDeposit                            = sdkerrors.Register(ModuleName, 69, "deposit amount less than minimum deposit")
	ErrInvalidDenomPath                      = sdkerrors.Register(ModuleName, 70, "denomPath invalid")
	ErrModuleDisabled                        = sdkerrors.Register(ModuleName, 71, "Module is not enabled/ disabled")
	ErrInvalidIntParse                       = sdkerrors.Register(ModuleName, 72, "unable to parse to sdk.Int")
	ErrICATxFailure                          = sdkerrors.Register(ModuleName, 73, "ica transaction failed")
	ErrCannotRemoveNonExistentDelegation     = sdkerrors.Register(ModuleName, 74, "Cannot remove delegation from a non existing delegation")
	ErrInValidAllowListedValidators          = sdkerrors.Register(ModuleName, 75, "invalid allow listed validators")
	ErrInvalidFee                            = sdkerrors.Register(ModuleName, 76, "invalid fee")
	ErrInvalidDeposit                        = sdkerrors.Register(ModuleName, 77, "invalid deposit")
	ErrNoHostChainDelegations                = sdkerrors.Register(ModuleName, 78, "no delegations on host chain")
	ErrCannotRemoveNonExistentUndelegation   = sdkerrors.Register(ModuleName, 79, "Cannot remove undelegation from a non existing undelegation")
	ErrBurnFailed                            = sdkerrors.Register(ModuleName, 80, "burn failed")
	ErrUndelegationEpochNotFound             = sdkerrors.Register(ModuleName, 81, "undelegation epoch not found")
	ErrTransientUndelegationTransferNotFound = sdkerrors.Register(ModuleName, 82, "Transient undelegation transfer not found")
	ErrHostChainDelegationsLTUndelegations   = sdkerrors.Register(ModuleName, 83, "Host chain delegated amount is less than undelegations requested.")
	ErrInvalidHostAccountOwnerIDs            = sdkerrors.Register(ModuleName, 84, "Host account owner ids are not set, was it present in default genesis?")
	ErrModuleAlreadyEnabled                  = sdkerrors.Register(ModuleName, 85, "Module is already enabled")
)
