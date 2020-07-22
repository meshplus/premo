package bxh_tester

import (
	"testing"

	"github.com/meshplus/bitxhub-kit/key"
	"github.com/meshplus/premo/internal/repo"

	"github.com/meshplus/bitxhub-kit/types"

	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/sirupsen/logrus"
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

type Snake struct {
	suite.Suite
	client rpcx.Client
	from   types.Address
	pk     crypto.PrivateKey
	to     types.Address
}

func TestTester(t *testing.T) {
	keyPath, err := repo.KeyPath()
	require.Nil(t, err)

	key, err := key.LoadKey(keyPath)
	require.Nil(t, err)

	pk, err := key.GetPrivateKey(repo.KeyPassword)
	require.Nil(t, err)

	from, err := pk.PublicKey().Address()
	require.Nil(t, err)

	pk1, err := asym.GenerateKey(asym.ECDSASecp256r1)
	require.Nil(t, err)

	to, err := pk1.PublicKey().Address()
	require.Nil(t, err)

	client, err := rpcx.New(
		rpcx.WithAddrs(cfg.addrs),
		rpcx.WithLogger(cfg.logger),
		rpcx.WithPrivateKey(pk),
	)
	require.Nil(t, err)

	suite.Run(t, &Snake{
		client: client,
		from:   from,
		pk:     pk,
		to:     to,
	})
}
