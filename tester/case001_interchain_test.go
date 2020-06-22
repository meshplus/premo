package tester

import (
	"strconv"
	"time"

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

	beforeBalance, err := suite.ethClient.GetBalance(username)
	suite.Nil(err)

	logger.Infof("before Aline's eth balance:%s", beforeBalance)
	fabricBeforeBalance, err := suite.fabricClient.GetBalance(username)
	suite.Nil(err)
	logger.Infof("before Aline's fabric balance:%s", fabricBeforeBalance)

	err = suite.ethClient.InterchainTransfer(suite.fabricClient.appchainId, username, username, amount)
	suite.Nil(err)
	afterBalance, err := suite.ethClient.GetBalance(username)
	logger.Infof("after Aline's eth balance:%s", afterBalance)
	suite.Assert(amount, beforeBalance, afterBalance)

	time.Sleep(5 * time.Second)

	fabricAfterBalance, err := suite.fabricClient.GetBalance(username)
	logger.Infof("after Aline's fabric balance:%s", fabricBeforeBalance)
	suite.Nil(err)

	suite.Assert(amount, fabricAfterBalance, fabricBeforeBalance)
}

func (suite *Interchain) TestFabric2Eth() {
	username := "Alice"
	amount := "1"

	fabricBeforeBalance, err := suite.fabricClient.GetBalance(username)
	suite.Nil(err)

	beforeBalance, err := suite.ethClient.GetBalance(username)
	suite.Nil(err)

	err = suite.fabricClient.InterchainTransfer(suite.ethClient.appchainId, suite.ethClient.contractAddr, username, username, amount)
	suite.Nil(err)

	fabricAfterBalance, err := suite.fabricClient.GetBalance(username)
	suite.Nil(err)

	time.Sleep(5 * time.Second)

	afterBalance, err := suite.ethClient.GetBalance(username)
	suite.Assert(amount, beforeBalance, afterBalance)

	suite.Assert(amount, fabricAfterBalance, fabricBeforeBalance)
}

func (suite *Interchain) Assert(expected, before, after string) {
	beforeI, err := strconv.Atoi(before)
	suite.Nil(err)
	afterI, err := strconv.Atoi(after)
	suite.Nil(err)

	expectedI, err := strconv.Atoi(expected)
	suite.Nil(err)

	suite.Equal(expectedI, beforeI-afterI)
}
