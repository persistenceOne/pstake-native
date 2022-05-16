package cosmos

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/persistenceOne/pstake-native/x/cosmos/keeper"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

// NewHandler returns a handler for "Gravity" type messages.
func NewHandler(k keeper.Keeper) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case *cosmosTypes.MsgMintTokensForAccount:
			res, err := msgServer.MintTokensForAccount(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *cosmosTypes.MsgSetOrchestrator:
			res, err := msgServer.SetOrchestrator(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *cosmosTypes.MsgMakeProposal:
			res, err := msgServer.MakeProposal(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *cosmosTypes.MsgVote:
			res, err := msgServer.Vote(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *cosmosTypes.MsgVoteWeighted:
			res, err := msgServer.VoteWeighted(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *cosmosTypes.MsgWithdrawStkAsset:
			res, err := msgServer.Withdraw(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized Cosmos Module Msg type: %v", sdk.MsgTypeURL(msg)))
		}
	}
}

func NewCosmosLiquidStakingParametersHandler(k keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *cosmosTypes.ChangeMultisigProposal:
			return keeper.HandleChangeMultisigProposal(ctx, k, c)
		case *cosmosTypes.EnableModuleProposal:
			return keeper.HandleEnableModuleProposal(ctx, k, c)
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized distr proposal content type: %T", c)
		}
	}
}
