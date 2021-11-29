package bxh_tester

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync/atomic"
	"time"

	"github.com/looplab/fsm"
	appchain_mgr "github.com/meshplus/bitxhub-core/appchain-mgr"
	"github.com/meshplus/bitxhub-core/governance"
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

type Rule struct {
	Address string                      `json:"address"`
	ChainId string                      `json:"chain_id"`
	Status  governance.GovernanceStatus `json:"status"`
	FSM     *fsm.FSM                    `json:"fsm"`
}

type ProposalType string
type ProposalStatus string
type EndReason string
type RoleType string
type Ballot struct {
	VoterAddr string `json:"voter_addr"`
	Approve   string `json:"approve"`
	Num       uint64 `json:"num"`
	Reason    string `json:"reason"`
	VoteTime  int64  `json:"vote_time"`
}
type Role struct {
	ID       string   `toml:"id" json:"id"`
	RoleType RoleType `toml:"role_type" json:"role_type"`
	// 	GovernanceAdmin info
	Weight uint64 `json:"weight" toml:"weight"`
	// AuditAdmin info
	NodePid string `toml:"pid" json:"pid"`
	// Appchain info
	AppchainID string                      `toml:"appchain_id" json:"appchain_id"`
	Status     governance.GovernanceStatus `toml:"status" json:"status"`
	FSM        *fsm.FSM                    `json:"fsm"`
}
type Proposal struct {
	Id            string                      `json:"id"`
	Typ           ProposalType                `json:"typ"`
	Status        ProposalStatus              `json:"status"`
	ObjId         string                      `json:"obj_id"`
	ObjLastStatus governance.GovernanceStatus `json:"obj_last_status"`
	// ballot information: voter address -> ballot
	BallotMap              map[string]Ballot    `json:"ballot_map"`
	ApproveNum             uint64               `json:"approve_num"`
	AgainstNum             uint64               `json:"against_num"`
	ElectorateList         []*Role              `json:"electorate_list"`
	InitialElectorateNum   uint64               `json:"initial_electorate_num"`
	AvailableElectorateNum uint64               `json:"available_electorate_num"`
	ThresholdElectorateNum uint64               `json:"threshold_electorate_num"`
	EventType              governance.EventType `json:"event_type"`
	EndReason              EndReason            `json:"end_reason"`
	LockProposalId         string               `json:"lock_proposal_id"`
	IsSpecial              bool                 `json:"is_special"`
	IsSuperAdminVoted      bool                 `json:"is_super_admin_voted"`
	SubmitReason           string               `json:"submit_reason"`
	WithdrawReason         string               `json:"withdraw_reason"`
	CreateTime             int64                `json:"create_time"`
	Extra                  []byte               `json:"extra"`
}

var nonce1 uint64
var nonce2 uint64
var nonce3 uint64
var nonce4 uint64

func (suite *Snake) SetupSuite() {
	node1, err := repo.Node1Path()
	suite.Require().Nil(err)

	key1, err := asym.RestorePrivateKey(node1, repo.KeyPassword)
	suite.Require().Nil(err)

	node1Addr, err := key1.PublicKey().Address()
	suite.Require().Nil(err)

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

	suite.sendTransaction(key1)
	suite.sendTransaction(key2)
	suite.sendTransaction(key3)
	suite.sendTransaction(key4)

	nonce, err := suite.client.GetPendingNonceByAccount(node1Addr.String())
	suite.Require().Nil(err)
	nonce1 = nonce - 1

	nonce, err = suite.client.GetPendingNonceByAccount(node2Addr.String())
	suite.Require().Nil(err)
	nonce2 = nonce - 1

	nonce, err = suite.client.GetPendingNonceByAccount(node3Addr.String())
	suite.Require().Nil(err)
	nonce3 = nonce - 1

	nonce, err = suite.client.GetPendingNonceByAccount(node4Addr.String())
	suite.Require().Nil(err)
	nonce4 = nonce - 1
}

func (suite *Snake) RegisterAppchain(pk crypto.PrivateKey, name, address string) error {
	client := suite.NewClient(pk)
	from, err := pk.PublicKey().Address()
	if err != nil {
		return err
	}
	args := []*pb.Arg{
		rpcx.String(from.String()),  //chainID
		rpcx.String(name),           //chainName
		rpcx.String("Flato V1.0.3"), //chainType
		rpcx.Bytes([]byte("")),      //trustRoot
		rpcx.String("0x857133c5C69e6Ce66F7AD46F200B9B3573e77582"), //broker
		rpcx.String("desc"),               //desc
		rpcx.String(address),              //masterRuleAddr
		rpcx.String("https://github.com"), //masterRuleUrl
		rpcx.String(from.String()),        //adminAddrs
		rpcx.String("reason"),             //reason
	}
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "RegisterAppchain", nil, args...)
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return errors.New(string(res.Ret))
	}
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	if err != nil {
		return err
	}
	err = suite.VotePass(result.ProposalID)
	if err != nil {
		return err
	}
	return nil
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
	node1, err := repo.Node1Path()
	if err != nil {
		return err
	}

	key, err := asym.RestorePrivateKey(node1, repo.KeyPassword)
	if err != nil {
		return err
	}

	_, err = suite.vote(key, atomic.AddUint64(&nonce1, 1), rpcx.String(id), rpcx.String("approve"), rpcx.String("Appchain Pass"))

	node2, err := repo.Node2Path()
	if err != nil {
		return err
	}

	key, err = asym.RestorePrivateKey(node2, repo.KeyPassword)
	if err != nil {
		return err
	}

	_, err = suite.vote(key, atomic.AddUint64(&nonce2, 1), rpcx.String(id), rpcx.String("approve"), rpcx.String("Appchain Pass"))

	node3, err := repo.Node3Path()
	if err != nil {
		return err
	}

	key, err = asym.RestorePrivateKey(node3, repo.KeyPassword)
	if err != nil {
		return err
	}

	_, err = suite.vote(key, atomic.AddUint64(&nonce3, 1), rpcx.String(id), rpcx.String("approve"), rpcx.String("Appchain Pass"))

	node4, err := repo.Node4Path()
	if err != nil {
		return err
	}

	key, err = asym.RestorePrivateKey(node4, repo.KeyPassword)
	if err != nil {
		return err
	}

	_, err = suite.vote(key, atomic.AddUint64(&nonce4, 1), rpcx.String(id), rpcx.String("approve"), rpcx.String("Appchain Pass"))
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

	_, err = suite.vote(key, atomic.AddUint64(&nonce2, 1), rpcx.String(id), rpcx.String("reject"), rpcx.String("Appchain Pass"))
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

	_, err = suite.vote(key, atomic.AddUint64(&nonce3, 1), rpcx.String(id), rpcx.String("reject"), rpcx.String("Appchain Pass"))

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

	_, err = suite.vote(key, atomic.AddUint64(&nonce4, 1), rpcx.String(id), rpcx.String("reject"), rpcx.String("Appchain Pass"))
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
		From:  address.String(),
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

