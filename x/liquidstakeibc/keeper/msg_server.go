package keeper

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

const (
	KeyValidatorWeight string = "validator_weight"
	KeyDepositFee      string = "deposit_fee"
	KeyRestakeFee      string = "restake_fee"
	KeyUnstakeFee      string = "unstake_fee"
	KeyRedemptionFee   string = "redemption_fee"
	KeyMinimumDeposit  string = "min_deposit"
	KeyActive          string = "active"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the liquidstakeibc MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// RegisterHostChain adds a new host chain to the protocol
func (k msgServer) RegisterHostChain(
	goCtx context.Context,
	msg *types.MsgRegisterHostChain,
) (*types.MsgRegisterHostChainResponse, error) {
	ctx := sdktypes.UnwrapSDKContext(goCtx)

	// authority needs to be either the gov module account (for proposals)
	// or the module admin account (for normal txs)
	if msg.Authority != k.authority && msg.Authority != k.GetParams(ctx).AdminAddress {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "tx signer is not a module authority")
	}

	// get the host chain id
	chainID, err := k.GetChainID(ctx, msg.ConnectionId)
	if err != nil {
		return nil, fmt.Errorf("chain id not found for connection \"%s\": \"%w\"", msg.ConnectionId, err)
	}

	// build the host chain params
	hostChainParams := &types.HostChainLSParams{
		DepositFee:    msg.DepositFee,
		RestakeFee:    msg.RestakeFee,
		UnstakeFee:    msg.UnstakeFee,
		RedemptionFee: msg.RedemptionFee,
	}

	hc := &types.HostChain{
		ChainId:         chainID,
		ConnectionId:    msg.ConnectionId,
		ChannelId:       msg.ChannelId,
		PortId:          msg.PortId,
		Params:          hostChainParams,
		HostDenom:       msg.HostDenom,
		MinimumDeposit:  msg.MinimumDeposit,
		CValue:          sdktypes.NewDec(1),
		NextValsetHash:  []byte{},
		UnbondingFactor: msg.UnbondingFactor,
		Active:          false,
	}

	// save the host chain
	k.SetHostChain(ctx, hc)

	// register delegate ICA
	if err = k.RegisterICAAccount(ctx, hc.ConnectionId, k.DelegateAccountPortOwner(chainID)); err != nil {
		return nil, errorsmod.Wrapf(
			types.ErrRegisterFailed,
			"error registering %s delegate ica: %s",
			chainID,
			err.Error(),
		)
	}

	// register reward ICA
	if err = k.RegisterICAAccount(ctx, hc.ConnectionId, k.RewardsAccountPortOwner(chainID)); err != nil {
		return nil, errorsmod.Wrapf(
			types.ErrRegisterFailed,
			"error registering %s reward ica: %s",
			chainID,
			err.Error(),
		)
	}

	// query the host chain for the validator set
	if err := k.QueryHostChainValidators(ctx, hc, stakingtypes.QueryValidatorsRequest{}); err != nil {
		return nil, errorsmod.Wrapf(
			types.ErrFailedICQRequest,
			"error submitting validators icq: %s",
			err.Error(),
		)
	}

	return &types.MsgRegisterHostChainResponse{}, nil
}

