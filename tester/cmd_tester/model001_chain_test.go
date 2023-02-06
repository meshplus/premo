package cmd_tester

func (suite *Snake) Test() {
	_, err := suite.ExecuteShell("", "pwd")
	suite.Require().Nil(err)
}
