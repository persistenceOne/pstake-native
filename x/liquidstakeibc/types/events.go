package types

const (
	EventTypeLiquidStake    = "liquid-stake"
	EventTypeLiquidStakeLSM = "liquid-stake-lsm"
	EventTypeLiquidUnstake  = "liquid-unstake"
	EventTypeRedeem         = "redeem"
	EventTypePacket         = "ics27_packet"
	EventTypeTimeout        = "timeout"
	EventTypeSlashing       = "slashing"
	EventTypeUpdateParams   = "update_params"
	EventTypeChainDisabled  = "chain_disabled"

	AttributeAmount             = "amount"
	AttributeAmountReceived     = "received"
	AttributeDelegatorAddress   = "address"
	AttributePstakeDepositFee   = "pstake-deposit-fee"
	AttributePstakeUnstakeFee   = "pstake-unstake-fee"
	AttributePstakeRedeemFee    = "pstake-redeem-fee"
	AttributeChainID            = "chain-id"
	AttributeCValue             = "c-value"
	AttributeUnstakeAmount      = "undelegation-amount"
	AttributeUnstakeEpoch       = "undelegation-epoch"
	AttributeValidatorAddress   = "validator-address"
	AttributeExistingDelegation = "existing-delegation"
	AttributeUpdatedDelegation  = "updated-delegation"
	AttributeSlashedAmount      = "slashed-amount"
	AttributeKeyAuthority       = "authority"
	AttributeKeyUpdatedParams   = "updated_params"
	AttributeKeyAck             = "acknowledgement"
	AttributeKeyAckSuccess      = "success"
	AttributeKeyAckError        = "error"
	AttributeValueCategory      = ModuleName
)
