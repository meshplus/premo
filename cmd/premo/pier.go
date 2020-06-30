package main

import (
	"fmt"

	"github.com/meshplus/premo/internal/repo"
	"github.com/meshplus/premo/pkg/constant"
	"github.com/meshplus/premo/pkg/exec"
	"github.com/urfave/cli/v2"
)

var pierCMD = &cli.Command{
	Name:  "pier",
	Usage: "Start or stop the pier",
	Subcommands: []*cli.Command{
		{
			Name:  "start",
			Usage: "Start the pier",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "type",
					Aliases: []string{"t"},
					Value:   constant.FABRIC,
					Usage:   "bring up the pier, one of the fabric or ethereum",
				},
				&cli.StringFlag{
					Name:    "version",
					Aliases: []string{"v"},
					Value:   "master",
					Usage:   "the version of the pier checkout",
				},
			},
			Action: startPier,
		},
		{
			Name:  "stop",
			Usage: "Stop the pier",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "type",
					Aliases: []string{"t"},
					Value:   constant.FABRIC,
					Usage:   "stop the pier , one of the fabric or ethereum",
				},
			},
			Action: stopPier,
		},
	},
}

func startPier(ctx *cli.Context) error {
	pierVersion := ctx.String("version")
	appchain := ctx.String("type")

	repoRoot, err := repo.PathRootWithDefault("")
	if err != nil {
		return fmt.Errorf("please 'premo init' first")
	}

	err = runPier(appchain, repoRoot, pierVersion)
	if err != nil {
		return err
	}
	return nil
}

func stopPier(ctx *cli.Context) error {
	appchain := ctx.String("type")

	repoRoot, err := repo.PathRoot()
	if err != nil {
		return fmt.Errorf("please 'premo init' first")
	}

	return downPier(repoRoot, appchain)
}

func downPier(repoRoot, appchain string) error {
	args := make([]string, 0)
	args = append(args, "run_pier.sh", "down", "-t", appchain)
	err := exec.ExecuteShell(repoRoot, args...)
	if err != nil {
		return err
	}
	return nil
}

func runPier(appchain, repoRoot, pierVersion string) error {
	args := make([]string, 0)
	switch appchain {
	case constant.FABRIC:
		args = append(args, "run_pier.sh", "up", "-t", constant.FABRIC, "-r", ".pier_fabric", "-v", pierVersion)
	case constant.ETHEREUM:
		args = append(args, "run_pier.sh", "up", "-t", constant.ETHEREUM, "-r", ".pier_ethereum", "-v", pierVersion)
	default:
		return fmt.Errorf("pier mode must be one of the FABRIC or ETHEREUM")
	}

	err := exec.ExecuteShell(repoRoot, args...)
	if err != nil {
		return err
	}
	return nil
}
