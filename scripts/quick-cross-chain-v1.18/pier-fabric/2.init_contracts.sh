#!/usr/bin/env bash
set -e
#set -x
CURRENT_PATH=$(pwd)
PROJECT_PATH=$(dirname "${CURRENT_PATH}")
source "$PROJECT_PATH"/x.sh

channelID="mychannel"
broker_address="broker"
transfer_address="transfer"
data_swapper_address="data_swapper"

function init_broker() {
  print_blue "init broker contract"
  goduck fabric contract invoke --config-path "$CURRENT_PATH"/fabric/config.yaml $broker_address initialize 1356,$Fabric_ChainID
  sleep 2
}

function register_transfer_contract() {
  print_blue "register transfer contract to broker"
  goduck fabric contract invoke --config-path "$CURRENT_PATH"/fabric/config.yaml $transfer_address register 
  sleep 1
  print_blue "audit transfer contract"
  goduck fabric contract invoke --config-path "$CURRENT_PATH"/fabric/config.yaml $broker_address audit $channelID,$transfer_address,1
  sleep 1
}

function register_data_swapper_contract() {
  print_blue "register data_swapper contract to broker"
  goduck fabric contract invoke --config-path "$CURRENT_PATH"/fabric/config.yaml $data_swapper_address register 
  sleep 1
  print_blue "audit data_swapper contract"
  goduck fabric contract invoke --config-path "$CURRENT_PATH"/fabric/config.yaml $broker_address audit $channelID,$data_swapper_address,1
  sleep 1
}

function init_transfer() {
  print_blue "set default account (Alice) balance"
  goduck fabric contract invoke --config-path "$CURRENT_PATH"/fabric/config.yaml $transfer_address setBalance Alice,10000
  sleep 1
  print_blue "get default account (Alice) balance after set"
  goduck fabric contract invoke --config-path "$CURRENT_PATH"/fabric/config.yaml $transfer_address getBalance Alice
  sleep 1
}

function init_data_swapper() {
  print_blue "set default key (fabric)"
  goduck fabric contract invoke --config-path "$CURRENT_PATH"/fabric/config.yaml $data_swapper_address set fabric,data-test
  sleep 1
  print_blue "get default key (fabric) after set"
  goduck fabric contract invoke --config-path "$CURRENT_PATH"/fabric/config.yaml $data_swapper_address get fabric
}

init_broker
register_transfer_contract
register_data_swapper_contract
init_transfer
init_data_swapper