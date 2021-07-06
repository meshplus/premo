package bitxhub

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/coreos/etcd/pkg/stringutil"
	appchain_mgr "github.com/meshplus/bitxhub-core/appchain-mgr"
	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-kit/hexutil"
	"github.com/meshplus/bitxhub-kit/types"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/meshplus/premo/internal/repo"
	"github.com/wonderivan/logger"
)

var counter int64
var sender int64
var delayer int64
var maxDelay int64
var ibtppd []byte
var proofHash [32]byte

type Bee struct {
	adminPrivKey  crypto.PrivateKey
	adminFrom     *types.Address
	normalPrivKey crypto.PrivateKey
	normalFrom    *types.Address
	client        rpcx.Client
	tps           int
	count         uint64
	adminSeqNo    uint64
	nonce         uint64
	ctx           context.Context
	config        *Config
}

type RegisterResult struct {
	ChainID    string `json:"chain_id"`
	ProposalID string `json:"proposal_id"`
}

func NewBee(tps int, adminPk crypto.PrivateKey, adminFrom *types.Address, expectedNonce uint64, config *Config, ctx context.Context) (*Bee, error) {
	normalPk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	if err != nil {
		return nil, err
	}
	normalFrom, err := normalPk.PublicKey().Address()
	if err != nil {
		return nil, err
	}

	node0 := &rpcx.NodeInfo{Addr: config.BitxhubAddr[0]}
	client, err := rpcx.New(
		rpcx.WithNodesInfo(node0),
		rpcx.WithLogger(cfg.logger),
		rpcx.WithPrivateKey(normalPk),
	)
	if err != nil {
		return nil, err
	}

	return &Bee{
		client:        client,
		adminPrivKey:  adminPk,
		adminFrom:     adminFrom,
		normalPrivKey: normalPk,
		normalFrom:    normalFrom,
		tps:           tps,
		ctx:           ctx,
		config:        config,
		adminSeqNo:    expectedNonce,
		nonce:         1,
	}, nil
}

func (bee *Bee) GetAddress() *types.Address {
	return bee.normalFrom
}

func (bee *Bee) start(typ string) error {
	ticker := time.NewTicker(time.Millisecond * 50)
	defer ticker.Stop()

	for {
		select {
		case <-bee.ctx.Done():
			return nil
		case <-ticker.C:
			for i := 0; i < bee.tps/20; i++ {
				bee.count++
				go func(count uint64) {
					select {
					case <-bee.ctx.Done():
						return
					default:
						_, err := bee.SendTx(typ, "", count)
						if err != nil {
							logger.Error(err)
						}
					}
				}(bee.count)
			}
		}
	}
}

func (bee *Bee) genTx(typ, key string, count uint64) (*pb.Transaction, error) {
	switch typ {
	case "interchain":
		return bee.genInterchainTx(count), nil
	case "getData":
		return bee.genBVMTx(key, true)
	case "setData":
		return bee.genBVMTx(key, false)
	case "transfer":
		fallthrough
	default:
		to := &types.Address{}
		return bee.genTransferTx(to, 1)
	}
}

func (bee *Bee) SendTx(typ, key string, count uint64) (string, error) {
	nonce := atomic.LoadUint64(&bee.nonce)
	atomic.AddUint64(&bee.nonce, 1)
	tx, err := bee.genTx(typ, key, count)
	if err != nil {
		return "", err
	}

	return bee.client.SendTransaction(tx, &rpcx.TransactOpts{
		Nonce: nonce,
	})
}

func (bee *Bee) SendTxSync(typ, key string, count uint64) (*pb.Receipt, error) {
	nonce := atomic.LoadUint64(&bee.nonce)
	atomic.AddUint64(&bee.nonce, 1)
	tx, err := bee.genTx(typ, key, count)
	if err != nil {
		return nil, err
	}

	return bee.client.SendTransactionSync(tx, &rpcx.TransactOpts{
		Nonce: nonce,
	})
}

func (bee *Bee) SendDoubleSpendTxs() ([]*pb.Receipt, error) {
	var (
		receipts []*pb.Receipt
		wg       sync.WaitGroup
		lock     sync.Mutex
	)
	wg.Add(2)
	for i := 0; i < 2; i++ {
		nonce := atomic.LoadUint64(&bee.nonce)
		atomic.AddUint64(&bee.nonce, 1)
		go func(nonce uint64, i int) {
			defer wg.Done()
			tx, err := bee.genTransferTx(types.NewAddress([]byte{byte(i)}), 100000)
			if err != nil {
				return
			}
			receipt, err := bee.client.SendTransactionSync(tx, &rpcx.TransactOpts{Nonce: nonce})
			if err != nil {
				return
			}
			lock.Lock()
			receipts = append(receipts, receipt)
			lock.Unlock()
		}(nonce, i)
	}

	wg.Wait()

	return receipts, nil
}

func (bee *Bee) stop() {
	bee.client.Stop()
	return
}

