package bxh_tester

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	appchain_mgr "github.com/meshplus/bitxhub-core/appchain-mgr"
	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-kit/types"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/meshplus/premo/internal/repo"
	"github.com/stretchr/testify/suite"
)

type TransactionMgrSuite struct {
	suite.Suite
	client  rpcx.Client
	from    *types.Address
	pk      crypto.PrivateKey
	client0 *ChainClient
	client1 *ChainClient
	client2 *ChainClient
}

type ChainClient struct {
	client rpcx.Client
	addr   string
	pk     crypto.PrivateKey
}

//init
func (suite *TransactionMgrSuite) SetupTest() {
	suite.client0 = suite.genChainClient()
	suite.client1 = suite.genChainClient()
	suite.client2 = suite.genChainClient()

	suite.RegisterAppchain(suite.client0)
	suite.RegisterAppchain(suite.client1)
	suite.RegisterAppchain(suite.client2)
	suite.RegisterRule(suite.client0, "./testdata/simple_rule.wasm")
	suite.RegisterRule(suite.client1, "./testdata/simple_rule.wasm")
}
func (suite *TransactionMgrSuite) SetupSuite() {
	keyPath, err := repo.Node1Path()
	suite.Require().Nil(err)
	pk, err := asym.RestorePrivateKey(keyPath, repo.KeyPassword)
	suite.Require().Nil(err)
	node0 := &rpcx.NodeInfo{Addr: cfg.addrs[0]}
	client, err := rpcx.New(
		rpcx.WithNodesInfo(node0),
		rpcx.WithLogger(cfg.logger),
		rpcx.WithPrivateKey(pk),
	)
	from, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	suite.client = client
	suite.from = from
	suite.pk = pk

	res, err := suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "Init", nil, rpcx.String("did:bitxhub:relayroot:"+suite.from.String()))
	suite.Require().Nil(err)

	node2, err := repo.Node2Path()
	suite.Require().Nil(err)

	key, err := asym.RestorePrivateKey(node2, repo.KeyPassword)
	suite.Require().Nil(err)

	node2Addr, err := key.PublicKey().Address()
	suite.Require().Nil(err)

	res, err = suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "AddAdmin", nil, rpcx.String("did:bitxhub:relayroot:"+suite.from.String()), pb.String("did:bitxhub:relayroot:"+node2Addr.String()))
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))

	node3, err := repo.Node3Path()
	suite.Require().Nil(err)

	key, err = asym.RestorePrivateKey(node3, repo.KeyPassword)
	suite.Require().Nil(err)

	node3Addr, err := key.PublicKey().Address()
	suite.Require().Nil(err)

	res, err = suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "AddAdmin", nil, rpcx.String("did:bitxhub:relayroot:"+suite.from.String()), pb.String("did:bitxhub:relayroot:"+node3Addr.String()))
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))

	node4, err := repo.Node4Path()
	suite.Require().Nil(err)

	key, err = asym.RestorePrivateKey(node4, repo.KeyPassword)
	suite.Require().Nil(err)

	node4Addr, err := key.PublicKey().Address()
	suite.Require().Nil(err)

	res, err = suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "AddAdmin", nil, rpcx.String("did:bitxhub:relayroot:"+suite.from.String()), pb.String("did:bitxhub:relayroot:"+node4Addr.String()))
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
}

