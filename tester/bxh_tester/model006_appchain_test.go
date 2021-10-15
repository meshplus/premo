package bxh_tester

import (
	"encoding/json"
	"sync/atomic"

	"github.com/meshplus/bitxhub-core/governance"
	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/meshplus/premo/internal/repo"
	"github.com/pkg/errors"
)

type Model6 struct {
	*Snake
}

func (suite Model6) SetupTest() {
	suite.T().Parallel()
}

//tc：通过正确的参数注册应用链，应用链注册成功
func (suite Model6) Test0601_RegisterAppchainIsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
}

//tc：通过被占用的应用链名称注册应用链，应用链注册失败
func (suite Model6) Test0602_RegisterAppchainWithUsedNameIsFail() {
	pk1, from1, address1, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk1, from1, address1)
	suite.Require().Nil(err)
	pk2, _, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk2, from1, address2)
	suite.Require().NotNil(err)
}

//tc：通过曾被占用的应用链名称注册应用链，应用链注册失败
func (suite Model6) Test0603_RegisterAppchainWithFreeNameIsSuccess() {
	pk1, from1, address1, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk1, from1, address1)
	suite.Require().Nil(err)
	ChainID, err := suite.GetChainIDByName(from1)
	suite.Require().Nil(err)
	err = suite.UpdateAppchain(pk1, ChainID, from1+"123", "desc", []byte(""), from1)
	suite.Require().Nil(err)
	pk2, _, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk2, from1, address2)
	suite.Require().Nil(err)
}

//tc：通过空的应用链名称注册应用链，应用链注册失败
func (suite Model6) Test0604_RegisterAppchainWithEmptyNameIsFail() {
	pk, _, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, "", address)
	suite.Require().NotNil(err)
}

//tc：通过空的broker合约地址注册应用链，应用链注册失败
func (suite Model6) Test0605_RegisterAppchainWithEmptyBrokerIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(from),           //chainName
		rpcx.String("Flato V1.0.3"), //chainType
		rpcx.Bytes([]byte("")),      //trustRoot
		rpcx.String("0x857133c5C69e6Ce66F7AD46F200B9B3573e77582"), //broker
		rpcx.String("desc"),               //desc
		rpcx.String(address),              //masterRuleAddr
		rpcx.String("https://github.com"), //masterRuleUrl
		rpcx.String(from),                 //adminAddrs
		rpcx.String("reason"),             //reason
	}
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "RegisterAppchain", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
}

//tc：通过不存在的验证规则地址注册应用链，应用链注册失败
func (suite Model6) Test0606_RegisterAppchainWithNoExistRuleIsFail() {
	pk, from, _, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, "0x857133c5C69e6Ce66F7AD46F200B9B3573e77582")
	suite.Require().NotNil(err)
}

//tc：通过空的验证规则地址注册应用链，应用链注册失败
func (suite Model6) Test0607_RegisterAppchainWithEmptyRuleIsFail() {
	pk, from, _, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, "")
	suite.Require().NotNil(err)
}

//tc：通过其他应用链的默认验证规则注册应用链，应用链注册失败
func (suite Model6) Test0608_RegisterAppchainWithOthersRuleIsFail() {
	pk, from, _, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, "0x00000000000000000000000000000000000000a0")
	suite.Require().NotNil(err)
}

//tc：通过被占用的管理员地址注册应用链，应用链注册失败
func (suite Model6) Test0609_RegisterAppchainWithRepeatedAdminIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(from + "123"),                                 //chainName
		rpcx.String("Flato V1.0.3"),                               //chainType
		rpcx.Bytes([]byte("")),                                    //trustRoot
		rpcx.String("0x857133c5C69e6Ce66F7AD46F200B9B3573e77582"), //broker
		rpcx.String("desc"),                                       //desc
		rpcx.String(address),                                      //masterRuleAddr
		rpcx.String("https://github.com"),                         //masterRuleUrl
		rpcx.String(from),                                         //adminAddrs
		rpcx.String("reason"),                                     //reason
	}
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "RegisterAppchain", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
}

