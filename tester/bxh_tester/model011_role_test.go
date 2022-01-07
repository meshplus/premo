package bxh_tester

import (
	"encoding/json"
	"errors"
	"sync/atomic"

	"github.com/meshplus/bitxhub-core/governance"

	crypto2 "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-kit/crypto/asym/ecdsa"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/meshplus/premo/internal/repo"
)

const (
	RegisterRole = "RegisterRole"
	FreezeRole   = "FreezeRole"
	ActivateRole = "ActivateRole"
	LogoutRole   = "LogoutRole"
)

type Model11 struct {
	*Snake
}

//TODO:update audit admin test
//tc：冻结超级管理员，冻结失败
func (suite *Model11) Test1101_FreezeSuperAdminIsFail() {
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
	suite.Require().Contains(err.Error(), "does not support freeze")
}

//tc：激活超级管理员，激活失败
func (suite *Model11) Test1102_ActivateSuperAdminIsFail() {
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
	suite.Require().Contains(err.Error(), "can not do activate")
}

//tc：注销超级管理员，注销失败
func (suite *Model11) Test1103_LogoutSuperAdminIsFail() {
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
	suite.Require().Contains(err.Error(), "does not support logout")
}

//tc：注册普通治理管理员，注册成功
func (suite *Model11) Test1104_RegisterAdminIsSuccess() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("governanceAdmin"),
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

//tc：冻结普通治理管理员，冻结成功，发起投票失败，发起提案失败
func (suite Model11) Test1105_FreezeAdminIsSuccess() {
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
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	from2, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	args = []*pb.Arg{
		rpcx.String(from2.String()),                               //chainID
		rpcx.String(from2.String()),                               //chainName
		rpcx.String("Flato V1.0.3"),                               //chainType
		rpcx.Bytes([]byte("")),                                    //trustRoot
		rpcx.String("0x857133c5C69e6Ce66F7AD46F200B9B3573e77582"), //broker
		rpcx.String("desc"),                                       //desc
		rpcx.String(HappyRuleAddr),                                //masterRuleAddr
		rpcx.String("https://github.com"),                         //masterRuleUrl
		rpcx.String(from2.String()),                               //adminAddrs
		rpcx.String("reason"),                                     //reason
	}
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "RegisterAppchain", nil, args...)
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
func (suite Model11) Test1106_FreezeFrozenAdminIsFail() {
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
func (suite Model11) Test1107_ActivateFrozenAdminIsSuccess() {
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
func (suite Model11) Test1108_ActivateAvailableAdminIsFail() {
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
func (suite Model11) Test1109_LogoutAdminIsSuccess() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("governanceAdmin"),
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
func (suite Model11) Test1110_LogoutFrozenAdminIsSuccess() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(address.String()),
		rpcx.String("governanceAdmin"),
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
func (suite Model11) Test1111_RegisterAuditAdminIsSuccess() {
	pk1, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	address1, err := pk1.PublicKey().Address()
	suite.Require().Nil(err)
	pid, err := suite.CreatePid()
	suite.Require().Nil(err)
	pk2, from2, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk2, from2, address2)
	suite.Require().Nil(err)
	err = suite.RegisterNode(pid, 0, address1.String(), "nvpNode", from2)
	suite.Require().Nil(err)
	pk3, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	address3, err := pk3.PublicKey().Address()
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(address3.String()),
		rpcx.String("auditAdmin"),
		rpcx.String(address1.String()),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(RegisterRole, args...)
	suite.Require().Nil(err)
	role, err := suite.GetRoleInfoById(address3.String())
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceAvailable, role.Status)
	//recover
	args = []*pb.Arg{
		rpcx.String(address3.String()),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(LogoutRole, args...)
	suite.Require().Nil(err)
	err = suite.LogoutNode(address1.String())
	suite.Require().Nil(err)
}

//tc：注销审计管理员，注销成功
func (suite Model11) Test1112_LogoutAuditAdminIsSuccess() {
	pk1, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	address1, err := pk1.PublicKey().Address()
	suite.Require().Nil(err)
	pid, err := suite.CreatePid()
	suite.Require().Nil(err)
	pk2, from2, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk2, from2, address2)
	suite.Require().Nil(err)
	err = suite.RegisterNode(pid, 0, address1.String(), "nvpNode", from2)
	suite.Require().Nil(err)
	pk3, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	address3, err := pk3.PublicKey().Address()
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(address3.String()),
		rpcx.String("auditAdmin"),
		rpcx.String(address1.String()),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(RegisterRole, args...)
	suite.Require().Nil(err)
	role, err := suite.GetRoleInfoById(address3.String())
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceAvailable, role.Status)
	args = []*pb.Arg{
		rpcx.String(address3.String()),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(LogoutRole, args...)
	suite.Require().Nil(err)
	err = suite.LogoutNode(address1.String())
	suite.Require().Nil(err)
}

//tc：审计管理员参与治理投票，投票失败
func (suite Model11) Test1113_AuditAdminVoteIsFail() {
	pk1, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	address1, err := pk1.PublicKey().Address()
	suite.Require().Nil(err)
	pid1, err := suite.CreatePid()
	suite.Require().Nil(err)
	pk2, from2, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk2, from2, address2)
	suite.Require().Nil(err)
	err = suite.RegisterNode(pid1, 0, address1.String(), "nvpNode", from2)
	suite.Require().Nil(err)
	pk3, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	address3, err := pk3.PublicKey().Address()
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(address3.String()),
		rpcx.String("auditAdmin"),
		rpcx.String(address1.String()),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(RegisterRole, args...)
	suite.Require().Nil(err)
	role, err := suite.GetRoleInfoById(address3.String())
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceAvailable, role.Status)
	pk4, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	address4, err := pk4.PublicKey().Address()
	suite.Require().Nil(err)
	pid2, err := suite.CreatePid()
	suite.Require().Nil(err)
	pk5, from5, address5, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk5, from5, address5)
	suite.Require().Nil(err)
	err = suite.RegisterNode(pid2, 0, address4.String(), "nvpNode", from5)
	suite.Require().Nil(err)
	pk6, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	address6, err := pk6.PublicKey().Address()
	suite.Require().Nil(err)
	args = []*pb.Arg{
		rpcx.String(address6.String()),
		rpcx.String("auditAdmin"),
		rpcx.String(address4.String()),
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
	nonce, err := client1.GetPendingNonceByAccount(address3.String())
	suite.Require().Nil(err)
	_, err = suite.vote(pk1, nonce, rpcx.String(result.ProposalID), rpcx.String("approve"), rpcx.String("reason"))
	suite.Require().NotNil(err)
	err = suite.VotePass(result.ProposalID)
	suite.Require().Nil(err)
	//recover
	args = []*pb.Arg{
		rpcx.String(address3.String()),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(LogoutRole, args...)
	suite.Require().Nil(err)
	err = suite.LogoutNode(address1.String())
	suite.Require().Nil(err)
	args = []*pb.Arg{
		rpcx.String(address6.String()),
		rpcx.String("reason"),
	}
	err = suite.InvokeRoleContract(LogoutRole, args...)
	suite.Require().Nil(err)
	err = suite.LogoutNode(address4.String())
	suite.Require().Nil(err)
}

//tc：调用GetRole接口，成功查询用户角色结构
func (suite Model11) Test1114_GetRole() {
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
func (suite Model11) Test1115_GetRoleByAddr() {
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
