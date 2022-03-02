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
	address := suite.DeployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.RegisterDapp(pk, address.String())
	suite.Require().Nil(err)
}

//tc：根据存在的合约地址更新dapp，dapp更新成功
func (suite Model14) Test1402_UpdateDappIsSuccess() {
	address1 := suite.DeployLedgerContract()
	address2 := suite.DeployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.RegisterDapp(pk, address1.String())
	suite.Require().Nil(err)
	err = suite.UpdateDapp(pk, address2.String())
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.DappMgrContractAddr.Address(), "GetDapp", nil, rpcx.String(suite.MockDappID(pk)))
	suite.Require().Nil(err)
	suite.Require().Contains(string(res.Ret), address2.String())
}

//tc：根据存在的合约地址冻结dapp，dapp冻结成功
func (suite Model14) Test1403_FreezeDappIsSuccess() {
	address := suite.DeployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.RegisterDapp(pk, address.String())
	suite.Require().Nil(err)
	err = suite.FreezeDapp(pk)
	suite.Require().Nil(err)
}

//tc：根据存在的合约地址激活dapp，dapp激活成功
func (suite Model14) Test1404_ActivateDappIsSuccess() {
	address := suite.DeployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.RegisterDapp(pk, address.String())
	suite.Require().Nil(err)
	err = suite.FreezeDapp(pk)
	suite.Require().Nil(err)
	err = suite.CheckDappStatus(suite.MockDappID(pk), governance.GovernanceFrozen)
	suite.Require().Nil(err)
	err = suite.ActivateDapp(pk)
	suite.Require().Nil(err)
}

//tc：根据存在的合约地址转让dapp，dapp转让成功
func (suite Model14) Test1405_TransferDappIsSuccess() {
	address := suite.DeployLedgerContract()
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

//tc：评价存在dapp，dapp评价成功
func (suite Model14) Test1406_EvaluateDappIsSuccess() {
	address := suite.DeployLedgerContract()
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.RegisterDapp(pk, address.String())
	suite.Require().Nil(err)
	err = suite.EvaluateDapp(pk, suite.MockDappID(pk), "good", 5.0)
	suite.Require().Nil(err)
}

// RegisterDapp register dapp
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
	suite.Require().Nil(err)
	if res.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf(string(res.Ret))
	}
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	suite.Require().Nil(err)
	err = suite.VotePass(result.ProposalID)
	suite.Require().Nil(err)
	return nil
}

// UpdateDapp update dapp
func (suite Model14) UpdateDapp(pk crypto.PrivateKey, conAddrs string) error {
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(suite.MockDappID(pk)),
		rpcx.String(conAddrs + "123"),
		rpcx.String("test"),
		rpcx.String("https://github.com/meshplus/bitxhub"),
		rpcx.String(conAddrs),
		rpcx.String(""),
		rpcx.String("reason"),
	}
	res, err := client.InvokeBVMContract(constant.DappMgrContractAddr.Address(), "UpdateDapp", nil, args...)
	suite.Require().Nil(err)
	if res.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf(string(res.Ret))
	}
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	suite.Require().Nil(err)
	err = suite.VotePass(result.ProposalID)
	suite.Require().Nil(err)
	return nil
}

// FreezeDapp freeze dapp
func (suite Model14) FreezeDapp(pk crypto.PrivateKey) error {
	node1pk, from, err := repo.Node1Priv()
	suite.Require().Nil(err)
	client := suite.NewClient(node1pk)
	res, err := client.InvokeBVMContract(constant.DappMgrContractAddr.Address(), "FreezeDapp", &rpcx.TransactOpts{
		From:  from.String(),
		Nonce: atomic.AddUint64(&nonce1, 1),
	}, rpcx.String(suite.MockDappID(pk)), rpcx.String("reason"))
	suite.Require().Nil(err)
	if res.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf(string(res.Ret))
	}
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	suite.Require().Nil(err)
	err = suite.VotePass(result.ProposalID)
	suite.Require().Nil(err)
	return nil
}

// ActivateDapp activate dapp
func (suite Model14) ActivateDapp(pk crypto.PrivateKey) error {
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.DappMgrContractAddr.Address(), "ActivateDapp", nil, rpcx.String(suite.MockDappID(pk)), rpcx.String("reason"))
	suite.Require().Nil(err)
	if res.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf(string(res.Ret))
	}
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	suite.Require().Nil(err)
	err = suite.VotePass(result.ProposalID)
	suite.Require().Nil(err)
	return nil
}

// TransferDapp transfer dapp from pk1 to pk2
func (suite Model14) TransferDapp(pk1, pk2 crypto.PrivateKey) error {
	client := suite.NewClient(pk1)
	address, err := pk2.PublicKey().Address()
	suite.Require().Nil(err)
	res, err := client.InvokeBVMContract(constant.DappMgrContractAddr.Address(), "TransferDapp", nil,
		rpcx.String(suite.MockDappID(pk1)), rpcx.String(address.String()), rpcx.String("reason"))
	suite.Require().Nil(err)
	if res.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf(string(res.Ret))
	}
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	suite.Require().Nil(err)
	err = suite.VotePass(result.ProposalID)
	suite.Require().Nil(err)
	err = suite.ConfirmTransfer(pk1, pk2)
	suite.Require().Nil(err)
	return nil
}

