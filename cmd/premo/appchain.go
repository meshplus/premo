package main

import (
	"fmt"

	"github.com/meshplus/premo/internal/repo"
	"github.com/meshplus/premo/pkg/constant"
	"github.com/meshplus/premo/pkg/execute"
	"github.com/urfave/cli/v2"
)

var appchainCMD = &cli.Command{
	Name:  "appchain",
	Usage: "Bring up the appchain network",
	Subcommands: []*cli.Command{
		{
			Name:  "start",
			Usage: "Start the appchain network",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "type",
					Aliases: []string{"t"},
					Value:   constant.FABRIC,
					Usage:   "start the pier, one of the fabric or ethereum",
				},
			},
			Action: startAppchain,
		},
		{
			Name:  "stop",
			Usage: "Stop the appchain network",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "type",
					Aliases: []string{"t"},
					Value:   constant.FABRIC,
					Usage:   "stop the pier, one of the fabric or ethereum",
				},
			},
			Action: stopAppchain,
		},
	},
}

func stopAppchain(ctx *cli.Context) error {
	appchain := ctx.String("type")
	repoRoot, err := repo.PathRoot()
	if err != nil {
		return fmt.Errorf("please 'premo init' first")
	}
	return downAppchain(repoRoot, appchain)
}

func downAppchain(repoRoot, appchain string) error {
	args := make([]string, 0)
	args = append(args, "run_appchain.sh", "down", appchain)

	err := execute.ExecuteShell(repoRoot, args...)
	if err != nil {
		return err
	}
	return nil
}

func startAppchain(ctx *cli.Context) error {
	appchain := ctx.String("type")
	repoRoot, err := repo.PathRoot()
	if err != nil {
		return fmt.Errorf("please 'premo init' first")
	}

	err = runAppchain(appchain, repoRoot)
	if err != nil {
		return err
	}
	return nil
}

func runAppchain(appchain, repoRoot string) error {
	args := make([]string, 0)
	args = append(args, "run_appchain.sh", "up", appchain)

	err := execute.ExecuteShell(repoRoot, args...)
	if err != nil {
		return err
	}
	return nil
}
