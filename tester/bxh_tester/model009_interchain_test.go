package bxh_tester

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/gobuffalo/packr/v2"
	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/meshplus/premo/internal/repo"
)

const (
	normalTxGasFee = 1050000000
	ibtpGasFee     = normalTxGasFee * 10
	leftFee        = 10
	gasError       = "insufficient balance"
)

type Model9 struct {
	*Snake
}

//tc：中继链收到正常ibtp后，事务状态为TransactionStatus_BEGIN
func (suite *Model9) Test0901_SendIBTPWithStatusTransactionStatus_BEGIN() {
	pk1, err := suite.PrepareServer()
	suite.Require().Nil(err)
	from1, err := pk1.PublicKey().Address()
	suite.Require().Nil(err)
	pk2, err := suite.PrepareServer()
	suite.Require().Nil(err)
	from2, err := pk2.PublicKey().Address()
	suite.Require().Nil(err)
	box := packr.New(repo.ConfigPath, repo.ConfigPath)
	proof, err := box.Find("proof_fabric")
	suite.Require().Nil(err)
	ibtp := suite.MockIBTP(1, "1356:"+from1.String()+":mychannel&transfer", "1356:"+from2.String()+":mychannel&transfer", pb.IBTP_INTERCHAIN, proof)
	payload := suite.MockContent(
		"interchainCharge",
		[][]byte{[]byte("Alice"), []byte("Alice"), []byte("1")},
	)
	err = suite.SendInterchainTx(pk1, ibtp, payload, proof)
	suite.Require().Nil(err)
	status, err := suite.GetStatus(ibtp.ID())
	suite.Require().Nil(err)
	suite.Require().Equal(pb.TransactionStatus_BEGIN, status)
}

//tc：中继链未在事务超时块高前收到目的链回执，事务状态为TransactionStatus_BEGIN_ROLLBACK
func (suite *Model9) Test0902_GetNoReceiptBeforeTimeOutWithStatusTransactionStatus_BEGIN_ROLLBACK() {
	pk1, err := suite.PrepareServer()
	suite.Require().Nil(err)
	from1, err := pk1.PublicKey().Address()
	suite.Require().Nil(err)
	pk2, err := suite.PrepareServer()
	suite.Require().Nil(err)
	from2, err := pk2.PublicKey().Address()
	suite.Require().Nil(err)
	box := packr.New(repo.ConfigPath, repo.ConfigPath)
	proof, err := box.Find("proof_fabric")
	suite.Require().Nil(err)
	ibtp := suite.MockIBTP(1, "1356:"+from1.String()+":mychannel&transfer", "1356:"+from2.String()+":mychannel&transfer", pb.IBTP_INTERCHAIN, proof)
	payload := suite.MockContent(
		"interchainCharge",
		[][]byte{[]byte("Alice"), []byte("Alice"), []byte("1")},
	)
	err = suite.SendInterchainTx(pk1, ibtp, payload, proof)
	suite.Require().Nil(err)
	for i := 0; i < 15; i++ {
		suite.SendTransaction(pk1)
	}
	status, err := suite.GetStatus(ibtp.ID())
	suite.Require().Nil(err)
	suite.Require().Equal(pb.TransactionStatus_BEGIN_ROLLBACK, status)
}

//tc：中继链收到目的链或服务等不存在IBTP后，事务状态为TransactionStatus_BEGIN_FAIL
func (suite *Model9) Test0903_SendIBTPNoExistChainWithStatusTransactionStatus_BEGIN_FAIL() {
	pk1, err := suite.PrepareServer()
	suite.Require().Nil(err)
	from1, err := pk1.PublicKey().Address()
	suite.Require().Nil(err)
	pk2, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	from2, err := pk2.PublicKey().Address()
	suite.Require().Nil(err)
	box := packr.New(repo.ConfigPath, repo.ConfigPath)
	proof, err := box.Find("proof_fabric")
	suite.Require().Nil(err)
	ibtp := suite.MockIBTP(1, "1356:"+from1.String()+":mychannel&transfer", "1356:"+from2.String()+":mychannel&transfer", pb.IBTP_INTERCHAIN, proof)
	payload := suite.MockContent(
		"interchainCharge",
		[][]byte{[]byte("Alice"), []byte("Alice"), []byte("1")},
	)
	err = suite.SendInterchainTx(pk1, ibtp, payload, proof)
	suite.Require().Nil(err)
	status, err := suite.GetStatus(ibtp.ID())
	suite.Require().Nil(err)
	suite.Require().Equal(pb.TransactionStatus_BEGIN_FAILURE, status)
}

