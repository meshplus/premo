package bxh_tester

import (
	"context"
	"encoding/json"
	"math/big"
	"strconv"
	"time"

	"github.com/meshplus/premo/internal/repo"

	"github.com/meshplus/bitxhub-kit/types"
	"github.com/meshplus/bitxhub-model/pb"
)

const (
	GetInfoTimeout = 2 * time.Second
)

type Account struct {
	Type          string     `json:"type"`
	Balance       big.Int    `json:"balance"`
	ContractCount uint64     `json:"contract_count"`
	CodeHash      types.Hash `json:"code_hash"`
}

type Model1 struct {
	*Snake
}

//tc：根据区块高度查询区块，返回正确的区块信息
func (suite *Model1) Test0101_GetBlockByHeightIsSuccess() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	// first block
	block, err := client.GetBlock("1", pb.GetBlockRequest_HEIGHT, false)
	suite.Require().Nil(err)
	suite.Require().Equal(block.BlockHeader.Number, uint64(1))
	// current block
	chainMeta, err := client.GetChainMeta()
	suite.Require().Nil(err)
	block, err = client.GetBlock(strconv.Itoa(int(chainMeta.Height)), pb.GetBlockRequest_HEIGHT, false)
	suite.Require().Nil(err)
	suite.Require().Equal(chainMeta.Height, block.BlockHeader.Number)
}

//tc：根据当前区块的父区块高度查询区块，返回正确的区块信息
func (suite *Model1) Test0102_GetBlockByParentHeightIsSuccess() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	// get current block height
	chainMeta, err := client.GetChainMeta()
	suite.Require().Nil(err)
	// parent height
	h := int(chainMeta.Height - 1)
	block, err := client.GetBlock(strconv.Itoa(h), pb.GetBlockRequest_HEIGHT, false)
	suite.Require().Nil(err)
	suite.Require().Equal(uint64(h), block.BlockHeader.Number)
}

//tc：根据区块哈希查询区块，返回正确的区块信息
func (suite *Model1) Test0103_GetBlockByHashIsSuccess() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	// get current chain meta
	chainMeta, err := client.GetChainMeta()
	suite.Require().Nil(err)
	block, err := client.GetBlock(chainMeta.BlockHash.String(), pb.GetBlockRequest_HASH, false)
	suite.Require().Nil(err)
	suite.Require().Equal(chainMeta.BlockHash.String(), block.BlockHash.String())
}

//tc：根据当前区块的父区块哈希查询区块，返回正确的区块信息
func (suite *Model1) Test0104_GetBlockByParentHashIsSuccess() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	// get current chain meta
	chainMeta, err := client.GetChainMeta()
	suite.Require().Nil(err)
	// get parent block
	parentBlock, err := client.GetBlock(strconv.Itoa(int(chainMeta.Height-1)), pb.GetBlockRequest_HEIGHT, false)
	suite.Require().Nil(err)
	block, err := client.GetBlock(parentBlock.BlockHash.String(), pb.GetBlockRequest_HASH, false)
	suite.Require().Nil(err)
	suite.Require().Equal(chainMeta.Height-1, block.BlockHeader.Number)
}

//tc：查询链的validators，返回中继链的validator信息
func (suite *Model1) Test0105_GetValidatorsIsSuccess() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	Validator, err := client.GetValidators()
	suite.Require().Nil(err)
	suite.Require().NotNil(Validator)
}

//tc：根据指定范围查询区块头，返回正确范围内的区块头信息
func (suite *Model1) Test0106_GetBlockHeaderIsSuccess() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	ctx, cancel := context.WithTimeout(context.Background(), GetInfoTimeout)
	defer cancel()
	ch := make(chan *pb.BlockHeader)
	err = client.GetBlockHeader(ctx, 1, 1, ch)
	suite.Require().Nil(err)
	head := <-ch
	suite.Require().Equal(uint64(1), head.Number)
	// get current chain meta
	chainMeta, err := client.GetChainMeta()
	suite.Require().Nil(err)
	ch2 := make(chan *pb.BlockHeader)
	err = client.GetBlockHeader(ctx, chainMeta.Height, chainMeta.Height, ch2)
	suite.Require().Nil(err)
	head = <-ch2
	suite.Require().Equal(chainMeta.Height, head.Number)
}

//tc：查询链的元数据，返回当前链的chain_meta信息
func (suite *Model1) Test0107_GetChainMetaIsSuccess() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	chainMeta, err := client.GetChainMeta()
	suite.Require().Nil(err)
	suite.Require().True(chainMeta.Height > 0)
	_, err = client.GetBlock(strconv.Itoa(int(chainMeta.Height+1)), pb.GetBlockRequest_HEIGHT, false)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "out of bounds")
}

//tc：查询指定区块高度范围内的所有区块，返回正确范围区块信息
func (suite *Model1) Test0108_GetBlocksIsSuccess() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	chainMeta, err := client.GetChainMeta()
	suite.Require().Nil(err)
	start := uint64(1)
	if chainMeta.Height > 10 {
		start = chainMeta.Height - 10
	}
	res, err := client.GetBlocks(start, chainMeta.Height, false)
	suite.Require().Nil(err)
	block, err := client.GetBlock(strconv.Itoa(int(start)), pb.GetBlockRequest_HEIGHT, false)
	suite.Require().Nil(err)
	suite.Require().Equal(block.BlockHash, res.Blocks[0].BlockHash)
	suite.Require().Equal(int(chainMeta.Height-start)+1, len(res.Blocks))
}

//tc：根据指定地址查询余额，返回正确余额信息
func (suite *Model1) Test0109_GetAccountBalanceIsSuccess() {
	pk, from, err := repo.Node1Priv()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.GetAccountBalance(from.String())
	suite.Require().Nil(err)
	data := Account{}
	err = json.Unmarshal(res.Data, &data)
	suite.Require().Nil(err)
	suite.Require().NotEqual("0", data.Balance.String())
}

//tc：查询链的共识状态，返回正确的状态信息
func (suite *Model1) Test0110_GetChainStatusIsSuccess() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.GetChainStatus()
	suite.Require().Nil(err)
	suite.Require().Equal("normal", string(res.Data))
}

//tc：查询链的网络状态，返回正确的状态信息
func (suite *Model1) Test0111_GetNetworkMetaIsSuccess() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	networkInfo, err := client.GetNetworkMeta()
	suite.Require().Nil(err)
	suite.Require().NotNil(networkInfo)
}
