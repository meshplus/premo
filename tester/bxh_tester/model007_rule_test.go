package bxh_tester

import (
	"encoding/json"

	"github.com/meshplus/bitxhub-core/governance"
	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-kit/types"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/pkg/errors"
)

const (
	RegisterRule     = "RegisterRule"
	UpdateMasterRule = "UpdateMasterRule"
	LogoutRule       = "LogoutRule"
)

const (
	FabricRuleAddr    = "0x00000000000000000000000000000000000000a0"
	SimFabricRuleAddr = "0x00000000000000000000000000000000000000a1"
	HappyRuleAddr     = "0x00000000000000000000000000000000000000a2"
)

type Model7 struct {
	*Snake
}

func (suite *Model7) SetupTest() {
	suite.T().Parallel()
}

//tc：正确部署验证规则,并返回地址
func (suite Model7) Test0701_DeployRuleIsSuccess() {
	address, err := suite.DeploySimpleRule()
	suite.Require().Nil(err)
	suite.Require().NotNil(address)
}

//tc：部署验证规则字段为空，并提示错误信息
func (suite Model7) Test0702_DeployRuleIsFail() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	address, err := client.DeployContract([]byte(""), nil)
	suite.Require().NotNil(err)
	suite.Require().Nil(address)
}

//tc：非应用链管理员调用注册验证规则，验证规则注册失败
func (suite Model7) Test0703_RegisterRuleWithNoAdminIsFail() {
	_, chainID, _, err := suite.RegisterRule()
	suite.Require().Nil(err)
	address, err := suite.DeploySimpleRule()
	suite.Require().Nil(err)
	pk2, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk2, chainID, address, RegisterRule)
	suite.Require().NotNil(err)
}

//tc：应用链管理员未注册调用注册验证规则，验证规则注册成功
func (suite Model7) Test0704_RegisterRuleIsSuccess() {
	_, _, _, err := suite.RegisterRule()
	suite.Require().Nil(err)
}

//tc：应用链管理员已注册调用注册验证规则，验证规则注册成功
func (suite Model7) Test0705_RegisterRuleWithRegisteredAdminIsSuccess() {
	pk, chainID, _, err := suite.RegisterRule()
	suite.Require().Nil(err)
	address, err := suite.DeploySimpleRule()
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, chainID, address, RegisterRule)
	suite.Require().Nil(err)
}

//tc：应用链管理员注册未部署的验证规则，验证规则部署失败
func (suite Model7) Test0706_RegisterRuleWithNoRegisterRuleIsFail() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, suite.GetChainID(pk), "0x000000000000000000000000000000000000001", RegisterRule)
	suite.Require().NotNil(err)
}

//tc：应用链管理员注册available的验证规则，验证规则部署失败
func (suite Model7) Test0708_RegisterRuleWithAvailableRuleIsFail() {
	address, err := suite.DeploySimpleRule()
	suite.Require().Nil(err)
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	chainID := suite.GetChainID(pk)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, chainID, address, RegisterRule)
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.CheckRuleStatus(pk, chainID, address, governance.GovernanceAvailable)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, chainID, address, RegisterRule)
	suite.Require().NotNil(err)
}

//tc：应用链管理员注册binding的验证规则，验证规则部署失败
//tc：应用链管理员注册unbinding的验证规则，验证规则部署失败
func (suite Model7) Test0709_RegisterRuleWithBindingRuleIsFail() {
	address1, err := suite.DeploySimpleRule()
	suite.Require().Nil(err)
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	chainID := suite.GetChainID(pk)
	err = suite.InvokeRuleContract(pk, chainID, address1, RegisterRule)
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address1)
	suite.Require().Nil(err)
	err = suite.CheckRuleStatus(pk, chainID, address1, governance.GovernanceAvailable)
	suite.Require().Nil(err)
	address2, err := suite.DeploySimpleRule()
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, chainID, address2, RegisterRule)
	suite.Require().Nil(err)
	err = suite.CheckRuleStatus(pk, chainID, address2, governance.GovernanceBindable)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(chainID),
		rpcx.String(address2),
		rpcx.String("reason"),
	}
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), UpdateMasterRule, nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	err = suite.CheckRuleStatus(pk, chainID, address2, governance.GovernanceBinding)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, chainID, address2, RegisterRule)
	suite.Require().NotNil(err)
	err = suite.CheckRuleStatus(pk, chainID, address1, governance.GovernanceUnbinding)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, chainID, address1, RegisterRule)
	suite.Require().NotNil(err)
}

