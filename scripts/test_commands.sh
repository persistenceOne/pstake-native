#!/bin/sh

pstaked tx cosmos set-orchestrator-address persistencevaloper183g695ap32wnds5k9xwd3yq997dqxudfz524f3 persistence12v9prjx8m5fdalryqd0t4mgwe20637ltek5m0h --chain-id test --from test --keyring-backend=test -y -b block
#pstaked tx cosmos set-orchestrator-address cosmosvaloper1hcqg5wj9t42zawqkqucs7la85ffyv08le09ljt cosmos1xruvjju28j0a5ud5325rfdak8f5a04h07afg3f --chain-id test --from test --keyring-backend=test -y -b block
pstaked tx bank send cosmos1hcqg5wj9t42zawqkqucs7la85ffyv08lum327c cosmos1lcck2cxh7dzgkrfk53kysg9ktdrsjj6jv6txhe 10000000stake --chain-id test --from test --keyring-backend=test -y -b block
#pstaked tx bank send cosmos1hcqg5wj9t42zawqkqucs7la85ffyv08lum327c cosmos1xruvjju28j0a5ud5325rfdak8f5a04h07afg3f 10000000stake --chain-id test --from test --keyring-backend=test -y -b block
pstaked tx cosmos incoming cosmos10khgeppewe4rgfrcy809r9h00aquwxxxd6um38 cosmos1lcck2cxh7dzgkrfk53kysg9ktdrsjj6jv6txhe 10000000uatom cosmoshub-4 AE9ADDF593D45DDB09C8371F534AA773EB8CF288F63B09C160110338D362177B 100000 --chain-id test --from test1 --keyring-backend=test -y -b block
#pstaked tx cosmos incoming cosmos10khgeppewe4rgfrcy809r9h00aquwxxxd6um38 cosmos1xruvjju28j0a5ud5325rfdak8f5a04h07afg3f 10000000uatom cosmoshub-4 AE9ADDF593D45DDB09C8371F534AA773EB8CF288F63B09C160110338D362177B 100000 --chain-id test --from test3 --keyring-backend=test -y -b block
pstaked tx cosmos incoming cosmos1xruvjju28j0a5ud5325rfdak8f5a04h07afg3f cosmos1lcck2cxh7dzgkrfk53kysg9ktdrsjj6jv6txhe 10000000uatom cosmoshub-4 AE9ADDF593D45DDB09C8371F534AA773EB8CF288F63B09C160110338D362177C 100000 --chain-id test --from test1 --keyring-backend=test -y -b block
