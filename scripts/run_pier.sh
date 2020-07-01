set -e

source x.sh

CURRENT_PATH=$(pwd)
GODUCK_REPO_PATH=~/.goduck
PIER_CLIENT_FABRIC_VERSION=master
PIER_CLIENT_ETHEREUM_VERSION=master
VERSION=master

function printHelp() {
  print_blue "Usage:  "
  echo "  run_pier.sh <OPT>"
  echo "    <OPT> - one of 'up', 'down', 'restart'"
  echo "      - 'up' - bring up a new pier"
  echo "      - 'down' - clear a new pier"
  echo "    -t <mode> - pier type (default \"fabric\")"
  echo "    -v <version> - pier code version (default \"master\")"
  echo "    -r <pier_root> - pier repo path (default \".pier_fabric\")"
  echo "    -b <bitxhub_addr> - bitxhub addr(default \"localhost:60011\")"
  echo "  run_pier.sh -h (print this message)"
}

function prepare() {
  cd "${CURRENT_PATH}"
  if ! type goduck >/dev/null 2>&1; then
    print_blue "===> Install goduck"
    go get github.com/meshplus/goduck/cmd/goduck
  fi

  if [ ! -d "$HOME/.goduck" ]; then
    goduck init
  fi

  if [ -f "${CURRENT_PATH}/pier-${MODE}.pid" ]; then
    print_red "pier-${MODE} is already running in the background service"
    cat "${CURRENT_PATH}/pier-${MODE}.pid"
    exit 0
  fi

  cd "${CURRENT_PATH}"
  if [ ! -d pier ]; then
    print_blue "===> Cloning meshplus/pier repo and checkout ${PIER_VERSION}"
    git clone https://github.com/meshplus/pier.git
  fi
  cd pier && git checkout -f master && git reset --hard HEAD
  git pull && git checkout ${VERSION}

  print_blue "===> Compiling meshplus/pier"
  cd "${CURRENT_PATH}"/pier
  make install

  if [ "$MODE" == "fabric" ]; then
    print_blue "===> Generate fabric pier configure"
    # generate config for fabric pier
    PIER_ROOT="${CURRENT_PATH}"/.pier_fabric
    cd "${CURRENT_PATH}"
    if [ ! -d .pier_fabric ]; then
      mkdir .pier_fabric
    fi

    goduck pier config \
      --mode "relay" \
      --bitxhub "localhost:60011" \
      --validators "0xe6f8c9cf6e38bd506fae93b73ee5e80cc8f73667" \
      --validators "0x8374bb1e41d4a4bb4ac465e74caa37d242825efc" \
      --validators "0x759801eab44c9a9bbc3e09cb7f1f85ac57298708" \
      --validators "0xf2d66e2c27e93ff083ee3999acb678a36bb349bb" \
      --appchain-type "fabric" \
      --appchain-IP "127.0.0.1" \
      --target "${PIER_ROOT}"
    # copy appchain crypto-config and modify config.yaml
    if [ ! -d "${GODUCK_REPO_PATH}"/crypto-config ]; then
      print_red "crypto-config not found, please start fabric network first"
      exit 1
    fi
    cp -r "${GODUCK_REPO_PATH}"/crypto-config "${PIER_ROOT}"/fabric

    if [ ! -d pier-client-fabric ]; then
      print_blue "===> Cloning meshplus/pier-client-fabric repo and checkout ${PIER_CLIENT_FABRIC_VERSION}"
      git clone https://github.com/meshplus/pier-client-fabric.git
    fi
    cd pier-client-fabric && git checkout -f master && git reset --hard HEAD
    git pull && git checkout ${PIER_CLIENT_FABRIC_VERSION}
    print_blue "===> Compiling meshplus/pier-client-fabric"
    cd "${CURRENT_PATH}"/pier-client-fabric
    make fabric1.4

    cd "${CURRENT_PATH}"
    if [ ! -f fabric_rule.wasm ]; then
      print_blue "===> Downloading fabric_rule.wasm"
      wget https://github.com/meshplus/bitxhub/raw/v1.0.0-rc3/scripts/quick_start/fabric_rule.wasm
    fi

    PEM_PATH="${PIER_ROOT}"/fabric/crypto-config/peerOrganizations/org2.example.com/peers/peer1.org2.example.com/msp/signcerts/peer1.org2.example.com-cert.pem
    if [ ! -f "${PEM_PATH}" ]; then
      PEM_PATH="${PIER_ROOT}"/fabric/crypto-config/peerOrganizations/org2.example1.com/peers/peer1.org2.example1.com/msp/signcerts/peer1.org2.example1.com-cert.pem
    fi
    cp "${PEM_PATH}" "${PIER_ROOT}"/fabric/fabric.validators
  fi

  if [ "$MODE" == "ethereum" ]; then
    print_blue "===> Generate ethereum pier configure"
    # generate config for ethereum pier
    PIER_ROOT="${CURRENT_PATH}"/.pier_ethereum
    cd "${CURRENT_PATH}"
    if [ ! -d ".pier_ethereum" ]; then
      mkdir .pier_ethereum
    fi
    cd "${PIER_ROOT}"
    goduck pier config \
      --mode "relay" \
      --ID 1 \
      --bitxhub "localhost:60011" \
      --validators "0xe6f8c9cf6e38bd506fae93b73ee5e80cc8f73667" \
      --validators "0x8374bb1e41d4a4bb4ac465e74caa37d242825efc" \
      --validators "0x759801eab44c9a9bbc3e09cb7f1f85ac57298708" \
      --validators "0xf2d66e2c27e93ff083ee3999acb678a36bb349bb" \
      --appchain-type "ethereum" \
      --appchain-IP "127.0.0.1" \
      --target "${PIER_ROOT}"
    cd "${CURRENT_PATH}"
    if [ ! -d pier-client-ethereum ]; then
      print_blue "===> Cloning meshplus/pier-client-ethereum repo and checkout ${PIER_CLIENT_ETHEREUM_VERSION}"
      git clone https://github.com/meshplus/pier-client-ethereum.git
    fi
    cd pier-client-ethereum && git checkout -f master && git reset --hard HEAD
    git pull && git checkout ${PIER_CLIENT_ETHEREUM_VERSION}
    print_blue "===> Compiling meshplus/pier-client-ethereum"
    cd "${CURRENT_PATH}"/pier-client-ethereum
    make eth
    cd "${CURRENT_PATH}"
    if [ ! -f ethereum_rule.wasm ]; then
      print_blue "===> Downloading ethereum_rule.wasm"
      wget https://github.com/meshplus/pier-client-ethereum/blob/master/config/validating.wasm
      mv validating.wasm ethereum_rule.wasm
    fi
  fi
}

