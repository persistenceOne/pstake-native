package keeper

import (
	"context"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	ibcTransferTypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	"github.com/persistenceOne/pstake-native/x/lscosmos/types"
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
		return nil, types.ErrInvalidMessage
	}

	ctx := sdkTypes.UnwrapSDKContext(goCtx)

	// check if ibc-denom is whitelisted
	ibcParams := m.GetCosmosIBCParams(ctx)

	expectedDenom := ibcTransferTypes.GetPrefixedDenom(ibcParams.TokenTransferPort, ibcParams.TokenTransferChannel, ibcParams.BaseDenom)
	givenDenom := msg.Amount.Denom

	if givenDenom != expectedDenom {
		return nil, types.ErrInvalidDenom
	}

	// check if address in message is correct or not
	mintAddress, err := sdkTypes.AccAddressFromBech32(msg.MintAddress)
	if err != nil {
		return nil, sdkErrors.ErrInvalidAddress
	}

	// sanity check for the arguments of message

	if ctx.IsZero() || !msg.Amount.IsValid() {
		return nil, types.ErrInvalidArgs
	}

	//send the deposit to the deposit-module account
	depositAmount := sdkTypes.NewCoins(msg.Amount)
	err = m.SendTokensToDepositModule(ctx, depositAmount, mintAddress)
	if err != nil {
		return nil, types.ErrFailedDeposit
	}

	// amount of stk tokens to be minted
	mintAmountDec := msg.Amount.Amount.ToDec().Mul(m.GetCValue(ctx))

	mintToken, residue := sdkTypes.NewDecCoinFromDec(ibcParams.MintDenom, mintAmountDec).TruncateDecimal()
	if residue.Amount.GT(sdkTypes.NewDec(0)) {
		m.SendResidueToCommunityPool(ctx, sdkTypes.NewDecCoins(residue))
	}

	err = m.MintTokens(ctx, mintToken, mintAddress)
	if err != nil {
		return nil, types.ErrMintFailed
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
