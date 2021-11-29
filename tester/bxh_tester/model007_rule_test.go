package bxh_tester

import (
	"encoding/json"
	"fmt"

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

//tc：注册Fabric V1.4.3类型的应用链，默认验证规则注册成功
func (suite Model7) Test0703_RegisterDefaultRuleWithFabricV143() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchainWithType(pk, "Fabric V1.4.3", address, "{\"channel_id\":\"mychannel\",\"chaincode_id\":\"broker\",\"broker_version\":\"1\"}")
	suite.Require().Nil(err)
	suite.Require().Equal(true, suite.RuleContains(from, SimFabricRuleAddr))
	suite.Require().Equal(true, suite.RuleContains(from, FabricRuleAddr))
	suite.Require().Equal(true, suite.RuleContains(from, HappyRuleAddr))
}

//tc：注册Fabric V1.4.4类型的应用链，默认验证规则注册成功
func (suite Model7) Test0704_RegisterDefaultRuleWithFabricV144() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchainWithType(pk, "Fabric V1.4.4", address, "{\"channel_id\":\"mychannel\",\"chaincode_id\":\"broker\",\"broker_version\":\"1\"}")
	suite.Require().Nil(err)
	suite.Require().Equal(true, suite.RuleContains(from, SimFabricRuleAddr))
	suite.Require().Equal(true, suite.RuleContains(from, FabricRuleAddr))
	suite.Require().Equal(true, suite.RuleContains(from, HappyRuleAddr))
}

//tc：注册Hyperchain V1.8.3类型的应用链，默认验证规则注册成功
func (suite Model7) Test0705_RegisterDefaultRuleWithHyperchainV183() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchainWithType(pk, "Hyperchain V1.8.3", address, "0x857133c5C69e6Ce66F7AD46F200B9B3573e77582")
	suite.Require().Nil(err)
	suite.Require().Equal(true, suite.RuleContains(from, HappyRuleAddr))
}

////tc：注册Hyperchain V1.8.6类型的应用链，默认验证规则注册成功
func (suite Model7) Test0706_RegisterDefaultRuleWithHyperchainV186() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchainWithType(pk, "Hyperchain V1.8.6", address, "0x857133c5C69e6Ce66F7AD46F200B9B3573e77582")
	suite.Require().Nil(err)
	suite.Require().Equal(true, suite.RuleContains(from, HappyRuleAddr))
}

//tc：注册Flato V1.0.0类型的应用链，默认验证规则注册成功
func (suite Model7) Test0707_RegisterDefaultRuleWithFlatoV100() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchainWithType(pk, "Flato V1.0.0", address, "0x857133c5C69e6Ce66F7AD46F200B9B3573e77582")
	suite.Require().Nil(err)
	suite.Require().Equal(true, suite.RuleContains(from, HappyRuleAddr))
}

//tc：注册Flato V1.0.3类型的应用链，默认验证规则注册成功
func (suite Model7) Test0708_RegisterDefaultRuleWithFlatoV103() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchainWithType(pk, "Flato V1.0.3", address, "0x857133c5C69e6Ce66F7AD46F200B9B3573e77582")
	suite.Require().Nil(err)
	suite.Require().Equal(true, suite.RuleContains(from, HappyRuleAddr))
}

//tc：注册Flato V1.0.6类型的应用链，默认验证规则注册成功
func (suite Model7) Test0709_RegisterDefaultRuleWithFlatoV106() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchainWithType(pk, "Flato V1.0.6", address, "0x857133c5C69e6Ce66F7AD46F200B9B3573e77582")
	suite.Require().Nil(err)
	suite.Require().Equal(true, suite.RuleContains(from, HappyRuleAddr))
}

//tc：注册BCOS V2.6.0类型的应用链，默认验证规则注册成功
func (suite Model7) Test0710_RegisterDefaultRuleWithBCOSV260() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchainWithType(pk, "BCOS V2.6.0", address, "0x857133c5C69e6Ce66F7AD46F200B9B3573e77582")
	suite.Require().Nil(err)
	suite.Require().Equal(true, suite.RuleContains(from, HappyRuleAddr))
}

//tc：注册CITA V20.2.2类型的应用链，默认验证规则注册成功
func (suite Model7) Test0711_RegisterDefaultRuleWithCITAV2022() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchainWithType(pk, "CITA V20.2.2", address, "0x857133c5C69e6Ce66F7AD46F200B9B3573e77582")
	suite.Require().Nil(err)
	suite.Require().Equal(true, suite.RuleContains(from, HappyRuleAddr))
}

