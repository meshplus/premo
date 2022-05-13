package bxh_tester

import (
	"encoding/json"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	crypto2 "github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
	"github.com/meshplus/premo/internal/repo"
	"github.com/onrik/ethrpc"
)

type Model18 struct {
	*Snake
	client *ethrpc.EthRPC
}

//tc：获取eth协议版本，获取成功
func (suite Model18) Test1801_GetProtocolVersionIsSuccess() {
	version, err := suite.client.EthProtocolVersion()
	suite.Require().Nil(err)
	suite.Require().Equal("0x41", version)
}

//tc：获取bitxhub chainID，获取成功
func (suite Model18) Test1802_GetChainIDIsSuccess() {
	resp, err := suite.client.Call("eth_chainId")
	suite.Require().Nil(err)
	var res string
	err = json.Unmarshal(resp, &res)
	suite.Require().Nil(err)
	suite.Require().Equal(hexutil.EncodeUint64(1356), res)
}

//tc：获取当前节点是否在挖矿，获取成功
func (suite Model18) Test1803_GetMiningStatusIsSuccess() {
	status, err := suite.client.EthMining()
	suite.Require().Nil(err)
	suite.Require().Equal(false, status)
}

//tc：获取当前节点的算力，获取成功
func (suite Model18) Test1804_GetHashRateIsSuccess() {
	hashRate, err := suite.client.EthHashrate()
	suite.Require().Nil(err)
	suite.Require().Equal(0, hashRate)
}

//tc：获取当前gas的价格，获取成功
func (suite Model18) Test1805_GetGasPriceIsSuccess() {
	price, err := suite.client.EthGasPrice()
	suite.Require().Nil(err)
	suite.Require().Equal("50000", price.String())
}

//tc：获取gas上限的建议，获取成功
func (suite Model18) Test1806_GetMaxPriorityFeePerGasIsSuccess() {
	resp, err := suite.client.Call("eth_maxPriorityFeePerGas")
	suite.Require().Nil(err)
	var res string
	err = json.Unmarshal(resp, &res)
	suite.Require().Nil(err)
	suite.Require().Equal("0x0", res)
}

//tc：获取当前的区块高度，获取成功
func (suite Model18) Test1807_GetBlockNumberIsSuccess() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	meta, err := client.GetChainMeta()
	suite.Require().Nil(err)
	number, err := suite.client.EthBlockNumber()
	suite.Require().Nil(err)
	suite.Require().Equal(int(meta.Height), number)
}

//tc：根据正确的账户地址获取账户金额，获取成功
func (suite Model18) Test1808_GetBalanceIsSuccess() {
	_, address, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.TransferFromAdmin(address.String(), "1")
	suite.Require().Nil(err)
	balance, err := suite.client.EthGetBalance(address.String(), "latest")
	suite.Require().Nil(err)
	suite.Require().Equal("1000000000000000000", balance.String())
}

//tc：根据正确的账户地址和正确的关键字获取存储，获取成功
func (suite Model18) Test1809_GetStorageAtIsSuccess() {
	pk, err := asym.GenerateKeyPair(crypto2.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.StoreContractAddr.Address(), "Set", nil, pb.String("key"), pb.String("value"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	data, err := suite.client.Call("eth_getStorageAt", constant.StoreContractAddr.String(), "key", "latest")
	suite.Require().Nil(err)
	str1 := strings.ReplaceAll(string(data), "\"", "")
	bytes, err := hexutil.Decode(str1)
	suite.Require().Nil(err)
	str2 := strings.ReplaceAll(string(bytes), "\"", "")
	suite.Require().Equal("value", str2)
}

//tc：根据正确的账户地址获取当前账户发生的交易数量，获取成功
func (suite Model18) Test1810_GetTransactionCountIsSuccess() {
	_, address, err := repo.KeyPriv()
	suite.Require().Nil(err)
	nonce, err := suite.client.EthGetTransactionCount(address.String(), "latest")
	suite.Require().Nil(err)
	suite.Require().Equal(0, nonce)
}

//tc：根据正确的区块hash获取当前区块的交易数量，获取成功
func (suite Model18) Test1811_GetBlockTransactionCountByHashIsSuccess() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	meta, err := client.GetChainMeta()
	suite.Require().Nil(err)
	block, err := client.GetBlock(meta.BlockHash.String(), pb.GetBlockRequest_HASH)
	suite.Require().Nil(err)
	count, err := suite.client.EthGetBlockTransactionCountByHash(meta.BlockHash.String())
	suite.Require().Nil(err)
	suite.Require().Equal(len(block.Transactions.Transactions), count)
}

//tc：根据不存在区块hash获取当前区块的交易数量，获取失败
func (suite Model18) Test1812_GetBlockTransactionCountByHashWithErrorBlockHashIsFail() {
	_, err := suite.client.EthGetBlockTransactionCountByHash("0x0000000000000000000000000000000000000000000000000000000000000000")
	suite.Require().NotNil(err)
}

//tc：根据正确的区块高度获取当前区块的交易数量，获取成功
func (suite Model18) Test1813_GetBlockTransactionCountByNumberIsSuccess() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	meta, err := client.GetChainMeta()
	suite.Require().Nil(err)
	block, err := client.GetBlock(strconv.FormatUint(meta.Height, 10), pb.GetBlockRequest_HEIGHT)
	suite.Require().Nil(err)
	count, err := suite.client.EthGetBlockTransactionCountByNumber(int(meta.Height))
	suite.Require().Nil(err)
	suite.Require().Equal(len(block.Transactions.Transactions), count)
}

