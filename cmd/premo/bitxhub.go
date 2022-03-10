package main

import (
	"fmt"

	"github.com/coreos/etcd/pkg/fileutil"
	"github.com/meshplus/premo/internal/repo"
	"github.com/meshplus/premo/pkg/exec"
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
				&cli.StringFlag{
					Name:    "version",
					Aliases: []string{"v"},
					Value:   "v1.18",
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
		return err
	}
	if !fileutil.Exist(repoRoot) {
		return fmt.Errorf("please run `premo init` first")
	}
	return downBitxhub(repoRoot)
}

func downBitxhub(repoRoot string) error {
	_, err := exec.ExecuteShell(repoRoot, "./stop_bitxhub.sh")
	if err != nil {
		return err
	}
	return nil
}

func startBitxhub(ctx *cli.Context) error {
	repoRoot, err := repo.PathRoot()
	if err != nil {
		return err
	}
	if !fileutil.Exist(repoRoot) {
		return fmt.Errorf("please run `premo init` first")
	}
	version := ctx.String("version")
	switch version {
	case "v1.18":
		repoRoot = repoRoot + "/quick-cross-chain-v1.18"
	}

	err = runBitXHub(repoRoot)
	if err != nil {
		return err
	}
	return nil
}

func runBitXHub(repoRoot string) error {
	_, err := exec.ExecuteShell(repoRoot, "./1.start_bitxhub.sh")
	if err != nil {
		return err
	}
	return nil
}
