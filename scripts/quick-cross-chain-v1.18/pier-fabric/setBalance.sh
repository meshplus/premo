#!/usr/bin/env bash
set -e
LD_LIBRARY_PATH=$(pwd)
CONFIG_PATH=$(pwd)/fabric
export LD_LIBRARY_PATH
export CONFIG_PATH
transfer_address="transfer"

echo "1.设置初始账户Alice余额为10000"
goduck fabric contract invoke --config-path "$CONFIG_PATH"/config.yaml $transfer_address setBalance Alice,10000
sleep 2

echo "2.查询Alice当前余额"
goduck fabric contract invoke --config-path "$CONFIG_PATH"/config.yaml $transfer_address getBalance Alice
