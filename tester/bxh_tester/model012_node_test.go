package bxh_tester

import (
	"encoding/json"
	"fmt"
	"sync/atomic"

	"github.com/looplab/fsm"
	"github.com/meshplus/bitxhub-core/governance"
	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/meshplus/premo/internal/repo"
	"github.com/pkg/errors"
)

type NodeType string
type Node struct {
	Pid      string   `toml:"pid" json:"pid"`
	Account  string   `toml:"account" json:"account"`
	NodeType NodeType `toml:"node_type" json:"node_type"`
	// VP Node Info
	VPNodeId uint64                      `toml:"id" json:"id"`
	Primary  bool                        `toml:"primary" json:"primary"`
	Status   governance.GovernanceStatus `toml:"status" json:"status"`
	FSM      *fsm.FSM                    `json:"fsm"`
}

type Model12 struct {
	*Snake
}

//tc：非中继链管理员，注册节点，节点注销失败
func (suite Model12) Test1201_RegisterNodeWithNoRelayAdminIsFail() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	from, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	pid, err := suite.CreatePid()
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(pid),
		rpcx.Uint64(5),
		rpcx.String(from.String()),
		rpcx.String("vpNode"),
		rpcx.String("reason"),
	}
	res, err := client.InvokeBVMContract(constant.NodeManagerContractAddr.Address(), "RegisterNode", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
}

//tc：中继链管理员，注册节点，节点注册失败
func (suite Model12) Test1202_RegisterNodeWithRelayNodeIsSuccess() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	from, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	pid, err := suite.CreatePid()
	suite.Require().Nil(err)
	err = suite.RegisterNode(pid, 5, from.String(), "vpNode", "")
	suite.Require().Nil(err)
	//recover
	err = suite.LogoutNode(from.String())
	suite.Require().Nil(err)
	err = suite.CheckNodeStatus(from.String(), governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

//tc：节点处于unavailable状态注册节点，节点注册成功
func (suite Model12) Test1202_RegisterNodeWithUnavailableNodeIsSuccess() {
	pid, err := suite.CreatePid()
	suite.Require().Nil(err)
	path, err := repo.Node1Path()
	suite.Require().Nil(err)
	pk, err := asym.RestorePrivateKey(path, repo.KeyPassword)
	suite.Require().Nil(err)
	from, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	pk2, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	from2, err := pk2.PublicKey().Address()
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(from2.String()), //nodeAccount
		rpcx.String("vpNode"),       //nodeType
		rpcx.String(pid),            //nodePid
		rpcx.Uint64(5),              //nodeVpId
		rpcx.String(from2.String()), //nodeName
		rpcx.String(""),             //permitStr
		rpcx.String("reason"),       //reason
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.NodeManagerContractAddr.Address(), "RegisterNode", &rpcx.TransactOpts{
		From:  from.String(),
		Nonce: atomic.AddUint64(&nonce1, 1),
	}, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	err = suite.VoteReject(result.ProposalID)
	suite.Require().Nil(err)
	err = suite.CheckNodeStatus(from2.String(), governance.GovernanceUnavailable)
	suite.Require().Nil(err)
	err = suite.RegisterNode(pid, 5, from2.String(), "vpNode", "")
	suite.Require().Nil(err)
	//recover
	err = suite.LogoutNode(from2.String())
	suite.Require().Nil(err)
}

//tc：节点处于available状态注册节点，节点注册失败
func (suite Model12) Test1203_RegisterNodeWithAvailableNodeIsFail() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	from, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	pid, err := suite.CreatePid()
	suite.Require().Nil(err)
	err = suite.RegisterNode(pid, 5, from.String(), "vpNode", "")
	suite.Require().Nil(err)
	err = suite.CheckNodeStatus(from.String(), governance.GovernanceAvailable)
	suite.Require().Nil(err)
	err = suite.RegisterNode(pid, 5, from.String(), "vpNode", "")
	suite.Require().NotNil(err)
	//recover
	err = suite.LogoutNode(from.String())
	suite.Require().Nil(err)
}

//tc：节点用已存在的pid注册节点，节点注册失败
func (suite Model12) Test1204_RegisterNodeWithSamePidIsFail() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	from, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	pid, err := suite.CreatePid()
	suite.Require().Nil(err)
	err = suite.RegisterNode(pid, 5, from.String(), "vpNode", "")
	suite.Require().Nil(err)
	err = suite.RegisterNode(pid, 6, from.String(), "vpNode", "")
	suite.Require().NotNil(err)
	//recover
	err = suite.LogoutNode(from.String())
	suite.Require().Nil(err)
}

//tc：节点用已存在的vpid注册节点，节点注册失败
func (suite Model12) Test1205_RegisterNodeWithSameVPIDIsFail() {
	pk1, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	from1, err := pk1.PublicKey().Address()
	suite.Require().Nil(err)
	pk2, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	from2, err := pk2.PublicKey().Address()
	suite.Require().Nil(err)
	pid1, err := suite.CreatePid()
	suite.Require().Nil(err)
	err = suite.RegisterNode(pid1, 5, from1.String(), "vpNode", "")
	suite.Require().Nil(err)
	pid2, err := suite.CreatePid()
	suite.Require().Nil(err)
	err = suite.RegisterNode(pid2, 5, from2.String(), "vpNode", "")
	suite.Require().NotNil(err)
	//recover
	err = suite.LogoutNode(from1.String())
	suite.Require().Nil(err)
}

//tc：节点用已存在的account注册节点，节点注册失败
//TODO:account now is useless

