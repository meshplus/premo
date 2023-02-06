package bxh_tester

import (
	"encoding/json"
	"fmt"
	"sync/atomic"

	"github.com/meshplus/bitxhub-core/governance"

	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/meshplus/premo/internal/repo"
)

const (
	RegisterRole    = "RegisterRole"
	FreezeRole      = "FreezeRole"
	ActivateRole    = "ActivateRole"
	LogoutRole      = "LogoutRole"
	BindRole        = "BindRole"
	GovernanceAdmin = "governanceAdmin"
	AuditAdmin      = "auditAdmin"
)

type Model11 struct {
	*Snake
}

//tc：注册治理管理员，管理员未注册，注册成功
func (suite *Model11) Test1101_RegisterAdminIsSuccess() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterRole(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	//recover
	err = suite.LogoutRole(from.String())
	suite.Require().Nil(err)
	err = suite.CheckRoleStatus(from.String(), governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

//tc：注册治理管理员，管理员处于unavailable，注册成功
func (suite *Model11) Test1102_RegisterAdminWithUnavailableAdminIsSuccess() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RoleToUnavailable(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	err = suite.RegisterRole(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	//recover
	err = suite.LogoutRole(from.String())
	suite.Require().Nil(err)
	err = suite.CheckRoleStatus(from.String(), governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

//tc：冻结治理管理员，管理员处于available，冻结成功
func (suite *Model11) Test1103_FreezeAdminWithAvailableAdminIsSuccess() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterRole(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	err = suite.FreezeRole(from.String())
	suite.Require().Nil(err)
	//recover
	err = suite.LogoutRole(from.String())
	suite.Require().Nil(err)
	err = suite.CheckRoleStatus(from.String(), governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

//tc：激活治理管理员，管理员处于frozen，激活成功
func (suite *Model11) Test1104_ActivateAdminWithFrozenAdminIsSuccess() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RoleToFrozen(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	err = suite.ActivateRole(from.String())
	suite.Require().Nil(err)
	//recover
	err = suite.LogoutRole(from.String())
	suite.Require().Nil(err)
	err = suite.CheckRoleStatus(from.String(), governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

//tc：注销治理管理员，管理员处于available，注销成功
func (suite *Model11) Test1105_LogoutAdminWithAvailableAdminIsSuccess() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterRole(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	err = suite.LogoutRole(from.String())
	suite.Require().Nil(err)
}

//tc：注销治理管理员，管理员处于freezing，注销成功
func (suite *Model11) Test1106_LogoutAdminWithFreezingAdminIsSuccess() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	_, err = suite.RoleToFreezing(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	err = suite.LogoutRole(from.String())
	suite.Require().Nil(err)
}

//tc：注销治理管理员，管理员处于frozen，注销成功
func (suite *Model11) Test1107_LogoutAdminWithFrozenAdminIsSuccess() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RoleToFrozen(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	err = suite.LogoutRole(from.String())
	suite.Require().Nil(err)
}

//tc：注销治理管理员，管理员处于activating，注销成功
func (suite *Model11) Test1108_LogoutAdminWithActivatingAdminIsSuccess() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	_, err = suite.RoleToActivating(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	err = suite.LogoutRole(from.String())
	suite.Require().Nil(err)
}

//tc：治理管理员处于available，管理员参与提案，参与成功
func (suite *Model11) Test1109_VoteWithAvailableAdminIsSuccess() {
	pk, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterRole(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	proposal := suite.GetProposal()
	client := suite.NewClient(pk)
	nonce, err := client.GetPendingNonceByAccount(from.String())
	suite.Require().Nil(err)
	_, err = suite.vote(pk, nonce, rpcx.String(proposal), rpcx.String("approve"), rpcx.String("Vote"))
	suite.Require().Nil(err)
	//recover
	err = suite.LogoutRole(from.String())
	suite.Require().Nil(err)
	err = suite.CheckRoleStatus(from.String(), governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

//tc：治理管理员处于freezing，管理员参与提案，参与成功
func (suite *Model11) Test1110_VoteWithFreezingAdminIsSuccess() {
	pk, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	proposal1, err := suite.RoleToFreezing(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	proposal2 := suite.GetProposal()
	client := suite.NewClient(pk)
	nonce, err := client.GetPendingNonceByAccount(from.String())
	suite.Require().Nil(err)
	_, err = suite.vote(pk, nonce, rpcx.String(proposal2), rpcx.String("approve"), rpcx.String("Vote"))
	suite.Require().Nil(err)
	//recover
	err = suite.VoteReject(proposal1)
	suite.Require().Nil(err)
	err = suite.LogoutRole(from.String())
	suite.Require().Nil(err)
	err = suite.CheckRoleStatus(from.String(), governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

//TODO tc：注销治理管理员导致提案不可能达成，提案自动放弃
//tc：注册审计管理员，管理员未注册，注册成功
func (suite *Model11) Test1111_RegisterAuditAdminWithNoRegisterAdminIsSuccess() {
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	pid, err := suite.MockPid()
	suite.Require().Nil(err)
	pk2, from2, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk2, from2, address2)
	suite.Require().Nil(err)
	err = suite.RegisterNode(from1.String(), "nvpNode", pid, 0, from1.String(), from2)
	suite.Require().Nil(err)
	_, from3, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterRole(from3.String(), AuditAdmin, from1.String())
	suite.Require().Nil(err)
}

//tc：注册审计管理员，管理员处于unavailable，注册成功
func (suite *Model11) Test1112_RegisterAuditAdminWithUnavailableAdminIsSuccess() {
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	pid, err := suite.MockPid()
	suite.Require().Nil(err)
	pk2, from2, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk2, from2, address2)
	suite.Require().Nil(err)
	err = suite.RegisterNode(from1.String(), "nvpNode", pid, 0, from1.String(), from2)
	suite.Require().Nil(err)
	_, from3, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RoleToUnavailable(from3.String(), AuditAdmin, from1.String())
	suite.Require().Nil(err)
	err = suite.RegisterRole(from3.String(), AuditAdmin, from1.String())
	suite.Require().Nil(err)
}

//tc：注销审计管理员，管理员处于available，注销成功
func (suite *Model11) Test1113_LogoutAuditAdminWithAvailableAdminIsSuccess() {
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	pid, err := suite.MockPid()
	suite.Require().Nil(err)
	pk2, from2, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk2, from2, address2)
	suite.Require().Nil(err)
	err = suite.RegisterNode(from1.String(), "nvpNode", pid, 0, from1.String(), from2)
	suite.Require().Nil(err)
	_, from3, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterRole(from3.String(), AuditAdmin, from1.String())
	suite.Require().Nil(err)
	err = suite.LogoutRole(from3.String())
	suite.Require().Nil(err)
}

//tc：注销审计管理员，管理员处于frozen，注销成功
func (suite *Model11) Test1114_LogoutAuditAdminWithFrozenAdminIsSuccess() {
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	pid, err := suite.MockPid()
	suite.Require().Nil(err)
	pk2, from2, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk2, from2, address2)
	suite.Require().Nil(err)
	err = suite.RegisterNode(from1.String(), "nvpNode", pid, 0, from1.String(), from2)
	suite.Require().Nil(err)
	_, from3, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterRole(from3.String(), AuditAdmin, from1.String())
	suite.Require().Nil(err)
	err = suite.LogoutNode(from1.String())
	suite.Require().Nil(err)
	err = suite.LogoutRole(from3.String())
	suite.Require().Nil(err)
}

//tc：注销审计管理员，管理员处于binding，注销成功
func (suite *Model11) Test1115_LogoutAuditAdminWithBindingAdminIsSuccess() {
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	pid, err := suite.MockPid()
	suite.Require().Nil(err)
	pk2, from2, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk2, from2, address2)
	suite.Require().Nil(err)
	err = suite.RegisterNode(from1.String(), "nvpNode", pid, 0, from1.String(), from2)
	suite.Require().Nil(err)
	_, from3, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RoleToBinding(from3.String(), AuditAdmin, from1.String())
	suite.Require().Nil(err)
	err = suite.LogoutRole(from3.String())
	suite.Require().Nil(err)
}

//tc：注销审计管理员，审计节点处于updating，注销成功，审计节点提案通过，审计节点处于available
func (suite *Model11) Test1116_LogoutAuditAdminWithUpdatingNodePassIsSuccess() {
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	pid, err := suite.MockPid()
	suite.Require().Nil(err)
	pk2, from2, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk2, from2, address2)
	suite.Require().Nil(err)
	err = suite.RegisterNode(from1.String(), "nvpNode", pid, 0, from1.String(), from2)
	suite.Require().Nil(err)
	_, from3, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterRole(from3.String(), AuditAdmin, from1.String())
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(from3.String()),
		rpcx.String(from3.String() + "123"),
		rpcx.String(from2),
		rpcx.String("reason"),
	}
	proposal, err := suite.InvokeNodeContract(UpdateNode, args...)
	suite.Require().Nil(err)
	err = suite.LogoutRole(from3.String())
	suite.Require().Nil(err)
	err = suite.VotePass(proposal)
	suite.Require().Nil(err)
	err = suite.CheckNodeStatus(from1.String(), governance.GovernanceAvailable)
	suite.Require().Nil(err)
}

//tc：注销审计管理员，审计节点处于updating，注销成功，审计节点提案不通过，审计节点处于available
func (suite *Model11) Test1117_LogoutAuditAdminWithUpdatingNodeRejectIsSuccess() {
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	pid, err := suite.MockPid()
	suite.Require().Nil(err)
	pk2, from2, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk2, from2, address2)
	suite.Require().Nil(err)
	err = suite.RegisterNode(from1.String(), "nvpNode", pid, 0, from1.String(), from2)
	suite.Require().Nil(err)
	_, from3, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterRole(from3.String(), AuditAdmin, from1.String())
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(from1.String()),
		rpcx.String(from1.String() + "123"),
		rpcx.String(from2),
		rpcx.String("reason"),
	}
	proposal, err := suite.InvokeNodeContract(UpdateNode, args...)
	suite.Require().Nil(err)
	err = suite.LogoutRole(from3.String())
	suite.Require().Nil(err)
	err = suite.VoteReject(proposal)
	suite.Require().Nil(err)
	err = suite.CheckNodeStatus(from1.String(), governance.GovernanceAvailable)
	suite.Require().Nil(err)
}

//tc：注销审计管理员，审计节点处于binding，注销成功，审计节点提案放弃，审计节点处于available
func (suite *Model11) Test1118_LogoutAuditAdminWithBindingNodeIsSuccess() {
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	pid, err := suite.MockPid()
	suite.Require().Nil(err)
	pk2, from2, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk2, from2, address2)
	suite.Require().Nil(err)
	err = suite.RegisterNode(from1.String(), "nvpNode", pid, 0, from1.String(), from2)
	suite.Require().Nil(err)
	_, from3, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RoleToBinding(from3.String(), AuditAdmin, from1.String())
	suite.Require().Nil(err)
	err = suite.LogoutRole(from3.String())
	suite.Require().Nil(err)
}

//tc：注销审计管理员，审计节点处于binded，注销成功
func (suite *Model11) Test1119_LogoutAuditAdminWithBindedNodeIsSuccess() {
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	pid, err := suite.MockPid()
	suite.Require().Nil(err)
	pk2, from2, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk2, from2, address2)
	suite.Require().Nil(err)
	err = suite.RegisterNode(from1.String(), "nvpNode", pid, 0, from1.String(), from2)
	suite.Require().Nil(err)
	_, from3, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterRole(from3.String(), AuditAdmin, from1.String())
	suite.Require().Nil(err)
	err = suite.LogoutRole(from3.String())
	suite.Require().Nil(err)
}

//tc：注销审计管理员，审计节点处于logouting，注销成功，审计节点提案通过，审计节点处于forbidden
func (suite *Model11) Test1120_LogoutRoleWithLogoutingNodePassIsSuccess() {
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	pid, err := suite.MockPid()
	suite.Require().Nil(err)
	pk2, from2, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk2, from2, address2)
	suite.Require().Nil(err)
	err = suite.RegisterNode(from1.String(), "nvpNode", pid, 0, from1.String(), from2)
	suite.Require().Nil(err)
	_, from3, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterRole(from3.String(), AuditAdmin, from1.String())
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(from1.String()),
		rpcx.String("reason"),
	}
	proposal, err := suite.InvokeNodeContract(LogoutNode, args...)
	suite.Require().Nil(err)
	err = suite.CheckNodeStatus(from1.String(), governance.GovernanceLogouting)
	suite.Require().Nil(err)
	err = suite.LogoutRole(from3.String())
	suite.Require().Nil(err)
	err = suite.VotePass(proposal)
	suite.Require().Nil(err)
	err = suite.CheckNodeStatus(from1.String(), governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

//tc：注销审计管理员，审计节点处于logouting，注销成功，审计节点提案不通过，审计节点处于available
func (suite *Model11) Test1121_LogoutAuditAdminWithLogoutingRejectIsSuccess() {
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	pid, err := suite.MockPid()
	suite.Require().Nil(err)
	pk2, from2, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk2, from2, address2)
	suite.Require().Nil(err)
	err = suite.RegisterNode(from1.String(), "nvpNode", pid, 0, from1.String(), from2)
	suite.Require().Nil(err)
	_, from3, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterRole(from3.String(), AuditAdmin, from1.String())
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(from1.String()),
		rpcx.String("reason"),
	}
	proposal, err := suite.InvokeNodeContract(LogoutNode, args...)
	suite.Require().Nil(err)
	err = suite.CheckNodeStatus(from1.String(), governance.GovernanceLogouting)
	suite.Require().Nil(err)
	err = suite.LogoutRole(from3.String())
	suite.Require().Nil(err)
	err = suite.VoteReject(proposal)
	suite.Require().Nil(err)
	err = suite.CheckNodeStatus(from1.String(), governance.GovernanceAvailable)
	suite.Require().Nil(err)
}

//tc：审计管理员重新绑定审计节点，审计节点处于available，绑定成功
func (suite *Model11) Test1122_BindRoleWithAvailableRoleIsSuccess() {
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	pid1, err := suite.MockPid()
	suite.Require().Nil(err)
	pk2, from2, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk2, from2, address2)
	suite.Require().Nil(err)
	err = suite.RegisterNode(from1.String(), "nvpNode", pid1, 0, from1.String(), from2)
	suite.Require().Nil(err)
	_, from3, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterRole(from3.String(), AuditAdmin, from1.String())
	suite.Require().Nil(err)
	err = suite.LogoutNode(from1.String())
	suite.Require().Nil(err)
	_, from4, err := repo.KeyPriv()
	suite.Require().Nil(err)
	pid2, err := suite.MockPid()
	suite.Require().Nil(err)
	pk5, from5, address5, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk5, from5, address5)
	suite.Require().Nil(err)
	err = suite.RegisterNode(from4.String(), "nvpNode", pid2, 0, from4.String(), from5)
	suite.Require().Nil(err)
	err = suite.BindRole(from3.String(), from4.String())
	suite.Require().Nil(err)
}