func (suite *TransactionMgrSuite) genChainClient() *ChainClient {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)

	addr, err := pk.PublicKey().Address()
	suite.Require().Nil(err)

	node0 := &rpcx.NodeInfo{Addr: cfg.addrs[0]}
	client, err := rpcx.New(
		rpcx.WithNodesInfo(node0),
		rpcx.WithLogger(cfg.logger),
		rpcx.WithPrivateKey(pk),
	)
	suite.Require().Nil(err)

	return &ChainClient{
		client: client,
		addr:   "did:bitxhub:appchain" + addr.String() + ":.",
		pk:     pk,
	}
}
func (suite *TransactionMgrSuite) VotePass(id string) error {
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

func (suite *TransactionMgrSuite) vote(key crypto.PrivateKey, args ...*pb.Arg) (*pb.Receipt, error) {
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

	tx := &pb.Transaction{
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

func (suite *TransactionMgrSuite) GetChainStatusById(id string) (*pb.Receipt, error) {
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

	tx := &pb.Transaction{
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

func (suite *TransactionMgrSuite) RegisterAppchain(client *ChainClient) {
	pubAddress, err := client.pk.PublicKey().Address()
	suite.Require().Nil(err)

	client.client.SetPrivateKey(client.pk)
	bytes, err := client.pk.PublicKey().Bytes()
	suite.Require().Nil(err)
	var pubKeyStr = base64.StdEncoding.EncodeToString(bytes)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":" + pubAddress.String()), //id
		rpcx.String("did:bitxhub:appchain" + pubAddress.String() + ":."),                      //ownerDID
		rpcx.String("AppChain"), //name
		rpcx.String("/ipfs/QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"), //docAddr
		rpcx.String("QmQVxzUqN2Yv2UHUQXYwH8dSNkM8ReJ9qPqwJsf8zzoNUi"),       //docHash
		rpcx.String(""),                 //validators
		rpcx.String("raft"),             //consensus_type
		rpcx.String("hyperchain"),       //chain_type
		rpcx.String("Appchain for tax"), //desc
		rpcx.String("1.8"),              //version
		rpcx.String(pubKeyStr),          //public key
	}
	res, err := client.client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "Register", nil, args...)
	suite.Require().Nil(err)
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	suite.Require().Nil(err)
	suite.Require().NotNil(result.ChainID)
	err = suite.VotePass(result.ProposalID)
	suite.Require().Nil(err)

	res, err = suite.GetChainStatusById(result.ChainID)
	suite.Require().Nil(err)
	appchain := &rpcx.Appchain{}
	err = json.Unmarshal(res.Ret, appchain)
	suite.Require().Nil(err)
	suite.Require().Equal(appchain_mgr.AppchainAvailable, appchain.Status)
}

func (suite *TransactionMgrSuite) RegisterRule(client *ChainClient, ruleFile string) {
	client.client.SetPrivateKey(client.pk)

	from, err := client.pk.PublicKey().Address()
	suite.Require().Nil(err)
	ChainID := "did:bitxhub:appchain" + from.String() + ":."

	// deploy rule
	bytes, err := ioutil.ReadFile(ruleFile)
	suite.Require().Nil(err)
	addr, err := client.client.DeployContract(bytes, nil)
	suite.Require().Nil(err)

	// register rule
	res, err := client.client.InvokeBVMContract(constant.RuleManagerContractAddr.Address(), "RegisterRule", nil, pb.String(ChainID), pb.String(addr.String()))
	suite.Require().Nil(err)
	suite.Require().True(res.IsSuccess())
}

//tc:一对一跨链交易，跨链交易成功执行，中继链事务管理合约中该事务状态为成功
func (suite *TransactionMgrSuite) Test1001_One2One_AssetExchange_HappyPath() {
	index := uint64(1)
	from := suite.client0.addr
	to := suite.client1.addr

	aei := pb.AssetExchangeInfo{
		Id:            from + "123456",
		SenderOnSrc:   "aaa",
		ReceiverOnSrc: "bbb",
		AssetOnSrc:    1,
		SenderOnDst:   "BBB",
		ReceiverOnDst: "AAA",
		AssetOnDst:    10,
	}
	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	content, err := aei.Marshal()
	suite.Require().Nil(err)

	ib0 := &pb.IBTP{From: from, To: to, Index: index, Type: pb.IBTP_ASSET_EXCHANGE_INIT, Proof: proofHash[:], Extra: content}
	tx, _ := suite.client0.client.GenerateIBTPTx(ib0)
	tx.Extra = []byte(proof)
	res0, err := suite.client0.client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res0.Status)

	res, err := suite.client0.client.InvokeBVMContract(constant.AssetExchangeContractAddr.Address(), "GetStatus", nil, pb.String(aei.Id))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	suite.Require().Equal("0", string(res.Ret))

	ib1 := &pb.IBTP{From: to, To: from, Index: index, Type: pb.IBTP_ASSET_EXCHANGE_REDEEM, Proof: proofHash[:], Extra: []byte(aei.Id)}
	tx, _ = suite.client1.client.GenerateIBTPTx(ib1)
	tx.Extra = []byte(proof)
	res1, err := suite.client1.client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res1.Status)

	res, err = suite.client1.client.InvokeBVMContract(constant.AssetExchangeContractAddr.Address(), "GetStatus", nil, pb.String(aei.Id))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	suite.Require().Equal("1", string(res.Ret))
}

