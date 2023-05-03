package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibckeeper "github.com/cosmos/ibc-go/v6/modules/core/keeper"
	ibctmtypes "github.com/cosmos/ibc-go/v6/modules/light-clients/07-tendermint/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey

	accountKeeper       types.AccountKeeper
	bankKeeper          types.BankKeeper
	epochsKeeper        types.EpochsKeeper
	icaControllerKeeper types.ICAControllerKeeper
	scopedKeeper        types.ScopedKeeper
	ibcKeeper           *ibckeeper.Keeper
	icqKeeper           types.ICQKeeper

	paramSpace paramtypes.Subspace

	msgRouter *baseapp.MsgServiceRouter

	authority string
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,

	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	epochsKeeper types.EpochsKeeper,
	icaControllerKeeper types.ICAControllerKeeper,
	scopedKeeper types.ScopedKeeper,
	ibcKeeper *ibckeeper.Keeper,
	icqKeeper types.ICQKeeper,

	paramSpace paramtypes.Subspace,

	msgRouter *baseapp.MsgServiceRouter,

	authority string,
) Keeper {
	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:                 cdc,
		accountKeeper:       accountKeeper,
		bankKeeper:          bankKeeper,
		epochsKeeper:        epochsKeeper,
		icaControllerKeeper: icaControllerKeeper,
		scopedKeeper:        scopedKeeper,
		ibcKeeper:           ibcKeeper,
		icqKeeper:           icqKeeper,
		storeKey:            storeKey,
		paramSpace:          paramSpace,
		msgRouter:           msgRouter,
		authority:           authority,
	}
}

// Logger returns a module-specific logger.
func (k *Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetParams gets the total set of liquidstakeibc parameters.
func (k *Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the total set of liquidstakeibc parameters.
func (k *Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// SetHostChain sets a host chain in the store
func (k *Keeper) SetHostChain(ctx sdk.Context, hostZone *types.HostChain) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.HostChainKey)
	bytes := k.cdc.MustMarshal(hostZone)
	store.Set([]byte(hostZone.ChainId), bytes)
}

// GetHostChain returns a host chain given its id
func (k *Keeper) GetHostChain(ctx sdk.Context, chainID string) (types.HostChain, bool) {
	hc := types.HostChain{}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.HostChainKey)
	bytes := store.Get([]byte(chainID))
	if len(bytes) == 0 {
		return hc, false
	}

	k.cdc.MustUnmarshal(bytes, &hc)
	return hc, true
}

// GetAllHostChains retrieves all registered host chains
func (k *Keeper) GetAllHostChains(ctx sdk.Context) []*types.HostChain {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.HostChainKey)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	hostChains := make([]*types.HostChain, 0)
	for ; iterator.Valid(); iterator.Next() {
		hc := types.HostChain{}
		k.cdc.MustUnmarshal(iterator.Value(), &hc)
		hostChains = append(hostChains, &hc)
	}

	return hostChains
}

// GetHostChainFromIbcDenom returns a host chain given its ibc denomination on Persistence
func (k *Keeper) GetHostChainFromIbcDenom(ctx sdk.Context, ibcDenom string) (types.HostChain, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.HostChainKey)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	found := false
	hc := types.HostChain{}
	for ; iterator.Valid(); iterator.Next() {
		chain := types.HostChain{}
		k.cdc.MustUnmarshal(iterator.Value(), &chain)

		if chain.IBCDenom() == ibcDenom {
			hc = chain
			found = true
			break
		}
	}

	return hc, found
}

func (k *Keeper) GetHostChainFromDelegatorAddress(ctx sdk.Context, delegatorAddress string) (types.HostChain, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.HostChainKey)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	found := false
	hc := types.HostChain{}
	for ; iterator.Valid(); iterator.Next() {
		chain := types.HostChain{}
		k.cdc.MustUnmarshal(iterator.Value(), &chain)

		if chain.DelegationAccount != nil && chain.DelegationAccount.Address == delegatorAddress {
			hc = chain
			found = true
			break
		}
	}

	return hc, found
}

// GetDepositModuleAccount returns deposit module account interface
func (k *Keeper) GetDepositModuleAccount(ctx sdk.Context) authtypes.ModuleAccountI {
	return k.accountKeeper.GetModuleAccount(ctx, types.DepositModuleAccount)
}

