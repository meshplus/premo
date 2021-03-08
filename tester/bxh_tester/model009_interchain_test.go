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

func (suite Snake) NewClient(pk crypto.PrivateKey) *rpcx.ChainClient {
	node0 := &rpcx.NodeInfo{Addr: cfg.addrs[0]}
	client, err := rpcx.New(
		rpcx.WithNodesInfo(node0),
		rpcx.WithLogger(cfg.logger),
		rpcx.WithPrivateKey(pk),
	)
	suite.Require().Nil(err)
	return client
}

func (suite *Snake) RegisterAppchain(pk crypto.PrivateKey, chainType string) {
	pubBytes, err := pk.PublicKey().Bytes()
	suite.Require().Nil(err)

	client := suite.NewClient(pk)

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
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "Register", nil, args...)
	suite.Require().Nil(err)
	appChain := &rpcx.Appchain{}
	err = json.Unmarshal(res.Ret, appChain)
	suite.Require().Nil(err)
	suite.Require().NotNil(appChain.ID)
}

func (suite *Snake) RegisterRule(pk crypto.PrivateKey, ruleFile string) {
	client := suite.NewClient(pk)

	from, err := pk.PublicKey().Address()
	suite.Require().Nil(err)

	// deploy rule
	bytes, err := ioutil.ReadFile(ruleFile)
	suite.Require().Nil(err)
	addr, err := client.DeployContract(bytes, nil)
	suite.Require().Nil(err)

	// register rule
	res, err := client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), "RegisterRule", nil, pb.String(from.String()), pb.String(addr.String()))
	suite.Require().Nil(err)
	suite.Require().True(res.IsSuccess())
}

// ------ interchain tests ------
//tc:发送跨链交易，正常发送跨链交易，返回回执状态成功
func (suite *Snake) Test0901_HandleIBTPShouldSucceed() {
	kA, kB, from, to := suite.prepare()
	suite.RegisterAppchain(kA, "hyperchain")
	suite.RegisterAppchain(kB, "fabric")
	suite.RegisterRule(kA, "./testdata/simple_rule.wasm")
	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	client := suite.NewClient(kA)
	ib := &pb.IBTP{From: from.String(), To: to.String(), Index: 1, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}

	tx, _ := client.GenerateIBTPTx(ib)
	tx.Extra = []byte(proof)
	res, err := client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From:      fmt.Sprintf("%s-%s-%d", ib.From, ib.To, ib.Category()),
		IBTPNonce: ib.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(res.Status, pb.Receipt_SUCCESS)
	fmt.Println(string(res.Ret))
}

//tc:发送跨链交易，来源链不存在，返回回执状态失败
func (suite *Snake) Test0902_HandleIBTPWithNonexistentFrom() {
	kA, kB, from, to := suite.prepare()
	suite.RegisterAppchain(kB, "fabric")
	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	client := suite.NewClient(kA)

	ib := &pb.IBTP{From: from.String(), To: to.String(), Index: 1, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}

	tx, _ := client.GenerateIBTPTx(ib)
	tx.Extra = []byte(proof)
	_, err := client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From:      fmt.Sprintf("%s-%s-%d", ib.From, ib.To, ib.Category()),
		IBTPNonce: ib.Index,
	})
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "not found in DB")
}

//tc:发送跨链交易，目的链不存在，返回回执状态失败
func (suite *Snake) Test0903_HandleIBTPWithNonexistentTo() {
	kA, _, from, to := suite.prepare()
	suite.RegisterAppchain(kA, "hyperchain")
	suite.RegisterRule(kA, "./testdata/simple_rule.wasm")
	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	client := suite.NewClient(kA)

	ib := &pb.IBTP{From: from.String(), To: to.String(), Index: 1, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}

	tx, _ := client.GenerateIBTPTx(ib)
	tx.Extra = []byte(proof)
	res, err := client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From:      fmt.Sprintf("%s-%s-%d", ib.From, ib.To, ib.Category()),
		IBTPNonce: ib.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(res.Status, pb.Receipt_SUCCESS)
}

