package interchain_tester

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Rican7/retry"
	"github.com/Rican7/retry/backoff"
	"github.com/Rican7/retry/strategy"
	"github.com/meshplus/premo/pkg/appchain/ethereum"
	"github.com/meshplus/premo/pkg/appchain/fabric"
	"github.com/stretchr/testify/suite"
)

type Helper struct {
	suite *suite.Suite
}

type EthClientHelper struct {
	Helper
	*ethereum.EthClient
	abiPath      string
	contractAddr string
	appchainId   string
}

type FabricClientHelper struct {
	Helper
	*fabric.FabricClient
	appchainId string
}

func (c *EthClientHelper) GetBalance(username string) string {
	response, err := c.Invoke(c.abiPath, c.contractAddr, "getBalance", username)
	c.suite.Nil(err)
	return response
}

func (c *EthClientHelper) SetBalance(username, amount string) string {
	response, err := c.Invoke(c.abiPath, c.contractAddr, "setBalance", username, amount)
	c.suite.Nil(err)
	return response
}

func (c *EthClientHelper) InterchainTransfer(targetAppId, from, to, amount string) {
	args := make([]string, 0)
	args = append(args, targetAppId, "mychannel&transfer", from, to, amount)

	var err error
	c.Retry(func(attempt uint) error {
		_, err := c.Invoke(c.abiPath, c.contractAddr, "transfer", args...)
		if err != nil {
			return err
		}
		return nil
	})
	c.suite.Nil(err)
}

func (c *FabricClientHelper) GetBalance(username string) string {
	args := make([][]byte, 0)
	args = append(args, []byte(username))

	var response string
	var err error
	c.Retry(func(attempt uint) error {
		response, err = c.Invoke("transfer", "getBalance", args...)
		if err != nil {
			return err
		}
		return nil
	})
	c.suite.Nil(err)
	return response
}

func (c *FabricClientHelper) SetBalance(username, amount string) string {
	args := make([][]byte, 0)
	args = append(args, []byte(username), []byte(amount))

	var response string
	var err error
	c.Retry(func(attempt uint) error {
		response, err = c.Invoke("transfer", "setBalance", args...)
		if err != nil {
			return err
		}
		return nil
	})
	c.suite.Nil(err)
	return response
}

func (c *FabricClientHelper) InterchainTransfer(targetAppId, contractAddr, from, to, amount string) {
	args := make([][]byte, 0)
	args = append(args, []byte(targetAppId), []byte(contractAddr), []byte(from), []byte(to), []byte(amount))
	_, err := c.Invoke("transfer", "transfer", args...)
	c.suite.Nil(err)
}

func (h *Helper) Retry(action retry.Action) {
	err := retry.Retry(action,
		strategy.Limit(5),
		strategy.Backoff(backoff.Fibonacci(500*time.Millisecond)),
	)
	h.suite.Nil(err)
}

func (h *Helper) AssertBalance(expected, before, after string) error {
	beforeI, err := strconv.Atoi(before)
	if err != nil {
		return err
	}
	afterI, err := strconv.Atoi(after)
	if err != nil {
		return err
	}
	expectedI, err := strconv.Atoi(expected)
	if err != nil {
		return err
	}
	if expectedI != beforeI-afterI {
		return fmt.Errorf("not equal, expected:%d, actual:%d", expectedI, beforeI-afterI)
	}
	return nil
}
