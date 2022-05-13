#!/usr/bin/env bash
set -e
source ../../x.sh
source ../config.sh
CURRENT_PATH=$(pwd)

function check_before() {
  print_blue "===> Check pier process after start"
  process_count=$(ps aux | grep "pier-flato" | grep -v "grep" | wc -l)
  if [ "$process_count" -eq 0 ]; then
    print_green "No pier-flato running"
  else
    print_red "pier-flato running, kill it"
    ps aux | grep "pier-flato" | grep -v "grep" | awk '{print $2}' | xargs kill -9
  fi
}

function start() {
  print_blue "===> Start pier-flato"
  nohup pier --repo "$CURRENT_PATH" start 2>gc.log 1>pier.log &
}

function check_after() {
  print_blue "===> Check pier process after start"
  process_count=$(ps aux | grep "pier-flato" | grep -v "grep" | wc -l)
  if [ "$process_count" -gt 1 ]; then
    print_green "Start pier-flato successful"
  else
    print_red "Start pier-flato failed"
    exit 2
  fi
}

check_before
start
sleep 3
check_after
