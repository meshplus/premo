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
	Module               string                      `json:"module"`
	Typ                  string                      `json:"typ"`
	ParticipateThreshold float64                     `json:"participate_threshold"`
	Extra                []byte                      `json:"extra"`
	Status               governance.GovernanceStatus `json:"status"`
	FSM                  *fsm.FSM                    `json:"fsm"`
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
	err = suite.UpdateStrategy(AppchainMgr, ZeroPermission, 0)
	suite.Require().Nil(err)
	strategy, err := suite.GetStrategyByType(AppchainMgr)
	suite.Require().Nil(err)
	suite.Require().Equal(ZeroPermission, strategy.Typ)
}

//tc：更新验证规则管理的策略为简单治理策略，策略更新成功
func (suite Model17) Test1702_UpdateRuleMgrZeroPermissionIsSuccess() {
	err := suite.StrategyToSimple(RuleMgr)
	suite.Require().Nil(err)
	err = suite.UpdateStrategy(RuleMgr, ZeroPermission, 0)
	suite.Require().Nil(err)
	strategy, err := suite.GetStrategyByType(RuleMgr)
	suite.Require().Nil(err)
	suite.Require().Equal(ZeroPermission, strategy.Typ)
}

//tc：更新节点管理的策略为简单治理策略，策略更新成功
func (suite Model17) Test1703_UpdateNodeMgrZeroPermissionIsSuccess() {
	err := suite.StrategyToSimple(NodeMgr)
	suite.Require().Nil(err)
	err = suite.UpdateStrategy(NodeMgr, ZeroPermission, 0)
	suite.Require().Nil(err)
	strategy, err := suite.GetStrategyByType(NodeMgr)
	suite.Require().Nil(err)
	suite.Require().Equal(ZeroPermission, strategy.Typ)
}

//tc：更新服务管理的策略为简单治理策略，策略更新成功
func (suite Model17) Test1704_UpdateServiceMgrZeroPermissionIsSuccess() {
	err := suite.StrategyToSimple(ServiceMgr)
	suite.Require().Nil(err)
	err = suite.UpdateStrategy(ServiceMgr, ZeroPermission, 0)
	suite.Require().Nil(err)
	strategy, err := suite.GetStrategyByType(ServiceMgr)
	suite.Require().Nil(err)
	suite.Require().Equal(ZeroPermission, strategy.Typ)
}

//tc：更新投票策略管理的策略为简单治理策略，策略更新成功
func (suite Model17) Test1705_UpdateStrategyMgrZeroPermissionIsSuccess() {
	err := suite.StrategyToSimple(StrategyMgr)
	suite.Require().Nil(err)
	err = suite.UpdateStrategy(StrategyMgr, ZeroPermission, 0)
	suite.Require().Nil(err)
	strategy, err := suite.GetStrategyByType(StrategyMgr)
	suite.Require().Nil(err)
	suite.Require().Equal(ZeroPermission, strategy.Typ)
}

//tc：更新身份管理的策略为简单治理策略，策略更新成功
func (suite Model17) Test1706_UpdateRoleMgrZeroPermissionIsSuccess() {
	err := suite.StrategyToSimple(RoleMgr)
	suite.Require().Nil(err)
	err = suite.UpdateStrategy(RoleMgr, ZeroPermission, 0)
	suite.Require().Nil(err)
	strategy, err := suite.GetStrategyByType(RoleMgr)
	suite.Require().Nil(err)
	suite.Require().Equal(ZeroPermission, strategy.Typ)
}

//tc：更新Dapp管理的策略为简单治理策略，策略更新成功
func (suite Model17) Test1707_UpdateDappMgrZeroPermissionIsSuccess() {
	err := suite.StrategyToSimple(DappMgr)
	suite.Require().Nil(err)
	err = suite.UpdateStrategy(DappMgr, ZeroPermission, 0)
	suite.Require().Nil(err)
	strategy, err := suite.GetStrategyByType(DappMgr)
	suite.Require().Nil(err)
	suite.Require().Equal(ZeroPermission, strategy.Typ)
}

