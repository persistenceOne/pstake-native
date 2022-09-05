package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
)

// GetMintedAmount gets minted amount
func (k Keeper) GetMintedAmount(ctx sdk.Context) sdk.Int {
	return k.bankKeeper.GetSupply(ctx, k.GetHostChainParams(ctx).MintDenom).Amount
}

func (k Keeper) GetDepositAccountAmount(ctx sdk.Context) sdk.Int {
	hostChainParams := k.GetHostChainParams(ctx)
	ibcDenom := ibctransfertypes.ParseDenomTrace(
		ibctransfertypes.GetPrefixedDenom(
			hostChainParams.TransferPort, hostChainParams.TransferChannel, hostChainParams.BaseDenom,
		),
	).IBCDenom()

	return k.bankKeeper.GetBalance(
		ctx,
		k.GetDepositModuleAccount(ctx).GetAddress(),
		ibcDenom,
	).Amount
}

func (k Keeper) GetDelegationAccountAmount(ctx sdk.Context) sdk.Int {
	hostChainParams := k.GetHostChainParams(ctx)
	ibcDenom := ibctransfertypes.ParseDenomTrace(
		ibctransfertypes.GetPrefixedDenom(
			hostChainParams.TransferPort, hostChainParams.TransferChannel, hostChainParams.BaseDenom,
		),
	).IBCDenom()
	return k.bankKeeper.GetBalance(
		ctx,
		k.GetDelegationModuleAccount(ctx).GetAddress(),
		ibcDenom,
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

// GetCValue gets the C cached C value if cache is valid or re-calculates if expired
// returns 1 in case where total staked amount is 0
func (k Keeper) GetCValue(ctx sdk.Context) sdk.Dec {
	stakedAmount := k.GetDepositAccountAmount(ctx).
		Add(k.GetDelegationAccountAmount(ctx)).
		Add(k.GetIBCTransferTransientAmount(ctx)).
		Add(k.GetDelegationTransientAmount(ctx)).
		Add(k.GetStakedAmount(ctx)).
		Add(k.GetHostDelegationAccountAmount(ctx))

	if stakedAmount.IsZero() || k.GetMintedAmount(ctx).IsZero() {
		return sdk.OneDec()
	}

	return sdk.NewDecFromInt(k.GetMintedAmount(ctx)).Quo(sdk.NewDecFromInt(stakedAmount))
}
