package bxh_tester

import (
	"fmt"
	"math/rand"

	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
)

func (suite *Snake) TestSet16MData() {
	rand10MBytes := make([]byte, 1<<24)
	_, err := rand.Read(rand10MBytes)
	suite.Nil(err)

	fmt.Println(string(rand10MBytes[:10000]))
	receipt, err := suite.client.InvokeBVMContract(rpcx.StoreContractAddr, "Set", pb.String("test-16m"), pb.String(string(rand10MBytes)))
	suite.Nil(err)

	fmt.Println(string(receipt.Ret))
	suite.Equal(pb.Receipt_FAILED, receipt.Status)
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

func (suite *Snake) TestSetNormal() {
	receipt, err := suite.client.InvokeBVMContract(rpcx.StoreContractAddr, "Set", pb.String("key_for_normal"), pb.String("value_for_normal"))
	suite.Nil(err)

	suite.Equal(pb.Receipt_SUCCESS, receipt.Status)
}

func (suite *Snake) TestGetNotExistingKey() {
	receipt, err := suite.client.InvokeBVMContract(rpcx.StoreContractAddr, "Get", pb.String("key_for_not_exist"))
	suite.Nil(err)

	suite.Equal(pb.Receipt_FAILED, receipt.Status)
}
