package bxh_tester

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/meshplus/bitxhub-core/governance"
	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
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

//tc：正确部署验证规则,并返回地址
func (suite Model7) Test0701_DeployRuleIsSuccess() {
	address, err := suite.DeploySimpleRule()
	suite.Require().Nil(err)
	suite.Require().NotNil(address)
}

//tc：注册Fabric V1.4.3类型的应用链，默认验证规则注册成功
func (suite Model7) Test0702_RegisterDefaultRuleWithFabricV143IsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchainWithType(pk, "Fabric V1.4.3", address, "{\"channel_id\":\"mychannel\",\"chaincode_id\":\"broker\",\"broker_version\":\"1\"}")
	suite.Require().Nil(err)
	suite.Require().Equal(true, suite.RuleContains(from, SimFabricRuleAddr))
	suite.Require().Equal(true, suite.RuleContains(from, FabricRuleAddr))
	suite.Require().Equal(true, suite.RuleContains(from, HappyRuleAddr))
}

//tc：注册Fabric V1.4.4类型的应用链，默认验证规则注册成功
func (suite Model7) Test0703_RegisterDefaultRuleWithFabricV144IsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchainWithType(pk, "Fabric V1.4.4", address, "{\"channel_id\":\"mychannel\",\"chaincode_id\":\"broker\",\"broker_version\":\"1\"}")
	suite.Require().Nil(err)
	suite.Require().Equal(true, suite.RuleContains(from, SimFabricRuleAddr))
	suite.Require().Equal(true, suite.RuleContains(from, FabricRuleAddr))
	suite.Require().Equal(true, suite.RuleContains(from, HappyRuleAddr))
}

//tc：注册Hyperchain V1.8.3类型的应用链，默认验证规则注册成功
func (suite Model7) Test0704_RegisterDefaultRuleWithHyperchainV183IsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchainWithType(pk, "Hyperchain V1.8.3", address, "0x857133c5C69e6Ce66F7AD46F200B9B3573e77582")
	suite.Require().Nil(err)
	suite.Require().Equal(true, suite.RuleContains(from, HappyRuleAddr))
}

////tc：注册Hyperchain V1.8.6类型的应用链，默认验证规则注册成功
func (suite Model7) Test0705_RegisterDefaultRuleWithHyperchainV186IsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchainWithType(pk, "Hyperchain V1.8.6", address, "0x857133c5C69e6Ce66F7AD46F200B9B3573e77582")
	suite.Require().Nil(err)
	suite.Require().Equal(true, suite.RuleContains(from, HappyRuleAddr))
}

//tc：注册Flato V1.0.0类型的应用链，默认验证规则注册成功
func (suite Model7) Test0706_RegisterDefaultRuleWithFlatoV100IsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchainWithType(pk, "Flato V1.0.0", address, "0x857133c5C69e6Ce66F7AD46F200B9B3573e77582")
	suite.Require().Nil(err)
	suite.Require().Equal(true, suite.RuleContains(from, HappyRuleAddr))
}

//tc：注册Flato V1.0.3类型的应用链，默认验证规则注册成功
func (suite Model7) Test0707_RegisterDefaultRuleWithFlatoV103IsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchainWithType(pk, "Flato V1.0.3", address, "0x857133c5C69e6Ce66F7AD46F200B9B3573e77582")
	suite.Require().Nil(err)
	suite.Require().Equal(true, suite.RuleContains(from, HappyRuleAddr))
}

//tc：注册Flato V1.0.6类型的应用链，默认验证规则注册成功
func (suite Model7) Test0708_RegisterDefaultRuleWithFlatoV106IsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchainWithType(pk, "Flato V1.0.6", address, "0x857133c5C69e6Ce66F7AD46F200B9B3573e77582")
	suite.Require().Nil(err)
	suite.Require().Equal(true, suite.RuleContains(from, HappyRuleAddr))
}

//tc：注册BCOS V2.6.0类型的应用链，默认验证规则注册成功
func (suite Model7) Test0709_RegisterDefaultRuleWithBCOSV260IsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchainWithType(pk, "BCOS V2.6.0", address, "0x857133c5C69e6Ce66F7AD46F200B9B3573e77582")
	suite.Require().Nil(err)
	suite.Require().Equal(true, suite.RuleContains(from, HappyRuleAddr))
}

//tc：注册CITA V20.2.2类型的应用链，默认验证规则注册成功
func (suite Model7) Test0710_RegisterDefaultRuleWithCITAV2022IsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchainWithType(pk, "CITA V20.2.2", address, "0x857133c5C69e6Ce66F7AD46F200B9B3573e77582")
	suite.Require().Nil(err)
	suite.Require().Equal(true, suite.RuleContains(from, HappyRuleAddr))
}

//tc：注册ETH类型的应用链，默认验证规则注册成功
func (suite Model7) Test0711_RegisterDefaultRuleWithETHIsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchainWithType(pk, "ETH", address, "0x857133c5C69e6Ce66F7AD46F200B9B3573e77582")
	suite.Require().Nil(err)
	suite.Require().Equal(true, suite.RuleContains(from, HappyRuleAddr))
}

