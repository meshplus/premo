package bitxhub

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-kit/fileutil"
	"github.com/meshplus/bitxhub-kit/types"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/meshplus/premo/internal/repo"
	"github.com/wonderivan/logger"
)

var counter int64
var sender int64
var delayer int64

type bee struct {
	xprivKey crypto.PrivateKey
	xfrom    types.Address
	xid      string
	client   rpcx.Client
	privKey  crypto.PrivateKey
	from     types.Address
	tps      int
	count    uint64
	ctx      context.Context
	cancel   context.CancelFunc
	config   *Config
}

func NewBee(tps int, keyPath string, addrs []string, config *Config) (*bee, error) {
	xpk, err := asym.GenerateKeyPair(crypto.Secp256k1)

	if err != nil {
		return nil, err
	}

	xfrom, err := xpk.PublicKey().Address()
	if err != nil {
		return nil, err
	}

	if !fileutil.Exist(keyPath) {
		keyPath, err = repo.KeyPath()
		if err != nil {
			return nil, err
		}
	}
	privKey, err := asym.RestorePrivateKey(keyPath, repo.KeyPassword)
	if err != nil {
		return nil, err
	}

	from, err := privKey.PublicKey().Address()
	if err != nil {
		return nil, err
	}

	client, err := rpcx.New(
		rpcx.WithAddrs(addrs),
		rpcx.WithLogger(cfg.logger),
		rpcx.WithPrivateKey(privKey),
	)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &bee{
		client:   client,
		xprivKey: xpk,
		xfrom:    xfrom,
		privKey:  privKey,
		from:     from,
		tps:      tps,
		ctx:      ctx,
		cancel:   cancel,
		xid:      "",
		config:   config,
	}, nil
}

func (bee *bee) start(typ string) error {
	d := time.Second / time.Duration(bee.tps*10/9)
	ticker := time.NewTicker(d)
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
			atomic.AddUint64(&bee.count, 1)
			switch typ {
			case "interchain":
				if err := bee.sendInterchainTx(bee.count); err != nil {
					logger.Error(err)
					return err
				}
			case "data":
				if err := bee.sendBVMTx(bee.count); err != nil {
					logger.Error(err)
					return err
				}
			case "transfer":
				fallthrough
			default:
				privkey, err := asym.GenerateKeyPair(crypto.Secp256k1)
				if err != nil {
					logger.Error(err)
					return err
				}

				to, err := privkey.PublicKey().Address()
				if err != nil {
					return err
				}

				if err := bee.sendTransferTx(to); err != nil {
					logger.Error(err)
					return err
				}
			}
		}
	}
}

func (bee *bee) stop() {
	bee.cancel()
}

func (bee *bee) sendBVMTx(i uint64) error {
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

	tx := &pb.Transaction{
		From:      bee.xfrom,
		To:        rpcx.StoreContractAddr,
		Data:      td,
		Timestamp: time.Now().UnixNano(),
		Nonce:     rand.Int63(),
	}

	if err := tx.Sign(bee.privKey); err != nil {
		return err
	}
	_, err = bee.client.SendTransaction(tx)
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
	receipt, err := bee.client.InvokeBVMContract(rpcx.AppchainMgrContractAddr, "Register", rpcx.String(validators),
		rpcx.Int32(1), rpcx.String(chainType), rpcx.String(name), rpcx.String(desc), rpcx.String(version), rpcx.String(string(pubKey)))
	if err != nil {
		return fmt.Errorf("register appchain error: %w", err)
	}
	appchain := &rpcx.Appchain{}
	if err := json.Unmarshal(receipt.Ret, appchain); err != nil {
		return err
	}
	ID := appchain.ID

	// Audit chain
	_, err = bee.client.InvokeBVMContract(rpcx.AppchainMgrContractAddr,
		"Audit", rpcx.String(ID), rpcx.Int32(1), rpcx.String(""))
	if err != nil {
		return fmt.Errorf("audit appchain error:%w", err)
	}

	ruleAddr := "0x00000000000000000000000000000000000000a1"
	// deploy rule
	if chainType == "hyperchain" {
		contractAddr, err := bee.client.DeployContract(contract)
		if err != nil {
			return fmt.Errorf("deploy contract error:%w", err)
		}

		ruleAddr = contractAddr.String()
	} else if chainType == "fabric:complex" {
		ruleAddr = "0x00000000000000000000000000000000000000a0"
	}

	_, err = bee.client.InvokeContract(pb.TransactionData_BVM, ValidationContractAddr, "RegisterRule", rpcx.String(ID), rpcx.String(ruleAddr))
	if err != nil {
		return fmt.Errorf("register rule error:%w", err)
	}
	return nil
}

func (bee *bee) sendTransferTx(to types.Address) error {
	atomic.AddInt64(&sender, 1)

	tx := &pb.Transaction{
		From:      bee.from,
		To:        to,
		Timestamp: time.Now().UnixNano(),
		Data: &pb.TransactionData{
			Type:   pb.TransactionData_NORMAL,
			Amount: 1,
		},
		Nonce: rand.Int63(),
	}

	err := tx.Sign(bee.privKey)
	if err != nil {
		return err
	}

	_, err = bee.client.SendTransaction(tx)
	if err != nil {
		return err
	}
	go bee.counterReceipt(tx)

	return nil
}

func (bee *bee) sendInterchainTx(i uint64) error {
	atomic.AddInt64(&sender, 1)
	ibtp := mockIBTP(i, bee.xfrom.String(), bee.xfrom.String(), bee.config.Proof)
	b, err := ibtp.Marshal()
	if err != nil {
		return err
	}

	args := make([]*pb.Arg, 0)
	args = append(args, rpcx.Bytes(b))

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

	tx := &pb.Transaction{
		From:      bee.xfrom,
		To:        rpcx.InterchainContractAddr,
		Data:      td,
		Timestamp: time.Now().UnixNano(),
		Nonce:     rand.Int63(),
		Extra:     bee.config.Proof,
	}

	if err := tx.Sign(bee.xprivKey); err != nil {
		return err
	}

	_, err = bee.client.SendTransaction(tx)
	if err != nil {
		return err
	}
	go bee.counterReceipt(tx)

	return nil
}

func mockIBTP(index uint64, from, to string, proof []byte) *pb.IBTP {
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

	ibtppd, _ := payload.Marshal()

	proofHash := sha256.Sum256(proof)

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