//tc：注册ETH类型的应用链，默认验证规则注册成功
func (suite Model7) Test0712_RegisterDefaultRuleWithETH() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchainWithType(pk, "ETH", address, "0x857133c5C69e6Ce66F7AD46F200B9B3573e77582")
	suite.Require().Nil(err)
	suite.Require().Equal(true, suite.RuleContains(from, HappyRuleAddr))
}

//tc：注册其他类型的应用链，默认验证规则注册成功
func (suite Model7) Test0713_RegisterDefaultRuleWithOthers() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchainWithType(pk, "Other", address, "0x857133c5C69e6Ce66F7AD46F200B9B3573e77582")
	suite.Require().Nil(err)
	suite.Require().Equal(true, suite.RuleContains(from, HappyRuleAddr))
}

//tc：非应用链管理员调用注册验证规则，验证规则注册失败
func (suite Model7) Test0714_RegisterRuleWithNoAdminIsFail() {
	pk1, from, address1, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchainWithType(pk1, "Hyperchain V1.8.3", address1, "0x857133c5C69e6Ce66F7AD46F200B9B3573e77582")
	suite.Require().Nil(err)
	pk2, _, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk2, from, address2, RegisterRule)
	suite.Require().NotNil(err)
}

//tc：应用链管理员调用注册验证规则，验证规则注册成功
func (suite Model7) Test0715_RegisterRuleIsSuccess() {
	pk, from, address1, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchainWithType(pk, "Hyperchain V1.8.3", address1, "0x857133c5C69e6Ce66F7AD46F200B9B3573e77582")
	suite.Require().Nil(err)
	_, _, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, from, address2, RegisterRule)
	suite.Require().Nil(err)
}

//tc：应用链管理员注册未部署的验证规则，验证规则部署失败
func (suite Model7) Test0716_RegisterRuleWithNoRegisterRuleIsFail() {
	pk, from, address1, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchainWithType(pk, "Hyperchain V1.8.3", address1, "0x857133c5C69e6Ce66F7AD46F200B9B3573e77582")
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, from, "0x000000000000000000000000000000000000001", RegisterRule)
	suite.Require().NotNil(err)
}

//tc：应用链管理员注册available的验证规则，验证规则部署失败
func (suite Model7) Test0717_RegisterRuleWithAvailableRuleIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchainWithType(pk, "Hyperchain V1.8.3", address, "0x857133c5C69e6Ce66F7AD46F200B9B3573e77582")
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, from, address, RegisterRule)
	suite.Require().NotNil(err)
}

//tc：应用链管理员注册binding的验证规则，验证规则注册失败
//tc：应用链管理员注册unbinding的验证规则，验证规则注册失败
func (suite Model7) Test0718_RegisterRuleWithBindingRuleIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchainWithType(pk, "Hyperchain V1.8.3", address, "0x857133c5C69e6Ce66F7AD46F200B9B3573e77582")
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(from),
		rpcx.String(HappyRuleAddr),
		rpcx.String("reason"),
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), UpdateMasterRule, nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	err = suite.CheckRuleStatus(pk, from, HappyRuleAddr, governance.GovernanceBinding)
	suite.Require().Nil(err)
	err = suite.CheckRuleStatus(pk, from, address, governance.GovernanceUnbinding)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, from, address, RegisterRule)
	suite.Require().NotNil(err)
	err = suite.InvokeRuleContract(pk, from, HappyRuleAddr, RegisterRule)
	suite.Require().NotNil(err)
}

//tc：应用链管理员注册bindable的验证规则，验证规则部署失败
func (suite Model7) Test0719_RegisterRuleWithBindableRuleIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchainWithType(pk, "Hyperchain V1.8.3", address, "0x857133c5C69e6Ce66F7AD46F200B9B3573e77582")
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, from, HappyRuleAddr, RegisterRule)
	suite.Require().NotNil(err)
}

