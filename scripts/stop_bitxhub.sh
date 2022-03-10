#!/usr/bin/env bash

source x.sh

process_count=$(ps aux | grep "bitxhub --repo" | grep -v grep | wc -l)
if [ $process_count == 0 ]; then
  print_green "No bitxhub node running"
else
  print_red "Bitxhub nodes running, kill it"
  kill -9 $(ps aux | grep "bitxhub --repo" | grep -v grep | awk '{print $2}')
fi
