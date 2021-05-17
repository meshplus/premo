package bxh_tester

import (
	"math/rand"

	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
)

type Model3 struct {
	*Snake
}

//tc:调用store合约，set 10M数据，交易回执显示失败
func (suite *Model3) Test0301_Set10MData() {
	rand10MBytes := make([]byte, 1<<23+1<<21)
	_, err := rand.Read(rand10MBytes)
	suite.Require().Nil(err)

	_, err = suite.client.InvokeBVMContract(constant.StoreContractAddr.Address(), "Set", nil, pb.String("test-10m"), pb.String(string(rand10MBytes)))
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "received message larger than max")
}

//tc:调用store合约，get的key为空，交易回执显示失败
func (suite *Model3) Test0302_GetEmptyKey() {
	receipt, err := suite.client.InvokeBVMContract(constant.StoreContractAddr.Address(), "Get", nil, pb.String("key_for_not_exist"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, receipt.Status)
	suite.Require().Contains(string(receipt.Ret), "there is not exist key")
}

//tc:调用store合约，set的key为空，交易回执显示失败
func (suite *Model3) Test0303_SetEmptyKey() {
	receipt, err := suite.client.InvokeBVMContract(constant.StoreContractAddr.Address(), "Set", nil, pb.String(""), pb.String("value_for_empty"))
	suite.Require().Nil(err)
	suite.Require().True(receipt.IsSuccess())
}

//tc:调用store合约，set的value为空，交易回执显示失败
func (suite *Model3) Test0304_SetEmptyValue() {
	receipt, err := suite.client.InvokeBVMContract(constant.StoreContractAddr.Address(), "Set", nil, pb.String("key_for_empty"), pb.String(""))
	suite.Require().Nil(err)
	suite.Require().True(receipt.IsSuccess())
}

//tc:调用store合约，set （a，b），交易回执显示成功
//tc:调用store合约，get（a），交易回执状态显示成功并且回执数据为b
func (suite *Model3) Test0305_SetAndGetNormal() {
	normalKey := "key_for_normal"
	normalValue := "value_for_normal"
	receipt1, err := suite.client.InvokeBVMContract(constant.StoreContractAddr.Address(), "Set", nil, pb.String(normalKey), pb.String(normalValue))
	suite.Require().Nil(err)
	suite.Require().True(receipt1.IsSuccess())

	receipt2, err := suite.client.InvokeBVMContract(constant.StoreContractAddr.Address(), "Get", nil, pb.String("key_for_normal"))
	suite.Require().Nil(err)

	suite.Require().True(receipt2.IsSuccess())
	suite.Require().Equal(normalValue, string(receipt2.Ret))
}
