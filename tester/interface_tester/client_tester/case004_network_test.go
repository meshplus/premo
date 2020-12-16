package interface_tester

func (suite *Snake) TestGetNetwork() {
	url := getURL("info?type=1")

	data, err := httpGet(url)
	suite.Require().Nil(err)

	ret, err := parseResponse(data)
	suite.Require().Nil(err)

	retJson, err := prettyJson(ret)
	suite.Require().Nil(err)
	suite.Require().NotNil(retJson)
	// network info
	suite.Require().Contains(retJson, "QmQUcDYCtqbpn5Nhaw4FAGxQaSSNvdWfAFcpQT9SPiezbS")
	suite.Require().Contains(retJson, "QmQW3bFn8XX1t4W14Pmn37bPJUpUVBrBjnPuBZwPog3Qdy")
	suite.Require().Contains(retJson, "QmXi58fp9ZczF3Z5iz1yXAez3Hy5NYo1R8STHWKEM9XnTL")
	suite.Require().Contains(retJson, "QmbmD1kzdsxRiawxu7bRrteDgW1ituXupR8GH6E2EUAHY4")

}
