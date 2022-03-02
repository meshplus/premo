package cmd_tester

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestTester(t *testing.T) {
	suite.Run(t, &Snake{})
}
