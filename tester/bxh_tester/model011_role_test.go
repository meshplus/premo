package bxh_tester

import (
	"encoding/base64"
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

type Role struct {
	ID       string `toml:"id" json:"id"`
	RoleType string `toml:"role_type" json:"role_type"`

	// 	GovernanceAdmin info
	Weight uint64 `json:"weight" toml:"weight"`

	// AuditAdmin info
	NodePid string `toml:"pid" json:"pid"`

	Status governance.GovernanceStatus `toml:"status" json:"status"`
	FSM    *fsm.FSM                    `json:"fsm"`
}

type Model11 struct {
	*Snake
}

//tc：更新超级治理管理员，更新失败
func (suite *Model11) Test001_UpdateSuperAdminIsFail() {
	node1, err := repo.Node1Path()
	suite.Require().Nil(err)
	pk, err := asym.RestorePrivateKey(node1, repo.KeyPassword)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	pid, err := suite.createPid()
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String(pid),
	}
	err = suite.InvokeRoleContract(UpdateAuditAdminNode, args...)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "the role is not a AuditAdmin")
}

//tc：冻结超级管理员，冻结失败
func (suite *Model11) Test002_FreezeSuperAdminIsFail() {
	node1, err := repo.Node1Path()
	suite.Require().Nil(err)
	pk, err := asym.RestorePrivateKey(node1, repo.KeyPassword)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(address.String()),
	}
	err = suite.InvokeRoleContract(FreezeRole, args...)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "super governance admin can not be freeze")
}

//tc：激活超级管理员，激活失败
func (suite *Model11) Test003_ActivateSuperAdminIsFail() {
	node1, err := repo.Node1Path()
	suite.Require().Nil(err)
	pk, err := asym.RestorePrivateKey(node1, repo.KeyPassword)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(address.String()),
	}
	err = suite.InvokeRoleContract(ActivateRole, args...)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "can not be activate")
}

//tc：注销超级管理员，注销失败
func (suite *Model11) Test004_LogoutSuperAdminIsFail() {
	node1, err := repo.Node1Path()
	suite.Require().Nil(err)
	pk, err := asym.RestorePrivateKey(node1, repo.KeyPassword)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(address.String()),
	}
	err = suite.InvokeRoleContract(LogoutRole, args...)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "super governance admin can not be logout")
}

//tc：注册普通治理管理员，注册成功
func (suite *Model11) Test005_RegisterAdminIsSuccess() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	pid, err := suite.createPid()
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("governanceAdmin"),
		rpcx.String(pid),
	}
	err = suite.InvokeRoleContract(RegisterRole, args...)
	suite.Require().Nil(err)
	role, err := suite.getRoleById(address.String())
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceAvailable, role.Status)

	//recover
	args = []*pb.Arg{
		rpcx.String(address.String()),
	}
	err = suite.InvokeRoleContract(LogoutRole, args...)
	suite.Require().Nil(err)
}

//tc：更新普通治理管理员，更新失败
func (suite Model11) Test006_UpdateAdminIsFail() {
	node3, err := repo.Node3Path()
	suite.Require().Nil(err)
	pk, err := asym.RestorePrivateKey(node3, repo.KeyPassword)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	pid, err := suite.createPid()
	suite.Require().Nil(err)

	args := []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String(pid),
	}
	err = suite.InvokeRoleContract(UpdateAuditAdminNode, args...)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "the role is not a AuditAdmin")
}

//tc：冻结普通治理管理员，冻结成功，发起投票失败，发起提案失败
func (suite Model11) Test007_FreezeAdminIsSuccess() {
	node3, err := repo.Node3Path()
	suite.Require().Nil(err)
	pk, err := asym.RestorePrivateKey(node3, repo.KeyPassword)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(address.String()),
	}
	err = suite.InvokeRoleContract(FreezeRole, args...)
	suite.Require().Nil(err)
	role, err := suite.getRoleById(address.String())
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceFrozen, role.Status)

	client := suite.NewClient(pk)
	pk, err = asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	bytes, err := pk.PublicKey().Bytes()
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(bytes)

	args = []*pb.Arg{
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
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "Register", &rpcx.TransactOpts{
		Nonce: atomic.AddUint64(&nonce3, 1),
	}, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)

	args = []*pb.Arg{
		rpcx.String(address.String()),
	}

	//recover
	err = suite.InvokeRoleContract(ActivateRole, args...)
	suite.Require().Nil(err)
	role, err = suite.getRoleById(address.String())
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceAvailable, role.Status)
}