//tc：通过曾被占用的管理员地址注册应用链，应用链注册失败
func (suite Model6) Test0610_RegisterAppchainWithFreeAdminIsSuccess() {
	pk1, from1, address1, err := suite.DeployRule()
	suite.Require().Nil(err)
	pk2, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	from2, err := pk2.PublicKey().Address()
	suite.Require().Nil(err)
	client := suite.NewClient(pk1)
	args := []*pb.Arg{
		rpcx.String(from1),
		rpcx.String(from1),          //chainName
		rpcx.String("Flato V1.0.3"), //chainType
		rpcx.Bytes([]byte("")),      //trustRoot
		rpcx.String("0x857133ce6Ce66F7AD46F200B9B3573e77582"), //broker
		rpcx.String("desc"),                       //desc
		rpcx.String(address1),                     //masterRuleAddr
		rpcx.String("https://github.com"),         //masterRuleUrl
		rpcx.String(from1 + "," + from2.String()), //adminAddrs
		rpcx.String("reason"),                     //reason
	}
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "RegisterAppchain", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	pk3, from3, address3, err := suite.DeployRule()
	suite.Require().Nil(err)
	client = suite.NewClient(pk3)
	args = []*pb.Arg{
		rpcx.String(from3),
		rpcx.String(from3),          //chainName
		rpcx.String("Flato V1.0.3"), //chainType
		rpcx.Bytes([]byte("")),      //trustRoot
		rpcx.String("0x857133c5C69e6Ce66F7AD46F200B9B3573e77582"), //broker
		rpcx.String("desc"),                       //desc
		rpcx.String(address3),                     //masterRuleAddr
		rpcx.String("https://github.com"),         //masterRuleUrl
		rpcx.String(from3 + "," + from2.String()), //adminAddrs
		rpcx.String("reason"),                     //reason
	}
	res, err = client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "RegisterAppchain", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
}

//tc：通过空的管理员地址注册应用链，应用链注册失败
func (suite Model6) Test0611_RegisterAppchainWithEmptyAdminIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(from),           //chainName
		rpcx.String("Flato V1.0.3"), //chainType
		rpcx.Bytes([]byte("")),      //trustRoot
		rpcx.String("0x857133c5C69e6Ce66F7AD46F200B9B3573e77582"), //broker
		rpcx.String("desc"),               //desc
		rpcx.String(address),              //masterRuleAddr
		rpcx.String("https://github.com"), //masterRuleUrl
		rpcx.String(""),                   //adminAddrs
		rpcx.String("reason"),             //reason
	}
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "RegisterAppchain", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
}

//tc：注册应用链，应用链管理员不包含发起人，应用链注册失败
func (suite Model6) Test0612_RegisterAppchainWithNoExistSelfIsFail() {
	pk, from1, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	pk2, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	from2, err := pk2.PublicKey().Address()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(from1),          //chainName
		rpcx.String("Flato V1.0.3"), //chainType
		rpcx.Bytes([]byte("")),      //trustRoot
		rpcx.String("0x857133c5C69e6Ce66F7AD46F200B9B3573e77582"), //broker
		rpcx.String("desc"),               //desc
		rpcx.String(address),              //masterRuleAddr
		rpcx.String("https://github.com"), //masterRuleUrl
		rpcx.String(from2.String()),       //adminAddrs
		rpcx.String("reason"),             //reason
	}
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "RegisterAppchain", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
}

//tc：注册应用链，应用链管理员包含中继链管理员，应用链注册失败
func (suite Model6) Test0613_RegisterAppchainWithRelayAdminIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(from),           //chainName
		rpcx.String("Flato V1.0.3"), //chainType
		rpcx.Bytes([]byte("")),      //trustRoot
		rpcx.String("0x857133c5C69e6Ce66F7AD46F200B9B3573e77582"), //broker
		rpcx.String("desc"),                                       //desc
		rpcx.String(address),                                      //masterRuleAddr
		rpcx.String("https://github.com"),                         //masterRuleUrl
		rpcx.String("0xc7F999b83Af6DF9e67d0a37Ee7e900bF38b3D013"), //adminAddrs
		rpcx.String("reason"),                                     //reason
	}
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "RegisterAppchain", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
}

