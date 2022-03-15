package bitxhub

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"github.com/Rican7/retry"
	"github.com/Rican7/retry/backoff"
	"github.com/Rican7/retry/strategy"
	appchainMgr "github.com/meshplus/bitxhub-core/appchain-mgr"
	"github.com/meshplus/bitxhub-core/governance"
	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-kit/types"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/meshplus/premo/internal/repo"
)

var maxDelay int64
var counter int64
var sender int64
var delayer int64
var ibtppd []byte

type bee struct {
	normalPrivKey crypto.PrivateKey
	normalFrom    *types.Address
	client        rpcx.Client
	tps           int
	count         uint64
	nonce         uint64
	ctx           context.Context
	config        *Config
}

type RegisterResult struct {
	Extra      []byte `json:"extra"`
	ProposalID string `json:"proposal_id"`
}

func NewBee(tps int, adminPk crypto.PrivateKey, adminFrom *types.Address, config *Config, ctx context.Context) (*bee, error) {
	normalPk, normalFrom, err := repo.KeyPriv()
	if err != nil {
		return nil, err
	}
	node0 := &rpcx.NodeInfo{Addr: config.BitxhubAddr[0]}
	client, err := rpcx.New(
		rpcx.WithNodesInfo(node0),
		rpcx.WithLogger(log),
		rpcx.WithPrivateKey(normalPk),
	)
	if err != nil {
		return nil, err
	}
	nonce, err := client.GetPendingNonceByAccount(normalFrom.String())
	if err != nil {
		return nil, err
	}
	err = TransferFromAdmin(config, adminPk, adminFrom, normalFrom, "100")
	if err != nil {
		return nil, err
	}
	return &bee{
		client:        client,
		normalPrivKey: normalPk,
		normalFrom:    normalFrom,
		tps:           tps,
		ctx:           ctx,
		config:        config,
		nonce:         nonce,
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
				nonce := atomic.AddUint64(&bee.nonce, 1) - 1
				go func(count, nonce uint64) {
					select {
					case <-bee.ctx.Done():
						return
					default:
						err := retry.Retry(func(attempt uint) error {
							err := bee.sendTx(typ, count, nonce)
							if err != nil {
								return err
							}
							return nil
						}, strategy.Limit(5), strategy.Backoff(backoff.Fibonacci(500*time.Millisecond)))
						if err != nil {
							log.Error(err)
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

func (bee *bee) stop() error {
	err := bee.client.Stop()
	if err != nil {
		return err
	}
	return nil
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

	_, err = bee.client.SendTransaction(tx, &rpcx.TransactOpts{
		Nonce: normalNo,
	})
	if err != nil {
		return err
	}
	return nil
}

func (bee *bee) prepareChain(typ, desc string) error {
	bee.client.SetPrivateKey(bee.normalPrivKey)
	// register chain
	broker := "0x857133c5C69e6Ce66F7AD46F200B9B3573e77582"
	address := "0x00000000000000000000000000000000000000a2"
	if typ == "Fabric V1.4.3" {
		address = "0x00000000000000000000000000000000000000a0"
		broker = "{\"channel_id\":\"mychannel\",\"chaincode_id\":\"broker\",\"broker_version\":\"1\"}"
	} else if typ == "Flato V1.0.3" {
		repoRoot, err := repo.PathRoot()
		contract, err := ioutil.ReadFile(filepath.Join(repoRoot, "rule.wasm"))
		if err != nil {
			return fmt.Errorf("read rule file err %w", err)
		}
		addr, err := bee.client.DeployContract(contract, nil)
		if err != nil {
			return fmt.Errorf("deploy rule err %w", err)
		}
		atomic.AddUint64(&bee.nonce, 1)
		address = addr.String()
	}
	args := []*pb.Arg{
		rpcx.String(bee.normalFrom.String()),     //chainID
		rpcx.String(bee.normalFrom.String()),     //chainName
		rpcx.String(typ),                         //chainType
		rpcx.Bytes([]byte(bee.config.Validator)), //trustRoot
		rpcx.String(broker),                      //broker
		rpcx.String(desc),                        //desc
		rpcx.String(address),                     //masterRuleAddr
		rpcx.String("https://github.com"),        //masterRuleUrl
		rpcx.String(bee.normalFrom.String()),     //adminAddrs
		rpcx.String("reason"),                    //reason
	}
	res, err := bee.client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "RegisterAppchain", &rpcx.TransactOpts{
		From:  bee.normalFrom.String(),
		Nonce: atomic.LoadUint64(&bee.nonce),
	}, args...)
	if err != nil {
		return fmt.Errorf("register appchain error: %w", err)
	}
	atomic.AddUint64(&bee.nonce, 1)
	// vote chain
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	if err != nil || result.ProposalID == "" {
		return fmt.Errorf("vote chain unmarshal error: %w", err)
	}
	err = bee.VotePass(result.ProposalID)
	if err != nil {
		return fmt.Errorf("vote chain error: %w", err)
	}
	res, err = bee.GetChainStatusById(bee.normalPrivKey, bee.normalFrom.String())
	if err != nil {
		return fmt.Errorf("getChainStatus111 error: %w", err)
	}
	atomic.AddUint64(&bee.nonce, 1)
	appchain := &appchainMgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	if err != nil || appchain.Status != governance.GovernanceAvailable {
		return fmt.Errorf("chain error: %w", err)
	}
	//register server
	args = []*pb.Arg{
		rpcx.String(bee.normalFrom.String()),
		rpcx.String("mychannel&transfer"),
		rpcx.String(bee.normalFrom.String()),
		rpcx.String("CallContract"),
		rpcx.String("test"),
		rpcx.Uint64(1),
		rpcx.String(""),
		rpcx.String("test"),
		rpcx.String("reason"),
	}
	res, err = bee.client.InvokeBVMContract(constant.ServiceMgrContractAddr.Address(), "RegisterService", &rpcx.TransactOpts{
		From:  bee.normalFrom.String(),
		Nonce: atomic.LoadUint64(&bee.nonce),
	}, args...)
	if err != nil {
		return fmt.Errorf("register server error %w", err)
	}
	atomic.AddUint64(&bee.nonce, 1)
	//vote server
	result = &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	if err != nil || result.ProposalID == "" {
		return fmt.Errorf("vote server unmarshal error: %w", err)
	}
	err = bee.VotePass(result.ProposalID)
	if err != nil {
		return fmt.Errorf("vote server error: %w", err)
	}
	prepareInterchainTx()
	return nil
}

func (bee *bee) sendTransferTx(to *types.Address, normalNo uint64) error {
	data := &pb.TransactionData{
		Type:   pb.TransactionData_NORMAL,
		VmType: pb.TransactionData_XVM,
		Amount: "0",
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

	_, err = bee.client.SendTransaction(tx, &rpcx.TransactOpts{
		From:  bee.normalFrom.String(),
		Nonce: normalNo,
	})
	if err != nil {
		if strings.Contains(err.Error(), rpcx.ErrBrokenNetwork.Error()) {
			err = bee.sendTransferTx(to, normalNo)
		}
		return err
	}
	return nil
}

func (bee *bee) sendInterchainTx(i uint64, nonce uint64) error {
	atomic.AddInt64(&sender, 1)
	ibtp := mockIBTP(i, "1356:"+bee.normalFrom.String()+":mychannel&transfer", bee.config.Proof)

	tx := &pb.BxhTransaction{
		From:      bee.normalFrom,
		To:        constant.InterchainContractAddr.Address(),
		Timestamp: time.Now().UnixNano(),
		Extra:     bee.config.Proof,
		IBTP:      ibtp,
	}

	_, err := bee.client.SendTransaction(tx, &rpcx.TransactOpts{
		From:  bee.normalFrom.String(),
		Nonce: nonce,
	})
	if err != nil {
		return err
	}
	return nil
}

func prepareInterchainTx() {
	if ibtppd != nil {
		return
	}

	transferAmount := make([]byte, 8)
	binary.BigEndian.PutUint64(transferAmount, 1)
	content := &pb.Content{
		Func: "interchainCharge",
		Args: [][]byte{[]byte("Alice"), []byte("Alice"), transferAmount},
	}

	bytes, _ := content.Marshal()

	payload := &pb.Payload{
		Encrypted: false,
		Content:   bytes,
	}

	ibtppd, _ = payload.Marshal()
}

func mockIBTP(index uint64, from string, proof []byte) *pb.IBTP {
	proofHash := sha256.Sum256(proof)
	return &pb.IBTP{
		From:          from,
		To:            "1356:" + To + ":mychannel&transfer",
		Payload:       ibtppd,
		Index:         index,
		Type:          pb.IBTP_INTERCHAIN,
		TimeoutHeight: 10,
		Proof:         proofHash[:],
	}
}

func (bee *bee) VotePass(id string) error {
	pk1, _, err := repo.Node1Priv()
	if err != nil {
		return err
	}
	_, err = bee.vote(pk1, atomic.AddUint64(&index1, 1), pb.String(id), pb.String("approve"), pb.String("Appchain Pass"))
	if err != nil {
		return err
	}

	pk2, _, err := repo.Node2Priv()
	if err != nil {
		return err
	}
	_, err = bee.vote(pk2, atomic.AddUint64(&index2, 1), pb.String(id), pb.String("approve"), pb.String("Appchain Pass"))
	if err != nil {
		return err
	}

	pk3, _, err := repo.Node3Priv()
	if err != nil {
		return err
	}
	_, err = bee.vote(pk3, atomic.AddUint64(&index3, 1), pb.String(id), pb.String("approve"), pb.String("Appchain Pass"))
	if err != nil {
		return err
	}
	return nil
}

func (bee *bee) vote(key crypto.PrivateKey, index uint64, args ...*pb.Arg) (*pb.Receipt, error) {
	client, err := rpcx.New(
		rpcx.WithNodesInfo(&rpcx.NodeInfo{Addr: bee.config.BitxhubAddr[0]}),
		rpcx.WithLogger(log),
		rpcx.WithPrivateKey(key),
	)
	address, err := key.PublicKey().Address()
	if err != nil {
		return nil, err
	}
	res, err := client.InvokeBVMContract(constant.GovernanceContractAddr.Address(), "Vote", &rpcx.TransactOpts{
		From:  address.String(),
		Nonce: index,
	}, args...)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (bee *bee) GetChainStatusById(pk crypto.PrivateKey, id string) (*pb.Receipt, error) {
	client, err := rpcx.New(
		rpcx.WithNodesInfo(&rpcx.NodeInfo{Addr: bee.config.BitxhubAddr[0]}),
		rpcx.WithLogger(log),
		rpcx.WithPrivateKey(pk),
	)
	if err != nil {
		return nil, err
	}
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "GetAppchain", nil, rpcx.String(id))
	if err != nil {
		return nil, err
	}
	return res, nil
}

func TransferFromAdmin(config *Config, adminPrivKey crypto.PrivateKey, adminFrom *types.Address, address *types.Address, amount string) error {
	node0 := &rpcx.NodeInfo{Addr: config.BitxhubAddr[0]}
	client, err := rpcx.New(
		rpcx.WithNodesInfo(node0),
		rpcx.WithLogger(log),
		rpcx.WithPrivateKey(adminPrivKey),
	)
	if err != nil {
		return err
	}
	data := &pb.TransactionData{
		Amount: amount + "000000000000000000",
	}
	payload, err := data.Marshal()
	if err != nil {
		return err
	}

	tx := &pb.BxhTransaction{
		From:      adminFrom,
		To:        address,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}

	ret, err := client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From:  adminFrom.String(),
		Nonce: atomic.AddUint64(&adminNonce, 1) - 1,
	})
	if err != nil {
		return err
	}
	if ret.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf(string(ret.Ret))
	}
	return nil
}
