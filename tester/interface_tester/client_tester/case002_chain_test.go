package interface_tester

import (
	"encoding/json"
	"strconv"
)

func (suite *Snake) TestGetChainMeta() {
	url := getURL("chain_meta")

	data, err := httpGet(url)
	suite.Require().Nil(err)

	var meta Meta
	err = json.Unmarshal(data, &meta)
	suite.Require().Nil(err)

	height, err := strconv.Atoi(meta.Height)
	suite.Require().True(height > 0)
	suite.Require().NotNil(meta.BlockHash)
}

func (suite *Snake) TestGetChainStatus() {
	url := getURL("info?type=0")

	data, err := httpGet(url)
	suite.Require().Nil(err)

	ret, err := parseResponse(data)
	suite.Require().Nil(err)
	suite.Require().Equal("normal", ret)
}
