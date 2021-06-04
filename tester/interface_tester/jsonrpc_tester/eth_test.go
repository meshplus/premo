package jsonrpc_tester

import (
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"strconv"
	"time"

	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/pkg/errors"

	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/premo/internal/repo"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
)

func (suite *Client) TestGetProtocolVersion() {
	var result hexutil.Uint
	err := suite.rpcClient.Call(&result, "eth_protocolVersion")
	suite.Require().Nil(err)
	suite.Require().Equal(hexutil.Uint(0x41), result)
	fmt.Println(result)
}

func (suite *Client) TestGetChainId() {
	var result hexutil.Uint
	err := suite.rpcClient.Call(&result, "eth_chainId")
	suite.Require().Nil(err)
	suite.Require().Equal(hexutil.Uint(1356), result)
}

func (suite *Client) TestMining() {
	// Always false
	var result bool
	err := suite.rpcClient.Call(&result, "eth_mining")
	suite.Require().Nil(err)
	suite.Require().False(result)
}

func (suite *Client) TestGetHashrate() {
	// Always zero
	var result hexutil.Uint
	err := suite.rpcClient.Call(&result, "eth_hashrate")
	suite.Require().Nil(err)
	suite.Require().Equal(hexutil.Uint(0), result)
	fmt.Println(result)
}

func (suite *Client) TestGetGasPrice() {
	var result hexutil.Big
	err := suite.rpcClient.Call(&result, "eth_gasPrice")
	suite.Require().Nil(err)
	suite.Require().Equal(int64(0), result.ToInt().Int64())
}

func (suite *Client) TestGetBlockNumber() {
	var result hexutil.Uint64
	err := suite.rpcClient.Call(&result, "eth_blockNumber")
	suite.Require().Nil(err)
	suite.Require().NotNil(result)
}

func (suite *Client) TestGetBalance() {
	address := "0xc7F999b83Af6DF9e67d0a37Ee7e900bF38b3D013"
	balance1, err := suite.GetBalance(address)
	suite.Require().Nil(err)

	seed := time.Now().Unix()
	rand.Seed(seed)

	num := rand.Int63n(10)
	err = suite.SendTxSuccess(int(num))
	suite.Require().Nil(err)
	balance2, err := suite.GetBalance(address)
	suite.Require().Nil(err)
	suite.Require().Equal(big.NewInt(0).Add(balance2.ToInt(), big.NewInt(num)), balance1.ToInt())

	num = rand.Int63n(10)
	err = suite.SendTxSuccess(int(num))
	suite.Require().Nil(err)
	balance3, err := suite.GetBalance(address)
	suite.Require().Nil(err)
	suite.Require().Equal(big.NewInt(0).Add(balance3.ToInt(), big.NewInt(num)), balance2.ToInt())

	num = rand.Int63n(10)
	err = suite.SendTxFail(int(num))
	suite.Require().Nil(err)
	balance4, err := suite.GetBalance(address)
	suite.Require().Nil(err)
	suite.Require().Equal(balance4.ToInt(), balance3.ToInt())
}

func (suite *Client) TestGetStorageAt() {
	blockNum1, err := suite.SetStore("k1", "v1")
	suite.Require().Nil(err)
	value, err := suite.GetStore("k1", blockNum1)
	suite.Require().Nil(err)
	suite.Require().Equal("v1", value)

	blockNum2, err := suite.SetStore("k2", "v2")
	suite.Require().Nil(err)
	value, err = suite.GetStore("k2", blockNum2)
	suite.Require().Nil(err)
	suite.Require().Equal("v2", value)

	blockNum3, err := suite.SetStore("k1", "v3")
	suite.Require().Nil(err)
	value, err = suite.GetStore("k1", blockNum3)
	suite.Require().Nil(err)
	suite.Require().Equal("v3", value)

	value, err = suite.GetStore("k1", blockNum1)
	suite.Require().Nil(err)
	suite.Require().Equal("v1", value)

}

func (suite *Client) TestGetTransactionCount() {
	var result hexutil.Uint64
	addrStr := "0xc7F999b83Af6DF9e67d0a37Ee7e900bF38b3D013"
	err := suite.rpcClient.Call(&result, "eth_getTransactionCount", addrStr, "0x12")
	suite.Require().Nil(err)
	suite.Require().NotNil(result)
}

//func (suite *Client) TestSendRawTransaction() {
//	var result common.Hash
//	tx := &types.EthTransaction{}
//	data, err := tx.MarshalBinary()
//	suite.Require().Nil(err)
//	err = suite.rpcClient.Call(&result, "eth_sendRawTransaction", data)
//	suite.Require().Nil(err)
//	suite.Require().NotNil(result)
//	fmt.Println(result)
//}