//tc：通过错误的broker合约地址注册fabric类型的应用链，应用链注册失败
func (suite Model6) Test0614_RegisterAppchainWithErrorBrokerIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(from),            //chainName
		rpcx.String("Fabric v1.4.3"), //chainType
		rpcx.Bytes([]byte("")),       //trustRoot
		rpcx.String("0x857133c5C69e6Ce66F7AD46F200B9B3573e77582"), //broker
		rpcx.String("desc"),               //desc
		rpcx.String(address),              //masterRuleAddr
		rpcx.String("https://github.com"), //masterRuleUrl
		rpcx.String(from),                 //adminAddrs
		rpcx.String("reason"),             //reason
	}
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "RegisterAppchain", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
}

//tc：非应用链管理员更新应用链，应用链更新失败
func (suite Model6) Test0615_UpdateAppchainWithNoAdminIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	pk2, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	ChainID, err := suite.GetChainIDByName(from)
	suite.Require().Nil(err)
	err = suite.UpdateAppchain(pk2, ChainID, from, "desc", []byte(""), from)
	suite.Require().NotNil(err)
}

//tc：通过正确的参数更新应用链，应用链更新成功
func (suite Model6) Test0616_UpdateAppchainIsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	ChainID, err := suite.GetChainIDByName(from)
	suite.Require().Nil(err)
	err = suite.UpdateAppchain(pk, ChainID, from+"123", "desc", []byte(""), from)
	suite.Require().Nil(err)
}

//tc：通过不存在的应用链id更新应用链，应用链更新失败
func (suite Model6) Test0617_UpdateAppchainWithNoExistIDIsFail() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	from, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	err = suite.UpdateAppchain(pk, "test", from.String(), "desc", []byte(""), from.String())
	suite.Require().NotNil(err)
}

//tc：通过空的应用链id更新应用链，应用链更新失败
func (suite Model6) Test0618_UpdateAppchainWithEmptyIDIsFail() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	from, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	err = suite.UpdateAppchain(pk, "", from.String(), "desc", []byte(""), from.String())
	suite.Require().NotNil(err)
}

//tc：通过被占用的应用链名称更新应用链，应用链更新失败
func (suite Model6) Test0619_UpdateAppchainWithUsedNameIsFail() {
	pk1, from1, address1, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk1, from1, address1)
	suite.Require().Nil(err)
	pk2, from2, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk2, from2, address2)
	suite.Require().Nil(err)
	chainID, err := suite.GetChainIDByName(from2)
	suite.Require().Nil(err)
	err = suite.UpdateAppchain(pk2, chainID, from1, "desc", []byte(""), from2)
	suite.Require().NotNil(err)
}

//tc：通过曾被占用的应用链名称更新应用链，应用链更新成功
func (suite Model6) Test0620_UpdateAppchainWithFreeNameIsSuccess() {
	pk1, from1, address1, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk1, from1, address1)
	suite.Require().Nil(err)
	chainID1, err := suite.GetChainIDByName(from1)
	suite.Require().Nil(err)
	err = suite.UpdateAppchain(pk1, chainID1, from1+"123", "desc", []byte(""), from1)
	suite.Require().Nil(err)
	pk2, from2, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk2, from2, address2)
	suite.Require().Nil(err)
	chainID2, err := suite.GetChainIDByName(from2)
	suite.Require().Nil(err)
	err = suite.UpdateAppchain(pk2, chainID2, from1, "desc", []byte(""), from2)
	suite.Require().Nil(err)
}

//tc：通过空的应用链名称更新应用链，应用链更新失败
func (suite Model6) Test0621_UpdateAppchainWithEmptyNameIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	chainID, err := suite.GetChainIDByName(from)
	suite.Require().Nil(err)
	err = suite.UpdateAppchain(pk, chainID, "", "desc", []byte(""), from)
	suite.Require().NotNil(err)
}

//tc：通过被占用的管理员地址更新应用链，应用链更新失败
func (suite Model6) Test0622_UpdateAppchainWithUsedAdminIsFail() {
	pk1, from1, address1, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk1, from1, address1)
	suite.Require().Nil(err)
	pk2, from2, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk2, from2, address2)
	suite.Require().Nil(err)
	chainID2, err := suite.GetChainIDByName(from2)
	suite.Require().Nil(err)
	err = suite.UpdateAppchain(pk2, chainID2, "", "desc", []byte(""), from1+","+from2)
	suite.Require().NotNil(err)
}

