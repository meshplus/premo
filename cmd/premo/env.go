package main

import (
	"fmt"
	"strconv"

	"github.com/meshplus/premo/internal/repo"
	"github.com/meshplus/premo/pkg/constant"
	"github.com/meshplus/premo/pkg/execute"
	"github.com/urfave/cli/v2"
)

func createEnvCMD() *cli.Command {
	return &cli.Command{
		Name:  "env",
		Usage: "create the interchain env",
		Subcommands: []*cli.Command{
			{
				Name:  "interchain",
				Usage: "create the interchain network",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "source_appchain",
						Value: constant.FABRIC,
						Usage: "bring up the source appchain network, one of the fabric or ethereum",
					},
					&cli.StringFlag{
						Name:  "target_appchain",
						Value: constant.ETHEREUM,
						Usage: "bring up the target appchain network, one of the fabric or ethereum",
					},
					&cli.IntFlag{
						Name:  "bitxhub_num",
						Value: 4,
						Usage: "the number of the bitxhub nodes",
					},
					&cli.StringFlag{
						Name:  "bitxhub_version",
						Value: "master",
						Usage: "the version of the bitxhub checkout",
					},
					&cli.StringFlag{
						Name:  "pier_version",
						Value: "master",
						Usage: "the version of the pier checkout",
					},
				},
				Action: createInterchainNetwork,
			},
		},
	}
}

func createInterchainNetwork(ctx *cli.Context) error {
	sourceAppchain := ctx.String("source_appchain")
	targetAppchain := ctx.String("target_appchain")
	num := ctx.Int("bitxhub_num")
	version := ctx.String("bitxhub_version")
	pierVersion := ctx.String("pier_version")

	repo, err := repo.PathRoot()
	if err != nil {
		return fmt.Errorf("please 'premo init' first")
	}
	if err := runBitXHub(num, repo, version); err != nil {
		return err
	}
	if err := runAppchain(sourceAppchain, repo); err != nil {
		return err
	}
	if err := runPier(sourceAppchain, repo, pierVersion); err != nil {
		return err
	}
	if err := runAppchain(targetAppchain, repo); err != nil {
		return err
	}
	if err := runPier(targetAppchain, repo, pierVersion); err != nil {
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

	task := execute.ExecTask{
		Command:      "/bin/bash",
		Args:         args,
		Repo:         repo,
		StreamStdio:  true,
		PrintCommand: true,
	}
	result, err := task.Execute()
	if err != nil {
		return fmt.Errorf("execute run_appchain.sh error:%w", err)
	}
	if result.ExitCode != 0 {
		return fmt.Errorf("execute run_appchain.sh error:%s", result.Stderr)
	}
	return nil
}

func runPier(appchain, repo, pierVersion string) error {
	args := make([]string, 0)
	switch appchain {
	case constant.FABRIC:
		args = append(args, "run_pier.sh", "up", "-t", constant.FABRIC, "-r", ".pier_fabric", "-v", pierVersion)
	case constant.ETHEREUM:
		args = append(args, "run_pier.sh", "up", "-t", constant.ETHEREUM, "-r", ".pier_ethereum", "-v", pierVersion)
	default:
		return fmt.Errorf("pier mode must be one of the FABRIC or ETHEREUM")
	}

	task := execute.ExecTask{
		Command:      "/bin/bash",
		Args:         args,
		Repo:         repo,
		StreamStdio:  true,
		PrintCommand: true,
	}
	result, err := task.Execute()
	if err != nil {
		return fmt.Errorf("execute run_pier.sh error:%w", err)
	}
	if result.ExitCode != 0 {
		return fmt.Errorf("execute run_appchain.sh error:%s", result.Stderr)
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
	task := execute.ExecTask{
		Command:      "/bin/bash",
		Args:         args,
		Repo:         repo,
		StreamStdio:  true,
		PrintCommand: true,
	}
	result, err := task.Execute()
	if err != nil {
		return fmt.Errorf("execute run_bitxhub.sh error:%w", err)
	}
	if result.ExitCode != 0 {
		return fmt.Errorf("execute run_appchain.sh error:%s", result.Stderr)
	}
	return nil
}
