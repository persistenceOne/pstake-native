package types

// IBC events
const (
	EventTypeTimeout       = "timeout"
	EventTypeLiquidStake   = "liquid-stake"
	EventTypeRedeem        = "redeem"
	EventTypeRewardBoost   = "reward-boost"
	EventTypeLiquidUnstake = "liquid-unstake"
	EventTypeClaim         = "claim"
	// this line is used by starport scaffolding # ibc/packet/event

	AttributeKeyAckSuccess    = "success"
	AttributeKeyAck           = "acknowledgement"
	AttributeKeyAckError      = "error"
	AttributeAmount           = "amount"
	AttributeAmountRecieved   = "received"
	AttributeUnstakeAmount    = "undelegation-amount"
	AttributePstakeDepositFee = "pstake-deposit-fee"
	AttributePstakeRedeemFee  = "pstake-redeem-fee"
	AttributePstakeUnstakeFee = "pstake-unstake-fee"
	AttributeDelegatorAddress = "address"
	AttributeRewarderAddress  = "rewarder-address"
	AttributeClaimedAmount    = "claimed-amount"
	AttributeValueCategory    = ModuleName
)
