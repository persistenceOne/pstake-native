package types

// IBC events
const (
	EventTypeTimeout     = "timeout"
	EventTypeLiquidStake = "liquid-stake"
	// this line is used by starport scaffolding # ibc/packet/event

	AttributeKeyAckSuccess    = "success"
	AttributeKeyAck           = "acknowledgement"
	AttributeKeyAckError      = "error"
	AttributeAmountMinted     = "amount"
	AttributeDelegatorAddress = "address"
	AttributeValueCategory    = ModuleName
)
