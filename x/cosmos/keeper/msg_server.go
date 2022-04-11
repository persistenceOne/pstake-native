package keeper

import (
	"context"
	"strconv"

	"github.com/armon/go-metrics"
	sdkTelemetry "github.com/cosmos/cosmos-sdk/telemetry"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the gov MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) cosmosTypes.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (k msgServer) SetOrchestrator(c context.Context, msg *cosmosTypes.MsgSetOrchestrator) (*cosmosTypes.MsgSetOrchestratorResponse, error) {
	err := msg.ValidateBasic()
	if err != nil {
		return nil, sdkErrors.Wrap(err, "Key not valid")
	}
	ctx := sdkTypes.UnwrapSDKContext(c)

	//Accept transaction if module is enabled
	if !k.GetParams(ctx).ModuleEnabled {
		return nil, cosmosTypes.ErrModuleNotEnabled
	}

	validator, e1 := sdkTypes.ValAddressFromBech32(msg.Validator)
	orchestrator, e2 := sdkTypes.AccAddressFromBech32(msg.Orchestrator)
	if e1 != nil || e2 != nil {
		return nil, sdkErrors.Wrap(err, "Key not valid")
	}
	_, foundExistingOrchestratorKey := k.GetOrchestratorValidator(ctx, orchestrator)

	if k.Keeper.stakingKeeper.Validator(ctx, validator) == nil {
		return nil, sdkErrors.Wrap(stakingTypes.ErrNoValidatorFound, validator.String())
	} else if foundExistingOrchestratorKey {
		return nil, sdkErrors.Wrap(cosmosTypes.ErrResetDelegateKeys, validator.String())
	}

	//TODO reverse key value (Important for unique orch addresses for each validator)
	delegateKeys := k.GetDelegateKeys(ctx)
	for i := range delegateKeys {
		if delegateKeys[i].Orchestrator == orchestrator.String() {
			return nil, sdkErrors.Wrap(err, "Duplicate Orchestrator Key")
		}
	}
	// set the orchestrator address
	k.SetOrchestratorValidator(ctx, validator, orchestrator)

	ctx.EventManager().EmitEvent(
		sdkTypes.NewEvent(
			sdkTypes.EventTypeMessage,
			sdkTypes.NewAttribute(sdkTypes.AttributeKeyModule, cosmosTypes.AttributeValueCategory),
			sdkTypes.NewAttribute(cosmosTypes.AttributeKeySetOperatorAddr, orchestrator.String()),
		),
	)

	return &cosmosTypes.MsgSetOrchestratorResponse{}, nil
}

// Send TODO Modify outgoing pool
func (k msgServer) Withdraw(c context.Context, msg *cosmosTypes.MsgWithdrawStkAsset) (*cosmosTypes.MsgWithdrawStkAssetResponse, error) {
	ctx := sdkTypes.UnwrapSDKContext(c)

	//Accept transaction if module is enabled
	if !k.GetParams(ctx).ModuleEnabled {
		return nil, cosmosTypes.ErrModuleNotEnabled
	}

	from, err := sdkTypes.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return nil, err
	}
	to, err := sdkTypes.AccAddressFromBech32(msg.ToAddress)
	if err != nil {
		return nil, err
	}

	if ctx.IsZero() || sdkTypes.VerifyAddressFormat(from) != nil || sdkTypes.VerifyAddressFormat(to) != nil ||
		!msg.Amount.IsValid() {
		return nil, sdkErrors.Wrap(cosmosTypes.ErrInvalid, "arguments")
	}

	if msg.Amount.GetDenom() != k.GetParams(ctx).MintDenom {
		return nil, cosmosTypes.ErrInvalidWithdrawDenom
	}

	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, from, cosmosTypes.ModuleName, sdkTypes.NewCoins(msg.Amount)); err != nil {
		return nil, err
	}

	err = k.addToWithdrawPool(ctx, *msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdkTypes.NewEvent(
			sdkTypes.EventTypeMessage,
			sdkTypes.NewAttribute(sdkTypes.AttributeKeyModule, cosmosTypes.AttributeValueCategory),
		),
	)
	return &cosmosTypes.MsgWithdrawStkAssetResponse{}, nil
}

