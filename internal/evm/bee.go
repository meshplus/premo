package evm

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	eth "github.com/meshplus/go-eth-client"
	"github.com/sirupsen/logrus"
)

type Bee struct {
	typ    string
	tps    int
	client *eth.EthRPC
	pk     *ecdsa.PrivateKey
	ctx    context.Context
	nonce  uint64
}

func NewBee(config *Config) (*Bee, error) {
	client, err := eth.New(eth.WithUrls([]string{config.JsonRpc}))
	if err != nil {
		return nil, err
	}
	pk, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	addr := crypto.PubkeyToAddress(pk.PublicKey)
	nonce, err := client.EthGetTransactionCount(addr, nil)
	if err != nil {
		return nil, err
	}
	return &Bee{
		typ:    config.Typ,
		client: client,
		pk:     pk,
		tps:    config.TPS / config.Concurrent,
		ctx:    config.Ctx,
		nonce:  nonce,
	}, nil
}

func (bee *Bee) Start() error {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			for i := 0; i < bee.tps; i++ {
				go func(nonce uint64) {
					err := bee.SendTx(nonce)
					if err != nil {
						log.WithFields(logrus.Fields{
							"error": err.Error(),
						}).Info("Error send evm tx")
					}
					atomic.AddInt64(&delayer, 1)
				}(bee.nonce)
				bee.nonce = bee.nonce + 1
			}
		case <-bee.ctx.Done():
			return nil
		}
	}
}

func (bee *Bee) Stop() error {
	return nil
}

func (bee *Bee) SendTx(nonce uint64) error {
	switch bee.typ {
	case "deploy":
		err := bee.DeployContract(nonce)
		if err != nil {
			return err
		}
	case "invoke":
		err := bee.Invoke(nonce)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unexpected type: %v", bee.typ)
	}
	return nil
}

func (bee *Bee) DeployContract(nonce uint64) error {
	if compileResult == nil {
		return fmt.Errorf("no compile result")
	}
	res, err := bee.client.Deploy(bee.pk, compileResult, args, eth.WithNonce(nonce))
	if err != nil {
		return err
	}
	fmt.Println(res)
	return nil
}

func (bee *Bee) Invoke(nonce uint64) error {
	_, err := bee.client.Invoke(bee.pk, &contractAbi, address, function, args, eth.WithNonce(nonce))
	if err != nil {
		return err
	}
	return nil
}
