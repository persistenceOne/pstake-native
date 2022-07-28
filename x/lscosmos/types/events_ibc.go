package types

// IBC events
const (
	EventTypeTimeout = "timeout"
	EventTypeMint    = "mint-tokens"
	// this line is used by starport scaffolding # ibc/packet/event

	AttributeKeyAckSuccess = "success"
	AttributeKeyAck        = "acknowledgement"
	AttributeKeyAckError   = "error"
	AttributeAmountMinted  = "amount"
	AttributeMintedAddress = "address"
)
