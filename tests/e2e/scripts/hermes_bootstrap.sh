#!/bin/bash

set -ex

# initialize Hermes relayer configuration
mkdir -p /root/.hermes/
touch /root/.hermes/config.toml

# setup Hermes relayer configuration
tee /root/.hermes/config.toml <<EOF
[global]
log_level = 'info'

[mode]

[mode.clients]
enabled = true
refresh = true
misbehaviour = true

[mode.connections]
enabled = false

[mode.channels]
enabled = false

[mode.packets]
enabled = true
clear_interval = 100
clear_on_start = true
tx_confirmation = true

[rest]
enabled = true
host = '0.0.0.0'
port = 3031

[telemetry]
enabled = true
host = '127.0.0.1'
port = 3001

[[chains]]
id = '$PSTAKE_A_E2E_CHAIN_ID'
rpc_addr = 'http://$PSTAKE_A_E2E_VAL_HOST:26657'
grpc_addr = 'http://$PSTAKE_A_E2E_VAL_HOST:9090'
websocket_addr = 'ws://$PSTAKE_A_E2E_VAL_HOST:26657/websocket'
rpc_timeout = '10s'
account_prefix = 'persistence'
key_name = 'val01-pstake-a'
store_prefix = 'ibc'
default_gas = 500000
max_gas = 900000
gas_price = { price = 0.001, denom = 'uxprt' }
gas_multiplier = 2
clock_drift = '1m' # to accomdate docker containers
trusting_period = '14days'
trust_threshold = { numerator = '1', denominator = '3' }

[[chains]]
id = '$PSTAKE_B_E2E_CHAIN_ID'
rpc_addr = 'http://$PSTAKE_B_E2E_VAL_HOST:26657'
grpc_addr = 'http://$PSTAKE_B_E2E_VAL_HOST:9090'
websocket_addr = 'ws://$PSTAKE_B_E2E_VAL_HOST:26657/websocket'
rpc_timeout = '10s'
account_prefix = 'persistence'
key_name = 'val01-pstake-b'
store_prefix = 'ibc'
default_gas = 500000
max_gas = 900000
gas_price = { price = 0.001, denom = 'uxprt' }
gas_multiplier = 2
clock_drift = '1m' # to accomdate docker containers
trusting_period = '14days'
trust_threshold = { numerator = '1', denominator = '3' }
EOF

echo ${PSTAKE_B_E2E_VAL_MNEMONIC} >> /root/.hermes/pstake_b_e2e_chain_mnemonic.txt
echo ${PSTAKE_A_E2E_VAL_MNEMONIC} >> /root/.hermes/pstake_a_e2e_chain_mnemonic.txt

# import keys
hermes keys add --chain ${PSTAKE_B_E2E_CHAIN_ID} --key-name "val01-pstake-b" --mnemonic-file /root/.hermes/pstake_b_e2e_chain_mnemonic.txt
hermes keys add --chain ${PSTAKE_A_E2E_CHAIN_ID} --key-name "val01-pstake-a" --mnemonic-file /root/.hermes/pstake_a_e2e_chain_mnemonic.txt

# start Hermes relayer
hermes start