func (bee *Bee) genBVMTx(key string, isGet bool) (*pb.Transaction, error) {
	atomic.AddInt64(&sender, 1)
	args := make([]*pb.Arg, 0)
	var pl *pb.InvokePayload

	if key == "" {
		keys := stringutil.RandomStrings(20, 1)
		key = keys[0]
	}

	if isGet {
		args = append(args, rpcx.String(key))
		pl = &pb.InvokePayload{
			Method: "Get",
			Args:   args,
		}
	} else {
		hash := sha256.Sum256([]byte(key))
		args = append(args, rpcx.String(key), rpcx.String(hexutil.Encode(hash[:])))
		pl = &pb.InvokePayload{
			Method: "Set",
			Args:   args,
		}
	}

	data, err := pl.Marshal()
	if err != nil {
		return nil, err
	}

	td := &pb.TransactionData{
		Type:    pb.TransactionData_INVOKE,
		VmType:  pb.TransactionData_BVM,
		Payload: data,
	}
	payload, err := td.Marshal()
	if err != nil {
		return nil, err
	}

	return &pb.Transaction{
		From:      bee.normalFrom,
		To:        constant.StoreContractAddr.Address(),
		Payload:   payload,
		Timestamp: time.Now().UnixNano(),
	}, nil
}

func (bee *Bee) prepareChain(chainType, name, validators, version, desc string, contract []byte) error {
	bee.client.SetPrivateKey(bee.normalPrivKey)
	// register chain
	pubKey, _ := bee.normalPrivKey.PublicKey().Bytes()
	receipt, err := bee.invokeContract(bee.normalFrom, constant.AppchainMgrContractAddr.Address(), atomic.LoadUint64(&bee.nonce),
		"Register", rpcx.String(validators), rpcx.String("raft"), rpcx.String(chainType),
		rpcx.String(name), rpcx.String(desc), rpcx.String(version), rpcx.String(string(pubKey)))
	if err != nil {
		return fmt.Errorf("register appchain error: %w", err)
	}

	atomic.AddUint64(&bee.nonce, 1)
	// vote chain
	result := &RegisterResult{}
	err = json.Unmarshal(receipt.Ret, result)
	if err != nil || result.ProposalID == "" {
		return fmt.Errorf("vote unmarshal error: %w", err)
	}
	err = bee.VotePass(result.ProposalID)
	if err != nil {
		return fmt.Errorf("vote chain error: %w", err)
	}
	res, err := bee.GetChainStatusById(result.ChainID)
	if err != nil {
		return fmt.Errorf("getChainStatus error: %w", err)
	}
	appchain := &rpcx.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	if err != nil || appchain.Status != appchain_mgr.AppchainAvailable {
		return fmt.Errorf("chain error: %w", err)
	}
	ID := result.ChainID

	ruleAddr := "0x00000000000000000000000000000000000000a1"
	// deploy rule
	bee.client.SetPrivateKey(bee.normalPrivKey)
	if chainType == "hyperchain" {
		contractAddr, err := bee.client.DeployContract(contract, nil)
		if err != nil {
			return fmt.Errorf("deploy contract error:%w", err)
		}
		atomic.AddUint64(&bee.nonce, 1)
		ruleAddr = contractAddr.String()
	} else if chainType == "fabric:complex" {
		ruleAddr = "0x00000000000000000000000000000000000000a0"
	}

	_, err = bee.invokeContract(bee.normalFrom, ValidationContractAddr, atomic.LoadUint64(&bee.nonce),
		"RegisterRule", rpcx.String(ID), rpcx.String(ruleAddr))
	if err != nil {
		return fmt.Errorf("register rule error:%w", err)
	}
	atomic.AddUint64(&bee.nonce, 1)

	prepareInterchainTx(bee.config.Proof)

	return nil
}

func (bee *Bee) invokeContract(from, to *types.Address, nonce uint64, method string, args ...*pb.Arg) (*pb.Receipt, error) {
	pl := &pb.InvokePayload{
		Method: method,
		Args:   args[:],
	}

	data, err := pl.Marshal()
	if err != nil {
		return nil, err
	}

	td := &pb.TransactionData{
		Type:    pb.TransactionData_INVOKE,
		VmType:  pb.TransactionData_BVM,
		Payload: data,
	}
	payload, err := td.Marshal()

	tx := &pb.Transaction{
		From:      from,
		To:        to,
		Payload:   payload,
		Timestamp: time.Now().UnixNano(),
	}

	return bee.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From:  from.String(),
		Nonce: nonce,
	})
}

func (bee *Bee) genTransferTx(to *types.Address, amount uint64) (*pb.Transaction, error) {
	atomic.AddInt64(&sender, 1)

	data := &pb.TransactionData{
		Type:   pb.TransactionData_NORMAL,
		VmType: pb.TransactionData_XVM,
		Amount: amount,
	}
	payload, err := data.Marshal()
	if err != nil {
		return nil, err
	}

	return &pb.Transaction{
		From:      bee.normalFrom,
		To:        to,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}, nil
}

