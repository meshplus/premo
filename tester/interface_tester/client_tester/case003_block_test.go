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

func (suite *Snake) TestGetBlockByHeightIsTrue() {
	height, err := getHeight()
	suite.Require().Nil(err)

	url, err := getURL(fmt.Sprintf("block?type=0&value=%s", height))
	suite.Require().Nil(err)

	data, err := httpGet(url)
	suite.Require().Nil(err)
	suite.Require().NotNil(data)
	suite.Require().NotContains(string(data), "error")
}

func (suite *Snake) TestGetBlockByHeightIsFalse() {
	height, err := getHeight()
	suite.Require().Nil(err)

	wrongHeight, err := strconv.Atoi(height)

	url, err := getURL(fmt.Sprintf("block?type=0&value=%d", wrongHeight+1))
	suite.Require().Nil(err)

	data, err := httpGet(url)
	suite.Require().Nil(err)

	var result Result
	err = json.Unmarshal(data, &result)
	suite.Require().Nil(err)
	suite.Require().Equal("out of bounds", result.Error)
	suite.Require().Equal(2, result.Code)
}

func (suite *Snake) TestGetBlockByHashIsTrue() {
	hash, err := getBlockHash()
	suite.Require().Nil(err)

	url, err := getURL(fmt.Sprintf("block?type=1&value=%s", hash))
	suite.Require().Nil(err)

	data, err := httpGet(url)
	suite.Require().Nil(err)
	suite.Require().NotNil(data)
	suite.Require().NotContains(string(data), "error")
}

func (suite *Snake) TestGetBlockByHashIsFalse() {
	hash, err := getBlockHash()
	suite.Require().Nil(err)

	url, err := getURL(fmt.Sprintf("block?type=1&value=%s", hash+"123"))
	suite.Require().Nil(err)

	data, err := httpGet(url)
	suite.Require().Nil(err)

	var result Result
	err = json.Unmarshal(data, &result)
	suite.Require().Nil(err)
	suite.Require().Equal("invalid format of block hash for querying block", result.Error)
	suite.Require().Equal(2, result.Code)
	suite.Require().Equal("invalid format of block hash for querying block", result.Message)
}
