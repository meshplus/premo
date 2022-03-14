#!/usr/bin/env bash
#set -x
set -e

CURRENT_PATH=$(pwd)
source x.sh

check_pier
cd "$Pier_Project_Path" && make install
cd "$CURRENT_PATH"/pier-ether
bash 1.generate_ether_config.sh
sleep 1
bash 2.init_contracts.sh
sleep 1
bash 3.register_appchain_and_service.sh
sleep 1
bash 4.start_pier.sh
