package server

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/meshplus/bitxhub-kit/types"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/meshplus/premo/internal/repo"
	"github.com/sirupsen/logrus"
)

const (
	NormalKey     = "key_for_normal"
	NormalValue   = "value_for_normal"
	HappyRuleAddr = "0x00000000000000000000000000000000000000a2"
)

type Server struct {
	remote     string
	port       int
	router     *gin.Engine
	clientPool []*Grpc
	log        *logrus.Logger
	toAddr     *types.Address
	hashMp     sync.Map
}

type Grpc struct {
	client  *rpcx.ChainClient
	address *types.Address
	nonce   uint64
	index   uint64
}

func NewServer(remote string, port, poolSize int) (*Server, error) {
	defaultRemote = remote
	err := initializeAdminNonce()
	if err != nil {
		return nil, err
	}
	pk, to, err := repo.KeyPriv()
	if err != nil {
		return nil, err
	}
	err = TransferFromAdmin(remote, to.String(), "100")
	if err != nil {
		return nil, err
	}
	err = prepareInterchain(pk)
	if err != nil {
		return nil, err
	}
	ibtpIdx := map[string]uint64{}
	clientPool := make([]*Grpc, poolSize)
	mutex := sync.Mutex{}
	for i := 0; i < poolSize; i++ {
		go func(i int) {
			pk, address, err := repo.KeyPriv()
			if err != nil {
				panic(err)
			}
			client, err := rpcx.New(
				rpcx.WithNodesInfo(&rpcx.NodeInfo{Addr: remote}),
				rpcx.WithPrivateKey(pk),
			)
			if err != nil {
				panic(err)
			}
			err = TransferFromAdmin(remote, address.String(), "100")
			if err != nil {
				panic(err)
			}
			err = prepareInterchain(pk)
			if err != nil {
				panic(err)
			}
			nonce, err := client.GetPendingNonceByAccount(address.String())
			if err != nil {
				panic(err)
			}
			clientPool[i] = &Grpc{client: client, address: address, nonce: nonce, index: 1}
			mutex.Lock()
			ibtpIdx[address.String()] = 1
			mutex.Unlock()
		}(i)
	}
	return &Server{
		remote:     remote,
		port:       port,
		router:     gin.Default(),
		clientPool: clientPool,
		toAddr:     to,
		log:        logrus.New(),
		hashMp:     sync.Map{},
	}, nil
}

func (server *Server) Start() {
	rand.Seed(time.Now().UnixNano())
	go func() {
		err := server.listenBlock()
		if err != nil {
			server.log.Error(err)
		}
	}()

	v1 := server.router.Group("v1")
	{
		v1.GET("/transfer", server.transfer)
		v1.GET("/interchain", server.interchain)
		v1.GET("/setData", server.setData)
		v1.GET("/getData", server.getData)
	}
	err := server.router.Run(fmt.Sprintf(":%d", server.port))
	if err != nil {
		server.log.Error(err)
		return
	}
}

