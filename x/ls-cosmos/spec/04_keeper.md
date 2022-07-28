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
    mintTokens(ctx sdk.Context, mintCoin sdk.Coin, mintAddress sdk.AccAddress) error
}
```