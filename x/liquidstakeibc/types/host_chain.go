package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctfrtypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
)

func (hc *HostChain) IBCDenom() string {
	return ibctfrtypes.ParseDenomTrace(ibctfrtypes.GetPrefixedDenom(hc.PortId, hc.ChannelId, hc.HostDenom)).IBCDenom()
}

func (hc *HostChain) MintDenom() string {
	return "stk" + "/" + hc.HostDenom
}

func (hc *HostChain) GetValidator(operatorAddress string) (*Validator, bool) {
	for _, validator := range hc.Validators {
		if validator.OperatorAddress == operatorAddress {
			return validator, true
		}
	}

	return nil, false
}

func (hc *HostChain) GetHostChainTotalDelegations() sdk.Int {
	totalDelegations := sdk.ZeroInt()
	for _, validator := range hc.Validators {
		totalDelegations.Add(validator.DelegatedAmount)
	}

	return totalDelegations
}
