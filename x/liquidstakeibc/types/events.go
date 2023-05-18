package types

const (
	EventTypeLiquidStake   = "liquid-stake"
	EventTypeLiquidUnstake = "liquid-unstake"
	EventTypeRedeem        = "redeem"
	EventTypePacket        = "ics27_packet"
	EventTypeTimeout       = "timeout"
	EventTypeSlashing      = "slashing"

	AttributeAmount             = "amount"
	AttributeAmountReceived     = "received"
	AttributeDelegatorAddress   = "address"
	AttributePstakeDepositFee   = "pstake-deposit-fee"
	AttributePstakeUnstakeFee   = "pstake-unstake-fee"
	AttributePstakeRedeemFee    = "pstake-redeem-fee"
	AttributeUnstakeAmount      = "undelegation-amount"
	AttributeUnstakeEpoch       = "undelegation-epoch"
	AttributeValidatorAddress   = "validator-address"
	AttributeExistingDelegation = "existing-delegation"
	AttributeUpdatedDelegation  = "updated-delegation"
	AttributeSlashedAmount      = "slashed-amount"
	AttributeKeyAck             = "acknowledgement"
	AttributeKeyAckSuccess      = "success"
	AttributeKeyAckError        = "error"
	AttributeValueCategory      = ModuleName
)
