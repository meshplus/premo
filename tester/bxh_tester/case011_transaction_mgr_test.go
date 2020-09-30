package bxh_tester

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"

	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
)

func (suite *TransactionMgrSuite) Test001_One2One_HappyPath() {
	index := uint64(1)
	from := suite.client0.addr
	to := suite.client1.addr
	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	ib0 := &pb.IBTP{From: from, To: to, Index: index, Type: pb.IBTP_INTERCHAIN, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}
	data0, err := ib0.Marshal()
	suite.Require().Nil(err)

	tx, _ := suite.client0.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data0))
	tx.Extra = []byte(proof)
	res0, err := suite.client0.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib0.From + ib0.To,
		IBTPNonce: ib0.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_Status(0), res0.Status, string(res0.Ret))

	ib1 := &pb.IBTP{From: from, To: to, Index: index, Type: pb.IBTP_RECEIPT_SUCCESS, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}
	data1, err := ib1.Marshal()
	suite.Require().Nil(err)

	tx, _ = suite.client1.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data1))
	tx.Extra = []byte(proof)
	res1, err := suite.client1.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib1.From + ib1.To,
		IBTPNonce: ib1.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_Status(0), res1.Status, string(res1.Ret))

	txId := fmt.Sprintf("%s-%s-%d", from, to, index)
	res, err := suite.client0.client.InvokeBVMContract(rpcx.TransactionMgrContractAddr, "GetStatus",nil, pb.String(txId))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_Status(0), res.Status, string(res.Ret))

	status, err := strconv.Atoi(string(res.Ret))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.TransactionStatus_SUCCESS, pb.TransactionStatus(status))
}

func (suite *TransactionMgrSuite) Test002_One2One_ToFail() {
	index := uint64(1)
	from := suite.client0.addr
	to := suite.client1.addr
	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	ib0 := &pb.IBTP{From: from, To: to, Index: index, Type: pb.IBTP_INTERCHAIN, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}
	data0, err := ib0.Marshal()
	suite.Require().Nil(err)

	tx, _ := suite.client0.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data0))
	tx.Extra = []byte(proof)
	res0, err := suite.client0.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib0.From + ib0.To,
		IBTPNonce: ib0.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_Status(0), res0.Status, string(res0.Ret))

	ib1 := &pb.IBTP{From: from, To: to, Index: index, Type: pb.IBTP_RECEIPT_FAILURE, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}
	data1, err := ib1.Marshal()
	suite.Require().Nil(err)

	tx, _ = suite.client1.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data1))
	tx.Extra = []byte(proof)
	res1, err := suite.client1.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib1.From + ib1.To,
		IBTPNonce: ib1.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_Status(0), res1.Status, string(res1.Ret))

	txId := fmt.Sprintf("%s-%s-%d", from, to, index)
	res, err := suite.client0.client.InvokeBVMContract(rpcx.TransactionMgrContractAddr, "GetStatus",nil, pb.String(txId))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_Status(0), res.Status, string(res.Ret))

	status, err := strconv.Atoi(string(res.Ret))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.TransactionStatus_FAILURE, pb.TransactionStatus(status))
}

func (suite *TransactionMgrSuite) Test003_One2One_Unfinished() {
	index := uint64(1)
	from := suite.client0.addr
	to := suite.client1.addr
	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	ib0 := &pb.IBTP{From: from, To: to, Index: index, Type: pb.IBTP_INTERCHAIN, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}
	data0, err := ib0.Marshal()
	suite.Require().Nil(err)

	tx, _ := suite.client0.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data0))
	tx.Extra = []byte(proof)
	res0, err := suite.client0.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib0.From + ib0.To,
		IBTPNonce: ib0.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_Status(0), res0.Status)

	txId := fmt.Sprintf("%s-%s-%d", from, to, index)
	res, err := suite.client0.client.InvokeBVMContract(rpcx.TransactionMgrContractAddr, "GetStatus",nil, pb.String(txId))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_Status(0), res.Status, string(res.Ret))

	status, err := strconv.Atoi(string(res.Ret))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.TransactionStatus_BEGIN, pb.TransactionStatus(status))
}