//tc：应用链管理员注册bindable的验证规则，验证规则部署失败
func (suite Model7) Test0710_RegisterRuleWithBindableRuleIsFail() {
	address1, err := suite.DeploySimpleRule()
	suite.Require().Nil(err)
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	chainID := suite.GetChainID(pk)
	err = suite.InvokeRuleContract(pk, chainID, address1, RegisterRule)
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address1)
	suite.Require().Nil(err)
	err = suite.CheckRuleStatus(pk, chainID, address1, governance.GovernanceAvailable)
	suite.Require().Nil(err)
	address2, err := suite.DeploySimpleRule()
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, chainID, address2, RegisterRule)
	suite.Require().Nil(err)
	err = suite.CheckRuleStatus(pk, chainID, address2, governance.GovernanceBindable)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, chainID, address2, RegisterRule)
	suite.Require().NotNil(err)
}

//tc：应用链管理员注册forbidden的验证规则，验证规则部署失败
func (suite Model7) Test0711_RegisterRuleWithForbiddenRuleIsFail() {
	address1, err := suite.DeploySimpleRule()
	suite.Require().Nil(err)
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	chainID := suite.GetChainID(pk)
	err = suite.InvokeRuleContract(pk, chainID, address1, RegisterRule)
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address1)
	suite.Require().Nil(err)
	err = suite.CheckRuleStatus(pk, chainID, address1, governance.GovernanceAvailable)
	suite.Require().Nil(err)
	address2, err := suite.DeploySimpleRule()
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, chainID, address2, RegisterRule)
	suite.Require().Nil(err)
	err = suite.CheckRuleStatus(pk, chainID, address2, governance.GovernanceBindable)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, chainID, address2, LogoutRule)
	suite.Require().Nil(err)
	err = suite.CheckRuleStatus(pk, chainID, address2, governance.GovernanceForbidden)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, chainID, address2, RegisterRule)
	suite.Require().NotNil(err)
}

//tc：非应用链管理员更新主验证规则，验证规则更新失败
//tc：非应用链管理员更新主验证规则，验证规则注销失败
func (suite Model7) Test0712_UpdateAndLogoutRuleWithNoAdminIsSuccess() {
	address1, err := suite.DeploySimpleRule()
	suite.Require().Nil(err)
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	chainID := suite.GetChainID(pk)
	err = suite.InvokeRuleContract(pk, chainID, address1, RegisterRule)
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address1)
	suite.Require().Nil(err)
	err = suite.CheckRuleStatus(pk, chainID, address1, governance.GovernanceAvailable)
	suite.Require().Nil(err)
	address2, err := suite.DeploySimpleRule()
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, chainID, address2, RegisterRule)
	suite.Require().Nil(err)
	err = suite.CheckRuleStatus(pk, chainID, address2, governance.GovernanceBindable)
	suite.Require().Nil(err)
	pk2, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk2, chainID, address2, UpdateMasterRule)
	suite.Require().NotNil(err)
	err = suite.InvokeRuleContract(pk2, chainID, address2, LogoutRule)
	suite.Require().NotNil(err)
}

