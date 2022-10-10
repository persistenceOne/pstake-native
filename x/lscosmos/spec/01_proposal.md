<!--
order: 1
-->

# Register Host Chain Proposal

Host chain register proposal takes the following parameter into account :

- `ConnectionID` : IBC Connection to be whitelisted for module use.
- `TransferChannel` : Channel whitelisted for token transfer.
- `TransferPort` : Port whitelisted for token transfer.
- `BaseDenom` : Base denom to be matched with the base denom present in IBC/Denom and stake only the whitelisted base
  denom.
- `MintDenom` : Mint denom which is to be used for minting the liquid staked representative of staked token.
- `MinDeposit` : It is the amount that a user needs to deposit in order to do a liquid stake transaction.
- `AllowListedValidators` : It is the list of validators with corresponding weights used for delegation strategy.
- `PstakeDepositFee` : It is the fraction of fees on liquid staking deposit.
- `PstakeRestakeFee` : It is the fraction of fees on auto restake of rewards.
- `PstakeUnstakeFee` : It is the fraction of fees on unstake transaction.
- `PstakeRedemptionFee` : It is the fraction of fees on instant redemption.
- `PstakeFeeAddress` : Address on which the above mentioned fees are sent to.

It is to whitelist these params before the modules is able to liquid stake atoms on Persistence Chain.

Example of a register host chain proposal :

```json
{
  "title": "register host chain proposal",
  "description": "this proposal register host chain params in the chain",
  "module_enabled": true,
  "chain_id": "test-1",
  "connection_id": "connection-0",
  "transfer_channel": "channel-0",
  "transfer_port": "transfer",
  "base_denom": "uatom",
  "mint_denom": "stk/uatom",
  "min_deposit": "1",
  "allow_listed_validators": {
    "allow_listed_validators": [
      {
        "validator_address": "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt",
        "target_weight": "1"
      }
    ]
  },
  "pstake_params": {
    "pstake_deposit_fee": "0.00",
    "pstake_restake_fee": "0.05",
    "pstake_unstake_fee": "0.00",
    "pstake_redemption_fee": "0.1",
    "pstake_fee_address": "persistence108cqtjz7gqasctvrw74kewg6642062kmfuujsd"
  }
  "deposit": "1000000stake"
}
```

Sample command to submit a host chain register proposal :

```
$ $BIN_NAME tx gov submit-proposal pstake-lscosmos-min-deposit-and-fee-change  <path/to/proposal.json> --from <key_or_address> --fees <1000stake> --gas <200000>
```

# Change Min Deposit and Fee Proposal

This proposal takes the following parameters into account :

- `MinDeposit` : of type sdk.Int
- `PstakeDepositFee` : of type sdk.Dec
- `PstakeRestakeFee` : of type sdk.Dec
- `PstakeUnstakeFee` : of type sdk.Dec
- `PstakeRedemptionFee` : of type sdk.Dec

Example proposal :

```json
{
  "title": "min-deposit and fee change proposal",
  "description": "this proposal changes min-deposit and protocol fee on chain",
  "min_deposit": "5",
  "pstake_deposit_fee": "0.1",
  "pstake_restake_fee": "0.1",
  "pstake_unstake_fee": "0.1",
  "pstake_redemption_fee": "0.1",
  "deposit": "1000000stake"
}
```

Sample command to submit proposal :

```
$ $BIN_NAME tx gov submit-proposal pstake-lscosmos-min-deposit-and-fee-change  <path/to/proposal.json> --from <key_or_address> --fees <1000stake> --gas <200000>
```

It used when there is a need to change all the different types of fees and minimum deposit needed by user to liquid
stake.

# Change Pstake Fee Address Proposal

This proposal take the following prameters into account :

- `PstakeFeeAddress` : It is the address in which pstake protocol fee is sent.

Example proposal :

```json
{
  "title": "change pstake fee address",
  "description": "this proposal changes pstake fee address in the chain",
  "pstake_fee_address": "persistence1pss7nxeh3f9md2vuxku8q99femnwdjtcpe9ky9"
  "deposit": "10000000stake"
}
```

Sample command to submit proposal :

```
$ $BIN_NAME tx gov submit-proposal pstake-lscosmos-change-pstake-fee-address <path/to/proposal.json> --from <key_or_address> --fees <1000stake> --gas <200000>
```

It is used to change the pstake fee address if old one is not needed anymore.

# Change Allow Listed Validators Proposal

This proposal take the following prameters into account :

- `AllowListedValidators` : List of validators and their corresponding weights.

Example proposal :

```json
{
  "title": "change pstake fee address",
  "description": "this proposal changes pstake fee address in the chain",
  "allow_listed_validators": {
    "allow_listed_validators": [
      {
        "validator_address": "cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt",
        "target_weight": "1"
      }
    ]
  },
  "deposit": "100000stake"
}
```

Sample command to submit proposal :

```
$ $BIN_NAME tx gov submit-proposal pstake-lscosmos-change-allow-listed-validator-set <path/to/proposal.json> --from <key_or_address> --fees <1000stake> --gas <200000>
```

It is to change the validator set of lscsomos module if in case the old validator set becomes stale.


