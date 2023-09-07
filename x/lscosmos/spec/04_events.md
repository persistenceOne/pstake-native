<!--
order: 4
-->

# Events

The lscosmos module emits the following events:

## Handlers

### MsgLiquidStake

Send liquid stake tokens from account to module (transfer IBC tokens to module account)

| Type     | Attribute Key | Attribute Value    |
|----------|---------------|--------------------|
| transfer | recipient     | {recipientAddress} |
| transfer | amount        | {amount}           |
| message  | action        | send               |
| transfer | sender        | {senderAddress}    |
| message  | sender        | {fromAddress}      |

Mint stk tokens

| Type     | Attribute Key | Attribute Value |
|----------|---------------|-----------------|
| coinbase | minter        | {minter}        |
| coinbase | amount        | {amount}        |

Send tokens from module to account (transfer mint stk tokens)

| Type     | Attribute Key | Attribute Value    |
|----------|---------------|--------------------|
| transfer | recipient     | {recipientAddress} |
| transfer | amount        | {amount}           |
| message  | action        | send               |
| transfer | sender        | {senderAddress}    |
| message  | sender        | {fromAddress}      |

Send tokens from module to account (transfer mint stk tokens as fees to declared fee address)

| Type     | Attribute Key | Attribute Value    |
|----------|---------------|--------------------|
| transfer | recipient     | {recipientAddress} |
| transfer | amount        | {amount}           |
| message  | action        | send               |
| transfer | sender        | {senderAddress}    |
| message  | sender        | {fromAddress}      |

Last event

| Type         | Attribute Key      | Attribute Value                     |
|--------------|--------------------|-------------------------------------|
| liquid-stake | address            | {recipientAddress}                  |
| liquid-stake | amount             | {mintSTKAmount}                     |
| liquid-stake | received           | {mintSTKAmount - protocolFeeAmount} |
| liquid-stake | pstake-deposit-fee | {protocolFeeAmount}                 |
| message      | module             | lscosmos                            |
| message      | sender             | {address}                           |

### MsgLiquidUnstake

Send liquid unstake tokens from account to module (transfer stkTokens to module account)

| Type     | Attribute Key | Attribute Value    |
|----------|---------------|--------------------|
| transfer | recipient     | {recipientAddress} |
| transfer | amount        | {amount}           |
| message  | action        | send               |
| transfer | sender        | {senderAddress}    |
| message  | sender        | {fromAddress}      |

Send tokens from module to account (transfer stk tokens from undelegation module account as fees to declared fee
address)

| Type     | Attribute Key | Attribute Value    |
|----------|---------------|--------------------|
| transfer | recipient     | {recipientAddress} |
| transfer | amount        | {amount}           |
| message  | action        | send               |
| transfer | sender        | {senderAddress}    |
| message  | sender        | {fromAddress}      |

Last event

| Type           | Attribute Key       | Attribute Value      |
|----------------|---------------------|----------------------|
| liquid-unstake | address             | {delegatiorAddress}  |
| liquid-unstake | received            | {amount}             |
| liquid-unstake | undelegation-amount | {undelegationAmount} |
| liquid-unstake | pstake-unstake-fee  | {protocolFeeAmount}  |
| message        | module              | lscosmos             |
| message        | sender              | {address}            |

### MsgRedeem

Send redeeem tokens from account to module (transfer stkTokens to module account)

| Type     | Attribute Key | Attribute Value    |
|----------|---------------|--------------------|
| transfer | recipient     | {recipientAddress} |
| transfer | amount        | {amount}           |
| message  | action        | send               |
| transfer | sender        | {senderAddress}    |
| message  | sender        | {fromAddress}      |

Send tokens from module to account (transfer stk tokens as fees to declared fee address)

| Type     | Attribute Key | Attribute Value    |
|----------|---------------|--------------------|
| transfer | recipient     | {recipientAddress} |
| transfer | amount        | {amount}           |
| message  | action        | send               |
| transfer | sender        | {senderAddress}    |
| message  | sender        | {fromAddress}      |

Send tokens from module to account (transfer IBC tokens from deposit module account)

| Type     | Attribute Key | Attribute Value    |
|----------|---------------|--------------------|
| transfer | recipient     | {recipientAddress} |
| transfer | amount        | {amount}           |
| message  | action        | send               |
| transfer | sender        | {senderAddress}    |
| message  | sender        | {fromAddress}      |

