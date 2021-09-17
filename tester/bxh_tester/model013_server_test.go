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
	"github.com/pkg/errors"
)

type Model13 struct {
	*Snake
}

func (suite *Model13) SetupTest() {
	suite.T().Parallel()
}

//tc：中继链管理员注册服务，服务注册失败
func (suite Model13) Test1301_RegisterServerWithRelayAdminIsFail() {
	path, err := repo.Node1Path()
	suite.Require().Nil(err)
	pk1, err := asym.RestorePrivateKey(path, repo.KeyPassword)
	suite.Require().Nil(err)
	from, err := pk1.PublicKey().Address()
	suite.Require().Nil(err)
	pk2, chainID, address, err := suite.RegisterRule()
	err = suite.RegisterAppchain(pk2, chainID, address)
	suite.Require().Nil(err)
	client := suite.NewClient(pk1)
	args := []*pb.Arg{
		rpcx.String(chainID),
		rpcx.String(suite.GetServerID(pk2)),
		rpcx.String("testServer"),
		rpcx.String("CallContract"),
		rpcx.String("test"),
		rpcx.Bool(true),
		rpcx.String(""),
		rpcx.String("test"),
		rpcx.String("reason"),
	}
	res, err := client.InvokeBVMContract(constant.ServiceMgrContractAddr.Address(), "RegisterService", &rpcx.TransactOpts{
		From:  from.String(),
		Nonce: atomic.AddUint64(&nonce1, 1),
	}, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
}

//tc：非应用链管理员注册服务，服务注册失败
func (suite Model13) Test1302_RegisterServerWithNoAdminIsFail() {
	pk1, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	pk2, chainID, address, err := suite.RegisterRule()
	err = suite.RegisterAppchain(pk2, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk1, chainID, suite.GetServerID(pk2), "CallContract")
	suite.Require().NotNil(err)
}

//tc：应用链管理员注册服务，服务注册成功
func (suite Model13) Test1303_RegisterServerIsSuccess() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	suite.Require().Nil(err)
}

//tc：应用链处于未注册状态注册服务，服务注册成功
func (suite Model13) Test1304_RegisterServerWithNoRegisterAdminIsFail() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, "test", suite.GetServerID(pk), "CallContract")
	suite.Require().NotNil(err)
}

//tc：服务处于registing状态注册服务，服务注册失败
func (suite Model13) Test1305_RegisterServerWithRegistingServerIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ServerToRegisting(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	suite.Require().NotNil(err)
}

//tc：服务处于available状态注册服务，服务注册失败
func (suite Model13) Test1306_RegisterServerWithAvailableServerIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	suite.Require().Nil(err)
	err = suite.CheckServerStatus(chainID+":"+suite.GetServerID(pk), governance.GovernanceAvailable)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	suite.Require().NotNil(err)
}

//tc：服务处于unavailable状态注册服务，服务注册成功
func (suite Model13) Test1307_RegisterServerWithUnavailableServerIsSuccess() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ServerToUnavailable(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	suite.Require().Nil(err)
}

//tc：服务处于updating状态注册服务，服务注册失败
func (suite Model13) Test1308_RegisterServerWithUpdatingServerIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ServerToUpdating(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	suite.Require().NotNil(err)
}

//tc：服务处于activating状态注册服务，服务注册失败
func (suite Model13) Test1309_RegisterServerWithActivatingServerIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ServerToActivating(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	suite.Require().NotNil(err)
}

//tc：服务处于freezing状态注册服务，服务注册失败
func (suite Model13) Test1310_RegisterServerWithFreezingServerIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ServerToFreezing(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	suite.Require().NotNil(err)
}

//tc：服务处于frozen状态注册服务，服务注册失败
func (suite Model13) Test1311_RegisterServerWithFrozenServerIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ServerToFrozen(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	suite.Require().NotNil(err)
}

//tc：服务处于logouting状态注册服务，服务注册失败
func (suite Model13) Test1312_RegisterServerWithLogoutingServerIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ServerToLogouting(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	suite.Require().NotNil(err)
}

//tc：服务处于forbidden状态注册服务，服务注册失败
func (suite Model13) Test1313_RegisterServerWithForbiddenServerIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ServerToForbidden(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	suite.Require().NotNil(err)
}

//tc：应用链处于unavailable状态注册服务，服务注册失败
func (suite Model13) Test1314_RegisterServerWithUnavailableChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ServerToUnavailable(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	suite.Require().Nil(err)
}

//tc：应用链处于activating状态注册服务，服务注册失败
func (suite Model13) Test1315_RegisterServerWithActivatingChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ServerToActivating(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	suite.Require().NotNil(err)
}

//tc：应用链处于freezing状态注册服务，服务注册成功
func (suite Model13) Test1316_RegisterServerWithFreezingChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ChainToFreezing(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	suite.Require().Nil(err)
}

//tc：应用链处于frozen状态注册服务，服务注册失败
func (suite Model13) Test1317_RegisterServerWithFrozenChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ChainToFrozen(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	suite.Require().NotNil(err)
}

