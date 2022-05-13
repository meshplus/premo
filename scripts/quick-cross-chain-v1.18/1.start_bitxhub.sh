#!/usr/bin/env bash
set -e
source ../x.sh
source ./config.sh
CURRENT_PATH=$(pwd)

cd "${BITXHUB_PROJECT_PATH}" && make install
cd "$CURRENT_PATH"/bitxhub
if [ "$BITXHUB_TYPE" == "solo" ]; then
  bash start_solo.sh
else
  bash start_raft.sh
fi
