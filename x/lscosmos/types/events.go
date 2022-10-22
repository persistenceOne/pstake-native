package types

// IBC events
const (
	EventTypePacket        = "ics27_packet"
	EventTypeTimeout       = "timeout"
	EventTypeLiquidStake   = "liquid-stake"
	EventTypeRedeem        = "redeem"
	EventTypeRewardBoost   = "reward-boost"
	EventTypeLiquidUnstake = "liquid-unstake"
	EventTypeClaim         = "claim"
	EventTypeJumpStart     = "jump-start"
	// this line is used by starport scaffolding # ibc/packet/event

	AttributeKeyAckSuccess    = "success"
	AttributeKeyAck           = "acknowledgement"
	AttributeKeyAckError      = "error"
	AttributeAmount           = "amount"
	AttributeAmountReceived   = "received"
	AttributeUnstakeAmount    = "undelegation-amount"
	AttributePstakeDepositFee = "pstake-deposit-fee"
	AttributePstakeRedeemFee  = "pstake-redeem-fee"
	AttributePstakeUnstakeFee = "pstake-unstake-fee"
	AttributeDelegatorAddress = "address"
	AttributeRewarderAddress  = "rewarder-address"
	AttributeClaimedAmount    = "claimed-amount"
	AttributePstakeAddress    = "pstake-address"
	AttributeValueCategory    = ModuleName
)
