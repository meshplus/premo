package main

import (
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "Premo"
	app.Usage = "Premo is a testing framework that can help to test BitXHub."
	app.Compiled = time.Now()

	cli.VersionPrinter = func(c *cli.Context) {
		printVersion()
	}

	// global flags
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "repo",
			Usage: "Premo storage repo path",
		},
	}

	app.Commands = []*cli.Command{
		initCMD,
		versionCMD,
		testCMD,
		pierCMD,
		bitxhubCMD,
		appchainCMD,
		interchainCMD,
		statusCMD,
		serverCMD,
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
