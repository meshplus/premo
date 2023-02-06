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

//tc：冻结超级管理员，冻结失败
func (suite *Model11) Test1101_FreezeSuperAdminIsFail() {
	_, from, err := repo.Node1Priv()
	suite.Require().Nil(err)
	err = suite.FreezeRole(from.String())
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "does not support freeze")
}

//tc：激活超级管理员，激活失败
func (suite *Model11) Test1102_ActivateSuperAdminIsFail() {
	_, from, err := repo.Node1Priv()
	suite.Require().Nil(err)
	err = suite.ActivateRole(from.String())
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "can not do activate")
}

//tc：注销超级管理员，注销失败
func (suite *Model11) Test1103_LogoutSuperAdminIsFail() {
	_, from, err := repo.Node1Priv()
	suite.Require().Nil(err)
	err = suite.LogoutRole(from.String())
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "does not support logout")
}

//tc：注册治理管理员，管理员未注册，注册成功
func (suite *Model11) Test1104_RegisterAdminIsSuccess() {
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

//tc：注册治理管理员，管理员处于available，注册失败
func (suite *Model11) Test1105_RegisterAdminWithAvailableAdminIsFail() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterRole(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	err = suite.RegisterRole(from.String(), GovernanceAdmin, "")
	suite.Require().NotNil(err)
	//recover
	err = suite.LogoutRole(from.String())
	suite.Require().Nil(err)
	err = suite.CheckRoleStatus(from.String(), governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

//tc：注册治理管理员，管理员处于registing，注册失败
func (suite *Model11) Test1106_RegisterAdminWithRegistingAdminIsFail() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	proposal, err := suite.RoleToRegisting(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	err = suite.RegisterRole(from.String(), GovernanceAdmin, "")
	suite.Require().NotNil(err)
	//recover
	err = suite.VoteReject(proposal)
	suite.Require().Nil(err)
	err = suite.CheckRoleStatus(from.String(), governance.GovernanceUnavailable)
	suite.Require().Nil(err)
}

//tc：注册治理管理员，管理员处于unavailable，注册成功
func (suite *Model11) Test1107_RegisterAdminWithUnavailableAdminIsSuccess() {
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

//tc：注册治理管理员，管理员处于freezing，注册失败
func (suite *Model11) Test1108_RegisterAdminWithFreezingAdminIsFail() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	proposal, err := suite.RoleToFreezing(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	err = suite.RegisterRole(from.String(), GovernanceAdmin, "")
	suite.Require().NotNil(err)
	//recover
	err = suite.VotePass(proposal)
	suite.Require().Nil(err)
	err = suite.LogoutRole(from.String())
	suite.Require().Nil(err)
	err = suite.CheckRoleStatus(from.String(), governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

//tc：注册治理管理员，管理员处于frozen，注册失败
func (suite *Model11) Test1109_RegisterAdminWithFrozenAdminIsFail() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RoleToFrozen(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	err = suite.RegisterRole(from.String(), GovernanceAdmin, "")
	suite.Require().NotNil(err)
	//recover
	err = suite.LogoutRole(from.String())
	suite.Require().Nil(err)
	err = suite.CheckRoleStatus(from.String(), governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

//tc：注册治理管理员，管理员处于activating，注册失败
func (suite *Model11) Test1110_RegisterAdminWithActivatingAdminIsFail() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	proposal, err := suite.RoleToActivating(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	err = suite.RegisterRole(from.String(), GovernanceAdmin, "")
	suite.Require().NotNil(err)
	//recover
	err = suite.VotePass(proposal)
	suite.Require().Nil(err)
	err = suite.LogoutRole(from.String())
	suite.Require().Nil(err)
	err = suite.CheckRoleStatus(from.String(), governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

//tc：注册治理管理员，管理员处于logouting，注册失败
func (suite *Model11) Test1111_RegisterAdminWithLogoutingAdminIsFail() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	proposal, err := suite.RoleToLogouting(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	err = suite.RegisterRole(from.String(), GovernanceAdmin, "")
	suite.Require().NotNil(err)
	//recover
	err = suite.VotePass(proposal)
	suite.Require().Nil(err)
	err = suite.CheckRoleStatus(from.String(), governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

//tc：注册治理管理员，管理员处于forbidden，注册失败
func (suite *Model11) Test1112_RegisterAdminWithForbiddenAdminIsFail() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RoleToForbidden(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	err = suite.RegisterRole(from.String(), GovernanceAdmin, "")
	suite.Require().NotNil(err)
	err = suite.CheckRoleStatus(from.String(), governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

//tc：冻结治理管理员，管理员未注册，冻结失败
func (suite *Model11) Test1113_FreezeAdminWithNoRegisterAdminIsFail() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.FreezeRole(from.String())
	suite.Require().NotNil(err)
}

//tc：冻结治理管理员，管理员处于available，冻结成功
func (suite *Model11) Test1114_FreezeAdminWithAvailableAdminIsSuccess() {
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

//tc：冻结治理管理员，管理员处于registing，冻结失败
func (suite *Model11) Test1115_FreezeAdminWithRegistingAdminIsFail() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	proposal, err := suite.RoleToRegisting(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	err = suite.FreezeRole(from.String())
	suite.Require().NotNil(err)
	//recover
	err = suite.VoteReject(proposal)
	suite.Require().Nil(err)
	err = suite.CheckRoleStatus(from.String(), governance.GovernanceUnavailable)
	suite.Require().Nil(err)
}

//tc：冻结治理管理员，管理员处于unavailable，冻结失败
func (suite *Model11) Test1116_FreezeAdminWithUnavailableAdminIsFail() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RoleToUnavailable(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	err = suite.FreezeRole(from.String())
	suite.Require().NotNil(err)
}

//tc：冻结治理管理员，管理员处于freezing，冻结失败
func (suite *Model11) Test1117_FreezeAdminWithFreezingAdminIsFail() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	proposal, err := suite.RoleToFreezing(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	err = suite.FreezeRole(from.String())
	suite.Require().NotNil(err)
	//recover
	err = suite.VoteReject(proposal)
	suite.Require().Nil(err)
	err = suite.LogoutRole(from.String())
	suite.Require().Nil(err)
	err = suite.CheckRoleStatus(from.String(), governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

//tc：冻结治理管理员，管理员处于frozen，冻结失败
func (suite *Model11) Test1118_FreezeAdminWithFrozenAdminIsFail() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RoleToFrozen(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	err = suite.FreezeRole(from.String())
	suite.Require().NotNil(err)
	//recover
	err = suite.LogoutRole(from.String())
	suite.Require().Nil(err)
	err = suite.CheckRoleStatus(from.String(), governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

//tc：冻结治理管理员，管理员处于activating，冻结失败
func (suite *Model11) Test1119_FreezeAdminWithActivatingAdminIsFail() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	proposal, err := suite.RoleToActivating(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	err = suite.FreezeRole(from.String())
	suite.Require().NotNil(err)
	//recover
	err = suite.VotePass(proposal)
	suite.Require().Nil(err)
	err = suite.LogoutRole(from.String())
	suite.Require().Nil(err)
	err = suite.CheckRoleStatus(from.String(), governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

//tc：冻结治理管理员，管理员处于logouting，冻结失败
func (suite *Model11) Test1120_FreezeAdminWithLogoutingAdminIsFail() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	proposal, err := suite.RoleToLogouting(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	err = suite.FreezeRole(from.String())
	suite.Require().NotNil(err)
	//recover
	err = suite.VotePass(proposal)
	suite.Require().Nil(err)
	err = suite.CheckRoleStatus(from.String(), governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

//tc：冻结治理管理员，管理员处于forbidden，冻结失败
func (suite *Model11) Test1121_FreezeAdminWithForbiddenAdminIsFail() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RoleToForbidden(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	err = suite.FreezeRole(from.String())
	suite.Require().NotNil(err)
	err = suite.CheckRoleStatus(from.String(), governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

//tc：治理管理员冻结自己，冻结失败
func (suite *Model11) Test1122_FreezeAdminWithAdminSelfIsFail() {
	pk, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterRole(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(from.String()),
		rpcx.String("reason"),
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.RoleContractAddr.Address(), FreezeRole, nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
}

//tc：激活治理管理员，管理员未注册，激活失败
func (suite *Model11) Test1123_ActivateAdminWithNoRegisterAdminIsFail() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.ActivateRole(from.String())
	suite.Require().NotNil(err)
}

//tc：激活治理管理员，管理员处于available，激活失败
func (suite *Model11) Test1124_ActivateAdminWithAvailableAdminIsFail() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterRole(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	err = suite.ActivateRole(from.String())
	suite.Require().NotNil(err)
	//recover
	err = suite.LogoutRole(from.String())
	suite.Require().Nil(err)
	err = suite.CheckRoleStatus(from.String(), governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

//tc：激活治理管理员，管理员处于registing，激活失败
func (suite *Model11) Test1125_ActivateAdminWithRegistingAdminIsFail() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	proposal, err := suite.RoleToRegisting(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	err = suite.ActivateRole(from.String())
	suite.Require().NotNil(err)
	//recover
	err = suite.VoteReject(proposal)
	suite.Require().Nil(err)
	err = suite.CheckRoleStatus(from.String(), governance.GovernanceUnavailable)
	suite.Require().Nil(err)
}

//tc：激活治理管理员，管理员处于unavailable，激活失败
func (suite *Model11) Test1126_ActivateAdminWithUnavailableAdminIsFail() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RoleToUnavailable(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	err = suite.ActivateRole(from.String())
	suite.Require().NotNil(err)
}

//tc：激活治理管理员，管理员处于freezing，激活失败
func (suite *Model11) Test1127_ActivateAdminWithFreezingAdminIsFail() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	proposal, err := suite.RoleToFreezing(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	err = suite.ActivateRole(from.String())
	suite.Require().NotNil(err)
	//recover
	err = suite.VoteReject(proposal)
	suite.Require().Nil(err)
	err = suite.LogoutRole(from.String())
	suite.Require().Nil(err)
	err = suite.CheckRoleStatus(from.String(), governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

//tc：激活治理管理员，管理员处于frozen，激活成功
func (suite *Model11) Test1128_ActivateAdminWithFrozenAdminIsSuccess() {
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

//tc：激活治理管理员，管理员处于activating，激活失败
func (suite *Model11) Test1129_ActivateAdminWithActivatingAdminIsFail() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	proposal, err := suite.RoleToActivating(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	err = suite.ActivateRole(from.String())
	suite.Require().NotNil(err)
	//recover
	err = suite.VoteReject(proposal)
	suite.Require().Nil(err)
	err = suite.LogoutRole(from.String())
	suite.Require().Nil(err)
	err = suite.CheckRoleStatus(from.String(), governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

//tc：激活治理管理员，管理员处于logouting，激活失败
func (suite *Model11) Test1130_ActivateAdminWithLogoutingAdminIsFail() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	proposal, err := suite.RoleToLogouting(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	err = suite.ActivateRole(from.String())
	suite.Require().NotNil(err)
	//recover
	err = suite.VotePass(proposal)
	suite.Require().Nil(err)
	err = suite.CheckRoleStatus(from.String(), governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

//tc：激活治理管理员，管理员处于forbidden，激活失败
func (suite *Model11) Test1131_ActivateAdminWithForbiddenAdminIsFail() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RoleToForbidden(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	err = suite.ActivateRole(from.String())
	suite.Require().NotNil(err)
}

//tc：注销治理管理员，管理员未注册，注销失败
func (suite *Model11) Test1132_LogoutAdminWithNoRegisterAdminIsFail() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.LogoutRole(from.String())
	suite.Require().NotNil(err)
}

//tc：注销治理管理员，管理员处于available，注销成功
func (suite *Model11) Test1133_LogoutAdminWithAvailableAdminIsSuccess() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterRole(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	err = suite.LogoutRole(from.String())
	suite.Require().Nil(err)
}

//tc：注销治理管理员，管理员处于registing，注销失败
func (suite *Model11) Test1134_LogoutAdminWithRegistingAdminIsFail() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	proposal, err := suite.RoleToRegisting(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	err = suite.LogoutRole(from.String())
	suite.Require().NotNil(err)
	//recover
	err = suite.VoteReject(proposal)
	suite.Require().Nil(err)
	err = suite.CheckRoleStatus(from.String(), governance.GovernanceUnavailable)
	suite.Require().Nil(err)
}

//tc：注销治理管理员，管理员处于unavailable，注销失败
func (suite *Model11) Test1135_LogoutAdminWithUnavailableAdminIsFail() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RoleToUnavailable(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	err = suite.LogoutRole(from.String())
	suite.Require().NotNil(err)
}

//tc：注销治理管理员，管理员处于freezing，注销成功
func (suite *Model11) Test1136_LogoutAdminWithFreezingAdminIsSuccess() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	_, err = suite.RoleToFreezing(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	err = suite.LogoutRole(from.String())
	suite.Require().Nil(err)
}

//tc：注销治理管理员，管理员处于frozen，注销成功
func (suite *Model11) Test1137_LogoutAdminWithFrozenAdminIsSuccess() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RoleToFrozen(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	err = suite.LogoutRole(from.String())
	suite.Require().Nil(err)
}

//tc：注销治理管理员，管理员处于activating，注销成功
func (suite *Model11) Test1138_LogoutAdminWithActivatingAdminIsSuccess() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	_, err = suite.RoleToActivating(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	err = suite.LogoutRole(from.String())
	suite.Require().Nil(err)
}

//tc：注销治理管理员，管理员处于logouting，注销失败
func (suite *Model11) Test1139_LogoutAdminWithLogoutingAdminIsFail() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	_, err = suite.RoleToLogouting(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	err = suite.LogoutRole(from.String())
	suite.Require().NotNil(err)
}

//tc：注销治理管理员，管理员处于forbidden，注销失败
func (suite *Model11) Test1140_LogoutAdminWithForbiddenAdminIsFail() {
	_, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RoleToForbidden(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	err = suite.LogoutRole(from.String())
	suite.Require().NotNil(err)
}

//tc：治理管理员处于available，管理员参与提案，参与成功
func (suite *Model11) Test1141_VoteWithAvailableAdminIsSuccess() {
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

//tc：治理管理员处于registing，管理员参与提案，参与失败
func (suite *Model11) Test1142_VoteWithRegistingAdminIsFail() {
	pk, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	proposal1, err := suite.RoleToRegisting(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	proposal2 := suite.GetProposal()
	client := suite.NewClient(pk)
	nonce, err := client.GetPendingNonceByAccount(from.String())
	suite.Require().Nil(err)
	_, err = suite.vote(pk, nonce, rpcx.String(proposal2), rpcx.String("approve"), rpcx.String("Vote"))
	suite.Require().NotNil(err)
	//recover
	err = suite.VoteReject(proposal1)
	suite.Require().Nil(err)
	err = suite.CheckRoleStatus(from.String(), governance.GovernanceUnavailable)
	suite.Require().Nil(err)
}

//tc：治理管理员处于unavailable，管理员参与提案，参与失败
func (suite *Model11) Test1143_VoteWithUnavailableAdminIsFail() {
	pk, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RoleToUnavailable(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	proposal := suite.GetProposal()
	client := suite.NewClient(pk)
	nonce, err := client.GetPendingNonceByAccount(from.String())
	suite.Require().Nil(err)
	_, err = suite.vote(pk, nonce, rpcx.String(proposal), rpcx.String("approve"), rpcx.String("Vote"))
	suite.Require().NotNil(err)
}

//tc：治理管理员处于freezing，管理员参与提案，参与成功
func (suite *Model11) Test1144_VoteWithFreezingAdminIsSuccess() {
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

//tc：治理管理员处于frozen，管理员参与提案，参与失败
func (suite *Model11) Test1145_VoteWithFrozenAdminIsFail() {
	pk, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RoleToFrozen(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	proposal2 := suite.GetProposal()
	client := suite.NewClient(pk)
	nonce, err := client.GetPendingNonceByAccount(from.String())
	suite.Require().Nil(err)
	_, err = suite.vote(pk, nonce, rpcx.String(proposal2), rpcx.String("approve"), rpcx.String("Vote"))
	suite.Require().NotNil(err)
	//recover
	err = suite.LogoutRole(from.String())
	suite.Require().Nil(err)
	err = suite.CheckRoleStatus(from.String(), governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

//tc：治理管理员处于activating，管理员参与提案，参与失败
func (suite *Model11) Test1146_VoteWithActivatingAdminIsFail() {
	pk, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	proposal1, err := suite.RoleToActivating(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	proposal2 := suite.GetProposal()
	client := suite.NewClient(pk)
	nonce, err := client.GetPendingNonceByAccount(from.String())
	suite.Require().Nil(err)
	_, err = suite.vote(pk, nonce, rpcx.String(proposal2), rpcx.String("approve"), rpcx.String("Vote"))
	suite.Require().NotNil(err)
	//recover
	err = suite.VoteReject(proposal1)
	suite.Require().Nil(err)
	err = suite.LogoutRole(from.String())
	suite.Require().Nil(err)
	err = suite.CheckRoleStatus(from.String(), governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

//tc：治理管理员处于logouting，管理员参与提案，参与失败
func (suite *Model11) Test1147_VoteWithLogoutingAdminIsFail() {
	pk, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	proposal1, err := suite.RoleToLogouting(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	proposal2 := suite.GetProposal()
	client := suite.NewClient(pk)
	nonce, err := client.GetPendingNonceByAccount(from.String())
	suite.Require().Nil(err)
	_, err = suite.vote(pk, nonce, rpcx.String(proposal2), rpcx.String("approve"), rpcx.String("Vote"))
	suite.Require().NotNil(err)
	//recover
	err = suite.VotePass(proposal1)
	suite.Require().Nil(err)
	err = suite.CheckRoleStatus(from.String(), governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

//tc：治理管理员处于forbidden，管理员参与提案，参与失败
func (suite *Model11) Test1148_VoteWithForbiddenAdminIsFail() {
	pk, from, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RoleToForbidden(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	proposal2 := suite.GetProposal()
	client := suite.NewClient(pk)
	nonce, err := client.GetPendingNonceByAccount(from.String())
	suite.Require().Nil(err)
	_, err = suite.vote(pk, nonce, rpcx.String(proposal2), rpcx.String("approve"), rpcx.String("Vote"))
	suite.Require().NotNil(err)
}

//TODO tc：注销治理管理员导致提案不可能达成，提案自动放弃
//tc：注册审计管理员，管理员未注册，注册成功
func (suite *Model11) Test1149_RegisterAuditAdminWithNoRegisterAdminIsSuccess() {
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

//tc：注册审计管理员，管理员处于available，注册失败
func (suite *Model11) Test1150_RegisterAuditAdminWithAvailableAdminIsFail() {
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
	err = suite.RegisterRole(from3.String(), AuditAdmin, from1.String())
	suite.Require().NotNil(err)
}

//tc：注册审计管理员，管理员处于registing，注册失败
func (suite *Model11) Test1151_RegisterAuditAdminWithRegistingAdminIsFail() {
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
	_, err = suite.RoleToRegisting(from3.String(), AuditAdmin, from1.String())
	suite.Require().Nil(err)
	err = suite.RegisterRole(from3.String(), AuditAdmin, from1.String())
	suite.Require().NotNil(err)
}

//tc：注册审计管理员，管理员处于unavailable，注册成功
func (suite *Model11) Test1152_RegisterAuditAdminWithUnavailableAdminIsSuccess() {
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

//tc：注册审计管理员，管理员处于frozen，注册失败
func (suite *Model11) Test1153_RegisterAuditAdminWithFrozenAdminIsFail() {
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
	err = suite.CheckRoleStatus(from3.String(), governance.GovernanceFrozen)
	suite.Require().Nil(err)
	err = suite.RegisterRole(from3.String(), AuditAdmin, from1.String())
	suite.Require().NotNil(err)
}

//tc：注册审计管理员，管理员处于binding，注册失败
func (suite *Model11) Test1154_RegisterAuditAdminWithBindingAdminIsFail() {
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
	err = suite.RegisterRole(from3.String(), AuditAdmin, from1.String())
	suite.Require().NotNil(err)
}

//tc：注册审计管理员，管理员处于logouting，注册失败
func (suite *Model11) Test1155_RegisterAuditAdminWithLogoutingAdminIsFail() {
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
	_, err = suite.RoleToLogouting(from3.String(), AuditAdmin, from1.String())
	suite.Require().Nil(err)
	err = suite.RegisterRole(from3.String(), AuditAdmin, from1.String())
	suite.Require().NotNil(err)
}

//tc：注册审计管理员，管理员处于forbidden，注册失败
func (suite *Model11) Test1156_RegisterAuditAdminWithForbiddenAdminIsFail() {
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
	err = suite.RoleToForbidden(from3.String(), AuditAdmin, from1.String())
	suite.Require().Nil(err)
	err = suite.RegisterRole(from3.String(), AuditAdmin, from1.String())
	suite.Require().NotNil(err)
}

//tc：注册审计管理员，审计节点未注册，注册失败
func (suite *Model11) Test1157_RegisterAuditAdminWithNoRegisterNodeIsFail() {
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	pk2, from2, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk2, from2, address2)
	suite.Require().Nil(err)
	_, from3, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterRole(from3.String(), AuditAdmin, from1.String())
	suite.Require().NotNil(err)
}

//tc：注册审计管理员，审计节点处于registing，注册失败
func (suite *Model11) Test1158_RegisterAuditAdminWithRegistingNodeIsFail() {
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	pid, err := suite.MockPid()
	suite.Require().Nil(err)
	pk2, from2, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk2, from2, address2)
	suite.Require().Nil(err)
	_, err = suite.NodeToRegisting(from1.String(), "nvpNode", pid, 0, from1.String(), from2)
	suite.Require().Nil(err)
	_, from3, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterRole(from3.String(), AuditAdmin, from1.String())
	suite.Require().NotNil(err)
}

//tc：注册审计管理员，审计节点处于unavailable，注册失败
func (suite *Model11) Test1159_RegisterAuditAdminWithUnavailableNodeIsFail() {
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	pid, err := suite.MockPid()
	suite.Require().Nil(err)
	pk2, from2, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk2, from2, address2)
	suite.Require().Nil(err)
	err = suite.NodeToUnavailable(from1.String(), "nvpNode", pid, 0, from1.String(), from2)
	suite.Require().Nil(err)
	_, from3, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterRole(from3.String(), AuditAdmin, from1.String())
	suite.Require().NotNil(err)
}

//tc：注册审计管理员，审计节点处于updating，注册失败
func (suite *Model11) Test1160_RegisterAuditAdminWithUpdatingNodeIsFail() {
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	pid, err := suite.MockPid()
	suite.Require().Nil(err)
	pk2, from2, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk2, from2, address2)
	suite.Require().Nil(err)
	_, err = suite.NodeToUpdating(from1.String(), "nvpNode", pid, 0, from1.String(), from2)
	suite.Require().Nil(err)
	_, from3, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterRole(from3.String(), AuditAdmin, from1.String())
	suite.Require().NotNil(err)
}

//tc：注册审计管理员，审计节点处于binding，注册失败
func (suite *Model11) Test1161_RegisterAuditAdminWithBindingNodeIsFail() {
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	pid, err := suite.MockPid()
	suite.Require().Nil(err)
	pk2, from2, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk2, from2, address2)
	suite.Require().Nil(err)
	_, err = suite.NodeToBinding(from1.String(), "nvpNode", pid, 0, from1.String(), from2)
	suite.Require().Nil(err)
	_, from3, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterRole(from3.String(), AuditAdmin, from1.String())
	suite.Require().NotNil(err)
}

//tc：注册审计管理员，审计节点处于binded，注册失败
func (suite *Model11) Test1162_RegisterAuditAdminWithBindedNodeIsFail() {
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	pid, err := suite.MockPid()
	suite.Require().Nil(err)
	pk2, from2, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk2, from2, address2)
	suite.Require().Nil(err)
	err = suite.NodeToBinded(from1.String(), "nvpNode", pid, 0, from1.String(), from2)
	suite.Require().Nil(err)
	_, from3, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterRole(from3.String(), AuditAdmin, from1.String())
	suite.Require().NotNil(err)
}

//tc：注册审计管理员，审计节点处于logouting，注册失败
func (suite *Model11) Test1163_RegisterAuditAdminWithLogoutingNodeIsFail() {
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	pid, err := suite.MockPid()
	suite.Require().Nil(err)
	pk2, from2, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk2, from2, address2)
	suite.Require().Nil(err)
	_, err = suite.NodeToLogouting(from1.String(), "nvpNode", pid, 0, from1.String(), from2)
	suite.Require().Nil(err)
	_, from3, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterRole(from3.String(), AuditAdmin, from1.String())
	suite.Require().NotNil(err)
}

//tc：注册审计管理员，审计节点处于forbidden，注册失败
func (suite *Model11) Test1164_RegisterAuditAdminWithForbiddenNodeIsFail() {
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	pid, err := suite.MockPid()
	suite.Require().Nil(err)
	pk2, from2, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk2, from2, address2)
	suite.Require().Nil(err)
	err = suite.NodeToForbidden(from1.String(), "nvpNode", pid, 0, from1.String(), from2)
	suite.Require().Nil(err)
	_, from3, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterRole(from3.String(), AuditAdmin, from1.String())
	suite.Require().NotNil(err)
}

//tc：注销审计管理员，管理员未注册，注销失败
func (suite *Model11) Test1165_LogoutAuditAdminWithNoRegisterAdminIsFail() {
	_, from1, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.LogoutRole(from1.String())
	suite.Require().NotNil(err)
}

//tc：注销审计管理员，管理员处于available，注销成功
func (suite *Model11) Test1166_LogoutAuditAdminWithAvailableAdminIsSuccess() {
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

//tc：注销审计管理员，管理员处于registing，注销失败
func (suite *Model11) Test1167_LogoutAuditAdminWithRegistingAdminIsFail() {
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
	_, err = suite.RoleToRegisting(from3.String(), AuditAdmin, from1.String())
	suite.Require().Nil(err)
	err = suite.LogoutRole(from3.String())
	suite.Require().NotNil(err)
}

//tc：注销审计管理员，管理员处于unavailable，注销失败
func (suite *Model11) Test1168_LogoutAuditAdminWithUnavailableAdminIsFail() {
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
	_, err = suite.RoleToRegisting(from3.String(), AuditAdmin, from1.String())
	suite.Require().Nil(err)
	err = suite.LogoutRole(from3.String())
	suite.Require().NotNil(err)
}

//tc：注销审计管理员，管理员处于frozen，注销成功
func (suite *Model11) Test1169_LogoutAuditAdminWithFrozenAdminIsSuccess() {
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
func (suite *Model11) Test1170_LogoutAuditAdminWithBindingAdminIsSuccess() {
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

//tc：注销审计管理员，管理员处于logouting，注销失败
func (suite *Model11) Test1171_LogoutAuditAdminWithLogoutingAdminIsFail() {
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
	_, err = suite.RoleToLogouting(from3.String(), AuditAdmin, from1.String())
	suite.Require().Nil(err)
	err = suite.LogoutRole(from3.String())
	suite.Require().NotNil(err)
}

//tc：注销审计管理员，管理员处于forbidden，注销失败
func (suite *Model11) Test1172_LogoutAuditAdminWithForbiddenAdminIsFail() {
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
	err = suite.RoleToForbidden(from3.String(), AuditAdmin, from1.String())
	suite.Require().Nil(err)
	err = suite.LogoutRole(from3.String())
	suite.Require().NotNil(err)
}

//tc：注销审计管理员，审计节点处于updating，注销成功，审计节点提案通过，审计节点处于available
func (suite *Model11) Test1173_LogoutAuditAdminWithUpdatingNodePassIsSuccess() {
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
func (suite *Model11) Test1174_LogoutAuditAdminWithUpdatingNodeRejectIsSuccess() {
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
func (suite *Model11) Test1175_LogoutAuditAdminWithBindingNodeIsSuccess() {
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
func (suite *Model11) Test1176_LogoutAuditAdminWithBindedNodeIsSuccess() {
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
func (suite *Model11) Test1177_LogoutRoleWithLogoutingNodePassIsSuccess() {
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
func (suite *Model11) Test1178_LogoutAuditAdminWithLogoutingRejectIsSuccess() {
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

//tc：审计管理员重新绑定审计节点，审计节点未注册，绑定失败
func (suite *Model11) Test1179_BindRoleWithNoRegisterRoleIsFail() {
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
	err = suite.BindRole(from3.String(), from4.String())
	suite.Require().NotNil(err)
}

//tc：审计管理员重新绑定审计节点，审计节点处于registing，绑定失败
func (suite *Model11) Test1180_BindRoleWithRegistingRoleIsFail() {
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
	_, err = suite.NodeToRegisting(from4.String(), "nvpNode", pid2, 0, from4.String(), from5)
	suite.Require().Nil(err)
	err = suite.BindRole(from3.String(), from4.String())
	suite.Require().NotNil(err)
}

//tc：审计管理员重新绑定审计节点，审计节点处于unavailable，绑定失败
func (suite *Model11) Test1181_BindRoleWithUnavailableRoleIsFail() {
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
	err = suite.NodeToUnavailable(from4.String(), "nvpNode", pid2, 0, from4.String(), from5)
	suite.Require().Nil(err)
	err = suite.BindRole(from3.String(), from4.String())
	suite.Require().NotNil(err)
}

//tc：审计管理员重新绑定审计节点，审计节点处于available，绑定成功
func (suite *Model11) Test1182_BindRoleWithAvailableRoleIsSuccess() {
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

//tc：审计管理员重新绑定审计节点，审计节点处于binding，绑定失败
func (suite *Model11) Test1183_BindRoleWithBindingRoleIsFail() {
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
	_, err = suite.NodeToBinding(from4.String(), "nvpNode", pid2, 0, from4.String(), from5)
	suite.Require().Nil(err)
	err = suite.BindRole(from3.String(), from4.String())
	suite.Require().NotNil(err)
}

//tc：审计管理员重新绑定审计节点，审计节点处于binded，绑定失败
func (suite *Model11) Test1184_BindRoleWithBindedRoleIsFail() {
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
	err = suite.NodeToBinded(from4.String(), "nvpNode", pid2, 0, from4.String(), from5)
	suite.Require().Nil(err)
	err = suite.BindRole(from3.String(), from4.String())
	suite.Require().NotNil(err)
}

//tc：审计管理员重新绑定审计节点，审计节点处于updating，绑定失败
func (suite *Model11) Test1185_BindRoleWithUpdatingRoleIsFail() {
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
	_, err = suite.NodeToUpdating(from4.String(), "nvpNode", pid2, 0, from4.String(), from5)
	suite.Require().Nil(err)
	err = suite.BindRole(from3.String(), from4.String())
	suite.Require().NotNil(err)
}

//tc：审计管理员重新绑定审计节点，审计节点处于logouting，绑定失败
func (suite *Model11) Test1186_BindRoleWithLogoutingRoleIsFail() {
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
	_, err = suite.NodeToLogouting(from4.String(), "nvpNode", pid2, 0, from4.String(), from5)
	suite.Require().Nil(err)
	err = suite.BindRole(from3.String(), from4.String())
	suite.Require().NotNil(err)
}

//tc：审计管理员重新绑定审计节点，审计节点处于forbidden，绑定失败
func (suite *Model11) Test1187_BindRoleWithForbiddenRoleIsFail() {
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
	err = suite.NodeToForbidden(from4.String(), "nvpNode", pid2, 0, from4.String(), from5)
	suite.Require().Nil(err)
	err = suite.BindRole(from3.String(), from4.String())
	suite.Require().NotNil(err)
}

//tc：审计管理员绑定中，注销审计节点，节点注销成功
func (suite *Model11) Test1188_BindRoleThenLogoutNodeIsSuccess() {
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
func (suite *Model11) Test1189_BindRoleThenLogoutRoleIsSuccess() {
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

//tc：审计管理员参与治理投票，投票失败
func (suite *Model11) Test1190_AuditAdminVoteIsFail() {
	pk1, from1, err := repo.KeyPriv()
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
	proposal := suite.GetProposal()
	client := suite.NewClient(pk1)
	nonce, err := client.GetPendingNonceByAccount(from1.String())
	suite.Require().Nil(err)
	_, err = suite.vote(pk1, nonce, rpcx.String(proposal), rpcx.String("approve"), rpcx.String("Vote"))
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