//tc:一对一跨链交易，来源链拒绝交易，中继链事务管理合约中该事务状态为错误
func (suite *TransactionMgrSuite) Test1002_One2One_AssetExchange_FromRefund() {
	index := uint64(1)
	from := suite.client0.addr
	to := suite.client1.addr

	aei := pb.AssetExchangeInfo{
		Id:            from + "123456",
		SenderOnSrc:   "aaa",
		ReceiverOnSrc: "bbb",
		AssetOnSrc:    1,
		SenderOnDst:   "BBB",
		ReceiverOnDst: "AAA",
		AssetOnDst:    10,
	}
	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	content, err := aei.Marshal()
	suite.Require().Nil(err)

	ib0 := &pb.IBTP{From: from, To: to, Index: index, Type: pb.IBTP_ASSET_EXCHANGE_INIT, Proof: proofHash[:], Extra: content}
	tx, _ := suite.client0.client.GenerateIBTPTx(ib0)
	tx.Extra = []byte(proof)
	res0, err := suite.client0.client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().True(res0.IsSuccess())

	ib1 := &pb.IBTP{From: from, To: to, Index: index + 1, Type: pb.IBTP_ASSET_EXCHANGE_REFUND, Proof: proofHash[:], Extra: []byte(aei.Id)}
	tx, _ = suite.client0.client.GenerateIBTPTx(ib1)
	tx.Extra = []byte(proof)
	res1, err := suite.client0.client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().True(res1.IsSuccess())

	res, err := suite.client0.client.InvokeBVMContract(constant.AssetExchangeContractAddr.Address(), "GetStatus", nil, pb.String(aei.Id))
	suite.Require().Nil(err)
	suite.Require().True(res.IsSuccess())
	suite.Require().Equal("2", string(res.Ret))

	ib1 = &pb.IBTP{From: to, To: from, Index: index, Type: pb.IBTP_ASSET_EXCHANGE_REDEEM, Proof: proofHash[:], Extra: []byte(aei.Id)}

	tx, _ = suite.client1.client.GenerateIBTPTx(ib1)
	tx.Extra = []byte(proof)
	res1, err = suite.client1.client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res1.Status)
}

//tc:一对一跨链交易，目的链拒绝交易，中继链事务管理合约中该事务状态为错误
func (suite *TransactionMgrSuite) Test1003_One2One_AssetExchange_ToRefund() {
	index := uint64(1)
	from := suite.client0.addr
	to := suite.client1.addr

	aei := pb.AssetExchangeInfo{
		Id:            from + "123456",
		SenderOnSrc:   "aaa",
		ReceiverOnSrc: "bbb",
		AssetOnSrc:    1,
		SenderOnDst:   "BBB",
		ReceiverOnDst: "AAA",
		AssetOnDst:    10,
	}
	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	content, err := aei.Marshal()
	suite.Require().Nil(err)

	ib0 := &pb.IBTP{From: from, To: to, Index: index, Type: pb.IBTP_ASSET_EXCHANGE_INIT, Proof: proofHash[:], Extra: content}
	tx, _ := suite.client0.client.GenerateIBTPTx(ib0)
	tx.Extra = []byte(proof)
	res0, err := suite.client0.client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().True(res0.IsSuccess())

	ib1 := &pb.IBTP{From: to, To: from, Index: index, Type: pb.IBTP_ASSET_EXCHANGE_REFUND, Proof: proofHash[:], Extra: []byte(aei.Id)}
	tx, _ = suite.client1.client.GenerateIBTPTx(ib1)
	tx.Extra = []byte(proof)
	res1, err := suite.client1.client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().True(res1.IsSuccess())

	res, err := suite.client0.client.InvokeBVMContract(constant.AssetExchangeContractAddr.Address(), "GetStatus", nil, pb.String(aei.Id))
	suite.Require().Nil(err)
	suite.Require().True(res.IsSuccess())
	suite.Require().Equal("2", string(res.Ret))

	ib1 = &pb.IBTP{From: to, To: from, Index: index + 1, Type: pb.IBTP_ASSET_EXCHANGE_REDEEM, Proof: proofHash[:], Extra: []byte(aei.Id)}
	tx, _ = suite.client1.client.GenerateIBTPTx(ib1)
	tx.Extra = []byte(proof)
	res1, err = suite.client1.client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res1.Status)
}

