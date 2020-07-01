# Premo

![build](https://github.com/meshplus/premo/workflows/build/badge.svg)[![Go Report Card](https://goreportcard.com/badge/github.com/meshplus/premo)](https://goreportcard.com/report/github.com/meshplus/premo)

BitXHub interchain transaction testing framework.

## Quick Start

### Installation

```shell
git clone git@github.com:meshplus/premo.git
cd premo
make install
```

### Initialization

```shell
premo init
```

It will create `~/.premo` directory on you computer.

### Start Premo

```shell
premo interchain up
```

It will start the following things: 

+ a ethereum chain and a fabric chain (both in the form of docker container)
+ a fabric pier and a ethereum pier
+ deploy necessary contracts on fabric (ethererum chain image was already deployed)

### Do Testing

```shell
make tester
```

## Usage

```shell
premo [global options] command [command options] [arguments...]
```

### command

+ `init`        init config home for premo
+ `version`     Premo version
+ `test`        test bitxhub function
+ `pier`        Start or stop the pier
+ `bitxhub`     Start or stop the bitxhub cluster
+ `appchain`    Bring up the appchain network
+ `interchain`  Start or Stop the interchain system
+ `help, h`     Shows a list of commands or help for one command

### global options

+ `--repo value`  Premo storage repo path
+ `--help, -h`    show help (default: false)