//tc：中继链收到成功的跨链交易回执后，事务状态为TransactionStatus_SUCCESS
func (suite *Model9) Test0904_GetReceiptSuccessWithStatusTransactionStatus_SUCCESS() {
	pk1, err := suite.PrepareServer()
	suite.Require().Nil(err)
	from1, err := pk1.PublicKey().Address()
	suite.Require().Nil(err)
	pk2, err := suite.PrepareServer()
	suite.Require().Nil(err)
	from2, err := pk2.PublicKey().Address()
	suite.Require().Nil(err)
	box := packr.New(repo.ConfigPath, repo.ConfigPath)
	proof, err := box.Find("proof_fabric")
	suite.Require().Nil(err)
	ibtp := suite.MockIBTP(1, "1356:"+from1.String()+":mychannel&transfer", "1356:"+from2.String()+":mychannel&transfer", pb.IBTP_INTERCHAIN, proof)
	payload := suite.MockContent(
		"interchainCharge",
		[][]byte{[]byte("Alice"), []byte("Alice"), []byte("1")},
	)
	err = suite.SendInterchainTx(pk1, ibtp, payload, proof)
	suite.Require().Nil(err)
	ibtp = suite.MockIBTP(1, "1356:"+from1.String()+":mychannel&transfer", "1356:"+from2.String()+":mychannel&transfer", pb.IBTP_RECEIPT_SUCCESS, proof)
	payload = suite.MockResult([][]byte(nil))
	err = suite.SendInterchainTx(pk1, ibtp, payload, proof)
	suite.Require().Nil(err)
	status, err := suite.GetStatus(ibtp.ID())
	suite.Require().Nil(err)
	suite.Require().Equal(pb.TransactionStatus_SUCCESS, status)
}

//tc：中继链收到失败的跨链交易回执后，事务状态为TransactionStatus_FAILURE
func (suite *Model9) Test0905_GetReceiptFailWithStatusTransactionStatus_FAILURE() {
	pk1, err := suite.PrepareServer()
	suite.Require().Nil(err)
	from1, err := pk1.PublicKey().Address()
	suite.Require().Nil(err)
	pk2, err := suite.PrepareServer()
	suite.Require().Nil(err)
	from2, err := pk2.PublicKey().Address()
	suite.Require().Nil(err)
	box := packr.New(repo.ConfigPath, repo.ConfigPath)
	proof, err := box.Find("proof_fabric")
	suite.Require().Nil(err)
	ibtp := suite.MockIBTP(1, "1356:"+from1.String()+":mychannel&transfer", "1356:"+from2.String()+":mychannel&transfer", pb.IBTP_INTERCHAIN, proof)
	payload := suite.MockContent(
		"interchainCharge",
		[][]byte{[]byte("Alice"), []byte("Alice"), []byte("1")},
	)
	err = suite.SendInterchainTx(pk1, ibtp, payload, proof)
	suite.Require().Nil(err)
	ibtp = suite.MockIBTP(1, "1356:"+from1.String()+":mychannel&transfer", "1356:"+from2.String()+":mychannel&transfer", pb.IBTP_RECEIPT_FAILURE, proof)
	payload = suite.MockResult([][]byte(nil))
	err = suite.SendInterchainTx(pk1, ibtp, payload, proof)
	suite.Require().Nil(err)
	status, err := suite.GetStatus(ibtp.ID())
	suite.Require().Nil(err)
	suite.Require().Equal(pb.TransactionStatus_FAILURE, status)
}

//tc：中继链在超时块高以后收到目的链回执后，事务状态为TransactionStatus_ROLLBACK
func (suite *Model9) Test0906_GetReceiptAfterTimeOutWithStatusTransactionStatus_ROLLBACK() {
	pk1, err := suite.PrepareServer()
	suite.Require().Nil(err)
	from1, err := pk1.PublicKey().Address()
	suite.Require().Nil(err)
	pk2, err := suite.PrepareServer()
	suite.Require().Nil(err)
	from2, err := pk2.PublicKey().Address()
	suite.Require().Nil(err)
	box := packr.New(repo.ConfigPath, repo.ConfigPath)
	proof, err := box.Find("proof_fabric")
	suite.Require().Nil(err)
	ibtp := suite.MockIBTP(1, "1356:"+from1.String()+":mychannel&transfer", "1356:"+from2.String()+":mychannel&transfer", pb.IBTP_INTERCHAIN, proof)
	payload := suite.MockContent(
		"interchainCharge",
		[][]byte{[]byte("Alice"), []byte("Alice"), []byte("1")},
	)
	err = suite.SendInterchainTx(pk1, ibtp, payload, proof)
	suite.Require().Nil(err)
	for i := 0; i < 15; i++ {
		suite.SendTransaction(pk1)
	}
	status, err := suite.GetStatus(ibtp.ID())
	suite.Require().Nil(err)
	suite.Require().Equal(pb.TransactionStatus_BEGIN_ROLLBACK, status)
	ibtp = suite.MockIBTP(1, "1356:"+from1.String()+":mychannel&transfer", "1356:"+from2.String()+":mychannel&transfer", pb.IBTP_RECEIPT_ROLLBACK, proof)
	payload = suite.MockResult([][]byte(nil))
	err = suite.SendInterchainTx(pk1, ibtp, payload, proof)
	suite.Require().Nil(err)
	status, err = suite.GetStatus(ibtp.ID())
	suite.Require().Nil(err)
	suite.Require().Equal(pb.TransactionStatus_ROLLBACK, status)
}

