package bxh_tester

import (
	"fmt"
	"time"

	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-kit/types"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
)

type Model5 struct {
	*Snake
}

//tc:向中继链发送只读交易查询交易余额
func (suite *Model5) Test0501_NormalReadOnlyIsSuccess() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	keyForNormal := "key_for_normal"
	valueForNormal := "value_for_normal"
	receipt, err := client.InvokeBVMContract(constant.StoreContractAddr.Address(), "Set", nil, pb.String(keyForNormal), pb.String(valueForNormal))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, receipt.Status)
	queryKey, err := genContractTransaction(pb.TransactionData_BVM, pk,
		constant.StoreContractAddr.Address(), "Get", pb.String(keyForNormal))
	suite.Require().Nil(err)
	queryKey.Nonce = 1
	receipt, err = client.SendView(queryKey)
	suite.Require().Nil(err)
	suite.Require().True(receipt.IsSuccess())
	suite.Require().Equal(valueForNormal, string(receipt.Ret))
}

// genContractTransaction generated tx by args
func genContractTransaction(vmType pb.TransactionData_VMType, pk crypto.PrivateKey, address *types.Address, method string, args ...*pb.Arg) (*pb.BxhTransaction, error) {
	from, err := pk.PublicKey().Address()
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
	if err != nil {
		return nil, err
	}
	tx := &pb.BxhTransaction{
		From:      from,
		To:        address,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}
	if err := tx.Sign(pk); err != nil {
		return nil, fmt.Errorf("tx sign: %w", err)
	}
	tx.TransactionHash = tx.Hash()
	return tx, nil
}