// UpdateHostChain updates a registered host chain
func (k msgServer) UpdateHostChain(
	goCtx context.Context,
	msg *types.MsgUpdateHostChain,
) (*types.MsgUpdateHostChainResponse, error) {
	ctx := sdktypes.UnwrapSDKContext(goCtx)

	// authority needs to be either the gov module account (for proposals)
	// or the module admin account (for normal txs)
	if msg.Authority != k.authority && msg.Authority != k.GetParams(ctx).AdminAddress {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "tx signer is not a module authority")
	}

	hc, found := k.GetHostChain(ctx, msg.ChainId)
	if !found {
		return nil, fmt.Errorf("invalid chain id \"%s\", host chain is not registered", msg.ChainId)
	}

	for _, update := range msg.Updates {
		switch update.Key {
		case KeyValidatorWeight:
			validator, weight, found := strings.Cut(update.Value, ",")
			if !found {
				return nil, fmt.Errorf("unable to parse validator update string")
			}

			if err := k.UpdateHostChainValidatorWeight(ctx, hc, validator, weight); err != nil {
				return nil, fmt.Errorf("invalid validator weight update values: %v", err)
			}
		case KeyDepositFee:
			fee, err := sdktypes.NewDecFromStr(update.Value)
			if err != nil {
				return nil, fmt.Errorf("unable to parse string to sdk.Dec: %w", err)
			}

			hc.Params.DepositFee = fee
			if fee.LT(sdktypes.NewDec(0)) {
				return nil, fmt.Errorf("invalid deposit fee value, less than zero")
			}
		case KeyRestakeFee:
			fee, err := sdktypes.NewDecFromStr(update.Value)
			if err != nil {
				return nil, fmt.Errorf("unable to parse string to sdk.Dec: %w", err)
			}

			hc.Params.RestakeFee = fee
			if fee.LT(sdktypes.NewDec(0)) {
				return nil, fmt.Errorf("invalid deposit fee value, less than zero")
			}
		case KeyRedemptionFee:
			fee, err := sdktypes.NewDecFromStr(update.Value)
			if err != nil {
				return nil, fmt.Errorf("unable to parse string to sdk.Dec: %w", err)
			}

			hc.Params.RedemptionFee = fee
			if fee.LT(sdktypes.NewDec(0)) {
				return nil, fmt.Errorf("invalid deposit fee value, less than zero")
			}
		case KeyUnstakeFee:
			fee, err := sdktypes.NewDecFromStr(update.Value)
			if err != nil {
				return nil, fmt.Errorf("unable to parse string to sdk.Dec: %w", err)
			}

			hc.Params.UnstakeFee = fee
			if fee.LT(sdktypes.NewDec(0)) {
				return nil, fmt.Errorf("invalid deposit fee value, less than zero")
			}
		case KeyMinimumDeposit:
			minimumDeposit, ok := sdktypes.NewIntFromString(update.Value)
			if !ok {
				return nil, fmt.Errorf("unable to parse string to sdk.Int")
			}

			hc.MinimumDeposit = minimumDeposit
			if minimumDeposit.LT(sdktypes.NewInt(0)) {
				return nil, fmt.Errorf("invalid minimum deposit value less than zero")
			}
		case KeyActive:
			active, err := strconv.ParseBool(update.Value)
			if err != nil {
				return nil, fmt.Errorf("unable to parse string to bool")
			}

			hc.Active = active
		default:
			return nil, fmt.Errorf("invalid or unexpected update key: %s", update.Key)
		}
	}

	k.SetHostChain(ctx, hc)

	return &types.MsgUpdateHostChainResponse{}, nil
}

