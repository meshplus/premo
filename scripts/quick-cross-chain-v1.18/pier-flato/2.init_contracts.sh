#!/usr/bin/env bash
set -e
#set -x
CURRENT_PATH=$(pwd)
PROJECT_PATH=$(dirname "${CURRENT_PATH}")
source "$PROJECT_PATH"/x.sh

broker_address=$(cat "$CURRENT_PATH"/flato/flato.toml | grep contract_address | awk {'print $3'} | awk -F '\"' {'print $2'})
print_green "broker contract address: $broker_address"
transfer_address=$(cat "$CURRENT_PATH"/flato/flato.toml | grep transfer.abi | awk {'print $1'})
print_green "transfer contract address: $transfer_address"
data_swapper_address=$(cat "$CURRENT_PATH"/flato/flato.toml | grep data_swapper.abi | awk {'print $1'})
print_green "data_swapper contract address: $data_swapper_address"

function init_broker() {
  print_blue "init broker contract"
  goduck hpc invoke --config-path "$CURRENT_PATH"/flato --abi-path "$CURRENT_PATH"/flato/broker.abi $broker_address initialize
  sleep 2
}

function register_transfer_contract() {
  print_blue "register and audit transfer contract"
  goduck hpc invoke --config-path "$CURRENT_PATH"/flato --abi-path "$CURRENT_PATH"/flato/broker.abi $broker_address register $transfer_address
  sleep 1
  goduck hpc invoke --config-path "$CURRENT_PATH"/flato --abi-path "$CURRENT_PATH"/flato/broker.abi $broker_address audit $transfer_address,1
  sleep 1
}

function register_data_swapper_contract() {
  print_blue "register and audit data_swapper contract"
  goduck hpc invoke --config-path "$CURRENT_PATH"/flato --abi-path "$CURRENT_PATH"/flato/broker.abi $broker_address register $data_swapper_address
  sleep 1
  goduck hpc invoke --config-path "$CURRENT_PATH"/flato --abi-path "$CURRENT_PATH"/flato/broker.abi $broker_address audit $data_swapper_address,1
  sleep 1
}

function init_transfer() {
  print_blue "set default account (Alice) balance"
  goduck hpc invoke --config-path "$CURRENT_PATH"/flato --abi-path "$CURRENT_PATH"/flato/transfer.abi $transfer_address setBalance Alice,10000
  sleep 1
  print_blue "get default account (Alice) balance after set"
  goduck hpc invoke --config-path "$CURRENT_PATH"/flato --abi-path "$CURRENT_PATH"/flato/transfer.abi $transfer_address getBalance Alice
  sleep 1
}

function init_data_swapper() {
  print_blue "set default key (ether)"
  goduck hpc invoke --config-path "$CURRENT_PATH"/flato --abi-path "$CURRENT_PATH"/flato/data_swapper.abi $data_swapper_address set ether,data-test
  sleep 1
  print_blue "get default key (ether) after set"
  goduck hpc invoke --config-path "$CURRENT_PATH"/flato --abi-path "$CURRENT_PATH"/flato/data_swapper.abi $data_swapper_address getData ether
}

init_broker
register_transfer_contract
register_data_swapper_contract
init_transfer
init_data_swapper
