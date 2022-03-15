#!/usr/bin/env bash
set -e
#set -x
export LD_LIBRARY_PATH=$(pwd)
export CONFIG_PATH=$(pwd)/fabric
transfer_address="transfer"
dst_appchainID=${1:-"fltaoappchain1"}
dst_transfer_address=${2:-"0xED35A2b46e5f8c89990262B636ed5E9e705C0FBb"}

echo "1. get Alice balance"
goduck fabric contract invoke --config-path "$CONFIG_PATH"/config.yaml $transfer_address getBalance Alice
sleep 2

echo "2. transfer"
goduck fabric contract invoke --config-path "$CONFIG_PATH"/config.yaml $transfer_address transfer 1356:$dst_appchainID:$dst_transfer_address,Alice,Alice,100
