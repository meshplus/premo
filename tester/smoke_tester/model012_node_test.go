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
)

const (
	RegisterNode = "RegisterNode"
	UpdateNode   = "UpdateNode"
	LogoutNode   = "LogoutNode"
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

//tc：中继链管理员，注册治理节点，节点注册成功
func (suite *Model12) Test1201_RegisterVpNodeWithRelayNodeIsSuccess() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	pid, err := suite.MockPid()
	suite.Require().Nil(err)
	err = suite.RegisterNode(from.String(), "vpNode", pid, 5, from.String(), "")
	suite.Require().Nil(err)
	//recover
	err = suite.LogoutNode(from.String())
	suite.Require().Nil(err)
	err = suite.CheckNodeStatus(from.String(), governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

//tc：治理节点未注册，节点注册成功
func (suite *Model12) Test1202_RegisterVpNodeWithNoRegisterIsSuccess() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	pid, err := suite.MockPid()
	suite.Require().Nil(err)
	err = suite.RegisterNode(from.String(), "vpNode", pid, 5, from.String(), "")
	suite.Require().Nil(err)
	//recover
	err = suite.LogoutNode(from.String())
	suite.Require().Nil(err)
	err = suite.CheckNodeStatus(from.String(), governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

//tc：治理节点处于unavailable状态注册治理节点，节点注册成功
func (suite *Model12) Test1203_RegisterVpNodeWithUnavailableNodeIsSuccess() {
	pid, err := suite.MockPid()
	suite.Require().Nil(err)
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.NodeToUnavailable(from.String(), "vpNode", pid, 5, from.String(), "")
	suite.Require().Nil(err)
	err = suite.RegisterNode(from.String(), "vpNode", pid, 5, from.String(), "")
	suite.Require().Nil(err)
	//recover
	err = suite.LogoutNode(from.String())
	suite.Require().Nil(err)
	err = suite.CheckNodeStatus(from.String(), governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

//tc：治理节点用已存在的name注册节点，节点注册成功
func (suite *Model12) Test1204_RegisterVpNodeWithSameNameIsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterNode(from1.String(), "nvpNode", "", 0, from1.String(), from)
	suite.Require().Nil(err)
	pid, err := suite.MockPid()
	suite.Require().Nil(err)
	_, from2, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterNode(from2.String(), "vpNode", pid, 5, from1.String(), "")
	suite.Require().Nil(err)
	//recover
	err = suite.LogoutNode(from1.String())
	suite.Require().Nil(err)
	err = suite.CheckNodeStatus(from1.String(), governance.GovernanceForbidden)
	suite.Require().Nil(err)
	err = suite.LogoutNode(from2.String())
	suite.Require().Nil(err)
	err = suite.CheckNodeStatus(from2.String(), governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

//tc：中继链管理员，注销治理节点，节点注销成功
func (suite *Model12) Test1205_LogoutVpNodeIsSuccess() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	pid, err := suite.MockPid()
	suite.Require().Nil(err)
	err = suite.RegisterNode(from.String(), "vpNode", pid, 5, from.String(), "")
	suite.Require().Nil(err)
	err = suite.LogoutNode(from.String())
	suite.Require().Nil(err)
	err = suite.CheckNodeStatus(from.String(), governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

//tc：治理节点处于available状态注销治理节点，节点注销成功
func (suite *Model12) Test1206_LogoutVpNodeIsSuccess() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	pid, err := suite.MockPid()
	suite.Require().Nil(err)
	err = suite.RegisterNode(from.String(), "vpNode", pid, 5, from.String(), "")
	suite.Require().Nil(err)
	err = suite.LogoutNode(from.String())
	suite.Require().Nil(err)
	err = suite.CheckNodeStatus(from.String(), governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

//tc：中继链管理员，注册审计节点，节点注册成功
func (suite *Model12) Test1207_RegisterNvpNodeWithRelayNodeIsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterNode(from1.String(), "nvpNode", "", 0, from1.String(), from)
	suite.Require().Nil(err)
}

//tc：审计节点未注册，节点注册成功
func (suite *Model12) Test1208_RegisterNvpNodeIsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterNode(from1.String(), "nvpNode", "", 0, from1.String(), from)
	suite.Require().Nil(err)
}

//tc：审计节点处于unavailable状态注册审计节点，节点注册成功
func (suite *Model12) Test1209_RegisterNvpNodeWithUnavailableNodeIsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.NodeToUnavailable(from1.String(), "nvpNode", "", 0, from1.String(), from)
	suite.Require().Nil(err)
	err = suite.RegisterNode(from1.String(), "nvpNode", "", 0, from1.String(), from)
	suite.Require().Nil(err)
}

//tc：审计节点用已存在的pid注册节点，节点注册成功
func (suite *Model12) Test1210_RegisterNvpNodeWithSamePidIsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	pid, err := suite.MockPid()
	suite.Require().Nil(err)
	err = suite.RegisterNode(from1.String(), "vpNode", pid, 5, from1.String(), "")
	suite.Require().Nil(err)
	_, from2, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterNode(from2.String(), "nvpNode", pid, 0, from2.String(), from)
	suite.Require().Nil(err)
	//recover
	err = suite.LogoutNode(from1.String())
	suite.Require().Nil(err)
}

//tc：审计节点用已存在的vpid注册节点，节点注册成功
func (suite *Model12) Test1211_RegisterNvpNodeWithSameVPIDIsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	_, from2, err := repo.KeyPriv()
	suite.Require().Nil(err)
	pid1, err := suite.MockPid()
	suite.Require().Nil(err)
	err = suite.RegisterNode(from1.String(), "vpNode", pid1, 5, from1.String(), "")
	suite.Require().Nil(err)
	err = suite.RegisterNode(from2.String(), "nvpNode", "", 5, from2.String(), from)
	suite.Require().Nil(err)
	//recover
	err = suite.LogoutNode(from1.String())
	suite.Require().Nil(err)
}

//tc：中继链管理员，更新审计节点，节点更新成功
func (suite *Model12) Test1212_UpdateNvpNodeWithRelayAdminIsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterNode(from1.String(), "nvpNode", "", 0, from1.String(), from)
	suite.Require().Nil(err)
	err = suite.UpdateNode(from1.String(), from1.String(), from)
	suite.Require().Nil(err)
}

//tc：审计管理员，更新审计节点，节点更新成功
func (suite *Model12) Test1213_UpdateNvpNodeWithAuditAdminIsSuccess() {
	pk1, from1, address1, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk1, from1, address1)
	suite.Require().Nil(err)
	pk2, from2, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterNode(from2.String(), "nvpNode", "", 0, from2.String(), from1)
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(from2.String()),
		rpcx.String(from2.String() + "123"),
		rpcx.String(from1),
	}
	client := suite.NewClient(pk2)
	_, err = client.InvokeBVMContract(constant.NodeManagerContractAddr.Address(), UpdateNode, nil, args...)
	suite.Require().NotNil(err)
}

