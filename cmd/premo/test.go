package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gobuffalo/packr"
	"github.com/meshplus/premo/internal/bitxhub"
	"github.com/urfave/cli/v2"
)

var testCMD = &cli.Command{
	Name:  "test",
	Usage: "test bitxhub function",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:    "concurrent",
			Aliases: []string{"c"},
			Value:   100,
			Usage:   "concurrent number",
		},
		&cli.IntFlag{
			Name:    "tps",
			Aliases: []string{"t"},
			Value:   500,
			Usage:   "all tx number",
		},
		&cli.IntFlag{
			Name:    "duration",
			Aliases: []string{"d"},
			Value:   60,
			Usage:   "test duration",
		},
		&cli.StringFlag{
			Name:    "key_path",
			Aliases: []string{"k"},
			Usage:   "Specific key path",
		},
		&cli.StringFlag{
			Name:  "type",
			Usage: "Specific tx type: interchain, data, transfer",
			Value: "transfer",
		},
	},
	Action: benchmark,
}

func benchmark(ctx *cli.Context) error {
	box := packr.NewBox(ConfigPath)
	val, err := box.Find("fabric.validators")
	if err != nil {
		return err
	}
	contract, err := box.Find("rule.wasm")
	if err != nil {
		return err
	}
	config := &bitxhub.Config{
		Concurrent: ctx.Int("concurrent"),
		TPS:        ctx.Int("tps"),
		Duration:   ctx.Int("duration"),
		Type:       ctx.String("type"),
		KeyPath:    ctx.String("key_path"),
		Validator:  string(val),
		Rule:       contract,
	}

	broker, err := bitxhub.New(config)
	if err != nil {
		return err
	}

	handleShutdown(broker)

	err = broker.Start(config.Type)
	if err != nil {
		return err
	}

	return nil
}

func handleShutdown(node *bitxhub.Broker) {
	var stop = make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM)
	signal.Notify(stop, syscall.SIGINT)
	go func() {
		<-stop
		fmt.Println("received interrupt signal, shutting down...")
		if err := node.Stop(); err != nil {
			panic(err)
		}
		os.Exit(0)
	}()
}