//tc：通过曾被占用的管理员地址更新应用链，应用链更新成功
func (suite Model6) Test0623_UpdateAppchainWithFreeAdminIsSuccess() {
	pk1, from1, address1, err := suite.DeployRule()
	suite.Require().Nil(err)
	pk2, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	from2, err := pk2.PublicKey().Address()
	suite.Require().Nil(err)
	client := suite.NewClient(pk1)
	args := []*pb.Arg{
		rpcx.String(from1),
		rpcx.String(from1),          //chainName
		rpcx.String("Flato V1.0.3"), //chainType
		rpcx.Bytes([]byte("")),      //trustRoot
		rpcx.String("0x857133c5C69e6Ce66F7AD46F200B9B3573e77582"), //broker
		rpcx.String("desc"),                       //desc
		rpcx.String(address1),                     //masterRuleAddr
		rpcx.String("https://github.com"),         //masterRuleUrl
		rpcx.String(from1 + "," + from2.String()), //adminAddrs
		rpcx.String("reason"),                     //reason
	}
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "RegisterAppchain", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	suite.Require().Nil(err)
	err = suite.VotePass(result.ProposalID)
	suite.Require().Nil(err)
	chainID1, err := suite.GetChainIDByName(from1)
	suite.Require().Nil(err)
	err = suite.UpdateAppchain(pk1, chainID1, from1, "desc", []byte(""), from1)
	suite.Require().Nil(err)
	pk3, from3, address3, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk3, from3, address3)
	suite.Require().Nil(err)
	chainID3, err := suite.GetChainIDByName(from3)
	suite.Require().Nil(err)
	err = suite.UpdateAppchain(pk3, chainID3, from3, "desc", []byte(""), from3+","+from2.String())
	suite.Require().Nil(err)
}

//tc：通过空的管理员地址更新应用链，应用链更新失败
func (suite Model6) Test0624_UpdateAppchainWithEmptyAdminIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	chainID, err := suite.GetChainIDByName(from)
	suite.Require().Nil(err)
	err = suite.UpdateAppchain(pk, chainID, from, "desc", []byte(""), "")
	suite.Require().NotNil(err)
}

//tc：更新应用链，应用链管理员不包含发起人，应用链更新失败
func (suite Model6) Test0625_UpdateAppchainWithNoExistSelfIsFail() {
	pk1, from1, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	pk2, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	from2, err := pk2.PublicKey().Address()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk1, from1, address)
	suite.Require().Nil(err)
	chainID, err := suite.GetChainIDByName(from1)
	suite.Require().Nil(err)
	err = suite.UpdateAppchain(pk1, chainID, from1, "desc", []byte(""), from2.String())
	suite.Require().NotNil(err)
}

//tc：应用链处于activating状态更新应用链，应用链更新失败
func (suite Model6) Test0626_UpdateAppchainWithActivatingChainIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToActivating(pk, from, address)
	suite.Require().Nil(err)
	chainID, err := suite.GetChainIDByName(from)
	suite.Require().Nil(err)
	err = suite.UpdateAppchain(pk, chainID, from, "desc", []byte(""), from)
	suite.Require().NotNil(err)
}

//tc：应用链处于freezing状态更新应用链，应用链更新失败
func (suite Model6) Test0627_UpdateAppchainWithFreezingChainIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToFreezing(pk, from, address)
	suite.Require().Nil(err)
	chainID, err := suite.GetChainIDByName(from)
	suite.Require().Nil(err)
	err = suite.UpdateAppchain(pk, chainID, from, "desc", []byte(""), from)
	suite.Require().NotNil(err)
}

//tc：应用链处于frozen状态更新应用链，应用链更新成功
func (suite Model6) Test0628_UpdateAppchainWithFrozenChainIsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToFrozen(pk, from, address)
	suite.Require().Nil(err)
	chainID, err := suite.GetChainIDByName(from)
	suite.Require().Nil(err)
	err = suite.UpdateAppchain(pk, chainID, from, "desc", []byte(""), from)
	suite.Require().Nil(err)
}

//tc：应用链处于logouting状态更新应用链，应用链更新失败
func (suite Model6) Test0629_UpdateAppchainWithLogoutingChainIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToLogouting(pk, from, address)
	suite.Require().Nil(err)
	chainID, err := suite.GetChainIDByName(from)
	suite.Require().Nil(err)
	err = suite.UpdateAppchain(pk, chainID, from, "desc", []byte(""), from)
	suite.Require().NotNil(err)
}

