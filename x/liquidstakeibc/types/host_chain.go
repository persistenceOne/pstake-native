package types

func (hc *HostChain) GetValidator(operatorAddress string) (*Validator, bool) {
	for _, validator := range hc.Validators {
		if validator.OperatorAddress == operatorAddress {
			return validator, true
		}
	}

	return nil, false
}
