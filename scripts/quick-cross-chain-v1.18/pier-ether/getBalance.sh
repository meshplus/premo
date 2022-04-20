#!/usr/bin/env bash
set -e
#set -x
LD_LIBRARY_PATH=$(pwd)
export LD_LIBRARY_PATH
transfer_address=$(awk '{print $1}' <transfer_address.info)

echo "get Alice balance"
goduck ether contract invoke --address http://172.16.30.84:8545 --key-path ethereum/account.key --psd-path ethereum/password --abi-path transfer.abi "$transfer_address" getBalance Alice
