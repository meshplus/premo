package bxh_tester

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-kit/types"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
)

func (suite *Snake) prepare() (crypto.PrivateKey, crypto.PrivateKey, types.Address, types.Address) {
	kA, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	kB, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)

	from, err := kA.PublicKey().Address()
	suite.Require().Nil(err)
	to, err := kB.PublicKey().Address()
	suite.Require().Nil(err)

	return kA, kB, from, to
}

func (suite *Snake) RegisterAppchain(pk crypto.PrivateKey, chainType string) {
	pubBytes, err := pk.PublicKey().Bytes()
	suite.Require().Nil(err)

	suite.client.SetPrivateKey(pk)
	var pubKeyStr = hex.EncodeToString(pubBytes)
	args := []*pb.Arg{
		rpcx.String(""),                 //validators
		rpcx.Int32(0),                   //consensus_type
		rpcx.String(chainType),          //chain_type
		rpcx.String("AppChain"),         //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	res, err := suite.client.InvokeBVMContract(rpcx.AppchainMgrContractAddr, "Register", args...)
	suite.Require().Nil(err)
	appChain := &rpcx.Appchain{}
	err = json.Unmarshal(res.Ret, appChain)
	suite.Require().Nil(err)
	suite.Require().NotNil(appChain.ID)
}

func (suite *Snake) RegisterRule(pk crypto.PrivateKey, ruleFile string) {
	suite.client.SetPrivateKey(pk)

	from, err := pk.PublicKey().Address()
	suite.Require().Nil(err)

	// deploy rule
	bytes, err := ioutil.ReadFile(ruleFile)
	suite.Require().Nil(err)
	addr, err := suite.client.DeployContract(bytes)
	suite.Require().Nil(err)

	// register rule
	res, err := suite.client.InvokeBVMContract(rpcx.RuleManagerContractAddr, "RegisterRule", pb.String(from.Hex()), pb.String(addr.Hex()))
	suite.Require().Nil(err)
	suite.Require().True(res.IsSuccess())
}

// ------ interchain tests ------
func (suite *Snake) TestHandleIBTPShouldSucceed() {
	kA, kB, from, to := suite.prepare()
	suite.RegisterAppchain(kA, "hyperchain")
	suite.RegisterAppchain(kB, "fabric")
	suite.RegisterRule(kA, "./testdata/simple_rule.wasm")
	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	suite.client.SetPrivateKey(kA)
	ib := &pb.IBTP{From: from.Hex(), To: to.Hex(), Index: 1, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}
	data, err := ib.Marshal()
	suite.Require().Nil(err)

	tx, _ := suite.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data))
	tx.Extra = []byte(proof)
	res, err := suite.client.SendTransactionWithReceipt(tx)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
}

func (suite *Snake) TestHandleIBTPWithNonexistentFrom() {
	kA, kB, from, to := suite.prepare()
	suite.RegisterAppchain(kB, "fabric")
	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	suite.client.SetPrivateKey(kA)
	ib := &pb.IBTP{From: from.Hex(), To: to.Hex(), Index: 1, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}
	data, err := ib.Marshal()
	suite.Require().Nil(err)

	tx, _ := suite.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data))
	tx.Extra = []byte(proof)
	res, err := suite.client.SendTransactionWithReceipt(tx)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
}

func (suite *Snake) TestHandleIBTPWithNonexistentTo() {
	kA, _, from, to := suite.prepare()
	suite.RegisterAppchain(kA, "hyperchain")
	suite.RegisterRule(kA, "./testdata/simple_rule.wasm")
	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	ib := &pb.IBTP{From: from.Hex(), To: to.Hex(), Index: 1, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}
	data, err := ib.Marshal()
	suite.Require().Nil(err)

	tx, _ := suite.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data))
	tx.Extra = []byte(proof)
	res, err := suite.client.SendTransactionWithReceipt(tx)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
}

func (suite *Snake) TestHandleIBTPWithNonexistentRule() {
	kA, kB, from, to := suite.prepare()
	suite.RegisterAppchain(kA, "hyperchain")
	suite.RegisterAppchain(kB, "fabric")
	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	suite.client.SetPrivateKey(kA)
	ib := &pb.IBTP{From: from.Hex(), To: to.Hex(), Index: 1, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}
	data, err := ib.Marshal()
	suite.Require().Nil(err)

	tx, _ := suite.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data))
	tx.Extra = []byte(proof)
	res, err := suite.client.SendTransactionWithReceipt(tx)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
}

func (suite *Snake) TestHandleIBTPWithWrongIBTPIndex() {
	kA, kB, from, to := suite.prepare()
	suite.RegisterAppchain(kA, "hyperchain")
	suite.RegisterAppchain(kB, "fabric")
	suite.RegisterRule(kA, "./testdata/simple_rule.wasm")
	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	suite.client.SetPrivateKey(kA)
	ib := &pb.IBTP{From: from.Hex(), To: to.Hex(), Index: 2, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}
	data, err := ib.Marshal()
	suite.Require().Nil(err)

	tx, _ := suite.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data))
	tx.Extra = []byte(proof)
	res, err := suite.client.SendTransactionWithReceipt(tx)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
}

func (suite *Snake) TestGetIBTPByID() {
	kA, kB, from, to := suite.prepare()
	suite.RegisterAppchain(kA, "hyperchain")
	suite.RegisterAppchain(kB, "fabric")
	suite.RegisterRule(kA, "./testdata/simple_rule.wasm")
	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	suite.client.SetPrivateKey(kA)
	ib := &pb.IBTP{From: from.Hex(), To: to.Hex(), Index: 1, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}
	data, err := ib.Marshal()
	suite.Require().Nil(err)

	tx, _ := suite.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data))
	tx.Extra = []byte(proof)
	res, err := suite.client.SendTransactionWithReceipt(tx)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	ib.Index = 2
	data, err = ib.Marshal()
	suite.Require().Nil(err)

	tx, _ = suite.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data))
	tx.Extra = []byte(proof)
	res, err = suite.client.SendTransactionWithReceipt(tx)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	ib.Index = 3
	data, err = ib.Marshal()
	suite.Require().Nil(err)
	tx, _ = suite.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data))
	tx.Extra = []byte(proof)
	res, err = suite.client.SendTransactionWithReceipt(tx)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	// get IBTP by ID
	ib.Index = 2
	res, err = suite.client.InvokeBVMContract(rpcx.InterchainContractAddr, "GetIBTPByID", pb.String(ib.ID()))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
}

func (suite *Snake) TestHandleIBTPWithWrongProof() {
	kA, kB, from, to := suite.prepare()
	suite.RegisterAppchain(kA, "hyperchain")
	suite.RegisterAppchain(kB, "fabric")
	suite.RegisterRule(kA, "./testdata/example.wasm")
	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	suite.client.SetPrivateKey(kA)
	ib := &pb.IBTP{From: from.Hex(), To: to.Hex(), Index: 1, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}
	data, err := ib.Marshal()
	suite.Require().Nil(err)

	tx, _ := suite.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data))
	tx.Extra = []byte(proof)
	res, err := suite.client.SendTransactionWithReceipt(tx)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
}
