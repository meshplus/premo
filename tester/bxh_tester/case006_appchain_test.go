package bxh_tester

import (
	"encoding/hex"
	"encoding/json"

	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/tidwall/gjson"
)

func (suite *Snake) TestRegisterAppchain() {
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

func (suite *Snake) TestRegisterAppchainLoseFields() {
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

func (suite *Snake) TestRegisterReplicaAppchain() {
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
	appchainID := gjson.Get(string(res.Ret), "appchainID").String()

	res1, err := suite.client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "Register", nil, args...)
	suite.Require().Nil(err)

	appchainID1 := gjson.Get(string(res1.Ret), "appchainID").String()
	suite.Require().Equal(appchainID, appchainID1)
}

func (suite *Snake) TestUpdateAppchain() {
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

func (suite *Snake) TestUpdateAppchainLoseFields() {
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

func (suite *Snake) TestAuditAppchain() {
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

func (suite *Snake) TestRepeatAuditAppchain() {
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

func (suite *Snake) TestFetchAuditRecord() {
	args := []*pb.Arg{
		rpcx.String(suite.from.String()),
	}
	res, err := suite.client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "FetchAuditRecords", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(res.Status, pb.Receipt_SUCCESS)
	suite.Require().NotNil(res.Ret)
}

func (suite *Snake) TestGetAppchain() {
	var args []*pb.Arg
	res, err := suite.client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "Appchain", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(res.Status, pb.Receipt_SUCCESS)
	suite.Require().NotNil(res.Ret)
}

func (suite *Snake) TestGetAppchainByID() {
	args := []*pb.Arg{
		rpcx.String(suite.from.String()),
	}
	res, err := suite.client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "GetAppchain", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(res.Status, pb.Receipt_SUCCESS)
	suite.Require().NotNil(res.Ret)
}
