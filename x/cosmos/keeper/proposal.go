package keeper

import (
	"fmt"

	sdkClient "github.com/cosmos/cosmos-sdk/client"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	sdkTx "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/x/authz"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

type KeyAndValueForProposal struct {
	Key   cosmosTypes.ProposalKey
	Value cosmosTypes.ProposalValue
}

// creates a new proposal with given voting period in genesis
func (k Keeper) createProposal(c sdk.Context, proposal KeyAndValueForProposal) error {
	proposalID, err := k.GetProposalID(c)
	if err != nil {
		return err
	}
	submitTime := proposal.Value.ProposalDetails.VotingStartTime
	votingPeriod := proposal.Value.ProposalDetails.VotingEndTime.Sub(proposal.Value.ProposalDetails.VotingStartTime) - k.GetParams(c).CosmosProposalParams.ReduceVotingPeriodBy

	newProposal, err := cosmosTypes.NewProposal(proposalID, proposal.Value.ProposalDetails.Title, proposal.Value.ProposalDetails.Description, submitTime, votingPeriod, proposal.Value.ProposalDetails.ProposalID)
	if err != nil {
		return err
	}

	k.SetProposal(c, newProposal)
	//k.InsertActiveProposalQueue(c, proposalID, newProposal.VotingEndTime)
	k.SetProposalID(c, proposalID+1)

	k.AfterProposalSubmission(c, proposalID)

	k.setProposalPosted(c, proposal)

	c.EventManager().EmitEvent(
		sdk.NewEvent(
			cosmosTypes.EventTypeSubmitProposal,
			sdk.NewAttribute(cosmosTypes.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
		),
	)

	return nil
}

func (k Keeper) generateOutgoingWeightedVoteEvent(ctx sdk.Context, result map[cosmosTypes.VoteOption]sdk.Dec, cosmosProposalID uint64) {
	nextID := k.autoIncrementID(ctx, []byte(cosmosTypes.KeyLastTXPoolID))
	params := k.GetParams(ctx)

	var voteMsgAny []*codecTypes.Any
	msg := govTypes.MsgVoteWeighted{
		ProposalId: cosmosProposalID,
		Voter:      params.CustodialAddress,
		Options: []govTypes.WeightedVoteOption{
			{
				Option: govTypes.OptionEmpty,
				Weight: result[cosmosTypes.OptionEmpty],
			},
			{
				Option: govTypes.OptionYes,
				Weight: result[cosmosTypes.OptionYes],
			},
			{
				Option: govTypes.OptionAbstain,
				Weight: result[cosmosTypes.OptionAbstain],
			},
			{
				Option: govTypes.OptionNo,
				Weight: result[cosmosTypes.OptionNo],
			},
			{
				Option: govTypes.OptionNoWithVeto,
				Weight: result[cosmosTypes.OptionNoWithVeto],
			},
		},
	}
	msgAny, err := codecTypes.NewAnyWithValue(&msg)
	if err != nil {
		panic(err)
	}

	voteMsgAny = append(voteMsgAny, msgAny)
	execMsg := authz.MsgExec{
		Grantee: k.getCurrentAddress(ctx).String(),
		Msgs:    voteMsgAny,
	}

	execMsgAny, err := codecTypes.NewAnyWithValue(&execMsg)
	if err != nil {
		panic(err)
	}

	tx := cosmosTypes.CosmosTx{
		Tx: sdkTx.Tx{
			Body: &sdkTx.TxBody{
				Messages:      []*codecTypes.Any{execMsgAny},
				Memo:          "",
				TimeoutHeight: 0,
			},
			AuthInfo: &sdkTx.AuthInfo{
				SignerInfos: nil,
				Fee: &sdkTx.Fee{
					Amount:   nil,
					GasLimit: 200000,
					Payer:    "",
				},
			},
			Signatures: nil,
		},
		EventEmitted:      false,
		Status:            "",
		TxHash:            "",
		ActiveBlockHeight: ctx.BlockHeight() + cosmosTypes.StorageWindow,
		SignerAddress:     k.getCurrentAddress(ctx).String(),
	}

	//Once event is emitted, store it in KV store for orchestrators to query transactions and sign them
	k.setNewTxnInOutgoingPool(ctx, nextID, tx)

	k.setNewInTransactionQueue(ctx, nextID)
}

// SetProposalID sets the new proposal ID to the store
func (k Keeper) SetProposalID(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(cosmosTypes.ProposalIDKey, cosmosTypes.GetProposalIDBytes(proposalID))
}

// GetProposal get proposal from store by ProposalID
func (k Keeper) GetProposal(ctx sdk.Context, proposalID uint64) (cosmosTypes.Proposal, bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(cosmosTypes.ProposalKey1(proposalID))
	if bz == nil {
		return cosmosTypes.Proposal{}, false
	}

	var proposal cosmosTypes.Proposal
	err := k.cdc.Unmarshal(bz, &proposal)
	if err != nil {
		return cosmosTypes.Proposal{}, false
	}

	return proposal, true
}

// SetProposal set a proposal to store
func (k Keeper) SetProposal(ctx sdk.Context, proposal cosmosTypes.Proposal) {
	store := ctx.KVStore(k.storeKey)

	bz, err := k.cdc.Marshal(&proposal)
	if err != nil {
		panic("error in marshaling proposal" + err.Error())
	}

	store.Set(cosmosTypes.ProposalKey1(proposal.ProposalId), bz)
}

func (k Keeper) SetProposalPassed(ctx sdk.Context, proposalID uint64, result map[cosmosTypes.VoteOption]sdk.Dec) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(cosmosTypes.ProposalKey1(proposalID))
	var proposal cosmosTypes.Proposal
	k.cdc.MustUnmarshal(bz, &proposal)

	proposal.Status = cosmosTypes.StatusPassed
	proposal.FinalTallyResult.Abstain = result[cosmosTypes.OptionAbstain].RoundInt()
	proposal.FinalTallyResult.Yes = result[cosmosTypes.OptionYes].RoundInt()
	proposal.FinalTallyResult.No = result[cosmosTypes.OptionNo].RoundInt()
	proposal.FinalTallyResult.NoWithVeto = result[cosmosTypes.OptionNoWithVeto].RoundInt()

	store.Set(cosmosTypes.ProposalKey1(proposalID), k.cdc.MustMarshal(&proposal))
}

func (k Keeper) GetProposalID(ctx sdk.Context) (proposalID uint64, err error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(cosmosTypes.ProposalIDKey)
	if bz == nil {
		return 0, sdkErrors.Wrap(cosmosTypes.ErrInvalidGenesis, "initial proposal ID hasn't been set")
	}

	proposalID = cosmosTypes.GetProposalIDFromBytes(bz)
	return proposalID, nil
}

func (k Keeper) GetProposals(ctx sdk.Context) (proposals cosmosTypes.Proposals) {
	k.IterateProposals(ctx, func(proposal cosmosTypes.Proposal) bool {
		proposals = append(proposals, proposal)
		return false
	})
	return
}

func (k Keeper) IterateProposals(ctx sdk.Context, cb func(proposal cosmosTypes.Proposal) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, cosmosTypes.ProposalsKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var proposal cosmosTypes.Proposal
		err := k.cdc.Unmarshal(iterator.Value(), &proposal)
		if err != nil {
			panic(err)
		}

		if cb(proposal) {
			break
		}
	}
}

