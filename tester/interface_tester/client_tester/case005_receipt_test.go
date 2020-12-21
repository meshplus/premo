package interface_tester

import (
	"time"
)

func (suite Snake) TestGetReceiptIsTrue() {
	hash, err := suite.sendInterchain()
	suite.Require().Nil(err)

	//wait for bitxhub
	time.Sleep(time.Second * 3)
	url := getURL("receipt/" + hash)
	suite.Require().Nil(err)

	data, err := httpGet(url)
	suite.Require().Nil(err)
	suite.Require().NotContains(string(data), "error")
}

func (suite Snake) TestGetReceiptIsFalse() {

	hash, err := suite.sendInterchain()
	suite.Require().Nil(err)

	//wait for bitxhub
	time.Sleep(time.Second * 3)
	hashByte := []byte(hash)
	hashByte[len(hash)-1] = hashByte[len(hash)-1] + 1

	url := getURL("receipt/" + string(hashByte))
	suite.Require().Nil(err)

	data, err := httpGet(url)
	suite.Require().Nil(err)
	suite.Require().Contains(string(data), "error")
}
