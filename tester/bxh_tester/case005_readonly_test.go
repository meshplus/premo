package bxh_tester

import (
	"fmt"
	"github.com/meshplus/bitxhub-model/constant"
	"math/rand"
	"time"

	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/types"
	"github.com/meshplus/bitxhub-model/pb"
)

func (suite *Snake) TestNormalReadOnly() {
	keyForNormal := "key_for_normal"
	valueForNormal := "value_for_normal"
	receipt, err := suite.client.InvokeBVMContract(constant.StoreContractAddr.Address(), "Set", nil, pb.String(keyForNormal), pb.String(valueForNormal))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, receipt.Status)

	queryKey, err := genContractTransaction(pb.TransactionData_BVM, suite.pk,
		constant.StoreContractAddr.Address(), "Get", pb.String(keyForNormal))
	queryKey.Nonce = 1

	receipt, err = suite.client.SendView(queryKey)
	suite.Require().Nil(err)
	suite.Require().True(receipt.Status == pb.Receipt_SUCCESS)
	suite.Require().Equal(valueForNormal, string(receipt.Ret))
}

func (suite *Snake) TestSendTx2ReadOnlyApi() {
	rand.Seed(time.Now().UnixNano())
	randKey := make([]byte, 20)
	_, err := rand.Read(randKey)
	suite.Require().Nil(err)

	valueForRand := "value_for_rand"

	// send tx to SendView api and value not set
	tx, err := genContractTransaction(pb.TransactionData_BVM, suite.pk,
		constant.StoreContractAddr.Address(), "Set", pb.String(string(randKey)), pb.String(valueForRand))
	tx.Nonce = 1
	queryKey, err := genContractTransaction(pb.TransactionData_BVM, suite.pk,
		constant.StoreContractAddr.Address(), "Get", pb.String(string(randKey)))
	queryKey.Nonce = 1

	receipt, err := suite.client.SendView(tx)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, receipt.Status)

	receipt, err = suite.client.SendView(queryKey)
	suite.Require().Nil(err)
	suite.Require().True(receipt.Status == pb.Receipt_FAILED)

	// send tx to SendTransactionWithReceipt api and value got set
	receipt, err = suite.client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, receipt.Status)

	receipt, err = suite.client.SendView(queryKey)
	suite.Require().Nil(err)
	suite.Require().True(receipt.Status == pb.Receipt_SUCCESS)
	suite.Require().Equal(valueForRand, string(receipt.Ret))
}

func genContractTransaction(
	vmType pb.TransactionData_VMType, privateKey crypto.PrivateKey,
	address *types.Address, method string, args ...*pb.Arg) (*pb.Transaction, error) {
	from, err := privateKey.PublicKey().Address()
	if err != nil {
		return nil, err
	}

	pl := &pb.InvokePayload{
		Method: method,
		Args:   args[:],
	}

	data, err := pl.Marshal()
	if err != nil {
		return nil, err
	}

	td := &pb.TransactionData{
		Type:    pb.TransactionData_INVOKE,
		VmType:  vmType,
		Payload: data,
	}

	payload, err := td.Marshal()

	tx := &pb.Transaction{
		From:      from,
		To:        address,
		Timestamp: time.Now().UnixNano(),
		Nonce:     rand.Uint64(),
		Payload:   payload,
	}

	if err := tx.Sign(privateKey); err != nil {
		return nil, fmt.Errorf("tx sign: %w", err)
	}

	tx.TransactionHash = tx.Hash()

	return tx, nil
}
