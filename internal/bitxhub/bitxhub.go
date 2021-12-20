package bitxhub

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	appchain_mgr "github.com/meshplus/bitxhub-core/appchain-mgr"
	"github.com/meshplus/bitxhub-core/governance"

	"github.com/meshplus/bitxhub-model/constant"

	"github.com/meshplus/bitxhub-model/pb"

	"github.com/meshplus/bitxhub-kit/types"

	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/meshplus/premo/internal/repo"
	"github.com/sirupsen/logrus"
	"github.com/wonderivan/logger"
)

var index1 uint64
var index2 uint64
var index3 uint64
var adminNonce uint64
var log = logrus.New()
var To string
var cfg = &config{
	addrs: []string{
		"localhost:60011",
	},
	logger: log,
}

type config struct {
	addrs  []string
	logger rpcx.Logger
}

type Broker struct {
	config     *Config
	bees       []*bee
	client     rpcx.Client
	adminNonce uint64
	ctx        context.Context
	cancel     context.CancelFunc
}

type Config struct {
	Concurrent  int
	TPS         int
	Duration    int // s uint
	Type        string
	Validator   string
	Proof       []byte
	Rule        []byte
	KeyPath     string
	BitxhubAddr []string
	Appchain    string
}

func New(config *Config) (*Broker, error) {
	log.WithFields(logrus.Fields{
		"concurrent": config.Concurrent,
		"tps":        config.TPS,
		"duration":   config.Duration,
		"type":       config.Type,
	}).Info("Premo configuration")
	adminPk, err := asym.RestorePrivateKey(config.KeyPath, repo.KeyPassword)
	if err != nil {
		return nil, err
	}

	adminFrom, err := adminPk.PublicKey().Address()
	if err != nil {
		return nil, err
	}

	node0 := &rpcx.NodeInfo{Addr: config.BitxhubAddr[0]}
	client, err := rpcx.New(
		rpcx.WithNodesInfo(node0),
		rpcx.WithLogger(cfg.logger),
		rpcx.WithPrivateKey(adminPk),
	)
	if err != nil {
		return nil, err
	}

	//query nodes nonce
	node1, err := repo.Node1Path()
	if err != nil {
		return nil, err
	}
	key, err := asym.RestorePrivateKey(node1, repo.KeyPassword)
	if err != nil {
		return nil, err
	}
	address, err := key.PublicKey().Address()
	if err != nil {
		return nil, err
	}
	index1, err = client.GetPendingNonceByAccount(address.String())
	if err != nil {
		return nil, err
	}

	node2, err := repo.Node2Path()
	if err != nil {
		return nil, err
	}
	key, err = asym.RestorePrivateKey(node2, repo.KeyPassword)
	if err != nil {
		return nil, err
	}
	address, err = key.PublicKey().Address()
	if err != nil {
		return nil, err
	}
	index2, err = client.GetPendingNonceByAccount(address.String())
	if err != nil {
		return nil, err
	}

	node3, err := repo.Node3Path()
	if err != nil {
		return nil, err
	}
	key, err = asym.RestorePrivateKey(node3, repo.KeyPassword)
	if err != nil {
		return nil, err
	}
	address, err = key.PublicKey().Address()
	if err != nil {
		return nil, err
	}
	index3, err = client.GetPendingNonceByAccount(address.String())
	if err != nil {
		return nil, err
	}
	index1 -= 1
	index2 -= 1
	index3 -= 1

	// query pending nonce for adminKey
	adminNonce, err = client.GetPendingNonceByAccount(adminFrom.String())
	if err != nil {
		return nil, err
	}
	// prepare to
	to, err := PrepareTo(config, adminPk, adminFrom)
	if err != nil {
		return nil, err
	}
	To = to

	bees := make([]*bee, 0, config.Concurrent)
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(config.Concurrent)
	for i := 0; i < config.Concurrent; i++ {
		go func() {
			defer wg.Done()
			bee, err := NewBee(config.TPS/config.Concurrent, adminPk, adminFrom, config, ctx)
			if err != nil {
				log.Error("New bee: ", err.Error())
				return
			}
			if config.Type == "interchain" {
				if err := bee.prepareChain(bee.config.Appchain, "fabric for law"); err != nil {
					log.Error(err)
					return
				}
			}

			bees = append(bees, bee)
		}()
	}

	wg.Wait()
	log.WithFields(logrus.Fields{
		"number": len(bees),
	}).Info("generate all bees")

	return &Broker{
		config:     config,
		bees:       bees,
		client:     client,
		adminNonce: adminNonce,
		ctx:        ctx,
		cancel:     cancel,
	}, nil
}

