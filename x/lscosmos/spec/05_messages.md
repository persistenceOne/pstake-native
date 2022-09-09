<!--
order: 5
-->

# Messages

### MsgLiquidStake

LiquidStake whitelisted IBC tokens and mint represntative stkTokens .

The message will fail under the following conditions:

- The coins are not whitelisted for liquid-staking
- The delegator address is restricted

### MsgJuice

Juice is a transaction to boost rewards on the protocol.

This message will fail under following conditions:
- The coins are not whitelisted for liquid-staking
- The rewarder address is restricted
