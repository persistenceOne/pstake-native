#!/bin/bash
interchain-queries keys restore --chain $(jq -r ".chains[1].name" /configs/keys.json) --home /icq cosKey0