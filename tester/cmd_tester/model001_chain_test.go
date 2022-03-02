package cmd_tester

func (suite Snake) Test() {
	_, err := suite.ExecuteShell("/home/jiuhuche120", "pwd")
	suite.Require().Nil(err)
}
