#!/usr/bin/env bash

source x.sh

process_count=$(ps aux | grep pier-ether | grep -v grep | wc -l)
if [ $process_count == 0 ]; then
  print_green "No pier running"
else
  print_red "Pier running, kill it"
  kill -9 $(ps aux | grep pier-ether | grep -v grep | awk '{print $2}')
fi
