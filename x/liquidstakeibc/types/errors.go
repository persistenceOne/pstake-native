package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrInvalidDenom      = errorsmod.Register(ModuleName, 2000, "invalid token denom")
	ErrInvalidHostChain  = errorsmod.Register(ModuleName, 2001, "host chain not registered")
	ErrMinDeposit        = errorsmod.Register(ModuleName, 2002, "deposit amount less than minimum deposit")
	ErrFailedDeposit     = errorsmod.Register(ModuleName, 2003, "deposit failed")
	ErrMintFailed        = errorsmod.Register(ModuleName, 2004, "minting failed")
	ErrRegisterFailed    = errorsmod.Register(ModuleName, 2005, "host chain register failed")
	ErrInvalidVersion    = errorsmod.Register(ModuleName, 2006, "invalid version")
	ErrFailedICQRequest  = errorsmod.Register(ModuleName, 2007, "icq failed")
	ErrDepositNotFound   = errorsmod.Register(ModuleName, 2008, "deposit record not found")
	ErrICATxFailure      = errorsmod.Register(ModuleName, 2009, "ica transaction failed")
	ErrInvalidMessages   = errorsmod.Register(ModuleName, 2010, "not enough messages")
	ErrInvalidChannelId  = errorsmod.Register(ModuleName, 2011, "invalid channel id")
	ErrInvalidResponses  = errorsmod.Register(ModuleName, 2012, "not enough message responses")
	ErrValidatorNotFound = errorsmod.Register(ModuleName, 2013, "validator not found")
)
