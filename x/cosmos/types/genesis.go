package types

import sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"

// NewGenesisState todo fill up incoming and outgoing trsactions array and maintain it with a length of 10000 txns
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

func (data GenesisState) Equal(other GenesisState) bool {
	return data.Params.Equal(other.Params)
}
