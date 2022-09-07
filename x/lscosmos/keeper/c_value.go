package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetMintedAmount gets minted amount
func (k Keeper) GetMintedAmount(ctx sdk.Context) sdk.Int {
	return k.bankKeeper.GetSupply(ctx, k.GetHostChainParams(ctx).MintDenom).Amount
}

func (k Keeper) GetDepositAccountAmount(ctx sdk.Context) sdk.Int {
	return k.bankKeeper.GetBalance(
		ctx,
		k.GetDepositModuleAccount(ctx).GetAddress(),
		k.GetIBCDenom(ctx),
	).Amount
}

func (k Keeper) GetDelegationAccountAmount(ctx sdk.Context) sdk.Int {
	return k.bankKeeper.GetBalance(
		ctx,
		k.GetDelegationModuleAccount(ctx).GetAddress(),
		k.GetIBCDenom(ctx),
	).Amount
}

func (k Keeper) GetIBCTransferTransientAmount(ctx sdk.Context) sdk.Int {
	// TODO get amount from transient state
	return sdk.ZeroInt()
}

func (k Keeper) GetDelegationTransientAmount(ctx sdk.Context) sdk.Int {
	// TODO get amount from transient state
	return sdk.ZeroInt()
}

func (k Keeper) GetStakedAmount(ctx sdk.Context) sdk.Int {
	sum := sdk.ZeroInt()
	for _, delegation := range k.GetDelegationState(ctx).HostAccountDelegations {
		sum = sum.Add(delegation.Amount.Amount)
	}
	return sum
}

func (k Keeper) GetHostDelegationAccountAmount(ctx sdk.Context) sdk.Int {
	return k.GetDelegationState(ctx).HostDelegationAccountBalance.AmountOf(k.GetHostChainParams(ctx).BaseDenom)
}

// GetCValue gets the C value after recalculating everytime when the
// function is called. Returns 1 if stakedAmount or mintedAmount is zero.
func (k Keeper) GetCValue(ctx sdk.Context) sdk.Dec {
	stakedAmount := k.GetDepositAccountAmount(ctx).
		Add(k.GetDelegationAccountAmount(ctx)).
		Add(k.GetIBCTransferTransientAmount(ctx)).
		Add(k.GetDelegationTransientAmount(ctx)).
		Add(k.GetStakedAmount(ctx)).
		Add(k.GetHostDelegationAccountAmount(ctx))

	mintedAmount := k.GetMintedAmount(ctx)
	if stakedAmount.IsZero() || mintedAmount.IsZero() {
		return sdk.OneDec()
	}

	return sdk.NewDecFromInt(mintedAmount).Quo(sdk.NewDecFromInt(stakedAmount))
}

func (k Keeper) ConvertStkToToken(ctx sdk.Context, stkCoin sdk.Coin) (sdk.Coin, sdk.DecCoin) {
	// calculate the current stkToken value
	tokenValue := stkCoin.Amount.ToDec().Mul(sdk.OneDec().Quo(k.GetCValue(ctx)))

	return sdk.NewDecCoinFromDec(k.GetIBCDenom(ctx), tokenValue).TruncateDecimal()
}

func (k Keeper) ConvertTokenToStk(ctx sdk.Context, token sdk.Coin) (sdk.Coin, sdk.DecCoin) {
	mintDenom := k.GetHostChainParams(ctx).MintDenom

	// calculate the current token value
	tokenValue := token.Amount.ToDec().Mul(k.GetCValue(ctx))

	return sdk.NewDecCoinFromDec(mintDenom, tokenValue).TruncateDecimal()
}
