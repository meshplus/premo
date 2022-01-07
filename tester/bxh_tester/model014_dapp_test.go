package bxh_tester

import (
	"encoding/json"
	"fmt"
	"sync/atomic"

	"github.com/looplab/fsm"
	"github.com/meshplus/bitxhub-core/governance"
	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/meshplus/premo/internal/repo"
	"github.com/pkg/errors"
)

type TransferRecord struct {
	From       string `json:"from"`
	To         string `json:"to"`
	Reason     string `json:"reason"`
	Confirm    bool   `json:"confirm"`
	CreateTime int64  `json:"create_time"`
}
type Dapp struct {
	DappID            string                                  `json:"dapp_id"` // first owner address + num
	Name              string                                  `json:"name"`
	Type              string                                  `json:"type"`
	Desc              string                                  `json:"desc"`
	ContractAddr      map[string]struct{}                     `json:"contract_addr"`
	Permission        map[string]struct{}                     `json:"permission"` // users which are not allowed to see the dapp
	OwnerAddr         string                                  `json:"owner_addr"`
	CreateTime        int64                                   `json:"create_time"`
	Score             float64                                 `json:"score"`
	EvaluationRecords map[string]*governance.EvaluationRecord `json:"evaluation_records"`
	TransferRecords   []*TransferRecord                       `json:"transfer_records"`
	Status            governance.GovernanceStatus             `json:"status"`
	FSM               *fsm.FSM                                `json:"fsm"`
}

type Model14 struct {
	*Snake
}

func (suite *Model14) SetupTest() {
	suite.T().Parallel()
}

//tc：根据存在的合约地址注册dapp，dapp注册成功
func (suite Model14) Test1401_RegisterDappIsSuccess() {
	address := suite.deployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.RegisterDapp(pk, address.String())
	suite.Require().Nil(err)
}

//tc：根据不存在的合约地址注册dapp，dapp注册失败
func (suite Model14) Test1402_RegisterDappWithNoExistAddrIsFail() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.RegisterDapp(pk, "0x0000000000000000000000000000000000000000")
	suite.Require().NotNil(err)
}

//tc：dapp使用已经绑定dapp的合约地址注册dapp，dapp注册失败
func (suite Model14) Test1403_RegisterDappWithUsedAddrIsFail() {
	address := suite.deployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.RegisterDapp(pk, address.String())
	suite.Require().Nil(err)
	err = suite.RegisterDapp(pk, address.String())
	suite.Require().NotNil(err)
}

//tc：根据存在的合约地址更新dapp，dapp更新成功
func (suite Model14) Test1404_UpdateDappIsSuccess() {
	address1 := suite.deployLedgerContract()
	address2 := suite.deployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.RegisterDapp(pk, address1.String())
	suite.Require().Nil(err)
	err = suite.UpdateDapp(pk, address2.String())
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.DappMgrContractAddr.Address(), "GetDapp", nil, rpcx.String(suite.GetDappID(pk)))
	suite.Require().Nil(err)
	suite.Require().Contains(string(res.Ret), address2.String())
}

//tc：根据不存在的合约地址更新dapp，dapp更新失败
func (suite Model14) Test1405_UpdateDappWithNoExistAddrIsFail() {
	address := suite.deployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.RegisterDapp(pk, address.String())
	suite.Require().Nil(err)
	err = suite.UpdateDapp(pk, "0x0000000000000000000000000000000000000000")
	suite.Require().NotNil(err)
}

//tc:根据已绑定的合约地址更新dapp，dapp更新失败
func (suite Model14) Test1406_UpdateDappWithUsedAddrIsFail() {
	address1 := suite.deployLedgerContract()
	pk1, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.RegisterDapp(pk1, address1.String())
	suite.Require().Nil(err)
	address2 := suite.deployLedgerContract()
	pk2, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.RegisterDapp(pk2, address2.String())
	suite.Require().Nil(err)
	err = suite.UpdateDapp(pk2, address1.String())
	suite.Require().NotNil(err)
	client := suite.NewClient(pk1)
	res, err := client.InvokeBVMContract(constant.DappMgrContractAddr.Address(), "GetDapp", nil, rpcx.String(suite.GetDappID(pk2)))
	suite.Require().Nil(err)
	suite.Require().Contains(string(res.Ret), address2.String())
}