//tc：审计管理员绑定中，注销审计节点，节点注销成功
func (suite *Model11) Test1123_BindRoleThenLogoutNodeIsSuccess() {
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	pid1, err := suite.MockPid()
	suite.Require().Nil(err)
	pk2, from2, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk2, from2, address2)
	suite.Require().Nil(err)
	err = suite.RegisterNode(from1.String(), "nvpNode", pid1, 0, from1.String(), from2)
	suite.Require().Nil(err)
	_, from3, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterRole(from3.String(), AuditAdmin, from1.String())
	suite.Require().Nil(err)
	err = suite.LogoutNode(from1.String())
	suite.Require().Nil(err)
	_, from4, err := repo.KeyPriv()
	suite.Require().Nil(err)
	pid2, err := suite.MockPid()
	suite.Require().Nil(err)
	pk5, from5, address5, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk5, from5, address5)
	suite.Require().Nil(err)
	err = suite.RegisterNode(from4.String(), "nvpNode", pid2, 0, from4.String(), from5)
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(from3.String()),
		rpcx.String(from4.String()),
		rpcx.String("reason"),
	}
	_, err = suite.InvokeRoleContract(BindRole, args...)
	suite.Require().Nil(err)
	err = suite.LogoutNode(from4.String())
	suite.Require().Nil(err)
}

