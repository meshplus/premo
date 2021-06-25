package bxh_tester

import (
	"io/ioutil"

	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"

	"github.com/meshplus/bitxhub-kit/types"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
)

type Model4 struct {
	*Snake
}

func (suite *Model4) Test0411_LegerSet() {
	address := suite.deployLedgerContract()

	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	client := suite.NewClient(pk)
	res, err := client.InvokeXVMContract(address, "state_test_set", nil, rpcx.String("Alice"), rpcx.String("111"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	suite.Require().Equal("1", string(res.Ret))
}

func (suite Model4) Test0412_LegerSetWithValueLoss() {
	address := suite.deployLedgerContract()

	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	client := suite.NewClient(pk)
	res, err := client.InvokeXVMContract(address, "state_test_set", nil, rpcx.String("Alice"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	suite.Require().Contains(string(res.Ret), "Missing 1 argument(s)")
}

func (suite Model4) Test0413_LegerSetWithKVLoss() {
	address := suite.deployLedgerContract()

	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	client := suite.NewClient(pk)
	res, err := client.InvokeXVMContract(address, "state_test_set", nil)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	suite.Require().Contains(string(res.Ret), "Missing 2 argument(s)")
}

func (suite Model4) Test0414_LegerSetWithErrorMethod() {
	address := suite.deployLedgerContract()

	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	client := suite.NewClient(pk)
	res, err := client.InvokeXVMContract(address, "state_test_set111", nil, rpcx.String("Alice"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	suite.Require().Contains(string(res.Ret), "wrong rule contract")
}

func (suite *Model4) Test0415_LegerSetRepeat() {
	address := suite.deployLedgerContract()

	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	client := suite.NewClient(pk)
	res, err := client.InvokeXVMContract(address, "state_test_set", nil, rpcx.String("Alice"), rpcx.String("111"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	suite.Require().Equal("1", string(res.Ret))

	res, err = client.InvokeXVMContract(address, "state_test_set", nil, rpcx.String("Alice"), rpcx.String("111"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	suite.Require().Equal("1", string(res.Ret))
}

func (suite *Model4) Test0416_LegerGetAliceWithoutSet() {
	address := suite.deployLedgerContract()

	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	client := suite.NewClient(pk)
	res, err := client.InvokeXVMContract(address, "state_test_get", nil, rpcx.String("Alice"), rpcx.String("111"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	suite.Require().Contains(string(res.Ret), "Failed to call the `state_test_get` exported function.")
}

func (suite *Model4) Test0417_GetNilWithoutSet() {
	address := suite.deployLedgerContract()

	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	client := suite.NewClient(pk)
	res, err := client.InvokeXVMContract(address, "state_test_get", nil)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	suite.Require().Contains(string(res.Ret), "Missing 2 argument(s)")
}

func (suite *Model4) Test0418_SetAliceGetAlice() {
	address := suite.deployLedgerContract()

	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	client := suite.NewClient(pk)
	res, err := client.InvokeXVMContract(address, "state_test_set", nil, rpcx.String("Alice"), rpcx.String("111"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	suite.Require().Equal("1", string(res.Ret))

	res, err = client.InvokeXVMContract(address, "state_test_get", nil, rpcx.String("Alice"), rpcx.String("111"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	suite.Require().Equal("1", string(res.Ret))
}

func (suite *Model4) Test0419_SetAliceGetBob() {
	address := suite.deployLedgerContract()

	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	client := suite.NewClient(pk)
	res, err := client.InvokeXVMContract(address, "state_test_set", nil, rpcx.String("Alice"), rpcx.String("111"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	suite.Require().Equal("1", string(res.Ret))

	res, err = client.InvokeXVMContract(address, "state_test_get", nil, rpcx.String("Bob"), rpcx.String("111"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	suite.Require().Contains(string(res.Ret), "Failed to call the `state_test_get` exported function.")
}

func (suite *Model4) Test0420_SetAliceGetNil() {
	address := suite.deployLedgerContract()

	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	client := suite.NewClient(pk)
	res, err := client.InvokeXVMContract(address, "state_test_set", nil, rpcx.String("Alice"), rpcx.String("111"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	suite.Require().Equal("1", string(res.Ret))

	res, err = client.InvokeXVMContract(address, "state_test_get", nil)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	suite.Require().Contains(string(res.Ret), "Missing 2 argument(s)")
}

func (suite Model4) Test0421_SetAliceGetAliceRepeat() {
	address := suite.deployLedgerContract()

	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	client := suite.NewClient(pk)
	res, err := client.InvokeXVMContract(address, "state_test_set", nil, rpcx.String("Alice"), rpcx.String("111"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	suite.Require().Equal("1", string(res.Ret))

	res, err = client.InvokeXVMContract(address, "state_test_get", nil, rpcx.String("Alice"), rpcx.String("111"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	suite.Require().Equal("1", string(res.Ret))

	res, err = client.InvokeXVMContract(address, "state_test_get", nil, rpcx.String("Alice"), rpcx.String("111"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	suite.Require().Equal("1", string(res.Ret))
}

func (suite *Model4) deployLedgerContract() *types.Address {
	contract, err := ioutil.ReadFile("testdata/ledger_test_gc.wasm")
	suite.Require().Nil(err)

	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	client := suite.NewClient(pk)
	address, err := client.DeployContract(contract, nil)
	suite.Require().Nil(err)
	suite.Require().NotNil(address)
	return address
}
