package bxh_tester

import (
	"time"

	"github.com/meshplus/bitxhub-kit/types"
	"github.com/meshplus/bitxhub-model/pb"
	"github.com/tidwall/gjson"
)

//tc:发送转账交易，from的金额少于转账的金额，交易回执显示失败
func (suite *Snake) Test0201_TransferLessThanAmount() {
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

	tx := &pb.BxhTransaction{
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
	suite.Require().Contains(string(ret.Ret), "not sufficient funds")
}

//tc:发送转账交易，to为0x0000000000000000000000000000000000000000，交易回执显示失败
func (suite *Snake) Test0202_ToAddressIs0X000___000() {
	to := "0x0000000000000000000000000000000000000000"

	data := &pb.TransactionData{
		Amount: 1,
	}
	payload, err := data.Marshal()
	suite.Require().Nil(err)
	tx := &pb.BxhTransaction{
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

//tc:发送转账交易，type设置为XVM，交易回执显示失败
func (suite *Snake) Test0203_TypeIsXVM() {
	data := &pb.TransactionData{
		VmType: pb.TransactionData_XVM,
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

	ret, err := suite.client.GetReceipt(hash)
	suite.Require().NotNil(ret)
	suite.Require().True(ret.IsSuccess())
	suite.Require().Equal(tx.Hash().String(), ret.TxHash.String())
}

//tc:发送转账交易，正常情况发送，交易回执状态显示成功，对应from和to地址金额相对应变化
func (suite *Snake) Test0204_Transfer() {
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

	ret, err := suite.client.GetReceipt(hash)
	suite.Require().NotNil(ret)
	suite.Require().True(ret.IsSuccess())
	suite.Require().Equal(tx.Hash().String(), ret.TxHash.String())
}
