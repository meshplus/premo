package did_tester

import (
	"encoding/json"
	"fmt"
	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
)

func (suite *Snake) Test001_DidInit() {
	//tc：初始化did_registry,发起调用gosdk绑定的私钥和初始化链上地址不一致
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.to.String()),
	}
	res, err := suite.client.InvokeBVMContract(constant.DIDRegistryContractAddr.Address(), "Init", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)

	//tc：初始化did_registry时，参数格式错误
	args = []*pb.Arg{
		rpcx.String("did:bitxhub:root" + suite.from.String()),
	}
	res, err = suite.client.InvokeBVMContract(constant.DIDRegistryContractAddr.Address(), "Init", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)

	//tc：正确初始化did_registry
	args = []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()),
	}
	res, err = suite.client.InvokeBVMContract(constant.DIDRegistryContractAddr.Address(), "Init", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	//tc：重复初始化did_registry
	args = []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()),
	}
	res, err = suite.client.InvokeBVMContract(constant.DIDRegistryContractAddr.Address(), "Init", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	suite.Require().Contains(string(res.Ret), "call error: init err, already init")
}

//tc：注册did，发起调用gosdk绑定的私钥和初始化链上地址不一致
func (suite *Snake) Test002_DidRegisterWithErrorKey() {
	key := suite.pk
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	suite.client.SetPrivateKey(pk)
	//create doc. must start ipfs service
	//put data
	docAddr, err := suite.client.IPFSPutFromLocal("./testdata/did.json")
	suite.Require().Nil(err)
	fmt.Println(string(docAddr.GetData()))
	//get data
	docData, err := suite.client.IPFSGet("/ipfs/" + string(docAddr.GetData()))
	suite.Require().Nil(err)
	fmt.Println(string(docData.GetData()))

	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain001:" + suite.from.String()), //admin
		rpcx.String(docAddr.String()),
		rpcx.Bytes(docData.GetData()),
		rpcx.Bytes([]byte{1, 2, 3}),
	}
	res, err := suite.client.InvokeBVMContract(constant.DIDRegistryContractAddr.Address(), "Register", nil, args...)
	suite.client.SetPrivateKey(key)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
}

//tc：注册did时，参数格式错误
func (suite *Snake) Test003_DidRegisterWithErrorArgs() {
	//create doc. must start ipfs service
	//put data
	docAddr, err := suite.client.IPFSPutFromLocal("./testdata/did.json")
	suite.Require().Nil(err)
	fmt.Println(string(docAddr.GetData()))
	//get data
	docData, err := suite.client.IPFSGet("/ipfs/" + string(docAddr.GetData()))
	suite.Require().Nil(err)
	fmt.Println(string(docData.GetData()))

	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain001" + suite.from.String()), //admin
		rpcx.String(docAddr.String()),
		rpcx.Bytes(docData.GetData()),
		rpcx.Bytes([]byte{1, 2, 3}),
	}
	res, err := suite.client.InvokeBVMContract(constant.DIDRegistryContractAddr.Address(), "Register", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
}

//tc：注册did时，注册其他链上的did
func (suite *Snake) Test004_DidRegisterWithOtherAppchain() {
	//create doc. must start ipfs service
	//put data
	docAddr, err := suite.client.IPFSPutFromLocal("./testdata/did.json")
	suite.Require().Nil(err)
	fmt.Println(string(docAddr.GetData()))
	//get data
	docData, err := suite.client.IPFSGet("/ipfs/" + string(docAddr.GetData()))
	suite.Require().Nil(err)
	fmt.Println(string(docData.GetData()))

	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain002:" + suite.from.String()), //admin
		rpcx.String(docAddr.String()),
		rpcx.Bytes(docData.GetData()),
		rpcx.Bytes([]byte{1, 2, 3}),
	}
	res, err := suite.client.InvokeBVMContract(constant.DIDRegistryContractAddr.Address(), "Register", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
}

