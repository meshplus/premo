package main

import (
	"fmt"

	"github.com/meshplus/premo/internal/repo"
	"github.com/meshplus/premo/pkg/exec"
	"github.com/urfave/cli/v2"
)

const (
	FABRIC   = "fabric"
	ETHEREUM = "ethereum"
	FLATO    = "flato"
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
					Value:   ETHEREUM,
					Usage:   "bring up the pier, fabric or ethereum or flato",
				},
				&cli.StringFlag{
					Name:    "version",
					Aliases: []string{"v"},
					Value:   "v1.18",
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
					Value:   ETHEREUM,
					Usage:   "stop the pier , fabric or ethereum or flato",
				},
			},
			Action: stopPier,
		},
	},
}

func startPier(ctx *cli.Context) error {
	repoRoot, err := repo.PathRootWithDefault()
	if err != nil {
		return fmt.Errorf("please 'premo init' first")
	}

	pierVersion := ctx.String("version")
	appchain := ctx.String("type")
	switch pierVersion {
	case "v1.18":
		repoRoot = repoRoot + "/quick-cross-chain-v1.18"
	}

	err = runPier(repoRoot, appchain)
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
	var arg string
	switch appchain {
	case ETHEREUM:
		arg = "./stop_pier_ether.sh"
	case FABRIC:
		arg = "./stop_pier_fabric.sh"
	case FLATO:
		arg = "./stop_pier_flato.sh"
	}
	_, err := exec.ExecuteShell(repoRoot, arg)
	if err != nil {
		return err
	}
	return nil
}

func runPier(repoRoot, appchain string) error {
	var arg string
	switch appchain {
	case ETHEREUM:
		arg = "./2.start_pier_ether.sh"
	case FABRIC:
		arg = "./2.start_pier_fabric.sh"
	case FLATO:
		arg = "./2.start_pier_flato.sh"
	}
	_, err := exec.ExecuteShell(repoRoot, arg)
	if err != nil {
		return err
	}
	return nil
}
