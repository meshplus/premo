package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/coreos/etcd/pkg/fileutil"
	"github.com/meshplus/premo/internal/repo"
	"github.com/urfave/cli/v2"
)

var initCMD = &cli.Command{
	Name:   "init",
	Usage:  "init config home for premo",
	Action: Initialize,
}

func Initialize(ctx *cli.Context) error {
	repoRoot := ctx.String("repo")
	if repoRoot == "" {
		root, err := repo.PathRoot()
		if err != nil {
			return err
		}
		repoRoot = root
	}
	if fileutil.Exist(repoRoot) {
		fmt.Println("premo configuration file already exists")
		fmt.Println("reinitializing would overwrite your configuration, Y/N?")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		if input.Text() != "Y" && input.Text() != "y" {
			return nil
		}
	}

	return repo.Initialize(repoRoot)
}
