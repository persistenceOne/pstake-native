#!/bin/bash
interchain-queries keys restore --chain $(jq -r ".chains[0].name" /configs/keys.json) --home /icq perKey0
