package api

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Rican7/retry"
	"github.com/Rican7/retry/strategy"
	"github.com/gin-gonic/gin"
	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/meshplus/premo/internal/bitxhub"
	"github.com/meshplus/premo/internal/repo"
	"github.com/sirupsen/logrus"
)

type Server struct {
	router   *gin.Engine
	logger   logrus.FieldLogger
	port     uint64
	config   *bitxhub.Config
	client   *rpcx.ChainClient
	txMap    sync.Map
	beeC     chan *bitxhub.Bee
	dBeeC    chan *bitxhub.Bee
	adminKey crypto.PrivateKey

	ctx    context.Context
	cancel context.CancelFunc
}

func NewServer(port uint64, config *bitxhub.Config, logger logrus.FieldLogger) (*Server, error) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	adminPk, err := asym.RestorePrivateKey(config.KeyPath, repo.KeyPassword)
	if err != nil {
		return nil, err
	}

	node0 := &rpcx.NodeInfo{Addr: config.BitxhubAddr[0]}
	client, err := rpcx.New(
		rpcx.WithNodesInfo(node0),
		rpcx.WithLogger(logger),
		rpcx.WithPrivateKey(adminPk),
	)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Server{
		router:   router,
		logger:   logger,
		port:     port,
		config:   config,
		client:   client,
		beeC:     make(chan *bitxhub.Bee, config.Concurrent),
		adminKey: adminPk,
		ctx:      ctx,
		cancel:   cancel,
	}, nil
}

func (g *Server) Start() error {
	g.router.Use(gin.Recovery())
	v1 := g.router.Group("/v1")
	v1.GET("sendTx", g.sendTx)
	g.logger.Infof("start prepare client")

	adminAdress, err := g.adminKey.PublicKey().Address()
	if err != nil {
		return err
	}
	bitxhub.AdminNonce, err = g.client.GetPendingNonceByAccount(adminAdress.String())
	if err != nil {
		return err
	}

	//query nodes nonce
	node1, err := repo.Node1Path()
	if err != nil {
		return err
	}
	key, err := asym.RestorePrivateKey(node1, repo.KeyPassword)
	if err != nil {
		return err
	}
	address, err := key.PublicKey().Address()
	if err != nil {
		return err
	}
	bitxhub.Index1, err = g.client.GetPendingNonceByAccount(address.String())

	node2, err := repo.Node2Path()
	if err != nil {
		return err
	}
	key, err = asym.RestorePrivateKey(node2, repo.KeyPassword)
	if err != nil {
		return err
	}
	address, err = key.PublicKey().Address()
	if err != nil {
		return err
	}
	bitxhub.Index2, err = g.client.GetPendingNonceByAccount(address.String())

	node3, err := repo.Node3Path()
	if err != nil {
		return err
	}
	key, err = asym.RestorePrivateKey(node3, repo.KeyPassword)
	if err != nil {
		return err
	}
	address, err = key.PublicKey().Address()
	if err != nil {
		return err
	}
	bitxhub.Index3, err = g.client.GetPendingNonceByAccount(address.String())

	bitxhub.Index1 -= 1
	bitxhub.Index2 -= 1
	bitxhub.Index3 -= 1

	// query pending nonce for adminKey
	bitxhub.AdminNonce, err = g.client.GetPendingNonceByAccount(adminAdress.String())

	var wg sync.WaitGroup
	wg.Add(g.config.Concurrent)

	for i := 0; i < g.config.Concurrent; i++ {
		go func() {
			defer wg.Done()

			bee, err := bitxhub.NewBee(g.config.TPS, g.adminKey, adminAdress, g.config, context.TODO())
			if err != nil {
				g.logger.Errorf("newBee err: %s", err)
				return
			}

			if err := bee.PrepareChain(g.config.Appchain, "检查链", g.config.Validator, "1.4.4", "fabric for law", g.config.Rule); err != nil {
				g.logger.Errorf("register appchain err: %s", err)
				return
			}
			g.beeC <- bee
		}()
	}

	wg.Wait()
	g.logger.WithFields(logrus.Fields{
		"number": g.config.Concurrent,
	}).Info("generate all bees")

	go func() {
		err := g.router.Run(fmt.Sprintf(":%d", g.port))
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		ch, err := g.client.Subscribe(context.TODO(), pb.SubscriptionRequest_BLOCK, nil)
		if err != nil {
			g.logger.WithField("error", err).Error("subscribe block")
			return
		}
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-g.ctx.Done():
				return
			case data, ok := <-ch:
				if !ok {
					g.logger.Warn("block subscription channel is closed")
					return
				}

				block := data.(*pb.Block)
				for _, tx := range block.Transactions.Transactions {
					bxhTx := tx.(*pb.BxhTransaction)
					go func(tx *pb.BxhTransaction) {
						var (
							ch chan struct{}
						)
						hash := tx.Hash().String()
						if err := retry.Retry(func(attempt uint) error {
							val, ok := g.txMap.Load(hash)
							if !ok {
								g.logger.Warnf("tx %s not found in map", hash)
								return fmt.Errorf("tx not found")
							}
							ch = val.(chan struct{})
							g.txMap.Delete(hash)
							return nil
						}, strategy.Wait(time.Millisecond*10), strategy.Limit(3)); err != nil {
							return
						}

						ch <- struct{}{}
					}(bxhTx)
				}
			}
		}
	}()
	<-g.ctx.Done()

	return nil
}

func (g *Server) Stop() error {
	g.cancel()
	g.logger.Infoln("gin service stop")
	return nil
}

func (g *Server) sendTx(c *gin.Context) {
	start := time.Now()
	var isGet bool
	typ := c.Query("type")
	if typ == "" {
		typ = "transfer"
	}

	key := ""
	if typ == "setData" || typ == "getData" {
		key = c.Query("key")
		if typ == "getData" {
			isGet = true
		}
	}

	bee := <-g.beeC

	hash, err := bee.SendTx(typ, 1, key, isGet)
	if err != nil {
		g.logger.Errorf("sendTx err: %s", err)
	}
	g.beeC <- bee

	g.waitForConfirm(hash)

	c.Status(http.StatusOK)
	g.logger.Infof("send tx and get receipt costs %d", time.Since(start).Milliseconds())
}

func (g *Server) waitForConfirm(txHash string) {
	ch := make(chan struct{})
	g.txMap.Store(txHash, ch)

	<-ch
}
