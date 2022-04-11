package bxh_tester

import (
	"encoding/json"
	"fmt"
	"sync/atomic"

	"github.com/meshplus/bitxhub-core/governance"
	serviceMgr "github.com/meshplus/bitxhub-core/service-mgr"
	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/meshplus/premo/internal/repo"
)

type Model13 struct {
	*Snake
}

//tc：通过曾被占用的服务名称注册服务，服务注册成功
func (suite Model13) Test1301_RegisterServerWithHaveUsedNameIsSuccess() {
	pk1, chainID1, address1, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk1, chainID1, address1)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk1, chainID1, chainID1, chainID1, "CallContract")
	suite.Require().Nil(err)
	err = suite.UpdateService(pk1, chainID1+":"+chainID1, chainID1+"123")
	suite.Require().Nil(err)
	pk2, chainID2, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk2, chainID2, address2)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk2, chainID2, chainID2, chainID1, "CallContract")
	suite.Require().Nil(err)
}

//tc：应用链管理员注册服务，服务注册成功
func (suite Model13) Test1302_RegisterServerIsSuccess() {
	pk, chainID, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, chainID, chainID, "CallContract")
	suite.Require().Nil(err)
}

//tc：服务处于unavailable状态注册服务，服务注册成功
func (suite Model13) Test1303_RegisterServerWithUnavailableServerIsSuccess() {
	pk, chainID, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ServerToUnavailable(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, chainID, chainID, "CallContract")
	suite.Require().Nil(err)
}

//tc：应用链处于freezing状态注册服务，服务注册成功
func (suite Model13) Test1304_RegisterServerWithFreezingChainIsSuccess() {
	pk, chainID, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToFreezing(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, chainID, chainID, "CallContract")
	suite.Require().Nil(err)
}

//tc：应用链管理员更新服务，服务更新成功
func (suite Model13) Test1305_UpdateServerIsSuccess() {
	pk, chainID, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, chainID, chainID, "CallContract")
	suite.Require().Nil(err)
	err = suite.UpdateService(pk, chainID+":"+chainID, chainID+"123")
	suite.Require().Nil(err)
}

//tc：服务处于frozen状态更新服务，服务更新成功
func (suite Model13) Test1306_UpdateServerWithFrozenServerIsSuccess() {
	pk, chainID, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ServerToFrozen(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.UpdateService(pk, chainID+":"+chainID, chainID+"123")
	suite.Require().Nil(err)
	err = suite.CheckServerStatus(chainID+":"+chainID, governance.GovernanceAvailable)
	suite.Require().Nil(err)
}

//tc：应用链处于freezing状态更新服务，服务更新成功
func (suite Model13) Test1307_UpdateServerWithFreezingChainIsSuccess() {
	pk, chainID, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, chainID, chainID, "CallContract")
	suite.Require().Nil(err)
	node1Key, pubAddress, err := repo.Node1Priv()
	suite.Require().Nil(err)
	client := suite.NewClient(node1Key)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "FreezeAppchain",
		&rpcx.TransactOpts{
			From:  pubAddress.String(),
			Nonce: atomic.AddUint64(&nonce1, 1),
		},
		rpcx.String(chainID), rpcx.String("reason"),
	)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	err = suite.CheckChainStatus(chainID, governance.GovernanceFreezing)
	suite.Require().Nil(err)
	err = suite.UpdateService(pk, chainID+":"+chainID, chainID+"123")
	suite.Require().Nil(err)
}

//tc：中继链管理员冻结服务，服务冻结成功
func (suite Model13) Test1308_FreezeServerWithRelayAdminIsSuccess() {
	pk, chainID, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, chainID, chainID, "CallContract")
	suite.Require().Nil(err)
	err = suite.FreezeService(chainID + ":" + chainID)
	suite.Require().Nil(err)
	err = suite.CheckServerStatus(chainID+":"+chainID, governance.GovernanceFrozen)
	suite.Require().Nil(err)
}

