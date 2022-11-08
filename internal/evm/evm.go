package evm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/go-eth-client/utils"
	"github.com/meshplus/premo/internal/common"

	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
	eth "github.com/meshplus/go-eth-client"
	"github.com/meshplus/premo/internal/repo"
	"github.com/sirupsen/logrus"
)

const MaxBlockSize = 2048

var log = logrus.New()
var lock = sync.Mutex{}
var maxDelay int64
var counter int64
var delayer int64

var compileResult *eth.CompileResult
var contractAbi abi.ABI
var abiEvent *utils.AbiEvent
var code string
var function string
var address string
var args []interface{}

type Config struct {
	Concurrent   int
	TPS          int
	Duration     int
	Typ          string
	ContractPath string
	ContractName string
	AbiPath      string
	CodePath     string
	Address      string
	Function     string
	Args         string
	KeyPath      string
	JsonRpc      string
	Grpc         string
	Ctx          context.Context
	CancelFunc   context.CancelFunc
}

type Evm struct {
	config *Config
	bees   []*Bee
	client *rpcx.ChainClient
}

func New(config *Config) (*Evm, error) {
	log.WithFields(logrus.Fields{
		"concurrent": config.Concurrent,
		"tps":        config.TPS,
		"duration":   config.Duration,
		"type":       config.Typ,
	}).Info("Premo configuration")
	evm := new(Evm)
	evm.config = config
	node0 := &rpcx.NodeInfo{Addr: config.Grpc}
	pk, addr, err := repo.Node1Priv()
	if err != nil {
		return nil, err
	}
	client, err := rpcx.New(
		rpcx.WithNodesInfo(node0),
		rpcx.WithLogger(log),
		rpcx.WithPrivateKey(pk),
	)
	adminPk, err := asym.RestorePrivateKey(config.KeyPath, repo.KeyPassword)
	if err != nil {
		return nil, err
	}

	adminFrom, err := adminPk.PublicKey().Address()
	if err != nil {
		return nil, err
	}

	common.AdminNonce, err = client.GetPendingNonceByAccount(addr.String())
	if err != nil {
		return nil, err
	}
	evm.client = client

	evm.bees = make([]*Bee, 0, config.Concurrent)
	var wg sync.WaitGroup
	wg.Add(config.Concurrent)
	for i := 0; i < config.Concurrent; i++ {
		go func() {
			defer wg.Done()
			bee, err := NewBee(config, client, adminPk, adminFrom)
			if err != nil {
				log.WithFields(logrus.Fields{
					"error": err.Error(),
				}).Error("Error new bee")
				return
			}
			lock.Lock()
			evm.bees = append(evm.bees, bee)
			lock.Unlock()
			return
		}()
	}
	wg.Wait()
	log.WithFields(logrus.Fields{
		"number": len(evm.bees),
	}).Info("generate all bees")

	switch evm.config.Typ {
	case deploy:
		if err = evm.prepareDeploy(); err != nil {
			return nil, err
		}
	case deployByCode:
		if err = evm.prepareDeployByCode(); err != nil {
			return nil, err
		}
	case invoke:
		if err = evm.prepareInvoke(); err != nil {
			return nil, err
		}
	case invokeWithByte:
		if err = evm.prepareInvokeWithByte(); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("not support config type:%s", evm.config.Typ)
	}

	return evm, nil
}

