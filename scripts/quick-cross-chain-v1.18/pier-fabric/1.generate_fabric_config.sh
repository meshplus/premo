#!/usr/bin/env bash
set -e
#set -x
CURRENT_PATH=$(pwd)
PROJECT_PATH=$(dirname "${CURRENT_PATH}")
source "$PROJECT_PATH"/x.sh

function prepare() {
    check_goduck
    if [ ! -d "$CURRENT_PATH"/contracts ]; then
        print_blue "===> copy chaincode"
        cp -r "$Fabric_Project_Path"/example/contracts.zip "$CURRENT_PATH"/
        unzip -q contracts.zip
        rm contracts.zip
    fi
    
    if [ ! -f "$CURRENT_PATH"/plugins/fabric-client-1.4 ]; then
        if [ ! -f "$Fabric_Project_Path"/build/fabric-client-1.4 ]; then
            print_red "please make install fabric-clent-1.4 first!!"
            exit 2
        fi
        if [ ! -f "$CURRENT_PATH"/plugins/ ]; then
            mkdir "$CURRENT_PATH"/plugins/
        fi
        print_blue "===> copy plugin"
        cp -r "$Fabric_Project_Path"/build/fabric-client-1.4 "$CURRENT_PATH"/plugins/
    fi
    
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
    x_sed "s/$(cat "$CURRENT_PATH"/pier.toml | grep "id.*appchain" | awk -F '\"' {'print $2'})/$Fabric_ChainID/g" "$CURRENT_PATH"/pier.toml
    x_sed "s/$(cat "$CURRENT_PATH"/pier.toml | grep plugin | awk -F '\"' {'print $2'})/fabric-client-1.4/g" "$CURRENT_PATH"/pier.toml
    x_sed "s/$(cat "$CURRENT_PATH"/pier.toml | grep config | awk -F '\"' {'print $2'})/fabric/g" "$CURRENT_PATH"/pier.toml
    x_sed "s/$(cat "$CURRENT_PATH"/pier.toml | grep http | awk {'print $3'})/24544/g" "$CURRENT_PATH"/pier.toml
    x_sed "s/$(cat "$CURRENT_PATH"/pier.toml | grep pprof | awk {'print $3'})/24555/g" "$CURRENT_PATH"/pier.toml
    
    print_blue "modify fabric.toml"
    x_sed "s/$(cat "$CURRENT_PATH"/fabric/fabric.toml | grep chain_id | awk -F '\"' {'print $2'})/$Fabric_ChainID/g" "$CURRENT_PATH"/fabric/fabric.toml
}

prepare
sleep 5
deploy_contracts
modify_pier_config
