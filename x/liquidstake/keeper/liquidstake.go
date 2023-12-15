package keeper

import (
	"encoding/json"
	"time"

	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/persistenceOne/pstake-native/v2/x/liquidstake/types"
)

func (k Keeper) LiquidBondDenom(ctx sdk.Context) string {
	return k.GetParams(ctx).LiquidBondDenom
}

// GetNetAmountState calculates the sum of bondedDenom balance, total delegation tokens(slash applied LiquidTokens), total remaining reward of types.LiquidStakeProxyAcc
// During liquid unstaking, stkxprt immediately burns and the unbonding queue belongs to the requester, so the liquid staker's unbonding values are excluded on netAmount
// It is used only for calculation and query and is not stored in kv.
func (k Keeper) GetNetAmountState(ctx sdk.Context) (nas types.NetAmountState) {
	totalRemainingRewards, totalDelShares, totalLiquidTokens := k.CheckDelegationStates(ctx, types.LiquidStakeProxyAcc)

	totalUnbondingBalance := sdk.ZeroInt()
	ubds := k.stakingKeeper.GetAllUnbondingDelegations(ctx, types.LiquidStakeProxyAcc)
	for _, ubd := range ubds {
		for _, entry := range ubd.Entries {
			// use Balance(slashing applied) not InitialBalance(without slashing)
			totalUnbondingBalance = totalUnbondingBalance.Add(entry.Balance)
		}
	}

	nas = types.NetAmountState{
		StkxprtTotalSupply:    k.bankKeeper.GetSupply(ctx, k.LiquidBondDenom(ctx)).Amount,
		TotalDelShares:        totalDelShares,
		TotalLiquidTokens:     totalLiquidTokens,
		TotalRemainingRewards: totalRemainingRewards,
		TotalUnbondingBalance: totalUnbondingBalance,
		ProxyAccBalance:       k.GetProxyAccBalance(ctx, types.LiquidStakeProxyAcc).Amount,
	}

	nas.NetAmount = nas.CalcNetAmount()
	nas.MintRate = nas.CalcMintRate()
	return
}

// LiquidStake mints stkXPRT worth of staking coin value according to NetAmount and performs LiquidDelegate.
func (k Keeper) LiquidStake(
	ctx sdk.Context, proxyAcc, liquidStaker sdk.AccAddress, stakingCoin sdk.Coin) (newShares math.LegacyDec, stkXPRTMintAmount math.Int, err error) {
	params := k.GetParams(ctx)

	// check minimum liquid stake amount
	if stakingCoin.Amount.LT(params.MinLiquidStakeAmount) {
		return sdk.ZeroDec(), sdk.ZeroInt(), types.ErrLessThanMinLiquidStakeAmount
	}

	// check bond denomination
	bondDenom := k.stakingKeeper.BondDenom(ctx)
	if stakingCoin.Denom != bondDenom {
		return sdk.ZeroDec(), sdk.ZeroInt(), errors.Wrapf(
			types.ErrInvalidBondDenom, "invalid coin denomination: got %s, expected %s", stakingCoin.Denom, bondDenom,
		)
	}

	whitelistedValsMap := types.GetWhitelistedValsMap(params.WhitelistedValidators)
	activeVals := k.GetActiveLiquidValidators(ctx, whitelistedValsMap)
	if activeVals.Len() == 0 || !activeVals.TotalWeight(whitelistedValsMap).IsPositive() {
		return sdk.ZeroDec(), sdk.ZeroInt(), types.ErrActiveLiquidValidatorsNotExists
	}

	// NetAmount must be calculated before send
	nas := k.GetNetAmountState(ctx)

	// send staking coin to liquid stake proxy account to proxy delegation, need sufficient spendable balances
	err = k.bankKeeper.SendCoins(ctx, liquidStaker, proxyAcc, sdk.NewCoins(stakingCoin))
	if err != nil {
		return sdk.ZeroDec(), sdk.ZeroInt(), err
	}

	// mint stkxprt, MintAmount = TotalSupply * StakeAmount/NetAmount
	liquidBondDenom := k.LiquidBondDenom(ctx)
	stkXPRTMintAmount = stakingCoin.Amount
	if nas.StkxprtTotalSupply.IsPositive() {
		stkXPRTMintAmount = types.NativeTokenToStkXPRT(stakingCoin.Amount, nas.StkxprtTotalSupply, nas.NetAmount)
	}

	if !stkXPRTMintAmount.IsPositive() {
		return sdk.ZeroDec(), sdk.ZeroInt(), types.ErrTooSmallLiquidStakeAmount
	}

	// mint on module acc and send
	mintCoin := sdk.NewCoins(sdk.NewCoin(liquidBondDenom, stkXPRTMintAmount))
	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, mintCoin)
	if err != nil {
		return sdk.ZeroDec(), stkXPRTMintAmount, err
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, liquidStaker, mintCoin)
	if err != nil {
		return sdk.ZeroDec(), stkXPRTMintAmount, err
	}

	newShares, err = k.LiquidDelegate(ctx, proxyAcc, activeVals, stakingCoin.Amount, whitelistedValsMap)
	return newShares, stkXPRTMintAmount, err
}

