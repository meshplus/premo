package bxh_tester

import (
	"io/ioutil"

	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"

	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
)

type Model7 struct {
	*Snake
}

//tc:注册规则，指定WASM合约地址与应用链ID绑定，返回回执状态成功
func (suite *Model7) Test0701_RegisterRuleShouldSucceed() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)

	contract, err := ioutil.ReadFile("./testdata/simple_rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	args := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.String(contractAddr.String()),
	}
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), "RegisterRule", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
}

//tc:注册规则，指定WASM合约地址与不存在的应用链ID绑定，返回回执状态失败
func (suite *Model7) Test0702_RegisterUnexistedAppchainShouldFail() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)

	contract, err := ioutil.ReadFile("./testdata/simple_rule.wasm")
	suite.Require().Nil(err)

	contractAddr, err := client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	args := []*pb.Arg{
		rpcx.String("1234"),
		rpcx.String(contractAddr.String()),
	}
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), "RegisterRule", nil, args...)

	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	suite.Require().Contains(string(res.Ret), "this appchain does not exist")
}

//tc:审核规则，指定WASM合约审核，返回回执状态成功
func (suite *Model7) Test0704_AuditRuleShouldSucceed() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)

	contract, err := ioutil.ReadFile("./testdata/simple_rule.wasm")
	suite.Require().Nil(err)
	contractAddr, err := client.DeployContract(contract, nil)
	suite.Require().Nil(err)

	args := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.String(contractAddr.String()),
	}
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), "RegisterRule", nil, args...)
	suite.Require().Nil(err)

	args2 := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.Int32(1),               //audit approve
		rpcx.String("Audit passed"), //desc
	}
	res, err = client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), "Audit", nil, args2...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

}

//tc:获取规则地址，根据应用链ID和链类型获取合约地址，应用链未绑定合约，返回回执失败
func (suite *Model7) Test0705_GetRuleAddressShouldFail() {
	// get validation rule contract address when appchain not bind rule
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)

	args := []*pb.Arg{
		rpcx.String(suite.to.String()),
		rpcx.String("ethereum"),
	}
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), "GetRuleAddress", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
}

func (suite *Model7) Test0706_GetFabricRuleAddressShouldSucceed() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)

	args := []*pb.Arg{
		rpcx.String(suite.to.String()),
		rpcx.String("fabric"),
	}
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), "GetRuleAddress", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	suite.Require().NotNil(res.Ret)
}
