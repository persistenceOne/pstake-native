package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// NewMintTokenStoreValue returns MintTokenStoreValue struct
func NewMintTokenStoreValue(msg MsgMintTokensForAccount, ratio sdkTypes.Dec, valAddress sdkTypes.ValAddress,
	activeBlockHeight int64) MintTokenStoreValue {
	return MintTokenStoreValue{
		MintTokens:         msg,
		Ratio:              ratio,
		ValidatorAddresses: []string{valAddress.String()},
		Counter:            1,
		Minted:             false,
		AddedToEpoch:       false,
		ActiveBlockHeight:  activeBlockHeight,
	}
}

// NewChainIDHeightAndTxHash returns ChainIDHeightAndTxHashKey struct
func NewChainIDHeightAndTxHash(chainID string, blockHeight int64, txHash string) ChainIDHeightAndTxHashKey {
	return ChainIDHeightAndTxHashKey{
		ChainID:     chainID,
		BlockHeight: blockHeight,
		TxHash:      txHash,
	}
}

// NewProposalKey returns ProposalKey struct
func NewProposalKey(chainID string, blockHeight int64, proposalID uint64) ProposalKey {
	return ProposalKey{
		ChainID:     chainID,
		BlockHeight: blockHeight,
		ProposalID:  proposalID,
	}
}

// NewProposalValue returns ProposalValue struct
func NewProposalValue(msg MsgMakeProposal, valAddress sdkTypes.ValAddress, ratio sdkTypes.Dec,
	blockHeight int64) ProposalValue {
	return ProposalValue{
		ProposalDetails:    msg,
		ValidatorAddresses: []string{valAddress.String()},
		ProposalPosted:     false,
		Ratio:              ratio,
		Counter:            1,
		ActiveBlockHeight:  blockHeight,
	}
}

// NewTxHashValue returns TxHashValue struct
func NewTxHashValue(msg MsgTxStatus, ratio sdkTypes.Dec, activeBlockHeight int64,
	valAddress sdkTypes.ValAddress) TxHashValue {
	return TxHashValue{
		TxStatus:           msg,
		ValidatorAddresses: []string{valAddress.String()},
		Status:             []string{msg.Status},
		Ratio:              ratio,
		TxCleared:          false,
		Counter:            1,
		ActiveBlockHeight:  activeBlockHeight,
	}
}

// NewWithdrawStoreValue returns WithdrawStoreValue struct
func NewWithdrawStoreValue(msg MsgWithdrawStkAsset) WithdrawStoreValue {
	return WithdrawStoreValue{
		WithdrawDetails: []MsgWithdrawStkAsset{msg},
		UnbondEmitFlag:  []bool{false},
	}
}

// NewValueOutgoingUnbondStore returns ValueOutgoingUnbondStore struct
func NewValueOutgoingUnbondStore(undelegateMessage []stakingTypes.MsgUndelegate,
	epochNumber int64, cValue sdkTypes.Dec) ValueOutgoingUnbondStore {
	return ValueOutgoingUnbondStore{
		EpochNumber:        epochNumber,
		UndelegateMessages: undelegateMessage,
		CValue:             cValue,
	}
}

// NewValueUndelegateSuccessStore returns ValueUndelegateSuccessStore struct
func NewValueUndelegateSuccessStore(msg MsgUndelegateSuccess, valAddress sdkTypes.ValAddress, ratio sdkTypes.Dec,
	activeBlockHeight int64) ValueUndelegateSuccessStore {
	return ValueUndelegateSuccessStore{
		UndelegateSuccess:  msg,
		ValidatorAddresses: []string{valAddress.String()},
		Ratio:              ratio,
		Counter:            1,
		ActiveBlockHeight:  activeBlockHeight,
	}
}

// NewValidatorStoreValue returns ValidatorStoreValue struct
func NewValidatorStoreValue(orchAddress sdkTypes.AccAddress) ValidatorStoreValue {
	return ValidatorStoreValue{
		OrchestratorAddresses: []string{orchAddress.String()},
	}
}

