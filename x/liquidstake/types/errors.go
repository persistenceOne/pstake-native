package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

// Sentinel errors for the liquidstake module.
var (
	ErrActiveLiquidValidatorsNotExists = sdkerrors.Register(ModuleName, 2, "active liquid validators not exists")
	ErrInvalidDenom                    = sdkerrors.Register(ModuleName, 3, "invalid denom")
	ErrInvalidBondDenom                = sdkerrors.Register(ModuleName, 4, "invalid bond denom")
	ErrInvalidLiquidBondDenom          = sdkerrors.Register(ModuleName, 5, "invalid liquid bond denom")
	ErrNotImplementedYet               = sdkerrors.Register(ModuleName, 6, "not implemented yet")
	ErrLessThanMinLiquidStakeAmount    = sdkerrors.Register(ModuleName, 7, "staking amount should be over params.min_liquid_stake_amount")
	ErrInvalidStkXPRTSupply            = sdkerrors.Register(ModuleName, 8, "invalid liquid bond denom supply")
	ErrInvalidActiveLiquidValidators   = sdkerrors.Register(ModuleName, 9, "invalid active liquid validators")
	ErrLiquidValidatorsNotExists       = sdkerrors.Register(ModuleName, 10, "liquid validators not exists")
	ErrInsufficientProxyAccBalance     = sdkerrors.Register(ModuleName, 11, "insufficient liquid tokens or balance of proxy account, need to wait for new liquid validator to be added or unbonding of proxy account to be completed")
	ErrTooSmallLiquidStakeAmount       = sdkerrors.Register(ModuleName, 12, "liquid staking amount is too small, the result becomes zero")
	ErrTooSmallLiquidUnstakingAmount   = sdkerrors.Register(ModuleName, 13, "liquid unstaking amount is too small, the result becomes zero")
	ErrNoLPContractAddress             = sdkerrors.Register(ModuleName, 14, "CW address of an LP contract is not set")
	ErrDisabledLSM                     = sdkerrors.Register(ModuleName, 15, "LSM delegation is disabled")
	ErrLSMTokenizeFailed               = sdkerrors.Register(ModuleName, 16, "LSM tokenization failed")
	ErrLSMRedeemFailed                 = sdkerrors.Register(ModuleName, 17, "LSM redemption failed")
	ErrLPContract                      = sdkerrors.Register(ModuleName, 18, "CW contract execution failed")
)