//tc：根据不存在的区块高度获取当前区块的交易数量，获取失败
func (suite Model18) Test1814_GetBlockTransactionCountByNumberWithNoExistHeightIsFail() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	meta, err := client.GetChainMeta()
	suite.Require().Nil(err)
	_, err = suite.client.EthGetBlockTransactionCountByNumber(int(meta.Height) + 1)
	suite.Require().NotNil(err)
}

//tc：根据正确的合约地址获取合约，获取成功
func (suite Model18) Test1815_GetCodeIsSuccess() {
	address, err := suite.DeploySimpleRule()
	suite.Require().Nil(err)
	code, err := suite.client.EthGetCode(address, "latest")
	suite.Require().Nil(err)
	suite.Require().NotNil(code)
}

//tc：根据不存在的的合约地址获取合约，获取失败
func (suite Model18) Test1816_GetCodeWithNoExistAddrIsFail() {
	code, err := suite.client.EthGetCode("0x0000000000000000000000000000000000000000", "latest")
	suite.Require().Nil(err)
	suite.Require().Equal("0x", code)
}

//tc：根据正确的交易hash获取日志，获取成功
func (suite Model18) Test1817_GetTransactionLogsIsSuccess() {

}

//tc：根据不存在的交易hash获取日志，获取失败
func (suite Model18) Test1818_GetTransactionLogsWithNoExistHashIsFail() {
	_, err := suite.client.Call("eth_getTransactionLogs", "0x0000000000000000000000000000000000000000000000000000000000000000")
	suite.Require().NotNil(err)
}

//tc：发送正确的交易，交易发送成功
func (suite Model18) Test1819_SendRawTransactionIsSuccess() {
	pk, from, err := repo.Node1Priv()
	suite.Require().Nil(err)
	_, to, err := repo.KeyPriv()
	suite.Require().Nil(err)
	too := common.HexToAddress(to.String())
	amount := big.NewInt(1000000000)
	gasLimit := uint64(21000)
	gasPrice, err := suite.client.EthGasPrice()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	nonce, err := client.GetPendingNonceByAccount(from.String())
	suite.Require().Nil(err)

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &too,
		Value:    amount,
		Gas:      gasLimit,
		GasPrice: &gasPrice,
		Data:     []byte{},
	})

	bytes, err := pk.Bytes()
	suite.Require().Nil(err)
	privateKey, err := crypto.HexToECDSA(common.Bytes2Hex(bytes))
	suite.Require().Nil(err)
	signTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(1356)), privateKey)
	suite.Require().Nil(err)
	data, err := signTx.MarshalBinary()
	rawTx := hexutil.Bytes(data).String()

	hash, err := suite.client.EthSendRawTransaction(rawTx)
	suite.Require().Nil(err)
	res, err := client.GetReceipt(hash)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	ret, err := client.GetAccountBalance(to.String())
	suite.Require().Nil(err)
	account := Account{}
	err = json.Unmarshal(ret.Data, &account)
	suite.Require().Nil(err)
	suite.Require().Equal("1000000000", account.Balance.String())
}

