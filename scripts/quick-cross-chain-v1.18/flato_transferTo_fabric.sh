#!/usr/bin/env bash
set -e

CURRENT_PATH=$(pwd)
source ../x.sh

dst_appchain="$Fabric_ChainID"
print_green "dst_appchain: $dst_appchain"

dst_service="mychannel&transfer"
print_green "dst_service: $dst_service"
amount=10

function get_flato_balance() {
  print_blue "get Alice balance on flato"
  cd "$CURRENT_PATH"/pier-flato
  bash getBalance.sh | grep receipt | awk {'print $7'}
}

function get_fabric_balance() {
  print_blue "get Alice balance on fabric"
  cd "$CURRENT_PATH"/pier-fabric
  bash getBalance.sh | grep result | awk {'print $3'}
}

function flato_transferTo_fabric() {
  flato_balance_before=$(get_flato_balance)
  fabric_balance_before=$(get_fabric_balance)
  print_green "before interchain, Alice balance on flato is: $flato_balance_before"
  print_green "before interchain, Alice balance on fabric is: $fabric_balance_before"

  cd "$CURRENT_PATH"/pier-flato
  bash interchain_transfer.sh "$dst_appchain" "$dst_service" $amount
  #sleep 2

  flato_balance_after=$(get_flato_balance)
  fabric_balance_after=$(get_fabric_balance)
  print_green "after interchain, Alice balance on flato is: $flato_balance_after"
  print_green "after interchain, Alice balance on fabric is: $fabric_balance_after"
}

flato_transferTo_fabric
