#!/bin/sh

export BINARY=pstaked

$BINARY q bank balances persistence10khgeppewe4rgfrcy809r9h00aquwxxxrk6glr \
--output json | jq

# Liquid Stake
$BINARY tx lspersistence liquid-stake 1000000000stake \
--chain-id localnet \
--from user2 \
--keyring-backend test \
--broadcast-mode block \
--yes \
--output json | jq

# Query liquid validators
$BINARY q lspersistence liquid-validators --output json | jq

# Query balance of user2, you can find 1000000000bstake balance
$BINARY q bank balances persistence10khgeppewe4rgfrcy809r9h00aquwxxxrk6glr \
--output json | jq

# Query liquid staking states including net amount, mint rate
#$BINARY q lspersistence states --output json | jq


# Liquid UnStake
$BINARY tx lspersistence liquid-unstake 500000000bstake \
--gas 400000 \
--chain-id localnet \
--from user2 \
--keyring-backend test \
--broadcast-mode block \
--yes \
--output json | jq

# Query balance of user2, you can find 500000000bstake balance left
$BINARY q bank balances persistence10khgeppewe4rgfrcy809r9h00aquwxxxrk6glr \
--output json | jq

# Query liquid validators, you can find del_shares, liquid_tokens 500000000.000000000000000000 + withdrawn and re-staked rewards + UnstakeFee (0.001)
$BINARY q lspersistence liquid-validators --output json | jq

# Query unbonding, 499500000(UnstakeFee(0.001) deducted) + rewards
$BINARY q staking unbonding-delegations persistence10khgeppewe4rgfrcy809r9h00aquwxxxrk6glr --output json | jq

# Query balance of liquidstaking proxy module account
$BINARY q bank balances persistence10khgeppewe4rgfrcy809r9h00aquwxxxrk6glr \
--output json | jq