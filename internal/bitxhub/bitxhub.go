package bitxhub

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/meshplus/premo/internal/repo"
	"github.com/sirupsen/logrus"
	"github.com/wonderivan/logger"
)

const (
	DefaultTo = "000000000000000000000000000000000000000a"
)

var index1 uint64
var index2 uint64
var index3 uint64
var log = logrus.New()
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
	index1 -= 1
	index2 -= 1
	index3 -= 1

	// query pending nonce for adminKey
	adminNonce, err := client.GetPendingNonceByAccount(adminFrom.String())
	if err != nil {
		return nil, err
	}

	bees := make([]*bee, 0, config.Concurrent)

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(config.Concurrent)

	for i := 0; i < config.Concurrent; i++ {
		go func() {
			defer wg.Done()

			expectedNonce := atomic.AddUint64(&adminNonce, 1) - 1
			bee, err := NewBee(config.TPS/config.Concurrent, adminPk, adminFrom, expectedNonce, config, ctx)
			if err != nil {
				logger.Error("New bee: ", err.Error())
				return
			}
			if config.Type == "interchain" {
				if err := bee.prepareChain(config.Appchain, "检查链", config.Validator, "1.4.4", "fabric for law", config.Rule); err != nil {
					logger.Error(err)
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
				for _, tx := range block.Transactions {
					cnt++
					counter++

					txDelay := now - tx.Timestamp
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
		broker.bees[i].stop()
	}
	err := broker.client.Stop()
	if err != nil {
		panic(err)
	}
	avgDelay := float64(delayer) / float64(counter)
	duration := time.Since(current).Seconds()
	log.WithFields(logrus.Fields{
		"number":    counter,
		"duration":  duration,
		"tps":       float64(counter) / duration,
		"avg delay": avgDelay / float64(time.Millisecond),
		"max delay": float64(maxDelay) / float64(time.Millisecond),
	}).Info("finish testing")
	return nil
}
