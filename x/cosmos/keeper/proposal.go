package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"time"

	sdkClient "github.com/cosmos/cosmos-sdk/client"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	sdkTx "github.com/cosmos/cosmos-sdk/types/tx"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

// creates a new proposal with given voting period in genesis
func (k Keeper) createProposal(c sdk.Context, proposal cosmosTypes.KeyAndValueForProposal) error {
	proposalID, err := k.GetProposalID(c)
	if err != nil {
		return err
	}
	submitTime := proposal.Value.VotingStartTime
	votingPeriod := proposal.Value.VotingEndTime.Sub(proposal.Value.VotingStartTime) - k.GetParams(c).CosmosProposalParams.ReduceVotingPeriodBy

	newProposal, err := cosmosTypes.NewProposal(proposalID, proposal.Value.Title, proposal.Value.Description, submitTime, votingPeriod, proposal.Value.CosmosProposalID)

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
		Voter:      "cosmos15vm0p2x990762txvsrpr26ya54p5qlz9xqlw5z",
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
		Grantee: params.CustodialAddress,
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
		EventEmitted:      true,
		Status:            "",
		TxHash:            "",
		NativeBlockHeight: ctx.BlockHeight(),
		ActiveBlockHeight: ctx.BlockHeight() + cosmosTypes.StorageWindow,
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			cosmosTypes.EventTypeOutgoing,
			sdk.NewAttribute(cosmosTypes.AttributeKeyOutgoingTXID, fmt.Sprint(nextID)),
		),
	)
	//Once event is emitted, store it in KV store for orchestrators to query transactions and sign them
	k.setNewTxnInOutgoingPool(ctx, nextID, tx)
}

// SetProposalID sets the new proposal ID to the store
func (k Keeper) SetProposalID(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(cosmosTypes.ProposalIDKey, cosmosTypes.GetProposalIDBytes(proposalID))
}

// GetProposal get proposal from store by ProposalID
func (keeper Keeper) GetProposal(ctx sdk.Context, proposalID uint64) (cosmosTypes.Proposal, bool) {
	store := ctx.KVStore(keeper.storeKey)

	bz := store.Get(cosmosTypes.ProposalKey1(proposalID))
	if bz == nil {
		return cosmosTypes.Proposal{}, false
	}

	var proposal cosmosTypes.Proposal
	err := proposal.Unmarshal(bz)
	if err != nil {
		return cosmosTypes.Proposal{}, false
	}

	return proposal, true
}

// SetProposal set a proposal to store
func (keeper Keeper) SetProposal(ctx sdk.Context, proposal cosmosTypes.Proposal) {
	store := ctx.KVStore(keeper.storeKey)

	bz, err := proposal.Marshal()
	if err != nil {
		panic("error in marshaling proposal" + err.Error())
	}

	store.Set(cosmosTypes.ProposalKey1(proposal.ProposalId), bz)
}

