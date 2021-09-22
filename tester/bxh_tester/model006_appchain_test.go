package bxh_tester

import (
	"encoding/json"
	"sync/atomic"

	"github.com/meshplus/bitxhub-core/governance"
	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/meshplus/premo/internal/repo"
	"github.com/pkg/errors"
)

type Model6 struct {
	*Snake
}

func (suite *Model6) SetupTest() {
	suite.T().Parallel()
}

//tc：非应用链管理员注册应用链，应用链注册失败
func (suite Model6) Test0601_RegisterAppchainWithNoAdminIsFail() {
	_, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().NotNil(err)
}

//tc：应用链管理员未注册，注册应用链。应用链注册成功
func (suite Model6) Test0602_RegisterAppchainWithNoRegisterAdminIsSuccess() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	chainID := suite.GetChainID(pk)
	err = suite.RegisterAppchain(pk, chainID, SimFabricRuleAddr)
	suite.Require().Nil(err)
}

//tc：应用链管理员已注册，注册应用链。应用链注册成功
func (suite Model6) Test0603_RegisterAppchainWithRegisteredAdminIsSuccess() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
}

//tc：应用链已注册，注册应用链，应用链注册失败
func (suite Model6) Test0604_RegisterAppchainWithRegisteredChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().NotNil(err)
}

//tc：应用链处于registing状态注册应用链，应用链注册失败
func (suite Model6) Test0605_RegisterAppchainWithRegistingChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ChainToRegisting(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().NotNil(err)
}

//tc：应用链处于unavailable状态注册应用链，应用链注册成功
func (suite Model6) Test0606_RegisterAppchainWithUnavailableChainIsSuccess() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ChainToUnavailable(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
}

//tc：应用链处于activating状态注册应用链，应用链注册失败
func (suite Model6) Test0607_RegisterAppchainWithActivatingChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ChainToActivating(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().NotNil(err)
}

//tc：应用链处于freezing状态注册应用链，应用链注册失败
func (suite Model6) Test0608_RegisterAppchainWithFreezingChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ChainToFreezing(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().NotNil(err)
}

//tc：应用链处于frozen状态注册应用链，应用链注册失败
func (suite Model6) Test0609_RegisterAppchainWithFrozenChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ChainToFrozen(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().NotNil(err)
}

//tc：应用链处于logouting状态注册应用链，应用链注册失败
func (suite Model6) Test0610_RegisterAppchainWithLogoutingChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ChainToLogouting(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().NotNil(err)
}

//tc：应用链处于forbidden状态注册应用链，应用链注册失败
func (suite Model6) Test0611_RegisterAppchainWithForbiddenChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ChainToForbidden(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().NotNil(err)
}

//tc：非应用链管理员更新应用链，应用链更新失败
func (suite Model6) Test0612_UpdateAppchainWithNoAdminIsFail() {
	pk1, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk1, chainID, address)
	suite.Require().Nil(err)
	pk2, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.UpdateAppchain(pk2, chainID, "test")
	suite.Require().NotNil(err)
}

//tc：应用链未注册，更新应用链。应用链更新失败
func (suite Model6) Test0613_UpdateAppchainWithNoRegisterChainIsFail() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.UpdateAppchain(pk, suite.GetChainID(pk), "test")
	suite.Require().NotNil(err)
}

//tc：应用链已注册，更新应用链，应用链更新成功
func (suite Model6) Test0614_UpdateAppchainIsSuccess() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.UpdateAppchain(pk, chainID, "test")
	suite.Require().Nil(err)
}

//tc：应用链处于registing状态更新应用链，应用链更新失败
func (suite Model6) Test0615_UpdateAppchainWithRegistingChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ChainToRegisting(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.UpdateAppchain(pk, chainID, "test")
	suite.Require().NotNil(err)
}

//tc：应用链处于unavailable状态更新应用链，应用链更新失败
func (suite Model6) Test0616_UpdateAppchainWithUnavailableChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ChainToUnavailable(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.UpdateAppchain(pk, chainID, "test")
	suite.Require().NotNil(err)
}

