package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	epochtypes "github.com/persistenceOne/persistence-sdk/v2/x/epochs/types"
	liquidstakeibctypes "github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
	"github.com/persistenceOne/pstake-native/v2/x/ratesync/types"
)

// Wrapper struct
type LiquidStakeIBCHooks struct {
	k Keeper
}

var _ liquidstakeibctypes.LiquidStakeIBCHooks = LiquidStakeIBCHooks{}

// Create new lsibc hooks
func (k Keeper) LiquidStakeIBCHooks() LiquidStakeIBCHooks {
	return LiquidStakeIBCHooks{k}
}

func (h LiquidStakeIBCHooks) PostCValueUpdate(ctx sdk.Context, mintDenom, hostDenom string, cValue sdk.Dec) error {
	h.k.Logger(ctx).Info("called ratesync hook for PostCValueUpdate")
	return h.k.PostCValueUpdate(ctx, mintDenom, hostDenom, cValue)
}
func (k Keeper) PostCValueUpdate(ctx sdk.Context, mintDenom, hostDenom string, cValue sdk.Dec) error {
	hcs := k.GetAllHostChain(ctx)
	for _, hc := range hcs {
		if hc.Features.LiquidStakeIBC.Enabled {
			err := k.ExecuteLiquidStakeRateTx(ctx, hc.Features.LiquidStakeIBC, mintDenom, hostDenom, cValue, hc.Id, hc.ConnectionId, hc.IcaAccount)
			if err != nil {
				k.Logger(ctx).Error("cannot ExecuteLiquidStakeRateTx for host chain ",
					"id", hc.Id,
					"mint-denom", mintDenom,
					"err:", err)
			}
		}
	}
	return nil
}

// Wrapper struct
type EpochHooks struct {
	k Keeper
}

var _ epochtypes.EpochHooks = EpochHooks{}

// Create new epoch hooks
func (k Keeper) EpochHooks() EpochHooks {
	return EpochHooks{k}
}

func (e EpochHooks) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	e.k.Logger(ctx).Info("called ratesync hook for AfterEpochEnd")
	//return e.k.AfterEpochEnd(ctx, epochIdentifier, epochNumber) //TODO uncomment after wiring liquidstakekeeper to keeper and app
	return nil
}
func (k Keeper) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	nas := k.liquidStakeKeeper.GetNetAmountState(ctx)
	liquidBondDenom := k.liquidStakeKeeper.LiquidBondDenom(ctx)
	bondDenom, found := liquidstakeibctypes.MintDenomToHostDenom(liquidBondDenom)
	if !found {
		return errorsmod.Wrapf(sdkerrors.ErrNotFound, "bondDenom could not be derived from host denom")
	}
	hcs := k.GetAllHostChain(ctx)
	for _, hc := range hcs {
		if hc.Features.LiquidStake.Enabled && epochIdentifier == types.LiquidStakeEpoch {
			// Add liquidstakekeeper and do stuff
			err := k.ExecuteLiquidStakeRateTx(ctx, hc.Features.LiquidStakeIBC, liquidBondDenom, bondDenom, nas.MintRate, hc.Id, hc.ConnectionId, hc.IcaAccount)
			if err != nil {
				k.Logger(ctx).Error("cannot ExecuteLiquidStakeRateTx for host chain ",
					"id", hc.Id,
					"mint-denom", liquidBondDenom,
					"err:", err)
			}
		}

	}
	return nil
}
func (e EpochHooks) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	e.k.Logger(ctx).Info("called ratesync hook for BeforeEpochStart")
	return nil
}
