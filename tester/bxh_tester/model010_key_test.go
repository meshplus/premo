package bxh_tester

import (
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/premo/internal/repo"
)

type Model10 struct {
	*Snake
}

func (suite *Model10) SetupTest() {
	suite.T().Parallel()
}

//tc：根据正确的密码从私钥文件中读取私钥，私钥读取成功
func (suite *Model10) Test1001_RestoreKeyIsSuccess() {
	keyPath, err := repo.Node1Path()
	suite.Require().Nil(err)
	_, err = asym.RestorePrivateKey(keyPath, repo.KeyPassword)
	suite.Require().Nil(err)
}

//tc：根据错误的密码从私钥文件中读取私钥，私钥读取失败
func (suite *Model10) Test1002_RestoreKeyIsFail() {
	keyPath, err := repo.Node1Path()
	suite.Require().Nil(err)
	_, err = asym.RestorePrivateKey(keyPath, "")
	suite.Require().NotNil(err)
}