func (suite *Snake) GetChainStatusByName(name string) (governance.GovernanceStatus, error) {
	key, err := asym.GenerateKeyPair(crypto.Secp256k1)
	if err != nil {
		return "", err
	}
	client := suite.NewClient(key)
	address, err := key.PublicKey().Address()
	if err != nil {
		return "", err
	}
	args := []*pb.Arg{
		rpcx.String(name),
	}
	invokePayload := &pb.InvokePayload{
		Method: "GetAppchainByName",
		Args:   args,
	}
	payload, err := invokePayload.Marshal()
	if err != nil {
		return "", err
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
		return "", err
	}
	res, err := client.SendTransactionWithReceipt(tx, nil)
	if err != nil {
		return "", err
	}
	if res.Status == pb.Receipt_FAILED {
		return "", errors.New(string(res.Ret))
	}
	appchain := appchain_mgr.Appchain{}
	err = json.Unmarshal(res.Ret, &appchain)
	if err != nil {
		return "", err
	}
	return appchain.Status, nil
}

func (suite Snake) GetChainIDByName(name string) (string, error) {
	key, err := asym.GenerateKeyPair(crypto.Secp256k1)
	if err != nil {
		return "", err
	}
	client := suite.NewClient(key)
	address, err := key.PublicKey().Address()
	if err != nil {
		return "", err
	}
	args := []*pb.Arg{
		rpcx.String(name),
	}
	invokePayload := &pb.InvokePayload{
		Method: "GetAppchainByName",
		Args:   args,
	}
	payload, err := invokePayload.Marshal()
	if err != nil {
		return "", err
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
		return "", err
	}
	res, err := client.SendTransactionWithReceipt(tx, nil)
	if err != nil {
		return "", err
	}
	appchain := appchain_mgr.Appchain{}
	err = json.Unmarshal(res.Ret, &appchain)
	if err != nil {
		return "", err
	}
	return appchain.ID, nil
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
		From:  from.String(),
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

func (suite Snake) GetChainID(pk crypto.PrivateKey) string {
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	return "Chain" + address.String()
}

func (suite Snake) GetServerID(pk crypto.PrivateKey) string {
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	return address.String()
}

func (suite Snake) DeploySimpleRule() (string, error) {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	if err != nil {
		return "", err
	}
	client := suite.NewClient(pk)
	contract, err := ioutil.ReadFile("testdata/simple_rule.wasm")
	if err != nil {
		return "", err
	}
	address, err := client.DeployContract(contract, nil)
	if err != nil {
		return "", err
	}
	return address.String(), nil
}

func (suite *Snake) GetRuleStatus(pk crypto.PrivateKey, ChainID string, contractAddr string) (status governance.GovernanceStatus, err error) {
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.String(contractAddr),
	}
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), "GetRuleByAddr", nil, args...)
	if err != nil {
		return "", err
	}
	if res.Status == pb.Receipt_FAILED {
		return "", errors.New(string(res.Ret))
	}
	rule := &Rule{}
	err = json.Unmarshal(res.Ret, rule)
	if err != nil {
		return "", err
	}
	return rule.Status, nil
}

func (suite *Snake) CheckChainStatus(name string, expectStatus governance.GovernanceStatus) error {
	status, err := suite.GetChainStatusByName(name)
	if err != nil {
		return err
	}
	if expectStatus != status {
		return errors.New(fmt.Sprintf("expect status is %s ,but get status %s", expectStatus, status))
	}
	return nil
}

func (suite *Snake) CheckRuleStatus(pk crypto.PrivateKey, chainID, address string, expectStatus governance.GovernanceStatus) error {
	status, err := suite.GetRuleStatus(pk, chainID, address)
	if err != nil {
		return err
	}
	if expectStatus != status {
		return errors.New(fmt.Sprintf("expect status is %s ,but get status %s", expectStatus, status))
	}
	return nil
}

func (suite Snake) CountAppchains() (string, error) {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	if err != nil {
		return "", err
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "CountAppchains", nil)
	if err != nil {
		return "", err
	}
	if res.Status == pb.Receipt_FAILED {
		return "", errors.New(string(res.Ret))
	}
	return string(res.Ret), nil
}

//func (suite Snake) GetProposalByID(id string) (Proposal, error) {
//	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
//	if err != nil {
//		return "", err
//	}
//	client := suite.NewClient(pk)
//
//}