//tc：审计管理员绑定中，注销审计管理员，管理员注销成功
func (suite *Model11) Test1124_BindRoleThenLogoutRoleIsSuccess() {
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	pid1, err := suite.MockPid()
	suite.Require().Nil(err)
	pk2, from2, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk2, from2, address2)
	suite.Require().Nil(err)
	err = suite.RegisterNode(from1.String(), "nvpNode", pid1, 0, from1.String(), from2)
	suite.Require().Nil(err)
	_, from3, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterRole(from3.String(), AuditAdmin, from1.String())
	suite.Require().Nil(err)
	err = suite.LogoutNode(from1.String())
	suite.Require().Nil(err)
	_, from4, err := repo.KeyPriv()
	suite.Require().Nil(err)
	pid2, err := suite.MockPid()
	suite.Require().Nil(err)
	pk5, from5, address5, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk5, from5, address5)
	suite.Require().Nil(err)
	err = suite.RegisterNode(from4.String(), "nvpNode", pid2, 0, from4.String(), from5)
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(from3.String()),
		rpcx.String(from4.String()),
		rpcx.String("reason"),
	}
	_, err = suite.InvokeRoleContract(BindRole, args...)
	suite.Require().Nil(err)
	err = suite.LogoutRole(from3.String())
	suite.Require().Nil(err)
}

