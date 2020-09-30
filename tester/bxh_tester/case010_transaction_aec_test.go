package bxh_tester

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-kit/types"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/stretchr/testify/suite"
)

type TransactionMgrSuite struct {
	suite.Suite
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

	suite.RegisterAppchain(suite.client0, "hyperchain")
	suite.RegisterAppchain(suite.client1, "hyperchain")
	suite.RegisterAppchain(suite.client2, "fabric")
	suite.RegisterRule(suite.client0, "./testdata/simple_rule.wasm")
	suite.RegisterRule(suite.client1, "./testdata/simple_rule.wasm")
}

func (suite *TransactionMgrSuite) genChainClient() *ChainClient {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)

	addr, err := pk.PublicKey().Address()
	suite.Require().Nil(err)

	client, err := rpcx.New(
		rpcx.WithAddrs(cfg.addrs),
		rpcx.WithLogger(cfg.logger),
		rpcx.WithPrivateKey(pk),
	)
	suite.Require().Nil(err)

	return &ChainClient{
		client: client,
		addr:   addr.String(),
		pk:     pk,
	}
}

func (suite *TransactionMgrSuite) RegisterAppchain(client *ChainClient, chainType string) {
	pubBytes, err := client.pk.PublicKey().Bytes()
	suite.Require().Nil(err)

	client.client.SetPrivateKey(client.pk)
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
	res, err := client.client.InvokeBVMContract(rpcx.AppchainMgrContractAddr, "Register",nil, args...)
	suite.Require().Nil(err)
	appChain := &rpcx.Appchain{}
	err = json.Unmarshal(res.Ret, appChain)
	suite.Require().Nil(err)
	suite.Require().NotNil(appChain.ID)
}

func (suite *TransactionMgrSuite) RegisterRule(client *ChainClient, ruleFile string) {
	client.client.SetPrivateKey(client.pk)

	from, err := client.pk.PublicKey().Address()
	suite.Require().Nil(err)

	// deploy rule
	bytes, err := ioutil.ReadFile(ruleFile)
	suite.Require().Nil(err)
	addr, err := client.client.DeployContract(bytes,nil)
	suite.Require().Nil(err)

	// register rule
	res, err := client.client.InvokeBVMContract(rpcx.RuleManagerContractAddr, "RegisterRule",nil, pb.String(from.Hex()), pb.String(addr.Hex()))
	suite.Require().Nil(err)
	suite.Require().True(res.IsSuccess())
}

func (suite *TransactionMgrSuite) Test001_One2One_AssetExchange_HappyPath() {
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
	data0, err := ib0.Marshal()
	suite.Require().Nil(err)

	tx, _ := suite.client0.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data0))
	tx.Extra = []byte(proof)
	res0, err := suite.client0.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib0.From + ib0.To,
		IBTPNonce: ib0.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res0.Status, string(res0.Ret))

	res, err := suite.client0.client.InvokeBVMContract(rpcx.AssetExchangeContractAddr, "GetStatus",nil, pb.String(aei.Id))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status, string(res.Ret))
	suite.Require().Equal("0", string(res.Ret))

	ib1 := &pb.IBTP{From: to, To: from, Index: index, Type: pb.IBTP_ASSET_EXCHANGE_REDEEM, Proof: proofHash[:], Extra: []byte(aei.Id)}
	data1, err := ib1.Marshal()
	suite.Require().Nil(err)

	tx, _ = suite.client1.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data1))
	tx.Extra = []byte(proof)
	res1, err := suite.client1.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib1.From + ib1.To,
		IBTPNonce: ib1.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res1.Status, string(res1.Ret))

	res, err = suite.client0.client.InvokeBVMContract(rpcx.AssetExchangeContractAddr, "GetStatus",nil, pb.String(aei.Id))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status, string(res.Ret))
	suite.Require().Equal("1", string(res.Ret))
}

func (suite *TransactionMgrSuite) Test002_One2One_AssetExchange_FromRefund() {
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
	data0, err := ib0.Marshal()
	suite.Require().Nil(err)

	tx, _ := suite.client0.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data0))
	tx.Extra = []byte(proof)
	res0, err := suite.client0.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib0.From + ib0.To,
		IBTPNonce: ib0.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res0.Status, string(res0.Ret))

	ib1 := &pb.IBTP{From: from, To: to, Index: index + 1, Type: pb.IBTP_ASSET_EXCHANGE_REFUND, Proof: proofHash[:], Extra: []byte(aei.Id)}
	data1, err := ib1.Marshal()
	suite.Require().Nil(err)

	tx, _ = suite.client0.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data1))
	tx.Extra = []byte(proof)
	res1, err := suite.client0.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib1.From + ib1.To,
		IBTPNonce: ib1.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res1.Status, string(res1.Ret))

	res, err := suite.client0.client.InvokeBVMContract(rpcx.AssetExchangeContractAddr, "GetStatus",nil, pb.String(aei.Id))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status, string(res.Ret))
	suite.Require().Equal("2", string(res.Ret))

	ib1 = &pb.IBTP{From: to, To: from, Index: index, Type: pb.IBTP_ASSET_EXCHANGE_REDEEM, Proof: proofHash[:], Extra: []byte(aei.Id)}
	data1, err = ib1.Marshal()
	suite.Require().Nil(err)

	tx, _ = suite.client1.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data1))
	tx.Extra = []byte(proof)
	res1, err = suite.client1.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib1.From + ib1.To,
		IBTPNonce: ib1.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res1.Status, string(res1.Ret))
}

