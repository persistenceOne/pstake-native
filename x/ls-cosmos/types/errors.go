package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/lscosmos module sentinel errors
var (
	ErrSample               = sdkerrors.Register(ModuleName, 1100, "sample error")
	ErrInvalidPacketTimeout = sdkerrors.Register(ModuleName, 1500, "invalid packet timeout")
	ErrInvalidVersion       = sdkerrors.Register(ModuleName, 1501, "invalid version")
	ErrInvalidMessage       = sdkerrors.Register(ModuleName, 64, "invalid message")
	ErrInvalidDenom         = sdkerrors.Register(ModuleName, 65, "denom not whitelisted/ invalid denom")
	ErrInvalidArgs          = sdkerrors.Register(ModuleName, 66, "invalid arguments")
	ErrFailedDeposit        = sdkerrors.Register(ModuleName, 67, "deposit failed")
	ErrMintFailed           = sdkerrors.Register(ModuleName, 68, "minting failed")
)