//tc：应用链处于activating状态更新应用链，应用链更新失败
func (suite Model6) Test0617_UpdateAppchainWithActivatingChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ChainToActivating(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.UpdateAppchain(pk, chainID, "test")
	suite.Require().NotNil(err)
}

//tc：应用链处于freezing状态更新应用链，应用链更新失败
func (suite Model6) Test0618_UpdateAppchainWithFreezingChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ChainToFreezing(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.UpdateAppchain(pk, chainID, "test")
	suite.Require().NotNil(err)
}

//tc：应用链处于frozen状态更新应用链，应用链更新成功
func (suite Model6) Test0619_UpdateAppchainWithFrozenChainIsSuccess() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ChainToFrozen(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.UpdateAppchain(pk, chainID, "test")
	suite.Require().Nil(err)
}

//tc：应用链处于logouting状态更新应用链，应用链更新失败
func (suite Model6) Test0620_UpdateAppchainWithLogoutingChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ChainToLogouting(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.UpdateAppchain(pk, chainID, "test")
	suite.Require().NotNil(err)
}

//tc：应用链处于forbidden状态更新应用链，应用链更新失败
func (suite Model6) Test0621_UpdateAppchainWithForbiddenChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ChainToForbidden(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.UpdateAppchain(pk, chainID, "test")
	suite.Require().NotNil(err)
}

//tc：非中继链管理员冻结应用链，应用链冻结失败
func (suite Model6) Test0622_FreezeAppchainWithNoAdminIsFail() {
	pk1, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk1, chainID, address)
	suite.Require().Nil(err)
	pk2, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk2)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "FreezeAppchain", nil, rpcx.String(chainID), rpcx.String("reason"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
}

//tc：应用链未注册，冻结应用链，应用链冻结失败
func (suite Model6) Test0623_FreezeAppchainWithNoRegisterChainIsFail() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "FreezeAppchain", nil, rpcx.String(suite.GetChainID(pk)), rpcx.String("reason"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
}

//tc：中继链管理员冻结应用链，应用链冻结成功
func (suite Model6) Test0624_FreezeAppchainIsSuccess() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.FreezeAppchain(chainID)
	suite.Require().Nil(err)
}

//tc：应用链处于registing状态冻结应用链，应用链冻结失败
func (suite Model6) Test0625_FreezeAppchainWithRegistingChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ChainToRegisting(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.FreezeAppchain(chainID)
	suite.Require().NotNil(err)
}

//tc：应用链处于unavailable状态冻结应用链，应用链冻结失败
func (suite Model6) Test0626_FreezeAppchainWithUnavailableChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ChainToUnavailable(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.FreezeAppchain(chainID)
	suite.Require().NotNil(err)
}

//tc：应用链处于activating状态冻结应用链，应用链冻结成功
func (suite Model6) Test0627_FreezeAppchainWithActivatingChainIsSuccess() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ChainToActivating(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.FreezeAppchain(chainID)
	suite.Require().Nil(err)
	err = suite.CheckChainStatus(chainID, governance.GovernanceFrozen)
	suite.Require().Nil(err)
}

//tc：应用链处于freezing状态冻结应用链，应用链冻结失败
func (suite Model6) Test0628_FreezeAppchainWithFreezingChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ChainToFreezing(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.FreezeAppchain(chainID)
	suite.Require().NotNil(err)
}

//tc：应用链处于frozen状态冻结应用链，应用链冻结失败
func (suite Model6) Test0629_FreezeAppchainWithFrozenChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.FreezeAppchain(chainID)
	suite.Require().Nil(err)
	err = suite.CheckChainStatus(chainID, governance.GovernanceFrozen)
	suite.Require().Nil(err)
	err = suite.FreezeAppchain(chainID)
	suite.Require().NotNil(err)
}

//tc：应用链处于logouting状态冻结应用链，应用链冻结失败
func (suite Model6) Test0630_FreezeAppchainWithLogoutingChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ChainToLogouting(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.FreezeAppchain(chainID)
	suite.Require().NotNil(err)
}

//tc：应用链处于forbidden状态冻结应用链，应用链冻结失败
func (suite Model6) Test0631_FreezeAppchainWithForbiddenChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ChainToForbidden(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.FreezeAppchain(chainID)
	suite.Require().NotNil(err)
}

