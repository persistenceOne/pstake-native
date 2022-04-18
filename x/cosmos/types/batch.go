package types

import (
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"time"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

func NewIncomingMintTx(orchestratorAddress sdkTypes.AccAddress, counter uint64) IncomingMintTx {
	return IncomingMintTx{
		OrchAddresses: []string{orchestratorAddress.String()},
		Counter:       counter,
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

func NewProposalValue(title string, description string, orchAddress string, ratio float32, votingStartTime time.Time, votingEndTime time.Time, cosmosProposalID uint64) ProposalValue {
	return ProposalValue{
		Title:                 title,
		Description:           description,
		OrchestratorAddresses: []string{orchAddress},
		Ratio:                 ratio,
		Counter:               1,
		ProposalPosted:        false,
		VotingStartTime:       votingStartTime,
		VotingEndTime:         votingEndTime,
		CosmosProposalID:      cosmosProposalID,
	}
}

func NewTxHashValue(txId uint64, orchestratorAddress sdkTypes.AccAddress, ratio float32, status string, nativeBlockHeight int64, activeBlockHeight int64) TxHashValue {
	return TxHashValue{
		TxID:                  txId,
		OrchestratorAddresses: []string{orchestratorAddress.String()},
		Status:                []string{status},
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

func NewValuOutgoingUnbondStore(undelegateMessage []stakingTypes.MsgUndelegate, epochNumber int64) ValuOutgoingUnbondStore {
	return ValuOutgoingUnbondStore{
		EpochNumber:        epochNumber,
		UndelegateMessages: undelegateMessage,
	}
}

func NewValidatorStoreValue(orchestratorAddress sdkTypes.AccAddress) ValidatorStoreValue {
	return ValidatorStoreValue{
		OrchestratorAddresses: []string{orchestratorAddress.String()},
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

func NewRewardsClaimedValue(orchestratorAddress sdkTypes.AccAddress, amount sdkTypes.Coin, ratio float32, nativeBlockHeight int64, activeBlockHeight int64) RewardsClaimedValue {
	return RewardsClaimedValue{
		OrchestratorAddresses: []string{orchestratorAddress.String()},
		Amount:                []sdkTypes.Coin{amount},
		Ratio:                 ratio,
		Counter:               1,
		AddedToCurrentEpoch:   false,
		NativeBlockHeight:     nativeBlockHeight,
		ActiveBlockHeight:     activeBlockHeight,
	}
}
