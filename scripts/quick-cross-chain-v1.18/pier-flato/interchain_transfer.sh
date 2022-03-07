#!/usr/bin/env bash
set -e
#set -x
CURRENT_PATH=$(pwd)
PROJECT_PATH=$(dirname "${CURRENT_PATH}")
source "$PROJECT_PATH"/x.sh

transfer_address=$(cat "$CURRENT_PATH"/flato/flato.toml|grep transfer.abi | awk {'print $1'})
print_green "transfer contract address: $transfer_address"

#dst_appchain=$1
#dst_service=$2
#amount=$3

print_blue "transfer from flato to other appchain"
goduck hpc invoke --config-path flato --abi-path flato/transfer.abi $transfer_address transfer 1356:"$1":"$2",Alice,Alice,"$3"
print_blue "wait 5s, please check pier.log"
sleep 5