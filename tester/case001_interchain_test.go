package tester

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Rican7/retry"
	"github.com/Rican7/retry/backoff"
	"github.com/Rican7/retry/strategy"
	"github.com/meshplus/bitxhub-kit/log"
	"github.com/stretchr/testify/suite"
)

var logger = log.NewWithModule("interchain_test")

type Interchain struct {
	suite.Suite
	repoRoot string

	ethClient    *EthClientHelper
	fabricClient *FabricClientHelper
}

func (suite *Interchain) SetupSuite() {
	suite.NotNil(suite.repoRoot)

}

func (suite *Interchain) TestEth2Fabric() {
	username := "Alice"
	amount := "1"

	ethBeforeBalance, err := suite.ethClient.GetBalance(username)
	suite.Nil(err)

	logger.Infof("before Aline's eth balance:%s", ethBeforeBalance)
	fabricBeforeBalance, err := suite.fabricClient.GetBalance(username)
	suite.Nil(err)
	logger.Infof("before Aline's fabric balance:%s", fabricBeforeBalance)

	err = suite.ethClient.InterchainTransfer(suite.fabricClient.appchainId, username, username, amount)
	suite.Nil(err)

	var fabricAfterBalance, ethAfterBalance string
	err = retry.Retry(func(attempt uint) error {
		ethAfterBalance, err = suite.ethClient.GetBalance(username)
		if err != nil {
			return err
		}
		err = AssertBalance(amount, ethBeforeBalance, ethAfterBalance)
		if err != nil {
			return err
		}

		return nil
	},
		strategy.Limit(5),
		strategy.Backoff(backoff.Fibonacci(1*time.Second)),
	)
	suite.Nil(err)
	logger.Infof("after Aline's eth balance:%s", ethAfterBalance)

	err = retry.Retry(func(attempt uint) error {
		fabricAfterBalance, err = suite.fabricClient.GetBalance(username)
		if err != nil {
			return err
		}
		err = AssertBalance(amount, fabricAfterBalance, fabricBeforeBalance)
		if err != nil {
			return err
		}

		return nil
	},
		strategy.Limit(5),
		strategy.Backoff(backoff.Fibonacci(1*time.Second)),
	)
	suite.Nil(err)
	logger.Infof("after Aline's fabric balance:%s", fabricAfterBalance)
}

func (suite *Interchain) TestFabric2Eth() {
	username := "Alice"
	amount := "1"

	fabricBeforeBalance, err := suite.fabricClient.GetBalance(username)
	suite.Nil(err)
	logger.Infof("before Aline's fabric balance:%s", fabricBeforeBalance)

	ethBeforeBalance, err := suite.ethClient.GetBalance(username)
	suite.Nil(err)
	logger.Infof("before Aline's eth balance:%s", ethBeforeBalance)

	err = suite.fabricClient.InterchainTransfer(suite.ethClient.appchainId, suite.ethClient.contractAddr, username, username, amount)
	suite.Nil(err)

	var fabricAfterBalance, ethAfterBalance string
	err = retry.Retry(func(attempt uint) error {
		fabricAfterBalance, err = suite.fabricClient.GetBalance(username)
		if err != nil {
			return err
		}
		err = AssertBalance(amount, fabricBeforeBalance, fabricAfterBalance)
		if err != nil {
			return err
		}

		return nil
	},
		strategy.Limit(5),
		strategy.Backoff(backoff.Fibonacci(1*time.Second)),
	)
	suite.Nil(err)
	logger.Infof("after Aline's fabric balance:%s", fabricAfterBalance)

	err = retry.Retry(func(attempt uint) error {
		ethAfterBalance, err = suite.ethClient.GetBalance(username)
		if err != nil {
			return err
		}
		err = AssertBalance(amount, ethAfterBalance, ethBeforeBalance)
		if err != nil {
			return err
		}

		return nil
	},
		strategy.Limit(5),
		strategy.Backoff(backoff.Fibonacci(1*time.Second)),
	)
	suite.Nil(err)
	logger.Infof("after Aline's eth balance:%s", ethAfterBalance)
}

func AssertBalance(expected, before, after string) error {
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