func (broker *Broker) Start(typ string) error {
	logger.Info("starting broker")
	var wg sync.WaitGroup
	wg.Add(len(broker.bees))

	current := time.Now()

	meta0, err := broker.client.GetChainMeta()
	if err != nil {
		return err
	}

	for i := 0; i < len(broker.bees); i++ {
		go func(i int) {
			wg.Done()
			err := broker.bees[i].start(typ)
			if err != nil {
				logger.Error(err)
				return
			}
			log.WithFields(logrus.Fields{
				"index": i + 1,
			}).Debug("start bee")
		}(i)
	}
	log.WithFields(logrus.Fields{
		"number": len(broker.bees),
	}).Info("start all bees")

	go func() {
		var (
			cnt  = int64(0)
			dly  = int64(0)
			mDly = int64(0)
		)
		ch, err := broker.client.Subscribe(context.TODO(), pb.SubscriptionRequest_BLOCK, nil)
		if err != nil {
			log.WithField("error", err).Error("subscribe block")
			return
		}
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-broker.ctx.Done():
				return
			case <-ticker.C:
				c := float64(cnt)
				d := float64(dly) / float64(time.Millisecond)
				md := float64(mDly) / float64(time.Millisecond)
				log.Infof("current tps is %d, average tx delay is %fms, max tx delay is %fms", cnt, d/c, md)

				if maxDelay < mDly {
					maxDelay = mDly
				}

				cnt = 0
				dly = 0
				mDly = 0

			case data, ok := <-ch:
				if !ok {
					log.Warn("block subscription channel is closed")
					return
				}

				block := data.(*pb.Block)
				now := time.Now().UnixNano()
				for _, tx := range block.Transactions.Transactions {
					cnt++
					counter++

					txDelay := now - tx.GetTimeStamp()
					dly += txDelay
					delayer += txDelay

					if mDly < txDelay {
						mDly = txDelay
					}
				}
			}
		}
	}()

	time.Sleep(time.Duration(broker.config.Duration) * time.Second)
	wg.Wait()

	_ = broker.Stop(current)

	meta1, err := broker.client.GetChainMeta()
	if err != nil {
		return err
	}
	logger.Info("Collecting tps info, please wait...")
	time.Sleep(20 * time.Second)

	skip := (meta1.Height - meta0.Height) / 8
	begin := meta0.Height + skip
	end := meta1.Height - skip
	tps, err := broker.client.GetTPS(begin, end)
	if err != nil {
		return err
	}
	log.Infof("the TPS from block %d to %d is %d", begin, end, tps)

	return nil
}

func (broker *Broker) Stop(current time.Time) error {
	broker.cancel()
	// wait for goroutines inside bees to stop
	time.Sleep(3 * time.Second)

	logger.Info("Bees are quiting, please wait...")
	for i := 0; i < len(broker.bees); i++ {
		err := broker.bees[i].stop()
		if err != nil {
			return err
		}
	}
	err := broker.client.Stop()
	if err != nil {
		log.Warn(err)
	}
	delayerAvg := float64(delayer) / float64(counter)
	log.WithFields(logrus.Fields{
		"number":   counter,
		"duration": time.Since(current).Seconds(),
		"tps":      float64(counter) / time.Since(current).Seconds(),
		"tx_delay": delayerAvg / float64(time.Millisecond),
	}).Info("finish testing")
	return nil
}

func PrepareTo(config *Config, adminPk crypto.PrivateKey, adminFrom *types.Address) (string, error) {
	node0 := &rpcx.NodeInfo{Addr: config.BitxhubAddr[0]}
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	if err != nil {
		log.Error(err)
	}
	from, err := pk.PublicKey().Address()
	if err != nil {
		log.Error(err)
	}
	client, err := rpcx.New(
		rpcx.WithNodesInfo(node0),
		rpcx.WithLogger(cfg.logger),
		rpcx.WithPrivateKey(pk),
	)
	if err != nil {
		return "", err
	}
	err = TransferFromAdmin(config, adminPk, adminFrom, from, "100")
	if err != nil {
		return "", err
	}
	args := []*pb.Arg{
		rpcx.String(from.String()),                                //chainID
		rpcx.String(from.String()),                                //chainName
		rpcx.String("Flato v1.0.3"),                               //chainType
		rpcx.Bytes([]byte("")),                                    //trustRoot
		rpcx.String("0x857133c5C69e6Ce66F7AD46F200B9B3573e77582"), //broker
		rpcx.String("desc"),                                       //desc
		rpcx.String("0x00000000000000000000000000000000000000a2"), //masterRuleAddr
		rpcx.String("https://github.com"),                         //masterRuleUrl
		rpcx.String(from.String()),                                //adminAddrs
		rpcx.String("reason"),                                     //reason
	}
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "RegisterAppchain", nil, args...)
	if err != nil {
		return "", err
	}
	//vote chain
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	if err != nil || result.ProposalID == "" {
		return "", fmt.Errorf("vote chain unmarshal error: %w", err)
	}
	b := bee{}
	b.config = config
	err = b.VotePass(result.ProposalID)
	if err != nil {
		return "", fmt.Errorf("vote chain error: %w", err)
	}
	res, err = b.GetChainStatusById(from.String())
	if err != nil {
		return "", fmt.Errorf("getChainStatus error: %w", err)
	}
	appchain := &appchain_mgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	if err != nil || appchain.Status != governance.GovernanceAvailable {
		return "", fmt.Errorf("chain error: %w", err)
	}
	//register server
	args = []*pb.Arg{
		rpcx.String(from.String()),
		rpcx.String("mychannel&transfer"),
		rpcx.String(from.String()),
		rpcx.String("CallContract"),
		rpcx.String("test"),
		rpcx.Uint64(1),
		rpcx.String(""),
		rpcx.String("test"),
		rpcx.String("reason"),
	}
	res, err = client.InvokeBVMContract(constant.ServiceMgrContractAddr.Address(), "RegisterService", nil, args...)
	if err != nil {
		return "", fmt.Errorf("register server error %w", err)
	}
	//vote server
	result = &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	if err != nil || result.ProposalID == "" {
		return "", fmt.Errorf("vote server unmarshal error: %w", err)
	}
	err = b.VotePass(result.ProposalID)
	if err != nil {
		return "", fmt.Errorf("vote server error: %w", err)
	}
	return from.String(), nil
}
