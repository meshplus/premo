package bxh_tester

import (
	"crypto/sha256"
	"time"

	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
	"github.com/pkg/errors"
)

type Model15 struct {
	*Snake
}

func (suite *Model15) SetupTest() {
	suite.T().Parallel()
}

//tc：proof正确，跨链交易执行成功
//func (suite Model15) Test1301_IBTPIsSuccess() {
//	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
//	suite.Require().Nil(err)
//	from, err := pk.PublicKey().Address()
//	suite.Require().Nil(err)
//	err = suite.RegisterAppchainWithType(pk, "Fabric V1.4.3", HappyRuleAddr, "{\"channel_id\":\"mychannel\",\"chaincode_id\":\"broker\",\"broker_version\":\"1\"}")
//	err = suite.RegisterServer(pk, from.String(), "mychannel&transfer", "CallContract")
//	suite.Require().Nil(err)
//	for i := 0; i < 10000; i++ {
//		box := packr.NewBox(repo.ConfigPath)
//		proof, err := box.Find("proof_1.0.0_rc_complex")
//		ibtp := suite.MockIBTP(uint64(i+1), "1356:"+from.String()+":mychannel&transfer", "1356:"+from.String()+":mychannel&transfer", pb.IBTP_INTERCHAIN, proof)
//		payload := suite.MockContent(
//			"interchainCharge",
//			[][]byte{[]byte("Alice"), []byte("Alice"), []byte("1")},
//		)
//		err = suite.SendInterchainTx(pk, ibtp, payload, proof)
//		suite.Require().Nil(err)
//		time.Sleep(time.Second * 5)
//	}
//}

