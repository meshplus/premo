package bxh_tester

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/meshplus/bitxhub-kit/types"
	"github.com/meshplus/bitxhub-model/pb"
)

const (
	GetInfoTimeout = 2 * time.Second
)

type Account struct {
	Type          string     `json:"type"`
	Balance       uint64     `json:"balance"`
	ContractCount uint64     `json:"contract_count"`
	CodeHash      types.Hash `json:"code_hash"`
}

type Model1 struct {
	*Snake
}

//tc: 根据区块高度查询区块，返回正确的区块信息
func (suite *Model1) Test0101_GetBlockByHeight() {
	// first block
	block, err := suite.client.GetBlock("1", pb.GetBlockRequest_HEIGHT)
	suite.Require().Nil(err)
	suite.Require().Equal(block.BlockHeader.Number, uint64(1))

	// current block
	chainMeta, err := suite.client.GetChainMeta()
	suite.Require().Nil(err)

	block, err = suite.client.GetBlock(strconv.Itoa(int(chainMeta.Height)), pb.GetBlockRequest_HEIGHT)
	suite.Require().Nil(err)
	suite.Require().Equal(chainMeta.Height, block.BlockHeader.Number)
}

//tc:根据不存在的区块高度查询区块，返回错误信息
func (suite *Model1) Test0102_GetBlockByNonexistentHeight() {
	// get current block height
	chainMeta, err := suite.client.GetChainMeta()
	suite.Require().Nil(err)

	_, err = suite.client.GetBlock(strconv.Itoa(int(chainMeta.Height+1)), pb.GetBlockRequest_HEIGHT)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "out of bounds")

	_, err = suite.client.GetBlock("0", pb.GetBlockRequest_HEIGHT)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "out of bounds")

	_, err = suite.client.GetBlock("-1", pb.GetBlockRequest_HEIGHT)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "wrong block number")
}

//tc:根据非法的区块高度查询区块，返回错误信息
func (suite *Model1) Test0103_GetBlockByWrongHeight() {
	_, err := suite.client.GetBlock("a", pb.GetBlockRequest_HEIGHT)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "wrong block number")

	_, err = suite.client.GetBlock("!我@#", pb.GetBlockRequest_HEIGHT)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "wrong block number")
}

//tc: 根据当前区块的父区块高度查询区块，返回正确的区块信息
func (suite *Model1) Test0104_GetBlockByParentHeight() {
	// get current block height
	chainMeta, err := suite.client.GetChainMeta()
	suite.Require().Nil(err)

	// parent height
	h := int(chainMeta.Height - 1)

	block, err := suite.client.GetBlock(strconv.Itoa(h), pb.GetBlockRequest_HEIGHT)
	suite.Require().Nil(err)
	suite.Require().Equal(uint64(h), block.BlockHeader.Number)
}

//tc:根据区块哈希查询区块，返回正确的区块信息
func (suite *Model1) Test0105_GetBlockByHash() {
	// get current chain meta
	chainMeta, err := suite.client.GetChainMeta()
	suite.Require().Nil(err)

	block, err := suite.client.GetBlock(chainMeta.BlockHash.String(), pb.GetBlockRequest_HASH)
	suite.Require().Nil(err)
	suite.Require().Equal(chainMeta.BlockHash.String(), block.BlockHash.String())
}

//tc:根据错误的区块哈希查询区块，返回错误信息
func (suite *Model1) Test0106_GetBlockByWrongHash() {
	_, err := suite.client.GetBlock(" ", pb.GetBlockRequest_HASH)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "invalid format of block hash for querying block")

	_, err = suite.client.GetBlock("0x0000000000000000000000000000000012345678900000000000000000000000", pb.GetBlockRequest_HASH)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "not found in DB")

}

//tc:根据当前区块的父区块哈希查询区块，返回正确的区块信息
func (suite *Model1) Test0107_GetBlockByParentHash() {
	// get current chain meta
	chainMeta, err := suite.client.GetChainMeta()
	suite.Require().Nil(err)

	// get parenrt block
	parentBlock, err := suite.client.GetBlock(strconv.Itoa(int(chainMeta.Height-1)), pb.GetBlockRequest_HEIGHT)
	suite.Require().Nil(err)

	block, err := suite.client.GetBlock(parentBlock.BlockHash.String(), pb.GetBlockRequest_HASH)
	suite.Require().Nil(err)
	suite.Require().Equal(chainMeta.Height-1, block.BlockHeader.Number)
}

//tc:查询链的validators，返回中继链的validator信息
func (suite *Model1) Test0108_GetValidators() {
	Validator, err := suite.client.GetValidators()
	suite.Require().Nil(err)
	suite.Require().NotNil(Validator)
}

