package keeper

import (
	"context"
	"fmt"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	ibcporttypes "github.com/cosmos/ibc-go/v3/modules/core/05-port/types"
	host "github.com/cosmos/ibc-go/v3/modules/core/24-host"

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
	ctx := sdktypes.UnwrapSDKContext(goCtx)

	// sanity check for the arguments of message
	if ctx.IsZero() || !msg.Amount.IsValid() {
		return nil, sdkerrors.Wrapf(types.ErrInvalidArgs, "got invalid amount or ctx")
	}
	if !m.GetModuleState(ctx) {
		return nil, types.ErrModuleDisabled
	}
	//GetParams
	hostChainParams := m.GetHostChainParams(ctx)

	//check for minimum deposit amount
	if msg.Amount.Amount.LT(hostChainParams.MinDeposit) {
		return nil, sdkerrors.Wrapf(
			types.ErrMinDeposit, "expected amount more than %s, got %s", hostChainParams.MinDeposit, msg.Amount.Amount,
		)
	}

	expectedIBCPrefix := ibctransfertypes.GetDenomPrefix(hostChainParams.TransferPort, hostChainParams.TransferChannel)

	denomTraceStr, err := m.ibcTransferKeeper.DenomPathFromHash(ctx, msg.Amount.Denom)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInvalidDenom, "got error : %s", err)
	}
	denomTrace := ibctransfertypes.ParseDenomTrace(denomTraceStr)

	// Check if ibc path matches allowlisted path.
	if expectedIBCPrefix != denomTrace.GetPrefix() {
		return nil, sdkerrors.Wrapf(
			types.ErrInvalidDenomPath, "expected %s, got %s", expectedIBCPrefix, denomTrace.GetPrefix(),
		)
	}
	//Check if base denom is valid (uatom) , this can be programmed further to accommodate for liquid staked vouchers.
	if denomTrace.BaseDenom != hostChainParams.BaseDenom {
		return nil, sdkerrors.Wrapf(
			types.ErrInvalidDenom, "expected %s, got %s", hostChainParams.BaseDenom, denomTrace.BaseDenom,
		)
	}

	// get the delegator address from bech32 string
	delegatorAddress, err := sdktypes.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "got error : %s", err)
	}

	// amount of stk tokens to be minted. We calculate this before depositing any amount so as to not affect minting c-value.
	// We do not care about residue here because it won't be minted and bank.TotalSupply invariant should not be affected
	cValue := m.GetCValue(ctx)
	mintToken, _ := m.ConvertTokenToStk(ctx, sdktypes.NewDecCoinFromCoin(msg.Amount), cValue)

	//send the deposit to the deposit-module account
	depositAmount := sdktypes.NewCoins(msg.Amount)
	err = m.SendTokensToDepositModule(ctx, depositAmount, delegatorAddress)
	if err != nil {
		return nil, sdkerrors.Wrapf(
			types.ErrFailedDeposit, "failed to deposit tokens to module account %s, got error : %s", types.DepositModuleAccount, err,
		)
	}

	//Mint staked representative tokens in lscosmos module account
	err = m.bankKeeper.MintCoins(ctx, types.ModuleName, sdktypes.NewCoins(mintToken))
	if err != nil {
		return nil, sdkerrors.Wrapf(
			types.ErrMintFailed, "failed to mint coins in module %s, got error %s", types.ModuleName, err,
		)
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
		return nil, sdkerrors.Wrapf(
			types.ErrMintFailed, "failed to send coins from module %s to account %s, got error : %s",
			types.ModuleName, delegatorAddress.String(), err,
		)
	}

	//Send protocol fee to protocol pool
	if protocolCoin.IsPositive() {
		err = m.SendProtocolFee(ctx, sdktypes.NewCoins(protocolCoin), types.ModuleName, hostChainParams.PstakeParams.PstakeFeeAddress)
		if err != nil {
			return nil, sdkerrors.Wrapf(
				types.ErrFailedDeposit, "failed to send protocol fee to pstake fee address %s, got error : %s",
				hostChainParams.PstakeParams.PstakeFeeAddress, err,
			)
		}
	}
	ctx.EventManager().EmitEvents(sdktypes.Events{
		sdktypes.NewEvent(
			types.EventTypeLiquidStake,
			sdktypes.NewAttribute(types.AttributeDelegatorAddress, delegatorAddress.String()),
			sdktypes.NewAttribute(types.AttributeAmount, mintToken.String()),
			sdktypes.NewAttribute(types.AttributeAmountReceived, mintToken.Sub(protocolCoin).String()),
			sdktypes.NewAttribute(types.AttributePstakeDepositFee, protocolCoin.String()),
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
	ctx := sdktypes.UnwrapSDKContext(goCtx)

	// sanity check for the arguments of message
	if ctx.IsZero() || !msg.Amount.IsValid() {
		return nil, sdkerrors.Wrapf(types.ErrInvalidArgs, "got invalid amount or ctx")
	}
	if !m.GetModuleState(ctx) {
		return nil, types.ErrModuleDisabled
	}

	//GetParams
	hostChainParams := m.GetHostChainParams(ctx)

	expectedIBCPrefix := ibctransfertypes.GetDenomPrefix(hostChainParams.TransferPort, hostChainParams.TransferChannel)

	denomTraceStr, err := m.ibcTransferKeeper.DenomPathFromHash(ctx, msg.Amount.Denom)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInvalidDenom, "got error : %s", err)
	}
	denomTrace := ibctransfertypes.ParseDenomTrace(denomTraceStr)

	// Check if ibc path matches allowlisted path.
	if expectedIBCPrefix != denomTrace.GetPrefix() {
		return nil, sdkerrors.Wrapf(
			types.ErrInvalidDenomPath, "expected %s, got %s", expectedIBCPrefix, denomTrace.GetPrefix(),
		)
	}
	//Check if base denom is valid (uatom) , this can be programmed further to accommodate for liquid staked vouchers.
	if denomTrace.BaseDenom != hostChainParams.BaseDenom {
		return nil, sdkerrors.Wrapf(
			types.ErrInvalidDenom, "expected %s, got %s", hostChainParams.BaseDenom, denomTrace.BaseDenom,
		)
	}

	// check if address in message is correct or not
	rewarderAddress, err := sdktypes.AccAddressFromBech32(msg.RewarderAddress)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "incorrect address, got error : %s", err)
	}

	//send the rewards boost amount  to the deposit-module account
	rewardsBoostAmount := sdktypes.NewCoins(msg.Amount)
	err = m.SendTokensToRewardBoosterModuleAccount(ctx, rewardsBoostAmount, rewarderAddress)
	if err != nil {
		return nil, sdkerrors.Wrapf(
			types.ErrFailedDeposit, "failed to deposit tokens to module account %s, got error : %s", types.RewardBoosterModuleAccount, err,
		)
	}

	ctx.EventManager().EmitEvents(sdktypes.Events{
		sdktypes.NewEvent(
			types.EventTypeRewardBoost,
			sdktypes.NewAttribute(types.AttributeRewarderAddress, rewarderAddress.String()),
			sdktypes.NewAttribute(types.AttributeAmountReceived, rewardsBoostAmount.String()),
		),
		sdktypes.NewEvent(
			sdktypes.EventTypeMessage,
			sdktypes.NewAttribute(sdktypes.AttributeKeyModule, types.AttributeValueCategory),
			sdktypes.NewAttribute(sdktypes.AttributeKeySender, msg.RewarderAddress),
		)},
	)
	return &types.MsgJuiceResponse{}, nil
}

