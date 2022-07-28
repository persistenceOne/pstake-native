package keeper

import (
	"context"
	"fmt"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	ibcTransferTypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	"github.com/persistenceOne/pstake-native/x/ls-cosmos/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (m msgServer) LiquidStake(goCtx context.Context, msg *types.MsgLiquidStake) (*types.MsgLiquidStakeResponse, error) {
	err := msg.ValidateBasic()
	if err != nil {
		return nil, sdkErrors.Wrap(err, "invalid message")
	}

	ctx := sdkTypes.UnwrapSDKContext(goCtx)

	// check if ibc-denom is whitelisted
	ibcParams := m.GetCosmosIBCParams(ctx)

	expectedDenom := ibcTransferTypes.GetPrefixedDenom(ibcParams.TokenTransferPort, ibcParams.TokenTransferChannel, ibcParams.BaseDenom)

	givenDenom := msg.Amount.Denom

	if givenDenom != expectedDenom {
		return nil, sdkErrors.Wrap(err, "denom not whitelisted/ invalid denom")
	}

	// check if address in message is correct or not
	mintAddress, err := sdkTypes.AccAddressFromBech32(msg.MintAddress)
	if err != nil {
		return nil, err
	}

	// sanity check for the arguments of message

	if ctx.IsZero() || !msg.Amount.IsValid() {
		return nil, sdkErrors.Wrap(fmt.Errorf("invalid"), " arguments")
	}
	// amount of stk tokens to be minted
	mintAmountDec := msg.Amount.Amount.ToDec().Mul(m.GetCValue(ctx))

	mintToken, _ := sdkTypes.NewDecCoinFromDec(ibcParams.MintDenom, mintAmountDec).TruncateDecimal()

	err = m.mintTokens(ctx, mintToken, mintAddress)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdkTypes.NewEvent(
			types.EventTypeMint,
			sdkTypes.NewAttribute(types.AttributeMintedAddress, mintAddress.String()),
			sdkTypes.NewAttribute(types.AttributeAmountMinted, mintAmountDec.String()),
		),
	)
	return &types.MsgLiquidStakeResponse{}, nil
}