//tc：中继链一对多场景下收到正常ibtp后，事务状态为TransactionStatus_BEGIN
func (suite *Model9) Test0907_SendIBTPSWithStatusTransactionStatus_BEGIN() {
	pk1, err := suite.PrepareServer()
	suite.Require().Nil(err)
	from1, err := pk1.PublicKey().Address()
	suite.Require().Nil(err)
	pk2, err := suite.PrepareServer()
	suite.Require().Nil(err)
	from2, err := pk2.PublicKey().Address()
	suite.Require().Nil(err)
	pk3, err := suite.PrepareServer()
	suite.Require().Nil(err)
	from3, err := pk3.PublicKey().Address()
	suite.Require().Nil(err)
	box := packr.New(repo.ConfigPath, repo.ConfigPath)
	proof, err := box.Find("proof_fabric")
	suite.Require().Nil(err)
	ibtp1 := suite.MockIBTP(1, "1356:"+from1.String()+":mychannel&transfer", "1356:"+from2.String()+":mychannel&transfer", pb.IBTP_INTERCHAIN, proof)
	ibtp2 := suite.MockIBTP(1, "1356:"+from1.String()+":mychannel&transfer", "1356:"+from3.String()+":mychannel&transfer", pb.IBTP_INTERCHAIN, proof)
	payload := suite.MockContent(
		"interchainCharge",
		[][]byte{[]byte("Alice"), []byte("Alice"), []byte("1")},
	)
	mp1 := &pb.StringUint64Map{Keys: []string{ibtp1.To, ibtp2.To}, Vals: []uint64{1, 1}}
	ibtp1.Group = mp1
	ibtp2.Group = mp1
	err = suite.SendInterchainTx(pk1, ibtp1, payload, proof)
	suite.Require().Nil(err)
	err = suite.SendInterchainTx(pk2, ibtp2, payload, proof)
	suite.Require().Nil(err)
	status, err := suite.GetStatus(ibtp1.ID())
	suite.Require().Nil(err)
	suite.Require().Equal(pb.TransactionStatus_BEGIN, status)
	status, err = suite.GetStatus(ibtp2.ID())
	suite.Require().Nil(err)
	suite.Require().Equal(pb.TransactionStatus_BEGIN, status)
}