//tc：应用链处于freezing状态冻结服务，服务冻结成功
func (suite Model13) Test1309_FreezeServerWithFreezingChainIsSuccess() {
	pk, chainID, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, chainID, chainID, "CallContract")
	suite.Require().Nil(err)
	node1Key, pubAddress, err := repo.Node1Priv()
	suite.Require().Nil(err)
	client := suite.NewClient(node1Key)
	nonce := atomic.AddUint64(&nonce1, 1)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "FreezeAppchain",
		&rpcx.TransactOpts{
			From:  pubAddress.String(),
			Nonce: nonce,
		},
		rpcx.String(chainID), rpcx.String("reason"),
	)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	err = suite.CheckChainStatus(chainID, governance.GovernanceFreezing)
	suite.Require().Nil(err)
	err = suite.FreezeService(chainID + ":" + chainID)
	suite.Require().Nil(err)
}

//tc：中继链管理员激活服务，服务激活成功
func (suite Model13) Test1310_ActivateServerWithRelayAdminIsSuccess() {
	pk, chainID, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ServerToFrozen(pk, chainID, address)
	suite.Require().Nil(err)
	node1Key, from, err := repo.Node1Priv()
	suite.Require().Nil(err)
	client := suite.NewClient(node1Key)
	res, err := client.InvokeBVMContract(constant.ServiceMgrContractAddr.Address(), "ActivateService", &rpcx.TransactOpts{
		From:  from.String(),
		Nonce: atomic.AddUint64(&nonce1, 1),
	}, rpcx.String(chainID+":"+chainID), rpcx.String("reason"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	suite.Require().Nil(err)
	err = suite.VotePass(result.ProposalID)
	suite.Require().Nil(err)
	err = suite.CheckServerStatus(chainID+":"+chainID, governance.GovernanceAvailable)
	suite.Require().Nil(err)
}

//tc：应用链管理员激活服务，服务激活成功
func (suite Model13) Test1311_ActivateServerIsSuccess() {
	pk, chainID, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, chainID, chainID, "CallContract")
	suite.Require().Nil(err)
	err = suite.FreezeService(chainID + ":" + chainID)
	suite.Require().Nil(err)
	err = suite.CheckServerStatus(chainID+":"+chainID, governance.GovernanceFrozen)
	suite.Require().Nil(err)
	err = suite.ActivateService(pk, chainID+":"+chainID)
	suite.Require().Nil(err)
}

//tc：应用链处于freezing状态激活服务，服务激活成功
func (suite Model13) Test1312_ActivateServerWithFreezingChainIsSuccess() {
	pk, chainID, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, chainID, chainID, "CallContract")
	suite.Require().Nil(err)
	err = suite.FreezeService(chainID + ":" + chainID)
	suite.Require().Nil(err)
	node1Key, pubAddress, err := repo.Node1Priv()
	suite.Require().Nil(err)
	client := suite.NewClient(node1Key)
	nonce := atomic.AddUint64(&nonce1, 1)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "FreezeAppchain",
		&rpcx.TransactOpts{
			From:  pubAddress.String(),
			Nonce: nonce,
		},
		rpcx.String(chainID), rpcx.String("reason"),
	)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	err = suite.CheckChainStatus(chainID, governance.GovernanceFreezing)
	suite.Require().Nil(err)
	err = suite.ActivateService(pk, chainID+":"+chainID)
	suite.Require().Nil(err)
	err = suite.CheckServerStatus(chainID+":"+chainID, governance.GovernanceAvailable)
	suite.Require().Nil(err)
}

//tc：应用链管理员注销服务，服务注销成功
func (suite Model13) Test1313_LogoutServerWithRelayAdminIsSuccess() {
	pk, chainID, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, chainID, chainID, "CallContract")
	suite.Require().Nil(err)
	err = suite.LogoutService(pk, chainID+":"+chainID)
	suite.Require().Nil(err)
}

