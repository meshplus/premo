#!/usr/bin/env bash
set -e
source ../../x.sh
source ../config.sh
CURRENT_PATH=$(pwd)

dst_appchain=$FABRIC_ID
dst_service="mychannel&data_swapper"
key="10"

data_swapper_address=$(grep data_swapper.abi <"$CURRENT_PATH"/flato/flato.toml | awk '{print $1}')
print_green "data_swapper contract address: $data_swapper_address"

print_blue "get data from other appchain"
goduck hpc invoke --config-path "$CURRENT_PATH"/flato --abi-path "$CURRENT_PATH"/flato/data_swapper.abi "$data_swapper_address" get 1356:"$dst_appchain":"$dst_service",$key
print_blue "wait 5s, please check pier.log"
sleep 5
