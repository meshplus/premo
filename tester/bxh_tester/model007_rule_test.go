package bxh_tester

import (
	"encoding/json"
	"io/ioutil"

	"github.com/looplab/fsm"
	"github.com/meshplus/bitxhub-core/governance"

	"github.com/meshplus/bitxhub-kit/crypto"

	"github.com/meshplus/bitxhub-kit/types"

	"github.com/pkg/errors"

	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
)

const (
	RegisterRule = "RegisterRule"
	BindRule     = "BindRule"
	UnbindRule   = "UnbindRule"
	FreezeRule   = "FreezeRule"
	ActivateRule = "ActivateRule"
	LogoutRule   = "LogoutRule"
)

type Rule struct {
	Address string                      `json:"address"`
	ChainId string                      `json:"chain_id"`
	Status  governance.GovernanceStatus `json:"status"`
	FSM     *fsm.FSM                    `json:"fsm"`
}

type Model7 struct {
	*Snake
}

//tc：正确部署验证规则
func (suite *Model7) Test0701_DeployRule() {
	contract, err := ioutil.ReadFile("./testdata/simple_rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)
	suite.Require().NotNil(contractAddr)
}

//tc：部署验证规则为空
func (suite *Model7) Test0702_DeployRuleIsEmpty() {
	_, err := suite.client.DeployContract(nil, nil)
	suite.Require().NotNil(err)
}

//tc：验证规则未绑定，绑定验证规则
func (suite *Model7) Test0703_BindRule() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().Nil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceAvailable, status)
}

//tc：验证规则未绑定，绑定验证规则
func (suite *Model7) Test0704_BindRuleWithReject() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContractWithReject(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().Nil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceBindable, status)
}

//tc：验证规则binding状态，绑定验证规则
func (suite *Model7) Test0705_BindRuleWithBinding() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.String(contractAddr.String()),
	}
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), RegisterRule, nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().NotNil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceBinding, status)
}

//tc：验证规则available状态，绑定验证规则
func (suite *Model7) Test0706_BindRuleWithAvailable() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().NotNil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceAvailable, status)
}

//tc：验证规则unbinding状态，绑定验证规则
func (suite *Model7) Test0707_BindRuleWithUnbinding() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().Nil(err)

	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.String(contractAddr.String()),
	}
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), UnbindRule, nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().NotNil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceUnbinding, status)
}

//tc：验证规则bindable状态，绑定验证规则
func (suite *Model7) Test0708_BindRuleWithBindable() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, UnbindRule)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, BindRule)
	suite.Require().Nil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceAvailable, status)
}

//tc：验证规则logouting状态，绑定验证规则
func (suite *Model7) Test0709_BindRuleWithLogouting() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().Nil(err)

	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.String(contractAddr.String()),
	}
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), LogoutRule, nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().NotNil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceLogouting, status)
}

//tc：验证规则forbidden状态，绑定验证规则
func (suite *Model7) Test0710_BindRuleWithForbidden() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, LogoutRule)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().NotNil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceForbidden, status)
}

//tc：验证规则freezing状态，绑定验证规则
func (suite *Model7) Test0711_BindRuleWithFreezing() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().Nil(err)

	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.String(contractAddr.String()),
	}
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), FreezeRule, nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().NotNil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceFreezing, status)
}

//tc：验证规则frozen状态，绑定验证规则
func (suite *Model7) Test0712_BindRuleWithFrozen() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, FreezeRule)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().NotNil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceFrozen, status)
}

//tc：验证规则未绑定，解绑验证规则
func (suite *Model7) Test0713_UnbindRuleWithNoRule() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	address := types.NewAddressByStr("0x64C5334AadE6c623ae829422C34B6f310b031aa0")
	err = suite.InvokeRuleContract(pk, ChainID, address, UnbindRule)
	suite.Require().NotNil(err)
}

//tc：验证规则binding状态，解绑验证规则
func (suite *Model7) Test0714_UnbindRuleWithBinding() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.String(contractAddr.String()),
	}
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), RegisterRule, nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, UnbindRule)
	suite.Require().NotNil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceBinding, status)
}

//tc：验证规则available状态，解绑验证规则
func (suite *Model7) Test0715_UnbindRule() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, UnbindRule)
	suite.Require().Nil(err)
}

//tc：验证规则available状态，解绑验证规则
func (suite *Model7) Test0716_UnbindRuleWithReject() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContractWithReject(pk, ChainID, contractAddr, UnbindRule)
	suite.Require().Nil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceAvailable, status)
}

//tc：验证规则unbinding状态，解绑验证规则
func (suite *Model7) Test0717_UnbindRuleWithUnbinding() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().Nil(err)

	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.String(contractAddr.String()),
	}
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), UnbindRule, nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, UnbindRule)
	suite.Require().NotNil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceUnbinding, status)
}

