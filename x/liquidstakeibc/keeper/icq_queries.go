package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/persistenceOne/pstake-native/v3/x/liquidstakeibc/types"
)

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

// QueryValidatorDelegationUpdate sends an ICQ query to get a validator delegation
func (k *Keeper) QueryValidatorDelegationUpdate(
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
		DelegationUpdate,
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

// QueryNonCompoundableRewardsHostChainAccountBalance sends an ICQ query to get the non-compoundable rewards host account balance
func (k *Keeper) QueryNonCompoundableRewardsHostChainAccountBalance(
	ctx sdk.Context,
	hc *types.HostChain,
) error {
	_, byteAddress, err := bech32.DecodeAndConvert(hc.RewardsAccount.Address)
	if err != nil {
		return err
	}

	key := banktypes.CreatePrefixedAccountStoreKey(byteAddress, []byte(hc.RewardParams.Denom))

	k.icqKeeper.MakeRequest(
		ctx,
		hc.ConnectionId,
		hc.ChainId,
		types.BankStoreQuery,
		key,
		sdk.NewInt(int64(-1)),
		types.ModuleName,
		NonCompoundableRewardAccountBalances,
		0,
	)

	return nil
}
