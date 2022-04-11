package bxh_tester

import (
	"encoding/json"
	"fmt"
	"strings"
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

type Model17 struct {
	*Snake
}
type ProposalStrategy struct {
	Module string                      `json:"module"`
	Typ    string                      `json:"typ"`
	Extra  string                      `json:"extra"`
	Status governance.GovernanceStatus `json:"status"`
	FSM    *fsm.FSM                    `json:"fsm"`
}

const (
	AppchainMgr    = "appchain_mgr"
	RuleMgr        = "rule_mgr"
	NodeMgr        = "node_mgr"
	ServiceMgr     = "service_mgr"
	StrategyMgr    = "proposal_strategy_mgr"
	DappMgr        = "dapp_mgr"
	RoleMgr        = "role_mgr"
	ZeroPermission = "ZeroPermission"
	SimpleMajority = "SimpleMajority"
)

//tc：更新应用链管理的策略为简单治理策略，策略更新成功
func (suite Model17) Test1701_UpdateAppchainMgrZeroPermissionIsSuccess() {
	err := suite.StrategyToSimple(AppchainMgr)
	suite.Require().Nil(err)
	err = suite.UpdateStrategy(AppchainMgr, ZeroPermission, "")
	suite.Require().Nil(err)
	strategy, err := suite.GetStrategyByType(AppchainMgr)
	suite.Require().Nil(err)
	suite.Require().Equal(ZeroPermission, strategy.Typ)
}

//tc：更新验证规则管理的策略为简单治理策略，策略更新成功
func (suite Model17) Test1702_UpdateRuleMgrZeroPermissionIsSuccess() {
	err := suite.StrategyToSimple(RuleMgr)
	suite.Require().Nil(err)
	err = suite.UpdateStrategy(RuleMgr, ZeroPermission, "")
	suite.Require().Nil(err)
	strategy, err := suite.GetStrategyByType(RuleMgr)
	suite.Require().Nil(err)
	suite.Require().Equal(ZeroPermission, strategy.Typ)
}

//tc：更新节点管理的策略为简单治理策略，策略更新成功
func (suite Model17) Test1703_UpdateNodeMgrZeroPermissionIsSuccess() {
	err := suite.StrategyToSimple(NodeMgr)
	suite.Require().Nil(err)
	err = suite.UpdateStrategy(NodeMgr, ZeroPermission, "")
	suite.Require().Nil(err)
	strategy, err := suite.GetStrategyByType(NodeMgr)
	suite.Require().Nil(err)
	suite.Require().Equal(ZeroPermission, strategy.Typ)
}

//tc：更新服务管理的策略为简单治理策略，策略更新成功
func (suite Model17) Test1704_UpdateServiceMgrZeroPermissionIsSuccess() {
	err := suite.StrategyToSimple(ServiceMgr)
	suite.Require().Nil(err)
	err = suite.UpdateStrategy(ServiceMgr, ZeroPermission, "")
	suite.Require().Nil(err)
	strategy, err := suite.GetStrategyByType(ServiceMgr)
	suite.Require().Nil(err)
	suite.Require().Equal(ZeroPermission, strategy.Typ)
}

//tc：更新投票策略管理的策略为简单治理策略，策略更新成功
func (suite Model17) Test1705_UpdateStrategyMgrZeroPermissionIsSuccess() {
	err := suite.StrategyToSimple(StrategyMgr)
	suite.Require().Nil(err)
	err = suite.UpdateStrategy(StrategyMgr, ZeroPermission, "")
	suite.Require().Nil(err)
	strategy, err := suite.GetStrategyByType(StrategyMgr)
	suite.Require().Nil(err)
	suite.Require().Equal(ZeroPermission, strategy.Typ)
}

//tc：更新身份管理的策略为简单治理策略，策略更新成功
func (suite Model17) Test1706_UpdateRoleMgrZeroPermissionIsSuccess() {
	err := suite.StrategyToSimple(RoleMgr)
	suite.Require().Nil(err)
	err = suite.UpdateStrategy(RoleMgr, ZeroPermission, "")
	suite.Require().Nil(err)
	strategy, err := suite.GetStrategyByType(RoleMgr)
	suite.Require().Nil(err)
	suite.Require().Equal(ZeroPermission, strategy.Typ)
}

//tc：更新Dapp管理的策略为简单治理策略，策略更新成功
func (suite Model17) Test1707_UpdateDappMgrZeroPermissionIsSuccess() {
	err := suite.StrategyToSimple(DappMgr)
	suite.Require().Nil(err)
	err = suite.UpdateStrategy(DappMgr, ZeroPermission, "")
	suite.Require().Nil(err)
	strategy, err := suite.GetStrategyByType(DappMgr)
	suite.Require().Nil(err)
	suite.Require().Equal(ZeroPermission, strategy.Typ)
}

//tc：更新全部模块的策略为简单治理策略，策略更新成功
func (suite Model17) Test1708_UpdateAllStrategyZeroIsSuccess() {
	err := suite.StrategyToSimple(RuleMgr)
	suite.Require().Nil(err)
	err = suite.UpdateAllStrategy(ZeroPermission, "")
	suite.Require().Nil(err)
}

//tc：更新应用链管理的策略为简单治理策略，策略更新成功
func (suite Model17) Test1709_UpdateAppchainMgrSimpleMajorityIsSuccess() {
	err := suite.StrategyToZero(AppchainMgr)
	suite.Require().Nil(err)
	err = suite.UpdateStrategy(AppchainMgr, SimpleMajority, "a > 0.5 * t")
	suite.Require().Nil(err)
	strategy, err := suite.GetStrategyByType(AppchainMgr)
	suite.Require().Nil(err)
	suite.Require().Equal(SimpleMajority, strategy.Typ)
}

//tc：更新验证规则管理的策略为简单多数策略，策略更新成功
func (suite Model17) Test1710_UpdateRuleMgrSimpleMajorityIsSuccess() {
	err := suite.StrategyToZero(RuleMgr)
	suite.Require().Nil(err)
	err = suite.UpdateStrategy(RuleMgr, SimpleMajority, "a > 0.5 * t")
	suite.Require().Nil(err)
	strategy, err := suite.GetStrategyByType(RuleMgr)
	suite.Require().Nil(err)
	suite.Require().Equal(SimpleMajority, strategy.Typ)
}

//tc：更新节点管理的策略为简单多数策略，策略更新成功
func (suite Model17) Test1711_UpdateNodeMgrSimpleMajorityIsSuccess() {
	err := suite.StrategyToZero(NodeMgr)
	suite.Require().Nil(err)
	err = suite.UpdateStrategy(NodeMgr, SimpleMajority, "a > 0.5 * t")
	suite.Require().Nil(err)
	strategy, err := suite.GetStrategyByType(NodeMgr)
	suite.Require().Nil(err)
	suite.Require().Equal(SimpleMajority, strategy.Typ)
}

//tc：更新服务管理的策略为简单多数策略，策略更新成功
func (suite Model17) Test1712_UpdateServiceMgrSimpleMajorityIsSuccess() {
	err := suite.StrategyToZero(ServiceMgr)
	suite.Require().Nil(err)
	err = suite.UpdateStrategy(ServiceMgr, SimpleMajority, "a > 0.5 * t")
	suite.Require().Nil(err)
	strategy, err := suite.GetStrategyByType(ServiceMgr)
	suite.Require().Nil(err)
	suite.Require().Equal(SimpleMajority, strategy.Typ)
}

//tc：更新投票策略管理的策略为简单多数策略，策略更新成功
func (suite Model17) Test1713_UpdateStrategyMgrSimpleMajorityIsSuccess() {
	err := suite.StrategyToZero(StrategyMgr)
	suite.Require().Nil(err)
	err = suite.UpdateStrategy(StrategyMgr, SimpleMajority, "a > 0.5 * t")
	suite.Require().Nil(err)
	strategy, err := suite.GetStrategyByType(StrategyMgr)
	suite.Require().Nil(err)
	suite.Require().Equal(SimpleMajority, strategy.Typ)
}

//tc：更新身份管理的策略为简单多数策略，策略更新成功
func (suite Model17) Test1714_UpdateRoleMgrSimpleMajorityIsSuccess() {
	err := suite.StrategyToZero(RoleMgr)
	suite.Require().Nil(err)
	err = suite.UpdateStrategy(RoleMgr, SimpleMajority, "a > 0.5 * t")
	suite.Require().Nil(err)
	strategy, err := suite.GetStrategyByType(RoleMgr)
	suite.Require().Nil(err)
	suite.Require().Equal(SimpleMajority, strategy.Typ)
}

//tc：更新Dapp管理的策略为简单多数策略，策略更新成功
func (suite Model17) Test1715_UpdateDappMgrSimpleMajorityIsSuccess() {
	err := suite.StrategyToZero(DappMgr)
	suite.Require().Nil(err)
	err = suite.UpdateStrategy(DappMgr, SimpleMajority, "a > 0.5 * t")
	suite.Require().Nil(err)
	strategy, err := suite.GetStrategyByType(DappMgr)
	suite.Require().Nil(err)
	suite.Require().Equal(SimpleMajority, strategy.Typ)
}

//tc：更新全部模块的策略为简单多数策略，策略更新成功
func (suite Model17) Test1716_UpdateAllStrategySimpleIsSuccess() {
	err := suite.StrategyToZero(RuleMgr)
	suite.Require().Nil(err)
	err = suite.UpdateAllStrategy(SimpleMajority, "a > 0.5 * t")
	suite.Require().Nil(err)
}

//tc：更新投票策略为简单多数策略，公式为t==4，增加治理管理员，公式变为默认公式
func (suite Model17) Test1717_UpdateStrategyBeforeAdminAdded() {
	err := suite.StrategyToSimple(StrategyMgr)
	suite.Require().Nil(err)
	err = suite.UpdateStrategy(StrategyMgr, SimpleMajority, "t==4")
	suite.Require().Nil(err)
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	from, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	err = suite.RegisterRole(from.String(), GovernanceAdmin, "")
	suite.Require().Nil(err)
	//recover
	err = suite.LogoutRole(from.String())
	suite.Require().Nil(err)
	err = suite.CheckRoleStatus(from.String(), governance.GovernanceForbidden)
	suite.Require().Nil(err)
}

// UpdateStrategy update strategy
func (suite Model17) UpdateStrategy(model, typ, extra string) error {
	pk, from, err := repo.Node1Priv()
	if err != nil {
		return err
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ProposalStrategyMgrContractAddr.Address(), "UpdateProposalStrategy", &rpcx.TransactOpts{
		From:  from.String(),
		Nonce: atomic.AddUint64(&nonce1, 1),
	}, rpcx.String(model), rpcx.String(typ), rpcx.String(extra), rpcx.String("reason"))
	if err != nil {
		return err
	}
	if pb.Receipt_SUCCESS != res.Status {
		return fmt.Errorf(string(res.Ret))
	}
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, &result)
	if err != nil {
		return err
	}
	err = suite.VotePass(result.ProposalID)
	if err != nil {
		return err
	}
	return nil
}