//tc：发送to为空的交易，交易发送失败
func (suite Model18) Test1820_SendRawTransactionWithEmptyToIsFail() {
	pk, from, err := repo.Node1Priv()
	suite.Require().Nil(err)
	amount := big.NewInt(1000000000)
	gasLimit := uint64(21000)
	gasPrice, err := suite.client.EthGasPrice()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	nonce, err := client.GetPendingNonceByAccount(from.String())
	suite.Require().Nil(err)

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       nil,
		Value:    amount,
		Gas:      gasLimit,
		GasPrice: &gasPrice,
		Data:     []byte{},
	})

	bytes, err := pk.Bytes()
	suite.Require().Nil(err)
	privateKey, err := crypto.HexToECDSA(common.Bytes2Hex(bytes))
	suite.Require().Nil(err)
	signTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(1356)), privateKey)
	suite.Require().Nil(err)
	data, err := signTx.MarshalBinary()
	rawTx := hexutil.Bytes(data).String()

	_, err = suite.client.EthSendRawTransaction(rawTx)
	suite.Require().NotNil(err)
}

//tc：发送签名为空的交易，交易发送失败
func (suite Model18) Test1821_SendRawTransactionWithEmptySignatureIsFail() {
	pk, from, err := repo.Node1Priv()
	suite.Require().Nil(err)
	_, to, err := repo.KeyPriv()
	suite.Require().Nil(err)
	too := common.HexToAddress(to.String())
	amount := big.NewInt(1000000000)
	gasLimit := uint64(21000)
	gasPrice, err := suite.client.EthGasPrice()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	nonce, err := client.GetPendingNonceByAccount(from.String())
	suite.Require().Nil(err)

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &too,
		Value:    amount,
		Gas:      gasLimit,
		GasPrice: &gasPrice,
		Data:     []byte{},
	})
	data, err := tx.MarshalBinary()
	rawTx := hexutil.Bytes(data).String()

	_, err = suite.client.EthSendRawTransaction(rawTx)
	suite.Require().NotNil(err)
}

//tc：发送签名错误的交易，交易发送失败
func (suite Model18) Test1822_SendRawTransactionWithErrorSignatureIsFail() {
	pk1, from, err := repo.Node1Priv()
	suite.Require().Nil(err)
	pk2, to, err := repo.KeyPriv()
	suite.Require().Nil(err)
	too := common.HexToAddress(to.String())
	amount := big.NewInt(1000000000)
	gasLimit := uint64(21000)
	gasPrice, err := suite.client.EthGasPrice()
	suite.Require().Nil(err)
	client := suite.NewClient(pk1)
	nonce, err := client.GetPendingNonceByAccount(from.String())
	suite.Require().Nil(err)

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &too,
		Value:    amount,
		Gas:      gasLimit,
		GasPrice: &gasPrice,
		Data:     []byte{},
	})
	bytes, err := pk2.Bytes()
	suite.Require().Nil(err)
	privateKey, err := crypto.HexToECDSA(common.Bytes2Hex(bytes))
	suite.Require().Nil(err)
	signTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(1356)), privateKey)
	suite.Require().Nil(err)
	data, err := signTx.MarshalBinary()
	rawTx := hexutil.Bytes(data).String()

	_, err = suite.client.EthSendRawTransaction(rawTx)
	suite.Require().NotNil(err)
}

//tc：发送price低于当前系统price的交易，交易发送失败
func (suite Model18) Test1823_SendRawTransactionWithLessPriceIsFail() {
	pk, from, err := repo.Node1Priv()
	suite.Require().Nil(err)
	_, to, err := repo.KeyPriv()
	suite.Require().Nil(err)
	too := common.HexToAddress(to.String())
	amount := big.NewInt(1000000000)
	gasLimit := uint64(21000)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	nonce, err := client.GetPendingNonceByAccount(from.String())
	suite.Require().Nil(err)

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &too,
		Value:    amount,
		Gas:      gasLimit,
		GasPrice: big.NewInt(49999),
		Data:     []byte{},
	})

	bytes, err := pk.Bytes()
	suite.Require().Nil(err)
	privateKey, err := crypto.HexToECDSA(common.Bytes2Hex(bytes))
	suite.Require().Nil(err)
	signTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(1356)), privateKey)
	suite.Require().Nil(err)
	data, err := signTx.MarshalBinary()
	rawTx := hexutil.Bytes(data).String()

	_, err = suite.client.EthSendRawTransaction(rawTx)
	suite.Require().NotNil(err)
}