//tc：更新应用链管理的策略为简单治理策略，提案未通过，其他管理员发起应用链管理投票策略变更的提案，提案发起失败
func (suite Model17) Test1708_UpdateAppchainMgrZeroPermissionWithUpdatingIsFail() {
	result, err := suite.StrategyToUpdating(AppchainMgr)
	suite.Require().Nil(err)
	//recover strategy
	defer func() {
		err = suite.VotePass(result.ProposalID)
		suite.Require().Nil(err)
	}()
	err = suite.UpdateStrategy(AppchainMgr, ZeroPermission, 0)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "is updating, can not do update")
}

//tc：更新验证规则的策略为简单治理策略，提案未通过，其他管理员发起应用链管理投票策略变更的提案，提案发起失败
func (suite Model17) Test1709_UpdateRuleMgrZeroPermissionWithUpdatingIsFail() {
	result, err := suite.StrategyToUpdating(RuleMgr)
	suite.Require().Nil(err)
	//recover strategy
	defer func() {
		err = suite.VotePass(result.ProposalID)
		suite.Require().Nil(err)
	}()
	err = suite.UpdateStrategy(RuleMgr, ZeroPermission, 0)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "is updating, can not do update")
}

//tc：更新节点管理的策略为简单治理策略，提案未通过，其他管理员发起应用链管理投票策略变更的提案，提案发起失败
func (suite Model17) Test1710_UpdateNodeMgrZeroPermissionWithUpdatingIsFail() {
	result, err := suite.StrategyToUpdating(NodeMgr)
	suite.Require().Nil(err)
	//recover strategy
	defer func() {
		err = suite.VotePass(result.ProposalID)
		suite.Require().Nil(err)
	}()
	err = suite.UpdateStrategy(NodeMgr, ZeroPermission, 0)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "is updating, can not do update")
}

//tc：更新服务管理的策略为简单治理策略，提案未通过，其他管理员发起应用链管理投票策略变更的提案，提案发起失败
func (suite Model17) Test1711_UpdateServiceMgrZeroPermissionWithUpdatingIsFail() {
	result, err := suite.StrategyToUpdating(ServiceMgr)
	suite.Require().Nil(err)
	//recover strategy
	defer func() {
		err = suite.VotePass(result.ProposalID)
		suite.Require().Nil(err)
	}()
	err = suite.UpdateStrategy(ServiceMgr, ZeroPermission, 0)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "is updating, can not do update")
}

//tc：更新投票管理的策略为简单治理策略，提案未通过，其他管理员发起应用链管理投票策略变更的提案，提案发起失败
func (suite Model17) Test1712_UpdateStrategyMgrZeroPermissionWithUpdatingIsFail() {
	result, err := suite.StrategyToUpdating(StrategyMgr)
	suite.Require().Nil(err)
	//recover strategy
	defer func() {
		err = suite.VotePass(result.ProposalID)
		suite.Require().Nil(err)
	}()
	err = suite.UpdateStrategy(StrategyMgr, ZeroPermission, 0)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "is updating, can not do update")
}

//tc：更新身份管理的策略为简单治理策略，提案未通过，其他管理员发起应用链管理投票策略变更的提案，提案发起失败
func (suite Model17) Test1713_UpdateRoleMgrZeroPermissionWithUpdatingIsFail() {
	result, err := suite.StrategyToUpdating(RoleMgr)
	suite.Require().Nil(err)
	//recover strategy
	defer func() {
		err = suite.VotePass(result.ProposalID)
		suite.Require().Nil(err)
	}()
	err = suite.UpdateStrategy(RoleMgr, ZeroPermission, 0)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "is updating, can not do update")
}

//tc：更新Dapp管理的策略为简单治理策略，提案未通过，其他管理员发起应用链管理投票策略变更的提案，提案发起失败
func (suite Model17) Test1714_UpdateDappMgrZeroPermissionWithUpdatingIsFail() {
	result, err := suite.StrategyToUpdating(DappMgr)
	suite.Require().Nil(err)
	//recover strategy
	defer func() {
		err = suite.VotePass(result.ProposalID)
		suite.Require().Nil(err)
	}()
	err = suite.UpdateStrategy(DappMgr, ZeroPermission, 0)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "is updating, can not do update")
}

//tc：更新全部模块的策略为简单治理策略，策略更新成功
func (suite Model17) Test1715_UpdateAllStrategyZeroIsSuccess() {
	err := suite.StrategyToSimple(RuleMgr)
	suite.Require().Nil(err)
	err = suite.UpdateAllStrategy(ZeroPermission, 0)
	suite.Require().Nil(err)
}

