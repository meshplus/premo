#!/usr/bin/env bash

RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

function print_blue() {
  printf "${BLUE}%s${NC}\n" "$1"
}

function print_green() {
  printf "${GREEN}%s${NC}\n" "$1"
}

function print_red() {
  printf "${RED}%s${NC}\n" "$1"
}

# The sed commend with system judging
# Examples:
# sed -i 's/a/b/g' bob.txt => x_sed 's/a/b/g' bob.txt
function x_sed() {
  system=$(uname)

  if [ "${system}" = "Linux" ]; then
    sed -i "$@"
  else
    sed -i '' "$@"
  fi
}

function check_goduck() {
  if ! type goduck >/dev/null 2>&1; then
    print_blue "===> Install goduck"
    go get github.com/meshplus/goduck/cmd/goduck
    goduck init
  fi
}

function check_pier() {
  if ! type pier >/dev/null 2>&1; then
    print_blue "===> Compileing pier"
    cd "$Pier_Project_Path" || exit
    make install
  fi
}

function check_bitxhub() {
  if ! type bitxhub >/dev/null 2>&1; then
    print_blue "===> Compileing bitxhub"
    cd "$BitXHub_Project_Path" || exit
    make install
  fi
}

# BitXHub Config
BitXHub_Project_Path="$GOPATH/src/meshplus/bitxhub"
BitXHub_Type="raft"
#BitXHub_Addrs=["localhost:60011", "localhost:60012", "localhost:60013", "localhost:60014"]

# Pier Config
Pier_Project_Path="$GOPATH/src/meshplus/pier"

# Pier-flato Config
Flato_Project_Path="$GOPATH/src/meshplus/pier-client-flato"
Flato_ChainID="flatoappchain1"
#默认先使用happy的验证规则
Flato_Rule_address="0x00000000000000000000000000000000000000a2"

# Pier-fabric Config
Fabric_Project_Path="$GOPATH/src/meshplus/pier-client-fabric"
Fabric_ChainID="fabricappchain1"
#默认先使用happy的验证规则
Flato_Rule_address="0x00000000000000000000000000000000000000a2"

# Pier-ether Config
Ether_Project_Path="$GOPATH/src/meshplus/pier-client-ethereum"
Ether_ChainID="etherappchain1"
#默认先使用happy的验证规则
Ether_Rule_address="0x00000000000000000000000000000000000000a2"