//tc：服务处于updating状态注销服务，服务注销成功
func (suite Model13) Test1314_LogoutServerWithUpdatingServerIsSuccess() {
	pk, chainID, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ServerToUpdating(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.LogoutService(pk, chainID+":"+chainID)
	suite.Require().Nil(err)
}

//tc：服务处于activating状态注销服务，服务注销成功
func (suite Model13) Test1315_LogoutServerWithActivatingServerIsSuccess() {
	pk, chainID, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ServerToActivating(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.LogoutService(pk, chainID+":"+chainID)
	suite.Require().Nil(err)
}

//tc：服务处于freezing状态注销服务，服务注销成功
func (suite Model13) Test1316_LogoutServerWithFreezingServerIsSuccess() {
	pk, chainID, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ServerToFreezing(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.LogoutService(pk, chainID+":"+chainID)
	suite.Require().Nil(err)
}

//tc：服务处于pause状态注销服务，服务注销成功
func (suite Model13) Test1317_LogoutServerWithPauseServerIsSuccess() {
	pk, chainID, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ServerToPause(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.LogoutService(pk, chainID+":"+chainID)
	suite.Require().Nil(err)
}

//tc：应用链处于activating状态注销服务，服务注销成功
func (suite Model13) Test1318_LogoutServerWithActivatingChainIsSuccess() {
	pk, chainID, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, chainID, chainID, "CallContract")
	suite.Require().Nil(err)
	err = suite.FreezeAppchain(chainID)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "ActivateAppchain", nil, rpcx.String(chainID), rpcx.String("reason"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	err = suite.CheckChainStatus(chainID, governance.GovernanceActivating)
	suite.Require().Nil(err)
	err = suite.LogoutService(pk, chainID+":"+chainID)
	suite.Require().Nil(err)
}

//tc：应用链处于freezing状态注销服务，服务注销成功
func (suite Model13) Test1319_LogoutServerWithFreezingChainIsSuccess() {
	pk, chainID, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, chainID, chainID, "CallContract")
	suite.Require().Nil(err)
	node1Key, pubAddress, err := repo.Node1Priv()
	suite.Require().Nil(err)
	client := suite.NewClient(node1Key)
	nonce := atomic.AddUint64(&nonce1, 1)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "FreezeAppchain",
		&rpcx.TransactOpts{
			From:  pubAddress.String(),
			Nonce: nonce,
		},
		rpcx.String(chainID), rpcx.String("reason"),
	)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	err = suite.CheckChainStatus(chainID, governance.GovernanceFreezing)
	suite.Require().Nil(err)
	err = suite.LogoutService(pk, chainID+":"+chainID)
	suite.Require().Nil(err)
}

//tc：应用链处于frozen状态注销服务，服务注销成功
func (suite Model13) Test1320_LogoutServerWithFrozenChainIsSuccess() {
	pk, chainID, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, chainID, chainID, "CallContract")
	suite.Require().Nil(err)
	err = suite.FreezeAppchain(chainID)
	suite.Require().Nil(err)
	err = suite.CheckChainStatus(chainID, governance.GovernanceFrozen)
	suite.Require().Nil(err)
	err = suite.LogoutService(pk, chainID+":"+chainID)
	suite.Require().Nil(err)
}

//tc：应用链处于logouting状态注销服务，服务注销成功
func (suite Model13) Test1321_LogoutServerWithLogoutingChainIsSuccess() {
	pk, chainID, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, chainID, chainID, "CallContract")
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "LogoutAppchain", nil, rpcx.String(chainID), rpcx.String("reason"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	err = suite.CheckChainStatus(chainID, governance.GovernanceLogouting)
	suite.Require().Nil(err)
	err = suite.LogoutService(pk, chainID+":"+chainID)
	suite.Require().Nil(err)
}

// RegisterServer register server
func (suite Snake) RegisterServer(pk crypto.PrivateKey, chainID, serviceID, name, typ string) error {
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(chainID),
		rpcx.String(serviceID),
		rpcx.String(name),
		rpcx.String(typ),
		rpcx.String("test"),
		rpcx.Uint64(1),
		rpcx.String(""),
		rpcx.String("test"),
		rpcx.String("reason"),
	}
	res, err := client.InvokeBVMContract(constant.ServiceMgrContractAddr.Address(), "RegisterService", nil, args...)
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf(string(res.Ret))
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

// FreezeService freeze server
func (suite Snake) FreezeService(chainServiceID string) error {
	node1Key, from, err := repo.Node1Priv()
	if err != nil {
		return err
	}
	client := suite.NewClient(node1Key)
	res, err := client.InvokeBVMContract(constant.ServiceMgrContractAddr.Address(), "FreezeService", &rpcx.TransactOpts{
		From:  from.String(),
		Nonce: atomic.AddUint64(&nonce1, 1),
	}, rpcx.String(chainServiceID), rpcx.String("reason"))
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf(string(res.Ret))
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

// ActivateService activate server
func (suite Model13) ActivateService(pk crypto.PrivateKey, chainServiceID string) error {
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ServiceMgrContractAddr.Address(), "ActivateService", nil, rpcx.String(chainServiceID), rpcx.String("reason"))
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf(string(res.Ret))
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

// UpdateService update server
func (suite Snake) UpdateService(pk crypto.PrivateKey, chainServiceID, name string) error {
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(chainServiceID),
		rpcx.String(name),
		rpcx.String("test"),
		rpcx.String(""),
		rpcx.String("test"),
		rpcx.String("reason"),
	}
	res, err := client.InvokeBVMContract(constant.ServiceMgrContractAddr.Address(), "UpdateService", nil, args...)
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf(string(res.Ret))
	}
	if string(res.Ret) == "" {
		return nil
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

// LogoutService logout server
func (suite Model13) LogoutService(pk crypto.PrivateKey, chainServiceID string) error {
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ServiceMgrContractAddr.Address(), "LogoutService", nil, rpcx.String(chainServiceID), rpcx.String("reason"))
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf(string(res.Ret))
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

// CheckServerStatus check server status
func (suite Snake) CheckServerStatus(serverID string, status governance.GovernanceStatus) error {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	if err != nil {
		return err
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ServiceMgrContractAddr.Address(), "GetServiceInfo", nil, rpcx.String(serverID))
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf(string(res.Ret))
	}
	server := &serviceMgr.Service{}
	err = json.Unmarshal(res.Ret, server)
	if err != nil {
		return err
	}
	if server.Status != status {
		return fmt.Errorf("expect status is %s, but got %s", status, server.Status)
	}
	return nil
}

// ServerToRegisting get a registing server
func (suite Model13) ServerToRegisting(pk crypto.PrivateKey, chainID, address string) error {
	err := suite.RegisterAppchain(pk, chainID, address)
	if err != nil {
		return err
	}
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(chainID),
		rpcx.String(chainID),
		rpcx.String(chainID),
		rpcx.String("CallContract"),
		rpcx.String("test"),
		rpcx.Uint64(1),
		rpcx.String(""),
		rpcx.String("test"),
		rpcx.String("reason"),
	}
	res, err := client.InvokeBVMContract(constant.ServiceMgrContractAddr.Address(), "RegisterService", nil, args...)
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf(string(res.Ret))
	}
	err = suite.CheckServerStatus(chainID+":"+chainID, governance.GovernanceRegisting)
	if err != nil {
		return err
	}
	return nil
}