// LockOnLP sends tokens to a CW contract (Superfluid LP) with time locking.
// It performs a CosmWasm execution through global message handler and may fail.
// Emits events on a successful call.
func (k Keeper) LockOnLP(ctx sdk.Context, delegator sdk.AccAddress, amount sdk.Coin) ([]byte, error) {
	params := k.GetParams(ctx)

	if len(params.CwLockedPoolAddress) == 0 {
		return nil, types.ErrNoLPContractAddress
	} else if amount.Denom != params.LiquidBondDenom {
		return nil, types.ErrInvalidDenom.Wrapf("cannot lock any denom on LP except liquid bond denom: %s", params.LiquidBondDenom)
	}

	msg := &LockLstAssetMsg{
		Asset: Asset{
			Amount: amount.Amount.String(),
			Info: AssetInfo{
				NativeToken: NativeTokenInfo{
					Denom: amount.Denom,
				},
			},
		},
	}

	callData, err := json.Marshal(&ExecMsg{
		LockLstAsset: msg,
	})
	if err != nil {
		panic("failed to marshal CW contract call LockLstAsset")
	}

	cwMsg := &wasmtypes.MsgExecuteContract{
		Sender:   delegator.String(),
		Contract: k.GetParams(ctx).CwLockedPoolAddress,
		Msg:      wasmtypes.RawContractMessage(callData),
		Funds:    sdk.NewCoins(amount),
	}

	handler := k.router.Handler(cwMsg)
	if handler == nil {
		return nil, sdkerrors.ErrUnknownRequest.Wrapf("unrecognized message route: %s", sdk.MsgTypeURL(cwMsg))
	}

	msgResp, err := handler(ctx, cwMsg)
	if err != nil {
		return nil, types.ErrLPContract.Wrapf("error: %s, message %v", err.Error(), cwMsg)
	}

	// emit the events from the dispatched actions
	events := msgResp.Events
	sdkEvents := make([]sdk.Event, 0, len(events))
	for _, event := range events {
		e := event
		sdkEvents = append(sdkEvents, sdk.Event(e))
	}

	ctx.EventManager().EmitEvents(sdkEvents)

	return msgResp.Data, nil
}

type ExecMsg struct {
	LockLstAsset *LockLstAssetMsg `json:"lock_lst_asset,omitempty"`
}

type LockLstAssetMsg struct {
	Asset Asset `json:"asset"`
}

type Asset struct {
	Amount string    `json:"amount"`
	Info   AssetInfo `json:"info"`
}

type AssetInfo struct {
	NativeToken NativeTokenInfo `json:"native_token"`
}

type NativeTokenInfo struct {
	Denom string `json:"denom"`
}

