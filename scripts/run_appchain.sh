set -e

source x.sh

CURRENT_PATH=$(pwd)
OPT=$1
MODE=$2

function printHelp() {
  print_blue "Usage:  "
  echo "  run_appchain.sh <OPT>"
  echo "    <OPT> - one of 'up', 'down'"
  echo "      - 'up' - bring up the appchain"
  echo "      - 'down' - shut down the appchain"
  echo "  run_appchain.sh -h (print this message)"
}

function prepare() {
  if ! type goduck >/dev/null 2>&1; then
    print_blue "===> Install goduck"
    go get github.com/meshplus/goduck/cmd/goduck
  fi

  if [ ! -d "$HOME/.goduck" ]; then
    goduck init
  fi
}

function appchain_up() {
  prepare

  if [ "$MODE" == "ethereum" ]; then
    print_blue "===> Start ethereum appchain"
    goduck ether start
  fi
  if [ "$MODE" == "fabric" ]; then
    print_blue "===> Start fabric appchain"
    goduck fabric start

    if [ "$(docker container ls | grep -c broker)" == 0 ]; then
        goduck fabric chaincode
    fi
  fi
}

function appchain_down() {
  prepare

  if [ "$MODE" == "ethereum" ]; then
    print_blue "===> Stop ethereum appchain"
    goduck ether stop
  fi
  if [ "$MODE" == "fabric" ]; then
    print_blue "===> Stop fabric appchain"
    goduck fabric stop
  fi
}

if [ "$OPT" == "up" ]; then
  appchain_up
elif [ "$OPT" == "down" ]; then
  appchain_down
else
  printHelp
  exit 1
fi
