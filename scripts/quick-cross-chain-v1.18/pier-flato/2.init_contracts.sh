#!/usr/bin/env bash
set -e
source ../../x.sh
source ../config.sh
CURRENT_PATH=$(pwd)

broker_address=$(grep contract_address <"$CURRENT_PATH"/flato/flato.toml | awk '{print $3}' | awk -F '\"' '{print $2}')
print_green "broker contract address: $broker_address"

transfer_address=$(awk '{print $1}' <"$CURRENT_PATH"/transfer_address.info)
print_green "transfer contract: $transfer_address"

data_swapper_address=$(awk '{print $1}' <"$CURRENT_PATH"/data_swapper_address.info)
print_green "data_swapper contract: $data_swapper_address"


function register_transfer_contract() {
  print_blue "register and audit transfer contract"
  goduck hpc invoke --config-path "$CURRENT_PATH"/flato --abi-path "$CURRENT_PATH"/flato/broker.abi "$broker_address" audit "$transfer_address"^1
  sleep 1
}

function register_data_swapper_contract() {
  print_blue "register and audit data_swapper contract"
  goduck hpc invoke --config-path "$CURRENT_PATH"/flato --abi-path "$CURRENT_PATH"/flato/broker.abi "$broker_address" audit "$data_swapper_address"^1
  sleep 1
}

function init_transfer() {
  print_blue "set default account (Alice) balance"
  goduck hpc invoke --config-path "$CURRENT_PATH"/flato --abi-path "$CURRENT_PATH"/flato/transfer.abi "$transfer_address" setBalance Alice^10000
  sleep 1
  print_blue "get default account (Alice) balance after set"
  goduck hpc invoke --config-path "$CURRENT_PATH"/flato --abi-path "$CURRENT_PATH"/flato/transfer.abi "$transfer_address" getBalance Alice
  sleep 1
}

function init_data_swapper() {
  print_blue "set default key (ether)"
  goduck hpc invoke --config-path "$CURRENT_PATH"/flato --abi-path "$CURRENT_PATH"/flato/data_swapper.abi "$data_swapper_address" set ether^data-test
  sleep 1
  print_blue "get default key (ether) after set"
  goduck hpc invoke --config-path "$CURRENT_PATH"/flato --abi-path "$CURRENT_PATH"/flato/data_swapper.abi "$data_swapper_address" getData ether
}

register_transfer_contract
register_data_swapper_contract
init_transfer
init_data_swapper