func (m msgServer) LiquidUnstake(goCtx context.Context, msg *types.MsgLiquidUnstake) (*types.MsgLiquidUnstakeResponse, error) {
	ctx := sdktypes.UnwrapSDKContext(goCtx)
	// sanity check for the arguments of message
	if ctx.IsZero() {
		return nil, types.ErrInvalidArgs
	}
	if !m.GetModuleState(ctx) {
		return nil, types.ErrModuleDisabled
	}

	hostChainParams := m.GetHostChainParams(ctx)

	if msg.Amount.Denom != hostChainParams.MintDenom {
		return nil, sdkerrors.Wrapf(types.ErrInvalidDenom, "Expected %s, got %s", hostChainParams.MintDenom, msg.Amount.Denom)
	}

	delegatorAddress, err := sdktypes.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, err
	}
	// take deposit into module acc
	err = m.bankKeeper.SendCoinsFromAccountToModule(ctx, delegatorAddress, types.UndelegationModuleAccount, sdktypes.NewCoins(msg.Amount))
	if err != nil {
		return nil, err
	}
	// take pstake fees
	unstakeCoin := msg.Amount
	pstakeFeeAmt := hostChainParams.PstakeParams.PstakeUnstakeFee.MulInt(msg.Amount.Amount).TruncateInt()
	pstakeFee := sdktypes.NewCoin(msg.Amount.Denom, pstakeFeeAmt)
	if pstakeFeeAmt.IsPositive() {
		err = m.SendProtocolFee(ctx, sdktypes.NewCoins(pstakeFee), types.UndelegationModuleAccount, hostChainParams.PstakeParams.PstakeFeeAddress)
		if err != nil {
			return nil, err
		}
		unstakeCoin = msg.Amount.Sub(pstakeFee)
	}

	// Add entry to unbonding db
	epoch := m.epochKeeper.GetEpochInfo(ctx, types.UndelegationEpochIdentifier)
	unbondingEpochNumber := types.CurrentUnbondingEpoch(epoch.CurrentEpoch)
	m.AddDelegatorUnbondingEpochEntry(ctx, delegatorAddress, unbondingEpochNumber, unstakeCoin)
	m.AddTotalUndelegationForEpoch(ctx, unbondingEpochNumber, unstakeCoin)

	// check is there are delegations worth the amount to be undelegated.
	// there are chances where the delegation epoch is not yet done so stkAtom are more than delegated amount
	// in this case users should just redeem tokens. (as tokens should be present as part of deposit tokens)
	delegationState := m.GetDelegationState(ctx)
	undelegations, err := m.GetHostAccountUndelegationForEpoch(ctx, unbondingEpochNumber)
	if err != nil {
		return nil, err
	}
	totalDelegations := delegationState.TotalDelegations(hostChainParams.BaseDenom)
	baseDenomUndelegations, _ := m.ConvertStkToToken(ctx, sdktypes.NewDecCoinFromCoin(undelegations.TotalUndelegationAmount), m.GetCValue(ctx))
	if totalDelegations.IsLT(sdktypes.NewCoin(hostChainParams.BaseDenom, baseDenomUndelegations.Amount)) {
		return nil, sdkerrors.Wrapf(types.ErrHostChainDelegationsLTUndelegations, "Delegated amount: %s is less than total undelegations for the epoch: %s", totalDelegations, undelegations.TotalUndelegationAmount)
	}

	ctx.EventManager().EmitEvents(sdktypes.Events{
		sdktypes.NewEvent(
			types.EventTypeLiquidUnstake,
			sdktypes.NewAttribute(types.AttributeDelegatorAddress, msg.GetDelegatorAddress()),
			sdktypes.NewAttribute(types.AttributeAmountReceived, msg.Amount.String()),
			sdktypes.NewAttribute(types.AttributePstakeUnstakeFee, pstakeFee.String()),
			sdktypes.NewAttribute(types.AttributeUnstakeAmount, unstakeCoin.String()),
		),
		sdktypes.NewEvent(
			sdktypes.EventTypeMessage,
			sdktypes.NewAttribute(sdktypes.AttributeKeyModule, types.AttributeValueCategory),
			sdktypes.NewAttribute(sdktypes.AttributeKeySender, msg.GetDelegatorAddress()),
		)},
	)
	return &types.MsgLiquidUnstakeResponse{}, nil
}

