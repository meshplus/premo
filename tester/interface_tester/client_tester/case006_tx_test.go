package interface_tester

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	appchain_mgr "github.com/meshplus/bitxhub-core/appchain-mgr"
	"github.com/meshplus/bitxhub-core/governance"
	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-kit/types"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/meshplus/premo/internal/repo"
)

type RegisterResult struct {
	ChainID    string `json:"chain_id"`
	ProposalID string `json:"proposal_id"`
}

func (suite *Snake) TestTxGetIsTrue() {
	hash, err := suite.sendInterchain()
	suite.Require().Nil(err)

	//wait for bitxhub
	time.Sleep(time.Second * 3)
	url := getURL("transaction/" + hash)

	data, err := httpGet(url)
	suite.Require().Nil(err)
	suite.Require().NotContains(string(data), "error")
	suite.Require().Contains(string(data), "tx_meta")
}

func (suite *Snake) TestTxGetWithNonexistent() {
	wrongHash := "0x0000000000000000000000000000000012345678900000000000000000000000"

	url := getURL("transaction/" + wrongHash)

	data, err := httpGet(url)
	suite.Require().Nil(err)
	suite.Require().Contains(string(data), "error")
	suite.Require().Contains(string(data), "not found in DB")
}

func (suite *Snake) TestTxGetWithInvalidFormat() {
	wrongHash := "0x0000000000000000000000000000000012345678900000000000000000000000"
	url := getURL("transaction/" + wrongHash + "123!@#")

	data, err := httpGet(url)
	suite.Require().Nil(err)
	suite.Require().Contains(string(data), "error")
	suite.Require().Contains(string(data), "invalid format of tx hash for querying transaction")
}

