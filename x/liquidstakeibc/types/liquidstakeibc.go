package types

// CurrentUnbondingEpoch computes and returns the current unbonding epoch to the next nearest
// multiple of the host chain Undelegation Factor
func CurrentUnbondingEpoch(factor, epochNumber int64) int64 {
	if epochNumber%factor == 0 {
		return epochNumber
	}
	return epochNumber + factor - epochNumber%factor
}

// PreviousUnbondingEpoch computes and returns the previous unbonding epoch to the previous nearest
// multiple of the host chain Undelegation Factor
func PreviousUnbondingEpoch(factor, epochNumber int64) int64 {
	if epochNumber%factor == 0 {
		return epochNumber - factor
	}
	return epochNumber - epochNumber%factor
}