func (bee *Bee) genInterchainTx(i uint64) *pb.Transaction {
	atomic.AddInt64(&sender, 1)
	ibtp := mockIBTP(i, bee.normalFrom.String(), bee.normalFrom.String(), bee.config.Proof)

	return &pb.Transaction{
		From:      bee.normalFrom,
		To:        constant.InterchainContractAddr.Address(),
		Timestamp: time.Now().UnixNano(),
		Extra:     bee.config.Proof,
		IBTP:      ibtp,
	}
}

func prepareInterchainTx(proof []byte) {
	if ibtppd != nil {
		return
	}

	content := &pb.Content{
		SrcContractId: "mychannel&transfer",
		DstContractId: "mychannel&transfer",
		Func:          "interchainCharge",
		Args:          [][]byte{[]byte("Alice"), []byte("Alice"), []byte("1")},
		Callback:      "interchainConfirm",
	}

	bytes, _ := content.Marshal()

	payload := &pb.Payload{
		Encrypted: false,
		Content:   bytes,
	}

	ibtppd, _ = payload.Marshal()
	proofHash = sha256.Sum256(proof)
}

func mockIBTP(index uint64, from, to string, proof []byte) *pb.IBTP {
	return &pb.IBTP{
		From:      from,
		To:        to,
		Payload:   ibtppd,
		Index:     index,
		Type:      pb.IBTP_INTERCHAIN,
		Timestamp: time.Now().UnixNano(),
		Proof:     proofHash[:],
	}
}

func (bee *Bee) VotePass(id string) error {
	node1, err := repo.Node1Path()
	if err != nil {
		return err
	}

	key, err := asym.RestorePrivateKey(node1, repo.KeyPassword)
	if err != nil {
		return err
	}

	_, err = bee.vote(key, atomic.AddUint64(&index1, 1), pb.String(id), pb.String("approve"), pb.String("Appchain Pass"))
	if err != nil {
		return err
	}

	node2, err := repo.Node2Path()
	if err != nil {
		return err
	}

	key, err = asym.RestorePrivateKey(node2, repo.KeyPassword)
	if err != nil {
		return err
	}

	_, err = bee.vote(key, atomic.AddUint64(&index2, 1), pb.String(id), pb.String("approve"), pb.String("Appchain Pass"))
	if err != nil {
		return err
	}

	node3, err := repo.Node3Path()
	if err != nil {
		return err
	}

	key, err = asym.RestorePrivateKey(node3, repo.KeyPassword)
	if err != nil {
		return err
	}

	_, err = bee.vote(key, atomic.AddUint64(&index3, 1), pb.String(id), pb.String("approve"), pb.String("Appchain Pass"))
	if err != nil {
		return err
	}
	return nil
}

func (bee *Bee) vote(key crypto.PrivateKey, index uint64, args ...*pb.Arg) (*pb.Receipt, error) {
	client, err := rpcx.New(
		rpcx.WithNodesInfo(&rpcx.NodeInfo{Addr: bee.config.BitxhubAddr[0]}),
		rpcx.WithLogger(cfg.logger),
		rpcx.WithPrivateKey(key),
	)
	address, err := key.PublicKey().Address()
	if err != nil {
		return nil, err
	}
	invokePayload := &pb.InvokePayload{
		Method: "Vote",
		Args:   args,
	}

	payload, err := invokePayload.Marshal()
	if err != nil {
		return nil, err
	}

	data := &pb.TransactionData{
		Type:    pb.TransactionData_INVOKE,
		VmType:  pb.TransactionData_BVM,
		Payload: payload,
	}
	payload, err = data.Marshal()

	tx := &pb.Transaction{
		From:      address,
		To:        constant.GovernanceContractAddr.Address(),
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}
	if err != nil {
		return nil, err
	}
	receipt, err := client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		Nonce: index,
	})
	if err != nil {
		return nil, err
	}
	return receipt, nil
}

func (bee *Bee) GetChainStatusById(id string) (*pb.Receipt, error) {
	node, err := repo.Node1Path()
	key, err := asym.RestorePrivateKey(node, repo.KeyPassword)
	if err != nil {
		return nil, err
	}
	client, err := rpcx.New(
		rpcx.WithNodesInfo(&rpcx.NodeInfo{Addr: bee.config.BitxhubAddr[0]}),
		rpcx.WithLogger(cfg.logger),
		rpcx.WithPrivateKey(key),
	)
	address, err := key.PublicKey().Address()
	if err != nil {
		return nil, err
	}
	args := []*pb.Arg{
		rpcx.String(id),
	}
	invokePayload := &pb.InvokePayload{
		Method: "GetAppchain",
		Args:   args,
	}

	payload, err := invokePayload.Marshal()
	if err != nil {
		return nil, err
	}

	data := &pb.TransactionData{
		Type:    pb.TransactionData_INVOKE,
		VmType:  pb.TransactionData_BVM,
		Payload: payload,
	}
	payload, err = data.Marshal()

	tx := &pb.Transaction{
		From:      address,
		To:        constant.AppchainMgrContractAddr.Address(),
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}
	if err != nil {
		return nil, err
	}
	receipt, err := client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		Nonce: atomic.AddUint64(&index1, 1),
	})
	if err != nil {
		return nil, err
	}
	return receipt, nil
}
