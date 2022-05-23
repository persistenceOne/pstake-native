package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func NewIncomingMintTx(orchestratorAddress sdkTypes.AccAddress, ratio sdkTypes.Dec) IncomingMintTx {
	return IncomingMintTx{
		OrchAddresses: []string{orchestratorAddress.String()},
		Counter:       1,
		Ratio:         ratio,
	}
}

func NewAddressAndAmount(destinationAddress sdkTypes.AccAddress, amount sdkTypes.Coin, nativeBlockHeight int64) AddressAndAmountKey {
	return AddressAndAmountKey{
		DestinationAddress: destinationAddress.String(),
		Amount:             amount,
		Acknowledgment:     false,
		Minted:             false,
		NativeBlockHeight:  nativeBlockHeight,
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

func NewProposalValue(msg MsgMakeProposal, orchAddress string, ratio sdkTypes.Dec) ProposalValue {
	return ProposalValue{
		ProposalDetails:       msg,
		OrchestratorAddresses: []string{orchAddress},
		ProposalPosted:        false,
		Ratio:                 ratio,
		Counter:               1,
	}
}

func NewTxHashValue(msg MsgTxStatus, ratio sdkTypes.Dec, nativeBlockHeight, activeBlockHeight int64) TxHashValue {
	return TxHashValue{
		TxStatus:              msg,
		OrchestratorAddresses: []string{msg.OrchestratorAddress},
		Status:                []string{msg.Status},
		Ratio:                 ratio,
		TxCleared:             false,
		Counter:               1,
		NativeBlockHeight:     nativeBlockHeight,
		ActiveBlockHeight:     activeBlockHeight,
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

func NewValueUndelegateSuccessStore(msg MsgUndelegateSuccess, orchestratorAddress string, ratio sdkTypes.Dec,
	nativeBlockHeight, activeBlockHeight int64) ValueUndelegateSuccessStore {
	return ValueUndelegateSuccessStore{
		UndelegateSuccess:     msg,
		OrchestratorAddresses: []string{orchestratorAddress},
		Ratio:                 ratio,
		Counter:               1,
		NativeBlockHeight:     nativeBlockHeight,
		ActiveBlockHeight:     activeBlockHeight,
	}
}

func NewStakingEpochValue(keyAndValue KeyAndValueForMinting) StakingEpochValue {
	return StakingEpochValue{
		EpochMintingTxns: []KeyAndValueForMinting{keyAndValue},
	}
}

func NewMintingEpochValue(txIDAndStatus MintingEpochValueMember) MintingEpochValue {
	return MintingEpochValue{
		TxIDAndStatus: []MintingEpochValueMember{txIDAndStatus},
	}
}

func NewRewardsClaimedValue(msg MsgRewardsClaimedOnCosmosChain, orchestratorAddress string, ratio sdkTypes.Dec,
	nativeBlockHeight, activeBlockHeight int64) RewardsClaimedValue {
	return RewardsClaimedValue{
		RewardsClaimed:        msg,
		OrchestratorAddresses: []string{orchestratorAddress},
		Ratio:                 ratio,
		Counter:               1,
		AddedToCurrentEpoch:   false,
		NativeBlockHeight:     nativeBlockHeight,
		ActiveBlockHeight:     activeBlockHeight,
	}
}

func NewValidatorStoreValue(orchAddress sdkTypes.AccAddress) ValidatorStoreValue {
	return ValidatorStoreValue{
		OrchestratorAddresses: []string{orchAddress.String()},
	}
}

func NewOutgoingSignaturePoolValue(singleSignature SingleSignatureDataForOutgoingPool, orchestratorAddress sdkTypes.AccAddress) OutgoingSignaturePoolValue {
	return OutgoingSignaturePoolValue{
		SingleSignatures:      []SingleSignatureDataForOutgoingPool{singleSignature},
		OrchestratorAddresses: []string{orchestratorAddress.String()},
		Counter:               1,
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

func NewSlashingStoreValue(msg MsgSlashingEventOnCosmosChain, ratio sdkTypes.Dec, orch string) SlashingStoreValue {
	return SlashingStoreValue{
		SlashingDetails:       msg,
		Ratio:                 ratio,
		OrchestratorAddresses: []string{orch},
		Counter:               1,
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
