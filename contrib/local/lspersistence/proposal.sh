#!/bin/sh

CHAIN_BIN=pstaked

# Submit a param-change governance proposal
$CHAIN_BIN tx gov submit-proposal param-change ../configs/proposal_whitelist_validator.json \
  --chain-id localnet \
  --from user1 \
  --keyring-backend test \
  --broadcast-mode block \
  --yes \
  --gas auto \
  --gas-adjustment 2.0 \
  --home /tmp/trash/.pstaked \
  --output json | jq

# Query the proposal to check the status PROPOSAL_STATUS_VOTING_PERIOD
$CHAIN_BIN q gov proposals --output json | jq

# Vote
$CHAIN_BIN tx gov vote 1 yes \
  --chain-id localnet \
  --from val1 \
  --keyring-backend test \
  --broadcast-mode block \
  --yes \
  --home /tmp/trash/.pstaked \
  --output json | jq

#
# Wait a while (30s) for the proposal to pass
#
sleep 30
# Query the proposal again to check the status PROPOSAL_STATUS_PASSED
$CHAIN_BIN q gov proposals --home /tmp/trash/.pstaked --output json | jq

# Query the values set as liquidstaking parameters and liquid-validators
$CHAIN_BIN q lspersistence params --home /tmp/trash/.pstaked --output json | jq
$CHAIN_BIN q lspersistence liquid-validators --home /tmp/trash/.pstaked --home /tmp/trash/.pstaked --output json | jq