//tc：非中继链管理员，注销节点，节点注销失败
func (suite Model12) Test1206_LogoutNodeWithNoRelayAdminIsFail() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	from, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	pid, err := suite.CreatePid()
	suite.Require().Nil(err)
	err = suite.RegisterNode(pid, 5, from.String(), "vpNode", "")
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(from.String()),
		rpcx.String("reason"),
	}
	res, err := client.InvokeBVMContract(constant.NodeManagerContractAddr.Address(), "LogoutNode", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	//recover
	err = suite.LogoutNode(from.String())
	suite.Require().Nil(err)
}

//tc：中继链管理员，注销节点，节点注销成功
func (suite Model12) Test1207_LogoutNodeIsSuccess() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	from, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	pid, err := suite.CreatePid()
	suite.Require().Nil(err)
	err = suite.RegisterNode(pid, 5, from.String(), "vpNode", "")
	suite.Require().Nil(err)
	err = suite.LogoutNode(from.String())
	suite.Require().Nil(err)
}

//tc：中继链管理员，注销不存在account的节点，节点注销失败
func (suite Model12) Test1208_LogoutNodeWithNoAccountIsFail() {
	err := suite.LogoutNode("0x79a1215469FaB6f9c63c1816b45183AD3624bE33")
	suite.Require().NotNil(err)
}

//tc：中继链管理员，节点数量4个注销节点，节点注销失败
func (suite Model12) Test1209_LogoutNodeWithFourNodeIsFail() {
	err := suite.LogoutNode("0x79a1215469FaB6f9c63c1816b45183AD3624bE34")
	suite.Require().NotNil(err)
}

//tc：根据节点类型统计可用的节点数量，正确返回节点数量
func (suite Model12) Test1210_CountAvailableNodesByTypeIsSuccess() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.NodeManagerContractAddr.Address(), "CountAvailableNodes", nil, rpcx.String("vpNode"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	suite.Require().Equal("4", string(res.Ret))
}

//tc：根据不存在的节点类型统计可用的节点数量，查询失败
func (suite Model12) Test1211_CountAvailableNodesByTypeWithNoTypeIsFail() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.NodeManagerContractAddr.Address(), "CountAvailableNodes", nil, rpcx.String("vpNode111"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	suite.Require().Equal("0", string(res.Ret))
}

//tc：查询所有节点信息，正确返回节点信息
func (suite Model12) Test1212_NodesIsSuccess() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.NodeManagerContractAddr.Address(), "Nodes", nil)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	suite.Require().NotNil(string(res.Ret))
}

//tc：根据account判断是否是可用节点，正确返回节点状态
func (suite Model12) Test1213_IsAvailableIsSuccess() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.NodeManagerContractAddr.Address(), "IsAvailable", nil, rpcx.String("0x79a1215469FaB6f9c63c1816b45183AD3624bE34"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	suite.Require().Equal("true", string(res.Ret))
}

//tc：根据错误的account判断是否是可用节点，正确返回节点状态
func (suite Model12) Test1214_IsAvailableWithNoAccountIsFail() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.NodeManagerContractAddr.Address(), "IsAvailable", nil, rpcx.String("111"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	suite.Require().Equal("false", string(res.Ret))
}

func (suite Snake) RegisterNode(nodePid string, nodeVpId uint64, nodeAccount, nodeType, permit string) error {
	path, err := repo.Node1Path()
	if err != nil {
		return err
	}
	pk, err := asym.RestorePrivateKey(path, repo.KeyPassword)
	if err != nil {
		return err
	}
	from, err := pk.PublicKey().Address()
	if err != nil {
		return err
	}
	args := []*pb.Arg{
		rpcx.String(nodeAccount), //nodeAccount
		rpcx.String(nodeType),    //nodeType
		rpcx.String(nodePid),     //nodePid
		rpcx.Uint64(nodeVpId),    //nodeVpId
		rpcx.String(nodeAccount), //nodeName
		rpcx.String(permit),      //permitStr
		rpcx.String("reason"),    //reason
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.NodeManagerContractAddr.Address(), "RegisterNode", &rpcx.TransactOpts{
		From:  from.String(),
		Nonce: atomic.AddUint64(&nonce1, 1),
	}, args...)
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return errors.New(string(res.Ret))
	}
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	if result.ProposalID == "" {
		return nil
	}
	err = suite.VotePass(result.ProposalID)
	if err != nil {
		return err
	}
	return nil
}

func (suite Snake) LogoutNode(account string) error {
	path, err := repo.Node1Path()
	if err != nil {
		return err
	}
	pk, err := asym.RestorePrivateKey(path, repo.KeyPassword)
	if err != nil {
		return err
	}
	from, err := pk.PublicKey().Address()
	if err != nil {
		return err
	}
	args := []*pb.Arg{
		rpcx.String(account),
		rpcx.String("reason"),
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.NodeManagerContractAddr.Address(), "LogoutNode", &rpcx.TransactOpts{
		From:  from.String(),
		Nonce: atomic.AddUint64(&nonce1, 1),
	}, args...)
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return errors.New(string(res.Ret))
	}
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	if result.ProposalID == "" {
		return nil
	}
	err = suite.VotePass(result.ProposalID)
	if err != nil {
		return err
	}
	return nil
}

func (suite Snake) CheckNodeStatus(account string, status governance.GovernanceStatus) error {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	if err != nil {
		return err
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.NodeManagerContractAddr.Address(), "GetNode", nil, rpcx.String(account))
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return errors.New(string(res.Ret))
	}
	node := &Node{}
	err = json.Unmarshal(res.Ret, node)
	if err != nil {
		return err
	}
	if node.Status != status {
		return errors.New(fmt.Sprintf("expect status is %s, but got %s", status, node.Status))
	}
	return nil
}