// ConfirmTransfer confirm transfer dapp
func (suite Model14) ConfirmTransfer(pk1, pk2 crypto.PrivateKey) error {
	client := suite.NewClient(pk2)
	res, err := client.InvokeBVMContract(constant.DappMgrContractAddr.Address(), "ConfirmTransfer", nil, rpcx.String(suite.MockDappID(pk1)))
	suite.Require().Nil(err)
	if res.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf(string(res.Ret))
	}
	return nil
}

// EvaluateDapp evaluate dapp [0-5]
func (suite Model14) EvaluateDapp(pk crypto.PrivateKey, id, desc string, score float64) error {
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.DappMgrContractAddr.Address(), "EvaluateDapp", nil,
		rpcx.String(id), rpcx.String(desc), rpcx.Float64(score))
	suite.Require().Nil(err)
	if res.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf(string(res.Ret))
	}
	return nil
}

// CheckDappStatus check dapp status
func (suite Snake) CheckDappStatus(id string, status governance.GovernanceStatus) error {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.DappMgrContractAddr.Address(), "GetDapp", nil, rpcx.String(id))
	suite.Require().Nil(err)
	if res.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf(string(res.Ret))
	}
	dapp := &Dapp{}
	err = json.Unmarshal(res.Ret, dapp)
	suite.Require().Nil(err)
	if dapp.Status != status {
		return fmt.Errorf("expect status is %s, but got %s", status, dapp.Status)
	}
	return nil
}

// MockDappID mock first dapp ID
func (suite Snake) MockDappID(pk crypto.PrivateKey) string {
	address, _ := pk.PublicKey().Address()
	return address.String() + "-0"
}

// DappToAvailable get an available dapp
func (suite Model14) DappToAvailable(pk crypto.PrivateKey, address string) error {
	err := suite.RegisterDapp(pk, address)
	suite.Require().Nil(err)
	err = suite.CheckDappStatus(suite.MockDappID(pk), governance.GovernanceAvailable)
	suite.Require().Nil(err)
	return nil
}

// DappToUnavailable get an unavailable dapp
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
	suite.Require().Nil(err)
	if res.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf(string(res.Ret))
	}
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	suite.Require().Nil(err)
	err = suite.VoteReject(result.ProposalID)
	suite.Require().Nil(err)
	err = suite.CheckDappStatus(suite.MockDappID(pk), governance.GovernanceUnavailable)
	suite.Require().Nil(err)
	return nil
}

// DappToActivating get an activating dapp
func (suite Model14) DappToActivating(pk crypto.PrivateKey, address string) error {
	err := suite.RegisterDapp(pk, address)
	suite.Require().Nil(err)
	err = suite.FreezeDapp(pk)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.DappMgrContractAddr.Address(), "ActivateDapp", nil, rpcx.String(suite.MockDappID(pk)), rpcx.String("reason"))
	suite.Require().Nil(err)
	if res.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf(string(res.Ret))
	}
	err = suite.CheckDappStatus(suite.MockDappID(pk), governance.GovernanceActivating)
	suite.Require().Nil(err)
	return nil
}

// DappToUpdating get an updating dapp
func (suite Model14) DappToUpdating(pk crypto.PrivateKey, address string) error {
	err := suite.RegisterDapp(pk, address)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(suite.MockDappID(pk)),
		rpcx.String(address + "123"),
		rpcx.String("test"),
		rpcx.String("https://github.com/meshplus/bitxhub"),
		rpcx.String(address),
		rpcx.String(""),
		rpcx.String("reason"),
	}
	res, err := client.InvokeBVMContract(constant.DappMgrContractAddr.Address(), "UpdateDapp", nil, args...)
	suite.Require().Nil(err)
	if res.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf(string(res.Ret))
	}
	err = suite.CheckDappStatus(suite.MockDappID(pk), governance.GovernanceUpdating)
	suite.Require().Nil(err)
	return nil
}

// DappToFreezing get a freezing dapp
func (suite Model14) DappToFreezing(pk crypto.PrivateKey, address string) error {
	err := suite.RegisterDapp(pk, address)
	suite.Require().Nil(err)
	node1pk, from, err := repo.Node1Priv()
	suite.Require().Nil(err)
	client := suite.NewClient(node1pk)
	res, err := client.InvokeBVMContract(constant.DappMgrContractAddr.Address(), "FreezeDapp", &rpcx.TransactOpts{
		From:  from.String(),
		Nonce: atomic.AddUint64(&nonce1, 1),
	}, rpcx.String(suite.MockDappID(pk)), rpcx.String("reason"))
	suite.Require().Nil(err)
	if res.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf(string(res.Ret))
	}
	err = suite.CheckDappStatus(suite.MockDappID(pk), governance.GovernanceFreezing)
	suite.Require().Nil(err)
	return nil
}

// DappToFrozen get a frozen dapp
func (suite Model14) DappToFrozen(pk crypto.PrivateKey, address string) error {
	err := suite.RegisterDapp(pk, address)
	suite.Require().Nil(err)
	err = suite.FreezeDapp(pk)
	suite.Require().Nil(err)
	err = suite.CheckDappStatus(suite.MockDappID(pk), governance.GovernanceFrozen)
	suite.Require().Nil(err)
	return nil
}
