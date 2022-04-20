#!/usr/bin/env bash
set -e

CURRENT_PATH=$(pwd)
source ../x.sh

dst_appchain="$Ether_ChainID"
print_green "dst_appchain: $dst_appchain"

dst_service=$(cat "$CURRENT_PATH"/pier-ether/transfer_address.info | awk {'print $1'})
print_green "dst_service: $dst_service"
amount=10

function get_flato_balance() {
  print_blue "get Alice balance on flato"
  cd "$CURRENT_PATH"/pier-flato
  bash getBalance.sh | grep receipt | awk {'print $7'}
}

function get_ether_balance() {
  print_blue "get Alice balance on ether"
  cd "$CURRENT_PATH"/pier-ether
  bash getBalance.sh | grep result | awk {'print $3'}
}

function flato_transferTo_ether() {
  flato_balance_before=$(get_flato_balance)
  ether_balance_before=$(get_ether_balance)
  print_green "before interchain, Alice balance on flato is: $flato_balance_before"
  print_green "before interchain, Alice balance on ether is: $ether_balance_before"

  cd "$CURRENT_PATH"/pier-flato
  bash interchain_transfer.sh "$dst_appchain" "$dst_service" $amount
  #sleep 2

  flato_balance_after=$(get_flato_balance)
  ether_balance_after=$(get_ether_balance)
  print_green "after interchain, Alice balance on flato is: $flato_balance_after"
  print_green "after interchain, Alice balance on ether is: $ether_balance_after"
}

flato_transferTo_ether