// LSMDelegate captures a staked amount from existing delegation using LSM, re-stakes from proxyAcc and
// mints stkXPRT worth of stk coin value according to NetAmount and performs LiquidDelegate.
func (k Keeper) LSMDelegate(
	ctx sdk.Context,
	delegator sdk.AccAddress,
	validator sdk.ValAddress,
	proxyAcc sdk.AccAddress,
	stakingCoin sdk.Coin,
) (newShares math.LegacyDec, stkXPRTMintAmount math.Int, err error) {
	params := k.GetParams(ctx)

	if params.LsmDisabled {
		return sdk.ZeroDec(), sdk.ZeroInt(), types.ErrDisabledLSM
	}

	// check minimum liquid stake amount
	if stakingCoin.Amount.LT(params.MinLiquidStakeAmount) {
		return sdk.ZeroDec(), sdk.ZeroInt(), types.ErrLessThanMinLiquidStakeAmount
	}

	// check bond denomination
	bondDenom := k.stakingKeeper.BondDenom(ctx)
	if stakingCoin.Denom != bondDenom {
		return sdk.ZeroDec(), sdk.ZeroInt(), errors.Wrapf(
			types.ErrInvalidBondDenom, "invalid coin denomination: got %s, expected %s", stakingCoin.Denom, bondDenom,
		)
	}

	whitelistedValsMap := types.GetWhitelistedValsMap(params.WhitelistedValidators)
	activeVals := k.GetActiveLiquidValidators(ctx, whitelistedValsMap)
	if activeVals.Len() == 0 || !activeVals.TotalWeight(whitelistedValsMap).IsPositive() {
		return sdk.ZeroDec(), sdk.ZeroInt(), types.ErrActiveLiquidValidatorsNotExists
	}

	if !whitelistedValsMap.IsListed(validator.String()) {
		return sdk.ZeroDec(), sdk.ZeroInt(), types.ErrLiquidValidatorsNotExists.Wrap("delegation from a non allowed validator")
	}

	// NetAmount must be calculated before send
	nas := k.GetNetAmountState(ctx)

	// perform an LSM tokenize->bank send->redeem flow: moving delegation from user's account onto proxyAcc

	lsmTokenizeMsg := &stakingtypes.MsgTokenizeShares{
		DelegatorAddress:    delegator.String(),
		ValidatorAddress:    validator.String(),
		Amount:              stakingCoin,
		TokenizedShareOwner: proxyAcc.String(),
	}

	handler := k.router.Handler(lsmTokenizeMsg)
	if handler == nil {
		return sdk.ZeroDec(), sdk.ZeroInt(), sdkerrors.ErrUnknownRequest.Wrapf("unrecognized message route: %s", sdk.MsgTypeURL(lsmTokenizeMsg))
	}

	// [1] tokenize delegation into LSM shares
	msgResp, err := handler(ctx, lsmTokenizeMsg)
	if err != nil {
		return sdk.ZeroDec(), sdk.ZeroInt(), types.ErrLSMTokenizeFailed.Wrapf("error: %s; message: %v", err.Error(), lsmTokenizeMsg)
	}

	var lsmTokenizeResp stakingtypes.MsgTokenizeSharesResponse
	if err := k.cdc.Unmarshal(msgResp.Data, &lsmTokenizeResp); err != nil {
		panic("wrong data type: " + err.Error())
	}

	// [2] send LSM shares to proxyAcc
	err = k.bankKeeper.SendCoins(ctx, delegator, proxyAcc, sdk.NewCoins(lsmTokenizeResp.Amount))
	if err != nil {
		return sdk.ZeroDec(), stkXPRTMintAmount, err
	}

	lsmRedeemMsg := &stakingtypes.MsgRedeemTokensForShares{
		DelegatorAddress: proxyAcc.String(),
		Amount:           lsmTokenizeResp.Amount,
	}

	handler = k.router.Handler(lsmRedeemMsg)
	if handler == nil {
		return sdk.ZeroDec(), sdk.ZeroInt(), sdkerrors.ErrUnknownRequest.Wrapf("unrecognized message route: %s", sdk.MsgTypeURL(lsmRedeemMsg))
	}

	// [3] redeem LSM shares from proxyAcc, to obtain a delegation
	msgResp, err = handler(ctx, lsmRedeemMsg)
	if err != nil {
		return sdk.ZeroDec(), sdk.ZeroInt(), types.ErrLSMRedeemFailed.Wrapf("error: %s; message: %v", err.Error(), lsmRedeemMsg)
	}

	var lsmRedeemResp stakingtypes.MsgRedeemTokensForSharesResponse
	if err := k.cdc.Unmarshal(msgResp.Data, &lsmRedeemResp); err != nil {
		panic("wrong data type: " + err.Error())
	}

	// obtained newShares from LSM
	newShares = lsmRedeemResp.Amount.Amount.ToLegacyDec()

	// mint stkxprt, MintAmount = TotalSupply * StakeAmount/NetAmount
	liquidBondDenom := k.LiquidBondDenom(ctx)
	stkXPRTMintAmount = stakingCoin.Amount
	if nas.StkxprtTotalSupply.IsPositive() {
		stkXPRTMintAmount = types.NativeTokenToStkXPRT(stakingCoin.Amount, nas.StkxprtTotalSupply, nas.NetAmount)
	}

	if !stkXPRTMintAmount.IsPositive() {
		return sdk.ZeroDec(), sdk.ZeroInt(), types.ErrTooSmallLiquidStakeAmount
	}

	// mint stkXPRT on module acc
	mintCoin := sdk.NewCoins(sdk.NewCoin(liquidBondDenom, stkXPRTMintAmount))
	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, mintCoin)
	if err != nil {
		return sdk.ZeroDec(), stkXPRTMintAmount, err
	}

	// send stkXPRT to delegator acc
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, delegator, mintCoin)
	if err != nil {
		return sdk.ZeroDec(), stkXPRTMintAmount, err
	}

	// but immediately lock new stkXPRT into LP on behalf of the delegator
	_, err = k.LockOnLP(ctx, delegator, sdk.NewCoin(liquidBondDenom, stkXPRTMintAmount))
	if err != nil {
		return sdk.ZeroDec(), stkXPRTMintAmount, err
	}

	return newShares, stkXPRTMintAmount, err
}

