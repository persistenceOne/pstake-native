package types

const (
	EventTypeLiquidStake    = "liquid-stake"
	EventTypeLiquidStakeLSM = "liquid-stake-lsm"
	EventTypeLiquidUnstake  = "liquid-unstake"
	EventTypeRedeem         = "redeem"
	EventTypePacket         = "ics27-packet"
	EventTypeTimeout        = "timeout"
	EventTypeSlashing       = "slashing"
	EventTypeUpdateParams   = "update-params"
	EventTypeChainDisabled  = "chain-disabled"

	AttributeInputAmount        = "input-amount"
	AttributeOutputAmount       = "output-amount"
	AttributeDelegatorAddress   = "address"
	AttributePstakeDepositFee   = "pstake-deposit-fee"
	AttributePstakeUnstakeFee   = "pstake-unstake-fee"
	AttributePstakeRedeemFee    = "pstake-redeem-fee"
	AttributeChainID            = "chain-id"
	AttributeNewCValue          = "new-c-value"
	AttributeOldCValue          = "old-c-value"
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
