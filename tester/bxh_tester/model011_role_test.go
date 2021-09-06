package bxh_tester

import (
	"encoding/json"
	"errors"
	"sync/atomic"

	crypto2 "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/looplab/fsm"
	"github.com/meshplus/bitxhub-core/governance"
	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-kit/crypto/asym/ecdsa"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/meshplus/premo/internal/repo"
)

const (
	RegisterRole         = "RegisterRole"
	UpdateAuditAdminNode = "UpdateAuditAdminNode"
	FreezeRole           = "FreezeRole"
	ActivateRole         = "ActivateRole"
	LogoutRole           = "LogoutRole"
)

type RoleType string
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

type Model11 struct {
	*Snake
}

//tc：更新超级治理管理员，更新失败
func (suite *Model11) Test1101_UpdateSuperAdminIsFail() {
	node1, err := repo.Node1Path()
	suite.Require().Nil(err)
	pk, err := asym.RestorePrivateKey(node1, repo.KeyPassword)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	pid, err := suite.CreatePid()
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String(pid),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(UpdateAuditAdminNode, args...)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "the role is not a AuditAdmin")
}

//tc：冻结超级管理员，冻结失败
func (suite *Model11) Test1102_FreezeSuperAdminIsFail() {
	node1, err := repo.Node1Path()
	suite.Require().Nil(err)
	pk, err := asym.RestorePrivateKey(node1, repo.KeyPassword)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(FreezeRole, args...)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "super governance admin can not be freeze")
}

//tc：激活超级管理员，激活失败
func (suite *Model11) Test1103_ActivateSuperAdminIsFail() {
	node1, err := repo.Node1Path()
	suite.Require().Nil(err)
	pk, err := asym.RestorePrivateKey(node1, repo.KeyPassword)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(ActivateRole, args...)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "can not be activate")
}

//tc：注销超级管理员，注销失败
func (suite *Model11) Test1104_LogoutSuperAdminIsFail() {
	node1, err := repo.Node1Path()
	suite.Require().Nil(err)
	pk, err := asym.RestorePrivateKey(node1, repo.KeyPassword)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(LogoutRole, args...)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "super governance admin can not be logout")
}

//tc：注册普通治理管理员，注册成功
func (suite *Model11) Test1105_RegisterAdminIsSuccess() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("governanceAdmin"),
		rpcx.String(""),
		rpcx.String(""),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(RegisterRole, args...)
	suite.Require().Nil(err)
	role, err := suite.GetRoleInfoById(address.String())
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceAvailable, role.Status)
	//recover
	args = []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(LogoutRole, args...)
	suite.Require().Nil(err)
}

//tc：更新普通治理管理员，更新失败
func (suite Model11) Test1106_UpdateAdminIsFail() {
	node3, err := repo.Node3Path()
	suite.Require().Nil(err)
	pk, err := asym.RestorePrivateKey(node3, repo.KeyPassword)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	pid, err := suite.CreatePid()
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String(pid),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(UpdateAuditAdminNode, args...)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "the role is not a AuditAdmin")
}

//tc：冻结普通治理管理员，冻结成功，发起投票失败，发起提案失败
func (suite Model11) Test1107_FreezeAdminIsSuccess() {
	node3, err := repo.Node3Path()
	suite.Require().Nil(err)
	node3pk, err := asym.RestorePrivateKey(node3, repo.KeyPassword)
	suite.Require().Nil(err)
	from, err := node3pk.PublicKey().Address()
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(from.String()),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(FreezeRole, args...)
	suite.Require().Nil(err)
	role, err := suite.GetRoleInfoById(from.String())
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceFrozen, role.Status)
	address, err := suite.DeploySimpleRule()
	suite.Require().Nil(err)
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	chainID := suite.GetChainID(pk)
	client := suite.NewClient(pk)
	args = []*pb.Arg{
		rpcx.String(chainID),
		rpcx.String(address),
		rpcx.String("reason"),
	}
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), RegisterRule, nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	suite.Require().Nil(err)
	_, err = suite.vote(node3pk, atomic.AddUint64(&nonce3, 1), rpcx.String(result.ProposalID), rpcx.String("approve"), rpcx.String("pass"))
	suite.Require().NotNil(err)
	//recover
	args = []*pb.Arg{
		rpcx.String(from.String()),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(ActivateRole, args...)
	suite.Require().Nil(err)
	role, err = suite.GetRoleInfoById(from.String())
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceAvailable, role.Status)
}

