package bxh_tester

import (
	"io/ioutil"
	"strconv"

	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-kit/types"
	"github.com/meshplus/bitxhub-model/pb"

	rpcx "github.com/meshplus/go-bitxhub-client"
)

type Model4 struct {
	*Snake
}

func (suite Model4) SetupTest() {
	suite.T().Parallel()
}

func (suite Model4) Test0401_LegerSet() {
	address := suite.DeployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeXVMContract(address, "state_test_set", nil, rpcx.String("Alice"), rpcx.String("111"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	suite.Require().Equal("1", string(res.Ret))
}

func (suite Model4) Test0402_LegerSetWithValueLoss() {
	address := suite.DeployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeXVMContract(address, "state_test_set", nil, rpcx.String("Alice"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	suite.Require().Contains(string(res.Ret), "Missing 1 argument(s)")
}

func (suite Model4) Test0403_LegerSetWithKVLoss() {
	address := suite.DeployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeXVMContract(address, "state_test_set", nil)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	suite.Require().Contains(string(res.Ret), "Missing 2 argument(s)")
}

func (suite Model4) Test0404_LegerSetWithErrorMethod() {
	address := suite.DeployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeXVMContract(address, "state_test_set111", nil, rpcx.String("Alice"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	suite.Require().Contains(string(res.Ret), "does not exist")
}

func (suite Model4) Test0405_LegerSetRepeat() {
	address := suite.DeployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
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

func (suite Model4) Test0406_LegerGetAliceWithoutSet() {
	address := suite.DeployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeXVMContract(address, "state_test_get", nil, rpcx.String("Alice"), rpcx.String("111"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	suite.Require().Equal("0", string(res.Ret))
}

func (suite Model4) Test0407_GetNilWithoutSet() {
	address := suite.DeployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeXVMContract(address, "state_test_get", nil)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	suite.Require().Contains(string(res.Ret), "Missing 2 argument(s)")
}

func (suite Model4) Test0408_SetAliceGetAlice() {
	address := suite.DeployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
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

func (suite Model4) Test0409_SetAliceGetBob() {
	address := suite.DeployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeXVMContract(address, "state_test_set", nil, rpcx.String("Alice"), rpcx.String("111"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	suite.Require().Equal("1", string(res.Ret))
	res, err = client.InvokeXVMContract(address, "state_test_get", nil, rpcx.String("Bob"), rpcx.String("111"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	suite.Require().Equal("0", string(res.Ret))
}

func (suite Model4) Test0410_SetAliceGetNil() {
	address := suite.DeployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
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

func (suite Model4) Test0411_SetAliceGetAliceRepeat() {
	address := suite.DeployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
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

func (suite Model4) Test0412_GetCurrentHeight() {
	address := suite.DeployResultContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeXVMContract(address, "get_current_height", nil)
	suite.Require().Nil(err)
	meta, err := client.GetChainMeta()
	suite.Require().Nil(err)
	suite.Require().LessOrEqual(string(res.Ret), strconv.FormatUint(meta.Height-1, 10))
}

func (suite Model4) Test0412_GetTxHash() {
	address := suite.DeployResultContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeXVMContract(address, "test_tx_hash", nil)
	suite.Require().Nil(err)
	suite.Require().Equal(string(res.Ret), res.TxHash.String())
}

func (suite *Snake) DeployLedgerContract() *types.Address {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	contract, err := ioutil.ReadFile("testdata/ledger_test_gc.wasm")
	suite.Require().Nil(err)
	address, err := client.DeployContract(contract, nil)
	suite.Require().Nil(err)
	suite.Require().NotNil(address)
	return address
}

func (suite Snake) DeployResultContract() *types.Address {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	contract, err := ioutil.ReadFile("testdata/result.wasm")
	suite.Require().Nil(err)
	address, err := client.DeployContract(contract, nil)
	suite.Require().Nil(err)
	suite.Require().NotNil(address)
	return address
}