//tc：冻结已冻结的普通治理管理员，冻结失败
func (suite Model11) Test008_FreezeFrozenAdminIsFail() {
	node3, err := repo.Node3Path()
	suite.Require().Nil(err)
	pk, err := asym.RestorePrivateKey(node3, repo.KeyPassword)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(address.String()),
	}
	err = suite.InvokeRoleContract(FreezeRole, args...)
	suite.Require().Nil(err)
	role, err := suite.getRoleById(address.String())
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
func (suite Model11) Test009_ActivateFrozenAdminIsSuccess() {
	node3, err := repo.Node3Path()
	suite.Require().Nil(err)
	pk, err := asym.RestorePrivateKey(node3, repo.KeyPassword)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(address.String()),
	}
	err = suite.InvokeRoleContract(FreezeRole, args...)
	suite.Require().Nil(err)
	role, err := suite.getRoleById(address.String())
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceFrozen, role.Status)

	err = suite.InvokeRoleContract(ActivateRole, args...)
	suite.Require().Nil(err)
	role, err = suite.getRoleById(address.String())
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceAvailable, role.Status)
}

//tc：激活未冻结的普通治理管理员，激活失败
func (suite Model11) Test010_ActivateAvailableAdminIsFail() {
	node3, err := repo.Node3Path()
	suite.Require().Nil(err)
	pk, err := asym.RestorePrivateKey(node3, repo.KeyPassword)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	role, err := suite.getRoleById(address.String())
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceAvailable, role.Status)

	args := []*pb.Arg{
		rpcx.String(address.String()),
	}
	err = suite.InvokeRoleContract(ActivateRole, args...)
	suite.Require().NotNil(err)
	role, err = suite.getRoleById(address.String())
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceAvailable, role.Status)
}

//tc：注销普通治理管理员，注销成功
func (suite Model11) Test011_LogoutAdminIsSuccess() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	pid, err := suite.createPid()
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("governanceAdmin"),
		rpcx.String(pid),
	}
	err = suite.InvokeRoleContract(RegisterRole, args...)
	suite.Require().Nil(err)
	role, err := suite.getRoleById(address.String())
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceAvailable, role.Status)

	args = []*pb.Arg{
		rpcx.String(address.String()),
	}
	err = suite.InvokeRoleContract(LogoutRole, args...)
	suite.Require().Nil(err)
}

//tc：注销冻结的普通治理管理员，注销成功
func (suite Model11) Test012_LogoutFrozenAdminIsSuccess() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	pid, err := suite.createPid()
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("governanceAdmin"),
		rpcx.String(pid),
	}
	err = suite.InvokeRoleContract(RegisterRole, args...)
	suite.Require().Nil(err)
	args = []*pb.Arg{
		rpcx.String(address.String()),
	}
	err = suite.InvokeRoleContract(FreezeRole, args...)
	suite.Require().Nil(err)
	role, err := suite.getRoleById(address.String())
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceFrozen, role.Status)

	args = []*pb.Arg{
		rpcx.String(address.String()),
	}
	err = suite.InvokeRoleContract(LogoutRole, args...)
	suite.Require().Nil(err)
}

//tc：注册审计管理员，注册成功
func (suite Model11) Test013_RegisterAuditAdminIsSuccess() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	pid, err := suite.createPid()
	suite.Require().Nil(err)
	err = suite.registerNode(pid, address.String(), "nvpNode")
	suite.Require().Nil(err)

	args := []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("auditAdmin"),
		rpcx.String(pid),
	}
	err = suite.InvokeRoleContract(RegisterRole, args...)
	suite.Require().Nil(err)
	role, err := suite.getRoleById(address.String())
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceAvailable, role.Status)

	//recover
	args = []*pb.Arg{
		rpcx.String(address.String()),
	}
	err = suite.InvokeRoleContract(LogoutRole, args...)
	suite.Require().Nil(err)
	err = suite.logoutNode(pid)
	suite.Require().Nil(err)
}

