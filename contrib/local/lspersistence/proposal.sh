#!/bin/sh

export CHAIN_BIN=pstaked

# Submit a param-change governance proposal
$CHAIN_BIN tx gov submit-proposal param-change ~/proposal_whitelist_validator.json \
  --chain-id localnet \
  --from user1 \
  --keyring-backend test \
  --broadcast-mode block \
  --yes \
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
  --output json | jq

#
# Wait a while (30s) for the proposal to pass
#
sleep 30
# Query the proposal again to check the status PROPOSAL_STATUS_PASSED
$CHAIN_BIN q gov proposals --output json | jq

# Query the values set as liquidstaking parameters and liquid-validators
$CHAIN_BIN q lspersistence params --output json | jq
$CHAIN_BIN q lspersistence liquid-validators --output json | jq

