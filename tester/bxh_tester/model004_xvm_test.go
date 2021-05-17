package bxh_tester

import (
	"io/ioutil"

	"github.com/meshplus/bitxhub-kit/types"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
)

type Model4 struct {
	*Snake
}

func (suite *Model4) Test0401_LegerSet() {
	address := suite.deployLedgerContract()

	res, err := suite.client.InvokeXVMContract(address, "state_test_set", nil, rpcx.String("Alice"), rpcx.String("111"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	suite.Require().Equal("1", string(res.Ret))
}

func (suite *Model4) Test0402_LegerSetWithValueLoss() {
	address := suite.deployLedgerContract()

	res, err := suite.client.InvokeXVMContract(address, "state_test_set", nil, rpcx.String("Alice"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	suite.Require().Contains(string(res.Ret), "Missing 1 argument(s)")
}

func (suite *Model4) Test0403_LegerSetWithKVLoss() {
	address := suite.deployLedgerContract()

	res, err := suite.client.InvokeXVMContract(address, "state_test_set", nil)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	suite.Require().Contains(string(res.Ret), "Missing 2 argument(s)")
}

func (suite *Model4) Test0404_LegerSetWithErrorMethod() {
	address := suite.deployLedgerContract()

	res, err := suite.client.InvokeXVMContract(address, "state_test_set111", nil, rpcx.String("Alice"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	suite.Require().Contains(string(res.Ret), "wrong rule contract")
}

func (suite *Model4) Test0405_LegerSetRepeat() {
	address := suite.deployLedgerContract()

	res, err := suite.client.InvokeXVMContract(address, "state_test_set", nil, rpcx.String("Alice"), rpcx.String("111"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	suite.Require().Equal("1", string(res.Ret))

	res, err = suite.client.InvokeXVMContract(address, "state_test_set", nil, rpcx.String("Alice"), rpcx.String("111"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	suite.Require().Equal("1", string(res.Ret))
}

func (suite *Model4) Test0406_LegerGetAliceWithoutSet() {
	address := suite.deployLedgerContract()

	res, err := suite.client.InvokeXVMContract(address, "state_test_get", nil, rpcx.String("Alice"), rpcx.String("111"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	suite.Require().Contains(string(res.Ret), "Failed to call the `state_test_get` exported function.")
}

func (suite *Model4) Test0407_GetNilWithoutSet() {
	address := suite.deployLedgerContract()

	res, err := suite.client.InvokeXVMContract(address, "state_test_get", nil)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	suite.Require().Contains(string(res.Ret), "Missing 2 argument(s)")
}

func (suite *Model4) Test0408_SetAliceGetAlice() {
	address := suite.deployLedgerContract()

	res, err := suite.client.InvokeXVMContract(address, "state_test_set", nil, rpcx.String("Alice"), rpcx.String("111"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	suite.Require().Equal("1", string(res.Ret))

	res, err = suite.client.InvokeXVMContract(address, "state_test_get", nil, rpcx.String("Alice"), rpcx.String("111"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	suite.Require().Equal("1", string(res.Ret))
}

func (suite *Model4) Test0409_SetAliceGetBob() {
	address := suite.deployLedgerContract()

	res, err := suite.client.InvokeXVMContract(address, "state_test_set", nil, rpcx.String("Alice"), rpcx.String("111"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	suite.Require().Equal("1", string(res.Ret))

	res, err = suite.client.InvokeXVMContract(address, "state_test_get", nil, rpcx.String("Bob"), rpcx.String("111"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	suite.Require().Contains(string(res.Ret), "Failed to call the `state_test_get` exported function.")
}

func (suite *Model4) Test0410_SetAliceGetNil() {
	address := suite.deployLedgerContract()

	res, err := suite.client.InvokeXVMContract(address, "state_test_set", nil, rpcx.String("Alice"), rpcx.String("111"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	suite.Require().Equal("1", string(res.Ret))

	res, err = suite.client.InvokeXVMContract(address, "state_test_get", nil)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	suite.Require().Contains(string(res.Ret), "Missing 2 argument(s)")
}

func (suite *Model4) Test0411_SetAliceGetAliceRepeat() {
	address := suite.deployLedgerContract()

	res, err := suite.client.InvokeXVMContract(address, "state_test_set", nil, rpcx.String("Alice"), rpcx.String("111"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	suite.Require().Equal("1", string(res.Ret))

	res, err = suite.client.InvokeXVMContract(address, "state_test_get", nil, rpcx.String("Alice"), rpcx.String("111"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	suite.Require().Equal("1", string(res.Ret))

	res, err = suite.client.InvokeXVMContract(address, "state_test_get", nil, rpcx.String("Alice"), rpcx.String("111"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	suite.Require().Equal("1", string(res.Ret))
}

func (suite *Snake) deployLedgerContract() *types.Address {
	contract, err := ioutil.ReadFile("testdata/ledger_test_gc.wasm")
	suite.Require().Nil(err)

	address, err := suite.client.DeployContract(contract, nil)
	suite.Require().Nil(err)
	suite.Require().NotNil(address)
	return address
}
