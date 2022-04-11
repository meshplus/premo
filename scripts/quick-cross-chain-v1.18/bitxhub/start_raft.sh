#!/usr/bin/env bash
set -e
source ../../x.sh
source ../config.sh
CURRENT_PATH=$(pwd)
REPO_PATH=${CURRENT_PATH}/repo_raft

function check_before() {
  print_blue "===> Check bitxhub process before start"
  process_count=$(ps aux | grep "bitxhub --repo" | grep -v "grep" | wc -l)
  if [ "$process_count" == 0 ]; then
    print_green "No bitxhub node running"
  else
    print_red "Bitxhub nodes running, kill it"
    ps aux | grep "bitxhub --repo" | grep -v "grep" | awk '{print $2}' | xargs kill -9
  fi
}

function config() {
  print_blue "===> Generate bitxhub config"
  rm -rf "${REPO_PATH}"
  for ((i = 1; i < 5; i = i + 1)); do
    root=${REPO_PATH}/node${i}
    mkdir -p "${root}"
    cp -rf "${BITXHUB_PROJECT_PATH}"/scripts/certs/node"${i}"/* "${root}"
    cp -rf "${BITXHUB_PROJECT_PATH}"/config/* "${root}"

    echo " #!/usr/bin/env bash" >"${root}"/start.sh
    echo "bitxhub --repo \$(pwd)" start >>"${root}"/start.sh

    x_replace "s/60011/6001${i}/g" "${root}/bitxhub.toml"
    x_replace "s/9091/909${i}/g" "${root}/bitxhub.toml"
    x_replace "s/53121/5312${i}/g" "${root}/bitxhub.toml"
    x_replace "s/40011/4001${i}/g" "${root}/bitxhub.toml"
    x_replace "s/8881/888${i}/g" "${root}/bitxhub.toml"
    x_replace "s/SimpleMajority/ZeroPermission/g" "${root}/bitxhub.toml"
    x_replace "1s/1/${i}/" "${root}/network.toml"
    x_replace "1,/true/s/true/false/" "${root}/order.toml"
  done
}

function start() {
  print_blue "===> Start raft bitxhub"
  for ((i = 1; i < 5; i = i + 1)); do
    cd "${REPO_PATH}/node${i}"
    nohup bash start.sh 2>gc.log 1>node.log &
  done
}

function check_after() {
  print_blue "===> Check bitxhub process after start"
  process_count=$(ps aux | grep "bitxhub --repo" | grep -v "grep" | wc -l)
  if [ "$process_count" -gt 2 ]; then
    print_green "Start bitxhub successful"
  else
    print_red "Start bitxhub failed"
    exit
  fi
}

check_before
config
start
sleep 5
check_after
