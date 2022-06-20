package types

const (
	EventTypeSubmitProposal = "submit_proposal"
	EventTypeOutgoing       = "outgoing_txn"
	EventTypeSignedOutgoing = "signed_tx"
	EventTypeProposalVote   = "proposal_vote"

	AttributeKeySetOperatorAddr = "set_operator_address"
	AttributeKeyOutgoingTXID    = "outgoing_tx_id"
	AttributeSender             = "sender"
	AttributeKeyProposalID      = "proposal_id"
	AttributeKeyOption          = "option"
	AttributeValueCategory      = ModuleName
)