//tc:一对一跨链交易，获取资产交换的多签
func (suite *TransactionMgrSuite) Test1004_AssetExchange_Signs() {
	index := uint64(1)
	from := suite.client0.addr
	to := suite.client1.addr

	aei := pb.AssetExchangeInfo{
		Id:            from + "123456",
		SenderOnSrc:   "aaa",
		ReceiverOnSrc: "bbb",
		AssetOnSrc:    1,
		SenderOnDst:   "BBB",
		ReceiverOnDst: "AAA",
		AssetOnDst:    10,
	}
	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	content, err := aei.Marshal()
	suite.Require().Nil(err)

	ib0 := &pb.IBTP{From: from, To: to, Index: index, Type: pb.IBTP_ASSET_EXCHANGE_INIT, Proof: proofHash[:], Extra: content}
	tx, _ := suite.client0.client.GenerateIBTPTx(ib0)
	tx.Extra = []byte(proof)
	res0, err := suite.client0.client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().True(res0.IsSuccess())

	resp, err := suite.client0.client.GetMultiSigns(aei.Id, pb.GetMultiSignsRequest_ASSET_EXCHANGE)
	suite.Require().Nil(err)
	suite.Require().NotNil(resp)
	suite.Require().Equal(4, len(resp.Sign))

	msg := fmt.Sprintf("%s-%d", aei.Id, 0)
	digest := sha256.Sum256([]byte(msg))

	for validator, sign := range resp.Sign {
		ok, err := asym.Verify(crypto.Secp256k1, sign, digest[:], *types.NewAddressByStr(validator))
		suite.Require().Nil(err)
		suite.Require().True(ok)
		fmt.Println(validator)
	}
}

