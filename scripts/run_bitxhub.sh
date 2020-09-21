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
  echo "      - 'down' - shut down the bitxhub network"
  echo "  run_bitxhub.sh -h (print this message)"
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

  if [ -f "${CURRENT_PATH}"/.bitxhub/bitxhub.pid ]; then
    print_red "bitxhub is already running in the background service"
    cat "${CURRENT_PATH}"/.bitxhub/bitxhub.pid
    exit 0
  fi

  if [ "$MODE" == "solo" ]; then
    print_blue "===> Generate bitxhub solo configure"
    goduck bitxhub config --mode solo --target "${CURRENT_PATH}"/.bitxhub --version "${VERSION}"
  fi
  if [ "$MODE" == "cluster" ]; then
    print_blue "===> Generate bitxhub cluster configure"
    goduck bitxhub config --num "${N}" --target "${CURRENT_PATH}"/.bitxhub --version "${VERSION}"
  fi

  if [ ! -d "bitxhub" ]; then
    git clone https://github.com/meshplus/bitxhub.git
  fi

  cd bitxhub
  git checkout -f master && git reset --hard HEAD
  git pull
  if [ -n "${VERSION}" ]; then
    print_blue "git checkout ${VERSION}"
    git checkout "${VERSION}"
  fi

  print_blue "===> Build bitxhub node"
  make prepare && make install
  cd internal/plugins && make plugins
}

function bitxhub_up() {
  prepare

  cd "${CURRENT_PATH}"/.bitxhub

  if [ "$MODE" == "solo" ]; then
    if [ ! -d nodeSolo/plugins ]; then
      mkdir nodeSolo/plugins
      cp -r "${CURRENT_PATH}"/bitxhub/internal/plugins/build/solo.so nodeSolo/plugins
    fi

    nohup bitxhub --repo="${CURRENT_PATH}"/.bitxhub/nodeSolo start >/dev/null 2>&1 &
    PID=$!
    sleep 3
    if [ -n "$(ps -p ${PID} -o pid=)" ]; then
      echo "===> Start bitxhub solo node successful"
      echo $PID >>"${CURRENT_PATH}"/.bitxhub/bitxhub.pid
    else
      print_red "===> Start bitxhub solo node fail"
    fi

  fi

  if [ "$MODE" == "cluster" ]; then
    for ((i = 1; i < N + 1; i = i + 1)); do
      if [ ! -d node${i}/plugins ]; then
        mkdir node${i}/plugins
        cp -r "${CURRENT_PATH}"/bitxhub/internal/plugins/build/raft.so node${i}/plugins
      fi

      nohup bitxhub --repo="${CURRENT_PATH}"/.bitxhub/node${i} start >/dev/null 2>&1 &
      PID=$!
      sleep 3
      if [ -n "$(ps -p ${PID} -o pid=)" ]; then
        echo "===> Start bitxhub node${i} successful"
        echo $PID >>"${CURRENT_PATH}"/.bitxhub/bitxhub.pid
      else
        print_red "===> Start bitxhub node${i} fail"
      fi
    done
  fi

}

function bitxhub_down() {
  set +e
  cd "${CURRENT_PATH}"/.bitxhub
  if [ -a bitxhub.pid ]; then
    list=$(cat bitxhub.pid)
    for pid in $list; do
      kill "$pid"
      if [ $? -eq 0 ]; then
        echo "node pid:$pid exit"
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