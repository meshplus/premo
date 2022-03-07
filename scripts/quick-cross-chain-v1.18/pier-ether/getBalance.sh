#!/usr/bin/env bash
set -e
#set -x
export LD_LIBRARY_PATH=$(pwd)
transfer_address=$(cat transfer_address.info | awk {'print $1'})

echo "get Alice balance"
goduck ether contract invoke --address http://172.16.30.84:8545 --key-path ethereum/account.key --psd-path ethereum/password --abi-path transfer.abi $transfer_address getBalance Alice