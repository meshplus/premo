package bxh_tester

import (
	"crypto/sha256"
	"time"

	"github.com/gobuffalo/packr"
	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
	"github.com/meshplus/premo/internal/repo"
	"github.com/pkg/errors"
)

type Model13 struct {
	*Snake
}

func (suite *Model13) SetupTest() {
	suite.T().Parallel()
}

//tc：proof正确，跨链交易执行成功
func (suite Model13) Test1301_IBTPIsSuccess() {
	box := packr.NewBox(repo.ConfigPath)
	proof, err := box.Find("proof_1.0.0_rc_complex")
	suite.Require().Nil(err)
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, suite.GetChainID(pk), FabricRuleAddr)
	suite.Require().Nil(err)
	ibtp := suite.MockIBTP(1, suite.GetChainID(pk), suite.GetChainID(pk), proof)
	payload := suite.MockPayload(
		"mychannel&transfer",
		"mychannel&transfer",
		"interchainCharge",
		"",
		"interchainRollback",
		[][]byte{[]byte("Alice"), []byte("Alice"), []byte("1")},
		[][]byte{},
		[][]byte{[]byte("Alice"), []byte("1")},
	)
	err = suite.SendInterchainTx(pk, ibtp, payload, proof)
	suite.Require().Nil(err)
}

//tc：ibtp的index不一致，跨链交易执行失败
func (suite Model13) Test1302_IBTPWithWrongIndexIsFail() {
	box := packr.NewBox(repo.ConfigPath)
	proof, err := box.Find("proof_1.0.0_rc_complex")
	suite.Require().Nil(err)
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, suite.GetChainID(pk), FabricRuleAddr)
	suite.Require().Nil(err)
	ibtp := suite.MockIBTP(2, suite.GetChainID(pk), suite.GetChainID(pk), proof)
	payload := suite.MockPayload(
		"mychannel&transfer",
		"mychannel&transfer",
		"interchainCharge",
		"",
		"interchainRollback",
		[][]byte{[]byte("Alice"), []byte("Alice"), []byte("1")},
		[][]byte{},
		[][]byte{[]byte("Alice"), []byte("1")},
	)
	err = suite.SendInterchainTx(pk, ibtp, payload, proof)
	suite.Require().NotNil(err)
}

//tc：ibtp的目的链dst_contract_did不一致，跨链交易执行失败
func (suite Model13) Test1303_IBTPWithWrongDSTIsFail() {
	box := packr.NewBox(repo.ConfigPath)
	proof, err := box.Find("proof_1.0.0_rc_complex")
	suite.Require().Nil(err)
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, suite.GetChainID(pk), FabricRuleAddr)
	suite.Require().Nil(err)
	ibtp := suite.MockIBTP(1, suite.GetChainID(pk), suite.GetChainID(pk), proof)
	payload := suite.MockPayload(
		"mychannel&transfer",
		"mychannel&transfer123",
		"interchainCharge",
		"",
		"interchainRollback",
		[][]byte{[]byte("Alice"), []byte("Alice"), []byte("1")},
		[][]byte{},
		[][]byte{[]byte("Alice"), []byte("1")},
	)
	err = suite.SendInterchainTx(pk, ibtp, payload, proof)
	suite.Require().NotNil(err)
}

//tc：ibtp的来源链src_contract_id不一致，跨链交易执行失败
func (suite Model13) Test1304_IBTPWithWrongSRCIsFail() {
	box := packr.NewBox(repo.ConfigPath)
	proof, err := box.Find("proof_1.0.0_rc_complex")
	suite.Require().Nil(err)
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, suite.GetChainID(pk), FabricRuleAddr)
	suite.Require().Nil(err)
	ibtp := suite.MockIBTP(1, suite.GetChainID(pk), suite.GetChainID(pk), proof)
	payload := suite.MockPayload(
		"mychannel&transfer123",
		"mychannel&transfer",
		"interchainCharge",
		"",
		"interchainRollback",
		[][]byte{[]byte("Alice"), []byte("Alice"), []byte("1")},
		[][]byte{},
		[][]byte{[]byte("Alice"), []byte("1")},
	)
	err = suite.SendInterchainTx(pk, ibtp, payload, proof)
	suite.Require().NotNil(err)
}

