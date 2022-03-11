#!/usr/bin/env bash

set -e

CURRENT_PATH=$(pwd)
source x.sh

BitXHub_Type="$(cat x.sh | grep BitXHub_Type | awk -F '\"' {'print $2'})"
check_bitxhub
cd "${BitXHub_Project_Path}" && make install
cd "$CURRENT_PATH"/bitxhub
if [ "$BitXHub_Type" == solo ]; then
  bash start_solo.sh
else
  bash start_raft.sh
fi
