package bitxhub

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	appchain_mgr "github.com/meshplus/bitxhub-core/appchain-mgr"
	"github.com/meshplus/bitxhub-core/governance"
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

type bee struct {
	adminPrivKey  crypto.PrivateKey
	adminFrom     *types.Address
	normalPrivKey crypto.PrivateKey
	normalFrom    *types.Address
	client        rpcx.Client
	tps           int
	count         uint64
	adminNonce    uint64
	nonce         uint64
	ctx           context.Context
	config        *Config
}

type RegisterResult struct {
	Extra      []byte `json:"extra"`
	ProposalID string `json:"proposal_id"`
}

func NewBee(tps int, adminPk crypto.PrivateKey, adminFrom *types.Address, expectedNonce uint64, config *Config, ctx context.Context) (*bee, error) {
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

	return &bee{
		client:        client,
		adminPrivKey:  adminPk,
		adminFrom:     adminFrom,
		normalPrivKey: normalPk,
		normalFrom:    normalFrom,
		tps:           tps,
		ctx:           ctx,
		config:        config,
		adminNonce:    expectedNonce,
		nonce:         0,
	}, nil
}

func (bee *bee) start(typ string) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-bee.ctx.Done():
			return nil
		case <-ticker.C:
			for i := 0; i < bee.tps; i++ {
				bee.count++
				nonce := atomic.LoadUint64(&bee.nonce)
				atomic.AddUint64(&bee.nonce, 1)
				go func(count, nonce uint64) {
					select {
					case <-bee.ctx.Done():
						return
					default:
						err := bee.sendTx(typ, count, nonce)
						if err != nil {
							logger.Error(err)
						}
					}
				}(bee.count, nonce)
			}
		}
	}
}

