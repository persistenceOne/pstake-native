package types

import sdk "github.com/cosmos/cosmos-sdk/types"

var _ DBHelper = &ProposalValue{}
var _ DBHelper = &TxHashValue{}
var _ DBHelper = &RewardsClaimedValue{}
var _ DBHelper = &ValueUndelegateSuccessStore{}
var _ DBHelper = &OutgoingSignaturePoolValue{}
var _ DBHelper = &SlashingStoreValue{}
var _ DBHelper = &MintTokenStoreValue{}

func (m *ProposalValue) Find(valAddress string) bool {
	for _, address := range m.ValidatorAddresses {
		if address == valAddress {
			return true
		}
	}
	return false
}

func (m *ProposalValue) UpdateValues(valAddress string, totalValidatorCount int64) {
	m.ValidatorAddresses = append(m.ValidatorAddresses, valAddress)
	m.Counter++
	m.Ratio = sdk.NewDec(m.Counter).Quo(sdk.NewDec(totalValidatorCount))
}

func (m *TxHashValue) Find(valAddress string) bool {
	for _, address := range m.ValidatorAddresses {
		if address == valAddress {
			return true
		}
	}
	return false
}

func (m *TxHashValue) UpdateValues(valAddress string, totalValidatorCount int64) {
	m.ValidatorAddresses = append(m.ValidatorAddresses, valAddress)
	m.Counter++
	m.Ratio = sdk.NewDec(m.Counter).Quo(sdk.NewDec(totalValidatorCount))
}

func (m *RewardsClaimedValue) Find(valAddress string) bool {
	for _, address := range m.ValidatorAddresses {
		if address == valAddress {
			return true
		}
	}
	return false
}

func (m *RewardsClaimedValue) UpdateValues(valAddress string, totalValidatorCount int64) {
	m.ValidatorAddresses = append(m.ValidatorAddresses, valAddress)
	m.Counter++
	m.Ratio = sdk.NewDec(m.Counter).Quo(sdk.NewDec(totalValidatorCount))
}

func (m *ValueUndelegateSuccessStore) Find(valAddress string) bool {
	for _, address := range m.ValidatorAddresses {
		if address == valAddress {
			return true
		}
	}
	return false
}

func (m *ValueUndelegateSuccessStore) UpdateValues(valAddress string, totalValidatorCount int64) {
	m.ValidatorAddresses = append(m.ValidatorAddresses, valAddress)
	m.Counter++
	m.Ratio = sdk.NewDec(m.Counter).Quo(sdk.NewDec(totalValidatorCount))
}

func (m *OutgoingSignaturePoolValue) Find(valAddress string) bool {
	for _, address := range m.ValidatorAddresses {
		if address == valAddress {
			return true
		}
	}
	return false
}

func (m *OutgoingSignaturePoolValue) UpdateValues(valAddress string, _ int64) {
	m.ValidatorAddresses = append(m.ValidatorAddresses, valAddress)
	m.Counter++
}

func (m *SlashingStoreValue) Find(valAddress string) bool {
	for _, address := range m.ValidatorAddresses {
		if address == valAddress {
			return true
		}
	}
	return false
}

func (m *SlashingStoreValue) UpdateValues(valAddress string, totalValidatorCount int64) {
	m.ValidatorAddresses = append(m.ValidatorAddresses, valAddress)
	m.Counter++
	m.Ratio = sdk.NewDec(m.Counter).Quo(sdk.NewDec(totalValidatorCount))
}

func (m *MintTokenStoreValue) Find(valAddress string) bool {
	for _, address := range m.ValidatorAddresses {
		if address == valAddress {
			return true
		}
	}
	return false
}

func (m *MintTokenStoreValue) UpdateValues(valAddress string, totalValidatorCount int64) {
	m.ValidatorAddresses = append(m.ValidatorAddresses, valAddress)
	m.Counter++
	m.Ratio = sdk.NewDec(m.Counter).Quo(sdk.NewDec(totalValidatorCount))
}
