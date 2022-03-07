#!/usr/bin/env bash
set -e
#set -x
export LD_LIBRARY_PATH=$(pwd)
export CONFIG_PATH=$(pwd)/fabric
transfer_address="transfer"

echo "1. get Alice balance"
goduck fabric contract invoke --config-path "$CONFIG_PATH"/config.yaml $transfer_address getBalance Alice