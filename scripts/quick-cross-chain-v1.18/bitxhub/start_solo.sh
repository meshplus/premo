#!/usr/bin/env bash

set -e

CURRENT_PATH=$(pwd)
PROJECT_PATH=$(dirname "${CURRENT_PATH}")
source "$PROJECT_PATH"/x.sh

REPO_PATH=${CURRENT_PATH}/repo_solo

function check_before() {
  print_blue "===> Check bitxhub process before start"
  process_count=$(ps aux |grep "bitxhub --repo" |grep -v grep|wc -l)
  if [ $process_count == 0 ]; then
    print_green "No bitxhub node running"
  else
    print_red "Bitxhub nodes running, kill it"
    kill -9 $(ps aux |grep "bitxhub --repo" |grep -v grep|awk '{print $2}')
  fi
}

function config() {
  print_blue "===> Generate bitxhub config"
  rm -rf "${REPO_PATH}"
  mkdir -p "${REPO_PATH}"
  cp -rf "${BitXHub_Project_Path}"/scripts/certs/node1/* "${REPO_PATH}"
  cp -rf "${BitXHub_Project_Path}"/config/* "${REPO_PATH}"

  echo " #!/usr/bin/env bash" >"${REPO_PATH}"/start.sh
  echo "bitxhub --repo \$(pwd)" start >>"${REPO_PATH}"/start.sh
  x_sed "s/solo = false/solo = true/g" ${REPO_PATH}/bitxhub.toml
  x_sed "s/raft/solo/g" ${REPO_PATH}/bitxhub.toml

}

function start() {
  print_blue "===> Start solo bitxhub"
  cd "${REPO_PATH}"
  nohup bash start.sh 2>gc.log 1>node.log &
}

function check_after() {
  print_blue "===> Check bitxhub process after start"
  process_count=$(ps -ef |grep "bitxhub --repo" |grep -v grep|wc -l)
  if [ $process_count -gt 0 ]; then
    print_green "Start bitxhub successed"
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
