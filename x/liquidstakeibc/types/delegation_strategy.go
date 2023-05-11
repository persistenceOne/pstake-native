package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ValidatorDelegation struct {
	ValidatorAddress string
	DelegationAmount sdk.Coin
}
