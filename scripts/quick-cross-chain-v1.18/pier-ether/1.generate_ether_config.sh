#!/usr/bin/env bash
set -e
source ../../x.sh
source ../config.sh
CURRENT_PATH=$(pwd)

function prepare() {
  check_goduck
  print_blue "===> copy contracts"
  rm -rf "$CURRENT_PATH/contracts"
  cp -r "$ETHER_PROJECT_PATH/example" "$CURRENT_PATH/contracts"
  if [ ! -f "$ETHER_PROJECT_PATH/build/eth-client" ]; then
    print_red "please make install eth-client first!!"
    exit 2
  fi
  print_blue "===> copy plugin"
  rm -rf "$CURRENT_PATH/plugins/"
  if [ ! -f "$CURRENT_PATH/plugins"/ ]; then
    mkdir "$CURRENT_PATH/plugins/"
  fi
  cp -r "$ETHER_PROJECT_PATH/build/eth-client" "$CURRENT_PATH/plugins/"
}

function deploy_contracts() {
  print_blue "Deploy broker contract"
  goduck ether contract deploy --address http://172.16.30.84:8545 --key-path "$CURRENT_PATH"/ethereum/account.key --psd-path "$CURRENT_PATH"/ethereum/password --code-path "$CURRENT_PATH"/contracts/broker.sol 1356^"$ETHER_ID"^["0xc7F999b83Af6DF9e67d0a37Ee7e900bF38b3D013","0x79a1215469FaB6f9c63c1816b45183AD3624bE34","0x97c8B516D19edBf575D72a172Af7F418BE498C37","0xc0Ff2e0b3189132D815b8eb325bE17285AC898f8"]^2^["0xc60ba75739b3492189d80c71ad0aebc0c57695ff"]^1 >"$CURRENT_PATH"/broker.abi
  broker_address=$(grep Deployed <"$CURRENT_PATH/broker.abi" | grep -o '0x.\{40\}')
  print_green "broker contract address: $broker_address"
  if [ -z "$broker_address" ]; then
    echo "broker_address id is empty"
    exit 2
  fi
  x_replace "1,4d" broker.abi
  x_replace "s/$(grep contract_address <"$CURRENT_PATH"/ethereum/ethereum.toml | awk '{print $3}')/\"$broker_address\"/g" "$CURRENT_PATH"/ethereum/ethereum.toml
  sleep 1

  print_blue "Deploy transfer contract"
  goduck ether contract deploy --address http://172.16.30.84:8545 --key-path "$CURRENT_PATH"/ethereum/account.key --psd-path "$CURRENT_PATH"/ethereum/password --code-path "$CURRENT_PATH"/contracts/transfer.sol "$broker_address" >"$CURRENT_PATH/transfer.abi"
  transfer_address=$(grep Deployed <"$CURRENT_PATH/transfer.abi" | grep -o '0x.\{40\}')
  print_green "transfer contract address: $transfer_address"
  if [ -z "$transfer_address" ]; then
    echo "transfer_address id is empty"
    exit 2
  fi
  echo "$transfer_address" >"$CURRENT_PATH"/transfer_address.info
  x_replace '1,4d' "$CURRENT_PATH/transfer.abi"
  sleep 1

  print_blue "Deploy data_swapper contract"
  goduck ether contract deploy --address http://172.16.30.84:8545 --key-path "$CURRENT_PATH"/ethereum/account.key --psd-path "$CURRENT_PATH"/ethereum/password --code-path "$CURRENT_PATH"/contracts/data_swapper.sol "$broker_address" >"$CURRENT_PATH"/data_swapper.abi
  data_swapper_address=$(grep Deployed <"$CURRENT_PATH/data_swapper.abi" | grep -o '0x.\{40\}')
  print_green "data_swapper contract address: $data_swapper_address"
  if [ -z "$data_swapper_address" ]; then
    echo "data_swapper_address id is empty"
    exit 2
  fi
  echo "$data_swapper_address" >"$CURRENT_PATH/data_swapper_address.info"
  x_replace '1,4d' "$CURRENT_PATH/data_swapper.abi"
  sleep 1
}

function modify_pier_config() {
  print_blue "modify pier.toml"
  x_replace "s/$(grep "id.*appchain" <"$CURRENT_PATH/pier.toml" | awk -F '\"' '{print $2}')/$ETHER_ID/g" "$CURRENT_PATH"/pier.toml
  x_replace "s/$(grep plugin <"$CURRENT_PATH/pier.toml" | awk -F '\"' '{print $2}')/eth-client/g" "$CURRENT_PATH"/pier.toml
  x_replace "s/$(grep config <"$CURRENT_PATH/pier.toml" | awk -F '\"' '{print $2}')/ethereum/g" "$CURRENT_PATH"/pier.toml

  x_replace "s/$(grep http <"$CURRENT_PATH/pier.toml" | awk '{print $3}')/34544/g" "$CURRENT_PATH"/pier.toml
  x_replace "s/$(grep pprof <"$CURRENT_PATH/pier.toml" | awk '{print $3}')/34555/g" "$CURRENT_PATH"/pier.toml

  print_blue "modify ethereum.toml"
  x_replace "s/$(grep name <"$CURRENT_PATH/ethereum/ethereum.toml" | awk -F '\"' '{print $2}')/$ETHER_ID/g" "$CURRENT_PATH"/ethereum/ethereum.toml
}

prepare
deploy_contracts
modify_pier_config
