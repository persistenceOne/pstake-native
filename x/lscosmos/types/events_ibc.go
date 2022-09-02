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
	AttributeAmountRecieved   = "received"
	AttributePstakeDepositFee = "pstake-deposit-fee"
	AttributeDelegatorAddress = "address"
	AttributeValueCategory    = ModuleName
)