func (suite *TransactionMgrSuite) Test004_One2Multi_HappyPath() {
	index := uint64(1)
	from := suite.client0.addr
	to1 := suite.client1.addr
	to2 := suite.client2.addr
	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	ib0 := &pb.IBTP{From: from, To: to1, Index: index, Type: pb.IBTP_INTERCHAIN, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}
	ib1 := &pb.IBTP{From: from, To: to2, Index: index, Type: pb.IBTP_INTERCHAIN, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}

	ibtps := pb.IBTPs{}
	ibtps.Ibtps = append(ibtps.Ibtps, ib0)
	ibtps.Ibtps = append(ibtps.Ibtps, ib1)

	data0, err := ibtps.Marshal()
	suite.Require().Nil(err)

	tx, _ := suite.client0.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTPs", pb.Bytes(data0))
	tx.Extra = []byte(proof)
	res0, err := suite.client0.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib0.From + ib0.To,
		IBTPNonce: ib0.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_Status(0), res0.Status, string(res0.Ret))

	ib1 = &pb.IBTP{From: from, To: to1, Index: index, Type: pb.IBTP_RECEIPT_SUCCESS, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}
	data1, err := ib1.Marshal()
	suite.Require().Nil(err)

	time.Sleep(time.Second)

	tx, _ = suite.client1.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data1))
	tx.Extra = []byte(proof)
	res1, err := suite.client1.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib1.From + ib1.To,
		IBTPNonce: ib1.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_Status(0), res1.Status, string(res1.Ret))

	ib2 := &pb.IBTP{From: from, To: to2, Index: index, Type: pb.IBTP_RECEIPT_SUCCESS, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}
	data2, err := ib2.Marshal()
	suite.Require().Nil(err)

	tx, _ = suite.client2.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data2))
	tx.Extra = []byte(proof)
	res2, err := suite.client2.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib2.From + ib2.To,
		IBTPNonce: ib2.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_Status(0), res2.Status, string(res2.Ret))

	globalTxId := fmt.Sprintf("%s-%s", from, res0.TxHash.String())
	res, err := suite.client0.client.InvokeBVMContract(rpcx.TransactionMgrContractAddr, "GetStatus",nil, pb.String(globalTxId))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_Status(0), res.Status, string(res.Ret))

	status, err := strconv.Atoi(string(res.Ret))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.TransactionStatus_SUCCESS, pb.TransactionStatus(status))
}

func (suite *TransactionMgrSuite) Test005_One2Multi_ToFail() {
	index := uint64(1)
	from := suite.client0.addr
	to1 := suite.client1.addr
	to2 := suite.client2.addr
	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	ib0 := &pb.IBTP{From: from, To: to1, Index: index, Type: pb.IBTP_INTERCHAIN, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}
	ib1 := &pb.IBTP{From: from, To: to2, Index: index, Type: pb.IBTP_INTERCHAIN, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}

	ibtps := pb.IBTPs{}
	ibtps.Ibtps = append(ibtps.Ibtps, ib0)
	ibtps.Ibtps = append(ibtps.Ibtps, ib1)

	data0, err := ibtps.Marshal()
	suite.Require().Nil(err)

	tx, _ := suite.client0.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTPs", pb.Bytes(data0))
	tx.Extra = []byte(proof)
	res0, err := suite.client0.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib0.From + ib0.To,
		IBTPNonce: ib0.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_Status(0), res0.Status, string(res0.Ret))

	ib1 = &pb.IBTP{From: from, To: to1, Index: index, Type: pb.IBTP_RECEIPT_SUCCESS, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}
	data1, err := ib1.Marshal()
	suite.Require().Nil(err)

	time.Sleep(time.Second)

	tx, _ = suite.client1.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data1))
	tx.Extra = []byte(proof)
	res1, err := suite.client1.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib1.From + ib1.To,
		IBTPNonce: ib1.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_Status(0), res1.Status, string(res1.Ret))

	ib2 := &pb.IBTP{From: from, To: to2, Index: index, Type: pb.IBTP_RECEIPT_FAILURE, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}
	data2, err := ib2.Marshal()
	suite.Require().Nil(err)

	tx, _ = suite.client2.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data2))
	tx.Extra = []byte(proof)
	res2, err := suite.client2.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib2.From + ib2.To,
		IBTPNonce: ib2.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_Status(0), res2.Status, string(res2.Ret))

	globalTxId := fmt.Sprintf("%s-%s", from, res0.TxHash.String())
	res, err := suite.client0.client.InvokeBVMContract(rpcx.TransactionMgrContractAddr, "GetStatus",nil, pb.String(globalTxId))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_Status(0), res.Status, string(res.Ret))

	status, err := strconv.Atoi(string(res.Ret))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.TransactionStatus_FAILURE, pb.TransactionStatus(status))
}

