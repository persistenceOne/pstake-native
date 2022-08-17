package keeper

import (
	"context"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	ibcTransferTypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"

	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (m msgServer) LiquidStake(goCtx context.Context, msg *types.MsgLiquidStake) (*types.MsgLiquidStakeResponse, error) {
	err := msg.ValidateBasic()
	if err != nil {
		return nil, types.ErrInvalidMessage
	}

	ctx := sdkTypes.UnwrapSDKContext(goCtx)

	// sanity check for the arguments of message
	if ctx.IsZero() || !msg.Amount.IsValid() {
		return nil, types.ErrInvalidArgs
	}
	if !m.GetModuleState(ctx) {
		return nil, types.ErrModuleDisabled
	}
	//GetParams
	ibcParams := m.GetCosmosIBCParams(ctx)

	//check for minimum deposit amount
	if msg.Amount.Amount.LT(ibcParams.MinDeposit) {
		return nil, types.ErrMinDeposit
	}

	expectedIBCPrefix := ibcTransferTypes.GetDenomPrefix(ibcParams.TokenTransferPort, ibcParams.TokenTransferChannel)

	denomTraceStr, err := m.ibcTransferKeeper.DenomPathFromHash(ctx, msg.Amount.Denom)
	if err != nil {
		return nil, err
	}
	denomTrace := ibcTransferTypes.ParseDenomTrace(denomTraceStr)

	// Check if ibc path matches allowlisted path.
	if expectedIBCPrefix != denomTrace.GetPrefix() {
		return nil, types.ErrInvalidDenomPath
	}
	//Check if base denom is valid (uatom) , this can be programmed further to accommodate for liquid staked vouchers.
	if denomTrace.BaseDenom != ibcParams.BaseDenom {
		return nil, types.ErrInvalidDenom
	}

	// check if address in message is correct or not
	delegatorAddress, err := sdkTypes.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, sdkErrors.ErrInvalidAddress
	}

	//send the deposit to the deposit-module account
	depositAmount := sdkTypes.NewCoins(msg.Amount)
	err = m.SendTokensToDepositModule(ctx, depositAmount, delegatorAddress)
	if err != nil {
		return nil, types.ErrFailedDeposit
	}

	// amount of stk tokens to be minted
	mintAmountDec := msg.Amount.Amount.ToDec().Mul(m.GetCValue(ctx))
	mintToken, residue := sdkTypes.NewDecCoinFromDec(ibcParams.MintDenom, mintAmountDec).TruncateDecimal()
	if residue.Amount.GT(sdkTypes.NewDec(0)) {
		m.SendResidueToCommunityPool(ctx, sdkTypes.NewDecCoins(residue))
	}

	//Mint staked representative tokens in lscosmos module account
	err = m.bankKeeper.MintCoins(ctx, types.ModuleName, sdkTypes.NewCoins(mintToken))
	if err != nil {
		return nil, types.ErrMintFailed
	}

	//Calculate protocol fee
	protocolFee := ibcParams.PStakeDepositFee
	protocolFeeAmount := protocolFee.MulInt(mintToken.Amount)
	protocolCoins, residue := sdkTypes.NewDecCoinFromDec(ibcParams.MintDenom, protocolFeeAmount).TruncateDecimal()

	if residue.Amount.GT(sdkTypes.NewDec(0)) {
		m.SendResidueToCommunityPool(ctx, sdkTypes.NewDecCoins(residue))
	}

	//Send (mintedTokens - protocolTokens) to delegator address
	err = m.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, delegatorAddress,
		sdkTypes.NewCoins(mintToken))
	if err != nil {
		return nil, types.ErrMintFailed
	}

	//Send protocol fee to protocol pool
	err = m.SendProtocolFee(ctx, sdkTypes.NewCoins(protocolCoins), delegatorAddress)
	if err != nil {
		return nil, types.ErrFailedDeposit
	}

	//TODO: emit ICA delegator module address?
	ctx.EventManager().EmitEvents(sdkTypes.Events{
		sdkTypes.NewEvent(
			types.EventTypeLiquidStake,
			sdkTypes.NewAttribute(types.AttributeDelegatorAddress, delegatorAddress.String()),
			sdkTypes.NewAttribute(types.AttributeAmountMinted, mintAmountDec.String()),
		),
		sdkTypes.NewEvent(
			sdkTypes.EventTypeMessage,
			sdkTypes.NewAttribute(sdkTypes.AttributeKeyModule, types.AttributeKeyAck),
			sdkTypes.NewAttribute(sdkTypes.AttributeKeySender, msg.DelegatorAddress),
		)},
	)
	return &types.MsgLiquidStakeResponse{}, nil
}
