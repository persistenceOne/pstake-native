package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
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

// QueryHostChainValidator sends an ICQ query to retrieve a specific host chain validator
func (k *Keeper) QueryHostChainValidator(
	ctx sdk.Context,
	hc *types.HostChain,
	validatorAddress string,
) error {

	_, byteAddress, err := bech32.DecodeAndConvert(validatorAddress)
	if err != nil {
		return err
	}

	k.icqKeeper.MakeRequest(
		ctx,
		hc.ConnectionId,
		hc.ChainId,
		types.StakingStoreQuery,
		stakingtypes.GetValidatorKey(byteAddress),
		sdk.NewInt(int64(-1)),
		types.ModuleName,
		Validator,
		0,
	)

	return nil
}

// QueryDelegationHostChainAccountBalance sends an ICQ query to get the delegation host account balance
func (k *Keeper) QueryDelegationHostChainAccountBalance(
	ctx sdk.Context,
	hc *types.HostChain,
) error {
	_, byteAddress, err := bech32.DecodeAndConvert(hc.DelegationAccount.Address)
	if err != nil {
		return err
	}

	key := banktypes.CreatePrefixedAccountStoreKey(byteAddress, []byte(hc.HostDenom))
	fmt.Println(key)

	k.icqKeeper.MakeRequest(
		ctx,
		hc.ConnectionId,
		hc.ChainId,
		types.BankStoreQuery,
		key,
		sdk.NewInt(int64(-1)),
		types.ModuleName,
		DelegationAccountBalances,
		0,
	)

	return nil
}

// QueryRewardsHostChainAccountBalance sends an ICQ query to get the rewards host account balance
func (k *Keeper) QueryRewardsHostChainAccountBalance(
	ctx sdk.Context,
	hc *types.HostChain,
) error {
	_, byteAddress, err := bech32.DecodeAndConvert(hc.RewardsAccount.Address)
	if err != nil {
		return err
	}

	key := banktypes.CreatePrefixedAccountStoreKey(byteAddress, []byte(hc.HostDenom))
	fmt.Println(key)

	k.icqKeeper.MakeRequest(
		ctx,
		hc.ConnectionId,
		hc.ChainId,
		types.BankStoreQuery,
		key,
		sdk.NewInt(int64(-1)),
		types.ModuleName,
		RewardAccountBalances,
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
	_, delegatorAddr, err := bech32.DecodeAndConvert(hc.DelegationAccount.Address)
	if err != nil {
		return err
	}

	_, validatorAddr, err := bech32.DecodeAndConvert(validator.OperatorAddress)
	if err != nil {
		return err
	}

	k.icqKeeper.MakeRequest(
		ctx,
		hc.ConnectionId,
		hc.ChainId,
		types.StakingStoreQuery,
		stakingtypes.GetDelegationKey(delegatorAddr, validatorAddr),
		sdk.NewInt(int64(-1)),
		types.ModuleName,
		Delegation,
		0,
	)

	return nil
}
