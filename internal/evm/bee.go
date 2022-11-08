package evm

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	crypto2 "github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/types"
	rpcx "github.com/meshplus/go-bitxhub-client"
	eth "github.com/meshplus/go-eth-client"
	"github.com/meshplus/go-eth-client/utils"
	"github.com/meshplus/premo/internal/common"
	"github.com/sirupsen/logrus"
)

const (
	deploy         = "deploy"
	deployByCode   = "deployByCode"
	invoke         = "invoke"
	invokeWithByte = "invokeWithByte"
)

type Bee struct {
	typ    string
	tps    int
	client *eth.EthRPC
	pk     *ecdsa.PrivateKey
	ctx    context.Context
	nonce  int64
}

func NewBee(config *Config, grpcCli *rpcx.ChainClient, adminPk crypto2.PrivateKey, adminFrom *types.Address) (*Bee, error) {
	pk, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	addr := crypto.PubkeyToAddress(pk.PublicKey)
	if err != nil {
		return nil, err
	}
	err = common.TransferFromAdmin(grpcCli, adminPk, adminFrom, types.NewAddress(addr.Bytes()), "100")
	if err != nil {
		return nil, err
	}

	client, err := eth.New(config.JsonRpc)
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
		nonce:  int64(nonce),
	}, nil
}

func (b *Bee) Start() error {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			for i := 0; i < b.tps; i++ {
				go func(nonce int64) {
					err := b.SendTx(nonce)
					if err != nil {
						log.WithFields(logrus.Fields{
							"error": err.Error(),
						}).Info("Error send evm tx")
					}
					atomic.AddInt64(&delayer, 1)
				}(b.nonce)
				b.nonce = b.nonce + 1
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
	case deploy:
		return b.DeployContract(nonce)
	case deployByCode:
		return b.DeployContractByCode(nonce)
	case invoke:
		fallthrough
	case invokeWithByte:
		return b.Invoke(nonce)
	default:
		return fmt.Errorf("unexpected type: %v", b.typ)
	}
}

func (b *Bee) DeployContract(nonce int64) error {
	if compileResult == nil {
		return fmt.Errorf("no compile result")
	}
	_, err := b.client.Deploy(b.pk, compileResult, args, eth.WithNonce(uint64(nonce)))
	if err != nil {
		return err
	}
	return nil
}

func (b *Bee) DeployContractByCode(nonce int64) error {
	_, _, err := b.client.DeployByCode(b.pk, contractAbi, code, args, eth.WithNonce(uint64(nonce)))
	if err != nil {
		return err
	}
	return nil
}

func (b *Bee) Invoke(nonce int64) error {
	var err error
	if function == "publish" {
		type creditPackage struct {
			Credit   *big.Int
			Quantity uint8
			Duration *big.Int
		}
		type signStruct struct {
			HashedMessage [32]byte
			V             uint8
			R             [32]byte
			S             [32]byte
		}

		type publishStruct struct {
			DataId      *big.Int
			Publisher   string
			Prices      []creditPackage
			AuthList    []*big.Int
			SharingMode uint8
			Extra       string
			Sign        signStruct
		}
		input := publishStruct{
			DataId:    big.NewInt(nonce + 111301),
			Publisher: crypto.PubkeyToAddress(b.pk.PublicKey).String(),
			Prices: []creditPackage{{
				Credit:   big.NewInt(10),
				Quantity: 0,
				Duration: big.NewInt(1000),
			}},
			AuthList:    nil,
			SharingMode: 0,
			Extra:       "",
		}

		//log.Errorf("dataId:%d, bee:%s", nonce, crypto.PubkeyToAddress(b.pk.PublicKey))
		inputBytes, err := json.Marshal(input)
		if err != nil {
			return err
		}
		args, err = utils.DecodeBytes(abiEvent, "publish", inputBytes)
		if err != nil {
			return err
		}
	}
	_, err = b.client.Invoke(b.pk, &contractAbi, address, function, args, eth.WithNonce(uint64(nonce)))
	if err != nil {
		return err
	}
	return nil
}
