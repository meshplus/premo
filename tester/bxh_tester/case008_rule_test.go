package bxh_tester

import (
	"io/ioutil"

	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
)

func (suite *Snake) TestRegisterRuleShouldSucceed() {
	contract, err := ioutil.ReadFile("./testdata/simple_rule.wasm")
	suite.Nil(err)

	contractAddr, err := suite.client.DeployContract(contract)
	suite.Nil(err)

	args := []*pb.Arg{
		rpcx.String(suite.from.String()),
		rpcx.String(contractAddr.String()),
	}
	res, err := suite.client.InvokeBVMContract(rpcx.RuleManagerContractAddr, "RegisterRule", args...)
	suite.Nil(err)
	suite.True(res.Status == pb.Receipt_SUCCESS)
}

func (suite *Snake) TestAuditRuleShouldSucceed() {
	contract, err := ioutil.ReadFile("./testdata/simple_rule.wasm")
	suite.Nil(err)

	contractAddr, err := suite.client.DeployContract(contract)
	suite.Nil(err)

	args := []*pb.Arg{
		rpcx.String(suite.from.String()),
		rpcx.String(contractAddr.String()),
	}
	res, err := suite.client.InvokeBVMContract(rpcx.RuleManagerContractAddr, "RegisterRule", args...)
	suite.Nil(err)

	args2 := []*pb.Arg{
		rpcx.String(suite.from.String()),
		rpcx.Int32(1),               //audit approve
		rpcx.String("Audit passed"), //desc
	}
	res, err = suite.client.InvokeBVMContract(rpcx.RuleManagerContractAddr, "Audit", args2...)
	suite.Nil(err)
	suite.True(res.Status == pb.Receipt_SUCCESS)

}

func (suite *Snake) TestGetRuleAddressShouldSucceed() {
	// get validation rule contract address when appchain binds rule
	args := []*pb.Arg{
		rpcx.String(suite.from.String()),
		rpcx.String("ethereum"),
	}
	res, err := suite.client.InvokeBVMContract(rpcx.RuleManagerContractAddr, "GetRuleAddress", args...)
	suite.Nil(err)
	suite.True(res.Status == pb.Receipt_SUCCESS)
}

func (suite *Snake) TestGetRuleAddressShouldFail() {
	// get validation rule contract address when appchain not bind rule
	args := []*pb.Arg{
		rpcx.String(suite.to.String()),
		rpcx.String("ethereum"),
	}
	res, err := suite.client.InvokeBVMContract(rpcx.RuleManagerContractAddr, "GetRuleAddress", args...)
	suite.Nil(err)
	suite.True(res.Status == pb.Receipt_FAILED)
}

func (suite *Snake) TestGetFabricRuleAddressShouldSucceed() {
	args := []*pb.Arg{
		rpcx.String(suite.to.String()),
		rpcx.String("fabric"),
	}
	res, err := suite.client.InvokeBVMContract(rpcx.RuleManagerContractAddr, "GetRuleAddress", args...)
	suite.Nil(err)
	suite.True(res.Status == pb.Receipt_SUCCESS)
}

func (suite *Snake) TestRegisterUnexistedWasmRuleShouldFail() {
	contractAddr := "0x1234"
	args := []*pb.Arg{
		rpcx.String(suite.from.String()),
		rpcx.String(contractAddr),
	}
	res, err := suite.client.InvokeBVMContract(rpcx.RuleManagerContractAddr, "RegisterRule", args...)
	suite.NotNil(err)
	suite.True(res.Status == pb.Receipt_FAILED)
}

func (suite *Snake) TestRegisterUnexistedAppchainShouldFail() {
	contract, err := ioutil.ReadFile("./testdata/simple_rule.wasm")
	suite.Nil(err)

	contractAddr, err := suite.client.DeployContract(contract)
	suite.Nil(err)

	args := []*pb.Arg{
		rpcx.String("1234"),
		rpcx.String(contractAddr.String()),
	}
	res, err := suite.client.InvokeBVMContract(rpcx.RuleManagerContractAddr, "RegisterRule", args...)
	suite.NotNil(err)
	suite.True(res.Status == pb.Receipt_FAILED)
}
