package types

import (
	errorsmod "cosmossdk.io/errors"
)

// Sentinel errors for the lspersistence module.
var (
	ErrActiveLiquidValidatorsNotExists = errorsmod.Register(ModuleName, 2, "active liquid validators not exists")
	ErrInvalidDenom                    = errorsmod.Register(ModuleName, 3, "invalid denom")
	ErrInvalidBondDenom                = errorsmod.Register(ModuleName, 4, "invalid bond denom")
	ErrInvalidLiquidBondDenom          = errorsmod.Register(ModuleName, 5, "invalid liquid bond denom")
	ErrNotImplementedYet               = errorsmod.Register(ModuleName, 6, "not implemented yet")
	ErrLessThanMinLiquidStakingAmount  = errorsmod.Register(ModuleName, 7, "staking amount should be over params.min_liquid_staking_amount")
	ErrInvalidBTokenSupply             = errorsmod.Register(ModuleName, 8, "invalid liquid bond denom supply")
	ErrInvalidActiveLiquidValidators   = errorsmod.Register(ModuleName, 9, "invalid active liquid validators")
	ErrLiquidValidatorsNotExists       = errorsmod.Register(ModuleName, 10, "liquid validators not exists")
	ErrInsufficientProxyAccBalance     = errorsmod.Register(ModuleName, 11, "insufficient liquid tokens or balance of proxy account, need to wait for new liquid validator to be added or unbonding of proxy account to be completed")
	ErrTooSmallLiquidStakingAmount     = errorsmod.Register(ModuleName, 12, "liquid staking amount is too small, the result becomes zero")
	ErrTooSmallLiquidUnstakingAmount   = errorsmod.Register(ModuleName, 13, "liquid unstaking amount is too small, the result becomes zero")
)
