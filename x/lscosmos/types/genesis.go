package types

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{

		// this line is used by starport scaffolding # genesis/types/default
		Params:                DefaultParams(),
		ModuleEnabled:         false,
		HostChainParams:       HostChainParams{},
		AllowListedValidators: AllowListedValidators{},
		DelegationState:       DelegationState{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {

	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
