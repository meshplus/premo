package evm

import (
	"context"
	"fmt"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"

	eth "github.com/meshplus/go-eth-client"
)

type Bee struct {
	typ     string
	tps     int
	jsonRpc string
	ctx     context.Context
}

func NewBee(config *Config) (*Bee, error) {
	return &Bee{
		typ:     config.Typ,
		jsonRpc: config.JsonRpc,
		tps:     config.TPS / config.Concurrent,
		ctx:     config.Ctx,
	}, nil
}

func (b *Bee) Start() error {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			for i := 0; i < b.tps; i++ {
				go func() {
					err := b.SendTx(0)
					if err != nil {
						log.WithFields(logrus.Fields{
							"error": err.Error(),
						}).Info("Error send evm tx")
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
	client, err := NewClient(b.jsonRpc)
	if err != nil {
		return err
	}
	_, err = client.Deploy(compileResult, args, eth.WithNonce(big.NewInt(nonce)))
	if err != nil {
		return err
	}
	return nil
}

func (b *Bee) Invoke(nonce int64) error {
	client, err := NewClient(b.jsonRpc)
	if err != nil {
		return err
	}
	_, err = client.Invoke(contractAbi, address, function, args, eth.WithTxNonce(uint64(nonce)))
	if err != nil {
		return err
	}
	return nil
}
