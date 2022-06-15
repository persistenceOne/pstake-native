package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func NewMintTokenStoreValue(msg MsgMintTokensForAccount, ratio sdkTypes.Dec, valAddress sdkTypes.ValAddress, activeBlockHeight int64) MintTokenStoreValue {
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

func NewChainIDHeightAndTxHash(chainID string, blockHeight int64, txHash string) ChainIDHeightAndTxHashKey {
	return ChainIDHeightAndTxHashKey{
		ChainID:     chainID,
		BlockHeight: blockHeight,
		TxHash:      txHash,
	}
}

func NewProposalKey(chainID string, blockHeight int64, proposalID uint64) ProposalKey {
	return ProposalKey{
		ChainID:     chainID,
		BlockHeight: blockHeight,
		ProposalID:  proposalID,
	}
}

func NewProposalValue(msg MsgMakeProposal, valAddress sdkTypes.ValAddress, ratio sdkTypes.Dec, blockHeight int64) ProposalValue {
	return ProposalValue{
		ProposalDetails:    msg,
		ValidatorAddresses: []string{valAddress.String()},
		ProposalPosted:     false,
		Ratio:              ratio,
		Counter:            1,
		ActiveBlockHeight:  blockHeight,
	}
}

func NewTxHashValue(msg MsgTxStatus, ratio sdkTypes.Dec, activeBlockHeight int64, valAddress sdkTypes.ValAddress) TxHashValue {
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

func NewWithdrawStoreValue(msg MsgWithdrawStkAsset) WithdrawStoreValue {
	return WithdrawStoreValue{
		WithdrawDetails: []MsgWithdrawStkAsset{msg},
		UnbondEmitFlag:  []bool{false},
	}
}

func NewValueOutgoingUnbondStore(undelegateMessage []stakingTypes.MsgUndelegate, epochNumber int64) ValueOutgoingUnbondStore {
	return ValueOutgoingUnbondStore{
		EpochNumber:        epochNumber,
		UndelegateMessages: undelegateMessage,
	}
}

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

func NewMintingEpochValue(txIDAndStatus MintingEpochValueMember) MintingEpochValue {
	return MintingEpochValue{
		TxIDAndStatus: []MintingEpochValueMember{txIDAndStatus},
	}
}

func NewRewardsClaimedValue(msg MsgRewardsClaimedOnCosmosChain, valAddress sdkTypes.ValAddress, ratio sdkTypes.Dec,
	activeBlockHeight int64) RewardsClaimedValue {
	return RewardsClaimedValue{
		RewardsClaimed:      msg,
		ValidatorAddresses:  []string{valAddress.String()},
		Ratio:               ratio,
		Counter:             1,
		AddedToCurrentEpoch: false,
		ActiveBlockHeight:   activeBlockHeight,
	}
}

func NewValidatorStoreValue(orchAddress sdkTypes.AccAddress) ValidatorStoreValue {
	return ValidatorStoreValue{
		OrchestratorAddresses: []string{orchAddress.String()},
	}
}

func NewOutgoingSignaturePoolValue(singleSignature SingleSignatureDataForOutgoingPool, valAddress sdkTypes.ValAddress) OutgoingSignaturePoolValue {
	return OutgoingSignaturePoolValue{
		SingleSignatures:   []SingleSignatureDataForOutgoingPool{singleSignature},
		ValidatorAddresses: []string{valAddress.String()},
		Counter:            1,
	}
}

func ConvertSingleSignatureDataToSingleSignatureDataForOutgoingPool(data signing.SingleSignatureData) SingleSignatureDataForOutgoingPool {
	return SingleSignatureDataForOutgoingPool{
		SignMode:  data.SignMode,
		Signature: data.Signature,
	}
}

func ConvertSingleSignatureDataForOutgoingPoolToSingleSignatureData(data SingleSignatureDataForOutgoingPool) signing.SingleSignatureData {
	return signing.SingleSignatureData{
		SignMode:  data.SignMode,
		Signature: data.Signature,
	}
}

func NewOutgoingQueueValue(active bool, retryCounter uint64) OutgoingQueueValue {
	return OutgoingQueueValue{
		Active:       active,
		RetryCounter: retryCounter,
	}
}

func NewSlashingStoreValue(msg MsgSlashingEventOnCosmosChain, ratio sdkTypes.Dec, valAddress sdkTypes.ValAddress, blockHeight int64) SlashingStoreValue {
	return SlashingStoreValue{
		SlashingDetails:    msg,
		Ratio:              ratio,
		ValidatorAddresses: []string{valAddress.String()},
		Counter:            1,
		AddedToCValue:      false,
		ActiveBlockHeight:  blockHeight,
	}
}

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

func (c *CosmosTx) setSignerInfos(infos []*tx.SignerInfo) {
	c.Tx.AuthInfo.SignerInfos = infos
}

func (c *CosmosTx) setSignatures(sigs [][]byte) {
	c.Tx.Signatures = sigs
}