//tc：应用链管理员更新主验证规则，验证规则更新成功
//tc：应用链管理员更新主验证规则，验证规则注销成功
func (suite Model7) Test0713_UpdateAndLogoutRuleIsSuccess() {
	address1, err := suite.DeploySimpleRule()
	suite.Require().Nil(err)
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	chainID := suite.GetChainID(pk)
	err = suite.InvokeRuleContract(pk, chainID, address1, RegisterRule)
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address1)
	suite.Require().Nil(err)
	err = suite.CheckRuleStatus(pk, chainID, address1, governance.GovernanceAvailable)
	suite.Require().Nil(err)
	address2, err := suite.DeploySimpleRule()
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, chainID, address2, RegisterRule)
	suite.Require().Nil(err)
	err = suite.CheckRuleStatus(pk, chainID, address2, governance.GovernanceBindable)
	suite.Require().Nil(err)
	err = suite.CheckChainStatus(chainID, governance.GovernanceAvailable)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, chainID, address2, UpdateMasterRule)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, chainID, address1, LogoutRule)
	suite.Require().Nil(err)
	err = suite.CheckChainStatus(chainID, governance.GovernanceAvailable)
	suite.Require().Nil(err)
}

//tc:应用链处于未注册状态更新主验证规则，验证规则更新失败
//tc:应用链处于未注册状态更新主注销规则，验证规则注销失败
func (suite Model7) Test0714_UpdateAndLogoutRuleWithNoRegisterAdmin() {
	address1, err := suite.DeploySimpleRule()
	suite.Require().Nil(err)
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	chainID := suite.GetChainID(pk)
	err = suite.InvokeRuleContract(pk, chainID, address1, RegisterRule)
	suite.Require().Nil(err)
	err = suite.CheckRuleStatus(pk, chainID, address1, governance.GovernanceBindable)
	suite.Require().Nil(err)
	address2, err := suite.DeploySimpleRule()
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, chainID, address2, RegisterRule)
	suite.Require().Nil(err)
	err = suite.CheckRuleStatus(pk, chainID, address2, governance.GovernanceBindable)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, chainID, address2, UpdateMasterRule)
	suite.Require().NotNil(err)
	err = suite.InvokeRuleContract(pk, chainID, address2, LogoutRule)
	suite.Require().Nil(err)
	err = suite.CheckRuleStatus(pk, chainID, address2, governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

//tc：应用链处于registering状态更新主验证规则，验证规则更新失败
//tc：应用链处于registering状态注销验证规则，验证规则注销成功
func (suite Model7) Test0715_UpdateAndLogoutRuleWithRegistingChain() {
	pk, chainID, address1, err := suite.RegisterRule()
	suite.Require().Nil(err)
	address2, err := suite.DeploySimpleRule()
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, suite.GetChainID(pk), address2, RegisterRule)
	suite.Require().Nil(err)
	err = suite.ChainToRegisting(pk, chainID, address1)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, chainID, address2, UpdateMasterRule)
	suite.Require().NotNil(err)
	err = suite.InvokeRuleContract(pk, chainID, address2, LogoutRule)
	suite.Require().Nil(err)
	err = suite.CheckChainStatus(chainID, governance.GovernanceRegisting)
	suite.Require().Nil(err)
}

//tc：应用链处于unavailable状态更新主验证规则，验证规则更新失败
//tc：应用链处于unavailable状态注销验证规则，验证规则注销成功
func (suite Model7) Test0716_UpdateAndLogoutRuleWithUnavailableChain() {
	pk, chainID, address1, err := suite.RegisterRule()
	suite.Require().Nil(err)
	address2, err := suite.DeploySimpleRule()
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, suite.GetChainID(pk), address2, RegisterRule)
	suite.Require().Nil(err)
	err = suite.ChainToUnavailable(pk, chainID, address1)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, chainID, address2, UpdateMasterRule)
	suite.Require().NotNil(err)
	err = suite.InvokeRuleContract(pk, chainID, address2, LogoutRule)
	suite.Require().Nil(err)
	err = suite.CheckChainStatus(chainID, governance.GovernanceUnavailable)
	suite.Require().Nil(err)
}