func (m msgServer) Redeem(goCtx context.Context, msg *types.MsgRedeem) (*types.MsgRedeemResponse, error) {
	ctx := sdktypes.UnwrapSDKContext(goCtx)

	// sanity check for the arguments of message
	if ctx.IsZero() || !msg.Amount.IsValid() {
		return nil, sdkerrors.Wrapf(types.ErrInvalidArgs, "got invalid amount or ctx")
	}
	if !m.GetModuleState(ctx) {
		return nil, types.ErrModuleDisabled
	}

	// take redeem address from msg address string
	redeemAddress, err := sdktypes.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "got error : %s", err)
	}

	// get the ibc denom and host chain params
	ibcDenom := m.GetIBCDenom(ctx)
	hostChainParams := m.GetHostChainParams(ctx)

	// check msg amount denom
	if msg.Amount.Denom != hostChainParams.MintDenom {
		return nil, sdkerrors.Wrapf(types.ErrInvalidDenom, "expected %s, got %s", hostChainParams.BaseDenom, msg.Amount.Denom)
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
		return nil, sdkerrors.Wrapf(
			types.ErrMintFailed, "failed to send coins from account %s to module %s, got error : %s",
			redeemAddress.String(), types.ModuleName, err,
		)
	}

	// send protocol fee to protocol pool
	if protocolCoin.IsPositive() {
		err = m.SendProtocolFee(ctx, sdktypes.NewCoins(protocolCoin), types.ModuleName, hostChainParams.PstakeParams.PstakeFeeAddress)
		if err != nil {
			return nil, sdkerrors.Wrapf(
				types.ErrFailedDeposit, "failed to send protocol fee to pstake fee address %s, got error : %s",
				hostChainParams.PstakeParams.PstakeFeeAddress, err,
			)
		}
	}
	// convert redeem amount to ibc/whitelisted-denom amount (sub protocolCoin) based on the current c-value
	redeemStk := msg.Amount.Sub(protocolCoin)
	redeemToken, _ := m.ConvertStkToToken(ctx, sdktypes.NewDecCoinFromCoin(redeemStk), m.GetCValue(ctx))

	// get all deposit account balances
	allDepositBalances := m.bankKeeper.GetAllBalances(ctx, authtypes.NewModuleAddress(types.DepositModuleAccount))
	delegationBalance := sdktypes.NewCoin(ibcDenom, allDepositBalances.AmountOf(ibcDenom))

	// check deposit account has sufficient funds
	if redeemToken.IsGTE(delegationBalance) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, "expected tokens under %s, got %s for redeem", delegationBalance.String(), redeemToken.String())
	}

	// send the ibc/Denom token from module to the account
	err = m.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.DepositModuleAccount, redeemAddress, sdktypes.NewCoins(redeemToken))
	if err != nil {
		return nil, sdkerrors.Wrapf(
			types.ErrMintFailed, "failed to send coins from module %s to account %s, got error : %s",
			types.DepositModuleAccount, redeemAddress.String(), err,
		)
	}

	// burn the redeemStk token
	err = m.bankKeeper.BurnCoins(ctx, types.ModuleName, sdktypes.NewCoins(redeemStk))
	if err != nil {
		return nil, sdkerrors.Wrapf(
			types.ErrBurnFailed, "failed to burn coins from module %s, got error %s", types.ModuleName, err,
		)
	}

	ctx.EventManager().EmitEvents(sdktypes.Events{
		sdktypes.NewEvent(
			types.EventTypeRedeem,
			sdktypes.NewAttribute(types.AttributeDelegatorAddress, redeemAddress.String()),
			sdktypes.NewAttribute(types.AttributeAmount, msg.Amount.String()),
			sdktypes.NewAttribute(types.AttributeAmountReceived, redeemToken.String()),
			sdktypes.NewAttribute(types.AttributePstakeRedeemFee, protocolCoin.String()),
		),
		sdktypes.NewEvent(
			sdktypes.EventTypeMessage,
			sdktypes.NewAttribute(sdktypes.AttributeKeyModule, types.AttributeValueCategory),
			sdktypes.NewAttribute(sdktypes.AttributeKeySender, msg.DelegatorAddress),
		)},
	)
	return &types.MsgRedeemResponse{}, nil
}

