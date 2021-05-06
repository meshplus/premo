package bxh_tester

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	appchain_mgr "github.com/meshplus/bitxhub-core/appchain-mgr"
	"github.com/meshplus/bitxhub-core/governance"
	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/pkg/errors"
)

//tc:注册信息缺失或错误
func (suite *Snake) Test0601_RegisterAppchainLoseFields() {
	args := []*pb.Arg{
		rpcx.String(""),     //validators
		rpcx.String("raft"), //consensus_type
		rpcx.String("1.8"),  //version
	}
	res, err := suite.client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "Register", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(res.Status, pb.Receipt_FAILED)
	suite.Require().Contains(string(res.Ret), "too few input arguments")
}

//tc:应用链未注册，注册应用链
func (suite *Snake) Test0602_RegisterAppchain() {
	_, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	res, err := suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchain_mgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceAvailable, appchain.Status)
}

//tc:应用链处于注册中状态，注册应用链
func (suite Snake) Test0603_RegisterAppchainWithRegisting() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	client, err := rpcx.New(
		rpcx.WithNodesInfo(&rpcx.NodeInfo{Addr: cfg.addrs[0]}),
		rpcx.WithLogger(cfg.logger),
		rpcx.WithPrivateKey(pk),
	)
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(pubAddress.Bytes())
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
	suite.Require().Contains(string(res.Ret), "under status: ApplyAudit")

	res, err = suite.GetChainStatusById(result.ChainID)
	suite.Require().Nil(err)
	appchain := &appchain_mgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceRegisting, appchain.Status)
}

//tc:应用链状态已注册，注册应用链
func (suite *Snake) Test0604_RegisterAppchainRepeat() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(pubAddress.Bytes())
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
	suite.Require().Contains(string(res.Ret), "under status: Normal")

	res, err = suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchain_mgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceAvailable, appchain.Status)
}

//tc:应用链处于更新中状态，注册应用链
func (suite *Snake) Test0605_RegisterAppchainWithUpdating() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(pubAddress.Bytes())
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":."),    //ownerDID
		rpcx.String("/ipfs/QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"), //docAddr
		rpcx.String("QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),       //docHash
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain111"),      //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "UpdateAppchain", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	res, err = suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchain_mgr.Appchain{}
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
	fmt.Println(string(res.Ret))
	suite.Require().Contains(string(res.Ret), "under status: Updating")

	res, err = suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceUpdating, appchain.Status)
}

//tc:应用链处于冻结中状态，注册应用链
func (suite *Snake) Test0606_RegisterAppchainWithFreezing() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	client := suite.NewClient(pk)
	res, err := suite.client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "FreezeAppchain", nil, rpcx.String(ChainID))
	suite.Require().Nil(err)

	res, err = suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchain_mgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceFreezing, appchain.Status)

	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(pubAddress.Bytes())
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
	suite.Require().Contains(string(res.Ret), "under status: Freezing")
	fmt.Println(string(res.Ret))

	res, err = suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceFreezing, appchain.Status)
}

//tc:应用链处于冻结状态，注册应用链
func (suite *Snake) Test0607_RegisterAppchainWithFrozen() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	err = suite.freezeAppchain(ChainID)
	suite.Require().Nil(err)

	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(pubAddress.Bytes())
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
	suite.Require().Contains(string(res.Ret), "under status: Frozen")

	res, err = suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchain_mgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceFrozen, appchain.Status)
}

//tc:应用链处于注销中状态，注册应用链
func (suite *Snake) Test0608_RegisterAppchainWithLogouting() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "LogoutAppchain", nil, rpcx.String(ChainID))
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))

	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(pubAddress.Bytes())
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
	suite.Require().Contains(string(res.Ret), "under status: Logouting")

	res, err = suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchain_mgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceLogouting, appchain.Status)
}

//tc:应用链处于注销状态，注册应用链
func (suite *Snake) Test0609_RegisterAppChainWithUnavailable() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	err = suite.logoutAppchain(pk)
	suite.Require().Nil(err)

	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(pubAddress.Bytes())
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
	suite.Require().Contains(string(res.Ret), "under status: Unavailable")

	res, err = suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchain_mgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceUnavailable, appchain.Status)
}

//tc:激活信息缺失或错误
func (suite *Snake) Test0610_ActivateAppchainLoseFields() {
	_, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	err = suite.freezeAppchain(ChainID)
	suite.Require().Nil(err)

	res, err := suite.client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "ActivateAppchain", nil)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	suite.Require().Contains(string(res.Ret), "too few input arguments")
}

//tc:应用链未注册，激活应用链
func (suite *Snake) Test0611_ActivateAppchainWithNoRegister() {
	err := suite.activateAppchain("did:bitxhub:appchain11111111111111111111111111111111111:.")
	suite.Require().NotNil(err)
}