func (suite *TransactionMgrSuite) Test006_One2Multi_Unfinished() {
	index := uint64(1)
	from := suite.client0.addr
	to1 := suite.client1.addr
	to2 := suite.client2.addr
	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	ib0 := &pb.IBTP{From: from, To: to1, Index: index, Type: pb.IBTP_INTERCHAIN, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}
	ib1 := &pb.IBTP{From: from, To: to2, Index: index, Type: pb.IBTP_INTERCHAIN, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}

	ibtps := pb.IBTPs{}
	ibtps.Ibtps = append(ibtps.Ibtps, ib0)
	ibtps.Ibtps = append(ibtps.Ibtps, ib1)

	data0, err := ibtps.Marshal()
	suite.Require().Nil(err)

	tx, _ := suite.client0.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTPs", pb.Bytes(data0))
	tx.Extra = []byte(proof)
	res0, err := suite.client0.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib0.From + ib0.To,
		IBTPNonce: ib0.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_Status(0), res0.Status, string(res0.Ret))

	ib1 = &pb.IBTP{From: from, To: to1, Index: index, Type: pb.IBTP_RECEIPT_SUCCESS, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}
	data1, err := ib1.Marshal()
	suite.Require().Nil(err)

	time.Sleep(time.Second)

	tx, _ = suite.client1.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data1))
	tx.Extra = []byte(proof)
	res1, err := suite.client1.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib1.From + ib1.To,
		IBTPNonce: ib1.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_Status(0), res1.Status, string(res1.Ret))

	globalTxId := fmt.Sprintf("%s-%s", from, res0.TxHash.String())
	res, err := suite.client0.client.InvokeBVMContract(rpcx.TransactionMgrContractAddr, "GetStatus",nil, pb.String(globalTxId))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_Status(0), res.Status, string(res.Ret))

	status, err := strconv.Atoi(string(res.Ret))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.TransactionStatus_BEGIN, pb.TransactionStatus(status))
}

func (suite *TransactionMgrSuite) Test007_One2One_FinishedStatus_Success() {
	index := uint64(1)
	from := suite.client0.addr
	to := suite.client1.addr
	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	ib0 := &pb.IBTP{From: from, To: to, Index: index, Type: pb.IBTP_INTERCHAIN, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}
	data0, err := ib0.Marshal()
	suite.Require().Nil(err)

	tx, _ := suite.client0.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data0))
	tx.Extra = []byte(proof)
	res0, err := suite.client0.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib0.From + ib0.To,
		IBTPNonce: ib0.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_Status(0), res0.Status, string(res0.Ret))

	ib1 := &pb.IBTP{From: from, To: to, Index: index, Type: pb.IBTP_RECEIPT_SUCCESS, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}
	data1, err := ib1.Marshal()
	suite.Require().Nil(err)

	tx, _ = suite.client1.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data1))
	tx.Extra = []byte(proof)
	res1, err := suite.client1.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib1.From + ib1.To,
		IBTPNonce: ib1.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_Status(0), res1.Status, string(res1.Ret))

	txId := fmt.Sprintf("%s-%s-%d", from, to, index)

	//test client0 GetStatus
	txRes0, err := suite.client0.client.InvokeBVMContract(rpcx.TransactionMgrContractAddr, "GetStatus",nil, pb.String(txId))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_Status(0), txRes0.Status, string(txRes0.Ret))

	//test client1 GetStatus
	txRes1, err := suite.client1.client.InvokeBVMContract(rpcx.TransactionMgrContractAddr, "GetStatus",nil, pb.String(txId))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_Status(0), txRes1.Status, string(txRes1.Ret))
}