//tc：应用链处于forbidden状态更新应用链，应用链更新失败
func (suite Model6) Test0630_UpdateAppchainWithForbiddenChainIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToForbidden(pk, from, address)
	suite.Require().Nil(err)
	chainID, err := suite.GetChainIDByName(from)
	suite.Require().Nil(err)
	err = suite.UpdateAppchain(pk, chainID, from, "desc", []byte(""), from)
	suite.Require().NotNil(err)
}

//tc：应用链更新名称字段，产生提案
func (suite Model6) Test0631_UpdateAppchainWithNameFieldHaveProposal() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	ChainID, err := suite.GetChainIDByName(from)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.String(from + "123"),
		rpcx.String("desc"),
		rpcx.Bytes([]byte("")),
		rpcx.String(from),
		rpcx.String("reason"),
	}
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "UpdateAppchain", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	suite.Require().NotEqual("", result.ProposalID)
}

//tc：应用链信任根名称字段，产生提案
func (suite Model6) Test0632_UpdateAppchainWithTrustRootFieldHaveProposal() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	ChainID, err := suite.GetChainIDByName(from)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.String(from),
		rpcx.String("desc"),
		rpcx.Bytes([]byte("123")),
		rpcx.String(from),
		rpcx.String("reason"),
	}
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "UpdateAppchain", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	suite.Require().NotEqual("", result.ProposalID)
}

//tc：应用链更新管理员地址字段，产生提案
func (suite Model6) Test0633_UpdateAppchainWithAdminsFieldHaveProposal() {
	pk1, from1, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk1, from1, address)
	suite.Require().Nil(err)
	ChainID, err := suite.GetChainIDByName(from1)
	suite.Require().Nil(err)
	pk2, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	from2, err := pk2.PublicKey().Address()
	suite.Require().Nil(err)
	client := suite.NewClient(pk1)
	args := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.String(from1),
		rpcx.String("desc"),
		rpcx.Bytes([]byte("")),
		rpcx.String(from1 + "," + from2.String()),
		rpcx.String("reason"),
	}
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "UpdateAppchain", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	suite.Require().NotEqual("", result.ProposalID)
}

//tc：应用链更新描述字段，不产生提案
func (suite Model6) Test0634_UpdateAppchainWithDescFieldNoProposal() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	ChainID, err := suite.GetChainIDByName(from)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(ChainID),
		rpcx.String(from),
		rpcx.String("desc123"),
		rpcx.Bytes([]byte("")),
		rpcx.String(from),
		rpcx.String("reason"),
	}
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "UpdateAppchain", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	suite.Require().Equal("", result.ProposalID)
}

//tc：非中继链管理员冻结应用链，应用链冻结失败
func (suite Model6) Test0635_FreezeAppchainWithNoAdminIsFail() {
	pk1, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk1, from, address)
	suite.Require().Nil(err)
	pk2, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk2)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "FreezeAppchain", nil, rpcx.String(from), rpcx.String("reason"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
}

//tc：应用链未注册，冻结应用链，应用链冻结失败
func (suite Model6) Test0636_FreezeAppchainWithNoRegisterChainIsFail() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "FreezeAppchain", nil, rpcx.String(suite.GetChainID(pk)), rpcx.String("reason"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
}

//tc：中继链管理员冻结应用链，应用链冻结成功
func (suite Model6) Test0637_FreezeAppchainIsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	err = suite.FreezeAppchain(from)
	suite.Require().Nil(err)
}

//tc：应用链处于activating状态冻结应用链，应用链冻结成功
func (suite Model6) Test0640_FreezeAppchainWithActivatingChainIsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToActivating(pk, from, address)
	suite.Require().Nil(err)
	err = suite.FreezeAppchain(from)
	suite.Require().Nil(err)
	err = suite.CheckChainStatus(from, governance.GovernanceFrozen)
	suite.Require().Nil(err)
}

//tc：应用链处于freezing状态冻结应用链，应用链冻结失败
func (suite Model6) Test0641_FreezeAppchainWithFreezingChainIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToFreezing(pk, from, address)
	suite.Require().Nil(err)
	err = suite.FreezeAppchain(from)
	suite.Require().NotNil(err)
}

