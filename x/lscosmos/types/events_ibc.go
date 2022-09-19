package types

// IBC events
const (
	EventTypeTimeout     = "timeout"
	EventTypeLiquidStake = "liquid-stake"
	EventTypeRewardBoost = "reward-boost"
	// this line is used by starport scaffolding # ibc/packet/event

	AttributeKeyAckSuccess     = "success"
	AttributeKeyAck            = "acknowledgement"
	AttributeKeyAckError       = "error"
	AttributeAmountMinted      = "amount"
	AttributeAmountRecieved    = "received"
	AttributePstakeDepositFee  = "pstake-deposit-fee"
	AttributeDelegatorAddress  = "address"
	AttributeRewarderAddress   = "rewarder-address"
	AttributeWithdrawerAddress = "withdraw-address"
	AttributeAmountWithdrawn   = "withdraw-amount"
	AttributeValueCategory     = ModuleName
)
