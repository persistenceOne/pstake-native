# Liquid Stake IBC

## Abstarct

The Persistence chain module `x/liquidstakeibc`, is the main implementation of the Persistence liquid staking protocol.

## Contents

1. [Concepts](#Concepts)
2. [State](#State)
3. [Proposals](#Proposals)
4. [Mesages](#Messages)
5. [Events](#Events)
6. [Queries](#Queries)
7. [Keepers](#Keepers)
8. [Parameters](#Parameters)

## Concepts

### Liquid Staking

`Liquid Staking`, as a concept, is a mechanism in which delegations are made liquid and can be transformed, traded, 
or otherwise utilised.

The goal of liquid staking is to allow delegators to maintain their staked position while simultaneously allowing 
them to seek out the best returns for their capital.

This is achieved by `minting` an asset representative of the native bonded token at the point of delegation, which 
can then in turn be used by DeFi protocols.

### Host Chain

A `Host Chain` in the Liquid Stake IBC module represents an IBC connected blockchain, whose base token can be liquid 
staked using the `x/liquidstakeibc` module. An example of that would be the `gaia` chain and its base asset `ATOM`.

### C Value

The `c_value` of an LST (Liquid Staked Token) is the effective ratio between the total amount of minted representative
tokens, `stkAssets`, and the total amount of native assets staked on the host chain. This ratio is used calculate the
amount of stkAssets which will be minted by the module when performing a liquid stake action, and the amount of 
stkAssets which will be burned when unbonding.

## State

### HostChain

A `HostChain` represents an IBC connected blockchain which is registered on the module and accepts liquid stake delegations.

```go
type HostChain struct {
    // host chain id
    ChainId string                                             `protobuf:"bytes,1,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
    // ibc connection id
    ConnectionId string                                        `protobuf:"bytes,2,opt,name=connection_id,json=connectionId,proto3" json:"connection_id,omitempty"`
    // module params
    Params *HostChainLSParams                                  `protobuf:"bytes,3,opt,name=params,proto3" json:"params,omitempty"`
    // native token denom
    HostDenom string                                           `protobuf:"bytes,4,opt,name=host_denom,json=hostDenom,proto3" json:"host_denom,omitempty"`
    // ibc connection channel id
    ChannelId string                                           `protobuf:"bytes,5,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
    // ibc connection port id
    PortId string                                              `protobuf:"bytes,6,opt,name=port_id,json=portId,proto3" json:"port_id,omitempty"`
    // delegation host account
    DelegationAccount *ICAAccount                              `protobuf:"bytes,7,opt,name=delegation_account,json=delegationAccount,proto3" json:"delegation_account,omitempty"`
    // reward host account
    RewardsAccount *ICAAccount                                 `protobuf:"bytes,8,opt,name=rewards_account,json=rewardsAccount,proto3" json:"rewards_account,omitempty"`
    // validator set
    Validators []*Validator                                    `protobuf:"bytes,9,rep,name=validators,proto3" json:"validators,omitempty"`
    // minimum ls amount
    MinimumDeposit github_com_cosmos_cosmos_sdk_types.Int      `protobuf:"bytes,10,opt,name=minimum_deposit,json=minimumDeposit,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"minimum_deposit"`
    // redemption rate
    CValue github_com_cosmos_cosmos_sdk_types.Dec              `protobuf:"bytes,11,opt,name=c_value,json=cValue,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"c_value"`
    // previous redemption rate
    LastCValue github_com_cosmos_cosmos_sdk_types.Dec          `protobuf:"bytes,12,opt,name=last_c_value,json=lastCValue,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"last_c_value"`
    // undelegation epoch factor
    UnbondingFactor int64                                      `protobuf:"varint,13,opt,name=unbonding_factor,json=unbondingFactor,proto3" json:"unbonding_factor,omitempty"`
    // whether the chain is ready to accept delegations or not
    Active bool                                                `protobuf:"varint,14,opt,name=active,proto3" json:"active,omitempty"`
    // factor limit for auto-compounding, daily periodic rate (APY / 365s)
    AutoCompoundFactor github_com_cosmos_cosmos_sdk_types.Dec  `protobuf:"bytes,15,opt,name=auto_compound_factor,json=autoCompoundFactor,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"auto_compound_factor"`
    // host chain flags
    Flags *HostChainFlags                                      `protobuf:"bytes,16,opt,name=flags,proto3" json:"flags,omitempty"`
}
```

### HostChainFlags

The `HostChainFlags` are used to determine arbitrary attributes of a host chain that can be added while the chain is
already registered.

```go
type HostChainFlags struct {
	// whether the chain accepts LSM delegations or not
    Lsm bool `protobuf:"varint,1,opt,name=lsm,proto3" json:"lsm,omitempty"`
}
```

### HostChainLSParams

The `HostChainLSParams` determine module wide params for the given host chain. They are mainly used for fee purposes.

```go
type HostChainLSParams struct {
    DepositFee    github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,1,opt,name=deposit_fee,json=depositFee,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"deposit_fee"`
    RestakeFee    github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,2,opt,name=restake_fee,json=restakeFee,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"restake_fee"`
    UnstakeFee    github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,3,opt,name=unstake_fee,json=unstakeFee,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"unstake_fee"`
    RedemptionFee github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,4,opt,name=redemption_fee,json=redemptionFee,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"redemption_fee"`
}
```

### ICAAccount

An `ICAAccount` represents an account in the host chain which the module has control over.

```go
type ICAAccount struct {
    // address of the ica on the controller chain
    Address string                       `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
    // token balance of the ica
    Balance types.Coin                   `protobuf:"bytes,2,opt,name=balance,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coin" json:"balance"`
    // owner string
    Owner        string                  `protobuf:"bytes,3,opt,name=owner,proto3" json:"owner,omitempty"`
	// the state of the ICA channel
    ChannelState ICAAccount_ChannelState `protobuf:"varint,4,opt,name=channel_state,json=channelState,proto3,enum=pstake.liquidstakeibc.v1beta1.ICAAccount_ChannelState" json:"channel_state,omitempty"`
}
```
```go
enum ChannelState {
    // ICA channel is being created
    ICA_CHANNEL_CREATING = 0;
    // ICA is established and the account can be used
    ICA_CHANNEL_CREATED = 1;
}
```

### Validator

A `Validator` represents a validator on the host chain which is tracked by the module and will receive its delegations.

```go
type Validator struct {
    // valoper address
    OperatorAddress string                                 `protobuf:"bytes,1,opt,name=operator_address,json=operatorAddress,proto3" json:"operator_address,omitempty"`
    // validator status
    Status string                                          `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
    // validator weight in the set
    Weight github_com_cosmos_cosmos_sdk_types.Dec          `protobuf:"bytes,3,opt,name=weight,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"weight"`
    // amount delegated by the module to the validator
    DelegatedAmount github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,4,opt,name=delegated_amount,json=delegatedAmount,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"delegated_amount"`
    // the validator token exchange rate, total bonded tokens divided by total shares issued
    ExchangeRate github_com_cosmos_cosmos_sdk_types.Dec    `protobuf:"bytes,5,opt,name=exchange_rate,json=exchangeRate,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"exchange_rate"`
    // the unbonding epoch number when the validator transitioned into the state
    UnbondingEpoch int64                                   `protobuf:"varint,6,opt,name=unbonding_epoch,json=unbondingEpoch,proto3" json:"unbonding_epoch,omitempty"`
}
```

### Deposit

A `Deposit` represents all the delegations that the module received within one epoch.

```go
type Deposit struct {
// deposit target chain
ChainId string             `protobuf:"bytes,1,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
Amount  types.Coin         `protobuf:"bytes,2,opt,name=amount,proto3" json:"amount"`
// epoch number of the deposit
Epoch int64                `protobuf:"varint,3,opt,name=epoch,proto3" json:"epoch,omitempty"`
// state
State Deposit_DepositState `protobuf:"varint,4,opt,name=state,proto3,enum=pstake.liquidstakeibc.v1beta1.Deposit_DepositState" json:"state,omitempty"`
// sequence id of the ibc transaction
IbcSequenceId string       `protobuf:"bytes,5,opt,name=ibc_sequence_id,json=ibcSequenceId,proto3" json:"ibc_sequence_id,omitempty"`
}
```
```go
const (
    // no action has been initiated on the deposit
    LSMDeposit_DEPOSIT_PENDING LSMDeposit_LSMDepositState = 0
    // deposit sent to the host chain delegator address
    LSMDeposit_DEPOSIT_SENT LSMDeposit_LSMDepositState = 1
    // deposit received by the host chain delegator address
    LSMDeposit_DEPOSIT_RECEIVED LSMDeposit_LSMDepositState = 2
    // deposit started the untokenization process
    LSMDeposit_DEPOSIT_UNTOKENIZING LSMDeposit_LSMDepositState = 3
)
```

### LSMDeposit

An `LSMDeposit` behaves the same way as a `Deposit` but for LSM delegations.

```go
type LSMDeposit struct {
    // deposit target chain
    ChainId string                                `protobuf:"bytes,1,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
    // this is calculated when liquid staking [lsm_shares * validator_exchange_rate]
    Amount github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,2,opt,name=amount,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"amount"`
    // LSM token shares, they are mapped 1:1 with the delegator shares that are tokenized
    // https://github.com/iqlusioninc/cosmos-sdk/pull/19
    Shares github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,3,opt,name=shares,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"shares"`
    // LSM token denom
    Denom string                                  `protobuf:"bytes,4,opt,name=denom,proto3" json:"denom,omitempty"`
    // LSM token ibc denom
    IbcDenom string                               `protobuf:"bytes,5,opt,name=ibc_denom,json=ibcDenom,proto3" json:"ibc_denom,omitempty"`
    // address of the delegator
    DelegatorAddress string                       `protobuf:"bytes,6,opt,name=delegator_address,json=delegatorAddress,proto3" json:"delegator_address,omitempty"`
    // state o the deposit
    State LSMDeposit_LSMDepositState              `protobuf:"varint,7,opt,name=state,proto3,enum=pstake.liquidstakeibc.v1beta1.LSMDeposit_LSMDepositState" json:"state,omitempty"`
    // sequence id of the ibc transaction
    IbcSequenceId string                          `protobuf:"bytes,8,opt,name=ibc_sequence_id,json=ibcSequenceId,proto3" json:"ibc_sequence_id,omitempty"`
}
```
```go
const (
    // no action has been initiated on the deposit
    LSMDeposit_DEPOSIT_PENDING LSMDeposit_LSMDepositState = 0
    // deposit sent to the host chain delegator address
    LSMDeposit_DEPOSIT_SENT LSMDeposit_LSMDepositState = 1
    // deposit received by the host chain delegator address
    LSMDeposit_DEPOSIT_RECEIVED LSMDeposit_LSMDepositState = 2
    // deposit started the untokenization process
    LSMDeposit_DEPOSIT_UNTOKENIZING LSMDeposit_LSMDepositState = 3
)
```

### Unbonding

An `Unbonding` represents all the undelegations that the module received within one epoch.

```go
type Unbonding struct {
    // unbonding target chain
    ChainId string                 `protobuf:"bytes,1,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
    // epoch number of the unbonding record
    EpochNumber int64              `protobuf:"varint,2,opt,name=epoch_number,json=epochNumber,proto3" json:"epoch_number,omitempty"`
    // time when the unbonding matures and can be collected
    MatureTime time.Time           `protobuf:"bytes,3,opt,name=mature_time,json=matureTime,proto3,stdtime" json:"mature_time"`
    // stk token amount that is burned with the unbonding
    BurnAmount types.Coin          `protobuf:"bytes,4,opt,name=burn_amount,json=burnAmount,proto3" json:"burn_amount"`
    // host token amount that is being unbonded
    UnbondAmount types.Coin        `protobuf:"bytes,5,opt,name=unbond_amount,json=unbondAmount,proto3" json:"unbond_amount"`
    // sequence id of the ibc transaction
    IbcSequenceId string           `protobuf:"bytes,6,opt,name=ibc_sequence_id,json=ibcSequenceId,proto3" json:"ibc_sequence_id,omitempty"`
    // state of the unbonding during the process
    State Unbonding_UnbondingState `protobuf:"varint,7,opt,name=state,proto3,enum=pstake.liquidstakeibc.v1beta1.Unbonding_UnbondingState" json:"state,omitempty"`
}
```
```go
const (
    // no action has been initiated on the unbonding
    Unbonding_UNBONDING_PENDING Unbonding_UnbondingState = 0
    // unbonding action has been sent to the host chain
    Unbonding_UNBONDING_INITIATED Unbonding_UnbondingState = 1
    // unbonding is waiting for the maturing period of the host chain
    Unbonding_UNBONDING_MATURING Unbonding_UnbondingState = 2
    // unbonding has matured and is ready to transfer from the host chain
    Unbonding_UNBONDING_MATURED Unbonding_UnbondingState = 3
    // unbonding is on the persistence chain and can be claimed
    Unbonding_UNBONDING_CLAIMABLE Unbonding_UnbondingState = 4
    // unbonding has failed
    Unbonding_UNBONDING_FAILED Unbonding_UnbondingState = 5
)
```

### UserUnbonding

A `UserUnbonding` maps a user specific unbonding to the corresponding `Unbonding` object.

```go
type UserUnbonding struct {
    // unbonding target chain
    ChainId string          `protobuf:"bytes,1,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
    // epoch when the unbonding started
    EpochNumber int64       `protobuf:"varint,2,opt,name=epoch_number,json=epochNumber,proto3" json:"epoch_number,omitempty"`
    // address which requested the unbonding
    Address string          `protobuf:"bytes,3,opt,name=address,proto3" json:"address,omitempty"`
    // stk token amount that is being unbonded
    StkAmount types.Coin    `protobuf:"bytes,4,opt,name=stk_amount,json=stkAmount,proto3" json:"stk_amount"`
    // host token amount that is being unbonded
    UnbondAmount types.Coin `protobuf:"bytes,5,opt,name=unbond_amount,json=unbondAmount,proto3" json:"unbond_amount"`
}
```

### ValidatorUnbonding

A `ValidatorUnbonding` represents a full validator unbonding, that is, all the bonded tokens on that validator
as a result of the validator transitioning into Unbonding/Unbonded state.

```go
type ValidatorUnbonding struct {
    // unbonding target chain
    ChainId string          `protobuf:"bytes,1,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
    // epoch when the unbonding started
    EpochNumber int64       `protobuf:"varint,2,opt,name=epoch_number,json=epochNumber,proto3" json:"epoch_number,omitempty"`
    // time when the unbonding matures and can be collected
    MatureTime time.Time    `protobuf:"bytes,3,opt,name=mature_time,json=matureTime,proto3,stdtime" json:"mature_time"`
    // address of the validator that is being unbonded
    ValidatorAddress string `protobuf:"bytes,4,opt,name=validator_address,json=validatorAddress,proto3" json:"validator_address,omitempty"`
    // amount unbonded from the validator
    Amount types.Coin       `protobuf:"bytes,5,opt,name=amount,proto3" json:"amount"`
    // sequence id of the ibc transaction
    IbcSequenceId string    `protobuf:"bytes,6,opt,name=ibc_sequence_id,json=ibcSequenceId,proto3" json:"ibc_sequence_id,omitempty"`
}
```

### KVUpdate

A `KVUpdate` represents a simple KV pair used to update a host chain.

```go
type KVUpdate struct {
    Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
    Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}
```

## Proposals

### register-host-chain

Proposal to register a host chain.

`pstaked tx gov submit-proposal [proposal-file]`

```json
{
  "messages": [{
    "@type": "/pstake.liquidstakeibc.v1beta1.MsgRegisterHostChain",
    "authority": "persistence10d07y265gmmuvt4z0w9aw880jnsr700j5w4kch",
    "connection_id": "connection-0",
    "channel_id": "channel-0",
    "port_id": "transfer",
    "deposit_fee": "0.00",
    "restake_fee": "0.05",
    "unstake_fee": "0.00",
    "redemption_fee": "0.005",
    "host_denom": "uatom",
    "minimum_deposit": "1",
    "unbonding_factor": "4"
  }],
  "deposit": "10000000uxprt",
  "proposer": "persistence1hcqg5wj9t42zawqkqucs7la85ffyv08ljhhesu",
  "title": "Register Cosmos",
  "summary": "Registers the Cosmos host chain to pStake",
  "metadata": ""
}
```

### update-host-chain

Proposal to update a host chain set of attributes.

`pstaked tx gov submit-proposal [proposal-file]`

```json
{
  "messages": [{
    "@type": "/pstake.liquidstakeibc.v1beta1.MsgUpdateHostChain",
    "authority": "persistence10d07y265gmmuvt4z0w9aw880jnsr700j5w4kch",
    "chain_id": "gaia-1",
    "updates": [
      {
        "key": "validator_weight",
        "value": "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt,0.5"
      },
      {
        "key": "validator_weight",
        "value": "cosmosvaloper1lcck2cxh7dzgkrfk53kysg9ktdrsjj6jfwlnm2,0.5"
      }]
  }],
  "deposit": "10000000uxprt",
  "proposer": "persistence1hcqg5wj9t42zawqkqucs7la85ffyv08ljhhesu",
  "title": "Update Cosmos validator weight",
  "summary": "Updates the Cosmos validator weights.",
  "metadata": ""
}
```

### update-params

Proposal to update the module params.

`pstaked tx gov submit-proposal [proposal-file]`

```json
{
  "messages": [{
    "@type": "/pstake.liquidstakeibc.v1beta1.MsgUpdateParams",
    "authority": "persistence10d07y265gmmuvt4z0w9aw880jnsr700j5w4kch",
    "params": {
      "admin_address": "persistence10khgeppewe4rgfrcy809r9h00aquwxxxrk6glr",
      "fee_address": "persistence1xruvjju28j0a5ud5325rfdak8f5a04h0s30mld"
    }
  }],
  "deposit": "10000000uxprt",
  "proposer": "persistence1hcqg5wj9t42zawqkqucs7la85ffyv08ljhhesu",
  "title": "Update module addresses",
  "summary": "Updates both the admin and the fee address of the module",
  "metadata": ""
}
```

## Messages

```protobuf
// Msg defines the liquidstakeibc services.
service Msg {
  rpc RegisterHostChain(MsgRegisterHostChain) returns (MsgRegisterHostChainResponse);
  rpc UpdateHostChain(MsgUpdateHostChain) returns (MsgUpdateHostChainResponse);

  rpc LiquidStake(MsgLiquidStake) returns (MsgLiquidStakeResponse) {
    option (google.api.http).post = "/pstake/liquidstakeibc/v1beta1/LiquidStake";
  }

  rpc LiquidStakeLSM(MsgLiquidStakeLSM) returns (MsgLiquidStakeLSMResponse) {
    option (google.api.http).post = "/pstake/liquidstakeibc/v1beta1/LiquidStakeLSM";
  }

  rpc LiquidUnstake(MsgLiquidUnstake) returns (MsgLiquidUnstakeResponse) {
    option (google.api.http).post = "/pstake/liquidstakeibc/v1beta1/LiquidUnstake";
  }

  rpc Redeem(MsgRedeem) returns (MsgRedeemResponse) {
    option (google.api.http).post = "/pstake/liquidstakeibc/v1beta1/Redeem";
  }

  rpc UpdateParams(MsgUpdateParams) returns (MsgUpdateParamsResponse);
}
```

### MsgRegisterHostChain

Saves the host chain object into the module store and initiates ICA channel creation.

It can only be executed by either the `gov` module account or the module admin account.

```go
type MsgRegisterHostChain struct {
    Authority          string                                 `protobuf:"bytes,1,opt,name=authority,proto3" json:"authority,omitempty"`
    ConnectionId       string                                 `protobuf:"bytes,2,opt,name=connection_id,json=connectionId,proto3" json:"connection_id,omitempty"`
    DepositFee         github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,3,opt,name=deposit_fee,json=depositFee,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"deposit_fee"`
    RestakeFee         github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,4,opt,name=restake_fee,json=restakeFee,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"restake_fee"`
    UnstakeFee         github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,5,opt,name=unstake_fee,json=unstakeFee,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"unstake_fee"`
    RedemptionFee      github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,6,opt,name=redemption_fee,json=redemptionFee,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"redemption_fee"`
    ChannelId          string                                 `protobuf:"bytes,7,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
    PortId             string                                 `protobuf:"bytes,8,opt,name=port_id,json=portId,proto3" json:"port_id,omitempty"`
    HostDenom          string                                 `protobuf:"bytes,9,opt,name=host_denom,json=hostDenom,proto3" json:"host_denom,omitempty"`
    MinimumDeposit     github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,10,opt,name=minimum_deposit,json=minimumDeposit,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"minimum_deposit"`
    UnbondingFactor    int64                                  `protobuf:"varint,11,opt,name=unbonding_factor,json=unbondingFactor,proto3" json:"unbonding_factor,omitempty"`
    AutoCompoundFactor int64                                  `protobuf:"varint,12,opt,name=auto_compound_factor,json=autoCompoundFactor,proto3" json:"auto_compound_factor,omitempty"`
}
```

### MsgUpdateHostChain

Updates different attributes of a host chain using KV pairs.

It can only be executed by either the `gov` module account or the module admin account.

```go
type MsgUpdateHostChain struct {
    Authority string      `protobuf:"bytes,1,opt,name=authority,proto3" json:"authority,omitempty"`
    ChainId   string      `protobuf:"bytes,2,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
    Updates   []*KVUpdate `protobuf:"bytes,3,rep,name=updates,proto3" json:"updates,omitempty"`
}
```
The available attributes to update are the following: 

```go
const (
    KeyAddValidator       string = "add_validator"
    KeyRemoveValidator    string = "remove_validator"
    KeyValidatorSlashing  string = "validator_slashing"
    KeyValidatorWeight    string = "validator_weight"
    KeyDepositFee         string = "deposit_fee"
    KeyRestakeFee         string = "restake_fee"
    KeyUnstakeFee         string = "unstake_fee"
    KeyRedemptionFee      string = "redemption_fee"
    KeyMinimumDeposit     string = "min_deposit"
    KeyActive             string = "active"
    KeySetWithdrawAddress string = "set_withdraw_address"
    KeyAutocompoundFactor string = "autocompound_factor"
    KeyFlags              string = "flags"
)
```

The `KeyValidatorSlashing` is used to update a specific validator exchange rate and status manually, which is done in
response to a slashing event.

### MsgLiquidStake

Adds the message amount to the current delegation epoch deposit and mints the corresponding stkAssets using the host
chain c value.

```go
type MsgLiquidStake struct {
    DelegatorAddress string     `protobuf:"bytes,1,opt,name=delegator_address,json=delegatorAddress,proto3" json:"delegator_address,omitempty"`
    Amount           types.Coin `protobuf:"bytes,2,opt,name=amount,proto3" json:"amount"`
}
```

### MsgLiquidStakeLSM

Untokenizes the given LSM delegations immediately, to avoid high price impact and mints the corresponding stkAssets using the host
chain c value.

```go
type MsgLiquidStakeLSM struct {
    DelegatorAddress string                                   `protobuf:"bytes,1,opt,name=delegator_address,json=delegatorAddress,proto3" json:"delegator_address,omitempty"`
    Delegations      github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,2,rep,name=delegations,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"delegations"`
}
```

### MsgLiquidUnstake

Adds the message amount to the current unbonding epoch record and burns the corresponding stkAssets using the host
chain c value.

```go
type MsgLiquidUnstake struct {
    DelegatorAddress string     `protobuf:"bytes,1,opt,name=delegator_address,json=delegatorAddress,proto3" json:"delegator_address,omitempty"`
    Amount           types.Coin `protobuf:"bytes,2,opt,name=amount,proto3" json:"amount"`
}
```

### MsgRedeem

Attempts to instantly redeem stkAssets by using the current epoch deposit amount. If there is not enough deposited amount
the message will fail.

```go
type MsgRedeem struct {
    DelegatorAddress string     `protobuf:"bytes,1,opt,name=delegator_address,json=delegatorAddress,proto3" json:"delegator_address,omitempty"`
    Amount           types.Coin `protobuf:"bytes,2,opt,name=amount,proto3" json:"amount"`
}
```

### MsgUpdateParams

Updates the current module params.

It can only be executed by either the `gov` module account or the module admin account.

```go
type MsgUpdateParams struct {
    Authority string `protobuf:"bytes,1,opt,name=authority,proto3" json:"authority,omitempty"`
    Params    Params `protobuf:"bytes,2,opt,name=params,proto3" json:"params"`
}
```

## Events

List of the events emitted by the module.

### LiquidStake

| Type         | Attribute Key      | Attribute Value     |
|:-------------|:-------------------|:--------------------|
| message      | module             | liquidstakeibc      |
| message      | sender             | {delegator_address} |
| liquid-stake | address            | {delegator_address} |
| liquid-stake | amount             | {staked_amount}     |
| liquid-stake | received           | {amount_received}   |
| liquid-stake | pstake-deposit-fee | {deposit_fee}       |

### LiquidUnstake

| Type            | Attribute Key       | Attribute Value      |
|:----------------|:--------------------|:---------------------|
| message         | module              | liquidstakeibc       |
| message         | sender              | {delegator_address}  |
| liquid-unstake  | address             | {delegator_address}  |
| liquid-unstake  | pstake-unstake-fee  | {unstake_fee}        |
| liquid-unstake  | received            | {amount_received}    |
| liquid-unstake  | undelegation-amount | {undelegated_amount} |
| liquid-unstake  | undelegation-epoch  | {undelegation_epoch} |

### Redeem

| Type            | Attribute Key      | Attribute Value     |
|:----------------|:-------------------|:--------------------|
| message         | module             | liquidstakeibc      |
| message         | sender             | {delegator_address} |
| liquid-unstake  | address            | {delegator_address} |
| liquid-unstake  | amount             | {redeem_amount}     |
| liquid-unstake  | received           | {amount_received}   |
| liquid-unstake  | pstake-redeem-fee  | {redeem_fee}        |

### UpdateParams

| Type            | Attribute Key     | Attribute Value   |
|:----------------|:------------------|:------------------|
| message         | module            | liquidstakeibc    |
| liquid-unstake  | authority         | {authority}       |
| liquid-unstake  | updated_params    | {updated_params}  |

## Queries

```protobuf
// Query defines the gRPC querier service.
service Query {
  // Queries the parameters of the module.
  rpc Params(QueryParamsRequest) returns (QueryParamsResponse) {
    option (google.api.http).get = "/pstake/liquidstakeibc/v1beta1/params";
  }

  // Queries a HostChain by id.
  rpc HostChain(QueryHostChainRequest) returns (QueryHostChainResponse) {
    option (google.api.http).get = "/pstake/liquidstakeibc/v1beta1/host_chain/{chain_id}";
  }

  // Queries for all the HostChains.
  rpc HostChains(QueryHostChainsRequest) returns (QueryHostChainsResponse) {
    option (google.api.http).get = "/pstake/liquidstakeibc/v1beta1/host_chains";
  }

  // Queries for all the deposits for a host chain.
  rpc Deposits(QueryDepositsRequest) returns (QueryDepositsResponse) {
    option (google.api.http).get = "/pstake/liquidstakeibc/v1beta1/deposits/{chain_id}";
  }

  // Queries for all the deposits for a host chain.
  rpc LSMDeposits(QueryLSMDepositsRequest) returns (QueryLSMDepositsResponse) {
    option (google.api.http).get = "/pstake/liquidstakeibc/v1beta1/lsm_deposits/{chain_id}";
  }

  // Queries all unbondings for a host chain.
  rpc Unbondings(QueryUnbondingsRequest) returns (QueryUnbondingsResponse) {
    option (google.api.http).get = "/pstake/liquidstakeibc/v1beta1/unbondings/{chain_id}";
  }

  // Queries an unbonding for a host chain.
  rpc Unbonding(QueryUnbondingRequest) returns (QueryUnbondingResponse) {
    option (google.api.http).get = "/pstake/liquidstakeibc/v1beta1/unbonding/{chain_id}/{epoch}";
  }

  // Queries all unbondings for a delegator address.
  rpc UserUnbondings(QueryUserUnbondingsRequest) returns (QueryUserUnbondingsResponse) {
    option (google.api.http).get = "/pstake/liquidstakeibc/v1beta1/user_unbondings/{address}";
  }

  // Queries all validator unbondings for a host chain.
  rpc ValidatorUnbondings(QueryValidatorUnbondingRequest) returns (QueryValidatorUnbondingResponse) {
    option (google.api.http).get = "/pstake/liquidstakeibc/v1beta1/validator_unbondings/{chain_id}";
  }

  // Queries for a host chain deposit account balance.
  rpc DepositAccountBalance(QueryDepositAccountBalanceRequest) returns (QueryDepositAccountBalanceResponse) {
    option (google.api.http).get = "/pstake/liquidstakeibc/v1beta1/deposit_account_balance/{chain_id}";
  }

  // Queries for a host chain exchange rate between the host token and the stk token.
  rpc ExchangeRate(QueryExchangeRateRequest) returns (QueryExchangeRateResponse) {
    option (google.api.http).get = "/pstake/liquidstakeibc/v1beta1/exchange_rate/{chain_id}";
  }
}
```

## Keepers

https://github.com/persistenceOne/pstake-native/blob/main/x/liquidstakeibc/keeper/keeper.go

## Parameters

Module parameters:

| Key                      | Type   | Default |
|:-------------------------|:-------|:--------|
| admin_address            | string | N/A     |
| fee_address              | string | N/A     |
| upper_c_value_limit      | string | "0.85"  |
| lower_c_value_limit      | string | "1.1"   |


Description of parameters:

* `admin_address` - admin account of the module, which is used to perform high privilege operations.
* `fee_address` - address that gathers fees on the module.
* `upper_c_value_limit` - module-wide c value upper hard limit.
* `lower_c_value_limit` - module-wide c value lower hard limit.