// SetHostChainValidator sets a validator on the target host chain
func (k *Keeper) SetHostChainValidator(
	ctx sdk.Context,
	hc *types.HostChain,
	validator *types.Validator,
) {
	found := false
	for i, val := range hc.Validators {
		if validator.OperatorAddress == val.OperatorAddress {
			hc.Validators[i] = validator
			found = true
			break
		}
	}

	if !found {
		hc.Validators = append(hc.Validators, validator)
	}

	k.SetHostChain(ctx, hc)
}

// SetHostChainValidators sets the validators on a host chain from an ICQ
func (k *Keeper) SetHostChainValidators(
	ctx sdk.Context,
	hc *types.HostChain,
	validators []stakingtypes.Validator,
) {
	for _, validator := range validators {
		val, found := hc.GetValidator(validator.OperatorAddress)

		switch {
		case !found:
			hc.Validators = append(
				hc.Validators,
				&types.Validator{
					OperatorAddress: validator.OperatorAddress,
					Status:          validator.Status.String(),
					Weight:          sdk.ZeroDec(),
					DelegatedAmount: sdk.ZeroInt(),
				},
			)
		case validator.Status.String() != val.Status:
			val.Status = validator.Status.String()
		}
	}

	k.SetHostChain(ctx, hc)
}

// SendProtocolFee to the community pool
func (k *Keeper) SendProtocolFee(ctx sdk.Context, protocolFee sdk.Coins, moduleAccount, feeAddress string) error {
	addr, err := sdk.AccAddressFromBech32(feeAddress)
	if err != nil {
		return err
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, moduleAccount, addr, protocolFee)
	if err != nil {
		return err
	}
	return nil
}

// GetClientState retrieves the client state given a connection id
func (k *Keeper) GetClientState(ctx sdk.Context, connectionID string) (*ibctmtypes.ClientState, error) {
	conn, found := k.ibcKeeper.ConnectionKeeper.GetConnection(ctx, connectionID)
	if !found {
		return nil, fmt.Errorf("invalid connection id, \"%s\" not found", connectionID)
	}

	clientState, found := k.ibcKeeper.ClientKeeper.GetClientState(ctx, conn.ClientId)
	if !found {
		return nil, fmt.Errorf("client id \"%s\" not found for connection \"%s\"", conn.ClientId, connectionID)
	}

	client, ok := clientState.(*ibctmtypes.ClientState)
	if !ok {
		return nil, fmt.Errorf("invalid client state for connection \"%s\"", connectionID)
	}

	return client, nil
}

// GetChainID gets the id of the host chain given a connection id
func (k *Keeper) GetChainID(ctx sdk.Context, connectionID string) (string, error) {
	clientState, err := k.GetClientState(ctx, connectionID)
	if err != nil {
		return "", fmt.Errorf("client state not found for connection \"%s\": \"%s\"", connectionID, err.Error())
	}

	return clientState.ChainId, nil
}

// RegisterICAAccount registers an ICA
func (k *Keeper) RegisterICAAccount(ctx sdk.Context, connectionId, owner string) error {
	return k.icaControllerKeeper.RegisterInterchainAccount(
		ctx,
		connectionId,
		owner,
		"",
	)
}

func (k *Keeper) GetEpochNumber(ctx sdk.Context, epoch string) int64 {
	return k.epochsKeeper.GetEpochInfo(ctx, epoch).CurrentEpoch
}

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

func (k *Keeper) QueryHostChainAccountBalance(
	ctx sdk.Context,
	hc *types.HostChain,
	address string,
	denom string,
) error {
	balanceQuery := banktypes.QueryBalanceRequest{
		Address: address,
		Denom:   denom,
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

func (k *Keeper) GetHostChainCValue(ctx sdk.Context, hc *types.HostChain) sdk.Dec {
	mintedAmount := k.bankKeeper.GetSupply(ctx, hc.MintDenom()).Amount
	totalDelegations := hc.GetHostChainTotalDelegations()
	delegationAccountBalance := hc.DelegationAccount.Balance.Amount
	moduleAccountBalance := k.bankKeeper.GetBalance(
		ctx,
		authtypes.NewModuleAddress(types.DepositModuleAccount),
		hc.IBCDenom(),
	).Amount

	liquidStakedAmount := totalDelegations.
		Add(delegationAccountBalance).
		Add(moduleAccountBalance)

	if mintedAmount.IsZero() || liquidStakedAmount.IsZero() {
		return sdk.OneDec()
	}

	return sdk.NewDecFromInt(mintedAmount).Quo(sdk.NewDecFromInt(liquidStakedAmount))
}