//tc：中继链一对多场景下收到一个receipt_failure，事务状态为TransactionStatus_BEGIN_FAIL
func (suite *Model9) Test0908_GetOneReceiptFailWithStatusTransactionStatus_BEGIN_FAIL() {
	pk1, err := suite.PrepareServer()
	suite.Require().Nil(err)
	from1, err := pk1.PublicKey().Address()
	suite.Require().Nil(err)
	pk2, err := suite.PrepareServer()
	suite.Require().Nil(err)
	from2, err := pk2.PublicKey().Address()
	suite.Require().Nil(err)
	pk3, err := suite.PrepareServer()
	suite.Require().Nil(err)
	from3, err := pk3.PublicKey().Address()
	suite.Require().Nil(err)
	box := packr.New(repo.ConfigPath, repo.ConfigPath)
	proof, err := box.Find("proof_fabric")
	suite.Require().Nil(err)
	ibtp1 := suite.MockIBTP(1, "1356:"+from1.String()+":mychannel&transfer", "1356:"+from2.String()+":mychannel&transfer", pb.IBTP_INTERCHAIN, proof)
	ibtp2 := suite.MockIBTP(1, "1356:"+from1.String()+":mychannel&transfer", "1356:"+from3.String()+":mychannel&transfer", pb.IBTP_INTERCHAIN, proof)
	payload := suite.MockContent(
		"interchainCharge",
		[][]byte{[]byte("Alice"), []byte("Alice"), []byte("1")},
	)
	mp1 := &pb.StringUint64Map{Keys: []string{ibtp1.To, ibtp2.To}, Vals: []uint64{1, 1}}
	ibtp1.Group = mp1
	ibtp2.Group = mp1
	err = suite.SendInterchainTx(pk1, ibtp1, payload, proof)
	suite.Require().Nil(err)
	err = suite.SendInterchainTx(pk2, ibtp2, payload, proof)
	suite.Require().Nil(err)
	status, err := suite.GetStatus(ibtp1.ID())
	suite.Require().Nil(err)
	suite.Require().Equal(pb.TransactionStatus_BEGIN, status)
	status, err = suite.GetStatus(ibtp2.ID())
	suite.Require().Nil(err)
	suite.Require().Equal(pb.TransactionStatus_BEGIN, status)
	ibtp1 = suite.MockIBTP(1, "1356:"+from1.String()+":mychannel&transfer", "1356:"+from2.String()+":mychannel&transfer", pb.IBTP_RECEIPT_FAILURE, proof)
	payload = suite.MockResult([][]byte(nil))
	err = suite.SendInterchainTx(pk1, ibtp1, payload, proof)
	suite.Require().Nil(err)
	status, err = suite.GetStatus(ibtp1.ID())
	suite.Require().Nil(err)
	suite.Require().Equal(pb.TransactionStatus_BEGIN_FAILURE, status)
	status, err = suite.GetStatus(ibtp2.ID())
	suite.Require().Nil(err)
	suite.Require().Equal(pb.TransactionStatus_BEGIN_FAILURE, status)
}

//tc：中继链一对多场景下超过超时块高，事务状态为TransactionStatus_BEGIN_ROLLBACK
func (suite *Model9) Test0909_GetNoAllReceiptBeforeTimeOutWithStatusTransactionStatus_BEGIN_ROLLBACK() {
	pk1, err := suite.PrepareServer()
	suite.Require().Nil(err)
	from1, err := pk1.PublicKey().Address()
	suite.Require().Nil(err)
	pk2, err := suite.PrepareServer()
	suite.Require().Nil(err)
	from2, err := pk2.PublicKey().Address()
	suite.Require().Nil(err)
	pk3, err := suite.PrepareServer()
	suite.Require().Nil(err)
	from3, err := pk3.PublicKey().Address()
	suite.Require().Nil(err)
	box := packr.New(repo.ConfigPath, repo.ConfigPath)
	proof, err := box.Find("proof_fabric")
	suite.Require().Nil(err)
	ibtp1 := suite.MockIBTP(1, "1356:"+from1.String()+":mychannel&transfer", "1356:"+from2.String()+":mychannel&transfer", pb.IBTP_INTERCHAIN, proof)
	ibtp2 := suite.MockIBTP(1, "1356:"+from1.String()+":mychannel&transfer", "1356:"+from3.String()+":mychannel&transfer", pb.IBTP_INTERCHAIN, proof)
	payload := suite.MockContent(
		"interchainCharge",
		[][]byte{[]byte("Alice"), []byte("Alice"), []byte("1")},
	)
	mp1 := &pb.StringUint64Map{Keys: []string{ibtp1.To, ibtp2.To}, Vals: []uint64{1, 1}}
	ibtp1.Group = mp1
	ibtp2.Group = mp1
	err = suite.SendInterchainTx(pk1, ibtp1, payload, proof)
	suite.Require().Nil(err)
	err = suite.SendInterchainTx(pk2, ibtp2, payload, proof)
	suite.Require().Nil(err)
	status, err := suite.GetStatus(ibtp1.ID())
	suite.Require().Nil(err)
	suite.Require().Equal(pb.TransactionStatus_BEGIN, status)
	status, err = suite.GetStatus(ibtp2.ID())
	suite.Require().Nil(err)
	suite.Require().Equal(pb.TransactionStatus_BEGIN, status)
	ibtp1 = suite.MockIBTP(1, "1356:"+from1.String()+":mychannel&transfer", "1356:"+from2.String()+":mychannel&transfer", pb.IBTP_RECEIPT_SUCCESS, proof)
	payload = suite.MockResult([][]byte(nil))
	err = suite.SendInterchainTx(pk1, ibtp1, payload, proof)
	suite.Require().Nil(err)
	for i := 0; i < 15; i++ {
		suite.SendTransaction(pk1)
	}
	status, err = suite.GetStatus(ibtp1.ID())
	suite.Require().Nil(err)
	suite.Require().Equal(pb.TransactionStatus_BEGIN_ROLLBACK, status)
	status, err = suite.GetStatus(ibtp2.ID())
	suite.Require().Nil(err)
	suite.Require().Equal(pb.TransactionStatus_BEGIN_ROLLBACK, status)
}