//tc：更新审计管理员，更新的id和原先一样，更新失败
func (suite Model11) Test013_UpdateAuditAdminWithSamePidIsFail() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	pid, err := suite.createPid()
	suite.Require().Nil(err)
	err = suite.registerNode(pid, address.String(), "nvpNode")
	suite.Require().Nil(err)

	args := []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("auditAdmin"),
		rpcx.String(pid),
	}
	err = suite.InvokeRoleContract(RegisterRole, args...)
	suite.Require().Nil(err)
	role, err := suite.getRoleById(address.String())
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceAvailable, role.Status)

	args = []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String(pid),
	}
	err = suite.InvokeRoleContract(UpdateAuditAdminNode, args...)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "the node ID is the same as before")

	//recover
	args = []*pb.Arg{
		rpcx.String(address.String()),
	}
	err = suite.InvokeRoleContract(LogoutRole, args...)
	suite.Require().Nil(err)
	err = suite.logoutNode(pid)
	suite.Require().Nil(err)
}

//tc：更新审计管理员，更新的id不存在，更新失败
func (suite Model11) Test013_UpdateAuditAdminWithNoExistPidIsFail() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	pid, err := suite.createPid()
	suite.Require().Nil(err)
	err = suite.registerNode(pid, address.String(), "nvpNode")
	suite.Require().Nil(err)

	args := []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("auditAdmin"),
		rpcx.String(pid),
	}
	err = suite.InvokeRoleContract(RegisterRole, args...)
	suite.Require().Nil(err)
	role, err := suite.getRoleById(address.String())
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceAvailable, role.Status)
	args = []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String(pid[0:len(pid)-1] + "1"),
	}
	err = suite.InvokeRoleContract(UpdateAuditAdminNode, args...)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), " this node does not exist")

	//recover
	args = []*pb.Arg{
		rpcx.String(address.String()),
	}
	err = suite.InvokeRoleContract(LogoutRole, args...)
	suite.Require().Nil(err)
	err = suite.logoutNode(pid)
	suite.Require().Nil(err)
}

//tc：更新审计管理员，更新id正确，更新成功
func (suite Model11) Test014_UpdateAuditAdminIsSuccess() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	pid, err := suite.createPid()
	suite.Require().Nil(err)
	err = suite.registerNode(pid, address.String(), "nvpNode")
	suite.Require().Nil(err)

	args := []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("auditAdmin"),
		rpcx.String(pid),
	}
	err = suite.InvokeRoleContract(RegisterRole, args...)
	suite.Require().Nil(err)
	role, err := suite.getRoleById(address.String())
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceAvailable, role.Status)
	newPid, err := suite.createPid()
	suite.Require().Nil(err)
	err = suite.registerNode(newPid, address.String(), "nvpNode")
	suite.Require().Nil(err)
	args = []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String(newPid),
	}
	err = suite.InvokeRoleContract(UpdateAuditAdminNode, args...)
	suite.Require().Nil(err)

	//recover
	args = []*pb.Arg{
		rpcx.String(address.String()),
	}
	err = suite.InvokeRoleContract(LogoutRole, args...)
	suite.Require().Nil(err)
	err = suite.logoutNode(pid)
	suite.Require().Nil(err)
}

//tc：冻结审计管理员，冻结成功（审计节点尚未实现）
func (suite Model11) Test015_FreezeAuditAdminIsSuccess() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	pid, err := suite.createPid()
	suite.Require().Nil(err)
	err = suite.registerNode(pid, address.String(), "nvpNode")
	suite.Require().Nil(err)

	args := []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("auditAdmin"),
		rpcx.String(pid),
	}
	err = suite.InvokeRoleContract(RegisterRole, args...)
	suite.Require().Nil(err)

	args = []*pb.Arg{
		rpcx.String(address.String()),
	}
	err = suite.InvokeRoleContract(FreezeRole, args...)
	suite.Require().Nil(err)
	role, err := suite.getRoleById(address.String())
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceFrozen, role.Status)

	//recover
	args = []*pb.Arg{
		rpcx.String(address.String()),
	}
	err = suite.InvokeRoleContract(LogoutRole, args...)
	suite.Require().Nil(err)
	err = suite.logoutNode(pid)
	suite.Require().Nil(err)
}