// InvokeRoleContract invoke role contract by methodName and args
func (suite *Snake) InvokeRoleContract(method string, args ...*pb.Arg) (string, error) {
	pk, _, err := repo.Node2Priv()
	if err != nil {
		return "", err
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.RoleContractAddr.Address(), method, &rpcx.TransactOpts{
		Nonce: atomic.AddUint64(&nonce2, 1),
	}, args...)
	if err != nil {
		return "", err
	}
	if res.Status == pb.Receipt_FAILED {
		return "", fmt.Errorf(string(res.Ret))
	}
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	if err != nil {
		return "", err
	}
	return result.ProposalID, nil
}

// RegisterRole register role
func (suite *Snake) RegisterRole(id, typ, account string) error {
	args := []*pb.Arg{
		rpcx.String(id),
		rpcx.String(typ),
		rpcx.String(account),
		rpcx.String("reason"),
	}
	proposal, err := suite.InvokeRoleContract(RegisterRole, args...)
	if err != nil {
		return err
	}
	return suite.VotePass(proposal)
}

// FreezeRole freeze role
func (suite *Snake) FreezeRole(id string) error {
	args := []*pb.Arg{
		rpcx.String(id),
		rpcx.String("reason"),
	}
	proposal, err := suite.InvokeRoleContract(FreezeRole, args...)
	if err != nil {
		return err
	}
	return suite.VotePass(proposal)
}

