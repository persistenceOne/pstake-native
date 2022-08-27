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

It is to whitelist these params before the modules is able to liquid stake atoms on Persistence Chain.

Example of a register host chain proposal :

```go
{
"title": "register host chain proposal",
"description": "this proposal register host chain params in the chain",
"connection_i_d": "test connection",
"transfer_channel": "test-channel-1",
"transfer_port": "test-transfer",
"base_denom": "uatom",
"mint_denom": "ustkatom",
"min_deposit": "5",
"pstake_deposit_fee": "0.1",
"pstake_restake_fee": "0.1",
"pstake_unstake_fee": "0.1",
"deposit": "100stake"
}
```

Sample command to submit a host chain register proposal :

```
$ $BIN_NAME tx gov submit-proposal register-host-chain <path/to/proposal.json> --from <key_or_address> --fees <1000stake> --gas <200000>
```