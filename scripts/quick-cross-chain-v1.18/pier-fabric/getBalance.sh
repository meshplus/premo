#!/usr/bin/env bash
set -e
LD_LIBRARY_PATH=$(pwd)
CONFIG_PATH=$(pwd)/fabric
export LD_LIBRARY_PATH
export CONFIG_PATH
transfer_address="transfer"

echo "1. get Alice balance"
goduck fabric contract invoke --config-path "$CONFIG_PATH"/config.yaml $transfer_address getBalance Alice