//tc：正确注册did
func (suite Snake) Test005_DidRegister() {
	//create doc. must start ipfs service
	//put data
	docAddr, err := suite.client.IPFSPutFromLocal("./testdata/did.json")
	suite.Require().Nil(err)
	fmt.Println(string(docAddr.GetData()))
	//get data
	docData, err := suite.client.IPFSGet("/ipfs/" + string(docAddr.GetData()))
	suite.Require().Nil(err)
	fmt.Println(string(docData.GetData()))

	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain001:" + suite.from.String()), //admin
		rpcx.String(docAddr.String()),
		rpcx.Bytes(docData.GetData()),
		rpcx.Bytes([]byte{1, 2, 3}),
	}
	res, err := suite.client.InvokeBVMContract(constant.DIDRegistryContractAddr.Address(), "Register", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
}

//tc：重复注册did
func (suite Snake) Test006_DidRegisterRepeat() {
	//create doc. must start ipfs service
	//put data
	docAddr, err := suite.client.IPFSPutFromLocal("./testdata/did.json")
	suite.Require().Nil(err)
	fmt.Println(string(docAddr.GetData()))
	//get data
	docData, err := suite.client.IPFSGet("/ipfs/" + string(docAddr.GetData()))
	suite.Require().Nil(err)
	fmt.Println(string(docData.GetData()))

	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain001:" + suite.from.String()), //admin
		rpcx.String(docAddr.String()),
		rpcx.Bytes(docData.GetData()),
		rpcx.Bytes([]byte{1, 2, 3}),
	}
	res, err := suite.client.InvokeBVMContract(constant.DIDRegistryContractAddr.Address(), "Register", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
}

//tc：正确解析得到did文档
func (suite Snake) Test007_DidResolve() {
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain001:" + suite.from.String()),
	}
	res, err := suite.client.InvokeBVMContract(constant.DIDRegistryContractAddr.Address(), "Resolve", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	mi := &Info{}
	err = Bytes2Struct(res.Ret, mi)
	suite.Require().Nil(err)
	fmt.Println(string(mi.DocHash))
	fmt.Println("Resolve res.Ret: ", mi)
}

//tc：解析不存在的did文档
func (suite Snake) Test008_DidResolveWithDocNotExist() {
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain002:" + suite.from.String()),
	}
	res, err := suite.client.InvokeBVMContract(constant.DIDRegistryContractAddr.Address(), "Resolve", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	suite.Require().Contains(string(res.Ret), "not existed")
}

//tc：更新did，发起调用gosdk绑定的私钥和初始化链上地址不一致
func (suite Snake) Test009_DidUpdateWithErrorKey() {
	key := suite.pk
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	suite.client.SetPrivateKey(pk)
	//create doc. must start ipfs service
	//put data
	docAddr, err := suite.client.IPFSPutFromLocal("./testdata/did_cp.json")
	suite.Require().Nil(err)
	fmt.Println(string(docAddr.GetData()))
	//get data
	docData, err := suite.client.IPFSGet("/ipfs/" + string(docAddr.GetData()))
	suite.Require().Nil(err)
	fmt.Println(string(docData.GetData()))

	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain001:" + suite.from.String()),
		rpcx.String(docAddr.String()),
		rpcx.Bytes(docData.GetData()),
		rpcx.Bytes([]byte{1, 2, 3}),
	}
	res, err := suite.client.InvokeBVMContract(constant.DIDRegistryContractAddr.Address(), "Update", nil, args...)
	suite.client.SetPrivateKey(key)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
}

//tc：更新did时，参数格式错误
func (suite *Snake) Test010_DidUpdateWithErrorArgs() {
	//create doc. must start ipfs service
	//put data
	docAddr, err := suite.client.IPFSPutFromLocal("./testdata/did_cp.json")
	suite.Require().Nil(err)
	fmt.Println(string(docAddr.GetData()))
	//get data
	docData, err := suite.client.IPFSGet("/ipfs/" + string(docAddr.GetData()))
	suite.Require().Nil(err)
	fmt.Println(string(docData.GetData()))

	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain001" + suite.from.String()),
		rpcx.String(docAddr.String()),
		rpcx.Bytes(docData.GetData()),
		rpcx.Bytes([]byte{1, 2, 3}),
	}
	res, err := suite.client.InvokeBVMContract(constant.DIDRegistryContractAddr.Address(), "Update", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
}