//tc：应用链处于frozen状态冻结应用链，应用链冻结失败
func (suite Model6) Test0642_FreezeAppchainWithFrozenChainIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	err = suite.FreezeAppchain(from)
	suite.Require().Nil(err)
	err = suite.CheckChainStatus(from, governance.GovernanceFrozen)
	suite.Require().Nil(err)
	err = suite.FreezeAppchain(from)
	suite.Require().NotNil(err)
}

//tc：应用链处于logouting状态冻结应用链，应用链冻结失败
func (suite Model6) Test0643_FreezeAppchainWithLogoutingChainIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToLogouting(pk, from, address)
	suite.Require().Nil(err)
	err = suite.FreezeAppchain(from)
	suite.Require().NotNil(err)
}

//tc：应用链处于forbidden状态冻结应用链，应用链冻结失败
func (suite Model6) Test0644_FreezeAppchainWithForbiddenChainIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToForbidden(pk, from, address)
	suite.Require().Nil(err)
	err = suite.FreezeAppchain(from)
	suite.Require().NotNil(err)
}

//tc：非中继链管理员非应用链管理员激活应用链，应用链激活失败
func (suite Model6) Test0645_ActivateAppchainWithNoAdminIsFail() {
	pk1, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToFrozen(pk1, from, address)
	suite.Require().Nil(err)
	pk2, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.ActivateAppchain(pk2, from)
	suite.Require().NotNil(err)
}

//tc：应用链管理员激活应用链，应用链激活成功
func (suite Model6) Test0646_ActivateAppchainIsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToFrozen(pk, from, address)
	suite.Require().Nil(err)
	err = suite.ActivateAppchain(pk, from)
	suite.Require().Nil(err)
}

//tc：中继链管理员激活应用链，应用链激活成功
func (suite Model6) Test0647_ActivateAppchainWithRelayAdminIsSuccess() {
	pk, chainID, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToFrozen(pk, chainID, address)
	suite.Require().Nil(err)
	path, err := repo.Node1Path()
	suite.Require().Nil(err)
	node1pk, err := asym.RestorePrivateKey(path, repo.KeyPassword)
	suite.Require().Nil(err)
	from, err := node1pk.PublicKey().Address()
	suite.Require().Nil(err)
	client := suite.NewClient(node1pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "ActivateAppchain", &rpcx.TransactOpts{
		From:  from.String(),
		Nonce: atomic.AddUint64(&nonce1, 1),
	}, rpcx.String(chainID), rpcx.String("reason"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	suite.Require().Nil(err)
	err = suite.VotePass(result.ProposalID)
	suite.Require().Nil(err)
}

//tc：应用链未注册，激活应用链，应用链激活失败
func (suite Model6) Test0648_ActivateAppchainWithNoRegisterChainIsFail() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.ActivateAppchain(pk, suite.GetChainID(pk))
	suite.Require().NotNil(err)
}

//tc：应用链处于activating状态激活应用链，应用链激活失败
func (suite Model6) Test0651_ActivateAppchainWithActivatingChainIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToActivating(pk, from, address)
	suite.Require().Nil(err)
	err = suite.ActivateAppchain(pk, from)
	suite.Require().NotNil(err)
}

//tc：应用链处于freezing状态激活应用链，应用链激活失败
func (suite Model6) Test0652_ActivateAppchainWithFreezingChainIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToFreezing(pk, from, address)
	suite.Require().Nil(err)
	err = suite.ActivateAppchain(pk, from)
	suite.Require().NotNil(err)
}

//tc：应用链处于logouting状态激活应用链，应用链激活失败
func (suite Model6) Test0653_ActivateAppchainWithLogoutingChainIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToLogouting(pk, from, address)
	suite.Require().Nil(err)
	err = suite.ActivateAppchain(pk, from)
	suite.Require().NotNil(err)
}

//tc：应用链处于forbidden状态激活应用链，应用链激活失败
func (suite Model6) Test0654_ActivateAppchainWithForbiddenChainIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToForbidden(pk, from, address)
	suite.Require().Nil(err)
	err = suite.ActivateAppchain(pk, from)
	suite.Require().NotNil(err)
}

//tc：非应用链管理员注销应用链，应用链注销失败
func (suite Model6) Test0655_LogoutAppchainWithNoAdminIsFail() {
	pk1, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk1, from, address)
	suite.Require().Nil(err)
	pk2, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.LogoutAppchain(pk2, from)
	suite.Require().NotNil(err)
}

