package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ValAddressAndAmountForStakingAndUndelegating struct {
	validator sdk.ValAddress
	amount    sdk.Coin
}

// gives a list of all validators having weighted amount for few and 1uatom for rest in order to auto claim all rewards accumulated in current epoch
func (k Keeper) fetchValidatorsToDelegate(ctx sdk.Context, amount sdk.Coin) []ValAddressAndAmountForStakingAndUndelegating {
	//TODO : Add pseudo code for filtering out validators to delegate
	return nil
}

// gives a list of validators having weighted amount for few validators
func (k Keeper) fetchValidatorsToUndelegate(ctx sdk.Context, amount sdk.Coin) []ValAddressAndAmountForStakingAndUndelegating {
	//TODO : Implement opposite of fetchValidatorsToDelegate
	return nil
}