//tc：ibtp的来源链func不一致，跨链交易执行失败
func (suite Model13) Test1305_IBTPWithWrongFuncIsFail() {
	box := packr.NewBox(repo.ConfigPath)
	proof, err := box.Find("proof_1.0.0_rc_complex")
	suite.Require().Nil(err)
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, suite.GetChainID(pk), FabricRuleAddr)
	suite.Require().Nil(err)
	ibtp := suite.MockIBTP(1, suite.GetChainID(pk), suite.GetChainID(pk), proof)
	payload := suite.MockPayload(
		"mychannel&transfer",
		"mychannel&transfer",
		"interchainCharge123",
		"",
		"interchainRollback",
		[][]byte{[]byte("Alice"), []byte("Alice"), []byte("1")},
		[][]byte{},
		[][]byte{[]byte("Alice"), []byte("1")},
	)
	err = suite.SendInterchainTx(pk, ibtp, payload, proof)
	suite.Require().NotNil(err)
}

//tc：ibtp的来源链args不一致，跨链交易执行失败
func (suite Model13) Test1306_IBTPWithWrongArgsIsFail() {
	box := packr.NewBox(repo.ConfigPath)
	proof, err := box.Find("proof_1.0.0_rc_complex")
	suite.Require().Nil(err)
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, suite.GetChainID(pk), FabricRuleAddr)
	suite.Require().Nil(err)
	ibtp := suite.MockIBTP(1, suite.GetChainID(pk), suite.GetChainID(pk), proof)
	payload := suite.MockPayload(
		"mychannel&transfer",
		"mychannel&transfer",
		"interchainCharge",
		"",
		"interchainRollback",
		[][]byte{[]byte("Alice"), []byte("Alice"), []byte("123")},
		[][]byte{},
		[][]byte{[]byte("Alice"), []byte("1")},
	)
	err = suite.SendInterchainTx(pk, ibtp, payload, proof)
	suite.Require().NotNil(err)
}

//tc：ibtp的来源链callback不一致，跨链交易执行失败
func (suite Model13) Test1307_IBTPWithWrongCallBackIsFail() {
	box := packr.NewBox(repo.ConfigPath)
	proof, err := box.Find("proof_1.0.0_rc_complex")
	suite.Require().Nil(err)
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, suite.GetChainID(pk), FabricRuleAddr)
	suite.Require().Nil(err)
	ibtp := suite.MockIBTP(1, suite.GetChainID(pk), suite.GetChainID(pk), proof)
	payload := suite.MockPayload(
		"mychannel&transfer",
		"mychannel&transfer",
		"interchainCharge",
		"123",
		"interchainRollback",
		[][]byte{[]byte("Alice"), []byte("Alice"), []byte("1")},
		[][]byte{},
		[][]byte{[]byte("Alice"), []byte("1")},
	)
	err = suite.SendInterchainTx(pk, ibtp, payload, proof)
	suite.Require().NotNil(err)
}

//tc：ibtp的来源链argscb不一致，跨链交易执行失败
func (suite Model13) Test1308_IBTPWithWrongArgsCbIsFail() {
	box := packr.NewBox(repo.ConfigPath)
	proof, err := box.Find("proof_1.0.0_rc_complex")
	suite.Require().Nil(err)
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, suite.GetChainID(pk), FabricRuleAddr)
	suite.Require().Nil(err)
	ibtp := suite.MockIBTP(1, suite.GetChainID(pk), suite.GetChainID(pk), proof)
	payload := suite.MockPayload(
		"mychannel&transfer",
		"mychannel&transfer",
		"interchainCharge",
		"",
		"interchainRollback",
		[][]byte{[]byte("Alice"), []byte("Alice"), []byte("1")},
		[][]byte{[]byte("123")},
		[][]byte{[]byte("Alice"), []byte("1")},
	)
	err = suite.SendInterchainTx(pk, ibtp, payload, proof)
	suite.Require().NotNil(err)
}

