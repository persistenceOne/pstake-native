package types

const (
	// keyvals
	HostChainKeyVal     string = "host_chain"
	WorkflowKeyVal      string = "workflow"
	EpochKeyVal         string = "epoch"
	ErrorKeyVal         string = "error"
	ValidatorKeyVal     string = "validator"
	SequenceIDKeyVal    string = "sequence_id"
	ChannelKeyVal       string = "channel"
	OwnerKeyVal         string = "owner"
	AddressKeyVal       string = "address"
	PortKeyVal          string = "port"
	DelegatorKeyVal     string = "delegator"
	HostDenomKeyVal     string = "host_denom"
	MessagesKeyVal      string = "messages"
	AmountKeyVal        string = "amount"
	FromValidatorKeyVal string = "from_validator"
	ToValidatorKeyVal   string = "to_validator"
	CValueKeyVal        string = "c_value"
	LowerLimitKeyVal    string = "lower_limit"
	UpperLimitKeyVal    string = "upper_limit"

	DelegateWorkflow                  string = "DoDelegate"
	ClaimWorkflow                     string = "DoClaim"
	ProcessMatureUndelegationWorkflow string = "DoProcessMatureUndelegations"
	RecreateICAWorkflow               string = "DoRecreateICA"
	RedeemLSMWorkflow                 string = "DoRedeemLSMTokens"
)
