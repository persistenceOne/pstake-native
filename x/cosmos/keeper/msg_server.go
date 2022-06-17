package keeper

import (
	"context"
	"fmt"
	"strconv"

	"github.com/armon/go-metrics"
	sdkTelemetry "github.com/cosmos/cosmos-sdk/telemetry"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	signing2 "github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the gov MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(k Keeper) cosmosTypes.MsgServer {
	return &msgServer{Keeper: k}
}

func (k msgServer) SetOrchestrator(c context.Context, msg *cosmosTypes.MsgSetOrchestrator) (*cosmosTypes.MsgSetOrchestratorResponse, error) {
	err := msg.ValidateBasic()
	if err != nil {
		return nil, sdkErrors.Wrap(err, "Key not valid")
	}
	ctx := sdkTypes.UnwrapSDKContext(c)

	validator, e1 := sdkTypes.ValAddressFromBech32(msg.Validator)
	orchestrator, e2 := sdkTypes.AccAddressFromBech32(msg.Orchestrator)
	if e1 != nil || e2 != nil {
		return nil, sdkErrors.Wrap(err, "Key not valid")
	}

	// check if that validator can set an orchestrator address
	valset := k.getAllOracleValidatorSet(ctx)

	found := false
	for _, val := range valset {
		if val.Address == validator.String() {
			found = true
		}
	}

	if !found {
		return nil, cosmosTypes.ErrValidatorNotAllowed
	}

	//check if orchestrator public key exist or not
	orchAccI := k.authKeeper.GetAccount(ctx, orchestrator)
	if orchAccI.GetPubKey() == nil {
		return nil, sdkErrors.Wrap(cosmosTypes.ErrPubKeyNotFound, orchestrator.String())
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

func (k msgServer) RemoveOrchestrator(c context.Context, msg *cosmosTypes.MsgRemoveOrchestrator) (*cosmosTypes.MsgRemoveOrchestratorResponse, error) {
	err := msg.ValidateBasic()
	if err != nil {
		return nil, sdkErrors.Wrap(err, "Key not valid")
	}
	ctx := sdkTypes.UnwrapSDKContext(c)

	validator, e1 := sdkTypes.ValAddressFromBech32(msg.Validator)
	orchestrator, e2 := sdkTypes.AccAddressFromBech32(msg.Orchestrator)
	if e1 != nil || e2 != nil {
		return nil, sdkErrors.Wrap(err, "Key not valid")
	}

	// removes orch address from validator mapping if it is not present in current multisig or it is the one and ony mapping
	err = k.RemoveValidatorOrchestrator(ctx, validator, orchestrator)
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

	return &cosmosTypes.MsgRemoveOrchestratorResponse{}, nil
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
	to, err := cosmosTypes.AccAddressFromBech32(msg.ToAddress, cosmosTypes.Bech32Prefix)
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

	// send amount from account to module
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

	// check if the address passed in the msg are correct or not
	destinationAddress, err := sdkTypes.AccAddressFromBech32(msg.AddressFromMemo)
	if err != nil {
		return nil, err
	}
	orchestratorAddress, err := sdkTypes.AccAddressFromBech32(msg.OrchestratorAddress)
	if err != nil {
		return nil, err
	}

	// sanity check for arguments passed in the message
	if ctx.IsZero() || sdkTypes.VerifyAddressFormat(destinationAddress) != nil || sdkTypes.VerifyAddressFormat(orchestratorAddress) != nil ||
		!msg.Amount.IsValid() {
		return nil, sdkErrors.Wrap(cosmosTypes.ErrInvalid, "arguments")
	}

	// check if the denom for staking matches or not
	uatomDenom, err := k.GetParams(ctx).GetBondDenomOf("uatom")
	if err != nil {
		return nil, err
	}
	if uatomDenom != msg.Amount.Denom {
		return nil, cosmosTypes.ErrInvalidBondDenom
	}

	val, found, err := k.getAllValidatorOrchestratorMappingAndFindIfExist(ctx, orchestratorAddress)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, fmt.Errorf("orchestrator not found")
	}

	validatorAddress, found := k.GetValidatorOrchestrator(ctx, val)
	if !found {
		return nil, cosmosTypes.ErrInvalidProposal
	}
	if validatorAddress == nil {
		return nil, fmt.Errorf("unauthorized to make proposal")
	}

	// update oracle height for both sides
	k.setOracleLastUpdateHeightCosmos(ctx, orchestratorAddress, msg.BlockHeight)
	k.setOracleLastUpdateHeightNative(ctx, orchestratorAddress, ctx.BlockHeight())

	k.addToMintTokenStore(ctx, *msg, validatorAddress)

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

	val, found, err := k.getAllValidatorOrchestratorMappingAndFindIfExist(ctx, orchestratorAddress)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, fmt.Errorf("orchestrator not found")
	}

	validatorAddress, found := k.GetValidatorOrchestrator(ctx, val)
	if !found {
		return nil, cosmosTypes.ErrInvalidProposal
	}
	if validatorAddress == nil {
		return nil, fmt.Errorf("unauthorized to make proposal")
	}

	// update oracle height for both sides
	k.setOracleLastUpdateHeightCosmos(ctx, orchestratorAddress, msg.BlockHeight)
	k.setOracleLastUpdateHeightNative(ctx, orchestratorAddress, ctx.BlockHeight())

	k.setProposalDetails(ctx, *msg, validatorAddress)

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

	val, found, err := k.getAllValidatorOrchestratorMappingAndFindIfExist(ctx, accAddr)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, fmt.Errorf("orchestrator not found")
	}

	validatorAddress, found := k.GetValidatorOrchestrator(ctx, val)
	if !found {
		return nil, cosmosTypes.ErrInvalidProposal
	}
	if validatorAddress == nil {
		return nil, fmt.Errorf("unauthorized to make proposal")
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

	val, found, err := k.getAllValidatorOrchestratorMappingAndFindIfExist(ctx, accAddr)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, fmt.Errorf("orchestrator not found")
	}

	validatorAddress, found := k.GetValidatorOrchestrator(ctx, val)
	if !found {
		return nil, cosmosTypes.ErrInvalidProposal
	}
	if validatorAddress == nil {
		return nil, fmt.Errorf("unauthorized to make proposal")
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

// TxStatus Accepts status as : "success" or "gas failure" or "sequence mismatch"
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

	orchestratorAddress, orchErr := sdkTypes.AccAddressFromBech32(msg.OrchestratorAddress)
	if orchErr != nil {
		return nil, orchErr
	}

	val, found, err := k.getAllValidatorOrchestratorMappingAndFindIfExist(ctx, orchestratorAddress)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, fmt.Errorf("orchestrator not found")
	}

	validatorAddress, found := k.GetValidatorOrchestrator(ctx, val)
	if !found {
		return nil, cosmosTypes.ErrInvalidProposal
	}
	if validatorAddress == nil {
		return nil, fmt.Errorf("unauthorized to make proposal")
	}

	// update oracle height for both sides
	k.setOracleLastUpdateHeightCosmos(ctx, orchestratorAddress, msg.BlockHeight)
	k.setOracleLastUpdateHeightNative(ctx, orchestratorAddress, ctx.BlockHeight())

	//TODO : add failure type for proposal transactions. (in case of chain upgrade on cosmos chain)
	if msg.Status == cosmosTypes.Success || msg.Status == cosmosTypes.GasFailure ||
		msg.Status == cosmosTypes.SequenceMismatch || msg.Status == cosmosTypes.KeeperFailure {
		k.setTxHashAndDetails(ctx, *msg, validatorAddress)
	} else {
		return nil, cosmosTypes.ErrInvalidStatus
	}

	ctx.EventManager().EmitEvent(
		sdkTypes.NewEvent(
			sdkTypes.EventTypeMessage,
			sdkTypes.NewAttribute(sdkTypes.AttributeKeyModule, cosmosTypes.AttributeValueCategory),
			sdkTypes.NewAttribute(cosmosTypes.AttributeSender, orchestratorAddress.String()),
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

	orchestratorAddress, orchErr := sdkTypes.AccAddressFromBech32(msg.OrchestratorAddress)
	if orchErr != nil {
		return nil, orchErr
	}

	//check if orchestrator address is present in a validator orchestrator mapping
	val, found, err := k.getAllValidatorOrchestratorMappingAndFindIfExist(ctx, orchestratorAddress)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, fmt.Errorf("orchestrator not found")
	}

	//check if validator exists on the network
	validatorAddress, found := k.GetValidatorOrchestrator(ctx, val)
	if !found {
		return nil, cosmosTypes.ErrInvalidProposal
	}
	if validatorAddress == nil {
		return nil, fmt.Errorf("unauthorized to make proposal")
	}

	// update oracle height for both sides
	k.setOracleLastUpdateHeightCosmos(ctx, orchestratorAddress, msg.BlockHeight)
	k.setOracleLastUpdateHeightNative(ctx, orchestratorAddress, ctx.BlockHeight())

	k.addToRewardsClaimedPool(ctx, *msg, validatorAddress)

	ctx.EventManager().EmitEvent(
		sdkTypes.NewEvent(
			sdkTypes.EventTypeMessage,
			sdkTypes.NewAttribute(sdkTypes.AttributeKeyModule, cosmosTypes.AttributeValueCategory),
			sdkTypes.NewAttribute(cosmosTypes.AttributeSender, orchestratorAddress.String()),
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
	_, err = cosmosTypes.ValAddressFromBech32(msg.ValidatorAddress, cosmosTypes.Bech32PrefixValAddr)
	if err != nil {
		return nil, err
	}
	custodialAddress, err := cosmosTypes.AccAddressFromBech32(msg.DelegatorAddress, cosmosTypes.Bech32Prefix)
	if err != nil {
		return nil, err
	}

	if custodialAddress.String() != params.CustodialAddress {
		return nil, cosmosTypes.ErrInvalidCustodialAddress
	}

	//check if orchestrator address is present in a validator orchestrator mapping
	val, found, err := k.getAllValidatorOrchestratorMappingAndFindIfExist(ctx, orchestratorAddress)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, fmt.Errorf("orchestrator not found")
	}

	//check if validator exists on the network
	validatorAddress, found := k.GetValidatorOrchestrator(ctx, val)
	if !found {
		return nil, cosmosTypes.ErrInvalidProposal
	}
	if validatorAddress == nil {
		return nil, fmt.Errorf("unauthorized to make proposal")
	}

	// update oracle height for both sides
	k.setOracleLastUpdateHeightCosmos(ctx, orchestratorAddress, msg.BlockHeight)
	k.setOracleLastUpdateHeightNative(ctx, orchestratorAddress, ctx.BlockHeight())

	k.setUndelegateSuccessDetails(ctx, *msg, validatorAddress)

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
	orchestratorAddress, err := sdkTypes.AccAddressFromBech32(msg.OrchestratorAddress)
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

	//check if orchestrator address is present in a validator orchestrator mapping
	val, found, err := k.getAllValidatorOrchestratorMappingAndFindIfExist(ctx, orchestratorAddress)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, fmt.Errorf("orchestrator not found")
	}

	//check if validator exists on the network
	validatorAddress, found := k.GetValidatorOrchestrator(ctx, val)
	if !found {
		return nil, cosmosTypes.ErrInvalidProposal
	}
	if validatorAddress == nil {
		return nil, fmt.Errorf("unauthorized to make proposal")
	}

	//verify signature
	custodialAddress, err := cosmosTypes.AccAddressFromBech32(outgoingTx.CosmosTxDetails.SignerAddress, cosmosTypes.Bech32Prefix)
	if err != nil {
		return nil, err
	}

	// get account state from module db
	multisigAccount := k.getAccountState(ctx, custodialAddress)
	if multisigAccount == nil {
		return nil, cosmosTypes.ErrMultiSigAddressNotFound
	}

	signerData := signing.SignerData{
		ChainID:       k.GetParams(ctx).CosmosProposalParams.ChainID,
		AccountNumber: multisigAccount.GetAccountNumber(),
		Sequence:      multisigAccount.GetSequence() + 1, // increment by 1 as it is the current sequence number is stored in the db
	}
	signatureData := signing2.SingleSignatureData{
		SignMode:  signing2.SignMode_SIGN_MODE_LEGACY_AMINO_JSON,
		Signature: msg.Signature,
	}

	account := k.authKeeper.GetAccount(ctx, orchestratorAddress)
	if account == nil {
		return nil, cosmosTypes.ErrOrchAddressNotFound
	}

	err = cosmosTypes.VerifySignature(account.GetPubKey(), signerData, signatureData, outgoingTx.CosmosTxDetails.Tx)
	if err != nil {
		return nil, err
	}

	// update oracle height for both sides
	k.setOracleLastUpdateHeightCosmos(ctx, orchestratorAddress, msg.BlockHeight)
	k.setOracleLastUpdateHeightNative(ctx, orchestratorAddress, ctx.BlockHeight())

	singleSignatureDataForOutgoingPool := cosmosTypes.ConvertSingleSignatureDataToSingleSignatureDataForOutgoingPool(signatureData)
	err = k.addToOutgoingSignaturePool(ctx, singleSignatureDataForOutgoingPool, msg.OutgoingTxID, orchestratorAddress, validatorAddress)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdkTypes.NewEvent(
			sdkTypes.EventTypeMessage,
			sdkTypes.NewAttribute(sdkTypes.AttributeKeyModule, cosmosTypes.AttributeValueCategory),
			sdkTypes.NewAttribute(sdkTypes.AttributeKeySender, orchestratorAddress.String()),
		),
	)
	return &cosmosTypes.MsgSetSignatureResponse{}, nil
}

