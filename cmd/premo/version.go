package main

import (
	"fmt"
	"github.com/meshplus/premo"
	"github.com/urfave/cli/v2"
)

func getVersionCMD() *cli.Command {
	return &cli.Command{
		Name:   "version",
		Usage:  "Premo version",
		Action: version,
	}
}

func version(ctx *cli.Context) error {
	printVersion()

	return nil
}

func printVersion() {
	fmt.Printf("Premo version: %s-%s-%s\n", premo.CurrentVersion, premo.CurrentBranch, premo.CurrentCommit)
	fmt.Printf("App build date: %s\n", premo.BuildDate)
	fmt.Printf("System version: %s\n", premo.Platform)
	fmt.Printf("Golang version: %s\n", premo.GoVersion)
}