//tc:应用链处于注册中状态，激活应用链
func (suite *Snake) Test0612_ActivateAppchainWithRegisting() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	client, err := rpcx.New(
		rpcx.WithNodesInfo(&rpcx.NodeInfo{Addr: cfg.addrs[0]}),
		rpcx.WithLogger(cfg.logger),
		rpcx.WithPrivateKey(pk),
	)
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(pubAddress.Bytes())
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

	err = suite.activateAppchain(result.ChainID)
	fmt.Println(err)
	suite.Require().NotNil(err)
}

//tc:应用链状态已注册，激活应用链
func (suite Snake) Test0613_ActivateAppchainWithAvailable() {
	_, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	err = suite.activateAppchain(ChainID)
	fmt.Println(err)
	suite.Require().NotNil(err)
}

//tc:应用链处于更新中状态，激活应用链
func (suite *Snake) Test0614_ActivateAppchainWithUpdating() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)

	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(pubAddress.Bytes())
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":."),    //ownerDID
		rpcx.String("/ipfs/QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"), //docAddr
		rpcx.String("QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),       //docHash
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain"),         //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "UpdateAppchain", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))

	res, err = suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchain_mgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceUpdating, appchain.Status)

	err = suite.activateAppchain(ChainID)
	suite.Require().NotNil(err)
}

//tc:应用链处于冻结中状态，激活应用链
func (suite *Snake) Test0615_ActivateAppchainWithFreezing() {
	_, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	res, err := suite.client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "FreezeAppchain", nil, rpcx.String(ChainID))
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))

	res, err = suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchain_mgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceFreezing, appchain.Status)

	err = suite.activateAppchain(ChainID)
	suite.Require().NotNil(err)
}

//tc:应用链处于冻结状态，激活应用链
func (suite *Snake) Test0616_ActivateAppchain() {
	_, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	err = suite.freezeAppchain(ChainID)
	suite.Require().Nil(err)

	err = suite.activateAppchain(ChainID)
	suite.Require().Nil(err)

	res, err := suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchain_mgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceAvailable, appchain.Status)
}

//tc:应用链处于注销中状态，激活应用链
func (suite Snake) Test0617_ActivateAppchainWithLogouting() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "LogoutAppchain", nil, rpcx.String(ChainID))
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))

	err = suite.activateAppchain(ChainID)
	suite.Require().NotNil(err)

	res, err = suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchain_mgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceLogouting, appchain.Status)
}

//tc:应用链处于注销状态，激活应用链
func (suite Snake) Test0618_ActivateAppchainWithUnavailable() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	err = suite.logoutAppchain(pk)
	suite.Require().Nil(err)

	res, err := suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchain_mgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceForbidden, appchain.Status)

	err = suite.activateAppchain(ChainID)
	suite.Require().NotNil(err)
}

//tc:更新信息缺失或错误
func (suite *Snake) Test0619_UpdateAppchainLoseFields() {
	pk, _, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(pubAddress.Bytes())
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
func (suite *Snake) Test0620_UpdateAppchainWithRegisting() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	client, err := rpcx.New(
		rpcx.WithNodesInfo(&rpcx.NodeInfo{Addr: cfg.addrs[0]}),
		rpcx.WithLogger(cfg.logger),
		rpcx.WithPrivateKey(pk),
	)
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(pubAddress.Bytes())
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
func (suite *Snake) Test0621_UpdateAppchain() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(pubAddress.Bytes())
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":."),    //ownerDID
		rpcx.String("/ipfs/QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"), //docAddr
		rpcx.String("QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),       //docHash
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain11111"),    //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "UpdateAppchain", nil, args...)
	suite.Require().Nil(err)
	err = suite.VotePass(string(res.Ret))
	suite.Require().Nil(err)

	res, err = suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchain_mgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceAvailable, appchain.Status)
}

//tc:应用链处于更新中的状态，更新应用链
func (suite *Snake) Test0622_UpdateAppchainWithUpdating() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(pubAddress.Bytes())
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":."),    //ownerDID
		rpcx.String("/ipfs/QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"), //docAddr
		rpcx.String("QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),       //docHash
		rpcx.String(""),                 //validators
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
	appchain := &appchain_mgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceUpdating, appchain.Status)

	err = suite.updateAppchain(pk, args...)
	suite.Require().NotNil(err)
}