//tc：应用链处于logouting状态注册服务，服务注册失败
func (suite Model13) Test1318_RegisterServerWithLogoutingChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ChainToLogouting(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	suite.Require().NotNil(err)
}

//tc：应用链处于forbidden状态注册服务，服务注册失败
func (suite Model13) Test1319_RegisterServerWithForbiddenChainsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ChainToForbidden(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	suite.Require().NotNil(err)
}

//tc：服务类型异常注册服务，服务注册失败
func (suite Model13) Test1320_RegisterServerWithErrorTypeIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract111")
	suite.Require().NotNil(err)
}

//tc：中继链管理员更新服务，服务更新失败
func (suite Model13) Test1321_UpdateServerWithRelayAdminIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	suite.Require().Nil(err)
	path, err := repo.Node1Path()
	suite.Require().Nil(err)
	node1pk, err := asym.RestorePrivateKey(path, repo.KeyPassword)
	suite.Require().Nil(err)
	from, err := node1pk.PublicKey().Address()
	suite.Require().Nil(err)
	client := suite.NewClient(node1pk)
	args := []*pb.Arg{
		rpcx.String(chainID + ":" + suite.GetServerID(pk)),
		rpcx.String("test111"),
		rpcx.String("test"),
		rpcx.Bool(true),
		rpcx.String(""),
		rpcx.String("test"),
		rpcx.String("reason"),
	}
	res, err := client.InvokeBVMContract(constant.ServiceMgrContractAddr.Address(), "UpdateService", &rpcx.TransactOpts{
		From:  from.String(),
		Nonce: atomic.AddUint64(&nonce1, 1),
	}, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
}

//tc：应用链管理员更新服务，服务更新成功
func (suite Model13) Test1322_UpdateServerIsSuccess() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	suite.Require().Nil(err)
	err = suite.UpdateService(pk, chainID+":"+suite.GetServerID(pk), "test111")
	suite.Require().Nil(err)
}

//tc：非应用链管理员更新服务，服务更新失败
func (suite Model13) Test1323_UpdateServerWithNoAdminIsFail() {
	pk1, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk1, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk1, chainID, suite.GetServerID(pk1), "CallContract")
	suite.Require().Nil(err)
	pk2, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.UpdateService(pk2, chainID+":"+suite.GetServerID(pk1), "test111")
	suite.Require().NotNil(err)
}

//tc：服务处于registing状态更新服务，服务更新失败
func (suite Model13) Test1324_UpdateServerWithRegistingServerIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ServerToRegisting(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.UpdateService(pk, chainID+":"+suite.GetServerID(pk), "test111")
	suite.Require().NotNil(err)
}

//tc：服务处于unavailable状态更新服务，服务更新成功
func (suite Model13) Test1325_UpdateServerWithUnavailableServerIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ServerToUnavailable(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.UpdateService(pk, chainID+":"+suite.GetServerID(pk), "test111")
	suite.Require().NotNil(err)
}

//tc：服务处于updating状态更新服务，服务更新失败
func (suite Model13) Test1326_UpdateServerWithUpdatingServerIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ServerToUpdating(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.UpdateService(pk, chainID+":"+suite.GetServerID(pk), "test111")
	suite.Require().NotNil(err)
}

//tc：服务处于activating状态更新服务，服务更新失败
func (suite Model13) Test1327_UpdateServerWithActivatingServerIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ServerToActivating(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.UpdateService(pk, chainID+":"+suite.GetServerID(pk), "test111")
	suite.Require().NotNil(err)
}

//tc：服务处于freezing状态更新服务，服务更新失败
func (suite Model13) Test1328_UpdateServerWithFreezingServerIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ServerToFreezing(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.UpdateService(pk, chainID+":"+suite.GetServerID(pk), "test111")
	suite.Require().NotNil(err)
}

//tc：服务处于frozen状态更新服务，服务更新成功
func (suite Model13) Test1329_UpdateServerWithFrozenServerIsSuccess() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ServerToFrozen(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.UpdateService(pk, chainID+":"+suite.GetServerID(pk), "test111")
	suite.Require().Nil(err)
	err = suite.CheckServerStatus(chainID+":"+suite.GetServerID(pk), governance.GovernanceAvailable)
	suite.Require().Nil(err)
}

//tc：服务处于logouting状态更新服务，服务更新失败
func (suite Model13) Test1330_UpdateServerWithLogoutingServerIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ServerToLogouting(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.UpdateService(pk, chainID+":"+suite.GetServerID(pk), "test111")
	suite.Require().NotNil(err)
}

//tc：服务处于forbidden状态更新服务，服务更新失败
func (suite Model13) Test1331_UpdateServerWithForbiddenServerIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ServerToForbidden(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.UpdateService(pk, chainID+":"+suite.GetServerID(pk), "test111")
	suite.Require().NotNil(err)
}

