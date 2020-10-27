package bxh_tester

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-kit/types"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
)

func (suite *Snake) prepare() (crypto.PrivateKey, crypto.PrivateKey, *types.Address, *types.Address) {
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
	res, err := suite.client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "Register", nil, args...)
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
	addr, err := suite.client.DeployContract(bytes, nil)
	suite.Require().Nil(err)

	// register rule
	res, err := suite.client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), "RegisterRule", nil, pb.String(from.String()), pb.String(addr.String()))
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
	ib := &pb.IBTP{From: from.String(), To: to.String(), Index: 1, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}

	tx, _ := suite.client.GenerateIBTPTx(ib)
	tx.Extra = []byte(proof)
	res, err := suite.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From:      fmt.Sprintf("%s-%s-%d", ib.From, ib.To, ib.Category()),
		IBTPNonce: ib.Index,
	})
	suite.Require().Nil(err)
	suite.Require().True(res.IsSuccess())
}

func (suite *Snake) TestHandleIBTPWithNonexistentFrom() {
	kA, kB, from, to := suite.prepare()
	suite.RegisterAppchain(kB, "fabric")
	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	suite.client.SetPrivateKey(kA)
	ib := &pb.IBTP{From: from.String(), To: to.String(), Index: 1, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}

	tx, _ := suite.client.GenerateIBTPTx(ib)
	tx.Extra = []byte(proof)
	_, err := suite.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From:      fmt.Sprintf("%s-%s-%d", ib.From, ib.To, ib.Category()),
		IBTPNonce: ib.Index,
	})
	suite.Require().NotNil(err)
}

func (suite *Snake) TestHandleIBTPWithNonexistentTo() {
	kA, _, from, to := suite.prepare()
	suite.RegisterAppchain(kA, "hyperchain")
	suite.RegisterRule(kA, "./testdata/simple_rule.wasm")
	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	ib := &pb.IBTP{From: from.String(), To: to.String(), Index: 1, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}

	tx, _ := suite.client.GenerateIBTPTx(ib)
	tx.Extra = []byte(proof)
	res, err := suite.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From:      fmt.Sprintf("%s-%s-%d", ib.From, ib.To, ib.Category()),
		IBTPNonce: ib.Index,
	})
	suite.Require().Nil(err)
	suite.Require().True(res.IsSuccess())
}

func (suite *Snake) TestHandleIBTPWithNonexistentRule() {
	kA, kB, from, to := suite.prepare()
	suite.RegisterAppchain(kA, "hyperchain")
	suite.RegisterAppchain(kB, "fabric")
	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	suite.client.SetPrivateKey(kA)
	ib := &pb.IBTP{From: from.String(), To: to.String(), Index: 1, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}

	tx, _ := suite.client.GenerateIBTPTx(ib)
	tx.Extra = []byte(proof)
	_, err := suite.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From:      fmt.Sprintf("%s-%s-%d", ib.From, ib.To, ib.Category()),
		IBTPNonce: ib.Index,
	})
	suite.Require().NotNil(err)
}

func (suite *Snake) TestHandleIBTPWithWrongIBTPIndex() {
	kA, kB, from, to := suite.prepare()
	suite.RegisterAppchain(kA, "hyperchain")
	suite.RegisterAppchain(kB, "fabric")
	suite.RegisterRule(kA, "./testdata/simple_rule.wasm")
	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	suite.client.SetPrivateKey(kA)
	ib := &pb.IBTP{From: from.String(), To: to.String(), Index: 2, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}

	tx, _ := suite.client.GenerateIBTPTx(ib)
	tx.Extra = []byte(proof)
	res, err := suite.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From:      fmt.Sprintf("%s-%s-%d", ib.From, ib.To, ib.Category()),
		IBTPNonce: ib.Index,
	})
	suite.Require().NotNil(err)
	suite.Require().Nil(res)
}

func (suite *Snake) TestGetIBTPByID() {
	kA, kB, from, to := suite.prepare()
	suite.RegisterAppchain(kA, "hyperchain")
	suite.RegisterAppchain(kB, "fabric")
	suite.RegisterRule(kA, "./testdata/simple_rule.wasm")
	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	suite.client.SetPrivateKey(kA)
	ib := &pb.IBTP{From: from.String(), To: to.String(), Index: 1, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}

	tx, _ := suite.client.GenerateIBTPTx(ib)
	tx.Extra = []byte(proof)
	res, err := suite.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From:      fmt.Sprintf("%s-%s-%d", ib.From, ib.To, ib.Category()),
		IBTPNonce: ib.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	ib.Index = 2
	tx, _ = suite.client.GenerateIBTPTx(ib)
	tx.Extra = []byte(proof)
	res, err = suite.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From:      fmt.Sprintf("%s-%s-%d", ib.From, ib.To, ib.Category()),
		IBTPNonce: ib.Index,
	})
	suite.Require().Nil(err)
	suite.Require().True(res.IsSuccess())

	ib.Index = 3
	tx, _ = suite.client.GenerateIBTPTx(ib)
	tx.Extra = []byte(proof)
	res, err = suite.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From:      fmt.Sprintf("%s-%s-%d", ib.From, ib.To, ib.Category()),
		IBTPNonce: ib.Index,
	})
	suite.Require().Nil(err)
	suite.Require().True(res.IsSuccess())

	// get IBTP by ID
	ib.Index = 2
	res, err = suite.client.InvokeBVMContract(constant.InterchainContractAddr.Address(), "GetIBTPByID", nil, pb.String(ib.ID()))
	suite.Require().Nil(err)
	suite.Require().True(res.IsSuccess())
	suite.client.SetPrivateKey(suite.pk)
}

func (suite *Snake) TestHandleIBTPWithWrongProof() {
	kA, kB, from, to := suite.prepare()
	suite.RegisterAppchain(kA, "hyperchain")
	suite.RegisterAppchain(kB, "fabric")
	suite.RegisterRule(kA, "./testdata/example.wasm")
	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	suite.client.SetPrivateKey(kA)
	ib := &pb.IBTP{From: from.String(), To: to.String(), Index: 1, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}
	data, err := ib.Marshal()
	suite.Require().Nil(err)

	tx, _ := suite.client.GenerateContractTx(pb.TransactionData_BVM, constant.InterchainContractAddr.Address(), "HandleIBTP", pb.Bytes(data))
	tx.Extra = []byte(proof)
	_, err = suite.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From:      ib.From + ib.To,
		IBTPNonce: ib.Index,
	})
	suite.Require().NotNil(err)
}