//tc：审计节点处于available状态更新审计节点，节点更新成功
func (suite *Model12) Test1214_UpdateNvpNodeWithAvailableNodeIsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterNode(from1.String(), "nvpNode", "", 0, from1.String(), from)
	suite.Require().Nil(err)
	err = suite.UpdateNode(from1.String(), from1.String(), from)
	suite.Require().Nil(err)
}

//tc：审计节点处于binded状态更新审计节点，节点更新成功
func (suite *Model12) Test1215_UpdateNvpNodeWithBindedNodeIsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.NodeToBinded(from1.String(), "nvpNode", "", 0, from1.String(), from)
	suite.Require().Nil(err)
	err = suite.UpdateNode(from1.String(), from1.String(), from)
	suite.Require().Nil(err)
}

//tc：中继链管理员，注销审计节点，节点注销成功
func (suite *Model12) Test1216_LogoutNvpNodeIsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterNode(from1.String(), "nvpNode", "", 0, from1.String(), from)
	suite.Require().Nil(err)
	err = suite.LogoutNode(from1.String())
	suite.Require().Nil(err)
}

//tc：审计节点处于available状态注销审计节点，节点注销成功
func (suite *Model12) Test1217_LogoutNvpNodeWithAvailableNodeIsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterNode(from1.String(), "nvpNode", "", 0, from1.String(), from)
	suite.Require().Nil(err)
	err = suite.LogoutNode(from1.String())
	suite.Require().Nil(err)
}

//tc：审计节点处于binding状态注销审计节点，节点注销成功
func (suite *Model12) Test1218_LogoutNvpNodeWithBindingNodeIsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	_, err = suite.NodeToBinding(from1.String(), "nvpNode", "", 0, from1.String(), from)
	suite.Require().Nil(err)
	err = suite.LogoutNode(from1.String())
	suite.Require().Nil(err)
}

