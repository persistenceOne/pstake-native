package types

import (
	"encoding/json"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Coin returns the sdk.Coin amount from WeightedAddressAmount struct
func (w WeightedAddressAmount) Coin() sdk.Coin {
	return sdk.NewCoin(w.Denom, w.Amount)
}

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

type WeightedAddressAmounts []WeightedAddressAmount

var _ sort.Interface = WeightedAddressAmounts{}

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
	return ws[i].Amount.LT(ws[j].Amount)
}
func (ws WeightedAddressAmounts) Swap(i, j int) {
	ws[i], ws[j] = ws[j], ws[i]
}
func (ws WeightedAddressAmounts) Sort() WeightedAddressAmounts {
	sort.Sort(ws)
	return ws
}

func (ws WeightedAddressAmounts) Marshal() ([]byte, error) {
	if ws == nil {
		return json.Marshal(WeightedAddressAmounts{})
	}
	return json.Marshal(ws)
}

func (ws WeightedAddressAmounts) Unmarshal(bz []byte) error {
	err := json.Unmarshal(bz, &ws)
	if err != nil {
		return err
	}
	return nil
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
		if w.Amount.IsNegative() {
			zeroValuedAddrAmts = append(zeroValuedAddrAmts, w)
		}
	}
	return zeroValuedAddrAmts
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
