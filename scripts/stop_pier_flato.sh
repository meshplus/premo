#!/usr/bin/env bash
source x.sh
process_count=$(ps aux | grep "pier-flato" | grep -v "grep" | wc -l)
if [ "$process_count" == 0 ]; then
  print_green "No bitxhub node running"
else
  print_red "Bitxhub nodes running, kill it"
  ps aux | grep "pier-flato" | grep -v "grep" | awk '{print $2}' | xargs kill -9
fi