// GetStrategyByType get strategy by model type
func (suite Model17) GetStrategyByType(typ string) (*ProposalStrategy, error) {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	if err != nil {
		return nil, err
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ProposalStrategyMgrContractAddr.Address(), "GetProposalStrategy", nil, rpcx.String(typ))
	if err != nil {
		return nil, err
	}
	strategy := ProposalStrategy{}
	err = json.Unmarshal(res.Ret, &strategy)
	if err != nil {
		return nil, err
	}
	return &strategy, nil
}

// StrategyToZero update strategy to ZeroPermission
func (suite Model17) StrategyToZero(model string) error {
	strategy, err := suite.GetStrategyByType(model)
	if err != nil {
		return err
	}
	if strategy.Typ == ZeroPermission {
		return nil
	}
	pk, from, err := repo.Node1Priv()
	if err != nil {
		return err
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ProposalStrategyMgrContractAddr.Address(), "UpdateProposalStrategy", &rpcx.TransactOpts{
		From:  from.String(),
		Nonce: atomic.AddUint64(&nonce1, 1),
	}, rpcx.String(model), rpcx.String(ZeroPermission), rpcx.String(""), rpcx.String("reason"))
	if err != nil {
		return err
	}
	if pb.Receipt_FAILED == res.Status {
		return fmt.Errorf(string(res.Ret))
	}
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, &result)
	if err != nil {
		return err
	}
	err = suite.VotePass(result.ProposalID)
	if err != nil {
		return err
	}
	return nil
}

