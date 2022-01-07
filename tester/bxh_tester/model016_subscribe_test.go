package bxh_tester

import (
	"context"
	"fmt"

	"github.com/gobuffalo/packr"
	"github.com/meshplus/bitxhub-model/pb"
	"github.com/meshplus/premo/internal/repo"
)

type Model16 struct {
	*Snake
}

//tc:订阅普通交易
func (suite Model16) Test1601_SubscribeAuditInfoWithNoIBTPIsSuccess() {
	pk1, from1, address, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk1, from1, address)
	suite.Require().Nil(err)
	pid, err := suite.CreatePid()
	suite.Require().Nil(err)
	pk2, from2, _, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterNode(pid, 0, from2, "nvpNode", from1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk2)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c, err := client.SubscribeAudit(ctx, pb.AuditSubscriptionRequest_AUDIT_NODE, 1, nil)
	go func() {
		for i := 0; i < 3; i++ {
			err := suite.FreezeAppchain(from1)
			suite.Require().Nil(err)
			err = suite.ActivateAppchain(pk1, from1)
			suite.Require().Nil(err)
		}
	}()
	var index = 0
	for {
		select {
		case infoData, ok := <-c:
			suite.Require().Equal(ok, true)
			suite.Require().NotNil(infoData)
			auditTxInfo := infoData.(*pb.AuditTxInfo)
			if !auditTxInfo.Tx.IsIBTP() {
				data := &pb.TransactionData{}
				err = data.Unmarshal(auditTxInfo.Tx.GetPayload())
				suite.Require().Nil(err)
				suite.Require().NotNil(data)
				payload := &pb.InvokePayload{}
				err = payload.Unmarshal(data.Payload)
				suite.Require().Nil(err)
				if payload.Method == "FreezeAppchain" {
					index++
				}
				if index == 3 {
					return
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

//订阅跨链交易
func (suite Model16) Test1602_SubscribeAuditInfoWithIBTPIsSuccess() {
	pk1, err := suite.PrepareServer()
	suite.Require().Nil(err)
	from1, err := pk1.PublicKey().Address()
	suite.Require().Nil(err)
	pk2, err := suite.PrepareServer()
	suite.Require().Nil(err)
	from2, err := pk2.PublicKey().Address()
	suite.Require().Nil(err)
	box := packr.NewBox(repo.ConfigPath)
	proof, err := box.Find("proof_1.0.0_rc_complex")
	ibtp := suite.MockIBTP(1, "1356:"+from1.String()+":mychannel&transfer", "1356:"+from2.String()+":mychannel&transfer", pb.IBTP_INTERCHAIN, proof)
	payload := suite.MockContent(
		"interchainCharge",
		[][]byte{[]byte("Alice"), []byte("Alice"), []byte("1")},
	)
	pid, err := suite.CreatePid()
	suite.Require().Nil(err)
	pk3, from3, _, err := suite.DeployRule()
	suite.Require().Nil(err)
	err = suite.RegisterNode(pid, 0, from3, "nvpNode", from1.String())
	suite.Require().Nil(err)
	client := suite.NewClient(pk3)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c, err := client.SubscribeAudit(ctx, pb.AuditSubscriptionRequest_AUDIT_NODE, 1, nil)
	go func() {
		err = suite.SendInterchainTx(pk1, ibtp, payload, proof)
		suite.Require().Nil(err)
	}()
	for {
		select {
		case infoData, ok := <-c:
			suite.Require().Equal(ok, true)
			suite.Require().NotNil(infoData)
			auditTxInfo := infoData.(*pb.AuditTxInfo)
			if auditTxInfo.Tx.IsIBTP() {
				fmt.Println(auditTxInfo.Tx.IBTP)
				return
			}
		case <-ctx.Done():
			return
		}
	}
}
