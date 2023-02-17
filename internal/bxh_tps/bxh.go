package bxh_tps

import (
	"context"

	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/meshplus/premo/internal/bitxhub"
	"github.com/meshplus/premo/internal/repo"
	"github.com/sirupsen/logrus"
)

type Client struct {
	cli    rpcx.Client
	config *Config
	logger *logrus.Logger
}

type Config struct {
	BitxhubAddr string
	Start       int
	End         int
}

func New(config *Config) (*Client, error) {
	privateKey, _, err := repo.KeyPriv()
	if err != nil {
		return nil, err
	}
	var log = logrus.New()
	opts := []rpcx.Option{
		rpcx.WithPrivateKey(privateKey),
		rpcx.WithLogger(log),
		rpcx.WithNodesInfo(&rpcx.NodeInfo{Addr: config.BitxhubAddr}),
	}

	cli, err := rpcx.New(opts...)
	if err != nil {
		return nil, err
	}
	return &Client{
		cli:    cli,
		config: config,
		logger: log,
	}, nil
}

func (c *Client) Start(ctx context.Context) error {
	var (
		tps     float32
		txCount uint64
		round   uint64
		err     error
	)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			tps, txCount, round, err = c.calTps()
			if err != nil {
				return err
			}
			c.logger.Infof("the total TPS from block %d to %d is %f, totoal tx count is %d",
				c.config.Start, c.config.End, tps/float32(round), txCount)
			return err
		}
	}
}

func (c *Client) Stop() error {
	err := c.cli.Stop()
	if err != nil {
		return err
	}
	c.logger.Info("stop calculate bitxhub tps")
	return nil
}

func (c *Client) calTps() (float32, uint64, uint64, error) {
	tmpBegin := c.config.Start
	end := c.config.End
	var (
		response     *pb.GetTPSResponse
		totalTps     float32
		totalTxCount uint64
		round        uint64
		err          error
	)
	for tmpBegin < end {
		if end-tmpBegin > bitxhub.MaxBlockSize {
			response, err = c.cli.GetTPS(uint64(tmpBegin), uint64(tmpBegin+bitxhub.MaxBlockSize))
			if err != nil {
				c.logger.Errorf("getTps err:%s", err)
				return 0, 0, 0, err
			}
			c.logger.Infof("the TPS from block %d to %d is %f", tmpBegin, tmpBegin+bitxhub.MaxBlockSize, response.Tps)
		} else {
			response, err = c.cli.GetTPS(uint64(tmpBegin), uint64(end))
			if err != nil {
				c.logger.Errorf("getTps err:%s", err)
				return 0, 0, 0, err
			}
			c.logger.Infof("the TPS from block %d to %d is %f", tmpBegin, end, response.Tps)
		}
		totalTxCount += response.TxCount
		totalTps += response.Tps
		round++
		tmpBegin = tmpBegin + bitxhub.MaxBlockSize
	}
	return totalTps, totalTxCount, round, nil
}
