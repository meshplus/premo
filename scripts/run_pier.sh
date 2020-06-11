set -e

source x.sh

CURRENT_PATH=$(pwd)
PIER_CLIENT_FABRIC_VERSION=master
PIER_CLIENT_ETHEREUM_VERSION=master
GODUCK_VERSION=master

function printHelp() {
  print_blue "Usage:  "
  echo "  run_pier.sh <OPT>"
  echo "    <OPT> - one of 'up', 'down', 'restart'"
  echo "      - 'up' - bring up a new pier"
  echo "      - 'down' - clear a new pier"
  echo "    -t <mode> - pier type (default \"fabric\")"
  echo "    -v <version> - pier code version (default \"master\")"
  echo "    -r <pier_root> - pier repo path (default \".pier-fabric\")"
  echo "    -b <bitxhub_addr> - bitxhub addr(default \"localhost:60011\")"
  echo "  run_pier.sh -h (print this message)"
}

function prepare() {
  cd "${CURRENT_PATH}"
  if ! type goduck >/dev/null 2>&1; then
    print_blue "===> Install goduck"
    go get github.com/meshplus/goduck &&
      cd goduck && git checkout ${GODUCK_VERSION}
    make install
  fi
  
  cd "${CURRENT_PATH}"
  if [ ! -d pier ]; then
    print_blue "===> Cloning meshplus/pier repo and checkout ${PIER_VERSION}"
    git clone https://github.com/meshplus/pier.git
  fi
  cd pier && git checkout ${VERSION}

  print_blue "===> Compiling meshplus/pier"
  cd "${CURRENT_PATH}"/pier
  make install
  
  if [ "$MODE" == "fabric" ]; then
    print_blue "===> Generate fabric pier configure"
    # generate config for fabric pier
    # ...
    cd "${CURRENT_PATH}"
    if [ ! -d pier-client-fabric ]; then
        print_blue "===> Cloning meshplus/pier-client-fabric repo and checkout ${PIER_CLIENT_FABRIC_VERSION}"
        git clone https://github.com/meshplus/pier-client-fabric.git
    fi
    cd pier-client-fabric && git checkout ${PIER_CLIENT_FABRIC_VERSION}
    print_blue "===> Compiling meshplus/pier-client-fabric"
    cd "${CURRENT_PATH}"/pier-client-fabric
    make fabric1.4
    
     cd "${CURRENT_PATH}"
    if [ ! -f fabric_rule.wasm ]; then
        print_blue "===> Downloading fabric_rule.wasm"
        wget https://raw.githubusercontent.com/meshplus/bitxhub/master/scripts/quick_start/fabric_rule.wasm
    fi
  fi

  if [ "$MODE" == "ethereum" ]; then
    print_blue "===> Generate ethereum pier configure"
    # generate config for ethereum pier
    # ...
    cd "${CURRENT_PATH}"
    if [ ! -d pier-client-ethereum ]; then
        print_blue "===> Cloning meshplus/pier-client-ethereum repo and checkout ${PIER_CLIENT_ETHEREUM_VERSION}"
        git clone https://github.com/meshplus/pier-client-ethereum.git
    fi
    cd pier-client-ethereum && git checkout ${PIER_CLIENT_ETHEREUM_VERSION}
    print_blue "===> Compiling meshplus/pier-client-ethereum"
    cd "${CURRENT_PATH}"/pier-client-ethereum
    make eth
    
    if [ ! -f ethereum_rule.wasm ]; then
        print_blue "===> Downloading ethereum_rule.wasm"
        # wget https://raw.githubusercontent.com/meshplus/bitxhub/master/scripts/quick_start/ethereum_rule.wasm
    fi
  fi
}

function appchain_register(){
    pier --repo "${PIER_ROOT}" appchain register \
    --name $1 \
    --type $2 \
    --desc $3 \
    --version $4 \
    --validators "${PIER_ROOT}/$5"
}

function rule_deploy(){
    print_blue "===> deploy path: ${CURRENT_PATH}/$1_rule.wasm"
    pier --repo "${PIER_ROOT}" rule deploy --path "${CURRENT_PATH}/$1_rule.wasm"
}

function pier_up() {
  prepare

  pier --repo="${PIER_ROOT}" init
  mkdir -p "${PIER_ROOT}"/plugins
  print_blue "===> pier_root: $PIER_ROOT, bitxhub_addr: $BITXHUB_ADDR"
  
  if [ "$MODE" == "fabric" ]; then
    print_blue "===> Copy fabric plugins"
    cp "${CURRENT_PATH}"/pier-client-fabric/build/fabric-client-1.4.so "${PIER_ROOT}"/plugins/
    print_blue "===> Register pier(fabric) to bitxhub"
    appchain_register chainA fabric chainA-description 1.4.3 fabric/fabric.validators
    print_blue "===> Deploy rule in bitxhub"
    rule_deploy fabric
    print_blue "===> Start pier"
    cd "${CURRENT_PATH}"
    export CONFIG_PATH=${PIER_ROOT}/fabric
  fi

  if [ "$MODE" == "ethereum" ]; then
    print_blue "===> Copy ethereum plugins"
    cp "${CURRENT_PATH}"/pier-client-ethereum/build/eth-client.so.so "${PIER_ROOT}"/plugins/
    print_blue "===> Register pier(ethereum) to bitxhub"
    # start ethereum pier
    # ...
  fi
    
  nohup pier --repo "${PIER_ROOT}" start >/dev/null 2>&1 &
  echo $! >>"${CURRENT_PATH}/pier-${MODE}.pid"
}

function pier_down() {
  set +e
  print_blue "===> Kill $MODE pier"

  if [ -a pier-$MODE.pid ]; then
    pid=$(cat pier-$MODE.pid)
    kill "$pid"
    if [ $? -eq 0 ]; then
      echo "pier-$MODE pid:$pid exit"
    else
      print_red "pier exit fail, try use kill -9 $pid"
    fi
    rm pier-$MODE.pid
  fi
}

PIER_ROOT=${CURRENT_PATH}/.pier-fabric
BITXHUB_ADDR="localhost:60011"
MODE="fabric"
VERSION="master"

OPT=$1
shift

while getopts "h?t:v:r:b:" opt; do
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
  b)
    BITXHUB_ADDR=$OPTARG
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
