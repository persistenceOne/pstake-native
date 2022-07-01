package types

import sdk "github.com/cosmos/cosmos-sdk/types"

var _ DBHelper = &ProposalValue{}
var _ DBHelper = &TxHashValue{}
var _ DBHelper = &ValueUndelegateSuccessStore{}
var _ DBHelper = &OutgoingSignaturePoolValue{}
var _ DBHelper = &SlashingStoreValue{}
var _ DBHelper = &MintTokenStoreValue{}

// Find returns if the valAddress passed is present or not
func (m *ProposalValue) Find(valAddress string) bool {
	for _, address := range m.ValidatorAddresses {
		if address == valAddress {
			return true
		}
	}
	return false
}

// UpdateValues updates validator addresses array and total validator count
func (m *ProposalValue) UpdateValues(valAddress string, totalValidatorCount int64) {
	m.ValidatorAddresses = append(m.ValidatorAddresses, valAddress)
	m.Counter++
	m.Ratio = sdk.NewDec(m.Counter).Quo(sdk.NewDec(totalValidatorCount))
}

// Find returns if the valAddress passed is present or not
func (m *TxHashValue) Find(valAddress string) bool {
	for _, address := range m.ValidatorAddresses {
		if address == valAddress {
			return true
		}
	}
	return false
}

// UpdateValues updates validator addresses array and total validator count
func (m *TxHashValue) UpdateValues(valAddress string, totalValidatorCount int64) {
	m.ValidatorAddresses = append(m.ValidatorAddresses, valAddress)
	m.Counter++
	m.Ratio = sdk.NewDec(m.Counter).Quo(sdk.NewDec(totalValidatorCount))
}

// Find returns if the valAddress passed is present or not
func (m *ValueUndelegateSuccessStore) Find(valAddress string) bool {
	for _, address := range m.ValidatorAddresses {
		if address == valAddress {
			return true
		}
	}
	return false
}

// UpdateValues updates validator addresses array and total validator count
func (m *ValueUndelegateSuccessStore) UpdateValues(valAddress string, totalValidatorCount int64) {
	m.ValidatorAddresses = append(m.ValidatorAddresses, valAddress)
	m.Counter++
	m.Ratio = sdk.NewDec(m.Counter).Quo(sdk.NewDec(totalValidatorCount))
}

// Find returns if the valAddress passed is present or not
func (m *OutgoingSignaturePoolValue) Find(valAddress string) bool {
	for _, address := range m.ValidatorAddresses {
		if address == valAddress {
			return true
		}
	}
	return false
}

// UpdateValues updates validator addresses array and total validator count
func (m *OutgoingSignaturePoolValue) UpdateValues(valAddress string, _ int64) {
	m.ValidatorAddresses = append(m.ValidatorAddresses, valAddress)
	m.Counter++
}

// Find returns if the valAddress passed is present or not
func (m *SlashingStoreValue) Find(valAddress string) bool {
	for _, address := range m.ValidatorAddresses {
		if address == valAddress {
			return true
		}
	}
	return false
}

// UpdateValues updates validator addresses array and total validator count
func (m *SlashingStoreValue) UpdateValues(valAddress string, totalValidatorCount int64) {
	m.ValidatorAddresses = append(m.ValidatorAddresses, valAddress)
	m.Counter++
	m.Ratio = sdk.NewDec(m.Counter).Quo(sdk.NewDec(totalValidatorCount))
}

// Find returns if the valAddress passed is present or not
func (m *MintTokenStoreValue) Find(valAddress string) bool {
	for _, address := range m.ValidatorAddresses {
		if address == valAddress {
			return true
		}
	}
	return false
}

// UpdateValues updates validator addresses array and total validator count
func (m *MintTokenStoreValue) UpdateValues(valAddress string, totalValidatorCount int64) {
	m.ValidatorAddresses = append(m.ValidatorAddresses, valAddress)
	m.Counter++
	m.Ratio = sdk.NewDec(m.Counter).Quo(sdk.NewDec(totalValidatorCount))
}
