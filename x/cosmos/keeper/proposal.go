package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	cosmosTypes "github.com/persistenceOne/pstake-native/x/cosmos/types"
)

func (k Keeper) CreateProposal(c types.Context, proposal cosmosTypes.KeyAndValueForProposal) error {
	proposalID, err := k.GetProposalID(c)
	if err != nil {
		return err
	}
	submitTime := c.BlockHeader().Time
	votingPeriod := k.GetParams(c).CosmosProposalParams.VotingPeriod

	newProposal, err := cosmosTypes.NewProposal(proposalID, proposal.Value.Title, proposal.Value.Description, submitTime, votingPeriod)

	k.SetProposal(c, newProposal)
	k.InsertActiveProposalQueue(c, proposalID, newProposal.VotingEndTime)
	k.SetProposalID(c, proposalID+1)

	k.AfterProposalSubmission(c, proposalID)

	k.setProposalPosted(c, proposal)

	c.EventManager().EmitEvent(
		types.NewEvent(
			cosmosTypes.EventTypeSubmitProposal,
			types.NewAttribute(cosmosTypes.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
		),
	)

	return nil
}

// SetProposalID sets the new proposal ID to the store
func (k Keeper) SetProposalID(ctx types.Context, proposalID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(cosmosTypes.ProposalIDKey, cosmosTypes.GetProposalIDBytes(proposalID))
}

// GetProposal get proposal from store by ProposalID
func (keeper Keeper) GetProposal(ctx types.Context, proposalID uint64) (cosmosTypes.Proposal, bool) {
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
func (keeper Keeper) SetProposal(ctx types.Context, proposal cosmosTypes.Proposal) {
	store := ctx.KVStore(keeper.storeKey)

	bz, err := proposal.Marshal()
	if err != nil {
		panic("error in marshaling proposal" + err.Error())
	}

	store.Set(cosmosTypes.ProposalKey1(proposal.ProposalId), bz)
}

func (k Keeper) GetProposalID(ctx types.Context) (proposalID uint64, err error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(cosmosTypes.ProposalIDKey)
	if bz == nil {
		return 0, sdkerrors.Wrap(cosmosTypes.ErrInvalidGenesis, "initial proposal ID hasn't been set")
	}

	proposalID = cosmosTypes.GetProposalIDFromBytes(bz)
	return proposalID, nil
}

func (k Keeper) setProposalDetails(ctx types.Context, chainID string, blockHeight int64, proposalID int64, title string, description string, orchestratorAddress types.AccAddress) {
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
			panic("error in marshaling proposalValue")
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
		newProposalValue := cosmosTypes.NewProposalValue(title, description, orchestratorAddress.String(), ratio)
		bz, err := newProposalValue.Marshal()
		if err != nil {
			panic("error in marshaling proposalValue")
		}
		proposalStore.Set(key, bz)
	}
}

func (k Keeper) setProposalPosted(ctx types.Context, proposal cosmosTypes.KeyAndValueForProposal) {
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
			panic("error in marshaling proposalValue")
		}
		proposalValue.ProposalPosted = true
		bz, err := proposalValue.Marshal()
		if err != nil {
			panic("error in marshaling proposalValue")
		}
		proposalStore.Set(key, bz)
	}
}

func (k Keeper) GetAllKeyAndValueForProposal(ctx types.Context) []cosmosTypes.KeyAndValueForProposal {
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
