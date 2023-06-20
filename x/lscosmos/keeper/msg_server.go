package keeper

import (
	"context"

	"github.com/persistenceOne/pstake-native/v2/x/lscosmos/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer { //nolint:staticcheck
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{} //nolint:staticcheck

// LiquidStake defines a method for liquid staking tokens
func (m msgServer) LiquidStake(goCtx context.Context, msg *types.MsgLiquidStake) (*types.MsgLiquidStakeResponse, error) {
	return nil, types.ErrDeprecated
}

// LiquidUnstake defines a method for unstaking the liquid staked tokens
func (m msgServer) LiquidUnstake(goCtx context.Context, msg *types.MsgLiquidUnstake) (*types.MsgLiquidUnstakeResponse, error) {
	return nil, types.ErrDeprecated
}

// Redeem defines a method for redeeming liquid staked tokens instantly
func (m msgServer) Redeem(goCtx context.Context, msg *types.MsgRedeem) (*types.MsgRedeemResponse, error) {
	return nil, types.ErrDeprecated
}

// Claim defines a method for claiming unstaked mature tokens or failed unbondings
func (m msgServer) Claim(goCtx context.Context, msg *types.MsgClaim) (*types.MsgClaimResponse, error) {
	return nil, types.ErrDeprecated
}

// JumpStart defines a method for jump-starting the module through fee address account.
func (m msgServer) JumpStart(goCtx context.Context, msg *types.MsgJumpStart) (*types.MsgJumpStartResponse, error) {
	return nil, types.ErrDeprecated
}

// RecreateICA defines a method for recreating closed ica channels
func (m msgServer) RecreateICA(goCtx context.Context, msg *types.MsgRecreateICA) (*types.MsgRecreateICAResponse, error) {
	return nil, types.ErrDeprecated
}

// ChangeModuleState defines an admin method for disabling or re-enabling module state
func (m msgServer) ChangeModuleState(goCtx context.Context, msg *types.MsgChangeModuleState) (*types.MsgChangeModuleStateResponse, error) {
	return nil, types.ErrDeprecated
}

// ReportSlashing defines an admin method for reporting slashing on a validator
func (m msgServer) ReportSlashing(goCtx context.Context, msg *types.MsgReportSlashing) (*types.MsgReportSlashingResponse, error) {
	return nil, types.ErrDeprecated
}