func (suite *Client) TestGetBlockTransactionCountByNumber() {
	var result hexutil.Uint
	err := suite.rpcClient.Call(&result, "eth_getBlockTransactionCountByNumber", 1)
	suite.Require().Nil(err)
	suite.Require().NotNil(result)
	fmt.Println(result)
}

func (suite *Client) TestGetBlockTransactionCountByHash() {
	var result1 hexutil.Uint64
	err := suite.rpcClient.Call(&result1, "eth_blockNumber")
	suite.Require().Nil(err)
	suite.Require().NotNil(result1)
	var result2 hexutil.Uint64
	height, err := hexutil.DecodeUint64(result1.String())
	suite.Require().Nil(err)
	block, err := suite.client.GetBlock(strconv.FormatUint(height, 10), pb.GetBlockRequest_HEIGHT)
	suite.Require().Nil(err)
	err = suite.rpcClient.Call(&result2, "eth_getBlockTransactionCountByHash", block.BlockHash)
	suite.Require().Nil(err)
	suite.Require().Equal(len(block.Transactions.Transactions), int(result2))
}

func (suite *Client) SendTxSuccess(num int) error {
	keyPath, err := repo.Node1Path()
	if err != nil {
		return err
	}
	pk, err := asym.RestorePrivateKey(keyPath, repo.KeyPassword)
	if err != nil {
		return err
	}
	from, err := pk.PublicKey().Address()
	if err != nil {
		return err
	}
	pk1, err := asym.GenerateKeyPair(crypto.Secp256k1)
	if err != nil {
		return err
	}
	to, err := pk1.PublicKey().Address()
	if err != nil {
		return err
	}
	for i := 0; i < num; i++ {
		data := &pb.TransactionData{
			Amount: (*pb.BigInt)(new(big.Int).SetInt64(1)),
		}
		payload, err := data.Marshal()
		if err != nil {
			return err
		}

		tx := &pb.BxhTransaction{
			From:      from,
			To:        to,
			Timestamp: time.Now().UnixNano(),
			Payload:   payload,
		}

		res, err := suite.client.SendTransactionWithReceipt(tx, nil)
		if err != nil {
			return err
		}
		if res.Status != pb.Receipt_SUCCESS {
			return errors.New(string(res.Ret))
		}
	}
	return nil
}

func (suite *Client) SendTxFail(num int) error {
	keyPath, err := repo.Node1Path()
	if err != nil {
		return err
	}
	pk, err := asym.RestorePrivateKey(keyPath, repo.KeyPassword)
	if err != nil {
		return err
	}
	from, err := pk.PublicKey().Address()
	if err != nil {
		return err
	}
	for i := 0; i < num; i++ {
		data := &pb.TransactionData{
			Amount: (*pb.BigInt)(new(big.Int).SetInt64(1)),
		}
		payload, err := data.Marshal()
		if err != nil {
			return err
		}

		tx := &pb.BxhTransaction{
			From:      from,
			Timestamp: time.Now().UnixNano(),
			Payload:   payload,
		}

		_, _ = suite.client.SendTransactionWithReceipt(tx, nil)
	}
	return nil
}

func (suite Client) GetBalance(address string) (hexutil.Big, error) {
	var blockNum hexutil.Uint64
	err := suite.rpcClient.Call(&blockNum, "eth_blockNumber")
	suite.Require().Nil(err)
	suite.Require().NotNil(blockNum)
	var balance hexutil.Big
	err = suite.rpcClient.Call(&balance, "eth_getBalance", address, blockNum)
	if err != nil {
		return hexutil.Big{}, err
	}
	return balance, nil
}

func (suite Client) SetStore(key, value string) (hexutil.Uint64, error) {
	res, err := suite.client.InvokeBVMContract(constant.StoreContractAddr.Address(), "Set", nil, pb.String(key), pb.String(value))
	if err != nil {
		return hexutil.Uint64(0), err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return hexutil.Uint64(0), errors.New(string(res.Ret))
	}
	var result hexutil.Uint64
	err = suite.rpcClient.Call(&result, "eth_blockNumber")
	if err != nil {
		return hexutil.Uint64(0), err
	}
	return result, nil
}

func (suite Client) GetStore(key string, blockNum hexutil.Uint64) (string, error) {
	var result hexutil.Bytes
	err := suite.rpcClient.Call(&result, "eth_getStorageAt", constant.StoreContractAddr.Address(), key, blockNum)
	if err != nil {
		return "", err
	}
	var value string
	err = json.Unmarshal(result, &value)
	suite.Require().Nil(err)
	return value, nil
}
