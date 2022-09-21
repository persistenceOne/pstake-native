export VAL_INDEX=${HOSTNAME##*-}
export VAL_NAME=$(jq -r ".validators[$VAL_INDEX].name" /configs/validators.json)
echo "Validator Index: $VAL_INDEX, Key name: $VAL_NAME"

# Wait for the node to be synced
export max_tries=10
while [[ $(gaiad status --home $GAIA_DIR 2>&1 | jq ".SyncInfo.catching_up") == true ]]
do
    if [[ max_tries -lt 0 ]]; then echo "Not able to sync with genesis node"; exit 1; fi
    echo "Still syncing... Sleeping for 15 secs. Tries left $max_tries"
    ((max_tries--))
    sleep 30
done

sleep 10

echo "Keys list"
gaiad keys list --home $GAIA_DIR --keyring-backend test

export VAL_ADDRESS=$(gaiad keys show $VAL_NAME --home $GAIA_DIR --bech val --keyring-backend test --output json | jq -r ".address")
export STATUS=$(gaiad q staking validator $VAL_ADDRESS --node http://gaia-genesis.dev-native.svc.cluster.local:26657 --output json | jq -r ".status")


echo "STATUS:" $STATUS
if [ "$STATUS" != "BOND_STATUS_BONDED" ]; then
    # Run create validator tx command
    echo "Running txn for create-validator"
    export VALIDATOR_PUBKEY=$(gaiad tendermint show-validator --home $GAIA_DIR)


    echo "VALIDATOR PUBKEY: " $VALIDATOR_PUBKEY
    gaiad tx staking create-validator \
        --home $GAIA_DIR \
        --pubkey=$VALIDATOR_PUBKEY \
        --moniker $VAL_NAME \
        --amount 80000000000000000uatom \
        --keyring-backend="test" \
        --chain-id $CHAIN_ID \
        --from $VAL_NAME \
        --commission-rate="0.10" \
        --commission-max-rate="0.20" \
        --commission-max-change-rate="0.01" \
        --min-self-delegation="1000000" \
        --gas="auto"\
        --gas-adjustment 1.5 \
        --yes 2>&1 | tee /validator.log
fi

exit 0