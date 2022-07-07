package bitxhub

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/wcharczuk/go-chart/v2"

	appchain_mgr "github.com/meshplus/bitxhub-core/appchain-mgr"
	"github.com/meshplus/bitxhub-core/governance"
	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-kit/types"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/meshplus/premo/internal/repo"
	"github.com/sirupsen/logrus"
)

var index1 uint64
var index2 uint64
var index3 uint64
var adminNonce uint64
var log = logrus.New()
var To string

type Broker struct {
	config     *Config
	bees       []*bee
	client     rpcx.Client
	adminNonce uint64
	ctx        context.Context
	cancel     context.CancelFunc
	lock       sync.Mutex
}

type Config struct {
	Concurrent  int
	TPS         int
	Duration    int // s uint
	Type        string
	Validator   string
	Proof       []byte
	KeyPath     string
	BitxhubAddr []string
	Appchain    string
	Graph       bool
}

func New(config *Config) (*Broker, error) {
	log.WithFields(logrus.Fields{
		"concurrent": config.Concurrent,
		"tps":        config.TPS,
		"duration":   config.Duration,
		"type":       config.Type,
	}).Info("Premo configuration")
	adminPk, err := asym.RestorePrivateKey(config.KeyPath, repo.KeyPassword)
	if err != nil {
		return nil, err
	}

	adminFrom, err := adminPk.PublicKey().Address()
	if err != nil {
		return nil, err
	}

	node0 := &rpcx.NodeInfo{Addr: config.BitxhubAddr[0]}
	client, err := rpcx.New(
		rpcx.WithNodesInfo(node0),
		rpcx.WithLogger(log),
		rpcx.WithPrivateKey(adminPk),
	)
	if err != nil {
		return nil, err
	}

	//query nodes nonce
	_, node1Address, err := repo.Node1Priv()
	if err != nil {
		return nil, err
	}
	index1, err = client.GetPendingNonceByAccount(node1Address.String())
	if err != nil {
		return nil, err
	}

	_, node2Address, err := repo.Node2Priv()
	if err != nil {
		return nil, err
	}
	index2, err = client.GetPendingNonceByAccount(node2Address.String())
	if err != nil {
		return nil, err
	}

	_, node3Address, err := repo.Node3Priv()
	if err != nil {
		return nil, err
	}
	index3, err = client.GetPendingNonceByAccount(node3Address.String())
	if err != nil {
		return nil, err
	}

	index1 -= 1
	index2 -= 1
	index3 -= 1

	// query pending nonce for adminKey
	adminNonce, err = client.GetPendingNonceByAccount(adminFrom.String())
	if err != nil {
		return nil, err
	}
	// prepare to
	to, err := PrepareTo(client, config, adminPk, adminFrom)
	if err != nil {
		return nil, err
	}
	To = to

	bees := make([]*bee, 0, config.Concurrent)
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(config.Concurrent)
	for i := 0; i < config.Concurrent; i++ {
		go func() {
			defer wg.Done()
			bee, err := NewBee(config.TPS/config.Concurrent, adminPk, adminFrom, config, ctx)
			if err != nil {
				log.Error("New bee: ", err.Error())
				return
			}
			if config.Type == "interchain" {
				if err := bee.prepareChain(bee.config.Appchain, "fabric for law"); err != nil {
					log.Error(err)
					return
				}
			}

			bees = append(bees, bee)
		}()
	}

	wg.Wait()
	log.WithFields(logrus.Fields{
		"number": len(bees),
	}).Info("generate all bees")

	return &Broker{
		config:     config,
		bees:       bees,
		client:     client,
		adminNonce: adminNonce,
		ctx:        ctx,
		cancel:     cancel,
	}, nil
}

