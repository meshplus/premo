#!/usr/bin/env bash
set -e
#set -x
CURRENT_PATH=$(pwd)
PROJECT_PATH=$(dirname "${CURRENT_PATH}")
source "$PROJECT_PATH"/x.sh

transfer_address=$(cat "$CURRENT_PATH"/flato/flato.toml|grep transfer.abi | awk {'print $1'})
print_green "transfer contract address: $transfer_address"

print_blue "set default account (Alice) balance 10000"
goduck hpc invoke --config-path "$CURRENT_PATH"/flato --abi-path "$CURRENT_PATH"/flato/transfer.abi $transfer_address setBalance Alice,10000
