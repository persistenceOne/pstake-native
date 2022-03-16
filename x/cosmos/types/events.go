package types

//TODO Events
const (
	EventTypeSubmitProposal    = "submit_proposal"
	EventTypeOutgoing          = "outgoing_txn"
	EventTypeOutgoingCancelled = "outgoing_txn_cancelled"
	EventTypeIncoming          = "incoming_txn"
	EventTypeIncomingCancelled = "incoming_txn_cancelled"
	EventTypeOutgoingVotes     = "outgoing_votes"
	EventTypeSlashIncoming     = "slash_incoming"
	EventTypeUnbondingComplete = "unbonding_complete"
	EventTypeProposalVote      = "proposal_vote"
	EventTypeAddToOutgoingPool = "add_to_outgoing_pool"

	AttributeKeySetOperatorAddr = "set_operator_address"
	AttributeKeyOutgoingTXID    = "outgoing_tx_id"
	AttributeKeyIncomingTXID    = "incoming_tx_id"
	AttributeKeyNonce           = "nonce"
	AttributeSender             = "sender"
	AttributeKeyProposalID      = "proposal_id"
	AttributeKeyOption          = "option"
	AttributeValueCategory      = "governance"
)