//tc：非中继链管理员非应用链管理员激活应用链，应用链激活失败
func (suite Model6) Test0632_ActivateAppchainWithNoAdminIsFail() {
	pk1, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ChainToFrozen(pk1, chainID, address)
	suite.Require().Nil(err)
	pk2, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.ActivateAppchain(pk2, chainID)
	suite.Require().NotNil(err)
}

//tc：应用链管理员激活应用链，应用链激活成功
func (suite Model6) Test0633_ActivateAppchainIsSuccess() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ChainToFrozen(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.ActivateAppchain(pk, chainID)
	suite.Require().Nil(err)
}

//tc：中继链管理员激活应用链，应用链激活成功
func (suite Model6) Test0634_ActivateAppchainWithRelayAdminIsSuccess() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ChainToFrozen(pk, chainID, address)
	suite.Require().Nil(err)
	path, err := repo.Node1Path()
	suite.Require().Nil(err)
	node1pk, err := asym.RestorePrivateKey(path, repo.KeyPassword)
	suite.Require().Nil(err)
	from, err := node1pk.PublicKey().Address()
	suite.Require().Nil(err)
	client := suite.NewClient(node1pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "ActivateAppchain", &rpcx.TransactOpts{
		From:  from.String(),
		Nonce: atomic.AddUint64(&nonce1, 1),
	}, rpcx.String(chainID), rpcx.String("reason"))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	result := &RegisterResult{}
	err = json.Unmarshal(res.Ret, result)
	suite.Require().Nil(err)
	err = suite.VotePass(result.ProposalID)
	suite.Require().Nil(err)
}

//tc：应用链未注册，激活应用链，应用链激活失败
func (suite Model6) Test0635_ActivateAppchainWithNoRegisterChainIsFail() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.ActivateAppchain(pk, suite.GetChainID(pk))
	suite.Require().NotNil(err)
}

//tc：应用链处于registing状态激活应用链，应用链激活失败
func (suite Model6) Test0636_ActivateAppchainWithRegistingChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ChainToRegisting(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.ActivateAppchain(pk, chainID)
	suite.Require().NotNil(err)
}

//tc：应用链处于unavailable状态激活应用链，应用链激活失败
func (suite Model6) Test0637_ActivateAppchainWithUnavailableChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ChainToUnavailable(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.ActivateAppchain(pk, chainID)
	suite.Require().NotNil(err)
}

//tc：应用链处于activating状态激活应用链，应用链激活失败
func (suite Model6) Test0638_ActivateAppchainWithActivatingChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ChainToActivating(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.ActivateAppchain(pk, chainID)
	suite.Require().NotNil(err)
}

//tc：应用链处于freezing状态激活应用链，应用链激活失败
func (suite Model6) Test0639_ActivateAppchainWithFreezingChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ChainToFreezing(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.ActivateAppchain(pk, chainID)
	suite.Require().NotNil(err)
}

//tc：应用链处于logouting状态激活应用链，应用链激活失败
func (suite Model6) Test0640_ActivateAppchainWithLogoutingChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ChainToLogouting(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.ActivateAppchain(pk, chainID)
	suite.Require().NotNil(err)
}

//tc：应用链处于forbidden状态激活应用链，应用链激活失败
func (suite Model6) Test0641_ActivateAppchainWithForbiddenChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ChainToForbidden(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.ActivateAppchain(pk, chainID)
	suite.Require().NotNil(err)
}

//tc：非应用链管理员注销应用链，应用链注销失败
func (suite Model6) Test0642_LogoutAppchainWithNoAdminIsFail() {
	pk1, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk1, chainID, address)
	suite.Require().Nil(err)
	pk2, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.LogoutAppchain(pk2, chainID)
	suite.Require().NotNil(err)
}

//tc：应用链管理员注销应用链，应用链注销成功
func (suite Model6) Test0643_LogoutAppchainIsSuccess() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.RegisterAppchain(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.LogoutAppchain(pk, chainID)
	suite.Require().Nil(err)
}

