package interface_tester

func (suite *Snake) TestGetAccount() {
	suite.registerAppchain(suite.pk, "fabric")
	from, err := suite.pk.PublicKey().Address()
	suite.Require().Nil(err)

	url := getURL("account_balance/" + from.String())

	data, err := httpGet(url)
	suite.Require().Nil(err)

	ret, err := parseResponse(data)
	suite.Require().Nil(err)

	retJson, err := prettyJson(ret)
	suite.Require().Nil(err)
	suite.Require().Contains(retJson, "normal")
	suite.Require().Contains(retJson, "100000000")
}