//tc：正确更新did
func (suite Snake) Test011_DidUpdate() {
	//create doc. must start ipfs service
	//put data
	docAddr, err := suite.client.IPFSPutFromLocal("./testdata/did_cp.json")
	suite.Require().Nil(err)
	fmt.Println(string(docAddr.GetData()))
	//get data
	docData, err := suite.client.IPFSGet("/ipfs/" + string(docAddr.GetData()))
	suite.Require().Nil(err)
	fmt.Println(string(docData.GetData()))

	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain001:" + suite.from.String()),
		rpcx.String(docAddr.String()),
		rpcx.Bytes(docData.GetData()),
		rpcx.Bytes([]byte{1, 2, 3}),
	}
	res, err := suite.client.InvokeBVMContract(constant.DIDRegistryContractAddr.Address(), "Update", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	//Resolve
	args = []*pb.Arg{
		rpcx.String("did:bitxhub:appchain001:" + suite.from.String()),
	}
	res, err = suite.client.InvokeBVMContract(constant.DIDRegistryContractAddr.Address(), "Resolve", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	mi := &Info{}
	err = Bytes2Struct(res.Ret, mi)
	suite.Require().Nil(err)
	fmt.Println(string(mi.DocHash))
	suite.Require().Equal(docData.GetData(), mi.DocHash)
}

//tc：冻结did，发起调用gosdk绑定的私钥和初始化链上地址不一致
func (suite *Snake) Test012_DidFreezeWithErrorKey() {
	key := suite.pk
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	suite.client.SetPrivateKey(pk)

	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()),
		rpcx.String("did:bitxhub:appchain001:" + suite.from.String()),
		rpcx.Bytes([]byte{1, 2, 3}),
	}
	res, err := suite.client.InvokeBVMContract(constant.DIDRegistryContractAddr.Address(), "Freeze", nil, args...)
	suite.Require().Nil(err)
	suite.client.SetPrivateKey(key)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)

	//check status
	status, err := suite.getDidStatus("did:bitxhub:appchain001:" + suite.from.String())
	suite.Require().Nil(err)
	suite.Require().Equal("Normal", status)
}

//tc：非管理员冻结did
func (suite *Snake) Test013_DidFreezeIsNotAdmin() {
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain001:" + suite.from.String()),
		rpcx.String("did:bitxhub:appchain001:" + suite.from.String()),
		rpcx.Bytes([]byte{1, 2, 3}),
	}
	res, err := suite.client.InvokeBVMContract(constant.DIDRegistryContractAddr.Address(), "Freeze", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)

	//check status
	status, err := suite.getDidStatus("did:bitxhub:appchain001:" + suite.from.String())
	suite.Require().Nil(err)
	suite.Require().Equal("Normal", status)
}

//tc：正确冻结did
func (suite Snake) Test014_DidFreeze() {
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()),
		rpcx.String("did:bitxhub:appchain001:" + suite.from.String()),
		rpcx.Bytes([]byte{1, 2, 3}),
	}
	res, err := suite.client.InvokeBVMContract(constant.DIDRegistryContractAddr.Address(), "Freeze", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	//check status
	status, err := suite.getDidStatus("did:bitxhub:appchain001:" + suite.from.String())
	suite.Require().Nil(err)
	suite.Require().Equal("Frozen", status)
}

//tc：重复冻结did
func (suite Snake) Test015_DidFreezeRepeat() {
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()),
		rpcx.String("did:bitxhub:appchain001:" + suite.from.String()),
		rpcx.Bytes([]byte{1, 2, 3}),
	}
	res, err := suite.client.InvokeBVMContract(constant.DIDRegistryContractAddr.Address(), "Freeze", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)

	//check status
	status, err := suite.getDidStatus("did:bitxhub:appchain001:" + suite.from.String())
	suite.Require().Nil(err)
	suite.Require().Equal("Frozen", status)
}

