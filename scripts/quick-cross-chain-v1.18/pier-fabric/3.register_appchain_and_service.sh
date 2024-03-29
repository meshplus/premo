#!/usr/bin/env bash
set -e
source ../../x.sh
source ../config.sh
CURRENT_PATH=$(pwd)
PROJECT_PATH=$(dirname "$CURRENT_PATH")

channelID="mychannel"
broker_address="broker"
broker_version="v1.0.0"
transfer_address="transfer"
data_swapper_address="data_swapper"

pier_fabric_id=$(pier --repo "$CURRENT_PATH" id)

function prepare() {
  if [ "$BITXHUB_TYPE" == solo ]; then
    bitxhub_node1_config="$PROJECT_PATH/bitxhub/repo_solo"
  else
    bitxhub_node1_config="$PROJECT_PATH/bitxhub/repo_raft/node1"
  fi

  print_blue "get balance from bitxhub"
  bitxhub --repo "$bitxhub_node1_config" client tx send --key "$bitxhub_node1_config"/key.json --to "$pier_fabric_id" --amount 10000000000000000
  sleep 1
}

function register_appchain() {
  print_blue "register appchain to bitxhub"
  pier --repo "$CURRENT_PATH" appchain register --appchain-id "$FABRIC_ID" --name "$FABRIC_ID" --type=fabric --trustroot="$CURRENT_PATH"/fabric/fabric.validators --broker-cid $channelID --broker-ccid $broker_address --broker-v $broker_version --desc="fabric appchain for test" --master-rule "$FABRIC_RULE_ADDRESS" --rule-url http://github.com --reason "test" >"$CURRENT_PATH"/register_appchain.info

  proposalId=$(grep successfully < "$CURRENT_PATH"/register_appchain.info | awk '{print $7}')
  if [ -z "$proposalId" ]; then
    print_red "proposal id is empty"
    exit 2
  fi
  print_green "register appchain successfully, proposalId is: $proposalId"
}

function register_transfer_service() {
  print_blue "register transfer service to bitxhub"
  pier --repo "$CURRENT_PATH" appchain service register --appchain-id "$FABRIC_ID" --service-id "$channelID&$transfer_address" --name fabric-transfer-test --intro "test" --type CallContract --details "test" --reason "test"
  sleep 1
}

function register_data_swapper_service() {
  print_blue "register data_swapper service to bitxhub"
  pier --repo ./ appchain service register --appchain-id "$FABRIC_ID" --service-id "$channelID&$data_swapper_address" --name fabric-data-test --intro "test" --type CallContract --details "test" --reason "test"
  sleep 1
}

prepare
register_appchain
register_transfer_service
register_data_swapper_service