//tc:一对一跨链交易，来源链和目的链都拒绝交易，中继链事务管理合约中该事务状态为错误
func (suite *TransactionMgrSuite) Test1005_One2One_AssetExchange_FromToRefund() {
	index := uint64(1)
	from := suite.client0.addr
	to := suite.client1.addr

	aei := pb.AssetExchangeInfo{
		Id:            from + "123456",
		SenderOnSrc:   "aaa",
		ReceiverOnSrc: "bbb",
		AssetOnSrc:    1,
		SenderOnDst:   "BBB",
		ReceiverOnDst: "AAA",
		AssetOnDst:    10,
	}
	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	content, err := aei.Marshal()
	suite.Require().Nil(err)

	ib0 := &pb.IBTP{From: from, To: to, Index: index, Type: pb.IBTP_ASSET_EXCHANGE_INIT, Proof: proofHash[:], Extra: content}
	tx, _ := suite.client0.client.GenerateIBTPTx(ib0)
	tx.Extra = []byte(proof)
	res0, err := suite.client0.client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().True(res0.IsSuccess())

	ib1 := &pb.IBTP{From: from, To: to, Index: index + 1, Type: pb.IBTP_ASSET_EXCHANGE_REFUND, Proof: proofHash[:], Extra: []byte(aei.Id)}
	tx, _ = suite.client0.client.GenerateIBTPTx(ib1)
	tx.Extra = []byte(proof)
	res1, err := suite.client0.client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().True(res1.IsSuccess())

	ib2 := &pb.IBTP{From: to, To: from, Index: index, Type: pb.IBTP_ASSET_EXCHANGE_REFUND, Proof: proofHash[:], Extra: []byte(aei.Id)}
	tx, _ = suite.client1.client.GenerateIBTPTx(ib2)
	tx.Extra = []byte(proof)
	res2, err := suite.client1.client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res2.Status, string(res2.Ret))

	res, err := suite.client0.client.InvokeBVMContract(constant.AssetExchangeContractAddr.Address(), "GetStatus", nil, pb.String(aei.Id))
	suite.Require().Nil(err)
	suite.Require().True(res.IsSuccess())
	suite.Require().Equal("2", string(res.Ret))

	ib1 = &pb.IBTP{From: to, To: from, Index: index + 1, Type: pb.IBTP_ASSET_EXCHANGE_REDEEM, Proof: proofHash[:], Extra: []byte(aei.Id)}

	tx, _ = suite.client1.client.GenerateIBTPTx(ib1)
	tx.Extra = []byte(proof)
	res1, err = suite.client1.client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res1.Status, string(res1.Ret))
}

//tc:一对一跨链交易，相同id重复注册，中继链事务管理合约中该事务状态为错误
func (suite *TransactionMgrSuite) Test1006_One2One_AssetExchange_SameId() {
	index := uint64(1)
	from := suite.client0.addr
	to := suite.client1.addr

	aei := pb.AssetExchangeInfo{
		Id:            from + "123456",
		SenderOnSrc:   "aaa",
		ReceiverOnSrc: "bbb",
		AssetOnSrc:    1,
		SenderOnDst:   "BBB",
		ReceiverOnDst: "AAA",
		AssetOnDst:    10,
	}
	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	content, err := aei.Marshal()
	suite.Require().Nil(err)

	ib0 := &pb.IBTP{From: from, To: to, Index: index, Type: pb.IBTP_ASSET_EXCHANGE_INIT, Proof: proofHash[:], Extra: content}
	tx, _ := suite.client0.client.GenerateIBTPTx(ib0)
	tx.Extra = []byte(proof)
	res0, err := suite.client0.client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().True(res0.IsSuccess())

	res, err := suite.client0.client.InvokeBVMContract(constant.AssetExchangeContractAddr.Address(), "GetStatus", nil, pb.String(aei.Id))
	suite.Require().Nil(err)
	suite.Require().True(res.IsSuccess())
	suite.Require().Equal("0", string(res.Ret))

	ib1 := &pb.IBTP{From: to, To: from, Index: index, Type: pb.IBTP_ASSET_EXCHANGE_REDEEM, Proof: proofHash[:], Extra: []byte(aei.Id)}
	tx, _ = suite.client1.client.GenerateIBTPTx(ib1)
	tx.Extra = []byte(proof)
	res1, err := suite.client1.client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res1.Status)

	res, err = suite.client0.client.InvokeBVMContract(constant.AssetExchangeContractAddr.Address(), "GetStatus", nil, pb.String(aei.Id))
	suite.Require().Nil(err)
	suite.Require().True(res.IsSuccess())
	suite.Require().Equal("1", string(res.Ret))

	//same id
	content, err = aei.Marshal()
	suite.Require().Nil(err)

	ib0 = &pb.IBTP{From: from, To: to, Index: index + 1, Type: pb.IBTP_ASSET_EXCHANGE_INIT, Proof: proofHash[:], Extra: content}
	tx, _ = suite.client0.client.GenerateIBTPTx(ib0)
	tx.Extra = []byte(proof)
	res0, err = suite.client0.client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res0.Status)

	res, err = suite.client0.client.InvokeBVMContract(constant.AssetExchangeContractAddr.Address(), "GetStatus", nil, pb.String(aei.Id))
	suite.Require().Nil(err)
	suite.Require().True(res.IsSuccess())

	ib1 = &pb.IBTP{From: to, To: from, Index: index + 1, Type: pb.IBTP_ASSET_EXCHANGE_REDEEM, Proof: proofHash[:], Extra: []byte(aei.Id)}
	tx, _ = suite.client1.client.GenerateIBTPTx(ib1)
	tx.Extra = []byte(proof)
	res1, err = suite.client1.client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res1.Status)

	res, err = suite.client0.client.InvokeBVMContract(constant.AssetExchangeContractAddr.Address(), "GetStatus", nil, pb.String(aei.Id))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

}

