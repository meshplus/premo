package bxh_tester

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gobuffalo/packr"
	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/meshplus/premo/internal/repo"
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
	box := packr.NewBox(repo.ConfigPath)
	proof, err := box.Find("proof_fabric")
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
	box := packr.NewBox(repo.ConfigPath)
	proof, err := box.Find("proof_fabric")
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
	box := packr.NewBox(repo.ConfigPath)
	proof, err := box.Find("proof_fabric")
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
	box := packr.NewBox(repo.ConfigPath)
	proof, err := box.Find("proof_fabric")
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
	box := packr.NewBox(repo.ConfigPath)
	proof, err := box.Find("proof_fabric")
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
	box := packr.NewBox(repo.ConfigPath)
	proof, err := box.Find("proof_fabric")
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
	box := packr.NewBox(repo.ConfigPath)
	proof, err := box.Find("proof_fabric")
	ibtp1 := suite.MockIBTP(1, "1356:"+from1.String()+":mychannel&transfer", "1356:"+from2.String()+":mychannel&transfer", pb.IBTP_INTERCHAIN, proof)
	ibtp2 := suite.MockIBTP(1, "1356:"+from1.String()+":mychannel&transfer", "1356:"+from3.String()+":mychannel&transfer", pb.IBTP_INTERCHAIN, proof)
	payload := suite.MockContent(
		"interchainCharge",
		[][]byte{[]byte("Alice"), []byte("Alice"), []byte("1")},
	)
	mp1 := &pb.StringUint64Map{Keys: []string{ibtp2.To}, Vals: []uint64{1}}
	mp2 := &pb.StringUint64Map{Keys: []string{ibtp1.To}, Vals: []uint64{1}}
	ibtp1.Group = mp1
	ibtp2.Group = mp2
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
	box := packr.NewBox(repo.ConfigPath)
	proof, err := box.Find("proof_fabric")
	ibtp1 := suite.MockIBTP(1, "1356:"+from1.String()+":mychannel&transfer", "1356:"+from2.String()+":mychannel&transfer", pb.IBTP_INTERCHAIN, proof)
	ibtp2 := suite.MockIBTP(1, "1356:"+from1.String()+":mychannel&transfer", "1356:"+from3.String()+":mychannel&transfer", pb.IBTP_INTERCHAIN, proof)
	payload := suite.MockContent(
		"interchainCharge",
		[][]byte{[]byte("Alice"), []byte("Alice"), []byte("1")},
	)
	mp1 := &pb.StringUint64Map{Keys: []string{ibtp2.To}, Vals: []uint64{1}}
	mp2 := &pb.StringUint64Map{Keys: []string{ibtp1.To}, Vals: []uint64{1}}
	ibtp1.Group = mp1
	ibtp2.Group = mp2
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
	box := packr.NewBox(repo.ConfigPath)
	proof, err := box.Find("proof_fabric")
	ibtp1 := suite.MockIBTP(1, "1356:"+from1.String()+":mychannel&transfer", "1356:"+from2.String()+":mychannel&transfer", pb.IBTP_INTERCHAIN, proof)
	ibtp2 := suite.MockIBTP(1, "1356:"+from1.String()+":mychannel&transfer", "1356:"+from3.String()+":mychannel&transfer", pb.IBTP_INTERCHAIN, proof)
	payload := suite.MockContent(
		"interchainCharge",
		[][]byte{[]byte("Alice"), []byte("Alice"), []byte("1")},
	)
	mp1 := &pb.StringUint64Map{Keys: []string{ibtp2.To}, Vals: []uint64{1}}
	mp2 := &pb.StringUint64Map{Keys: []string{ibtp1.To}, Vals: []uint64{1}}
	ibtp1.Group = mp1
	ibtp2.Group = mp2
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
	box := packr.NewBox(repo.ConfigPath)
	proof, err := box.Find("proof_fabric")
	ibtp1 := suite.MockIBTP(1, "1356:"+from1.String()+":mychannel&transfer", "1356:"+from2.String()+":mychannel&transfer", pb.IBTP_INTERCHAIN, proof)
	ibtp2 := suite.MockIBTP(1, "1356:"+from1.String()+":mychannel&transfer", "1356:"+from3.String()+":mychannel&transfer", pb.IBTP_INTERCHAIN, proof)
	payload := suite.MockContent(
		"interchainCharge",
		[][]byte{[]byte("Alice"), []byte("Alice"), []byte("1")},
	)
	mp1 := &pb.StringUint64Map{Keys: []string{ibtp2.To}, Vals: []uint64{1}}
	mp2 := &pb.StringUint64Map{Keys: []string{ibtp1.To}, Vals: []uint64{1}}
	ibtp1.Group = mp1
	ibtp2.Group = mp2
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
	box := packr.NewBox(repo.ConfigPath)
	proof, err := box.Find("proof_fabric")
	ibtp1 := suite.MockIBTP(1, "1356:"+from1.String()+":mychannel&transfer", "1356:"+from2.String()+":mychannel&transfer", pb.IBTP_INTERCHAIN, proof)
	ibtp2 := suite.MockIBTP(1, "1356:"+from1.String()+":mychannel&transfer", "1356:"+from3.String()+":mychannel&transfer", pb.IBTP_INTERCHAIN, proof)
	payload := suite.MockContent(
		"interchainCharge",
		[][]byte{[]byte("Alice"), []byte("Alice"), []byte("1")},
	)
	mp1 := &pb.StringUint64Map{Keys: []string{ibtp2.To}, Vals: []uint64{1}}
	mp2 := &pb.StringUint64Map{Keys: []string{ibtp1.To}, Vals: []uint64{1}}
	ibtp1.Group = mp1
	ibtp2.Group = mp2
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
	box := packr.NewBox(repo.ConfigPath)
	proof, err := box.Find("proof_fabric")
	ibtp1 := suite.MockIBTP(1, "1356:"+from1.String()+":mychannel&transfer", "1356:"+from2.String()+":mychannel&transfer", pb.IBTP_INTERCHAIN, proof)
	ibtp2 := suite.MockIBTP(1, "1356:"+from1.String()+":mychannel&transfer", "1356:"+from3.String()+":mychannel&transfer", pb.IBTP_INTERCHAIN, proof)
	payload := suite.MockContent(
		"interchainCharge",
		[][]byte{[]byte("Alice"), []byte("Alice"), []byte("1")},
	)
	mp1 := &pb.StringUint64Map{Keys: []string{ibtp2.To}, Vals: []uint64{1}}
	mp2 := &pb.StringUint64Map{Keys: []string{ibtp1.To}, Vals: []uint64{1}}
	ibtp1.Group = mp1
	ibtp2.Group = mp2
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