//tc：冻结已冻结的审计管理员，冻结失败
func (suite Model11) Test016_FreezeFrozenAuditAdminIsFail() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	pid, err := suite.createPid()
	suite.Require().Nil(err)
	err = suite.registerNode(pid, address.String(), "nvpNode")
	suite.Require().Nil(err)

	args := []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("auditAdmin"),
		rpcx.String(pid),
	}
	err = suite.InvokeRoleContract(RegisterRole, args...)
	suite.Require().Nil(err)

	args = []*pb.Arg{
		rpcx.String(address.String()),
	}
	err = suite.InvokeRoleContract(FreezeRole, args...)
	suite.Require().Nil(err)
	role, err := suite.getRoleById(address.String())
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceFrozen, role.Status)
	args = []*pb.Arg{
		rpcx.String(address.String()),
	}
	err = suite.InvokeRoleContract(FreezeRole, args...)
	suite.Require().NotNil(err)

	//recover
	args = []*pb.Arg{
		rpcx.String(address.String()),
	}
	err = suite.InvokeRoleContract(LogoutRole, args...)
	suite.Require().Nil(err)
	err = suite.logoutNode(pid)
	suite.Require().Nil(err)
}

//tc：激活冻结的审计管理员，激活成功
func (suite Model11) Test017_ActivateFrozenAuditAdminIsSuccess() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	pid, err := suite.createPid()
	suite.Require().Nil(err)
	err = suite.registerNode(pid, address.String(), "nvpNode")
	suite.Require().Nil(err)

	args := []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("auditAdmin"),
		rpcx.String(pid),
	}
	err = suite.InvokeRoleContract(RegisterRole, args...)
	suite.Require().Nil(err)

	args = []*pb.Arg{
		rpcx.String(address.String()),
	}
	err = suite.InvokeRoleContract(FreezeRole, args...)
	suite.Require().Nil(err)
	role, err := suite.getRoleById(address.String())
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceFrozen, role.Status)
	args = []*pb.Arg{
		rpcx.String(address.String()),
	}
	err = suite.InvokeRoleContract(ActivateRole, args...)
	suite.Require().Nil(err)

	//recover
	args = []*pb.Arg{
		rpcx.String(address.String()),
	}
	err = suite.InvokeRoleContract(LogoutRole, args...)
	suite.Require().Nil(err)
	err = suite.logoutNode(pid)
	suite.Require().Nil(err)
}

//tc：激活未冻结的审计管理员，激活失败
func (suite Model11) Test018_ActivateNoFrozenAuditAdminIsFail() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	pid, err := suite.createPid()
	suite.Require().Nil(err)
	err = suite.registerNode(pid, address.String(), "nvpNode")
	suite.Require().Nil(err)

	args := []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("auditAdmin"),
		rpcx.String(pid),
	}
	err = suite.InvokeRoleContract(RegisterRole, args...)
	suite.Require().Nil(err)

	role, err := suite.getRoleById(address.String())
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceAvailable, role.Status)
	args = []*pb.Arg{
		rpcx.String(address.String()),
	}
	err = suite.InvokeRoleContract(ActivateRole, args...)
	suite.Require().NotNil(err)

	//recover
	args = []*pb.Arg{
		rpcx.String(address.String()),
	}
	err = suite.InvokeRoleContract(LogoutRole, args...)
	suite.Require().Nil(err)
	err = suite.logoutNode(pid)
	suite.Require().Nil(err)
}

//tc：注销审计管理员，注销成功
func (suite Model11) Test019_LogoutAuditAdminIsSuccess() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	pid, err := suite.createPid()
	suite.Require().Nil(err)
	err = suite.registerNode(pid, address.String(), "nvpNode")
	suite.Require().Nil(err)

	args := []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("auditAdmin"),
		rpcx.String(pid),
	}
	err = suite.InvokeRoleContract(RegisterRole, args...)
	suite.Require().Nil(err)

	role, err := suite.getRoleById(address.String())
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceAvailable, role.Status)
	args = []*pb.Arg{
		rpcx.String(address.String()),
	}
	err = suite.InvokeRoleContract(LogoutRole, args...)
	suite.Require().Nil(err)
	err = suite.logoutNode(pid)
	suite.Require().Nil(err)
}