func (k Keeper) GetProposalsFiltered(ctx sdk.Context, params cosmosTypes.QueryProposalsRequest) cosmosTypes.Proposals {
	proposals := k.GetProposals(ctx)
	filteredProposals := make([]cosmosTypes.Proposal, 0, len(proposals))

	for _, p := range proposals {
		matchStatus := true

		// match status (if supplied/valid)
		if cosmosTypes.ValidProposalStatus(params.ProposalStatus) {
			matchStatus = p.Status == params.ProposalStatus
		}

		if matchStatus {
			filteredProposals = append(filteredProposals, p)
		}
	}

	start, end := sdkClient.Paginate(len(filteredProposals), 10, 10, 100) //TODO : Add Page and limit
	if start < 0 || end < 0 {
		filteredProposals = []cosmosTypes.Proposal{}
	} else {
		filteredProposals = filteredProposals[start:end]
	}

	return filteredProposals
}

func (k Keeper) setProposalDetails(ctx sdk.Context, msg cosmosTypes.MsgMakeProposal, validatorAddress sdk.ValAddress) {
	proposalStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.ProposalStoreKey)
	proposalKey := cosmosTypes.NewProposalKey(msg.ChainID, msg.BlockHeight, msg.ProposalID)
	key := k.cdc.MustMarshal(&proposalKey)
	totalValidatorCount := k.GetTotalValidatorOrchestratorCount(ctx)

	// store has the key in it or not
	if !proposalStore.Has(key) {
		ratio := sdk.NewDec(1).Quo(sdk.NewDec(totalValidatorCount))
		newProposalValue := cosmosTypes.NewProposalValue(msg, validatorAddress, ratio, ctx.BlockHeight()+cosmosTypes.StorageWindow)
		proposalStore.Set(key, k.cdc.MustMarshal(&newProposalValue))
		return
	}

	var proposalValue cosmosTypes.ProposalValue
	k.cdc.MustUnmarshal(proposalStore.Get(key), &proposalValue)

	// Match if the message value and stored value are same
	// if not equal then initialize by new value in store
	if !StoreValueEqualOrNotProposalEvent(proposalValue, msg) {
		ratio := sdk.NewDec(1).Quo(sdk.NewDec(totalValidatorCount))
		newProposalValue := cosmosTypes.NewProposalValue(msg, validatorAddress, ratio, ctx.BlockHeight()+cosmosTypes.StorageWindow)
		proposalStore.Set(key, k.cdc.MustMarshal(&newProposalValue))
		return
	}

	if !proposalValue.Find(validatorAddress.String()) {
		proposalValue.UpdateValues(validatorAddress.String(), totalValidatorCount)
		proposalStore.Set(key, k.cdc.MustMarshal(&proposalValue))
		return
	}
}

