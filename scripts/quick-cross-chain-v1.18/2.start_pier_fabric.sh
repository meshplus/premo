#!/usr/bin/env bash
set -e

CURRENT_PATH=$(pwd)
source x.sh

check_pier
cd "$CURRENT_PATH"/pier-fabric
bash 1.generate_fabric_config.sh
sleep 2
bash 2.init_contracts.sh
sleep 2
bash 3.register_appchain_and_service.sh
sleep 2
bash 4.start_pier.sh
