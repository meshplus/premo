package api

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/meshplus/premo/internal/bitxhub"
	"github.com/sirupsen/logrus"
)

type Server struct {
	router *gin.Engine
	logger logrus.FieldLogger
	port   uint64
	config *bitxhub.Config
	client *rpcx.ChainClient
	txMap  sync.Map

	beeC chan *bitxhub.Bee

	ctx    context.Context
	cancel context.CancelFunc
}

func NewServer(port uint64, config *bitxhub.Config, logger logrus.FieldLogger) (*Server, error) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	if err != nil {
		return nil, err
	}

	node0 := &rpcx.NodeInfo{Addr: config.BitxhubAddr[0]}
	client, err := rpcx.New(
		rpcx.WithNodesInfo(node0),
		rpcx.WithLogger(logger),
		rpcx.WithPrivateKey(pk),
	)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Server{
		router: router,
		logger: logger,
		port:   port,
		config: config,
		client: client,
		beeC:   make(chan *bitxhub.Bee, config.Concurrent),
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

func (g *Server) Start() error {
	g.router.Use(gin.Recovery())
	v1 := g.router.Group("/v1")
	{
		v1.GET("sendTx", g.sendTx)
	}

	for i := 0; i < g.config.Concurrent; i++ {
		bee, err := bitxhub.NewBee(g.config.TPS, nil, nil, 0, g.config, context.TODO())
		if err != nil {
			return err
		}

		g.beeC <- bee
	}

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
				for i := 0; i < len(block.Transactions); {
					hash := block.Transactions[i].Hash().String()
					val, ok := g.txMap.LoadAndDelete(hash)
					if !ok {
						g.logger.Warnf("tx %s not found in map", hash)
						time.Sleep(time.Millisecond * 10)
						continue
					}

					ch := val.(chan struct{})
					ch <- struct{}{}
					i++
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
	typ := c.Query("type")
	if typ == "" {
		typ = "transfer"
	}

	bee := <-g.beeC

	txHash, err := bee.SendTx(typ, 1)
	if err != nil {
		return
	}
	g.beeC <- bee

	g.waitForConfirm(txHash)
}

func (g *Server) waitForConfirm(txHash string) {
	ch := make(chan struct{})
	g.txMap.Store(txHash, ch)

	<-ch
}
