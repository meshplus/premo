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
    goduck pier start --chain ethereum --type docker
  fi
  if [ "$MODE" == "fabric" ]; then
    print_blue "===> Start fabric appchain"
    goduck pier start --chain fabric --type docker
  fi
}

function appchain_down() {
  prepare

  if [ "$MODE" == "ethereum" ]; then
    print_blue "===> Stop ethereum appchain"
    goduck pier stop --chain ethereum
  fi
  if [ "$MODE" == "fabric" ]; then
    print_blue "===> Start fabric appchain"
    goduck pier stop --chain fabric
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
