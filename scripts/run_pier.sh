set -e

source x.sh

CURRENT_PATH=$(pwd)
PIER_CLIENT_FABRIC_VERSION=v1.0.0-rc1

function printHelp() {
  print_blue "Usage:  "
  echo "  run_pier.sh <OPT>"
  echo "    <OPT> - one of 'up', 'down', 'restart'"
  echo "      - 'up' - bring up a new pier"
  echo "      - 'down' - clear a new pier"
  echo "    -t <mode> - pier type (default \".fabric\")"
  echo "    -v <version> - pier code version (default \".master\")"
  echo "    -r <pier_root> - pier repo path (default \".pier\")"
  echo "    -p <pier_port> - pier port (default \"8987\")"
  echo "    -b <bitxhub_addr> - bitxhub addr(default \"localhost:60011\")"
  echo "    -o <pprof_port> - pier pprof port(default \"44555\")"
  echo "  run_pier.sh -h (print this message)"
}

function prepare() {
  cd "${CURRENT_PATH}"
  if ! type goduck >/dev/null 2>&1; then
    print_blue "===> Install goduck"
    go get github.com/meshplus/goduck
  fi

  if [ "$MODE" == "fabric" ]; then
    print_blue "===> Generate fabric pier configure"
    # generate config for fabric pier
    # ...
  fi
  if [ "$MODE" == "ethereum" ]; then
    print_blue "===> Generate ethereum pier configure"
    # generate config for ethereum pier
    # ...
  fi

  if [ ! -d pier ]; then
    print_blue "===> Cloning meshplus/pier repo and checkout ${PIER_VERSION}"
    git clone https://github.com/meshplus/pier.git &&
      cd pier && git checkout ${VERSION}
  fi

  print_blue "===> Compiling meshplus/pier"
  cd "${CURRENT_PATH}"/pier
  make install
  
  cd "${CURRENT_PATH}"
  if [ ! -f fabric_rule.wasm ]; then
    print_blue "===> Downloading fabric_rule.wasm"
    wget https://raw.githubusercontent.com/meshplus/bitxhub/master/scripts/quick_start/fabric_rule.wasm
  fi
}

function pier_up() {
  prepare
  
  pier --repo="${PIER_ROOT}" init
  
  print_blue "===> pier_root: $PIER_ROOT, pier_port: $PIER_PORT, bitxhub_addr: $BITXHUB_ADDR, pprof: $PPROF_PORT"
  
  if [ "$MODE" == "fabric" ]; then
    print_blue "===> Register pier(fabric) to bitxhub"
    pier --repo "${PIER_ROOT}" appchain register \
    --name chainA \
    --type fabric \
    --desc chainA-description \
    --version 1.4.3 \
    --validators "${PIER_ROOT}"/fabric/fabric.validators
  
    print_blue "===> Deploy rule in bitxhub"
    pier --repo "${PIER_ROOT}" rule deploy --path "${CURRENT_PATH}"/fabric_rule.wasm
    print_blue "===> Start pier"
    cd "${CURRENT_PATH}"
    export CONFIG_PATH=${PIER_ROOT}/fabric
  fi

  if [ "$MODE" == "ethereum" ]; then
    print_blue "===> Register pier(ethereum) to bitxhub"
    # start ethereum pier
    # ...
  fi
    
  pier --repo "${PIER_ROOT}" start
}

function pier_down() {
  print_blue "pier PID to kill: $PIER_PID"
  kill -9 $PIER_PID
  print_blue "kill result: $?"
}


PIER_ROOT=${CURRENT_PATH}/.pier
PIER_PORT=8987
BITXHUB_ADDR="localhost:60011"
MODE="fabric"
VERSION="master"

OPT=$1
shift

while getopts "h?t:v:r:p:b:o:i:" opt; do
  case "$opt" in
  h | \?)
    printHelp
    exit 0
    ;;
  t)
    MODE=$OPTARG
    ;;
  v)
    VERSION=$OPTARG
    ;;
  r)
    PIER_ROOT=$OPTARG
    ;;
  p)
    PIER_PORT=$OPTARG
    ;;
  b)
    BITXHUB_ADDR=$OPTARG
    ;;
  o)
    PPROF_PORT=$OPTARG
    ;;
  i)
    PIER_PID=$OPTARG
    ;;
  esac
done

if [ "$OPT" == "up" ]; then
  pier_up
elif [ "$OPT" == "down" ]; then
  pier_down
else
  printHelp
  exit 1
fi