//tc：冻结已冻结的普通治理管理员，冻结失败
func (suite Model11) Test1108_FreezeFrozenAdminIsFail() {
	node3, err := repo.Node3Path()
	suite.Require().Nil(err)
	pk, err := asym.RestorePrivateKey(node3, repo.KeyPassword)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(FreezeRole, args...)
	suite.Require().Nil(err)
	role, err := suite.GetRoleInfoById(address.String())
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceFrozen, role.Status)
	//repeat
	err = suite.InvokeRoleContract(FreezeRole, args...)
	suite.Require().NotNil(err)
	//recover
	err = suite.InvokeRoleContract(ActivateRole, args...)
	suite.Require().Nil(err)
}

//tc：激活冻结的普通治理管理员，激活成功
func (suite Model11) Test1109_ActivateFrozenAdminIsSuccess() {
	node3, err := repo.Node3Path()
	suite.Require().Nil(err)
	pk, err := asym.RestorePrivateKey(node3, repo.KeyPassword)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(FreezeRole, args...)
	suite.Require().Nil(err)
	role, err := suite.GetRoleInfoById(address.String())
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceFrozen, role.Status)
	err = suite.InvokeRoleContract(ActivateRole, args...)
	suite.Require().Nil(err)
	role, err = suite.GetRoleInfoById(address.String())
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceAvailable, role.Status)
}

//tc：激活未冻结的普通治理管理员，激活失败
func (suite Model11) Test1110_ActivateAvailableAdminIsFail() {
	node3, err := repo.Node3Path()
	suite.Require().Nil(err)
	pk, err := asym.RestorePrivateKey(node3, repo.KeyPassword)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	role, err := suite.GetRoleInfoById(address.String())
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceAvailable, role.Status)
	args := []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(ActivateRole, args...)
	suite.Require().NotNil(err)
	role, err = suite.GetRoleInfoById(address.String())
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceAvailable, role.Status)
}

//tc：注销普通治理管理员，注销成功
func (suite Model11) Test1111_LogoutAdminIsSuccess() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("governanceAdmin"),
		rpcx.String(""),
		rpcx.String(""),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(RegisterRole, args...)
	suite.Require().Nil(err)
	role, err := suite.GetRoleInfoById(address.String())
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceAvailable, role.Status)
	args = []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(LogoutRole, args...)
	suite.Require().Nil(err)
}

//tc：注销冻结的普通治理管理员，注销成功
func (suite Model11) Test1112_LogoutFrozenAdminIsSuccess() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("governanceAdmin"),
		rpcx.String(""),
		rpcx.String(""),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(RegisterRole, args...)
	suite.Require().Nil(err)
	args = []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(FreezeRole, args...)
	suite.Require().Nil(err)
	role, err := suite.GetRoleInfoById(address.String())
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceFrozen, role.Status)
	args = []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(LogoutRole, args...)
	suite.Require().Nil(err)
}

//tc：注册审计管理员，注册成功
func (suite Model11) Test1113_RegisterAuditAdminIsSuccess() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	pid, err := suite.CreatePid()
	suite.Require().Nil(err)
	err = suite.RegisterNode(pid, 0, address.String(), "nvpNode")
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("auditAdmin"),
		rpcx.String(pid),
		rpcx.String(""),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(RegisterRole, args...)
	suite.Require().Nil(err)
	role, err := suite.GetRoleInfoById(address.String())
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceAvailable, role.Status)

	//recover
	args = []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(LogoutRole, args...)
	suite.Require().Nil(err)
	err = suite.LogoutNode(pid)
	suite.Require().Nil(err)
}

//tc：更新审计管理员，更新的id和原先一样，更新失败
func (suite Model11) Test1114_UpdateAuditAdminWithSamePidIsFail() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	pid, err := suite.CreatePid()
	suite.Require().Nil(err)
	err = suite.RegisterNode(pid, 0, address.String(), "nvpNode")
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("auditAdmin"),
		rpcx.String(pid),
		rpcx.String(""),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(RegisterRole, args...)
	suite.Require().Nil(err)
	role, err := suite.GetRoleInfoById(address.String())
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceAvailable, role.Status)
	args = []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String(pid),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(UpdateAuditAdminNode, args...)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "the node ID is the same as before")
	//recover
	args = []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(LogoutRole, args...)
	suite.Require().Nil(err)
	err = suite.LogoutNode(pid)
	suite.Require().Nil(err)
}

