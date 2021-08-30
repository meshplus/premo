package bxh_tester

import (
	"encoding/json"
	"io/ioutil"

	"github.com/looplab/fsm"
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

type Rule struct {
	Address string                      `json:"address"`
	ChainId string                      `json:"chain_id"`
	Status  governance.GovernanceStatus `json:"status"`
	FSM     *fsm.FSM                    `json:"fsm"`
}

type Model7 struct {
	*Snake
}

func (suite *Model7) SetupTest() {
	suite.T().Parallel()
}

//tc：正确部署验证规则,并返回地址
func (suite *Model7) Test0701_DeployRule() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	contract, err := ioutil.ReadFile("./testdata/simple_rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := client.DeployContract(contract, nil)
	suite.Require().Nil(err)
	suite.Require().NotNil(contractAddr)
}

//tc：部署验证规则为空，并提示错误信息
func (suite *Model7) Test0702_DeployRuleIsEmpty() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	_, err = client.DeployContract(nil, nil)
	suite.Require().NotNil(err)
}

//tc：验证规则不存在，注册验证规则
func (suite *Model7) Test0703_BindRuleWithNoRule() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	address := types.NewAddressByStr("0x64C5334AadE6c623ae829422C34B6f310b031aa0")

	err = suite.InvokeRuleContract(pk, ChainID, address, RegisterRule)
	suite.Require().NotNil(err)
}

//tc:发起正确绑定规则的请求，中继链管理员投票通过，验证规则状态为available
func (suite *Model7) Test0704_BindRule() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().Nil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceAvailable, status)
}

//tc:发起正确绑定规则的请求，中继链管理员投票不通过，验证规则状态为bindable
func (suite *Model7) Test0705_BindRuleWithReject() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContractWithReject(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().Nil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceBindable, status)
}

//tc:发起正确绑定规则的请求，中继链管理员投票过程中，验证规则状态为binding
func (suite *Model7) Test0706_BindRuleWithDoing() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	args := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.String(contractAddr.String()),
		rpcx.String("reason"),
	}
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), RegisterRule, nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceBinding, status)
}

//tc:验证规则状态为binding，发起绑定请求，提示对应错误信息
func (suite *Model7) Test0707_BindRuleWithBinding() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	args := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.String(contractAddr.String()),
		rpcx.String("reason"),
	}
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), RegisterRule, nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().NotNil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceBinding, status)
}

//tc:验证规则状态为bindable，发起绑定请求，提示对应错误信息
func (suite *Model7) Test0708_BindRuleWithBindable() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr1, err := client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	contractAddr2, err := client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContractWithReject(pk, ChainID, contractAddr1, RegisterRule)
	suite.Require().Nil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr1)
	suite.Require().Equal(governance.GovernanceBindable, status)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr1, RegisterRule)
	suite.Require().NotNil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr2, RegisterRule)
	suite.Require().Nil(err)
}

//tc:验证规则状态为available，发起绑定请求，提示对应错误信息
func (suite *Model7) Test0709_BindRuleWithAvailable() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().NotNil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceAvailable, status)
}

//tc:验证规则状态为forbidden，发起绑定请求，提示对应错误信息
func (suite *Model7) Test0710_BindRuleWithForbidden() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr1, err := client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	contractAddr2, err := client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr1, RegisterRule)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr2, RegisterRule)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr2, LogoutRule)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr2, RegisterRule)
	suite.Require().NotNil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr1)
	suite.Require().Equal(governance.GovernanceAvailable, status)
	status, err = suite.getRuleStatus(pk, ChainID, contractAddr2)
	suite.Require().Equal(governance.GovernanceForbidden, status)
}

//tc：验证规则不存在，更新验证规则
func (suite *Model7) Test0711_UpdateRuleWithNoRule() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr1, err := client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	contractAddr2, err := client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr1, RegisterRule)
	suite.Require().Nil(err)

	args := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.String(contractAddr2.String()),
		rpcx.String("reason"),
	}
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), UpdateMasterRule, nil, args...)
	suite.Require().Nil(err)
	suite.Require().Contains(string(res.Ret), "the rule does not exist")
}

//tc:发起正确更新规则的请求，中继链管理员投票通过，验证规则状态为available
func (suite *Model7) Test0712_UpdateRule() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr1, err := client.DeployContract(contract, nil)
	suite.Require().Nil(err)
	contractAddr2, err := client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr1, RegisterRule)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, ChainID, contractAddr2, RegisterRule)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, ChainID, contractAddr2, UpdateMasterRule)
	suite.Require().Nil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr1)
	suite.Require().Equal(governance.GovernanceBindable, status)
	status, err = suite.getRuleStatus(pk, ChainID, contractAddr2)
	suite.Require().Equal(governance.GovernanceAvailable, status)
}

//tc:发起正确更新规则的请求，中继链管理员投票不通过，验证规则状态为bindable
func (suite *Model7) Test0713_UpdateRuleWithReject() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr1, err := client.DeployContract(contract, nil)
	suite.Require().Nil(err)
	contractAddr2, err := client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr1, RegisterRule)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, ChainID, contractAddr2, RegisterRule)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContractWithReject(pk, ChainID, contractAddr2, UpdateMasterRule)
	suite.Require().Nil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr1)
	suite.Require().Equal(governance.GovernanceAvailable, status)
	status, err = suite.getRuleStatus(pk, ChainID, contractAddr2)
	suite.Require().Equal(governance.GovernanceBindable, status)
}

