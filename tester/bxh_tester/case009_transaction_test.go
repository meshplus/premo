package bxh_tester

import (
	"math/rand"
	"time"

	"github.com/Rican7/retry"
	"github.com/Rican7/retry/backoff"
	"github.com/Rican7/retry/strategy"
	"github.com/meshplus/bitxhub-kit/crypto"
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
	}

	err := tx.Sign(suite.pk)
	suite.Require().Nil(err)

	_, err = suite.client.SendTransaction(tx, nil)
	suite.Require().NotNil(err)
}

func (suite *Snake) TestTXEmptyTo() {
	tx := &pb.Transaction{
		From: suite.from,
		Data: &pb.TransactionData{
			Amount: 1,
		},
		Timestamp: time.Now().UnixNano(),
	}

	err := tx.Sign(suite.pk)
	suite.Require().Nil(err)

	_, err = suite.client.SendTransaction(tx, nil)
	suite.Require().NotNil(err)
}

/*增加form和to都为空*/
func (suite *Snake) TestTXEmptyFromAndTo() {
	tx := &pb.Transaction{
		Data: &pb.TransactionData{
			Amount: 1,
		},
		Timestamp: time.Now().UnixNano(),
	}
	err := tx.Sign(suite.pk)
	suite.Require().Nil(err)
	_, err = suite.client.SendTransaction(tx, nil)
	suite.Require().NotNil(err)
}

/*增加from和to相同*/
func (suite *Snake) TestTXSameFromAndTo() {
	tx := &pb.Transaction{
		From: suite.from,
		To:   suite.from,
		Data: &pb.TransactionData{
			Amount: 1,
		},
		Timestamp: time.Now().UnixNano(),
	}
	err := tx.Sign(suite.pk)
	suite.Require().Nil(err)

	_, err = suite.client.SendTransaction(tx, nil)
	suite.Require().NotNil(err)
}

func (suite *Snake) TestTXEmptySig() {
	tx := &pb.Transaction{
		From: suite.from,
		To:   suite.to,
		Data: &pb.TransactionData{
			Amount: 1,
		},
		Timestamp: time.Now().UnixNano(),
	}
	tx.Nonce = 1
	_, err := suite.client.SendTransaction(tx, nil)
	suite.Require().NotNil(err)
}

func (suite *Snake) TestTXWrongSigPrivateKey() {
	tx := &pb.Transaction{
		From: suite.from,
		To:   suite.to,
		Data: &pb.TransactionData{
			Amount: 1,
		},
		Timestamp: time.Now().UnixNano(),
	}

	pk1, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)

	err = tx.Sign(pk1)
	suite.Require().Nil(err)

	hash, err := suite.client.SendTransaction(tx, nil)
	suite.Require().Nil(err)

	res, err := suite.client.GetReceipt(hash)
	suite.Require().NotNil(res)
	//suite.Require().NotNil(err)
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
		Extra:     MB10,
	}

	err := tx.Sign(suite.pk)
	suite.Require().Nil(err)

	_, err = suite.client.SendTransaction(tx, nil)
	suite.Require().NotNil(err)
}

func (suite *Snake) TestGetTxByHash() {
	tx := &pb.Transaction{
		From: suite.from,
		To:   suite.to,
		Data: &pb.TransactionData{
			Amount: 1,
		},
		Timestamp: time.Now().UnixNano(),
	}

	err := tx.Sign(suite.pk)
	suite.Require().Nil(err)

	hash, err := suite.client.SendTransaction(tx, nil)
	suite.Require().Nil(err)

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
	suite.Require().Nil(err1)
	suite.Require().Nil(err)
	suite.Require().NotNil(ret)
}

func (suite *Snake) TestGetReceiptByHash() {
	tx := &pb.Transaction{
		From: suite.from,
		To:   suite.to,
		Data: &pb.TransactionData{
			Amount: 1,
		},
		Timestamp: time.Now().UnixNano(),
	}

	hash, err := suite.client.SendTransaction(tx, nil)
	suite.Require().Nil(err)

	ret, err := suite.client.GetReceipt(hash)
	suite.Require().NotNil(ret)
	suite.Require().True(ret.Status == pb.Receipt_SUCCESS)
	suite.Require().Equal(tx.Hash().String(), ret.TxHash.String())
}

/*通过错误的hash值进行查询*/
func (suite *Snake) TestGetReceiptByWrongHash() {
	tx := &pb.Transaction{
		From: suite.from,
		To:   suite.to,
		Data: &pb.TransactionData{
			Amount: 1,
		},
		Timestamp: time.Now().UnixNano(),
	}

	err := tx.Sign(suite.pk)
	suite.Require().Nil(err)

	hash, err := suite.client.SendTransaction(tx, nil)
	suite.Require().Nil(err)
	hash = hash + "1"
	ret, err := suite.client.GetReceipt(hash)
	suite.Require().Nil(ret)
}
