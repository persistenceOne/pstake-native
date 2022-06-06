package types

import (
	"encoding/json"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (w WeightedAddressAmount) Coin() sdk.Coin {
	return sdk.NewCoin(w.Denom, w.Amount)
}

func NewWeightedAddressAmount(address string, weight sdk.Dec, coin sdk.Coin) WeightedAddressAmount {
	return WeightedAddressAmount{
		Address: address,
		Weight:  weight,
		Denom:   coin.Denom,
		Amount:  coin.Amount,
	}
}

type WeightedAddressAmounts []WeightedAddressAmount

var _ sort.Interface = WeightedAddressAmounts{}

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

func (ws WeightedAddressAmounts) GetZeroWeighted() WeightedAddressAmounts {
	zeroWeightedAddrAmts := WeightedAddressAmounts{}
	for _, w := range ws {
		if w.Weight.IsZero() {
			zeroWeightedAddrAmts = append(zeroWeightedAddrAmts, w)
		}
	}
	return zeroWeightedAddrAmts
}

func (ws WeightedAddressAmounts) GetZeroValued() WeightedAddressAmounts {
	zeroValuedAddrAmts := WeightedAddressAmounts{}
	for _, w := range ws {
		if w.Amount.IsNegative() {
			zeroValuedAddrAmts = append(zeroValuedAddrAmts, w)
		}
	}
	return zeroValuedAddrAmts
}

func GetWeightedAddressMap(ws WeightedAddressAmounts) map[string]sdk.Dec {
	addressMap := map[string]sdk.Dec{}
	for _, w := range ws {
		addressMap[w.Address] = w.Weight
	}
	return addressMap
}

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
