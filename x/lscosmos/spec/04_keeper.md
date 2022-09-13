<!--
order: 4
-->

# Keeper

## Keeper functions

`ls-cosmos` keeper module provides utility functions to implement liquid staking

```go
// Keeper is the interface for ls-cosmos module keeper
type Keeper interface{
	
    // MintTokens for a given whitelisted IBC Token
    mintTokens(ctx sdk.Context, mintCoin sdk.Coin, delegatorAddress sdk.AccAddress) error 
    //SendTokens to the DepositModuleAccount
    SendTokensToDepositModule(ctx sdk.Context, depositCoin sdk.Coins, senderAddress sdk.AccAddress) 
    //Send residue coins to CommunityPool
    SendResidueToCommunityPool(ctx sdk.Context, residue []sdk.DecCoin)
    //Send ProtocolFee to protocol community pool
    SendProtocolFee(ctx sdk.Context, protocolFee sdk.Coins, delegatorAddr sdk.AccAddress) error
}


```