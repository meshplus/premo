package bxh_tester

import (
	"crypto/sha256"
	"time"

	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
)

type Model9 struct {
	*Snake
}

// ------ interchain tests ------
func (suite *Model9) Test0901_HandleIBTPShouldSucceed() {
	kA, ChainID1, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	_, ChainID2, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	suite.BindRule(kA, "./testdata/simple_rule.wasm", ChainID1)

	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	client := suite.NewClient(kA)

	ib := &pb.IBTP{From: ChainID1, To: ChainID2, Index: 1, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}

	tx, _ := client.GenerateIBTPTx(ib)
	tx.Extra = []byte(proof)
	res, err := client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().Equal(res.Status, pb.Receipt_SUCCESS)
}

func (suite *Model9) Test0902_HandleIBTPWithNonexistentFrom() {
	kB, ChainID2, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	suite.BindRule(kB, "./testdata/simple_rule.wasm", ChainID2)

	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))
	kA, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	from, err := kA.PublicKey().Address()
	suite.Require().Nil(err)
	ChainID1 := "did:bitxhub:appchain" + from.String() + ":."

	client := suite.NewClient(kA)

	ib := &pb.IBTP{From: ChainID1, To: ChainID2, Index: 1, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}

	tx, _ := client.GenerateIBTPTx(ib)
	tx.Extra = []byte(proof)
	res, err := client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().Equal(res.Status, pb.Receipt_FAILED)
	suite.Require().Contains(string(res.Ret), "cannot get registered appchain")
}

func (suite *Model9) Test0903_HandleIBTPWithNonexistentTo() {
	kA, ChainID1, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	suite.BindRule(kA, "./testdata/simple_rule.wasm", ChainID1)

	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))
	kB, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	to, err := kB.PublicKey().Address()
	suite.Require().Nil(err)
	ChainID2 := "did:bitxhub:appchain" + to.String() + ":."

	client := suite.NewClient(kA)

	ib := &pb.IBTP{From: ChainID1, To: ChainID2, Index: 1, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}

	tx, _ := client.GenerateIBTPTx(ib)
	tx.Extra = []byte(proof)
	res, err := client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().Equal(res.Status, pb.Receipt_FAILED)
	suite.Require().Contains(string(res.Ret), "is not registered")
}

func (suite *Model9) Test0904_HandleIBTPWithNonexistentRule() {
	kA, ChainID1, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	_, ChainID2, err := suite.RegisterAppchain()
	suite.Require().Nil(err)

	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	client := suite.NewClient(kA)

	ib := &pb.IBTP{From: ChainID1, To: ChainID2, Index: 1, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}

	tx, _ := client.GenerateIBTPTx(ib)
	tx.Extra = []byte(proof)
	res, err := client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().Equal(res.Status, pb.Receipt_FAILED)
	suite.Require().Contains(string(res.Ret), "appchain didn't register rule")
}

func (suite *Model9) Test0905_HandleIBTPWithWrongIBTPIndex() {
	kA, ChainID1, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	_, ChainID2, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	suite.BindRule(kA, "./testdata/simple_rule.wasm", ChainID1)

	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	client := suite.NewClient(kA)
	ib := &pb.IBTP{From: ChainID1, To: ChainID2, Index: 2, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}

	tx, _ := client.GenerateIBTPTx(ib)
	tx.Extra = []byte(proof)
	res, err := client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().Contains(string(res.Ret), "wrong index")
}

func (suite *Model9) Test0906_GetIBTPByID() {
	kA, ChainID1, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	_, ChainID2, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	suite.BindRule(kA, "./testdata/simple_rule.wasm", ChainID1)

	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	client := suite.NewClient(kA)
	ib := &pb.IBTP{From: ChainID1, To: ChainID2, Index: 1, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}

	tx, _ := client.GenerateIBTPTx(ib)
	tx.Extra = []byte(proof)
	res, err := client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	ib.Index = 2
	tx, _ = client.GenerateIBTPTx(ib)
	tx.Extra = []byte(proof)
	res, err = client.SendTransactionWithReceipt(tx, nil)

	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	ib.Index = 3
	tx, _ = client.GenerateIBTPTx(ib)
	tx.Extra = []byte(proof)
	res, err = client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	// get IBTP by ID
	ib.Index = 2
	res, err = client.InvokeBVMContract(constant.InterchainContractAddr.Address(), "GetIBTPByID", nil, pb.String(ib.ID()))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	suite.Require().NotNil(res.Ret)
}

func (suite *Model9) Test0907_HandleIBTPWithWrongProof() {
	kA, ChainID1, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	_, ChainID2, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	suite.BindRule(kA, "./testdata/simple_rule.wasm", ChainID1)

	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	client := suite.NewClient(kA)
	ib := &pb.IBTP{From: ChainID1, To: ChainID2, Index: 1, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}
	data, err := ib.Marshal()
	suite.Require().Nil(err)

	tx, _ := client.GenerateContractTx(pb.TransactionData_BVM, constant.InterchainContractAddr.Address(), "HandleIBTP", pb.Bytes(data))
	tx.Extra = []byte(proof)
	res, err := client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().Contains(string(res.Ret), "Call using []uint8 as type *pb.IBTP")
}

func (suite *Model9) Test0908_HandleIBTPWithTxInBlock() {
	kA, ChainID1, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	_, ChainID2, err := suite.RegisterAppchain()
	suite.Require().Nil(err)
	suite.BindRule(kA, "./testdata/simple_rule.wasm", ChainID1)

	client := suite.NewClient(kA)

	ib := &pb.IBTP{From: ChainID1, To: ChainID2, Index: 1, Timestamp: time.Now().UnixNano()}
	tx, _ := client.GenerateIBTPTx(ib)

	hash, err := client.SendTransaction(tx, nil)
	suite.Require().Nil(err)

	time.Sleep(time.Second * 2)
	transaction, err := client.GetTransaction(hash)
	suite.Require().Nil(err)
	suite.Require().Equal(transaction.Tx.TransactionHash.String(), hash)

}