//tc：更新审计管理员，更新的id不存在，更新失败
func (suite Model11) Test1115_UpdateAuditAdminWithNoExistPidIsFail() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	pid, err := suite.CreatePid()
	suite.Require().Nil(err)
	err = suite.RegisterNode(pid, 0, address.String(), "nvpNode")
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("auditAdmin"),
		rpcx.String(pid),
		rpcx.String(""),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(RegisterRole, args...)
	suite.Require().Nil(err)
	role, err := suite.GetRoleInfoById(address.String())
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceAvailable, role.Status)
	args = []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String(pid[0:len(pid)-1] + "1"),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(UpdateAuditAdminNode, args...)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), " this node does not exist")
	//recover
	args = []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(LogoutRole, args...)
	suite.Require().Nil(err)
	err = suite.LogoutNode(pid)
	suite.Require().Nil(err)
}

//tc：更新审计管理员，更新id正确，更新成功
func (suite Model11) Test1116_UpdateAuditAdminIsSuccess() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	pid, err := suite.CreatePid()
	suite.Require().Nil(err)
	err = suite.RegisterNode(pid, 0, address.String(), "nvpNode")
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("auditAdmin"),
		rpcx.String(pid),
		rpcx.String(""),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(RegisterRole, args...)
	suite.Require().Nil(err)
	role, err := suite.GetRoleInfoById(address.String())
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceAvailable, role.Status)
	newPid, err := suite.CreatePid()
	suite.Require().Nil(err)
	err = suite.RegisterNode(newPid, 0, address.String(), "nvpNode")
	suite.Require().Nil(err)
	args = []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String(newPid),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(UpdateAuditAdminNode, args...)
	suite.Require().Nil(err)
	//recover
	args = []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(LogoutRole, args...)
	suite.Require().Nil(err)
	err = suite.LogoutNode(pid)
	suite.Require().Nil(err)
}

//tc：冻结审计管理员，冻结成功（审计节点尚未实现）
func (suite Model11) Test1117_FreezeAuditAdminIsSuccess() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	pid, err := suite.CreatePid()
	suite.Require().Nil(err)
	err = suite.RegisterNode(pid, 0, address.String(), "nvpNode")
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("auditAdmin"),
		rpcx.String(pid),
		rpcx.String(""),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(RegisterRole, args...)
	suite.Require().Nil(err)
	args = []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(FreezeRole, args...)
	suite.Require().Nil(err)
	role, err := suite.GetRoleInfoById(address.String())
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceFrozen, role.Status)
	//recover
	args = []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(LogoutRole, args...)
	suite.Require().Nil(err)
	err = suite.LogoutNode(pid)
	suite.Require().Nil(err)
}

//tc：冻结已冻结的审计管理员，冻结失败
func (suite Model11) Test1118_FreezeFrozenAuditAdminIsFail() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	pid, err := suite.CreatePid()
	suite.Require().Nil(err)
	err = suite.RegisterNode(pid, 0, address.String(), "nvpNode")
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("auditAdmin"),
		rpcx.String(pid),
		rpcx.String(""),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(RegisterRole, args...)
	suite.Require().Nil(err)
	args = []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(FreezeRole, args...)
	suite.Require().Nil(err)
	role, err := suite.GetRoleInfoById(address.String())
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceFrozen, role.Status)
	args = []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(FreezeRole, args...)
	suite.Require().NotNil(err)
	//recover
	args = []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(LogoutRole, args...)
	suite.Require().Nil(err)
	err = suite.LogoutNode(pid)
	suite.Require().Nil(err)
}

//tc：激活冻结的审计管理员，激活成功
func (suite Model11) Test1119_ActivateFrozenAuditAdminIsSuccess() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	pid, err := suite.CreatePid()
	suite.Require().Nil(err)
	err = suite.RegisterNode(pid, 0, address.String(), "nvpNode")
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("auditAdmin"),
		rpcx.String(pid),
		rpcx.String(""),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(RegisterRole, args...)
	suite.Require().Nil(err)
	args = []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(FreezeRole, args...)
	suite.Require().Nil(err)
	role, err := suite.GetRoleInfoById(address.String())
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceFrozen, role.Status)
	args = []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(ActivateRole, args...)
	suite.Require().Nil(err)
	//recover
	args = []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(LogoutRole, args...)
	suite.Require().Nil(err)
	err = suite.LogoutNode(pid)
	suite.Require().Nil(err)
}

