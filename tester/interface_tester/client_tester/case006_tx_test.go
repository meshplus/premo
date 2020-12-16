package interface_tester

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"

	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-kit/types"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
)

func (suite *Snake) TestTxGetIsTrue() {
	hash, err := suite.sendInterchain()
	suite.Require().Nil(err)

	//wait for bitxhub
	time.Sleep(time.Second * 3)
	url := getURL("transaction/" + hash)

	data, err := httpGet(url)
	suite.Require().Nil(err)
	suite.Require().NotContains(string(data), "error")
}

func (suite *Snake) TestTxGetIsFalse() {
	hash, err := suite.sendInterchain()
	suite.Require().Nil(err)

	//wait for bitxhub
	time.Sleep(time.Second * 3)
	hashByte := []byte(hash)
	hashByte[len(hash)-1] = hashByte[len(hash)-1] + 1

	url := getURL("transaction/" + string(hashByte))

	data, err := httpGet(url)
	suite.Require().Nil(err)
	suite.Require().Contains(string(data), "error")
}

func (suite Snake) TestTxSendIsTrue() {
	txType := 0
	amount := uint64(1)

	kA, kB, from, to := suite.prepare()
	suite.registerAppchain(kA, "hyperchain")
	suite.registerAppchain(kB, "fabric")
	suite.registerRule(kA, "../../../config/rule.wasm")

	data := &pb.TransactionData{
		Type:   pb.TransactionData_Type(txType),
		Amount: amount,
	}

	payload, err := data.Marshal()
	suite.Require().Nil(err)

	tx := &pb.Transaction{
		From:      from,
		To:        to,
		Timestamp: time.Now().UnixNano(),
		Nonce:     rand.Uint64(),
		Payload:   payload,
	}

	err = tx.Sign(kA)
	suite.Require().Nil(err)

	reqData, err := json.Marshal(tx)
	suite.Require().Nil(err)

	url := getURL("transaction")

	ret, err := httpPost(url, reqData)
	suite.Require().Nil(err)
	suite.Require().Contains(string(ret), "tx_hash")
}

func (suite Snake) TestTxSendIsFalse() {
	txType := 0
	amount := uint64(1)

	kA, kB, _, to := suite.prepare()
	suite.registerAppchain(kA, "hyperchain")
	suite.registerAppchain(kB, "fabric")
	suite.registerRule(kA, "../../../config/rule.wasm")

	data := &pb.TransactionData{
		Type:   pb.TransactionData_Type(txType),
		Amount: amount,
	}

	payload, err := data.Marshal()
	suite.Require().Nil(err)

	tx := &pb.Transaction{
		//From: from,
		To:        to,
		Timestamp: time.Now().UnixNano(),
		Nonce:     rand.Uint64(),
		Payload:   payload,
	}

	err = tx.Sign(kA)
	suite.Require().Nil(err)

	reqData, err := json.Marshal(tx)
	suite.Require().Nil(err)

	url := getURL("transaction")

	ret, err := httpPost(url, reqData)
	suite.Require().Nil(err)
	suite.Require().Contains(string(ret), "tx from address is nil")
}

func (suite *Snake) prepare() (crypto.PrivateKey, crypto.PrivateKey, *types.Address, *types.Address) {
	kA := suite.pk
	//suite.Require().Nil(err)
	kB, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)

	from, err := kA.PublicKey().Address()
	suite.Require().Nil(err)
	to, err := kB.PublicKey().Address()
	suite.Require().Nil(err)

	return kA, kB, from, to
}

func (suite *Snake) registerAppchain(pk crypto.PrivateKey, chainType string) {
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

func (suite *Snake) registerRule(pk crypto.PrivateKey, ruleFile string) {
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

func (suite Snake) sendInterchainWithReceipt() (crypto.PrivateKey, crypto.PrivateKey, *types.Address, *types.Address, *pb.Receipt, error) {
	//sendInterchain
	kA, kB, from, to := suite.prepare()
	suite.registerAppchain(kA, "hyperchain")
	suite.registerAppchain(kB, "fabric")
	suite.registerRule(kA, "../../../config/rule.wasm")
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
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	return kA, kB, from, to, res, nil
}

func (suite Snake) sendInterchain() (string, error) {
	//sendInterchain
	kA, kB, from, to := suite.prepare()
	suite.registerAppchain(kA, "hyperchain")
	suite.registerAppchain(kB, "fabric")
	suite.registerRule(kA, "../../../config/rule.wasm")
	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	suite.client.SetPrivateKey(kA)
	ib := &pb.IBTP{From: from.String(), To: to.String(), Index: 1, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}

	tx, _ := suite.client.GenerateIBTPTx(ib)
	tx.Extra = []byte(proof)
	hash, err := suite.client.SendTransaction(tx, &rpcx.TransactOpts{
		From:      fmt.Sprintf("%s-%s-%d", ib.From, ib.To, ib.Category()),
		IBTPNonce: ib.Index,
	})
	if err != nil {
		return "", err
	}
	return hash, err
}