func (k Keeper) SetProposalPassed(ctx sdk.Context, proposalID uint64, result map[cosmosTypes.VoteOption]sdk.Dec) error {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(cosmosTypes.ProposalKey1(proposalID))
	var proposal cosmosTypes.Proposal
	err := proposal.Unmarshal(bz)
	if err != nil {
		return err
	}
	proposal.Status = cosmosTypes.StatusPassed
	proposal.FinalTallyResult.Abstain = result[cosmosTypes.OptionAbstain].RoundInt()
	proposal.FinalTallyResult.Yes = result[cosmosTypes.OptionYes].RoundInt()
	proposal.FinalTallyResult.No = result[cosmosTypes.OptionNo].RoundInt()
	proposal.FinalTallyResult.NoWithVeto = result[cosmosTypes.OptionNoWithVeto].RoundInt()
	bz, err = proposal.Marshal()
	store.Set(cosmosTypes.ProposalKey1(proposalID), bz)
	return nil
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

func (keeper Keeper) GetProposals(ctx sdk.Context) (proposals cosmosTypes.Proposals) {
	keeper.IterateProposals(ctx, func(proposal cosmosTypes.Proposal) bool {
		proposals = append(proposals, proposal)
		return false
	})
	return
}

func (keeper Keeper) IterateProposals(ctx sdk.Context, cb func(proposal cosmosTypes.Proposal) (stop bool)) {
	store := ctx.KVStore(keeper.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, cosmosTypes.ProposalsKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var proposal cosmosTypes.Proposal
		err := proposal.Unmarshal(iterator.Value())
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

func (k Keeper) setProposalDetails(ctx sdk.Context, chainID string, blockHeight int64, proposalID uint64, title string,
	description string, orchestratorAddress sdk.AccAddress, votingStartTime time.Time, votingEndTime time.Time) {
	store := ctx.KVStore(k.storeKey)
	proposalStore := prefix.NewStore(store, []byte(cosmosTypes.ProposalStoreKey))
	proposalKey := cosmosTypes.NewProposalKey(chainID, blockHeight, proposalID)
	key, err := proposalKey.Marshal()
	if err != nil {
		panic("error in marshaling proposalKey")
	}
	if proposalStore.Has(key) {
		var proposalValue cosmosTypes.ProposalValue
		err := proposalValue.Unmarshal(proposalStore.Get(key))
		if err != nil {
			panic("error in unmarshalling proposalValue")
		}
		if !proposalValue.Find(orchestratorAddress.String()) {
			proposalValue.OrchestratorAddresses = append(proposalValue.OrchestratorAddresses, orchestratorAddress.String())
			proposalValue.Counter++
			proposalValue.Ratio = float32(proposalValue.Counter) / float32(k.getTotalValidatorOrchestratorCount(ctx))
			bz, err := proposalValue.Marshal()
			if err != nil {
				panic("error in marshaling proposalValue")
			}
			proposalStore.Set(key, bz)
		}
	} else {
		ratio := float32(1) / float32(k.getTotalValidatorOrchestratorCount(ctx))
		newProposalValue := cosmosTypes.NewProposalValue(title, description, orchestratorAddress.String(), ratio, votingStartTime, votingEndTime, proposalID)
		bz, err := newProposalValue.Marshal()
		if err != nil {
			panic("error in marshaling proposalValue")
		}
		proposalStore.Set(key, bz)
	}
}

func (k Keeper) setProposalPosted(ctx sdk.Context, proposal cosmosTypes.KeyAndValueForProposal) {
	store := ctx.KVStore(k.storeKey)
	proposalStore := prefix.NewStore(store, []byte(cosmosTypes.ProposalStoreKey))
	proposalKey := cosmosTypes.NewProposalKey(proposal.Key.ChainID, proposal.Key.BlockHeight, proposal.Key.ProposalID)
	key, err := proposalKey.Marshal()
	if err != nil {
		panic("error in marshaling proposalKey")
	}
	if proposalStore.Has(key) {
		var proposalValue cosmosTypes.ProposalValue
		err := proposalValue.Unmarshal(proposalStore.Get(key))
		if err != nil {
			panic("error in unmarshalling proposalValue")
		}
		proposalValue.ProposalPosted = true
		bz, err := proposalValue.Marshal()
		if err != nil {
			panic("error in marshaling proposalValue")
		}
		proposalStore.Set(key, bz)
	}
}

func (k Keeper) getAllKeyAndValueForProposal(ctx sdk.Context) []cosmosTypes.KeyAndValueForProposal {
	store := ctx.KVStore(k.storeKey)
	proposalStore := prefix.NewStore(store, []byte(cosmosTypes.ProposalStoreKey))
	var list []cosmosTypes.KeyAndValueForProposal
	iterator := proposalStore.Iterator(nil, nil)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var key cosmosTypes.ProposalKey
		err := key.Unmarshal(iterator.Key())
		if err != nil {
			panic("error in unmarshalling proposal key")
		}
		var value cosmosTypes.ProposalValue
		err = value.Unmarshal(iterator.Value())
		if err != nil {
			panic("error in unmarshalling proposal value")
		}
		list = append(list, cosmosTypes.KeyAndValueForProposal{
			Key:   key,
			Value: value,
		})
	}
	return list
}

func (k Keeper) IterateProposalsForEmittingVotingTxn(ctx sdk.Context) {
	proposals := k.GetProposals(ctx)
	for _, proposal := range proposals {
		if !(proposal.Status == cosmosTypes.StatusPassed) {
			if proposal.VotingEndTime.Before(ctx.BlockTime()) {
				passes, tallResults := k.Tally(ctx, proposal)
				if passes {
					err := k.SetProposalPassed(ctx, proposal.ProposalId, tallResults)
					if err != nil {
						panic(err)
					}
					k.generateOutgoingWeightedVoteEvent(ctx, tallResults, proposal.CosmosProposalId)
				}
			}
		}
	}
}

func (k Keeper) ProcessProposals(ctx sdk.Context) error {
	list := k.getAllKeyAndValueForProposal(ctx)
	for _, element := range list {
		if element.Value.Ratio > 0.66 && !element.Value.ProposalPosted {
			err := k.createProposal(ctx, element)
			if err != nil {
				return err
			}
		}
	}

	fmt.Println(ctx.BlockTime(), "Current time")
	k.IterateProposalsForEmittingVotingTxn(ctx)

	return nil
}
