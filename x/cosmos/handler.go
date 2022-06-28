package cosmos

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"

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
		case *cosmosTypes.MsgSetSignature:
			res, err := msgServer.SetSignature(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *cosmosTypes.MsgRemoveOrchestrator:
			res, err := msgServer.RemoveOrchestrator(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *cosmosTypes.MsgTxStatus:
			res, err := msgServer.TxStatus(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *cosmosTypes.MsgRewardsClaimedOnCosmosChain:
			res, err := msgServer.RewardsClaimed(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *cosmosTypes.MsgSlashingEventOnCosmosChain:
			res, err := msgServer.SlashingEvent(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			return nil, sdkErrors.Wrap(sdkErrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized Cosmos Module Msg type: %v", sdk.MsgTypeURL(msg)))
		}
	}
}

func NewCosmosLiquidStakingParametersHandler(k keeper.Keeper) govTypes.Handler {
	return func(ctx sdk.Context, content govTypes.Content) error {
		switch c := content.(type) {
		case *cosmosTypes.ChangeMultisigProposal:
			return keeper.HandleChangeMultisigProposal(ctx, k, c)
		case *cosmosTypes.EnableModuleProposal:
			return keeper.HandleEnableModuleProposal(ctx, k, c)
		case *cosmosTypes.ChangeCosmosValidatorWeightsProposal:
			return keeper.HandleChangeCosmosValidatorWeightsProposal(ctx, k, c)
		case *cosmosTypes.ChangeOracleValidatorWeightsProposal:
			return keeper.HandleChangeOracleValidatorWeightsProposal(ctx, k, c)
		default:
			return sdkErrors.Wrapf(sdkErrors.ErrUnknownRequest, "unrecognized distr proposal content type: %T", c)
		}
	}
}