//tc：应用链处于activating状态更新主验证规则，验证规则更新失败
//tc：应用链处于activating状态注销验证规则，验证规则注销成功
func (suite Model7) Test0717_UpdateAndLogoutRuleWithActivatingChain() {
	pk, chainID, address1, err := suite.RegisterRule()
	suite.Require().Nil(err)
	address2, err := suite.DeploySimpleRule()
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, suite.GetChainID(pk), address2, RegisterRule)
	suite.Require().Nil(err)
	err = suite.ChainToActivating(pk, chainID, address1)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, chainID, address2, UpdateMasterRule)
	suite.Require().NotNil(err)
	err = suite.InvokeRuleContract(pk, chainID, address2, LogoutRule)
	suite.Require().Nil(err)
	err = suite.CheckChainStatus(chainID, governance.GovernanceActivating)
	suite.Require().Nil(err)
}

//tc：应用链处于freezing状态更新主验证规则，验证规则更新失败
//tc：应用链处于freezing状态注销验证规则，验证规则注销成功
func (suite Model7) Test0717_UpdateAndLogoutRuleWithFreezingChain() {
	pk, chainID, address1, err := suite.RegisterRule()
	suite.Require().Nil(err)
	address2, err := suite.DeploySimpleRule()
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, suite.GetChainID(pk), address2, RegisterRule)
	suite.Require().Nil(err)
	err = suite.ChainToFreezing(pk, chainID, address1)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, chainID, address2, UpdateMasterRule)
	suite.Require().NotNil(err)
	err = suite.InvokeRuleContract(pk, chainID, address2, LogoutRule)
	suite.Require().Nil(err)
	err = suite.CheckChainStatus(chainID, governance.GovernanceFreezing)
	suite.Require().Nil(err)
}

//tc：应用链处于frozen状态更新主验证规则，验证规则更新成功
//tc：应用链处于frozen状态注销验证规则，验证规则注销成功
func (suite Model7) Test0718_UpdateAndLogoutRuleWithFrozenChain() {
	pk, chainID, address1, err := suite.RegisterRule()
	suite.Require().Nil(err)
	address2, err := suite.DeploySimpleRule()
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, suite.GetChainID(pk), address2, RegisterRule)
	suite.Require().Nil(err)
	err = suite.ChainToFrozen(pk, chainID, address1)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, chainID, address2, UpdateMasterRule)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, chainID, address2, LogoutRule)
	suite.Require().NotNil(err)
	err = suite.CheckChainStatus(chainID, governance.GovernanceFrozen)
	suite.Require().Nil(err)
}

//tc：应用链处于logouting状态更新主验证规则，验证规则更新失败
//tc：应用链处于logouting状态注销验证规则，验证规则注销成功
func (suite Model7) Test0719_UpdateAndLogoutRuleWithLogoutingChainIsFail() {
	pk, chainID, address1, err := suite.RegisterRule()
	suite.Require().Nil(err)
	address2, err := suite.DeploySimpleRule()
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, suite.GetChainID(pk), address2, RegisterRule)
	suite.Require().Nil(err)
	err = suite.ChainToLogouting(pk, chainID, address1)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, chainID, address2, UpdateMasterRule)
	suite.Require().NotNil(err)
	err = suite.InvokeRuleContract(pk, chainID, address2, LogoutRule)
	suite.Require().Nil(err)
	err = suite.CheckChainStatus(chainID, governance.GovernanceLogouting)
	suite.Require().Nil(err)
}

