package bxh_tester

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"

	crypto2 "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/looplab/fsm"
	"github.com/meshplus/bitxhub-core/governance"
	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-kit/crypto/asym/ecdsa"
	"github.com/meshplus/bitxhub-kit/types"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/meshplus/premo/internal/repo"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

var cfg = &config{
	addrs: []string{
		"localhost:60011",
		"localhost:60012",
		"localhost:60013",
		"localhost:60014",
	},
	logger: logrus.New(),
}

type config struct {
	addrs  []string
	logger rpcx.Logger
}
type Snake struct {
	suite.Suite
}
type RegisterResult struct {
	Extra      []byte `json:"extra"`
	ProposalID string `json:"proposal_id"`
}

type Rule struct {
	Address string                      `json:"address"`
	ChainId string                      `json:"chain_id"`
	Status  governance.GovernanceStatus `json:"status"`
	FSM     *fsm.FSM                    `json:"fsm"`
}
type RoleType string
type Role struct {
	ID         string                      `toml:"id" json:"id"`
	RoleType   RoleType                    `toml:"role_type" json:"role_type"`
	Weight     uint64                      `json:"weight" toml:"weight"`
	NodePid    string                      `toml:"pid" json:"pid"`
	AppchainID string                      `toml:"appchain_id" json:"appchain_id"`
	Status     governance.GovernanceStatus `toml:"status" json:"status"`
	FSM        *fsm.FSM                    `json:"fsm"`
}

var nonce1 uint64
var nonce2 uint64
var nonce3 uint64
var nonce4 uint64

func (suite *Snake) SetupSuite() {
	key1, node1Addr, err := repo.Node1Priv()
	suite.Require().Nil(err)
	key2, node2Addr, err := repo.Node2Priv()
	suite.Require().Nil(err)
	key3, node3Addr, err := repo.Node3Priv()
	suite.Require().Nil(err)
	key4, node4Addr, err := repo.Node4Priv()
	suite.Require().Nil(err)
	suite.SendTransaction(key1)
	suite.SendTransaction(key2)
	suite.SendTransaction(key3)
	suite.SendTransaction(key4)
	pk, _, err := repo.KeyPriv()
	node0 := &rpcx.NodeInfo{Addr: cfg.addrs[0]}
	client, err := rpcx.New(
		rpcx.WithNodesInfo(node0),
		rpcx.WithLogger(cfg.logger),
		rpcx.WithPrivateKey(pk),
	)
	suite.Require().Nil(err)
	nonce, err := client.GetPendingNonceByAccount(node1Addr.String())
	suite.Require().Nil(err)
	nonce1 = nonce - 1
	nonce, err = client.GetPendingNonceByAccount(node2Addr.String())
	suite.Require().Nil(err)
	nonce2 = nonce - 1
	nonce, err = client.GetPendingNonceByAccount(node3Addr.String())
	suite.Require().Nil(err)
	nonce3 = nonce - 1
	nonce, err = client.GetPendingNonceByAccount(node4Addr.String())
	suite.Require().Nil(err)
	nonce4 = nonce - 1

	rand.Seed(time.Now().UnixNano())
}

// NewClient return client by privateKey
func (suite *Snake) NewClient(pk crypto.PrivateKey) *rpcx.ChainClient {
	node0 := &rpcx.NodeInfo{Addr: cfg.addrs[0]}
	client, err := rpcx.New(
		rpcx.WithNodesInfo(node0),
		rpcx.WithLogger(cfg.logger),
		rpcx.WithPrivateKey(pk),
		rpcx.WithTimeoutLimit(time.Second*60),
	)
	suite.Require().Nil(err)
	from, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	err = suite.TransferFromAdmin(from.String(), "1")
	suite.Require().Nil(err)
	return client
}

// VotePass vote pass proposal by id with four admin
func (suite *Snake) VotePass(id string) error {
	return suite.Vote(id, "approve")
}

// VoteReject vote reject proposal by id with four admin
func (suite *Snake) VoteReject(id string) error {
	return suite.Vote(id, "reject")
}