//tc：解冻did，发起调用gosdk绑定的私钥和初始化链上地址不一致
func (suite *Snake) Test016_DidUnFreezeWithErrorKey()  {
	key := suite.pk
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	suite.client.SetPrivateKey(pk)

	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()),
		rpcx.String("did:bitxhub:appchain001:" + suite.from.String()),
		rpcx.Bytes([]byte{1, 2, 3}),
	}
	res, err := suite.client.InvokeBVMContract(constant.DIDRegistryContractAddr.Address(), "UnFreeze", nil, args...)
	suite.Require().Nil(err)
	suite.client.SetPrivateKey(key)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)

	//check status
	status, err := suite.getDidStatus("did:bitxhub:appchain001:" + suite.from.String())
	suite.Require().Nil(err)
	suite.Require().Equal("Frozen", status)
}

//tc：非管理员解冻did
func (suite *Snake) Test017_DidUnFreezeIsNotAdmin() {
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain001:" + suite.from.String()),
		rpcx.String("did:bitxhub:appchain001:" + suite.from.String()),
		rpcx.Bytes([]byte{1, 2, 3}),
	}
	res, err := suite.client.InvokeBVMContract(constant.DIDRegistryContractAddr.Address(), "UnFreeze", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)

	//check status
	status, err := suite.getDidStatus("did:bitxhub:appchain001:" + suite.from.String())
	suite.Require().Nil(err)
	suite.Require().Equal("Frozen", status)
}

//tc：正确解冻did
func (suite Snake) Test018_DidUnFreeze() {
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()),
		rpcx.String("did:bitxhub:appchain001:" + suite.from.String()),
		rpcx.Bytes([]byte{1, 2, 3}),
	}
	res, err := suite.client.InvokeBVMContract(constant.DIDRegistryContractAddr.Address(), "UnFreeze", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	//check status
	status, err := suite.getDidStatus("did:bitxhub:appchain001:" + suite.from.String())
	suite.Require().Nil(err)
	suite.Require().Equal("Normal", status)
}

//tc：重复解冻did
func (suite Snake) Test019_DidUnFreezeRepeat() {
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()),
		rpcx.String("did:bitxhub:appchain001:" + suite.from.String()),
		rpcx.Bytes([]byte{1, 2, 3}),
	}
	res, err := suite.client.InvokeBVMContract(constant.DIDRegistryContractAddr.Address(), "UnFreeze", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)

	//check status
	status, err := suite.getDidStatus("did:bitxhub:appchain001:" + suite.from.String())
	suite.Require().Nil(err)
	suite.Require().Equal("Normal", status)
}

//tc：tc：删除did，发起调用gosdk绑定的私钥和初始化链上地址不一致
func (suite *Snake) Test020_DidDeleteWithErrorKey() {
	key := suite.pk
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	suite.client.SetPrivateKey(pk)

	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain001:" + suite.from.String()),
		rpcx.Bytes([]byte{1, 2, 3}),
	}
	res, err := suite.client.InvokeBVMContract(constant.DIDRegistryContractAddr.Address(), "Delete", nil, args...)
	suite.client.SetPrivateKey(key)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
}

//tc：删除did，did不存在
func (suite *Snake) Test021_DidDeleteWithNotExist() {
	key := suite.pk
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	suite.client.SetPrivateKey(pk)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)

	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain001:" + address.String()),
		rpcx.Bytes([]byte{1, 2, 3}),
	}
	res, err := suite.client.InvokeBVMContract(constant.DIDRegistryContractAddr.Address(), "Delete", nil, args...)
	suite.client.SetPrivateKey(key)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
}

//tc：正确删除did
func (suite *Snake) Test022_DidDelete() {
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()),
		rpcx.String("did:bitxhub:appchain001:" + suite.from.String()),
		rpcx.Bytes([]byte{1, 2, 3}),
	}
	res, err := suite.client.InvokeBVMContract(constant.DIDRegistryContractAddr.Address(), "Delete", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
}

//tc：重复删除did
func (suite *Snake) Test023_DidDeleteRepeat() {
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain001:" + suite.from.String()),
		rpcx.Bytes([]byte{1, 2, 3}),
	}
	res, err := suite.client.InvokeBVMContract(constant.DIDRegistryContractAddr.Address(), "Delete", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
}