//tc：注销已冻结的审计管理员，冻结成功
func (suite Model11) Test020_LogoutFrozenAuditAdminIsSuccess() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	pid, err := suite.createPid()
	suite.Require().Nil(err)
	err = suite.registerNode(pid, address.String(), "nvpNode")
	suite.Require().Nil(err)

	args := []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("auditAdmin"),
		rpcx.String(pid),
	}
	err = suite.InvokeRoleContract(RegisterRole, args...)
	suite.Require().Nil(err)

	args = []*pb.Arg{
		rpcx.String(address.String()),
	}
	err = suite.InvokeRoleContract(FreezeRole, args...)
	suite.Require().Nil(err)
	role, err := suite.getRoleById(address.String())
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceFrozen, role.Status)

	args = []*pb.Arg{
		rpcx.String(address.String()),
	}
	err = suite.InvokeRoleContract(LogoutRole, args...)
	suite.Require().Nil(err)
	err = suite.logoutNode(pid)
	suite.Require().Nil(err)
}

//tc：审计管理员参与治理投票，投票失败
func (suite Model11) Test021_AuditAdminVoteIsFail() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	pid, err := suite.createPid()
	suite.Require().Nil(err)
	err = suite.registerNode(pid, address.String(), "nvpNode")
	suite.Require().Nil(err)

	args := []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("auditAdmin"),
		rpcx.String(pid),
	}
	err = suite.InvokeRoleContract(RegisterRole, args...)
	suite.Require().Nil(err)

	role, err := suite.getRoleById(address.String())
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceAvailable, role.Status)

	args = []*pb.Arg{
		rpcx.String(address.String()),
	}
	node2, err := repo.Node2Path()
	suite.Require().Nil(err)
	privateKey, err := asym.RestorePrivateKey(node2, repo.KeyPassword)
	suite.Require().Nil(err)
	client := suite.NewClient(privateKey)
	res, err := client.InvokeBVMContract(constant.RoleContractAddr.Address(), FreezeRole, &rpcx.TransactOpts{
		Nonce: atomic.AddUint64(&nonce2, 1),
	}, args...)
	suite.Require().Nil(err)
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	suite.Require().Nil(err)
	nonce, err := suite.client.GetPendingNonceByAccount(address.String())
	suite.Require().Nil(err)
	res, err = suite.vote(pk, nonce, pb.String(result.ProposalID), pb.String("approve"), pb.String("Appchain Pass"))
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "the administrator can not vote to the proposal")

	//recover
	args = []*pb.Arg{
		rpcx.String(address.String()),
	}
	err = suite.InvokeRoleContract(LogoutRole, args...)
	suite.Require().Nil(err)
	err = suite.logoutNode(pid)
	suite.Require().Nil(err)
}

//tc：调用GetRole接口，成功查询用户角色结构
func (suite Model11) Test022_GetRole() {
	node1, err := repo.Node1Path()
	suite.Require().Nil(err)
	privateKey, err := asym.RestorePrivateKey(node1, repo.KeyPassword)
	suite.Require().Nil(err)
	client := suite.NewClient(privateKey)
	res, err := client.InvokeBVMContract(constant.RoleContractAddr.Address(), "GetRole", &rpcx.TransactOpts{
		Nonce: atomic.AddUint64(&nonce1, 1),
	})
	suite.Require().Nil(err)
	suite.Require().Equal("governanceAdmin(super)", string(res.Ret))
}

//tc：调用GetRoleById接口，成功根据正确的iD查询用户角色结构
func (suite Model11) Test023_GetRoleById() {
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
	res, err := client.InvokeBVMContract(constant.RoleContractAddr.Address(), "GetRoleById", &rpcx.TransactOpts{
		Nonce: atomic.AddUint64(&nonce1, 1),
	}, args...)
	suite.Require().Nil(err)
	role := &Role{}
	err = json.Unmarshal(res.Ret, role)
	suite.Require().Nil(err)
	suite.Require().Equal("governanceAdmin", role.RoleType)
}

//tc：调用GetAdminRoles接口，成功获取所有治理管理员列表
func (suite Model11) Test024_GetAdminRoles() {
	node1, err := repo.Node1Path()
	suite.Require().Nil(err)
	privateKey, err := asym.RestorePrivateKey(node1, repo.KeyPassword)
	suite.Require().Nil(err)
	client := suite.NewClient(privateKey)
	res, err := client.InvokeBVMContract(constant.RoleContractAddr.Address(), "GetAdminRoles", &rpcx.TransactOpts{
		Nonce: atomic.AddUint64(&nonce1, 1),
	})
	suite.Require().Nil(err)
	var roles []Role
	err = json.Unmarshal(res.Ret, &roles)
	suite.Require().Nil(err)
	suite.Require().Greater(len(roles), 0)
}