//tc：应用链管理员注册forbidden的验证规则，验证规则部署失败
func (suite Model7) Test0720_RegisterRuleWithForbiddenRuleIsFail() {
	pk, from, address1, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchainWithType(pk, "Hyperchain V1.8.3", address1, "0x857133c5C69e6Ce66F7AD46F200B9B3573e77582")
	suite.Require().Nil(err)
	_, _, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, from, address2, RegisterRule)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, from, address2, LogoutRule)
	suite.Require().Nil(err)
	err = suite.CheckRuleStatus(pk, from, address2, governance.GovernanceForbidden)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, from, address2, RegisterRule)
	suite.Require().NotNil(err)
}

//tc：应用链不存在注册验证规则，验证规则注册失败
func (suite Model7) Test0721_RegisterRuleWithUnRegisteredChainIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, from, address, RegisterRule)
	suite.Require().NotNil(err)
}

//tc：应用链处于updating注册验证规则，验证规则注册成功
func (suite Model7) Test0722_RegisterRuleWithUpdatingChainIsSuccess() {
	pk, from, address1, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToUpdating(pk, from, address1)
	suite.Require().Nil(err)
	_, _, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, from, address2, RegisterRule)
	suite.Require().Nil(err)
}

//tc：应用链处于activating注册验证规则，验证规则注册成功
func (suite Snake) Test0723_RegisterRuleWithActivatingChainIsSuccess() {
	pk, from, address1, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToActivating(pk, from, address1)
	suite.Require().Nil(err)
	_, _, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, from, address2, RegisterRule)
	suite.Require().Nil(err)
}

//tc：应用链处于freezing注册验证规则，验证规则注册成功
func (suite Snake) Test0724_RegisterRuleWithFreezingChainIsSuccess() {
	pk, from, address1, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToFreezing(pk, from, address1)
	suite.Require().Nil(err)
	_, _, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, from, address2, RegisterRule)
	suite.Require().Nil(err)
}

//tc：应用链处于Frozen注册验证规则，验证规则注册成功
func (suite Snake) Test0725_RegisterRuleWithFrozenChainIsSuccess() {
	pk, from, address1, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToFrozen(pk, from, address1)
	suite.Require().Nil(err)
	_, _, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, from, address2, RegisterRule)
	suite.Require().Nil(err)
}

//tc：应用链处于Logouting注册验证规则，验证规则注册成功
func (suite Snake) Test0726_RegisterRuleWithLogoutingChainIsSuccess() {
	pk, from, address1, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToLogouting(pk, from, address1)
	suite.Require().Nil(err)
	_, _, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, from, address2, RegisterRule)
	suite.Require().Nil(err)
}

//tc：应用链处于Forbidden注册验证规则，验证规则注册成功
func (suite Snake) Test0726_RegisterRuleWithForbiddenChainIsSuccess() {
	pk, from, address1, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToForbidden(pk, from, address1)
	suite.Require().Nil(err)
	_, _, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, from, address2, RegisterRule)
	suite.Require().Nil(err)
}

//tc：非应用链管理员更新主验证规则，验证规则更新失败
//tc：非应用链管理员注销主验证规则，验证规则注销失败
func (suite Model7) Test0712_UpdateAndLogoutRuleWithNoAdminIsSuccess() {
	pk1, from, address1, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchainWithType(pk1, "Hyperchain V1.8.3", address1, "0x857133c5C69e6Ce66F7AD46F200B9B3573e77582")
	suite.Require().Nil(err)
	pk2, _, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk1, from, address2, RegisterRule)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk2, from, address2, UpdateMasterRule)
	suite.Require().NotNil(err)
	err = suite.InvokeRuleContract(pk2, from, address2, LogoutRule)
	suite.Require().NotNil(err)
}

//tc：应用链管理员更新主验证规则，验证规则更新成功
//tc：应用链管理员更新主验证规则，验证规则注销成功
func (suite Model7) Test0713_UpdateAndLogoutRuleIsSuccess() {
	pk, from, address1, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchainWithType(pk, "Hyperchain V1.8.3", address1, "0x857133c5C69e6Ce66F7AD46F200B9B3573e77582")
	suite.Require().Nil(err)
	_, _, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, from, address2, RegisterRule)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, from, address2, UpdateMasterRule)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, from, address1, LogoutRule)
	suite.Require().Nil(err)
}

//tc:应用链处于未注册状态更新主验证规则，验证规则更新失败
//tc:应用链处于未注册状态更新主注销规则，验证规则注销失败
func (suite Model7) Test0714_UpdateAndLogoutRuleWithNoRegisterAdmin() {
	pk, from, address1, err := suite.DeployRule()
	suite.Require().Nil(err)
	_, _, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, from, address2, UpdateMasterRule)
	suite.Require().NotNil(err)
	err = suite.InvokeRuleContract(pk, from, address1, LogoutRule)
	suite.Require().NotNil(err)
}