//tc：注册其他类型的应用链，默认验证规则注册成功
func (suite Model7) Test0712_RegisterDefaultRuleWithOthersIsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchainWithType(pk, "Other", address, "0x857133c5C69e6Ce66F7AD46F200B9B3573e77582")
	suite.Require().Nil(err)
	suite.Require().Equal(true, suite.RuleContains(from, HappyRuleAddr))
}

//tc：应用链管理员调用注册验证规则，验证规则注册成功
func (suite Model7) Test0713_RegisterRuleIsSuccess() {
	pk, from, address1, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchainWithType(pk, "Hyperchain V1.8.3", address1, "0x857133c5C69e6Ce66F7AD46F200B9B3573e77582")
	suite.Require().Nil(err)
	_, _, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterRule(pk, from, address2)
	suite.Require().Nil(err)
}

//tc：应用链处于updating注册验证规则，验证规则注册成功
func (suite Model7) Test0714_RegisterRuleWithUpdatingChainIsSuccess() {
	pk, from, address1, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToUpdating(pk, from, address1)
	suite.Require().Nil(err)
	_, _, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterRule(pk, from, address2)
	suite.Require().Nil(err)
}

//tc：应用链处于activating注册验证规则，验证规则注册成功
func (suite Model7) Test0715_RegisterRuleWithActivatingChainIsSuccess() {
	pk, from, address1, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToActivating(pk, from, address1)
	suite.Require().Nil(err)
	_, _, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterRule(pk, from, address2)
	suite.Require().Nil(err)
}

//tc：应用链处于freezing注册验证规则，验证规则注册成功
func (suite Model7) Test0716_RegisterRuleWithFreezingChainIsSuccess() {
	pk, from, address1, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToFreezing(pk, from, address1)
	suite.Require().Nil(err)
	_, _, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterRule(pk, from, address2)
	suite.Require().Nil(err)
}

//tc：应用链处于Frozen注册验证规则，验证规则注册成功
func (suite Model7) Test0717_RegisterRuleWithFrozenChainIsSuccess() {
	pk, from, address1, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToFrozen(pk, from, address1)
	suite.Require().Nil(err)
	_, _, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterRule(pk, from, address2)
	suite.Require().Nil(err)
}

//tc：应用链处于Logouting注册验证规则，验证规则注册成功
func (suite Model7) Test0718_RegisterRuleWithLogoutingChainIsSuccess() {
	pk, from, address1, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToLogouting(pk, from, address1)
	suite.Require().Nil(err)
	_, _, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterRule(pk, from, address2)
	suite.Require().Nil(err)
}

//tc：应用链管理员更新验证规则，验证规则更新成功
func (suite Model7) Test0719_UpdateRuleIsSuccess() {
	pk, from, address1, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchainWithType(pk, "Hyperchain V1.8.3", address1, "0x857133c5C69e6Ce66F7AD46F200B9B3573e77582")
	suite.Require().Nil(err)
	_, _, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterRule(pk, from, address2)
	suite.Require().Nil(err)
	err = suite.UpdateMasterRule(pk, from, address2)
	suite.Require().Nil(err)
}

//tc：应用链处于frozen状态更新验证规则，验证规则更新成功
func (suite Model7) Test0720_UpdateRuleWithFrozenChainIsSuccess() {
	pk, from, address1, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToFrozen(pk, from, address1)
	suite.Require().Nil(err)
	_, _, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterRule(pk, from, address2)
	suite.Require().Nil(err)
	err = suite.UpdateMasterRule(pk, from, address2)
	suite.Require().Nil(err)
}

//tc：应用链管理员注销验证规则，验证规则注销成功
func (suite Model7) Test0721_LogoutRuleIsSuccess() {
	pk, from, address1, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchainWithType(pk, "Hyperchain V1.8.3", address1, "0x857133c5C69e6Ce66F7AD46F200B9B3573e77582")
	suite.Require().Nil(err)
	_, _, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterRule(pk, from, address2)
	suite.Require().Nil(err)
	err = suite.LogoutRule(pk, from, address2)
	suite.Require().Nil(err)
}

//tc：应用链处于activating状态注销验证规则，验证规则注销成功
func (suite Model7) Test0722_LogoutRuleWithActivatingChainIsSuccess() {
	pk, from, address1, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToActivating(pk, from, address1)
	suite.Require().Nil(err)
	_, _, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterRule(pk, from, address2)
	suite.Require().Nil(err)
	err = suite.LogoutRule(pk, from, address2)
	suite.Require().Nil(err)
}

//tc：应用链处于freezing状态注销验证规则，验证规则注销成功
func (suite Model7) Test0723_LogoutRuleWithFreezingChainIsSuccess() {
	pk, from, address1, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToFreezing(pk, from, address1)
	suite.Require().Nil(err)
	_, _, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterRule(pk, from, address2)
	suite.Require().Nil(err)
	err = suite.LogoutRule(pk, from, address2)
	suite.Require().Nil(err)
}

