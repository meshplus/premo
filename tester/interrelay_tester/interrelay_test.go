package interrelay_tester

import (
	"crypto/sha256"
	"fmt"
	"github.com/bitxhub/bitxid"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"time"
)

var cfg2 = &config{
	addrs: []string{
		"localhost:50011",
	},
	logger: logrus.New(),
}

func (suite *Snake) Test001_InterRelay_Init_Relay1() {
	Relay1 := "0x454e2569dD093D09E5E8B4aB764692780D795C9a"
	ruleFile := "./testdata/simple_rule.wasm"
	bytes, err := ioutil.ReadFile(ruleFile)
	suite.Require().Nil(err)
	addr, err := suite.client.DeployContract(bytes, nil)
	suite.Require().Nil(err)
	res, err := suite.client.InvokeBVMContract(
		constant.RuleManagerContractAddr.Address(),
		"RegisterRule",
		nil,
		pb.String(Relay1),
		pb.String(addr.String()))
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().True(res.IsSuccess())

	adminDID := "did:bitxhub:relayroot:" + suite.from.String()
	res, err = suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "Init", nil, pb.String(adminDID))
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
}

func (suite *Snake) Test002_InterRelay_Init_Relay2() {
	kA, _, from, _ := suite.prepare()
	node0 := &rpcx.NodeInfo{Addr: cfg2.addrs[0]}
	client, err := rpcx.New(
		rpcx.WithNodesInfo(node0),
		rpcx.WithLogger(cfg2.logger),
		rpcx.WithPrivateKey(kA),
	)
	suite.Require().Nil(err)
	/****************************************************/
	Relay1 := "0x703b22368195d5063C5B5C26019301Cf2EbC83e2"
	ruleFile := "./testdata/simple_rule.wasm"
	bytes, err := ioutil.ReadFile(ruleFile)
	suite.Require().Nil(err)
	addr, err := client.DeployContract(bytes, nil)
	suite.Require().Nil(err)
	res, err := client.InvokeBVMContract(
		constant.RuleManagerContractAddr.Address(),
		"RegisterRule",
		nil,
		pb.String(Relay1),
		pb.String(addr.String()))
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().True(res.IsSuccess())
	/****************************************************/
	adminAddrStr := from.String()
	adminDID := "did:bitxhub:appchain001:" + adminAddrStr
	res, err = client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "Init", nil, pb.String(adminDID))
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
}

func (suite *Snake) Test003_InterRelay_HandleIBTP() {
	Relay1 := "0x454e2569dD093D09E5E8B4aB764692780D795C9a"
	Relay2 := "0x703b22368195d5063C5B5C26019301Cf2EbC83e2"
	proof := "test"
	proofHash := sha256.Sum256([]byte(proof))

	item := &bitxid.MethodItem{
		BasicItem: bitxid.BasicItem{
			ID:      bitxid.DID("did:bitxhub:appchain001:."),
			DocAddr: "/ipfs/1/...",
			DocHash: []byte{1},
			Status:  bitxid.Normal,
		},
	}
	itemBytes, err := bitxid.Struct2Bytes(item)

	content := pb.Content{
		SrcContractId: constant.MethodRegistryContractAddr.String(),
		DstContractId: constant.MethodRegistryContractAddr.String(),
		Func:          "Synchronize",
		Args:          [][]byte{[]byte("did:bitxhub:relayroot:."), itemBytes},
		Callback:      "",
	}
	contentBytes, err := content.Marshal()
	suite.Require().Nil(err)

	payload := pb.Payload{
		Encrypted: false,
		Content:   contentBytes,
	}
	payloadBytes, err := payload.Marshal()
	suite.Require().Nil(err)

	ib := &pb.IBTP{
		From:      Relay1,
		To:        Relay2,
		Payload:   payloadBytes,
		Index:     1,
		Timestamp: time.Now().UnixNano(),
		Proof:     proofHash[:],
	}

	tx, err := suite.client.GenerateIBTPTx(ib)
	suite.Require().Nil(err)
	tx.Extra = []byte(proof)

	res, err := suite.client.SendTransactionWithReceipt(tx, &rpcx.TransactOpts{
		From: fmt.Sprintf(	"%s-%s-%d", ib.From, ib.To, ib.Category()),
	})
	suite.Require().Nil(err)
	fmt.Println("res.Ret:", string(res.Ret))
}