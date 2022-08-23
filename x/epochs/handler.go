package epochs

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/persistenceOne/pstake-native/x/epochs/keeper"
	"github.com/persistenceOne/pstake-native/x/epochs/types"
)

// NewHandler returns a handler for epochs module messages
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		//nolint:staticcheck,unused_ctx
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		// switch case is not present as it throws a linting error
		errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
	}
}