//tc：审计节点处于binded状态注销审计节点，节点注销成功
func (suite *Model12) Test1219_LogoutNvpNodeWithBindedNodeIsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.NodeToBinded(from1.String(), "nvpNode", "", 0, from1.String(), from)
	suite.Require().Nil(err)
	err = suite.LogoutNode(from1.String())
	suite.Require().Nil(err)
}

//tc：审计节点处于updating状态注销审计节点，节点注销成功
func (suite *Model12) Test1220_LogoutNvpNodeWithUpdatingNodeIsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	_, err = suite.NodeToUpdating(from1.String(), "nvpNode", "", 0, from1.String(), from)
	suite.Require().Nil(err)
	err = suite.LogoutNode(from1.String())
	suite.Require().Nil(err)
}

// InvokeNodeContract invoke node contract by method and args
func (suite *Snake) InvokeNodeContract(method string, args ...*pb.Arg) (string, error) {
	pk, from, err := repo.Node1Priv()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.NodeManagerContractAddr.Address(), method, &rpcx.TransactOpts{
		From:  from.String(),
		Nonce: atomic.AddUint64(&nonce1, 1),
	}, args...)
	suite.Require().Nil(err)
	if res.Status != pb.Receipt_SUCCESS {
		return "", fmt.Errorf(string(res.Ret))
	}
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	suite.Require().Nil(err)
	return result.ProposalID, nil
}

// RegisterNode register node
func (suite *Snake) RegisterNode(nodeAccount, nodeType, nodePid string, nodeVpId uint64, nodeName, permit string) error {
	args := []*pb.Arg{
		rpcx.String(nodeAccount), //nodeAccount
		rpcx.String(nodeType),    //nodeType
		rpcx.String(nodePid),     //nodePid
		rpcx.Uint64(nodeVpId),    //nodeVpId
		rpcx.String(nodeName),    //nodeName
		rpcx.String(permit),      //permitStr
		rpcx.String("reason"),    //reason
	}
	proposal, err := suite.InvokeNodeContract(RegisterNode, args...)
	suite.Require().Nil(err)
	err = suite.VotePass(proposal)
	suite.Require().Nil(err)
	err = suite.CheckNodeStatus(nodeAccount, governance.GovernanceAvailable)
	suite.Require().Nil(err)
	return nil
}

//UpdateNode update node
func (suite *Snake) UpdateNode(nodeAccount, nodeName, permit string) error {
	args := []*pb.Arg{
		rpcx.String(nodeAccount),
		rpcx.String(nodeName),
		rpcx.String(permit),
		rpcx.String("reason"),
	}
	proposal, err := suite.InvokeNodeContract(UpdateNode, args...)
	suite.Require().Nil(err)
	err = suite.VotePass(proposal)
	suite.Require().Nil(err)
	return nil
}

// LogoutNode logout node
func (suite *Snake) LogoutNode(account string) error {
	args := []*pb.Arg{
		rpcx.String(account),
		rpcx.String("reason"),
	}
	proposal, err := suite.InvokeNodeContract(LogoutNode, args...)
	suite.Require().Nil(err)
	err = suite.VotePass(proposal)
	suite.Require().Nil(err)
	err = suite.CheckNodeStatus(account, governance.GovernanceForbidden)
	suite.Require().Nil(err)
	return nil
}

// NodeToRegisting get a registing node
func (suite *Snake) NodeToRegisting(nodeAccount, nodeType, nodePid string, nodeVpId uint64, nodeName, permit string) (string, error) {
	args := []*pb.Arg{
		rpcx.String(nodeAccount), //nodeAccount
		rpcx.String(nodeType),    //nodeType
		rpcx.String(nodePid),     //nodePid
		rpcx.Uint64(nodeVpId),    //nodeVpId
		rpcx.String(nodeName),    //nodeName
		rpcx.String(permit),      //permitStr
		rpcx.String("reason"),    //reason
	}
	proposal, err := suite.InvokeNodeContract(RegisterNode, args...)
	suite.Require().Nil(err)
	err = suite.CheckNodeStatus(nodeAccount, governance.GovernanceRegisting)
	suite.Require().Nil(err)
	return proposal, nil
}

