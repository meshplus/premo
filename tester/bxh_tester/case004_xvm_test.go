package bxh_tester

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"

	"github.com/meshplus/bitxhub-kit/hexutil"
	"github.com/meshplus/bitxhub-kit/types"
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
	deployExampleContract(suite)
}

func (suite *Snake) TestInvokeContract() {
	address := deployExampleContract(suite)

	result, err := suite.client.InvokeXVMContract(address, "a", rpcx.Int32(1), rpcx.Int32(2))
	suite.Nil(err)
	suite.True(result.Status == pb.Receipt_SUCCESS)
	suite.True("336" == string(result.Ret))
}

func (suite *Snake) TestInvokeContractNotExistMethod() {
	address := deployExampleContract(suite)

	result, err := suite.client.InvokeXVMContract(address, "bbb", rpcx.Int32(1), rpcx.Int32(2))
	suite.Nil(err)
	suite.True(result.Status == pb.Receipt_FAILED)
}

func (suite *Snake) TestInvokeRandomAddressContract() {
	bs := hexutil.Encode([]byte("random contract address"))
	fmt.Println(bs)
	fakeAddr := types.String2Address(bs)

	result, err := suite.client.InvokeXVMContract(fakeAddr, "bbb", rpcx.Int32(1))
	suite.Nil(err)
	suite.True(result.Status == pb.Receipt_FAILED)
}

func (suite *Snake) TestInvokeContractEmptyMethod() {
	address := deployExampleContract(suite)

	result, err := suite.client.InvokeXVMContract(address, "")
	suite.Nil(err)
	suite.True(result.Status == pb.Receipt_FAILED)
}

func (suite *Snake) TestDeploy10MContract() {
	// todo: wait for bitxhub to limit contract size
}

func (suite *Snake) TestDeployContractWrongArg() {
	address := deployExampleContract(suite)

	result, err := suite.client.InvokeXVMContract(address, "a", rpcx.String("1"), rpcx.Int32(2))
	suite.Nil(err)
	suite.True(result.Status == pb.Receipt_FAILED)

	// incorrect function params
	result, err = suite.client.InvokeXVMContract(address, "a", rpcx.Int32(1), rpcx.String("2"))
	suite.Nil(err)
	suite.True(result.Status == pb.Receipt_FAILED)

	result, err = suite.client.InvokeXVMContract(address, "a", rpcx.String("1"), rpcx.String("2"))
	suite.Nil(err)
	suite.True(result.Status == pb.Receipt_FAILED)
}

func (suite *Snake) TestDeployContractWrongNumberArg() {
	address := deployExampleContract(suite)

	result, err := suite.client.InvokeXVMContract(address, "a", rpcx.Int32(1), rpcx.Int32(2), rpcx.Int32(3))
	suite.Nil(err)
	suite.True(result.Status == pb.Receipt_FAILED)
}

func deployExampleContract(suite *Snake) types.Address {
	contract, err := ioutil.ReadFile("testdata/example.wasm")
	suite.Nil(err)

	address, err := suite.client.DeployContract(contract)
	suite.Nil(err)
	suite.NotNil(address)
	return address
}