//tc：应用链未注册，注销应用链，应用链注销失败
func (suite Model6) Test0644_LogoutAppchainWithNoRegisterChainIsFail() {
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	err = suite.LogoutAppchain(pk, suite.GetChainID(pk))
	suite.Require().NotNil(err)
}

//tc：应用链处于registing状态注销应用链，应用链注销失败
func (suite Model6) Test0645_LogoutAppchainWithRegistingChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ChainToRegisting(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.LogoutAppchain(pk, chainID)
	suite.Require().NotNil(err)
}

//tc：应用链处于unavailable状态注销应用链，应用链注销失败
func (suite Model6) Test0646_LogoutAppchainWithUnavailableChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ChainToUnavailable(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.LogoutAppchain(pk, chainID)
	suite.Require().NotNil(err)
}

//tc：应用链处于activating状态注销应用链，应用链注销成功
func (suite Model6) Test0647_LogoutAppchainWithActivatingChainIsSuccess() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ChainToActivating(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.LogoutAppchain(pk, chainID)
	suite.Require().Nil(err)
}

//tc：应用链处于freezing状态注销应用链，应用链激活成功
func (suite Model6) Test0648_LogoutAppchainWithFreezingChainIsSuccess() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ChainToFreezing(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.LogoutAppchain(pk, chainID)
	suite.Require().Nil(err)
}

//tc：应用链处于frozen状态注销应用链，应用链注销成功
func (suite Model6) Test0649_LogoutAppchainWithFrozenChainIsSuccess() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ChainToFrozen(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.LogoutAppchain(pk, chainID)
	suite.Require().Nil(err)
}

//tc：应用链处于logouting状态注销应用链，应用链注销失败
func (suite Model6) Test0650_LogoutAppchainWithLogoutingChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ChainToLogouting(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.LogoutAppchain(pk, chainID)
	suite.Require().NotNil(err)
}

//tc：应用链处于forbidden状态注销应用链，应用链注销失败
func (suite Model6) Test0651_LogoutAppchainWithForbiddenChainIsFail() {
	pk, chainID, address, err := suite.RegisterRule()
	suite.Require().Nil(err)
	err = suite.ChainToForbidden(pk, chainID, address)
	suite.Require().Nil(err)
	err = suite.LogoutAppchain(pk, chainID)
	suite.Require().NotNil(err)
}

func (suite *Snake) FreezeAppchain(chainID string) error {
	path, err := repo.Node1Path()
	if err != nil {
		return err
	}
	node1Key, err := asym.RestorePrivateKey(path, repo.KeyPassword)
	if err != nil {
		return err
	}
	client := suite.NewClient(node1Key)
	from, err := node1Key.PublicKey().Address()
	if err != nil {
		return err
	}
	nonce := atomic.AddUint64(&nonce1, 1)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "FreezeAppchain",
		&rpcx.TransactOpts{
			From:  from.String(),
			Nonce: nonce,
		},
		rpcx.String(chainID), rpcx.String("reason"),
	)
	if err != nil {
		return err
	}
	if res.Status == pb.Receipt_FAILED {
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

func (suite *Snake) UpdateAppchain(pk crypto.PrivateKey, chainID, desc string) error {
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(chainID),
		rpcx.String(desc),
	}

	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "UpdateAppchain", nil, args...)
	if err != nil {
		return err
	}
	if res.Status == pb.Receipt_FAILED {
		return errors.New(string(res.Ret))
	}
	return nil
}

func (suite *Snake) ActivateAppchain(pk crypto.PrivateKey, chainID string) error {
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "ActivateAppchain", nil, rpcx.String(chainID), rpcx.String("reason"))
	if err != nil {
		return err
	}
	if res.Status == pb.Receipt_FAILED {
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

func (suite *Snake) LogoutAppchain(pk crypto.PrivateKey, chainID string) error {
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "LogoutAppchain", nil, rpcx.String(chainID), rpcx.String("reason"))
	if err != nil {
		return err
	}
	if res.Status == pb.Receipt_FAILED {
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

func (suite Snake) ChainToRegisting(pk crypto.PrivateKey, chainID, address string) error {
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(chainID),   //ID
		rpcx.Bytes([]byte("")), //trustRoot
		rpcx.String("0x857133c5C69e6Ce66F7AD46F200B9B3573e77582"), //broker
		rpcx.String("desc"),   //desc
		rpcx.String(address),  //masterRule
		rpcx.String("reason"), //reason
	}
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "RegisterAppchain", nil, args...)
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return err
	}
	err = suite.CheckChainStatus(chainID, governance.GovernanceRegisting)
	if err != nil {
		return err
	}
	return nil
}