// StrategyToSimple update strategy to SimpleMajority
func (suite Model17) StrategyToSimple(model string) error {
	strategy, err := suite.GetStrategyByType(model)
	if err != nil {
		return err
	}
	if strategy.Typ == SimpleMajority {
		return nil
	}
	pk, from, err := repo.Node1Priv()
	if err != nil {
		return err
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ProposalStrategyMgrContractAddr.Address(), "UpdateProposalStrategy", &rpcx.TransactOpts{
		From:  from.String(),
		Nonce: atomic.AddUint64(&nonce1, 1),
	}, rpcx.String(model), rpcx.String(SimpleMajority), rpcx.String("a > 0.5 * t"), rpcx.String("reason"))
	if err != nil {
		return err
	}
	if pb.Receipt_FAILED == res.Status {
		return fmt.Errorf(string(res.Ret))
	}
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, &result)
	if err != nil {
		return err
	}
	err = suite.VotePass(result.ProposalID)
	if err != nil {
		return err
	}
	return nil
}

// StrategyToUpdating update strategy to updating status
func (suite Model17) StrategyToUpdating(model string) (*RegisterResult, error) {
	err := suite.StrategyToSimple(StrategyMgr)
	if err != nil {
		return nil, err
	}
	err = suite.StrategyToSimple(model)
	if err != nil {
		return nil, err
	}
	pk, from, err := repo.Node1Priv()
	if err != nil {
		return nil, err
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ProposalStrategyMgrContractAddr.Address(), "UpdateProposalStrategy", &rpcx.TransactOpts{
		From:  from.String(),
		Nonce: atomic.AddUint64(&nonce1, 1),
	}, rpcx.String(model), rpcx.String(ZeroPermission), rpcx.String(""), rpcx.String("reason"))
	if err != nil {
		return nil, err
	}
	if pb.Receipt_SUCCESS != res.Status {
		return nil, fmt.Errorf(string(res.Ret))
	}
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, &result)
	if err != nil {
		return nil, err
	}
	strategy, err := suite.GetStrategyByType(model)
	if err != nil {
		return nil, err
	}
	if governance.GovernanceUpdating != strategy.Status {
		return nil, fmt.Errorf("%v is %v", strategy.Module, strategy.Status)
	}
	return result, nil
}

