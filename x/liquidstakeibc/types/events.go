package types

const (
	EventTypeLiquidStake = "liquid-stake"
	EventTypePacket      = "ics27_packet"
	EventTypeTimeout     = "timeout"

	AttributeAmount           = "amount"
	AttributeAmountReceived   = "received"
	AttributeDelegatorAddress = "address"
	AttributePstakeDepositFee = "pstake-deposit-fee"
	AttributeKeyAck           = "acknowledgement"
	AttributeKeyAckSuccess    = "success"
	AttributeKeyAckError      = "error"
	AttributeValueCategory    = ModuleName
)