//tc：根据正确的参数调用内置合约，调用成功
func (suite Model18) Test1824_CallIsSuccess() {
	hash, err := suite.client.EthSendRawTransaction("0xf901c28082c3508301eaed8080b90170608060405234801561001057600080fd5b50610150806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80632e64cec11461003b5780636057361d14610059575b600080fd5b610043610075565b60405161005091906100d9565b60405180910390f35b610073600480360381019061006e919061009d565b61007e565b005b60008054905090565b8060008190555050565b60008135905061009781610103565b92915050565b6000602082840312156100b3576100b26100fe565b5b60006100c184828501610088565b91505092915050565b6100d3816100f4565b82525050565b60006020820190506100ee60008301846100ca565b92915050565b6000819050919050565b600080fd5b61010c816100f4565b811461011757600080fd5b5056fea26469706673582212209a159a4f3847890f10bfb87871a61eba91c5dbf5ee3cf6398207e292eee22a1664736f6c63430008070033820abca004262a31c3f0ef18950fb1e2d53d894bd2d1096a5140e3575eff53924baac054a07cbc9db70bd61fb4b6e4979b8f35a53675a745dd4dfedb595bd27fdc9aa0fb91")
	suite.Require().Nil(err)
	pk, from, err := repo.Node1Priv()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.GetReceipt(hash)
	suite.Require().Nil(err)

	// set value 100
	_, err = suite.client.EthSendRawTransaction("0xf8870782c350825d74941dfcd0c15e69737b346a08fe3ecd1ca7d9d482bb80a46057361d0000000000000000000000000000000000000000000000000000000000000458820abca048e58c7dc5932d08ef0184b59c7b3470c45c1f32e151b20f8f12f3085e1e68aca06993ef8e147063467be9bed9aba8075922f36b378f6869d6fd03d09e6db47fb2")
	suite.Require().Nil(err)

	// get value
	tx := ethrpc.T{
		From: from.String(),
		To:   res.ContractAddress.String(),
		Data: "0x2e64cec1",
	}
	value, err := suite.client.EthCall(tx, "latest")
	suite.Require().Nil(err)
	suite.Require().Equal("0x0000000000000000000000000000000000000000000000000000000000000064", value)
}

//tc：根据错误的参数调用内置合约，调用失败
func (suite Model18) Test1825_CallWithErrorArgsIsFail() {
	hash, err := suite.client.EthSendRawTransaction("0xf901c28082c3508301eaed8080b90170608060405234801561001057600080fd5b50610150806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80632e64cec11461003b5780636057361d14610059575b600080fd5b610043610075565b60405161005091906100d9565b60405180910390f35b610073600480360381019061006e919061009d565b61007e565b005b60008054905090565b8060008190555050565b60008135905061009781610103565b92915050565b6000602082840312156100b3576100b26100fe565b5b60006100c184828501610088565b91505092915050565b6100d3816100f4565b82525050565b60006020820190506100ee60008301846100ca565b92915050565b6000819050919050565b600080fd5b61010c816100f4565b811461011757600080fd5b5056fea26469706673582212209a159a4f3847890f10bfb87871a61eba91c5dbf5ee3cf6398207e292eee22a1664736f6c63430008070033820abca004262a31c3f0ef18950fb1e2d53d894bd2d1096a5140e3575eff53924baac054a07cbc9db70bd61fb4b6e4979b8f35a53675a745dd4dfedb595bd27fdc9aa0fb91")
	suite.Require().Nil(err)
	pk, from, err := repo.Node1Priv()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.GetReceipt(hash)
	suite.Require().Nil(err)

	tx := ethrpc.T{
		From: from.String(),
		To:   res.ContractAddress.String(),
		Data: "0x2e64cec2",
	}
	_, err = suite.client.EthCall(tx, "latest")
	suite.Require().NotNil(err)
}

