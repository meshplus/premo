package bxh_tester

import (
	"sync/atomic"
	"time"

	"github.com/meshplus/bitxhub-kit/types"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/meshplus/premo/internal/repo"
)

type Model2 struct {
	*Snake
}

//tc：发送转账交易，to为0x0000000000000000000000000000000000000000，交易回执显示成功
func (suite *Model2) Test0201_ToAddressIs0X000IsSuccess() {
	pk, from, err := repo.Node2Priv()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	data := &pb.TransactionData{
		Amount: "1",
	}
	payload, err := data.Marshal()
	suite.Require().Nil(err)
	tx := &pb.BxhTransaction{
		From:      from,
		To:        types.NewAddressByStr("0x0000000000000000000000000000000000000000"),
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}

	ret, err := client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		Nonce: atomic.AddUint64(&nonce2, 1),
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, ret.Status)
}

//tc：发送转账交易，正常情况发送，交易回执状态显示成功，对应from和to地址金额相对应变化
func (suite *Model2) Test0202_TransferIsSuccess() {
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
	ret, err := client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		Nonce: atomic.AddUint64(&nonce2, 1),
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, ret.Status)
}
