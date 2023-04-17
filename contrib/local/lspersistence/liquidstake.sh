#!/bin/sh

export CHAIN_BIN=pstaked
export CHAIN_DATA_DIR="${CHAIN_DATA_DIR:=.pstaked}"

$CHAIN_BIN q bank balances persistence152nvc6f096v6n6tr5lg50xq22ak0chsr0ru7xc \
  --home /tmp/trash/.pstaked \
  --output json | jq

# Liquid Stake
$CHAIN_BIN tx lspersistence liquid-stake 1000000000stake \
  --chain-id localnet \
  --from user1 \
  --keyring-backend test \
  --broadcast-mode block \
  --yes \
  --home /tmp/trash/.pstaked \
  --output json | jq
# Query liquid validators
$CHAIN_BIN q lspersistence liquid-validators --home /tmp/trash/.pstaked --output json | jq


# Query balance of user1, you can find 1000000000bstake balance
$CHAIN_BIN q bank balances persistence152nvc6f096v6n6tr5lg50xq22ak0chsr0ru7xc \
  --home /tmp/trash/.pstaked \
  --output json | jq \


# Query liquid staking states including net amount, mint rate
#$CHAIN_BIN q lspersistence states --output json | jq

# Liquid UnStake
$CHAIN_BIN tx lspersistence liquid-unstake 500000000bstake \
  --gas 400000 \
  --chain-id localnet \
  --from user1 \
  --keyring-backend test \
  --broadcast-mode block \
  --yes \
  --home /tmp/trash/.pstaked \
  --output json | jq \


# Query balance of user1, you can find 500000000bstake balance left
$CHAIN_BIN q bank balances persistence152nvc6f096v6n6tr5lg50xq22ak0chsr0ru7xc \
  --home /tmp/trash/.pstaked \
  --output json | jq \

# Query liquid validators, you can find del_shares, liquid_tokens 500000000.000000000000000000 + withdrawn and re-staked rewards + UnstakeFee (0.001)
$CHAIN_BIN q lspersistence liquid-validators --home /tmp/trash/.pstaked --output json | jq


# Query unbonding, 499500000(UnstakeFee(0.001) deducted) + rewards
$CHAIN_BIN q staking unbonding-delegations persistence152nvc6f096v6n6tr5lg50xq22ak0chsr0ru7xc --home /tmp/trash/.pstaked --output json | jq


# Query balance of liquidstaking proxy module account
$CHAIN_BIN q bank balances persistence152nvc6f096v6n6tr5lg50xq22ak0chsr0ru7xc --home /tmp/trash/.pstaked \
  --output json | jq \