//tc:根据指定范围查询区块头，返回正确范围内的区块头信息
func (suite *Model1) Test0109_GetBlockHeader() {
	ctx, cancel := context.WithTimeout(context.Background(), GetInfoTimeout)
	defer cancel()

	ch := make(chan *pb.BlockHeader)

	err := suite.client.GetBlockHeader(ctx, 1, 1, ch)
	suite.Require().Nil(err)

	head := <-ch
	suite.Require().Equal(uint64(1), head.Number)

	// get current chain meta
	chainMeta, err := suite.client.GetChainMeta()
	suite.Require().Nil(err)

	ch2 := make(chan *pb.BlockHeader)

	err = suite.client.GetBlockHeader(ctx, chainMeta.Height, chainMeta.Height, ch2)
	suite.Require().Nil(err)

	head = <-ch2
	suite.Require().Equal(chainMeta.Height, head.Number)
}

//tc:根据不存在的范围查询区块头，返回区块头为空
func (suite *Model1) Test0110_GetNonexistentBlockHeader() {
	// get current chain meta
	chainMeta, err := suite.client.GetChainMeta()
	suite.Require().Nil(err)

	ctx, cancel := context.WithTimeout(context.Background(), GetInfoTimeout)
	defer cancel()

	ch := make(chan *pb.BlockHeader)

	err = suite.client.GetBlockHeader(ctx, chainMeta.Height+1, chainMeta.Height+1, ch)
	suite.Require().Nil(err)

	_, ok := <-ch
	suite.Require().Equal(false, ok)
}

//tc:查询链的元数据，返回当前链的chain_meta信息
func (suite *Model1) Test0111_GetChainMeta() {
	chainMeta, err := suite.client.GetChainMeta()
	suite.Require().Nil(err)
	suite.Require().True(chainMeta.Height > 0)

	_, err = suite.client.GetBlock(strconv.Itoa(int(chainMeta.Height+1)), pb.GetBlockRequest_HEIGHT)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "out of bounds")
}

//tc:查询指定区块高度范围内的所有区块，返回正确范围区块信息
func (suite *Model1) Test0112_GetBlocks() {
	chainMeta, err := suite.client.GetChainMeta()
	suite.Require().Nil(err)

	start := uint64(1)
	if chainMeta.Height > 10 {
		start = chainMeta.Height - 10
	}
	res, err := suite.client.GetBlocks(start, chainMeta.Height)
	suite.Require().Nil(err)

	block, err := suite.client.GetBlock(strconv.Itoa(int(start)), pb.GetBlockRequest_HEIGHT)
	suite.Require().Nil(err)

	suite.Require().Equal(block.BlockHash, res.Blocks[0].BlockHash)
	suite.Require().Equal(int(chainMeta.Height-start)+1, len(res.Blocks))
}

//tc:查询不存在的高度范围的所有区块，返回区块信息为空
func (suite *Model1) Test0113_GetBlocksByNonexistentRange() {
	chainMeta, err := suite.client.GetChainMeta()
	suite.Require().Nil(err)

	res, err := suite.client.GetBlocks(chainMeta.Height+1, chainMeta.Height+1)
	suite.Require().Nil(err)

	suite.Require().Equal(0, len(res.Blocks))
}

//tc:根据指定地址查询余额，返回正确余额信息
func (suite *Model1) Test0114_GetAccountBalance() {
	res, err := suite.client.GetAccountBalance(suite.from.String())
	suite.Require().Nil(err)

	data := Account{}
	err = json.Unmarshal(res.Data, &data)
	suite.Require().Nil(err)
	suite.Require().NotEqual(0, data.Balance)
}

//tc：根据空的地址查询余额，返回余额为0
func (suite *Model1) Test0115_GetAccountBalanceByNilAddress() {
	res, err := suite.client.GetAccountBalance("0x0000000000000000000000000000000000000000")
	suite.Require().Nil(err)

	data := Account{}
	err = json.Unmarshal(res.Data, &data)
	suite.Require().Nil(err)
	suite.Require().Equal(uint64(0), data.Balance)
}

func (suite *Model1) Test0116_GetAccountBalanceByWrongAddress() {
	_, err := suite.client.GetAccountBalance("ABC")
	suite.Require().NotNil(err)

	_, err = suite.client.GetAccountBalance("0x123")
	suite.Require().NotNil(err)

	_, err = suite.client.GetAccountBalance("__ _~~+——*/")
	suite.Require().NotNil(err)
}

//tc:查询链的共识状态，返回正确的状态信息
func (suite *Model1) Test0117_GetChainStatus() {
	res, err := suite.client.GetChainStatus()
	suite.Require().Nil(err)
	suite.Require().Equal("normal", string(res.Data))
}

//tc:查询链的网络状态，返回正确的状态信息
func (suite *Model1) Test0118_GetNetworkMeta() {
	networkInfo, err := suite.client.GetNetworkMeta()
	suite.Require().Nil(err)
	suite.Require().NotNil(networkInfo)
}
