package bxh_tester

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/tidwall/gjson"
)

//tc:正常注册应用链，返回回执状态成功
func (suite *Snake) Test0601_RegisterAppchain() {
	pubAddress, err := suite.pk.PublicKey().Address()
	suite.Require().Nil(err)

	var pubKeyStr = hex.EncodeToString(pubAddress.Bytes())
	args := []*pb.Arg{
		rpcx.String(""),                 //validators
		rpcx.Int32(0),                   //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain1"),        //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	res, err := suite.client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "Register", nil, args...)
	suite.Require().Nil(err)
	appChain := &rpcx.Appchain{}
	err = json.Unmarshal(res.Ret, appChain)
	suite.Require().Nil(err)
	suite.Require().NotNil(appChain.ID)
}

//tc:必填字段测试，返回回执状态失败
func (suite *Snake) Test0602_RegisterAppchainLoseFields() {
	args := []*pb.Arg{
		rpcx.String(""),    //validators
		rpcx.Int32(0),      //consensus_type
		rpcx.String("1.8"), //version
	}
	res, err := suite.client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "Register", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(res.Status, pb.Receipt_FAILED)
	suite.Require().Contains(string(res.Ret), "too few input arguments")
}

//tc:重复注册应用链，返回回执状态成功
func (suite *Snake) Test0603_RegisterReplicaAppchain() {
	pubBytes, err := suite.pk.PublicKey().Bytes()
	suite.Require().Nil(err)

	var pubKeyStr = hex.EncodeToString(pubBytes)
	args := []*pb.Arg{
		rpcx.String(""),                 //validators
		rpcx.Int32(0),                   //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain1"),        //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	res, err := suite.client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "Register", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	appchainID := gjson.Get(string(res.Ret), "appchainID").String()

	res1, err := suite.client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "Register", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res1.Status)
	appchainID1 := gjson.Get(string(res1.Ret), "appchainID").String()
	suite.Require().Equal(appchainID, appchainID1)
}

//tc:正常更新应用链
func (suite *Snake) Test0604_UpdateAppchain() {
	pubBytes, err := suite.pk.PublicKey().Bytes()
	suite.Require().Nil(err)

	var pubKeyStr = hex.EncodeToString(pubBytes)
	args := []*pb.Arg{
		rpcx.String(""),                 //validators
		rpcx.Int32(0),                   //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain1"),        //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	_, err = suite.client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "Register", nil, args...)
	suite.Require().Nil(err)

	args1 := []*pb.Arg{
		rpcx.String(suite.from.String()),
		rpcx.Int32(1),               //audit approve
		rpcx.String("Audit passed"), //desc
	}
	res1, err := suite.client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "Audit", nil, args1...)
	suite.Require().Nil(err)
	suite.Require().Equal(res1.Status, pb.Receipt_SUCCESS)

	args[2] = rpcx.String("hyperchain11111")
	res2, err := suite.client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "UpdateAppchain", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(res2.Status, pb.Receipt_SUCCESS)
}

//tc:必填字段测试，返回回执状态失败
func (suite *Snake) Test0605_UpdateAppchainLoseFields() {
	pubBytes, err := suite.pk.PublicKey().Bytes()
	suite.Require().Nil(err)

	pubKeyStr := hex.EncodeToString(pubBytes)
	args := []*pb.Arg{
		rpcx.String(""),                 //validators
		rpcx.Int32(0),                   //consensus_type
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	res, err := suite.client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "UpdateAppchain", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(res.Status, pb.Receipt_FAILED)
	suite.Require().Contains(string(res.Ret), "too few input arguments")
}

//tc:应用链审核状态改变，返回回执状态成功
func (suite *Snake) Test0606_AuditAppchain() {
	pubBytes, err := suite.pk.PublicKey().Bytes()
	suite.Require().Nil(err)

	var pubKeyStr = hex.EncodeToString(pubBytes)
	args := []*pb.Arg{
		rpcx.String(""),                 //validators
		rpcx.Int32(0),                   //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain1"),        //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	_, err = suite.client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "Register", nil, args...)
	suite.Require().Nil(err)

	args1 := []*pb.Arg{
		rpcx.String(suite.from.String()),
		rpcx.Int32(1),               //audit approve
		rpcx.String("Audit passed"), //desc
	}
	res, err := suite.client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "Audit", nil, args1...)
	suite.Require().Nil(err)
	suite.Require().Equal(res.Status, pb.Receipt_SUCCESS)
	suite.Require().Contains(string(res.Ret), "audit")
	suite.Require().Contains(string(res.Ret), "successfully")
}

//tc:多次审核，返回回执状态成功
func (suite *Snake) Test0607_RepeatAuditAppchain() {
	pubBytes, err := suite.pk.PublicKey().Bytes()
	suite.Require().Nil(err)

	var pubKeyStr = hex.EncodeToString(pubBytes)
	args := []*pb.Arg{
		rpcx.String(""),                 //validators
		rpcx.Int32(0),                   //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("AppChain1"),        //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	_, err = suite.client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "Register", nil, args...)
	suite.Require().Nil(err)

	args1 := []*pb.Arg{
		rpcx.String(suite.from.String()),
		rpcx.Int32(0),                 //audit approve
		rpcx.String("Audit rejected"), //desc
	}
	res1, err := suite.client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "Audit", nil, args1...)
	suite.Require().Nil(err)
	suite.Require().Equal(res1.Status, pb.Receipt_SUCCESS)

	args2 := []*pb.Arg{
		rpcx.String(suite.from.String()),
		rpcx.Int32(1),               //audit approve
		rpcx.String("Audit passed"), //desc
	}
	res2, err := suite.client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "Audit", nil, args2...)
	suite.Require().Nil(err)
	suite.Require().Equal(res2.Status, pb.Receipt_SUCCESS)
}

//tc:正确获取审核记录
func (suite *Snake) Test0608_FetchAuditRecord() {
	args := []*pb.Arg{
		rpcx.String(suite.from.String()),
	}
	res, err := suite.client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "FetchAuditRecords", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(res.Status, pb.Receipt_SUCCESS)
	suite.Require().NotNil(res.Ret)
}

//tc:查询调用方的应用链信息
func (suite *Snake) Test0609_GetAppchain() {
	var args []*pb.Arg
	res, err := suite.client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "Appchain", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(res.Status, pb.Receipt_SUCCESS)
	suite.Require().NotNil(res.Ret)
	fmt.Println(string(res.Ret))
}

//tc:根据指定ID查询应用链信息
func (suite *Snake) Test0610_GetAppchainByID() {
	args := []*pb.Arg{
		rpcx.String(suite.from.String()),
	}
	res, err := suite.client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "GetAppchain", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(res.Status, pb.Receipt_SUCCESS)
	suite.Require().NotNil(res.Ret)
}

//tc:根据错误的ID查询应用链信息
func (suite *Snake) Test0611_GetAppchainByErrorID() {
	args := []*pb.Arg{
		rpcx.String(suite.from.String() + "123"),
	}
	res, err := suite.client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "GetAppchain", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	suite.Require().Equal("call error: this appchain does not exist", string(res.Ret))
}
