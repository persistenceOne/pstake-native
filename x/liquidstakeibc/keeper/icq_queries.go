package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

// QueryHostChainValidators sends an ICQ query to retrieve the host chain validator set
func (k *Keeper) QueryHostChainValidators(
	ctx sdk.Context,
	hc *types.HostChain,
	req stakingtypes.QueryValidatorsRequest,
) error {
	bz, err := k.cdc.Marshal(&req)
	if err != nil {
		return err
	}

	k.icqKeeper.MakeRequest(
		ctx,
		hc.ConnectionId,
		hc.ChainId,
		"cosmos.staking.v1beta1.Query/Validators",
		bz,
		sdk.NewInt(int64(-1)),
		types.ModuleName,
		ValidatorSet,
		0,
	)

	return nil
}

// QueryHostChainAccountBalance sends an ICQ query to get a host account balance
func (k *Keeper) QueryHostChainAccountBalance(
	ctx sdk.Context,
	hc *types.HostChain,
	address string,
) error {
	balanceQuery := banktypes.QueryBalanceRequest{
		Address: address,
		Denom:   hc.HostDenom,
	}
	bz, err := k.cdc.Marshal(&balanceQuery)
	if err != nil {
		return err
	}

	k.icqKeeper.MakeRequest(
		ctx,
		hc.ConnectionId,
		hc.ChainId,
		"cosmos.bank.v1beta1.Query/Balance",
		bz,
		sdk.NewInt(int64(-1)),
		types.ModuleName,
		Balances,
		0,
	)

	return nil
}

// QueryValidatorDelegation sends an ICQ query to get a validator delegation
func (k *Keeper) QueryValidatorDelegation(
	ctx sdk.Context,
	hc *types.HostChain,
	validator *types.Validator,
) error {
	delegationRequest := stakingtypes.QueryDelegationRequest{
		DelegatorAddr: hc.DelegationAccount.Address,
		ValidatorAddr: validator.OperatorAddress,
	}
	bz, err := k.cdc.Marshal(&delegationRequest)
	if err != nil {
		return err
	}

	k.icqKeeper.MakeRequest(
		ctx,
		hc.ConnectionId,
		hc.ChainId,
		"cosmos.staking.v1beta1.Query/Delegation",
		bz,
		sdk.NewInt(int64(-1)),
		types.ModuleName,
		Delegation,
		0,
	)

	return nil
}