//tc：中继链收到所有回执为成功，事务状态为TransactionStatus_SUCCESS
func (suite *Model9) Test0910_GetAllReceiptSuccessWithStatusTransactionStatus_SUCCESS() {
	pk1, err := suite.PrepareServer()
	suite.Require().Nil(err)
	from1, err := pk1.PublicKey().Address()
	suite.Require().Nil(err)
	pk2, err := suite.PrepareServer()
	suite.Require().Nil(err)
	from2, err := pk2.PublicKey().Address()
	suite.Require().Nil(err)
	pk3, err := suite.PrepareServer()
	suite.Require().Nil(err)
	from3, err := pk3.PublicKey().Address()
	suite.Require().Nil(err)
	box := packr.New(repo.ConfigPath, repo.ConfigPath)
	proof, err := box.Find("proof_fabric")
	suite.Require().Nil(err)
	ibtp1 := suite.MockIBTP(1, "1356:"+from1.String()+":mychannel&transfer", "1356:"+from2.String()+":mychannel&transfer", pb.IBTP_INTERCHAIN, proof)
	ibtp2 := suite.MockIBTP(1, "1356:"+from1.String()+":mychannel&transfer", "1356:"+from3.String()+":mychannel&transfer", pb.IBTP_INTERCHAIN, proof)
	payload := suite.MockContent(
		"interchainCharge",
		[][]byte{[]byte("Alice"), []byte("Alice"), []byte("1")},
	)
	mp1 := &pb.StringUint64Map{Keys: []string{ibtp1.To, ibtp2.To}, Vals: []uint64{1, 1}}
	ibtp1.Group = mp1
	ibtp2.Group = mp1
	err = suite.SendInterchainTx(pk1, ibtp1, payload, proof)
	suite.Require().Nil(err)
	err = suite.SendInterchainTx(pk2, ibtp2, payload, proof)
	suite.Require().Nil(err)
	status, err := suite.GetStatus(ibtp1.ID())
	suite.Require().Nil(err)
	suite.Require().Equal(pb.TransactionStatus_BEGIN, status)
	status, err = suite.GetStatus(ibtp2.ID())
	suite.Require().Nil(err)
	suite.Require().Equal(pb.TransactionStatus_BEGIN, status)
	ibtp1 = suite.MockIBTP(1, "1356:"+from1.String()+":mychannel&transfer", "1356:"+from2.String()+":mychannel&transfer", pb.IBTP_RECEIPT_SUCCESS, proof)
	payload = suite.MockResult([][]byte(nil))
	err = suite.SendInterchainTx(pk1, ibtp1, payload, proof)
	suite.Require().Nil(err)
	ibtp2 = suite.MockIBTP(1, "1356:"+from1.String()+":mychannel&transfer", "1356:"+from3.String()+":mychannel&transfer", pb.IBTP_RECEIPT_SUCCESS, proof)
	payload = suite.MockResult([][]byte(nil))
	err = suite.SendInterchainTx(pk1, ibtp2, payload, proof)
	suite.Require().Nil(err)
	status, err = suite.GetStatus(ibtp1.ID())
	suite.Require().Nil(err)
	suite.Require().Equal(pb.TransactionStatus_SUCCESS, status)
	status, err = suite.GetStatus(ibtp2.ID())
	suite.Require().Nil(err)
	suite.Require().Equal(pb.TransactionStatus_SUCCESS, status)
}