function appchain_register() {
  pier --repo "${PIER_ROOT}" appchain register \
    --name $1 \
    --type $2 \
    --desc $3 \
    --version $4 \
    --validators "${PIER_ROOT}/$5"
}

function rule_deploy() {
  print_blue "===> deploy path: ${CURRENT_PATH}/$1_rule.wasm"
  pier --repo "${PIER_ROOT}" rule deploy --path "${CURRENT_PATH}/$1_rule.wasm"
}

function pier_up() {
  prepare

  cd "${CURRENT_PATH}"
  mkdir -p "${PIER_ROOT}"/plugins
  print_blue "===> pier_root: $PIER_ROOT, bitxhub_addr: $BITXHUB_ADDR"

  if [ "$MODE" == "fabric" ]; then
    print_blue "===> Copy fabric plugins"
    cp "${CURRENT_PATH}"/pier-client-fabric/build/fabric-client-1.4.so "${PIER_ROOT}"/plugins/
    print_blue "===> Register pier(fabric) to bitxhub"
    appchain_register chainA fabric chainA-description 1.4.3 fabric/fabric.validators
    print_blue "===> Deploy rule in bitxhub"
    rule_deploy fabric
    cd "${CURRENT_PATH}"
    FABRIC_CONFIG_PATH="${PIER_ROOT}"/fabric
    x_replace "s:\${CONFIG_PATH}:$FABRIC_CONFIG_PATH:g" "${PIER_ROOT}"/fabric/config.yaml
  fi

  if [ "$MODE" == "ethereum" ]; then
    print_blue "===> Copy ethereum plugins"
    cp "${CURRENT_PATH}"/pier-client-ethereum/build/eth-client.so "${PIER_ROOT}"/plugins/
    print_blue "===> Register pier(ethereum) to bitxhub"
    appchain_register chainB ether chainB-description 1.0 ethereum/ether.validators
    print_blue "===> Deploy rule in bitxhub"
    rule_deploy ethereum
  fi

  print_blue "===> Start pier..."
  nohup pier --repo "${PIER_ROOT}" start >/dev/null 2>&1 &
  echo $! >>"${PIER_ROOT}/pier-${MODE}.pid"
  print_green "===> Start pier successfully!!!"
}

function pier_down() {
  set +e
  print_blue "===> Kill $MODE pier"

  cd "${CURRENT_PATH}"/.pier_$MODE
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

PIER_ROOT="${CURRENT_PATH}"/.pier_fabric
BITXHUB_ADDR="localhost:60011"
MODE="fabric"

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
