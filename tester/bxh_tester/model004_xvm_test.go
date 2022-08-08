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

func (suite *Model4) SetupTest() {
	suite.T().Parallel()
}

//tc：部署账本合约后调用state_test_set方法设置键值对为（Alice，111），合约调用成功
func (suite *Model4) Test0401_LegerSetIsSuccess() {
	address := suite.DeployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeXVMContract(address, "state_test_set", nil, rpcx.String("Alice"), rpcx.String("111"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
}

//tc：部署账本合约后调用state_test_set方法设置键值对为（Alice，），合约调用失败
func (suite *Model4) Test0402_LegerSetWithValueLossIsFail() {
	address := suite.DeployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeXVMContract(address, "state_test_set", nil, rpcx.String("Alice"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	suite.Require().Contains(string(res.Ret), "expected 2 arguments, got 1")
}

//tc：部署账本合约后调用state_test_set方法设置键值对为（，），合约调用失败
func (suite *Model4) Test0403_LegerSetWithKVLossIsFail() {
	address := suite.DeployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeXVMContract(address, "state_test_set", nil)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	suite.Require().Contains(string(res.Ret), "expected 2 arguments, got 0")
}

//tc：部署账本合约后调用state_test_set111方法设置键值对为（Alice，111），合约调用失败
func (suite *Model4) Test0404_LegerSetWithErrorMethodIsFail() {
	address := suite.DeployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeXVMContract(address, "state_test_set111", nil, rpcx.String("Alice"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	suite.Require().Contains(string(res.Ret), "no such method")
}

//tc：部署账本合约后调用state_test_set方法设置键值对为（Alice，111）,重复调用，合约调用成功
func (suite *Model4) Test0405_LegerSetRepeatIsSuccess() {
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

//tc：部署账本合约后直接调用state_test_get方法获取Alice的值，合约调用失败
func (suite *Model4) Test0406_LegerGetAliceWithoutSetIsFail() {
	address := suite.DeployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeXVMContract(address, "state_test_get", nil, rpcx.String("Alice"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
}

//tc：部署账本合约后直接调用state_test_get方法获取nil的值，合约调用失败
func (suite *Model4) Test0407_GetNilWithoutSetIsFail() {
	address := suite.DeployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeXVMContract(address, "state_test_get", nil)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	suite.Require().Contains(string(res.Ret), "expected 1 arguments, got 0")
}

//tc：部署账本合约后设置键值对为（Alice，111），调用state_test_get方法获取Alice的值,合约调用成功
func (suite *Model4) Test0408_SetAliceGetAliceIsSuccess() {
	address := suite.DeployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeXVMContract(address, "state_test_set", nil, rpcx.String("Alice"), rpcx.String("111"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	res, err = client.InvokeXVMContract(address, "state_test_get", nil, rpcx.String("Alice"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	suite.Require().Equal("111", string(res.Ret))
}

//tc：部署账本合约后设置键值对为（Alice，111），调用state_test_get方法获取Bob的值，合约调用失败
func (suite *Model4) Test0409_SetAliceGetBobIsFail() {
	address := suite.DeployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeXVMContract(address, "state_test_set", nil, rpcx.String("Alice"), rpcx.String("111"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	res, err = client.InvokeXVMContract(address, "state_test_get", nil, rpcx.String("Bob"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
}

//tc：部署账本合约后设置键值对为（Alice，111），调用state_test_get方法获取nil的值，合约调用失败
func (suite *Model4) Test0410_SetAliceGetNilIsFail() {
	address := suite.DeployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeXVMContract(address, "state_test_set", nil, rpcx.String("Alice"), rpcx.String("111"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	res, err = client.InvokeXVMContract(address, "state_test_get", nil)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	suite.Require().Contains(string(res.Ret), "expected 1 arguments, got 0")
}

//tc：部署账本合约后设置键值对为（Alice，111），调用state_test_get方法获取Alice的值，重复调用，合约调用成功
func (suite *Model4) Test0411_SetAliceGetAliceRepeatIsSuccess() {
	address := suite.DeployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeXVMContract(address, "state_test_set", nil, rpcx.String("Alice"), rpcx.String("111"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	res, err = client.InvokeXVMContract(address, "state_test_get", nil, rpcx.String("Alice"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	suite.Require().Equal("111", string(res.Ret))
	res, err = client.InvokeXVMContract(address, "state_test_get", nil, rpcx.String("Alice"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	suite.Require().Equal("111", string(res.Ret))
}

//tc：部署结果合约，获取当前的块高，合约调用成功
func (suite *Model4) Test0412_GetCurrentHeightIsSuccess() {
	address := suite.DeployResultContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeXVMContract(address, "test_current_height", nil)
	suite.Require().Nil(err)
	meta, err := client.GetChainMeta()
	suite.Require().Nil(err)
	height, err := strconv.ParseUint(string(res.Ret), 10, 64)
	suite.Require().Nil(err)
	suite.Require().LessOrEqual(height, meta.Height-1)
}

//tc：部署结果合约，获取当前交易的交易hash，合约调用成功
func (suite *Model4) Test0413_GetTxHashIsSuccess() {
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

func (suite *Snake) DeployResultContract() *types.Address {
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