// NodeToUnavailable get a unavailable node
func (suite *Snake) NodeToUnavailable(nodeAccount, nodeType, nodePid string, nodeVpId uint64, nodeName, permit string) error {
	args := []*pb.Arg{
		rpcx.String(nodeAccount), //nodeAccount
		rpcx.String(nodeType),    //nodeType
		rpcx.String(nodePid),     //nodePid
		rpcx.Uint64(nodeVpId),    //nodeVpId
		rpcx.String(nodeName),    //nodeName
		rpcx.String(permit),      //permitStr
		rpcx.String("reason"),    //reason
	}
	proposal, err := suite.InvokeNodeContract(RegisterNode, args...)
	suite.Require().Nil(err)
	err = suite.VoteReject(proposal)
	suite.Require().Nil(err)
	err = suite.CheckNodeStatus(nodeAccount, governance.GovernanceUnavailable)
	suite.Require().Nil(err)
	return nil
}

// NodeToLogouting get a logouting node
func (suite *Snake) NodeToLogouting(nodeAccount, nodeType, nodePid string, nodeVpId uint64, nodeName, permit string) (string, error) {
	err := suite.RegisterNode(nodeAccount, nodeType, nodePid, nodeVpId, nodeName, permit)
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(nodeAccount),
		rpcx.String("reason"),
	}
	proposal, err := suite.InvokeNodeContract(LogoutNode, args...)
	suite.Require().Nil(err)
	err = suite.CheckNodeStatus(nodeAccount, governance.GovernanceLogouting)
	suite.Require().Nil(err)
	return proposal, nil
}

// NodeToForbidden get a forbidden node
func (suite *Snake) NodeToForbidden(nodeAccount, nodeType, nodePid string, nodeVpId uint64, nodeName, permit string) error {
	err := suite.RegisterNode(nodeAccount, nodeType, nodePid, nodeVpId, nodeName, permit)
	suite.Require().Nil(err)
	err = suite.LogoutNode(nodeAccount)
	suite.Require().Nil(err)
	err = suite.CheckNodeStatus(nodeAccount, governance.GovernanceForbidden)
	suite.Require().Nil(err)
	return nil
}

// NodeToBinding get a binding nvp node
func (suite *Snake) NodeToBinding(nodeAccount, nodeType, nodePid string, nodeVpId uint64, nodeName, permit string) (string, error) {
	err := suite.RegisterNode(nodeAccount, nodeType, nodePid, nodeVpId, nodeName, permit)
	suite.Require().Nil(err)
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(from.String()),
		rpcx.String(AuditAdmin),
		rpcx.String(nodeAccount),
		rpcx.String("reason"),
	}
	proposal, err := suite.InvokeRoleContract(RegisterRole, args...)
	suite.Require().Nil(err)
	err = suite.CheckNodeStatus(nodeAccount, governance.GovernanceBinding)
	suite.Require().Nil(err)
	return proposal, nil
}

// NodeToBinded get a binded node
func (suite *Snake) NodeToBinded(nodeAccount, nodeType, nodePid string, nodeVpId uint64, nodeName, permit string) error {
	err := suite.RegisterNode(nodeAccount, nodeType, nodePid, nodeVpId, nodeName, permit)
	suite.Require().Nil(err)
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterRole(from.String(), AuditAdmin, nodeAccount)
	suite.Require().Nil(err)
	err = suite.CheckNodeStatus(nodeAccount, "binded")
	suite.Require().Nil(err)
	return nil
}

// NodeToUpdating get a updating node
func (suite *Snake) NodeToUpdating(nodeAccount, nodeType, nodePid string, nodeVpId uint64, nodeName, permit string) (string, error) {
	err := suite.RegisterNode(nodeAccount, nodeType, nodePid, nodeVpId, nodeName, permit)
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(nodeAccount),
		rpcx.String(nodeName + "123"),
		rpcx.String(permit),
		rpcx.String("reason"),
	}
	proposal, err := suite.InvokeNodeContract(UpdateNode, args...)
	suite.Require().Nil(err)
	err = suite.CheckNodeStatus(nodeAccount, governance.GovernanceUpdating)
	suite.Require().Nil(err)
	return proposal, nil
}

// CheckNodeStatus check node status
func (suite *Snake) CheckNodeStatus(account string, status governance.GovernanceStatus) error {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.NodeManagerContractAddr.Address(), "GetNode", nil, rpcx.String(account))
	suite.Require().Nil(err)
	if res.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf(string(res.Ret))
	}
	node := &Node{}
	err = json.Unmarshal(res.Ret, node)
	suite.Require().Nil(err)
	if node.Status != status {
		return fmt.Errorf("expect status is %s, but got %s", status, node.Status)
	}
	return nil
}