// PrepareServer prepare a server and return privateKey
func (suite *Snake) PrepareServer() (crypto.PrivateKey, error) {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	from, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	err = suite.RegisterAppchainWithType(pk, "Fabric V1.4.3", HappyRuleAddr, "{\"channel_id\":\"mychannel\",\"chaincode_id\":\"broker\",\"broker_version\":\"1\"}")
	suite.Require().Nil(err)
	err = suite.RegisterServer(pk, from.String(), "mychannel&transfer", from.String(), "CallContract")
	suite.Require().Nil(err)
	return pk, nil
}

// GetStatus get tx status
func (suite *Model9) GetStatus(txId string) (pb.TransactionStatus, error) {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	tx, err := client.GenerateContractTx(pb.TransactionData_BVM, constant.TransactionMgrContractAddr.Address(), "GetStatus", rpcx.String(txId))
	suite.Require().Nil(err)
	res, err := client.SendView(tx)
	suite.Require().Nil(err)
	if res.Status != pb.Receipt_SUCCESS {
		return -1, fmt.Errorf(string(res.Ret))
	}
	status, err := strconv.ParseInt(string(res.Ret), 10, 64)
	suite.Require().Nil(err)
	return pb.TransactionStatus(status), nil
}

// SendInterchainTx send interchain tx
func (suite *Snake) SendInterchainTx(pk crypto.PrivateKey, ibtp *pb.IBTP, payload, proof []byte) error {
	ibtp.Payload = payload
	from, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	tx := &pb.BxhTransaction{
		From:      from,
		To:        constant.InterchainContractAddr.Address(),
		Timestamp: time.Now().UnixNano(),
		Extra:     proof,
		IBTP:      ibtp,
	}
	client := suite.NewClient(pk)
	res, err := client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	if res.Status == pb.Receipt_FAILED {
		return fmt.Errorf(string(res.Ret))
	}
	return nil
}
