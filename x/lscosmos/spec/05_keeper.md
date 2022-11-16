<!--
order: 5
-->

# Keeper

## Keeper functions

`ls-cosmos` keeper module provides these utility functions to implement liquid staking

```go
// Keeper is the interface for ls-cosmos module keeper
type Keeper interface{
	// Queries
	HostChainParams(c context.Context, in *types.QueryHostChainParamsRequest) (*types.QueryHostChainParamsResponse, error)
	Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error)
	DelegationState(c context.Context, request *types.QueryDelegationStateRequest) (*types.QueryDelegationStateResponse, error)
	AllowListedValidators(c context.Context, request *types.QueryAllowListedValidatorsRequest) (*types.QueryAllowListedValidatorsResponse, error)
	CValue(c context.Context, request *types.QueryCValueRequest) (*types.QueryCValueResponse, error)
	ModuleState(c context.Context, request *types.QueryModuleStateRequest) (*types.QueryModuleStateResponse, error)
	IBCTransientStore(c context.Context, request *types.QueryIBCTransientStoreRequest) (*types.QueryIBCTransientStoreResponse, error)
	Unclaimed(c context.Context, request *types.QueryUnclaimedRequest) (*types.QueryUnclaimedResponse, error)
	FailedUnbondings(c context.Context, request *types.QueryFailedUnbondingsRequest) (*types.QueryFailedUnbondingsResponse, error)
	PendingUnbondings(c context.Context, request *types.QueryPendingUnbondingsRequest) (*types.QueryPendingUnbondingsResponse, error)
	UnbondingEpochCValue(c context.Context, request *types.QueryUnbondingEpochCValueRequest) (*types.QueryUnbondingEpochCValueResponse, error)
	HostAccountUndelegation(c context.Context, request *types.QueryHostAccountUndelegationRequest) (*types.QueryHostAccountUndelegationResponse, error)
	DelegatorUnbondingEpochEntry(c context.Context, request *types.QueryDelegatorUnbondingEpochEntryRequest) (*types.QueryDelegatorUnbondingEpochEntryResponse, error)
	RewardsBoosterAccount(c context.Context, request *types.QueryRewardBoosterAccountRequest) (*types.QueryRewardBoosterAccountResponse, error)
	HostAccounts(c context.Context, request *types.QueryHostAccountsRequest) (*types.QueryHostAccountsResponse, error)
	DepositModuleAccount(c context.Context, request *types.QueryDepositModuleAccountRequest) (*types.QueryDepositModuleAccountResponse, error)
	DelegatorUnbondingEpochEntries(c context.Context, request *types.QueryAllDelegatorUnbondingEpochEntriesRequest) (*types.QueryAllDelegatorUnbondingEpochEntriesResponse, error)
	
	// Params 
	GetParams(ctx types.Context) types.Params
	SetParams(ctx types.Context, params types.Params)
	
	// Allow listed validators 
	SetAllowListedValidators(ctx types.Context, allowlistedValidators types.AllowListedValidators)
	GetAllowListedValidators(ctx types.Context) types.AllowListedValidators

	// C value 
	GetMintedAmount(ctx types.Context) types.Int
	GetDepositAccountAmount(ctx types.Context) types.Int
	GetDelegationAccountAmount(ctx types.Context) types.Int
	GetIBCTransferTransientAmount(ctx types.Context) types.Int
	GetDelegationTransientAmount(ctx types.Context) types.Int
	GetStakedAmount(ctx types.Context) types.Int
	GetHostDelegationAccountAmount(ctx types.Context) types.Int
	GetCValue(ctx types.Context) types.Dec
	ConvertStkToToken(ctx types.Context, stkCoin types.DecCoin, cValue types.Dec) (types.Coin, types.DecCoin)
	ConvertTokenToStk(ctx types.Context, token types.DecCoin, cValue types.Dec) (types.Coin, types.DecCoin)
	
	// Host chain params
	SetHostChainParams(ctx types.Context, hostChainParams types.HostChainParams)
	GetHostChainParams(ctx types.Context) types.HostChainParams
	GetIBCDenom(ctx types.Context) string
	
	// Delegation state
	SetDelegationState(ctx types.Context, delegationState types.DelegationState)
	GetDelegationState(ctx types.Context) types.DelegationState
	AddBalanceToDelegationState(ctx types.Context, coin types.Coin)
	RemoveBalanceFromDelegationState(ctx types.Context, coins types.Coins)
	SetHostChainDelegationAddress(ctx types.Context, addr string) error
	AddHostAccountDelegation(ctx types.Context, delegation types.HostAccountDelegation)
	SubtractHostAccountDelegation(ctx types.Context, delegation types.HostAccountDelegation) error
	AddHostAccountUndelegation(ctx types.Context, undelegationEntry types.HostAccountUndelegation)
	AddTotalUndelegationForEpoch(ctx types.Context, epochNumber int64, amount types.Coin)
	AddEntriesForUndelegationEpoch(ctx types.Context, epochNumber int64, entries []types.UndelegationEntry)
	UpdateCompletionTimeForUndelegationEpoch(ctx types.Context, epochNumber int64, completionTime time.Time)
	RemoveHostAccountUndelegation(ctx types.Context, epochNumber int64) error
	GetHostAccountUndelegationForEpoch(ctx types.Context, epochNumber int64) (types.HostAccountUndelegation, error)
	GetHostAccountMaturedUndelegations(ctx types.Context) []types.HostAccountUndelegation
	
	
	// Keeper
	Logger(ctx types.Context) log.Logger
	ChanCloseInit(ctx types.Context, portID string, channelID string) error
	IsBound(ctx types.Context, portID string) bool
	BindPort(ctx types.Context, portID string) error
	AuthenticateCapability(ctx types.Context, cap *types.Capability, name string) bool
	ClaimCapability(ctx types.Context, cap *types.Capability, name string) error
	NewCapability(ctx types.Context, name string) error
	GetDepositModuleAccount(ctx types.Context) types.ModuleAccountI
	GetDelegationModuleAccount(ctx types.Context) types.ModuleAccountI
	GetRewardModuleAccount(ctx types.Context) types.ModuleAccountI
	GetUndelegationModuleAccount(ctx types.Context) types.ModuleAccountI
	GetRewardBoosterModuleAccount(ctx types.Context) types.ModuleAccountI
	MintTokens(ctx types.Context, mintCoin types.Coin, delegatorAddress types.AccAddress) error
	SendTokensToDepositModule(ctx types.Context, depositCoin types.Coins, senderAddress types.AccAddress) error
	SendTokensToRewardBoosterModuleAccount(ctx types.Context, rewardsBoostCoin types.Coins, senderAddress types.AccAddress) error
	SendProtocolFee(ctx types.Context, protocolFee types.Coins, moduleAccount string, pstakeFeeAddressString string) error
	
	// Delegator unbonding epoch entry
	SetDelegatorUnbondingEpochEntry(ctx types.Context, unbondingEpochEntry types.DelegatorUnbondingEpochEntry)
	GetDelegatorUnbondingEpochEntry(ctx types.Context, delegatorAddress types.AccAddress, epochNumber int64) types.DelegatorUnbondingEpochEntry
	RemoveDelegatorUnbondingEpochEntry(ctx types.Context, delegatorAddress types.AccAddress, epochNumber int64)
	IterateDelegatorUnbondingEpochEntry(ctx types.Context, delegatorAddress types.AccAddress) []types.DelegatorUnbondingEpochEntry
	IterateAllDelegatorUnbondingEpochEntry(ctx types.Context) []types.DelegatorUnbondingEpochEntry
	AddDelegatorUnbondingEpochEntry(ctx types.Context, delegatorAddress types.AccAddress, epochNumber int64, amount types.Coin)
	
	// Module State
	SetModuleState(ctx types.Context, enable bool)
	GetModuleState(ctx types.Context) bool
	
	// Epoch
	BeforeEpochStart(ctx types.Context, epochIdentifier string, epochNumber int64) error
	AfterEpochEnd(ctx types.Context, epochIdentifier string, epochNumber int64) error
	NewEpochHooks() EpochsHooks
	NewEpochHooks() EpochsHooks
	
	// Epoch work flows
	DelegationEpochWorkFlow(ctx types.Context, hostChainParams types.HostChainParams) error
	RewardEpochEpochWorkFlow(ctx types.Context, hostChainParams types.HostChainParams) error
	UndelegationEpochWorkFlow(ctx types.Context, hostChainParams types.HostChainParams, epochNumber int64) error
	
	// IBC packet related
	OnRecvIBCTransferPacket(ctx types.Context, packet types.Packet, relayer types.AccAddress, transferAck exported.Acknowledgement) error
	OnAcknowledgementIBCTransferPacket(ctx types.Context, packet types.Packet, acknowledgement []byte, relayer types.AccAddress, transferAckErr error) error
	OnTimeoutIBCTransferPacket(ctx types.Context, packet types.Packet, relayer types.AccAddress, transferTimeoutErr error) error
	
	// Hooks
	NewIBCTransferHooks() IBCTransferHooks
	
	// ABCI
	BeginBlock(ctx types.Context)
	
	// ABCI helpers
	DoDelegate(ctx types.Context) error
	ProcessMaturedUndelegation(ctx types.Context) error
	
	// Unboding Epoch C value
	SetUnbondingEpochCValue(ctx types.Context, unbondingEpochCValue types.UnbondingEpochCValue)
	GetUnbondingEpochCValue(ctx types.Context, epochNumber int64) types.UnbondingEpochCValue
	IterateAllUnbondingEpochCValues(ctx types.Context) []types.UnbondingEpochCValue
	MatureUnbondingEpochCValue(ctx types.Context, epochNumber int64)
	FailUnbondingEpochCValue(ctx types.Context, epochNumber int64, undelegationAmount types.Coin)
	
	// Host Accounts
	SetHostAccounts(ctx types.Context, hostAccounts types.HostAccounts)
	GetHostAccounts(ctx types.Context) types.HostAccounts
	
	// Generate and execute ICA
	GenerateAndExecuteICATx(ctx types.Context, connectionID string, portID string, msgs []types.Msg) error
	
	// IBC transient store helpers
	SetIBCTransientStore(ctx types.Context, ibcAmountTransientStore types.IBCAmountTransientStore)
	GetIBCTransientStore(ctx types.Context) types.IBCAmountTransientStore
	AddIBCTransferToTransientStore(ctx types.Context, amount types.Coin)
	RemoveIBCTransferFromTransientStore(ctx types.Context, amount types.Coin)
	AddICADelegateToTransientStore(ctx types.Context, amount types.Coin)
	RemoveICADelegateFromTransientStore(ctx types.Context, amount types.Coin)
	AddUndelegationTransferToTransientStore(ctx types.Context, undelegationTransfer types.TransientUndelegationTransfer)
	RemoveUndelegationTransferFromTransientStore(ctx types.Context, amount types.Coin) (types.TransientUndelegationTransfer, error)
	
	// Hanshake
	OnChanOpenInit(ctx types.Context, order types.Order, connectionHops []string, portID string, channelID string, chanCap *types.Capability, counterparty types.Counterparty, version string) error
	OnChanOpenTry(ctx types.Context, order types.Order, connectionHops []string, portID string, channelID string, chanCap *types.Capability, counterparty types.Counterparty, counterpartyVersion string) (string, error)
	OnChanOpenAck(ctx types.Context, portID string, channelID string, counterpartyChannelID string, counterpartyVersion string) error
	OnChanOpenConfirm(ctx types.Context, portID string, channelID string) error
	OnChanCloseInit(ctx types.Context, portID string, channelID string) error
	OnChanCloseConfirm(ctx types.Context, portID string, channelID string) error
	OnRecvPacket(ctx types.Context, modulePacket types.Packet, relayer types.AccAddress) exported.Acknowledgement
	OnAcknowledgementPacket(ctx types.Context, modulePacket types.Packet, acknowledgement []byte, relayer types.AccAddress) error
	OnTimeoutPacket(ctx types.Context, modulePacket types.Packet, relayer types.AccAddress) error
	
	// Delegation Strategy
	DelegateMsgs(ctx types.Context, delegatorAddr string, amount types.Int, denom string) ([]types.Msg, error)
	UndelegateMsgs(ctx types.Context, delegatorAddr string, amount types.Int, denom string) ([]types.Msg, []types.UndelegationEntry, error)
	GetAllValidatorsState(ctx types.Context) (types.AllowListedVals, types.HostAccountDelegations)
	
	// ICQ callbacks
	CallbackHandler() Callbacks
	HandleRewardsAccountBalanceCallback(ctx types.Context, response []byte, query types.Query) error
	
	// Host chain reward address
	SetHostChainRewardAddress(ctx types.Context, hostChainRewardAddress types.HostChainRewardAddress)
	GetHostChainRewardAddress(ctx types.Context) types.HostChainRewardAddress
	SetHostChainRewardAddressIfEmpty(ctx types.Context, hostChainRewardAddress types.HostChainRewardAddress) error
}


```