//tc：更新某一模块的投票策略，提案未通过，更新全部投票策略失败
func (suite Model17) Test1716_UpdateAllStrategyWithUpdatingRuleMgrIsFail() {
	result, err := suite.StrategyToUpdating(RuleMgr)
	suite.Require().Nil(err)
	//recover strategy
	defer func() {
		err := suite.VotePass(result.ProposalID)
		suite.Require().Nil(err)
	}()
	err = suite.UpdateAllStrategy(ZeroPermission, 0)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "is updating, can not do update")
}

//tc：更新全部投票策略，提案未通过，更新某一模块的投票策略失败
func (suite Model17) Test1717_UpdateStrategyWithUpdatingAllStrategyIsFail() {
	result, err := suite.AllStrategyToUpdating()
	suite.Require().Nil(err)
	//recover strategy
	defer func() {
		err := suite.VotePass(result.ProposalID)
		suite.Require().Nil(err)
	}()
	err = suite.UpdateStrategy(RuleMgr, ZeroPermission, 0)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "is updating, can not do update")
}

//tc：更新投票策略为简单治理策略，投票的阈值为-1，策略更新失败
func (suite Model17) Test1718_UpdateStrategyZeroWithThreshold_1() {
	err := suite.StrategyToSimple(RuleMgr)
	suite.Require().Nil(err)
	err = suite.UpdateStrategy(RuleMgr, ZeroPermission, -1)
	suite.Require().NotNil(err)
}

//tc：更新投票策略为简单治理策略，投票的阈值为0.5，策略更新失败
func (suite Model17) Test1719_UpdateStrategyZeroWithThreshold_2() {
	err := suite.StrategyToSimple(RuleMgr)
	suite.Require().Nil(err)
	err = suite.UpdateStrategy(RuleMgr, ZeroPermission, 0.5)
	suite.Require().NotNil(err)
}

//tc：更新投票策略为简单治理策略，投票的阈值为2，策略更新失败
func (suite Model17) Test1720_UpdateStrategyZeroWithThreshold_3() {
	err := suite.StrategyToSimple(RuleMgr)
	suite.Require().Nil(err)
	err = suite.UpdateStrategy(RuleMgr, ZeroPermission, 2)
	suite.Require().NotNil(err)
}

//tc：更新应用链管理的策略为简单治理策略，策略更新成功
func (suite Model17) Test1721_UpdateAppchainMgrSimpleMajorityIsSuccess() {
	err := suite.StrategyToZero(AppchainMgr)
	suite.Require().Nil(err)
	err = suite.UpdateStrategy(AppchainMgr, SimpleMajority, 0.75)
	suite.Require().Nil(err)
	strategy, err := suite.GetStrategyByType(AppchainMgr)
	suite.Require().Nil(err)
	suite.Require().Equal(SimpleMajority, strategy.Typ)
}

//tc：更新验证规则管理的策略为简单多数策略，策略更新成功
func (suite Model17) Test1722_UpdateRuleMgrSimpleMajorityIsSuccess() {
	err := suite.StrategyToZero(RuleMgr)
	suite.Require().Nil(err)
	err = suite.UpdateStrategy(RuleMgr, SimpleMajority, 0.75)
	suite.Require().Nil(err)
	strategy, err := suite.GetStrategyByType(RuleMgr)
	suite.Require().Nil(err)
	suite.Require().Equal(SimpleMajority, strategy.Typ)
}

//tc：更新节点管理的策略为简单多数策略，策略更新成功
func (suite Model17) Test1723_UpdateNodeMgrSimpleMajorityIsSuccess() {
	err := suite.StrategyToZero(NodeMgr)
	suite.Require().Nil(err)
	err = suite.UpdateStrategy(NodeMgr, SimpleMajority, 0.75)
	suite.Require().Nil(err)
	strategy, err := suite.GetStrategyByType(NodeMgr)
	suite.Require().Nil(err)
	suite.Require().Equal(SimpleMajority, strategy.Typ)
}

//tc：更新服务管理的策略为简单多数策略，策略更新成功
func (suite Model17) Test1724_UpdateServiceMgrSimpleMajorityIsSuccess() {
	err := suite.StrategyToZero(ServiceMgr)
	suite.Require().Nil(err)
	err = suite.UpdateStrategy(ServiceMgr, SimpleMajority, 0.75)
	suite.Require().Nil(err)
	strategy, err := suite.GetStrategyByType(ServiceMgr)
	suite.Require().Nil(err)
	suite.Require().Equal(SimpleMajority, strategy.Typ)
}

