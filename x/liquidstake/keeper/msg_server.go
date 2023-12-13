package keeper

// DONTCOVER

// Although written in msg_server_test.go, it is approached at the keeper level rather than at the msgServer level
// so is not included in the coverage.

import (
	"context"
	"time"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstake/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the liquidstake MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) LiquidStake(goCtx context.Context, msg *types.MsgLiquidStake) (*types.MsgLiquidStakeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	newShares, stkXPRTMintAmount, err := k.Keeper.LiquidStake(ctx, types.LiquidStakeProxyAcc, msg.GetDelegator(), msg.Amount)
	if err != nil {
		return nil, err
	}

	liquidBondDenom := k.LiquidBondDenom(ctx)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
		sdk.NewEvent(
			types.EventTypeMsgLiquidStake,
			sdk.NewAttribute(types.AttributeKeyDelegator, msg.DelegatorAddress),
			sdk.NewAttribute(types.AttributeKeyLiquidAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyNewShares, newShares.String()),
			sdk.NewAttribute(types.AttributeKeyStkXPRTMintedAmount, sdk.Coin{Denom: liquidBondDenom, Amount: stkXPRTMintAmount}.String()),
		),
	})
	return &types.MsgLiquidStakeResponse{}, nil
}

func (k msgServer) StakeToLP(goCtx context.Context, msg *types.MsgStakeToLP) (*types.MsgStakeToLPResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	newShares, stkXPRTMintAmount, err := k.LSMDelegate(
		ctx,
		msg.GetDelegator(),
		msg.GetValidator(),
		types.LiquidStakeProxyAcc,
		msg.StakedAmount,
	)
	if err != nil {
		return nil, err
	}

	liquidBondDenom := k.LiquidBondDenom(ctx)
	stkXPRTMinted := sdk.Coin{
		Denom:  liquidBondDenom,
		Amount: stkXPRTMintAmount,
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
		sdk.NewEvent(
			types.EventTypeMsgStakeToLP,
			sdk.NewAttribute(types.AttributeKeyDelegator, msg.DelegatorAddress),
			sdk.NewAttribute(types.AttributeKeyStakedAmount, msg.LiquidAmount.String()),
			sdk.NewAttribute(types.AttributeKeyNewShares, newShares.String()),
			sdk.NewAttribute(types.AttributeKeyStkXPRTMintedAmount, stkXPRTMinted.String()),
		),
	})

	if msg.LiquidAmount.Amount.IsPositive() {
		newShares, stkXPRTMintAmount, err := k.Keeper.LiquidStake(ctx, types.LiquidStakeProxyAcc, msg.GetDelegator(), msg.LiquidAmount)
		if err != nil {
			return nil, err
		}

		stkXPRTMinted := sdk.Coin{
			Denom:  liquidBondDenom,
			Amount: stkXPRTMintAmount,
		}

		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			),
			sdk.NewEvent(
				types.EventTypeMsgStakeToLP,
				sdk.NewAttribute(types.AttributeKeyDelegator, msg.DelegatorAddress),
				sdk.NewAttribute(types.AttributeKeyLiquidAmount, msg.LiquidAmount.String()),
				sdk.NewAttribute(types.AttributeKeyNewShares, newShares.String()),
				sdk.NewAttribute(types.AttributeKeyStkXPRTMintedAmount, stkXPRTMinted.String()),
			),
		})

		_, err = k.LockOnLP(ctx, msg.GetDelegator(), stkXPRTMinted)
		if err != nil {
			return nil, err
		}
	}

	return &types.MsgStakeToLPResponse{}, nil
}

func (k msgServer) LiquidUnstake(goCtx context.Context, msg *types.MsgLiquidUnstake) (*types.MsgLiquidUnstakeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	completionTime, unbondingAmount, _, unbondedAmount, err := k.Keeper.LiquidUnstake(ctx, types.LiquidStakeProxyAcc, msg.GetDelegator(), msg.Amount)
	if err != nil {
		return nil, err
	}

	bondDenom := k.stakingKeeper.BondDenom(ctx)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
		sdk.NewEvent(
			types.EventTypeMsgLiquidUnstake,
			sdk.NewAttribute(types.AttributeKeyDelegator, msg.DelegatorAddress),
			sdk.NewAttribute(types.AttributeKeyStakedAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyUnbondingAmount, sdk.Coin{Denom: bondDenom, Amount: unbondingAmount}.String()),
			sdk.NewAttribute(types.AttributeKeyUnbondedAmount, sdk.Coin{Denom: bondDenom, Amount: unbondedAmount}.String()),
			sdk.NewAttribute(types.AttributeKeyCompletionTime, completionTime.Format(time.RFC3339)),
		),
	})
	return &types.MsgLiquidUnstakeResponse{
		CompletionTime: completionTime,
	}, nil
}

func (k msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if msg.Authority != k.authority {
		return nil, errorsmod.Wrapf(sdkerrors.ErrorInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	err := k.SetParams(ctx, msg.Params)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
		sdk.NewEvent(
			types.EventTypeMsgUpdateParams,
			sdk.NewAttribute(types.AttributeKeyAuthority, msg.Authority),
			sdk.NewAttribute(types.AttributeKeyUpdatedParams, msg.Params.String()),
		),
	})

	return &types.MsgUpdateParamsResponse{}, nil
}
