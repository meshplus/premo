package bxh_tester

import (
	"testing"

	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/meshplus/premo/internal/repo"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestTester(t *testing.T) {
	keyPath, err := repo.Node1Path()
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
	var clients []*rpcx.ChainClient
	for i := 0; i < 10; i++ {
		client, err := rpcx.New(
			rpcx.WithNodesInfo(node0),
			rpcx.WithLogger(cfg.logger),
			rpcx.WithPrivateKey(pk),
		)
		require.Nil(t, err)
		clients = append(clients, client)
	}
	if len(clients) == 10 {
		suite.Run(t, &Model1{&Snake{client: clients[0], from: from, pk: pk, to: to}})
		suite.Run(t, &Model2{&Snake{client: clients[1], from: from, pk: pk, to: to}})
		suite.Run(t, &Model3{&Snake{client: clients[2], from: from, pk: pk, to: to}})
		suite.Run(t, &Model4{&Snake{client: clients[3], from: from, pk: pk, to: to}})
		suite.Run(t, &Model5{&Snake{client: clients[4], from: from, pk: pk, to: to}})
		suite.Run(t, &Model6{&Snake{client: clients[5], from: from, pk: pk, to: to}})
		suite.Run(t, &Model7{&Snake{client: clients[6], from: from, pk: pk, to: to}})
		suite.Run(t, &Model8{&Snake{client: clients[7], from: from, pk: pk, to: to}})
		suite.Run(t, &Model9{&Snake{client: clients[8], from: from, pk: pk, to: to}})
	}
}
