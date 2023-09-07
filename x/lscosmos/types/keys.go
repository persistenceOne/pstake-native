package types

const (
	// ModuleName defines the module name
	ModuleName = "lscosmos"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// MsgTypeLiquidStake is the type of message to liquid stake
	MsgTypeLiquidStake = "msg_liquid_stake"

	// MsgTypeLiquidUnstake is the type of message liquid unstake
	MsgTypeLiquidUnstake = "msg_liquid_unstake"

	// MsgTypeRedeem is the type of message redeem
	MsgTypeRedeem = "msg_redeem"

	// MsgTypeClaim is the type of message claim
	MsgTypeClaim = "msg_claim"

	// MsgTypeRecreateICA is the type of message RecreateICA
	MsgTypeRecreateICA = "msg_recreate_ica"

	// MsgTypeJumpStart is the type of message Jump start
	MsgTypeJumpStart = "msg_jump_start"

	// MsgTypeChangeModuleState is the type of message Change Module State
	MsgTypeChangeModuleState = "msg_change_module_state"

	// MsgTypeReportSlashing is the type of message Report Slashing
	MsgTypeReportSlashing = "msg_report_slashing"
)
