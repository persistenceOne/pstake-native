package keeper

import (
	"context"
	"fmt"
	signing2 "github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
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

	//check if the validator is present and return error if not present
	if k.Keeper.stakingKeeper.Validator(ctx, validator) == nil {
		return nil, sdkErrors.Wrap(stakingTypes.ErrNoValidatorFound, validator.String())
	}

	// set the orchestrator address
	err = k.SetValidatorOrchestrator(ctx, validator, orchestrator)
	if err != nil {
		return nil, err
	}

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
	err := msg.ValidateBasic()
	if err != nil {
		return nil, sdkErrors.Wrap(err, "Key not valid")
	}
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
	err := msg.ValidateBasic()
	if err != nil {
		return nil, sdkErrors.Wrap(err, "Key not valid")
	}
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

	uatomDenom, err := params.GetBondDenomOf("uatom")
	if err != nil {
		return nil, err
	}
	uatomAmount := msg.Amount.AmountOf(uatomDenom)
	uStkXprtCoin := sdkTypes.NewCoin(params.MintDenom, uatomAmount)

	k.setMintAddressAndAmount(ctx, msg.ChainID, msg.BlockHeight, msg.TxHash, destinationAddress, uStkXprtCoin)

	_, val, _, err := k.getAllValidartorOrchestratorMappingAndFindIfExist(ctx, orchestratorAddress)
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, fmt.Errorf("validator address not found")
	}

	_, found := k.GetValidatorOrchestrator(ctx, val)
	if !found {
		return nil, sdkErrors.Wrap(cosmosTypes.ErrOrchAddressNotFound, "No orchestrator validator mapping found")
	}
	err = k.addToMintingPoolTx(ctx, msg.TxHash, destinationAddress, orchestratorAddress, msg.Amount)
	if err != nil {
		return nil, err
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
	err := msg.ValidateBasic()
	if err != nil {
		return nil, sdkErrors.Wrap(err, "Key not valid")
	}
	ctx := sdkTypes.UnwrapSDKContext(c)

	//Accept transaction if module is enabled
	if !k.GetParams(ctx).ModuleEnabled {
		return nil, cosmosTypes.ErrModuleNotEnabled
	}

	orchestratorAddress, err := sdkTypes.AccAddressFromBech32(msg.OrchestratorAddress)
	if err != nil {
		return nil, err
	}

	_, val, _, err := k.getAllValidartorOrchestratorMappingAndFindIfExist(ctx, orchestratorAddress)
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, fmt.Errorf("validator address not found")
	}

	validatorAddress, found := k.GetValidatorOrchestrator(ctx, val)
	if validatorAddress == nil {
		return nil, fmt.Errorf("unauthorized to make proposal")
	}
	if !found {
		return nil, cosmosTypes.ErrInvalidProposal
	}

	k.setProposalDetails(ctx, msg.ChainID, msg.BlockHeight, msg.ProposalID, msg.Title, msg.Description, orchestratorAddress, msg.VotingStartTime, msg.VotingEndTime)

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
	err := msg.ValidateBasic()
	if err != nil {
		return nil, sdkErrors.Wrap(err, "Key not valid")
	}
	ctx := sdkTypes.UnwrapSDKContext(c)

	//Accept transaction if module is enabled
	if !k.GetParams(ctx).ModuleEnabled {
		return nil, cosmosTypes.ErrModuleNotEnabled
	}

	accAddr, accErr := sdkTypes.AccAddressFromBech32(msg.Voter)
	if accErr != nil {
		return nil, accErr
	}

	_, val, _, err := k.getAllValidartorOrchestratorMappingAndFindIfExist(ctx, accAddr)
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, fmt.Errorf("validator address not found")
	}
	validatorAddress, found := k.GetValidatorOrchestrator(ctx, val)
	if validatorAddress == nil {
		return nil, fmt.Errorf("unauthorized to vote")
	}
	if !found {
		return nil, cosmosTypes.ErrInvalidVote
	}

	err = k.Keeper.AddVote(ctx, msg.ProposalId, accAddr, cosmosTypes.NewNonSplitVoteOption(msg.Option))
	if err != nil {
		return nil, err
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
	err := msg.ValidateBasic()
	if err != nil {
		return nil, sdkErrors.Wrap(err, "Key not valid")
	}
	ctx := sdkTypes.UnwrapSDKContext(c)

	//Accept transaction if module is enabled
	if !k.GetParams(ctx).ModuleEnabled {
		return nil, cosmosTypes.ErrModuleNotEnabled
	}

	accAddr, accErr := sdkTypes.AccAddressFromBech32(msg.Voter)
	if accErr != nil {
		return nil, accErr
	}

	_, val, _, err := k.getAllValidartorOrchestratorMappingAndFindIfExist(ctx, accAddr)
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, fmt.Errorf("validator address not found")
	}

	validatorAddress, found := k.GetValidatorOrchestrator(ctx, val)
	if validatorAddress == nil {
		return nil, fmt.Errorf("unauthorized to vote")
	}
	if !found {
		return nil, cosmosTypes.ErrInvalidVote
	}

	err = k.Keeper.AddVote(ctx, msg.ProposalId, accAddr, msg.Options)
	if err != nil {
		return nil, err
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
	err := msg.ValidateBasic()
	if err != nil {
		return nil, sdkErrors.Wrap(err, "Key not valid")
	}
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
	if !(txn.CosmosTxDetails.TxHash == "") {
		return nil, cosmosTypes.ErrTxnDetailsAlreadySent
	}

	err = k.setTxDetailsSignedByOrchestrator(ctx, msg.TxID, txHash, msg.Tx)
	if err != nil {
		return nil, err
	}

	k.setTxHashAndDetails(ctx, orchAddr, msg.TxID, txHash, "pending")

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
	err := msg.ValidateBasic()
	if err != nil {
		return nil, sdkErrors.Wrap(err, "Key not valid")
	}
	ctx := sdkTypes.UnwrapSDKContext(c)

	//Accept transaction if module is enabled
	if !k.GetParams(ctx).ModuleEnabled {
		return nil, cosmosTypes.ErrModuleNotEnabled
	}

	orchAddr, orchErr := sdkTypes.AccAddressFromBech32(msg.OrchestratorAddress)
	if orchErr != nil {
		return nil, orchErr
	}

	_, val, _, err := k.getAllValidartorOrchestratorMappingAndFindIfExist(ctx, orchAddr)
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, fmt.Errorf("validator address not found")
	}

	validatorAddress, found := k.GetValidatorOrchestrator(ctx, val)
	if validatorAddress == nil {
		return nil, fmt.Errorf("unauthorized to send tx status")
	}
	if !found {
		return nil, fmt.Errorf("validator address does not exit")
	}

	if msg.Status == "success" || msg.Status == "failure" {
		k.setTxHashAndDetails(ctx, orchAddr, 0, msg.TxHash, msg.Status)
	} else {
		return nil, cosmosTypes.ErrInvalidStatus
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

func (k msgServer) RewardsClaimed(c context.Context, msg *cosmosTypes.MsgRewardsClaimedOnCosmosChain) (*cosmosTypes.MsgRewardsClaimedOnCosmosChainResponse, error) {
	ctx := sdkTypes.UnwrapSDKContext(c)

	//Accept transaction if module is enabled
	if !k.GetParams(ctx).ModuleEnabled {
		return nil, cosmosTypes.ErrModuleNotEnabled
	}

	orchAddr, orchErr := sdkTypes.AccAddressFromBech32(msg.OrchestratorAddress)
	if orchErr != nil {
		return nil, orchErr
	}

	//check if orchestrator address is present in a validator orchestrator mapping
	_, val, _, err := k.getAllValidartorOrchestratorMappingAndFindIfExist(ctx, orchAddr)
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, fmt.Errorf("validator address not found")
	}

	//check if validator exists on the network
	validatorAddress, found := k.GetValidatorOrchestrator(ctx, val)
	if validatorAddress == nil {
		return nil, fmt.Errorf("unauthorized to send tx status")
	}

	if !found {
		return nil, fmt.Errorf("validator address does not exit")
	}

	err = k.addToRewardsClaimedPool(ctx, orchAddr, msg.AmountClaimed, msg.ChainID, msg.BlockHeight)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdkTypes.NewEvent(
			sdkTypes.EventTypeMessage,
			sdkTypes.NewAttribute(sdkTypes.AttributeKeyModule, cosmosTypes.AttributeValueCategory),
			sdkTypes.NewAttribute(cosmosTypes.AttributeSender, orchAddr.String()),
		),
	)

	return &cosmosTypes.MsgRewardsClaimedOnCosmosChainResponse{}, nil
}

