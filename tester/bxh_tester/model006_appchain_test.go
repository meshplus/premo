package bxh_tester

import (
	"encoding/base64"
	"encoding/json"

	appchainmgr "github.com/meshplus/bitxhub-core/appchain-mgr"
	"github.com/meshplus/bitxhub-core/governance"
	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/pkg/errors"
)

type Model6 struct {
	*Snake
}

//tc:注册信息缺失或错误
func (suite *Model6) Test0601_RegisterAppchainLoseFields() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(""),     //validators
		rpcx.String("raft"), //consensus_type
		rpcx.String("1.8"),  //version
	}
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "Register", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(res.Status, pb.Receipt_FAILED)
	suite.Require().Contains(string(res.Ret), "too few input arguments")
}

//tc:应用链未注册，注册应用链
func (suite *Model6) Test0602_RegisterAppchain() {
	_, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	res, err := suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchainmgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceAvailable, appchain.Status)
}

//tc:应用链未注册，注册应用链，投票不通过
func (suite *Model6) Test0603_RegisterAppchainWithReject() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	bytes, err := pk.PublicKey().Bytes()
	suite.Require().Nil(err)
	client, err := rpcx.New(
		rpcx.WithNodesInfo(&rpcx.NodeInfo{Addr: cfg.addrs[0]}),
		rpcx.WithLogger(cfg.logger),
		rpcx.WithPrivateKey(pk),
	)
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(bytes)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":" + pubAddress.String()), //id
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":."),                      //ownerDID
		rpcx.String("/ipfs/QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),                   //docAddr
		rpcx.String("QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),                         //docHash
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain"),         //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "Register", nil, args...)
	suite.Require().Nil(err)
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	suite.Require().Nil(err)
	err = suite.VoteReject(result.ProposalID)
	suite.Require().Nil(err)

	res, err = suite.GetChainStatusById(string(result.Extra))
	suite.Require().Nil(err)
	appchain := &appchainmgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceUnavailable, appchain.Status)
}

//tc:应用链处于注册中状态，注册应用链
func (suite *Model6) Test0604_RegisterAppchainWithRegisting() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	bytes, err := pk.PublicKey().Bytes()
	suite.Require().Nil(err)
	client, err := rpcx.New(
		rpcx.WithNodesInfo(&rpcx.NodeInfo{Addr: cfg.addrs[0]}),
		rpcx.WithLogger(cfg.logger),
		rpcx.WithPrivateKey(pk),
	)
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(bytes)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":" + pubAddress.String()), //id
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":."),                      //ownerDID
		rpcx.String("/ipfs/QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),                   //docAddr
		rpcx.String("QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),                         //docHash
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain"),         //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "Register", nil, args...)
	suite.Require().Nil(err)
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	suite.Require().Nil(err)

	res, err = client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "Register", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	suite.Require().Contains(string(res.Ret), "appchain has registered")

	res, err = suite.GetChainStatusById(string(result.Extra))
	suite.Require().Nil(err)
	appchain := &appchainmgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceRegisting, appchain.Status)
}

//tc:应用链状态已注册，注册应用链
func (suite *Model6) Test0605_RegisterAppchainRepeat() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	bytes, err := pk.PublicKey().Bytes()
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(bytes)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":" + pubAddress.String()), //id
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":."),                      //ownerDID
		rpcx.String("/ipfs/QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),                   //docAddr
		rpcx.String("QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),                         //docHash
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain"),         //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "Register", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Contains(string(res.Ret), "appchain has registered")

	res, err = suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchainmgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceAvailable, appchain.Status)
}