//tc：激活未冻结的审计管理员，激活失败
func (suite Model11) Test1120_ActivateNoFrozenAuditAdminIsFail() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	pid, err := suite.CreatePid()
	suite.Require().Nil(err)
	err = suite.RegisterNode(pid, 0, address.String(), "nvpNode")
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("auditAdmin"),
		rpcx.String(pid),
		rpcx.String(""),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(RegisterRole, args...)
	suite.Require().Nil(err)
	role, err := suite.GetRoleInfoById(address.String())
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceAvailable, role.Status)
	args = []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(ActivateRole, args...)
	suite.Require().NotNil(err)
	//recover
	args = []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(LogoutRole, args...)
	suite.Require().Nil(err)
	err = suite.LogoutNode(pid)
	suite.Require().Nil(err)
}

//tc：注销审计管理员，注销成功
func (suite Model11) Test1121_LogoutAuditAdminIsSuccess() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	pid, err := suite.CreatePid()
	suite.Require().Nil(err)
	err = suite.RegisterNode(pid, 0, address.String(), "nvpNode")
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("auditAdmin"),
		rpcx.String(pid),
		rpcx.String(""),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(RegisterRole, args...)
	suite.Require().Nil(err)
	role, err := suite.GetRoleInfoById(address.String())
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceAvailable, role.Status)
	args = []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(LogoutRole, args...)
	suite.Require().Nil(err)
	err = suite.LogoutNode(pid)
	suite.Require().Nil(err)
}

//tc：注销已冻结的审计管理员，冻结成功
func (suite Model11) Test1122_LogoutFrozenAuditAdminIsSuccess() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	pid, err := suite.CreatePid()
	suite.Require().Nil(err)
	err = suite.RegisterNode(pid, 0, address.String(), "nvpNode")
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("auditAdmin"),
		rpcx.String(pid),
		rpcx.String(""),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(RegisterRole, args...)
	suite.Require().Nil(err)
	args = []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(FreezeRole, args...)
	suite.Require().Nil(err)
	role, err := suite.GetRoleInfoById(address.String())
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceFrozen, role.Status)
	args = []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(LogoutRole, args...)
	suite.Require().Nil(err)
	err = suite.LogoutNode(pid)
	suite.Require().Nil(err)
}