////tc：ibtp的index不一致，跨链交易执行失败
//func (suite Model15) Test1302_IBTPWithWrongIndexIsFail() {
//	box := packr.NewBox(repo.ConfigPath)
//	proof, err := box.Find("proof_1.0.0_rc_complex")
//	suite.Require().Nil(err)
//	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
//	suite.Require().Nil(err)
//	err = suite.RegisterAppchain(pk, suite.GetChainID(pk), FabricRuleAddr)
//	suite.Require().Nil(err)
//	ibtp := suite.MockIBTP(2, suite.GetChainID(pk), suite.GetChainID(pk), proof)
//	payload := suite.MockContent(
//		"interchainCharge",
//		[][]byte{[]byte("Alice"), []byte("Alice"), []byte("1")},
//	)
//	err = suite.SendInterchainTx(pk, ibtp, payload, proof)
//	suite.Require().NotNil(err)
//}
//
////tc：ibtp的目的链dst_contract_did不一致，跨链交易执行失败
//func (suite Model15) Test1303_IBTPWithWrongDSTIsFail() {
//	box := packr.NewBox(repo.ConfigPath)
//	proof, err := box.Find("proof_1.0.0_rc_complex")
//	suite.Require().Nil(err)
//	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
//	suite.Require().Nil(err)
//	err = suite.RegisterAppchain(pk, suite.GetChainID(pk), FabricRuleAddr)
//	suite.Require().Nil(err)
//	ibtp := suite.MockIBTP(1, suite.GetChainID(pk), suite.GetChainID(pk), proof)
//	payload := suite.MockContent(
//		"interchainCharge",
//		[][]byte{[]byte("Alice"), []byte("Alice"), []byte("1")},
//	)
//	err = suite.SendInterchainTx(pk, ibtp, payload, proof)
//	suite.Require().NotNil(err)
//}
//
////tc：ibtp的来源链src_contract_id不一致，跨链交易执行失败
//func (suite Model15) Test1304_IBTPWithWrongSRCIsFail() {
//	box := packr.NewBox(repo.ConfigPath)
//	proof, err := box.Find("proof_1.0.0_rc_complex")
//	suite.Require().Nil(err)
//	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
//	suite.Require().Nil(err)
//	err = suite.RegisterAppchain(pk, suite.GetChainID(pk), FabricRuleAddr)
//	suite.Require().Nil(err)
//	ibtp := suite.MockIBTP(1, suite.GetChainID(pk), suite.GetChainID(pk), proof)
//	payload := suite.MockContent(
//		"interchainCharge",
//		[][]byte{[]byte("Alice"), []byte("Alice"), []byte("1")},
//	)
//	err = suite.SendInterchainTx(pk, ibtp, payload, proof)
//	suite.Require().NotNil(err)
//}
//
////tc：ibtp的来源链func不一致，跨链交易执行失败
//func (suite Model15) Test1305_IBTPWithWrongFuncIsFail() {
//	box := packr.NewBox(repo.ConfigPath)
//	proof, err := box.Find("proof_1.0.0_rc_complex")
//	suite.Require().Nil(err)
//	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
//	suite.Require().Nil(err)
//	err = suite.RegisterAppchain(pk, suite.GetChainID(pk), FabricRuleAddr)
//	suite.Require().Nil(err)
//	ibtp := suite.MockIBTP(1, suite.GetChainID(pk), suite.GetChainID(pk), proof)
//	payload := suite.MockContent(
//		"interchainCharge",
//		[][]byte{[]byte("Alice"), []byte("Alice"), []byte("1")},
//	)
//	err = suite.SendInterchainTx(pk, ibtp, payload, proof)
//	suite.Require().NotNil(err)
//}
//
////tc：ibtp的来源链args不一致，跨链交易执行失败
//func (suite Model15) Test1306_IBTPWithWrongArgsIsFail() {
//	box := packr.NewBox(repo.ConfigPath)
//	proof, err := box.Find("proof_1.0.0_rc_complex")
//	suite.Require().Nil(err)
//	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
//	suite.Require().Nil(err)
//	err = suite.RegisterAppchain(pk, suite.GetChainID(pk), FabricRuleAddr)
//	suite.Require().Nil(err)
//	ibtp := suite.MockIBTP(1, suite.GetChainID(pk), suite.GetChainID(pk), proof)
//	payload := suite.MockContent(
//		"interchainCharge",
//		[][]byte{[]byte("Alice"), []byte("Alice"), []byte("1")},
//	)
//	err = suite.SendInterchainTx(pk, ibtp, payload, proof)
//	suite.Require().NotNil(err)
//}
//
////tc：ibtp的来源链callback不一致，跨链交易执行失败
//func (suite Model15) Test1307_IBTPWithWrongCallBackIsFail() {
//	box := packr.NewBox(repo.ConfigPath)
//	proof, err := box.Find("proof_1.0.0_rc_complex")
//	suite.Require().Nil(err)
//	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
//	suite.Require().Nil(err)
//	err = suite.RegisterAppchain(pk, suite.GetChainID(pk), FabricRuleAddr)
//	suite.Require().Nil(err)
//	ibtp := suite.MockIBTP(1, suite.GetChainID(pk), suite.GetChainID(pk), proof)
//	payload := suite.MockContent(
//		"interchainCharge",
//		[][]byte{[]byte("Alice"), []byte("Alice"), []byte("1")},
//	)
//	err = suite.SendInterchainTx(pk, ibtp, payload, proof)
//	suite.Require().NotNil(err)
//}
//
////tc：ibtp的来源链argscb不一致，跨链交易执行失败
//func (suite Model15) Test1308_IBTPWithWrongArgsCbIsFail() {
//	box := packr.NewBox(repo.ConfigPath)
//	proof, err := box.Find("proof_1.0.0_rc_complex")
//	suite.Require().Nil(err)
//	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
//	suite.Require().Nil(err)
//	err = suite.RegisterAppchain(pk, suite.GetChainID(pk), FabricRuleAddr)
//	suite.Require().Nil(err)
//	ibtp := suite.MockIBTP(1, suite.GetChainID(pk), suite.GetChainID(pk), proof)
//	payload := suite.MockContent(
//		"interchainCharge",
//		[][]byte{[]byte("Alice"), []byte("Alice"), []byte("1")},
//	)
//	err = suite.SendInterchainTx(pk, ibtp, payload, proof)
//	suite.Require().NotNil(err)
//}
//
////tc：ibtp的来源链rollback不一致，跨链交易执行失败
//func (suite Model15) Test1309_IBTPWithWrongRollBackIsFail() {
//	box := packr.NewBox(repo.ConfigPath)
//	proof, err := box.Find("proof_1.0.0_rc_complex")
//	suite.Require().Nil(err)
//	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
//	suite.Require().Nil(err)
//	err = suite.RegisterAppchain(pk, suite.GetChainID(pk), FabricRuleAddr)
//	suite.Require().Nil(err)
//	ibtp := suite.MockIBTP(1, suite.GetChainID(pk), suite.GetChainID(pk), proof)
//	payload := suite.MockContent(
//		"interchainCharge",
//		[][]byte{[]byte("Alice"), []byte("Alice"), []byte("1")},
//	)
//	err = suite.SendInterchainTx(pk, ibtp, payload, proof)
//	suite.Require().NotNil(err)
//}
//
////tc：ibtp的来源链argsrb不一致，跨链交易执行失败
//func (suite Model15) Test1311_IBTPWithWrongArgsRbIsFail() {
//	box := packr.NewBox(repo.ConfigPath)
//	proof, err := box.Find("proof_1.0.0_rc_complex")
//	suite.Require().Nil(err)
//	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
//	suite.Require().Nil(err)
//	err = suite.RegisterAppchain(pk, suite.GetChainID(pk), FabricRuleAddr)
//	suite.Require().Nil(err)
//	ibtp := suite.MockIBTP(1, suite.GetChainID(pk), suite.GetChainID(pk), proof)
//	payload := suite.MockContent(
//		"interchainCharge",
//		[][]byte{[]byte("Alice"), []byte("Alice"), []byte("1")},
//	)
//	err = suite.SendInterchainTx(pk, ibtp, payload, proof)
//	suite.Require().NotNil(err)
//}
//
////tc：proof非法，跨链交易执行失败
//func (suite Model15) Test1312_IBTPWithWrongProofIsFail() {
//	box := packr.NewBox(repo.ConfigPath)
//	proof, err := box.Find("proof_1.0.0_rc_complex_error")
//	suite.Require().Nil(err)
//	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
//	suite.Require().Nil(err)
//	err = suite.RegisterAppchain(pk, suite.GetChainID(pk), FabricRuleAddr)
//	suite.Require().Nil(err)
//	ibtp := suite.MockIBTP(1, suite.GetChainID(pk), suite.GetChainID(pk), proof)
//	payload := suite.MockContent(
//		"interchainCharge",
//		[][]byte{[]byte("Alice"), []byte("Alice"), []byte("1")},
//	)
//	err = suite.SendInterchainTx(pk, ibtp, payload, proof)
//	suite.Require().NotNil(err)
//}
func (suite Snake) MockIBTP(index uint64, from, to string, typ pb.IBTP_Type, proof []byte) *pb.IBTP {
	proofHash := sha256.Sum256(proof)
	return &pb.IBTP{
		From:          from,
		To:            to,
		Index:         index,
		Type:          typ,
		TimeoutHeight: 10,
		Proof:         proofHash[:],
	}
}

