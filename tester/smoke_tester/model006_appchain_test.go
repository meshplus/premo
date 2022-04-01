package bxh_tester

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"

	appchainmgr "github.com/meshplus/bitxhub-core/appchain-mgr"
	"github.com/meshplus/bitxhub-core/governance"
	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/meshplus/premo/internal/repo"
)

type Model6 struct {
	*Snake
}

//tc：通过正确的参数注册应用链，应用链注册成功
func (suite Model6) Test0601_RegisterAppchainIsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
}

//tc：通过曾被占用的应用链名称注册应用链，应用链注册成功
func (suite Model6) Test0602_RegisterAppchainWithFreeNameIsSuccess() {
	pk1, from1, address1, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk1, from1, address1)
	suite.Require().Nil(err)
	err = suite.UpdateAppchain(pk1, from1, from1+"123", "desc", []byte(""), from1)
	suite.Require().Nil(err)
	pk2, _, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk2, from1, address2)
	suite.Require().Nil(err)
}

//tc：通过曾被占用的管理员地址注册应用链，应用链注册成功
func (suite Model6) Test0603_RegisterAppchainWithFreeAdminIsSuccess() {
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

//tc：通过正确的参数更新应用链，应用链更新成功
func (suite Model6) Test0604_UpdateAppchainIsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	err = suite.UpdateAppchain(pk, from, from+"123", "desc", []byte(""), from)
	suite.Require().Nil(err)
}

//tc：通过曾被占用的应用链名称更新应用链，应用链更新成功
func (suite Model6) Test0605_UpdateAppchainWithFreeNameIsSuccess() {
	pk1, from1, address1, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk1, from1, address1)
	suite.Require().Nil(err)
	err = suite.UpdateAppchain(pk1, from1, from1+"123", "desc", []byte(""), from1)
	suite.Require().Nil(err)
	pk2, from2, address2, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk2, from2, address2)
	suite.Require().Nil(err)
	err = suite.UpdateAppchain(pk2, from2, from1, "desc", []byte(""), from2)
	suite.Require().Nil(err)
}

//tc：通过曾被占用的管理员地址更新应用链，应用链更新成功
func (suite Model6) Test0606_UpdateAppchainWithFreeAdminIsSuccess() {
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
	err = suite.UpdateAppchain(pk1, from1, from1, "desc", []byte(""), from1)
	suite.Require().Nil(err)
	pk3, from3, address3, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk3, from3, address3)
	suite.Require().Nil(err)
	err = suite.UpdateAppchain(pk3, from3, from3, "desc", []byte(""), from3+","+from2.String())
	suite.Require().Nil(err)
}

//tc：应用链处于frozen状态更新应用链，应用链更新成功
func (suite Model6) Test0607_UpdateAppchainWithFrozenChainIsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToFrozen(pk, from, address)
	suite.Require().Nil(err)
	err = suite.UpdateAppchain(pk, from, from, "desc", []byte(""), from)
	suite.Require().Nil(err)
}

//tc：应用链更新名称字段，产生提案
func (suite Model6) Test0608_UpdateAppchainWithNameFieldHaveProposalIsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(from),
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
func (suite Model6) Test0609_UpdateAppchainWithTrustRootFieldHaveProposalIsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(from),
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
func (suite Model6) Test0610_UpdateAppchainWithAdminsFieldHaveProposalIsSuccess() {
	pk1, from1, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk1, from1, address)
	suite.Require().Nil(err)
	pk2, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	from2, err := pk2.PublicKey().Address()
	suite.Require().Nil(err)
	client := suite.NewClient(pk1)
	args := []*pb.Arg{
		rpcx.String(from1),
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
func (suite Model6) Test0611_UpdateAppchainWithDescFieldNoProposalIsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(from),
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

//tc：中继链管理员冻结应用链，应用链冻结成功
func (suite Model6) Test0612_FreezeAppchainIsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	err = suite.FreezeAppchain(from)
	suite.Require().Nil(err)
}

//tc：应用链管理员激活应用链，应用链激活成功
func (suite Model6) Test0613_ActivateAppchainIsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToFrozen(pk, from, address)
	suite.Require().Nil(err)
	err = suite.ActivateAppchain(pk, from)
	suite.Require().Nil(err)
}

//tc：中继链管理员激活应用链，应用链激活成功
func (suite Model6) Test0614_ActivateAppchainWithRelayAdminIsSuccess() {
	pk, chainID, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToFrozen(pk, chainID, address)
	suite.Require().Nil(err)
	node1pk, from, err := repo.Node1Priv()
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

//tc：应用链管理员注销应用链，应用链注销成功
func (suite Model6) Test0615_LogoutAppchainIsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	err = suite.LogoutAppchain(pk, from)
	suite.Require().Nil(err)
}

