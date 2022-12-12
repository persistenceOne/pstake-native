package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

// GetMintedAmount gets minted amount
func (k Keeper) GetMintedAmount(ctx sdk.Context) sdk.Int {
	return k.bankKeeper.GetSupply(ctx, k.GetHostChainParams(ctx).MintDenom).Amount
}

// GetDepositAccountAmount returns the deposit account amount of the IBC denom
func (k Keeper) GetDepositAccountAmount(ctx sdk.Context) sdk.Int {
	return k.bankKeeper.GetBalance(
		ctx,
		authtypes.NewModuleAddress(types.DepositModuleAccount),
		k.GetIBCDenom(ctx),
	).Amount
}

// GetIBCTransferTransientAmount returns the IBC transfer transient amount of the IBC denom
func (k Keeper) GetIBCTransferTransientAmount(ctx sdk.Context) sdk.Int {
	transferAmount := k.GetIBCTransientStore(ctx).IBCTransfer

	sum := sdk.ZeroInt()
	for _, coin := range transferAmount {
		sum = sum.Add(coin.Amount)
	}

	return sum
}

// GetDelegationTransientAmount returns the delegation transient amount of the IBC denom
func (k Keeper) GetDelegationTransientAmount(ctx sdk.Context) sdk.Int {
	icaDelegateAmount := k.GetIBCTransientStore(ctx).ICADelegate.Amount
	if icaDelegateAmount.IsNil() {
		return sdk.ZeroInt()
	}

	return icaDelegateAmount
}

// GetStakedAmount returns the total staked amount stored in delegation state
func (k Keeper) GetStakedAmount(ctx sdk.Context) sdk.Int {
	sum := sdk.ZeroInt()
	for _, delegation := range k.GetDelegationState(ctx).HostAccountDelegations {
		sum = sum.Add(delegation.Amount.Amount)
	}
	return sum
}

// GetHostDelegationAccountAmount returns the host account delegation account amount of IBC denom
func (k Keeper) GetHostDelegationAccountAmount(ctx sdk.Context) sdk.Int {
	return k.GetDelegationState(ctx).HostDelegationAccountBalance.AmountOf(k.GetHostChainParams(ctx).BaseDenom)
}

// GetCValue gets the C value after recalculating everytime when the
// function is called. Returns 1 if stakedAmount or mintedAmount is zero.
func (k Keeper) GetCValue(ctx sdk.Context) sdk.Dec {
	stakedAmount := k.GetDepositAccountAmount(ctx).
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

// ConvertStkToToken converts a stk coin to token with the c value provided in the input.
func (k Keeper) ConvertStkToToken(ctx sdk.Context, stkCoin sdk.DecCoin, cValue sdk.Dec) (sdk.Coin, sdk.DecCoin) {

	// calculate the current stkToken value
	tokenValue := stkCoin.Amount.Mul(sdk.OneDec().Quo(cValue))

	return sdk.NewDecCoinFromDec(k.GetIBCDenom(ctx), tokenValue).TruncateDecimal()
}

// ConvertTokenToStk converts a token coin to stk coin with the c value provided in the input.
func (k Keeper) ConvertTokenToStk(ctx sdk.Context, token sdk.DecCoin, cValue sdk.Dec) (sdk.Coin, sdk.DecCoin) {
	mintDenom := k.GetHostChainParams(ctx).MintDenom

	// calculate the current token value
	tokenValue := token.Amount.Mul(cValue)

	return sdk.NewDecCoinFromDec(mintDenom, tokenValue).TruncateDecimal()
}
