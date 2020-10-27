package bitxhub

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
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
var ibtppd []byte
var proofHash [32]byte
var lock sync.Mutex

type bee struct {
	xprivKey    crypto.PrivateKey
	xfrom       *types.Address
	xid         string
	client      rpcx.Client
	tps         int
	count       uint64
	norMalSeqNo uint64
	ibtpSeqNo   uint64
	ctx         context.Context
	cancel      context.CancelFunc
	config      *Config
}

func NewBee(tps int, addrs []string, config *Config) (*bee, error) {
	xpk, err := asym.RestorePrivateKey(config.KeyPath, repo.KeyPassword)
	if err != nil {
		return nil, err
	}

	xfrom, err := xpk.PublicKey().Address()
	if err != nil {
		return nil, err
	}

	node0 := &rpcx.NodeInfo{Addr: config.BitxhubAddr[0]}
	client, err := rpcx.New(
		rpcx.WithNodesInfo(node0),
		rpcx.WithLogger(cfg.logger),
		rpcx.WithPrivateKey(xpk),
	)
	if err != nil {
		return nil, err
	}

	// query pending nonce for privKey
	normalNonce, err := client.GetPendingNonceByAccount(xfrom.String())
	if err != nil {
		return nil, err
	}

	// query ibtp nonce for init in case ibtp has been sent to bitxhub before
	ibtp := mockIBTP(1, xfrom.String(), xfrom.String(), config.Proof)
	ibtpNonce, err := client.GetPendingNonceByAccount(
		fmt.Sprintf("%s-%s-%d", xfrom.String(), xfrom.String(), ibtp.Category()))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &bee{
		client:      client,
		xprivKey:    xpk,
		xfrom:       xfrom,
		tps:         tps,
		ctx:         ctx,
		cancel:      cancel,
		xid:         "",
		config:      config,
		ibtpSeqNo:   ibtpNonce,
		norMalSeqNo: normalNonce,
	}, nil
}

func (bee *bee) start(typ string) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-bee.ctx.Done():
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
				go func(count, ibtpNo, normalNo uint64) {
					err := bee.sendTx(typ, count, ibtpNo, normalNo)
					if err != nil {
						logger.Error(err)
					}
				}(bee.count, ibtpNo, normalNo)
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
		From:      bee.xfrom,
		To:        constant.StoreContractAddr.Address(),
		Payload:   payload,
		Timestamp: time.Now().UnixNano(),
		Nonce:     normalNo,
	}

	_, err = bee.client.SendTransaction(tx, &rpcx.TransactOpts{
		NormalNonce: normalNo,
	})
	if err != nil {
		return err
	}

	go bee.counterReceipt(tx)
	return nil
}

func (bee *bee) prepareChain(chainType, name, validators, version, desc string, contract []byte) error {
	bee.client.SetPrivateKey(bee.xprivKey)
	// register chain
	pubKey, _ := bee.xprivKey.PublicKey().Bytes()
	receipt, err := bee.invokeContract(constant.AppchainMgrContractAddr.Address(), atomic.LoadUint64(&bee.norMalSeqNo),
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

	// Audit chain
	receipt, err = bee.invokeContract(constant.AppchainMgrContractAddr.Address(), atomic.LoadUint64(&bee.norMalSeqNo),
		"Audit", rpcx.String(ID), rpcx.Int32(1), rpcx.String(""))
	if err != nil {
		return fmt.Errorf("audit appchain error:%w", err)
	}
	atomic.AddUint64(&bee.norMalSeqNo, 1)

	ruleAddr := "0x00000000000000000000000000000000000000a1"
	// deploy rule
	if chainType == "hyperchain" {
		contractAddr, err := bee.client.DeployContract(contract, nil)
		if err != nil {
			return fmt.Errorf("deploy contract error:%w", err)
		}

		ruleAddr = contractAddr.String()
	} else if chainType == "fabric:complex" {
		ruleAddr = "0x00000000000000000000000000000000000000a0"
	}

	_, err = bee.invokeContract(ValidationContractAddr, atomic.LoadUint64(&bee.norMalSeqNo),
		"RegisterRule", rpcx.String(ID), rpcx.String(ruleAddr))
	if err != nil {
		return fmt.Errorf("register rule error:%w", err)
	}
	atomic.AddUint64(&bee.norMalSeqNo, 1)

	prepareInterchainTx(bee.config.Proof)

	return nil
}

func (bee *bee) invokeContract(address *types.Address, nonce uint64, method string, args ...*pb.Arg) (*pb.Receipt, error) {
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
		From:      bee.xfrom,
		To:        address,
		Payload:   payload,
		Timestamp: time.Now().UnixNano(),
		Nonce:     nonce,
	}

	return bee.client.SendTransactionWithReceipt(tx, nil)
}

func (bee *bee) sendTransferTx(to *types.Address, normalNo uint64) error {
	atomic.AddInt64(&sender, 1)

	data := &pb.TransactionData{
		Type:   pb.TransactionData_NORMAL,
		Amount: 1,
	}
	payload, err := data.Marshal()
	if err != nil {
		return err
	}
	tx := &pb.Transaction{
		From:      bee.xfrom,
		To:        to,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}

	_, err = bee.client.SendTransaction(tx, &rpcx.TransactOpts{
		From:        bee.xfrom.String(),
		NormalNonce: normalNo,
	})
	if err != nil {
		return err
	}
	go bee.counterReceipt(tx)

	return nil
}

func (bee *bee) sendInterchainTx(i uint64, ibtpNo uint64) error {
	atomic.AddInt64(&sender, 1)
	ibtp := mockIBTP(i, bee.xfrom.String(), bee.xfrom.String(), bee.config.Proof)
	b, err := ibtp.Marshal()
	if err != nil {
		return err
	}

	pl := &pb.InvokePayload{
		Method: "HandleIBTP",
		Args:   []*pb.Arg{pb.Bytes(b)}[:],
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

	tx := &pb.Transaction{
		From:      bee.xfrom,
		To:        constant.InterchainContractAddr.Address(),
		Payload:   payload,
		Timestamp: time.Now().UnixNano(),
		Extra:     bee.config.Proof,
	}

	_, err = bee.client.SendTransaction(tx, &rpcx.TransactOpts{
		From:      fmt.Sprintf("%s-%s-%d", ibtp.From, ibtp.To, ibtp.Category()),
		IBTPNonce: ibtpNo,
	})
	if err != nil {
		return err
	}
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
			if err.Error() == "not found in DB" {
				continue
			}
			logger.Error(err)
			return
		}

		if receipt.Status.String() == "FAILED" {
			logger.Error(string(receipt.Ret))
			return
		}
		break
	}
	atomic.AddInt64(&delayer, time.Now().UnixNano()-tx.Timestamp)
	atomic.AddInt64(&counter, 1)
}
