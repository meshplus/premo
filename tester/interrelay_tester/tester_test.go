package interrelay_tester

import (
	"testing"

	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/meshplus/premo/internal/repo"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var cfg = &config{
	addrs: []string{
		"localhost:60011",
	},
	logger: logrus.New(),
}

type config struct {
	addrs  []string
	logger rpcx.Logger
}

func TestTester(t *testing.T) {
	keyPath, err := repo.KeyPath()
	require.Nil(t, err)

	pk, err := asym.RestorePrivateKey(keyPath, repo.KeyPassword)
	require.Nil(t, err)

	from, err := pk.PublicKey().Address()
	require.Nil(t, err)

	pk1, err := asym.GenerateKeyPair(crypto.Secp256k1)
	require.Nil(t, err)

	to, err := pk1.PublicKey().Address()
	require.Nil(t, err)
	node0 := &rpcx.NodeInfo{Addr: cfg.addrs[0]}
	//node1 := &rpcx.NodeInfo{Addr: cfg.addrs[1]}
	//node2 := &rpcx.NodeInfo{Addr: cfg.addrs[2]}
	//node3 := &rpcx.NodeInfo{Addr: cfg.addrs[3]}
	client, err := rpcx.New(
		rpcx.WithNodesInfo(node0),
		rpcx.WithLogger(cfg.logger),
		rpcx.WithPrivateKey(pk),
	)
	require.Nil(t, err)

	suite.Run(t, &Snake{client: client, from: from, pk: pk, to: to})
}