func (k msgServer) MintTokensForAccount(c context.Context, msg *cosmosTypes.MsgMintTokensForAccount) (*cosmosTypes.MsgMintTokensForAccountResponse, error) {
	ctx := sdkTypes.UnwrapSDKContext(c)

	//Accept transaction if module is enabled
	if !k.GetParams(ctx).ModuleEnabled {
		return nil, cosmosTypes.ErrModuleNotEnabled
	}

	destinationAddress, err := sdkTypes.AccAddressFromBech32(msg.AddressFromMemo)
	if err != nil {
		return nil, err
	}

	orchestratorAddress, err := sdkTypes.AccAddressFromBech32(msg.OrchestratorAddress)
	if err != nil {
		return nil, err
	}

	if ctx.IsZero() || sdkTypes.VerifyAddressFormat(destinationAddress) != nil || sdkTypes.VerifyAddressFormat(orchestratorAddress) != nil ||
		!msg.Amount.IsValid() {
		return nil, sdkErrors.Wrap(cosmosTypes.ErrInvalid, "arguments")
	}

	params := k.GetParams(ctx)
	uatomAmount := msg.Amount.AmountOf(params.BondDenom)
	uStkXprtCoin := sdkTypes.NewCoin(params.MintDenom, uatomAmount)
	newAmount := sdkTypes.NewCoins(uStkXprtCoin)
	if err != nil {
		return nil, err
	}

	k.setMintAddressAndAmount(ctx, msg.ChainID, msg.BlockHeight, msg.TxHash, destinationAddress, newAmount)

	_, found := k.GetOrchestratorValidator(ctx, orchestratorAddress)
	if found {
		err = k.addToMintingPoolTx(ctx, msg.TxHash, destinationAddress, orchestratorAddress, msg.Amount)
		if err != nil {
			return nil, err
		}
	}

	ctx.EventManager().EmitEvent(
		sdkTypes.NewEvent(
			sdkTypes.EventTypeMessage,
			sdkTypes.NewAttribute(sdkTypes.AttributeKeyModule, cosmosTypes.AttributeValueCategory),
			sdkTypes.NewAttribute(cosmosTypes.AttributeSender, orchestratorAddress.String()),
		),
	)
	return &cosmosTypes.MsgMintTokensForAccountResponse{}, nil
}

func (k msgServer) MakeProposal(c context.Context, msg *cosmosTypes.MsgMakeProposal) (*cosmosTypes.MsgMakeProposalResponse, error) {
	ctx := sdkTypes.UnwrapSDKContext(c)

	//Accept transaction if module is enabled
	if !k.GetParams(ctx).ModuleEnabled {
		return nil, cosmosTypes.ErrModuleNotEnabled
	}

	orchestratorAddress, err := sdkTypes.AccAddressFromBech32(msg.OrchestratorAddress)
	if err != nil {
		return nil, err
	}

	_, found := k.GetOrchestratorValidator(ctx, orchestratorAddress)
	if found {
		k.setProposalDetails(ctx, msg.ChainID, msg.BlockHeight, msg.ProposalID, msg.Title, msg.Description, orchestratorAddress, msg.VotingStartTime, msg.VotingEndTime)
	}

	ctx.EventManager().EmitEvent(
		sdkTypes.NewEvent(
			sdkTypes.EventTypeMessage,
			sdkTypes.NewAttribute(sdkTypes.AttributeKeyModule, cosmosTypes.AttributeValueCategory),
			sdkTypes.NewAttribute(cosmosTypes.AttributeSender, orchestratorAddress.String()),
		),
	)
	return &cosmosTypes.MsgMakeProposalResponse{}, nil
}

func (k msgServer) Vote(c context.Context, msg *cosmosTypes.MsgVote) (*cosmosTypes.MsgVoteResponse, error) {
	ctx := sdkTypes.UnwrapSDKContext(c)

	//Accept transaction if module is enabled
	if !k.GetParams(ctx).ModuleEnabled {
		return nil, cosmosTypes.ErrModuleNotEnabled
	}

	accAddr, accErr := sdkTypes.AccAddressFromBech32(msg.Voter)
	if accErr != nil {
		return nil, accErr
	}

	_, found := k.GetOrchestratorValidator(ctx, accAddr)
	if found {
		err := k.Keeper.AddVote(ctx, msg.ProposalId, accAddr, cosmosTypes.NewNonSplitVoteOption(msg.Option))
		if err != nil {
			return nil, err
		}
	} else {
		return nil, cosmosTypes.ErrInvalidVote
	}

	defer sdkTelemetry.IncrCounterWithLabels(
		[]string{cosmosTypes.ModuleName, "vote"},
		1,
		[]metrics.Label{
			sdkTelemetry.NewLabel("proposal_id", strconv.Itoa(int(msg.ProposalId))),
		},
	)

	ctx.EventManager().EmitEvent(
		sdkTypes.NewEvent(
			sdkTypes.EventTypeMessage,
			sdkTypes.NewAttribute(sdkTypes.AttributeKeyModule, cosmosTypes.AttributeValueCategory),
			sdkTypes.NewAttribute(sdkTypes.AttributeKeySender, msg.Voter),
		),
	)
	return &cosmosTypes.MsgVoteResponse{}, nil
}

