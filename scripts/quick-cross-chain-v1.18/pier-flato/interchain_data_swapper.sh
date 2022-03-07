#!/usr/bin/env bash
set -e
#set -x
CURRENT_PATH=$(pwd)
PROJECT_PATH=$(dirname "${CURRENT_PATH}")
source "$PROJECT_PATH"/x.sh

dst_appchain=${1-'Fabric_ChainID'}
dst_service=${2-'mychannel&data_swapper'}
key=${3-'10'}

data_swapper_address=$(cat "$CURRENT_PATH"/flato/flato.toml|grep data_swapper.abi | awk {'print $1'})
print_green "data_swapper contract address: $data_swapper_address"

print_blue "get data from other appchain"
goduck hpc invoke --config-path "$CURRENT_PATH"/flato --abi-path "$CURRENT_PATH"/flato/data_swapper.abi $data_swapper_address get 1356:"$dst_appchain":"$dst_service",$key
print_blue "wait 5s, please check pier.log"
sleep 5