//tc：应用链处于updating状态注销应用链，应用链注销成功
func (suite Model6) Test0616_LogoutAppchainWithUpdatingCainIsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToUpdating(pk, from, address)
	suite.Require().Nil(err)
	err = suite.LogoutAppchain(pk, from)
	suite.Require().Nil(err)
}

//tc：应用链处于activating状态注销应用链，应用链注销成功
func (suite Model6) Test0617_LogoutAppchainWithActivatingChainIsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToActivating(pk, from, address)
	suite.Require().Nil(err)
	err = suite.LogoutAppchain(pk, from)
	suite.Require().Nil(err)
}

//tc：应用链处于freezing状态注销应用链，应用链注销成功
func (suite Model6) Test0618_LogoutAppchainWithFreezingChainIsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToFreezing(pk, from, address)
	suite.Require().Nil(err)
	err = suite.LogoutAppchain(pk, from)
	suite.Require().Nil(err)
}

//tc：应用链处于frozen状态注销应用链，应用链注销成功
func (suite Model6) Test0619_LogoutAppchainWithFrozenChainIsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.ChainToFrozen(pk, from, address)
	suite.Require().Nil(err)
	err = suite.LogoutAppchain(pk, from)
	suite.Require().Nil(err)
}

//tc：更新主验证规则中，应用链处于frozen状态，注销应用链，应用链注销成功
func (suite Model6) Test0620_LogoutAppchainWithUpdateMasterRuleIsSuccess() {
	pk, from, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, from, address)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(from),
		rpcx.String(HappyRuleAddr),
		rpcx.String("reason"),
	}
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), UpdateMasterRule, nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	err = suite.CheckChainStatus(from, governance.GovernanceFrozen)
	suite.Require().Nil(err)
	err = suite.LogoutAppchain(pk, from)
	suite.Require().Nil(err)
}

// RegisterAppchain register appchain
func (suite *Snake) RegisterAppchain(pk crypto.PrivateKey, name, address string) error {
	client := suite.NewClient(pk)
	from, err := pk.PublicKey().Address()
	if err != nil {
		return err
	}
	args := []*pb.Arg{
		rpcx.String(from.String()),  //chainID
		rpcx.String(name),           //chainName
		rpcx.String("Flato V1.0.3"), //chainType
		rpcx.Bytes([]byte("")),      //trustRoot
		rpcx.String("0x857133c5C69e6Ce66F7AD46F200B9B3573e77582"), //broker
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
		return fmt.Errorf(string(res.Ret))
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

// RegisterAppchainWithType register appchain with type
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
		return fmt.Errorf(string(res.Ret))
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

// FreezeAppchain freeze appchain
func (suite *Snake) FreezeAppchain(chainID string) error {
	node1Key, _, err := repo.Node1Priv()
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
		return fmt.Errorf(string(res.Ret))
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

// UpdateAppchain updated appchain
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
		return fmt.Errorf(string(res.Ret))
	}
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	if err != nil {
		return err
	}
	if result.ProposalID != "" {
		err = suite.VotePass(result.ProposalID)
		if err != nil {
			return err
		}
	}
	return nil
}

// ActivateAppchain activated appchain
func (suite *Snake) ActivateAppchain(pk crypto.PrivateKey, chainID string) error {
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "ActivateAppchain", nil, rpcx.String(chainID), rpcx.String("reason"))
	if err != nil {
		return err
	}
	if res.Status == pb.Receipt_FAILED {
		return fmt.Errorf(string(res.Ret))
	}
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	if err != nil {
		return err
	}
	if result.ProposalID != "" {
		err = suite.VotePass(result.ProposalID)
		if err != nil {
			return err
		}
	}
	return nil
}

// LogoutAppchain logout appchain
func (suite *Snake) LogoutAppchain(pk crypto.PrivateKey, chainID string) error {
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "LogoutAppchain", nil, rpcx.String(chainID), rpcx.String("reason"))
	if err != nil {
		return err
	}
	if res.Status == pb.Receipt_FAILED {
		return fmt.Errorf(string(res.Ret))
	}
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	if err != nil {
		return err
	}
	if result.ProposalID != "" {
		err = suite.VotePass(result.ProposalID)
		if err != nil {
			return err
		}
	}
	return nil
}

// ChainToActivating get an activating appchain
func (suite Snake) ChainToActivating(pk crypto.PrivateKey, name, address string) error {
	err := suite.RegisterAppchain(pk, name, address)
	if err != nil {
		return err
	}
	err = suite.FreezeAppchain(name)
	if err != nil {
		return err
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "ActivateAppchain", nil, rpcx.String(name), rpcx.String("reason"))
	if err != nil {
		return err
	}
	if res.Status == pb.Receipt_FAILED {
		return fmt.Errorf(string(res.Ret))
	}
	err = suite.CheckChainStatus(name, governance.GovernanceActivating)
	if err != nil {
		return err
	}
	return nil
}