//tc:一对一跨链交易，目的账户和来源账户相同，中继链事务管理合约中该事务状态为成功
func (suite *TransactionMgrSuite) Test1007_One2One_AssetExchange_FromToFrom() {
	index := uint64(1)
	from := suite.client0.addr
	to := suite.client1.addr

	aei := pb.AssetExchangeInfo{
		Id:            from + "123456",
		SenderOnSrc:   "aaa",
		ReceiverOnSrc: "aaa",
		AssetOnSrc:    1,
		SenderOnDst:   "AAA",
		ReceiverOnDst: "AAA",
		AssetOnDst:    10,
	}
	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	content, err := aei.Marshal()
	suite.Require().Nil(err)

	ib0 := &pb.IBTP{From: from, To: to, Index: index, Type: pb.IBTP_ASSET_EXCHANGE_INIT, Proof: proofHash[:], Extra: content}
	tx, _ := suite.client0.client.GenerateIBTPTx(ib0)
	tx.Extra = []byte(proof)
	res0, err := suite.client0.client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().True(res0.IsSuccess())

	res, err := suite.client0.client.InvokeBVMContract(constant.AssetExchangeContractAddr.Address(), "GetStatus", nil, pb.String(aei.Id))
	suite.Require().Nil(err)
	suite.Require().True(res.IsSuccess())
	suite.Require().Equal("0", string(res.Ret))

	ib1 := &pb.IBTP{From: to, To: from, Index: index, Type: pb.IBTP_ASSET_EXCHANGE_REDEEM, Proof: proofHash[:], Extra: []byte(aei.Id)}
	tx, _ = suite.client1.client.GenerateIBTPTx(ib1)
	tx.Extra = []byte(proof)
	res1, err := suite.client1.client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().True(res1.IsSuccess())

	res, err = suite.client0.client.InvokeBVMContract(constant.AssetExchangeContractAddr.Address(), "GetStatus", nil, pb.String(aei.Id))
	suite.Require().Nil(err)
	suite.Require().True(res.IsSuccess())
	suite.Require().Equal("1", string(res.Ret))
}

//tc:一对一跨链交易，缺少id字段，中继链事务管理合约中该事务状态为失败
func (suite *TransactionMgrSuite) Test1008_One2One_AssetExchange_LoseFieldId() {
	index := uint64(1)
	from := suite.client0.addr
	to := suite.client1.addr

	aei := pb.AssetExchangeInfo{
		SenderOnSrc:   "aaa",
		ReceiverOnSrc: "bbb",
		AssetOnSrc:    1,
		SenderOnDst:   "BBB",
		ReceiverOnDst: "AAA",
		AssetOnDst:    10,
	}
	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	content, err := aei.Marshal()
	suite.Require().Nil(err)

	ib0 := &pb.IBTP{From: from, To: to, Index: index, Type: pb.IBTP_ASSET_EXCHANGE_INIT, Proof: proofHash[:], Extra: content}
	tx, _ := suite.client0.client.GenerateIBTPTx(ib0)
	tx.Extra = []byte(proof)
	res0, err := suite.client0.client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res0.Status)

	res, err := suite.client0.client.InvokeBVMContract(constant.AssetExchangeContractAddr.Address(), "GetStatus", nil, pb.String(aei.Id))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)

	ib1 := &pb.IBTP{From: to, To: from, Index: index, Type: pb.IBTP_ASSET_EXCHANGE_REDEEM, Proof: proofHash[:], Extra: []byte(aei.Id)}
	tx, _ = suite.client1.client.GenerateIBTPTx(ib1)
	tx.Extra = []byte(proof)
	res1, err := suite.client1.client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res1.Status)

	res, err = suite.client0.client.InvokeBVMContract(constant.AssetExchangeContractAddr.Address(), "GetStatus", nil, pb.String(aei.Id))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
}