// UpdateAllStrategy update all model strategy
func (suite Model17) UpdateAllStrategy(typ, extra string) error {
	pk, from, err := repo.Node1Priv()
	if err != nil {
		return err
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ProposalStrategyMgrContractAddr.Address(), "UpdateAllProposalStrategy", &rpcx.TransactOpts{
		From:  from.String(),
		Nonce: atomic.AddUint64(&nonce1, 1),
	}, rpcx.String(typ), rpcx.String(extra), rpcx.String("reason"))
	if err != nil {
		return err
	}
	if pb.Receipt_SUCCESS != res.Status {
		return fmt.Errorf(string(res.Ret))
	}
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, &result)
	if err != nil {
		return err
	}
	err = suite.VotePass(result.ProposalID)
	if err != nil {
		return err
	}
	return nil
}

// AllStrategyToUpdating update all strategy to updating status
func (suite Model17) AllStrategyToUpdating() (*RegisterResult, error) {
	err := suite.UpdateAllStrategy(SimpleMajority, "a > 0.5 * t")
	if err != nil && !strings.Contains(err.Error(), "no strategy information is updated") {
		return nil, err
	}
	pk, from, err := repo.Node1Priv()
	if err != nil {
		return nil, err
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ProposalStrategyMgrContractAddr.Address(), "UpdateAllProposalStrategy", &rpcx.TransactOpts{
		From:  from.String(),
		Nonce: atomic.AddUint64(&nonce1, 1),
	}, rpcx.String(ZeroPermission), rpcx.String(""), rpcx.String("reason"))
	if err != nil {
		return nil, err
	}
	if pb.Receipt_SUCCESS != res.Status {
		return nil, fmt.Errorf(string(res.Ret))
	}
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