//tc:应用链处于冻结中的状态，更新应用链
func (suite Snake) Test0623_UpdateAppchainWithFreezing() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	res, err := suite.client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "FreezeAppchain", nil, rpcx.String(ChainID))
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))

	res, err = suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchain_mgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceFreezing, appchain.Status)

	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(pubAddress.Bytes())
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
func (suite Snake) Test0624_UpdateAppchainWithFrozen() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	err = suite.freezeAppchain(ChainID)

	res, err := suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchain_mgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceFrozen, appchain.Status)

	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(pubAddress.Bytes())
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
func (suite Snake) Test0625_UpdateAppchainWithWithLogouting() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	_, err = suite.client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "LogoutAppchain", nil, rpcx.String(ChainID))
	suite.Require().Nil(err)

	res, err := suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchain_mgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceLogouting, appchain.Status)

	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(pubAddress.Bytes())
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
func (suite Snake) Test0626_UpdateAppchainWithUnavailable() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	err = suite.logoutAppchain(pk)
	suite.Require().Nil(err)
	res, err := suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchain_mgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceUnavailable, appchain.Status)

	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(pubAddress.Bytes())
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
func (suite *Snake) Test0627_FreezeAppchainLoseFields() {
	pk, _, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "FreezeAppchain", nil)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	suite.Require().Contains(string(res.Ret), "too few input arguments")
}

//tc:应用链处于注册中的状态，冻结应用链
func (suite *Snake) Test0628_FreezeAppchainWithRegisting() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	client, err := rpcx.New(
		rpcx.WithNodesInfo(&rpcx.NodeInfo{Addr: cfg.addrs[0]}),
		rpcx.WithLogger(cfg.logger),
		rpcx.WithPrivateKey(pk),
	)
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(pubAddress.Bytes())
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

	err = suite.freezeAppchain(result.ChainID)
	suite.Require().NotNil(err)
}

//tc:应用链状态已注册，冻结应用链
func (suite *Snake) Test0629_FreezeAppchain() {
	_, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	err = suite.freezeAppchain(ChainID)
	suite.Require().Nil(err)

	res, err := suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchain_mgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceFrozen, appchain.Status)
}

//tc:应用链处于更新中的状态，冻结应用链
func (suite Snake) Test0630_FreezeAppchainWithUpdating() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(pubAddress.Bytes())
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":."),    //ownerDID
		rpcx.String("/ipfs/QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"), //docAddr
		rpcx.String("QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),       //docHash
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain"),         //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "UpdateAppchain", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))

	res, err = suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchain_mgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceUpdating, appchain.Status)

	err = suite.freezeAppchain(ChainID)
	suite.Require().NotNil(err)
}

//tc:应用链处于冻结中的状态，冻结应用链
func (suite *Snake) Test0631_FreezeAppchainWithFreezing() {
	_, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	res, err := suite.client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "FreezeAppchain", nil, rpcx.String(ChainID))
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))

	res, err = suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchain_mgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceFreezing, appchain.Status)

	err = suite.freezeAppchain(ChainID)
	suite.Require().NotNil(err)
}

//tc:应用链处于冻结的状态，冻结应用链
func (suite *Snake) Test0632_FreezeAppchainWithFrozen() {
	_, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	err = suite.freezeAppchain(ChainID)
	suite.Require().Nil(err)

	err = suite.freezeAppchain(ChainID)
	suite.Require().NotNil(err)

	res, err := suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchain_mgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceFrozen, appchain.Status)
}

//tc:应用链处于注销中状态，冻结应用链
func (suite Snake) Test0633_FreezeAppchainWithWithLogouting() {
	_, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	res, err := suite.client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "LogoutAppchain", nil, rpcx.String(ChainID))
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))

	res, err = suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchain_mgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceLogouting, appchain.Status)

	err = suite.freezeAppchain(ChainID)
	suite.Require().NotNil(err)
}

//tc:应用链处于注销状态，冻结应用链
func (suite Snake) Test0634_FreezeAppchainWithUnavailable() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	err = suite.logoutAppchain(pk)
	suite.Require().Nil(err)
	res, err := suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchain_mgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceUnavailable, appchain.Status)

	err = suite.freezeAppchain(ChainID)
	suite.Require().NotNil(err)
}

//tc:应用链处于注册中的状态，注销应用链
func (suite *Snake) Test0635_LogoutAppchainWithRegisting() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	client, err := rpcx.New(
		rpcx.WithNodesInfo(&rpcx.NodeInfo{Addr: cfg.addrs[0]}),
		rpcx.WithLogger(cfg.logger),
		rpcx.WithPrivateKey(pk),
	)
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(pubAddress.Bytes())
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
func (suite Snake) Test0636_LogoutAppchain() {
	pk, _, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	err = suite.logoutAppchain(pk)
	suite.Require().Nil(err)
}