// LiquidStake defines a method for liquid staking tokens
func (k msgServer) LiquidStake(
	goCtx context.Context,
	msg *types.MsgLiquidStake,
) (*types.MsgLiquidStakeResponse, error) {
	ctx := sdktypes.UnwrapSDKContext(goCtx)

	// retrieve the host chain
	hostChain, found := k.GetHostChainFromIbcDenom(ctx, msg.Amount.Denom)
	if !found {
		return nil, errorsmod.Wrapf(
			types.ErrInvalidHostChain,
			"host chain with ibc denom %s not registered",
			msg.Amount.Denom,
		)
	}

	// check for minimum deposit amount
	if msg.Amount.Amount.LT(hostChain.MinimumDeposit) {
		return nil, errorsmod.Wrapf(
			types.ErrMinDeposit,
			"expected amount more than %s, got %s",
			hostChain.MinimumDeposit,
			msg.Amount.Amount,
		)
	}

	// get the delegator address from the bech32 string
	delegatorAddress, err := sdktypes.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "error parsing delegator address: %s", err)
	}

	// amount of stk tokens to be minted
	mintDenom := hostChain.MintDenom()
	mintAmount := sdktypes.NewDecCoinFromCoin(msg.Amount).Amount.Mul(hostChain.CValue)
	mintToken, _ := sdktypes.NewDecCoinFromDec(mintDenom, mintAmount).TruncateDecimal()

	// send the deposit to the deposit-module account
	depositAmount := sdktypes.NewCoins(msg.Amount)
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, delegatorAddress, types.DepositModuleAccount, depositAmount)
	if err != nil {
		return nil, errorsmod.Wrapf(
			types.ErrFailedDeposit,
			"failed to deposit tokens to module account %s: %s",
			types.DepositModuleAccount,
			err,
		)
	}

	// add the deposit amount to the deposit record for that chain/epoch
	currentEpoch := k.GetEpochNumber(ctx, types.DelegationEpoch)
	deposit, found := k.GetDepositForChainAndEpoch(ctx, hostChain.ChainId, currentEpoch)
	if !found {
		return nil, errorsmod.Wrapf(
			types.ErrDepositNotFound,
			"deposit not found for chain %s and epoch %v",
			hostChain.ChainId,
			currentEpoch,
		)
	}
	deposit.Amount.Amount = deposit.Amount.Amount.Add(msg.Amount.Amount)
	k.SetDeposit(ctx, deposit)

	// mint stk tokens in the module account
	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, sdktypes.NewCoins(mintToken))
	if err != nil {
		return nil, errorsmod.Wrapf(
			types.ErrMintFailed,
			"failed to mint coins in module %s: %s",
			types.ModuleName, err,
		)
	}

	// calculate protocol fee
	protocolFeeAmount := hostChain.Params.DepositFee.MulInt(mintToken.Amount)
	protocolFee, _ := sdktypes.NewDecCoinFromDec(mintDenom, protocolFeeAmount).TruncateDecimal()

	// send stk tokens to the delegator address
	err = k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx,
		types.ModuleName,
		delegatorAddress,
		sdktypes.NewCoins(mintToken.Sub(protocolFee)),
	)
	if err != nil {
		return nil, errorsmod.Wrapf(
			types.ErrMintFailed,
			"failed to send coins from module %s to account %s: %s",
			types.ModuleName,
			delegatorAddress.String(),
			err,
		)
	}

	// retrieve the module params
	params := k.GetParams(ctx)

	// send the protocol fee to the protocol pool
	if protocolFee.IsPositive() {
		err = k.SendProtocolFee(ctx, sdktypes.NewCoins(protocolFee), types.ModuleName, params.FeeAddress)
		if err != nil {
			return nil, errorsmod.Wrapf(
				types.ErrFailedDeposit,
				"failed to send protocol fee to pStake fee address %s: %s",
				params.FeeAddress,
				err,
			)
		}
	}
	ctx.EventManager().EmitEvents(sdktypes.Events{
		sdktypes.NewEvent(
			types.EventTypeLiquidStake,
			sdktypes.NewAttribute(types.AttributeDelegatorAddress, delegatorAddress.String()),
			sdktypes.NewAttribute(types.AttributeAmount, mintToken.String()),
			sdktypes.NewAttribute(types.AttributeAmountReceived, mintToken.Sub(protocolFee).String()),
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

// LiquidUnstake defines a method for unstaking liquid staked tokens
func (k msgServer) LiquidUnstake(
	goCtx context.Context,
	msg *types.MsgLiquidUnstake,
) (*types.MsgLiquidUnstakeResponse, error) {
	ctx := sdktypes.UnwrapSDKContext(goCtx)

	// get the host chain we need to unstake from
	hc, found := k.GetHostChainFromHostDenom(ctx, msg.HostDenom)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrInvalidHostChain,
			"host chain with host denom %s not registered",
			msg.HostDenom,
		)
	}

	// check if the message amount has the correct denom
	if msg.Amount.Denom != hc.MintDenom() {
		return nil, errorsmod.Wrapf(types.ErrInvalidDenom,
			"expected %s, got %s",
			hc.MintDenom(),
			msg.Amount.Denom,
		)
	}

	// parse the delegator address
	delegatorAddress, err := sdktypes.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	// send the tokens from the delegator address to the undelegation module account
	err = k.bankKeeper.SendCoinsFromAccountToModule(
		ctx,
		delegatorAddress,
		types.UndelegationModuleAccount,
		sdktypes.NewCoins(msg.Amount),
	)
	if err != nil {
		return nil, err
	}

	// send the unstake fee to the module fee address and subtract it from the total to unstake
	unstakeAmount := msg.Amount
	feeAmount := hc.Params.UnstakeFee.MulInt(unstakeAmount.Amount).TruncateInt()
	if feeAmount.IsPositive() {
		fee := sdktypes.NewCoin(msg.Amount.Denom, feeAmount)

		err = k.SendProtocolFee(
			ctx,
			sdktypes.NewCoins(fee),
			types.UndelegationModuleAccount,
			k.GetParams(ctx).FeeAddress)
		if err != nil {
			return nil, err
		}

		unstakeAmount = msg.Amount.Sub(fee)
	}

	// calculate the host chain token unbond amount from the stk amount
	decTokenAmount := sdktypes.NewDecCoinFromCoin(unstakeAmount).Amount.Mul(sdktypes.OneDec().Quo(hc.CValue))
	tokenAmount, _ := sdktypes.NewDecCoinFromDec(hc.HostDenom, decTokenAmount).TruncateDecimal()
	unbondAmount := sdktypes.NewCoin(hc.HostDenom, tokenAmount.Amount)

	// calculate the current unbonding epoch
	epoch := k.epochsKeeper.GetEpochInfo(ctx, types.UndelegationEpoch)
	unbondingEpoch := types.CurrentUnbondingEpoch(hc.UnbondingFactor, epoch.CurrentEpoch)

	// increase the unbonding value for the epoch both for the user record and the module record
	k.IncreaseUserUnbondingAmountForEpoch(ctx, hc.ChainId, msg.DelegatorAddress, unbondingEpoch, unstakeAmount, unbondAmount)
	k.IncreaseUndelegatingAmountForEpoch(ctx, hc.ChainId, unbondingEpoch, unstakeAmount, unbondAmount)

	// check if the total unbonding amount for the next unbonding epoch is less than what is currently staked
	totalUnbondings, _ := k.GetUnbonding(ctx, hc.ChainId, unbondingEpoch)
	totalDelegations := hc.GetHostChainTotalDelegations()
	if totalDelegations.LT(unbondAmount.Amount) {
		return nil, errorsmod.Wrapf(
			types.ErrNotEnoughDelegations,
			"delegated amount %s is less than the total undelegation %s for epoch %d",
			totalDelegations,
			totalUnbondings,
			unbondingEpoch,
		)
	}

	ctx.EventManager().EmitEvents(sdktypes.Events{
		sdktypes.NewEvent(
			types.EventTypeLiquidUnstake,
			sdktypes.NewAttribute(types.AttributeDelegatorAddress, msg.GetDelegatorAddress()),
			sdktypes.NewAttribute(types.AttributeAmountReceived, msg.Amount.String()),
			sdktypes.NewAttribute(types.AttributePstakeUnstakeFee, feeAmount.String()),
			sdktypes.NewAttribute(types.AttributeUnstakeAmount, unbondAmount.String()),
			sdktypes.NewAttribute(types.AttributeUnstakeEpoch, strconv.FormatInt(unbondingEpoch, 10)),
		),
		sdktypes.NewEvent(
			sdktypes.EventTypeMessage,
			sdktypes.NewAttribute(sdktypes.AttributeKeyModule, types.AttributeValueCategory),
			sdktypes.NewAttribute(sdktypes.AttributeKeySender, msg.GetDelegatorAddress()),
		)},
	)

	return &types.MsgLiquidUnstakeResponse{}, nil
}

// Redeem defines a method for instantly redeem liquid staked tokens
func (k msgServer) Redeem(
	goCtx context.Context,
	msg *types.MsgRedeem,
) (*types.MsgRedeemResponse, error) {
	ctx := sdktypes.UnwrapSDKContext(goCtx)

	// get the host chain we need to unstake from
	hc, found := k.GetHostChainFromHostDenom(ctx, msg.HostDenom)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrInvalidHostChain,
			"host chain with host denom %s not registered",
			msg.HostDenom,
		)
	}

	// check the msg amount denom is the host chain mint denom
	if msg.Amount.Denom != hc.MintDenom() {
		return nil, errorsmod.Wrapf(
			types.ErrInvalidDenom,
			"expected %s, got %s",
			hc.MintDenom(),
			msg.Amount.Denom,
		)
	}

	// get the redeem address
	redeemAddress, err := sdktypes.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "got error : %s", err)
	}

	// send the redeem amount to the module account
	err = k.bankKeeper.SendCoinsFromAccountToModule(
		ctx,
		redeemAddress,
		types.ModuleName,
		sdktypes.NewCoins(msg.Amount))
	if err != nil {
		return nil, errorsmod.Wrapf(
			types.ErrMintFailed,
			"failed to send instant redeemed coins from account %s to module %s: %s",
			redeemAddress.String(),
			types.ModuleName,
			err.Error(),
		)
	}

	// calculate the instant redemption fee
	fee, _ := sdktypes.NewDecCoinFromDec(
		hc.MintDenom(),
		hc.Params.RedemptionFee.MulInt(msg.Amount.Amount),
	).TruncateDecimal()

	// send the protocol fee to the module fee address
	if fee.IsPositive() {
		err = k.SendProtocolFee(
			ctx,
			sdktypes.NewCoins(fee),
			types.ModuleName,
			k.GetParams(ctx).FeeAddress,
		)
		if err != nil {
			return nil, errorsmod.Wrapf(
				types.ErrFailedDeposit,
				"failed to send instant redemption fee to module fee address %s: %s",
				k.GetParams(ctx).FeeAddress,
				err.Error(),
			)
		}
	}

	// amount of tokens to be redeemed
	stkAmount := msg.Amount.Sub(fee)
	redeemAmount := sdktypes.NewDecCoinFromCoin(stkAmount).Amount.Mul(hc.CValue)
	redeemToken, _ := sdktypes.NewDecCoinFromDec(hc.IBCDenom(), redeemAmount).TruncateDecimal()

	// check if there is enough deposits to fulfill the instant redemption request
	depositAccountBalance := k.bankKeeper.GetBalance(
		ctx,
		authtypes.NewModuleAddress(types.DepositModuleAccount),
		hc.IBCDenom(),
	)
	if redeemToken.IsGTE(depositAccountBalance) {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrInsufficientFunds,
			"can't instant redeem %s tokens, only %s is available",
			redeemToken.String(),
			depositAccountBalance.Amount.String(),
		)
	}

	// subtract the redemption amount from the deposits
	if err := k.AdjustDepositsForRedemption(ctx, hc, redeemToken); err != nil {
		return nil, errorsmod.Wrapf(
			types.ErrRedeemFailed,
			"could not adjust current deposits for redemption",
		)
	}

	// send the instant redeemed token from module to the account
	err = k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx,
		types.DepositModuleAccount,
		redeemAddress,
		sdktypes.NewCoins(redeemToken),
	)
	if err != nil {
		return nil, errorsmod.Wrapf(
			types.ErrRedeemFailed,
			"failed to send instant redeemed coins from module %s to account %s: %s",
			types.DepositModuleAccount,
			redeemAddress.String(),
			err.Error(),
		)
	}

	// burn the stk tokens
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdktypes.NewCoins(stkAmount))
	if err != nil {
		return nil, errorsmod.Wrapf(
			types.ErrBurnFailed,
			"failed to burn instant redeemed coins on module %s: %s",
			types.ModuleName,
			err.Error(),
		)
	}

	ctx.EventManager().EmitEvents(sdktypes.Events{
		sdktypes.NewEvent(
			types.EventTypeRedeem,
			sdktypes.NewAttribute(types.AttributeDelegatorAddress, redeemAddress.String()),
			sdktypes.NewAttribute(types.AttributeAmount, msg.Amount.String()),
			sdktypes.NewAttribute(types.AttributeAmountReceived, redeemToken.String()),
			sdktypes.NewAttribute(types.AttributePstakeRedeemFee, fee.String()),
		),
		sdktypes.NewEvent(
			sdktypes.EventTypeMessage,
			sdktypes.NewAttribute(sdktypes.AttributeKeyModule, types.AttributeValueCategory),
			sdktypes.NewAttribute(sdktypes.AttributeKeySender, msg.DelegatorAddress),
		)},
	)

	return &types.MsgRedeemResponse{}, nil
}