//tc：应用链处于activating状态更新服务，服务更新失败
func (suite Model13) Test1332_UpdateServerWithActivatingChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	suite.Require().Nil(err)
	err = suite.FreezeAppchain(chainID)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "ActivateAppchain", nil, rpcx.String(chainID), rpcx.String("reason"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	err = suite.CheckChainStatus(chainID, governance.GovernanceActivating)
	suite.Require().Nil(err)
	err = suite.UpdateService(pk, chainID+":"+suite.GetServerID(pk), "test111")
	suite.Require().NotNil(err)
}

//tc：应用链处于freezing状态更新服务，服务更新成功
func (suite Model13) Test1333_UpdateServerWithFreezingChainIsSuccess() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	suite.Require().Nil(err)
	path, err := repo.Node1Path()
	suite.Require().Nil(err)
	node1Key, err := asym.RestorePrivateKey(path, repo.KeyPassword)
	suite.Require().Nil(err)
	client := suite.NewClient(node1Key)
	pubAddress, err := node1Key.PublicKey().Address()
	suite.Require().Nil(err)
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
	err = suite.UpdateService(pk, chainID+":"+suite.GetServerID(pk), "test111")
	suite.Require().Nil(err)
}

//tc：应用链处于frozen状态更新服务，服务更新失败
func (suite Model13) Test1334_UpdateServerWithFrozenChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	suite.Require().Nil(err)
	err = suite.FreezeAppchain(chainID)
	suite.Require().Nil(err)
	err = suite.CheckChainStatus(chainID, governance.GovernanceFrozen)
	suite.Require().Nil(err)
	err = suite.UpdateService(pk, chainID+":"+suite.GetServerID(pk), "test111")
	suite.Require().NotNil(err)
}

//tc：应用链处于logouting状态更新服务，服务更新失败
func (suite Model13) Test1335_UpdateServerWithLogoutingChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "LogoutAppchain", nil, rpcx.String(chainID), rpcx.String("reason"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	err = suite.CheckChainStatus(chainID, governance.GovernanceLogouting)
	suite.Require().Nil(err)
	err = suite.UpdateService(pk, chainID+":"+suite.GetServerID(pk), "test111")
	suite.Require().NotNil(err)
}

//tc：应用链处于forbidden状态更新服务，服务更新失败
func (suite Model13) Test1336_UpdateServerWithForbiddenChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	suite.Require().Nil(err)
	err = suite.LogoutAppchain(pk, chainID)
	suite.Require().Nil(err)
	err = suite.CheckChainStatus(chainID, governance.GovernanceForbidden)
	suite.Require().Nil(err)
	err = suite.UpdateService(pk, chainID+":"+suite.GetServerID(pk), "test111")
	suite.Require().NotNil(err)
}

//tc：中继链管理员冻结服务，服务冻结成功
func (suite Model13) Test1337_FreezeServerWithRelayAdminIsSuccess() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	suite.Require().Nil(err)
	err = suite.FreezeService(chainID + ":" + suite.GetServerID(pk))
	suite.Require().Nil(err)
	err = suite.CheckServerStatus(chainID+":"+suite.GetServerID(pk), governance.GovernanceFrozen)
	suite.Require().Nil(err)
}

//tc：应用链管理员冻结服务，服务冻结失败
func (suite Model13) Test1338_FreezeServerWithAdminIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ServiceMgrContractAddr.Address(), "FreezeService", nil, rpcx.String(chainID+":"+suite.GetServerID(pk)), rpcx.String("reason"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
}

//tc：服务处于registing状态冻结服务，服务冻结失败
func (suite Model13) Test1339_FreezeServerWithRegistingServerIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ServerToRegisting(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.FreezeService(chainID + ":" + suite.GetServerID(pk))
	suite.Require().NotNil(err)
}

//tc：服务处于unavailable状态冻结服务，服务冻结失败
func (suite Model13) Test1340_FreezeServerWithUnavailableServerIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ServerToUnavailable(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.FreezeService(chainID + ":" + suite.GetServerID(pk))
	suite.Require().NotNil(err)
}

//tc：服务处于updating状态冻结服务，服务冻结成功
func (suite Model13) Test1341_FreezeServerWithUpdatingServerIsSuccess() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ServerToUpdating(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.FreezeService(chainID + ":" + suite.GetServerID(pk))
	suite.Require().Nil(err)
	err = suite.CheckServerStatus(chainID+":"+suite.GetServerID(pk), governance.GovernanceFrozen)
	suite.Require().Nil(err)
}

//tc：服务处于activating状态冻结服务，服务冻结成功
func (suite Model13) Test1342_FreezeServerWithActivatingServerIsSuccess() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ServerToActivating(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.FreezeService(chainID + ":" + suite.GetServerID(pk))
	suite.Require().Nil(err)
	err = suite.CheckServerStatus(chainID+":"+suite.GetServerID(pk), governance.GovernanceFrozen)
	suite.Require().Nil(err)
}

//tc：服务处于freezing状态冻结服务，服务冻结失败
func (suite Model13) Test1343_FreezeServerWithFreezingServerIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ServerToFreezing(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.FreezeService(chainID + ":" + suite.GetServerID(pk))
	suite.Require().NotNil(err)
}

