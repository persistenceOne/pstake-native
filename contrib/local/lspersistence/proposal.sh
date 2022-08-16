#!/bin/sh

export BINARY=pstaked

# Submit a param-change governance proposal
$BINARY tx gov submit-proposal param-change ~/proposal_whitelist_validator.json \
--chain-id localnet \
--from user1 \
--keyring-backend test \
--broadcast-mode block \
--yes \
--output json | jq

# Query the proposal to check the status PROPOSAL_STATUS_VOTING_PERIOD
$BINARY q gov proposals --output json | jq

# Vote
$BINARY tx gov vote 1 yes \
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
$BINARY q gov proposals --output json | jq

# Query the values set as liquidstaking parameters and liquid-validators
$BINARY q lspersistence params --output json | jq
$BINARY q lspersistence liquid-validators --output json | jq

