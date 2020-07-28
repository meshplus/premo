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

func (suite *Snake) TestGetBlockByHeight() {
	// first block
	block, err := suite.client.GetBlock("1", pb.GetBlockRequest_HEIGHT)
	suite.Nil(err)
	suite.Equal(block.BlockHeader.Number, uint64(1))

	// current block
	chainMeta, err := suite.client.GetChainMeta()
	suite.Nil(err)

	block, err = suite.client.GetBlock(strconv.Itoa(int(chainMeta.Height)), pb.GetBlockRequest_HEIGHT)
	suite.Nil(err)
	suite.Equal(chainMeta.Height, block.BlockHeader.Number)
}

func (suite *Snake) TestGetBlockByNonexistentHeight() {
	// get current block height
	chainMeta, err := suite.client.GetChainMeta()
	suite.Nil(err)

	_, err = suite.client.GetBlock(strconv.Itoa(int(chainMeta.Height+1)), pb.GetBlockRequest_HEIGHT)
	suite.NotNil(err)
	suite.Contains(err.Error(), "not found in DB")

	_, err = suite.client.GetBlock("0", pb.GetBlockRequest_HEIGHT)
	suite.NotNil(err)
	suite.Contains(err.Error(), "not found in DB")
}

func (suite *Snake) TestGetBlockByWrongHeight() {
	_, err := suite.client.GetBlock("a", pb.GetBlockRequest_HEIGHT)
	suite.NotNil(err)
	suite.Contains(err.Error(), "wrong block number")
}

func (suite *Snake) TestGetBlockByParentHeight() {
	// get current block height
	chainMeta, err := suite.client.GetChainMeta()
	suite.Nil(err)

	// parent height
	h := int(chainMeta.Height - 1)

	block, err := suite.client.GetBlock(strconv.Itoa(h), pb.GetBlockRequest_HEIGHT)
	suite.Nil(err)
	suite.Equal(uint64(h), block.BlockHeader.Number)
}

func (suite *Snake) TestGetBlockByHash() {
	// get current chain meta
	chainMeta, err := suite.client.GetChainMeta()
	suite.Nil(err)

	block, err := suite.client.GetBlock(chainMeta.BlockHash.String(), pb.GetBlockRequest_HASH)
	suite.Nil(err)
	suite.Equal(chainMeta.BlockHash, block.BlockHash)
}

func (suite *Snake) TestGetBlockByWrongHash() {
	_, err := suite.client.GetBlock(" ", pb.GetBlockRequest_HASH)
	suite.NotNil(err)
	suite.Contains(err.Error(), "not found in DB")
}

func (suite *Snake) TestGetBlockByParentHash() {
	// get current chain meta
	chainMeta, err := suite.client.GetChainMeta()
	suite.Nil(err)

	// get parenrt block
	parentBlock, err := suite.client.GetBlock(strconv.Itoa(int(chainMeta.Height-1)), pb.GetBlockRequest_HEIGHT)
	suite.Nil(err)

	block, err := suite.client.GetBlock(parentBlock.BlockHash.String(), pb.GetBlockRequest_HASH)
	suite.Nil(err)
	suite.Equal(chainMeta.Height-1, block.BlockHeader.Number)
}

func (suite *Snake) TestGetValidators() {
	_, err := suite.client.GetValidators()
	suite.Nil(err)
}

func (suite *Snake) TestGetBlockHeader() {
	ctx, cancel := context.WithTimeout(context.Background(), GetInfoTimeout)
	defer cancel()

	ch := make(chan *pb.BlockHeader)

	err := suite.client.GetBlockHeader(ctx, 1, 1, ch)
	suite.Nil(err)

	head := <-ch
	suite.Equal(uint64(1), head.Number)

	// get current chain meta
	chainMeta, err := suite.client.GetChainMeta()
	suite.Nil(err)

	ch2 := make(chan *pb.BlockHeader)

	err = suite.client.GetBlockHeader(ctx, chainMeta.Height, chainMeta.Height, ch2)
	suite.Nil(err)

	head = <-ch2
	suite.Equal(chainMeta.Height, head.Number)
}

func (suite *Snake) TestGetNonexistentBlockHeader() {
	// get current chain meta
	chainMeta, err := suite.client.GetChainMeta()
	suite.Nil(err)

	ctx, cancel := context.WithTimeout(context.Background(), GetInfoTimeout)
	defer cancel()

	ch := make(chan *pb.BlockHeader)

	err = suite.client.GetBlockHeader(ctx, chainMeta.Height+1, chainMeta.Height+1, ch)
	suite.Nil(err)

	_, ok := <-ch
	suite.Equal(false, ok)
}

func (suite *Snake) TestGetChainMeta() {
	chainMeta, err := suite.client.GetChainMeta()
	suite.Nil(err)

	_, err = suite.client.GetBlock(strconv.Itoa(int(chainMeta.Height+1)), pb.GetBlockRequest_HEIGHT)
	suite.NotNil(err)
	suite.Contains(err.Error(), "not found in DB")
}

func (suite *Snake) TestGetBlocks() {
	chainMeta, err := suite.client.GetChainMeta()
	suite.Nil(err)

	start := uint64(1)
	if chainMeta.Height > 10 {
		start = chainMeta.Height - 10
	}
	res, err := suite.client.GetBlocks(start, chainMeta.Height)
	suite.Nil(err)

	block, err := suite.client.GetBlock(strconv.Itoa(int(start)), pb.GetBlockRequest_HEIGHT)
	suite.Nil(err)

	suite.Equal(block.BlockHash, res.Blocks[0].BlockHash)
	suite.Equal(int(chainMeta.Height-start)+1, len(res.Blocks))
}

func (suite *Snake) TestGetBlocksByNonexistentRange() {
	chainMeta, err := suite.client.GetChainMeta()
	suite.Nil(err)

	res, err := suite.client.GetBlocks(chainMeta.Height+1, chainMeta.Height+1)
	suite.Nil(err)

	suite.Equal(0, len(res.Blocks))
}

func (suite *Snake) TestGetAccountBalance() {
	res, err := suite.client.GetAccountBalance(suite.from.String())
	suite.Nil(err)

	data := Account{}
	err = json.Unmarshal(res.Data, &data)
	suite.Nil(err)
	suite.NotEqual(0, data.Balance)
}

func (suite *Snake) TestGetAccountBalanceByNilAddress() {
	res, err := suite.client.GetAccountBalance("0x0000000000000000000000000000000000000000")
	suite.Nil(err)

	data := Account{}
	err = json.Unmarshal(res.Data, &data)
	suite.Nil(err)
	suite.Equal(uint64(0), data.Balance)
}

func (suite *Snake) TestGetAccountBalanceByWrongAddress() {
	_, err := suite.client.GetAccountBalance("ABC")
	suite.NotNil(err)

	_, err = suite.client.GetAccountBalance("0x123")
	suite.NotNil(err)

	_, err = suite.client.GetAccountBalance("__ _~~+——*/")
	suite.NotNil(err)
}

func (suite *Snake) TestGetChainStatus() {
	res, err := suite.client.GetChainStatus()
	suite.Nil(err)
	suite.Equal("normal", string(res.Data))
}

func (suite *Snake) TestGetNetworkMeta() {
	_, err := suite.client.GetNetworkMeta()
	suite.Nil(err)
}
