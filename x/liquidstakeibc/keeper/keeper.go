package keeper

import (
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibckeeper "github.com/cosmos/ibc-go/v6/modules/core/keeper"
	ibctmtypes "github.com/cosmos/ibc-go/v6/modules/light-clients/07-tendermint/types"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey

	accountKeeper       types.AccountKeeper
	bankKeeper          types.BankKeeper
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

// GetHostChainFromIbcDenom returns a host chain given its ibc denomination on Persistence
func (k *Keeper) GetHostChainFromIbcDenom(ctx sdk.Context, ibcDenom string) (types.HostChain, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.HostChainKey)
	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	found := false
	hc := types.HostChain{}
	hash := sha256.New()
	for ; iterator.Valid(); iterator.Next() {
		chain := types.HostChain{}
		k.cdc.MustUnmarshal(iterator.Value(), &chain)

		hash.Write([]byte(chain.PortId + "/" + chain.ChannelId + "/" + chain.HostDenom))
		chainIbcDenom := "ibc/" + strings.ToUpper(fmt.Sprintf("%x", hash.Sum(nil)))

		if chainIbcDenom == ibcDenom {
			hc = chain
			found = true
			break
		}
	}

	return hc, found
}

// GetDepositModuleAccount returns deposit module account interface
func (k Keeper) GetDepositModuleAccount(ctx sdk.Context) authtypes.ModuleAccountI {
	return k.accountKeeper.GetModuleAccount(ctx, types.DepositModuleAccount)
}

// SetHostChainValidators sets the validators on a host chain from an ICQ
func (k *Keeper) SetHostChainValidators(
	ctx sdk.Context,
	hs *types.HostChain,
	response *stakingtypes.QueryValidatorsResponse,
) {
	for _, validator := range response.Validators {
		val, found := hs.GetValidator(validator.OperatorAddress)

		switch {
		case !found:
			hs.Validators = append(
				hs.Validators,
				&types.Validator{
					OperatorAddress: validator.OperatorAddress,
					Status:          validator.Status.String(),
					CommissionRate:  validator.Commission.Rate,
				},
			)
		case validator.Status.String() != val.Status:
			val.Status = validator.Status.String()
		case validator.Commission.Rate != val.CommissionRate:
			val.CommissionRate = validator.Commission.Rate
		}
	}

	k.SetHostChain(ctx, hs)
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

// MintDenom generates a ls token denom based on the host token denom
func (k *Keeper) MintDenom(hostDenom string) string {
	return "stk" + "/" + hostDenom
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