//tc：应用链管理员注销应用链，应用链注销成功
func (suite Model6) Test0656_LogoutAppchainIsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	err = suite.LogoutAppchain(pk, from)
	suite.Require().Nil(err)
}

//tc：应用链未注册，注销应用链，应用链注销失败
func (suite Model6) Test0657_LogoutAppchainWithNoRegisterChainIsFail() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.LogoutAppchain(pk, suite.GetChainID(pk))
	suite.Require().NotNil(err)
}

//tc：应用链处于activating状态注销应用链，应用链注销成功
func (suite Model6) Test0660_LogoutAppchainWithActivatingChainIsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToActivating(pk, from, address)
	suite.Require().Nil(err)
	err = suite.LogoutAppchain(pk, from)
	suite.Require().Nil(err)
}

//tc：应用链处于freezing状态注销应用链，应用链激活成功
func (suite Model6) Test0661_LogoutAppchainWithFreezingChainIsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToFreezing(pk, from, address)
	suite.Require().Nil(err)
	err = suite.LogoutAppchain(pk, from)
	suite.Require().Nil(err)
}

//tc：应用链处于frozen状态注销应用链，应用链注销成功
func (suite Model6) Test0662_LogoutAppchainWithFrozenChainIsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToFrozen(pk, from, address)
	suite.Require().Nil(err)
	err = suite.LogoutAppchain(pk, from)
	suite.Require().Nil(err)
}

//tc：应用链处于logouting状态注销应用链，应用链注销失败
func (suite Model6) Test0663_LogoutAppchainWithLogoutingChainIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToLogouting(pk, from, address)
	suite.Require().Nil(err)
	err = suite.LogoutAppchain(pk, from)
	suite.Require().NotNil(err)
}

//tc：应用链处于forbidden状态注销应用链，应用链注销失败
func (suite Model6) Test0664_LogoutAppchainWithForbiddenChainIsFail() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToForbidden(pk, from, address)
	suite.Require().Nil(err)
	err = suite.LogoutAppchain(pk, from)
	suite.Require().NotNil(err)
}

func (suite *Snake) FreezeAppchain(chainID string) error {
	path, err := repo.Node1Path()
	if err != nil {
		return err
	}
	node1Key, err := asym.RestorePrivateKey(path, repo.KeyPassword)
	if err != nil {
		return err
	}
	client := suite.NewClient(node1Key)
	from, err := node1Key.PublicKey().Address()
	if err != nil {
		return err
	}
	nonce := atomic.AddUint64(&nonce1, 1)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "FreezeAppchain",
		&rpcx.TransactOpts{
			From:  from.String(),
			Nonce: nonce,
		},
		rpcx.String(chainID), rpcx.String("reason"),
	)
	if err != nil {
		return err
	}
	if res.Status == pb.Receipt_FAILED {
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

func (suite *Snake) UpdateAppchain(pk crypto.PrivateKey, id, name, desc string, trustRoot []byte, admins string) error {
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(id),
		rpcx.String(name),
		rpcx.String(desc),
		rpcx.Bytes(trustRoot),
		rpcx.String(admins),
		rpcx.String("reason"),
	}
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "UpdateAppchain", nil, args...)
	if err != nil {
		return err
	}
	if res.Status == pb.Receipt_FAILED {
		return errors.New(string(res.Ret))
	}
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	if err != nil {
		return err
	}
	if result.ProposalID != "" {
		if err := suite.VotePass(result.ProposalID); err != nil {
			return err
		}
	}
	return nil
}

func (suite *Snake) ActivateAppchain(pk crypto.PrivateKey, chainID string) error {
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "ActivateAppchain", nil, rpcx.String(chainID), rpcx.String("reason"))
	if err != nil {
		return err
	}
	if res.Status == pb.Receipt_FAILED {
		return errors.New(string(res.Ret))
	}
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	if err != nil {
		return err
	}
	if result.ProposalID != "" {
		if err := suite.VotePass(result.ProposalID); err != nil {
			return err
		}
	}
	return nil
}