//tc：更新投票策略管理的策略为简单多数策略，策略更新成功
func (suite Model17) Test1725_UpdateStrategyMgrSimpleMajorityIsSuccess() {
	err := suite.StrategyToZero(StrategyMgr)
	suite.Require().Nil(err)
	err = suite.UpdateStrategy(StrategyMgr, SimpleMajority, 0.75)
	suite.Require().Nil(err)
	strategy, err := suite.GetStrategyByType(StrategyMgr)
	suite.Require().Nil(err)
	suite.Require().Equal(SimpleMajority, strategy.Typ)
}

//tc：更新身份管理的策略为简单多数策略，策略更新成功
func (suite Model17) Test1726_UpdateRoleMgrSimpleMajorityIsSuccess() {
	err := suite.StrategyToZero(RoleMgr)
	suite.Require().Nil(err)
	err = suite.UpdateStrategy(RoleMgr, SimpleMajority, 0.75)
	suite.Require().Nil(err)
	strategy, err := suite.GetStrategyByType(RoleMgr)
	suite.Require().Nil(err)
	suite.Require().Equal(SimpleMajority, strategy.Typ)
}

//tc：更新Dapp管理的策略为简单多数策略，策略更新成功
func (suite Model17) Test1727_UpdateDappMgrSimpleMajorityIsSuccess() {
	err := suite.StrategyToZero(DappMgr)
	suite.Require().Nil(err)
	err = suite.UpdateStrategy(DappMgr, SimpleMajority, 0.75)
	suite.Require().Nil(err)
	strategy, err := suite.GetStrategyByType(DappMgr)
	suite.Require().Nil(err)
	suite.Require().Equal(SimpleMajority, strategy.Typ)
}

//tc：更新应用链管理的策略为简单多数策略，提案未通过，其他管理员发起应用链管理投票策略变更的提案，提案发起失败
func (suite Model17) Test1728_UpdateAppchainMgrSimpleMajorityWithUpdatingIsFail() {
	result, err := suite.StrategyToUpdating(AppchainMgr)
	suite.Require().Nil(err)
	//recover strategy
	defer func() {
		err = suite.VotePass(result.ProposalID)
		suite.Require().Nil(err)
	}()
	err = suite.UpdateStrategy(AppchainMgr, SimpleMajority, 0.75)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "is updating, can not do update")
}

//tc：更新验证规则的策略为简单治理策略，提案未通过，其他管理员发起应用链管理投票策略变更的提案，提案发起失败
func (suite Model17) Test1729_UpdateRuleMgrSimpleMajorityWithUpdatingIsFail() {
	result, err := suite.StrategyToUpdating(RuleMgr)
	suite.Require().Nil(err)
	//recover strategy
	defer func() {
		err = suite.VotePass(result.ProposalID)
		suite.Require().Nil(err)
	}()
	err = suite.UpdateStrategy(RuleMgr, SimpleMajority, 0.75)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "is updating, can not do update")
}

//tc：更新节点管理的策略为简单多数策略，提案未通过，其他管理员发起应用链管理投票策略变更的提案，提案发起失败
func (suite Model17) Test1730_UpdateNodeMgrSimpleMajorityWithUpdatingIsFail() {
	result, err := suite.StrategyToUpdating(NodeMgr)
	suite.Require().Nil(err)
	//recover strategy
	defer func() {
		err = suite.VotePass(result.ProposalID)
		suite.Require().Nil(err)
	}()
	err = suite.UpdateStrategy(NodeMgr, SimpleMajority, 0.75)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "is updating, can not do update")
}

//tc：更新服务管理的策略为简单多数策略，提案未通过，其他管理员发起应用链管理投票策略变更的提案，提案发起失败
func (suite Model17) Test1731_UpdateServiceMgrSimpleMajorityWithUpdatingIsFail() {
	result, err := suite.StrategyToUpdating(ServiceMgr)
	suite.Require().Nil(err)
	//recover strategy
	defer func() {
		err = suite.VotePass(result.ProposalID)
		suite.Require().Nil(err)
	}()
	err = suite.UpdateStrategy(ServiceMgr, SimpleMajority, 0.75)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "is updating, can not do update")
}

