<!--
order: 0
-->

# Introduction
In a proof-of-stake blockchain, asset holders wishing to stake their tokens in order to earn staking rewards are required to lock their capital. This locked capital, known as a delegation, acts as a slashable security deposit in the event that the chosen validator should misbehave.

Staking with the Cosmos ecosystem requires 2 parties: a validator and a delegator. The validator runs the hardware and a Cosmos SDK-based blockchain application to propose and validate new blocks on the chain in consensus with other validators. In order to prevent misbehaviour by validators, delegators put up capital in the form of the chain’s native staking tokens by way of security deposit. In the event of misbehaviour (accidental or malicious), 0.1% of the staked deposit for persistent downtime and 5% for a double signing violation (whereby the validator signs more than one time for a given block height), is burned or slashed as punishment.

To compensate for providing security deposits, delegators earn staking rewards proportional to the value of their staked assets for each block validated. Validators in turn charge a commission upon these rewards for providing the validator service.

In order to protect against certain classes of attack, whereby validators can commit offences and either they themselves, or their delegators, are able to unbond their deposit before the misbehaviour is reported, an unbonding period (usually on the order of 3 weeks) is instituted. During this unbonding period, the staked assets do not accrue rewards and are returned to the delegator as liquid, unlocked assets at the end of the unbonding period.

Hence, the maximum amount that a user may be slashed for any misbehaviour on most, if not all, Cosmos SDK networks is five (5) percent, leaving almost 95% of capital untouched.
