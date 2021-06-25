package bxh_tester

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"sync/atomic"
	"time"

	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-kit/types"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/meshplus/premo/internal/repo"
	"github.com/pkg/errors"
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
	//client0   rpcx.ChainClient
	client    rpcx.Client
	from      *types.Address
	fromIndex uint64
	pk        crypto.PrivateKey
	toIndex   uint64
	to        *types.Address
}
type RegisterResult struct {
	Extra      []byte `json:"extra"`
	ProposalID string `json:"proposal_id"`
}

var nonce2 uint64
var nonce3 uint64
var nonce4 uint64

// SetupTest init
func (suite *Snake) SetupTest() {
	suite.T().Parallel()
}

func (suite *Snake) SetupSuite() {
	node2, err := repo.Node2Path()
	suite.Require().Nil(err)

	key2, err := asym.RestorePrivateKey(node2, repo.KeyPassword)
	suite.Require().Nil(err)

	node2Addr, err := key2.PublicKey().Address()
	suite.Require().Nil(err)

	node3, err := repo.Node3Path()
	suite.Require().Nil(err)

	key3, err := asym.RestorePrivateKey(node3, repo.KeyPassword)
	suite.Require().Nil(err)

	node3Addr, err := key3.PublicKey().Address()
	suite.Require().Nil(err)

	node4, err := repo.Node4Path()
	suite.Require().Nil(err)

	key4, err := asym.RestorePrivateKey(node4, repo.KeyPassword)
	suite.Require().Nil(err)

	node4Addr, err := key4.PublicKey().Address()
	suite.Require().Nil(err)

	suite.sendTransaction(key2)
	suite.sendTransaction(key3)
	suite.sendTransaction(key4)

	nonce, err := suite.client.GetPendingNonceByAccount(node2Addr.String())
	suite.Require().Nil(err)
	nonce2 = nonce - 1

	nonce, err = suite.client.GetPendingNonceByAccount(node3Addr.String())
	suite.Require().Nil(err)
	nonce3 = nonce - 1

	nonce, err = suite.client.GetPendingNonceByAccount(node4Addr.String())
	suite.Require().Nil(err)
	nonce4 = nonce - 1
}

func (suite *Snake) RegisterAppchain() (crypto.PrivateKey, string, error) {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	if err != nil {
		return nil, "", err
	}
	pubAddress, err := pk.PublicKey().Address()
	if err != nil {
		return nil, "", err
	}
	client := suite.NewClient(pk)
	bytes, err := pk.PublicKey().Bytes()
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(bytes)

	args := []*pb.Arg{
		rpcx.String("appchain" + pubAddress.String()),                       //method
		rpcx.String("/ipfs/QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"), //docAddr
		rpcx.String("QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),       //docHash
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain"),         //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "Register", nil, args...)
	if err != nil {
		return nil, "", err
	}
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	if err != nil {
		return nil, "", err
	}
	err = suite.VotePass(result.ProposalID)
	if err != nil {
		return nil, "", err
	}
	return pk, string(result.Extra), nil
}

func (suite *Snake) RegisterRule(pk crypto.PrivateKey, ruleFile string, ChainID string) {
	client := suite.NewClient(pk)

	// deploy rule
	bytes, err := ioutil.ReadFile(ruleFile)
	suite.Require().Nil(err)
	addr, err := client.DeployContract(bytes, nil)
	suite.Require().Nil(err)

	// register rule
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), "RegisterRule", nil, pb.String(ChainID), pb.String(addr.String()))
	suite.Require().Nil(err)
	suite.Require().True(res.IsSuccess())
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	suite.Require().Nil(err)
	err = suite.VotePass(result.ProposalID)
	suite.Require().Nil(err)
}

func (suite *Snake) NewClient(pk crypto.PrivateKey) *rpcx.ChainClient {
	node0 := &rpcx.NodeInfo{Addr: cfg.addrs[0]}
	client, err := rpcx.New(
		rpcx.WithNodesInfo(node0),
		rpcx.WithLogger(cfg.logger),
		rpcx.WithPrivateKey(pk),
	)
	suite.Require().Nil(err)
	from, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	err = suite.TransferFromAdmin(from.String(), "1")
	suite.Require().Nil(err)
	return client
}