//tc：服务处于frozen状态冻结服务，服务冻结失败
func (suite Model13) Test1344_FreezeServerWithFrozenServerIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ServerToFrozen(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.FreezeService(chainID + ":" + suite.GetServerID(pk))
	suite.Require().NotNil(err)
}

//tc：服务处于logouting状态冻结服务，服务冻结失败
func (suite Model13) Test1345_FreezeServerWithLogoutingServerIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ServerToLogouting(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.FreezeService(chainID + ":" + suite.GetServerID(pk))
	suite.Require().NotNil(err)
}

//tc：服务处于forbidden状态冻结服务，服务冻结失败
func (suite Model13) Test1346_FreezeServerWithForbiddenServerIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ServerToForbidden(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.FreezeService(chainID + ":" + suite.GetServerID(pk))
	suite.Require().NotNil(err)
}

//tc：应用链处于activating状态冻结服务，服务冻结失败
func (suite Model13) Test1347_FreezeServerWithActivatingChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	suite.Require().Nil(err)
	err = suite.FreezeAppchain(chainID)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "ActivateAppchain", nil, rpcx.String(chainID), rpcx.String("reason"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	err = suite.CheckChainStatus(chainID, governance.GovernanceActivating)
	suite.Require().Nil(err)
	err = suite.FreezeService(chainID + ":" + suite.GetServerID(pk))
	suite.Require().NotNil(err)
}

//tc：应用链处于freezing状态冻结服务，服务冻结成功
func (suite Model13) Test1348_FreezeServerWithFreezingChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	suite.Require().Nil(err)
	path, err := repo.Node1Path()
	suite.Require().Nil(err)
	node1Key, err := asym.RestorePrivateKey(path, repo.KeyPassword)
	suite.Require().Nil(err)
	client := suite.NewClient(node1Key)
	pubAddress, err := node1Key.PublicKey().Address()
	suite.Require().Nil(err)
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
	err = suite.FreezeService(chainID + ":" + suite.GetServerID(pk))
	suite.Require().Nil(err)
}

//tc：应用链处于frozen状态更新服务，服务更新失败
func (suite Model13) Test1349_FreezeServerWithFrozenChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	suite.Require().Nil(err)
	err = suite.FreezeAppchain(chainID)
	suite.Require().Nil(err)
	err = suite.CheckChainStatus(chainID, governance.GovernanceFrozen)
	suite.Require().Nil(err)
	err = suite.FreezeService(chainID + ":" + suite.GetServerID(pk))
	suite.Require().NotNil(err)
}

//tc：应用链处于logouting状态更新服务，服务更新失败
func (suite Model13) Test1350_FreezeServerWithLogoutingChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "LogoutAppchain", nil, rpcx.String(chainID), rpcx.String("reason"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	err = suite.CheckChainStatus(chainID, governance.GovernanceLogouting)
	suite.Require().Nil(err)
	err = suite.FreezeService(chainID + ":" + suite.GetServerID(pk))
	suite.Require().NotNil(err)
}

//tc：应用链处于forbidden状态更新服务，服务更新失败
func (suite Model13) Test1351_FreezeServerWithForbiddenChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	suite.Require().Nil(err)
	err = suite.LogoutAppchain(pk, chainID)
	suite.Require().Nil(err)
	err = suite.CheckChainStatus(chainID, governance.GovernanceForbidden)
	suite.Require().Nil(err)
	err = suite.FreezeService(chainID + ":" + suite.GetServerID(pk))
	suite.Require().NotNil(err)
}

//tc：中继链管理员激活服务，服务激活成功
func (suite Model13) Test1352_ActivateServerWithRelayAdminIsSuccess() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ServerToFrozen(pk, chainID, address)
	suite.Require().Nil(err)
	path, err := repo.Node1Path()
	suite.Require().Nil(err)
	node1pk, err := asym.RestorePrivateKey(path, repo.KeyPassword)
	suite.Require().Nil(err)
	from, err := node1pk.PublicKey().Address()
	suite.Require().Nil(err)
	client := suite.NewClient(node1pk)
	res, err := client.InvokeBVMContract(constant.ServiceMgrContractAddr.Address(), "ActivateService", &rpcx.TransactOpts{
		From:  from.String(),
		Nonce: atomic.AddUint64(&nonce1, 1),
	}, rpcx.String(chainID+":"+suite.GetServerID(pk)), rpcx.String("reason"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	suite.Require().Nil(err)
	err = suite.VotePass(result.ProposalID)
	suite.Require().Nil(err)
	err = suite.CheckServerStatus(chainID+":"+suite.GetServerID(pk), governance.GovernanceAvailable)
	suite.Require().Nil(err)
}

//tc：应用链管理员激活服务，服务激活成功
func (suite Model13) Test1353_ActivateServerIsSuccess() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	suite.Require().Nil(err)
	err = suite.FreezeService(chainID + ":" + suite.GetServerID(pk))
	suite.Require().Nil(err)
	err = suite.CheckServerStatus(chainID+":"+suite.GetServerID(pk), governance.GovernanceFrozen)
	suite.Require().Nil(err)
	err = suite.ActivateService(pk, chainID+":"+suite.GetServerID(pk))
	suite.Require().Nil(err)
}