// ServerToUnavailable get a unavailable server
func (suite Model13) ServerToUnavailable(pk crypto.PrivateKey, chainID, address string) error {
	err := suite.RegisterAppchain(pk, chainID, address)
	if err != nil {
		return err
	}
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(chainID),
		rpcx.String(chainID),
		rpcx.String(chainID),
		rpcx.String("CallContract"),
		rpcx.String("test"),
		rpcx.Uint64(1),
		rpcx.String(""),
		rpcx.String("test"),
		rpcx.String("reason"),
	}
	res, err := client.InvokeBVMContract(constant.ServiceMgrContractAddr.Address(), "RegisterService", nil, args...)
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf(string(res.Ret))
	}
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	if err != nil {
		return err
	}
	err = suite.VoteReject(result.ProposalID)
	if err != nil {
		return err
	}
	err = suite.CheckServerStatus(string(result.Extra), governance.GovernanceUnavailable)
	if err != nil {
		return err
	}
	return nil
}

// ServerToUpdating get an updating server
func (suite Model13) ServerToUpdating(pk crypto.PrivateKey, chainID, address string) error {
	err := suite.RegisterAppchain(pk, chainID, address)
	if err != nil {
		return err
	}
	err = suite.RegisterServer(pk, chainID, chainID, chainID, "CallContract")
	if err != nil {
		return err
	}
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(chainID + ":" + chainID),
		rpcx.String(chainID + "123"),
		rpcx.String("test"),
		rpcx.String(""),
		rpcx.String("test"),
		rpcx.String("reason"),
	}
	res, err := client.InvokeBVMContract(constant.ServiceMgrContractAddr.Address(), "UpdateService", nil, args...)
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf(string(res.Ret))
	}
	err = suite.CheckServerStatus(chainID+":"+chainID, governance.GovernanceUpdating)
	if err != nil {
		return err
	}
	return nil
}