//tc：更新投票管理的策略为简单多数策略，提案未通过，其他管理员发起应用链管理投票策略变更的提案，提案发起失败
func (suite Model17) Test1732_UpdateStrategyMgrSimpleMajorityWithUpdatingIsFail() {
	result, err := suite.StrategyToUpdating(StrategyMgr)
	suite.Require().Nil(err)
	//recover strategy
	defer func() {
		err = suite.VotePass(result.ProposalID)
		suite.Require().Nil(err)
	}()
	err = suite.UpdateStrategy(StrategyMgr, SimpleMajority, 0.75)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "is updating, can not do update")
}

//tc：更新身份管理的策略为简单多数策略，提案未通过，其他管理员发起应用链管理投票策略变更的提案，提案发起失败
func (suite Model17) Test1733_UpdateRoleMgrSimpleMajorityWithUpdatingIsFail() {
	result, err := suite.StrategyToUpdating(RoleMgr)
	suite.Require().Nil(err)
	//recover strategy
	defer func() {
		err = suite.VotePass(result.ProposalID)
		suite.Require().Nil(err)
	}()
	err = suite.UpdateStrategy(RoleMgr, SimpleMajority, 0.75)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "is updating, can not do update")
}

//tc：更新Dapp管理的策略为简单多数策略，提案未通过，其他管理员发起应用链管理投票策略变更的提案，提案发起失败
func (suite Model17) Test1734_UpdateDappMgrSimpleMajorityWithUpdatingIsFail() {
	result, err := suite.StrategyToUpdating(DappMgr)
	suite.Require().Nil(err)
	//recover strategy
	defer func() {
		err = suite.VotePass(result.ProposalID)
		suite.Require().Nil(err)
	}()
	err = suite.UpdateStrategy(DappMgr, SimpleMajority, 0.75)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "is updating, can not do update")
}

//tc：更新全部模块的策略为简单多数策略，策略更新成功
func (suite Model17) Test1735_UpdateAllStrategySimpleIsSuccess() {
	err := suite.StrategyToSimple(RuleMgr)
	suite.Require().Nil(err)
	err = suite.UpdateAllStrategy(SimpleMajority, 0.75)
	suite.Require().Nil(err)
}

//tc：更新投票策略为简单多数策略，投票的阈值为-1，策略更新失败
func (suite Model17) Test1736_UpdateStrategySimpleWithThreshold_1() {
	err := suite.StrategyToSimple(RuleMgr)
	suite.Require().Nil(err)
	err = suite.UpdateStrategy(RuleMgr, SimpleMajority, -1)
	suite.Require().NotNil(err)
}

//tc：更新投票策略为简单多数策略，投票的阈值为0，策略更新失败
func (suite Model17) Test1737_UpdateStrategySimpleWithThreshold_2() {
	err := suite.StrategyToSimple(RuleMgr)
	suite.Require().Nil(err)
	err = suite.UpdateStrategy(RuleMgr, SimpleMajority, 0)
	suite.Require().NotNil(err)
}

//tc：更新投票策略为简单多数策略，投票的阈值为2，策略更新失败
func (suite Model17) Test1738_UpdateStrategySimpleWithThreshold_3() {
	err := suite.StrategyToSimple(RuleMgr)
	suite.Require().Nil(err)
	err = suite.UpdateStrategy(RuleMgr, SimpleMajority, 2)
	suite.Require().NotNil(err)
}