func (suite Snake) MockContent(funcName string, args [][]byte) []byte {
	content := &pb.Content{
		Func: funcName,
		Args: args,
	}
	bytes, _ := content.Marshal()
	payload := &pb.Payload{
		Encrypted: false,
		Content:   bytes,
	}
	ibtppd, _ := payload.Marshal()
	return ibtppd
}

func (suite Snake) MockResult(data [][]byte) []byte {
	result := &pb.Result{Data: data}
	bytes, _ := result.Marshal()
	payload := &pb.Payload{
		Encrypted: false,
		Content:   bytes,
	}
	ibtppd, _ := payload.Marshal()
	return ibtppd
}

func (suite Snake) SendInterchainTx(pk crypto.PrivateKey, ibtp *pb.IBTP, payload, proof []byte) error {
	ibtp.Payload = payload
	from, err := pk.PublicKey().Address()
	if err != nil {
		return err
	}
	tx := &pb.BxhTransaction{
		From:      from,
		To:        constant.InterchainContractAddr.Address(),
		Timestamp: time.Now().UnixNano(),
		Extra:     proof,
		IBTP:      ibtp,
	}
	client := suite.NewClient(pk)
	res, err := client.SendTransactionWithReceipt(tx, nil)
	if err != nil {
		return err
	}
	if res.Status == pb.Receipt_FAILED {
		return errors.New(string(res.Ret))
	}
	return nil
}