// LiquidDelegate delegates staking amount to active validators by proxy account.
func (k Keeper) LiquidDelegate(ctx sdk.Context, proxyAcc sdk.AccAddress, activeVals types.ActiveLiquidValidators, stakingAmt math.Int, whitelistedValsMap types.WhitelistedValsMap) (newShares math.LegacyDec, err error) {
	totalNewShares := sdk.ZeroDec()
	// crumb may occur due to a decimal point error in dividing the staking amount into the weight of liquid validators, It added on first active liquid validator
	weightedAmt, crumb := types.DivideByWeight(activeVals, stakingAmt, whitelistedValsMap)
	if len(weightedAmt) == 0 {
		return sdk.ZeroDec(), types.ErrInvalidActiveLiquidValidators
	}
	weightedAmt[0] = weightedAmt[0].Add(crumb)
	for i, val := range activeVals {
		if !weightedAmt[i].IsPositive() {
			continue
		}
		validator, _ := k.stakingKeeper.GetValidator(ctx, val.GetOperator())
		newShares, err = k.stakingKeeper.Delegate(ctx, proxyAcc, weightedAmt[i], stakingtypes.Unbonded, validator, true)
		if err != nil {
			return sdk.ZeroDec(), err
		}
		totalNewShares = totalNewShares.Add(newShares)
	}
	return totalNewShares, nil
}