//tc:一对一跨链交易，缺少sender字段，中继链事务管理合约中该事务状态为失败
func (suite *TransactionMgrSuite) Test1009_One2One_AssetExchange_LoseFieldSender() {
	index := uint64(1)
	from := suite.client0.addr
	to := suite.client1.addr

	aei := pb.AssetExchangeInfo{
		Id:            from + "123456",
		ReceiverOnSrc: "bbb",
		AssetOnSrc:    1,
		SenderOnDst:   "BBB",
		ReceiverOnDst: "AAA",
		AssetOnDst:    10,
	}
	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	content, err := aei.Marshal()
	suite.Require().Nil(err)

	ib0 := &pb.IBTP{From: from, To: to, Index: index, Type: pb.IBTP_ASSET_EXCHANGE_INIT, Proof: proofHash[:], Extra: content}
	tx, _ := suite.client0.client.GenerateIBTPTx(ib0)
	tx.Extra = []byte(proof)
	res0, err := suite.client0.client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res0.Status)

	res, err := suite.client0.client.InvokeBVMContract(constant.AssetExchangeContractAddr.Address(), "GetStatus", nil, pb.String(aei.Id))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)

	ib1 := &pb.IBTP{From: to, To: from, Index: index, Type: pb.IBTP_ASSET_EXCHANGE_REDEEM, Proof: proofHash[:], Extra: []byte(aei.Id)}
	tx, _ = suite.client1.client.GenerateIBTPTx(ib1)
	tx.Extra = []byte(proof)
	res1, err := suite.client1.client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res1.Status)

	res, err = suite.client0.client.InvokeBVMContract(constant.AssetExchangeContractAddr.Address(), "GetStatus", nil, pb.String(aei.Id))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
}

//tc:一对一跨链交易，缺少asset字段，中继链事务管理合约中该事务状态为失败
func (suite *TransactionMgrSuite) Test1010_One2One_AssetExchange_LoseFieldAsset() {
	index := uint64(1)
	from := suite.client0.addr
	to := suite.client1.addr

	aei := pb.AssetExchangeInfo{
		Id:            from + "123456",
		SenderOnSrc:   "aaa",
		ReceiverOnSrc: "bbb",
		SenderOnDst:   "BBB",
		ReceiverOnDst: "AAA",
		AssetOnDst:    10,
	}
	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	content, err := aei.Marshal()
	suite.Require().Nil(err)

	ib0 := &pb.IBTP{From: from, To: to, Index: index, Type: pb.IBTP_ASSET_EXCHANGE_INIT, Proof: proofHash[:], Extra: content}
	tx, _ := suite.client0.client.GenerateIBTPTx(ib0)
	tx.Extra = []byte(proof)
	res0, err := suite.client0.client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res0.Status)

	res, err := suite.client0.client.InvokeBVMContract(constant.AssetExchangeContractAddr.Address(), "GetStatus", nil, pb.String(aei.Id))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)

	ib1 := &pb.IBTP{From: to, To: from, Index: index, Type: pb.IBTP_ASSET_EXCHANGE_REDEEM, Proof: proofHash[:], Extra: []byte(aei.Id)}
	tx, _ = suite.client1.client.GenerateIBTPTx(ib1)
	tx.Extra = []byte(proof)
	res1, err := suite.client1.client.SendTransactionWithReceipt(tx, nil)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res1.Status)

	res, err = suite.client0.client.InvokeBVMContract(constant.AssetExchangeContractAddr.Address(), "GetStatus", nil, pb.String(aei.Id))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
}