//tc：应用链处于forbidden状态更新主验证规则，验证规则更新失败
//tc：应用链处于forbidden状态注销验证规则，验证规则注销成功
func (suite Model7) Test0720_UpdateAndLogoutRuleWithForbiddenChainIsFail() {
	pk, chainID, address1, err := suite.RegisterRule()
	suite.Require().Nil(err)
	address2, err := suite.DeploySimpleRule()
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, suite.GetChainID(pk), address2, RegisterRule)
	suite.Require().Nil(err)
	err = suite.ChainToForbidden(pk, chainID, address1)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, chainID, address2, UpdateMasterRule)
	suite.Require().NotNil(err)
	err = suite.InvokeRuleContract(pk, chainID, address2, LogoutRule)
	suite.Require().Nil(err)
	err = suite.CheckChainStatus(chainID, governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

//tc：应用链更新未注册的主验证规则，验证规则更新失败
//tc：应用链更新未注册的主验证规则，验证规则注销失败
func (suite Model7) Test0721_UpdateAndLogoutRuleWithNoRegisterRuleIsFail() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	chainID := suite.GetChainID(pk)
	address, err := suite.DeploySimpleRule()
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, chainID, address, RegisterRule)
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, chainID, "0x000000000000000000000000000000000000001", UpdateMasterRule)
	suite.Require().NotNil(err)
	err = suite.InvokeRuleContract(pk, chainID, "0x000000000000000000000000000000000000001", LogoutRule)
	suite.Require().NotNil(err)
}

//tc：应用链管理员更新available状态的主验证规则，验证规则更新失败
//tc：应用链管理员更新available状态的主验证规则，验证规则注销失败
func (suite Model7) Test0722_UpdateAndLogoutRuleWithAvailableRuleIsFail() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	chainID := suite.GetChainID(pk)
	address, err := suite.DeploySimpleRule()
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, chainID, address, RegisterRule)
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.CheckRuleStatus(pk, chainID, address, governance.GovernanceAvailable)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, chainID, address, UpdateMasterRule)
	suite.Require().NotNil(err)
	err = suite.InvokeRuleContract(pk, chainID, address, LogoutRule)
	suite.Require().NotNil(err)
	err = suite.CheckRuleStatus(pk, chainID, address, governance.GovernanceAvailable)
	suite.Require().Nil(err)
}

//tc：应用链管理员更新binding状态的主验证规则，验证规则更新失败
//tc：应用链管理员更新binding状态的主验证规则，验证规则注销失败
func (suite Model7) Test0723_UpdateAndLogoutRuleWithBindingRuleIsSuccess() {
	address1, err := suite.DeploySimpleRule()
	suite.Require().Nil(err)
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	chainID := suite.GetChainID(pk)
	err = suite.InvokeRuleContract(pk, chainID, address1, RegisterRule)
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address1)
	suite.Require().Nil(err)
	err = suite.CheckRuleStatus(pk, chainID, address1, governance.GovernanceAvailable)
	suite.Require().Nil(err)
	address2, err := suite.DeploySimpleRule()
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, chainID, address2, RegisterRule)
	suite.Require().Nil(err)
	err = suite.CheckRuleStatus(pk, chainID, address2, governance.GovernanceBindable)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(chainID),
		rpcx.String(address2),
		rpcx.String("reason"),
	}
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), UpdateMasterRule, nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	err = suite.CheckRuleStatus(pk, chainID, address2, governance.GovernanceBinding)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, chainID, address2, RegisterRule)
	suite.Require().NotNil(err)
	err = suite.CheckRuleStatus(pk, chainID, address1, governance.GovernanceUnbinding)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, chainID, address2, UpdateMasterRule)
	suite.Require().NotNil(err)
	err = suite.InvokeRuleContract(pk, chainID, address2, LogoutRule)
	suite.Require().NotNil(err)
	err = suite.CheckRuleStatus(pk, chainID, address1, governance.GovernanceUnbinding)
	suite.Require().Nil(err)
}

