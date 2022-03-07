#!/usr/bin/env bash
set -e
#set -x

CURRENT_PATH=$(pwd)
PROJECT_PATH=$(dirname "${CURRENT_PATH}")
source "$PROJECT_PATH"/x.sh

function prepare() {
  check_goduck
  if [ ! -d "$CURRENT_PATH"/contracts ]; then
    print_blue "===> copy contracts"
    cp -r "$Flato_Project_Path"/example/contracts "$CURRENT_PATH"/
  fi
  if [ ! -f "$CURRENT_PATH"/plugins/flt-client ]; then
    if [ ! -f "$Flato_Project_Path"/build/flt-client ]; then
      print_red "please make install flt-clent first!!"
      exit 2
    fi
    print_blue "===> copy plugin"
    cp -r "$Flato_Project_Path"/build/flt-client "$CURRENT_PATH"/plugins/
  fi
}

function deploy_contracts() {
  print_blue "Deploy contracts"
  print_blue "1. Deploy broker contract"
  goduck hpc deploy --config-path "$CURRENT_PATH"/flato --code-path "$CURRENT_PATH"/contracts/broker.sol 1356-"$Flato_ChainID"-["0xc7F999b83Af6DF9e67d0a37Ee7e900bF38b3D013","0x79a1215469FaB6f9c63c1816b45183AD3624bE34","0x97c8B516D19edBf575D72a172Af7F418BE498C37","0xc0Ff2e0b3189132D815b8eb325bE17285AC898f8"]-2-["0x000f1a7a08ccc48e5d30f80850cf1cf283aa3abd"]-1 > "$CURRENT_PATH"/broker.abi

  broker_address=$(cat "$CURRENT_PATH"/broker.abi | grep address: | grep -o '0x.\{40\}')
  print_green "broker contract address: $broker_address"
  if [ -z "$broker_address" ]; then
	  print_red "broker_address id is empty"
	  exit 2
  fi
  x_sed '1,2d' "$CURRENT_PATH"/broker.abi
  x_sed "s/$(cat "$CURRENT_PATH"/flato/flato.toml|grep contract_address | awk {'print $3'})/\"$broker_address\"/g" "$CURRENT_PATH"/flato/flato.toml
  sleep 1

  print_blue "2. Deploy transfer contract"
  goduck hpc deploy --config-path "$CURRENT_PATH"/flato --code-path "$CURRENT_PATH"/contracts/transfer.sol $broker_address > "$CURRENT_PATH"/transfer.abi
  transfer_address=$(cat "$CURRENT_PATH"/transfer.abi | grep address: | grep -o '0x.\{40\}')
  print_green "transfer contract address: $transfer_address"
  if [ -z "$transfer_address" ]; then
	  print_red "transfer_address id is empty"
	  exit 2
  fi
  x_sed '1,2d' "$CURRENT_PATH"/transfer.abi
  x_sed '/emitInterchainEvent/d' "$CURRENT_PATH"/transfer.abi
  x_sed "s/$(cat "$CURRENT_PATH"/flato/flato.toml|grep transfer.abi | awk {'print $1'})/$transfer_address/g" "$CURRENT_PATH"/flato/flato.toml
  sleep 1

  print_blue "3. Deploy data_swapper contract"
  goduck hpc deploy --config-path "$CURRENT_PATH"/flato --code-path "$CURRENT_PATH"/contracts/data_swapper.sol $broker_address > "$CURRENT_PATH"/data_swapper.abi
  #data_swapper_address=$(cat data_swapper.abi|grep address:| awk {'print $4'})
  data_swapper_address=$(cat "$CURRENT_PATH"/data_swapper.abi|grep address:| grep -o '0x.\{40\}')
  print_green "data_swapper contract address: $data_swapper_address"
  if [ -z "$data_swapper_address" ]; then
	  print_red "data_swapper_address id is empty"
	  exit 2
  fi
  x_sed '1,2d' "$CURRENT_PATH"/data_swapper.abi
  x_sed '/emitInterchainEvent/d' "$CURRENT_PATH"/data_swapper.abi
  x_sed "s/$(cat "$CURRENT_PATH"/flato/flato.toml|grep data_swapper.abi | awk {'print $1'})/$data_swapper_address/g" "$CURRENT_PATH"/flato/flato.toml
  sleep 1

  mv "$CURRENT_PATH"/*.abi "$CURRENT_PATH"/flato/
}

function register_mq() {
  print_blue "register mq"
  goduck mq register --config-path "$CURRENT_PATH"/flato $broker_address $(git config user.name)-$(date +"%Y%m%d-%s") > "$CURRENT_PATH"/mq.info
  queue_name=$(cat "$CURRENT_PATH"/mq.info | grep queue | awk {'print $3'} | awk -F ',' {'print $1'})
  print_green "queue name: $queue_name"
  exchanger=$(cat "$CURRENT_PATH"/mq.info | grep exchanger| awk {'print $5'})
  print_green "exchanger info: $exchanger"
  x_sed "s/$(cat "$CURRENT_PATH"/flato/flato.toml | grep queue_name | awk {'print $3'})/\"$queue_name\"/g" "$CURRENT_PATH"/flato/flato.toml
  x_sed "s/$(cat "$CURRENT_PATH"/flato/flato.toml | grep exchange | awk {'print $3'})/\"$exchanger\"/g" "$CURRENT_PATH"/flato/flato.toml
}

function modify_pier_config() {
  print_blue "modify pier.toml"
  x_sed "s/$(cat "$CURRENT_PATH"/pier.toml |grep "id.*appchain"|awk -F '\"' {'print $2'})/$Flato_ChainID/g" "$CURRENT_PATH"/pier.toml
  x_sed "s/$(cat "$CURRENT_PATH"/pier.toml |grep plugin|awk -F '\"' {'print $2'})/flt-client/g" "$CURRENT_PATH"/pier.toml
  x_sed "s/$(cat "$CURRENT_PATH"/pier.toml |grep config|awk -F '\"' {'print $2'})/flato/g" "$CURRENT_PATH"/pier.toml

  x_sed "s/$(cat "$CURRENT_PATH"/pier.toml |grep http|awk {'print $3'})/14544/g" "$CURRENT_PATH"/pier.toml
  x_sed "s/$(cat "$CURRENT_PATH"/pier.toml |grep pprof|awk {'print $3'})/14555/g" "$CURRENT_PATH"/pier.toml
}

prepare
deploy_contracts
register_mq
modify_pier_config
