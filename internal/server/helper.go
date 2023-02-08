package server

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"

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
	err = RegisterServer(pk)
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
	bytes, err := pk.PublicKey().Bytes()
	if err != nil {
		return err
	}
	args := []*pb.Arg{
		rpcx.String(from.String()),   //chainID
		rpcx.String(from.String()),   //chainName
		rpcx.Bytes(bytes),            //pubKey
		rpcx.String("Fabric V1.4.3"), //chainType
		rpcx.Bytes([]byte("")),       //trustRoot
		rpcx.String("{\"channel_id\":\"mychannel\",\"chaincode_id\":\"broker\",\"broker_version\":\"1\"}"), //broker
		rpcx.String("desc"),               //desc
		rpcx.String(HappyRuleAddr),        //masterRuleAddr
		rpcx.String("https://github.com"), //masterRuleUrl
		rpcx.String(from.String()),        //adminAddrs
		rpcx.String("reason"),             //reason
	}
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "RegisterAppchain", nil, args...)
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

func RegisterServer(pk crypto.PrivateKey) error {
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
		rpcx.String(from.String()),
		rpcx.String(from.String()),
		rpcx.String(from.String()),
		rpcx.String("CallContract"),
		rpcx.String("test"),
		rpcx.Uint64(1),
		rpcx.String(""),
		rpcx.String("test"),
		rpcx.String("reason"),
	}
	res, err := client.InvokeBVMContract(constant.ServiceMgrContractAddr.Address(), "RegisterService", nil, args...)
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

func MockIBTP(from, to *types.Address, index uint64) *pb.IBTP {
	proofHash := sha256.Sum256([]byte("mock ibtp"))
	return &pb.IBTP{
		From:          "1356:" + from.String() + ":" + from.String(),
		To:            "1356:" + to.String() + ":" + to.String(),
		Index:         index,
		Type:          pb.IBTP_INTERCHAIN,
		TimeoutHeight: 10,
		Proof:         proofHash[:],
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
	tx := &pb.BxhTransaction{
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

func TransferFromAdmin(remote, address, amount string) error {
	pk, node1, err := repo.Node1Priv()
	if err != nil {
		return err
	}
	node0 := &rpcx.NodeInfo{Addr: remote}
	client, err := rpcx.New(
		rpcx.WithNodesInfo(node0),
		rpcx.WithPrivateKey(pk),
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
		From:      node1,
		To:        types.NewAddressByStr(address),
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}
	nonce := atomic.AddUint64(&nonce1, 1)
	ret, err := client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From:    node1.String(),
		Nonce:   nonce - 1,
		PrivKey: nil,
	})
	if err != nil {
		return err
	}
	if ret.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf(string(ret.Ret))
	}
	return nil
}