func (broker *Broker) Start(typ string) error {
	log.Info("starting broker")
	var wg sync.WaitGroup
	wg.Add(len(broker.bees))

	current := time.Now()

	meta0, err := broker.client.GetChainMeta()
	if err != nil {
		return err
	}

	for i := 0; i < len(broker.bees); i++ {
		go func(i int) {
			wg.Done()
			err := broker.bees[i].start(typ)
			if err != nil {
				log.Error(err)
				return
			}
			log.WithFields(logrus.Fields{
				"index": i + 1,
			}).Debug("start bee")
		}(i)
	}
	log.WithFields(logrus.Fields{
		"number": len(broker.bees),
	}).Info("start all bees")
	// 64*3600*24*7/8/1024/1024<= 5MB
	var x []time.Time
	var tpsY []float64
	var latencyY []float64
	maxTps := 0.0
	maxLatency := 0.0

	go func() {
		var (
			cnt  = int64(0)
			dly  = int64(0)
			mDly = int64(0)
		)
		ch, err := broker.client.Subscribe(context.TODO(), pb.SubscriptionRequest_BLOCK, nil)
		if err != nil {
			log.WithField("error", err).Error("subscribe block")
			return
		}
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-broker.ctx.Done():
				return
			case <-ticker.C:
				c := float64(cnt)
				d := float64(dly) / float64(time.Millisecond)
				md := float64(mDly) / float64(time.Millisecond)
				log.Infof("current tps is %d, average tx delay is %fms, max tx delay is %fms", cnt, d/c, md)
				if c == 0 {
					continue
				}
				x = append(x, time.Now())
				tpsY = append(tpsY, float64(cnt))
				if maxTps < float64(cnt) {
					maxTps = float64(cnt)
				}
				latencyY = append(latencyY, d/c)
				if maxLatency < d/c {
					maxLatency = d / c
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
	}()

	time.Sleep(time.Duration(broker.config.Duration) * time.Second)
	wg.Wait()

	_ = broker.Stop(current)

	meta1, err := broker.client.GetChainMeta()
	if err != nil {
		return err
	}
	log.Info("Collecting tps info, please wait...")
	time.Sleep(20 * time.Second)

	skip := (meta1.Height - meta0.Height) / 8
	begin := meta0.Height + skip
	end := meta1.Height - skip
	tps, err := broker.client.GetTPS(begin, end)
	if err != nil {
		return err
	}
	log.Infof("the TPS from block %d to %d is %d", begin, end, tps)
	err = broker.client.Stop()
	if err != nil {
		return err
	}
	if broker.config.Graph {
		Graph(x, tpsY, latencyY, maxTps, maxLatency)
	}
	return nil
}

func (broker *Broker) Stop(current time.Time) error {
	// prevent stop function is repeatedly called
	broker.lock.Lock()
	broker.cancel()
	// wait for goroutines inside bees to stop
	time.Sleep(10 * time.Second)

	log.Info("Bees are quiting, please wait...")
	for i := 0; i < len(broker.bees); i++ {
		err := broker.bees[i].stop()
		if err != nil {
			return err
		}
	}
	//err := broker.client.Stop()
	//if err != nil {
	//	log.Warn(err)
	//}
	delayerAvg := float64(delayer) / float64(counter)
	log.WithFields(logrus.Fields{
		"number":   counter,
		"duration": time.Since(current).Seconds(),
		"tps":      float64(counter) / time.Since(current).Seconds(),
		"tx_delay": delayerAvg / float64(time.Millisecond),
	}).Info("finish testing")
	return nil
}

func PrepareTo(client *rpcx.ChainClient, config *Config, adminPk crypto.PrivateKey, adminFrom *types.Address) (string, error) {
	pk, from, err := repo.KeyPriv()
	if err != nil {
		return "", err
	}
	bytes, err := pk.PublicKey().Bytes()
	if err != nil {
		return "", err
	}
	err = TransferFromAdmin(client, adminPk, adminFrom, from, "100")
	if err != nil {
		return "", err
	}
	args := []*pb.Arg{
		rpcx.String(from.String()),                                //chainID
		rpcx.String(from.String()),                                //chainName
		rpcx.Bytes(bytes),                                         //pubKey
		rpcx.String("Flato v1.0.3"),                               //chainType
		rpcx.Bytes([]byte("")),                                    //trustRoot
		rpcx.String("0x857133c5C69e6Ce66F7AD46F200B9B3573e77582"), //broker
		rpcx.String("desc"),                                       //desc
		rpcx.String("0x00000000000000000000000000000000000000a2"), //masterRuleAddr
		rpcx.String("https://github.com"),                         //masterRuleUrl
		rpcx.String(from.String()),                                //adminAddrs
		rpcx.String("reason"),                                     //reason
	}

	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "RegisterAppchain", &rpcx.TransactOpts{
		From:    from.String(),
		Nonce:   0,
		PrivKey: pk,
	}, args...)
	if err != nil {
		return "", err
	}
	//vote chain
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	if err != nil || result.ProposalID == "" {
		return "", fmt.Errorf("vote chain unmarshal error: %w", err)
	}
	b := bee{}
	b.config = config
	err = b.VotePass(client, result.ProposalID)
	if err != nil {
		return "", fmt.Errorf("vote chain error: %w", err)
	}
	res, err = b.GetChainStatusById(client, pk, from.String())
	if err != nil {
		return "", fmt.Errorf("getChainStatus error: %w", err)
	}
	appchain := &appchain_mgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	if err != nil || appchain.Status != governance.GovernanceAvailable {
		return "", fmt.Errorf("chain error: %w", err)
	}
	//register server
	args = []*pb.Arg{
		rpcx.String(from.String()),
		rpcx.String("mychannel&transfer"),
		rpcx.String(from.String()),
		rpcx.String("CallContract"),
		rpcx.String("test"),
		rpcx.Uint64(1),
		rpcx.String(""),
		rpcx.String("test"),
		rpcx.String("reason"),
	}
	res, err = client.InvokeBVMContract(constant.ServiceMgrContractAddr.Address(), "RegisterService", &rpcx.TransactOpts{
		From:    from.String(),
		PrivKey: pk,
	}, args...)
	if err != nil {
		return "", fmt.Errorf("register service error %w", err)
	}
	//vote server
	result = &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	if err != nil || result.ProposalID == "" {
		return "", fmt.Errorf("vote server unmarshal error: %w", err)
	}
	err = b.VotePass(client, result.ProposalID)
	if err != nil {
		return "", fmt.Errorf("vote server error: %w", err)
	}
	return from.String(), nil
}

