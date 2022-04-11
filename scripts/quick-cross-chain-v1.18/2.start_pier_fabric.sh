#!/usr/bin/env bash
set -e
source ../x.sh
source ./config.sh
CURRENT_PATH=$(pwd)

cd "$PIER_PROJECT_PATH" && make install
cd "$CURRENT_PATH"/pier-fabric
bash 1.generate_fabric_config.sh
bash 2.init_contracts.sh
bash 3.register_appchain_and_service.sh
bash 4.start_pier.sh
