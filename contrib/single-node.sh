#!/bin/sh

set -o errexit -o nounset

CHAINID=$1
GENACCT=$2

if [ -z "$1" ]; then
  echo "Need to input chain id..."
  exit 1
fi

if [ -z "$2" ]; then
  echo "Need to input genesis account address..."
  exit 1
fi

# Build genesis file incl account for passed address
coins="100000000000000000stake"
pstaked init --chain-id $CHAINID $CHAINID
pstaked keys add validator --keyring-backend="test"
pstaked add-genesis-account $(pstaked keys show validator -a --keyring-backend="test") $coins
pstaked add-genesis-account $GENACCT $coins
pstaked gentx validator 5000000000stake --keyring-backend="test" --chain-id $CHAINID
pstaked collect-gentxs

# Set proper defaults and change ports
sed -i 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:26657"#g' ~/.pstaked/config/config.toml
sed -i 's/timeout_commit = "5s"/timeout_commit = "1s"/g' ~/.gaia/config/config.toml
sed -i 's/timeout_propose = "3s"/timeout_propose = "1s"/g' ~/.gaia/config/config.toml
sed -i 's/index_all_keys = false/index_all_keys = true/g' ~/.pstaked/config/config.toml

# Start the gaia
pstaked start --pruning=nothing