func (suite *TransactionMgrSuite) Test008_One2One_FinishedStatus_Failure() {
	index := uint64(1)
	from := suite.client0.addr
	to := suite.client1.addr
	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	ib0 := &pb.IBTP{From: from, To: to, Index: index, Type: pb.IBTP_INTERCHAIN, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}
	data0, err := ib0.Marshal()
	suite.Require().Nil(err)

	tx, _ := suite.client0.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data0))
	tx.Extra = []byte(proof)
	res0, err := suite.client0.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib0.From + ib0.To,
		IBTPNonce: ib0.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_Status(0), res0.Status, string(res0.Ret))

	ib1 := &pb.IBTP{From: from, To: to, Index: index, Type: pb.IBTP_RECEIPT_FAILURE, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}
	data1, err := ib1.Marshal()
	suite.Require().Nil(err)

	tx, _ = suite.client1.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data1))
	tx.Extra = []byte(proof)
	res1, err := suite.client1.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib1.From + ib1.To,
		IBTPNonce: ib1.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_Status(0), res1.Status, string(res1.Ret))

	txId := fmt.Sprintf("%s-%s-%d", from, to, index)

	//test client0 GetStatus
	txRes0, err := suite.client0.client.InvokeBVMContract(rpcx.TransactionMgrContractAddr, "GetStatus",nil, pb.String(txId))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_Status(0), txRes0.Status, string(txRes0.Ret))

	//test client1 GetStatus
	txRes1, err := suite.client1.client.InvokeBVMContract(rpcx.TransactionMgrContractAddr, "GetStatus",nil, pb.String(txId))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_Status(0), txRes1.Status, string(txRes1.Ret))
}

func (suite *TransactionMgrSuite) Test009_One2Multi_FinishedStatus_Success() {
	index := uint64(1)
	from := suite.client0.addr
	to0 := suite.client1.addr
	to1 := suite.client2.addr
	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	ib0 := &pb.IBTP{From: from, To: to0, Index: index, Type: pb.IBTP_INTERCHAIN, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}
	ib1 := &pb.IBTP{From: from, To: to1, Index: index, Type: pb.IBTP_INTERCHAIN, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}

	ibtps := pb.IBTPs{}
	ibtps.Ibtps = append(ibtps.Ibtps, ib0)
	ibtps.Ibtps = append(ibtps.Ibtps, ib1)

	data0, err := ibtps.Marshal()
	suite.Require().Nil(err)

	tx, _ := suite.client0.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTPs", pb.Bytes(data0))
	tx.Extra = []byte(proof)
	res0, err := suite.client0.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib0.From + ib0.To,
		IBTPNonce: ib0.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_Status(0), res0.Status, string(res0.Ret))

	//first
	ib1 = &pb.IBTP{From: from, To: to0, Index: index, Type: pb.IBTP_RECEIPT_SUCCESS, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}
	data1, err := ib1.Marshal()
	suite.Require().Nil(err)

	time.Sleep(time.Second)

	tx, _ = suite.client1.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data1))
	tx.Extra = []byte(proof)
	res1, err := suite.client1.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib1.From + ib1.To,
		IBTPNonce: ib1.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_Status(0), res1.Status, string(res1.Ret))

	//second
	ib2 := &pb.IBTP{From: from, To: to1, Index: index, Type: pb.IBTP_RECEIPT_SUCCESS, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}
	data2, err := ib2.Marshal()
	suite.Require().Nil(err)

	time.Sleep(time.Second)

	tx, _ = suite.client2.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data2))
	tx.Extra = []byte(proof)
	res2, err := suite.client2.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib2.From + ib2.To,
		IBTPNonce: ib2.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_Status(0), res2.Status, string(res2.Ret))

	globalTxId := fmt.Sprintf("%s-%s", from, res0.TxHash.String())

	//client0
	client0res, err := suite.client0.client.InvokeBVMContract(rpcx.TransactionMgrContractAddr, "GetStatus",nil, pb.String(globalTxId))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_Status(0), client0res.Status, string(client0res.Ret))

	client0status, err := strconv.Atoi(string(client0res.Ret))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.TransactionStatus_SUCCESS, pb.TransactionStatus(client0status))

	//client1
	client1res, err := suite.client1.client.InvokeBVMContract(rpcx.TransactionMgrContractAddr, "GetStatus",nil, pb.String(globalTxId))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_Status(0), client1res.Status, string(client1res.Ret))

	client1status, err := strconv.Atoi(string(client1res.Ret))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.TransactionStatus_SUCCESS, pb.TransactionStatus(client1status))

	//client2
	client2res, err := suite.client2.client.InvokeBVMContract(rpcx.TransactionMgrContractAddr, "GetStatus",nil, pb.String(globalTxId))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_Status(0), client2res.Status, string(client2res.Ret))

	client2status, err := strconv.Atoi(string(client2res.Ret))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.TransactionStatus_SUCCESS, pb.TransactionStatus(client2status))
}