//tc：根据正确的参数获取调用内置合约的gas limit，获取成功
func (suite Model18) Test1826_EstimateGasIsSuccess() {
	hash, err := suite.client.EthSendRawTransaction("0xf901c28082c3508301eaed8080b90170608060405234801561001057600080fd5b50610150806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80632e64cec11461003b5780636057361d14610059575b600080fd5b610043610075565b60405161005091906100d9565b60405180910390f35b610073600480360381019061006e919061009d565b61007e565b005b60008054905090565b8060008190555050565b60008135905061009781610103565b92915050565b6000602082840312156100b3576100b26100fe565b5b60006100c184828501610088565b91505092915050565b6100d3816100f4565b82525050565b60006020820190506100ee60008301846100ca565b92915050565b6000819050919050565b600080fd5b61010c816100f4565b811461011757600080fd5b5056fea26469706673582212209a159a4f3847890f10bfb87871a61eba91c5dbf5ee3cf6398207e292eee22a1664736f6c63430008070033820abca004262a31c3f0ef18950fb1e2d53d894bd2d1096a5140e3575eff53924baac054a07cbc9db70bd61fb4b6e4979b8f35a53675a745dd4dfedb595bd27fdc9aa0fb91")
	suite.Require().Nil(err)
	pk, from, err := repo.Node1Priv()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.GetReceipt(hash)
	suite.Require().Nil(err)

	tx := ethrpc.T{
		From: from.String(),
		To:   res.ContractAddress.String(),
		Data: "0x2e64cec1",
	}
	value, err := suite.client.EthEstimateGas(tx)
	suite.Require().Nil(err)
	suite.Require().Greater(value, 21000)
}

//tc：根据错误的参数获取调用内置合约的gas limit，获取失败
func (suite Model18) Test1827_EstimateGasWithErrorArgsIsFall() {
	hash, err := suite.client.EthSendRawTransaction("0xf901c28082c3508301eaed8080b90170608060405234801561001057600080fd5b50610150806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80632e64cec11461003b5780636057361d14610059575b600080fd5b610043610075565b60405161005091906100d9565b60405180910390f35b610073600480360381019061006e919061009d565b61007e565b005b60008054905090565b8060008190555050565b60008135905061009781610103565b92915050565b6000602082840312156100b3576100b26100fe565b5b60006100c184828501610088565b91505092915050565b6100d3816100f4565b82525050565b60006020820190506100ee60008301846100ca565b92915050565b6000819050919050565b600080fd5b61010c816100f4565b811461011757600080fd5b5056fea26469706673582212209a159a4f3847890f10bfb87871a61eba91c5dbf5ee3cf6398207e292eee22a1664736f6c63430008070033820abca004262a31c3f0ef18950fb1e2d53d894bd2d1096a5140e3575eff53924baac054a07cbc9db70bd61fb4b6e4979b8f35a53675a745dd4dfedb595bd27fdc9aa0fb91")
	suite.Require().Nil(err)
	pk, from, err := repo.Node1Priv()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.GetReceipt(hash)
	suite.Require().Nil(err)

	tx := ethrpc.T{
		From: from.String(),
		To:   res.ContractAddress.String(),
		Data: "0x2e64cec2",
	}
	_, err = suite.client.EthEstimateGas(tx)
	suite.Require().NotNil(err)
}

//tc：根据正确的区块hash获取区块内全部的交易hash，获取成功
func (suite Model18) Test1828_GetBlockByHashIsSuccess() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	meta, err := client.GetChainMeta()
	suite.Require().Nil(err)
	txs, err := suite.client.EthGetBlockByHash(meta.BlockHash.String(), false)
	suite.Require().Nil(err)
	suite.Require().NotNil(txs)
}

//tc：根据错误的区块hash获取区块内全部的交易hash，获取失败
func (suite Model18) Test1829_GetBlockByHashWithErrorHashIsFail() {
	_, err := suite.client.EthGetBlockByHash("0x0000000000000000000000000000000000000000", false)
	suite.Require().NotNil(err)
}

//tc：根据正确的区块hash获取区块内完整的交易，获取成功
func (suite Model18) Test1830_GetBlockByHasFullIsSuccess() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	meta, err := client.GetChainMeta()
	suite.Require().Nil(err)
	txs, err := suite.client.EthGetBlockByHash(meta.BlockHash.String(), true)
	suite.Require().Nil(err)
	suite.Require().NotNil(txs)
}

//tc：根据错误的区块hash获取区块内完整的交易，获取失败
func (suite Model18) Test1831_GetBlockByHashFullWithErrorHashIsSuccess() {
	_, err := suite.client.EthGetBlockByHash("0x0000000000000000000000000000000000000000", true)
	suite.Require().NotNil(err)
}

