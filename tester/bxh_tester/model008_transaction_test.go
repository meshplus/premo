package bxh_tester

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Rican7/retry"
	"github.com/Rican7/retry/backoff"
	"github.com/Rican7/retry/strategy"
	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-model/pb"
)

func (suite *Snake) Test0801_TXEmptyFrom() {
	data := &pb.TransactionData{
		Amount: 1,
	}
	payload, err := data.Marshal()
	suite.Require().Nil(err)

	tx := &pb.BxhTransaction{
		To:        suite.to,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}
	_, err = suite.client.SendTransaction(tx, nil)
	suite.Require().NotNil(err)
}

func (suite *Snake) Test0802_TXEmptyTo() {
	data := &pb.TransactionData{
		Amount: 1,
	}
	payload, err := data.Marshal()
	suite.Require().Nil(err)

	tx := &pb.BxhTransaction{
		From:      suite.from,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}
	_, err = suite.client.SendTransaction(tx, nil)
	suite.Require().NotNil(err)
}

/*增加form和to都为空*/
func (suite *Snake) Test0803_TXEmptyFromAndTo() {
	data := &pb.TransactionData{
		Amount: 1,
	}
	payload, err := data.Marshal()
	suite.Require().Nil(err)

	tx := &pb.BxhTransaction{
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}
	_, err = suite.client.SendTransaction(tx, nil)
	suite.Require().NotNil(err)
}

/*增加from和to相同*/
func (suite *Snake) Test0804_TXSameFromAndTo() {
	data := &pb.TransactionData{
		Amount: 1,
	}
	payload, err := data.Marshal()
	suite.Require().Nil(err)

	tx := &pb.BxhTransaction{
		From:      suite.from,
		To:        suite.from,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}
	_, err = suite.client.SendTransaction(tx, nil)
	suite.Require().NotNil(err)
}

func (suite *Snake) Test0805_TXEmptySig() {
	data := &pb.TransactionData{
		Amount: 1,
	}
	payload, err := data.Marshal()
	suite.Require().Nil(err)

	tx := &pb.BxhTransaction{
		From:      suite.from,
		To:        suite.to,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}
	_, err = suite.client.SendTransaction(tx, nil)
	suite.Require().Nil(err)
}

func (suite *Snake) Test0806_TXWrongSigPrivateKey() {
	data := &pb.TransactionData{
		Amount: 1,
	}
	payload, err := data.Marshal()
	suite.Require().Nil(err)

	tx := &pb.BxhTransaction{
		From:      suite.from,
		To:        suite.to,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}

	pk1, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)

	client := suite.NewClient(pk1)

	hash, err := client.SendTransaction(tx, nil)
	suite.Require().NotNil(err)

	_, err = suite.client.GetReceipt(hash)
	suite.Require().NotNil(err)

}

func (suite *Snake) Test0807_TXWrongSigAlgorithm() {
	// K1
}

func (suite *Snake) Test0808_TXExtra10MB() {
	MB10 := make([]byte, 10*1024*1024) // 10MB
	for i := 0; i < len(MB10); i++ {
		MB10[i] = uint8(rand.Intn(255))
	}

	data := &pb.TransactionData{
		Amount: 1,
	}
	payload, err := data.Marshal()
	suite.Require().Nil(err)

	tx := &pb.BxhTransaction{
		From:      suite.from,
		To:        suite.to,
		Timestamp: time.Now().UnixNano(),
		Extra:     MB10,
		Payload:   payload,
	}

	_, err = suite.client.SendTransaction(tx, nil)
	suite.Require().NotNil(err)
}

func (suite *Snake) Test0809_GetTxByHash() {
	data := &pb.TransactionData{
		Amount: 1,
	}
	payload, err := data.Marshal()
	suite.Require().Nil(err)

	tx := &pb.BxhTransaction{
		From:      suite.from,
		To:        suite.to,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}

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

func (suite *Snake) Test0810_GetReceiptByHash() {
	data := &pb.TransactionData{
		Amount: 1,
	}
	payload, err := data.Marshal()
	suite.Require().Nil(err)

	tx := &pb.BxhTransaction{
		From:      suite.from,
		To:        suite.to,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}

	hash, err := suite.client.SendTransaction(tx, nil)
	suite.Require().Nil(err)
	fmt.Println(hash)

	ret, err := suite.client.GetReceipt(hash)
	suite.Require().Nil(err)
	suite.Require().NotNil(ret)
	suite.Require().True(ret.IsSuccess())
	suite.Require().Equal(tx.Hash().String(), ret.TxHash.String())
}

/*通过错误的hash值进行查询*/
func (suite *Snake) Test0811_GetReceiptByWrongHash() {
	data := &pb.TransactionData{
		Amount: 1,
	}
	payload, err := data.Marshal()
	suite.Require().Nil(err)

	tx := &pb.BxhTransaction{
		From:      suite.from,
		To:        suite.to,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}

	hash, err := suite.client.SendTransaction(tx, nil)
	suite.Require().Nil(err)

	hash = hash[0:len(hash)-5] + "12345"
	ret, err := suite.client.GetReceipt(hash)
	suite.Require().Nil(ret)
}