func (k Keeper) deleteProposalDetails(ctx sdk.Context, key cosmosTypes.ProposalKey) {
	proposalStore := prefix.NewStore(ctx.KVStore(k.storeKey), cosmosTypes.ProposalStoreKey)
	proposalStore.Delete(k.cdc.MustMarshal(&key))
}

func (k Keeper) setProposalPosted(ctx sdk.Context, proposal KeyAndValueForProposal) {
	store := ctx.KVStore(k.storeKey)
	proposalStore := prefix.NewStore(store, cosmosTypes.ProposalStoreKey)
	proposalKey := cosmosTypes.NewProposalKey(proposal.Key.ChainID, proposal.Key.BlockHeight, proposal.Key.ProposalID)
	key, err := k.cdc.Marshal(&proposalKey)
	if err != nil {
		panic("error in marshaling proposalKey")
	}
	if proposalStore.Has(key) {
		var proposalValue cosmosTypes.ProposalValue
		err := k.cdc.Unmarshal(proposalStore.Get(key), &proposalValue)
		if err != nil {
			panic("error in unmarshalling proposalValue")
		}
		proposalValue.ProposalPosted = true
		bz, err := k.cdc.Marshal(&proposalValue)
		if err != nil {
			panic("error in marshaling proposalValue")
		}
		proposalStore.Set(key, bz)
	}
}

func (k Keeper) getAllKeyAndValueForProposal(ctx sdk.Context) []KeyAndValueForProposal {
	store := ctx.KVStore(k.storeKey)
	proposalStore := prefix.NewStore(store, cosmosTypes.ProposalStoreKey)
	var list []KeyAndValueForProposal
	iterator := proposalStore.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var key cosmosTypes.ProposalKey
		err := k.cdc.Unmarshal(iterator.Key(), &key)
		if err != nil {
			panic("error in unmarshalling proposal key")
		}
		var value cosmosTypes.ProposalValue
		err = k.cdc.Unmarshal(iterator.Value(), &value)
		if err != nil {
			panic("error in unmarshalling proposal value")
		}
		list = append(list, KeyAndValueForProposal{
			Key:   key,
			Value: value,
		})
	}
	return list
}

func (k Keeper) IterateProposalsForEmittingVotingTxn(ctx sdk.Context) {
	proposals := k.GetProposals(ctx)
	for _, proposal := range proposals {
		if proposal.Status == cosmosTypes.StatusPassed && proposal.VotingEndTime.After(ctx.BlockTime()) {
			continue
		}
		passes, tallyResults := k.Tally(ctx, proposal)
		if passes {
			k.SetProposalPassed(ctx, proposal.ProposalId, tallyResults)
			k.generateOutgoingWeightedVoteEvent(ctx, tallyResults, proposal.CosmosProposalId)
		}
	}
}

func (k Keeper) ProcessProposals(ctx sdk.Context) {
	list := k.getAllKeyAndValueForProposal(ctx)
	for _, element := range list {
		if element.Value.Ratio.GT(cosmosTypes.MinimumRatioForMajority) && !element.Value.ProposalPosted {
			err := k.createProposal(ctx, element)
			if err != nil {
				panic(err)
			}
		}
		if element.Value.ActiveBlockHeight < ctx.BlockHeight() && element.Value.ProposalPosted {
			k.deleteProposalDetails(ctx, element.Key)
		}
	}

	fmt.Println(ctx.BlockTime(), "Current time")
	k.IterateProposalsForEmittingVotingTxn(ctx)
}

func StoreValueEqualOrNotProposalEvent(storeValue cosmosTypes.ProposalValue, msgValue cosmosTypes.MsgMakeProposal) bool {
	if storeValue.ProposalDetails.Title != msgValue.Title {
		return false
	}
	if storeValue.ProposalDetails.Description != msgValue.Description {
		return false
	}
	if storeValue.ProposalDetails.ProposalID != msgValue.ProposalID {
		return false
	}
	if storeValue.ProposalDetails.ChainID != msgValue.ChainID {
		return false
	}
	if storeValue.ProposalDetails.BlockHeight != msgValue.BlockHeight {
		return false
	}
	if !storeValue.ProposalDetails.VotingStartTime.Equal(msgValue.VotingStartTime) {
		return false
	}
	if !storeValue.ProposalDetails.VotingEndTime.Equal(msgValue.VotingEndTime) {
		return false
	}
	return true
}