//tc：根据正确的区块高度获取区块内全部的交易hash，获取成功
func (suite Model18) Test1832_GetBlockByNumberIsSuccess() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	meta, err := client.GetChainMeta()
	suite.Require().Nil(err)
	txs, err := suite.client.EthGetBlockByNumber(int(meta.Height), false)
	suite.Require().Nil(err)
	suite.Require().NotNil(txs)
}

//tc：根据错误的区块高度获取区块内全部的交易hash，获取失败
func (suite Model18) Test1833_GetBlockByNumberWithErrorHeightIsFail() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	meta, err := client.GetChainMeta()
	suite.Require().Nil(err)
	_, err = suite.client.EthGetBlockByNumber(int(meta.Height)+1, false)
	suite.Require().NotNil(err)
}

//tc：根据正确的区块高度获取区块内完整的交易，获取成功
func (suite Model18) Test1834_GetBlockByNumberFullIsSuccess() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	meta, err := client.GetChainMeta()
	suite.Require().Nil(err)
	txs, err := suite.client.EthGetBlockByNumber(int(meta.Height), true)
	suite.Require().Nil(err)
	suite.Require().NotNil(txs)
}

//tc：根据错误的区块高度获取区块内完整的交易，获取失败
func (suite Model18) Test1835_GetBlockByNumberFullWithErrorHeightIsFail() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	meta, err := client.GetChainMeta()
	suite.Require().Nil(err)
	_, err = suite.client.EthGetBlockByNumber(int(meta.Height)+1, false)
	suite.Require().NotNil(err)
}

//tc：根据正确的交易hash获取完整的交易，获取成功
func (suite Model18) Test1836_GetTransactionByHashIsSuccess() {
	pk, from, err := repo.Node1Priv()
	suite.Require().Nil(err)
	_, to, err := repo.KeyPriv()
	suite.Require().Nil(err)
	too := common.HexToAddress(to.String())
	amount := big.NewInt(1000000000)
	gasLimit := uint64(21000)
	gasPrice, err := suite.client.EthGasPrice()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	nonce, err := client.GetPendingNonceByAccount(from.String())
	suite.Require().Nil(err)

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &too,
		Value:    amount,
		Gas:      gasLimit,
		GasPrice: &gasPrice,
		Data:     []byte{},
	})

	bytes, err := pk.Bytes()
	suite.Require().Nil(err)
	privateKey, err := crypto.HexToECDSA(common.Bytes2Hex(bytes))
	suite.Require().Nil(err)
	signTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(1356)), privateKey)
	suite.Require().Nil(err)
	data, err := signTx.MarshalBinary()
	rawTx := hexutil.Bytes(data).String()

	hash, err := suite.client.EthSendRawTransaction(rawTx)
	suite.Require().Nil(err)
	time.Sleep(1 * time.Second)
	res, err := suite.client.EthGetTransactionByHash(hash)
	suite.Require().Nil(err)
	suite.Require().NotNil(res)
}

//tc：根据错误的交易hash获取完整的交易，获取失败
func (suite Model18) Test1837_GetTransactionByHashWithErrorHashIsSuccess() {
	_, err := suite.client.EthGetTransactionByHash("0x0000000000000000000000000000000000000000000000000000000000000000")
	suite.Require().NotNil(err)
}

//tc：根据正确的区块hash和正确的交易index获取完整交易，获取成功
func (suite Model18) Test1838_GetTransactionByBlockHashAndIndexIsSuccess() {
	pk, from, err := repo.Node1Priv()
	suite.Require().Nil(err)
	_, to, err := repo.KeyPriv()
	suite.Require().Nil(err)
	too := common.HexToAddress(to.String())
	amount := big.NewInt(1000000000)
	gasLimit := uint64(21000)
	gasPrice, err := suite.client.EthGasPrice()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	nonce, err := client.GetPendingNonceByAccount(from.String())
	suite.Require().Nil(err)

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &too,
		Value:    amount,
		Gas:      gasLimit,
		GasPrice: &gasPrice,
		Data:     []byte{},
	})

	bytes, err := pk.Bytes()
	suite.Require().Nil(err)
	privateKey, err := crypto.HexToECDSA(common.Bytes2Hex(bytes))
	suite.Require().Nil(err)
	signTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(1356)), privateKey)
	suite.Require().Nil(err)
	data, err := signTx.MarshalBinary()
	rawTx := hexutil.Bytes(data).String()

	hash, err := suite.client.EthSendRawTransaction(rawTx)
	suite.Require().Nil(err)
	time.Sleep(1 * time.Second)
	res1, err := suite.client.EthGetTransactionByHash(hash)
	suite.Require().Nil(err)
	res2, err := suite.client.EthGetTransactionByBlockHashAndIndex(res1.BlockHash, *res1.TransactionIndex)
	suite.Require().Nil(err)
	suite.Require().Equal(res1, res2)
}

