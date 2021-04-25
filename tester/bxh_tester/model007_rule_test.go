package bxh_tester

import (
	"fmt"
	"io/ioutil"

	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
)

//tc:注册规则，指定WASM合约地址与应用链ID绑定，返回回执状态成功
func (suite *Snake) Test0701_RegisterRuleShouldSucceed() {
	_, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("./testdata/simple_rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	args := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.String(contractAddr.String()),
	}
	res, err := suite.client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), "RegisterRule", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
}

//tc:注册规则，指定WASM合约地址与不存在的应用链ID绑定，返回回执状态失败
func (suite *Snake) Test0702_RegisterUnexistedAppchainShouldFail() {
	_, _, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("./testdata/simple_rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	args := []*pb.Arg{
		rpcx.String("1234"),
		rpcx.String(contractAddr.String()),
	}
	res, err := suite.client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), "RegisterRule", nil, args...)

	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	suite.Require().Contains(string(res.Ret), "this appchain does not exist")
}

//tc:审核规则，指定WASM合约审核，返回回执状态成功
func (suite *Snake) Test0703_AuditRuleShouldSucceed() {
	_, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	contract, err := ioutil.ReadFile("./testdata/simple_rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	args := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.String(contractAddr.String()),
	}
	res, err := suite.client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), "RegisterRule", nil, args...)
	suite.Require().Nil(err)

	args2 := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.Int32(1),               //audit approve
		rpcx.String("Audit passed"), //desc
	}
	res, err = suite.client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), "Audit", nil, args2...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

}

//tc:获取规则地址，根据应用链ID和链类型获取合约地址，返回回执状态成功
func (suite *Snake) Test0704_GetRuleAddressShouldSucceed() {
	// get validation rule contract address when appchain binds rule
	_, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	contract, err := ioutil.ReadFile("./testdata/simple_rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	args := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.String(contractAddr.String()),
	}
	res, err := suite.client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), "RegisterRule", nil, args...)
	suite.Require().Nil(err)

	args = []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.String("ethereum"),
	}
	res, err = suite.client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), "GetRuleAddress", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	suite.Require().NotNil(res.Ret)
}

//tc:获取规则地址，根据应用链ID和链类型获取合约地址，应用链未绑定合约，返回回执失败
func (suite *Snake) Test0705_GetRuleAddressShouldFail() {
	// get validation rule contract address when appchain not bind rule
	_, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	args := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.String("ethereum"),
	}
	res, err := suite.client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), "GetRuleAddress", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
}

func (suite *Snake) Test0706_GetFabricRuleAddressShouldSucceed() {
	_, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	args := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.String("fabric"),
	}
	res, err := suite.client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), "GetRuleAddress", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	suite.Require().NotNil(res.Ret)
}