//tc：应用链处于activating状态更新主验证规则，验证规则更新失败
//tc：应用链处于activating状态注销验证规则，验证规则注销成功
func (suite Model7) Test0717_UpdateAndLogoutRuleWithActivatingChain() {
	pk, from, address1, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToActivating(pk, from, address1)
	suite.Require().Nil(err)
	_, _, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, from, address2, RegisterRule)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, from, address2, UpdateMasterRule)
	suite.Require().NotNil(err)
	err = suite.InvokeRuleContract(pk, from, address2, LogoutRule)
	suite.Require().Nil(err)
}

//tc：应用链处于freezing状态更新主验证规则，验证规则更新失败
//tc：应用链处于freezing状态注销验证规则，验证规则注销成功
func (suite Model7) Test0717_UpdateAndLogoutRuleWithFreezingChain() {
	pk, from, address1, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToFreezing(pk, from, address1)
	suite.Require().Nil(err)
	_, _, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, from, address2, RegisterRule)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, from, address2, UpdateMasterRule)
	suite.Require().NotNil(err)
	err = suite.InvokeRuleContract(pk, from, address2, LogoutRule)
	suite.Require().Nil(err)
}

//tc：应用链处于frozen状态更新主验证规则，验证规则更新成功
//tc：应用链处于frozen状态注销验证规则，验证规则注销成功
func (suite Model7) Test0718_UpdateAndLogoutRuleWithFrozenChain() {
	pk, from, address1, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToFrozen(pk, from, address1)
	suite.Require().Nil(err)
	_, _, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, from, address2, RegisterRule)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, from, address2, UpdateMasterRule)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, from, address1, LogoutRule)
	suite.Require().Nil(err)
}

//tc：应用链处于logouting状态更新主验证规则，验证规则更新失败
//tc：应用链处于logouting状态注销验证规则，验证规则注销成功
func (suite Model7) Test0719_UpdateAndLogoutRuleWithLogoutingChainIsFail() {
	pk, from, address1, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToLogouting(pk, from, address1)
	suite.Require().Nil(err)
	_, _, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, from, address2, RegisterRule)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, from, address2, UpdateMasterRule)
	suite.Require().NotNil(err)
	err = suite.InvokeRuleContract(pk, from, address2, LogoutRule)
	suite.Require().Nil(err)
}

//tc：应用链处于forbidden状态更新主验证规则，验证规则更新失败
//tc：应用链处于forbidden状态注销验证规则，验证规则注销失败
func (suite Model7) Test0720_UpdateAndLogoutRuleWithForbiddenChainIsFail() {
	pk, from, address1, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToForbidden(pk, from, address1)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, from, HappyRuleAddr, UpdateMasterRule)
	suite.Require().NotNil(err)
	err = suite.InvokeRuleContract(pk, from, HappyRuleAddr, LogoutRule)
	suite.Require().NotNil(err)
}

//tc：应用链更新未注册的主验证规则，验证规则更新失败
//tc：应用链更新未注册的主验证规则，验证规则注销失败
func (suite Model7) Test0721_UpdateAndLogoutRuleWithNoRegisterRuleIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchainWithType(pk, "Hyperchain V1.8.3", address, "0x857133c5C69e6Ce66F7AD46F200B9B3573e77582")
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, from, "0x000000000000000000000000000000000000001", UpdateMasterRule)
	suite.Require().NotNil(err)
	err = suite.InvokeRuleContract(pk, from, "0x000000000000000000000000000000000000001", LogoutRule)
	suite.Require().NotNil(err)
}

//tc：应用链管理员更新available状态的主验证规则，验证规则更新失败
//tc：应用链管理员更新available状态的主验证规则，验证规则注销失败
func (suite Model7) Test0722_UpdateAndLogoutRuleWithAvailableRuleIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchainWithType(pk, "Hyperchain V1.8.3", address, "0x857133c5C69e6Ce66F7AD46F200B9B3573e77582")
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, from, address, UpdateMasterRule)
	suite.Require().NotNil(err)
	err = suite.InvokeRuleContract(pk, from, address, LogoutRule)
	suite.Require().NotNil(err)
}

