package types

import sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"

// NewGenesisState creates a new GenesisState object. todo fill up outgoing transactions array and maintain it with a length of 10000 txns
func NewGenesisState(params Params, outgoingTx OutgoingTx) *GenesisState {
	return &GenesisState{
		Params:      params,
		OutgoingTxn: outgoingTx,
	}
}

// ValidateGenesis validates the provided genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return sdkErrors.Wrap(err, "params")
	}
	return nil
}

// DefaultGenesisState creates a default GenesisState object.
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:      DefaultParams(),
		OutgoingTxn: OutgoingTx{},
	}
}

// Equal check if two genesis states are equal
func (data GenesisState) Equal(other GenesisState) bool {
	return data.Params.Equal(other.Params)
}
