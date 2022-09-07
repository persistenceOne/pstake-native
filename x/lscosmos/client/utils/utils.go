package utils

import (
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

type RegisterHostChainProposalJSON struct {
	Title                 string                      `json:"title" yaml:"title"`
	Description           string                      `json:"description" yaml:"description"`
	ModuleEnabled         bool                        `json:"module_enabled" yaml:"module_enabled"`
	ChainID               string                      `json:"chain_id" yaml:"chain_id"`
	ConnectionID          string                      `json:"connection_id" yaml:"connection_id"`
	TransferChannel       string                      `json:"transfer_channel" yaml:"transfer_channel"`
	TransferPort          string                      `json:"transfer_port" yaml:"transfer_port"`
	BaseDenom             string                      `json:"base_denom" yaml:"base_denom"`
	MintDenom             string                      `json:"mint_denom" yaml:"mint_denom"`
	MinDeposit            string                      `json:"min_deposit" yaml:"min_deposit"`
	AllowListedValidators types.AllowListedValidators `json:"allow_listed_validators" yaml:"allow_listed_validators"`
	PstakeDepositFee      string                      `json:"pstake_deposit_fee" yaml:"pstake_deposit_fee"`
	PstakeRestakeFee      string                      `json:"pstake_restake_fee" yaml:"pstake_restake_fee"`
	PstakeUnstakeFee      string                      `json:"pstake_unstake_fee" yaml:"pstake_unstake_fee"`
	PstakeFeeAddress      string                      `json:"pstake_fee_address" yaml:"pstake_fee_address"`
	Deposit               string                      `json:"deposit" yaml:"deposit"`
}

func NewRegisterChainJSON(title, description string, moduleEnabled bool, chainID, connectionID, transferChannel, transferPort,
	baseDenom, mintDenom, minDeposit, pstakeFeeAddress, pstakeDepositFee, pstakeRestakeFee, pstakeUnstakeFee, deposit string, allowListedValidators types.AllowListedValidators) RegisterHostChainProposalJSON {
	return RegisterHostChainProposalJSON{
		Title:                 title,
		Description:           description,
		ModuleEnabled:         moduleEnabled,
		ChainID:               chainID,
		ConnectionID:          connectionID,
		TransferChannel:       transferChannel,
		TransferPort:          transferPort,
		BaseDenom:             baseDenom,
		MintDenom:             mintDenom,
		MinDeposit:            minDeposit,
		AllowListedValidators: allowListedValidators,
		PstakeDepositFee:      pstakeDepositFee,
		PstakeRestakeFee:      pstakeRestakeFee,
		PstakeUnstakeFee:      pstakeUnstakeFee,
		PstakeFeeAddress:      pstakeFeeAddress,
		Deposit:               deposit,
	}
}

// ParseRegisterHostChainProposalJSON reads and parses a RegisterHostChainProposalJSON from
// file.
func ParseRegisterHostChainProposalJSON(cdc *codec.LegacyAmino, proposalFile string) (RegisterHostChainProposalJSON, error) {
	proposal := RegisterHostChainProposalJSON{}

	contents, err := os.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}

type MinDepositAndFeeChangeProposalJSON struct {
	Title            string `json:"title" yaml:"title"`
	Description      string `json:"description" yaml:"description"`
	MinDeposit       string `json:"min_deposit" yaml:"min_deposit"`
	PstakeDepositFee string `json:"pstake_deposit_fee" yaml:"pstake_deposit_fee"`
	PstakeRestakeFee string `json:"pstake_restake_fee" yaml:"pstake_restake_fee"`
	PstakeUnstakeFee string `json:"pstake_unstake_fee" yaml:"pstake_unstake_fee"`
	Deposit          string `json:"deposit" yaml:"deposit"`
}

func NewMinDepositAndFeeChangeJSON(title, description, minDeposit, pstakeDepositFee, pstakeRestakeFee,
	pstakeUnstakeFee, deposit string) MinDepositAndFeeChangeProposalJSON {
	return MinDepositAndFeeChangeProposalJSON{
		Title:            title,
		Description:      description,
		MinDeposit:       minDeposit,
		PstakeDepositFee: pstakeDepositFee,
		PstakeRestakeFee: pstakeRestakeFee,
		PstakeUnstakeFee: pstakeUnstakeFee,
		Deposit:          deposit,
	}

}

// ParseMinDepositAndFeeChangeProposalJSON reads and parses a MinDepositAndFeeChangeProposal from
// file.
func ParseMinDepositAndFeeChangeProposalJSON(cdc *codec.LegacyAmino, proposalFile string) (MinDepositAndFeeChangeProposalJSON, error) {
	proposal := MinDepositAndFeeChangeProposalJSON{}

	contents, err := os.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}
	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}

type PstakeFeeAddressChangeProposalJSON struct {
	Title            string `json:"title" yaml:"title"`
	Description      string `json:"description" yaml:"description"`
	PstakeFeeAddress string `json:"pstake_fee_address" yaml:"pstake_fee_address"`
	Deposit          string `json:"deposit" yaml:"deposit"`
}

func NewPstakeFeeAddressChangeProposalJSON(title, description, pstakeFeeAddress, deposit string) PstakeFeeAddressChangeProposalJSON {
	return PstakeFeeAddressChangeProposalJSON{
		Title:            title,
		Description:      description,
		PstakeFeeAddress: pstakeFeeAddress,
		Deposit:          deposit,
	}

}

// ParsePstakeFeeAddressChangeProposalJSON reads and parses a PstakeFeeAddressChangeProposal  from
// file.
func ParsePstakeFeeAddressChangeProposalJSON(cdc *codec.LegacyAmino, proposalFile string) (PstakeFeeAddressChangeProposalJSON, error) {
	proposal := PstakeFeeAddressChangeProposalJSON{}

	contents, err := os.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}
	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
