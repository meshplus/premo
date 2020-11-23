package bitxhub

import "C"
import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-kit/types"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/wonderivan/logger"
)

var counter int64
var sender int64
var delayer int64
var ibtppd []byte
var proofHash [32]byte
var lock sync.Mutex

type bee struct {
	adminPrivKey  crypto.PrivateKey
	adminFrom     *types.Address
	normalPrivKey crypto.PrivateKey
	normalFrom    *types.Address
	client        rpcx.Client
	tps           int
	count         uint64
	adminSeqNo    uint64
	norMalSeqNo   uint64
	ibtpSeqNo     uint64
	ctx           context.Context
	cancel        context.CancelFunc
	config        *Config
}

func NewBee(tps int, adminPk crypto.PrivateKey, adminFrom *types.Address, expectedNonce uint64, config *Config) (*bee, error) {
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

	// query ibtp nonce for init in case ibtp has been sent to bitxhub before
	ibtp := mockIBTP(1, normalFrom.String(), normalFrom.String(), config.Proof)
	ibtpAccount := fmt.Sprintf("%s-%s-%d", normalFrom.String(), normalFrom.String(), ibtp.Category())
	ibtpNonce, err := client.GetPendingNonceByAccount(ibtpAccount)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &bee{
		client:        client,
		adminPrivKey:  adminPk,
		adminFrom:     adminFrom,
		normalPrivKey: normalPk,
		normalFrom:    normalFrom,
		tps:           tps,
		ctx:           ctx,
		cancel:        cancel,
		config:        config,
		adminSeqNo:    expectedNonce,
		ibtpSeqNo:     ibtpNonce,
		norMalSeqNo:   1,
	}, nil
}
func (bee *bee) start(typ string) error {
	var wg sync.WaitGroup
	wg.Add(bee.tps)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-bee.ctx.Done():
			wg.Wait()
			err := bee.client.Stop()
			if err != nil {
				panic(err)
			}
			return nil
		case <-ticker.C:
			for i := 0; i < bee.tps; i++ {
				bee.count++
				var (
					ibtpNo   uint64
					normalNo uint64
				)
				if typ == "interchain" {
					ibtpNo = atomic.LoadUint64(&bee.ibtpSeqNo)
					atomic.AddUint64(&bee.ibtpSeqNo, 1)
				} else {
					normalNo = atomic.LoadUint64(&bee.norMalSeqNo)
					atomic.AddUint64(&bee.norMalSeqNo, 1)
				}
				go func(ctx context.Context, count, ibtpNo, normalNo uint64) {
					select {
					case <-bee.ctx.Done():
						wg.Done()
					default:
						err := bee.sendTx(typ, count, ibtpNo, normalNo)
						if err != nil {
							logger.Error(err)
						}
					}
				}(bee.ctx, bee.count, ibtpNo, normalNo)
			}
		}
	}
}

func (bee *bee) sendTx(typ string, count, ibtpNo, normalNo uint64) error {
	switch typ {
	case "interchain":
		if err := bee.sendInterchainTx(count, ibtpNo); err != nil {
			return err
		}
	case "data":
		if err := bee.sendBVMTx(normalNo); err != nil {
			return err
		}
	case "transfer":
		fallthrough
	default:
		privkey, err := asym.GenerateKeyPair(crypto.Secp256k1)
		if err != nil {
			return err
		}

		to, err := privkey.PublicKey().Address()
		if err != nil {
			return err
		}

		if err := bee.sendTransferTx(to, normalNo); err != nil {
			return err
		}
	}
	return nil
}

func (bee *bee) stop() {
	bee.cancel()
}

func (bee *bee) sendBVMTx(normalNo uint64) error {
	atomic.AddInt64(&sender, 1)
	args := make([]*pb.Arg, 0)
	args = append(args, rpcx.String("a"), rpcx.String("10"))

	pl := &pb.InvokePayload{
		Method: "Set",
		Args:   args,
	}

	data, err := pl.Marshal()
	if err != nil {
		return err
	}

	td := &pb.TransactionData{
		Type:    pb.TransactionData_INVOKE,
		VmType:  pb.TransactionData_BVM,
		Payload: data,
	}
	payload, err := td.Marshal()
	if err != nil {
		return err
	}

	tx := &pb.Transaction{
		From:      bee.normalFrom,
		To:        constant.StoreContractAddr.Address(),
		Payload:   payload,
		Timestamp: time.Now().UnixNano(),
		Nonce:     normalNo,
	}

	txHash, err := bee.client.SendTransaction(tx, &rpcx.TransactOpts{
		NormalNonce: normalNo,
	})
	if err != nil {
		return err
	}
	tx.TransactionHash = types.NewHashByStr(txHash)

	go bee.counterReceipt(tx)
	return nil
}

