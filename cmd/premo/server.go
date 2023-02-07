package main

import (
	"github.com/meshplus/premo/internal/server"
	"github.com/urfave/cli/v2"
)

var serverCMD = &cli.Command{
	Name:  "server",
	Usage: "Start test server",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "remote_bitxhub_addr",
			Aliases: []string{"r"},
			Usage:   "Specify remote bitxhub address",
			Value:   "localhost:60011",
		},
		&cli.IntFlag{
			Name:    "port",
			Aliases: []string{"p"},
			Usage:   "Specify server's port",
			Value:   9999,
		},
	},
	Action: serverBenchmark,
}

func serverBenchmark(ctx *cli.Context) error {
	remote := ctx.String("remote_bitxhub_addr")
	port := ctx.Int("port")
	newServer, err := server.NewServer(remote, port, 10)
	if err != nil {
		return err
	}
	newServer.Start()
	return nil
}
