#!/bin/sh

pstaked config chain-id test
pstaked config keyring-backend test
pstaked config
pstaked tx gov submit-proposal change-cosmos-validator-weights scripts/data/cosmos_validator_set_proposal.json --from test --gas auto -y -b block
pstaked tx gov vote 1 yes --from test -y -b block
pstaked tx gov submit-proposal change-oracle-validator-weights scripts/data/oracle_validator_set_proposal.json --from test --gas auto -y -b block
pstaked tx gov vote 2 yes --from test  -y -b block
pstaked tx bank send persistence1hcqg5wj9t42zawqkqucs7la85ffyv08ljhhesu persistence1lcck2cxh7dzgkrfk53kysg9ktdrsjj6jzkd4ea 10000000stake --chain-id test --from test --keyring-backend=test -y -b block
pstaked tx bank send persistence1lcck2cxh7dzgkrfk53kysg9ktdrsjj6jzkd4ea persistence1hcqg5wj9t42zawqkqucs7la85ffyv08ljhhesu 10stake --chain-id test --from test1 --keyring-backend=test -y -b block
sleep 25
pstaked q gov proposals
pstaked tx cosmos set-orchestrator-address persistencevaloper1hcqg5wj9t42zawqkqucs7la85ffyv08lmnhye9 persistence1lcck2cxh7dzgkrfk53kysg9ktdrsjj6jzkd4ea --from test -y -b block --gas 400000

pstaked tx gov submit-proposal enable-module scripts/data/module_enable_proposal.json --from test -y -b block
pstaked tx gov vote 3 yes --from test -y -b block

sleep 25
pstaked tx cosmos incoming persistence10khgeppewe4rgfrcy809r9h00aquwxxxrk6glr persistence1lcck2cxh7dzgkrfk53kysg9ktdrsjj6jzkd4ea 10000000stake cosmoshub-4 AE9ADDF593D45DDB09C8371F534AA773EB8CF288F63B09C160110338D362177B 100000 --from test1 -y -b block --gas 400000
pstaked q bank balances persistence10khgeppewe4rgfrcy809r9h00aquwxxxrk6glr

pstaked tx cosmos rewards-claimed persistence1lcck2cxh7dzgkrfk53kysg9ktdrsjj6jzkd4ea 500stake cosmoshub-4 100000 --from test1
