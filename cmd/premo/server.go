package main

import (
	"github.com/gobuffalo/packr"
	"github.com/meshplus/bitxhub-kit/log"
	"github.com/meshplus/premo/internal/api"
	"github.com/meshplus/premo/internal/bitxhub"
	"github.com/meshplus/premo/internal/repo"
	"github.com/urfave/cli/v2"
)

var serverCMD = &cli.Command{
	Name:  "server",
	Usage: "start bitxhub as http server",
	Flags: []cli.Flag{
		&cli.Uint64Flag{
			Name:    "port",
			Aliases: []string{"p"},
			Value:   9999,
			Usage:   "http server port",
		},
		&cli.IntFlag{
			Name:    "concurrent",
			Aliases: []string{"c"},
			Value:   100,
			Usage:   "concurrent number",
		},
		&cli.StringFlag{
			Name:    "key_path",
			Aliases: []string{"k"},
			Usage:   "Specify key path",
		},
		&cli.IntFlag{
			Name:    "tps",
			Aliases: []string{"t"},
			Value:   500,
			Usage:   "all tx number",
		},
		&cli.StringSliceFlag{
			Name:    "remote_bitxhub_addr",
			Aliases: []string{"r"},
			Usage:   "Specify remote bitxhub address",
			Value:   cli.NewStringSlice("localhost:60011"),
		},
	},
	Action: server,
}

func server(ctx *cli.Context) error {
	box := packr.NewBox(repo.ConfigPath)
	//val, err := box.Find("fabric.validators")
	val, err := box.Find("single_validator")
	if err != nil {
		return err
	}
	proof, err := box.Find("proof_1.0.0_rc")
	if err != nil {
		return err
	}

	contract, err := box.Find("rule.wasm")
	if err != nil {
		return err
	}

	keyPath := ctx.String("key_path")
	if keyPath == "" {
		keyPath, err = repo.Node4Path()
		if err != nil {
			return err
		}
	}

	port := ctx.Uint64("port")

	config := &bitxhub.Config{
		Concurrent:  ctx.Int("concurrent"),
		BitxhubAddr: ctx.StringSlice("remote_bitxhub_addr"),
		Validator:   string(val),
		KeyPath:     keyPath,
		Proof:       proof,
		Rule:        contract,
		Appchain:    "fabric:simple",
	}

	server, err := api.NewServer(port, config, log.NewWithModule("server"))
	if err != nil {
		return err
	}

	if err := server.Start(); err != nil {
		return err
	}

	return nil
}
