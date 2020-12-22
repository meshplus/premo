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
}

func (suite *Snake) TestGetAccountWithInvalidAddress01() {
	suite.registerAppchain(suite.pk, "fabric")
	from, err := suite.pk.PublicKey().Address()
	suite.Require().Nil(err)
	account := from.String() + "123"

	url := getURL("account_balance/" + account)

	data, err := httpGet(url)
	suite.Require().Nil(err)
	suite.Require().Contains(string(data), "invalid account address")
}
func (suite *Snake) TestGetAccountWithInvalidAddress02() {
	suite.registerAppchain(suite.pk, "fabric")
	from, err := suite.pk.PublicKey().Address()
	suite.Require().Nil(err)
	account := from.String() + "我@#"

	url := getURL("account_balance/" + account)

	data, err := httpGet(url)
	suite.Require().Nil(err)
	suite.Require().Contains(string(data), "invalid account address")
}
