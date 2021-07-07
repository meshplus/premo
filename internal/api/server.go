package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"
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

	addr, err := g.adminKey.PublicKey().Address()
	if err != nil {
		return err
	}

	nonce, err := g.client.GetPendingNonceByAccount(addr.String())
	if err != nil {
		return err
	}

	data := &pb.TransactionData{
		Type:   pb.TransactionData_NORMAL,
		VmType: pb.TransactionData_XVM,
		Amount: 100000,
	}

	payload, err := data.Marshal()
	if err != nil {
		return err
	}

	tx := &pb.Transaction{
		From:      addr,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}

	for i := 0; i < g.config.Concurrent; i++ {
		bee, err := bitxhub.NewBee(g.config.TPS, nil, nil, 0, g.config, context.TODO())
		if err != nil {
			return err
		}

		tx.To = bee.GetAddress()

		_, err = g.client.SendTransaction(tx, &rpcx.TransactOpts{Nonce: nonce})
		if err != nil {
			return err
		}
		nonce++

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
				for _, tx := range block.Transactions {
					go func(tx *pb.Transaction) {
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
					}(tx)
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
	typ := c.Query("type")
	if typ == "" {
		typ = "transfer"
	}

	key := ""
	if typ == "setData" || typ == "getData" {
		key = c.Query("key")
	}

	bee := <-g.beeC

	if strings.Compare(typ, "doubleSpend") == 0 {
		receipts, err := bee.SendDoubleSpendTxs()
		g.beeC <- bee
		if err != nil {
			g.logger.Errorf("doubleSpend get error: %v", err)
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if len(receipts) != 2 {
			g.logger.Errorf("doubleSpend get receipt size: %v", len(receipts))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if receipts[0].Status == receipts[1].Status {
			g.logger.Errorf("doubleSpend get receipt: %v, %v", receipts[0].Status, receipts[1].Status)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.Status(http.StatusOK)
	} else {
		hash, err := bee.SendTx(typ, key, 1)
		g.beeC <- bee
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		g.waitForConfirm(hash)

		c.Status(http.StatusOK)
	}
	g.logger.Infof("send tx and get receipt costs %d", time.Since(start).Milliseconds())
}

func (g *Server) waitForConfirm(txHash string) {
	ch := make(chan struct{})
	g.txMap.Store(txHash, ch)

	<-ch
}
