package types

const (
	EventTypePacket                          = "ics27_packet"
	EventTypeTimeout                         = "timeout"
	EventTypeUpdateParams                    = "update_params"
	EventTypeCreateHostChain                 = "create_host_chain"
	EventTypeUpdateHostChain                 = "update_host_chain"
	EventTypeDeleteHostChain                 = "delete_host_chain"
	EventTypeCValueUpdate                    = "c_value_update"
	EventTypeUnsuccessfulInstantiateContract = "unsuccessful_instantiate_contract"
	EventTypeUnsuccessfulExecuteContract     = "unsuccessful_execute_contract"
	EventICAChannelCreated                   = "ica_channel_created"

	AttributeID               = "id"
	AttributeChainID          = "chain_id"
	AttributeConnectionID     = "connection_id"
	AttributeUpdates          = "connection_id"
	AttributeNewCValue        = "new_c_value"
	AttributeOldCValue        = "old_c_value"
	AttributeEpoch            = "epoch_number"
	AttributeKeyAuthority     = "authority"
	AttributeKeyUpdatedParams = "updated_params"
	AttributeKeyAck           = "acknowledgement"
	AttributeKeyAckSuccess    = "success"
	AttributeKeyAckError      = "error"
	AttributeICAMessages      = "ica_messages"
	AttributeICAPortOwner     = "ica_port_owner"
	AttributeICAChannelID     = "ica_channel_id"
	AttributeICAAddress       = "ica_address"
	AttributeSender           = "msg_sender"

	AttributeValueCategory = ModuleName
)
