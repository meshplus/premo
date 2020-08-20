package bitxhub

import (
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/sirupsen/logrus"
)

func NewBxhClient(addr, path string, logger *logrus.Entry) (*rpcx.ChainClient, error) {
	privKey, err := asym.RestorePrivateKey(path, "bitxhub")
	if err != nil {
		return nil, err
	}

	cli, err := rpcx.New(
		rpcx.WithAddrs([]string{addr}),
		rpcx.WithLogger(logger),
		rpcx.WithPrivateKey(privKey),
	)
	if err != nil {
		return nil, err
	}
	return cli, nil
}