//tc：调用GetAuditAdminRoles接口，成功获取所有审计管理员列表
func (suite Model11) Test025_GetAuditAdminRoles() {
	node1, err := repo.Node1Path()
	suite.Require().Nil(err)
	privateKey, err := asym.RestorePrivateKey(node1, repo.KeyPassword)
	suite.Require().Nil(err)
	client := suite.NewClient(privateKey)
	res, err := client.InvokeBVMContract(constant.RoleContractAddr.Address(), "GetAuditAdminRoles", &rpcx.TransactOpts{
		Nonce: atomic.AddUint64(&nonce1, 1),
	})
	suite.Require().Nil(err)
	var roles []Role
	err = json.Unmarshal(res.Ret, &roles)
	suite.Require().Nil(err)
	suite.Require().Greater(len(roles), 0)
}

//tc：调用IsAvailable接口，可用管理员返回true
func (suite Model11) Test026_IsAvailableWithTrue() {
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
	res, err := client.InvokeBVMContract(constant.RoleContractAddr.Address(), "IsAvailable", &rpcx.TransactOpts{
		Nonce: atomic.AddUint64(&nonce1, 1),
	}, args...)
	suite.Require().Nil(err)
	suite.Require().Equal("true", string(res.Ret))
}

//tc：调用IsAvailable接口，不可用管理员返回false
func (suite Model11) Test026_IsAvailableWithFalse() {
	node1, err := repo.Node1Path()
	suite.Require().Nil(err)
	privateKey, err := asym.RestorePrivateKey(node1, repo.KeyPassword)
	suite.Require().Nil(err)

	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	pid, err := suite.createPid()
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("governanceAdmin"),
		rpcx.String(pid),
	}
	err = suite.InvokeRoleContract(RegisterRole, args...)
	suite.Require().Nil(err)
	role, err := suite.getRoleById(address.String())
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceAvailable, role.Status)

	args = []*pb.Arg{
		rpcx.String(address.String()),
	}
	err = suite.InvokeRoleContract(LogoutRole, args...)
	suite.Require().Nil(err)

	client := suite.NewClient(privateKey)
	args = []*pb.Arg{
		rpcx.String(address.String()),
	}
	res, err := client.InvokeBVMContract(constant.RoleContractAddr.Address(), "IsAvailable", &rpcx.TransactOpts{
		Nonce: atomic.AddUint64(&nonce1, 1),
	}, args...)
	suite.Require().Nil(err)
	suite.Require().Equal("false", string(res.Ret))
}

//tc：调用IsSuperAdmin接口，超级治理管理员返回true
func (suite Model11) Test027_IsSuperAdminWithTrue() {
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
	res, err := client.InvokeBVMContract(constant.RoleContractAddr.Address(), "IsSuperAdmin", &rpcx.TransactOpts{
		Nonce: atomic.AddUint64(&nonce1, 1),
	}, args...)
	suite.Require().Nil(err)
	suite.Require().Equal("true", string(res.Ret))
}

//tc：调用IsSuperAdmin接口，非超级治理管理员返回false
func (suite Model11) Test028_IsSuperAdminWithFalse() {
	node2, err := repo.Node2Path()
	suite.Require().Nil(err)
	privateKey, err := asym.RestorePrivateKey(node2, repo.KeyPassword)
	suite.Require().Nil(err)
	address, err := privateKey.PublicKey().Address()
	suite.Require().Nil(err)
	client := suite.NewClient(privateKey)
	args := []*pb.Arg{
		rpcx.String(address.String()),
	}
	res, err := client.InvokeBVMContract(constant.RoleContractAddr.Address(), "IsSuperAdmin", &rpcx.TransactOpts{
		Nonce: atomic.AddUint64(&nonce2, 1),
	}, args...)
	suite.Require().Nil(err)
	suite.Require().Equal("false", string(res.Ret))
}