func (evm *Evm) prepareDeploy() error {
	client, err := NewClient(evm.config.JsonRpc)
	if err != nil {
		return err
	}
	compileResult, err = client.Compile(evm.config.ContractPath)
	if err != nil {
		return err
	}
	for idx, compileAbi := range compileResult.Abi {
		if strings.Contains(compileResult.Names[idx], evm.config.ContractName) {
			contractAbi, err = abi.JSON(bytes.NewReader([]byte(compileAbi)))
			if err != nil {
				return err
			}
			compileResult = &eth.CompileResult{
				Abi:   []string{compileResult.Abi[idx]},
				Bin:   []string{compileResult.Bin[idx]},
				Names: []string{compileResult.Names[idx]},
			}
			break
		}
	}
	if len(evm.config.Args) != 0 {
		argSplits := strings.Split(evm.config.Args, "^")
		var argArr []interface{}
		for _, arg := range argSplits {
			if strings.Index(arg, "[") == 0 && strings.LastIndex(arg, "]") == len(arg)-1 {
				if len(arg) == 2 {
					argArr = append(argArr, make([]string, 0))
					continue
				}
				// deal with slice
				argSp := strings.Split(arg[1:len(arg)-1], ",")
				argArr = append(argArr, argSp)
				continue
			}
			argArr = append(argArr, arg)
		}
		args, err = utils.Decode(&contractAbi, "", argArr...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (evm *Evm) prepareDeployByCode() error {
	var err error
	contractAbi, err = utils.LoadAbi(evm.config.AbiPath)
	if err != nil {
		return err
	}
	abiEvent, err = utils.InitializeParameter(evm.config.AbiPath)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadFile(evm.config.CodePath)
	if err != nil {
		return err
	}
	code = string(data)
	return nil
}

func (evm *Evm) prepareInvoke() error {
	var err error
	address = evm.config.Address
	if address == "" {
		return fmt.Errorf("address must be specified")
	}
	function = evm.config.Function
	if function == "" {
		return fmt.Errorf("function must be specified")
	}
	contractAbi, err = utils.LoadAbi(evm.config.AbiPath)
	if err != nil {
		return err
	}
	if len(evm.config.Args) != 0 {
		argSplits := strings.Split(evm.config.Args, "^")
		var argArr []interface{}
		for _, arg := range argSplits {
			if strings.Index(arg, "[") == 0 && strings.LastIndex(arg, "]") == len(arg)-1 {
				if len(arg) == 2 {
					argArr = append(argArr, make([]string, 0))
					continue
				}
				// deal with slice
				argSp := strings.Split(arg[1:len(arg)-1], ",")
				argArr = append(argArr, argSp)
				continue
			}
			argArr = append(argArr, arg)
		}
		args, err = utils.Decode(&contractAbi, evm.config.Function, argArr...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (evm *Evm) prepareInvokeWithByte() error {
	var err error
	function = evm.config.Function
	if function == "" {
		return fmt.Errorf("function must be specified")
	}
	contractAbi, err = utils.LoadAbi(evm.config.AbiPath)
	if err != nil {
		return err
	}
	abiEvent, err = utils.InitializeParameter(evm.config.AbiPath)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadFile(evm.config.CodePath)
	if err != nil {
		return err
	}
	code = string(data)
	address, err = prepareContract(evm.bees, contractAbi, code, abiEvent)
	if err != nil {
		return err
	}
	return nil
}

func prepareContract(bees []*Bee, contractAbi abi.ABI, code string, abiEvent *utils.AbiEvent) (string, error) {
	fmt.Println("deploy contract")
	contractAddr, num, err := bees[0].client.DeployByCode(bees[0].pk, contractAbi, code, nil)
	if err != nil {
		return "", err
	}
	bees[0].nonce++
	fmt.Printf("contract deployed at %s on block %d\n", contractAddr, num)
	time.Sleep(5 * time.Millisecond)

	fmt.Println("register org")

	orgInput := orgStruct{
		OrgId: big.NewInt(123),
		BxmId: "test-123",
		Extra: "",
	}
	bodyBytes1, err := json.Marshal(orgInput)
	if err != nil {
		return "", err
	}
	input1, err := utils.DecodeBytes(abiEvent, "registerOrg", bodyBytes1)
	if err != nil {
		return "", err
	}
	_, err = bees[0].client.Invoke(bees[0].pk, &contractAbi, contractAddr, "registerOrg", input1)
	if err != nil {
		return "", err
	}
	bees[0].nonce++

	for i, bee := range bees {
		fmt.Printf("register %d user, userAddr:%s\n", i, crypto.PubkeyToAddress(bee.pk.PublicKey).String())
		userInstance := userInput{
			User: userStruct{
				UserAddr: crypto.PubkeyToAddress(bee.pk.PublicKey).String(),
				OrgId:    big.NewInt(123),
				Credit:   big.NewInt(1000),
				Extra:    "",
			},
		}
		bodyBytes2, err := json.Marshal(userInstance)
		if err != nil {
			return "", err
		}
		input2, err := utils.DecodeBytes(abiEvent, "registerUser", bodyBytes2)
		if err != nil {
			return "", err
		}
		_, err = bee.client.Invoke(bee.pk, &contractAbi, contractAddr, "registerUser", input2)
		if err != nil {
			return "", err
		}
		bee.nonce++
	}
	return contractAddr, nil
}

func (evm *Evm) Start() error {
	log.Info("starting evm")
	meta0, err := evm.client.GetChainMeta()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(len(evm.bees))

	for _, bee := range evm.bees {
		go func(bee *Bee) {
			wg.Done()
			err := bee.Start()
			if err != nil {
				log.WithFields(logrus.Fields{
					"error": err.Error(),
				}).Error("Error start bee")
			}
		}(bee)
	}

	wg.Wait()
	log.WithFields(logrus.Fields{
		"number": len(evm.bees),
	}).Info("start all bees")

	// listen from bitxhub block
	go evm.listenBlock()

	ticker := time.NewTicker(time.Second * time.Duration(evm.config.Duration))
	select {
	case <-ticker.C:
		err = evm.calTps(meta0)
		if err != nil {
			log.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("Error get TPS")
			return err
		}
	case <-evm.config.Ctx.Done():
		err := evm.Stop()
		if err != nil {
			log.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("Error start evm")
			return err
		}
	}
	return nil
}

func (evm *Evm) Stop() error {
	evm.config.CancelFunc()
	return nil
}

func (evm *Evm) listenBlock() {
	var (
		cnt  = int64(0)
		dly  = int64(0)
		mDly = int64(0)
	)
	ch, err := evm.client.Subscribe(context.TODO(), pb.SubscriptionRequest_BLOCK, nil)
	if err != nil {
		log.WithField("error", err).Error("subscribe block")
		return
	}
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-evm.config.Ctx.Done():
			return
		case <-ticker.C:
			c := float64(cnt)
			d := float64(dly) / float64(time.Millisecond)
			md := float64(mDly) / float64(time.Millisecond)
			log.Infof("current tps is %d, average tx delay is %fms, max tx delay is %fms", cnt, d/c, md)
			if c == 0 {
				continue
			}
			if maxDelay < mDly {
				maxDelay = mDly
			}

			cnt = 0
			dly = 0
			mDly = 0

		case data, ok := <-ch:
			if !ok {
				log.Warn("block subscription channel is closed")
				return
			}
			block := data.(*pb.Block)
			now := time.Now().UnixNano()
			for _, tx := range block.Transactions.Transactions {
				cnt++
				counter++

				txDelay := now - tx.GetTimeStamp()
				dly += txDelay
				delayer += txDelay

				if mDly < txDelay {
					mDly = txDelay
				}
			}
		}
	}
}

func (evm *Evm) calTps(meta0 *pb.ChainMeta) error {
	_ = evm.Stop()

	meta1, err := evm.client.GetChainMeta()
	if err != nil {
		return err
	}
	log.Info("Collecting tps info, please wait...")
	time.Sleep(20 * time.Second)

	skip := (meta1.Height - meta0.Height) / 8
	begin := meta0.Height + skip
	end := meta1.Height - skip

	var (
		tps      uint64
		totalTps uint64
		count    uint64
		tmpBegin = begin
	)

	for tmpBegin < end {
		if end-tmpBegin > MaxBlockSize {
			tps, err = evm.client.GetTPS(tmpBegin, tmpBegin+MaxBlockSize)
			if err != nil {
				return err
			}
			log.Infof("the TPS from block %d to %d is %d", tmpBegin, tmpBegin+MaxBlockSize, tps)
		} else {
			tps, err = evm.client.GetTPS(tmpBegin, end)
			if err != nil {
				return err
			}
			log.Infof("the TPS from block %d to %d is %d", tmpBegin, end, tps)
		}
		totalTps += tps
		count++
		tmpBegin = tmpBegin + MaxBlockSize
	}

	log.Infof("the total TPS from block %d to %d is %d", begin, end, totalTps/count)
	err = evm.client.Stop()
	if err != nil {
		return err
	}
	return nil
}

func NewClient(jsonRpc string) (*eth.EthRPC, error) {
	client, err := eth.New(jsonRpc)
	if err != nil {
		return nil, err
	}
	return client, nil
}