func (suite *TransactionMgrSuite) Test003_One2One_AssetExchange_ToRefund() {
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
	data0, err := ib0.Marshal()
	suite.Require().Nil(err)

	tx, _ := suite.client0.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data0))
	tx.Extra = []byte(proof)
	res0, err := suite.client0.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib0.From + ib0.To,
		IBTPNonce: ib0.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res0.Status, string(res0.Ret))

	ib1 := &pb.IBTP{From: to, To: from, Index: index, Type: pb.IBTP_ASSET_EXCHANGE_REFUND, Proof: proofHash[:], Extra: []byte(aei.Id)}
	data1, err := ib1.Marshal()
	suite.Require().Nil(err)

	tx, _ = suite.client1.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data1))
	tx.Extra = []byte(proof)
	res1, err := suite.client1.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib1.From + ib1.To,
		IBTPNonce: ib1.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res1.Status, string(res1.Ret))

	res, err := suite.client0.client.InvokeBVMContract(rpcx.AssetExchangeContractAddr, "GetStatus",nil, pb.String(aei.Id))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status, string(res.Ret))
	suite.Require().Equal("2", string(res.Ret))

	ib1 = &pb.IBTP{From: to, To: from, Index: index + 1, Type: pb.IBTP_ASSET_EXCHANGE_REDEEM, Proof: proofHash[:], Extra: []byte(aei.Id)}
	data1, err = ib1.Marshal()
	suite.Require().Nil(err)

	tx, _ = suite.client1.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data1))
	tx.Extra = []byte(proof)
	res1, err = suite.client1.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib1.From + ib1.To,
		IBTPNonce: ib1.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res1.Status, string(res1.Ret))
}

func (suite *TransactionMgrSuite) Test004_AssetExchange_Signs() {
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
	data0, err := ib0.Marshal()
	suite.Require().Nil(err)

	tx, _ := suite.client0.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data0))
	tx.Extra = []byte(proof)
	res0, err := suite.client0.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib0.From + ib0.To,
		IBTPNonce: ib0.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res0.Status, string(res0.Ret))

	resp, err := suite.client0.client.GetMultiSigns(aei.Id, pb.GetMultiSignsRequest_ASSET_EXCHANGE)
	suite.Require().Nil(err)
	suite.Require().NotNil(resp)
	suite.Require().Equal(4, len(resp.Sign))

	msg := fmt.Sprintf("%s-%d", aei.Id, 0)
	digest := sha256.Sum256([]byte(msg))

	for validator, sign := range resp.Sign {
		ok, err := asym.Verify(crypto.Secp256k1, sign, digest[:], types.String2Address(validator))
		suite.Require().Nil(err)
		suite.Require().True(ok)
		fmt.Println(validator)
	}
}

func (suite *TransactionMgrSuite) Test005_One2One_AssetExchange_FromToRefund() {
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
	data0, err := ib0.Marshal()
	suite.Require().Nil(err)

	tx, _ := suite.client0.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data0))
	tx.Extra = []byte(proof)
	res0, err := suite.client0.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib0.From + ib0.To,
		IBTPNonce: ib0.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res0.Status, string(res0.Ret))

	ib1 := &pb.IBTP{From: from, To: to, Index: index + 1, Type: pb.IBTP_ASSET_EXCHANGE_REFUND, Proof: proofHash[:], Extra: []byte(aei.Id)}
	data1, err := ib1.Marshal()
	suite.Require().Nil(err)

	tx, _ = suite.client0.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data1))
	tx.Extra = []byte(proof)
	res1, err := suite.client0.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib1.From + ib1.To,
		IBTPNonce: ib1.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res1.Status, string(res1.Ret))

	ib2 := &pb.IBTP{From: to, To: from, Index: index, Type: pb.IBTP_ASSET_EXCHANGE_REFUND, Proof: proofHash[:], Extra: []byte(aei.Id)}
	data2, err := ib2.Marshal()
	suite.Require().Nil(err)

	tx, _ = suite.client0.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data2))
	tx.Extra = []byte(proof)
	res2, err := suite.client0.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib2.From + ib2.To,
		IBTPNonce: ib2.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res2.Status, string(res2.Ret))

	res, err := suite.client0.client.InvokeBVMContract(rpcx.AssetExchangeContractAddr, "GetStatus",nil, pb.String(aei.Id))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status, string(res.Ret))
	suite.Require().Equal("2", string(res.Ret))

	ib1 = &pb.IBTP{From: to, To: from, Index: index, Type: pb.IBTP_ASSET_EXCHANGE_REDEEM, Proof: proofHash[:], Extra: []byte(aei.Id)}
	data1, err = ib1.Marshal()
	suite.Require().Nil(err)

	tx, _ = suite.client1.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data1))
	tx.Extra = []byte(proof)
	res1, err = suite.client1.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib1.From + ib1.To,
		IBTPNonce: ib1.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res1.Status, string(res1.Ret))
}

