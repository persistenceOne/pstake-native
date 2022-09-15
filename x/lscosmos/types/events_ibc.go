package types

// IBC events
const (
	EventTypeTimeout       = "timeout"
	EventTypeLiquidStake   = "liquid-stake"
	EventTypeRedeem      = "redeem"
	EventTypeRewardBoost   = "reward-boost"
	EventTypeLiquidUnstake = "liquid-unstake"
	// this line is used by starport scaffolding # ibc/packet/event

	AttributeKeyAckSuccess    = "success"
	AttributeKeyAck           = "acknowledgement"
	AttributeKeyAckError      = "error"
	AttributeAmount           = "amount"
	AttributeAmountRecieved   = "received"
	AttributePstakeDepositFee = "pstake-deposit-fee"
	AttributeDelegatorAddress = "address"
	AttributeRewarderAddress  = "rewarder-address"
	AttributeValueCategory    = ModuleName
)