//tc：验证规则bindable状态，解绑验证规则
func (suite *Model7) Test0718_UnbindRuleWithBindable() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, UnbindRule)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, UnbindRule)
	suite.Require().NotNil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceBindable, status)
}

//tc：验证规则logouting状态，解绑验证规则
func (suite *Model7) Test0719_UnbindRuleWithLogouting() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().Nil(err)

	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.String(contractAddr.String()),
	}
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), LogoutRule, nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, UnbindRule)
	suite.Require().NotNil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceLogouting, status)
}

//tc：验证规则forbidden状态，解绑验证规则
func (suite *Model7) Test0720_UnbindRuleWithForbidden() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, LogoutRule)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, UnbindRule)
	suite.Require().NotNil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceForbidden, status)
}

//tc：验证规则freezing状态，解绑验证规则
func (suite *Model7) Test0721_UnbindRuleWithFreezing() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().Nil(err)

	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.String(contractAddr.String()),
	}
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), FreezeRule, nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, UnbindRule)
	suite.Require().NotNil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceFreezing, status)
}

//tc：验证规则frozen状态，解绑验证规则
func (suite *Model7) Test0722_UnbindRuleWithFrozen() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, FreezeRule)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, UnbindRule)
	suite.Require().NotNil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceFrozen, status)
}

//tc：验证规则未绑定，冻结验证规则
func (suite *Model7) Test0723_FreezeRuleWithNoRule() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	address := types.NewAddressByStr("0x64C5334AadE6c623ae829422C34B6f310b031aa0")
	err = suite.InvokeRuleContract(pk, ChainID, address, FreezeRule)
	suite.Require().NotNil(err)
}

//tc：验证规则binding状态，冻结验证规则
func (suite *Model7) Test0724_FreezeRuleWithBinding() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.String(contractAddr.String()),
	}
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), RegisterRule, nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, FreezeRule)
	suite.Require().Nil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceFrozen, status)
}

//tc：验证规则available状态，冻结验证规则
func (suite *Model7) Test0725_FreezeRule() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, FreezeRule)
	suite.Require().Nil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceFrozen, status)
}

//tc：验证规则available状态，冻结验证规则
func (suite *Model7) Test0726_FreezeRuleWithReject() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContractWithReject(pk, ChainID, contractAddr, FreezeRule)
	suite.Require().Nil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceAvailable, status)
}

//tc：验证规则unbinding状态，冻结验证规则
func (suite *Model7) Test0727_FreezeRuleWithUnbinding() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().Nil(err)

	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.String(contractAddr.String()),
	}
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), UnbindRule, nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, FreezeRule)
	suite.Require().Nil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceFrozen, status)
}

//tc：验证规则bindable状态，冻结验证规则
func (suite *Model7) Test0728_FreezeRuleWithBindable() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, UnbindRule)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, FreezeRule)
	suite.Require().Nil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceFrozen, status)
}

//tc：验证规则logouting状态，冻结验证规则
func (suite *Model7) Test0729_FreezeRuleWithLogouting() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().Nil(err)

	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.String(contractAddr.String()),
	}
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), LogoutRule, nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, FreezeRule)
	suite.Require().NotNil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceLogouting, status)
}

//tc：验证规则forbidden状态，冻结验证规则
func (suite *Model7) Test0730_FreezeRuleWithForbidden() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, LogoutRule)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, FreezeRule)
	suite.Require().NotNil(err)
}

//tc：验证规则freezing状态，解冻验证规则
func (suite *Model7) Test0731_FreezeRuleWithFreezing() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().Nil(err)

	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.String(contractAddr.String()),
	}
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), FreezeRule, nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, FreezeRule)
	suite.Require().NotNil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceFreezing, status)
}

//tc：验证规则frozen状态，解冻验证规则
func (suite *Model7) Test0732_FreezeRuleWithFrozen() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, FreezeRule)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, FreezeRule)
	suite.Require().NotNil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceFrozen, status)
}

//tc：验证规则未绑定，解冻验证规则
func (suite *Model7) Test0733_ActivateRuleWithNoRule() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	address := types.NewAddressByStr("0x64C5334AadE6c623ae829422C34B6f310b031aa0")
	err = suite.InvokeRuleContract(pk, ChainID, address, ActivateRule)
	suite.Require().NotNil(err)
}

//tc：验证规则binding状态，解冻验证规则
func (suite *Model7) Test0734_ActivateRuleWithBinding() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.String(contractAddr.String()),
	}
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), RegisterRule, nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, ActivateRule)
	suite.Require().NotNil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceBinding, status)
}

//tc：验证规则available状态，解冻验证规则
func (suite *Model7) Test0735_ActivateRuleWithAvailable() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, ActivateRule)
	suite.Require().NotNil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceAvailable, status)
}

