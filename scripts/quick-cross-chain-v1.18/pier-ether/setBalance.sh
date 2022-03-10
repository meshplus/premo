#!/usr/bin/env bash
set -e
#set -x
export LD_LIBRARY_PATH=$(pwd)
transfer_address=$(cat flato/flato.toml | grep transfer.abi | awk {'print $1'})

echo "1.设置初始账户Alice余额为10000"
goduck hpc invoke --config-path flato --abi-path transfer.abi $transfer_address setBalance Alice,10000

echo "2.查询Alice当前余额"
goduck hpc invoke --config-path flato --abi-path transfer.abi $transfer_address getBalance Alice