//tc：dapp处于unavailable状态更新dapp，dapp更新失败
func (suite Model14) Test1407_UpdateDappWithUnavailableDappIsFail() {
	address := suite.deployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.DappToUnavailable(pk, address.String())
	suite.Require().Nil(err)
	err = suite.UpdateDapp(pk, address.String())
	suite.Require().NotNil(err)
}

//tc：dapp处于activating状态更新dapp，dapp更新失败
func (suite Model14) Test1408_UpdateDappWithActivatingDappIsFail() {
	address := suite.deployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.DappToActivating(pk, address.String())
	suite.Require().Nil(err)
	err = suite.UpdateDapp(pk, address.String())
	suite.Require().NotNil(err)
}

//tc：dapp处于updating状态更新dapp，dapp更新失败
func (suite Model14) Test1409_UpdateDappWithUpdatingDappIsFail() {
	address := suite.deployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.DappToUpdating(pk, address.String())
	suite.Require().Nil(err)
	err = suite.UpdateDapp(pk, address.String())
	suite.Require().NotNil(err)
}

//tc：dapp处于freezing状态更新dapp，dapp更新失败
func (suite Model14) Test1410_UpdateDappWithFreezingDappIsFail() {
	address := suite.deployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.DappToFreezing(pk, address.String())
	suite.Require().Nil(err)
	err = suite.UpdateDapp(pk, address.String())
	suite.Require().NotNil(err)
}

//tc：dapp处于frozen状态更新dapp，dapp更新失败
func (suite Model14) Test14011_UpdateDappWithFrozenDappIsFail() {
	address := suite.deployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.DappToFrozen(pk, address.String())
	suite.Require().Nil(err)
	err = suite.UpdateDapp(pk, address.String())
	suite.Require().Nil(err)
}

//tc：根据存在的合约地址冻结dapp，dapp冻结成功
func (suite Model14) Test1412_FreezeDappIsSuccess() {
	address := suite.deployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.RegisterDapp(pk, address.String())
	suite.Require().Nil(err)
	err = suite.FreezeDapp(pk)
	suite.Require().Nil(err)
}

//tc：根据不存在的合约地址冻结dapp，dapp冻结失败
func (suite Model14) Test1413_FreezeDappWithNoExistAddrIsFail() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.FreezeDapp(pk)
	suite.Require().NotNil(err)
}

//tc：dapp处于Unavailable状态冻结dapp，dapp冻结成功
func (suite Model14) Test1414_FreezeDappWithUnavailableDappIsFail() {
	address := suite.deployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.DappToUnavailable(pk, address.String())
	suite.Require().Nil(err)
	err = suite.FreezeDapp(pk)
	suite.Require().NotNil(err)
}

//tc：dapp处于activating状态冻结dapp，dapp冻结失败
func (suite Model14) Test1415_FreezeDappWithActivatingDappIsFail() {
	address := suite.deployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.DappToActivating(pk, address.String())
	suite.Require().Nil(err)
	err = suite.FreezeDapp(pk)
	suite.Require().NotNil(err)
}

//tc：dapp处于updating状态冻结dapp，dapp冻结失败
func (suite Model14) Test1416_FreezeDappWithUpdatingDappIsFail() {
	address := suite.deployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.DappToUpdating(pk, address.String())
	suite.Require().Nil(err)
	err = suite.FreezeDapp(pk)
	suite.Require().NotNil(err)
}

//tc：dapp处于freezing状态冻结dapp，dapp冻结失败
func (suite Model14) Test1417_FreezeDappWithFreezingDappIsFail() {
	address := suite.deployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.DappToFreezing(pk, address.String())
	suite.Require().Nil(err)
	err = suite.FreezeDapp(pk)
	suite.Require().NotNil(err)
}

