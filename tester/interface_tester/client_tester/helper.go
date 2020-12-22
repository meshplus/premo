package interface_tester

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/tidwall/gjson"
)

const (
	empty = ""
	tab   = "  "
)

type Meta struct {
	Height            string `json:"height"`
	BlockHash         string `json:"block_hash"`
	InterchainTxCount string `json:"interchain_tx_count"`
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

func parseResponse(data []byte) (string, error) {
	res := gjson.Get(string(data), "data")

	ret, err := base64.StdEncoding.DecodeString(res.String())
	if err != nil {
		return "", fmt.Errorf("wrong data: %w", err)
	}

	return string(ret), nil
}

func prettyJson(data string) (string, error) {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(data), empty, tab)
	if err != nil {
		return "", fmt.Errorf("wrong data: %w", err)
	}
	return out.String(), nil
}
