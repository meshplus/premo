package bxh_tester

import (
	"strings"
	"time"

	"github.com/meshplus/bitxhub-kit/types"
	"github.com/meshplus/bitxhub-model/pb"
	"github.com/tidwall/gjson"
)

func (suite *Snake) TestTransferLessThanAmount() {
	res, err := suite.client.GetAccountBalance(suite.from.String())
	suite.Require().Nil(err)

	balance := gjson.Get(string(res.Data), "balance").Uint()
	suite.Require().Nil(err)
	amount := balance + 1

	data := &pb.TransactionData{
		Amount: amount,
	}
	payload, err := data.Marshal()
	suite.Require().Nil(err)

	tx := &pb.Transaction{
		From:      suite.from,
		To:        suite.to,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}

	hash, err := suite.client.SendTransaction(tx, nil)
	suite.Require().Nil(err)

	ret, err := suite.client.GetReceipt(hash)
	suite.Require().NotNil(ret)
	suite.Require().True(ret.Status == pb.Receipt_FAILED)
	suite.Require().Equal(tx.Hash().String(), ret.TxHash.String())
	suite.Require().True(strings.Contains(string(ret.Ret), "not sufficient funds"))
}

func (suite *Snake) TestToAddressIs0X000___000() {
	to := "0x0000000000000000000000000000000000000000"

	data := &pb.TransactionData{
		Amount: 1,
	}
	payload, err := data.Marshal()
	suite.Require().Nil(err)
	tx := &pb.Transaction{
		From:      suite.from,
		To:        types.NewAddress([]byte(to)),
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}

	hash, err := suite.client.SendTransaction(tx, nil)
	suite.Require().Nil(err)

	ret, err := suite.client.GetReceipt(hash)
	suite.Require().NotNil(ret)
	suite.Require().True(ret.IsSuccess())
	suite.Require().Equal(tx.Hash().String(), ret.TxHash.String())
}

func (suite *Snake) TestTypeIsXVM() {
	data := &pb.TransactionData{
		VmType: pb.TransactionData_XVM,
		Amount: 1,
	}
	payload, err := data.Marshal()
	suite.Require().Nil(err)
	tx := &pb.Transaction{
		From:      suite.from,
		To:        suite.to,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}

	hash, err := suite.client.SendTransaction(tx, nil)
	suite.Require().Nil(err)

	ret, err := suite.client.GetReceipt(hash)
	suite.Require().NotNil(ret)
	suite.Require().True(ret.IsSuccess())
	suite.Require().Equal(tx.Hash().String(), ret.TxHash.String())
}

func (suite *Snake) TestTransfer() {
	data := &pb.TransactionData{
		Amount: 1,
	}
	payload, err := data.Marshal()
	suite.Require().Nil(err)
	tx := &pb.Transaction{
		From:      suite.from,
		To:        suite.to,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}

	hash, err := suite.client.SendTransaction(tx, nil)
	suite.Require().Nil(err)

	ret, err := suite.client.GetReceipt(hash)
	suite.Require().NotNil(ret)
	suite.Require().True(ret.IsSuccess())
	suite.Require().Equal(tx.Hash().String(), ret.TxHash.String())
}