// LiquidUnstake burns unstakingStkXPRT and performs LiquidUnbond to active liquid validators with del shares worth of shares according to NetAmount with each validators current weight.
func (k Keeper) LiquidUnstake(
	ctx sdk.Context, proxyAcc, liquidStaker sdk.AccAddress, unstakingStkXPRT sdk.Coin,
) (time.Time, math.Int, []stakingtypes.UnbondingDelegation, math.Int, error) {

	// check bond denomination
	params := k.GetParams(ctx)
	liquidBondDenom := k.LiquidBondDenom(ctx)
	if unstakingStkXPRT.Denom != liquidBondDenom {
		return time.Time{}, sdk.ZeroInt(), []stakingtypes.UnbondingDelegation{}, sdk.ZeroInt(), errors.Wrapf(
			types.ErrInvalidLiquidBondDenom, "invalid coin denomination: got %s, expected %s", unstakingStkXPRT.Denom, liquidBondDenom,
		)
	}

	// Get NetAmount states
	nas := k.GetNetAmountState(ctx)

	if unstakingStkXPRT.Amount.GT(nas.StkxprtTotalSupply) {
		return time.Time{}, sdk.ZeroInt(), []stakingtypes.UnbondingDelegation{}, sdk.ZeroInt(), types.ErrInvalidStkXPRTSupply
	}

	// UnstakeAmount = NetAmount * StkXPRTAmount/TotalSupply * (1-UnstakeFeeRate)
	unbondingAmount := types.StkXPRTToNativeToken(unstakingStkXPRT.Amount, nas.StkxprtTotalSupply, nas.NetAmount)
	unbondingAmount = types.DeductFeeRate(unbondingAmount, params.UnstakeFeeRate)
	unbondingAmountInt := unbondingAmount.TruncateInt()

	if !unbondingAmountInt.IsPositive() {
		return time.Time{}, sdk.ZeroInt(), []stakingtypes.UnbondingDelegation{}, sdk.ZeroInt(), types.ErrTooSmallLiquidUnstakingAmount
	}

	// burn stkxprt
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, liquidStaker, types.ModuleName, sdk.NewCoins(unstakingStkXPRT))
	if err != nil {
		return time.Time{}, sdk.ZeroInt(), []stakingtypes.UnbondingDelegation{}, sdk.ZeroInt(), err
	}
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(liquidBondDenom, unstakingStkXPRT.Amount)))
	if err != nil {
		return time.Time{}, sdk.ZeroInt(), []stakingtypes.UnbondingDelegation{}, sdk.ZeroInt(), err
	}

	liquidVals := k.GetAllLiquidValidators(ctx)
	totalLiquidTokens, liquidTokenMap := liquidVals.TotalLiquidTokens(ctx, k.stakingKeeper, false)

	// if no totalLiquidTokens, withdraw directly from balance of proxy acc
	if !totalLiquidTokens.IsPositive() {
		if nas.ProxyAccBalance.GTE(unbondingAmountInt) {
			err = k.bankKeeper.SendCoins(
				ctx,
				types.LiquidStakeProxyAcc,
				liquidStaker,
				sdk.NewCoins(sdk.NewCoin(
					k.stakingKeeper.BondDenom(ctx),
					unbondingAmountInt,
				)),
			)
			if err != nil {
				return time.Time{}, sdk.ZeroInt(), []stakingtypes.UnbondingDelegation{}, sdk.ZeroInt(), err
			}

			return time.Time{}, sdk.ZeroInt(), []stakingtypes.UnbondingDelegation{}, unbondingAmountInt, nil
		}

		// error case where there is a quantity that are unbonding balance or remaining rewards that is not re-stake or withdrawn in netAmount.
		return time.Time{}, sdk.ZeroInt(), []stakingtypes.UnbondingDelegation{}, sdk.ZeroInt(), types.ErrInsufficientProxyAccBalance
	}

	// fail when no liquid validators to unbond
	if liquidVals.Len() == 0 {
		return time.Time{}, sdk.ZeroInt(), []stakingtypes.UnbondingDelegation{}, sdk.ZeroInt(), types.ErrLiquidValidatorsNotExists
	}

	// crumb may occur due to a decimal error in dividing the unstaking stkXPRT into the weight of liquid validators, it will remain in the NetAmount
	unbondingAmounts, crumb := types.DivideByCurrentWeight(liquidVals, unbondingAmount, totalLiquidTokens, liquidTokenMap)
	if !unbondingAmount.Sub(crumb).IsPositive() {
		return time.Time{}, sdk.ZeroInt(), []stakingtypes.UnbondingDelegation{}, sdk.ZeroInt(), types.ErrTooSmallLiquidUnstakingAmount
	}

	totalReturnAmount := sdk.ZeroInt()

	var ubdTime time.Time
	ubds := make([]stakingtypes.UnbondingDelegation, 0, len(liquidVals))
	for i, val := range liquidVals {
		// skip zero weight liquid validator
		if !unbondingAmounts[i].IsPositive() {
			continue
		}

		var ubd stakingtypes.UnbondingDelegation
		var returnAmount math.Int
		var weightedShare math.LegacyDec

		// calculate delShares from tokens with validation
		weightedShare, err = k.stakingKeeper.ValidateUnbondAmount(ctx, proxyAcc, val.GetOperator(), unbondingAmounts[i].TruncateInt())
		if err != nil {
			return time.Time{}, sdk.ZeroInt(), []stakingtypes.UnbondingDelegation{}, sdk.ZeroInt(), err
		}

		if !weightedShare.IsPositive() {
			continue
		}

		// unbond with weightedShare
		ubdTime, returnAmount, ubd, err = k.LiquidUnbond(ctx, proxyAcc, liquidStaker, val.GetOperator(), weightedShare, true)
		if err != nil {
			return time.Time{}, sdk.ZeroInt(), []stakingtypes.UnbondingDelegation{}, sdk.ZeroInt(), err
		}

		ubds = append(ubds, ubd)
		totalReturnAmount = totalReturnAmount.Add(returnAmount)
	}

	return ubdTime, totalReturnAmount, ubds, sdk.ZeroInt(), nil
}