// ChainToUpdating get a updating appchain
func (suite Snake) ChainToUpdating(pk crypto.PrivateKey, name, address string) error {
	err := suite.RegisterAppchain(pk, name, address)
	if err != nil {
		return err
	}
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(name),
		rpcx.String(name + "123"),
		rpcx.String("desc"),
		rpcx.Bytes([]byte("")),
		rpcx.String(name),
		rpcx.String("reason"),
	}
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "UpdateAppchain", nil, args...)
	if err != nil {
		return err
	}
	if res.Status == pb.Receipt_FAILED {
		return fmt.Errorf(string(res.Ret))
	}
	err = suite.CheckChainStatus(name, governance.GovernanceUpdating)
	if err != nil {
		return err
	}
	return nil
}

// ChainToFreezing get a freezing appchain
func (suite Snake) ChainToFreezing(pk crypto.PrivateKey, name, address string) error {
	err := suite.RegisterAppchain(pk, name, address)
	if err != nil {
		return err
	}
	node1, node1Addr, err := repo.Node1Priv()
	if err != nil {
		return err
	}
	client := suite.NewClient(node1)
	nonce := atomic.AddUint64(&nonce1, 1)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "FreezeAppchain",
		&rpcx.TransactOpts{
			From:  node1Addr.String(),
			Nonce: nonce,
		},
		rpcx.String(name), rpcx.String("reason"),
	)
	if err != nil {
		return err
	}
	if res.Status == pb.Receipt_FAILED {
		return fmt.Errorf(string(res.Ret))
	}
	err = suite.CheckChainStatus(name, governance.GovernanceFreezing)
	if err != nil {
		return err
	}
	return nil
}

// ChainToFrozen get a frozen appchain
func (suite Snake) ChainToFrozen(pk crypto.PrivateKey, name, address string) error {
	err := suite.RegisterAppchain(pk, name, address)
	if err != nil {
		return err
	}
	err = suite.FreezeAppchain(name)
	if err != nil {
		return err
	}
	err = suite.CheckChainStatus(name, governance.GovernanceFrozen)
	if err != nil {
		return err
	}
	return nil
}

// ChainToLogouting get a logouting appchain
func (suite Snake) ChainToLogouting(pk crypto.PrivateKey, name, address string) error {
	err := suite.RegisterAppchain(pk, name, address)
	if err != nil {
		return err
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "LogoutAppchain", nil, rpcx.String(name), rpcx.String("reason"))
	if err != nil {
		return err
	}
	if res.Status == pb.Receipt_FAILED {
		return fmt.Errorf(string(res.Ret))
	}
	err = suite.CheckChainStatus(name, governance.GovernanceLogouting)
	if err != nil {
		return err
	}
	return nil
}

// ChainToForbidden get a forbidden appchain
func (suite Snake) ChainToForbidden(pk crypto.PrivateKey, name, address string) error {
	err := suite.RegisterAppchain(pk, name, address)
	if err != nil {
		return err
	}
	err = suite.LogoutAppchain(pk, name)
	if err != nil {
		return err
	}
	err = suite.CheckChainStatus(name, governance.GovernanceForbidden)
	if err != nil {
		return err
	}
	return nil
}

// CheckChainStatus check chain status
func (suite *Snake) CheckChainStatus(name string, expectStatus governance.GovernanceStatus) error {
	status, err := suite.GetChainStatusByName(name)
	if err != nil {
		return err
	}
	if expectStatus != status {
		return fmt.Errorf("expect status is %s ,but get status %s", expectStatus, status)
	}
	return nil
}

// GetChainStatusByName return chain status by name
func (suite *Snake) GetChainStatusByName(name string) (governance.GovernanceStatus, error) {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	if err != nil {
		return "", err
	}
	client := suite.NewClient(pk)
	from, err := pk.PublicKey().Address()
	if err != nil {
		return "", err
	}
	args := []*pb.Arg{
		rpcx.String(name),
	}
	invokePayload := &pb.InvokePayload{
		Method: "GetAppchainByName",
		Args:   args,
	}
	payload, err := invokePayload.Marshal()
	if err != nil {
		return "", err
	}
	data := &pb.TransactionData{
		Type:    pb.TransactionData_INVOKE,
		VmType:  pb.TransactionData_BVM,
		Payload: payload,
	}
	payload, err = data.Marshal()
	if err != nil {
		return "", err
	}
	tx := &pb.BxhTransaction{
		From:      from,
		To:        constant.AppchainMgrContractAddr.Address(),
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}
	res, err := client.SendTransactionWithReceipt(tx, nil)
	if err != nil {
		return "", err
	}
	if res.Status == pb.Receipt_FAILED {
		return "", fmt.Errorf(string(res.Ret))
	}
	appchain := appchainmgr.Appchain{}
	err = json.Unmarshal(res.Ret, &appchain)
	if err != nil {
		return "", err
	}
	return appchain.Status, nil
}
