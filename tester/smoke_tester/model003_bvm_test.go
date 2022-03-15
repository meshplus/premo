package bxh_tester

import (
	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
)

type Model3 struct {
	*Snake
}

//tc：调用store合约，set （a，b），交易回执显示成功
//tc：调用store合约，get（a），交易回执状态显示成功并且回执数据为b
func (suite *Model3) Test0301_SetAndGetNormalIsSuccess() {
	normalKey := "key_for_normal"
	normalValue := "value_for_normal"
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res1, err := client.InvokeBVMContract(constant.StoreContractAddr.Address(), "Set", nil, pb.String(normalKey), pb.String(normalValue))
	suite.Require().Nil(err)
	suite.Require().True(res1.IsSuccess())
	res2, err := client.InvokeBVMContract(constant.StoreContractAddr.Address(), "Get", nil, pb.String("key_for_normal"))
	suite.Require().Nil(err)
	suite.Require().True(res2.IsSuccess())
	suite.Require().Equal(normalValue, string(res2.Ret))
}