//tc：审计管理员参与治理投票，投票失败
func (suite Model11) Test1123_AuditAdminVoteIsFail() {
	pk1, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	address1, err := pk1.PublicKey().Address()
	suite.Require().Nil(err)
	pid1, err := suite.CreatePid()
	suite.Require().Nil(err)
	err = suite.RegisterNode(pid1, 0, address1.String(), "nvpNode")
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(address1.String()),
		rpcx.String("auditAdmin"),
		rpcx.String(pid1),
		rpcx.String(""),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(RegisterRole, args...)
	suite.Require().Nil(err)
	role, err := suite.GetRoleInfoById(address1.String())
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceAvailable, role.Status)
	pk2, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	address2, err := pk2.PublicKey().Address()
	suite.Require().Nil(err)
	pid2, err := suite.CreatePid()
	suite.Require().Nil(err)
	err = suite.RegisterNode(pid2, 0, address2.String(), "nvpNode")
	suite.Require().Nil(err)
	args = []*pb.Arg{
		rpcx.String(address2.String()),
		rpcx.String("auditAdmin"),
		rpcx.String(pid2),
		rpcx.String(""),
		rpcx.String("reason"),
	}
	node2, err := repo.Node2Path()
	suite.Require().Nil(err)
	privateKey, err := asym.RestorePrivateKey(node2, repo.KeyPassword)
	suite.Require().Nil(err)
	client := suite.NewClient(privateKey)
	res, err := client.InvokeBVMContract(constant.RoleContractAddr.Address(), RegisterRole, &rpcx.TransactOpts{
		Nonce: atomic.AddUint64(&nonce2, 1),
	}, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	client1 := suite.NewClient(pk1)
	nonce, err := client1.GetPendingNonceByAccount(address1.String())
	suite.Require().Nil(err)
	_, err = suite.vote(pk1, nonce, rpcx.String(result.ProposalID), rpcx.String("approve"), rpcx.String("reason"))
	suite.Require().NotNil(err)
	err = suite.VotePass(result.ProposalID)
	suite.Require().Nil(err)
	//recover
	args = []*pb.Arg{
		rpcx.String(address1.String()),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(LogoutRole, args...)
	suite.Require().Nil(err)
	err = suite.LogoutNode(pid1)
	suite.Require().Nil(err)
	args = []*pb.Arg{
		rpcx.String(address2.String()),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(LogoutRole, args...)
	suite.Require().Nil(err)
	err = suite.LogoutNode(pid2)
	suite.Require().Nil(err)
}

//tc：调用GetRole接口，成功查询用户角色结构
func (suite Model11) Test1124_GetRole() {
	node1, err := repo.Node1Path()
	suite.Require().Nil(err)
	privateKey, err := asym.RestorePrivateKey(node1, repo.KeyPassword)
	suite.Require().Nil(err)
	client := suite.NewClient(privateKey)
	res, err := client.InvokeBVMContract(constant.RoleContractAddr.Address(), "GetRole", &rpcx.TransactOpts{
		Nonce: atomic.AddUint64(&nonce1, 1),
	})
	suite.Require().Nil(err)
	suite.Require().Equal("superGovernanceAdmin", string(res.Ret))
}

//tc：调用GetRoleById接口，成功根据正确的iD查询用户角色结构
func (suite Model11) Test1125_GetRoleByAddr() {
	node1, err := repo.Node1Path()
	suite.Require().Nil(err)
	privateKey, err := asym.RestorePrivateKey(node1, repo.KeyPassword)
	suite.Require().Nil(err)
	address, err := privateKey.PublicKey().Address()
	suite.Require().Nil(err)
	client := suite.NewClient(privateKey)
	args := []*pb.Arg{
		rpcx.String(address.String()),
	}
	res, err := client.InvokeBVMContract(constant.RoleContractAddr.Address(), "GetRoleByAddr", &rpcx.TransactOpts{
		Nonce: atomic.AddUint64(&nonce1, 1),
	}, args...)
	suite.Require().Nil(err)
	suite.Require().Equal("superGovernanceAdmin", string(res.Ret))
}

func (suite Model11) InvokeRoleContract(method string, args ...*pb.Arg) error {
	node2, err := repo.Node2Path()
	if err != nil {
		return err
	}
	privateKey, err := asym.RestorePrivateKey(node2, repo.KeyPassword)
	if err != nil {
		return err
	}
	client := suite.NewClient(privateKey)
	res, err := client.InvokeBVMContract(constant.RoleContractAddr.Address(), method, &rpcx.TransactOpts{
		Nonce: atomic.AddUint64(&nonce2, 1),
	}, args...)
	if err != nil {
		return err
	}
	if res.Status == pb.Receipt_FAILED {
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

func (suite Model11) InvokeNodeContract(method string, args ...*pb.Arg) error {
	node2, err := repo.Node2Path()
	if err != nil {
		return err
	}
	privateKey, err := asym.RestorePrivateKey(node2, repo.KeyPassword)
	if err != nil {
		return err
	}
	client := suite.NewClient(privateKey)
	res, err := client.InvokeBVMContract(constant.NodeManagerContractAddr.Address(), method, &rpcx.TransactOpts{
		Nonce: atomic.AddUint64(&nonce2, 1),
	}, args...)
	if err != nil {
		return err
	}
	if res.Status == pb.Receipt_FAILED {
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

func (suite Model11) GetRoleInfoById(roleId string) (*Role, error) {
	args := []*pb.Arg{
		rpcx.String(roleId),
	}
	node2, err := repo.Node2Path()
	if err != nil {
		return nil, err
	}
	privateKey, err := asym.RestorePrivateKey(node2, repo.KeyPassword)
	if err != nil {
		return nil, err
	}
	client := suite.NewClient(privateKey)
	res, err := client.InvokeBVMContract(constant.RoleContractAddr.Address(), "GetRoleInfoById", &rpcx.TransactOpts{
		Nonce: atomic.AddUint64(&nonce2, 1),
	}, args...)
	if err != nil {
		return nil, err
	}
	role := &Role{}
	err = json.Unmarshal(res.Ret, role)
	suite.Require().Nil(err)
	return role, nil
}

func (suite Snake) CreatePid() (string, error) {
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
