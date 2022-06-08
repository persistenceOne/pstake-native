package types

import (
	"errors"
	"time"
)

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

func NewGenesisState(epochs []EpochInfo) *GenesisState {
	return &GenesisState{Epochs: epochs}
}

// DefaultGenesis returns the default Capability genesis state
// keep stake and rewards epoch time duration same to keep cosmos module code from breaking
func DefaultGenesis() *GenesisState {
	epochs := []EpochInfo{
		{
			Identifier:              "stake",
			StartTime:               time.Time{},
			Duration:                time.Minute * 2, //TODO revert
			CurrentEpoch:            0,
			CurrentEpochStartHeight: 0,
			CurrentEpochStartTime:   time.Time{},
			EpochCountingStarted:    false,
		},
		{
			Identifier:              "reward",
			StartTime:               time.Time{},
			Duration:                time.Minute * 2, //TODO revert
			CurrentEpoch:            0,
			CurrentEpochStartHeight: 0,
			CurrentEpochStartTime:   time.Time{},
			EpochCountingStarted:    false,
		},
		{
			Identifier:              "undelegate",
			StartTime:               time.Time{},
			Duration:                time.Minute * 15, //TODO revert
			CurrentEpoch:            0,
			CurrentEpochStartHeight: 0,
			CurrentEpochStartTime:   time.Time{},
			EpochCountingStarted:    false,
		},
	}
	return NewGenesisState(epochs)
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// TODO: Epochs identifiers should be unique
	epochIdentifiers := map[string]bool{}
	for _, epoch := range gs.Epochs {
		if epoch.Identifier == "" {
			return errors.New("epoch identifier should NOT be empty")
		}
		if epochIdentifiers[epoch.Identifier] {
			return errors.New("epoch identifier should be unique")
		}
		if epoch.Duration == 0 {
			return errors.New("epoch duration should NOT be 0")
		}
		epochIdentifiers[epoch.Identifier] = true
	}
	return nil
}
