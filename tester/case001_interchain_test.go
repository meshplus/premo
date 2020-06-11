package tester

import (
	"github.com/stretchr/testify/suite"
)

type Interchain struct {
	suite.Suite
	repo string
}

func (suite *Interchain) SetupSuite() {
	suite.NotNil(suite.repo)

}

func (suite *Interchain) TestEth2Fabric() {

}

func (suite *Interchain) TestFabric2Eth() {

}
