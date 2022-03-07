#!/usr/bin/env bash
set -e
#set -x
export LD_LIBRARY_PATH=$(pwd)
transfer_address=$(cat transfer_address.info | awk {'print $1'})

#echo "1. set Alice balance"
#goduck ether contract invoke --address http://172.16.30.84:8545 --key-path ethereum/account.key --psd-path ethereum/password --abi-path transfer.abi $transfer_address setBalance Alice,10000 
#sleep 2

echo "1. get Alice balance"
goduck ether contract invoke --address http://172.16.30.84:8545 --key-path ethereum/account.key --psd-path ethereum/password --abi-path transfer.abi $transfer_address getBalance Alice
sleep 2

echo "2. transfer"
goduck ether contract invoke --address http://172.16.30.84:8545 --key-path ethereum/account.key --psd-path ethereum/password --abi-path transfer.abi $transfer_address transfer 1356:fabricappchain1:"mychannel&transfer",Alice,Alice,10
