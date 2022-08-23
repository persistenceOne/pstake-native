package utils

import (
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
)

type RegisterCosmosChainProposalReq struct {
	BaseReq     rest.BaseReq `json:"base_req" yaml:"base_req"`
	Title       string       `json:"title" yaml:"title"`
	Description string       `json:"description" yaml:"description"`
	// TODO update
	ModuleEnabled         bool                        `json:"module_enabled" yaml:"module_enabled"`
	IBCConnection         string                      `json:"ibc_connection" yaml:"ibc_connection"`
	TokenTransferChannel  string                      `json:"token_transfer_channel" yaml:"token_transfer_channel"`
	TokenTransferPort     string                      `json:"token_transfer_port" yaml:"token_transfer_port"`
	BaseDenom             string                      `json:"base_denom" yaml:"base_denom"`
	MintDenom             string                      `json:"mint_denom" yaml:"mint_denom"`
	MinDeposit            string                      `json:"min_deposit" yaml:"min_deposit"`
	AllowListedValidators types.AllowListedValidators `json:"allow_listed_validators" yaml:"allow_listed_validators"`
	PStakeDepositFee      string                      `json:"p_stake_deposit_fee" yaml:"p_stake_deposit_fee"`
	PStakeRestakeFee      string                      `json:"p_stake_restake_fee" yaml:"p_stake_restake_fee"`
	PStakeUnstakeFee      string                      `json:"p_stake_unstake_fee" yaml:"p_stake_unstake_fee"`
	Proposer              sdk.AccAddress              `json:"proposer" yaml:"proposer"`
	Deposit               sdk.Coins                   `json:"deposit" yaml:"deposit"`
}

type RegisterCosmosChainProposalJSON struct {
	Title                 string                      `json:"title" yaml:"title"`
	Description           string                      `json:"description" yaml:"description"`
	ModuleEnabled         bool                        `json:"module_enabled" yaml:"module_enabled"`
	IBCConnection         string                      `json:"ibc_connection" yaml:"ibc_connection"`
	TokenTransferChannel  string                      `json:"token_transfer_channel" yaml:"token_transfer_channel"`
	TokenTransferPort     string                      `json:"token_transfer_port" yaml:"token_transfer_port"`
	BaseDenom             string                      `json:"base_denom" yaml:"base_denom"`
	MintDenom             string                      `json:"mint_denom" yaml:"mint_denom"`
	MinDeposit            string                      `json:"min_deposit" yaml:"min_deposit"`
	AllowListedValidators types.AllowListedValidators `json:"allow_listed_validators" yaml:"allow_listed_validators"`
	PStakeDepositFee      string                      `json:"p_stake_deposit_fee" yaml:"p_stake_deposit_fee"`
	PStakeRestakeFee      string                      `json:"p_stake_restake_fee" yaml:"p_stake_restake_fee"`
	PStakeUnstakeFee      string                      `json:"p_stake_unstake_fee" yaml:"p_stake_unstake_fee"`
	Deposit               string                      `json:"deposit" yaml:"deposit"`
}

func NewRegisterChainJSON(title, description string, moduleEnabled bool, ibcConnectionID, tokenTransferChannel, tokenTransferPort,
	baseDenom, mintDenom, minDeposit string, allowListedValidators types.AllowListedValidators, pstakeDepositFee, pstakeRestakeFee, pstakeUnstakeFee, deposit string) RegisterCosmosChainProposalJSON {
	return RegisterCosmosChainProposalJSON{
		Title:                 title,
		Description:           description,
		ModuleEnabled:         moduleEnabled,
		IBCConnection:         ibcConnectionID,
		TokenTransferChannel:  tokenTransferChannel,
		TokenTransferPort:     tokenTransferPort,
		BaseDenom:             baseDenom,
		MintDenom:             mintDenom,
		MinDeposit:            minDeposit,
		AllowListedValidators: allowListedValidators,
		PStakeDepositFee:      pstakeDepositFee,
		PStakeRestakeFee:      pstakeRestakeFee,
		PStakeUnstakeFee:      pstakeUnstakeFee,
		Deposit:               deposit,
	}
}

// ParseRegisterCosmosChainProposalJSON reads and parses a RegisterCosmosChainProposalJSON from
// file.
func ParseRegisterCosmosChainProposalJSON(cdc *codec.LegacyAmino, proposalFile string) (RegisterCosmosChainProposalJSON, error) {
	proposal := RegisterCosmosChainProposalJSON{}

	contents, err := os.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err := cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
