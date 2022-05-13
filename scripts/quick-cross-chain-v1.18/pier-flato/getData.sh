#!/usr/bin/env bash
set -e
source ../../x.sh
source ../config.sh
CURRENT_PATH=$(pwd)

data_swapper_address=$(grep data_swapper.abi <"$CURRENT_PATH"/flato/flato.toml | awk '{print $1}')
print_green "data_swapper contract address: $data_swapper_address"

print_blue "set default key (ether)"
goduck hpc invoke --config-path "$CURRENT_PATH"/flato --abi-path "$CURRENT_PATH"/flato/data_swapper.abi "$data_swapper_address" getData ether
