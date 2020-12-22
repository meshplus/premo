package interface_tester

func (suite *Snake) TestGetNetwork() {
	url := getURL("info?type=1")

	data, err := httpGet(url)
	suite.Require().Nil(err)

	ret, err := parseResponse(data)
	suite.Require().Nil(err)

	suite.Require().Contains(ret, "pid")
	suite.Require().Contains(ret, "account")
	suite.Require().Contains(ret, "hosts")
}
