package tester

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/meshplus/premo/pkg/constant"

	"github.com/meshplus/premo/pkg/exec"
	"github.com/stretchr/testify/suite"
)

func TestTester(t *testing.T) {
	//pathRoot, err := repo.PathRoot()
	//require.Nil(t, err)
	//
	//err = setupAppchain(constant.FABRIC, pathRoot)
	//require.Nil(t, err)
	//
	//err = setupAppchain(constant.ETHEREUM, pathRoot)
	//require.Nil(t, err)
	//
	//err = setupBitxhub(4, pathRoot, "master")
	//require.Nil(t, err)
	//
	//err = setupPier(constant.ETHEREUM, pathRoot)
	//require.Nil(t, err)
	//
	//err = setupPier(constant.FABRIC, pathRoot)
	//require.Nil(t, err)

	suite.Run(t, &Interchain{ethRepo: "test_data/ethereum"})
}

func setupBitxhub(num int, repoRoot, version string) error {
	if err := runBitXHub(num, repoRoot, version); err != nil {
		return err
	}
	return nil
}

func setupPier(appchain, repoRoot string) error {
	if err := runPier(appchain, repoRoot); err != nil {
		return err
	}
	return nil
}

func setupAppchain(appchain, repoRoot string) error {
	if err := runAppchain(appchain, repoRoot); err != nil {
		return err
	}
	return nil
}

func runAppchain(appchain, repo string) error {
	args := make([]string, 0)
	switch appchain {
	case constant.FABRIC:
		args = append(args, "run_appchain.sh", "up", constant.FABRIC)
	case constant.ETHEREUM:
		args = append(args, "run_appchain.sh", "up", constant.ETHEREUM)
	default:
		return fmt.Errorf("appchain must be one of the FABRIC or ETHEREUM")
	}
	err := exec.ExecCmd(args, repo)
	if err != nil {
		return fmt.Errorf("execute run_appchain.sh error:%w", err)
	}
	return nil
}

func runPier(appchain, repo string) error {
	args := make([]string, 0)
	switch appchain {
	case constant.FABRIC:
		args = append(args, "run_pier.sh", "up", constant.FABRIC)
	case constant.ETHEREUM:
		args = append(args, "run_pier.sh", "up", constant.ETHEREUM)
	default:
		return fmt.Errorf("pier mode must be one of the FABRIC or ETHEREUM")
	}

	err := exec.ExecCmd(args, repo)
	if err != nil {
		return fmt.Errorf("execute run_pier.sh error:%w", err)
	}
	return nil
}

func runBitXHub(num int, repo, version string) error {
	var mode string
	if num > 1 {
		mode = "cluster"
	} else {
		mode = "solo"
	}
	args := make([]string, 0)
	args = append(args, "run_bitxhub.sh", "up", mode, strconv.Itoa(num), version)
	err := exec.ExecCmd(args, repo)
	if err != nil {
		return fmt.Errorf("execute run_bitxhub.sh error:%w", err)
	}
	return nil
}
