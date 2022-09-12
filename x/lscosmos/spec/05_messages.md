<!--
order: 5
-->

# Messages

### MsgLiquidStake

LiquidStake whitelisted IBC tokens and mint represntative stkTokens .

It does the following operations : 

- Validate the message and returns error if any check fails
- Checks if the module is enabled or disabled. If disabled and returns.
- Amount and IBC denoms are validated. If denom trace does not match the ones present in the proposal then the transaction fails.
- If all the above conditions pass then deposited tokens are sent to deposit module account.
- Once the tokens are submitted, some fees is cut from those tokens and the remaining tokens are minted in liquid staker account depending on the current C-Vaule.

Inputs for this message : 

- `DelegateAddress` : Address of the liquid staker.
- `Amount` : It the amount of IBC tokens submitted by the liquid staker.

Example of a liquid staking transaction :

```
$ $BIN_NAME tx lscosmos pstake-lscosmos-liquid-stake 2000000ibc/DENOM_HASH  --from <delegator_key_name> --chain-id <chain_id> --keyring-backend <keyring_backend>
```

### MsgJuice

Juice is a transaction to boost rewards on the protocol.

This message will fail under following conditions:
- The coins are not whitelisted for liquid-staking
- The rewarder address is restricted
