<!--
order: 6
-->

# Messages

In this section we describe process of the liquid staking module messages and the corresponding updates to the state. All created/modified state objects specified by each message are defined within the state section. 

### MsgLiquidStake

LiquidStake allowlisted IBC tokens and mint represntative stkTokens .

It does the following operations : 

- Validate if the amount is valid and returns an error of invalid arguments.
- Checks if the module is active and returns an error that the module is disabled if condition is not matched.
- Fetches host chain params and check if the amount deposited is more than min deposit set in store. If less, returns an error that the expected amount is more than certain limit set in store.
- Computes expected IBC prefix and checks if prefix from user and prefix in store matches. If it does not match then it returns an error of invalid denom path.
- Similar step for checking denom trace from user and stored value. If not equal, returns an error of invalid denom.
- Delegator address is checked and returns if address is invalid.
- Current C value is fetched and used to calculate amount of stk tokens to be minted corresponding to it.
- IBC token deposit is sent to deposit module account from the delegation account. If there is an error, tokens are sent back to user.
- Once the tokens are transferred to deposit module account, the above calculated stk tokens are minted and sent from module account to user account
- Protocol fees is calculated using already set parameters through governance proposal and sent to the pStake fee address.

Inputs for this message : 

- `DelegateAddress` : Address of the liquid staker.
- `Amount` : It is the amount of IBC tokens submitted by the liquid staker.

Example of a liquid staking transaction :

```
$ pstaked tx lscosmos liquid-stake 2000000ibc/DENOM_HASH  --from <delegator_key_name> --chain-id <chain_id> --keyring-backend <keyring_backend>
```

### MsgLiquidUnstake

LiquidUnstake is a transactions to unstake liquid staked assets.

It performs the following operations : 

- Checks if the module is active and returns an error that the module is disabled if condition is not matched.
- Validates if the mint denom stored matches the denom submitted by the user. If not, returns and error of invalid denom.
- Delegator address is checked and returns if address is invalid.
- Transfer tokens user wants to liquid unstake into undelegation module account.
- Protocol fees is calculated using already set parameters through governance proposal and sent to the pStake fee address.
- Entry is written in the KV store for the current unbonding epoch. 
- Another check is made to make sure that amount to be unbonded in the current unboding epoch does not overtake the amount currently staked. If that is the case then an error is returned.

Inputs for this message :

- `DelegatorAddress` : Address of the delegator wanting to liquid unstake.
- `Amount` : It is the amount of stk tokens submitted by the delegator.

```
$ pstaked tx lscosmos liquid-unstake 50000000stk/uatom --from <delegator_address> --chain-id <chain-id> --keyring-backend <keyring_backend>
```

### MsgRedeem

Redeem is a transaction to instantly withdraw tokens based on the current c value and redeem fees.

It performs the following operations : 

- Validate if the amount is valid and returns an error of invalid arguments.
- Checks if the module is active and returns an error that the module is disabled if condition is not matched.
- Delegator address is checked and returns if address is invalid.
- Validates if the mint denom stored matches the denom submitted by the user. If not, returns and error of invalid denom.
- Transfer tokens user wants to liquid unstake into lscosmos module account.
- Protocol fees is calculated using already set parameters through governance proposal and sent to the pStake fee address.
- Redeemable IBC tokens after deducting fees are transferred to user account.
- stk tokens after deduction of fees are burnt. If not burnt, an error is returned.

Inputs for this message :

- `DelegatorAddress` : Address of the delegator wanting to redeem instantly.
- `Amount` : It is the amount of stk tokens submitted by the delegator.

```
$ pstaked tx lscosmos redeem 10000000stk/uatom --from <delegator_address> --chain-id <chain-id> --keyring-backend <keyring_backend>
```

### MsgClaim

Claim is a transaction to claim all the matured unstakings or failed unstakings

It performs the following operations :

- Checks if the module is active and returns an error that the module is disabled if condition is not matched.
- Delegator address is checked and returns if address is invalid.
- All the delegator unbonding epoch entried are iterated to check status of all the entries.
- If an entry is matured then corresponding tokens are transferred to user account.
- Tokens in case of failed unstaking are also transferred back to the user's account.
- In case of any error, none of the claims are valid.

Inputs for this message : 

- `DelegatorAddress` : Address of the delegator wanting to claim matured undelegations or failed undelegations.

```
$ pstaked tx lscosmos claim --from <delegator_address> --chain-id <chain-id> --keyring-backend <keyring_backend>
```

### MsgJumpStart

JumpStart is a transactions reserved for the pstake fee address to restart the module in case of an emergency.

It performs the following operations : 

- Checks if the address in the transaction matches the pstake fee address. If not, an error is returned.
- Checks if the module is active and returns an error that the module is disabled if condition is not matched.
- Empty delegation state and host chain rewards address are set.
- Host accounts present in the message are validated and set if no error is present.
- All the pstake params are validated. If not correct then an invalid params error is returned.
- Some check pertaining to channels and ports are made to ensure correct module state.
- New capabilities are claimed for the provided channel and port.
- Delegations InterChainAccount is registered. If not then an error is returned.
- All the new host chain params and allow listed validators are set in respective stores.

Inputs for this message :

- `PstakeAddress` : Address set as pstake fee address
- `ChainID` : ChainID of blockchain on which liquid staking is aimed.
- `ConnectionID` : Connection ID for the IBC channel made for liquid staking module.
- `TransferChannel` : Transfer Channel specific to the module.
- `TransferPort` : Transfer Port the specific to the module. 
- `BaseDenom` : BaseDenom is denom of the host chain.
- `MintDenom` : Mint denom is denom to be minted.
- `MinDeposit` : Min deposit is the lower cap of the deposit that can be made while liquid staking.
- `AllowListedValidators` : Set of validators allowlisted for delegation on host chain. It also consists of weights given to each validator. 
- `PstakeParams` : Pstake params consists of different types of fees and pstake fee address
- `HostAccounts` : It is made of names for the ICA accounts to be created on host chain.

### MsgRecreateICA

RecreateICA is a transaction for recreating interchain account channels.

It performs  the following operations : 

- Checks if the module is active and returns an error that the module is disabled if condition is not matched.
- Recreates delgation account channel and rewards account channel.

Inputs for this message :

- `FromAddress` : Address from which this transaction is being sent.

```
$ pstaked tx lscosmos recreate-ica --from <from_address> --chain-id <chain-id> --keyring-backend <keyring_backend>
```

### MsgChangeModuleState

ChangeModuleState is a transaction for disabling/ reenabling the module.

It performs  the following operations :

- Checks if the module was initiated before, if no returns error
- Checks if the sender is admin
- Checks if the state is being changed.

Inputs for this message :

- `PstakeAddress` : Address from which this transaction is being sent (should be pstakeAddress).
- `ModuleState` : The boolean value true/false to which the module state is to be set.

```
$ pstaked tx lscosmos change-module-state false --from <from_address> --chain-id <chain-id> --keyring-backend <keyring_backend>
```