//tc:应用链处于不可用状态，注册应用链
func (suite *Model6) Test0606_RegisterAppchainWithUnavailable() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	bytes, err := pk.PublicKey().Bytes()
	suite.Require().Nil(err)
	client, err := rpcx.New(
		rpcx.WithNodesInfo(&rpcx.NodeInfo{Addr: cfg.addrs[0]}),
		rpcx.WithLogger(cfg.logger),
		rpcx.WithPrivateKey(pk),
	)
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(bytes)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":" + pubAddress.String()), //id
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":."),                      //ownerDID
		rpcx.String("/ipfs/QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),                   //docAddr
		rpcx.String("QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),                         //docHash
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain"),         //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "Register", nil, args...)
	suite.Require().Nil(err)
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	suite.Require().Nil(err)
	err = suite.VoteReject(result.ProposalID)
	suite.Require().Nil(err)

	res, err = suite.GetChainStatusById(string(result.Extra))
	suite.Require().Nil(err)
	appchain := &appchainmgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceUnavailable, appchain.Status)

	_, ChainID, err := suite.RegisterAppchain()
	res, err = suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceAvailable, appchain.Status)
}

//tc:应用链处于更新中状态，注册应用链
func (suite *Model6) Test0607_RegisterAppchainWithUpdating() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	bytes, err := pk.PublicKey().Bytes()
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(bytes)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":."),    //ownerDID
		rpcx.String("/ipfs/QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"), //docAddr
		rpcx.String("QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),       //docHash
		rpcx.String("123"),              //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain"),         //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "UpdateAppchain", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	res, err = suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchainmgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceUpdating, appchain.Status)

	args = []*pb.Arg{
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":" + pubAddress.String()), //id
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":."),                      //ownerDID
		rpcx.String("/ipfs/QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),                   //docAddr
		rpcx.String("QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),                         //docHash
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain111"),      //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	res, err = client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "Register", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Contains(string(res.Ret), "appchain has registered")

	res, err = suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceUpdating, appchain.Status)
}

//tc:应用链处于冻结中状态，注册应用链
func (suite *Model6) Test0608_RegisterAppchainWithFreezing() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "FreezeAppchain", nil, rpcx.String(ChainID))
	suite.Require().Nil(err)

	res, err = suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchainmgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceFreezing, appchain.Status)

	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	bytes, err := pk.PublicKey().Bytes()
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(bytes)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":" + pubAddress.String()), //id
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":."),                      //ownerDID
		rpcx.String("/ipfs/QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),                   //docAddr
		rpcx.String("QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),                         //docHash
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain"),         //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	res, err = client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "Register", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	suite.Require().Contains(string(res.Ret), "appchain has registered")

	res, err = suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceFreezing, appchain.Status)
}

//tc:应用链处于冻结状态，注册应用链
func (suite *Model6) Test0609_RegisterAppchainWithFrozen() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	err = suite.freezeAppchain(pk)
	suite.Require().Nil(err)

	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	bytes, err := pk.PublicKey().Bytes()
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(bytes)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":" + pubAddress.String()), //id
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":."),                      //ownerDID
		rpcx.String("/ipfs/QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),                   //docAddr
		rpcx.String("QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),                         //docHash
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain"),         //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "Register", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Contains(string(res.Ret), "appchain has registered")

	res, err = suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchainmgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceFrozen, appchain.Status)
}

//tc:应用链处于注销中状态，注册应用链
func (suite *Model6) Test0610_RegisterAppchainWithLogouting() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "LogoutAppchain", nil, rpcx.String(ChainID))
	suite.Require().Nil(err)

	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	bytes, err := pk.PublicKey().Bytes()
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(bytes)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":" + pubAddress.String()), //id
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":."),                      //ownerDID
		rpcx.String("/ipfs/QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),                   //docAddr
		rpcx.String("QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),                         //docHash
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain"),         //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	res, err = client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "Register", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Contains(string(res.Ret), "appchain has registered")

	res, err = suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchainmgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceLogouting, appchain.Status)
}

//tc:应用链处于注销状态，注册应用链
func (suite *Model6) Test0611_RegisterAppChainWithForbidden() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	err = suite.logoutAppchain(pk)
	suite.Require().Nil(err)

	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	bytes, err := pk.PublicKey().Bytes()
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(bytes)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":" + pubAddress.String()), //id
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":."),                      //ownerDID
		rpcx.String("/ipfs/QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),                   //docAddr
		rpcx.String("QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),                         //docHash
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain"),         //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "Register", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Contains(string(res.Ret), "appchain has registered")

	res, err = suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchainmgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceForbidden, appchain.Status)
}

