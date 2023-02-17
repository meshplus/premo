package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/meshplus/premo/internal/bxh_tps"
	"github.com/urfave/cli/v2"
)

var calBxhTpsCMD = &cli.Command{
	Name:  "cal_bxh",
	Usage: "calculate bitxhub tps",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "remote_bitxhub_addr",
			Aliases: []string{"addr"},
			Value:   "localhost:60011",
			Usage:   "Specify remote bitxhub address",
		},
		&cli.IntFlag{
			Name:    "startNum",
			Aliases: []string{"start"},
			Value:   1,
			Usage:   "start block number",
		},
		&cli.IntFlag{
			Name:    "endNum",
			Aliases: []string{"end"},
			Value:   10,
			Usage:   "end block number",
		},
	},
	Action: calTps,
}

func calTps(ctx *cli.Context) error {
	config := &bxh_tps.Config{
		BitxhubAddr: ctx.String("remote_bitxhub_addr"),
		Start:       ctx.Int("startNum"),
		End:         ctx.Int("endNum"),
	}

	bxhCli, err := bxh_tps.New(config)
	if err != nil {
		return err
	}

	bxhCtx, cancel := context.WithCancel(context.Background())
	handleBxhCliShutdown(bxhCli, cancel)
	err = bxhCli.Start(bxhCtx)
	if err != nil {
		return err
	}

	return nil
}

func handleBxhCliShutdown(bxhCli *bxh_tps.Client, cancel context.CancelFunc) {
	var stop = make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM)
	signal.Notify(stop, syscall.SIGINT)
	go func() {
		<-stop
		fmt.Println("received interrupt signal, shutting down...")
		cancel()
		err := bxhCli.Stop()
		if err != nil {
			fmt.Printf("stop client err:%s", err)
		}
		os.Exit(0)
	}()
}