//tc：服务处于registing状态激活服务，服务激活失败
func (suite Model13) Test1354_ActivateServerWithRegistingServerIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ServerToRegisting(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.ActivateService(pk, chainID+":"+suite.GetServerID(pk))
	suite.Require().NotNil(err)
}

//tc：服务处于unavailable状态激活服务，服务激活失败
func (suite Model13) Test1355_ActivateServerWithUnavailableServerIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ServerToUnavailable(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.ActivateService(pk, chainID+":"+suite.GetServerID(pk))
	suite.Require().NotNil(err)
}

//tc：服务处于updating状态激活服务，服务激活失败
func (suite Model13) Test1356_ActivateServerWithUpdatingServerIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ServerToUpdating(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.ActivateService(pk, chainID+":"+suite.GetServerID(pk))
	suite.Require().NotNil(err)
}

//tc：服务处于activating状态激活服务，服务激活失败
func (suite Model13) Test1357_ActivateServerWithActivatingServerIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ServerToActivating(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.ActivateService(pk, chainID+":"+suite.GetServerID(pk))
	suite.Require().NotNil(err)
}

//tc：服务处于freezing状态激活服务，服务激活失败
func (suite Model13) Test1358_ActivateServerWithFreezingServerIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ServerToFreezing(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.ActivateService(pk, chainID+":"+suite.GetServerID(pk))
	suite.Require().NotNil(err)
}

//tc：服务处于logouting状态激活服务，服务激活失败
func (suite Model13) Test1359_ActivateServerWithLogoutingServerIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ServerToLogouting(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.ActivateService(pk, chainID+":"+suite.GetServerID(pk))
	suite.Require().NotNil(err)
}

//tc：服务处于forbidden状态激活服务，服务激活失败
func (suite Model13) Test1360_ActivateServerWithForbiddenServerIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ServerToForbidden(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.ActivateService(pk, chainID+":"+suite.GetServerID(pk))
	suite.Require().NotNil(err)
}

//tc：应用链处于activating状态激活服务，服务激活失败
func (suite Model13) Test1361_ActivateServerWithActivatingChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	suite.Require().Nil(err)
	err = suite.FreezeAppchain(chainID)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "ActivateAppchain", nil, rpcx.String(chainID), rpcx.String("reason"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	err = suite.CheckChainStatus(chainID, governance.GovernanceActivating)
	suite.Require().Nil(err)
	err = suite.ActivateService(pk, chainID+":"+suite.GetServerID(pk))
	suite.Require().NotNil(err)
}

//tc：应用链处于freezing状态激活服务，服务激活成功
func (suite Model13) Test1362_ActivateServerWithFreezingChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	suite.Require().Nil(err)
	err = suite.FreezeService(chainID + ":" + suite.GetServerID(pk))
	suite.Require().Nil(err)
	path, err := repo.Node1Path()
	suite.Require().Nil(err)
	node1Key, err := asym.RestorePrivateKey(path, repo.KeyPassword)
	suite.Require().Nil(err)
	client := suite.NewClient(node1Key)
	pubAddress, err := node1Key.PublicKey().Address()
	suite.Require().Nil(err)
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
	err = suite.ActivateService(pk, chainID+":"+suite.GetServerID(pk))
	suite.Require().Nil(err)
	err = suite.CheckServerStatus(chainID+":"+suite.GetServerID(pk), governance.GovernanceAvailable)
	suite.Require().Nil(err)
}

//tc：应用链处于frozen状态激活服务，服务激活失败
func (suite Model13) Test1363_ActivateServerWithFrozenChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	suite.Require().Nil(err)
	err = suite.FreezeAppchain(chainID)
	suite.Require().Nil(err)
	err = suite.CheckChainStatus(chainID, governance.GovernanceFrozen)
	suite.Require().Nil(err)
	err = suite.ActivateService(pk, chainID+":"+suite.GetServerID(pk))
	suite.Require().NotNil(err)
}

//tc：应用链处于logouting状态激活服务，服务激活失败
func (suite Model13) Test1364_ActivateServerWithLogoutingChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "LogoutAppchain", nil, rpcx.String(chainID), rpcx.String("reason"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	err = suite.CheckChainStatus(chainID, governance.GovernanceLogouting)
	suite.Require().Nil(err)
	err = suite.ActivateService(pk, chainID+":"+suite.GetServerID(pk))
	suite.Require().NotNil(err)
}

//tc：应用链处于forbidden状态激活服务，服务激活失败
func (suite Model13) Test1365_ActivateServerWithForbiddenChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	suite.Require().Nil(err)
	err = suite.LogoutAppchain(pk, chainID)
	suite.Require().Nil(err)
	err = suite.CheckChainStatus(chainID, governance.GovernanceForbidden)
	suite.Require().Nil(err)
	err = suite.ActivateService(pk, chainID+":"+suite.GetServerID(pk))
	suite.Require().NotNil(err)
}

