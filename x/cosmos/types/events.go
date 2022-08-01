package types

const (
	EventTypeSubmitProposal   = "submit_proposal"
	EventTypeOutgoing         = "outgoing_txn"
	EventTypeSignedOutgoing   = "signed_tx"
	EventTypeProposalVote     = "proposal_vote"
	EventModuleEnableProposal = "module_enable_proposal"

	AttributeKeySetOperatorAddr = "set_operator_address"
	AttributeKeyOutgoingTXID    = "outgoing_tx_id"
	AttributeSender             = "sender"
	AttributeKeyProposalID      = "proposal_id"
	AttributeKeyOption          = "option"
	AttributeMultisigAddress    = "multisig_address"
	AttributeValueCategory      = ModuleName
)