//tc：应用链处于frozen状态注销验证规则，验证规则注销成功
func (suite Model7) Test0724_LogoutRuleWithFrozenChainIsSuccess() {
	pk, from, address1, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToFrozen(pk, from, address1)
	suite.Require().Nil(err)
	_, _, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterRule(pk, from, address2)
	suite.Require().Nil(err)
	err = suite.LogoutRule(pk, from, address2)
	suite.Require().Nil(err)
}

//tc：应用链处于logouting状态注销验证规则，验证规则注销成功
func (suite Model7) Test0725_LogoutRuleWithLogoutingChainIsSuccess() {
	pk, from, address1, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToLogouting(pk, from, address1)
	suite.Require().Nil(err)
	_, _, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterRule(pk, from, address2)
	suite.Require().Nil(err)
	err = suite.LogoutRule(pk, from, address2)
	suite.Require().Nil(err)
}

//tc：应用链注销bindable状态的验证规则，验证规则注销成功
func (suite Model7) Test0726_LogoutRuleWithBindableRuleIsSuccess() {
	pk, from, address1, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchainWithType(pk, "Hyperchain V1.8.3", address1, "0x857133c5C69e6Ce66F7AD46F200B9B3573e77582")
	suite.Require().Nil(err)
	_, _, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterRule(pk, from, address2)
	suite.Require().Nil(err)
	err = suite.LogoutRule(pk, from, address2)
	suite.Require().Nil(err)
}

// DeploySimpleRule deploy simple rule
func (suite Snake) DeploySimpleRule() (string, error) {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	contract, err := ioutil.ReadFile("testdata/simple_rule.wasm")
	suite.Require().Nil(err)
	address, err := client.DeployContract(contract, nil)
	suite.Require().Nil(err)
	return address.String(), nil
}

// DeployRule deploy rule and return address
func (suite Snake) DeployRule() (crypto.PrivateKey, string, string, error) {
	address, err := suite.DeploySimpleRule()
	suite.Require().Nil(err)
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	from, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	return pk, from.String(), address, nil
}

// RegisterRule register rule
func (suite Snake) RegisterRule(pk crypto.PrivateKey, ChainID, contractAddr string) error {
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.String(contractAddr),
		rpcx.String("https://github.com/meshplus/bitxhub"),
	}
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), RegisterRule, nil, args...)
	suite.Require().Nil(err)
	if res.Status == pb.Receipt_FAILED {
		return fmt.Errorf(string(res.Ret))
	}
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	if result.ProposalID == "" {
		return nil
	}
	err = suite.VotePass(result.ProposalID)
	suite.Require().Nil(err)
	return nil
}

// UpdateMasterRule update master rule
func (suite Snake) UpdateMasterRule(pk crypto.PrivateKey, ChainID, contractAddr string) error {
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.String(contractAddr),
		rpcx.String("reason"),
	}
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), UpdateMasterRule, nil, args...)
	suite.Require().Nil(err)
	if res.Status == pb.Receipt_FAILED {
		return fmt.Errorf(string(res.Ret))
	}
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	suite.Require().Nil(err)
	if result.ProposalID == "" {
		return nil
	}
	err = suite.VotePass(result.ProposalID)
	suite.Require().Nil(err)
	return nil
}

// LogoutRule logout rule
func (suite Snake) LogoutRule(pk crypto.PrivateKey, ChainID, contractAddr string) error {
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.String(contractAddr),
	}
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), LogoutRule, nil, args...)
	suite.Require().Nil(err)
	if res.Status == pb.Receipt_FAILED {
		return fmt.Errorf(string(res.Ret))
	}
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	suite.Require().Nil(err)
	if result.ProposalID == "" {
		return nil
	}
	err = suite.VotePass(result.ProposalID)
	suite.Require().Nil(err)
	return nil
}

// Rules return all rules
func (suite Snake) Rules(chainID string) ([]Rule, error) {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), "Rules", nil, rpcx.String(chainID))
	suite.Require().Nil(err)
	if res.Status != pb.Receipt_SUCCESS {
		return nil, fmt.Errorf(string(res.Ret))
	}
	var rules []Rule
	err = json.Unmarshal(res.Ret, &rules)
	suite.Require().Nil(err)
	return rules, nil
}

// RuleContains check whether the rule contains
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

// CheckRuleStatus check rule status
func (suite *Snake) CheckRuleStatus(pk crypto.PrivateKey, chainID, address string, expectStatus governance.GovernanceStatus) error {
	status, err := suite.GetRuleStatus(pk, chainID, address)
	suite.Require().Nil(err)
	if expectStatus != status {
		return fmt.Errorf("expect status is %s ,but get status %s", expectStatus, status)
	}
	return nil
}

// GetRuleStatus get rule status by chainID and address
func (suite *Snake) GetRuleStatus(pk crypto.PrivateKey, ChainID string, contractAddr string) (governance.GovernanceStatus, error) {
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.String(contractAddr),
	}
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), "GetRuleByAddr", nil, args...)
	suite.Require().Nil(err)
	if res.Status == pb.Receipt_FAILED {
		return "", fmt.Errorf(string(res.Ret))
	}
	rule := &Rule{}
	err = json.Unmarshal(res.Ret, rule)
	suite.Require().Nil(err)
	return rule.Status, nil
}
