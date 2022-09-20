package types

// IBC events
const (
	EventTypeTimeout     = "timeout"
	EventTypeLiquidStake = "liquid-stake"
	EventTypeRedeem      = "redeem"
	EventTypeRewardBoost = "reward-boost"
	// this line is used by starport scaffolding # ibc/packet/event

	AttributeKeyAckSuccess    = "success"
	AttributeKeyAck           = "acknowledgement"
	AttributeKeyAckError      = "error"
	AttributeAmountMinted     = "amount"
	AttributeAmountRecieved   = "received"
	AttributePstakeDepositFee = "pstake-deposit-fee"
	AttributeDelegatorAddress = "address"
	AttributeRewarderAddress  = "rewarder-address"
	AttributeRedeemAddress    = "redeem-address"
	AttributeAmountRedeemed   = "redeem-amount"
	AttributeValueCategory    = ModuleName
)
