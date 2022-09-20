package keeper

import (
	"context"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"

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

	ctx := sdktypes.UnwrapSDKContext(goCtx)

	// sanity check for the arguments of message
	if ctx.IsZero() || !msg.Amount.IsValid() {
		return nil, types.ErrInvalidArgs
	}
	if !m.GetModuleState(ctx) {
		return nil, types.ErrModuleDisabled
	}
	//GetParams
	hostChainParams := m.GetHostChainParams(ctx)

	//check for minimum deposit amount
	if msg.Amount.Amount.LT(hostChainParams.MinDeposit) {
		return nil, types.ErrMinDeposit
	}

	expectedIBCPrefix := ibctransfertypes.GetDenomPrefix(hostChainParams.TransferPort, hostChainParams.TransferChannel)

	denomTraceStr, err := m.ibcTransferKeeper.DenomPathFromHash(ctx, msg.Amount.Denom)
	if err != nil {
		return nil, err
	}
	denomTrace := ibctransfertypes.ParseDenomTrace(denomTraceStr)

	// Check if ibc path matches allowlisted path.
	if expectedIBCPrefix != denomTrace.GetPrefix() {
		return nil, types.ErrInvalidDenomPath
	}
	//Check if base denom is valid (uatom) , this can be programmed further to accommodate for liquid staked vouchers.
	if denomTrace.BaseDenom != hostChainParams.BaseDenom {
		return nil, types.ErrInvalidDenom
	}

	// check if address in message is correct or not
	delegatorAddress, err := sdktypes.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress
	}

	// amount of stk tokens to be minted. We calculate this before depositing any amount so as to not affect minting c-value.
	// We do not care about residue here because it won't be minted and bank.TotalSupply invariant should not be affected
	cValue := m.GetCValue(ctx)
	mintToken, _ := m.ConvertTokenToStk(ctx, sdktypes.NewDecCoinFromCoin(msg.Amount), cValue)

	//send the deposit to the deposit-module account
	depositAmount := sdktypes.NewCoins(msg.Amount)
	err = m.SendTokensToDepositModule(ctx, depositAmount, delegatorAddress)
	if err != nil {
		return nil, types.ErrFailedDeposit
	}

	//Mint staked representative tokens in lscosmos module account
	err = m.bankKeeper.MintCoins(ctx, types.ModuleName, sdktypes.NewCoins(mintToken))
	if err != nil {
		return nil, types.ErrMintFailed
	}

	//Calculate protocol fee
	protocolFee := hostChainParams.PstakeParams.PstakeDepositFee
	protocolFeeAmount := protocolFee.MulInt(mintToken.Amount)
	// We do not care about residue, as to not break Total calculation invariant.
	protocolCoin, _ := sdktypes.NewDecCoinFromDec(hostChainParams.MintDenom, protocolFeeAmount).TruncateDecimal()

	//Send (mintedTokens - protocolTokens) to delegator address
	err = m.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, delegatorAddress,
		sdktypes.NewCoins(mintToken.Sub(protocolCoin)))
	if err != nil {
		return nil, types.ErrMintFailed
	}

	//Send protocol fee to protocol pool
	err = m.SendProtocolFee(ctx, sdktypes.NewCoins(protocolCoin), types.ModuleName, hostChainParams.PstakeParams.PstakeFeeAddress)
	if err != nil {
		return nil, types.ErrFailedDeposit
	}

	ctx.EventManager().EmitEvents(sdktypes.Events{
		sdktypes.NewEvent(
			types.EventTypeLiquidStake,
			sdktypes.NewAttribute(types.AttributeDelegatorAddress, delegatorAddress.String()),
			sdktypes.NewAttribute(types.AttributeAmount, mintToken.String()),
			sdktypes.NewAttribute(types.AttributeAmountRecieved, mintToken.Sub(protocolCoin).String()),
			sdktypes.NewAttribute(types.AttributePstakeDepositFee, protocolFee.String()),
		),
		sdktypes.NewEvent(
			sdktypes.EventTypeMessage,
			sdktypes.NewAttribute(sdktypes.AttributeKeyModule, types.AttributeValueCategory),
			sdktypes.NewAttribute(sdktypes.AttributeKeySender, msg.DelegatorAddress),
		)},
	)
	return &types.MsgLiquidStakeResponse{}, nil
}