//Vote `vote` proposal by id and info with four admin
func (suite *Snake) Vote(id, info string) error {
	key1, _, err := repo.Node1Priv()
	if err != nil {
		return err
	}
	_, err = suite.vote(key1, atomic.AddUint64(&nonce1, 1), rpcx.String(id), rpcx.String(info), rpcx.String("Vote"))
	key2, _, err := repo.Node2Priv()
	if err != nil {
		return err
	}
	_, err = suite.vote(key2, atomic.AddUint64(&nonce2, 1), rpcx.String(id), rpcx.String(info), rpcx.String("Vote"))
	key3, _, err := repo.Node3Priv()
	if err != nil {
		return err
	}
	_, err = suite.vote(key3, atomic.AddUint64(&nonce3, 1), rpcx.String(id), rpcx.String(info), rpcx.String("Vote"))
	key4, _, err := repo.Node4Priv()
	if err != nil {
		return err
	}
	_, err = suite.vote(key4, atomic.AddUint64(&nonce4, 1), rpcx.String(id), rpcx.String(info), rpcx.String("Vote"))
	return nil
}

// vote `vote` proposal
func (suite *Snake) vote(key crypto.PrivateKey, nonce uint64, args ...*pb.Arg) (*pb.Receipt, error) {
	client, err := rpcx.New(
		rpcx.WithNodesInfo(&rpcx.NodeInfo{Addr: cfg.addrs[0]}),
		rpcx.WithLogger(cfg.logger),
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

// SendTransaction send a normal tx to bitxhub
func (suite *Snake) SendTransaction(pk crypto.PrivateKey) {
	node0 := &rpcx.NodeInfo{Addr: cfg.addrs[0]}
	client, err := rpcx.New(
		rpcx.WithNodesInfo(node0),
		rpcx.WithLogger(cfg.logger),
		rpcx.WithPrivateKey(pk),
	)
	suite.Require().Nil(err)
	from, err := pk.PublicKey().Address()
	data := &pb.TransactionData{
		Amount: "1",
	}
	payload, err := data.Marshal()
	suite.Require().Nil(err)
	_, to, err := repo.KeyPriv()
	suite.Require().Nil(err)
	tx := &pb.BxhTransaction{
		From:      from,
		To:        to,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}
	res, err := client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
}

// TransferFromAdmin transfer amount from admin
func (suite *Snake) TransferFromAdmin(address string, amount string) error {
	pk, node4Addr, err := repo.Node4Priv()
	if err != nil {
		return err
	}
	node0 := &rpcx.NodeInfo{Addr: cfg.addrs[0]}
	client, err := rpcx.New(
		rpcx.WithNodesInfo(node0),
		rpcx.WithLogger(cfg.logger),
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
		From:      node4Addr,
		To:        types.NewAddressByStr(address),
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}
	ret, err := client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From:  node4Addr.String(),
		Nonce: atomic.AddUint64(&nonce4, 1),
	})
	if err != nil {
		return err
	}
	if ret.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf(string(ret.Ret))
	}
	return nil
}

// MockIBTP mock a ibtp
func (suite *Snake) MockIBTP(index uint64, from, to string, typ pb.IBTP_Type, proof []byte) *pb.IBTP {
	proofHash := sha256.Sum256(proof)
	return &pb.IBTP{
		From:          from,
		To:            to,
		Index:         index,
		Type:          typ,
		TimeoutHeight: 10,
		Proof:         proofHash[:],
	}
}

// MockContent mock a content
func (suite *Snake) MockContent(funcName string, args [][]byte) []byte {
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

// MockResult mock a result
func (suite *Snake) MockResult(data [][]byte) []byte {
	result := &pb.Result{Data: data}
	bytes, _ := result.Marshal()
	payload := &pb.Payload{
		Encrypted: false,
		Content:   bytes,
	}
	ibtppd, _ := payload.Marshal()
	return ibtppd
}

// MockPid mock a pid
func (suite *Snake) MockPid() (string, error) {
	pk, err := asym.GenerateKeyPair(crypto.ECDSA_P256)
	if err != nil {
		return "", err
	}
	data, err := pk.Bytes()
	if err != nil {
		return "", err
	}
	key, err := ecdsa.UnmarshalPrivateKey(data, crypto.ECDSA_P256)
	if err != nil {
		return "", err
	}
	_, k, err := crypto2.KeyPairFromStdKey(key.K)
	if err != nil {
		return "", err
	}
	pid, err := peer.IDFromPublicKey(k)
	if err != nil {
		return "", err
	}
	return pid.String(), nil
}
