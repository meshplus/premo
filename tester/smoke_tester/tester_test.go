package bxh_tester

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestTester(t *testing.T) {
	suite.Run(t, &Model1{&Snake{}})
	suite.Run(t, &Model2{&Snake{}})
	suite.Run(t, &Model3{&Snake{}})
	suite.Run(t, &Model4{&Snake{}})
	suite.Run(t, &Model5{&Snake{}})
	suite.Run(t, &Model6{&Snake{}})
	suite.Run(t, &Model7{&Snake{}})
	suite.Run(t, &Model8{&Snake{}})
	suite.Run(t, &Model9{&Snake{}})
	suite.Run(t, &Model10{&Snake{}})
	suite.Run(t, &Model11{&Snake{}})
	suite.Run(t, &Model12{&Snake{}})
	suite.Run(t, &Model13{&Snake{}})
	suite.Run(t, &Model14{&Snake{}})
	//suite.Run(t, &Model15{&Snake{}})
	suite.Run(t, &Model16{&Snake{}})
	suite.Run(t, &Model17{&Snake{}})
}