//tc：中继链管理员注销服务，服务注销失败
func (suite Model13) Test1366_LogoutServerWithRelayAdminIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	suite.Require().Nil(err)
	path, err := repo.Node1Path()
	suite.Require().Nil(err)
	node1pk, err := asym.RestorePrivateKey(path, repo.KeyPassword)
	suite.Require().Nil(err)
	from, err := node1pk.PublicKey().Address()
	suite.Require().Nil(err)
	client := suite.NewClient(node1pk)
	res, err := client.InvokeBVMContract(constant.ServiceMgrContractAddr.Address(), "LogoutService", &rpcx.TransactOpts{
		From:  from.String(),
		Nonce: atomic.AddUint64(&nonce1, 1),
	}, rpcx.String(chainID+":"+suite.GetServerID(pk)), rpcx.String("reason"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
}

//tc：应用链管理员注销服务，服务注销成功
func (suite Model13) Test1367_LogoutServerWithRelayAdminIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	suite.Require().Nil(err)
	err = suite.LogoutService(pk, chainID+":"+suite.GetServerID(pk))
	suite.Require().Nil(err)
}

//tc：服务处于registing状态注销服务，服务注销失败
func (suite Model13) Test1368_LogoutServerWithRegistingServerIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ServerToRegisting(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.LogoutService(pk, chainID+":"+suite.GetServerID(pk))
	suite.Require().NotNil(err)
}

//tc：服务处于unavailable状态注销服务，服务注销失败
func (suite Model13) Test1369_LogoutServerWithUnavailableServerIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ServerToUnavailable(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.LogoutService(pk, chainID+":"+suite.GetServerID(pk))
	suite.Require().NotNil(err)
}

//tc：服务处于updating状态注销服务，服务注销成功
func (suite Model13) Test1370_LogoutServerWithUpdatingServerIsSuccess() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ServerToUpdating(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.LogoutService(pk, chainID+":"+suite.GetServerID(pk))
	suite.Require().Nil(err)
}

//tc：服务处于activating状态注销服务，服务注销成功
func (suite Model13) Test1371_LogoutServerWithActivatingServerIsSuccess() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ServerToActivating(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.LogoutService(pk, chainID+":"+suite.GetServerID(pk))
	suite.Require().Nil(err)
}

//tc：服务处于freezing状态注销服务，服务注销成功
func (suite Model13) Test1372_LogoutServerWithFreezingServerIsSuccess() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ServerToFreezing(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.LogoutService(pk, chainID+":"+suite.GetServerID(pk))
	suite.Require().Nil(err)
}

//tc：服务处于frozen状态注销服务，服务注销失败
func (suite Model13) Test1373_LogoutServerWithFrozenServerIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ServerToFrozen(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.LogoutService(pk, chainID+":"+suite.GetServerID(pk))
	suite.Require().Nil(err)
}

//tc：服务处于logouting状态注销服务，服务注销失败
func (suite Model13) Test1374_LogoutServerWithLogoutingServerIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ServerToLogouting(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.LogoutService(pk, chainID+":"+suite.GetServerID(pk))
	suite.Require().NotNil(err)
}

//tc：服务处于forbidden状态注销服务，服务注销失败
func (suite Model13) Test1375_LogoutServerWithForbiddenServerIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ServerToForbidden(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.LogoutService(pk, chainID+":"+suite.GetServerID(pk))
	suite.Require().NotNil(err)
}

//tc：应用链处于activating状态注销服务，服务注销失败
func (suite Model13) Test1376_LogoutServerWithActivatingChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	suite.Require().Nil(err)
	err = suite.FreezeAppchain(chainID)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "ActivateAppchain", nil, rpcx.String(chainID), rpcx.String("reason"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	err = suite.CheckChainStatus(chainID, governance.GovernanceActivating)
	suite.Require().Nil(err)
	err = suite.LogoutService(pk, chainID+":"+suite.GetServerID(pk))
	suite.Require().NotNil(err)
}

//tc：应用链处于freezing状态注销服务，服务注销成功
func (suite Model13) Test1377_LogoutServerWithFreezingChainIsSuccess() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	suite.Require().Nil(err)
	path, err := repo.Node1Path()
	suite.Require().Nil(err)
	node1Key, err := asym.RestorePrivateKey(path, repo.KeyPassword)
	suite.Require().Nil(err)
	client := suite.NewClient(node1Key)
	pubAddress, err := node1Key.PublicKey().Address()
	suite.Require().Nil(err)
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
	err = suite.LogoutService(pk, chainID+":"+suite.GetServerID(pk))
	suite.Require().Nil(err)
}

