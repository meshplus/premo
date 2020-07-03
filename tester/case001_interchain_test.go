package tester

import (
	"github.com/meshplus/bitxhub-kit/log"
	"github.com/stretchr/testify/suite"
)

var logger = log.NewWithModule("interchain_test")

type Interchain struct {
	suite.Suite
	repoRoot     string
	helper       *Helper
	ethClient    *EthClientHelper
	fabricClient *FabricClientHelper
}

func (suite *Interchain) SetupSuite() {
	suite.NotNil(suite.repoRoot)
	suite.helper = &Helper{&suite.Suite}
	suite.ethClient.Helper = *suite.helper
	suite.fabricClient.Helper = *suite.helper
}

func (suite *Interchain) TestEth2Fabric() {
	username := "Alice"
	amount := "1"

	fabricBeforeBalance := suite.fabricClient.GetBalance(username)
	logger.Infof("before Aline's fabric balance:%s", fabricBeforeBalance)

	ethBeforeBalance := suite.ethClient.GetBalance(username)
	logger.Infof("before Aline's eth balance:%s", ethBeforeBalance)

	logger.Infof("%s is sending %s coin from Ethereum to Fabric", username, amount)
	suite.ethClient.InterchainTransfer(suite.fabricClient.appchainId, username, username, amount)

	var fabricAfterBalance, ethAfterBalance string
	suite.helper.Retry(
		func(attempt uint) error {
			ethAfterBalance = suite.ethClient.GetBalance(username)
			if err := suite.helper.AssertBalance(amount, ethBeforeBalance, ethAfterBalance); err != nil {
				return err
			}
			return nil
		})
	logger.Infof("after Aline's eth balance:%s", ethAfterBalance)

	suite.helper.Retry(
		func(attempt uint) error {
			fabricAfterBalance = suite.fabricClient.GetBalance(username)
			if err := suite.helper.AssertBalance(amount, fabricAfterBalance, fabricBeforeBalance); err != nil {
				return err
			}
			return nil
		})
	logger.Infof("after Aline's fabric balance:%s", fabricAfterBalance)
}

func (suite *Interchain) TestFabric2Eth() {
	username := "Alice"
	amount := "1"

	fabricBeforeBalance := suite.fabricClient.GetBalance(username)
	logger.Infof("before Aline's fabric balance:%s", fabricBeforeBalance)

	ethBeforeBalance := suite.ethClient.GetBalance(username)
	logger.Infof("before Aline's eth balance:%s", ethBeforeBalance)

	logger.Infof("%s is sending %s coin from Fabric to Ethereum", username, amount)
	suite.fabricClient.InterchainTransfer(suite.ethClient.appchainId, suite.ethClient.contractAddr, username, username, amount)

	var fabricAfterBalance, ethAfterBalance string
	suite.helper.Retry(
		func(attempt uint) error {
			fabricAfterBalance = suite.fabricClient.GetBalance(username)
			if err := suite.helper.AssertBalance(amount, fabricBeforeBalance, fabricAfterBalance); err != nil {
				return err
			}
			return nil
		})
	logger.Infof("after Aline's fabric balance:%s", fabricAfterBalance)

	suite.helper.Retry(
		func(attempt uint) error {
			ethAfterBalance = suite.ethClient.GetBalance(username)
			if err := suite.helper.AssertBalance(amount, ethAfterBalance, ethBeforeBalance); err != nil {
				return err
			}
			return nil
		})

	logger.Infof("after Aline's eth balance:%s", ethAfterBalance)
}