//tc：根据错误的区块hash和正确的交易index获取完整交易，获取失败
func (suite Model18) Test1839_GetTransactionByBlockHashAndIndexWithErrorHashIsFail() {
	_, err := suite.client.EthGetTransactionByBlockHashAndIndex("0x0000000000000000000000000000000000000000", 0)
	suite.Require().NotNil(err)
}

//tc：根据正确的区块hash和错误的交易index获取完整交易，获取失败
func (suite Model18) Test1840_GetTransactionByBlockHashAndIndexWithErrorIndexIsFail() {
	pk, from, err := repo.Node1Priv()
	suite.Require().Nil(err)
	_, to, err := repo.KeyPriv()
	suite.Require().Nil(err)
	too := common.HexToAddress(to.String())
	amount := big.NewInt(1000000000)
	gasLimit := uint64(21000)
	gasPrice, err := suite.client.EthGasPrice()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	nonce, err := client.GetPendingNonceByAccount(from.String())
	suite.Require().Nil(err)

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &too,
		Value:    amount,
		Gas:      gasLimit,
		GasPrice: &gasPrice,
		Data:     []byte{},
	})

	bytes, err := pk.Bytes()
	suite.Require().Nil(err)
	privateKey, err := crypto.HexToECDSA(common.Bytes2Hex(bytes))
	suite.Require().Nil(err)
	signTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(1356)), privateKey)
	suite.Require().Nil(err)
	data, err := signTx.MarshalBinary()
	rawTx := hexutil.Bytes(data).String()

	hash, err := suite.client.EthSendRawTransaction(rawTx)
	suite.Require().Nil(err)
	time.Sleep(1 * time.Second)
	res1, err := suite.client.EthGetTransactionByHash(hash)
	suite.Require().Nil(err)
	_, err = suite.client.EthGetTransactionByBlockHashAndIndex(res1.BlockHash, 2001)
	suite.Require().NotNil(err)
}

//tc：根据错误的区块hash和错误的交易index获取完整交易，获取失败
func (suite Model18) Test1841_GetTransactionByBlockHashAndIndexWithErrorHashAndErrorIndexIsFail() {
	_, err := suite.client.EthGetTransactionByBlockHashAndIndex("0x0000000000000000000000000000000000000000", 2001)
	suite.Require().NotNil(err)
}

//tc：根据正确的区块高度和正确的交易index获取完整交易，获取成功
func (suite Model18) Test1842_GetTransactionByBlockNumberAndIndexIsSuccess() {
	pk, from, err := repo.Node1Priv()
	suite.Require().Nil(err)
	_, to, err := repo.KeyPriv()
	suite.Require().Nil(err)
	too := common.HexToAddress(to.String())
	amount := big.NewInt(1000000000)
	gasLimit := uint64(21000)
	gasPrice, err := suite.client.EthGasPrice()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	nonce, err := client.GetPendingNonceByAccount(from.String())
	suite.Require().Nil(err)

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &too,
		Value:    amount,
		Gas:      gasLimit,
		GasPrice: &gasPrice,
		Data:     []byte{},
	})

	bytes, err := pk.Bytes()
	suite.Require().Nil(err)
	privateKey, err := crypto.HexToECDSA(common.Bytes2Hex(bytes))
	suite.Require().Nil(err)
	signTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(1356)), privateKey)
	suite.Require().Nil(err)
	data, err := signTx.MarshalBinary()
	rawTx := hexutil.Bytes(data).String()

	hash, err := suite.client.EthSendRawTransaction(rawTx)
	suite.Require().Nil(err)
	time.Sleep(1 * time.Second)
	res1, err := suite.client.EthGetTransactionByHash(hash)
	suite.Require().Nil(err)
	res2, err := suite.client.EthGetTransactionByBlockNumberAndIndex(*res1.BlockNumber, *res1.TransactionIndex)
	suite.Require().Nil(err)
	suite.Require().Equal(res1, res2)
}

