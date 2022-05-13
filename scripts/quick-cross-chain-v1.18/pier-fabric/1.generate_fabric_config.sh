#!/usr/bin/env bash
set -e
source ../../x.sh
source ../config.sh
CURRENT_PATH=$(pwd)

function prepare() {
  check_goduck
  print_blue "===> copy chaincode"
  rm -rf contracts.zip contracts
  cp -r "$FABRIC_PROJECT_PATH/example/contracts.zip" "$CURRENT_PATH"
  unzip -q contracts.zip

  if [ ! -f "$FABRIC_PROJECT_PATH/build/fabric-client-1.4" ]; then
    print_red "please make install fabric-client-1.4 first!!"
    exit 2
  fi
  rm -rf "$CURRENT_PATH/plugins/"
  if [ ! -f "$CURRENT_PATH/plugins/" ]; then
    mkdir "$CURRENT_PATH/plugins/"
  fi
  print_blue "===> copy plugin"
  cp -r "$FABRIC_PROJECT_PATH/build/fabric-client-1.4" "$CURRENT_PATH/plugins/"

  goduck fabric clean
  sleep 2
  goduck fabric start

  print_blue "get config file and certs from fabric install repo"
  cp -rf ~/.goduck/fabric/config.yaml "$CURRENT_PATH"/fabric/
  cp -rf ~/.goduck/fabric/crypto-config "$CURRENT_PATH"/fabric/
  cp "$CURRENT_PATH"/fabric/crypto-config/peerOrganizations/org2.example.com/peers/peer1.org2.example.com/msp/signcerts/peer1.org2.example.com-cert.pem "$CURRENT_PATH"/fabric//fabric.validators
  cp -rf contracts "$GOPATH"/src/
}

function deploy_contracts() {
  print_blue "Deploy broker contract"
  goduck fabric contract deploy --config-path "$CURRENT_PATH"/fabric/config.yaml --gopath "$GOPATH" --ccp contracts/src/broker --ccid broker --version v1.0.0
  sleep 2

  print_blue "Deploy transfer contract"
  goduck fabric contract deploy --config-path "$CURRENT_PATH"/fabric/config.yaml --gopath "$GOPATH" --ccp contracts/src/transfer --ccid transfer --version v1.0.0
  sleep 2

  print_blue "Deploy data_swapper contract"
  goduck fabric contract deploy --config-path "$CURRENT_PATH"/fabric/config.yaml --gopath "$GOPATH" --ccp contracts/src/data_swapper --ccid data_swapper --version v1.0.0
  sleep 2
}

function modify_pier_config() {
  print_blue "modify pier.toml"
  x_replace "s/$(grep "id.*appchain" <"$CURRENT_PATH/pier.toml" | awk -F '\"' '{print $2}')/$FABRIC_ID/g" "$CURRENT_PATH"/pier.toml
  x_replace "s/$(grep plugin <"$CURRENT_PATH/pier.toml" | awk -F '\"' '{print $2}')/fabric-client-1.4/g" "$CURRENT_PATH"/pier.toml
  x_replace "s/$(grep config <"$CURRENT_PATH/pier.toml" | awk -F '\"' '{print $2}')/fabric/g" "$CURRENT_PATH"/pier.toml
  x_replace "s/$(grep http <"$CURRENT_PATH/pier.toml" | awk '{print $3}')/24544/g" "$CURRENT_PATH"/pier.toml
  x_replace "s/$(grep pprof <"$CURRENT_PATH/pier.toml" | awk '{print $3}')/24555/g" "$CURRENT_PATH"/pier.toml

  print_blue "modify fabric.toml"
  x_replace "s/$(grep chain_id <"$CURRENT_PATH"/fabric/fabric.toml | awk -F '\"' '{print $2}')/$FABRIC_ID/g" "$CURRENT_PATH"/fabric/fabric.toml
}

prepare
sleep 5
deploy_contracts
modify_pier_config