// ActivateRole activate role
func (suite *Snake) ActivateRole(id string) error {
	args := []*pb.Arg{
		rpcx.String(id),
		rpcx.String("reason"),
	}
	proposal, err := suite.InvokeRoleContract(ActivateRole, args...)
	if err != nil {
		return err
	}
	return suite.VotePass(proposal)
}

// LogoutRole logout role
func (suite *Snake) LogoutRole(id string) error {
	args := []*pb.Arg{
		rpcx.String(id),
		rpcx.String("reason"),
	}
	proposal, err := suite.InvokeRoleContract(LogoutRole, args...)
	if err != nil {
		return err
	}
	return suite.VotePass(proposal)
}

// BindRole bind role
func (suite *Snake) BindRole(roleId, nodeAccount string) error {
	args := []*pb.Arg{
		rpcx.String(roleId),
		rpcx.String(nodeAccount),
		rpcx.String("reason"),
	}
	proposal, err := suite.InvokeRoleContract(BindRole, args...)
	if err != nil {
		return err
	}
	return suite.VotePass(proposal)
}

// RoleToRegisting get a registing role
func (suite *Snake) RoleToRegisting(id, typ, account string) (string, error) {
	args := []*pb.Arg{
		rpcx.String(id),
		rpcx.String(typ),
		rpcx.String(account),
		rpcx.String("reason"),
	}
	proposal, err := suite.InvokeRoleContract(RegisterRole, args...)
	if err != nil {
		return "", err
	}
	err = suite.CheckRoleStatus(id, governance.GovernanceRegisting)
	if err != nil {
		return "", err
	}
	return proposal, nil
}