func (suite Snake) ChainToUnavailable(pk crypto.PrivateKey, chainID, address string) error {
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(chainID),   //ID
		rpcx.Bytes([]byte("")), //trustRoot
		rpcx.String("0x857133c5C69e6Ce66F7AD46F200B9B3573e77582"), //broker
		rpcx.String("desc"),   //desc
		rpcx.String(address),  //masterRule
		rpcx.String("reason"), //reason
	}
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "RegisterAppchain", nil, args...)
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return err
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
	err = suite.CheckChainStatus(chainID, governance.GovernanceUnavailable)
	if err != nil {
		return err
	}
	return nil
}

func (suite Snake) ChainToActivating(pk crypto.PrivateKey, chainID, address string) error {
	err := suite.RegisterAppchain(pk, chainID, address)
	if err != nil {
		return err
	}
	err = suite.FreezeAppchain(chainID)
	if err != nil {
		return err
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "ActivateAppchain", nil, rpcx.String(chainID), rpcx.String("reason"))
	if err != nil {
		return err
	}
	if res.Status == pb.Receipt_FAILED {
		return errors.New(string(res.Ret))
	}
	err = suite.CheckChainStatus(chainID, governance.GovernanceActivating)
	if err != nil {
		return err
	}
	return nil
}

func (suite Snake) ChainToFreezing(pk crypto.PrivateKey, chainID, address string) error {
	err := suite.RegisterAppchain(pk, chainID, address)
	if err != nil {
		return err
	}
	path, err := repo.Node1Path()
	if err != nil {
		return err
	}
	node1Key, err := asym.RestorePrivateKey(path, repo.KeyPassword)
	if err != nil {
		return err
	}
	client := suite.NewClient(node1Key)
	pubAddress, err := node1Key.PublicKey().Address()
	if err != nil {
		return err
	}
	nonce := atomic.AddUint64(&nonce1, 1)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "FreezeAppchain",
		&rpcx.TransactOpts{
			From:  pubAddress.String(),
			Nonce: nonce,
		},
		rpcx.String(chainID), rpcx.String("reason"),
	)
	if err != nil {
		return err
	}
	if res.Status == pb.Receipt_FAILED {
		return errors.New(string(res.Ret))
	}
	err = suite.CheckChainStatus(chainID, governance.GovernanceFreezing)
	if err != nil {
		return err
	}
	return nil
}

func (suite Snake) ChainToFrozen(pk crypto.PrivateKey, chainID, address string) error {
	err := suite.RegisterAppchain(pk, chainID, address)
	if err != nil {
		return err
	}
	err = suite.FreezeAppchain(chainID)
	if err != nil {
		return err
	}
	err = suite.CheckChainStatus(chainID, governance.GovernanceFrozen)
	if err != nil {
		return err
	}
	return nil
}

func (suite Snake) ChainToLogouting(pk crypto.PrivateKey, chainID, address string) error {
	err := suite.RegisterAppchain(pk, chainID, address)
	if err != nil {
		return err
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "LogoutAppchain", nil, rpcx.String(chainID), rpcx.String("reason"))
	if err != nil {
		return err
	}
	if res.Status == pb.Receipt_FAILED {
		return errors.New(string(res.Ret))
	}
	err = suite.CheckChainStatus(chainID, governance.GovernanceLogouting)
	if err != nil {
		return err
	}
	return nil
}

func (suite Snake) ChainToForbidden(pk crypto.PrivateKey, chainID, address string) error {
	err := suite.RegisterAppchain(pk, chainID, address)
	if err != nil {
		return err
	}
	err = suite.LogoutAppchain(pk, chainID)
	if err != nil {
		return err
	}
	err = suite.CheckChainStatus(chainID, governance.GovernanceForbidden)
	if err != nil {
		return err
	}
	return nil
}
