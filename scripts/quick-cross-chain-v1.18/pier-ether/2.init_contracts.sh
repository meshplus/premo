#!/usr/bin/env bash
set -e
#set -x
CURRENT_PATH=$(pwd)
PROJECT_PATH=$(dirname "${CURRENT_PATH}")
source "$PROJECT_PATH"/x.sh

broker_address=$(cat "$CURRENT_PATH"/ethereum/ethereum.toml | grep contract_address | awk -F '\"' {'print $2'})
print_green "broker contract: $broker_address"

transfer_address=$(cat "$CURRENT_PATH"/transfer_address.info | awk {'print $1'})
print_green "transfer contract: $transfer_address"

data_swapper_address=$(cat "$CURRENT_PATH"/data_swapper_address.info | awk {'print $1'})
print_green "data_swapper contract: $data_swapper_address"

function register_transfer_contract() {
  print_blue "audit transfer contract"
  goduck ether contract invoke --address http://172.16.30.84:8545 --key-path "$CURRENT_PATH"/ethereum/account.key --psd-path "$CURRENT_PATH"/ethereum/password --abi-path "$CURRENT_PATH"/broker.abi $broker_address audit $transfer_address,1
  sleep 1
}

function register_data_swapper_contract() {
  print_blue "audit data_swapper contract"
  goduck ether contract invoke --address http://172.16.30.84:8545 --key-path ethereum/account.key --psd-path ethereum/password --abi-path broker.abi $broker_address audit $data_swapper_address,1
  sleep 1
}

function init_transfer() {
  print_blue "set default account (Alice) balance"
  goduck ether contract invoke --address http://172.16.30.84:8545 --key-path ethereum/account.key --psd-path ethereum/password --abi-path transfer.abi $transfer_address setBalance Alice,10000
  sleep 1
  print_blue "get default account (Alice) balance after set"
  goduck ether contract invoke --address http://172.16.30.84:8545 --key-path ethereum/account.key --psd-path ethereum/password --abi-path transfer.abi $transfer_address getBalance Alice
  sleep 1
}

function init_data_swapper() {
  print_blue "set default key (ether)"
  goduck ether contract invoke --address http://172.16.30.84:8545 --key-path ethereum/account.key --psd-path ethereum/password --abi-path data_swapper.abi $data_swapper_address set ether,data-test
  sleep 1
  print_blue "get default key (ether) after set"
  goduck ether contract invoke --address http://172.16.30.84:8545 --key-path ethereum/account.key --psd-path ethereum/password --abi-path data_swapper.abi $data_swapper_address getData ether
}

register_transfer_contract
register_data_swapper_contract
init_transfer
init_data_swapper