func (suite *TransactionMgrSuite) Test006_One2One_AssetExchange_SameId() {
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
	data0, err := ib0.Marshal()
	suite.Require().Nil(err)

	tx, _ := suite.client0.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data0))
	tx.Extra = []byte(proof)
	res0, err := suite.client0.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib0.From + ib0.To,
		IBTPNonce: ib0.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res0.Status, string(res0.Ret))

	res, err := suite.client0.client.InvokeBVMContract(rpcx.AssetExchangeContractAddr, "GetStatus",nil, pb.String(aei.Id))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status, string(res.Ret))
	suite.Require().Equal("0", string(res.Ret))

	ib1 := &pb.IBTP{From: to, To: from, Index: index, Type: pb.IBTP_ASSET_EXCHANGE_REDEEM, Proof: proofHash[:], Extra: []byte(aei.Id)}
	data1, err := ib1.Marshal()
	suite.Require().Nil(err)

	tx, _ = suite.client1.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data1))
	tx.Extra = []byte(proof)
	res1, err := suite.client1.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib1.From + ib1.To,
		IBTPNonce: ib1.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res1.Status, string(res1.Ret))

	res, err = suite.client0.client.InvokeBVMContract(rpcx.AssetExchangeContractAddr, "GetStatus",nil, pb.String(aei.Id))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status, string(res.Ret))
	suite.Require().Equal("1", string(res.Ret))

	//same id
	content, err = aei.Marshal()
	suite.Require().Nil(err)

	ib0 = &pb.IBTP{From: from, To: to, Index: index+1, Type: pb.IBTP_ASSET_EXCHANGE_INIT, Proof: proofHash[:], Extra: content}
	data0, err = ib0.Marshal()
	suite.Require().Nil(err)


	tx, _ = suite.client0.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data0))
	tx.Extra = []byte(proof)
	res0, err = suite.client0.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib0.From + ib0.To,
		IBTPNonce: ib0.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res0.Status, string(res0.Ret))

	res, err = suite.client0.client.InvokeBVMContract(rpcx.AssetExchangeContractAddr, "GetStatus",nil, pb.String(aei.Id))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status, string(res.Ret))

	ib1 = &pb.IBTP{From: to, To: from, Index: index+1, Type: pb.IBTP_ASSET_EXCHANGE_REDEEM, Proof: proofHash[:], Extra: []byte(aei.Id)}
	data1, err = ib1.Marshal()
	suite.Require().Nil(err)

	tx, _ = suite.client1.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data1))
	tx.Extra = []byte(proof)
	res1, err = suite.client1.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib1.From + ib1.To,
		IBTPNonce: ib1.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res1.Status, string(res1.Ret))

	res, err = suite.client0.client.InvokeBVMContract(rpcx.AssetExchangeContractAddr, "GetStatus",nil, pb.String(aei.Id))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status, string(res.Ret))

}

func (suite *TransactionMgrSuite) Test007_One2One_AssetExchange_FromToFrom() {
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
	data0, err := ib0.Marshal()
	suite.Require().Nil(err)

	tx, _ := suite.client0.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data0))
	tx.Extra = []byte(proof)
	res0, err := suite.client0.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib0.From + ib0.To,
		IBTPNonce: ib0.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res0.Status, string(res0.Ret))

	res, err := suite.client0.client.InvokeBVMContract(rpcx.AssetExchangeContractAddr, "GetStatus",nil, pb.String(aei.Id))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status, string(res.Ret))
	suite.Require().Equal("0", string(res.Ret))

	ib1 := &pb.IBTP{From: to, To: from, Index: index, Type: pb.IBTP_ASSET_EXCHANGE_REDEEM, Proof: proofHash[:], Extra: []byte(aei.Id)}
	data1, err := ib1.Marshal()
	suite.Require().Nil(err)

	tx, _ = suite.client1.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data1))
	tx.Extra = []byte(proof)
	res1, err := suite.client1.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib1.From + ib1.To,
		IBTPNonce: ib1.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res1.Status, string(res1.Ret))

	res, err = suite.client0.client.InvokeBVMContract(rpcx.AssetExchangeContractAddr, "GetStatus",nil, pb.String(aei.Id))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status, string(res.Ret))
	suite.Require().Equal("1", string(res.Ret))
}

