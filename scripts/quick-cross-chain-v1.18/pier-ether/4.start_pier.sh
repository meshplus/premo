#!/usr/bin/env bash
set -e
#set -x
CURRENT_PATH=$(pwd)
PROJECT_PATH=$(dirname "${CURRENT_PATH}")
source "$PROJECT_PATH"/x.sh

function check_before() {
  print_blue "===> Check pier process after start"
  process_count=$(ps -ef | grep "pier-ether" | grep -v grep | wc -l)
  if [ $process_count -eq 0 ]; then
    print_green "No pier-ether running"
  else
    print_red "pier-ether running, kill it"
    kill -9 $(ps aux | grep "pier-ether" | grep -v grep | awk '{print $2}')
  fi
}

function start() {
  print_blue "===> Start pier-ether"
  nohup pier --repo "$CURRENT_PATH" start 2>gc.log 1>pier.log &
}

function check_after() {
  print_blue "===> Check pier process after start"
  process_count=$(ps -ef | grep "pier-ether" | grep -v grep | wc -l)
  if [ $process_count -gt 1 ]; then
    print_green "Start pier-ether successed"
  else
    print_red "Start pier-ether failed"
    exit 2
  fi
}

check_before
start
sleep 3
check_after