//tc:激活信息缺失或错误
func (suite *Model6) Test0612_ActivateAppchainLoseFields() {
	pk, _, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)

	err = suite.freezeAppchain(pk)
	suite.Require().Nil(err)

	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "ActivateAppchain", nil)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	suite.Require().Contains(string(res.Ret), "too few input arguments")
}

//tc:应用链未注册，激活应用链
func (suite *Model6) Test0613_ActivateAppchainWithNoRegister() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.activateAppchain(pk)
	suite.Require().NotNil(err)
}

//tc:应用链处于注册中状态，激活应用链
func (suite *Model6) Test0614_ActivateAppchainWithRegisting() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	bytes, err := pk.PublicKey().Bytes()
	suite.Require().Nil(err)
	client, err := rpcx.New(
		rpcx.WithNodesInfo(&rpcx.NodeInfo{Addr: cfg.addrs[0]}),
		rpcx.WithLogger(cfg.logger),
		rpcx.WithPrivateKey(pk),
	)
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(bytes)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":" + pubAddress.String()), //id
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":."),                      //ownerDID
		rpcx.String("/ipfs/QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),                   //docAddr
		rpcx.String("QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),                         //docHash
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain"),         //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "Register", nil, args...)
	suite.Require().Nil(err)
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	suite.Require().Nil(err)

	err = suite.activateAppchain(pk)
	suite.Require().NotNil(err)
}

//tc:应用链状态已注册，激活应用链
func (suite *Model6) Test0615_ActivateAppchainWithAvailable() {
	pk, _, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	err = suite.activateAppchain(pk)
	suite.Require().NotNil(err)
}

//tc:应用链处于更新中状态，激活应用链
func (suite *Model6) Test0616_ActivateAppchainWithUpdating() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)

	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	bytes, err := pk.PublicKey().Bytes()
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(bytes)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":."),    //ownerDID
		rpcx.String("/ipfs/QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"), //docAddr
		rpcx.String("QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),       //docHash
		rpcx.String("111"),              //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain"),         //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "UpdateAppchain", nil, args...)
	suite.Require().Nil(err)

	res, err = suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchainmgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceUpdating, appchain.Status)

	err = suite.activateAppchain(pk)
	suite.Require().NotNil(err)
}

//tc:应用链处于冻结中状态，激活应用链
func (suite *Model6) Test0617_ActivateAppchainWithFreezing() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)

	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "FreezeAppchain", nil, rpcx.String(ChainID))
	suite.Require().Nil(err)

	res, err = suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchainmgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceFreezing, appchain.Status)

	err = suite.activateAppchain(pk)
	suite.Require().NotNil(err)
}

//tc:应用链处于冻结状态，激活应用链
func (suite *Model6) Test0618_ActivateAppchain() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	err = suite.freezeAppchain(pk)
	suite.Require().Nil(err)

	err = suite.activateAppchain(pk)
	suite.Require().Nil(err)

	res, err := suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchainmgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceAvailable, appchain.Status)
}

//tc:tc:应用链处于冻结状态，激活应用链,投票不通过
func (suite *Model6) Test0619_ActivateAppchainWithReject() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)

	err = suite.freezeAppchain(pk)
	suite.Require().Nil(err)

	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "ActivateAppchain", nil, rpcx.String(ChainID))
	suite.Require().Nil(err)

	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	suite.Require().Nil(err)
	err = suite.VoteReject(result.ProposalID)
	suite.Require().Nil(err)

	res, err = suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchainmgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceFrozen, appchain.Status)
}

//tc:应用链处于注销中状态，激活应用链
func (suite *Model6) Test0620_ActivateAppchainWithLogouting() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "LogoutAppchain", nil, rpcx.String(ChainID))
	suite.Require().Nil(err)

	err = suite.activateAppchain(pk)
	suite.Require().NotNil(err)

	res, err = suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchainmgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceLogouting, appchain.Status)
}

//tc:应用链处于注销状态，激活应用链
func (suite *Model6) Test0621_ActivateAppchainWithForbidden() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	err = suite.logoutAppchain(pk)
	suite.Require().Nil(err)

	res, err := suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchainmgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceForbidden, appchain.Status)

	err = suite.activateAppchain(pk)
	suite.Require().NotNil(err)
}