//tc：验证规则unbinding状态，解冻验证规则
func (suite *Model7) Test0736_ActivateRuleWithUnbinding() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().Nil(err)

	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.String(contractAddr.String()),
	}
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), UnbindRule, nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, ActivateRule)
	suite.Require().NotNil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceUnbinding, status)
}

//tc：验证规则bindable状态，解冻验证规则
func (suite *Model7) Test0737_ActivateRuleWithBindable() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, UnbindRule)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, ActivateRule)
	suite.Require().NotNil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceBindable, status)
}

//tc：验证规则logouting状态，解冻验证规则
func (suite *Model7) Test0738_ActivateRuleWithLogouting() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().Nil(err)

	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.String(contractAddr.String()),
	}
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), LogoutRule, nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, ActivateRule)
	suite.Require().NotNil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceLogouting, status)
}

//tc：验证规则forbidden状态，解冻验证规则
func (suite *Model7) Test0739_ActivateRuleWithForbidden() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, LogoutRule)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, ActivateRule)
	suite.Require().NotNil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceForbidden, status)
}

//tc：验证规则freezing状态，解冻验证规则
func (suite *Model7) Test0740_ActivateRuleWithFreezing() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().Nil(err)

	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.String(contractAddr.String()),
	}
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), FreezeRule, nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, ActivateRule)
	suite.Require().NotNil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceFreezing, status)
}

//tc：验证规则frozen状态，解冻验证规则
func (suite *Model7) Test0741_ActivateRuleWithFrozen() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, FreezeRule)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, ActivateRule)
	suite.Require().Nil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceBindable, status)
}

//tc：验证规则未绑定，注销验证规则
func (suite *Model7) Test0742_LogoutRuleWithNoRule() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	address := types.NewAddressByStr("0x64C5334AadE6c623ae829422C34B6f310b031aa0")
	err = suite.InvokeRuleContract(pk, ChainID, address, LogoutRule)
	suite.Require().NotNil(err)
}

//tc：验证规则binding状态，注销验证规则
func (suite *Model7) Test0743_LogoutRuleWithBinding() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.String(contractAddr.String()),
	}
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), RegisterRule, nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, LogoutRule)
	suite.Require().Nil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceForbidden, status)
}

//tc：验证规则available状态，注销验证规则
func (suite *Model7) Test0744_LogoutRuleWithAvailable() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, LogoutRule)
	suite.Require().Nil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceForbidden, status)
}

//tc：验证规则unbinding状态，注销验证规则
func (suite *Model7) Test0745_LogoutRuleWithUnbinding() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().Nil(err)

	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.String(contractAddr.String()),
	}
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), UnbindRule, nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, LogoutRule)
	suite.Require().Nil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceForbidden, status)
}

//tc：验证规则bindable状态，注销验证规则
func (suite *Model7) Test0746_LogoutRuleWithBindable() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, UnbindRule)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, LogoutRule)
	suite.Require().Nil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceForbidden, status)
}

//tc：验证规则logouting状态，注销验证规则
func (suite *Model7) Test0747_LogoutRuleWithLogouting() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().Nil(err)

	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.String(contractAddr.String()),
	}
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), LogoutRule, nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, LogoutRule)
	suite.Require().NotNil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceLogouting, status)
}

//tc：验证规则forbidden状态，注销验证规则
func (suite *Model7) Test0748_LogoutRuleWithForbidden() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, LogoutRule)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, LogoutRule)
	suite.Require().NotNil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceForbidden, status)
}

//tc：验证规则freezing状态，注销验证规则
func (suite *Model7) Test0749_LogoutRuleWithFreezing() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().Nil(err)

	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.String(contractAddr.String()),
	}
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), FreezeRule, nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, LogoutRule)
	suite.Require().Nil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceForbidden, status)
}

//tc：验证规则frozen状态，注销验证规则
func (suite *Model7) Test0750_LogoutRuleWithFrozen() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, FreezeRule)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, LogoutRule)
	suite.Require().Nil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceForbidden, status)
}

func (suite *Snake) InvokeRuleContract(pk crypto.PrivateKey, ChainID string, contractAddr *types.Address, method string) error {
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.String(contractAddr.String()),
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
	err = suite.VoteReject(result.ProposalID)
	if err != nil {
		return err
	}
	return nil
}

func (suite *Snake) getRuleStatus(pk crypto.PrivateKey, ChainID string, contractAddr *types.Address) (status governance.GovernanceStatus, err error) {
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.String(contractAddr.String()),
	}
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), "GetRuleByAddr", nil, args...)
	if err != nil {
		return "", err
	}
	if res.Status == pb.Receipt_FAILED {
		return "", errors.New(string(res.Ret))
	}
	rule := &Rule{}
	err = json.Unmarshal(res.Ret, rule)
	if err != nil {
		return "", err
	}
	return rule.Status, nil
}
