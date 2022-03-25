package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

//TODO : Impl functions related to batch
func NewIncomingMintTx(orchestratorAddress sdk.AccAddress, counter uint64) IncomingMintTx {
	return IncomingMintTx{
		OrchAddresses: []string{orchestratorAddress.String()},
		Counter:       counter,
	}
}

func NewAddressAndAmount(destinationAddress sdk.AccAddress, amount sdk.Coins, nativeBlockHeight int64) AddressAndAmountKey {
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

func NewProposalKey(chainID string, blockHeight int64, proposalID int64) ProposalKey {
	return ProposalKey{
		ChainID:     chainID,
		BlockHeight: blockHeight,
		ProposalID:  proposalID,
	}
}

func NewProposalValue(title string, description string, orchAddress string, ratio float32, votingStartTime time.Time, votingEndTime time.Time) ProposalValue {
	return ProposalValue{
		Title:                 title,
		Description:           description,
		OrchestratorAddresses: []string{orchAddress},
		Ratio:                 ratio,
		Counter:               1,
		ProposalPosted:        false,
		VotingStartTime:       votingStartTime,
		VotingEndTime:         votingEndTime,
	}
}

func NewTxHashValue(txId uint64, orchestratorAddress sdk.AccAddress, ratio float32, status string, nativeBlockHeight int64, activeBlockHeight int64) TxHashValue {
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
