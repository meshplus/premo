package tester

import (
	"io/ioutil"
	"path/filepath"
	"strconv"
	"time"

	"github.com/meshplus/premo/pkg/appchain/fabric"

	"github.com/tidwall/gjson"

	"github.com/meshplus/premo/pkg/appchain/ethereum"
	"github.com/stretchr/testify/suite"
)

type Interchain struct {
	suite.Suite
	ethRepo          string
	ethAppchainId    string
	fabricAppchainId string
	ethClient        *EthClient
	fabricClient     *FabricClient
}

func (suite *Interchain) SetupSuite() {
	suite.NotNil(suite.ethRepo)

	contractAddrData, err := ioutil.ReadFile("test_data/ethereum/address.json")
	suite.Nil(err)
	result := gjson.GetBytes(contractAddrData, "transfer")
	contractAddr := result.String()

	ethClient, err := ethereum.New("http://localhost:8545", filepath.Join(suite.ethRepo, "account.key"))
	suite.Nil(err)

	suite.ethClient = &EthClient{
		EthClient:    ethClient,
		abiPath:      "test_data/ethereum/transfer.abi",
		contractAddr: contractAddr,
	}

	fabricClient, err := fabric.New("test_data/fabric")
	suite.Nil(err)
	suite.fabricClient = &FabricClient{fabricClient}
}

func (suite *Interchain) TestEth2Fabric() {
	username := "Alice"
	amount := "1"

	beforeBalance, err := suite.ethClient.GetBalance(username)
	suite.Nil(err)

	fabricBeforeBalance, err := suite.fabricClient.GetBalance(username)
	suite.Nil(err)

	err = suite.ethClient.InterchainTransfer(suite.fabricAppchainId, username, username, amount)
	suite.Nil(err)
	afterBalance, err := suite.ethClient.GetBalance(username)
	suite.Assert(amount, beforeBalance, afterBalance)

	time.Sleep(5 * time.Second)

	fabricAfterBalance, err := suite.fabricClient.GetBalance(username)
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

	err = suite.fabricClient.InterchainTransfer(suite.ethAppchainId, suite.ethClient.contractAddr, username, username, amount)
	suite.Nil(err)

	fabricAfterBalance, err := suite.fabricClient.GetBalance(username)
	suite.Nil(err)

	time.Sleep(5 * time.Second)

	afterBalance, err := suite.ethClient.GetBalance(username)
	suite.Assert(amount, beforeBalance, afterBalance)

	suite.Assert(amount, fabricAfterBalance, fabricBeforeBalance)
}

func (suite *Interchain) TestEth2Eth() {

}

func (suite *Interchain) TestFabric2Fabric() {

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