func (suite Snake) TestTxSendIsTrue() {
	txType := 0
	amount := uint64(1)

	kA, kB, from, to := suite.prepare()
	suite.registerAppchain(kA, "hyperchain")
	suite.registerAppchain(kB, "fabric")
	suite.BindRule(kA, "../../../config/rule.wasm")

	data := &pb.TransactionData{
		Type:   pb.TransactionData_Type(txType),
		Amount: amount,
	}

	payload, err := data.Marshal()
	suite.Require().Nil(err)

	// get nonce for this account
	nonce, err := suite.client.GetPendingNonceByAccount(from.String())
	suite.Require().Nil(err)
	tx := &pb.BxhTransaction{
		From:      from,
		To:        to,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
		Nonce:     nonce,
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

func (suite Snake) TestTxSendWithFromAddressIsNil() {
	txType := 0
	amount := uint64(1)

	kA, kB, _, to := suite.prepare()
	suite.registerAppchain(kA, "hyperchain")
	suite.registerAppchain(kB, "fabric")
	suite.BindRule(kA, "../../../config/rule.wasm")

	data := &pb.TransactionData{
		Type:   pb.TransactionData_Type(txType),
		Amount: amount,
	}

	payload, err := data.Marshal()
	suite.Require().Nil(err)

	tx := &pb.BxhTransaction{
		//From: from,
		To:        to,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
		Nonce:     1,
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

func (suite Snake) TestTxSendWithToAddressIsNil() {
	txType := 0
	amount := uint64(1)

	kA, kB, from, _ := suite.prepare()
	suite.registerAppchain(kA, "hyperchain")
	suite.registerAppchain(kB, "fabric")
	suite.BindRule(kA, "../../../config/rule.wasm")

	data := &pb.TransactionData{
		Type:   pb.TransactionData_Type(txType),
		Amount: amount,
	}

	payload, err := data.Marshal()
	suite.Require().Nil(err)

	tx := &pb.BxhTransaction{
		From: from,
		//To:        to,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
		Nonce:     1,
	}

	err = tx.Sign(kA)
	suite.Require().Nil(err)

	reqData, err := json.Marshal(tx)
	suite.Require().Nil(err)

	url := getURL("transaction")

	ret, err := httpPost(url, reqData)
	suite.Require().Nil(err)
	suite.Require().Contains(string(ret), "tx to address is nil")
}

func (suite Snake) TestTxSendWithEmptySign() {
	txType := 0
	amount := uint64(1)

	kA, kB, from, to := suite.prepare()
	suite.registerAppchain(kA, "hyperchain")
	suite.registerAppchain(kB, "fabric")
	suite.BindRule(kA, "../../../config/rule.wasm")

	data := &pb.TransactionData{
		Type:   pb.TransactionData_Type(txType),
		Amount: amount,
	}

	payload, err := data.Marshal()
	suite.Require().Nil(err)

	tx := &pb.BxhTransaction{
		From:      from,
		To:        to,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
		Nonce:     1,
	}

	//err = tx.Sign(kA)
	//suite.Require().Nil(err)

	reqData, err := json.Marshal(tx)
	suite.Require().Nil(err)

	url := getURL("transaction")

	ret, err := httpPost(url, reqData)
	suite.Require().Nil(err)
	suite.Require().Contains(string(ret), "signature can't be empty")
}

func (suite Snake) TestTxSendWithInvalidSign() {
	txType := 0
	amount := uint64(1)

	kA, kB, from, to := suite.prepare()
	suite.registerAppchain(kA, "hyperchain")
	suite.registerAppchain(kB, "fabric")
	suite.BindRule(kA, "../../../config/rule.wasm")

	data := &pb.TransactionData{
		Type:   pb.TransactionData_Type(txType),
		Amount: amount,
	}

	payload, err := data.Marshal()
	suite.Require().Nil(err)

	tx := &pb.BxhTransaction{
		From:      from,
		To:        to,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
		Nonce:     1,
	}

	err = tx.Sign(kB)
	suite.Require().Nil(err)

	reqData, err := json.Marshal(tx)
	suite.Require().Nil(err)

	url := getURL("transaction")

	ret, err := httpPost(url, reqData)
	suite.Require().Nil(err)
	suite.Require().Contains(string(ret), "invalid signature")
}

func (suite Snake) TestTxSendWithEmptyTimestamp() {
	txType := 0
	amount := uint64(1)

	kA, kB, from, to := suite.prepare()
	suite.registerAppchain(kA, "hyperchain")
	suite.registerAppchain(kB, "fabric")
	suite.BindRule(kA, "../../../config/rule.wasm")

	data := &pb.TransactionData{
		Type:   pb.TransactionData_Type(txType),
		Amount: amount,
	}

	payload, err := data.Marshal()
	suite.Require().Nil(err)

	tx := &pb.BxhTransaction{
		From: from,
		To:   to,
		//Timestamp: time.Now().UnixNano(),
		Payload: payload,
		Nonce:   1,
	}

	err = tx.Sign(kA)
	suite.Require().Nil(err)

	reqData, err := json.Marshal(tx)
	suite.Require().Nil(err)

	url := getURL("transaction")

	ret, err := httpPost(url, reqData)
	suite.Require().Nil(err)
	suite.Require().Contains(string(ret), "timestamp is illegal")
}

func (suite Snake) TestTxSendWithErrorTimestamp() {
	txType := 0
	amount := uint64(1)

	kA, kB, from, to := suite.prepare()
	suite.registerAppchain(kA, "hyperchain")
	suite.registerAppchain(kB, "fabric")
	suite.BindRule(kA, "../../../config/rule.wasm")

	data := &pb.TransactionData{
		Type:   pb.TransactionData_Type(txType),
		Amount: amount,
	}

	payload, err := data.Marshal()
	suite.Require().Nil(err)

	tx := &pb.BxhTransaction{
		From:      from,
		To:        to,
		Timestamp: 1608624000, // 2020/12/22 16:00:00
		Payload:   payload,
		Nonce:     1,
	}

	err = tx.Sign(kA)
	suite.Require().Nil(err)

	reqData, err := json.Marshal(tx)
	suite.Require().Nil(err)

	url := getURL("transaction")

	ret, err := httpPost(url, reqData)
	suite.Require().Nil(err)
	suite.Require().Contains(string(ret), "timestamp is illegal")
}

func (suite Snake) TestTxSendWithEmptyNonce() {
	txType := 0
	amount := uint64(1)

	kA, kB, from, to := suite.prepare()
	suite.registerAppchain(kA, "hyperchain")
	suite.registerAppchain(kB, "fabric")
	suite.BindRule(kA, "../../../config/rule.wasm")

	data := &pb.TransactionData{
		Type:   pb.TransactionData_Type(txType),
		Amount: amount,
	}

	payload, err := data.Marshal()
	suite.Require().Nil(err)

	tx := &pb.BxhTransaction{
		From:      from,
		To:        to,
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
		//Nonce:     1,
	}

	err = tx.Sign(kA)
	suite.Require().Nil(err)

	reqData, err := json.Marshal(tx)
	suite.Require().Nil(err)

	url := getURL("transaction")

	ret, err := httpPost(url, reqData)
	suite.Require().Nil(err)
	suite.Require().Contains(string(ret), "nonce is illegal")
}

func (suite Snake) TestTxSendWithEmptyPayload() {

	kA, kB, from, to := suite.prepare()
	suite.registerAppchain(kA, "hyperchain")
	suite.registerAppchain(kB, "fabric")
	suite.BindRule(kA, "../../../config/rule.wasm")

	tx := &pb.BxhTransaction{
		From:      from,
		To:        to,
		Timestamp: time.Now().UnixNano(),
		//Payload:   payload,
		Nonce: 1,
	}

	err := tx.Sign(kA)
	suite.Require().Nil(err)

	reqData, err := json.Marshal(tx)
	suite.Require().Nil(err)

	url := getURL("transaction")

	ret, err := httpPost(url, reqData)
	suite.Require().Nil(err)
	suite.Require().Contains(string(ret), "tx payload and ibtp can't both be nil")
}

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

func (suite *Snake) registerAppchain(pk crypto.PrivateKey, chainType string) {
	pubBytes, err := pk.PublicKey().Bytes()
	suite.Require().Nil(err)

	client := suite.NewClient(pk)

	var pubKeyStr = base64.StdEncoding.EncodeToString(pubBytes)
	args := []*pb.Arg{
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String(chainType),          //chain_type
		rpcx.String("AppChain"),         //name
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "Register", nil, args...)
	suite.Require().Nil(err)
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	suite.Require().Nil(err)
	suite.Require().NotNil(result.ChainID)
	err = suite.VotePass(result.ProposalID)
	suite.Require().Nil(err)

	res, err = suite.GetChainStatusById(result.ChainID)
	suite.Require().Nil(err)
	appchain := &appchain_mgr.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(governance.GovernanceAvailable, appchain.Status)
}

func (suite *Snake) BindRule(pk crypto.PrivateKey, ruleFile string) {
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

func (suite *Snake) VotePass(id string) error {
	node1, err := repo.Node1Path()
	if err != nil {
		return err
	}

	key, err := asym.RestorePrivateKey(node1, repo.KeyPassword)
	if err != nil {
		return err
	}

	_, err = suite.vote(key, pb.String(id), pb.String("approve"), pb.String("Appchain Pass"))
	if err != nil {
		return err
	}

	node2, err := repo.Node2Path()
	if err != nil {
		return err
	}

	key, err = asym.RestorePrivateKey(node2, repo.KeyPassword)
	if err != nil {
		return err
	}

	_, err = suite.vote(key, pb.String(id), pb.String("approve"), pb.String("Appchain Pass"))
	if err != nil {
		return err
	}

	node3, err := repo.Node3Path()
	if err != nil {
		return err
	}

	key, err = asym.RestorePrivateKey(node3, repo.KeyPassword)
	if err != nil {
		return err
	}

	_, err = suite.vote(key, pb.String(id), pb.String("approve"), pb.String("Appchain Pass"))
	if err != nil {
		return err
	}
	return nil
}

func (suite *Snake) vote(key crypto.PrivateKey, args ...*pb.Arg) (*pb.Receipt, error) {
	client, err := rpcx.New(
		rpcx.WithNodesInfo(&rpcx.NodeInfo{Addr: cfg.addrs[0]}),
		rpcx.WithLogger(cfg.logger),
		rpcx.WithPrivateKey(key),
	)
	address, err := key.PublicKey().Address()
	if err != nil {
		return nil, err
	}
	invokePayload := &pb.InvokePayload{
		Method: "Vote",
		Args:   args,
	}

	payload, err := invokePayload.Marshal()
	if err != nil {
		return nil, err
	}

	data := &pb.TransactionData{
		Type:    pb.TransactionData_INVOKE,
		VmType:  pb.TransactionData_BVM,
		Payload: payload,
	}
	payload, err = data.Marshal()

	tx := &pb.BxhTransaction{
		From:      address,
		To:        constant.GovernanceContractAddr.Address(),
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}
	if err != nil {
		return nil, err
	}
	receipt, err := client.SendTransactionWithReceipt(tx, nil)
	if err != nil {
		return nil, err
	}
	return receipt, nil
}

func (suite *Snake) GetChainStatusById(id string) (*pb.Receipt, error) {
	node, err := repo.Node1Path()
	key, err := asym.RestorePrivateKey(node, repo.KeyPassword)
	if err != nil {
		return nil, err
	}
	client, err := rpcx.New(
		rpcx.WithNodesInfo(&rpcx.NodeInfo{Addr: cfg.addrs[0]}),
		rpcx.WithLogger(cfg.logger),
		rpcx.WithPrivateKey(key),
	)
	address, err := key.PublicKey().Address()
	if err != nil {
		return nil, err
	}
	args := []*pb.Arg{
		rpcx.String(id),
	}
	invokePayload := &pb.InvokePayload{
		Method: "GetAppchain",
		Args:   args,
	}

	payload, err := invokePayload.Marshal()
	if err != nil {
		return nil, err
	}

	data := &pb.TransactionData{
		Type:    pb.TransactionData_INVOKE,
		VmType:  pb.TransactionData_BVM,
		Payload: payload,
	}
	payload, err = data.Marshal()

	tx := &pb.BxhTransaction{
		From:      address,
		To:        constant.AppchainMgrContractAddr.Address(),
		Timestamp: time.Now().UnixNano(),
		Payload:   payload,
	}
	if err != nil {
		return nil, err
	}
	receipt, err := client.SendTransactionWithReceipt(tx, nil)
	if err != nil {
		return nil, err
	}
	return receipt, nil
}

func (suite Snake) sendInterchainWithReceipt() (crypto.PrivateKey, crypto.PrivateKey, *types.Address, *types.Address, *pb.Receipt, error) {
	//sendInterchain
	kA, kB, from, to := suite.prepare()
	suite.registerAppchain(kA, "hyperchain")
	suite.registerAppchain(kB, "fabric")
	suite.BindRule(kA, "../../../config/rule.wasm")
	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	suite.client.SetPrivateKey(kA)
	ib := &pb.IBTP{From: from.String(), To: to.String(), Index: 1, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}

	tx, _ := suite.client.GenerateIBTPTx(ib)
	tx.Extra = []byte(proof)
	res, err := suite.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: fmt.Sprintf("%s-%s-%d", ib.From, ib.To, ib.Category()),
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
	suite.BindRule(kA, "../../../config/rule.wasm")
	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	suite.client.SetPrivateKey(kA)
	ib := &pb.IBTP{From: from.String(), To: to.String(), Index: 1, Timestamp: time.Now().UnixNano(), Proof: proofHash[:]}

	tx, _ := suite.client.GenerateIBTPTx(ib)
	tx.Extra = []byte(proof)
	hash, err := suite.client.SendTransaction(tx, &rpcx.TransactOpts{
		From: fmt.Sprintf("%s-%s-%d", ib.From, ib.To, ib.Category()),
	})
	if err != nil {
		return "", err
	}
	return hash, err
}
