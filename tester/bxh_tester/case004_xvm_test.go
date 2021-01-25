package bxh_tester

import (
	"io/ioutil"
	"time"

	"github.com/meshplus/bitxhub-kit/hexutil"
	"github.com/meshplus/bitxhub-kit/types"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
)

func (suite *Snake) TestDeployContractIsNull() {
	bytes := make([]byte, 0)
	_, err := suite.client.DeployContract(bytes, nil)

	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "can't deploy empty contract")
}

func (suite *Snake) TestDeployContractWithToAddress() {
	contract, err := ioutil.ReadFile("testdata/example.wasm")
	suite.Require().Nil(err)

	td := &pb.TransactionData{
		Type:    pb.TransactionData_INVOKE,
		VmType:  pb.TransactionData_XVM,
		Payload: contract,
	}

	payload, err := td.Marshal()
	suite.Require().Nil(err)
	tx := &pb.Transaction{
		From:      suite.from,
		To:        suite.to,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}

	receipt, err := suite.client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().True(receipt.Status == pb.Receipt_FAILED)
	suite.Require().Contains(string(receipt.Ret), "contract byte not correct")
}

func (suite *Snake) TestDeployContract() {
	deployExampleContract(suite)
}

func (suite *Snake) TestInvokeContract() {
	address := deployExampleContract(suite)

	result, err := suite.client.InvokeXVMContract(address, "a", nil, rpcx.Int32(1), rpcx.Int32(2))
	suite.Require().Nil(err)
	suite.Require().True(result.Status == pb.Receipt_SUCCESS)
	suite.Require().True("336" == string(result.Ret))
}

func (suite *Snake) TestInvokeContractNotExistMethod() {
	address := deployExampleContract(suite)

	result, err := suite.client.InvokeXVMContract(address, "bbb", nil, rpcx.Int32(1), rpcx.Int32(2))
	suite.Require().Nil(err)
	suite.Require().True(result.Status == pb.Receipt_FAILED)
	suite.Require().Contains(string(result.Ret), "wrong rule contract")
}

func (suite *Snake) TestInvokeRandomAddressContract() {
	// random addr len should be 42
	bs := hexutil.Encode([]byte("random contract addr"))
	fakeAddr := types.NewAddressByStr(bs)

	result, err := suite.client.InvokeXVMContract(fakeAddr, "bbb", nil, rpcx.Int32(1))
	suite.Require().Nil(err)
	suite.Require().True(result.Status == pb.Receipt_FAILED)
	suite.Require().Contains(string(result.Ret), "contract byte not correct")
}

func (suite *Snake) TestInvokeContractEmptyMethod() {
	address := deployExampleContract(suite)

	result, err := suite.client.InvokeXVMContract(address, "", nil)
	suite.Require().Nil(err)
	suite.Require().True(result.Status == pb.Receipt_FAILED)
	suite.Require().Contains(string(result.Ret), "lack of method name")
}

func (suite *Snake) TestDeploy10MContract() {
	// todo: wait for bitxhub to limit contract size
}

func (suite *Snake) TestInvokeContractWrongArg() {
	address := deployExampleContract(suite)

	result, err := suite.client.InvokeXVMContract(address, "a", nil, rpcx.String("1"), rpcx.Int32(2))
	suite.Require().Nil(err)
	suite.Require().True(result.Status == pb.Receipt_FAILED)
	suite.Require().Contains(string(result.Ret), "not found allocate method")

	// incorrect function params
	result, err = suite.client.InvokeXVMContract(address, "a", nil, rpcx.Int32(1), rpcx.String("2"))
	suite.Require().Nil(err)
	suite.Require().True(result.Status == pb.Receipt_FAILED)

	result, err = suite.client.InvokeXVMContract(address, "a", nil, rpcx.String("1"), rpcx.String("2"))
	suite.Require().Nil(err)
	suite.Require().True(result.Status == pb.Receipt_FAILED)
}

func (suite *Snake) TestDeployContractWrongNumberArg() {
	address := deployExampleContract(suite)

	result, err := suite.client.InvokeXVMContract(address, "a", nil, rpcx.Int32(1), rpcx.Int32(2), rpcx.Int32(3))
	suite.Require().Nil(err)
	suite.Require().True(result.Status == pb.Receipt_FAILED)
}

func deployExampleContract(suite *Snake) *types.Address {
	contract, err := ioutil.ReadFile("testdata/example.wasm")
	suite.Require().Nil(err)

	address, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)
	suite.Require().NotNil(address)
	return address
}