//tc：超管增加管理员，发起调用gosdk绑定的私钥和初始化链上地址不一致
func (suite *Snake) Test024_DidAddAdminWithErrorKey() {
	key := suite.pk
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	suite.client.SetPrivateKey(pk)
	admins, err := suite.getDidAdmins()
	suite.Require().Nil(err)
	num1 := len(admins)
	fmt.Println(admins)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()),
		rpcx.String("did:bitxhub:appchain001:" + suite.from.String()),
	}
	res, err := suite.client.InvokeBVMContract(constant.DIDRegistryContractAddr.Address(), "AddAdmin", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)

	admins, err = suite.getDidAdmins()
	suite.client.SetPrivateKey(key)
	suite.Require().Nil(err)
	num2 := len(admins)
	fmt.Println(admins)
	suite.Require().Equal(num1, num2)
}

//tc：超管增加管理员
func (suite *Snake) Test025_DidAddAdmin() {
	admins, err := suite.getDidAdmins()
	suite.Require().Nil(err)
	num1 := len(admins)
	fmt.Println(admins)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()),
		rpcx.String("did:bitxhub:appchain001:" + suite.from.String()),
	}
	res, err := suite.client.InvokeBVMContract(constant.DIDRegistryContractAddr.Address(), "AddAdmin", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	admins, err = suite.getDidAdmins()
	suite.Require().Nil(err)
	num2 := len(admins)
	fmt.Println(admins)
	suite.Require().Equal(num1+1, num2)
}

//tc：超管重复增加管理员
func (suite *Snake) Test026_DidAddAdminRepeat() {
	admins, err := suite.getDidAdmins()
	suite.Require().Nil(err)
	num1 := len(admins)
	fmt.Println(admins)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()),
		rpcx.String("did:bitxhub:appchain001:" + suite.from.String()),
	}
	res, err := suite.client.InvokeBVMContract(constant.DIDRegistryContractAddr.Address(), "AddAdmin", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)

	admins, err = suite.getDidAdmins()
	suite.Require().Nil(err)
	num2 := len(admins)
	fmt.Println(admins)
	suite.Require().Equal(num1, num2)
}

//tc：非超管增加管理员
func (suite *Snake) Test027_DidAddAdminWithNoAdmin() {
	key := suite.pk
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	suite.client.SetPrivateKey(pk)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)

	admins, err := suite.getDidAdmins()
	suite.Require().Nil(err)
	num1 := len(admins)
	fmt.Println(admins)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + address.String()),
		rpcx.String("did:bitxhub:appchain001:" + suite.from.String()),
	}
	res, err := suite.client.InvokeBVMContract(constant.DIDRegistryContractAddr.Address(), "AddAdmin", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)

	admins, err = suite.getDidAdmins()
	suite.client.SetPrivateKey(key)
	suite.Require().Nil(err)
	num2 := len(admins)
	fmt.Println(admins)
	suite.Require().Equal(num1, num2)
}

//tc：超管增加自己为管理员
func (suite *Snake) Test028_DidAddAdminWithSuperAdminSelf() {
	admins, err := suite.getDidAdmins()
	suite.Require().Nil(err)
	num1 := len(admins)
	fmt.Println(admins)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()),
		rpcx.String("did:bitxhub:root:" + suite.from.String()),
	}
	res, err := suite.client.InvokeBVMContract(constant.DIDRegistryContractAddr.Address(), "AddAdmin", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)

	admins, err = suite.getDidAdmins()
	suite.Require().Nil(err)
	num2 := len(admins)
	fmt.Println(admins)
	suite.Require().Equal(num1, num2)
}