//tc：应用链处于frozen状态注销服务，服务注销失败
func (suite Model13) Test1378_LogoutServerWithFrozenChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	suite.Require().Nil(err)
	err = suite.FreezeAppchain(chainID)
	suite.Require().Nil(err)
	err = suite.CheckChainStatus(chainID, governance.GovernanceFrozen)
	suite.Require().Nil(err)
	err = suite.LogoutService(pk, chainID+":"+suite.GetServerID(pk))
	suite.Require().NotNil(err)
}

//tc：应用链处于logouting状态注销服务，服务注销失败
func (suite Model13) Test1379_LogoutServerWithLogoutingChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "LogoutAppchain", nil, rpcx.String(chainID), rpcx.String("reason"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	err = suite.CheckChainStatus(chainID, governance.GovernanceLogouting)
	suite.Require().Nil(err)
	err = suite.LogoutService(pk, chainID+":"+suite.GetServerID(pk))
	suite.Require().NotNil(err)
}

//tc：应用链处于forbidden状态注销服务，服务注销失败
func (suite Model13) Test1380_LogoutServerWithForbiddenChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	suite.Require().Nil(err)
	err = suite.LogoutAppchain(pk, chainID)
	suite.Require().Nil(err)
	err = suite.CheckChainStatus(chainID, governance.GovernanceForbidden)
	suite.Require().Nil(err)
	err = suite.LogoutService(pk, chainID+":"+suite.GetServerID(pk))
	suite.Require().NotNil(err)
}

func (suite Snake) RegisterServer(pk crypto.PrivateKey, chainID, serviceID, typ string) error {
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(chainID),
		rpcx.String(serviceID),
		rpcx.String("testServer"),
		rpcx.String(typ),
		rpcx.String("test"),
		rpcx.Bool(true),
		rpcx.String(""),
		rpcx.String("test"),
		rpcx.String("reason"),
	}
	res, err := client.InvokeBVMContract(constant.ServiceMgrContractAddr.Address(), "RegisterService", nil, args...)
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

func (suite Snake) FreezeService(chainServiceID string) error {
	path, err := repo.Node1Path()
	if err != nil {
		return err
	}
	node1pk, err := asym.RestorePrivateKey(path, repo.KeyPassword)
	if err != nil {
		return err
	}
	from, err := node1pk.PublicKey().Address()
	if err != nil {
		return err
	}
	client := suite.NewClient(node1pk)
	res, err := client.InvokeBVMContract(constant.ServiceMgrContractAddr.Address(), "FreezeService", &rpcx.TransactOpts{
		From:  from.String(),
		Nonce: atomic.AddUint64(&nonce1, 1),
	}, rpcx.String(chainServiceID), rpcx.String("reason"))
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

func (suite Model13) ActivateService(pk crypto.PrivateKey, chainServiceID string) error {
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ServiceMgrContractAddr.Address(), "ActivateService", nil, rpcx.String(chainServiceID), rpcx.String("reason"))
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

func (suite Snake) UpdateService(pk crypto.PrivateKey, chainServiceID, name string) error {
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(chainServiceID),
		rpcx.String(name),
		rpcx.String("test"),
		rpcx.Bool(true),
		rpcx.String(""),
		rpcx.String("test"),
		rpcx.String("reason"),
	}
	res, err := client.InvokeBVMContract(constant.ServiceMgrContractAddr.Address(), "UpdateService", nil, args...)
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return errors.New(string(res.Ret))
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

func (suite Model13) LogoutService(pk crypto.PrivateKey, chainServiceID string) error {
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ServiceMgrContractAddr.Address(), "LogoutService", nil, rpcx.String(chainServiceID), rpcx.String("reason"))
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

func (suite Snake) CheckServerStatus(serverID string, status governance.GovernanceStatus) error {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ServiceMgrContractAddr.Address(), "GetServiceInfo", nil, rpcx.String(serverID))
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return errors.New(string(res.Ret))
	}
	server := &serviceMgr.Service{}
	err = json.Unmarshal(res.Ret, server)
	if err != nil {
		return err
	}
	if server.Status != status {
		return errors.New(fmt.Sprintf("expect status is %s, but got %s", status, server.Status))
	}
	return nil
}

func (suite Model13) ServerToRegisting(pk crypto.PrivateKey, chainID, address string) error {
	err := suite.RegisterAppchain(pk, chainID, address)
	if err != nil {
		return err
	}
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(chainID),
		rpcx.String(suite.GetServerID(pk)),
		rpcx.String("testServer"),
		rpcx.String("CallContract"),
		rpcx.String("test"),
		rpcx.Bool(true),
		rpcx.String(""),
		rpcx.String("test"),
		rpcx.String("reason"),
	}
	res, err := client.InvokeBVMContract(constant.ServiceMgrContractAddr.Address(), "RegisterService", nil, args...)
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return errors.New(string(res.Ret))
	}
	err = suite.CheckServerStatus(chainID+":"+suite.GetServerID(pk), governance.GovernanceRegisting)
	if err != nil {
		return err
	}
	return nil
}

