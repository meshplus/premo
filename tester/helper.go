package tester

import (
	"context"
	"net/rpc"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/meshplus/bitxhub-model/pb"
)

type EthPier struct {
	abi       abi.ABI
	config    *Config
	ctx       context.Context
	ethClient *ethclient.Client
	broker    *Broker
	session   *BrokerSession
	conn      *rpc.Client
	eventC    chan *pb.IBTP
	pierID    string
}