// NewOutgoingSignaturePoolValue returns OutgoingSignaturePoolValue struct
func NewOutgoingSignaturePoolValue(singleSignature SingleSignatureDataForOutgoingPool, valAddress sdkTypes.ValAddress,
	orchestratorAddress sdkTypes.AccAddress, activeBlockHeight int64) OutgoingSignaturePoolValue {
	return OutgoingSignaturePoolValue{
		SingleSignatures:      []SingleSignatureDataForOutgoingPool{singleSignature},
		ValidatorAddresses:    []string{valAddress.String()},
		Counter:               1,
		OrchestratorAddresses: []string{orchestratorAddress.String()},
		SignedEventEmitted:    false,
		ActiveBlockHeight:     activeBlockHeight,
	}
}

func NewEpochWithdrawSuccessStoreValue(cValue sdkTypes.Dec) EpochWithdrawSuccessStoreValue {
	return EpochWithdrawSuccessStoreValue{
		Status: false,
		CValue: cValue,
	}
}

// ConvertSingleSignatureDataToSingleSignatureDataForOutgoingPool returns SingleSignatureDataForOutgoingPool struct
func ConvertSingleSignatureDataToSingleSignatureDataForOutgoingPool(data signing.SingleSignatureData) SingleSignatureDataForOutgoingPool {
	return SingleSignatureDataForOutgoingPool{
		SignMode:  data.SignMode,
		Signature: data.Signature,
	}
}

// ConvertSingleSignatureDataForOutgoingPoolToSingleSignatureData returns signing.SingleSignatureData struct
func ConvertSingleSignatureDataForOutgoingPoolToSingleSignatureData(data SingleSignatureDataForOutgoingPool) signing.SingleSignatureData {
	return signing.SingleSignatureData{
		SignMode:  data.SignMode,
		Signature: data.Signature,
	}
}

// NewOutgoingQueueValue returns OutgoingQueueValue struct
func NewOutgoingQueueValue(active bool, retryCounter uint64) OutgoingQueueValue {
	return OutgoingQueueValue{
		Active:       active,
		RetryCounter: retryCounter,
	}
}

// NewSlashingStoreValue returns SlashingStoreValue struct
func NewSlashingStoreValue(msg MsgSlashingEventOnCosmosChain, ratio sdkTypes.Dec, valAddress sdkTypes.ValAddress,
	blockHeight int64) SlashingStoreValue {
	return SlashingStoreValue{
		SlashingDetails:    msg,
		Ratio:              ratio,
		ValidatorAddresses: []string{valAddress.String()},
		Counter:            1,
		AddedToCValue:      false,
		ActiveBlockHeight:  blockHeight,
	}
}

// SetSignatures sets signatures for CosmosTx, takes array of signatures as input and aggregate them together
func (c *CosmosTx) SetSignatures(signatures ...signing.SignatureV2) error {
	n := len(signatures)
	signerInfos := make([]*tx.SignerInfo, n)
	rawSigs := make([][]byte, n)

	for i, sig := range signatures {
		var modeInfo *tx.ModeInfo
		modeInfo, rawSigs[i] = authtx.SignatureDataToModeInfoAndSig(sig.Data)
		any, err := codectypes.NewAnyWithValue(sig.PubKey)
		if err != nil {
			return err
		}
		signerInfos[i] = &tx.SignerInfo{
			PublicKey: any,
			ModeInfo:  modeInfo,
			Sequence:  sig.Sequence,
		}
	}

	c.setSignerInfos(signerInfos)
	c.setSignatures(rawSigs)

	return nil
}

// setSignerInfos Sets signer infos
func (c *CosmosTx) setSignerInfos(infos []*tx.SignerInfo) {
	c.Tx.AuthInfo.SignerInfos = infos
}

// setSignatures Sets signatures
func (c *CosmosTx) setSignatures(sigs [][]byte) {
	c.Tx.Signatures = sigs
}

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
	m.Ratio = sdkTypes.NewDec(m.Counter).Quo(sdkTypes.NewDec(totalValidatorCount))
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
	m.Ratio = sdkTypes.NewDec(m.Counter).Quo(sdkTypes.NewDec(totalValidatorCount))
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
	m.Ratio = sdkTypes.NewDec(m.Counter).Quo(sdkTypes.NewDec(totalValidatorCount))
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
	m.Ratio = sdkTypes.NewDec(m.Counter).Quo(sdkTypes.NewDec(totalValidatorCount))
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
	m.Ratio = sdkTypes.NewDec(m.Counter).Quo(sdkTypes.NewDec(totalValidatorCount))
}
