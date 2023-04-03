package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

// x/lscosmos module sentinel errors
var (
	ErrSample                                = errorsmod.Register(ModuleName, 1100, "sample error")
	ErrInvalidPacketTimeout                  = errorsmod.Register(ModuleName, 1500, "invalid packet timeout")
	ErrInvalidVersion                        = errorsmod.Register(ModuleName, 1501, "invalid version")
	ErrInvalidMessage                        = errorsmod.Register(ModuleName, 64, "invalid message")
	ErrInvalidDenom                          = errorsmod.Register(ModuleName, 65, "denom not allow-listed/ invalid denom")
	ErrInvalidArgs                           = errorsmod.Register(ModuleName, 66, "invalid arguments")
	ErrFailedDeposit                         = errorsmod.Register(ModuleName, 67, "deposit failed")
	ErrMintFailed                            = errorsmod.Register(ModuleName, 68, "minting failed")
	ErrMinDeposit                            = errorsmod.Register(ModuleName, 69, "deposit amount less than minimum deposit")
	ErrInvalidDenomPath                      = errorsmod.Register(ModuleName, 70, "denomPath invalid")
	ErrModuleDisabled                        = errorsmod.Register(ModuleName, 71, "Module is not enabled/ disabled")
	ErrInvalidIntParse                       = errorsmod.Register(ModuleName, 72, "unable to parse to sdk.Int")
	ErrICATxFailure                          = errorsmod.Register(ModuleName, 73, "ica transaction failed")
	ErrCannotRemoveNonExistentDelegation     = errorsmod.Register(ModuleName, 74, "Cannot remove delegation from a non existing delegation")
	ErrInValidAllowListedValidators          = errorsmod.Register(ModuleName, 75, "invalid allow listed validators")
	ErrInvalidFee                            = errorsmod.Register(ModuleName, 76, "invalid fee")
	ErrInvalidDeposit                        = errorsmod.Register(ModuleName, 77, "invalid deposit")
	ErrNoHostChainDelegations                = errorsmod.Register(ModuleName, 78, "no delegations on host chain")
	ErrCannotRemoveNonExistentUndelegation   = errorsmod.Register(ModuleName, 79, "Cannot remove undelegation from a non existing undelegation")
	ErrBurnFailed                            = errorsmod.Register(ModuleName, 80, "burn failed")
	ErrUndelegationEpochNotFound             = errorsmod.Register(ModuleName, 81, "undelegation epoch not found")
	ErrTransientUndelegationTransferNotFound = errorsmod.Register(ModuleName, 82, "Transient undelegation transfer not found")
	ErrHostChainDelegationsLTUndelegations   = errorsmod.Register(ModuleName, 83, "Host chain delegated amount is less than undelegations requested.")
	ErrInvalidHostAccountOwnerIDs            = errorsmod.Register(ModuleName, 84, "Host account owner ids are not set, was it present in default genesis?")
	ErrModuleAlreadyEnabled                  = errorsmod.Register(ModuleName, 85, "Module is already enabled")
	ErrInvalidMsgs                           = errorsmod.Register(ModuleName, 86, "Invalid Msgs")
	ErrEqualBaseAndMintDenom                 = errorsmod.Register(ModuleName, 87, "BaseDenom and mintDenom cannot be same")
	ErrInsufficientFundsToUndelegate         = errorsmod.Register(ModuleName, 88, "undelegation amount greater than already staked")
	ErrInvalidMintDenom                      = errorsmod.Register(ModuleName, 89, "InvalidMintDenom, MintDenom should be stk/BaseDenom")
	ErrModuleNotInitialised                  = errorsmod.Register(ModuleName, 90, "ErrModuleNotInitialised, Module was never initialised")
	ErrModuleAlreadyInExpectedState          = errorsmod.Register(ModuleName, 91, "ModuleAlreadyInExpectedState, Module is already in expected state")
)
