package bxh_tester

import (
	"io/ioutil"

	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
)

func (suite *Snake) TestRegisterRuleShouldSucceed() {
	suite.RegisterAppchain(suite.pk, "hyperchain")

	contract, err := ioutil.ReadFile("./testdata/simple_rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract,nil)
	suite.Require().Nil(err)

	args := []*pb.Arg{
		rpcx.String(suite.from.String()),
		rpcx.String(contractAddr.String()),
	}
	res, err := suite.client.InvokeBVMContract(rpcx.RuleManagerContractAddr, "RegisterRule",nil, args...)
	suite.Require().Nil(err)
	suite.Require().True(res.Status == pb.Receipt_SUCCESS)
}

func (suite *Snake) TestAuditRuleShouldSucceed() {
	suite.RegisterAppchain(suite.pk, "hyperchain")

	contract, err := ioutil.ReadFile("./testdata/simple_rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract,nil)
	suite.Require().Nil(err)

	args := []*pb.Arg{
		rpcx.String(suite.from.String()),
		rpcx.String(contractAddr.String()),
	}
	res, err := suite.client.InvokeBVMContract(rpcx.RuleManagerContractAddr, "RegisterRule",nil, args...)
	suite.Require().Nil(err)

	args2 := []*pb.Arg{
		rpcx.String(suite.from.String()),
		rpcx.Int32(1),               //audit approve
		rpcx.String("Audit passed"), //desc
	}
	res, err = suite.client.InvokeBVMContract(rpcx.RuleManagerContractAddr, "Audit",nil, args2...)
	suite.Require().Nil(err)
	suite.Require().True(res.Status == pb.Receipt_SUCCESS)

}

func (suite *Snake) TestGetRuleAddressShouldSucceed() {
	// get validation rule contract address when appchain binds rule
	args := []*pb.Arg{
		rpcx.String(suite.from.String()),
		rpcx.String("ethereum"),
	}
	res, err := suite.client.InvokeBVMContract(rpcx.RuleManagerContractAddr, "GetRuleAddress",nil, args...)
	suite.Require().Nil(err)
	suite.Require().True(res.Status == pb.Receipt_SUCCESS)
}

func (suite *Snake) TestGetRuleAddressShouldFail() {
	// get validation rule contract address when appchain not bind rule
	args := []*pb.Arg{
		rpcx.String(suite.to.String()),
		rpcx.String("ethereum"),
	}
	res, err := suite.client.InvokeBVMContract(rpcx.RuleManagerContractAddr, "GetRuleAddress",nil, args...)
	suite.Require().Nil(err)
	suite.Require().True(res.Status == pb.Receipt_FAILED)
}

func (suite *Snake) TestGetFabricRuleAddressShouldSucceed() {
	args := []*pb.Arg{
		rpcx.String(suite.to.String()),
		rpcx.String("fabric"),
	}
	res, err := suite.client.InvokeBVMContract(rpcx.RuleManagerContractAddr, "GetRuleAddress",nil, args...)
	suite.Require().Nil(err)
	suite.Require().True(res.Status == pb.Receipt_SUCCESS)
}

func (suite *Snake) TestRegisterUnexistedAppchainShouldFail() {
	contract, err := ioutil.ReadFile("./testdata/simple_rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := suite.client.DeployContract(contract,nil)
	suite.Require().Nil(err)

	args := []*pb.Arg{
		rpcx.String("1234"),
		rpcx.String(contractAddr.String()),
	}
	res, err := suite.client.InvokeBVMContract(rpcx.RuleManagerContractAddr, "RegisterRule",nil, args...)

	suite.Require().Nil(err)
	suite.Require().True(res.Status == pb.Receipt_FAILED)
}
