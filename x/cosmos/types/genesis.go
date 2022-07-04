package types

import sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"

// NewGenesisState creates a new GenesisState object.
func NewGenesisState(params Params) *GenesisState {
	return &GenesisState{
		Params: params,
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
		Params: DefaultParams(),
	}
}

// Equal check if two genesis states are equal
func (data GenesisState) Equal(other GenesisState) bool {
	return data.Params.Equal(other.Params)
}