func (server *Server) transfer(ctx *gin.Context) {
	grpc := server.getClient()
	data := &pb.TransactionData{
		Amount: "1",
	}
	payload, err := data.Marshal()
	if err != nil {
		server.log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	tx := &pb.BxhTransaction{
		From:      grpc.address,
		To:        types.NewAddressByStr("0xc7F999b83Af6DF9e67d0a37Ee7e900bF38b3D013"),
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}
	nonce := atomic.AddUint64(&grpc.nonce, 1)
	hash, err := grpc.client.SendTransaction(tx, &rpcx.TransactOpts{
		From:  grpc.address.String(),
		Nonce: nonce - 1,
	})
	if err != nil {
		server.log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	server.waitConfirm(hash)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func (server *Server) interchain(ctx *gin.Context) {
	grpc := server.getClient()
	nonce := atomic.AddUint64(&grpc.nonce, 1)
	index := atomic.AddUint64(&grpc.index, 1)
	ibtp := MockIBTP(grpc.address, server.toAddr, index-1)
	payload := MockContent(
		"interchainCharge",
		[][]byte{[]byte("Alice"), []byte("Alice"), []byte("1")},
	)
	ibtp.Payload = payload
	tx := &pb.BxhTransaction{
		From:      grpc.address,
		To:        constant.InterchainContractAddr.Address(),
		Timestamp: time.Now().UnixNano(),
		IBTP:      ibtp,
		Extra:     []byte("mock ibtp"),
	}
	hash, err := grpc.client.SendTransaction(tx, &rpcx.TransactOpts{
		From:  grpc.address.String(),
		Nonce: nonce - 1,
	})
	if err != nil {
		server.log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	server.waitConfirm(hash)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
func (server *Server) setData(ctx *gin.Context) {
	grpc := server.getClient()
	nonce := atomic.AddUint64(&grpc.nonce, 1)

	args := []*pb.Arg{
		pb.String(NormalKey),
		pb.String(NormalValue),
	}
	pl := &pb.InvokePayload{
		Method: "Set",
		Args:   args,
	}
	data, err := pl.Marshal()
	if err != nil {
		server.log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	td := &pb.TransactionData{
		Type:    pb.TransactionData_INVOKE,
		VmType:  pb.TransactionData_BVM,
		Payload: data,
	}
	payload, err := td.Marshal()
	if err != nil {
		server.log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	tx := &pb.BxhTransaction{
		From:      grpc.address,
		To:        constant.StoreContractAddr.Address(),
		Payload:   payload,
		Timestamp: time.Now().UnixNano(),
	}
	hash, err := grpc.client.SendTransaction(tx, &rpcx.TransactOpts{
		From:  grpc.address.String(),
		Nonce: nonce - 1,
	})
	if err != nil {
		server.log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	server.waitConfirm(hash)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
func (server *Server) getData(ctx *gin.Context) {
	grpc := server.getClient()
	nonce := atomic.AddUint64(&grpc.nonce, 1)

	args := []*pb.Arg{
		pb.String(NormalKey),
	}
	pl := &pb.InvokePayload{
		Method: "Get",
		Args:   args,
	}
	data, err := pl.Marshal()
	if err != nil {
		server.log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	td := &pb.TransactionData{
		Type:    pb.TransactionData_INVOKE,
		VmType:  pb.TransactionData_BVM,
		Payload: data,
	}
	payload, err := td.Marshal()
	if err != nil {
		server.log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	tx := &pb.BxhTransaction{
		From:      grpc.address,
		To:        constant.StoreContractAddr.Address(),
		Payload:   payload,
		Timestamp: time.Now().UnixNano(),
	}
	hash, err := grpc.client.SendTransaction(tx, &rpcx.TransactOpts{
		From:  grpc.address.String(),
		Nonce: nonce - 1,
	})
	if err != nil {
		server.log.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	server.waitConfirm(hash)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func (server *Server) listenBlock() error {
	pk, _, err := repo.KeyPriv()
	if err != nil {
		return err
	}
	client, err := rpcx.New(
		rpcx.WithNodesInfo(&rpcx.NodeInfo{Addr: server.remote}),
		rpcx.WithPrivateKey(pk),
	)
	if err != nil {
		return err
	}
	ch, err := client.Subscribe(context.TODO(), pb.SubscriptionRequest_BLOCK, nil)
	if err != nil {
		return err
	}
	for data := range ch {
		block := data.(*pb.Block)
		for _, tx := range block.Transactions.Transactions {
			server.hashMp.Store(tx.GetHash().String(), true)
		}
	}
	return nil
}

func (server *Server) getClient() *Grpc {
	idx := rand.Intn(len(server.clientPool))
	return server.clientPool[idx]
}

func (server *Server) waitConfirm(hash string) {
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		_, ok := server.hashMp.Load(hash)
		if ok {
			return
		}
	}
}