//tc：调用IsAdmin接口，治理管理员返回true
func (suite Model11) Test029_IsAdminWithTrue() {
	node2, err := repo.Node2Path()
	suite.Require().Nil(err)
	privateKey, err := asym.RestorePrivateKey(node2, repo.KeyPassword)
	suite.Require().Nil(err)
	address, err := privateKey.PublicKey().Address()
	suite.Require().Nil(err)
	client := suite.NewClient(privateKey)
	args := []*pb.Arg{
		rpcx.String(address.String()),
	}
	res, err := client.InvokeBVMContract(constant.RoleContractAddr.Address(), "IsAdmin", &rpcx.TransactOpts{
		Nonce: atomic.AddUint64(&nonce2, 1),
	}, args...)
	suite.Require().Nil(err)
	suite.Require().Equal("true", string(res.Ret))
}

//tc：调用IsAdmin接口，非治理管理员返回false
func (suite Model11) Test030_IsAdminWithFalse() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	pid, err := suite.createPid()
	suite.Require().Nil(err)
	err = suite.registerNode(pid, address.String(), "nvpNode")
	suite.Require().Nil(err)

	args := []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("auditAdmin"),
		rpcx.String(pid),
	}
	err = suite.InvokeRoleContract(RegisterRole, args...)
	suite.Require().Nil(err)

	node2, err := repo.Node2Path()
	suite.Require().Nil(err)
	privateKey, err := asym.RestorePrivateKey(node2, repo.KeyPassword)
	suite.Require().Nil(err)
	client := suite.NewClient(privateKey)
	args = []*pb.Arg{
		rpcx.String(address.String()),
	}
	res, err := client.InvokeBVMContract(constant.RoleContractAddr.Address(), "IsAdmin", &rpcx.TransactOpts{
		Nonce: atomic.AddUint64(&nonce2, 1),
	}, args...)
	suite.Require().Nil(err)
	suite.Require().Equal("false", string(res.Ret))
}

//tc：调用IsAuditAdmin接口，审计管理员返回true
func (suite Model11) Test031_IsAuditAdminWithTrue() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	pid, err := suite.createPid()
	suite.Require().Nil(err)
	err = suite.registerNode(pid, address.String(), "nvpNode")
	suite.Require().Nil(err)

	args := []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("auditAdmin"),
		rpcx.String(pid),
	}
	err = suite.InvokeRoleContract(RegisterRole, args...)
	suite.Require().Nil(err)

	node2, err := repo.Node2Path()
	suite.Require().Nil(err)
	privateKey, err := asym.RestorePrivateKey(node2, repo.KeyPassword)
	suite.Require().Nil(err)
	client := suite.NewClient(privateKey)
	args = []*pb.Arg{
		rpcx.String(address.String()),
	}
	res, err := client.InvokeBVMContract(constant.RoleContractAddr.Address(), "IsAuditAdmin", &rpcx.TransactOpts{
		Nonce: atomic.AddUint64(&nonce2, 1),
	}, args...)
	suite.Require().Nil(err)
	suite.Require().Equal("true", string(res.Ret))
}

//tc：调用IsAuditAdmin接口，非审计管理员返回false
func (suite Model11) Test032_IsAuditAdminWithFalse() {
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
	res, err := client.InvokeBVMContract(constant.RoleContractAddr.Address(), "IsAuditAdmin", &rpcx.TransactOpts{
		Nonce: atomic.AddUint64(&nonce1, 1),
	}, args...)
	suite.Require().Nil(err)
	suite.Require().Equal("false", string(res.Ret))
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

func (suite Model11) getRoleById(roleId string) (*Role, error) {
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
	res, err := client.InvokeBVMContract(constant.RoleContractAddr.Address(), "GetRoleById", &rpcx.TransactOpts{
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

func (suite Model11) createPid() (string, error) {
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

func (suite Model11) registerNode(pid, address, nodeType string) error {
	args := []*pb.Arg{
		rpcx.String(pid),
		rpcx.Uint64(0),
		rpcx.String(address),
		rpcx.String(nodeType),
	}
	err := suite.InvokeNodeContract("RegisterNode", args...)
	return err
}

func (suite Model11) logoutNode(pid string) error {
	args := []*pb.Arg{
		rpcx.String(pid),
	}
	err := suite.InvokeNodeContract("LogoutNode", args...)
	return err
}
