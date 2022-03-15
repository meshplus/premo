package bxh_tester

import (
	"sync/atomic"
	"time"

	"github.com/Rican7/retry"
	"github.com/Rican7/retry/backoff"
	"github.com/Rican7/retry/strategy"
	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/meshplus/premo/internal/repo"
)

type Model8 struct {
	*Snake
}

//tc：根据正确交易hash获取交易，交易获取成功
func (suite *Model8) Test0801_GetTxByHashIsSuccess() {
	pk, from, err := repo.Node2Priv()
	suite.Require().Nil(err)
	_, to, err := repo.KeyPriv()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	data := &pb.TransactionData{
		Amount: "1",
	}
	payload, err := data.Marshal()
	suite.Require().Nil(err)
	tx := &pb.BxhTransaction{
		From:      from,
		To:        to,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}
	nonce := atomic.AddUint64(&nonce2, 1)
	hash, err := client.SendTransaction(tx, &rpcx.TransactOpts{
		Nonce: nonce,
	})
	suite.Require().Nil(err)
	var res *pb.GetTransactionResponse
	err1 := retry.Retry(func(attempt uint) error {
		pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
		suite.Require().Nil(err)
		client1 := suite.NewClient(pk)
		res, err = client1.GetTransaction(hash)
		if err != nil {
			return err
		}
		return nil
	},
		strategy.Limit(10),
		strategy.Backoff(backoff.Fibonacci(500*time.Millisecond)),
	)
	suite.Require().Nil(err1)
	suite.Require().NotNil(res)
}

//tc：根据正确交易hash获取回执，回执获取成功
func (suite *Model8) Test0802_GetReceiptByHashIsSuccess() {
	pk, from, err := repo.Node2Priv()
	suite.Require().Nil(err)
	_, to, err := repo.KeyPriv()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	data := &pb.TransactionData{
		Amount: "1",
	}
	payload, err := data.Marshal()
	suite.Require().Nil(err)
	tx := &pb.BxhTransaction{
		From:      from,
		To:        to,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}
	nonce := atomic.AddUint64(&nonce2, 1)
	res, err := client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		Nonce: nonce,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
}