//tc:发送跨链交易，来源链验证规则不存在，返回回执状态失败
func (suite *Snake) Test0904_HandleIBTPWithNonexistentRule() {
	kA, kB, from, to := suite.prepare()
	suite.RegisterAppchain(kA, "hyperchain")
	suite.RegisterAppchain(kB, "fabric")
	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	client := suite.NewClient(kA)
	ib := &pb.IBTP{From: from.String(), To: to.String(), Index: 1, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}

	tx, _ := client.GenerateIBTPTx(ib)
	tx.Extra = []byte(proof)
	_, err := client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From:      fmt.Sprintf("%s-%s-%d", ib.From, ib.To, ib.Category()),
		IBTPNonce: ib.Index,
	})
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "not found in DB")
}

//tc:发送跨链交易，IBTP index不匹配，返回回执状态失败
func (suite *Snake) Test0905_HandleIBTPWithWrongIBTPIndex() {
	kA, kB, from, to := suite.prepare()
	suite.RegisterAppchain(kA, "hyperchain")
	suite.RegisterAppchain(kB, "fabric")
	suite.RegisterRule(kA, "./testdata/simple_rule.wasm")
	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	client := suite.NewClient(kA)
	ib := &pb.IBTP{From: from.String(), To: to.String(), Index: 2, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}

	tx, _ := client.GenerateIBTPTx(ib)
	tx.Extra = []byte(proof)
	res, err := client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From:      fmt.Sprintf("%s-%s-%d", ib.From, ib.To, ib.Category()),
		IBTPNonce: ib.Index,
	})
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "not found in DB")
	suite.Require().Nil(res)
}

//tc:查询跨链交易，根据IBTP ID查询对应的IBTP，返回回执状态成功
func (suite *Snake) Test0906_GetIBTPByID() {
	kA, kB, from, to := suite.prepare()
	suite.RegisterAppchain(kA, "hyperchain")
	suite.RegisterAppchain(kB, "fabric")
	suite.RegisterRule(kA, "./testdata/simple_rule.wasm")
	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	client := suite.NewClient(kA)
	ib := &pb.IBTP{From: from.String(), To: to.String(), Index: 1, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}

	tx, _ := client.GenerateIBTPTx(ib)
	tx.Extra = []byte(proof)
	res, err := client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From:      fmt.Sprintf("%s-%s-%d", ib.From, ib.To, ib.Category()),
		IBTPNonce: ib.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	ib.Index = 2
	tx, _ = client.GenerateIBTPTx(ib)
	tx.Extra = []byte(proof)
	res, err = client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From:      fmt.Sprintf("%s-%s-%d", ib.From, ib.To, ib.Category()),
		IBTPNonce: ib.Index,
	})
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	ib.Index = 3
	tx, _ = client.GenerateIBTPTx(ib)
	tx.Extra = []byte(proof)
	res, err = client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From:      fmt.Sprintf("%s-%s-%d", ib.From, ib.To, ib.Category()),
		IBTPNonce: ib.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	// get IBTP by ID
	ib.Index = 2
	res, err = client.InvokeBVMContract(constant.InterchainContractAddr.Address(), "GetIBTPByID", nil, pb.String(ib.ID()))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	suite.Require().NotNil(res.Ret)
}

//tc:发送跨链交易，验证规则无法验证proof，返回回执状态失败
func (suite *Snake) Test0907_HandleIBTPWithWrongProof() {
	kA, kB, from, to := suite.prepare()
	suite.RegisterAppchain(kA, "hyperchain")
	suite.RegisterAppchain(kB, "fabric")
	suite.RegisterRule(kA, "./testdata/example.wasm")
	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	client := suite.NewClient(kA)
	ib := &pb.IBTP{From: from.String(), To: to.String(), Index: 1, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}
	data, err := ib.Marshal()
	suite.Require().Nil(err)

	tx, _ := client.GenerateContractTx(pb.TransactionData_BVM, constant.InterchainContractAddr.Address(), "HandleIBTP", pb.Bytes(data))
	tx.Extra = []byte(proof)
	_, err = client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From:      ib.From + ib.To,
		IBTPNonce: ib.Index,
	})
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "not found in DB")
}