//tc：应用链管理员更新binding状态的主验证规则，验证规则更新失败
//tc：应用链管理员更新binding状态的主验证规则，验证规则注销失败
func (suite Model7) Test0723_UpdateAndLogoutRuleWithBindingRuleIsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchainWithType(pk, "Hyperchain V1.8.3", address, "0x857133c5C69e6Ce66F7AD46F200B9B3573e77582")
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(from),
		rpcx.String(HappyRuleAddr),
		rpcx.String("reason"),
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), UpdateMasterRule, nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	err = suite.CheckRuleStatus(pk, from, HappyRuleAddr, governance.GovernanceBinding)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, from, HappyRuleAddr, UpdateMasterRule)
	suite.Require().NotNil(err)
	err = suite.InvokeRuleContract(pk, from, HappyRuleAddr, LogoutRule)
	suite.Require().NotNil(err)
}

//tc：应用链管理员更新unbinding状态的主验证规则，验证规则更新失败
//tc：应用链管理员更新unbinding状态的主验证规则，验证规则注销失败
func (suite Model7) Test0724_UpdateAndLogoutRuleWithUnbindingRuleIsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchainWithType(pk, "Hyperchain V1.8.3", address, "0x857133c5C69e6Ce66F7AD46F200B9B3573e77582")
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(from),
		rpcx.String(HappyRuleAddr),
		rpcx.String("reason"),
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), UpdateMasterRule, nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	err = suite.CheckRuleStatus(pk, from, address, governance.GovernanceUnbinding)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, from, address, UpdateMasterRule)
	suite.Require().NotNil(err)
	err = suite.InvokeRuleContract(pk, from, address, LogoutRule)
	suite.Require().NotNil(err)
}

//tc：应用链管理员更新forbidden状态的主验证规则，验证规则更新失败
//tc：应用链管理员更新forbidden状态的主验证规则，验证规则注销失败
func (suite Model7) Test0725_UpdateAndLogoutRuleWithForbiddenRuleIsFail() {
	pk, from, address1, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchainWithType(pk, "Hyperchain V1.8.3", address1, "0x857133c5C69e6Ce66F7AD46F200B9B3573e77582")
	suite.Require().Nil(err)
	_, _, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, from, address2, RegisterRule)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, from, address2, LogoutRule)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, from, address2, UpdateMasterRule)
	suite.Require().NotNil(err)
	err = suite.InvokeRuleContract(pk, from, address2, LogoutRule)
	suite.Require().NotNil(err)
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

func (suite Snake) DeployRule() (crypto.PrivateKey, string, string, error) {
	address, err := suite.DeploySimpleRule()
	if err != nil {
		return nil, "", "", err
	}
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	if err != nil {
		return nil, "", "", err
	}
	from, err := pk.PublicKey().Address()
	if err != nil {
		return nil, "", "", err
	}
	return pk, from.String(), address, nil
}

func (suite Snake) RegisterAppchainWithType(pk crypto.PrivateKey, typ, address, broker string) error {
	client := suite.NewClient(pk)
	from, err := pk.PublicKey().Address()
	if err != nil {
		return err
	}
	args := []*pb.Arg{
		rpcx.String(from.String()),        //chainID
		rpcx.String(from.String()),        //chainName
		rpcx.String(typ),                  //chainType
		rpcx.Bytes([]byte("")),            //trustRoot
		rpcx.String(broker),               //broker
		rpcx.String("desc"),               //desc
		rpcx.String(address),              //masterRuleAddr
		rpcx.String("https://github.com"), //masterRuleUrl
		rpcx.String(from.String()),        //adminAddrs
		rpcx.String("reason"),             //reason
	}
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "RegisterAppchain", nil, args...)
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

func (suite Snake) Rules(chainID string) ([]Rule, error) {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	if err != nil {
		return nil, err
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), "Rules", nil, rpcx.String(chainID))
	if err != nil {
		return nil, err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return nil, errors.New(string(res.Ret))
	}
	var rules []Rule
	err = json.Unmarshal(res.Ret, &rules)
	if err != nil {
		return nil, err
	}
	return rules, nil
}

func (suite Snake) RuleContains(chainID, address string) bool {
	rules, err := suite.Rules(chainID)
	if err != nil {
		return false
	}
	for i := 0; i < len(rules); i++ {
		if rules[i].Address == address {
			return true
		}
	}
	return false
}
