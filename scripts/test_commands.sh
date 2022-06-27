#!/bin/sh

pstaked config chain-id test
pstaked config keyring-backend test
pstaked config

pstaked tx gov submit-proposal change-cosmos-validator-weights scripts/data/cosmos_validator_set_proposal.json --from test --gas auto -y -b block
pstaked tx gov vote 1 yes --from test -y -b block

pstaked tx gov submit-proposal change-oracle-validator-weights scripts/data/oracle_validator_set_proposal.json --from test --gas auto -y -b block
pstaked tx gov vote 2 yes --from test  -y -b block

pstaked tx bank send persistence183g695ap32wnds5k9xwd3yq997dqxudfts2gqg persistence12v9prjx8m5fdalryqd0t4mgwe20637ltek5m0h 10000000stake --chain-id test --from test --keyring-backend=test -y
pstaked tx bank send persistence12v9prjx8m5fdalryqd0t4mgwe20637ltek5m0h persistence183g695ap32wnds5k9xwd3yq997dqxudfts2gqg 10stake --chain-id test --from test1 --keyring-backend=test -y

#pstaked tx bank send persistence183g695ap32wnds5k9xwd3yq997dqxudfts2gqg cosmos1xruvjju28j0a5ud5325rfdak8f5a04h07afg3f 10000000stake --chain-id test --from test --keyring-backend=test -y -b block
pstaked tx cosmos set-orchestrator-address persistencevaloper183g695ap32wnds5k9xwd3yq997dqxudfz524f3 persistence12v9prjx8m5fdalryqd0t4mgwe20637ltek5m0h --from test -y -b block

pstaked tx gov submit-proposal enable-module scripts/data/module_enable_proposal.json --from test -y -b block
pstaked tx gov vote 3 yes --from test -y -b block

#pstaked tx cosmos set-orchestrator-address cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt cosmos1xruvjju28j0a5ud5325rfdak8f5a04h07afg3f --chain-id test --from test --keyring-backend=test -y -b block
pstaked tx cosmos incoming persistence1g5lz0gq98y8tav477dltxgpdft0wr9rmqt7mvu persistence12v9prjx8m5fdalryqd0t4mgwe20637ltek5m0h 10000000uatom cosmoshub-4 AE9ADDF593D45DDB09C8371F534AA773EB8CF288F63B09C160110338D362177B 100000 --from test1 -y -b block
pstaked q bank balances persistence1g5lz0gq98y8tav477dltxgpdft0wr9rmqt7mvu

pstaked tx cosmos rewards-claimed persistence12v9prjx8m5fdalryqd0t4mgwe20637ltek5m0h 500uatom cosmoshub-4 100000 --from test1
