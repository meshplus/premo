package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/meshplus/premo/internal/evm"
	"github.com/meshplus/premo/internal/repo"
	"github.com/urfave/cli/v2"
)

var evmCMD = &cli.Command{
	Name:  "evm",
	Usage: "test bitxhub evm function",
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
			Name:  "contract_path",
			Usage: "Specify contract path",
		},
		&cli.StringFlag{
			Name:  "abi_path",
			Usage: "Specify abi path",
		},
		&cli.StringFlag{
			Name:  "code_path",
			Usage: "Specify abiBin path",
		},
		&cli.StringFlag{
			Name:    "remote_bitxhub_addr",
			Aliases: []string{"r"},
			Usage:   "Specify remote bitxhub address",
			Value:   "localhost:8881",
		},
		&cli.StringFlag{
			Name:  "address",
			Usage: "Specify contract address",
		},
		&cli.StringFlag{
			Name:     "type",
			Usage:    "Specify test type: deploy, deployByCode, invoke, invokeWithByte",
			Required: true,
		},
		&cli.StringFlag{
			Name:    "function",
			Aliases: []string{"f"},
			Usage:   "Specify invoke function(Only use in invoke type)",
		},
		&cli.StringFlag{
			Name:    "args",
			Aliases: []string{"a"},
			Usage:   "Specify args(both deploy and invoke)",
		},
	},
	Action: evmBenchmark,
}

func evmBenchmark(ctx *cli.Context) error {
	concurrent := ctx.Int("concurrent")
	tps := ctx.Int("tps")
	duration := ctx.Int("duration")
	contractPath := ctx.String("contract_path")
	strs := strings.Split(contractPath, "/")
	abiPath := ctx.String("abi_path")
	codePath := ctx.String("code_path")
	addr := ctx.String("remote_bitxhub_addr")
	address := ctx.String("address")
	typ := ctx.String("type")
	function := ctx.String("function")
	args := ctx.String("args")
	keyPath, err := repo.Node1Path()
	if err != nil {
		return err
	}
	split := strings.Split(addr, ":")
	if len(split) != 2 {
		return err
	}
	grpc := split[0] + ":6001" + string(addr[len(addr)-1])
	c, cancelFunc := context.WithCancel(context.Background())
	config := &evm.Config{
		Concurrent:   concurrent,
		TPS:          tps,
		Duration:     duration,
		Typ:          typ,
		ContractPath: contractPath,
		ContractName: strs[len(strs)-1],
		AbiPath:      abiPath,
		CodePath:     codePath,
		Address:      address,
		Function:     function,
		Args:         args,
		KeyPath:      keyPath,
		JsonRpc:      "http://" + addr,
		Grpc:         grpc,
		Ctx:          c,
		CancelFunc:   cancelFunc,
	}
	e, err := evm.New(config)
	if err != nil {
		return err
	}
	handleEvmShutdown(e)

	err = e.Start()
	if err != nil {
		return err
	}
	return nil
}

func handleEvmShutdown(e *evm.Evm) {
	var stop = make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM)
	signal.Notify(stop, syscall.SIGINT)
	go func() {
		<-stop
		fmt.Println("received interrupt signal, shutting down...")
		if err := e.Stop(); err != nil {
			panic(err)
		}
		os.Exit(0)
	}()
}
