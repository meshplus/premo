#!/usr/bin/env bash
set -e
#set -x
CURRENT_PATH=$(pwd)
PROJECT_PATH=$(dirname "${CURRENT_PATH}")
source "$PROJECT_PATH"/x.sh

function check_before() {
  print_blue "===> Check pier process after start"
  process_count=$(ps -ef | grep "pier-fabric" | grep -v grep | wc -l)
  if [ $process_count -eq 0 ]; then
    print_green "No pier-fabric running"
  else
    print_red "pier-flato running, kill it"
    kill -9 $(ps aux | grep "pier-fabric" | grep -v grep | awk '{print $2}')
  fi
}

function start() {
  print_blue "===> Start pier-fabric"
  nohup pier --repo "$CURRENT_PATH" start 2>gc.log 1>pier.log &
}

function check_after() {
  print_blue "===> Check pier process after start"
  process_count=$(ps -ef | grep "pier-fabric" | grep -v grep | wc -l)
  if [ $process_count -gt 1 ]; then
    print_green "Start pier-fabric successed"
  else
    print_red "Start pier-fabric failed"
    exit 2
  fi
}

check_before
start
sleep 3
check_after