// ServerToActivating get an activating server
func (suite Model13) ServerToActivating(pk crypto.PrivateKey, chainID, address string) error {
	err := suite.RegisterAppchain(pk, chainID, address)
	if err != nil {
		return err
	}
	err = suite.RegisterServer(pk, chainID, chainID, chainID, "CallContract")
	if err != nil {
		return err
	}
	err = suite.FreezeService(chainID + ":" + chainID)
	if err != nil {
		return err
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ServiceMgrContractAddr.Address(), "ActivateService", nil, rpcx.String(chainID+":"+chainID), rpcx.String("reason"))
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf(string(res.Ret))
	}
	err = suite.CheckServerStatus(chainID+":"+chainID, governance.GovernanceActivating)
	if err != nil {
		return err
	}
	return nil
}

// ServerToFreezing get a freezing server
func (suite Model13) ServerToFreezing(pk crypto.PrivateKey, chainID, address string) error {
	err := suite.RegisterAppchain(pk, chainID, address)
	if err != nil {
		return err
	}
	err = suite.RegisterServer(pk, chainID, chainID, chainID, "CallContract")
	if err != nil {
		return err
	}
	node1Key, from, err := repo.Node1Priv()
	if err != nil {
		return err
	}
	client := suite.NewClient(node1Key)
	res, err := client.InvokeBVMContract(constant.ServiceMgrContractAddr.Address(), "FreezeService", &rpcx.TransactOpts{
		From:  from.String(),
		Nonce: atomic.AddUint64(&nonce1, 1),
	}, rpcx.String(chainID+":"+chainID), rpcx.String("reason"))
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf(string(res.Ret))
	}
	err = suite.CheckServerStatus(chainID+":"+chainID, governance.GovernanceFreezing)
	if err != nil {
		return err
	}
	return nil
}

// ServerToFrozen get a frozen server
func (suite Model13) ServerToFrozen(pk crypto.PrivateKey, chainID, address string) error {
	err := suite.RegisterAppchain(pk, chainID, address)
	if err != nil {
		return err
	}
	err = suite.RegisterServer(pk, chainID, chainID, chainID, "CallContract")
	if err != nil {
		return err
	}
	err = suite.FreezeService(chainID + ":" + chainID)
	if err != nil {
		return err
	}
	err = suite.CheckServerStatus(chainID+":"+chainID, governance.GovernanceFrozen)
	if err != nil {
		return err
	}
	return nil
}

// ServerToLogouting get a logouting server
func (suite Model13) ServerToLogouting(pk crypto.PrivateKey, chainID, address string) error {
	err := suite.RegisterAppchain(pk, chainID, address)
	if err != nil {
		return err
	}
	err = suite.RegisterServer(pk, chainID, chainID, chainID, "CallContract")
	if err != nil {
		return err
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ServiceMgrContractAddr.Address(), "LogoutService", nil, rpcx.String(chainID+":"+chainID), rpcx.String("reason"))
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf(string(res.Ret))
	}
	err = suite.CheckServerStatus(chainID+":"+chainID, governance.GovernanceLogouting)
	if err != nil {
		return err
	}
	return nil
}

// ServerToForbidden get a forbidden server
func (suite Model13) ServerToForbidden(pk crypto.PrivateKey, chainID, address string) error {
	err := suite.RegisterAppchain(pk, chainID, address)
	if err != nil {
		return err
	}
	err = suite.RegisterServer(pk, chainID, chainID, chainID, "CallContract")
	if err != nil {
		return err
	}
	err = suite.LogoutService(pk, chainID+":"+chainID)
	if err != nil {
		return err
	}
	err = suite.CheckServerStatus(chainID+":"+chainID, governance.GovernanceForbidden)
	if err != nil {
		return err
	}
	return nil
}

// ServerToPause get a pause server
func (suite Model13) ServerToPause(pk crypto.PrivateKey, chainID, address string) error {
	err := suite.RegisterAppchain(pk, chainID, address)
	if err != nil {
		return err
	}
	err = suite.RegisterServer(pk, chainID, chainID, chainID, "CallContract")
	if err != nil {
		return err
	}
	err = suite.FreezeAppchain(chainID)
	if err != nil {
		return err
	}
	err = suite.CheckServerStatus(chainID+":"+chainID, governance.GovernancePause)
	if err != nil {
		return err
	}
	return nil
}