// LiquidUnbond unbond delegation shares to active validators by proxy account.
func (k Keeper) LiquidUnbond(
	ctx sdk.Context, proxyAcc, liquidStaker sdk.AccAddress, valAddr sdk.ValAddress, shares math.LegacyDec, checkMaxEntries bool,
) (time.Time, math.Int, stakingtypes.UnbondingDelegation, error) {
	validator, found := k.stakingKeeper.GetValidator(ctx, valAddr)
	if !found {
		return time.Time{}, sdk.ZeroInt(), stakingtypes.UnbondingDelegation{}, stakingtypes.ErrNoDelegatorForAddress
	}

	// If checkMaxEntries is true, perform a maximum limit unbonding entries check.
	if checkMaxEntries && k.stakingKeeper.HasMaxUnbondingDelegationEntries(ctx, liquidStaker, valAddr) {
		return time.Time{}, sdk.ZeroInt(), stakingtypes.UnbondingDelegation{}, stakingtypes.ErrMaxUnbondingDelegationEntries
	}

	// unbond from proxy account
	returnAmount, err := k.stakingKeeper.Unbond(ctx, proxyAcc, valAddr, shares)
	if err != nil {
		return time.Time{}, sdk.ZeroInt(), stakingtypes.UnbondingDelegation{}, err
	}

	// transfer the validator tokens to the not bonded pool
	if validator.IsBonded() {
		coins := sdk.NewCoins(sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), returnAmount))
		if err = k.bankKeeper.SendCoinsFromModuleToModule(ctx, stakingtypes.BondedPoolName, stakingtypes.NotBondedPoolName, coins); err != nil {
			panic(err)
		}
	}

	// Unbonding from proxy account, but queues to liquid staker.
	completionTime := ctx.BlockHeader().Time.Add(k.stakingKeeper.UnbondingTime(ctx))
	ubd := k.stakingKeeper.SetUnbondingDelegationEntry(ctx, liquidStaker, valAddr, ctx.BlockHeight(), completionTime, returnAmount)
	k.stakingKeeper.InsertUBDQueue(ctx, ubd, completionTime)

	return completionTime, returnAmount, ubd, nil
}

