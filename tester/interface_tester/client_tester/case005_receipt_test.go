package interface_tester

import (
	"time"
)

func (suite Snake) TestGetReceipt() {
	hash, err := suite.sendInterchain()
	suite.Require().Nil(err)

	//wait for bitxhub
	time.Sleep(time.Second * 3)
	url := getURL("receipt/" + hash)
	suite.Require().Nil(err)

	data, err := httpGet(url)
	suite.Require().Nil(err)
	suite.Require().Contains(string(data), "SUCCESS")
	suite.Require().NotContains(string(data), "error")
}

func (suite Snake) TestGetReceiptWithNonexistent() {
	wrongHash := "0x0000000000000000000000000000000012345678900000000000000000000000"

	url := getURL("receipt/" + wrongHash)

	data, err := httpGet(url)
	suite.Require().Nil(err)
	suite.Require().Contains(string(data), "error")
	suite.Require().Contains(string(data), "not found in DB")
}
func (suite Snake) TestGetReceiptWithInvalidFormat() {
	wrongHash := "0x0000000000000000000000000000000012345678900000000000000000000000"

	url := getURL("receipt/" + wrongHash + "123!@#")

	data, err := httpGet(url)
	suite.Require().Nil(err)
	suite.Require().Contains(string(data), "error")
	suite.Require().Contains(string(data), "invalid format of receipt hash for querying receipt")
}