func (bee *bee) prepareChain(chainType, name, validators, version, desc string, contract []byte) error {
	bee.client.SetPrivateKey(bee.normalPrivKey)
	// register chain
	pubKey, _ := bee.normalPrivKey.PublicKey().Bytes()
	receipt, err := bee.invokeContract(bee.normalFrom, constant.AppchainMgrContractAddr.Address(), atomic.LoadUint64(&bee.norMalSeqNo),
		"Register", rpcx.String(validators), rpcx.Int32(1), rpcx.String(chainType),
		rpcx.String(name), rpcx.String(desc), rpcx.String(version), rpcx.String(string(pubKey)))
	if err != nil {
		return fmt.Errorf("register appchain error: %w", err)
	}

	atomic.AddUint64(&bee.norMalSeqNo, 1)

	appchain := &rpcx.Appchain{}
	if err := json.Unmarshal(receipt.Ret, appchain); err != nil {
		return err
	}
	ID := appchain.ID

	// Audit chain and set adminPrivateKey for auditing
	bee.client.SetPrivateKey(bee.adminPrivKey)
	receipt, err = bee.invokeContract(bee.adminFrom, constant.AppchainMgrContractAddr.Address(), bee.adminSeqNo,
		"Audit", rpcx.String(ID), rpcx.Int32(1), rpcx.String("Audit passed"))
	if err != nil {
		return fmt.Errorf("audit appchain error:%w", err)
	}

	ruleAddr := "0x00000000000000000000000000000000000000a1"
	// deploy rule
	bee.client.SetPrivateKey(bee.normalPrivKey)
	if chainType == "hyperchain" {
		contractAddr, err := bee.client.DeployContract(contract, nil)
		if err != nil {
			return fmt.Errorf("deploy contract error:%w", err)
		}

		ruleAddr = contractAddr.String()
	} else if chainType == "fabric:complex" {
		ruleAddr = "0x00000000000000000000000000000000000000a0"
	}

	_, err = bee.invokeContract(bee.normalFrom, ValidationContractAddr, atomic.LoadUint64(&bee.norMalSeqNo),
		"RegisterRule", rpcx.String(ID), rpcx.String(ruleAddr))
	if err != nil {
		return fmt.Errorf("register rule error:%w", err)
	}
	atomic.AddUint64(&bee.norMalSeqNo, 1)

	prepareInterchainTx(bee.config.Proof)

	return nil
}

func (bee *bee) invokeContract(from, to *types.Address, nonce uint64, method string, args ...*pb.Arg) (*pb.Receipt, error) {
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
		From:        from.String(),
		NormalNonce: nonce,
	})
}

func (bee *bee) sendTransferTx(to *types.Address, normalNo uint64) error {
	atomic.AddInt64(&sender, 1)

	data := &pb.TransactionData{
		Type:   pb.TransactionData_NORMAL,
		VmType: pb.TransactionData_XVM,
		Amount: 0,
	}
	payload, err := data.Marshal()
	if err != nil {
		return err
	}
	tx := &pb.Transaction{
		From:      bee.normalFrom,
		To:        to,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}

	txHash, err := bee.client.SendTransaction(tx, &rpcx.TransactOpts{
		From:        bee.normalFrom.String(),
		NormalNonce: normalNo,
	})
	if err != nil {
		return err
	}
	tx.TransactionHash = types.NewHashByStr(txHash)
	go bee.counterReceipt(tx)

	return nil
}

func (bee *bee) sendInterchainTx(i uint64, ibtpNo uint64) error {
	atomic.AddInt64(&sender, 1)
	ibtp := mockIBTP(i, bee.normalFrom.String(), bee.normalFrom.String(), bee.config.Proof)

	tx := &pb.Transaction{
		From:      bee.normalFrom,
		To:        constant.InterchainContractAddr.Address(),
		Timestamp: time.Now().UnixNano(),
		Extra:     bee.config.Proof,
		IBTP:      ibtp,
	}

	txHash, err := bee.client.SendTransaction(tx, &rpcx.TransactOpts{
		From:      fmt.Sprintf("%s-%s-%d", ibtp.From, ibtp.To, ibtp.Category()),
		IBTPNonce: ibtpNo,
	})
	if err != nil {
		return err
	}
	tx.TransactionHash = types.NewHashByStr(txHash)
	go bee.counterReceipt(tx)

	return nil
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

func (bee *bee) counterReceipt(tx *pb.Transaction) {
	for {
		receipt, err := bee.client.GetReceipt(tx.Hash().String())
		if err != nil {
			if strings.Contains(err.Error(), "not found in DB") {
				continue
			}
			logger.Error(err)
			return
		}

		if !receipt.IsSuccess() {
			logger.Error("receipt for tx %s is failed, error msg: %s", tx.TransactionHash.String(), string(receipt.Ret))
			return
		}
		break
	}
	atomic.AddInt64(&delayer, time.Now().UnixNano()-tx.Timestamp)
	atomic.AddInt64(&counter, 1)
}