//tc:更新信息缺失或错误
func (suite *Model6) Test0622_UpdateAppchainLoseFields() {
	pk, _, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	bytes, err := pk.PublicKey().Bytes()
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(bytes)
	args := []*pb.Arg{
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "UpdateAppchain", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	suite.Require().Contains(string(res.Ret), "too few input arguments")
}

//tc:应用链处于注册中的状态，更新应用链
func (suite *Model6) Test0623_UpdateAppchainWithRegisting() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	bytes, err := pk.PublicKey().Bytes()
	suite.Require().Nil(err)
	client, err := rpcx.New(
		rpcx.WithNodesInfo(&rpcx.NodeInfo{Addr: cfg.addrs[0]}),
		rpcx.WithLogger(cfg.logger),
		rpcx.WithPrivateKey(pk),
	)
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(bytes)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":" + pubAddress.String()), //id
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":."),                      //ownerDID
		rpcx.String("/ipfs/QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),                   //docAddr
		rpcx.String("QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),                         //docHash
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain"),         //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "Register", nil, args...)
	suite.Require().Nil(err)
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	suite.Require().Nil(err)

	args[6] = rpcx.String("AppChain111111111")
	err = suite.updateAppchain(pk, args...)
	suite.Require().NotNil(err)
}

//tc:应用链状态已注册，更新应用链
func (suite *Model6) Test0624_UpdateAppchain() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	bytes, err := pk.PublicKey().Bytes()
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(bytes)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":."),    //ownerDID
		rpcx.String("/ipfs/QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"), //docAddr
		rpcx.String("QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),       //docHash
		rpcx.String("111"),              //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain"),         //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "UpdateAppchain", nil, args...)
	suite.Require().Nil(err)
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	suite.Require().Nil(err)
	err = suite.VotePass(result.ProposalID)
	suite.Require().Nil(err)

	res, err = suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchainmgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceAvailable, appchain.Status)
	suite.Require().Equal("AppChain", appchain.Name)
}

//tc:应用链状态已注册，更新应用链,投票不通过
func (suite *Model6) Test0625_UpdateAppchainWithReject() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	bytes, err := pk.PublicKey().Bytes()
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(bytes)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":."),    //ownerDID
		rpcx.String("/ipfs/QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"), //docAddr
		rpcx.String("QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),       //docHash
		rpcx.String("111"),              //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain"),         //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "UpdateAppchain", nil, args...)
	suite.Require().Nil(err)
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	suite.Require().Nil(err)
	err = suite.VoteReject(result.ProposalID)
	suite.Require().Nil(err)

	res, err = suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchainmgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceAvailable, appchain.Status)
	suite.Require().Equal("AppChain", appchain.Name)
}

//tc:应用链处于更新中的状态，更新应用链
func (suite *Model6) Test0626_UpdateAppchainWithUpdating() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	bytes, err := pk.PublicKey().Bytes()
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(bytes)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":."),    //ownerDID
		rpcx.String("/ipfs/QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"), //docAddr
		rpcx.String("QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),       //docHash
		rpcx.String("111"),              //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain"),         //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	_, err = client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "UpdateAppchain", nil, args...)
	suite.Require().Nil(err)

	res, err := suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchainmgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceUpdating, appchain.Status)

	err = suite.updateAppchain(pk, args...)
	suite.Require().NotNil(err)
}

