package types

import (
	"sort"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// WeightedAddressAmount defines address and their corresponding weight, amount, denom
// unbonding tokens
type WeightedAddressAmount struct {
	Address         string
	Weight          sdk.Dec
	Amount          math.Int
	Denom           string
	UnbondingTokens sdk.Coin
}

// ValAddressAmount defines validator address and it's corresponding amount
type ValAddressAmount struct {
	ValidatorAddr string
	Amount        sdk.Coin
}

type (
	WeightedAddressAmounts []WeightedAddressAmount
	ValAddressAmounts      []ValAddressAmount
)

// NewWeightedAddressAmount returns WeightedAddressAmount struct populated with given details
func NewWeightedAddressAmount(address string, weight sdk.Dec, coin sdk.Coin, unbondingTokens sdk.Coin) WeightedAddressAmount {
	return WeightedAddressAmount{
		Address:         address,
		Weight:          weight,
		Denom:           coin.Denom,
		Amount:          coin.Amount,
		UnbondingTokens: unbondingTokens,
	}
}

var (
	_ sort.Interface = WeightedAddressAmounts{}
	_ sort.Interface = ValAddressAmounts{}
)

// NewWeightedAddressAmounts returns WeightedAddressAmounts array
func NewWeightedAddressAmounts(w []WeightedAddressAmount) WeightedAddressAmounts {
	ws := WeightedAddressAmounts{}
	for _, element := range w {
		ws = append(ws, element)
	}
	return ws
}

func (ws WeightedAddressAmounts) Len() int {
	return len(ws)
}
func (ws WeightedAddressAmounts) Less(i, j int) bool {
	if ws[i].Amount.LT(ws[j].Amount) {
		return true
	}
	if ws[i].Amount.GT(ws[j].Amount) {
		return false
	}
	return ws[i].Address < ws[j].Address
}
func (ws WeightedAddressAmounts) Swap(i, j int) {
	ws[i], ws[j] = ws[j], ws[i]
}

// Coin returns the sdk.Coin amount from WeightedAddressAmount struct
func (w WeightedAddressAmount) Coin() sdk.Coin {
	return sdk.NewCoin(w.Denom, w.Amount)
}

// TotalAmount returns the total amount for a given denom
func (ws WeightedAddressAmounts) TotalAmount(denom string) sdk.Coin {
	total := sdk.NewCoin(denom, sdk.ZeroInt())

	for _, weightedAddr := range ws {
		if weightedAddr.Denom == denom {
			total.Amount = total.Amount.Add(weightedAddr.Amount)
		}
	}
	return total
}

// GetZeroWeighted returns the list of WeightedAddressAmount with zero weights
func (ws WeightedAddressAmounts) GetZeroWeighted() WeightedAddressAmounts {
	zeroWeightedAddrAmts := WeightedAddressAmounts{}
	for _, w := range ws {
		if w.Weight.IsZero() {
			zeroWeightedAddrAmts = append(zeroWeightedAddrAmts, w)
		}
	}
	return zeroWeightedAddrAmts
}

// GetZeroValued returns the list of WeightedAddressAmount with zero amount
func (ws WeightedAddressAmounts) GetZeroValued() WeightedAddressAmounts {
	zeroValuedAddrAmts := WeightedAddressAmounts{}
	for _, w := range ws {
		if w.Amount.IsZero() {
			zeroValuedAddrAmts = append(zeroValuedAddrAmts, w)
		}
	}
	return zeroValuedAddrAmts
}

func (ws ValAddressAmounts) Len() int {
	return len(ws)
}
func (ws ValAddressAmounts) Less(i, j int) bool {
	return ws[i].ValidatorAddr < ws[j].ValidatorAddr
}
func (ws ValAddressAmounts) Swap(i, j int) {
	ws[i], ws[j] = ws[j], ws[i]
}

// GetHostAccountDelegationMap returns the map of address as key and delegations as the value
func GetHostAccountDelegationMap(hostAccountDelegations []HostAccountDelegation) map[string]sdk.Coin {
	delegationMap := map[string]sdk.Coin{}
	for _, del := range hostAccountDelegations {
		delegationMap[del.ValidatorAddress] = del.Amount
	}
	return delegationMap
}

// GetWeightedAddressMap returns the map of address as key and weights as value
func GetWeightedAddressMap(ws WeightedAddressAmounts) map[string]sdk.Dec {
	addressMap := map[string]sdk.Dec{}
	for _, w := range ws {
		addressMap[w.Address] = w.Weight
	}
	return addressMap
}

// GetZeroNonZeroWightedAddrAmts returns a list of WeightedAddressAmount zero weights and non zero weights elements
func GetZeroNonZeroWightedAddrAmts(ws WeightedAddressAmounts) (zeroWeighted, nonZeroWeighted WeightedAddressAmounts) {
	for _, w := range ws {
		if w.Weight.IsZero() {
			zeroWeighted = append(zeroWeighted, w)
		} else {
			nonZeroWeighted = append(nonZeroWeighted, w)
		}
	}
	return zeroWeighted, nonZeroWeighted
}