func (k msgServer) VoteWeighted(c context.Context, msg *cosmosTypes.MsgVoteWeighted) (*cosmosTypes.MsgVoteWeightedResponse, error) {
	ctx := sdkTypes.UnwrapSDKContext(c)

	//Accept transaction if module is enabled
	if !k.GetParams(ctx).ModuleEnabled {
		return nil, cosmosTypes.ErrModuleNotEnabled
	}

	accAddr, accErr := sdkTypes.AccAddressFromBech32(msg.Voter)
	if accErr != nil {
		return nil, accErr
	}

	if _, found := k.GetOrchestratorValidator(ctx, accAddr); found {
		err := k.Keeper.AddVote(ctx, msg.ProposalId, accAddr, msg.Options)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, cosmosTypes.ErrInvalidVote
	}

	defer sdkTelemetry.IncrCounterWithLabels(
		[]string{cosmosTypes.ModuleName, "vote"},
		1,
		[]metrics.Label{
			sdkTelemetry.NewLabel("proposal_id", strconv.Itoa(int(msg.ProposalId))),
		},
	)

	ctx.EventManager().EmitEvent(
		sdkTypes.NewEvent(
			sdkTypes.EventTypeMessage,
			sdkTypes.NewAttribute(sdkTypes.AttributeKeyModule, cosmosTypes.AttributeValueCategory),
			sdkTypes.NewAttribute(sdkTypes.AttributeKeySender, msg.Voter),
		),
	)

	return &cosmosTypes.MsgVoteWeightedResponse{}, nil
}

// SignedTxFromOrchestrator Receives a signed txn from orchestrator and updates the details
func (k msgServer) SignedTxFromOrchestrator(c context.Context, msg *cosmosTypes.MsgSignedTx) (*cosmosTypes.MsgSignedTxResponse, error) {
	ctx := sdkTypes.UnwrapSDKContext(c)

	//Accept transaction if module is enabled
	if !k.GetParams(ctx).ModuleEnabled {
		return nil, cosmosTypes.ErrModuleNotEnabled
	}

	orchAddr, orchErr := sdkTypes.AccAddressFromBech32(msg.OrchestratorAddress)
	if orchErr != nil {
		return nil, orchErr
	}

	txBytes, err := msg.Tx.Marshal()
	if err != nil {
		return nil, err
	}
	txHash := cosmosTypes.BytesToHexUpper(txBytes)

	txn, err := k.getTxnFromOutgoingPoolByID(ctx, msg.TxID)
	if err != nil {
		return nil, err
	}
	if txn.CosmosTxDetails.TxHash == "" {
		err = k.setTxDetailsSignedByOrchestrator(ctx, msg.TxID, txHash, msg.Tx)
		if err != nil {
			return nil, err
		}

		k.setTxHashAndDetails(ctx, orchAddr, msg.TxID, txHash, "pending")
	} else {
		return nil, cosmosTypes.ErrTxnDetailsAlreadySent
	}

	ctx.EventManager().EmitEvent(
		sdkTypes.NewEvent(
			sdkTypes.EventTypeMessage,
			sdkTypes.NewAttribute(sdkTypes.AttributeKeyModule, msg.Type()),
			sdkTypes.NewAttribute(sdkTypes.AttributeKeySender, msg.OrchestratorAddress),
		),
	)
	return &cosmosTypes.MsgSignedTxResponse{}, nil
}

// TxStatus Accepts status as : "success" or "failure"
// Failure only to be sent when transaction fails due to insufficient fees
func (k msgServer) TxStatus(c context.Context, msg *cosmosTypes.MsgTxStatus) (*cosmosTypes.MsgTxStatusResponse, error) {
	ctx := sdkTypes.UnwrapSDKContext(c)

	//Accept transaction if module is enabled
	if !k.GetParams(ctx).ModuleEnabled {
		return nil, cosmosTypes.ErrModuleNotEnabled
	}

	orchAddr, orchErr := sdkTypes.AccAddressFromBech32(msg.OrchestratorAddress)
	if orchErr != nil {
		return nil, orchErr
	}
	_, found := k.GetOrchestratorValidator(ctx, orchAddr)
	if found {
		if msg.Status == "success" || msg.Status == "failure" {
			k.setTxHashAndDetails(ctx, orchAddr, 0, msg.TxHash, msg.Status)
		} else {
			return nil, cosmosTypes.ErrInvalidStatus
		}
	}

	ctx.EventManager().EmitEvent(
		sdkTypes.NewEvent(
			sdkTypes.EventTypeMessage,
			sdkTypes.NewAttribute(sdkTypes.AttributeKeyModule, cosmosTypes.AttributeValueCategory),
			sdkTypes.NewAttribute(cosmosTypes.AttributeSender, orchAddr.String()),
		),
	)
	return &cosmosTypes.MsgTxStatusResponse{}, nil
}