//tc:应用链处于冻结中的状态，更新应用链
func (suite *Model6) Test0627_UpdateAppchainWithFreezing() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)

	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "FreezeAppchain", nil, rpcx.String(ChainID))
	suite.Require().Nil(err)

	res, err = suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchainmgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceFreezing, appchain.Status)

	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	bytes, err := pk.PublicKey().Bytes()
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(bytes)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":" + pubAddress.String()), //id
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":."),                      //ownerDID
		rpcx.String("/ipfs/QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),                   //docAddr
		rpcx.String("QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),                         //docHash
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain"),         //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	err = suite.updateAppchain(pk, args...)
	suite.Require().NotNil(err)
}

//tc:应用链处于冻结的状态，更新应用链
func (suite *Model6) Test0628_UpdateAppchainWithFrozen() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	err = suite.freezeAppchain(pk)
	suite.Require().Nil(err)

	res, err := suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchainmgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceFrozen, appchain.Status)

	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	bytes, err := pk.PublicKey().Bytes()
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(bytes)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":" + pubAddress.String()), //id
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":."),                      //ownerDID
		rpcx.String("/ipfs/QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),                   //docAddr
		rpcx.String("QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),                         //docHash
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain"),         //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	err = suite.updateAppchain(pk, args...)
	suite.Require().NotNil(err)
}

//tc:应用链处于注销中状态，更新应用链
func (suite *Model6) Test0629_UpdateAppchainWithWithLogouting() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "LogoutAppchain", nil, rpcx.String(ChainID))
	suite.Require().Nil(err)

	res, err = suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchainmgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceLogouting, appchain.Status)

	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	bytes, err := pk.PublicKey().Bytes()
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(bytes)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":" + pubAddress.String()), //id
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":."),                      //ownerDID
		rpcx.String("/ipfs/QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),                   //docAddr
		rpcx.String("QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),                         //docHash
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain"),         //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	err = suite.updateAppchain(pk, args...)
	suite.Require().NotNil(err)
}

//tc:应用链处于注销状态，更新应用链
func (suite *Model6) Test0630_UpdateAppchainWithUnavailable() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	err = suite.logoutAppchain(pk)
	suite.Require().Nil(err)
	res, err := suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchainmgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceForbidden, appchain.Status)

	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	bytes, err := pk.PublicKey().Bytes()
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(bytes)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":" + pubAddress.String()), //id
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":."),                      //ownerDID
		rpcx.String("/ipfs/QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),                   //docAddr
		rpcx.String("QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),                         //docHash
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain"),         //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	err = suite.updateAppchain(pk, args...)
	suite.Require().NotNil(err)
}

//tc:冻结信息缺失或错误
func (suite *Model6) Test0631_FreezeAppchainLoseFields() {
	pk, _, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "FreezeAppchain", nil)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	suite.Require().Contains(string(res.Ret), "too few input arguments")
}

//tc:应用链处于注册中的状态，冻结应用链
func (suite *Model6) Test0632_FreezeAppchainWithRegisting() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	bytes, err := pk.PublicKey().Bytes()
	suite.Require().Nil(err)
	client, err := rpcx.New(
		rpcx.WithNodesInfo(&rpcx.NodeInfo{Addr: cfg.addrs[0]}),
		rpcx.WithLogger(cfg.logger),
		rpcx.WithPrivateKey(pk),
	)
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(bytes)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":" + pubAddress.String()), //id
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":."),                      //ownerDID
		rpcx.String("/ipfs/QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),                   //docAddr
		rpcx.String("QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),                         //docHash
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain"),         //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "Register", nil, args...)
	suite.Require().Nil(err)
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	suite.Require().Nil(err)

	err = suite.freezeAppchain(pk)
	suite.Require().NotNil(err)
}

//tc:应用链状态已注册，冻结应用链
func (suite *Model6) Test0633_FreezeAppchain() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	err = suite.freezeAppchain(pk)
	suite.Require().Nil(err)

	res, err := suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchainmgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceFrozen, appchain.Status)
}

//tc:应用链状态已注册，冻结应用链，投票不通过
func (suite *Model6) Test0634_FreezeAppchainWithReject() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)

	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "FreezeAppchain", nil, rpcx.String(ChainID))
	suite.Require().Nil(err)
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	suite.Require().Nil(err)
	err = suite.VoteReject(result.ProposalID)
	suite.Require().Nil(err)

	res, err = suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchainmgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceAvailable, appchain.Status)
}

//tc:应用链处于更新中的状态，冻结应用链
func (suite *Model6) Test0635_FreezeAppchainWithUpdating() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	bytes, err := pk.PublicKey().Bytes()
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(bytes)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":."),    //ownerDID
		rpcx.String("/ipfs/QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"), //docAddr
		rpcx.String("QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),       //docHash
		rpcx.String("111"),              //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain"),         //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "UpdateAppchain", nil, args...)
	suite.Require().Nil(err)

	res, err = suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchainmgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceUpdating, appchain.Status)

	err = suite.freezeAppchain(pk)
	suite.Require().Nil(err)
}