//tc：dapp处于frozen状态冻结dapp，dapp冻结失败
func (suite Model14) Test1418_FreezeDappWithFrozenDappIsFail() {
	address := suite.deployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.DappToFrozen(pk, address.String())
	suite.Require().Nil(err)
	err = suite.FreezeDapp(pk)
	suite.Require().NotNil(err)
}

//tc：根据存在的合约地址激活dapp，dapp激活成功
func (suite Model14) Test1419_ActivateDappIsSuccess() {
	address := suite.deployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.RegisterDapp(pk, address.String())
	suite.Require().Nil(err)
	err = suite.FreezeDapp(pk)
	suite.Require().Nil(err)
	err = suite.CheckDappStatus(suite.GetDappID(pk), governance.GovernanceFrozen)
	suite.Require().Nil(err)
	err = suite.ActivateDapp(pk)
	suite.Require().Nil(err)
}

//tc：根据不存在的合约地址激活dapp，dapp激活成功
func (suite Model14) Test1420_ActivateDappWithNoExistAddrIsFail() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.ActivateDapp(pk)
	suite.Require().NotNil(err)
}

//tc：dapp处于available状态激活dapp，dapp激活失败
func (suite Model14) Test1421_ActivateDappWithAvailableDappIsFail() {
	address := suite.deployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.DappToAvailable(pk, address.String())
	suite.Require().Nil(err)
	err = suite.ActivateDapp(pk)
	suite.Require().NotNil(err)
}

//tc：dapp处于unavailable状态激活dapp，dapp激活失败
func (suite Model14) Test1422_ActivateDappWithUnavailableDappIsFail() {
	address := suite.deployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.DappToUnavailable(pk, address.String())
	suite.Require().Nil(err)
	err = suite.ActivateDapp(pk)
	suite.Require().NotNil(err)
}

//tc：dapp处于activating状态激活dapp，dapp激活失败
func (suite Model14) Test1423_ActivateDappWithActivatingDappIsFail() {
	address := suite.deployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.DappToActivating(pk, address.String())
	suite.Require().Nil(err)
	err = suite.ActivateDapp(pk)
	suite.Require().NotNil(err)
}

//tc：dapp处于updating状态激活dapp，dapp激活失败
func (suite Model14) Test1424_ActivateDappWithUpdatingDappIsFail() {
	address := suite.deployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.DappToUpdating(pk, address.String())
	suite.Require().Nil(err)
	err = suite.ActivateDapp(pk)
	suite.Require().NotNil(err)
}

//tc：dapp处于freezing状态激活dapp，dapp激活失败
func (suite Model14) Test1425_ActivateDappWithFreezingDappIsFail() {
	address := suite.deployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.DappToFreezing(pk, address.String())
	suite.Require().Nil(err)
	err = suite.ActivateDapp(pk)
	suite.Require().NotNil(err)
}

//tc：根据存在的合约地址转让dapp，dapp转让成功
func (suite Model14) Test1426_TransferDappIsSuccess() {
	address := suite.deployLedgerContract()
	pk1, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.RegisterDapp(pk1, address.String())
	suite.Require().Nil(err)
	pk2, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.TransferDapp(pk1, pk2)
	suite.Require().Nil(err)
	err = suite.TransferDapp(pk1, pk2)
	suite.Require().NotNil(err)
}

//tc：根据不存在的合约地址转让dapp，dapp转让失败
func (suite Model14) Test1427_TransferDappWithNoExistDappIsFail() {
	pk1, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	pk2, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.TransferDapp(pk1, pk2)
	suite.Require().NotNil(err)
}

//tc：dapp处于unavailable状态转让dapp，dapp转让失败
func (suite Model14) Test1428_TransferDappWithUnavailableDappIsFail() {
	address := suite.deployLedgerContract()
	pk1, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.DappToUnavailable(pk1, address.String())
	suite.Require().Nil(err)
	pk2, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.TransferDapp(pk1, pk2)
	suite.Require().NotNil(err)
}

//tc：dapp处于activating状态转让dapp，dapp转让失败
func (suite Model14) Test1429_TransferDappWithActivatingDappIsFail() {
	address := suite.deployLedgerContract()
	pk1, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.DappToActivating(pk1, address.String())
	suite.Require().Nil(err)
	pk2, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.TransferDapp(pk1, pk2)
	suite.Require().NotNil(err)
}