//tc：ibtp的来源链rollback不一致，跨链交易执行失败
func (suite Model13) Test1309_IBTPWithWrongRollBackIsFail() {
	box := packr.NewBox(repo.ConfigPath)
	proof, err := box.Find("proof_1.0.0_rc_complex")
	suite.Require().Nil(err)
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, suite.GetChainID(pk), FabricRuleAddr)
	suite.Require().Nil(err)
	ibtp := suite.MockIBTP(1, suite.GetChainID(pk), suite.GetChainID(pk), proof)
	payload := suite.MockPayload(
		"mychannel&transfer",
		"mychannel&transfer",
		"interchainCharge",
		"",
		"interchainRollback123",
		[][]byte{[]byte("Alice"), []byte("Alice"), []byte("1")},
		[][]byte{},
		[][]byte{[]byte("Alice"), []byte("1")},
	)
	err = suite.SendInterchainTx(pk, ibtp, payload, proof)
	suite.Require().NotNil(err)
}

//tc：ibtp的来源链argsrb不一致，跨链交易执行失败
func (suite Model13) Test1311_IBTPWithWrongArgsRbIsFail() {
	box := packr.NewBox(repo.ConfigPath)
	proof, err := box.Find("proof_1.0.0_rc_complex")
	suite.Require().Nil(err)
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, suite.GetChainID(pk), FabricRuleAddr)
	suite.Require().Nil(err)
	ibtp := suite.MockIBTP(1, suite.GetChainID(pk), suite.GetChainID(pk), proof)
	payload := suite.MockPayload(
		"mychannel&transfer",
		"mychannel&transfer",
		"interchainCharge",
		"",
		"interchainRollback123",
		[][]byte{[]byte("Alice"), []byte("Alice"), []byte("1")},
		[][]byte{},
		[][]byte{[]byte("Alice"), []byte("123")},
	)
	err = suite.SendInterchainTx(pk, ibtp, payload, proof)
	suite.Require().NotNil(err)
}

//tc：proof非法，跨链交易执行失败
func (suite Model13) Test1312_IBTPWithWrongProofIsFail() {
	box := packr.NewBox(repo.ConfigPath)
	proof, err := box.Find("proof_1.0.0_rc_complex_error")
	suite.Require().Nil(err)
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, suite.GetChainID(pk), FabricRuleAddr)
	suite.Require().Nil(err)
	ibtp := suite.MockIBTP(1, suite.GetChainID(pk), suite.GetChainID(pk), proof)
	payload := suite.MockPayload(
		"mychannel&transfer",
		"mychannel&transfer",
		"interchainCharge",
		"",
		"interchainRollback123",
		[][]byte{[]byte("Alice"), []byte("Alice"), []byte("1")},
		[][]byte{},
		[][]byte{[]byte("Alice"), []byte("1")},
	)
	err = suite.SendInterchainTx(pk, ibtp, payload, proof)
	suite.Require().NotNil(err)
}
func (suite Model13) MockIBTP(index uint64, from, to string, proof []byte) *pb.IBTP {
	proofHash := sha256.Sum256(proof)
	return &pb.IBTP{
		From:      from,
		To:        to,
		Index:     index,
		Type:      pb.IBTP_INTERCHAIN,
		Timestamp: time.Now().UnixNano(),
		Proof:     proofHash[:],
	}
}

func (suite Model13) MockPayload(srcContractId, dstContractId, funcName, callback, rollback string, args, argscb, argsrb [][]byte) []byte {
	content := &pb.Content{
		SrcContractId: srcContractId,
		DstContractId: dstContractId,
		Func:          funcName,
		Args:          args,
		Callback:      callback,
		ArgsCb:        argscb,
		Rollback:      rollback,
		ArgsRb:        argsrb,
	}
	bytes, _ := content.Marshal()
	payload := &pb.Payload{
		Encrypted: false,
		Content:   bytes,
	}
	ibtppd, _ := payload.Marshal()
	return ibtppd
}

func (suite Model13) SendInterchainTx(pk crypto.PrivateKey, ibtp *pb.IBTP, payload, proof []byte) error {
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