//tc:应用链处于冻结中的状态，冻结应用链
func (suite *Model6) Test0636_FreezeAppchainWithFreezing() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)

	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "FreezeAppchain", nil, rpcx.String(ChainID))
	suite.Require().Nil(err)

	res, err = suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchainmgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceFreezing, appchain.Status)

	err = suite.freezeAppchain(pk)
	suite.Require().NotNil(err)
}

//tc:应用链处于冻结的状态，冻结应用链
func (suite *Model6) Test0637_FreezeAppchainWithFrozen() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	err = suite.freezeAppchain(pk)
	suite.Require().Nil(err)

	err = suite.freezeAppchain(pk)
	suite.Require().NotNil(err)

	res, err := suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchainmgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceFrozen, appchain.Status)
}

//tc:应用链处于注销中状态，冻结应用链
func (suite *Model6) Test0638_FreezeAppchainWithWithLogouting() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "LogoutAppchain", nil, rpcx.String(ChainID))
	suite.Require().Nil(err)

	res, err = suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchainmgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceLogouting, appchain.Status)

	err = suite.freezeAppchain(pk)
	suite.Require().NotNil(err)
}

//tc:应用链处于注销状态，冻结应用链
func (suite *Model6) Test0639_FreezeAppchainWithForbidden() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	err = suite.logoutAppchain(pk)
	suite.Require().Nil(err)
	res, err := suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchainmgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceForbidden, appchain.Status)

	err = suite.freezeAppchain(pk)
	suite.Require().NotNil(err)
}

//tc:应用链处于注册中的状态，注销应用链
func (suite *Model6) Test0640_LogoutAppchainWithRegisting() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	bytes, err := pk.PublicKey().Bytes()
	suite.Require().Nil(err)
	client, err := rpcx.New(
		rpcx.WithNodesInfo(&rpcx.NodeInfo{Addr: cfg.addrs[0]}),
		rpcx.WithLogger(cfg.logger),
		rpcx.WithPrivateKey(pk),
	)
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(bytes)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":" + pubAddress.String()), //id
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":."),                      //ownerDID
		rpcx.String("/ipfs/QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),                   //docAddr
		rpcx.String("QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),                         //docHash
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain"),         //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "Register", nil, args...)
	suite.Require().Nil(err)
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	suite.Require().Nil(err)

	err = suite.logoutAppchain(pk)
	suite.Require().NotNil(err)
}

//tc:应用链状态已注册，注销应用链
func (suite *Model6) Test0641_LogoutAppchain() {
	pk, _, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	err = suite.logoutAppchain(pk)
	suite.Require().Nil(err)
}

//tc:应用链状态已注册，注销应用链,投票不通过
func (suite *Model6) Test0642_LogoutAppchainWithReject() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "LogoutAppchain", nil, rpcx.String(ChainID))
	suite.Require().Nil(err)
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	suite.Require().Nil(err)
	err = suite.VoteReject(result.ProposalID)
	suite.Require().Nil(err)

	res, err = suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchainmgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceAvailable, appchain.Status)
}

//tc:应用链处于更新中的状态，注销应用链
func (suite *Model6) Test0643_LogoutAppchainWithUpdating() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	bytes, err := pk.PublicKey().Bytes()
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(bytes)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":."),    //ownerDID
		rpcx.String("/ipfs/QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"), //docAddr
		rpcx.String("QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),       //docHash
		rpcx.String("111"),              //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain"),         //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	_, err = client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "UpdateAppchain", nil, args...)
	suite.Require().Nil(err)

	res, err := suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchainmgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceUpdating, appchain.Status)

	err = suite.logoutAppchain(pk)
	suite.Require().Nil(err)
}

//tc:应用链处于冻结中的状态，注销应用链
func (suite *Model6) Test0644_LogoutAppchainWithFreezing() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)

	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "FreezeAppchain", nil, rpcx.String(ChainID))
	suite.Require().Nil(err)

	res, err = suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchainmgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceFreezing, appchain.Status)

	err = suite.logoutAppchain(pk)
	suite.Require().Nil(err)
}

