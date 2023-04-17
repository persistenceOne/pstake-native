<!-- order: 6 -->

# Hooks

## SetAdditionalVotingPowers (Liquid Governance)

SetAdditionalVotingPowers calculates the voter's voting power who owns `stkToken` that considers the following factors:

- Balance of stkToken
- Balance of PoolCoins including stkToken

The calculated voting power is added, deducted, or overwritten with `AdditionalVotingPowers` inside the tally logic
of `cosmos-sdk/x/gov` module. It is called in `govHooks.SetAdditionalVotingPowers`.

Each voting power of `AdditionalVotingPowers` is distributed to liquid validators by their weight of **bonded**
liquidTokens each liquid validators has **bonded** status of `cosmos-sdk/x/staking` module states     
