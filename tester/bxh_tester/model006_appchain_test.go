package bxh_tester

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	appchain_mgr "github.com/meshplus/bitxhub-core/appchain-mgr"
	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/meshplus/premo/internal/repo"
	"github.com/pkg/errors"
)

type RegisterResult struct {
	ChainID    string `json:"chain_id"`
	ProposalID string `json:"proposal_id"`
}

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
	_, ChainID, err := suite.registerAppchain()
	suite.Require().Nil(err)
	res, err := suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &rpcx.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(appchain_mgr.AppchainAvailable, appchain.Status)
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
	var pubKeyStr = hex.EncodeToString(pubAddress.Bytes())
	args := []*pb.Arg{
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain1"),        //name
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

	res, err = suite.GetChainStatusById(result.ChainID)
	suite.Require().Nil(err)
	appchain := &rpcx.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(appchain_mgr.AppchainRegisting, appchain.Status)
}

//tc:应用链状态已注册，注册应用链
func (suite *Snake) Test0604_RegisterAppchainRepeat() {
	pk, ChainID, err := suite.registerAppchain()
	suite.Require().Nil(err)

	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	var pubKeyStr = hex.EncodeToString(pubAddress.Bytes())
	args := []*pb.Arg{
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain1"),        //name
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
	appchain := &rpcx.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(appchain_mgr.AppchainAvailable, appchain.Status)
}

//tc:应用链处于更新中状态，注册应用链
func (suite *Snake) Test0605_RegisterAppchainWithUpdating() {
	pk, ChainID, err := suite.registerAppchain()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	var pubKeyStr = hex.EncodeToString(pubAddress.Bytes())
	args := []*pb.Arg{
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain111"),    //chain_type
		rpcx.String("AppChain1"),        //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	_, err = client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "UpdateAppchain", nil, args...)
	suite.Require().Nil(err)

	res, err := suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &rpcx.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(appchain_mgr.AppchainUpdating, appchain.Status)

	res, err = client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "Register", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Contains(string(res.Ret), "appchain has registered")

	res, err = suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(appchain_mgr.AppchainUpdating, appchain.Status)
}

//tc:应用链处于冻结中状态，注册应用链
func (suite *Snake) Test0606_RegisterAppchainWithFreezing() {
	pk, ChainID, err := suite.registerAppchain()
	suite.Require().Nil(err)

	client := suite.NewClient(pk)
	_, err = client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "FreezeAppchain", nil, rpcx.String(ChainID))
	suite.Require().Nil(err)

	res, err := suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &rpcx.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(appchain_mgr.AppchainFreezing, appchain.Status)

	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	var pubKeyStr = hex.EncodeToString(pubAddress.Bytes())
	args := []*pb.Arg{
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain1"),        //name
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
	suite.Require().Equal(appchain_mgr.AppchainFreezing, appchain.Status)
}

//tc:应用链处于冻结状态，注册应用链
func (suite *Snake) Test0607_RegisterAppchainWithFrozen() {
	pk, ChainID, err := suite.registerAppchain()
	suite.Require().Nil(err)

	err = suite.freezeAppchain(ChainID)
	suite.Require().Nil(err)

	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	var pubKeyStr = hex.EncodeToString(pubAddress.Bytes())
	args := []*pb.Arg{
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain1"),        //name
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
	appchain := &rpcx.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(appchain_mgr.AppchainFrozen, appchain.Status)
}

//tc:应用链处于注销中状态，注册应用链
func (suite *Snake) Test0608_RegisterAppchainWithLogouting() {
	pk, ChainID, err := suite.registerAppchain()
	suite.Require().Nil(err)

	client := suite.NewClient(pk)
	_, err = client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "LogoutAppchain", nil)
	suite.Require().Nil(err)

	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	var pubKeyStr = hex.EncodeToString(pubAddress.Bytes())
	args := []*pb.Arg{
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain1"),        //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "Register", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Contains(string(res.Ret), "appchain has registered")

	res, err = suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &rpcx.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(appchain_mgr.AppchainLogouting, appchain.Status)
}

//tc:应用链处于注销状态，注册应用链
func (suite *Snake) Test0609_RegisterAppChainWithUnavailable() {
	pk, ChainID, err := suite.registerAppchain()
	suite.Require().Nil(err)

	err = suite.logoutAppchain(pk)
	suite.Require().Nil(err)

	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	var pubKeyStr = hex.EncodeToString(pubAddress.Bytes())
	args := []*pb.Arg{
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain1"),        //name
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
	appchain := &rpcx.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(appchain_mgr.AppchainUnavailable, appchain.Status)
}

//tc:激活信息缺失或错误
func (suite *Snake) Test0610_ActivateAppchainLoseFields() {
	_, ChainID, err := suite.registerAppchain()
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
	err := suite.activateAppchain("0x111111111111111111111111111111111111111")
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
	var pubKeyStr = hex.EncodeToString(pubAddress.Bytes())
	args := []*pb.Arg{
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain1"),        //name
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
	_, ChainID, err := suite.registerAppchain()
	suite.Require().Nil(err)

	err = suite.activateAppchain(ChainID)
	suite.Require().NotNil(err)
}

//tc:应用链处于更新中状态，激活应用链
func (suite *Snake) Test0614_ActivateAppchainWithUpdating() {
	pk, ChainID, err := suite.registerAppchain()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	var pubKeyStr = hex.EncodeToString(pubAddress.Bytes())
	args := []*pb.Arg{
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain111"),    //chain_type
		rpcx.String("AppChain1"),        //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	_, err = client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "UpdateAppchain", nil, args...)
	suite.Require().Nil(err)

	res, err := suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &rpcx.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(appchain_mgr.AppchainUpdating, appchain.Status)

	err = suite.activateAppchain(ChainID)
	suite.Require().NotNil(err)
}

//tc:应用链处于冻结中状态，激活应用链
func (suite *Snake) Test0615_ActivateAppchainWithFreezing() {
	pk, ChainID, err := suite.registerAppchain()
	suite.Require().Nil(err)

	client := suite.NewClient(pk)
	_, err = client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "FreezeAppchain", nil, rpcx.String(ChainID))
	suite.Require().Nil(err)

	res, err := suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &rpcx.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(appchain_mgr.AppchainFreezing, appchain.Status)

	err = suite.activateAppchain(ChainID)
	suite.Require().NotNil(err)
}

//tc:应用链处于冻结状态，注册应用链
func (suite *Snake) Test0616_ActivateAppchain() {
	_, ChainID, err := suite.registerAppchain()
	suite.Require().Nil(err)

	err = suite.freezeAppchain(ChainID)
	suite.Require().Nil(err)

	err = suite.activateAppchain(ChainID)
	suite.Require().Nil(err)

	res, err := suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &rpcx.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(appchain_mgr.AppchainAvailable, appchain.Status)
}

//tc:应用链处于注销中状态，激活应用链
func (suite Snake) Test0617_ActivateAppchainWithLogouting() {
	pk, ChainID, err := suite.registerAppchain()
	suite.Require().Nil(err)

	client := suite.NewClient(pk)
	_, err = client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "LogoutAppchain", nil)
	suite.Require().Nil(err)

	err = suite.activateAppchain(ChainID)
	suite.Require().NotNil(err)

	res, err := suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &rpcx.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(appchain_mgr.AppchainLogouting, appchain.Status)
}

//tc:应用链处于注销状态，激活应用链
func (suite Snake) Test0618_ActivateAppchainWithUnavailable() {
	pk, ChainID, err := suite.registerAppchain()
	suite.Require().Nil(err)

	err = suite.logoutAppchain(pk)
	suite.Require().Nil(err)

	res, err := suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &rpcx.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(appchain_mgr.AppchainUnavailable, appchain.Status)

	err = suite.activateAppchain(ChainID)
	suite.Require().NotNil(err)
}

//tc:更新信息缺失或错误
func (suite *Snake) Test0619_UpdateAppchainLoseFields() {
	pk, _, err := suite.registerAppchain()
	suite.Require().Nil(err)

	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	var pubKeyStr = hex.EncodeToString(pubAddress.Bytes())
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
	var pubKeyStr = hex.EncodeToString(pubAddress.Bytes())
	args := []*pb.Arg{
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain1"),        //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "Register", nil, args...)
	suite.Require().Nil(err)
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	suite.Require().Nil(err)

	args[3] = rpcx.String("AppChain111111111")
	err = suite.updateAppchain(pk, args...)
	suite.Require().NotNil(err)
}

//tc:应用链状态已注册，更新应用链
func (suite *Snake) Test0621_UpdateAppchain() {
	pk, ChainID, err := suite.registerAppchain()
	suite.Require().Nil(err)

	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	var pubKeyStr = hex.EncodeToString(pubAddress.Bytes())
	args := []*pb.Arg{
		rpcx.String(""),                    //validators
		rpcx.String("raft"),                //consensus_type
		rpcx.String("hyperchain"),          //chain_type
		rpcx.String("AppChain11111111111"), //name
		rpcx.String("Appchain for tax"),    //desc
		rpcx.String("1.8"),                 //version
		rpcx.String(pubKeyStr),             //public key
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "UpdateAppchain", nil, args...)
	suite.Require().Nil(err)
	err = suite.VotePass(string(res.Ret))
	suite.Require().Nil(err)

	res, err = suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &rpcx.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(appchain_mgr.AppchainAvailable, appchain.Status)
}

//tc:应用链处于更新中的状态，更新应用链
func (suite *Snake) Test0622_UpdateAppchainWithUpdating() {
	pk, ChainID, err := suite.registerAppchain()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	var pubKeyStr = hex.EncodeToString(pubAddress.Bytes())
	args := []*pb.Arg{
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain111"),    //chain_type
		rpcx.String("AppChain1"),        //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	_, err = client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "UpdateAppchain", nil, args...)
	suite.Require().Nil(err)

	res, err := suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &rpcx.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(appchain_mgr.AppchainUpdating, appchain.Status)

	err = suite.updateAppchain(pk, args...)
	suite.Require().NotNil(err)
}

//tc:应用链处于冻结中的状态，更新应用链
func (suite Snake) Test0623_UpdateAppchainWithFreezing() {
	pk, ChainID, err := suite.registerAppchain()
	suite.Require().Nil(err)

	client := suite.NewClient(pk)
	_, err = client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "FreezeAppchain", nil, rpcx.String(ChainID))
	suite.Require().Nil(err)

	res, err := suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &rpcx.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(appchain_mgr.AppchainFreezing, appchain.Status)

	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	var pubKeyStr = hex.EncodeToString(pubAddress.Bytes())
	args := []*pb.Arg{
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain111"),    //chain_type
		rpcx.String("AppChain1"),        //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	err = suite.updateAppchain(pk, args...)
	suite.Require().NotNil(err)
}

//tc:应用链处于冻结的状态，更新应用链
func (suite Snake) Test0624_UpdateAppchainWithFrozen() {
	pk, ChainID, err := suite.registerAppchain()
	suite.Require().Nil(err)

	err = suite.freezeAppchain(ChainID)

	res, err := suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &rpcx.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(appchain_mgr.AppchainFrozen, appchain.Status)

	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	var pubKeyStr = hex.EncodeToString(pubAddress.Bytes())
	args := []*pb.Arg{
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain111"),    //chain_type
		rpcx.String("AppChain1"),        //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	err = suite.updateAppchain(pk, args...)
	suite.Require().NotNil(err)
}

//tc:应用链处于注销中状态，更新应用链
func (suite Snake) Test0625_UpdateAppchainWithWithLogouting() {
	pk, ChainID, err := suite.registerAppchain()
	suite.Require().Nil(err)

	client := suite.NewClient(pk)
	_, err = client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "LogoutAppchain", nil)
	suite.Require().Nil(err)

	res, err := suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &rpcx.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(appchain_mgr.AppchainLogouting, appchain.Status)

	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	var pubKeyStr = hex.EncodeToString(pubAddress.Bytes())
	args := []*pb.Arg{
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain111"),    //chain_type
		rpcx.String("AppChain1"),        //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	err = suite.updateAppchain(pk, args...)
	suite.Require().NotNil(err)
}

//tc:应用链处于注销状态，更新应用链
func (suite Snake) Test0626_UpdateAppchainWithUnavailable() {
	pk, ChainID, err := suite.registerAppchain()
	suite.Require().Nil(err)

	err = suite.logoutAppchain(pk)
	suite.Require().Nil(err)
	res, err := suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &rpcx.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(appchain_mgr.AppchainUnavailable, appchain.Status)

	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	var pubKeyStr = hex.EncodeToString(pubAddress.Bytes())
	args := []*pb.Arg{
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain111"),    //chain_type
		rpcx.String("AppChain1"),        //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	err = suite.updateAppchain(pk, args...)
	suite.Require().NotNil(err)
}

//tc:冻结信息缺失或错误
func (suite *Snake) Test0627_FreezeAppchainLoseFields() {
	pk, _, err := suite.registerAppchain()
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
	var pubKeyStr = hex.EncodeToString(pubAddress.Bytes())
	args := []*pb.Arg{
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain1"),        //name
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
	_, ChainID, err := suite.registerAppchain()
	suite.Require().Nil(err)

	err = suite.freezeAppchain(ChainID)
	suite.Require().Nil(err)

	res, err := suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &rpcx.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(appchain_mgr.AppchainFrozen, appchain.Status)
}

//tc:应用链处于更新中的状态，冻结应用链
func (suite Snake) Test0630_FreezeAppchainWithUpdating() {
	pk, ChainID, err := suite.registerAppchain()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	var pubKeyStr = hex.EncodeToString(pubAddress.Bytes())
	args := []*pb.Arg{
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain111"),    //chain_type
		rpcx.String("AppChain1"),        //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	_, err = client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "UpdateAppchain", nil, args...)
	suite.Require().Nil(err)

	res, err := suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &rpcx.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(appchain_mgr.AppchainUpdating, appchain.Status)

	err = suite.freezeAppchain(ChainID)
	suite.Require().NotNil(err)
}

//tc:应用链处于冻结中的状态，冻结应用链
func (suite *Snake) Test0631_FreezeAppchainWithFreezing() {
	pk, ChainID, err := suite.registerAppchain()
	suite.Require().Nil(err)

	client := suite.NewClient(pk)
	_, err = client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "FreezeAppchain", nil, rpcx.String(ChainID))
	suite.Require().Nil(err)

	res, err := suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &rpcx.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(appchain_mgr.AppchainFreezing, appchain.Status)

	err = suite.freezeAppchain(ChainID)
	suite.Require().NotNil(err)
}

//tc:应用链处于冻结的状态，冻结应用链
func (suite *Snake) Test0632_FreezeAppchainWithFrozen() {
	_, ChainID, err := suite.registerAppchain()
	suite.Require().Nil(err)

	err = suite.freezeAppchain(ChainID)
	suite.Require().Nil(err)

	err = suite.freezeAppchain(ChainID)
	suite.Require().NotNil(err)

	res, err := suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &rpcx.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(appchain_mgr.AppchainFrozen, appchain.Status)
}

//tc:应用链处于注销中状态，冻结应用链
func (suite Snake) Test0633_FreezeAppchainWithWithLogouting() {
	pk, ChainID, err := suite.registerAppchain()
	suite.Require().Nil(err)

	client := suite.NewClient(pk)
	_, err = client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "LogoutAppchain", nil)
	suite.Require().Nil(err)

	res, err := suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &rpcx.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(appchain_mgr.AppchainLogouting, appchain.Status)

	err = suite.freezeAppchain(ChainID)
	suite.Require().NotNil(err)
}

//tc:应用链处于注销状态，冻结应用链
func (suite Snake) Test0634_FreezeAppchainWithUnavailable() {
	pk, ChainID, err := suite.registerAppchain()
	suite.Require().Nil(err)

	err = suite.logoutAppchain(pk)
	suite.Require().Nil(err)
	res, err := suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &rpcx.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(appchain_mgr.AppchainUnavailable, appchain.Status)

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
	var pubKeyStr = hex.EncodeToString(pubAddress.Bytes())
	args := []*pb.Arg{
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain1"),        //name
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
	pk, _, err := suite.registerAppchain()
	suite.Require().Nil(err)

	err = suite.logoutAppchain(pk)
	suite.Require().Nil(err)
}

//tc:应用链处于更新中的状态，注销应用链
func (suite *Snake) Test0637_LogoutAppchainWithUpdating() {
	pk, ChainID, err := suite.registerAppchain()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	pubAddress, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	var pubKeyStr = hex.EncodeToString(pubAddress.Bytes())
	args := []*pb.Arg{
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain111"),    //chain_type
		rpcx.String("AppChain1"),        //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	_, err = client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "UpdateAppchain", nil, args...)
	suite.Require().Nil(err)

	res, err := suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &rpcx.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(appchain_mgr.AppchainUpdating, appchain.Status)

	err = suite.logoutAppchain(pk)
	suite.Require().NotNil(err)
}

//tc:应用链处于冻结中的状态，注销应用链
func (suite *Snake) Test0638_LogoutAppchainWithFreezing() {
	pk, ChainID, err := suite.registerAppchain()
	suite.Require().Nil(err)

	client := suite.NewClient(pk)
	_, err = client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "FreezeAppchain", nil, rpcx.String(ChainID))
	suite.Require().Nil(err)

	res, err := suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &rpcx.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(appchain_mgr.AppchainFreezing, appchain.Status)

	err = suite.logoutAppchain(pk)
	suite.Require().NotNil(err)
}

//tc:应用链处于冻结的状态，注销应用链
func (suite *Snake) Test0639_LogoutAppchainWithFrozen() {
	pk, ChainID, err := suite.registerAppchain()
	suite.Require().Nil(err)

	err = suite.freezeAppchain(ChainID)

	res, err := suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &rpcx.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(appchain_mgr.AppchainFrozen, appchain.Status)

	err = suite.logoutAppchain(pk)
	suite.Require().NotNil(err)
}

//tc:应用链处于注销中状态，注销应用链
func (suite *Snake) Test0640_LogoutAppchainWithLogouting() {
	pk, ChainID, err := suite.registerAppchain()
	suite.Require().Nil(err)

	client := suite.NewClient(pk)
	_, err = client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "LogoutAppchain", nil)
	suite.Require().Nil(err)

	res, err := suite.GetChainStatusById(ChainID)
	suite.Require().Nil(err)
	appchain := &rpcx.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(appchain_mgr.AppchainLogouting, appchain.Status)

	err = suite.logoutAppchain(pk)
	suite.Require().NotNil(err)
}

//tc:应用链处于注销状态，注销应用链
func (suite Snake) Test0641_LogoutAppchainWithUnavailable() {
	pk, _, err := suite.registerAppchain()
	suite.Require().Nil(err)

	err = suite.logoutAppchain(pk)
	suite.Require().Nil(err)

	err = suite.logoutAppchain(pk)
	suite.Require().NotNil(err)
}

//tc:查询调用方的应用链信息
func (suite *Snake) Test0642_GetAppchain() {
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

	var pubKeyStr = hex.EncodeToString(pubAddress.Bytes())
	args := []*pb.Arg{
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain1"),        //name
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
	appchain := &rpcx.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(appchain_mgr.AppchainAvailable, appchain.Status)

	args = []*pb.Arg{}
	res, err = client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "Appchain", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(res.Status, pb.Receipt_SUCCESS)
	suite.Require().NotNil(res.Ret)
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

	var pubKeyStr = hex.EncodeToString(pubAddress.Bytes())
	args := []*pb.Arg{
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain1"),        //name
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
	appchain := &rpcx.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(appchain_mgr.AppchainAvailable, appchain.Status)

	args = []*pb.Arg{
		rpcx.String(pubAddress.String()),
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

func (suite *Snake) registerAppchain() (crypto.PrivateKey, string, error) {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	if err != nil {
		return nil, "", err
	}
	pubAddress, err := pk.PublicKey().Address()
	if err != nil {
		return nil, "", err
	}
	client, err := rpcx.New(
		rpcx.WithNodesInfo(&rpcx.NodeInfo{Addr: cfg.addrs[0]}),
		rpcx.WithLogger(cfg.logger),
		rpcx.WithPrivateKey(pk),
	)
	if err != nil {
		return nil, "", err
	}
	var pubKeyStr = hex.EncodeToString(pubAddress.Bytes())
	args := []*pb.Arg{
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain1"),        //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "Register", nil, args...)
	if err != nil {
		return nil, "", err
	}
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	if err != nil {
		return nil, "", err
	}
	err = suite.VotePass(result.ProposalID)
	if err != nil {
		return nil, "", err
	}
	return pk, result.ChainID, nil
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
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "LogoutAppchain", nil)
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

func (suite *Snake) VotePass(id string) error {
	node1, err := repo.Node1Path()
	if err != nil {
		return err
	}

	key, err := asym.RestorePrivateKey(node1, repo.KeyPassword)
	if err != nil {
		return err
	}

	_, err = suite.vote(key, pb.String(id), pb.String("approve"), pb.String("Appchain Pass"))
	if err != nil {
		return err
	}

	node2, err := repo.Node2Path()
	if err != nil {
		return err
	}

	key, err = asym.RestorePrivateKey(node2, repo.KeyPassword)
	if err != nil {
		return err
	}

	_, err = suite.vote(key, pb.String(id), pb.String("approve"), pb.String("Appchain Pass"))
	if err != nil {
		return err
	}

	node3, err := repo.Node3Path()
	if err != nil {
		return err
	}

	key, err = asym.RestorePrivateKey(node3, repo.KeyPassword)
	if err != nil {
		return err
	}

	_, err = suite.vote(key, pb.String(id), pb.String("approve"), pb.String("Appchain Pass"))
	if err != nil {
		return err
	}
	return nil
}

func (suite *Snake) vote(key crypto.PrivateKey, args ...*pb.Arg) (*pb.Receipt, error) {
	client, err := rpcx.New(
		rpcx.WithNodesInfo(&rpcx.NodeInfo{Addr: cfg.addrs[0]}),
		rpcx.WithLogger(cfg.logger),
		rpcx.WithPrivateKey(key),
	)
	address, err := key.PublicKey().Address()
	if err != nil {
		return nil, err
	}
	invokePayload := &pb.InvokePayload{
		Method: "Vote",
		Args:   args,
	}

	payload, err := invokePayload.Marshal()
	if err != nil {
		return nil, err
	}

	data := &pb.TransactionData{
		Type:    pb.TransactionData_INVOKE,
		VmType:  pb.TransactionData_BVM,
		Payload: payload,
	}
	payload, err = data.Marshal()

	tx := &pb.Transaction{
		From:      address,
		To:        constant.GovernanceContractAddr.Address(),
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}
	if err != nil {
		return nil, err
	}
	res, err := client.SendTransactionWithReceipt(tx, nil)
	if err != nil {
		return nil, err
	}
	if res.Status == pb.Receipt_FAILED {
		return nil, errors.New(string(res.Ret))
	}
	return res, nil
}

func (suite *Snake) GetChainStatusById(id string) (*pb.Receipt, error) {
	node, err := repo.Node1Path()
	key, err := asym.RestorePrivateKey(node, repo.KeyPassword)
	if err != nil {
		return nil, err
	}
	client, err := rpcx.New(
		rpcx.WithNodesInfo(&rpcx.NodeInfo{Addr: cfg.addrs[0]}),
		rpcx.WithLogger(cfg.logger),
		rpcx.WithPrivateKey(key),
	)
	address, err := key.PublicKey().Address()
	if err != nil {
		return nil, err
	}
	args := []*pb.Arg{
		rpcx.String(id),
	}
	invokePayload := &pb.InvokePayload{
		Method: "GetAppchain",
		Args:   args,
	}

	payload, err := invokePayload.Marshal()
	if err != nil {
		return nil, err
	}

	data := &pb.TransactionData{
		Type:    pb.TransactionData_INVOKE,
		VmType:  pb.TransactionData_BVM,
		Payload: payload,
	}
	payload, err = data.Marshal()

	tx := &pb.Transaction{
		From:      address,
		To:        constant.AppchainMgrContractAddr.Address(),
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}
	if err != nil {
		return nil, err
	}
	res, err := client.SendTransactionWithReceipt(tx, nil)
	if err != nil {
		return nil, err
	}
	if res.Status == pb.Receipt_FAILED {
		return nil, errors.New(string(res.Ret))
	}
	return res, nil
}
