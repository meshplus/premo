package interface_tester

import (
	"encoding/json"
	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
)

type account struct {
	Type    string `json:"type"`
	Balance uint64 `json:"balance"`
}

func (suite *Snake) TestGetAccount() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.registerAppchain(pk, "fabric")
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
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	suite.registerAppchain(pk, "fabric")
	from, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	account := from.String() + "123"

	url := getURL("account_balance/" + account)

	data, err := httpGet(url)
	suite.Require().Nil(err)
	suite.Require().Contains(string(data), "invalid account address")
}
func (suite *Snake) TestGetAccountWithInvalidAddress02() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	suite.registerAppchain(pk, "fabric")
	from, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	account := from.String() + "我@#"

	url := getURL("account_balance/" + account)

	data, err := httpGet(url)
	suite.Require().Nil(err)
	suite.Require().Contains(string(data), "invalid account address")
}