// RoleToUnavailable get an Unavailable role
func (suite *Snake) RoleToUnavailable(id, typ, account string) error {
	args := []*pb.Arg{
		rpcx.String(id),
		rpcx.String(typ),
		rpcx.String(account),
		rpcx.String("reason"),
	}
	proposal, err := suite.InvokeRoleContract(RegisterRole, args...)
	if err != nil {
		return err
	}
	err = suite.VoteReject(proposal)
	if err != nil {
		return err
	}
	err = suite.CheckRoleStatus(id, governance.GovernanceUnavailable)
	if err != nil {
		return err
	}
	return nil
}

// RoleToFreezing get a freezing role
func (suite *Snake) RoleToFreezing(id, typ, account string) (string, error) {
	err := suite.RegisterRole(id, typ, account)
	if err != nil {
		return "", err
	}
	args := []*pb.Arg{
		rpcx.String(id),
		rpcx.String("reason"),
	}
	proposal, err := suite.InvokeRoleContract(FreezeRole, args...)
	if err != nil {
		return "", err
	}
	err = suite.CheckRoleStatus(id, governance.GovernanceFreezing)
	if err != nil {
		return "", err
	}
	return proposal, nil
}

// RoleToFrozen get a frozen role
func (suite *Snake) RoleToFrozen(id, typ, account string) error {
	err := suite.RegisterRole(id, typ, account)
	if err != nil {
		return err
	}
	err = suite.FreezeRole(id)
	if err != nil {
		return err
	}
	err = suite.CheckRoleStatus(id, governance.GovernanceFrozen)
	if err != nil {
		return err
	}
	return nil
}

// RoleToActivating get an activating role
func (suite *Snake) RoleToActivating(id, typ, account string) (string, error) {
	err := suite.RegisterRole(id, typ, account)
	if err != nil {
		return "", err
	}
	err = suite.FreezeRole(id)
	if err != nil {
		return "", err
	}
	args := []*pb.Arg{
		rpcx.String(id),
		rpcx.String("reason"),
	}
	proposal, err := suite.InvokeRoleContract(ActivateRole, args...)
	if err != nil {
		return "", err
	}
	err = suite.CheckRoleStatus(id, governance.GovernanceActivating)
	if err != nil {
		return "", err
	}
	return proposal, nil
}