func (k msgServer) UndelegateSuccess(c context.Context, msg *cosmosTypes.MsgUndelegateSuccess) (*cosmosTypes.MsgUndelegateSuccessResponse, error) {
	err := msg.ValidateBasic()
	if err != nil {
		return nil, sdkErrors.Wrap(err, "Key not valid")
	}
	if err != nil {
		return nil, err
	}
	ctx := sdkTypes.UnwrapSDKContext(c)

	//Accept transaction if module is enabled
	params := k.GetParams(ctx)
	if !params.ModuleEnabled {
		return nil, cosmosTypes.ErrModuleNotEnabled
	}

	orchestratorAddress, err := sdkTypes.AccAddressFromBech32(msg.OrchestratorAddress)
	if err != nil {
		return nil, err
	}
	validatorAddress, err := sdkTypes.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}
	custodialAddress, err := sdkTypes.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	if custodialAddress.String() != params.CustodialAddress {
		return nil, cosmosTypes.ErrInvalidCustodialAddress
	}

	//check if orchestrator address is present in a validator orchestrator mapping
	_, val, _, err := k.getAllValidartorOrchestratorMappingAndFindIfExist(ctx, orchestratorAddress)
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, fmt.Errorf("validator address not found")
	}

	_, found := k.GetValidatorOrchestrator(ctx, val)
	if found {
		err = k.setUndelegateSuccessDetails(ctx, validatorAddress, orchestratorAddress, msg.Amount, msg.TxHash, msg.ChainID, msg.BlockHeight)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, cosmosTypes.ErrInvalid
	}

	ctx.EventManager().EmitEvent(
		sdkTypes.NewEvent(
			sdkTypes.EventTypeMessage,
			sdkTypes.NewAttribute(sdkTypes.AttributeKeyModule, cosmosTypes.AttributeValueCategory),
			sdkTypes.NewAttribute(cosmosTypes.AttributeSender, orchestratorAddress.String()),
		),
	)

	return &cosmosTypes.MsgUndelegateSuccessResponse{}, nil
}