//tc：dapp处于updating状态转让dapp，dapp转让失败
func (suite Model14) Test1430_TransferDappWithUpdatingDappIsFail() {
	address := suite.deployLedgerContract()
	pk1, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.DappToUpdating(pk1, address.String())
	suite.Require().Nil(err)
	pk2, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.TransferDapp(pk1, pk2)
	suite.Require().NotNil(err)
}

//tc：dapp处于freezing状态转让dapp，dapp转让失败
func (suite Model14) Test1431_TransferDappWithFreezingDappIsFail() {
	address := suite.deployLedgerContract()
	pk1, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.DappToFreezing(pk1, address.String())
	suite.Require().Nil(err)
	pk2, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.TransferDapp(pk1, pk2)
	suite.Require().NotNil(err)
}

//tc：dapp处于frozen状态转让dapp，dapp转让失败
func (suite Model14) Test1432_TransferDappWithFrozenDappIsFail() {
	address := suite.deployLedgerContract()
	pk1, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.DappToFrozen(pk1, address.String())
	suite.Require().Nil(err)
	pk2, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.TransferDapp(pk1, pk2)
	suite.Require().NotNil(err)
}

//tc：dapp转让给自己，dapp转让失败
func (suite Model14) Test1433_TransferDappWithSelfIsFail() {
	address := suite.deployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.RegisterDapp(pk, address.String())
	suite.Require().Nil(err)
	err = suite.TransferDapp(pk, pk)
	suite.Require().NotNil(err)
}

//tc：dapp接收确认不存在的dapp转移，确认失败
func (suite Model14) Test1433_ConfirmTransferWithNoExistTransferIsFail() {
	address := suite.deployLedgerContract()
	pk1, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.RegisterDapp(pk1, address.String())
	suite.Require().Nil(err)
	pk2, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.ConfirmTransfer(pk1, pk2)
	suite.Require().NotNil(err)
}

//tc：评价存在dapp，dapp评价成功
func (suite Model14) Test1434_EvaluateDappIsSuccess() {
	address := suite.deployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.RegisterDapp(pk, address.String())
	suite.Require().Nil(err)
	err = suite.EvaluateDapp(pk, suite.GetDappID(pk), "good", 5.0)
	suite.Require().Nil(err)
}

//tc：评价不存在dapp，dapp评价失败
func (suite Model14) Test1435_EvaluateDappWithNoExistDappIsFail() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.EvaluateDapp(pk, suite.GetDappID(pk), "good", 5.0)
	suite.Require().NotNil(err)
}

//tc：评价评分不在[0-5]，dapp评价失败
func (suite Model14) Test1436_EvaluateDappWithErrorScoreIsFail() {
	address := suite.deployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.RegisterDapp(pk, address.String())
	suite.Require().Nil(err)
	err = suite.EvaluateDapp(pk, suite.GetDappID(pk), "good", 50.0)
	suite.Require().NotNil(err)
}

//tc：重复评价dapp，dapp评价失败
func (suite Model14) Test1437_EvaluateDappWithRepeatEvaluateIsFail() {
	address := suite.deployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.RegisterDapp(pk, address.String())
	suite.Require().Nil(err)
	err = suite.EvaluateDapp(pk, suite.GetDappID(pk), "good", 5.0)
	suite.Require().Nil(err)
	err = suite.EvaluateDapp(pk, suite.GetDappID(pk), "good", 5.0)
	suite.Require().NotNil(err)
}

func (suite Snake) RegisterDapp(pk crypto.PrivateKey, conAddrs string) error {
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(conAddrs),
		rpcx.String("application"),
		rpcx.String("test"),
		rpcx.String("https://github.com/meshplus/bitxhub"),
		rpcx.String(conAddrs),
		rpcx.String(""),
		rpcx.String("reason"),
	}
	res, err := client.InvokeBVMContract(constant.DappMgrContractAddr.Address(), "RegisterDapp", nil, args...)
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return errors.New(string(res.Ret))
	}
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	if err != nil {
		return err
	}
	err = suite.VotePass(result.ProposalID)
	if err != nil {
		return err
	}
	return nil
}