Burn stk tokens

| Type | Attribute Key | Attribute Value |
|------|---------------|-----------------|
| burn | burner        | {burner}        |
| burn | amount        | {amount}        |

Last event

| Type    | Attribute Key    | Attribute Value     |
|---------|------------------|---------------------|
| redeem  | address          | {delegatiorAddress} |
| redeem  | amount           | {amount}            |
| redeem  | received         | {redeemAmount}      |
| redeem  | pstake-redem-fee | {protocolFeeAmount} |
| message | module           | lscosmos            |
| message | sender           | {address}           |

### MsgRedeem

Send redeeem tokens from account to module (transfer stkTokens to module account)

| Type     | Attribute Key | Attribute Value    |
|----------|---------------|--------------------|
| transfer | recipient     | {recipientAddress} |
| transfer | amount        | {amount}           |
| message  | action        | send               |
| transfer | sender        | {senderAddress}    |
| message  | sender        | {fromAddress}      |

Send tokens from module to account (transfer stk tokens as fees to declared fee address)

| Type     | Attribute Key | Attribute Value    |
|----------|---------------|--------------------|
| transfer | recipient     | {recipientAddress} |
| transfer | amount        | {amount}           |
| message  | action        | send               |
| transfer | sender        | {senderAddress}    |
| message  | sender        | {fromAddress}      |

Send tokens from module to account (transfer IBC tokens from deposit module account)

| Type     | Attribute Key | Attribute Value    |
|----------|---------------|--------------------|
| transfer | recipient     | {recipientAddress} |
| transfer | amount        | {amount}           |
| message  | action        | send               |
| transfer | sender        | {senderAddress}    |
| message  | sender        | {fromAddress}      |

Burn stk tokens

| Type | Attribute Key | Attribute Value |
|------|---------------|-----------------|
| burn | burner        | {burner}        |
| burn | amount        | {amount}        |

Last event

| Type    | Attribute Key    | Attribute Value     |
|---------|------------------|---------------------|
| redeem  | address          | {delegatiorAddress} |
| redeem  | amount           | {amount}            |
| redeem  | received         | {redeemAmount}      |
| redeem  | pstake-redem-fee | {protocolFeeAmount} |
| message | module           | lscosmos            |
| message | sender           | {address}           |

### MsgClaim

Case 1 : Send redeeem tokens from module to account (transfer stkTokens to module account)

| Type     | Attribute Key | Attribute Value           |
|----------|---------------|---------------------------|
| transfer | recipient     | {recipientAddress}        |
| transfer | amount        | {amount}                  |
| message  | action        | send                      |
| transfer | sender        | undelegationModuleAccount |
| message  | sender        | {fromAddress}             |

After sending tokens from module to account emits a general claim event

| Type     | Attribute Key  | Attribute Value    |
|----------|----------------|--------------------|
| claim    | address        | {recipientAddress} |
| claim    | amount         | {amount}           |
| claim    | claimed-amount | send               |

Case 2 :
Send tokens from module to account (transfer stk tokens from undelegation module account)

| Type     | Attribute Key | Attribute Value           |
|----------|---------------|---------------------------|
| transfer | recipient     | {recipientAddress}        |
| transfer | amount        | {amount}                  |
| message  | action        | send                      |
| transfer | sender        | undelegationModuleAccount |
| message  | sender        | {fromAddress}             |

Last event

| Type    | Attribute Key    | Attribute Value     |
|---------|------------------|---------------------|
| message | module           | lscosmos            |
| message | sender           | {address}           |

### MsgJumpStart

| Type       | Attribute Key    | Attribute Value     |
|------------|------------------|---------------------|
| jump-start | pstake-address   | {address}           |
| message    | module           | lscosmos            |
| message    | sender           | {address}           |

### MsgRecreateICA

| Type        | Attribute Key           | Attribute Value          |
|-------------|-------------------------|--------------------------|
| recreat-ica | pstake-address          | {address}                |
| recreat-ica | recreate-delegation-ica | {delegatorAccountPortID} |
| recreat-ica | recreate-rewards-ica    | {rewardsAccountPortID}   |
| message     | module                  | lscosmos                 |
| message     | sender                  | {address}                |
