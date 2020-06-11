package tester

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/meshplus/premo/pkg/constant"

	"github.com/meshplus/bitxhub-kit/fileutil"
	"github.com/meshplus/premo/pkg/exec"
	"github.com/meshplus/premo/repo"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestTester(t *testing.T) {
	pathRoot, err := repo.PathRoot()
	require.Nil(t, err)

	err = setupAppchain(constant.FABRIC, pathRoot)
	require.Nil(t, err)

	err = setupAppchain(constant.ETHEREUM, pathRoot)
	require.Nil(t, err)

	suite.Run(t, &Interchain{repo: pathRoot})
}

func setupAppchain(appchain, repoPath string) error {
	repoRoot, err := repo.PathRootWithDefault(repoPath)
	if err != nil {
		return err
	}

	var pierPath string
	switch appchain {
	case constant.FABRIC:
		pierPath = filepath.Join(repoRoot, ".pier_"+appchain)
	case constant.ETHEREUM:
		pierPath = filepath.Join(repoRoot, ".pier_"+appchain)
	default:
		return fmt.Errorf("pier mode must be one of the fabric or constant.ETHEREUM")
	}

	if !fileutil.Exist(pierPath) {
		return fmt.Errorf("not found pier config directory:%s", pierPath)
	}
	//TODO: register appchain to bitxhub
	args := make([]string, 0)
	args = append(args, "run_pier.sh", pierPath, "")

	//TODO: deploy rule to bitxhub
	return exec.ExecCmd(args, repoRoot)
}
