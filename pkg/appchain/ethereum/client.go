package ethereum

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthClient struct {
	Client *ethclient.Client
	key    *ecdsa.PrivateKey
}

// New returns EthClient with ethereum addr and private key path
func New(configPath string) (*EthClient, error) {
	config, err := UnmarshalConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("unmarshal ether config error:%w", err)
	}
	etherCli, err := ethclient.Dial(config.Ether.Addr)
	if err != nil {
		return nil, err
	}

	keyByte, err := ioutil.ReadFile(filepath.Join(configPath, config.Ether.KeyPath))
	if err != nil {
		return nil, err
	}
	unlockedKey, err := keystore.DecryptKey(keyByte, "")
	if err != nil {
		return nil, err
	}

	return &EthClient{
		Client: etherCli,
		key:    unlockedKey.PrivateKey,
	}, nil
}

func (client *EthClient) Deploy(codePath string) (string, string, error) {
	// compile solidity first
	compileResult, err := compileSolidityCode(codePath)
	if err != nil {
		return "", "", err
	}

	if len(compileResult.Abi) == 0 || len(compileResult.Bins) == 0 || len(compileResult.Types) == 0 {
		return "", "", fmt.Errorf("empty contract")
	}
	// deploy a contract
	auth := bind.NewKeyedTransactor(client.key)

	for i, bin := range compileResult.Bins {
		if bin == "0x" {
			continue
		}
		parsed, err := abi.JSON(strings.NewReader(compileResult.Abi[i]))
		if err != nil {
			return "", "", err
		}

		code := strings.TrimPrefix(strings.TrimSpace(bin), "0x")
		addr, _, _, err := bind.DeployContract(auth, parsed, common.FromHex(code), client.Client)
		if err != nil {
			return "", "", err
		}
		return addr.Hex(), compileResult.Abi[i], nil
	}

	return "", "", nil
}

func (client *EthClient) Invoke(abiPath string, dstAddr string, function string, args ...string) (string, error) {

	file, err := ioutil.ReadFile(abiPath)
	if err != nil {
		return "", err
	}

	ab, err := abi.JSON(bytes.NewReader(file))
	if err != nil {
		return "", err
	}

	etherSession := &EtherSession{
		privateKey: client.key,
		etherCli:   client.Client,
		ctx:        context.Background(),
		ab:         ab,
	}

	// prepare for invoke parameters
	var argx []interface{}
	if len(args) != 0 {
		var argArr [][]byte
		for _, arg := range args {
			argArr = append(argArr, []byte(arg))
		}

		argx, err = ABIUnmarshal(ab, argArr, function)
		if err != nil {
			return "", err
		}
	}

	packed, err := ab.Pack(function, argx...)
	if err != nil {
		return "", err
	}

	invokerAddr := crypto.PubkeyToAddress(client.key.PublicKey)
	to := common.HexToAddress(dstAddr)

	if ab.Methods[function].IsConstant() {
		// for read only eth calls
		result, err := etherSession.ethCall(&invokerAddr, &to, function, packed)
		if err != nil {
			return "", err
		}

		if result == nil {
			fmt.Printf("\n======= invoke function %s =======\n", function)
			fmt.Println("no result")
			return "", nil
		}

		str := ""
		for _, r := range result {
			if r != nil {
				if reflect.TypeOf(r).String() == "[32]uint8" {
					v, ok := r.([32]byte)
					if ok {
						r = string(v[:])
					}
				}
			}
			str = fmt.Sprintf("%s,%v", str, r)
		}

		str = strings.Trim(str, ",")
		return str, nil
	}

	// for write only eth transaction
	signedTx, err := etherSession.ethTx(&invokerAddr, &to, packed)
	if err != nil {
		return "", err
	}

	return signedTx.Hash().String(), nil
}