func (suite Model17) UpdateStrategy(model, typ string, threshold float64) error {
	path, err := repo.Node1Path()
	if err != nil {
		return err
	}
	pk, err := asym.RestorePrivateKey(path, repo.KeyPassword)
	if err != nil {
		return err
	}
	from, err := pk.PublicKey().Address()
	if err != nil {
		return err
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ProposalStrategyMgrContractAddr.Address(), "UpdateProposalStrategy", &rpcx.TransactOpts{
		From:  from.String(),
		Nonce: atomic.AddUint64(&nonce1, 1),
	}, rpcx.String(model), rpcx.String(typ), rpcx.Float64(threshold), rpcx.String("reason"))
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

func (suite Model17) GetStrategyByType(typ string) (ProposalStrategy, error) {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	if err != nil {
		return ProposalStrategy{}, err
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ProposalStrategyMgrContractAddr.Address(), "GetProposalStrategy", nil, rpcx.String(typ))
	if err != nil {
		return ProposalStrategy{}, err
	}
	strategy := ProposalStrategy{}
	err = json.Unmarshal(res.Ret, &strategy)
	if err != nil {
		return ProposalStrategy{}, err
	}
	return strategy, nil
}

func (suite Model17) StrategyToZero(model string) error {
	strategy, err := suite.GetStrategyByType(model)
	if err != nil {
		return err
	}
	if strategy.Typ == ZeroPermission {
		return nil
	}
	path, err := repo.Node1Path()
	if err != nil {
		return err
	}
	pk, err := asym.RestorePrivateKey(path, repo.KeyPassword)
	if err != nil {
		return err
	}
	from, err := pk.PublicKey().Address()
	if err != nil {
		return err
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ProposalStrategyMgrContractAddr.Address(), "UpdateProposalStrategy", &rpcx.TransactOpts{
		From:  from.String(),
		Nonce: atomic.AddUint64(&nonce1, 1),
	}, rpcx.String(model), rpcx.String(ZeroPermission), rpcx.Float64(0), rpcx.String("reason"))
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

func (suite Model17) StrategyToSimple(model string) error {
	strategy, err := suite.GetStrategyByType(model)
	if err != nil {
		return err
	}
	if strategy.Typ == SimpleMajority {
		return nil
	}
	path, err := repo.Node1Path()
	if err != nil {
		return err
	}
	pk, err := asym.RestorePrivateKey(path, repo.KeyPassword)
	if err != nil {
		return err
	}
	from, err := pk.PublicKey().Address()
	if err != nil {
		return err
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ProposalStrategyMgrContractAddr.Address(), "UpdateProposalStrategy", &rpcx.TransactOpts{
		From:  from.String(),
		Nonce: atomic.AddUint64(&nonce1, 1),
	}, rpcx.String(model), rpcx.String(SimpleMajority), rpcx.Float64(0.75), rpcx.String("reason"))
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

func (suite Model17) StrategyToUpdating(model string) (*RegisterResult, error) {
	err := suite.StrategyToSimple(StrategyMgr)
	if err != nil {
		return nil, err
	}
	err = suite.StrategyToSimple(model)
	if err != nil {
		return nil, err
	}
	path, err := repo.Node1Path()
	if err != nil {
		return nil, err
	}
	pk, err := asym.RestorePrivateKey(path, repo.KeyPassword)
	if err != nil {
		return nil, err
	}
	from, err := pk.PublicKey().Address()
	if err != nil {
		return nil, err
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ProposalStrategyMgrContractAddr.Address(), "UpdateProposalStrategy", &rpcx.TransactOpts{
		From:  from.String(),
		Nonce: atomic.AddUint64(&nonce1, 1),
	}, rpcx.String(model), rpcx.String(ZeroPermission), rpcx.Float64(0), rpcx.String("reason"))
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
	suite.Require().Equal(governance.GovernanceUpdating, strategy.Status)
	return result, nil
}

func (suite Model17) UpdateAllStrategy(typ string, threshold float64) error {
	path, err := repo.Node1Path()
	suite.Require().Nil(err)
	pk, err := asym.RestorePrivateKey(path, repo.KeyPassword)
	suite.Require().Nil(err)
	from, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ProposalStrategyMgrContractAddr.Address(), "UpdateAllProposalStrategy", &rpcx.TransactOpts{
		From:  from.String(),
		Nonce: atomic.AddUint64(&nonce1, 1),
	}, rpcx.String(typ), rpcx.Float64(threshold), rpcx.String("reason"))
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

func (suite Model17) AllStrategyToUpdating() (*RegisterResult, error) {
	err := suite.UpdateAllStrategy(SimpleMajority, 0.75)
	if err != nil && !strings.Contains(err.Error(), "no strategy information is updated") {
		return nil, err
	}
	path, err := repo.Node1Path()
	if err != nil {
		return nil, err
	}
	pk, err := asym.RestorePrivateKey(path, repo.KeyPassword)
	if err != nil {
		return nil, err
	}
	from, err := pk.PublicKey().Address()
	if err != nil {
		return nil, err
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ProposalStrategyMgrContractAddr.Address(), "UpdateAllProposalStrategy", &rpcx.TransactOpts{
		From:  from.String(),
		Nonce: atomic.AddUint64(&nonce1, 1),
	}, rpcx.String(ZeroPermission), rpcx.Float64(0), rpcx.String("reason"))
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
