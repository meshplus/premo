set -e

source x.sh

CURRENT_PATH=$(pwd)
OPT=$1
MODE=$2
N=$3
VERSION=$4

function printHelp() {
  print_blue "Usage:  "
  echo "  run_bitxhub.sh <OPT>"
  echo "    <OPT> - one of 'up', 'down', 'restart'"
  echo "      - 'up' - bring up the bitxhub network"
  echo "      - 'down' - clear the bitxhub network"
  echo "  run_bitxhub.sh -h (print this message)"
}

function prepare() {
  cd "${CURRENT_PATH}"
  if ! type goduck >/dev/null 2>&1; then
    print_blue "===> Install goduck"
    go get github.com/meshplus/goduck
  fi

  if [ "$MODE" == "solo" ]; then
    print_blue "===> Generate bitxhub solo configure"
    goduck bitxhub config --mode solo --target "${CURRENT_PATH}"/bitxhub_config
  fi
  if [ "$MODE" == "cluster" ]; then
    print_blue "===> Generate bitxhub cluster configure"
    goduck bitxhub config --num "${N}" --target "${CURRENT_PATH}"/bitxhub_config
  fi

  if [ ! -d "bitxhub" ]; then
    git clone https://github.com/meshplus/bitxhub.git
  fi

  print_blue "===> Build bitxhub node"
  cd bitxhub
  git checkout -f master && git reset --hard HEAD
  git pull
  if [ -n "${VERSION}" ]; then
    print_blue "git checkout ${VERSION}"
    git checkout "${VERSION}"
  fi

  make prepare && make build
  cd internal/plugins && make plugins
}

function bitxhub_up() {
  prepare

  cd "${CURRENT_PATH}"/bitxhub_config

  if [ "$MODE" == "solo" ]; then
    if [ ! -d nodeSolo/plugins ]; then
      mkdir nodeSolo/plugins
      cp -r "${CURRENT_PATH}"/bitxhub/internal/plugins/build/solo.so nodeSolo/plugins
    fi
    echo "Start bitxhub solo node"
    nohup "${CURRENT_PATH}"/bitxhub/bin/bitxhub --repo="${CURRENT_PATH}"/bitxhub_config/nodeSolo start &
    echo $! >>"${CURRENT_PATH}"/bitxhub_config/bitxhub.pid
  fi

  if [ "$MODE" == "cluster" ]; then
    for ((i = 1; i < N + 1; i = i + 1)); do
      if [ ! -d node${i}/plugins ]; then
        mkdir node${i}/plugins
        cp -r "${CURRENT_PATH}"/bitxhub/internal/plugins/build/raft.so node${i}/plugins
      fi
      echo "Start bitxhub node${i}"
      nohup "${CURRENT_PATH}"/bitxhub/bin/bitxhub --repo="${CURRENT_PATH}"/bitxhub_config/node${i} start &
      echo $! >>"${CURRENT_PATH}"/bitxhub_config/bitxhub.pid
    done
  fi

}

function bitxhub_down() {
  set +e
  cd "${CURRENT_PATH}"/bitxhub_config
  if [ -a bitxhub.pid ]; then
    list=$(cat bitxhub.pid)
    for pid in $list; do
      kill "$pid"
      if [ $? -eq 0 ]; then
        echo "node pid:$pid exit"
      else
        print_red "program exit fail, try use kill -9 $pid"
      fi
    done
    rm bitxhub.pid
  fi
}

if [ "$OPT" == "up" ]; then
  bitxhub_up
elif [ "$OPT" == "down" ]; then
  bitxhub_down
else
  printHelp
  exit 1
fi
