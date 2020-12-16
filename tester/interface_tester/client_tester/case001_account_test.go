package interface_tester

import (
	"github.com/meshplus/bitxhub-model/pb"
)

func (suite *Snake) TestGetAccount() {
	//sendInterchain
	_, _, from, _, receipt, err := suite.sendInterchainWithReceipt()
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, receipt.Status)

	url, err := getURL("account_balance/" + from.Address)
	suite.Require().Nil(err)

	data, err := httpGet(url)
	suite.Require().Nil(err)

	ret, err := parseResponse(data)
	suite.Require().Nil(err)

	retJson, err := prettyJson(ret)
	suite.Require().Nil(err)
	suite.Require().Contains(retJson, "normal")
	suite.Require().Contains(retJson, "100000000")
}