func (m msgServer) Juice(goCtx context.Context, msg *types.MsgJuice) (*types.MsgJuiceResponse, error) {
	err := msg.ValidateBasic()
	if err != nil {
		return nil, types.ErrInvalidMessage
	}

	ctx := sdktypes.UnwrapSDKContext(goCtx)

	// sanity check for the arguments of message
	if ctx.IsZero() || !msg.Amount.IsValid() {
		return nil, types.ErrInvalidArgs
	}
	if !m.GetModuleState(ctx) {
		return nil, types.ErrModuleDisabled
	}
	//GetParams
	hostChainParams := m.GetHostChainParams(ctx)

	expectedIBCPrefix := ibctransfertypes.GetDenomPrefix(hostChainParams.TransferPort, hostChainParams.TransferChannel)

	denomTraceStr, err := m.ibcTransferKeeper.DenomPathFromHash(ctx, msg.Amount.Denom)
	if err != nil {
		return nil, err
	}
	denomTrace := ibctransfertypes.ParseDenomTrace(denomTraceStr)

	// Check if ibc path matches allowlisted path.
	if expectedIBCPrefix != denomTrace.GetPrefix() {
		return nil, types.ErrInvalidDenomPath
	}
	//Check if base denom is valid (uatom) , this can be programmed further to accommodate for liquid staked vouchers.
	if denomTrace.BaseDenom != hostChainParams.BaseDenom {
		return nil, types.ErrInvalidDenom
	}

	// check if address in message is correct or not
	rewarderAddress, err := sdktypes.AccAddressFromBech32(msg.RewarderAddress)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress
	}

	//send the rewards boost amount  to the deposit-module account
	rewardsBoostAmount := sdktypes.NewCoins(msg.Amount)
	err = m.SendTokensToRewardBoosterModuleAccount(ctx, rewardsBoostAmount, rewarderAddress)
	if err != nil {
		return nil, types.ErrFailedDeposit
	}

	ctx.EventManager().EmitEvents(sdktypes.Events{
		sdktypes.NewEvent(
			types.EventTypeRewardBoost,
			sdktypes.NewAttribute(types.AttributeRewarderAddress, rewarderAddress.String()),
			sdktypes.NewAttribute(types.AttributeAmountRecieved, rewardsBoostAmount.String()),
		),
		sdktypes.NewEvent(
			sdktypes.EventTypeMessage,
			sdktypes.NewAttribute(sdktypes.AttributeKeyModule, types.AttributeValueCategory),
			sdktypes.NewAttribute(sdktypes.AttributeKeySender, msg.RewarderAddress),
		)},
	)
	return &types.MsgJuiceResponse{}, nil
}

func (m msgServer) LiquidUnstake(goCtx context.Context, unstake *types.MsgLiquidUnstake) (*types.MsgLiquidUnstakeResponse, error) {
	// TODO implement this
	return nil, nil
}

func (m msgServer) Redeem(goCtx context.Context, msg *types.MsgRedeem) (*types.MsgRedeemResponse, error) {
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	ctx := sdktypes.UnwrapSDKContext(goCtx)

	// sanity check for the arguments of message
	if ctx.IsZero() || !msg.Amount.IsValid() {
		return nil, types.ErrInvalidArgs
	}
	if !m.GetModuleState(ctx) {
		return nil, types.ErrModuleDisabled
	}

	// take redeem address from msg address string
	redeemAddress, err := sdktypes.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress
	}

	// get the ibc denom and host chain params
	ibcDenom := m.GetIBCDenom(ctx)
	hostChainParams := m.GetHostChainParams(ctx)

	// check msg amount denom
	if msg.Amount.Denom != hostChainParams.MintDenom {
		return nil, types.ErrInvalidDenom
	}

	// We do not care about residue, as to not break Total calculation invariant.
	// protocolCoin is the redemption fee
	protocolCoin, _ := sdktypes.NewDecCoinFromDec(
		hostChainParams.MintDenom,
		hostChainParams.PstakeParams.PstakeRedemptionFee.MulInt(msg.Amount.Amount),
	).TruncateDecimal()

	// send redeem tokens to module account from redeem account
	err = m.bankKeeper.SendCoinsFromAccountToModule(ctx, redeemAddress, types.ModuleName, sdktypes.NewCoins(msg.Amount))
	if err != nil {
		return nil, err
	}

	// send protocol fee to protocol pool
	err = m.SendProtocolFee(ctx, sdktypes.NewCoins(protocolCoin), types.ModuleName, hostChainParams.PstakeParams.PstakeFeeAddress)
	if err != nil {
		return nil, types.ErrFailedDeposit
	}

	// convert redeem amount to ibc/whitelisted-denom amount (sub protocolCoin) based on the current c-value
	redeemStk := msg.Amount.Sub(protocolCoin)
	redeemToken, _ := m.ConvertStkToToken(ctx, sdktypes.NewDecCoinFromCoin(redeemStk), m.GetCValue(ctx))

	// get all deposit account balances
	allDepositBalances := m.bankKeeper.GetAllBalances(ctx, authtypes.NewModuleAddress(types.DepositModuleAccount))
	delegationBalance := sdktypes.NewCoin(ibcDenom, allDepositBalances.AmountOf(ibcDenom))

	// check deposit account has sufficient funds
	if redeemToken.IsGTE(delegationBalance) {
		return nil, sdkerrors.ErrInsufficientFunds
	}

	// send the ibc/Denom token from module to the account
	err = m.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.DepositModuleAccount, redeemAddress, sdktypes.NewCoins(redeemToken))
	if err != nil {
		return nil, err
	}

	// burn the redeemStk token
	err = m.bankKeeper.BurnCoins(ctx, types.ModuleName, sdktypes.NewCoins(redeemStk))
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdktypes.Events{
		sdktypes.NewEvent(
			types.EventTypeRedeem,
			sdktypes.NewAttribute(types.AttributeDelegatorAddress, redeemAddress.String()),
			sdktypes.NewAttribute(types.AttributeAmount, msg.Amount.String()),
			sdktypes.NewAttribute(types.AttributeAmountRecieved, redeemToken.String()),
			sdktypes.NewAttribute(types.AttributePstakeDepositFee, protocolCoin.String()),
		),
		sdktypes.NewEvent(
			sdktypes.EventTypeMessage,
			sdktypes.NewAttribute(sdktypes.AttributeKeyModule, types.AttributeValueCategory),
			sdktypes.NewAttribute(sdktypes.AttributeKeySender, msg.DelegatorAddress),
		)},
	)
	return &types.MsgRedeemResponse{}, nil
}