func (m msgServer) Claim(goCtx context.Context, msg *types.MsgClaim) (*types.MsgClaimResponse, error) {
	ctx := sdktypes.UnwrapSDKContext(goCtx)
	// sanity check for the arguments of message
	if ctx.IsZero() {
		return nil, types.ErrInvalidArgs
	}
	if !m.GetModuleState(ctx) {
		return nil, types.ErrModuleDisabled
	}

	// get AccAddress from bech32 string
	delegatorAddress, err := sdktypes.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	// get all the entries corresponding to the delegator address
	delegatorUnbondingEntries := m.IterateDelegatorUnbondingEpochEntry(ctx, delegatorAddress)

	// loop through all the epoch and send tokens if an entry has matured.
	for _, unbondingEntry := range delegatorUnbondingEntries {
		unbondingEpochCValue := m.GetUnbondingEpochCValue(ctx, unbondingEntry.EpochNumber)
		if unbondingEpochCValue.IsMatured {
			// get c value from the UnbondingEpochCValue struct
			// calculate claimable amount from un inverse c value
			claimableAmount := unbondingEntry.Amount.Amount.ToDec().Quo(unbondingEpochCValue.GetUnbondingEpochCValue())

			// calculate claimable coin and community coin to be sent to delegator account and community pool respectively
			claimableCoin, _ := sdktypes.NewDecCoinFromDec(m.GetIBCDenom(ctx), claimableAmount).TruncateDecimal()

			// send coin to delegator address from undelegation module account
			err = m.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.UndelegationModuleAccount, delegatorAddress, sdktypes.NewCoins(claimableCoin))
			if err != nil {
				return nil, err
			}

			ctx.EventManager().EmitEvents(sdktypes.Events{
				sdktypes.NewEvent(
					types.EventTypeClaim,
					sdktypes.NewAttribute(types.AttributeDelegatorAddress, delegatorAddress.String()),
					sdktypes.NewAttribute(types.AttributeAmount, unbondingEntry.Amount.String()),
					sdktypes.NewAttribute(types.AttributeClaimedAmount, claimableAmount.String()),
				)},
			)

			// remove entry from unbonding epoch entry
			m.RemoveDelegatorUnbondingEpochEntry(ctx, delegatorAddress, unbondingEntry.EpochNumber)
		}
		if unbondingEpochCValue.IsFailed {
			err = m.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.UndelegationModuleAccount, delegatorAddress, sdktypes.NewCoins(unbondingEntry.Amount))
			if err != nil {
				return nil, err
			}

			// remove entry from unbonding epoch entry
			m.RemoveDelegatorUnbondingEpochEntry(ctx, delegatorAddress, unbondingEntry.EpochNumber)
		}
	}

	// emit event
	ctx.EventManager().EmitEvents(sdktypes.Events{
		sdktypes.NewEvent(
			sdktypes.EventTypeMessage,
			sdktypes.NewAttribute(sdktypes.AttributeKeyModule, types.AttributeValueCategory),
			sdktypes.NewAttribute(sdktypes.AttributeKeySender, msg.DelegatorAddress),
		)},
	)

	return &types.MsgClaimResponse{}, nil
}

