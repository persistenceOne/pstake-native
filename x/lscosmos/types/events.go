package types

// IBC events
const (
	EventTypePacket            = "ics27_packet"
	EventTypeTimeout           = "timeout"
	EventTypeLiquidStake       = "liquid-stake"
	EventTypeRedeem            = "redeem"
	EventTypeLiquidUnstake     = "liquid-unstake"
	EventTypeClaim             = "claim"
	EventTypeJumpStart         = "jump-start"
	EventTypeRecreateICA       = "recreate-ica"
	EventTypeChangeModuleState = "change-module-state"
	EventTypeReportSlashing    = "report-slashing"
	// this line is used by starport scaffolding # ibc/packet/event

	AttributeKeyAckSuccess         = "success"
	AttributeKeyAck                = "acknowledgement"
	AttributeKeyAckError           = "error"
	AttributeAmount                = "amount"
	AttributeAmountReceived        = "received"
	AttributeUnstakeAmount         = "undelegation-amount"
	AttributePstakeDepositFee      = "pstake-deposit-fee"
	AttributePstakeRedeemFee       = "pstake-redeem-fee"
	AttributePstakeUnstakeFee      = "pstake-unstake-fee"
	AttributeDelegatorAddress      = "address"
	AttributeRewarderAddress       = "rewarder-address"
	AttributeClaimedAmount         = "claimed-amount"
	AttributePstakeAddress         = "pstake-address"
	AttributeFromAddress           = "from-address"
	AttributeRecreateDelegationICA = "recreate-delegation-ica"
	AttributeRecreateRewardsICA    = "recreate-rewards-ica"
	AttributeChangedModuleState    = "module-state"
	AtttibuteValidatorAddress      = "validator-address"
	AttributeValueCategory         = ModuleName
)
