package bxh_tester

import (
	"math/rand"
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

func (suite *Model8) SetupTest() {
	suite.T().Parallel()
}

//tc：发送交易，from为空，交易发送失败
func (suite *Model8) Test0801_TXEmptyFromIsFail() {
	pk, _, err := repo.Node2Priv()
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
		To:        to,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}
	nonce := atomic.LoadUint64(&nonce2)
	_, err = client.SendTransaction(tx, &rpcx.TransactOpts{
		Nonce: nonce,
	})
	suite.Require().NotNil(err)
}

//tc：发送交易，to为空，交易发送失败
func (suite *Model8) Test0802_TXEmptyToIsFail() {
	pk, from, err := repo.Node2Priv()
	client := suite.NewClient(pk)
	data := &pb.TransactionData{
		Amount: "1",
	}
	payload, err := data.Marshal()
	suite.Require().Nil(err)
	tx := &pb.BxhTransaction{
		From:      from,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}
	nonce := atomic.LoadUint64(&nonce2)
	_, err = client.SendTransaction(tx, &rpcx.TransactOpts{
		Nonce: nonce,
	})
	suite.Require().NotNil(err)
}

//tc：发送交易，from、to为空，交易发送失败
func (suite *Model8) Test0803_TXEmptyFromAndToIsFail() {
	pk, _, err := repo.Node2Priv()
	client := suite.NewClient(pk)
	data := &pb.TransactionData{
		Amount: "1",
	}
	payload, err := data.Marshal()
	suite.Require().Nil(err)

	tx := &pb.BxhTransaction{
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}
	nonce := atomic.LoadUint64(&nonce2)
	_, err = client.SendTransaction(tx, &rpcx.TransactOpts{
		Nonce: nonce,
	})
	suite.Require().NotNil(err)
}

//tc：发送交易，from、to相同，交易发送失败
func (suite *Model8) Test0804_TXSameFromAndToIsFail() {
	pk, from, err := repo.Node2Priv()
	client := suite.NewClient(pk)
	data := &pb.TransactionData{
		Amount: "1",
	}
	payload, err := data.Marshal()
	suite.Require().Nil(err)
	tx := &pb.BxhTransaction{
		From:      from,
		To:        from,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}
	nonce := atomic.LoadUint64(&nonce2)
	_, err = client.SendTransaction(tx, &rpcx.TransactOpts{
		Nonce: nonce,
	})
	suite.Require().NotNil(err)
}

//tc：发送交易，签名非法，交易发送失败
func (suite *Model8) Test0806_TXWrongSigPrivateKeyIsFail() {
	_, from, err := repo.Node2Priv()
	suite.Require().Nil(err)
	_, to, err := repo.KeyPriv()
	suite.Require().Nil(err)
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
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	_, err = client.SendTransaction(tx, nil)
	suite.Require().NotNil(err)
}

//tc：发送交易，交易超过10MB，交易发送失败
func (suite *Model8) Test0808_TXExtra10MBIsFail() {
	pk, from, err := repo.Node2Priv()
	suite.Require().Nil(err)
	_, to, err := repo.KeyPriv()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	MB10 := make([]byte, 10*1024*1024) // 10MB
	for i := 0; i < len(MB10); i++ {
		MB10[i] = uint8(rand.Intn(255))
	}
	data := &pb.TransactionData{
		Amount: "1",
	}
	payload, err := data.Marshal()
	suite.Require().Nil(err)
	tx := &pb.BxhTransaction{
		From:      from,
		To:        to,
		Timestamp: time.Now().UnixNano(),
		Extra:     MB10,
		Payload:   payload,
	}
	nonce := atomic.LoadUint64(&nonce2)
	_, err = client.SendTransaction(tx, &rpcx.TransactOpts{
		Nonce: nonce,
	})
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "larger than max")
}

//tc：根据正确交易hash获取交易，交易获取成功
func (suite *Model8) Test0809_GetTxByHashIsSuccess() {
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

//tc：根据错误交易hash获取交易，交易获取失败
func (suite Model8) Test0810_GetTxByWrongHashIsFail() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	_, err = client.GetTransaction("0xc7F999b83Af6DF9e67d0a37Ee7e900bF38b3D014")
	suite.Require().NotNil(err)
}

//tc：根据正确交易hash获取回执，回执获取成功
func (suite *Model8) Test0810_GetReceiptByHashIsSuccess() {
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

//tc：根据错误交易hash获取回执，回执获取失败
func (suite *Model8) Test0811_GetReceiptByWrongHash() {
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
	hash = hash[0:len(hash)-5] + "12345"
	_, err = client.GetReceipt(hash)
	suite.Require().NotNil(err)
}
