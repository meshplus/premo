# Premo
BitXHub interchain transaction testing framework

## Install

Install binary:

```shell
make install
```

## Start the system

In order to test the cross-chain system, first start it. If you want to start the system by hand...

Run bitxhub(relay chain):

```shell
cd scripts
bash run_bitxhub.sh up cluster 4 master

## shut down:
bash run_bitxhub.sh down
```

Run appchains(one ethereum chain and one fabric chain)：

```shell
# fabric appchain
bash run_appchain.sh up fabric
# ethereum appchain
bash run_appchain.sh up ethereum

## shut down:
bash run_appchain.sh down fabric
bash run_appchain.sh down ethereum
```

Run piers(one ethereum pier and one fabric pier):

```shell
# fabric pier
bash run_pier.sh up -t fabric -v master -r '.pier_fabric' -b 'localhost:60011'
# ethereum pier
bash run_pier.sh up -t ethereum -v master -r '.pier_ethereum' -b 'localhost:60011'

## shut down:
bash run_pier.sh down -t fabric
bash run_pier.sh down -t ethereum
```

## Test the system

See `premo -h` for details.

Initialize premo:

```shell
premo init
```

it will create `~/.premo` dir.