//tc:应用链处于更新中的状态，注销应用链
func (suite *Snake) Test0637_LogoutAppchainWithUpdating() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(pubAddress.Bytes())
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":."),    //ownerDID
		rpcx.String("/ipfs/QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"), //docAddr
		rpcx.String("QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),       //docHash
		rpcx.String(""),                 //validators
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
	appchain := &appchain_mgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceUpdating, appchain.Status)

	err = suite.logoutAppchain(pk)
	suite.Require().NotNil(err)
}

//tc:应用链处于冻结中的状态，注销应用链
func (suite *Snake) Test0638_LogoutAppchainWithFreezing() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	res, err := suite.client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "FreezeAppchain", nil, rpcx.String(ChainID))
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))

	res, err = suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchain_mgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceFreezing, appchain.Status)

	err = suite.logoutAppchain(pk)
	suite.Require().NotNil(err)
}

//tc:应用链处于冻结的状态，注销应用链
func (suite *Snake) Test0639_LogoutAppchainWithFrozen() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	err = suite.freezeAppchain(ChainID)

	res, err := suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchain_mgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceFrozen, appchain.Status)

	err = suite.logoutAppchain(pk)
	suite.Require().NotNil(err)
}

//tc:应用链处于注销中状态，注销应用链
func (suite *Snake) Test0640_LogoutAppchainWithLogouting() {
	pk, ChainID, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	res, err := suite.client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "LogoutAppchain", nil, rpcx.String(ChainID))
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))

	res, err = suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &appchain_mgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceLogouting, appchain.Status)

	err = suite.logoutAppchain(pk)
	suite.Require().NotNil(err)
}

//tc:应用链处于注销状态，注销应用链
func (suite Snake) Test0641_LogoutAppchainWithUnavailable() {
	pk, _, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	err = suite.logoutAppchain(pk)
	suite.Require().Nil(err)

	err = suite.logoutAppchain(pk)
	suite.Require().NotNil(err)
}

//tc:根据指定ID查询应用链信息
func (suite *Snake) Test0643_GetAppchainByID() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)

	client, err := rpcx.New(
		rpcx.WithNodesInfo(&rpcx.NodeInfo{Addr: cfg.addrs[0]}),
		rpcx.WithLogger(cfg.logger),
		rpcx.WithPrivateKey(pk),
	)
	suite.Require().Nil(err)

	var pubKeyStr = base64.StdEncoding.EncodeToString(pubAddress.Bytes())
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
	suite.Require().NotNil(result.ChainID)
	err = suite.VotePass(result.ProposalID)
	suite.Require().Nil(err)

	res, err = suite.GetChainStatusById(result.ChainID)
	suite.Require().Nil(err)
	appchain := &appchain_mgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceAvailable, appchain.Status)

	args = []*pb.Arg{
		rpcx.String(result.ChainID),
	}
	res, err = client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "GetAppchain", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(res.Status, pb.Receipt_SUCCESS)
	suite.Require().NotNil(res.Ret)
}

//tc:根据错误的ID查询应用链信息
func (suite *Snake) Test0644_GetAppchainByErrorID() {
	args := []*pb.Arg{
		rpcx.String(suite.from.String() + "123"),
	}
	res, err := suite.client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "GetAppchain", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	suite.Require().Equal("call error: this appchain does not exist", string(res.Ret))
}

func (suite *Snake) freezeAppchain(ChainID string) error {
	res, err := suite.client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "FreezeAppchain", nil, rpcx.String(ChainID))
	if err != nil {
		return err
	}
	if res.Status == pb.Receipt_FAILED {
		return errors.New(string(res.Ret))
	}
	err = suite.VotePass(string(res.Ret))
	if err != nil {
		return err
	}
	return nil
}

func (suite Snake) updateAppchain(pk crypto.PrivateKey, args ...*pb.Arg) error {
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "UpdateAppchain", nil, args...)
	if err != nil {
		return err
	}
	if res.Status == pb.Receipt_FAILED {
		return errors.New(string(res.Ret))
	}
	err = suite.VotePass(string(res.Ret))
	if err != nil {
		return err
	}
	return nil
}

func (suite *Snake) activateAppchain(ChainID string) error {
	res, err := suite.client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "ActivateAppchain", nil, rpcx.String(ChainID))
	if err != nil {
		return err
	}
	if res.Status == pb.Receipt_FAILED {
		return errors.New(string(res.Ret))
	}
	err = suite.VotePass(string(res.Ret))
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
	did := "did:bitxhub:appchain" + pubAddress.String() + ":."
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "LogoutAppchain", nil, rpcx.String(did))
	if err != nil {
		return err
	}
	if res.Status == pb.Receipt_FAILED {
		return errors.New(string(res.Ret))
	}
	err = suite.VotePass(string(res.Ret))
	if err != nil {
		return err
	}
	return nil
}