func (k msgServer) SetSignature(c context.Context, msg *cosmosTypes.MsgSetSignature) (*cosmosTypes.MsgSetSignatureResponse, error) {
	ctx := sdkTypes.UnwrapSDKContext(c)
	orchestratorAddr, err := sdkTypes.AccAddressFromBech32(msg.OrchestratorAddress)
	if err != nil {
		return nil, err
	}

	outgoingTx, err := k.getTxnFromOutgoingPoolByID(ctx, msg.OutgoingTxID)
	if err != nil {
		return nil, err
	}
	if len(outgoingTx.CosmosTxDetails.Tx.AuthInfo.SignerInfos) != 1 {
		return nil, sdkErrors.Wrap(sdkErrors.ErrorInvalidSigner, "there should be exactly one signer info.")
	}
	//verify orchestator
	//TODO verify orchestratorAddress

	//verify signature
	custodialAddress, err := cosmosTypes.AccAddressFromBech32(k.GetParams(ctx).CustodialAddress, cosmosTypes.Bech32Prefix)

	multisigAccount := k.authKeeper.GetAccount(ctx, custodialAddress)
	account := k.authKeeper.GetAccount(ctx, orchestratorAddr)
	signerData := signing.SignerData{
		ChainID:       k.GetParams(ctx).CosmosProposalParams.ChainID,
		AccountNumber: multisigAccount.GetAccountNumber(),
		Sequence:      outgoingTx.CosmosTxDetails.Tx.AuthInfo.SignerInfos[0].GetSequence(),
	}
	signatureData := signing2.SingleSignatureData{
		SignMode:  signing2.SignMode_SIGN_MODE_LEGACY_AMINO_JSON,
		Signature: msg.Signature,
	}
	err = cosmosTypes.VerifySignature(account.GetPubKey(), signerData, signatureData, outgoingTx.CosmosTxDetails.Tx)
	if err != nil {
		return nil, err
	}

	//TODO add signatures to DB

	ctx.EventManager().EmitEvent(
		sdkTypes.NewEvent(
			sdkTypes.EventTypeMessage,
			sdkTypes.NewAttribute(sdkTypes.AttributeKeyModule, cosmosTypes.AttributeValueCategory),
			sdkTypes.NewAttribute(sdkTypes.AttributeKeySender, msg.OrchestratorAddress),
		),
	)
	return &cosmosTypes.MsgSetSignatureResponse{}, nil
}
