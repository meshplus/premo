#!/usr/bin/env bash
set -e
#set -x
CURRENT_PATH=$(pwd)
PROJECT_PATH=$(dirname "${CURRENT_PATH}")
source "$PROJECT_PATH"/x.sh

data_swapper_address=$(cat "$CURRENT_PATH"/flato/flato.toml|grep data_swapper.abi | awk {'print $1'})
print_green "data_swapper contract address: $data_swapper_address"

print_blue "set default key (ether)"
goduck hpc invoke --config-path "$CURRENT_PATH"/flato --abi-path "$CURRENT_PATH"/flato/data_swapper.abi $data_swapper_address set ether,data-test