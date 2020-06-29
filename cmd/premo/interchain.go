package main

import (
	"fmt"

	"github.com/meshplus/premo/internal/repo"
	"github.com/meshplus/premo/pkg/constant"
	"github.com/urfave/cli/v2"
)

var interchainCMD = &cli.Command{
	Name:  "interchain",
	Usage: "Bring up the interchain network",
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
}

func createInterchainNetwork(ctx *cli.Context) error {
	sourceAppchain := ctx.String("source_appchain")
	targetAppchain := ctx.String("target_appchain")
	num := ctx.Int("bitxhub_num")
	version := ctx.String("bitxhub_version")
	pierVersion := ctx.String("pier_version")

	repoRoot, err := repo.PathRootWithDefault("")
	if err != nil {
		return fmt.Errorf("please 'premo init' first")
	}
	if err := runBitXHub(num, repoRoot, version); err != nil {
		return err
	}
	if err := runAppchain(sourceAppchain, repoRoot); err != nil {
		return err
	}
	if err := runPier(sourceAppchain, repoRoot, pierVersion); err != nil {
		return err
	}
	if err := runAppchain(targetAppchain, repoRoot); err != nil {
		return err
	}
	if err := runPier(targetAppchain, repoRoot, pierVersion); err != nil {
		return err
	}

	return nil
}