func (suite Model13) ServerToUnavailable(pk crypto.PrivateKey, chainID, address string) error {
	err := suite.RegisterAppchain(pk, chainID, address)
	if err != nil {
		return err
	}
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(chainID),
		rpcx.String(suite.GetServerID(pk)),
		rpcx.String("testServer"),
		rpcx.String("CallContract"),
		rpcx.String("test"),
		rpcx.Bool(true),
		rpcx.String(""),
		rpcx.String("test"),
		rpcx.String("reason"),
	}
	res, err := client.InvokeBVMContract(constant.ServiceMgrContractAddr.Address(), "RegisterService", nil, args...)
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

func (suite Model13) ServerToUpdating(pk crypto.PrivateKey, chainID, address string) error {
	err := suite.RegisterAppchain(pk, chainID, address)
	if err != nil {
		return err
	}
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	if err != nil {
		return err
	}
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(chainID + ":" + suite.GetServerID(pk)),
		rpcx.String("name"),
		rpcx.String("test"),
		rpcx.Bool(true),
		rpcx.String(""),
		rpcx.String("test"),
		rpcx.String("reason"),
	}
	res, err := client.InvokeBVMContract(constant.ServiceMgrContractAddr.Address(), "UpdateService", nil, args...)
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return errors.New(string(res.Ret))
	}
	err = suite.CheckServerStatus(chainID+":"+suite.GetServerID(pk), governance.GovernanceUpdating)
	if err != nil {
		return err
	}
	return nil
}

func (suite Model13) ServerToActivating(pk crypto.PrivateKey, chainID, address string) error {
	err := suite.RegisterAppchain(pk, chainID, address)
	if err != nil {
		return err
	}
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	if err != nil {
		return err
	}
	err = suite.FreezeService(chainID + ":" + suite.GetServerID(pk))
	if err != nil {
		return err
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ServiceMgrContractAddr.Address(), "ActivateService", nil, rpcx.String(chainID+":"+suite.GetServerID(pk)), rpcx.String("reason"))
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return errors.New(string(res.Ret))
	}
	err = suite.CheckServerStatus(chainID+":"+suite.GetServerID(pk), governance.GovernanceActivating)
	if err != nil {
		return err
	}
	return nil
}

func (suite Model13) ServerToFreezing(pk crypto.PrivateKey, chainID, address string) error {
	err := suite.RegisterAppchain(pk, chainID, address)
	if err != nil {
		return err
	}
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	if err != nil {
		return err
	}
	path, err := repo.Node1Path()
	if err != nil {
		return err
	}
	node1pk, err := asym.RestorePrivateKey(path, repo.KeyPassword)
	if err != nil {
		return err
	}
	from, err := node1pk.PublicKey().Address()
	if err != nil {
		return err
	}
	client := suite.NewClient(node1pk)
	res, err := client.InvokeBVMContract(constant.ServiceMgrContractAddr.Address(), "FreezeService", &rpcx.TransactOpts{
		From:  from.String(),
		Nonce: atomic.AddUint64(&nonce1, 1),
	}, rpcx.String(chainID+":"+suite.GetServerID(pk)), rpcx.String("reason"))
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return errors.New(string(res.Ret))
	}
	err = suite.CheckServerStatus(chainID+":"+suite.GetServerID(pk), governance.GovernanceFreezing)
	if err != nil {
		return err
	}
	return nil
}

func (suite Model13) ServerToFrozen(pk crypto.PrivateKey, chainID, address string) error {
	err := suite.RegisterAppchain(pk, chainID, address)
	if err != nil {
		return err
	}
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	if err != nil {
		return err
	}
	err = suite.FreezeService(chainID + ":" + suite.GetServerID(pk))
	if err != nil {
		return err
	}
	err = suite.CheckServerStatus(chainID+":"+suite.GetServerID(pk), governance.GovernanceFrozen)
	if err != nil {
		return err
	}
	return nil
}

func (suite Model13) ServerToLogouting(pk crypto.PrivateKey, chainID, address string) error {
	err := suite.RegisterAppchain(pk, chainID, address)
	if err != nil {
		return err
	}
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	if err != nil {
		return err
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ServiceMgrContractAddr.Address(), "LogoutService", nil, rpcx.String(chainID+":"+suite.GetServerID(pk)), rpcx.String("reason"))
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return errors.New(string(res.Ret))
	}
	err = suite.CheckServerStatus(chainID+":"+suite.GetServerID(pk), governance.GovernanceLogouting)
	if err != nil {
		return err
	}
	return nil
}

func (suite Model13) ServerToForbidden(pk crypto.PrivateKey, chainID, address string) error {
	err := suite.RegisterAppchain(pk, chainID, address)
	if err != nil {
		return err
	}
	err = suite.RegisterServer(pk, chainID, suite.GetServerID(pk), "CallContract")
	if err != nil {
		return err
	}
	err = suite.LogoutService(pk, chainID+":"+suite.GetServerID(pk))
	if err != nil {
		return err
	}
	err = suite.CheckServerStatus(chainID+":"+suite.GetServerID(pk), governance.GovernanceForbidden)
	if err != nil {
		return err
	}
	return nil
}
