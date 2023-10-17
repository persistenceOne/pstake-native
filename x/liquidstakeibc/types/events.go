package types

const (
	EventTypeLiquidStake                   = "liquid_stake"
	EventTypeLiquidStakeLSM                = "liquid_stake_lsm"
	EventTypeLiquidUnstake                 = "liquid_unstake"
	EventTypeRedeem                        = "redeem"
	EventTypePacket                        = "ics27_packet"
	EventTypeTimeout                       = "timeout"
	EventTypeSlashing                      = "slashing"
	EventTypeUpdateParams                  = "update_params"
	EventTypeChainDisabled                 = "chain_disabled"
	EventTypeValidatorStatusUpdate         = "validator_status_update"
	EventTypeValidatorExchangeRateUpdate   = "validator_exchange_rate_update"
	EventTypeValidatorDelegableStateUpdate = "validator_delegable_state_update"
	EventTypeDoDelegation                  = "send_delegation"

	AttributeInputAmount                 = "input_amount"
	AttributeOutputAmount                = "output_amount"
	AttributeDelegatorAddress            = "address"
	AttributePstakeDepositFee            = "pstake_deposit_fee"
	AttributePstakeUnstakeFee            = "pstake_unstake_fee"
	AttributePstakeRedeemFee             = "pstake_redeem_fee"
	AttributeChainID                     = "chain_id"
	AttributeNewCValue                   = "new_c_value"
	AttributeOldCValue                   = "old_c_value"
	AttributeEpoch                       = "epoch_number"
	AttributeValidatorAddress            = "validator_address"
	AttributeExistingDelegation          = "existing_delegation"
	AttributeUpdatedDelegation           = "updated_delegation"
	AttributeSlashedAmount               = "slashed_amount"
	AttributeKeyAuthority                = "authority"
	AttributeKeyUpdatedParams            = "updated_params"
	AttributeKeyAck                      = "acknowledgement"
	AttributeKeyAckSuccess               = "success"
	AttributeKeyAckError                 = "error"
	AttributeKeyValidatorNewStatus       = "validator_new_status"
	AttributeKeyValidatorOldStatus       = "validator_old_status"
	AttributeKeyValidatorNewExchangeRate = "validator_new_exchange_rate"
	AttributeKeyValidatorOldExchangeRate = "validator_old_exchange_rate"
	AttributeKeyValidatorDelegable       = "validator_delegable"
	AttributeTotalDelegatedAmount        = "total_delegated_amount"
	AttributeIBCSequenceID               = "ibc_sequence_id"
	AttributeICAMessages                 = "ica_messages"
	AttributeValueCategory               = ModuleName
)