//tc:应用链处于冻结的状态，注销应用链
func (suite *Model6) Test0645_LogoutAppchainWithFrozen() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	err = suite.freezeAppchain(pk)

	res, err := suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchainmgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceFrozen, appchain.Status)

	err = suite.logoutAppchain(pk)
	suite.Require().Nil(err)
}

//tc:应用链处于注销中状态，注销应用链
func (suite *Model6) Test0646_LogoutAppchainWithLogouting() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "LogoutAppchain", nil, rpcx.String(ChainID))
	suite.Require().Nil(err)

	res, err = suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchainmgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceLogouting, appchain.Status)

	err = suite.logoutAppchain(pk)
	suite.Require().NotNil(err)
}

//tc:应用链处于注销状态，注销应用链
func (suite *Model6) Test0647_LogoutAppchainWithUnavailable() {
	pk, _, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	err = suite.logoutAppchain(pk)
	suite.Require().Nil(err)

	err = suite.logoutAppchain(pk)
	suite.Require().NotNil(err)
}

//tc:根据指定ID查询应用链信息
func (suite *Model6) Test0648_GetAppchainByID() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	bytes, err := pk.PublicKey().Bytes()
	suite.Require().Nil(err)

	client, err := rpcx.New(
		rpcx.WithNodesInfo(&rpcx.NodeInfo{Addr: cfg.addrs[0]}),
		rpcx.WithLogger(cfg.logger),
		rpcx.WithPrivateKey(pk),
	)
	suite.Require().Nil(err)

	var pubKeyStr = base64.StdEncoding.EncodeToString(bytes)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":" + pubAddress.String()), //id
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":."),                      //ownerDID
		rpcx.String("/ipfs/QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),                   //docAddr
		rpcx.String("QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),                         //docHash
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain"),         //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "Register", nil, args...)
	suite.Require().Nil(err)
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	suite.Require().Nil(err)
	suite.Require().NotNil(string(result.Extra))
	err = suite.VotePass(result.ProposalID)
	suite.Require().Nil(err)

	res, err = suite.GetChainStatusById(string(result.Extra))
	suite.Require().Nil(err)
	appchain := &appchainmgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceAvailable, appchain.Status)

	args = []*pb.Arg{
		rpcx.String(string(result.Extra)),
	}
	res, err = client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "GetAppchain", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(res.Status, pb.Receipt_SUCCESS)
	suite.Require().NotNil(res.Ret)
}

//tc:根据错误的ID查询应用链信息
func (suite *Model6) Test0649_GetAppchainByErrorID() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(suite.from.String() + "123"),
	}
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "GetAppchain", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	suite.Require().Equal("call error: this appchain does not exist", string(res.Ret))
}

func (suite *Snake) freezeAppchain(pk crypto.PrivateKey) error {
	pubAddress, err := pk.PublicKey().Address()
	if err != nil {
		return err
	}
	client := suite.NewClient(pk)
	ChainID := "did:bitxhub:appchain" + pubAddress.String() + ":."
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "FreezeAppchain", nil, rpcx.String(ChainID))
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

func (suite *Snake) updateAppchain(pk crypto.PrivateKey, args ...*pb.Arg) error {
	client := suite.NewClient(pk)
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
	err = suite.VotePass(result.ProposalID)
	if err != nil {
		return err
	}
	return nil
}

func (suite *Snake) activateAppchain(pk crypto.PrivateKey) error {
	pubAddress, err := pk.PublicKey().Address()
	if err != nil {
		return err
	}
	client := suite.NewClient(pk)
	ChainID := "did:bitxhub:appchain" + pubAddress.String() + ":."
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "ActivateAppchain", nil, rpcx.String(ChainID))
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

func (suite *Snake) logoutAppchain(pk crypto.PrivateKey) error {
	pubAddress, err := pk.PublicKey().Address()
	if err != nil {
		return err
	}
	client := suite.NewClient(pk)
	ChainID := "did:bitxhub:appchain" + pubAddress.String() + ":."
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "LogoutAppchain", nil, rpcx.String(ChainID))
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
