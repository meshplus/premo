package interface_tester

import (
	"encoding/json"
	"fmt"
)

type Meta struct {
	Height            string `json:"height"`
	BlockHash         string `json:"block_hash"`
	InterchainTxCount string `json:"interchain_tx_count"`
}

func (suite *Snake) TestGetChainMeta() {
	url := getURL("chain_meta")

	data, err := httpGet(url)
	suite.Require().Nil(err)

	var meta Meta
	err = json.Unmarshal(data, &meta)
	suite.Require().Nil(err)
	suite.Require().NotNil(meta)
}

func (suite *Snake) TestGetChainStatus() {
	url := getURL("info?type=0")

	data, err := httpGet(url)
	suite.Require().Nil(err)

	ret, err := parseResponse(data)
	suite.Require().Nil(err)
	suite.Require().Equal("normal", ret)
}

func getMeta() (Meta, error) {
	var meta Meta
	url := getURL("chain_meta")

	data, err := httpGet(url)
	if err != nil {
		return meta, fmt.Errorf("get data error: %w", err)
	}

	err = json.Unmarshal(data, &meta)
	if err != nil {
		return meta, fmt.Errorf("json unmarshal error: %w", err)
	}
	return meta, nil
}

func getHeight() (string, error) {
	meta, err := getMeta()
	if err != nil {
		return "", fmt.Errorf("get meta error: %w", err)
	}
	return meta.Height, nil
}

func getBlockHash() (string, error) {
	meta, err := getMeta()
	if err != nil {
		return "", fmt.Errorf("get meta error: %w", err)
	}
	return meta.BlockHash, nil
}
