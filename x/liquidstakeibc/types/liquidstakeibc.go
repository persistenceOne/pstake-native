package types

// CurrentUnbondingEpoch computes and returns the current unbonding epoch to the next nearest
// multiple of UndelegationEpochNumberFactor
func CurrentUnbondingEpoch(epochNumber int64) int64 {
	if epochNumber%UndelegationEpochNumberFactor == 0 {
		return epochNumber
	}
	return epochNumber + UndelegationEpochNumberFactor - epochNumber%UndelegationEpochNumberFactor
}

// PreviousUnbondingEpoch computes and returns the previous unbonding epoch to the previous nearest
// multiple of UndelegationEpochNumberFactor
func PreviousUnbondingEpoch(epochNumber int64) int64 {
	if epochNumber%UndelegationEpochNumberFactor == 0 {
		return epochNumber - UndelegationEpochNumberFactor
	}
	return epochNumber - epochNumber%UndelegationEpochNumberFactor
}
