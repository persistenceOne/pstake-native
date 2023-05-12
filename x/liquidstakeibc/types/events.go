package types

const (
	EventTypeLiquidStake   = "liquid-stake"
	EventTypeLiquidUnstake = "liquid-unstake"
	EventTypePacket        = "ics27_packet"
	EventTypeTimeout       = "timeout"

	AttributeAmount           = "amount"
	AttributeAmountReceived   = "received"
	AttributeDelegatorAddress = "address"
	AttributePstakeDepositFee = "pstake-deposit-fee"
	AttributePstakeUnstakeFee = "pstake-unstake-fee"
	AttributeUnstakeAmount    = "undelegation-amount"
	AttributeKeyAck           = "acknowledgement"
	AttributeKeyAckSuccess    = "success"
	AttributeKeyAckError      = "error"
	AttributeValueCategory    = ModuleName
)
