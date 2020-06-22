package tester

import (
	"github.com/meshplus/premo/pkg/appchain/ethereum"
	"github.com/meshplus/premo/pkg/appchain/fabric"
)

type EthClientHelper struct {
	*ethereum.EthClient
	abiPath      string
	contractAddr string
	appchainId   string
}

type FabricClientHelper struct {
	*fabric.FabricClient
	appchainId string
}

func (c *EthClientHelper) GetBalance(username string) (string, error) {
	response, err := c.Invoke(c.abiPath, c.contractAddr, "getBalance", username)
	if err != nil {
		return "", err
	}
	return response, nil
}

func (c *EthClientHelper) InterchainTransfer(targetAppId, from, to, amount string) error {
	args := make([]string, 0)
	args = append(args, targetAppId, "mychannel&transfer", from, to, amount)
	_, err := c.Invoke(c.abiPath, c.contractAddr, "transfer", args...)
	if err != nil {
		return err
	}
	return nil
}

func (c *FabricClientHelper) GetBalance(username string) (string, error) {
	args := make([][]byte, 0)
	args = append(args, []byte(username))
	response, err := c.Invoke("transfer", "getBalance", args...)
	if err != nil {
		return "", err
	}
	return response, nil
}

func (c *FabricClientHelper) InterchainTransfer(targetAppId, contractAddr, from, to, amount string) error {
	args := make([][]byte, 0)
	args = append(args, []byte(targetAppId), []byte(contractAddr), []byte("transfer"), []byte(from), []byte(to), []byte(amount))
	_, err := c.Invoke("transfer", "transfer", args...)
	if err != nil {
		return err
	}
	return nil
}
