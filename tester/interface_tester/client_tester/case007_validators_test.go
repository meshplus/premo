package interface_tester

func (suite Snake) TestGetValidators() {
	url := getURL("info?type=2")

	data, err := httpGet(url)
	suite.Require().Nil(err)

	ret, err := parseResponse(data)
	suite.Require().Nil(err)

	retJson, err := prettyJson(ret)
	suite.Require().Nil(err)
	suite.Require().NotNil(retJson)
}
