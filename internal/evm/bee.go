package evm

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"sync/atomic"
	"time"

	eth "github.com/meshplus/go-eth-client"
	"github.com/sirupsen/logrus"
)

type Bee struct {
	client *eth.EthRPC
	pk     *ecdsa.PrivateKey
	typ    string
	tps    int
	ctx    context.Context
}

func NewBee(config *Config) (*Bee, error) {
	client, pk, err := NewClient(config)
	if err != nil {
		return nil, err
	}
	return &Bee{
		client: client,
		pk:     pk,
		typ:    config.Typ,
		tps:    config.TPS / config.Concurrent,
		ctx:    config.Ctx,
	}, nil
}

func (b *Bee) Start() error {
	ticker := time.NewTicker(time.Second)
	nonce := int64(-1)
	for {
		select {
		case <-ticker.C:
			for i := 0; i < b.tps; i++ {
				go func() {
					txNonce := atomic.AddInt64(&nonce, 1)
					err := b.SendTx(txNonce)
					if err != nil {
						log.WithFields(logrus.Fields{
							"error": err.Error(),
						}).Info("Error send evm tx")
						return
					}
					atomic.AddInt64(&delayer, 1)
				}()
			}
		case <-b.ctx.Done():
			return nil
		}
	}
}

func (b *Bee) Stop() error {
	return nil
}

func (b *Bee) SendTx(nonce int64) error {
	switch b.typ {
	case "deploy":
		err := b.DeployContract(nonce)
		if err != nil {
			return err
		}
	case "invoke":
		err := b.Invoke(nonce)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unexpected type: %v", b.typ)
	}
	return nil
}

func (b *Bee) DeployContract(nonce int64) error {
	if compileResult == nil {
		return fmt.Errorf("no compile result")
	}
	_, err := b.client.Deploy(compileResult, args, eth.WithNonce(big.NewInt(nonce)))
	if err != nil {
		return err
	}
	return nil
}

func (b *Bee) Invoke(nonce int64) error {
	_, err := b.client.Invoke(contractAbi, address, function, args, eth.WithTxNonce(uint64(nonce)))
	if err != nil {
		return err
	}
	return nil
}
