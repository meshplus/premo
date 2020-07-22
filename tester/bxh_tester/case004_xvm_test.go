package bxh_tester

import (
	"io/ioutil"
	"math/rand"
	"time"

	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
)

func (suite *Snake) TestDeployContractIsNull() {
	bytes := make([]byte, 0)
	_, err := suite.client.DeployContract(bytes)
	suite.NotNil(err)
}

func (suite *Snake) TestDeployContractWithToAddress() {
	contract, err := ioutil.ReadFile("testdata/example.wasm")
	suite.Nil(err)

	td := &pb.TransactionData{
		Type:    pb.TransactionData_INVOKE,
		VmType:  pb.TransactionData_XVM,
		Payload: contract,
	}

	tx := &pb.Transaction{
		From:      suite.from,
		To:        suite.to,
		Data:      td,
		Timestamp: time.Now().UnixNano(),
		Nonce:     rand.Int63(),
	}

	err = tx.Sign(suite.pk)
	suite.Nil(err)
	receipt, err := suite.client.SendTransactionWithReceipt(tx)
	suite.Nil(err)
	suite.True(receipt.Status == pb.Receipt_FAILED)
}

func (suite *Snake) TestDeployContract() {
	contract, err := ioutil.ReadFile("testdata/example.wasm")
	suite.Nil(err)

	address, err := suite.client.DeployContract(contract)
	suite.Nil(err)
	suite.NotNil(address)
}

func (suite *Snake) TestInvokeContract() {
	contract, err := ioutil.ReadFile("testdata/example.wasm")
	suite.Nil(err)

	address, err := suite.client.DeployContract(contract)
	suite.Nil(err)
	suite.NotNil(address)

	result, err := suite.client.InvokeXVMContract(address, "a", rpcx.Int32(1), rpcx.Int32(2))
	suite.Nil(err)
	suite.True(result.Status == pb.Receipt_SUCCESS)
	suite.True("336" == string(result.Ret))
}

func (suite *Snake) TestInvokeContractNotExistMethod() {
	contract, err := ioutil.ReadFile("testdata/example.wasm")
	suite.Nil(err)

	address, err := suite.client.DeployContract(contract)
	suite.Nil(err)
	suite.NotNil(address)

	result, err := suite.client.InvokeXVMContract(address, "bbb", rpcx.Int32(1), rpcx.Int32(2))
	suite.Nil(err)
	suite.True(result.Status == pb.Receipt_FAILED)
}