func (suite *Snake) VotePass(id string) error {
	node2, err := repo.Node2Path()
	if err != nil {
		return err
	}

	key, err := asym.RestorePrivateKey(node2, repo.KeyPassword)
	if err != nil {
		return err
	}

	_, err = suite.vote(key, atomic.AddUint64(&nonce2, 1), pb.String(id), pb.String("approve"), pb.String("Appchain Pass"))
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

	_, err = suite.vote(key, atomic.AddUint64(&nonce3, 1), pb.String(id), pb.String("approve"), pb.String("Appchain Pass"))

	if err != nil {
		return err
	}

	node4, err := repo.Node4Path()
	if err != nil {
		return err
	}

	key, err = asym.RestorePrivateKey(node4, repo.KeyPassword)
	if err != nil {
		return err
	}

	_, err = suite.vote(key, atomic.AddUint64(&nonce4, 1), pb.String(id), pb.String("approve"), pb.String("Appchain Pass"))
	if err != nil {
		return err
	}

	return nil
}

func (suite *Snake) VoteReject(id string) error {
	node2, err := repo.Node2Path()
	if err != nil {
		return err
	}

	key, err := asym.RestorePrivateKey(node2, repo.KeyPassword)
	if err != nil {
		return err
	}

	_, err = suite.vote(key, atomic.AddUint64(&nonce2, 1), pb.String(id), pb.String("reject"), pb.String("Appchain Pass"))
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

	_, err = suite.vote(key, atomic.AddUint64(&nonce3, 1), pb.String(id), pb.String("reject"), pb.String("Appchain Pass"))

	if err != nil {
		return err
	}

	node4, err := repo.Node4Path()
	if err != nil {
		return err
	}

	key, err = asym.RestorePrivateKey(node4, repo.KeyPassword)
	if err != nil {
		return err
	}

	_, err = suite.vote(key, atomic.AddUint64(&nonce4, 1), pb.String(id), pb.String("reject"), pb.String("Appchain Pass"))
	if err != nil {
		return err
	}

	return nil
}

func (suite *Snake) vote(key crypto.PrivateKey, nonce uint64, args ...*pb.Arg) (*pb.Receipt, error) {
	client, err := rpcx.New(
		rpcx.WithNodesInfo(&rpcx.NodeInfo{Addr: cfg.addrs[0]}),
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
	res, err := client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		Nonce: nonce,
	})
	if err != nil {
		return nil, err
	}
	if res.Status == pb.Receipt_FAILED {
		return nil, errors.New(string(res.Ret))
	}
	return res, nil
}

func (suite *Snake) sendTransaction(pk crypto.PrivateKey) {
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

	tx := &pb.BxhTransaction{
		From:      from,
		To:        suite.to,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}

	res, err := client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
}

func (suite *Snake) GetChainStatusById(id string) (*pb.Receipt, error) {
	key, err := asym.GenerateKeyPair(crypto.Secp256k1)
	if err != nil {
		return nil, err
	}
	client := suite.NewClient(key)
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
	res, err := client.SendTransactionWithReceipt(tx, nil)
	if err != nil {
		return nil, err
	}
	if res.Status == pb.Receipt_FAILED {
		return nil, errors.New(string(res.Ret))
	}
	return res, nil
}

func (suite Snake) TransferFromAdmin(address string, amount string) error {
	node4, err := repo.Node4Path()
	if err != nil {
		return err
	}
	pk, err := asym.RestorePrivateKey(node4, repo.KeyPassword)
	if err != nil {
		return err
	}
	from, err := pk.PublicKey().Address()
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
		From:      from,
		To:        types.NewAddressByStr(address),
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}

	ret, err := client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		Nonce: atomic.AddUint64(&nonce4, 1),
	})
	if err != nil {
		return err
	}
	if ret.Status != pb.Receipt_SUCCESS {
		return errors.New(string(ret.Ret))
	}
	return nil
}