func (suite Model14) UpdateDapp(pk crypto.PrivateKey, conAddrs string) error {
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(suite.GetDappID(pk)),
		rpcx.String(conAddrs + "123"),
		rpcx.String("test"),
		rpcx.String("https://github.com/meshplus/bitxhub"),
		rpcx.String(conAddrs),
		rpcx.String(""),
		rpcx.String("reason"),
	}
	res, err := client.InvokeBVMContract(constant.DappMgrContractAddr.Address(), "UpdateDapp", nil, args...)
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return errors.New(string(res.Ret))
	}
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	if err != nil {
		return err
	}
	err = suite.VotePass(result.ProposalID)
	if err != nil {
		return err
	}
	return nil
}

func (suite Model14) FreezeDapp(pk crypto.PrivateKey) error {
	path, err := repo.Node1Path()
	suite.Require().Nil(err)
	node1pk, err := asym.RestorePrivateKey(path, repo.KeyPassword)
	suite.Require().Nil(err)
	from, err := node1pk.PublicKey().Address()
	suite.Require().Nil(err)
	client := suite.NewClient(node1pk)
	res, err := client.InvokeBVMContract(constant.DappMgrContractAddr.Address(), "FreezeDapp", &rpcx.TransactOpts{
		From:  from.String(),
		Nonce: atomic.AddUint64(&nonce1, 1),
	}, rpcx.String(suite.GetDappID(pk)), rpcx.String("reason"))
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return errors.New(string(res.Ret))
	}
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	if err != nil {
		return err
	}
	err = suite.VotePass(result.ProposalID)
	if err != nil {
		return err
	}
	return nil
}

func (suite Model14) ActivateDapp(pk crypto.PrivateKey) error {
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.DappMgrContractAddr.Address(), "ActivateDapp", nil, rpcx.String(suite.GetDappID(pk)), rpcx.String("reason"))
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return errors.New(string(res.Ret))
	}
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	if err != nil {
		return err
	}
	err = suite.VotePass(result.ProposalID)
	if err != nil {
		return err
	}
	return nil
}

func (suite Model14) TransferDapp(pk1, pk2 crypto.PrivateKey) error {
	client := suite.NewClient(pk1)
	address, err := pk2.PublicKey().Address()
	if err != nil {
		return err
	}
	res, err := client.InvokeBVMContract(constant.DappMgrContractAddr.Address(), "TransferDapp", nil,
		rpcx.String(suite.GetDappID(pk1)), rpcx.String(address.String()), rpcx.String("reason"))
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return errors.New(string(res.Ret))
	}
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	if err != nil {
		return err
	}
	err = suite.VotePass(result.ProposalID)
	if err != nil {
		return err
	}
	err = suite.ConfirmTransfer(pk1, pk2)
	if err != nil {
		return err
	}
	return nil
}

func (suite Model14) ConfirmTransfer(pk1, pk2 crypto.PrivateKey) error {
	client := suite.NewClient(pk2)
	res, err := client.InvokeBVMContract(constant.DappMgrContractAddr.Address(), "ConfirmTransfer", nil, rpcx.String(suite.GetDappID(pk1)))
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return errors.New(string(res.Ret))
	}
	return nil
}

func (suite Model14) EvaluateDapp(pk crypto.PrivateKey, id, desc string, score float64) error {
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.DappMgrContractAddr.Address(), "EvaluateDapp", nil,
		rpcx.String(id), rpcx.String(desc), rpcx.Float64(score))
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return errors.New(string(res.Ret))
	}
	return nil
}

func (suite Snake) CheckDappStatus(id string, status governance.GovernanceStatus) error {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.DappMgrContractAddr.Address(), "GetDapp", nil, rpcx.String(id))
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return errors.New(string(res.Ret))
	}
	dapp := &Dapp{}
	err = json.Unmarshal(res.Ret, dapp)
	if err != nil {
		return err
	}
	if dapp.Status != status {
		return errors.New(fmt.Sprintf("expect status is %s, but got %s", status, dapp.Status))
	}
	return nil
}