//tc：应用链管理员更新unbinding状态的主验证规则，验证规则更新失败
//tc：应用链管理员更新unbinding状态的主验证规则，验证规则注销失败
func (suite Model7) Test0724_UpdateAndLogoutRuleWithUnbindingRuleIsSuccess() {
	address1, err := suite.DeploySimpleRule()
	suite.Require().Nil(err)
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	chainID := suite.GetChainID(pk)
	err = suite.InvokeRuleContract(pk, chainID, address1, RegisterRule)
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address1)
	suite.Require().Nil(err)
	err = suite.CheckRuleStatus(pk, chainID, address1, governance.GovernanceAvailable)
	suite.Require().Nil(err)
	address2, err := suite.DeploySimpleRule()
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, chainID, address2, RegisterRule)
	suite.Require().Nil(err)
	err = suite.CheckRuleStatus(pk, chainID, address2, governance.GovernanceBindable)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(chainID),
		rpcx.String(address2),
		rpcx.String("reason"),
	}
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), UpdateMasterRule, nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	err = suite.CheckRuleStatus(pk, chainID, address2, governance.GovernanceBinding)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, chainID, address2, RegisterRule)
	suite.Require().NotNil(err)
	err = suite.CheckRuleStatus(pk, chainID, address1, governance.GovernanceUnbinding)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, chainID, address1, UpdateMasterRule)
	suite.Require().NotNil(err)
	err = suite.InvokeRuleContract(pk, chainID, address1, LogoutRule)
	suite.Require().NotNil(err)
	err = suite.CheckRuleStatus(pk, chainID, address1, governance.GovernanceUnbinding)
	suite.Require().Nil(err)
}

//tc：应用链管理员更新forbidden状态的主验证规则，验证规则更新失败
//tc：应用链管理员更新forbidden状态的主验证规则，验证规则注销失败
func (suite Model7) Test0725_UpdateAndLogoutRuleWithForbiddenRuleIsFail() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	chainID := suite.GetChainID(pk)
	address1, err := suite.DeploySimpleRule()
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, chainID, address1, RegisterRule)
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address1)
	suite.Require().Nil(err)
	address2, err := suite.DeploySimpleRule()
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, chainID, address2, RegisterRule)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, chainID, address2, LogoutRule)
	suite.Require().Nil(err)
	err = suite.CheckRuleStatus(pk, chainID, address2, governance.GovernanceForbidden)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, chainID, address2, UpdateMasterRule)
	suite.Require().NotNil(err)
	err = suite.InvokeRuleContract(pk, chainID, address2, LogoutRule)
	suite.Require().NotNil(err)
	err = suite.CheckRuleStatus(pk, chainID, address2, governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

func (suite *Snake) InvokeRuleContract(pk crypto.PrivateKey, ChainID string, contractAddr string, method string) error {
	client := suite.NewClient(pk)
	var args []*pb.Arg
	if method == LogoutRule {
		args = []*pb.Arg{
			rpcx.String(ChainID),
			rpcx.String(contractAddr),
		}
	} else {
		args = []*pb.Arg{
			rpcx.String(ChainID),
			rpcx.String(contractAddr),
			rpcx.String("reason"),
		}
	}

	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), method, nil, args...)
	if err != nil {
		return err
	}
	if res.Status == pb.Receipt_FAILED {
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

func (suite *Snake) InvokeRuleContractWithReject(pk crypto.PrivateKey, ChainID string, contractAddr *types.Address, method string) error {
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.String(contractAddr.String()),
		rpcx.String("reason"),
	}
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), method, nil, args...)
	if err != nil {
		return err
	}
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	if res.Status == pb.Receipt_FAILED {
		return errors.New(string(res.Ret))
	}
	if result.ProposalID == "" {
		return nil
	}
	err = suite.VoteReject(result.ProposalID)
	if err != nil {
		return err
	}
	return nil
}

func (suite Snake) RegisterRule() (crypto.PrivateKey, string, string, error) {
	address, err := suite.DeploySimpleRule()
	if err != nil {
		return nil, "", "", err
	}
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	if err != nil {
		return nil, "", "", err
	}
	chainID := suite.GetChainID(pk)
	err = suite.InvokeRuleContract(pk, chainID, address, RegisterRule)
	if err != nil {
		return nil, "", "", err
	}
	return pk, chainID, address, nil
}