func (k msgServer) SlashingEvent(c context.Context, msg *cosmosTypes.MsgSlashingEventOnCosmosChain) (*cosmosTypes.MsgSlashingEventOnCosmosChainResponse, error) {
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
	_, err = cosmosTypes.ValAddressFromBech32(msg.ValidatorAddress, cosmosTypes.Bech32PrefixValAddr)
	if err != nil {
		return nil, err
	}

	//check if orchestrator address is present in a validator orchestrator mapping
	val, found, err := k.getAllValidatorOrchestratorMappingAndFindIfExist(ctx, orchestratorAddress)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, fmt.Errorf("orchestrator not found")
	}

	//check if validator exists on the network
	validatorAddress, found := k.GetValidatorOrchestrator(ctx, val)
	if !found {
		return nil, cosmosTypes.ErrInvalidProposal
	}
	if validatorAddress == nil {
		return nil, fmt.Errorf("unauthorized to make proposal")
	}

	// update oracle height for both sides
	k.setOracleLastUpdateHeightCosmos(ctx, orchestratorAddress, msg.BlockHeight)
	k.setOracleLastUpdateHeightNative(ctx, orchestratorAddress, ctx.BlockHeight())

	k.setSlashingEventDetails(ctx, *msg, validatorAddress)

	ctx.EventManager().EmitEvent(
		sdkTypes.NewEvent(
			sdkTypes.EventTypeMessage,
			sdkTypes.NewAttribute(sdkTypes.AttributeKeyModule, cosmosTypes.AttributeValueCategory),
			sdkTypes.NewAttribute(cosmosTypes.AttributeSender, orchestratorAddress.String()),
		),
	)

	return &cosmosTypes.MsgSlashingEventOnCosmosChainResponse{}, nil
}