//tc：中继链一对多场景收到所有失败的跨链交易回执后，事务状态为TransactionStatus_FAILURE
func (suite *Model9) Test0911_GetAllReceiptFailWithStatusTransactionStatus_FAILURE() {
	pk1, err := suite.PrepareServer()
	suite.Require().Nil(err)
	from1, err := pk1.PublicKey().Address()
	suite.Require().Nil(err)
	pk2, err := suite.PrepareServer()
	suite.Require().Nil(err)
	from2, err := pk2.PublicKey().Address()
	suite.Require().Nil(err)
	pk3, err := suite.PrepareServer()
	suite.Require().Nil(err)
	from3, err := pk3.PublicKey().Address()
	suite.Require().Nil(err)
	box := packr.New(repo.ConfigPath, repo.ConfigPath)
	proof, err := box.Find("proof_fabric")
	suite.Require().Nil(err)
	ibtp1 := suite.MockIBTP(1, "1356:"+from1.String()+":mychannel&transfer", "1356:"+from2.String()+":mychannel&transfer", pb.IBTP_INTERCHAIN, proof)
	ibtp2 := suite.MockIBTP(1, "1356:"+from1.String()+":mychannel&transfer", "1356:"+from3.String()+":mychannel&transfer", pb.IBTP_INTERCHAIN, proof)
	payload := suite.MockContent(
		"interchainCharge",
		[][]byte{[]byte("Alice"), []byte("Alice"), []byte("1")},
	)
	mp1 := &pb.StringUint64Map{Keys: []string{ibtp1.To, ibtp2.To}, Vals: []uint64{1, 1}}
	ibtp1.Group = mp1
	ibtp2.Group = mp1
	err = suite.SendInterchainTx(pk1, ibtp1, payload, proof)
	suite.Require().Nil(err)
	err = suite.SendInterchainTx(pk2, ibtp2, payload, proof)
	suite.Require().Nil(err)
	status, err := suite.GetStatus(ibtp1.ID())
	suite.Require().Nil(err)
	suite.Require().Equal(pb.TransactionStatus_BEGIN, status)
	status, err = suite.GetStatus(ibtp2.ID())
	suite.Require().Nil(err)
	suite.Require().Equal(pb.TransactionStatus_BEGIN, status)
	ibtp1 = suite.MockIBTP(1, "1356:"+from1.String()+":mychannel&transfer", "1356:"+from2.String()+":mychannel&transfer", pb.IBTP_RECEIPT_FAILURE, proof)
	payload = suite.MockResult([][]byte(nil))
	err = suite.SendInterchainTx(pk1, ibtp1, payload, proof)
	suite.Require().Nil(err)
	ibtp2 = suite.MockIBTP(1, "1356:"+from1.String()+":mychannel&transfer", "1356:"+from3.String()+":mychannel&transfer", pb.IBTP_RECEIPT_FAILURE, proof)
	payload = suite.MockResult([][]byte(nil))
	err = suite.SendInterchainTx(pk1, ibtp2, payload, proof)
	suite.Require().Nil(err)
	status, err = suite.GetStatus(ibtp1.ID())
	suite.Require().Nil(err)
	suite.Require().Equal(pb.TransactionStatus_FAILURE, status)
	status, err = suite.GetStatus(ibtp2.ID())
	suite.Require().Nil(err)
	suite.Require().Equal(pb.TransactionStatus_FAILURE, status)
}

//tc：中继链一对多场景超时块高以后收到所有目的链回执后，事务状态为TransactionStatus_ROLLBACK
func (suite *Model9) Test0912_GetAllReceiptTimeOutWithStatusTransactionStatus_ROLLBACK() {
	pk1, err := suite.PrepareServer()
	suite.Require().Nil(err)
	from1, err := pk1.PublicKey().Address()
	suite.Require().Nil(err)
	pk2, err := suite.PrepareServer()
	suite.Require().Nil(err)
	from2, err := pk2.PublicKey().Address()
	suite.Require().Nil(err)
	pk3, err := suite.PrepareServer()
	suite.Require().Nil(err)
	from3, err := pk3.PublicKey().Address()
	suite.Require().Nil(err)
	box := packr.New(repo.ConfigPath, repo.ConfigPath)
	proof, err := box.Find("proof_fabric")
	suite.Require().Nil(err)
	ibtp1 := suite.MockIBTP(1, "1356:"+from1.String()+":mychannel&transfer", "1356:"+from2.String()+":mychannel&transfer", pb.IBTP_INTERCHAIN, proof)
	ibtp2 := suite.MockIBTP(1, "1356:"+from1.String()+":mychannel&transfer", "1356:"+from3.String()+":mychannel&transfer", pb.IBTP_INTERCHAIN, proof)
	payload := suite.MockContent(
		"interchainCharge",
		[][]byte{[]byte("Alice"), []byte("Alice"), []byte("1")},
	)
	mp1 := &pb.StringUint64Map{Keys: []string{ibtp1.To, ibtp2.To}, Vals: []uint64{1, 1}}
	ibtp1.Group = mp1
	ibtp2.Group = mp1
	err = suite.SendInterchainTx(pk1, ibtp1, payload, proof)
	suite.Require().Nil(err)
	err = suite.SendInterchainTx(pk2, ibtp2, payload, proof)
	suite.Require().Nil(err)
	status, err := suite.GetStatus(ibtp1.ID())
	suite.Require().Nil(err)
	suite.Require().Equal(pb.TransactionStatus_BEGIN, status)
	status, err = suite.GetStatus(ibtp2.ID())
	suite.Require().Nil(err)
	suite.Require().Equal(pb.TransactionStatus_BEGIN, status)
	for i := 0; i < 15; i++ {
		suite.SendTransaction(pk1)
	}
	ibtp1 = suite.MockIBTP(1, "1356:"+from1.String()+":mychannel&transfer", "1356:"+from2.String()+":mychannel&transfer", pb.IBTP_RECEIPT_ROLLBACK, proof)
	payload = suite.MockResult([][]byte(nil))
	err = suite.SendInterchainTx(pk1, ibtp1, payload, proof)
	suite.Require().Nil(err)
	ibtp2 = suite.MockIBTP(1, "1356:"+from1.String()+":mychannel&transfer", "1356:"+from3.String()+":mychannel&transfer", pb.IBTP_RECEIPT_ROLLBACK, proof)
	payload = suite.MockResult([][]byte(nil))
	err = suite.SendInterchainTx(pk1, ibtp2, payload, proof)
	suite.Require().Nil(err)
	status, err = suite.GetStatus(ibtp1.ID())
	suite.Require().Nil(err)
	suite.Require().Equal(pb.TransactionStatus_ROLLBACK, status)
	status, err = suite.GetStatus(ibtp2.ID())
	suite.Require().Nil(err)
	suite.Require().Equal(pb.TransactionStatus_ROLLBACK, status)
}

