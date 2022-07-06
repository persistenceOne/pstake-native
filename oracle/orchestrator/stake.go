package orchestrator

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/relayer/relayer"
)

func MsgDelegate(cosmos *relayer.Chain, amount sdk.Coin, delegatorAddr string, validatorAddr string) sdk.Msg {
	msg := stakingtypes.MsgDelegate{
		DelegatorAddress: delegatorAddr,
		ValidatorAddress: validatorAddr,
		Amount:           amount,
	}
	return &msg
}

func MsgUndelegate(cosmos *relayer.Chain, amount sdk.Coin, delegatorAddr string, validatorAddr string) sdk.Msg {
	msg := stakingtypes.MsgUndelegate{
		DelegatorAddress: delegatorAddr,
		ValidatorAddress: validatorAddr,
		Amount:           amount,
	}

	return &msg
}
