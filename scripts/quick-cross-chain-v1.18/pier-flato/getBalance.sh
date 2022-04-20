#!/usr/bin/env bash
set -e
source ../../x.sh
source ../config.sh
CURRENT_PATH=$(pwd)

transfer_address=$(grep transfer.abi <"$CURRENT_PATH"/flato/flato.toml | awk '{print $1}')
print_green "transfer contract address: $transfer_address"

print_blue "get default account (Alice) balance"
goduck hpc invoke --config-path "$CURRENT_PATH"/flato --abi-path "$CURRENT_PATH"/flato/transfer.abi "$transfer_address" getBalance Alice
