package jsonrpc_tester

import (
	"testing"

	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/premo/internal/repo"

	"github.com/sirupsen/logrus"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var cfg = &config{
	addrs: []string{
		"localhost:60011",
		"localhost:60012",
		"localhost:60013",
		"localhost:60014",
	},
	logger: logrus.New(),
}

type config struct {
	addrs  []string
	logger rpcx.Logger
}
type Client struct {
	suite.Suite
	rpcClient *rpc.Client
	EthClient *ethclient.Client
	client    rpcx.Client
}

var host = "http://localhost:8881"

func TestTester(t *testing.T) {
	rpcClient, err := rpc.Dial(host)
	require.Nil(t, err)
	ethClient := ethclient.NewClient(rpcClient)
	require.Nil(t, err)

	keyPath, err := repo.Node1Path()
	require.Nil(t, err)
	pk, err := asym.RestorePrivateKey(keyPath, repo.KeyPassword)
	require.Nil(t, err)
	node0 := &rpcx.NodeInfo{Addr: cfg.addrs[0]}
	client, err := rpcx.New(
		rpcx.WithNodesInfo(node0),
		rpcx.WithLogger(cfg.logger),
		rpcx.WithPrivateKey(pk),
	)
	suite.Run(t, &Client{rpcClient: rpcClient, EthClient: ethClient, client: client})
}
