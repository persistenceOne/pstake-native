package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrInvalidDenom         = errorsmod.Register(ModuleName, 2000, "invalid token denom")
	ErrInvalidHostChain     = errorsmod.Register(ModuleName, 2001, "host chain not registered")
	ErrMinDeposit           = errorsmod.Register(ModuleName, 2002, "deposit amount less than minimum deposit")
	ErrFailedDeposit        = errorsmod.Register(ModuleName, 2003, "deposit failed")
	ErrMintFailed           = errorsmod.Register(ModuleName, 2004, "minting failed")
	ErrRegisterFailed       = errorsmod.Register(ModuleName, 2005, "host chain register failed")
	ErrFailedICQRequest     = errorsmod.Register(ModuleName, 2006, "icq failed")
	ErrDepositNotFound      = errorsmod.Register(ModuleName, 2007, "deposit record not found")
	ErrICATxFailure         = errorsmod.Register(ModuleName, 2008, "ica transaction failed")
	ErrInvalidMessages      = errorsmod.Register(ModuleName, 2009, "not enough messages")
	ErrInvalidResponses     = errorsmod.Register(ModuleName, 2010, "not enough message responses")
	ErrValidatorNotFound    = errorsmod.Register(ModuleName, 2011, "validator not found")
	ErrNotEnoughDelegations = errorsmod.Register(ModuleName, 2012, "delegated amount is less than undelegation amount requested")
	ErrRedeemFailed         = errorsmod.Register(ModuleName, 2013, "an error occurred while instant redeeming tokens")
	ErrBurnFailed           = errorsmod.Register(ModuleName, 2014, "burn failed")
	ErrParsingAmount        = errorsmod.Register(ModuleName, 2015, "could not parse message amount")
	ErrHostChainInactive    = errorsmod.Register(ModuleName, 2016, "host chain is not active")
)