// CheckDelegationStates returns total remaining rewards, delshares, liquid tokens of delegations by proxy account
func (k Keeper) CheckDelegationStates(ctx sdk.Context, proxyAcc sdk.AccAddress) (math.LegacyDec, math.LegacyDec, math.Int) {
	bondDenom := k.stakingKeeper.BondDenom(ctx)
	totalRewards := sdk.ZeroDec()
	totalDelShares := sdk.ZeroDec()
	totalLiquidTokens := sdk.ZeroInt()

	// Cache ctx for calculate rewards
	cachedCtx, _ := ctx.CacheContext()
	k.stakingKeeper.IterateDelegations(
		cachedCtx, proxyAcc,
		func(_ int64, del stakingtypes.DelegationI) (stop bool) {
			valAddr := del.GetValidatorAddr()
			val := k.stakingKeeper.Validator(cachedCtx, valAddr)
			endingPeriod := k.distrKeeper.IncrementValidatorPeriod(cachedCtx, val)
			delReward := k.distrKeeper.CalculateDelegationRewards(cachedCtx, val, del, endingPeriod)
			delShares := del.GetShares()
			if delShares.IsPositive() {
				totalDelShares = totalDelShares.Add(delShares)
				liquidTokens := val.TokensFromSharesTruncated(delShares).TruncateInt()
				totalLiquidTokens = totalLiquidTokens.Add(liquidTokens)
				totalRewards = totalRewards.Add(delReward.AmountOf(bondDenom))
			}
			return false
		},
	)

	return totalRewards, totalDelShares, totalLiquidTokens
}

func (k Keeper) WithdrawLiquidRewards(ctx sdk.Context, proxyAcc sdk.AccAddress) math.Int {
	totalRewards := sdk.ZeroInt()
	bondDenom := k.stakingKeeper.BondDenom(ctx)
	k.stakingKeeper.IterateDelegations(
		ctx, proxyAcc,
		func(_ int64, del stakingtypes.DelegationI) (stop bool) {
			valAddr := del.GetValidatorAddr()
			reward, err := k.distrKeeper.WithdrawDelegationRewards(ctx, proxyAcc, valAddr)
			if err != nil {
				panic(err)
			}
			totalRewards = totalRewards.Add(reward.AmountOf(bondDenom))
			return false
		},
	)
	return totalRewards
}

// GetLiquidValidator get a single liquid validator
func (k Keeper) GetLiquidValidator(ctx sdk.Context, addr sdk.ValAddress) (val types.LiquidValidator, found bool) {
	store := ctx.KVStore(k.storeKey)

	value := store.Get(types.GetLiquidValidatorKey(addr))
	if value == nil {
		return val, false
	}

	val = types.MustUnmarshalLiquidValidator(k.cdc, value)
	return val, true
}

// SetLiquidValidator set the main record holding liquid validator details
func (k Keeper) SetLiquidValidator(ctx sdk.Context, val types.LiquidValidator) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalLiquidValidator(k.cdc, &val)
	store.Set(types.GetLiquidValidatorKey(val.GetOperator()), bz)
}

// RemoveLiquidValidator remove a liquid validator on kv store
func (k Keeper) RemoveLiquidValidator(ctx sdk.Context, val types.LiquidValidator) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetLiquidValidatorKey(val.GetOperator()))
}

// GetAllLiquidValidators get the set of all liquid validators with no limits, used during genesis dump
func (k Keeper) GetAllLiquidValidators(ctx sdk.Context) (vals types.LiquidValidators) {
	store := ctx.KVStore(k.storeKey)
	vals = types.LiquidValidators{}
	iterator := sdk.KVStorePrefixIterator(store, types.LiquidValidatorsKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		val := types.MustUnmarshalLiquidValidator(k.cdc, iterator.Value())
		vals = append(vals, val)
	}

	return vals
}

