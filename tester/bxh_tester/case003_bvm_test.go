package bxh_tester

import (
	"math/rand"

	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
)

func (suite *Snake) TestSet10MData() {
	rand10MBytes := make([]byte, 1<<23+1<<21)
	_, err := rand.Read(rand10MBytes)
	suite.Nil(err)

	_, err = suite.client.InvokeBVMContract(rpcx.StoreContractAddr, "Set", pb.String("test-10m"), pb.String(string(rand10MBytes)))
	suite.Require().NotNil(err)
}

func (suite *Snake) TestGetEmptyKey() {
	receipt, err := suite.client.InvokeBVMContract(rpcx.StoreContractAddr, "Get", pb.String(""))
	suite.Nil(err)

	suite.Equal(pb.Receipt_FAILED, receipt.Status)
}

func (suite *Snake) TestSetEmptyKey() {
	receipt, err := suite.client.InvokeBVMContract(rpcx.StoreContractAddr, "Set", pb.String(""), pb.String("value_for_empty"))
	suite.Nil(err)

	suite.Equal(pb.Receipt_FAILED, receipt.Status)
}

func (suite *Snake) TestSetEmptyValue() {
	receipt, err := suite.client.InvokeBVMContract(rpcx.StoreContractAddr, "Set", pb.String("key_for_empty"), pb.String(""))
	suite.Nil(err)

	suite.Equal(pb.Receipt_FAILED, receipt.Status)
}

func (suite *Snake) TestSetAndGetNormal() {
	normalKey := "key_for_normal"
	normalValue := "value_for_normal"
	receipt, err := suite.client.InvokeBVMContract(rpcx.StoreContractAddr, "Set", pb.String(normalKey), pb.String(normalValue))
	suite.Nil(err)

	suite.Equal(pb.Receipt_SUCCESS, receipt.Status)

	receipt, err = suite.client.InvokeBVMContract(rpcx.StoreContractAddr, "Get", pb.String("key_for_normal"))
	suite.Nil(err)

	suite.Equal(pb.Receipt_SUCCESS, receipt.Status)
	suite.Equal(normalValue, string(receipt.Ret))
}

func (suite *Snake) TestGetNotExistingKey() {
	receipt, err := suite.client.InvokeBVMContract(rpcx.StoreContractAddr, "Get", pb.String("key_for_not_exist"))
	suite.Nil(err)

	suite.Equal(pb.Receipt_FAILED, receipt.Status)
}
