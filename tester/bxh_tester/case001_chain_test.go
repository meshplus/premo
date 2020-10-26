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
	suite.Require().Nil(err)
	suite.Require().Equal(block.BlockHeader.Number, uint64(1))

	// current block
	chainMeta, err := suite.client.GetChainMeta()
	suite.Require().Nil(err)

	block, err = suite.client.GetBlock(strconv.Itoa(int(chainMeta.Height)), pb.GetBlockRequest_HEIGHT)
	suite.Require().Nil(err)
	suite.Require().Equal(chainMeta.Height, block.BlockHeader.Number)
}

func (suite *Snake) TestGetBlockByNonexistentHeight() {
	// get current block height
	chainMeta, err := suite.client.GetChainMeta()
	suite.Require().Nil(err)

	_, err = suite.client.GetBlock(strconv.Itoa(int(chainMeta.Height+1)), pb.GetBlockRequest_HEIGHT)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "not found in DB")

	_, err = suite.client.GetBlock("0", pb.GetBlockRequest_HEIGHT)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "not found in DB")
}

func (suite *Snake) TestGetBlockByWrongHeight() {
	_, err := suite.client.GetBlock("a", pb.GetBlockRequest_HEIGHT)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "wrong block number")
}

func (suite *Snake) TestGetBlockByParentHeight() {
	// get current block height
	chainMeta, err := suite.client.GetChainMeta()
	suite.Require().Nil(err)

	// parent height
	h := int(chainMeta.Height - 1)

	block, err := suite.client.GetBlock(strconv.Itoa(h), pb.GetBlockRequest_HEIGHT)
	suite.Require().Nil(err)
	suite.Require().Equal(uint64(h), block.BlockHeader.Number)
}

func (suite *Snake) TestGetBlockByHash() {
	// get current chain meta
	chainMeta, err := suite.client.GetChainMeta()
	suite.Require().Nil(err)

	block, err := suite.client.GetBlock(chainMeta.BlockHash.String(), pb.GetBlockRequest_HASH)
	suite.Require().Nil(err)
	suite.Require().Equal(chainMeta.BlockHash.String(), block.BlockHash.String())
}

func (suite *Snake) TestGetBlockByWrongHash() {
	_, err := suite.client.GetBlock(" ", pb.GetBlockRequest_HASH)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "invalid format of block hash for querying block")
}

func (suite *Snake) TestGetBlockByParentHash() {
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

func (suite *Snake) TestGetValidators() {
	_, err := suite.client.GetValidators()
	suite.Require().Nil(err)
}

func (suite *Snake) TestGetBlockHeader() {
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

func (suite *Snake) TestGetNonexistentBlockHeader() {
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

func (suite *Snake) TestGetChainMeta() {
	chainMeta, err := suite.client.GetChainMeta()
	suite.Require().Nil(err)

	_, err = suite.client.GetBlock(strconv.Itoa(int(chainMeta.Height+1)), pb.GetBlockRequest_HEIGHT)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "not found in DB")
}

func (suite *Snake) TestGetBlocks() {
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

func (suite *Snake) TestGetBlocksByNonexistentRange() {
	chainMeta, err := suite.client.GetChainMeta()
	suite.Require().Nil(err)

	res, err := suite.client.GetBlocks(chainMeta.Height+1, chainMeta.Height+1)
	suite.Require().Nil(err)

	suite.Require().Equal(0, len(res.Blocks))
}

func (suite *Snake) TestGetAccountBalance() {
	res, err := suite.client.GetAccountBalance(suite.from.String())
	suite.Require().Nil(err)

	data := Account{}
	err = json.Unmarshal(res.Data, &data)
	suite.Require().Nil(err)
	suite.Require().NotEqual(0, data.Balance)
}

func (suite *Snake) TestGetAccountBalanceByNilAddress() {
	res, err := suite.client.GetAccountBalance("0x0000000000000000000000000000000000000000")
	suite.Require().Nil(err)

	data := Account{}
	err = json.Unmarshal(res.Data, &data)
	suite.Require().Nil(err)
	suite.Require().Equal(uint64(0), data.Balance)
}

func (suite *Snake) TestGetAccountBalanceByWrongAddress() {
	_, err := suite.client.GetAccountBalance("ABC")
	suite.Require().NotNil(err)

	_, err = suite.client.GetAccountBalance("0x123")
	suite.Require().NotNil(err)

	_, err = suite.client.GetAccountBalance("__ _~~+——*/")
	suite.Require().NotNil(err)
}

func (suite *Snake) TestGetChainStatus() {
	res, err := suite.client.GetChainStatus()
	suite.Require().Nil(err)
	suite.Require().Equal("normal", string(res.Data))
}

func (suite *Snake) TestGetNetworkMeta() {
	_, err := suite.client.GetNetworkMeta()
	suite.Require().Nil(err)
}
