package bxh_tester

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
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
	res, err := suite.client.Call("eth_chainId")
	suite.Require().Nil(err)
	suite.Require().Equal("0x54c", string(res))
}

//tc：获取当前节点是否在挖矿，获取成功
func (suite Model18) Test1803_GetMiningStatusIsSuccess() {
	status, err := suite.client.EthMining()
	suite.Require().Nil(err)
	suite.Require().Equal(false, status)
}

//tc：获取当前节点的算力，获取成功
func (suite Model18) Test1804_GetHashrateIsSuccess() {
	hashrate, err := suite.client.EthHashrate()
	suite.Require().Nil(err)
	suite.Require().Equal(0, hashrate)
}

//tc：获取当前gas的价格，获取成功
func (suite Model18) Test1805_GetGasPriceIsSuccess() {
	price, err := suite.client.EthGasPrice()
	suite.Require().Nil(err)
	suite.Require().Equal("50000", price.String())
}

//tc：获取gas上限的建议，获取成功
func (suite Model18) Test1806_GetMaxPriorityFeePerGasIsSuccess() {
	res, err := suite.client.Call("eth_maxPriorityFeePerGas")
	suite.Require().Nil(err)
	suite.Require().Equal("0x0", string(res))
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

//tc：根据正确的区块高度获取当前区块的交易数量，获取成功
func (suite Model18) Test1812_GetBlockTransactionCountByNumberIsSuccess() {
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

//tc：根据正确的合约地址获取合约，获取成功
func (suite Model18) Test1813_GetCodeIsSuccess() {
	address, err := suite.DeploySimpleRule()
	suite.Require().Nil(err)
	fmt.Println(address)
	code, err := suite.client.EthGetCode(address, "latest")
	suite.Require().Nil(err)
	fmt.Println(code)
}

//tc：根据正确的交易hash获取日志，获取成功
func (suite Model18) Test1814_GetTransactionLogsIsSuccess() {

}

//tc：发送正确的交易，交易发送成功
func (suite Model18) Test1815_SendRawTransactionIsSuccess() {
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
	rawTx := common.Bytes2Hex(data)

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

//tc：根据正确的参数调用内置合约，调用成功
func (suite Model18) Test1816_CallIsSuccess() {

}

//tc：根据正确的参数获取调用内置合约的gas limit，获取成功
func (suite Model18) Test1817_EstimateGasIsSuccess() {

}

//tc：根据正确的区块hash获取区块内全部的交易hash，获取成功
func (suite Model18) Test1818_GetBlockByHashIsSuccess() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	meta, err := client.GetChainMeta()
	suite.Require().Nil(err)
	txs, err := suite.client.EthGetBlockByHash(meta.BlockHash.String(), false)
	suite.Require().Nil(err)
	suite.Require().NotNil(txs)
}

//tc：根据正确的区块hash获取区块内完整的交易，获取成功
func (suite Model18) Test1819_GetBlockByHasFullIsSuccess() {
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
func (suite Model18) Test1820_GetBlockByHashFullWithErrorHashIsSuccess() {
	_, err := suite.client.EthGetBlockByHash("0x0000000000000000000000000000000000000000", true)
	suite.Require().NotNil(err)
}

//tc：根据正确的区块高度获取区块内全部的交易hash，获取成功
func (suite Model18) Test1821_GetBlockByNumberIsSuccess() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	meta, err := client.GetChainMeta()
	suite.Require().Nil(err)
	txs, err := suite.client.EthGetBlockByNumber(int(meta.Height), false)
	suite.Require().Nil(err)
	suite.Require().NotNil(txs)
}

//tc：根据正确的区块高度获取区块内完整的交易，获取成功
func (suite Model18) Test1822_GetBlockByNumberFullIsSuccess() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	meta, err := client.GetChainMeta()
	suite.Require().Nil(err)
	txs, err := suite.client.EthGetBlockByNumber(int(meta.Height), true)
	suite.Require().Nil(err)
	suite.Require().NotNil(txs)
}

//tc：根据正确的交易hash获取完整的交易，获取成功
func (suite Model18) Test1823_GetTransactionByHashIsSuccess() {
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
	rawTx := common.Bytes2Hex(data)

	hash, err := suite.client.EthSendRawTransaction(rawTx)
	suite.Require().Nil(err)
	res, err := suite.client.EthGetTransactionByHash(hash)
	suite.Require().Nil(err)
	suite.Require().NotNil(res)
}

//tc：根据错误的交易hash获取完整的交易，获取失败
func (suite Model18) Test1824_GetTransactionByHashWithErrorHashIsSuccess() {
	_, err := suite.client.EthGetTransactionByHash("0x0000000000000000000000000000000000000000000000000000000000000000")
	suite.Require().NotNil(err)
}

//tc：根据正确的区块hash和正确的交易index获取完整交易，获取成功
func (suite Model18) Test1825_GetTransactionByBlockHashAndIndexIsSuccess() {
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
	rawTx := common.Bytes2Hex(data)

	hash, err := suite.client.EthSendRawTransaction(rawTx)
	suite.Require().Nil(err)
	res1, err := suite.client.EthGetTransactionByHash(hash)
	suite.Require().Nil(err)
	res2, err := suite.client.EthGetTransactionByBlockHashAndIndex(res1.BlockHash, *res1.TransactionIndex)
	suite.Require().Nil(err)
	suite.Require().Equal(res1, res2)
}

//tc：根据正确的区块高度和正确的交易index获取完整交易，获取成功
func (suite Model18) Test1826_GetTransactionByBlockNumberAndIndexIsSuccess() {
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
	rawTx := common.Bytes2Hex(data)

	hash, err := suite.client.EthSendRawTransaction(rawTx)
	suite.Require().Nil(err)
	res1, err := suite.client.EthGetTransactionByHash(hash)
	suite.Require().Nil(err)
	res2, err := suite.client.EthGetTransactionByBlockNumberAndIndex(*res1.BlockNumber, *res1.TransactionIndex)
	suite.Require().Nil(err)
	suite.Require().Equal(res1, res2)
}

//tc：根据正确的交易hash获取交易回执，获取成功
func (suite Model18) Test1827_GetTransactionReceiptIsSuccess() {
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
	rawTx := common.Bytes2Hex(data)

	hash, err := suite.client.EthSendRawTransaction(rawTx)
	suite.Require().Nil(err)
	time.Sleep(time.Second * 2)
	res, err := suite.client.EthGetTransactionReceipt(hash)
	suite.Require().Nil(err)
	suite.Require().NotNil(res)
}