func (suite Snake) GetDappID(pk crypto.PrivateKey) string {
	address, _ := pk.PublicKey().Address()
	return address.String() + "-0"
}

func (suite Model14) DappToAvailable(pk crypto.PrivateKey, address string) error {
	err := suite.RegisterDapp(pk, address)
	if err != nil {
		return err
	}
	err = suite.CheckDappStatus(suite.GetDappID(pk), governance.GovernanceAvailable)
	if err != nil {
		return err
	}
	return nil
}

func (suite Model14) DappToUnavailable(pk crypto.PrivateKey, address string) error {
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(address),
		rpcx.String("application"),
		rpcx.String("test"),
		rpcx.String("https://github.com/meshplus/bitxhub"),
		rpcx.String(address),
		rpcx.String(""),
		rpcx.String("reason"),
	}
	res, err := client.InvokeBVMContract(constant.DappMgrContractAddr.Address(), "RegisterDapp", nil, args...)
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return errors.New(string(res.Ret))
	}
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	if err != nil {
		return err
	}
	err = suite.VoteReject(result.ProposalID)
	if err != nil {
		return err
	}
	err = suite.CheckDappStatus(suite.GetDappID(pk), governance.GovernanceUnavailable)
	if err != nil {
		return err
	}
	return nil
}

func (suite Model14) DappToActivating(pk crypto.PrivateKey, address string) error {
	err := suite.RegisterDapp(pk, address)
	if err != nil {
		return err
	}
	err = suite.FreezeDapp(pk)
	if err != nil {
		return err
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.DappMgrContractAddr.Address(), "ActivateDapp", nil, rpcx.String(suite.GetDappID(pk)), rpcx.String("reason"))
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return errors.New(string(res.Ret))
	}
	err = suite.CheckDappStatus(suite.GetDappID(pk), governance.GovernanceActivating)
	if err != nil {
		return err
	}
	return nil
}

func (suite Model14) DappToUpdating(pk crypto.PrivateKey, address string) error {
	err := suite.RegisterDapp(pk, address)
	if err != nil {
		return err
	}
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(suite.GetDappID(pk)),
		rpcx.String(address + "123"),
		rpcx.String("test"),
		rpcx.String("https://github.com/meshplus/bitxhub"),
		rpcx.String(address),
		rpcx.String(""),
		rpcx.String("reason"),
	}
	res, err := client.InvokeBVMContract(constant.DappMgrContractAddr.Address(), "UpdateDapp", nil, args...)
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return errors.New(string(res.Ret))
	}
	err = suite.CheckDappStatus(suite.GetDappID(pk), governance.GovernanceUpdating)
	if err != nil {
		return err
	}
	return nil
}

func (suite Model14) DappToFreezing(pk crypto.PrivateKey, address string) error {
	err := suite.RegisterDapp(pk, address)
	if err != nil {
		return err
	}
	path, err := repo.Node1Path()
	suite.Require().Nil(err)
	node1pk, err := asym.RestorePrivateKey(path, repo.KeyPassword)
	suite.Require().Nil(err)
	from, err := node1pk.PublicKey().Address()
	suite.Require().Nil(err)
	client := suite.NewClient(node1pk)
	res, err := client.InvokeBVMContract(constant.DappMgrContractAddr.Address(), "FreezeDapp", &rpcx.TransactOpts{
		From:  from.String(),
		Nonce: atomic.AddUint64(&nonce1, 1),
	}, rpcx.String(suite.GetDappID(pk)), rpcx.String("reason"))
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return errors.New(string(res.Ret))
	}
	err = suite.CheckDappStatus(suite.GetDappID(pk), governance.GovernanceFreezing)
	if err != nil {
		return err
	}
	return nil
}

func (suite Model14) DappToFrozen(pk crypto.PrivateKey, address string) error {
	err := suite.RegisterDapp(pk, address)
	if err != nil {
		return err
	}
	err = suite.FreezeDapp(pk)
	if err != nil {
		return err
	}
	err = suite.CheckDappStatus(suite.GetDappID(pk), governance.GovernanceFrozen)
	if err != nil {
		return err
	}
	return nil
}
