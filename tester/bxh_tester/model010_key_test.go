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

func (suite *Model10) Test1001_RestoreKeyIsSuccess() {
	keyPath, err := repo.Node2Path()
	suite.Require().Nil(err)
	_, err = asym.RestorePrivateKey(keyPath, repo.KeyPassword)
	suite.Require().Nil(err)
}

func (suite *Model10) Test1002_RestoreKeyIsFail() {
	keyPath, err := repo.Node2Path()
	suite.Require().Nil(err)
	_, err = asym.RestorePrivateKey(keyPath, "")
	suite.Require().NotNil(err)
}