func Graph(x []time.Time, tpsY []float64, latencyY []float64, maxTps, maxLatency float64) {
	graph := chart.Chart{
		XAxis: chart.XAxis{
			Name: "Time",
		},
		YAxis: chart.YAxis{
			Range: &chart.ContinuousRange{Min: 0, Max: maxTps},
		},
		YAxisSecondary: chart.YAxis{
			Range: &chart.ContinuousRange{Min: 0, Max: maxLatency},
		},
		Series: []chart.Series{
			chart.TimeSeries{
				Name:    "TPS",
				XValues: x,
				YValues: tpsY,
				Style: chart.Style{
					StrokeColor: chart.GetDefaultColor(0).WithAlpha(64),
					FillColor:   chart.GetDefaultColor(0).WithAlpha(64),
				},
			},
			chart.TimeSeries{
				Name:    "Latency",
				XValues: x,
				YValues: latencyY,
				YAxis:   chart.YAxisSecondary,
				Style: chart.Style{
					StrokeColor: chart.GetDefaultColor(1).WithAlpha(64),
					FillColor:   chart.GetDefaultColor(1).WithAlpha(64),
				},
			},
			chart.AnnotationSeries{
				Annotations: []chart.Value2{
					{XValue: chart.TimeToFloat64(x[len(x)-1]), YValue: tpsY[len(tpsY)-1], Label: "TPS", Style: chart.Style{StrokeColor: chart.ColorBlue}},
				},
			},
			chart.AnnotationSeries{
				YAxis: chart.YAxisSecondary,
				Annotations: []chart.Value2{
					{XValue: chart.TimeToFloat64(x[len(x)-1]), YValue: latencyY[len(latencyY)-1], Label: "Latency", Style: chart.Style{StrokeColor: chart.ColorGreen}},
				},
			},
		},
	}
	f, _ := os.Create("graph.png")
	defer f.Close()
	graph.Render(chart.PNG, f)
}
