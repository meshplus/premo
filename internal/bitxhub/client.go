package bitxhub

import (
	"github.com/meshplus/bitxhub-kit/key"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/sirupsen/logrus"
)

func NewBxhClient(addr, path string, logger *logrus.Entry) (*rpcx.ChainClient, error) {
	k, err := key.LoadKey(path)

	if err != nil {
		return nil, err
	}

	privKey, err := k.GetPrivateKey("bitxhub")
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