//tc：中继链因为gas费不足，需要回滚跨链交易
func (suite *Model9) Test0913_AccountHaveInsufficientGas() {
	pk1, err := suite.PrepareServer()
	suite.Require().Nil(err)
	from1, err := pk1.PublicKey().Address()
	suite.Require().Nil(err)
	pk2, err := suite.PrepareServer()
	suite.Require().Nil(err)
	from2, err := pk2.PublicKey().Address()
	suite.Require().Nil(err)
	client1 := suite.NewClient(pk1)
	res, err := client1.GetAccountBalance(from1.String())
	suite.Require().Nil(err)
	account := Account{}
	err = json.Unmarshal(res.Data, &account)
	suite.Require().Nil(err)
	fmt.Println(from1.String(), account.Balance)

	// left 3 tx balance, first tx for deducting balance,
	// second tx for sending success ibtp1,
	// third tx for getting interchain,
	// if send fourth tx for sending ibtp2, receive failed receipt for insufficeient balance(left gas)
	gasUsed := normalTxGasFee + ibtpGasFee*2
	balance := new(big.Int).Sub(big.NewInt(account.Balance.Int64()), big.NewInt(int64(gasUsed)+leftFee))

	data := &pb.TransactionData{
		Amount: balance.String(),
	}
	payload, err := data.Marshal()
	suite.Require().Nil(err)
	tx := &pb.BxhTransaction{
		From:      from1,
		To:        from2,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}
	receipt, err := client1.SendTransactionWithReceipt(tx, nil)
	suite.Require().True(receipt.IsSuccess())
	suite.Require().Nil(err)

	// mock ibtp1, get successful receipt
	box := packr.New(repo.ConfigPath, repo.ConfigPath)
	proof, err := box.Find("proof_fabric")
	suite.Require().Nil(err)
	ibtp1 := suite.MockIBTP(1, "1356:"+from1.String()+":mychannel&transfer", "1356:"+from2.String()+":mychannel&transfer", pb.IBTP_INTERCHAIN, proof)
	payload = suite.MockContent(
		"interchainCharge",
		[][]byte{[]byte("Alice"), []byte("Alice"), []byte("1")},
	)
	ibtp1.Payload = payload

	tx = &pb.BxhTransaction{
		From:      from1,
		To:        constant.InterchainContractAddr.Address(),
		Timestamp: time.Now().UnixNano(),
		Extra:     proof,
		IBTP:      ibtp1,
	}
	receipt, err = client1.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().True(receipt.IsSuccess())

	// get interchain to proof ibtp store bitxhub ledger successfully
	receipt1, err := client1.InvokeBVMContract(constant.InterchainContractAddr.Address(), "GetInterchain", nil, pb.String(ibtp1.From))
	suite.Require().Nil(err)
	interchain := &pb.Interchain{}
	err = interchain.Unmarshal(receipt1.Ret)
	suite.Require().Nil(err)
	suite.Require().True(receipt1.IsSuccess())
	suite.Require().Equal(ibtp1.From, interchain.ID)
	suite.Require().Equal(uint64(1), interchain.InterchainCounter[ibtp1.To])

	// ensure have insufficient gas
	res, err = client1.GetAccountBalance(from1.String())
	suite.Require().Nil(err)
	account = Account{}
	err = json.Unmarshal(res.Data, &account)
	suite.Require().Nil(err)
	fmt.Println(from1.String(), account.Balance)
	suite.Require().Equal(int64(leftFee), account.Balance.Int64())

	// mock ibtp2, send ibtp2 failed because insufficient gas
	ibtp2 := suite.MockIBTP(2, "1356:"+from1.String()+":mychannel&transfer", "1356:"+from2.String()+":mychannel&transfer", pb.IBTP_INTERCHAIN, proof)
	payload = suite.MockContent(
		"interchainCharge",
		[][]byte{[]byte("Alice"), []byte("Alice"), []byte("2")},
	)
	ibtp2.Payload = payload

	tx2 := &pb.BxhTransaction{
		From:      from1,
		To:        constant.InterchainContractAddr.Address(),
		Timestamp: time.Now().UnixNano(),
		Extra:     proof,
		IBTP:      ibtp2,
	}
	receipt2, err := client1.SendTransactionWithReceipt(tx2, nil)
	suite.Require().Nil(err)
	suite.Require().False(receipt2.IsSuccess())
	suite.Require().Contains(string(receipt2.Ret), gasError)

	client2 := suite.NewClient(pk2)
	res2, err := client2.GetAccountBalance(from1.String())
	suite.Require().Nil(err)

	// from1 have insufficient gas, so paying all left balance
	err = json.Unmarshal(res2.Data, &account)
	suite.Require().Nil(err)
	suite.Require().Equal(int64(0), account.Balance.Int64())

	// query interchain, bitxhub had already rollback interchain
	receipt, err = client2.InvokeBVMContract(constant.InterchainContractAddr.Address(), "GetInterchain", nil, pb.String(ibtp2.From))
	suite.Require().Nil(err)
	interchain = &pb.Interchain{}
	err = interchain.Unmarshal(receipt.Ret)
	suite.Require().Nil(err)
	suite.Require().True(receipt.IsSuccess())
	suite.Require().Equal(ibtp2.From, interchain.ID)
	suite.Require().Equal(uint64(1), interchain.InterchainCounter[ibtp2.To])
}