//tc:发起正确绑定规则的请求，中继链管理员投票过程中，验证规则状态为binding
func (suite *Model7) Test0714_UpdateRuleWithDoing() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr1, err := client.DeployContract(contract, nil)
	suite.Require().Nil(err)
	contractAddr2, err := client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr1, RegisterRule)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, ChainID, contractAddr2, RegisterRule)
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.String(contractAddr2.String()),
		rpcx.String("reason"),
	}
	_, err = client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), UpdateMasterRule, nil, args...)
	suite.Require().Nil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr1)
	suite.Require().Equal(governance.GovernanceUnbinding, status)
	status, err = suite.getRuleStatus(pk, ChainID, contractAddr2)
	suite.Require().Equal(governance.GovernanceBinding, status)
}

//tc:验证规则状态为binding，发起更新请求，提示对应错误信息
func (suite *Model7) Test0715_UpdateRuleWithBinding() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr1, err := client.DeployContract(contract, nil)
	suite.Require().Nil(err)
	contractAddr2, err := client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr1, RegisterRule)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, ChainID, contractAddr2, RegisterRule)
	suite.Require().Nil(err)
	args := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.String(contractAddr2.String()),
		rpcx.String("reason"),
	}
	_, err = client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), UpdateMasterRule, nil, args...)
	suite.Require().Nil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr1)
	suite.Require().Equal(governance.GovernanceUnbinding, status)
	status, err = suite.getRuleStatus(pk, ChainID, contractAddr2)
	suite.Require().Equal(governance.GovernanceBinding, status)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr2, UpdateMasterRule)
	suite.Require().NotNil(err)
}

//tc:验证规则状态为available，发起更新请求，提示对应错误信息
func (suite *Model7) Test0716_UpdateRuleWithAvailable() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, RegisterRule)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, ChainID, contractAddr, UpdateMasterRule)
	suite.Require().NotNil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr)
	suite.Require().Equal(governance.GovernanceAvailable, status)
}

//tc:验证规则状态为forbidden，发起更新请求，提示对应错误信息
func (suite *Model7) Test0717_UpdateRuleWithForbidden() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr1, err := client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	contractAddr2, err := client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr1, RegisterRule)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, ChainID, contractAddr2, RegisterRule)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, ChainID, contractAddr2, LogoutRule)
	suite.Require().Nil(err)
	err = suite.InvokeRuleContract(pk, ChainID, contractAddr2, UpdateMasterRule)
	suite.Require().NotNil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr1)
	suite.Require().Equal(governance.GovernanceAvailable, status)
	status, err = suite.getRuleStatus(pk, ChainID, contractAddr2)
	suite.Require().Equal(governance.GovernanceForbidden, status)
}

//tc：验证规则不存在，注销验证规则
func (suite *Model7) Test0718_LogoutRuleWithNoRule() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	address := types.NewAddressByStr("0x64C5334AadE6c623ae829422C34B6f310b031aa0")
	err = suite.InvokeRuleContract(pk, ChainID, address, LogoutRule)
	suite.Require().NotNil(err)
}

//tc:验证规则从available状态发起注销的请求，中继链管理员投票不通过，注销失败，验证规则状态不变
func (suite *Model7) Test0719_LogoutRuleWithBinding() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr1, err := client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	contractAddr2, err := client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr1, RegisterRule)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr2, RegisterRule)
	suite.Require().Nil(err)

	args := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.String(contractAddr2.String()),
		rpcx.String("reason"),
	}
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), UpdateMasterRule, nil, args...)
	suite.Require().Nil(err)
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr2)
	suite.Require().Equal(governance.GovernanceBinding, status)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr2, LogoutRule)
	suite.Require().NotNil(err)
}

//tc:验证规则状态为bindable，发起注销请求，提示对应错误信息
func (suite *Model7) Test0720_LogoutRule() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr1, err := client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	contractAddr2, err := client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr1, RegisterRule)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr2, RegisterRule)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr2, LogoutRule)
	suite.Require().Nil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr2)
	suite.Require().Equal(governance.GovernanceForbidden, status)
}

//tc:验证规则状态为available，发起绑定请求，提示对应错误信息
func (suite *Model7) Test0721_LogoutRuleWithAvailable() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr1, err := client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	contractAddr2, err := client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr1, RegisterRule)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr2, RegisterRule)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr1, LogoutRule)
	suite.Require().NotNil(err)
}

//tc:验证规则状态为forbidden，发起注销请求，提示对应错误信息
func (suite *Model7) Test0722_LogoutRuleWithForbidden() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)

	contract, err := ioutil.ReadFile("../../config/rule.wasm")
	suite.Require().Nil(err)

	contractAddr1, err := client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	contractAddr2, err := client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr1, RegisterRule)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr2, RegisterRule)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr2, LogoutRule)
	suite.Require().Nil(err)

	err = suite.InvokeRuleContract(pk, ChainID, contractAddr2, LogoutRule)
	suite.Require().NotNil(err)

	status, err := suite.getRuleStatus(pk, ChainID, contractAddr2)
	suite.Require().Equal(governance.GovernanceForbidden, status)
}

func (suite *Snake) InvokeRuleContract(pk crypto.PrivateKey, ChainID string, contractAddr *types.Address, method string) error {
	client := suite.NewClient(pk)
	var args []*pb.Arg
	if method == LogoutRule {
		args = []*pb.Arg{
			rpcx.String(ChainID),
			rpcx.String(contractAddr.String()),
		}
	} else {
		args = []*pb.Arg{
			rpcx.String(ChainID),
			rpcx.String(contractAddr.String()),
			rpcx.String("reason"),
		}
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
