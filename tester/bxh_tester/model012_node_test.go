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

//tc：非中继链管理员，注册治理节点，节点注册失败
func (suite *Model12) Test1201_RegisterVpNodeWithNoRelayAdminIsFail() {
	pk, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	pid, err := suite.MockPid()
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(from.String()), //nodeAccount
		rpcx.String("vpNode"),      //nodeType
		rpcx.String(pid),           //nodePid
		rpcx.Uint64(5),             // nodeVpId
		rpcx.String(from.String()), //nodeName
		rpcx.String(""),            //permitStr
		rpcx.String("reason"),
	}
	res, err := client.InvokeBVMContract(constant.NodeManagerContractAddr.Address(), "RegisterNode", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
}

//tc：中继链管理员，注册治理节点，节点注册成功
func (suite *Model12) Test1202_RegisterVpNodeWithRelayNodeIsSuccess() {
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
func (suite *Model12) Test1203_RegisterVpNodeWithNoRegisterIsSuccess() {
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

//tc：治理节点处于registing状态注册治理节点，节点注册失败
func (suite *Model12) Test1204_RegisterVpNodeWithRegistingNodeIsFail() {
	pid, err := suite.MockPid()
	suite.Require().Nil(err)
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	proposal, err := suite.NodeToRegisting(from.String(), "vpNode", pid, 5, from.String(), "")
	suite.Require().Nil(err)
	err = suite.RegisterNode(from.String(), "vpNode", pid, 5, from.String(), "")
	suite.Require().NotNil(err)
	//recover
	err = suite.VoteReject(proposal)
	suite.Require().Nil(err)
	err = suite.CheckNodeStatus(from.String(), governance.GovernanceUnavailable)
	suite.Require().Nil(err)
}

//tc：治理节点处于unavailable状态注册治理节点，节点注册成功
func (suite *Model12) Test1205_RegisterVpNodeWithUnavailableNodeIsSuccess() {
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

//tc：治理节点处于available状态注册治理节点，节点注册失败
func (suite *Model12) Test1206_RegisterVpNodeWithAvailableNodeIsFail() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	pid, err := suite.MockPid()
	suite.Require().Nil(err)
	err = suite.RegisterNode(from.String(), "vpNode", pid, 5, from.String(), "")
	suite.Require().Nil(err)
	err = suite.RegisterNode(from.String(), "vpNode", pid, 5, from.String(), "")
	suite.Require().NotNil(err)
	//recover
	err = suite.LogoutNode(from.String())
	suite.Require().Nil(err)
	err = suite.CheckNodeStatus(from.String(), governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

//tc：治理节点处于logouting状态注册治理节点，节点注册失败
func (suite *Model12) Test1207_RegisterVpNodeWithLogoutingNodeIsFail() {
	pid, err := suite.MockPid()
	suite.Require().Nil(err)
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	proposal, err := suite.NodeToLogouting(from.String(), "vpNode", pid, 5, from.String(), "")
	suite.Require().Nil(err)
	err = suite.RegisterNode(from.String(), "vpNode", pid, 5, from.String(), "")
	suite.Require().NotNil(err)
	//recover
	err = suite.VotePass(proposal)
	suite.Require().Nil(err)
	err = suite.CheckNodeStatus(from.String(), governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

//tc：治理节点处于forbidden状态注册治理节点，节点注册失败
func (suite *Model12) Test1208_RegisterVpNodeWithForbiddenNodeIsFail() {
	pid, err := suite.MockPid()
	suite.Require().Nil(err)
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.NodeToForbidden(from.String(), "vpNode", pid, 5, from.String(), "")
	suite.Require().Nil(err)
	err = suite.RegisterNode(from.String(), "vpNode", pid, 5, from.String(), "")
	suite.Require().NotNil(err)
	err = suite.CheckNodeStatus(from.String(), governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

//tc：治理节点用已存在的name注册节点，节点注册成功
func (suite *Model12) Test1209_RegisterVpNodeWithSameNameIsSuccess() {
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

//tc：治理节点用已存在的pid注册节点，节点注册失败
func (suite *Model12) Test1210_RegisterVpNodeWithSamePidIsFail() {
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	pid, err := suite.MockPid()
	suite.Require().Nil(err)
	err = suite.RegisterNode(from1.String(), "vpNode", pid, 5, from1.String(), "")
	suite.Require().Nil(err)
	_, from2, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterNode(from2.String(), "vpNode", pid, 5, from2.String(), "")
	suite.Require().NotNil(err)
	//recover
	err = suite.LogoutNode(from1.String())
	suite.Require().Nil(err)
	err = suite.CheckNodeStatus(from1.String(), governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

//tc：治理节点用已存在的vpid注册节点，节点注册失败
func (suite *Model12) Test1211_RegisterVpNodeWithSameVPIDIsFail() {
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	_, from2, err := repo.KeyPriv()
	suite.Require().Nil(err)
	pid1, err := suite.MockPid()
	suite.Require().Nil(err)
	err = suite.RegisterNode(from1.String(), "vpNode", pid1, 5, from1.String(), "")
	suite.Require().Nil(err)
	pid2, err := suite.MockPid()
	suite.Require().Nil(err)
	err = suite.RegisterNode(from2.String(), "vpNode", pid2, 5, from2.String(), "")
	suite.Require().NotNil(err)
	//recover
	err = suite.LogoutNode(from1.String())
	suite.Require().Nil(err)
	err = suite.CheckNodeStatus(from1.String(), governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

//tc：治理节点用较小vpid注册节点，节点注册失败
func (suite *Model12) Test1212_RegisterVpNodeWithSmallVpIDIsFail() {
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	pid1, err := suite.MockPid()
	suite.Require().Nil(err)
	err = suite.RegisterNode(from1.String(), "vpNode", pid1, 1, from1.String(), "")
	suite.Require().NotNil(err)
}

//tc：治理节点用已存在的account注册节点，节点注册失败
func (suite *Model12) Test1213_RegisterVpNodeWithSameAccountIsFail() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	pid1, err := suite.MockPid()
	suite.Require().Nil(err)
	err = suite.RegisterNode(from.String(), "vpNode", pid1, 5, from.String(), "")
	suite.Require().Nil(err)
	pid2, err := suite.MockPid()
	suite.Require().Nil(err)
	err = suite.RegisterNode(from.String(), "vpNode", pid2, 5, from.String(), "")
	suite.Require().NotNil(err)
	//recover
	err = suite.LogoutNode(from.String())
	suite.Require().Nil(err)
	err = suite.CheckNodeStatus(from.String(), governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

//tc：其他治理节点正在注册，注册节点，节点注册失败
func (suite *Model12) Test1214_RegisterVpNodeWithOtherRegistingNodeIsFail() {
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	_, from2, err := repo.KeyPriv()
	suite.Require().Nil(err)
	pid1, err := suite.MockPid()
	suite.Require().Nil(err)
	proposal, err := suite.NodeToRegisting(from1.String(), "vpNode", pid1, 5, from1.String(), "")
	suite.Require().Nil(err)
	pid2, err := suite.MockPid()
	suite.Require().Nil(err)
	err = suite.RegisterNode(from2.String(), "vpNode", pid2, 5, from2.String(), "")
	suite.Require().NotNil(err)
	//recover
	err = suite.VoteReject(proposal)
	suite.Require().Nil(err)
	err = suite.CheckNodeStatus(from1.String(), governance.GovernanceUnavailable)
	suite.Require().Nil(err)
}

//tc：其他治理节点正在注销，注册节点，节点注册失败
func (suite *Model12) Test1215_RegisterVpNodeWIthOtherLogoutingNodeIsFail() {
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	_, from2, err := repo.KeyPriv()
	suite.Require().Nil(err)
	pid1, err := suite.MockPid()
	suite.Require().Nil(err)
	proposal, err := suite.NodeToLogouting(from1.String(), "vpNode", pid1, 5, from1.String(), "")
	suite.Require().Nil(err)
	pid2, err := suite.MockPid()
	suite.Require().Nil(err)
	err = suite.RegisterNode(from2.String(), "vpNode", pid2, 5, from2.String(), "")
	suite.Require().NotNil(err)
	//recover
	err = suite.VotePass(proposal)
	suite.Require().Nil(err)
	err = suite.CheckNodeStatus(from1.String(), governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

//tc：非中继链管理员，注销治理节点，节点注销失败
func (suite *Model12) Test1216_LogoutVpNodeWithNoRelayAdminIsFail() {
	pk, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	pid, err := suite.MockPid()
	suite.Require().Nil(err)
	err = suite.RegisterNode(from.String(), "vpNode", pid, 5, from.String(), "")
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
	err = suite.CheckNodeStatus(from.String(), governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

//tc：中继链管理员，注销治理节点，节点注销成功
func (suite *Model12) Test1217_LogoutVpNodeIsSuccess() {
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

//tc：治理节点未注册，节点注销失败
func (suite *Model12) Test1218_LogoutVpNodeWithNoRegisterIsFail() {
	err := suite.LogoutNode("0x79a1215469FaB6f9c63c1816b45183AD3624bE31")
	suite.Require().NotNil(err)
}

//tc：治理节点处于registing状态注销治理节点，节点注销失败
func (suite *Model12) Test1219_LogoutVpNodeWithRegistingNodeIsFail() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	pid, err := suite.MockPid()
	suite.Require().Nil(err)
	proposal, err := suite.NodeToRegisting(from.String(), "vpNode", pid, 5, from.String(), "")
	suite.Require().Nil(err)
	err = suite.LogoutNode(from.String())
	suite.Require().NotNil(err)
	//recover
	err = suite.VoteReject(proposal)
	suite.Require().Nil(err)
	err = suite.CheckNodeStatus(from.String(), governance.GovernanceUnavailable)
	suite.Require().Nil(err)
}

//tc：治理节点处于unavailable状态注销治理节点，节点注销失败
func (suite *Model12) Test1220_LogoutVpNodeWithUnavailableNodeIsFail() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	pid, err := suite.MockPid()
	suite.Require().Nil(err)
	err = suite.NodeToUnavailable(from.String(), "vpNode", pid, 5, from.String(), "")
	suite.Require().Nil(err)
	err = suite.LogoutNode(from.String())
	suite.Require().NotNil(err)
}

//tc：治理节点处于available状态注销治理节点，节点注销成功
func (suite *Model12) Test1221_LogoutVpNodeIsSuccess() {
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

//tc：治理节点处于logouting状态注销治理节点，节点注销失败
func (suite *Model12) Test1222_LogoutVpNodeWithLogoutingNodeIsFail() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	pid, err := suite.MockPid()
	suite.Require().Nil(err)
	proposal, err := suite.NodeToLogouting(from.String(), "vpNode", pid, 5, from.String(), "")
	suite.Require().Nil(err)
	err = suite.LogoutNode(from.String())
	suite.Require().NotNil(err)
	//recover
	err = suite.VotePass(proposal)
	suite.Require().Nil(err)
	err = suite.CheckNodeStatus(from.String(), governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

//tc：治理节点处于forbidden状态注销治理节点，节点注销失败
func (suite *Model12) Test1223_LogoutVpNodeWithForbiddenNodeIsFail() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	pid, err := suite.MockPid()
	suite.Require().Nil(err)
	err = suite.NodeToForbidden(from.String(), "vpNode", pid, 5, from.String(), "")
	suite.Require().Nil(err)
	err = suite.LogoutNode(from.String())
	suite.Require().NotNil(err)
}

//tc：中继链管理员，注销创世节点，节点注销失败
func (suite *Model12) Test1224_LogoutVpNodeWithFourNodeIsFail() {
	err := suite.LogoutNode("0x79a1215469FaB6f9c63c1816b45183AD3624bE34")
	suite.Require().NotNil(err)
}

//tc：其他治理节点正在注册，注销节点，节点注销失败
func (suite *Model12) Test1225_LogoutVpNodeWithOtherRegistingNodeIsFail() {
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	_, from2, err := repo.KeyPriv()
	suite.Require().Nil(err)
	pid1, err := suite.MockPid()
	suite.Require().Nil(err)
	err = suite.RegisterNode(from1.String(), "vpNode", pid1, 5, from1.String(), "")
	suite.Require().Nil(err)
	pid2, err := suite.MockPid()
	suite.Require().Nil(err)
	proposal, err := suite.NodeToRegisting(from2.String(), "vpNode", pid2, 6, from2.String(), "")
	suite.Require().Nil(err)
	err = suite.LogoutNode(from1.String())
	suite.Require().NotNil(err)
	//recover
	err = suite.VoteReject(proposal)
	suite.Require().Nil(err)
	err = suite.CheckNodeStatus(from2.String(), governance.GovernanceUnavailable)
	suite.Require().Nil(err)
	err = suite.LogoutNode(from1.String())
	suite.Require().Nil(err)
	err = suite.CheckNodeStatus(from1.String(), governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

//tc：其他治理节点正在注销，注销节点，节点注销失败
func (suite *Model12) Test1226_LogoutVpNodeWithOtherLogoutingNodeIsFail() {
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	_, from2, err := repo.KeyPriv()
	suite.Require().Nil(err)
	pid1, err := suite.MockPid()
	suite.Require().Nil(err)
	err = suite.RegisterNode(from1.String(), "vpNode", pid1, 5, from1.String(), "")
	suite.Require().Nil(err)
	pid2, err := suite.MockPid()
	suite.Require().Nil(err)
	proposal, err := suite.NodeToLogouting(from2.String(), "vpNode", pid2, 6, from2.String(), "")
	suite.Require().Nil(err)
	err = suite.LogoutNode(from1.String())
	suite.Require().NotNil(err)
	//recover
	err = suite.VotePass(proposal)
	suite.Require().Nil(err)
	err = suite.CheckNodeStatus(from2.String(), governance.GovernanceForbidden)
	suite.Require().Nil(err)
	err = suite.LogoutNode(from1.String())
	suite.Require().Nil(err)
	err = suite.CheckNodeStatus(from1.String(), governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

//tc：非中继链管理员，注册审计节点，节点注册失败
func (suite *Model12) Test1227_RegisterNvpNodeWithNoRelayAdminIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	pk1, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	client := suite.NewClient(pk1)
	args := []*pb.Arg{
		rpcx.String(from1.String()), //nodeAccount
		rpcx.String("nvpNode"),      //nodeType
		rpcx.String(""),             //nodePid
		rpcx.Uint64(0),              // nodeVpId
		rpcx.String(from1.String()), //nodeName
		rpcx.String(from),           //permitStr
		rpcx.String("reason"),
	}
	res, err := client.InvokeBVMContract(constant.NodeManagerContractAddr.Address(), "RegisterNode", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
}

//tc：中继链管理员，注册审计节点，节点注册成功
func (suite *Model12) Test1228_RegisterNvpNodeWithRelayNodeIsSuccess() {
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
func (suite *Model12) Test1229_RegisterNvpNodeIsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterNode(from1.String(), "nvpNode", "", 0, from1.String(), from)
	suite.Require().Nil(err)
}

//tc：审计节点处于registing状态注册审计节点，节点注册失败
func (suite *Model12) Test1230_RegisterNvpNodeWithRegistingNodeIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	_, err = suite.NodeToRegisting(from1.String(), "nvpNode", "", 0, from1.String(), from)
	suite.Require().Nil(err)
	err = suite.RegisterNode(from1.String(), "nvpNode", "", 0, from1.String(), from)
	suite.Require().NotNil(err)
}

//tc：审计节点处于unavailable状态注册审计节点，节点注册成功
func (suite *Model12) Test1231_RegisterNvpNodeWithUnavailableNodeIsSuccess() {
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

//tc：审计节点处于available状态注册审计节点，节点注册失败
func (suite *Model12) Test1232_RegisterNvpNodeWithAvailableNodeIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterNode(from1.String(), "nvpNode", "", 0, from1.String(), from)
	suite.Require().Nil(err)
	err = suite.RegisterNode(from1.String(), "nvpNode", "", 0, from1.String(), from)
	suite.Require().NotNil(err)
}

//tc：审计节点处于binding状态注册审计节点，节点注册失败
func (suite *Model12) Test1233_RegisterNvpNodeWithBindingNodeIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	_, err = suite.NodeToBinding(from1.String(), "nvpNode", "", 0, from1.String(), from)
	suite.Require().Nil(err)
	err = suite.RegisterNode(from1.String(), "nvpNode", "", 0, from1.String(), from)
	suite.Require().NotNil(err)
}

//tc：审计节点处于binded状态注册审计节点，节点注册失败
func (suite *Model12) Test1234_RegisterNvpNodeWithBindedNodeIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.NodeToBinded(from1.String(), "nvpNode", "", 0, from1.String(), from)
	suite.Require().Nil(err)
	err = suite.RegisterNode(from1.String(), "nvpNode", "", 0, from1.String(), from)
	suite.Require().NotNil(err)
}

//tc：审计节点处于updating状态注册审计节点，节点注册失败
func (suite *Model12) Test1235_RegisterNvpNodeWithUpdatingNodeIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	_, err = suite.NodeToUpdating(from1.String(), "nvpNode", "", 0, from1.String(), from)
	suite.Require().Nil(err)
	err = suite.RegisterNode(from1.String(), "nvpNode", "", 0, from1.String(), from)
	suite.Require().NotNil(err)
}

//tc：审计节点处于logouting状态注册审计节点，节点注册失败
func (suite *Model12) Test1236_RegisterNvpNodeWithLogoutingNodeIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	_, err = suite.NodeToLogouting(from1.String(), "nvpNode", "", 0, from1.String(), from)
	suite.Require().Nil(err)
	err = suite.RegisterNode(from1.String(), "nvpNode", "", 0, from1.String(), from)
	suite.Require().NotNil(err)
}

//tc：审计节点处于forbidden状态注册审计节点，节点注册失败
func (suite *Model12) Test1237_RegisterNvpNodeWithForbiddenNodeIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.NodeToForbidden(from1.String(), "nvpNode", "", 0, from1.String(), from)
	suite.Require().Nil(err)
	err = suite.RegisterNode(from1.String(), "nvpNode", "", 0, from1.String(), from)
	suite.Require().NotNil(err)
}

//tc：审计节点用已存在的name注册节点，节点注册失败
func (suite *Model12) Test1238_RegisterNvpNodeWithSameNameIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterNode(from1.String(), "nvpNode", "", 0, from1.String(), from)
	suite.Require().Nil(err)
	_, from2, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterNode(from2.String(), "nvpNode", "", 0, from1.String(), from)
	suite.Require().NotNil(err)
}

//tc：审计节点用已存在的pid注册节点，节点注册成功
func (suite *Model12) Test1239_RegisterNvpNodeWithSamePidIsSuccess() {
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
func (suite *Model12) Test1240_RegisterNvpNodeWithSameVPIDIsSuccess() {
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

//tc：审计节点用已存在的account注册节点，节点注册失败
func (suite *Model12) Test1241_RegisterNvpNodeWithSameAccountIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterNode(from1.String(), "nvpNode", "", 0, from1.String(), from)
	suite.Require().Nil(err)
	err = suite.RegisterNode(from1.String(), "nvpNode", "", 0, from1.String()+"1", from)
	suite.Require().NotNil(err)
}

//tc：审计节点权限为空注册节点，节点注册失败
func (suite *Model12) Test1242_RegisterNvpNodeWithEmptyPermitIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterNode(from1.String(), "nvpNode", "", 0, from1.String(), "")
	suite.Require().NotNil(err)
}

//tc：非中继链管理员，更新审计节点，节点更新失败
func (suite *Model12) Test1243_UpdateNvpNodeWithNoRelayAdminIsFail() {
	pk1, from1, address1, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk1, from1, address1)
	suite.Require().Nil(err)
	_, from2, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterNode(from2.String(), "nvpNode", "", 0, from2.String(), from1)
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(from2.String()),
		rpcx.String(from2.String() + "123"),
		rpcx.String(from1),
	}
	pk3, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	client := suite.NewClient(pk3)
	res, err := client.InvokeBVMContract(constant.NodeManagerContractAddr.Address(), UpdateNode, nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
}

//tc：中继链管理员，更新审计节点，节点更新成功
func (suite *Model12) Test1244_UpdateNvpNodeWithRelayAdminIsSuccess() {
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
func (suite *Model12) Test1245_UpdateNvpNodeWithAuditAdminIsSuccess() {
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

//tc：审计节点未注册，节点更新失败
func (suite *Model12) Test1246_UpdateNvpNodeWithNoRegisterNodeIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.UpdateNode(from1.String(), from1.String(), from)
	suite.Require().NotNil(err)
}

//tc：审计节点处于registing状态更新审计节点，节点更新失败
func (suite *Model12) Test1247_UpdateNvpNodeWithRegistingNodeIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	_, err = suite.NodeToRegisting(from1.String(), "nvpNode", "", 0, from1.String(), from)
	suite.Require().Nil(err)
	err = suite.UpdateNode(from1.String(), from1.String(), from)
	suite.Require().NotNil(err)
}

//tc：审计节点处于unavailable状态更新审计节点，节点更新失败
func (suite *Model12) Test1248_UpdateNvpNodeWithUnavailableNodeIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.NodeToUnavailable(from1.String(), "nvpNode", "", 0, from1.String(), from)
	suite.Require().Nil(err)
	err = suite.UpdateNode(from1.String(), from1.String(), from)
	suite.Require().NotNil(err)
}

//tc：审计节点处于available状态更新审计节点，节点更新成功
func (suite *Model12) Test1249_UpdateNvpNodeWithAvailableNodeIsSuccess() {
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

//tc：审计节点处于binding状态更新审计节点，节点更新失败
func (suite *Model12) Test1250_UpdateNvpNodeWithBindingNodeIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	_, err = suite.NodeToBinding(from1.String(), "nvpNode", "", 0, from1.String(), from)
	suite.Require().Nil(err)
	err = suite.UpdateNode(from1.String(), from1.String(), from)
	suite.Require().NotNil(err)
}

//tc：审计节点处于binded状态更新审计节点，节点更新成功
func (suite *Model12) Test1251_UpdateNvpNodeWithBindedNodeIsSuccess() {
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

//tc：审计节点处于updating状态更新审计节点，节点更新失败
func (suite *Model12) Test1252_UpdateNvpNodeWithUpdatingNodeIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	_, err = suite.NodeToUpdating(from1.String(), "nvpNode", "", 0, from1.String(), from)
	suite.Require().Nil(err)
	err = suite.UpdateNode(from1.String(), from1.String(), from)
	suite.Require().NotNil(err)
}

//tc：审计节点处于logouting状态更新审计节点，节点更新失败
func (suite *Model12) Test1253_UpdateNvpNodeWithLogoutingNodeIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	_, err = suite.NodeToLogouting(from1.String(), "nvpNode", "", 0, from1.String(), from)
	suite.Require().Nil(err)
	err = suite.UpdateNode(from1.String(), from1.String(), from)
	suite.Require().NotNil(err)
}

//tc：审计节点处于forbidden状态更新审计节点，节点更新失败
func (suite *Model12) Test1254_UpdateNvpNodeWithForbiddenNodeIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.NodeToForbidden(from1.String(), "nvpNode", "", 0, from1.String(), from)
	suite.Require().Nil(err)
	err = suite.UpdateNode(from1.String(), from1.String(), from)
	suite.Require().NotNil(err)
}

//tc：审计点点用已存在的name更新节点，节点更新失败
func (suite *Model12) Test1255_UpdateNvpNodeWithSameNameIsFail() {
	pk1, from1, address1, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk1, from1, address1)
	suite.Require().Nil(err)
	pk2, from2, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk2, from2, address2)
	suite.Require().Nil(err)
	_, from3, err := repo.KeyPriv()
	suite.Require().Nil(err)
	_, from4, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterNode(from3.String(), "nvpNode", "", 0, from3.String(), from1)
	suite.Require().Nil(err)
	err = suite.RegisterNode(from4.String(), "nvpNode", "", 0, from4.String(), from2)
	suite.Require().Nil(err)
	err = suite.UpdateNode(from4.String(), from3.String(), from2)
	suite.Require().NotNil(err)
}

//tc：审计点点用空的name更新节点，节点更新失败
func (suite *Model12) Test1256_UpdateNvpNodeWithEmptyNameIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterNode(from1.String(), "nvpNode", "", 0, from1.String(), from)
	suite.Require().Nil(err)
	err = suite.UpdateNode(from1.String(), "", from)
	suite.Require().NotNil(err)
}

//tc：审计节点权限为空更新节点，节点更新失败
func (suite *Model12) Test1257_UpdateNvpNodeWithEmptyPermitIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterNode(from1.String(), "nvpNode", "", 0, from1.String(), from)
	suite.Require().Nil(err)
	err = suite.UpdateNode(from1.String(), from1.String(), "")
	suite.Require().NotNil(err)
}

//tc：非中继链管理员，注销审计节点，节点注销失败
func (suite *Model12) Test1258_LogoutNvpNodeWithNoRelayAdminIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	err = suite.RegisterNode(from1.String(), "nvpNode", "", 0, from1.String(), from)
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(from1.String()),
		rpcx.String("reason"),
	}
	res, err := client.InvokeBVMContract(constant.NodeManagerContractAddr.Address(), "LogoutNode", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	//recover
	err = suite.LogoutNode(from1.String())
	suite.Require().Nil(err)
}

//tc：中继链管理员，注销审计节点，节点注销成功
func (suite *Model12) Test1259_LogoutNvpNodeIsSuccess() {
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

//tc：审计节点未注册，节点注销失败
func (suite *Model12) Test1260_LogoutNvpNodeWithNoRegisterNodeIsFail() {
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.LogoutNode(from1.String())
	suite.Require().NotNil(err)
}

//tc：审计节点处于registing状态注销审计节点，节点注销失败
func (suite *Model12) Test1261_LogoutNvpNodeWithRegistingNodeIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	_, err = suite.NodeToRegisting(from1.String(), "nvpNode", "", 0, from1.String(), from)
	suite.Require().Nil(err)
	err = suite.LogoutNode(from1.String())
	suite.Require().NotNil(err)
}

//tc：审计节点处于unavailable状态注销审计节点，节点注销失败
func (suite *Model12) Test1262_LogoutNvpNodeWithUnavailableNodeIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.NodeToUnavailable(from1.String(), "nvpNode", "", 0, from1.String(), from)
	suite.Require().Nil(err)
	err = suite.LogoutNode(from1.String())
	suite.Require().NotNil(err)
}

//tc：审计节点处于available状态注销审计节点，节点注销成功
func (suite *Model12) Test1263_LogoutNvpNodeWithAvailableNodeIsSuccess() {
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
func (suite *Model12) Test1264_LogoutNvpNodeWithBindingNodeIsSuccess() {
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
func (suite *Model12) Test1265_LogoutNvpNodeWithBindedNodeIsSuccess() {
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
func (suite *Model12) Test1266_LogoutNvpNodeWithUpdatingNodeIsSuccess() {
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

//tc：审计节点处于logouting状态注销审计节点，节点注销失败
func (suite *Model12) Test1267_LogoutNvpNodeWithLogoutingNodeIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	_, err = suite.NodeToLogouting(from1.String(), "nvpNode", "", 0, from1.String(), from)
	suite.Require().Nil(err)
	err = suite.LogoutNode(from1.String())
	suite.Require().NotNil(err)
}

//tc：审计节点处于forbidden状态注销审计节点，节点注销失败
func (suite *Model12) Test1268_LogoutNvpNodeWithForbiddenNodeIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.NodeToForbidden(from1.String(), "nvpNode", "", 0, from1.String(), from)
	suite.Require().Nil(err)
	err = suite.LogoutNode(from1.String())
	suite.Require().NotNil(err)
}

// InvokeNodeContract invoke node contract by method and args
func (suite *Snake) InvokeNodeContract(method string, args ...*pb.Arg) (string, error) {
	pk, from, err := repo.Node1Priv()
	if err != nil {
		return "", err
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.NodeManagerContractAddr.Address(), method, &rpcx.TransactOpts{
		From:  from.String(),
		Nonce: atomic.AddUint64(&nonce1, 1),
	}, args...)
	if err != nil {
		return "", err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return "", fmt.Errorf(string(res.Ret))
	}
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	if err != nil {
		return "", err
	}
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
	if err != nil {
		return err
	}
	err = suite.VotePass(proposal)
	if err != nil {
		return err
	}
	err = suite.CheckNodeStatus(nodeAccount, governance.GovernanceAvailable)
	if err != nil {
		return err
	}
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
	if err != nil {
		return err
	}
	err = suite.VotePass(proposal)
	if err != nil {
		return err
	}
	return nil
}

// LogoutNode logout node
func (suite *Snake) LogoutNode(account string) error {
	args := []*pb.Arg{
		rpcx.String(account),
		rpcx.String("reason"),
	}
	proposal, err := suite.InvokeNodeContract(LogoutNode, args...)
	if err != nil {
		return err
	}
	err = suite.VotePass(proposal)
	if err != nil {
		return err
	}
	err = suite.CheckNodeStatus(account, governance.GovernanceForbidden)
	if err != nil {
		return err
	}
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
	if err != nil {
		return "", err
	}
	err = suite.CheckNodeStatus(nodeAccount, governance.GovernanceRegisting)
	if err != nil {
		return "", err
	}
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
	if err != nil {
		return err
	}
	err = suite.VoteReject(proposal)
	if err != nil {
		return err
	}
	err = suite.CheckNodeStatus(nodeAccount, governance.GovernanceUnavailable)
	if err != nil {
		return err
	}
	return nil
}

// NodeToLogouting get a logouting node
func (suite *Snake) NodeToLogouting(nodeAccount, nodeType, nodePid string, nodeVpId uint64, nodeName, permit string) (string, error) {
	err := suite.RegisterNode(nodeAccount, nodeType, nodePid, nodeVpId, nodeName, permit)
	if err != nil {
		return "", err
	}
	args := []*pb.Arg{
		rpcx.String(nodeAccount),
		rpcx.String("reason"),
	}
	proposal, err := suite.InvokeNodeContract(LogoutNode, args...)
	if err != nil {
		return "", err
	}
	err = suite.CheckNodeStatus(nodeAccount, governance.GovernanceLogouting)
	if err != nil {
		return "", err
	}
	return proposal, nil
}

// NodeToForbidden get a forbidden node
func (suite *Snake) NodeToForbidden(nodeAccount, nodeType, nodePid string, nodeVpId uint64, nodeName, permit string) error {
	err := suite.RegisterNode(nodeAccount, nodeType, nodePid, nodeVpId, nodeName, permit)
	if err != nil {
		return err
	}
	err = suite.LogoutNode(nodeAccount)
	if err != nil {
		return err
	}
	err = suite.CheckNodeStatus(nodeAccount, governance.GovernanceForbidden)
	if err != nil {
		return err
	}
	return nil
}

// NodeToBinding get a binding nvp node
func (suite *Snake) NodeToBinding(nodeAccount, nodeType, nodePid string, nodeVpId uint64, nodeName, permit string) (string, error) {
	err := suite.RegisterNode(nodeAccount, nodeType, nodePid, nodeVpId, nodeName, permit)
	if err != nil {
		return "", err
	}
	_, from, err := repo.KeyPriv()
	if err != nil {
		return "", err
	}
	args := []*pb.Arg{
		rpcx.String(from.String()),
		rpcx.String(AuditAdmin),
		rpcx.String(nodeAccount),
		rpcx.String("reason"),
	}
	proposal, err := suite.InvokeRoleContract(RegisterRole, args...)
	if err != nil {
		return "", err
	}
	err = suite.CheckNodeStatus(nodeAccount, governance.GovernanceBinding)
	if err != nil {
		return "", err
	}
	return proposal, nil
}

// NodeToBinded get a binded node
func (suite *Snake) NodeToBinded(nodeAccount, nodeType, nodePid string, nodeVpId uint64, nodeName, permit string) error {
	err := suite.RegisterNode(nodeAccount, nodeType, nodePid, nodeVpId, nodeName, permit)
	if err != nil {
		return err
	}
	_, from, err := repo.KeyPriv()
	if err != nil {
		return err
	}
	err = suite.RegisterRole(from.String(), AuditAdmin, nodeAccount)
	if err != nil {
		return err
	}
	err = suite.CheckNodeStatus(nodeAccount, "binded")
	if err != nil {
		return err
	}
	return nil
}

// NodeToUpdating get a updating node
func (suite *Snake) NodeToUpdating(nodeAccount, nodeType, nodePid string, nodeVpId uint64, nodeName, permit string) (string, error) {
	err := suite.RegisterNode(nodeAccount, nodeType, nodePid, nodeVpId, nodeName, permit)
	if err != nil {
		return "", err
	}
	args := []*pb.Arg{
		rpcx.String(nodeAccount),
		rpcx.String(nodeName + "123"),
		rpcx.String(permit),
		rpcx.String("reason"),
	}
	proposal, err := suite.InvokeNodeContract(UpdateNode, args...)
	if err != nil {
		return "", err
	}
	err = suite.CheckNodeStatus(nodeAccount, governance.GovernanceUpdating)
	if err != nil {
		return "", err
	}
	return proposal, nil
}

// CheckNodeStatus check node status
func (suite *Snake) CheckNodeStatus(account string, status governance.GovernanceStatus) error {
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
		return fmt.Errorf(string(res.Ret))
	}
	node := &Node{}
	err = json.Unmarshal(res.Ret, node)
	if err != nil {
		return err
	}
	if node.Status != status {
		return fmt.Errorf("expect status is %s, but got %s", status, node.Status)
	}
	return nil
}
