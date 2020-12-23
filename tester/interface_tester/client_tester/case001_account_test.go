package interface_tester

import "encoding/json"

type account struct {
	Type    string `json:"type"`
	Balance uint64 `json:"balance"`
}

func (suite *Snake) TestGetAccount() {
	suite.registerAppchain(suite.pk, "fabric")
	from, err := suite.pk.PublicKey().Address()
	suite.Require().Nil(err)

	url := getURL("account_balance/" + from.String())

	data, err := httpGet(url)
	suite.Require().Nil(err)

	ret, err := parseResponse(data)
	suite.Require().Nil(err)

	accountInfo := &account{}
	suite.Require().Nil(json.Unmarshal([]byte(ret), accountInfo))
	suite.Equal(accountInfo.Type, "normal")
	suite.True(accountInfo.Balance > 0)
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