func (suite *TransactionMgrSuite) Test008_One2One_AssetExchange_LoseFieldId() {
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
	data0, err := ib0.Marshal()
	suite.Require().Nil(err)

	tx, _ := suite.client0.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data0))
	tx.Extra = []byte(proof)
	res0, err := suite.client0.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib0.From + ib0.To,
		IBTPNonce: ib0.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res0.Status, string(res0.Ret))

	res, err := suite.client0.client.InvokeBVMContract(rpcx.AssetExchangeContractAddr, "GetStatus",nil, pb.String(aei.Id))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status, string(res.Ret))

	ib1 := &pb.IBTP{From: to, To: from, Index: index, Type: pb.IBTP_ASSET_EXCHANGE_REDEEM, Proof: proofHash[:], Extra: []byte(aei.Id)}
	data1, err := ib1.Marshal()
	suite.Require().Nil(err)

	tx, _ = suite.client1.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data1))
	tx.Extra = []byte(proof)
	res1, err := suite.client1.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib1.From + ib1.To,
		IBTPNonce: ib1.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res1.Status, string(res1.Ret))

	res, err = suite.client0.client.InvokeBVMContract(rpcx.AssetExchangeContractAddr, "GetStatus",nil, pb.String(aei.Id))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status, string(res.Ret))
}

func (suite *TransactionMgrSuite) Test008_One2One_AssetExchange_LoseFieldSender() {
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
	data0, err := ib0.Marshal()
	suite.Require().Nil(err)

	tx, _ := suite.client0.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data0))
	tx.Extra = []byte(proof)
	res0, err := suite.client0.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib0.From + ib0.To,
		IBTPNonce: ib0.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res0.Status, string(res0.Ret))

	res, err := suite.client0.client.InvokeBVMContract(rpcx.AssetExchangeContractAddr, "GetStatus",nil, pb.String(aei.Id))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status, string(res.Ret))

	ib1 := &pb.IBTP{From: to, To: from, Index: index, Type: pb.IBTP_ASSET_EXCHANGE_REDEEM, Proof: proofHash[:], Extra: []byte(aei.Id)}
	data1, err := ib1.Marshal()
	suite.Require().Nil(err)

	tx, _ = suite.client1.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data1))
	tx.Extra = []byte(proof)
	res1, err := suite.client1.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib1.From + ib1.To,
		IBTPNonce: ib1.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res1.Status, string(res1.Ret))

	res, err = suite.client0.client.InvokeBVMContract(rpcx.AssetExchangeContractAddr, "GetStatus",nil, pb.String(aei.Id))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status, string(res.Ret))
}

func (suite *TransactionMgrSuite) Test009_One2One_AssetExchange_LoseFieldAsset() {
	index := uint64(1)
	from := suite.client0.addr
	to := suite.client1.addr

	aei := pb.AssetExchangeInfo{
		Id:            from + "123456",
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
	data0, err := ib0.Marshal()
	suite.Require().Nil(err)

	tx, _ := suite.client0.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data0))
	tx.Extra = []byte(proof)
	res0, err := suite.client0.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib0.From + ib0.To,
		IBTPNonce: ib0.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res0.Status, string(res0.Ret))

	res, err := suite.client0.client.InvokeBVMContract(rpcx.AssetExchangeContractAddr, "GetStatus",nil, pb.String(aei.Id))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status, string(res.Ret))

	ib1 := &pb.IBTP{From: to, To: from, Index: index, Type: pb.IBTP_ASSET_EXCHANGE_REDEEM, Proof: proofHash[:], Extra: []byte(aei.Id)}
	data1, err := ib1.Marshal()
	suite.Require().Nil(err)

	tx, _ = suite.client1.client.GenerateContractTx(pb.TransactionData_BVM, rpcx.InterchainContractAddr, "HandleIBTP", pb.Bytes(data1))
	tx.Extra = []byte(proof)
	res1, err := suite.client1.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: ib1.From + ib1.To,
		IBTPNonce: ib1.Index,
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res1.Status, string(res1.Ret))

	res, err = suite.client0.client.InvokeBVMContract(rpcx.AssetExchangeContractAddr, "GetStatus",nil, pb.String(aei.Id))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status, string(res.Ret))
}
