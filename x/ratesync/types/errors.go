package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

// x/ratesync module sentinel errors
var (
	ErrRegisterFailed = errorsmod.Register(ModuleName, 3001, "host chain register failed")
	ErrInvalid        = errorsmod.Register(ModuleName, 3002, "Invalid data")
)