//tc：根据错误的区块高度和正确的交易index获取完整交易，获取失败
func (suite Model18) Test1843_GetTransactionByBlockNumberAndIndexWithErrorNumberIsFail() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	meta, err := client.GetChainMeta()
	suite.Require().Nil(err)
	_, err = suite.client.EthGetTransactionByBlockNumberAndIndex(int(meta.Height)+1, 0)
	suite.Require().NotNil(err)
}

//tc：根据正确的区块高度和错误的交易index获取完整交易，获取失败
func (suite Model18) Test1844_GetTransactionByBlockNumberAndIndexWithErrorIndexIsFail() {
	pk, from, err := repo.Node1Priv()
	suite.Require().Nil(err)
	_, to, err := repo.KeyPriv()
	suite.Require().Nil(err)
	too := common.HexToAddress(to.String())
	amount := big.NewInt(1000000000)
	gasLimit := uint64(21000)
	gasPrice, err := suite.client.EthGasPrice()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	nonce, err := client.GetPendingNonceByAccount(from.String())
	suite.Require().Nil(err)

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &too,
		Value:    amount,
		Gas:      gasLimit,
		GasPrice: &gasPrice,
		Data:     []byte{},
	})

	bytes, err := pk.Bytes()
	suite.Require().Nil(err)
	privateKey, err := crypto.HexToECDSA(common.Bytes2Hex(bytes))
	suite.Require().Nil(err)
	signTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(1356)), privateKey)
	suite.Require().Nil(err)
	data, err := signTx.MarshalBinary()
	rawTx := hexutil.Bytes(data).String()

	hash, err := suite.client.EthSendRawTransaction(rawTx)
	suite.Require().Nil(err)
	time.Sleep(1 * time.Second)
	res1, err := suite.client.EthGetTransactionByHash(hash)
	suite.Require().Nil(err)
	_, err = suite.client.EthGetTransactionByBlockNumberAndIndex(*res1.BlockNumber, 2001)
	suite.Require().NotNil(err)
}

//tc：根据错误的区块高度和错误的交易index获取完整交易，获取失败
func (suite Model18) Test1845_GetTransactionByBlockNumberAndIndexWithErrorNumberAndIndexIsFail() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	meta, err := client.GetChainMeta()
	suite.Require().Nil(err)
	_, err = suite.client.EthGetTransactionByBlockNumberAndIndex(int(meta.Height)+1, 2001)
	suite.Require().NotNil(err)
}

//tc：根据正确的交易hash获取交易回执，获取成功
func (suite Model18) Test1846_GetTransactionReceiptIsSuccess() {
	pk, from, err := repo.Node1Priv()
	suite.Require().Nil(err)
	_, to, err := repo.KeyPriv()
	suite.Require().Nil(err)
	too := common.HexToAddress(to.String())
	amount := big.NewInt(1000000000)
	gasLimit := uint64(21000)
	gasPrice, err := suite.client.EthGasPrice()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	nonce, err := client.GetPendingNonceByAccount(from.String())
	suite.Require().Nil(err)

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &too,
		Value:    amount,
		Gas:      gasLimit,
		GasPrice: &gasPrice,
		Data:     []byte{},
	})

	bytes, err := pk.Bytes()
	suite.Require().Nil(err)
	privateKey, err := crypto.HexToECDSA(common.Bytes2Hex(bytes))
	suite.Require().Nil(err)
	signTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(1356)), privateKey)
	suite.Require().Nil(err)
	data, err := signTx.MarshalBinary()
	rawTx := hexutil.Bytes(data).String()

	hash, err := suite.client.EthSendRawTransaction(rawTx)
	suite.Require().Nil(err)
	time.Sleep(time.Second * 2)
	res, err := suite.client.EthGetTransactionReceipt(hash)
	suite.Require().Nil(err)
	suite.Require().NotNil(res)
}

//tc：根据错误的交易hash获取交易回执，获取失败
func (suite Model18) Test1847_GetTransactionReceiptWithErrorHashIsFail() {
	_, err := suite.client.EthGetTransactionReceipt("0x0000000000000000000000000000000000000000000000000000000000000000")
	suite.Require().NotNil(err)
}
