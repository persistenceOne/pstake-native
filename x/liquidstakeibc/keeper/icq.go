package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	q "github.com/cosmos/cosmos-sdk/types/query"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	icqtypes "github.com/persistenceOne/persistence-sdk/v2/x/interchainquery/types"

	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

const (
	ValidatorSet = "validatorset"
)

type CallbackFn func(Keeper, sdk.Context, []byte, icqtypes.Query) error

type Callbacks struct {
	k         Keeper
	callbacks map[string]CallbackFn
}

var _ icqtypes.QueryCallbacks = Callbacks{}

func (k *Keeper) CallbackHandler() Callbacks {
	return Callbacks{*k, make(map[string]CallbackFn)}
}
func (c Callbacks) AddCallback(id string, fn interface{}) icqtypes.QueryCallbacks {
	c.callbacks[id] = fn.(CallbackFn)
	return c
}

func (c Callbacks) Call(ctx sdk.Context, id string, args []byte, query icqtypes.Query) error {
	return c.callbacks[id](c.k, ctx, args, query)
}

func (c Callbacks) Has(id string) bool {
	_, found := c.callbacks[id]
	return found
}

func (c Callbacks) RegisterCallbacks() icqtypes.QueryCallbacks {
	a := c.
		AddCallback(ValidatorSet, CallbackFn(ValidatorSetCallback))

	return a.(Callbacks)
}

// Callbacks

func ValidatorSetCallback(k Keeper, ctx sdk.Context, data []byte, query icqtypes.Query) error {
	hc, found := k.GetHostChain(ctx, query.ChainId)
	if !found {
		return fmt.Errorf("host chain with id %s is not registered", query.ChainId)
	}

	response := stakingtypes.QueryValidatorsResponse{}
	err := k.cdc.Unmarshal(data, &response)
	if err != nil {
		return fmt.Errorf("could not unmarshall ICQ validatorset response: %w", err)
	}

	// if the result is not complete, submit an ICQ query to gather the next chunk
	if response.Pagination != nil && len(response.Pagination.NextKey) > 0 {
		request := stakingtypes.QueryValidatorsRequest{}
		err = k.cdc.Unmarshal(query.Request, &request)
		if err != nil {
			return fmt.Errorf("could not unmarshall ICQ validatorset request: %w", err)
		}

		request.Pagination = new(q.PageRequest)
		request.Pagination.Key = response.Pagination.NextKey
		if err = k.QueryHostChainValidators(ctx, &hc, request); err != nil {
			return errorsmod.Wrapf(types.ErrFailedICQRequest, "error submitting validators icq: %w", err)
		}
	}

	k.SetHostChainValidators(ctx, &hc, response.Validators)

	return nil
}