//tc：超管删除管理员，发起调用gosdk绑定的私钥和初始化链上地址不一致
func (suite *Snake) Test029_DidRemoveAdminWithErrorKey() {
	key := suite.pk
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	suite.client.SetPrivateKey(pk)
	admins, err := suite.getDidAdmins()
	suite.Require().Nil(err)
	num1 := len(admins)
	fmt.Println(admins)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()),
		rpcx.String("did:bitxhub:appchain001:" + suite.from.String()),
	}
	res, err := suite.client.InvokeBVMContract(constant.DIDRegistryContractAddr.Address(), "RemoveAdmin", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)

	admins, err = suite.getDidAdmins()
	suite.client.SetPrivateKey(key)
	suite.Require().Nil(err)
	num2 := len(admins)
	fmt.Println(admins)
	suite.Require().Equal(num1, num2)

}

//tc：超管正确删除管理员
func (suite *Snake) Test030_DidRemoveAdmin() {
	admins, err := suite.getDidAdmins()
	suite.Require().Nil(err)
	num1 := len(admins)
	fmt.Println(admins)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()),
		rpcx.String("did:bitxhub:appchain001:" + suite.from.String()),
	}
	res, err := suite.client.InvokeBVMContract(constant.DIDRegistryContractAddr.Address(), "RemoveAdmin", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	admins, err = suite.getDidAdmins()
	suite.Require().Nil(err)
	num2 := len(admins)
	fmt.Println(admins)
	suite.Require().Equal(num1-1, num2)
}

//tc：超管重复删除管理员
func (suite *Snake) Test031_DidRemoveAdminRepeat() {
	admins, err := suite.getDidAdmins()
	suite.Require().Nil(err)
	num1 := len(admins)
	fmt.Println(admins)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()),
		rpcx.String("did:bitxhub:appchain001:" + suite.from.String()),
	}
	res, err := suite.client.InvokeBVMContract(constant.DIDRegistryContractAddr.Address(), "RemoveAdmin", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)

	admins, err = suite.getDidAdmins()
	suite.Require().Nil(err)
	num2 := len(admins)
	fmt.Println(admins)
	suite.Require().Equal(num1, num2)
}

//tc：非超管删除管理员
func (suite *Snake) Test032_DidRemoveAdminWithNoAdmin() {
	key := suite.pk
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	suite.client.SetPrivateKey(pk)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)

	admins, err := suite.getDidAdmins()
	suite.Require().Nil(err)
	num1 := len(admins)
	fmt.Println(admins)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + address.String()),
		rpcx.String("did:bitxhub:appchain001:" + suite.from.String()),
	}
	res, err := suite.client.InvokeBVMContract(constant.DIDRegistryContractAddr.Address(), "RemoveAdmin", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)

	admins, err = suite.getDidAdmins()
	suite.client.SetPrivateKey(key)
	suite.Require().Nil(err)
	num2 := len(admins)
	fmt.Println(admins)
	suite.Require().Equal(num1, num2)
}

//tc：超管删除自己为管理员
func (suite *Snake) Test033_DidRemoveAdminWithSuperAdminSelf() {
	admins, err := suite.getDidAdmins()
	suite.Require().Nil(err)
	num1 := len(admins)
	fmt.Println(admins)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()),
		rpcx.String("did:bitxhub:root:" + suite.from.String()),
	}
	res, err := suite.client.InvokeBVMContract(constant.DIDRegistryContractAddr.Address(), "RemoveAdmin", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)

	admins, err = suite.getDidAdmins()
	suite.Require().Nil(err)
	num2 := len(admins)
	fmt.Println(admins)
	suite.Require().Equal(num1, num2)
}

func (suite *Snake) getDidAdmins() ([]string, error) {
	res, err := suite.client.InvokeBVMContract(constant.DIDRegistryContractAddr.Address(), "GetAdmins", nil)
	if err != nil {
		return nil, err
	}
	var admins []string
	err = json.Unmarshal(res.Ret, &admins)
	if err != nil {
		return nil, err
	}
	return admins, nil
}

func (suite Snake) getDidStatus(Did string) (string, error) {
	args := []*pb.Arg{
		rpcx.String(Did),
	}
	res, err := suite.client.InvokeBVMContract(constant.DIDRegistryContractAddr.Address(), "Resolve", nil, args...)
	if err != nil {
		return "", err
	}
	mi := &Info{}
	err = Bytes2Struct(res.Ret, mi)
	if err != nil {
		return "", err
	}
	return mi.Status, nil
}
