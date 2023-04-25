package types

import (
	ibctfrtypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
)

func (hc *HostChain) GetIBCDenom() string {
	return ibctfrtypes.ParseDenomTrace(ibctfrtypes.GetPrefixedDenom(hc.PortId, hc.ChannelId, hc.HostDenom)).IBCDenom()
}

func (hc *HostChain) GetValidator(operatorAddress string) (*Validator, bool) {
	for _, validator := range hc.Validators {
		if validator.OperatorAddress == operatorAddress {
			return validator, true
		}
	}

	return nil, false
}