func (m msgServer) JumpStart(goCtx context.Context, msg *types.MsgJumpStart) (*types.MsgJumpStartResponse, error) {
	ctx := sdktypes.UnwrapSDKContext(goCtx)
	// sanity check for the arguments of message
	if ctx.IsZero() {
		return nil, types.ErrInvalidArgs
	}

	// check pstake fee address == from addr
	hostChainParams := m.GetHostChainParams(ctx)
	if msg.PstakeAddress != hostChainParams.PstakeParams.PstakeFeeAddress {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("msg.pstakeAddress should be equal to msg.PstakeParams.PstakeFeeAddress, got %s expected %s", msg.PstakeAddress, hostChainParams.PstakeParams.PstakeFeeAddress))
	}
	// check module disabled
	if m.GetModuleState(ctx) {
		return nil, types.ErrModuleAlreadyEnabled
	}
	// reset db, no need to release capability
	m.SetDelegationState(ctx, types.DelegationState{})
	m.SetHostChainRewardAddress(ctx, types.HostChainRewardAddress{})
	if err := msg.HostAccounts.Validate(); err != nil {
		return nil, err
	}
	m.SetHostAccounts(ctx, msg.HostAccounts)
	// do proposal things
	if msg.TransferPort != ibctransfertypes.PortID {
		return nil, sdkerrors.Wrap(ibcporttypes.ErrInvalidPort, "Only acceptable TransferPort is \"transfer\"")
	}

	// checks for valid and active channel
	channel, found := m.channelKeeper.GetChannel(ctx, msg.TransferPort, msg.TransferChannel)
	if !found {
		return nil, sdkerrors.Wrap(ibcchanneltypes.ErrChannelNotFound, fmt.Sprintf("channel for ibc transfer: %s not found", msg.TransferChannel))
	}
	if channel.State != ibcchanneltypes.OPEN {
		return nil, sdkerrors.Wrapf(
			ibcchanneltypes.ErrInvalidChannelState,
			"channel state is not OPEN (got %s)", channel.State.String(),
		)
	}
	// TODO Understand capabilities and see if it has to be/ should be claimed in lsscopedkeeper. If it even matters.
	_, err := m.lscosmosScopedKeeper.NewCapability(ctx, host.ChannelCapabilityPath(msg.TransferPort, msg.TransferChannel))
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "Failed to create and claim capability for ibc transfer port and channel")
	}

	hostAccounts := m.GetHostAccounts(ctx)
	// This checks for channel being active
	err = m.icaControllerKeeper.RegisterInterchainAccount(ctx, msg.ConnectionID, hostAccounts.DelegatorAccountOwnerID)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "Could not register ica delegation Address")
	}

	newHostChainParams := types.NewHostChainParams(msg.ChainID, msg.ConnectionID, msg.TransferChannel,
		msg.TransferPort, msg.BaseDenom, msg.MintDenom, msg.PstakeParams.PstakeFeeAddress,
		msg.MinDeposit, msg.PstakeParams.PstakeDepositFee, msg.PstakeParams.PstakeRestakeFee,
		msg.PstakeParams.PstakeUnstakeFee, msg.PstakeParams.PstakeRedemptionFee)

	m.SetHostChainParams(ctx, newHostChainParams)

	if !msg.AllowListedValidators.Valid() {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Allow listed validators is invalid")
	}
	m.SetAllowListedValidators(ctx, msg.AllowListedValidators)

	ctx.EventManager().EmitEvents(sdktypes.Events{
		sdktypes.NewEvent(
			types.EventTypeJumpStart,
			sdktypes.NewAttribute(types.AttributePstakeAddress, msg.PstakeAddress),
		),
		sdktypes.NewEvent(
			sdktypes.EventTypeMessage,
			sdktypes.NewAttribute(sdktypes.AttributeKeyModule, types.AttributeValueCategory),
			sdktypes.NewAttribute(sdktypes.AttributeKeySender, msg.PstakeAddress),
		)},
	)

	return &types.MsgJumpStartResponse{}, nil
}
