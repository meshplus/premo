package bxh_tester

import (
	"math/rand"
	"strings"
	"time"

	"github.com/meshplus/bitxhub-kit/types"
	"github.com/meshplus/bitxhub-model/pb"
	"github.com/tidwall/gjson"
)

func (suite *Snake) TestTransferLessThanAmount() {
	res, err := suite.client.GetAccountBalance(suite.from.String())
	suite.Nil(err)

	balance := gjson.Get(string(res.Data), "balance").Uint()
	suite.Nil(err)
	amount := balance + 1

	tx := &pb.Transaction{
		From: suite.from,
		To:   suite.to,
		Data: &pb.TransactionData{
			Amount: amount,
		},
		Timestamp: time.Now().UnixNano(),
		Nonce:     rand.Int63(),
	}

	err = tx.Sign(suite.pk)
	suite.Nil(err)

	hash, err := suite.client.SendTransaction(tx)
	suite.Nil(err)

	ret, err := suite.client.GetReceipt(hash)
	suite.NotNil(ret)
	suite.True(ret.Status == pb.Receipt_FAILED)
	suite.Equal(tx.Hash().String(), ret.TxHash.String())
	suite.True(strings.Contains(string(ret.Ret), "not sufficient funds"))
}

func (suite *Snake) TestToAddressIs0X000___000() {
	to := "0x0000000000000000000000000000000000000000"
	tx := &pb.Transaction{
		From: suite.from,
		To:   types.Bytes2Address([]byte(to)),
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

func (suite *Snake) TestTypeIsXVM() {
	tx := &pb.Transaction{
		From: suite.from,
		To:   suite.to,
		Data: &pb.TransactionData{
			VmType: pb.TransactionData_XVM,
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

func (suite *Snake) TestTransfer() {
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