// PrepareServer prepare a server and return privateKey
func (suite *Snake) PrepareServer() (crypto.PrivateKey, error) {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	if err != nil {
		return nil, err
	}
	from, err := pk.PublicKey().Address()
	if err != nil {
		return nil, err
	}
	err = suite.RegisterAppchainWithType(pk, "Fabric V1.4.3", HappyRuleAddr, "{\"channel_id\":\"mychannel\",\"chaincode_id\":\"broker\",\"broker_version\":\"1\"}")
	if err != nil {
		return nil, err
	}
	err = suite.RegisterServer(pk, from.String(), "mychannel&transfer", from.String(), "CallContract")
	if err != nil {
		return nil, err
	}
	return pk, nil
}

// GetStatus get tx status
func (suite *Model9) GetStatus(txId string) (pb.TransactionStatus, error) {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	if err != nil {
		return -1, err
	}
	client := suite.NewClient(pk)
	tx, err := client.GenerateContractTx(pb.TransactionData_BVM, constant.TransactionMgrContractAddr.Address(), "GetStatus", rpcx.String(txId))
	if err != nil {
		return -1, err
	}
	res, err := client.SendView(tx)
	if err != nil {
		return -1, err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return -1, fmt.Errorf(string(res.Ret))
	}
	status, err := strconv.ParseInt(string(res.Ret), 10, 64)
	if err != nil {
		return 0, err
	}
	return pb.TransactionStatus(status), nil
}

// SendInterchainTx send interchain tx
func (suite *Snake) SendInterchainTx(pk crypto.PrivateKey, ibtp *pb.IBTP, payload, proof []byte) error {
	ibtp.Payload = payload
	from, err := pk.PublicKey().Address()
	if err != nil {
		return err
	}
	tx := &pb.BxhTransaction{
		From:      from,
		To:        constant.InterchainContractAddr.Address(),
		Timestamp: time.Now().UnixNano(),
		Extra:     proof,
		IBTP:      ibtp,
	}
	client := suite.NewClient(pk)
	res, err := client.SendTransactionWithReceipt(tx, nil)
	if err != nil {
		return err
	}
	if res.Status == pb.Receipt_FAILED {
		return fmt.Errorf(string(res.Ret))
	}
	return nil
}
