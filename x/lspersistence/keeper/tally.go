package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/persistenceOne/pstake-native/x/lspersistence/types"
)

// GetVoterBalanceByDenom return map of balance amount of voter by denom
func (k Keeper) GetVoterBalanceByDenom(ctx sdk.Context, votes govtypes.Votes) map[string]map[string]sdk.Int {
	denomAddrBalanceMap := map[string]map[string]sdk.Int{}
	for _, vote := range votes {
		voter, err := sdk.AccAddressFromBech32(vote.Voter)
		if err != nil {
			continue
		}
		balances := k.bankKeeper.SpendableCoins(ctx, voter)
		for _, coin := range balances {
			if _, ok := denomAddrBalanceMap[coin.Denom]; !ok {
				denomAddrBalanceMap[coin.Denom] = map[string]sdk.Int{}
			}
			if coin.Amount.IsPositive() {
				denomAddrBalanceMap[coin.Denom][vote.Voter] = coin.Amount
			}
		}
	}
	return denomAddrBalanceMap
}

func (k Keeper) GetVotingPower(ctx sdk.Context, addr sdk.AccAddress) types.VotingPower {
	val, found := k.stakingKeeper.GetValidator(ctx, addr.Bytes())
	validatorVotingPower := sdk.ZeroInt()
	if found {
		validatorVotingPower = val.BondedTokens()
	}
	return types.VotingPower{
		Voter:                    addr.String(),
		StakingVotingPower:       k.CalcStakingVotingPower(ctx, addr),
		LiquidStakingVotingPower: k.CalcLiquidStakingVotingPower(ctx, addr),
		ValidatorVotingPower:     validatorVotingPower,
	}
}

// CalcStakingVotingPower returns voting power of the addr by normal delegations except self-delegation
func (k Keeper) CalcStakingVotingPower(ctx sdk.Context, addr sdk.AccAddress) sdk.Int {
	totalVotingPower := sdk.ZeroInt()
	k.stakingKeeper.IterateDelegations(
		ctx, addr,
		func(_ int64, del stakingtypes.DelegationI) (stop bool) {
			valAddr := del.GetValidatorAddr()
			val := k.stakingKeeper.Validator(ctx, valAddr)
			delShares := del.GetShares()
			// if the validator not bonded, bonded token and voting power is zero, and except self-delegation power
			if delShares.IsPositive() && val.IsBonded() && !valAddr.Equals(addr) {
				votingPower := val.TokensFromSharesTruncated(delShares).TruncateInt()
				if votingPower.IsPositive() {
					totalVotingPower = totalVotingPower.Add(votingPower)
				}
			}
			return false
		},
	)
	return totalVotingPower
}

// CalcLiquidStakingVotingPower returns voting power of the addr by liquid bond denom
func (k Keeper) CalcLiquidStakingVotingPower(ctx sdk.Context, addr sdk.AccAddress) sdk.Int {
	liquidBondDenom := k.LiquidBondDenom(ctx)

	// skip when no liquid bond token supply
	bTokenTotalSupply := k.bankKeeper.GetSupply(ctx, liquidBondDenom).Amount
	if !bTokenTotalSupply.IsPositive() {
		return sdk.ZeroInt()
	}

	// skip when no active validators, liquid tokens
	liquidVals := k.GetAllLiquidValidators(ctx)
	if len(liquidVals) == 0 {
		return sdk.ZeroInt()
	}

	// using only liquid tokens of bonded liquid validators to ensure voting power doesn't exceed delegation shares on x/gov tally
	totalBondedLiquidTokens, _ := liquidVals.TotalLiquidTokens(ctx, k.stakingKeeper, true)
	if !totalBondedLiquidTokens.IsPositive() {
		return sdk.ZeroInt()
	}

	bTokenAmount := sdk.ZeroInt()

	balances := k.bankKeeper.SpendableCoins(ctx, addr)
	for _, coin := range balances {
		// add balance of bToken
		if coin.Denom == liquidBondDenom {
			bTokenAmount = bTokenAmount.Add(coin.Amount)
		}
	}


	if bTokenAmount.IsPositive() {
		return types.BTokenToNativeToken(bTokenAmount, bTokenTotalSupply, totalBondedLiquidTokens.ToDec()).TruncateInt()
	} else {
		return sdk.ZeroInt()
	}
}
