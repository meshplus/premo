package server

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/types"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/meshplus/premo/internal/repo"
)

var (
	nonce1        uint64 = 0
	nonce2        uint64 = 0
	nonce3        uint64 = 0
	defaultRemote        = "localhost:60011"
)

type RegisterResult struct {
	Extra      []byte `json:"extra"`
	ProposalID string `json:"proposal_id"`
}

func initializeAdminNonce() error {
	pk, _, err := repo.Node4Priv()
	if err != nil {
		return err
	}
	client, err := rpcx.New(
		rpcx.WithNodesInfo(&rpcx.NodeInfo{Addr: defaultRemote}),
		rpcx.WithPrivateKey(pk),
	)
	if err != nil {
		return err
	}
	_, node1, err := repo.Node1Priv()
	if err != nil {
		return err
	}
	nonce1, err = client.GetPendingNonceByAccount(node1.String())
	if err != nil {
		return err
	}
	_, node2, err := repo.Node2Priv()
	if err != nil {
		return err
	}
	nonce2, err = client.GetPendingNonceByAccount(node2.String())
	if err != nil {
		return err
	}
	_, node3, err := repo.Node3Priv()
	if err != nil {
		return err
	}
	nonce3, err = client.GetPendingNonceByAccount(node3.String())
	if err != nil {
		return err
	}
	return nil
}

func prepareInterchain(pk crypto.PrivateKey) error {
	err := RegisterAppchain(pk)
	if err != nil {
		return err
	}
	address, err := pk.PublicKey().Address()
	if err != nil {
		return err
	}
	path, err := repo.RulePath()
	if err != nil {
		return err
	}
	err = RegisterRule(pk, path, address.String())
	if err != nil {
		return err
	}
	return nil
}

func RegisterAppchain(pk crypto.PrivateKey) error {
	node0 := &rpcx.NodeInfo{Addr: defaultRemote}
	client, err := rpcx.New(
		rpcx.WithNodesInfo(node0),
		rpcx.WithPrivateKey(pk),
	)
	if err != nil {
		return err
	}
	from, err := pk.PublicKey().Address()
	if err != nil {
		return err
	}
	args := []*pb.Arg{
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain1"),        //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(from.String()),      //public key
	}
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "Register", nil, args...)
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf(string(res.Ret))
	}
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	if err != nil {
		return err
	}
	err = VotePass(result.ProposalID)
	if err != nil {
		return err
	}
	return nil
}

func RegisterRule(pk crypto.PrivateKey, ruleFile string, ChainID string) error {
	node0 := &rpcx.NodeInfo{Addr: defaultRemote}
	client, err := rpcx.New(
		rpcx.WithNodesInfo(node0),
		rpcx.WithPrivateKey(pk),
	)

	// deploy rule
	bytes, err := ioutil.ReadFile(ruleFile)
	if err != nil {
		return err
	}
	addr, err := client.DeployContract(bytes, nil)
	if err != nil {
		return err
	}

	// register rule
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), "RegisterRule", nil, pb.String(ChainID), pb.String(addr.String()))
	if err != nil {
		return err
	}
	if !res.IsSuccess() {
		return fmt.Errorf(string(res.Ret))
	}
	return nil
}

func MockIBTP(from, to *types.Address, index uint64) *pb.IBTP {
	proofHash := sha256.Sum256([]byte("mock ibtp"))
	return &pb.IBTP{
		From:  from.String(),
		To:    to.String(),
		Index: index,
		Type:  pb.IBTP_INTERCHAIN,
		Proof: proofHash[:],
	}
}

func MockContent(funcName string, args [][]byte) []byte {
	content := &pb.Content{
		Func: funcName,
		Args: args,
	}
	bytes, _ := content.Marshal()
	payload := &pb.Payload{
		Encrypted: false,
		Content:   bytes,
	}
	ibtppd, _ := payload.Marshal()
	return ibtppd
}

func VotePass(id string) error {
	return Vote(id, "approve")
}

//Vote `vote` proposal by id and info with four admin
func Vote(id, info string) error {
	key1, _, err := repo.Node1Priv()
	if err != nil {
		return err
	}
	nonce := atomic.AddUint64(&nonce1, 1)
	res, err := vote(key1, nonce-1, rpcx.String(id), rpcx.String(info), rpcx.String("Vote"))
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf("vote err: %s", string(res.Ret))
	}
	key2, _, err := repo.Node2Priv()
	if err != nil {
		return err
	}
	nonce = atomic.AddUint64(&nonce2, 1)
	res, err = vote(key2, nonce-1, rpcx.String(id), rpcx.String(info), rpcx.String("Vote"))
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf("vote err: %s", string(res.Ret))
	}
	key3, _, err := repo.Node3Priv()
	if err != nil {
		return err
	}
	nonce = atomic.AddUint64(&nonce3, 1)
	res, err = vote(key3, nonce-1, rpcx.String(id), rpcx.String(info), rpcx.String("Vote"))
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf("vote err: %s", string(res.Ret))
	}
	return nil
}

// vote `vote` proposal
func vote(key crypto.PrivateKey, nonce uint64, args ...*pb.Arg) (*pb.Receipt, error) {
	client, err := rpcx.New(
		rpcx.WithNodesInfo(&rpcx.NodeInfo{Addr: defaultRemote}),
		rpcx.WithPrivateKey(key),
	)
	if err != nil {
		return nil, err
	}
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
	if err != nil {
		return nil, err
	}
	tx := &pb.Transaction{
		From:      address,
		To:        constant.GovernanceContractAddr.Address(),
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}
	res, err := client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From:  address.String(),
		Nonce: nonce,
	})
	if err != nil {
		return nil, err
	}
	if res.Status == pb.Receipt_FAILED {
		return nil, fmt.Errorf(string(res.Ret))
	}
	return res, nil
}

func TransferFromAdmin(remote, address string, amount uint64) error {
	pk, node1, err := repo.Node1Priv()
	if err != nil {
		return err
	}
	node0 := &rpcx.NodeInfo{Addr: remote}
	client, err := rpcx.New(
		rpcx.WithNodesInfo(node0),
		rpcx.WithPrivateKey(pk),
		rpcx.WithLogger(logrus.New()),
	)
	if err != nil {
		return err
	}
	data := &pb.TransactionData{
		Amount: amount,
	}
	payload, err := data.Marshal()
	if err != nil {
		return err
	}
	tx := &pb.Transaction{
		From:      node1,
		To:        types.NewAddressByStr(address),
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}
	nonce := atomic.AddUint64(&nonce1, 1)
	ret, err := client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From:  node1.String(),
		Nonce: nonce - 1,
	})
	if err != nil {
		return err
	}
	if ret.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf(string(ret.Ret))
	}
	return nil
}
