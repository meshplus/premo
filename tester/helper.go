package tester

import (
	"github.com/meshplus/premo/pkg/appchain/ethereum"
	"github.com/meshplus/premo/pkg/appchain/fabric"
)

type EthClient struct {
	*ethereum.EthClient
	abiPath      string
	contractAddr string
}

type FabricClient struct {
	*fabric.FabricClient
}

func (c *EthClient) GetBalance(username string) (string, error) {
	response, err := c.Invoke(c.abiPath, c.contractAddr, "getBalance", username)
	if err != nil {
		return "", err
	}
	return response, nil
}

func (c *EthClient) InterchainTransfer(targetAppId, from, to, amount string) error {
	args := make([]string, 0)
	args = append(args, targetAppId, "mychannel&transfer", "transfer", from, to, amount)
	_, err := c.Invoke(c.abiPath, c.contractAddr, "transfer", args...)
	if err != nil {
		return err
	}
	return nil
}

func (c *FabricClient) GetBalance(username string) (string, error) {
	args := make([][]byte, 0)
	args = append(args, []byte(username))
	response, err := c.Invoke("transfer", "getBalance", args...)
	if err != nil {
		return "", err
	}
	return response, nil
}

func (c *FabricClient) InterchainTransfer(targetAppId, contractAddr, from, to, amount string) error {
	args := make([][]byte, 0)
	args = append(args, []byte(targetAppId), []byte(contractAddr), []byte("transfer"), []byte(from), []byte(to), []byte(amount))
	_, err := c.Invoke("transfer", "transfer", args...)
	if err != nil {
		return err
	}
	return nil
}