func (suite *Snake) LogoutAppchain(pk crypto.PrivateKey, chainID string) error {
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "LogoutAppchain", nil, rpcx.String(chainID), rpcx.String("reason"))
	if err != nil {
		return err
	}
	if res.Status == pb.Receipt_FAILED {
		return errors.New(string(res.Ret))
	}
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	if err != nil {
		return err
	}
	if result.ProposalID != "" {
		if err := suite.VotePass(result.ProposalID); err != nil {
			return err
		}
	}
	return nil
}

func (suite Snake) ChainToActivating(pk crypto.PrivateKey, name, address string) error {
	err := suite.RegisterAppchain(pk, name, address)
	if err != nil {
		return err
	}
	chainID, err := suite.GetChainIDByName(name)
	if err != nil {
		return err
	}
	err = suite.FreezeAppchain(chainID)
	if err != nil {
		return err
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "ActivateAppchain", nil, rpcx.String(chainID), rpcx.String("reason"))
	if err != nil {
		return err
	}
	if res.Status == pb.Receipt_FAILED {
		return errors.New(string(res.Ret))
	}
	err = suite.CheckChainStatus(chainID, governance.GovernanceActivating)
	if err != nil {
		return err
	}
	return nil
}

func (suite Snake) ChainToUpdating(pk crypto.PrivateKey, chainID, address string) error {
	err := suite.RegisterAppchain(pk, chainID, address)
	if err != nil {
		return err
	}
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(chainID),
		rpcx.String(chainID + "123"),
		rpcx.String("desc"),
		rpcx.Bytes([]byte("")),
		rpcx.String(chainID),
		rpcx.String("reason"),
	}
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "UpdateAppchain", nil, args...)
	if err != nil {
		return err
	}
	if res.Status == pb.Receipt_FAILED {
		return errors.New(string(res.Ret))
	}
	err = suite.CheckChainStatus(chainID, governance.GovernanceUpdating)
	if err != nil {
		return err
	}
	return nil
}

func (suite Snake) ChainToFreezing(pk crypto.PrivateKey, chainID, address string) error {
	err := suite.RegisterAppchain(pk, chainID, address)
	if err != nil {
		return err
	}
	path, err := repo.Node1Path()
	if err != nil {
		return err
	}
	node1Key, err := asym.RestorePrivateKey(path, repo.KeyPassword)
	if err != nil {
		return err
	}
	client := suite.NewClient(node1Key)
	pubAddress, err := node1Key.PublicKey().Address()
	if err != nil {
		return err
	}
	nonce := atomic.AddUint64(&nonce1, 1)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "FreezeAppchain",
		&rpcx.TransactOpts{
			From:  pubAddress.String(),
			Nonce: nonce,
		},
		rpcx.String(chainID), rpcx.String("reason"),
	)
	if err != nil {
		return err
	}
	if res.Status == pb.Receipt_FAILED {
		return errors.New(string(res.Ret))
	}
	err = suite.CheckChainStatus(chainID, governance.GovernanceFreezing)
	if err != nil {
		return err
	}
	return nil
}

func (suite Snake) ChainToFrozen(pk crypto.PrivateKey, chainID, address string) error {
	err := suite.RegisterAppchain(pk, chainID, address)
	if err != nil {
		return err
	}
	err = suite.FreezeAppchain(chainID)
	if err != nil {
		return err
	}
	err = suite.CheckChainStatus(chainID, governance.GovernanceFrozen)
	if err != nil {
		return err
	}
	return nil
}

func (suite Snake) ChainToLogouting(pk crypto.PrivateKey, chainID, address string) error {
	err := suite.RegisterAppchain(pk, chainID, address)
	if err != nil {
		return err
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "LogoutAppchain", nil, rpcx.String(chainID), rpcx.String("reason"))
	if err != nil {
		return err
	}
	if res.Status == pb.Receipt_FAILED {
		return errors.New(string(res.Ret))
	}
	err = suite.CheckChainStatus(chainID, governance.GovernanceLogouting)
	if err != nil {
		return err
	}
	return nil
}

func (suite Snake) ChainToForbidden(pk crypto.PrivateKey, chainID, address string) error {
	err := suite.RegisterAppchain(pk, chainID, address)
	if err != nil {
		return err
	}
	err = suite.LogoutAppchain(pk, chainID)
	if err != nil {
		return err
	}
	err = suite.CheckChainStatus(chainID, governance.GovernanceForbidden)
	if err != nil {
		return err
	}
	return nil
}
