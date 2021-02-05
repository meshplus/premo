package did_tester

import (
	"encoding/hex"
	"encoding/json"

	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-kit/types"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/stretchr/testify/suite"
)

type Snake struct {
	suite.Suite
	//client    rpcx.ChainClient
	client    rpcx.Client
	from      *types.Address
	fromIndex uint64
	pk        crypto.PrivateKey
	toIndex   uint64
	to        *types.Address
}

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

func (suite *Snake) prepare() (crypto.PrivateKey, crypto.PrivateKey, *types.Address, *types.Address) {
	kA, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	kB, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)

	from, err := kA.PublicKey().Address()
	suite.Require().Nil(err)
	to, err := kB.PublicKey().Address()
	suite.Require().Nil(err)

	return kA, kB, from, to
}

func (suite *Snake) RegisterAppchainWithReturn(pk crypto.PrivateKey, chainType string) string {
	pubBytes, err := pk.PublicKey().Bytes()
	suite.Require().Nil(err)

	suite.client.SetPrivateKey(pk)
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
	res, err := suite.client.InvokeBVMContract(constant.AppchainMgrContractAddr.Address(), "Register", nil, args...)
	suite.Require().Nil(err)
	appChain := &rpcx.Appchain{}
	err = json.Unmarshal(res.Ret, appChain)
	suite.Require().Nil(err)
	suite.Require().NotNil(appChain.ID)
	return appChain.ID
}
