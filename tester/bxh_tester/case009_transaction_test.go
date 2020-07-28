package bxh_tester

import (
	"math/rand"
	"time"

	"github.com/Rican7/retry"
	"github.com/Rican7/retry/backoff"
	"github.com/Rican7/retry/strategy"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-model/pb"
)

func (suite *Snake) TestTXEmptyFrom() {
	tx := &pb.Transaction{
		To: suite.to,
		Data: &pb.TransactionData{
			Amount: 1,
		},
		Timestamp: time.Now().UnixNano(),
		Nonce:     rand.Int63(),
	}

	err := tx.Sign(suite.pk)
	suite.Nil(err)

	_, err = suite.client.SendTransaction(tx)
	suite.NotNil(err)
}

func (suite *Snake) TestTXEmptyTo() {
	tx := &pb.Transaction{
		From: suite.from,
		Data: &pb.TransactionData{
			Amount: 1,
		},
		Timestamp: time.Now().UnixNano(),
		Nonce:     rand.Int63(),
	}

	err := tx.Sign(suite.pk)
	suite.Nil(err)

	_, err = suite.client.SendTransaction(tx)
	suite.NotNil(err)
}

func (suite *Snake) TestTXEmptySig() {
	tx := &pb.Transaction{
		From: suite.from,
		To:   suite.to,
		Data: &pb.TransactionData{
			Amount: 1,
		},
		Timestamp: time.Now().UnixNano(),
		Nonce:     rand.Int63(),
	}

	_, err := suite.client.SendTransaction(tx)
	suite.NotNil(err)
}

func (suite *Snake) TestTXWrongSigPrivateKey() {
	tx := &pb.Transaction{
		From: suite.from,
		To:   suite.to,
		Data: &pb.TransactionData{
			Amount: 1,
		},
		Timestamp: time.Now().UnixNano(),
		Nonce:     rand.Int63(),
	}

	pk1, err := asym.GenerateKey(asym.ECDSASecp256r1)
	suite.Nil(err)

	err = tx.Sign(pk1)
	suite.Nil(err)

	hash, err := suite.client.SendTransaction(tx)
	suite.Nil(err)

	ret, err := suite.client.GetReceipt(hash)
	suite.NotNil(ret)
	suite.True(ret.Status == pb.Receipt_FAILED)
}

func (suite *Snake) TestTXWrongSigAlgorithm() {
	// K1
}

func (suite *Snake) TestTXExtra10MB() {
	MB10 := make([]byte, 10*1024*1024) // 10MB
	for i := 0; i < len(MB10); i++ {
		MB10[i] = uint8(rand.Intn(255))
	}

	tx := &pb.Transaction{
		From: suite.from,
		To:   suite.to,
		Data: &pb.TransactionData{
			Amount: 1,
		},
		Timestamp: time.Now().UnixNano(),
		Nonce:     rand.Int63(),
		Extra:     MB10,
	}

	err := tx.Sign(suite.pk)
	suite.Nil(err)

	_, err = suite.client.SendTransaction(tx)
	suite.NotNil(err)
}

func (suite *Snake) TestGetTxByHash() {
	tx := &pb.Transaction{
		From: suite.from,
		To:   suite.to,
		Data: &pb.TransactionData{
			Amount: 1,
		},
		Timestamp: time.Now().UnixNano(),
		Nonce:     rand.Int63(),
	}

	err := tx.Sign(suite.pk)
	suite.Nil(err)

	hash, err := suite.client.SendTransaction(tx)
	suite.Nil(err)

	var ret *pb.GetTransactionResponse
	err1 := retry.Retry(func(attempt uint) error {
		ret, err = suite.client.GetTransaction(hash)
		if err != nil {
			return err
		}
		return nil
	},
		strategy.Limit(5),
		strategy.Backoff(backoff.Fibonacci(500*time.Millisecond)),
	)
	suite.Nil(err1)
	suite.Nil(err)
	suite.NotNil(ret)
}

func (suite *Snake) TestGetReceiptByHash() {
	tx := &pb.Transaction{
		From: suite.from,
		To:   suite.to,
		Data: &pb.TransactionData{
			Amount: 1,
		},
		Timestamp: time.Now().UnixNano(),
		Nonce:     rand.Int63(),
	}

	err := tx.Sign(suite.pk)
	suite.Nil(err)

	hash, err := suite.client.SendTransaction(tx)
	suite.Nil(err)

	ret, err := suite.client.GetReceipt(hash)
	suite.NotNil(ret)
	suite.True(ret.Status == pb.Receipt_SUCCESS)
	suite.Equal(tx.Hash().String(), ret.TxHash.String())
}