func (suite *TransactionMgrSuite) Test010_One2Multi_FinishedStatus_Failure() {
	index := uint64(1)
	from := suite.client0.addr
	to0 := suite.client1.addr
	to1 := suite.client2.addr
	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	ib0 := &pb.IBTP{From: from, To: to0, Index: index, Type: pb.IBTP_INTERCHAIN, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}
	ib1 := &pb.IBTP{From: from, To: to1, Index: index, Type: pb.IBTP_INTERCHAIN, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}

	ibtps := pb.IBTPs{}
	ibtps.Ibtps = append(ibtps.Ibtps, ib0)
	ibtps.Ibtps = append(ibtps.Ibtps, ib1)

	data0, err := ibtps.Marshal()
	suite.Require().Nil(err)

	tx, _ := suite.client0.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTPs", pb.Bytes(data0))
	tx.Extra = []byte(proof)
	res0, err := suite.client0.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib0.From + ib0.To,
		IBTPNonce: ib0.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_Status(0), res0.Status, string(res0.Ret))

	//first
	ib1 = &pb.IBTP{From: from, To: to0, Index: index, Type: pb.IBTP_RECEIPT_SUCCESS, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}
	data1, err := ib1.Marshal()
	suite.Require().Nil(err)

	time.Sleep(time.Second)

	tx, _ = suite.client1.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data1))
	tx.Extra = []byte(proof)
	res1, err := suite.client1.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib1.From + ib1.To,
		IBTPNonce: ib1.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_Status(0), res1.Status, string(res1.Ret))

	//second
	ib2 := &pb.IBTP{From: from, To: to1, Index: index, Type: pb.IBTP_RECEIPT_FAILURE, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}
	data2, err := ib2.Marshal()
	suite.Require().Nil(err)

	time.Sleep(time.Second)

	tx, _ = suite.client2.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data2))
	tx.Extra = []byte(proof)
	res2, err := suite.client2.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib2.From + ib2.To,
		IBTPNonce: ib2.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_Status(0), res2.Status, string(res2.Ret))

	globalTxId := fmt.Sprintf("%s-%s", from, res0.TxHash.String())

	//client0
	client0res, err := suite.client0.client.InvokeBVMContract(rpcx.TransactionMgrContractAddr, "GetStatus",nil, pb.String(globalTxId))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_Status(0), client0res.Status, string(client0res.Ret))

	client0status, err := strconv.Atoi(string(client0res.Ret))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.TransactionStatus_FAILURE, pb.TransactionStatus(client0status))

	//client1
	client1res, err := suite.client1.client.InvokeBVMContract(rpcx.TransactionMgrContractAddr, "GetStatus",nil, pb.String(globalTxId))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_Status(0), client1res.Status, string(client1res.Ret))

	client1status, err := strconv.Atoi(string(client1res.Ret))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.TransactionStatus_FAILURE, pb.TransactionStatus(client1status))

	//client2
	client2res, err := suite.client2.client.InvokeBVMContract(rpcx.TransactionMgrContractAddr, "GetStatus",nil, pb.String(globalTxId))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_Status(0), client2res.Status, string(client2res.Ret))

	client2status, err := strconv.Atoi(string(client2res.Ret))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.TransactionStatus_FAILURE, pb.TransactionStatus(client2status))
}
