#!/usr/bin/env bash
set -e
source ../../x.sh
source ../config.sh
CURRENT_PATH=$(pwd)
PROJECT_PATH=$(dirname "$CURRENT_PATH")

pier_ether_id=$(pier --repo "$CURRENT_PATH" id)
print_green "pier-ether id: $pier_ether_id"

broker_address=$(grep contract_address <ethereum/ethereum.toml | awk -F '\"' '{print $2}')
print_green "broker contract: $broker_address"

transfer_address=$(awk '{print $1}' <"$CURRENT_PATH/transfer_address.info")
print_green "transfer contract: $transfer_address"

data_swapper_address=$(awk '{print $1}' <"$CURRENT_PATH/data_swapper_address.info")
print_green "data_swapper contract: $data_swapper_address"

function prepare() {
  if [ "$BITXHUB_TYPE" == solo ]; then
    bitxhub_node1_config="$PROJECT_PATH/bitxhub/repo_solo"
  else
    bitxhub_node1_config="$PROJECT_PATH/bitxhub/repo_raft/node1"
  fi

  print_blue "get balance from bitxhub"
  bitxhub --repo "$bitxhub_node1_config" client tx send --key "$bitxhub_node1_config"/key.json --to "$pier_ether_id" --amount 10000000000000000
}

function register_appchain() {
  print_blue "register appchain to bitxhub"
  pier --repo "$CURRENT_PATH" appchain register --appchain-id "$ETHER_ID" --name "$ETHER_ID" --type=ether --trustroot="$CURRENT_PATH"/ethereum/ether.validators --broker "$broker_address" --desc="ether appchain for test" --master-rule "$ETHER_RULE_ADDRESS" --rule-url http://github.com --admin 0xc60ba75739b3492189d80c71ad0aebc0c57695ff --reason "test" >"$CURRENT_PATH/register_appchain.info"

  proposalId=$(grep successfully <"$CURRENT_PATH"/register_appchain.info | awk '{print $7}')
  if [ -z "$proposalId" ]; then
    print_red "proposal id is empty"
    exit 2
  fi
  print_green "register appchain successfully, proposalId is: $proposalId"
  sleep 1
}

function register_transfer_service() {
  print_blue "register transfer service to bitxhub"
  pier --repo "$CURRENT_PATH" appchain service register --appchain-id "$ETHER_ID" --service-id "$transfer_address" --name ether-transfer-test --intro "test" --type CallContract --details "test" --reason "test"
  sleep 1
}

function register_data_swapper_service() {
  print_blue "register data_swapper service to bitxhub"
  pier --repo "$CURRENT_PATH" appchain service register --appchain-id "$ETHER_ID" --service-id "$data_swapper_address" --name ether-data-test --intro "test" --type CallContract --details "test" --reason "test"
  sleep 1
}

prepare
sleep 2
register_appchain
register_transfer_service
register_data_swapper_service
