package main

import (
	"fmt"
	"strconv"

	"github.com/meshplus/premo/internal/repo"
	"github.com/meshplus/premo/pkg/execute"
	"github.com/urfave/cli/v2"
)

var bitxhubCMD = &cli.Command{
	Name:  "bitxhub",
	Usage: "Start or stop the bitxhub cluster",
	Subcommands: []*cli.Command{
		{
			Name:  "start",
			Usage: "Start the bitxhub cluster",
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:    "num",
					Aliases: []string{"n"},
					Value:   4,
					Usage:   "the number of the bitxhub nodes",
				},
				&cli.StringFlag{
					Name:    "version",
					Aliases: []string{"v"},
					Value:   "master",
					Usage:   "the version of the bitxhub checkout",
				},
			},
			Action: startBitxhub,
		},
		{
			Name:   "stop",
			Usage:  "Stop the bitxhub cluster",
			Action: stopBitxhub,
		},
	},
}

func stopBitxhub(ctx *cli.Context) error {
	repoRoot, err := repo.PathRoot()
	if err != nil {
		return fmt.Errorf("please 'premo init' first")
	}
	return downBitxhub(repoRoot)
}

func downBitxhub(repoRoot string) error {
	args := make([]string, 0)
	args = append(args, "run_bitxhub.sh", "down")
	err := execute.ExecuteShell(repoRoot, args...)
	if err != nil {
		return err
	}
	return nil
}

func startBitxhub(ctx *cli.Context) error {
	num := ctx.Int("num")
	version := ctx.String("version")

	repoRoot, err := repo.PathRoot()
	if err != nil {
		return fmt.Errorf("please 'premo init' first")
	}
	err = runBitXHub(num, repoRoot, version)
	if err != nil {
		return err
	}
	return nil
}

func runBitXHub(num int, repoRoot, version string) error {
	var mode string
	if num > 1 {
		mode = "cluster"
	} else {
		mode = "solo"
	}

	args := make([]string, 0)
	args = append(args, "run_bitxhub.sh", "up", mode, strconv.Itoa(num), version)

	err := execute.ExecuteShell(repoRoot, args...)
	if err != nil {
		return err
	}
	return nil
}
