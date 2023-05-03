package keeper

import (
	"context"
	"fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
	// check if the message authority is the module authority (normally the gov account)
	if k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "expected %s got %s", k.authority, msg.Authority)
	}

	// unwrap context
	ctx := sdktypes.UnwrapSDKContext(goCtx)

	// get the host chain id
	chainId, err := k.GetChainID(ctx, msg.ConnectionId)
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
		ChainId:        chainId,
		ConnectionId:   msg.ConnectionId,
		ChannelId:      msg.ChannelId,
		PortId:         msg.PortId,
		Params:         hostChainParams,
		HostDenom:      msg.HostDenom,
		MinimumDeposit: msg.MinimumDeposit,
		CValue:         sdktypes.NewDec(1),
	}

	// save the host chain
	k.SetHostChain(ctx, hc)

	// register delegate ICA
	delegateAccount := chainId + "." + types.DelegateICAType
	if err = k.RegisterICAAccount(ctx, hc.ConnectionId, delegateAccount); err != nil {
		return nil, errorsmod.Wrapf(types.ErrRegisterFailed, "error registering %s delegate ica: %w", chainId, err)
	}

	// register reward ICA
	rewardAccount := chainId + "." + types.RewardsICAType
	if err = k.RegisterICAAccount(ctx, hc.ConnectionId, rewardAccount); err != nil {
		return nil, errorsmod.Wrapf(types.ErrRegisterFailed, "error registering %s reward ica: %w", chainId, err)
	}

	// query the host chain for the validator set
	if err := k.QueryHostChainValidators(ctx, hc, stakingtypes.QueryValidatorsRequest{}); err != nil {
		return nil, errorsmod.Wrapf(types.ErrFailedICQRequest, "error submitting validators icq: %w", err)
	}

	return &types.MsgRegisterHostChainResponse{}, nil
}

// UpdateHostChain updates a registered host chain
func (k msgServer) UpdateHostChain(
	goCtx context.Context,
	msg *types.MsgUpdateHostChain,
) (*types.MsgUpdateHostChainResponse, error) {
	if k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "expected %s got %s", k.authority, msg.Authority)
	}

	ctx := sdktypes.UnwrapSDKContext(goCtx)

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

	// update the host chain c value
	hostChain.CValue = k.GetHostChainCValue(ctx, hostChain)
	k.SetHostChain(ctx, hostChain)

	// add the deposit amount to the deposit record for that chain/epoch
	currentEpoch := k.GetEpochNumber(ctx, types.DelegationEpoch)
	deposit, found := k.GetDepositForChainAndEpoch(ctx, hostChain.ChainId, currentEpoch)
	if !found {
		return nil, errorsmod.Wrapf(
			types.ErrDepositNotFound,
			"deposit not found for chain %s and epoch %s",
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
