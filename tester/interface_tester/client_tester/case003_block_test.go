package interface_tester

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Result struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (suite *Snake) TestGetBlockByHeight() {
	height, err := getHeight()
	suite.Require().Nil(err)

	url := getURL(fmt.Sprintf("block?type=0&value=%s", height))

	data, err := httpGet(url)
	suite.Require().Nil(err)
	suite.Require().NotNil(data)
	suite.Require().NotContains(string(data), "error")
	suite.Require().Contains(string(data), "block_header")
}

func (suite *Snake) TestGetBlockByHeightWithHeightOutOfBounds() {
	height, err := getHeight()
	suite.Require().Nil(err)

	wrongHeight, err := strconv.Atoi(height)

	url := getURL(fmt.Sprintf("block?type=0&value=%d", wrongHeight+5))

	data, err := httpGet(url)
	suite.Require().Nil(err)

	var result Result
	err = json.Unmarshal(data, &result)
	suite.Require().Nil(err)
	suite.Require().Equal("out of bounds", result.Error)
	suite.Require().Equal(2, result.Code)
}

func (suite *Snake) TestGetBlockByHeightWithHeightIsNegative() {

	url := getURL(fmt.Sprintf("block?type=0&value=%d", -1))

	data, err := httpGet(url)
	suite.Require().Nil(err)

	var result Result
	err = json.Unmarshal(data, &result)
	suite.Require().Nil(err)
	suite.Require().Contains(result.Error, "wrong block number")
	suite.Require().Equal(2, result.Code)
}

func (suite *Snake) TestGetBlockByHeightWithHeightIsString() {

	url := getURL(fmt.Sprintf("block?type=0&value=%s", "!2#我"))

	data, err := httpGet(url)
	suite.Require().Nil(err)

	var result Result
	err = json.Unmarshal(data, &result)
	suite.Require().Nil(err)
	suite.Require().Contains(result.Error, "wrong block number")
	suite.Require().Equal(2, result.Code)
}

func (suite *Snake) TestGetBlockByHash() {
	hash, err := getBlockHash()
	suite.Require().Nil(err)

	url := getURL(fmt.Sprintf("block?type=1&value=%s", hash))

	data, err := httpGet(url)
	suite.Require().Nil(err)
	suite.Require().NotNil(data)
	suite.Require().NotContains(string(data), "error")
}

func (suite *Snake) TestGetBlockByHashWithInvalidFormat() {
	hash, err := getBlockHash()
	suite.Require().Nil(err)

	url := getURL(fmt.Sprintf("block?type=1&value=%s", hash+"123"))

	data, err := httpGet(url)
	suite.Require().Nil(err)

	var result Result
	err = json.Unmarshal(data, &result)
	suite.Require().Nil(err)
	suite.Require().Equal("invalid format of block hash for querying block", result.Error)
	suite.Require().Equal(2, result.Code)
}

func (suite *Snake) TestGetBlockByHashWithNonexistent() {
	wrongHash := "0x0000000000000000000000000000000012345678900000000000000000000000"

	url := getURL(fmt.Sprintf("block?type=1&value=%s", wrongHash))

	data, err := httpGet(url)
	suite.Require().Nil(err)

	var result Result
	err = json.Unmarshal(data, &result)
	suite.Require().Nil(err)
	suite.Require().Equal("not found in DB", result.Error)
	suite.Require().Equal(2, result.Code)
}