// GetActiveLiquidValidators get the set of active liquid validators.
func (k Keeper) GetActiveLiquidValidators(ctx sdk.Context, whitelistedValsMap types.WhitelistedValsMap) (vals types.ActiveLiquidValidators) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.LiquidValidatorsKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		val := types.MustUnmarshalLiquidValidator(k.cdc, iterator.Value())
		if k.IsActiveLiquidValidator(ctx, val, whitelistedValsMap) {
			vals = append(vals, val)
		}
	}
	return vals
}

func (k Keeper) GetAllLiquidValidatorStates(ctx sdk.Context) (liquidValidatorStates []types.LiquidValidatorState) {
	lvs := k.GetAllLiquidValidators(ctx)
	whitelistedValsMap := k.GetParams(ctx).WhitelistedValsMap()
	for _, lv := range lvs {
		active := k.IsActiveLiquidValidator(ctx, lv, whitelistedValsMap)
		lvState := types.LiquidValidatorState{
			OperatorAddress: lv.OperatorAddress,
			Weight:          lv.GetWeight(whitelistedValsMap, active),
			Status:          lv.GetStatus(active),
			DelShares:       lv.GetDelShares(ctx, k.stakingKeeper),
			LiquidTokens:    lv.GetLiquidTokens(ctx, k.stakingKeeper, false),
		}
		liquidValidatorStates = append(liquidValidatorStates, lvState)
	}
	return
}

func (k Keeper) GetLiquidValidatorState(ctx sdk.Context, addr sdk.ValAddress) (liquidValidatorState types.LiquidValidatorState, found bool) {
	lv, found := k.GetLiquidValidator(ctx, addr)
	if !found {
		return types.LiquidValidatorState{
			OperatorAddress: addr.String(),
			Weight:          sdk.ZeroInt(),
			Status:          types.ValidatorStatusUnspecified,
			DelShares:       sdk.ZeroDec(),
			LiquidTokens:    sdk.ZeroInt(),
		}, false
	}
	whitelistedValsMap := k.GetParams(ctx).WhitelistedValsMap()
	active := k.IsActiveLiquidValidator(ctx, lv, whitelistedValsMap)
	return types.LiquidValidatorState{
		OperatorAddress: lv.OperatorAddress,
		Weight:          lv.GetWeight(whitelistedValsMap, active),
		Status:          lv.GetStatus(active),
		DelShares:       lv.GetDelShares(ctx, k.stakingKeeper),
		LiquidTokens:    lv.GetLiquidTokens(ctx, k.stakingKeeper, false),
	}, true
}

func (k Keeper) IsActiveLiquidValidator(ctx sdk.Context, lv types.LiquidValidator, whitelistedValsMap types.WhitelistedValsMap) bool {
	val, found := k.stakingKeeper.GetValidator(ctx, lv.GetOperator())
	if !found {
		return false
	}
	return types.ActiveCondition(val, whitelistedValsMap.IsListed(lv.OperatorAddress), k.IsTombstoned(ctx, val))
}

func (k Keeper) IsTombstoned(ctx sdk.Context, val stakingtypes.Validator) bool {
	consPk, err := val.ConsPubKey()
	if err != nil {
		return false
	}
	return k.slashingKeeper.IsTombstoned(ctx, sdk.ConsAddress(consPk.Address()))
}

func (k Keeper) GetWeightMap(ctx sdk.Context, liquidVals types.LiquidValidators, whitelistedValsMap types.WhitelistedValsMap) (map[string]math.Int, math.Int) {
	weightMap := map[string]math.Int{}
	totalWeight := sdk.ZeroInt()
	for _, val := range liquidVals {
		weight := val.GetWeight(whitelistedValsMap, k.IsActiveLiquidValidator(ctx, val, whitelistedValsMap))
		totalWeight = totalWeight.Add(weight)
		weightMap[val.OperatorAddress] = weight
	}
	return weightMap, totalWeight
}
