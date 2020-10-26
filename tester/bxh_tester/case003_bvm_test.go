package bxh_tester

import (
	"github.com/meshplus/bitxhub-model/constant"
	"math/rand"
	"strconv"

	"github.com/meshplus/bitxhub-model/pb"
)

func (suite *Snake) TestSet10MData() {
	rand10MBytes := make([]byte, 1<<23+1<<21)
	_, err := rand.Read(rand10MBytes)
	suite.Require().Nil(err)
	_, err = suite.client.InvokeBVMContract(constant.StoreContractAddr.Address(), "Set", nil, pb.String("test-10m"), pb.String(string(rand10MBytes)))
	suite.Require().NotNil(err)
}

func (suite *Snake) TestGetEmptyKey() {
	receipt1, err := suite.client.InvokeBVMContract(constant.StoreContractAddr.Address(), "Get", nil, pb.String(strconv.FormatUint(uint64(rand.Int()), 10)))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, receipt1.Status)
}

func (suite *Snake) TestSetEmptyKey() {
	receipt, err := suite.client.InvokeBVMContract(constant.StoreContractAddr.Address(), "Set", nil, pb.String(""), pb.String("value_for_empty"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, receipt.Status)
}

func (suite *Snake) TestSetEmptyValue() {
	receipt, err := suite.client.InvokeBVMContract(constant.StoreContractAddr.Address(), "Set", nil, pb.String("key_for_empty"), pb.String(""))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, receipt.Status)
}

func (suite *Snake) TestSetAndGetNormal() {
	normalKey := "key_for_normal"
	normalValue := "value_for_normal"
	receipt, err := suite.client.InvokeBVMContract(constant.StoreContractAddr.Address(), "Set", nil, pb.String(normalKey), pb.String(normalValue))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, receipt.Status)

	receipt, err = suite.client.InvokeBVMContract(constant.StoreContractAddr.Address(), "Get", nil, pb.String("key_for_normal"))
	suite.Require().Nil(err)

	suite.Require().Equal(pb.Receipt_SUCCESS, receipt.Status)
	suite.Require().Equal(normalValue, string(receipt.Ret))
}

func (suite *Snake) TestGetNotExistingKey() {
	receipt, err := suite.client.InvokeBVMContract(constant.StoreContractAddr.Address(), "Get", nil, pb.String("key_for_not_exist"))
	suite.Require().Nil(err)

	suite.Require().Equal(pb.Receipt_FAILED, receipt.Status)
}
