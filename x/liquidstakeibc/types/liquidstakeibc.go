package types

func IsUnbondingEpoch(factor, epochNumber int64) bool {
	return epochNumber%factor == 0
}

// CurrentUnbondingEpoch computes and returns the current unbonding epoch to the next nearest
// multiple of the host chain Undelegation Factor
func CurrentUnbondingEpoch(factor, epochNumber int64) int64 {
	if epochNumber%factor == 0 {
		return epochNumber
	}
	return epochNumber + factor - epochNumber%factor
}
