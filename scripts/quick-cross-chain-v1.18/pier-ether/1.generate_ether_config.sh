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
    cp -r "$Ether_Project_Path"/example "$CURRENT_PATH"/contracts
  fi
  if [ ! -f "$CURRENT_PATH"/plugins/eth-client ]; then
    if [ ! -f "$Ether_Project_Path"/build/eth-client ]; then
      print_red "please make install eth-clent first!!"
      exit 2
    fi
    print_blue "===> copy plugin"
    cp -r "$Ether_Project_Path"/build/eth-client "$CURRENT_PATH"/plugins/
  fi
}

function deploy_contracts() {
  print_blue "Deploy broker contract"
  goduck ether contract deploy --address http://172.16.30.84:8545 --key-path "$CURRENT_PATH"/ethereum/account.key --psd-path "$CURRENT_PATH"/ethereum/password --code-path "$CURRENT_PATH"/contracts/broker.sol 1356-"$Ether_ChainID"-["0xc7F999b83Af6DF9e67d0a37Ee7e900bF38b3D013","0x79a1215469FaB6f9c63c1816b45183AD3624bE34","0x97c8B516D19edBf575D72a172Af7F418BE498C37","0xc0Ff2e0b3189132D815b8eb325bE17285AC898f8"]-2-["0xc60ba75739b3492189d80c71ad0aebc0c57695ff"]-1 > "$CURRENT_PATH"/broker.abi
  broker_address=$(cat "$CURRENT_PATH"/broker.abi | grep Deployed | grep -o '0x.\{40\}')
  print_green "broker contract address: $broker_address"
  if [ -z "$broker_address" ]; then
	  echo "broker_address id is empty"
	  exit 2
  fi
  x_sed '1,4d' broker.abi
  x_sed "s/$(cat "$CURRENT_PATH"/ethereum/ethereum.toml|grep contract_address | awk {'print $3'})/\"$broker_address\"/g" "$CURRENT_PATH"/ethereum/ethereum.toml
  sleep 1

  print_blue "Deploy transfer contract"
  goduck ether contract deploy --address http://172.16.30.84:8545 --key-path "$CURRENT_PATH"/ethereum/account.key --psd-path "$CURRENT_PATH"/ethereum/password --code-path "$CURRENT_PATH"/contracts/transfer.sol $broker_address > "$CURRENT_PATH"/transfer.abi
  transfer_address=$(cat "$CURRENT_PATH"/transfer.abi | grep Deployed | grep -o '0x.\{40\}')
  print_green "transfer contract address: $transfer_address" 
  if [ -z "$transfer_address" ]; then
	  echo "transfer_address id is empty"
	  exit 2
  fi 
  echo $transfer_address > "$CURRENT_PATH"/transfer_address.info
  x_sed '1,4d' "$CURRENT_PATH"/transfer.abi
  sleep 1
  
  print_blue "Deploy data_swapper contract"
  goduck ether contract deploy --address http://172.16.30.84:8545 --key-path "$CURRENT_PATH"/ethereum/account.key --psd-path "$CURRENT_PATH"/ethereum/password --code-path "$CURRENT_PATH"/contracts/data_swapper.sol $broker_address > "$CURRENT_PATH"/data_swapper.abi
  data_swapper_address=$(cat "$CURRENT_PATH"/data_swapper.abi|grep Deployed| grep -o '0x.\{40\}')
  print_green "data_swapper contract address: $data_swapper_address"
  if [ -z "$data_swapper_address" ]; then
	  echo "data_swapper_address id is empty"
	  exit 2
  fi
  echo $data_swapper_address > "$CURRENT_PATH"/data_swapper_address.info
  x_sed '1,4d' "$CURRENT_PATH"/data_swapper.abi
  sleep 1
}

function modify_pier_config() {
  print_blue "modify pier.toml"
  x_sed "s/$(cat "$CURRENT_PATH"/pier.toml |grep "id.*appchain"|awk -F '\"' {'print $2'})/$Ether_ChainID/g" "$CURRENT_PATH"/pier.toml
  x_sed "s/$(cat "$CURRENT_PATH"/pier.toml |grep plugin|awk -F '\"' {'print $2'})/eth-client/g" "$CURRENT_PATH"/pier.toml
  x_sed "s/$(cat "$CURRENT_PATH"/pier.toml |grep config|awk -F '\"' {'print $2'})/ethereum/g" "$CURRENT_PATH"/pier.toml

  x_sed "s/$(cat "$CURRENT_PATH"/pier.toml |grep http|awk {'print $3'})/34544/g" "$CURRENT_PATH"/pier.toml
  x_sed "s/$(cat "$CURRENT_PATH"/pier.toml |grep pprof|awk {'print $3'})/34555/g" "$CURRENT_PATH"/pier.toml

  print_blue "modify ethereum.toml"
  x_sed "s/$(cat "$CURRENT_PATH"/ethereum/ethereum.toml | grep name |awk -F '\"' {'print $2'})/$Ether_ChainID/g" "$CURRENT_PATH"/ethereum/ethereum.toml
}

prepare
deploy_contracts
modify_pier_config