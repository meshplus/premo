package bitxhub

import (
	"fmt"
	"io/ioutil"

	"github.com/meshplus/bitxhub-kit/types"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/meshplus/premo/pkg/utils"
)

var ValidationContractAddr = types.String2Address("000000000000000000000000000000000000000c")

func DeployRules(path, key, addr string) error {
	privKey, err := utils.PrivKeyFromKey(key)
	if err != nil {
		return err
	}

	cli, err := rpcx.New(
		rpcx.WithAddrs([]string{addr}),
		rpcx.WithLogger(cfg.logger),
		rpcx.WithPrivateKey(privKey),
	)
	if err != nil {
		return err
	}
	contract, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	contractAddr, err := cli.DeployContract(contract, nil)
	if err != nil {
		return err
	}
	chainAddr, err := privKey.PublicKey().Address()
	if err != nil {
		return err
	}
	_, err = cli.InvokeContract(pb.TransactionData_BVM, ValidationContractAddr,
		"RegisterRule", nil, rpcx.String(chainAddr.String()), rpcx.String(contractAddr.String()))
	if err != nil {
		return err
	}

	fmt.Println("Deploy rule to bitxhub successfully")
	return nil
}