func (bee *bee) sendTx(typ string, count, nonce uint64) error {
	switch typ {
	case "interchain":
		if err := bee.sendInterchainTx(count, nonce); err != nil {
			return err
		}
	case "data":
		if err := bee.sendBVMTx(nonce); err != nil {
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

		if err := bee.sendTransferTx(to, nonce); err != nil {
			return err
		}
	}
	return nil
}

func (bee *bee) stop() {
	bee.client.Stop()
	return
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

	tx := &pb.BxhTransaction{
		From:      bee.normalFrom,
		To:        constant.StoreContractAddr.Address(),
		Payload:   payload,
		Timestamp: time.Now().UnixNano(),
		Nonce:     normalNo,
	}

	txHash, err := bee.client.SendTransaction(tx, &rpcx.TransactOpts{
		Nonce: normalNo,
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
	var pubKeyStr = base64.StdEncoding.EncodeToString(pubKey)
	var ct = ""
	if chainType == "fabric:simple" || chainType == "fabric:complex" {
		ct = "fabric"
	} else {
		ct = chainType
	}
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain" + bee.normalFrom.String() + ":" + bee.normalFrom.String()), //id
		rpcx.String("did:bitxhub:appchain" + bee.normalFrom.String() + ":."),                          //ownerDID
		rpcx.String("/ipfs/QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),                           //docAddr
		rpcx.String("QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),                                 //docHash
		rpcx.String(validators), //validators
		rpcx.String("raft"),     //consensus_type
		rpcx.String(ct),         //chain_type
		rpcx.String(name),       //name
		rpcx.String(desc),       //desc
		rpcx.String(version),    //version
		rpcx.String(pubKeyStr),  //public key
	}
	receipt, err := bee.invokeContract(bee.normalFrom, constant.AppchainMgrContractAddr.Address(), atomic.LoadUint64(&bee.nonce),
		"Register", args...)
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
	res, err := bee.GetChainStatusById(string(result.Extra))
	if err != nil {
		return fmt.Errorf("getChainStatus error: %w", err)
	}
	appchain := &appchain_mgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	if err != nil || appchain.Status != governance.GovernanceAvailable {
		return fmt.Errorf("chain error: %w", err)
	}
	ID := string(result.Extra)

	ruleAddr := "0x00000000000000000000000000000000000000a0"
	// deploy rule
	bee.client.SetPrivateKey(bee.normalPrivKey)
	if chainType == "hyperchain" {
		contractAddr, err := bee.client.DeployContract(contract, nil)
		if err != nil {
			return fmt.Errorf("deploy contract error:%w", err)
		}
		atomic.AddUint64(&bee.nonce, 1)
		ruleAddr = contractAddr.String()
		res, err = bee.invokeContract(bee.normalFrom, constant.RuleManagerContractAddr.Address(), atomic.LoadUint64(&bee.nonce),
			"RegisterRule", rpcx.String(ID), rpcx.String(ruleAddr))
		if err != nil {
			return fmt.Errorf("register rule error:%w", err)
		}
		result := &RegisterResult{}
		err = json.Unmarshal(res.Ret, result)
		if err != nil {
			return err
		}
		err = bee.VotePass(result.ProposalID)
		if err != nil {
			return fmt.Errorf("contract vote error:%w", err)
		}
		atomic.AddUint64(&bee.nonce, 1)
	} else if chainType == "fabric:simple" {
		res, err = bee.invokeContract(bee.normalFrom, constant.RuleManagerContractAddr.Address(), atomic.LoadUint64(&bee.nonce),
			"UnbindRule", rpcx.String(ID), rpcx.String(ruleAddr))
		if err != nil {
			return fmt.Errorf("register rule error:%w", err)
		}
		result := &RegisterResult{}
		err = json.Unmarshal(res.Ret, result)
		if err != nil {
			return err
		}
		err = bee.VotePass(result.ProposalID)
		if err != nil {
			return fmt.Errorf("contract vote error:%w", err)
		}
		atomic.AddUint64(&bee.nonce, 1)
		ruleAddr = "0x00000000000000000000000000000000000000a1"
		res, err = bee.invokeContract(bee.normalFrom, constant.RuleManagerContractAddr.Address(), atomic.LoadUint64(&bee.nonce),
			"BindRule", rpcx.String(ID), rpcx.String(ruleAddr))
		if err != nil {
			return fmt.Errorf("register rule error:%w", err)
		}
		result = &RegisterResult{}
		err = json.Unmarshal(res.Ret, result)
		if err != nil {
			return err
		}
		err = bee.VotePass(result.ProposalID)
		if err != nil {
			return fmt.Errorf("contract vote error:%w", err)
		}
		atomic.AddUint64(&bee.nonce, 1)
	}
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

	tx := &pb.BxhTransaction{
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
	tx := &pb.BxhTransaction{
		From:      bee.normalFrom,
		To:        to,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}

	txHash, err := bee.client.SendTransaction(tx, &rpcx.TransactOpts{
		From:  bee.normalFrom.String(),
		Nonce: normalNo,
	})
	if err != nil {
		return err
	}
	tx.TransactionHash = types.NewHashByStr(txHash)
	go bee.counterReceipt(tx)

	return nil
}

func (bee *bee) sendInterchainTx(i uint64, nonce uint64) error {
	atomic.AddInt64(&sender, 1)
	ibtp := mockIBTP(i, "did:bitxhub:appchain"+bee.normalFrom.String()+":.", "did:bitxhub:appchain"+bee.normalFrom.String()+":.", bee.config.Proof)

	tx := &pb.BxhTransaction{
		From:      bee.normalFrom,
		To:        constant.InterchainContractAddr.Address(),
		Timestamp: time.Now().UnixNano(),
		Extra:     bee.config.Proof,
		IBTP:      ibtp,
	}

	txHash, err := bee.client.SendTransaction(tx, &rpcx.TransactOpts{
		From:  bee.normalFrom.String(),
		Nonce: nonce,
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

func (bee *bee) counterReceipt(tx *pb.BxhTransaction) {
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

func (bee *bee) VotePass(id string) error {
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

func (bee *bee) vote(key crypto.PrivateKey, index uint64, args ...*pb.Arg) (*pb.Receipt, error) {
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

	tx := &pb.BxhTransaction{
		From:      address,
		To:        constant.GovernanceContractAddr.Address(),
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}
	if err != nil {
		return nil, err
	}
	receipt, err := client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From:  address.String(),
		Nonce: index,
	})
	if err != nil {
		return nil, err
	}
	return receipt, nil
}

func (bee *bee) GetChainStatusById(id string) (*pb.Receipt, error) {
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

	tx := &pb.BxhTransaction{
		From:      address,
		To:        constant.AppchainMgrContractAddr.Address(),
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}
	if err != nil {
		return nil, err
	}
	receipt, err := client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From:  bee.adminFrom.String(),
		Nonce: atomic.AddUint64(&index1, 1),
	})
	if err != nil {
		return nil, err
	}
	return receipt, nil
}
