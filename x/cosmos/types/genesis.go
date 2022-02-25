package types

import sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"

func NewGenesisState(params Params, delegationCosmos []DelegationCosmos, incomingTx []IncomingTx, outgoingTx OutgoingTx) *GenesisState {
	return &GenesisState{
		Params:            params,
		CosmosDelegations: delegationCosmos,
		IncomingTxn:       incomingTx,
		OutgoingTxn:       outgoingTx,
	}
}

func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return sdkErrors.Wrap(err, "params")
	}
	return nil
}

func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:            DefaultParams(),
		CosmosDelegations: []DelegationCosmos{},
		IncomingTxn:       []IncomingTx{},
		OutgoingTxn:       OutgoingTx{},
	}
}