// RoleToLogouting get a logouting role
func (suite *Snake) RoleToLogouting(id, typ, account string) (string, error) {
	err := suite.RegisterRole(id, typ, account)
	if err != nil {
		return "", err
	}
	args := []*pb.Arg{
		rpcx.String(id),
		rpcx.String("reason"),
	}
	proposal, err := suite.InvokeRoleContract(LogoutRole, args...)
	if err != nil {
		return "", err
	}
	err = suite.CheckRoleStatus(id, governance.GovernanceLogouting)
	if err != nil {
		return "", err
	}
	return proposal, nil
}

// RoleToForbidden get a forbidden role
func (suite *Snake) RoleToForbidden(id, typ, account string) error {
	err := suite.RegisterRole(id, typ, account)
	if err != nil {
		return err
	}
	err = suite.LogoutRole(id)
	if err != nil {
		return err
	}
	err = suite.CheckRoleStatus(id, governance.GovernanceForbidden)
	if err != nil {
		return err
	}
	return nil
}

// RoleToBinding get a binding role
func (suite *Snake) RoleToBinding(id, typ, account string) error {
	err := suite.RegisterRole(id, typ, account)
	if err != nil {
		return err
	}
	err = suite.LogoutNode(account)
	if err != nil {
		return err
	}
	_, from1, err := repo.KeyPriv()
	if err != nil {
		return err
	}
	pid, err := suite.MockPid()
	if err != nil {
		return err
	}
	pk2, from2, address2, err := suite.DeployRule()
	if err != nil {
		return err
	}
	err = suite.RegisterAppchain(pk2, from2, address2)
	if err != nil {
		return err
	}
	err = suite.RegisterNode(from1.String(), "nvpNode", pid, 0, from1.String(), from2)
	if err != nil {
		return err
	}
	args := []*pb.Arg{
		rpcx.String(id),
		rpcx.String(from1.String()),
		rpcx.String("reason"),
	}
	_, err = suite.InvokeRoleContract(BindRole, args...)
	if err != nil {
		return err
	}
	err = suite.CheckRoleStatus(id, governance.GovernanceBinding)
	if err != nil {
		return err
	}
	return nil
}

// CheckRoleStatus check role status by id
func (suite *Snake) CheckRoleStatus(id string, expectStatus governance.GovernanceStatus) error {
	status, err := suite.GetRoleStatus(id)
	if err != nil {
		return err
	}
	if expectStatus != status {
		return fmt.Errorf("expect status is %s ,but get status %s", expectStatus, status)
	}
	return nil
}

// GetRoleStatus get role status by id
func (suite *Snake) GetRoleStatus(id string) (governance.GovernanceStatus, error) {
	args := []*pb.Arg{
		rpcx.String(id),
	}
	pk, _, err := repo.Node2Priv()
	if err != nil {
		return "", err
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.RoleContractAddr.Address(), "GetRoleInfoById", &rpcx.TransactOpts{
		Nonce: atomic.AddUint64(&nonce2, 1),
	}, args...)
	if err != nil {
		return "", err
	}
	role := &Role{}
	err = json.Unmarshal(res.Ret, role)
	if err != nil {
		return "", err
	}
	return role.Status, nil
}

// GetProposal return a proposal id
func (suite *Snake) GetProposal() string {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	bytes, err := pk.PublicKey().Bytes()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(from),           //chainID
		rpcx.String(from),           //chainName
		rpcx.Bytes(bytes),           //pubKey
		rpcx.String("Flato V1.0.3"), //chainType
		rpcx.Bytes([]byte("")),      //trustRoot
		rpcx.String("0x857133c5C69e6Ce66F7AD46F200B9B3573e77582"), //broker
		rpcx.String("desc"),               //desc
		rpcx.String(address),              //masterRuleAddr
		rpcx.String("https://github.com"), //masterRuleUrl
		rpcx.String(from),                 //adminAddrs
		rpcx.String("reason"),             //reason
	}
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "RegisterAppchain", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	suite.Require().Nil(err)
	return result.ProposalID
}
