package did_tester

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/crypto/asym"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
)

type Info struct {
	Method  string // method name
	Owner   string // owner of the method, is a did
	DocAddr string // address where the doc file stored
	DocHash []byte // hash of the doc file
	Status  string // status of method
}

// Bytes2Struct .
func Bytes2Struct(b []byte, s interface{}) error {
	buf := bytes.NewBuffer(b)
	err := gob.NewDecoder(buf).Decode(s)
	if err != nil {
		return fmt.Errorf("gob decode err: %w", err)
	}
	return nil
}

//TODO:add some wrong test

func (suite *Snake) Test001_MethodInit() {
	//tc：初始化method_registry,发起调用gosdk绑定的私钥和初始化链上地址不一致
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.to.String()),
	}
	res, err := suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "Init", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)

	//tc：初始化method_registry时，参数格式错误
	args = []*pb.Arg{
		rpcx.String("did:bitxhub:root" + suite.from.String()),
	}
	res, err = suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "Init", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)

	//tc：正确初始化method_registry
	args = []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()),
	}
	res, err = suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "Init", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	//tc：重复初始化method_registry
	args = []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()),
	}
	res, err = suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "Init", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Contains(string(res.Ret), "call error: init err, already init")
}

//tc：申请区块链Method标识，发起调用gosdk绑定的私钥和初始化链上地址不一致
func (suite *Snake) Test002_MethodApplyWithErrorKey() {
	key := suite.pk
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	suite.client.SetPrivateKey(pk)

	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()),
		rpcx.String("did:bitxhub:appchain001:."),
		rpcx.Bytes([]byte{1, 2, 3}),
	}
	res, err := suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "Apply", nil, args...)
	suite.client.SetPrivateKey(key)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	fmt.Println(string(res.Ret))
}

//tc：申请非法的Method标识
func (suite *Snake) Test003_MethodApplyWithErrorArgs() {
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()),
		rpcx.String("did:bitxhub:appchain001:"),
		rpcx.Bytes([]byte{1, 2, 3}),
	}
	res, err := suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "Apply", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	fmt.Println(string(res.Ret))
}

//tc：正确申请区块链Method标识
func (suite *Snake) Test004_MethodApply() {
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()),
		rpcx.String("did:bitxhub:appchain001:."),
		rpcx.Bytes([]byte{1, 2, 3}),
	}
	res, err := suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "Apply", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	fmt.Println(string(res.Ret))
}

//tc：重复申请区块链Method标识
func (suite *Snake) Test005_MethodApplyRepeat() {
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()),
		rpcx.String("did:bitxhub:appchain001:."),
		rpcx.Bytes([]byte{1, 2, 3}),
	}
	res, err := suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "Apply", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	fmt.Println(string(res.Ret))
}

//tc：管理员申请区块链Method标识审核不通过
func (suite *Snake) Test006_MethodAuditApplyIsFalse() {
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()), // admin
		rpcx.String("did:bitxhub:appchain001:."),
		rpcx.Int32(0),
		rpcx.Bytes([]byte{1, 2, 3}),
	}
	res, err := suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "AuditApply", nil, args...)
	suite.Require().Nil(err)
	fmt.Println("AuditApply res.Ret: ", string(res.Ret))
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
}

//tc：非管理员进行审核
func (suite *Snake) Test007_MethodAuditApplyWithNotAdmin() {
	key := suite.pk
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	suite.client.SetPrivateKey(pk)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + address.String()), // admin is did init
		rpcx.String("did:bitxhub:appchain001:."),
		rpcx.Int32(1),
		rpcx.Bytes([]byte{1, 2, 3}),
	}
	res, err := suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "AuditApply", nil, args...)
	suite.client.SetPrivateKey(key)
	suite.Require().Nil(err)
	fmt.Println("AuditApply res.Ret: ", string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
}

//tc：管理员申请区块链Method标识审核通过
func (suite *Snake) Test008_MethodAuditApplyIsTrue() {
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()), // admin is did init
		rpcx.String("did:bitxhub:appchain001:."),
		rpcx.Int32(1),
		rpcx.Bytes([]byte{1, 2, 3}),
	}
	res, err := suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "AuditApply", nil, args...)
	suite.Require().Nil(err)
	fmt.Println("AuditApply res.Ret: ", string(res.Ret))
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
}

//tc：管理员重复审核请求
func (suite *Snake) Test009_MethodAuditApplyIsTrueRepeat() {
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()), // admin is did init
		rpcx.String("did:bitxhub:appchain001:."),
		rpcx.Int32(1),
		rpcx.Bytes([]byte{1, 2, 3}),
	}
	res, err := suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "AuditApply", nil, args...)
	suite.Require().Nil(err)
	fmt.Println("AuditApply res.Ret: ", string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
}

//tc：管理员审核不存在的请求
func (suite *Snake) Test010_MethodAuditApplyWithNotApply() {
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()), // admin
		rpcx.String("did:bitxhub:appchain002:."),
		rpcx.Int32(0),
		rpcx.Bytes([]byte{1, 2, 3}),
	}
	res, err := suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "AuditApply", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	fmt.Println("AuditApply res.Ret: ", string(res.Ret))
	suite.Require().Contains(string(res.Ret), "audit apply err, auditapply did:bitxhub:appchain002:. not existed")
}

//tc：注册Method，发起调用gosdk绑定的私钥和初始化链上地址不一致
func (suite *Snake) Test011_MethodRegisterWithErrorKey() {
	//create doc. must start ipfs service
	//put data
	docAddr, err := suite.client.IPFSPutFromLocal("./testdata/method.json")
	suite.Require().Nil(err)
	fmt.Println(string(docAddr.GetData()))
	//get data
	docData, err := suite.client.IPFSGet("/ipfs/" + string(docAddr.GetData()))
	suite.Require().Nil(err)
	fmt.Println(string(docData.GetData()))

	key := suite.pk
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	suite.client.SetPrivateKey(pk)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()), //admin
		rpcx.String("did:bitxhub:appchain001:."),
		rpcx.String(docAddr.String()),
		rpcx.Bytes(docData.GetData()),
		rpcx.Bytes([]byte{1, 2, 3}),
	}
	res, err := suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "Register", nil, args...)
	suite.client.SetPrivateKey(key)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	fmt.Println(string(res.Ret))
}

//tc：注册Method，参数格式格式错误
func (suite *Snake) Test012_MethodRegisterWithErrorArgs() {
	//create doc. must start ipfs service
	//put data
	docAddr, err := suite.client.IPFSPutFromLocal("./testdata/method.json")
	suite.Require().Nil(err)
	fmt.Println(string(docAddr.GetData()))
	//get data
	docData, err := suite.client.IPFSGet("/ipfs/" + string(docAddr.GetData()))
	suite.Require().Nil(err)
	fmt.Println(string(docData.GetData()))

	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()), //admin
		rpcx.String("did:bitxhub:appchain001:"),
		rpcx.String(docAddr.String()),
		rpcx.Bytes(docData.GetData()),
		rpcx.Bytes([]byte{1, 2, 3}),
	}
	res, err := suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "Register", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
}

//tc：正确注册Method
func (suite *Snake) Test013_MethodRegister() {
	//create doc. must start ipfs service
	//put data
	docAddr, err := suite.client.IPFSPutFromLocal("./testdata/method.json")
	suite.Require().Nil(err)
	fmt.Println(string(docAddr.GetData()))
	//get data
	docData, err := suite.client.IPFSGet("/ipfs/" + string(docAddr.GetData()))
	suite.Require().Nil(err)
	fmt.Println(string(docData.GetData()))

	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()), //admin
		rpcx.String("did:bitxhub:appchain001:."),
		rpcx.String(docAddr.String()),
		rpcx.Bytes(docData.GetData()),
		rpcx.Bytes([]byte{1, 2, 3}),
	}
	res, err := suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "Register", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
}

//tc：重复注册Method
func (suite *Snake) Test014_MethodRegisterRepeat() {
	//create doc. must start ipfs service
	//put data
	docAddr, err := suite.client.IPFSPutFromLocal("./testdata/method.json")
	suite.Require().Nil(err)
	fmt.Println(string(docAddr.GetData()))
	//get data
	docData, err := suite.client.IPFSGet("/ipfs/" + string(docAddr.GetData()))
	suite.Require().Nil(err)
	fmt.Println(string(docData.GetData()))

	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()), //admin
		rpcx.String("did:bitxhub:appchain001:."),
		rpcx.String(docAddr.String()),
		rpcx.Bytes(docData.GetData()),
		rpcx.Bytes([]byte{1, 2, 3}),
	}
	res, err := suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "Register", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
}

//tc：正确解析得到Method文档
func (suite *Snake) Test015_MethodResolve() {
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:appchain001:."),
	}
	res, err := suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "Resolve", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	mi := &Info{}
	err = Bytes2Struct(res.Ret, mi)
	suite.Require().Nil(err)
	fmt.Println("Resolve res.Ret: ", mi)
	fmt.Println(string(mi.DocHash))
}

//tc：正确更新Method
func (suite *Snake) Test016_MethodUpdate() {
	//create doc. must start ipfs service
	//put data
	docAddr, err := suite.client.IPFSPutFromLocal("./testdata/method_cp.json")
	suite.Require().Nil(err)
	fmt.Println(string(docAddr.GetData()))
	//get data
	docData, err := suite.client.IPFSGet("/ipfs/" + string(docAddr.GetData()))
	suite.Require().Nil(err)
	fmt.Println(string(docData.GetData()))

	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()), //admin
		rpcx.String("did:bitxhub:appchain001:."),
		rpcx.String(docAddr.String()),
		rpcx.Bytes(docData.GetData()),
		rpcx.Bytes([]byte{1, 2, 3}),
	}
	res, err := suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "Update", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	fmt.Println(string(res.Ret))

	//Resolve
	args = []*pb.Arg{
		rpcx.String("did:bitxhub:appchain001:."),
	}
	res, err = suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "Resolve", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	mi := &Info{}
	err = Bytes2Struct(res.Ret, mi)
	suite.Require().Nil(err)
	suite.Require().Equal(docData.GetData(), mi.DocHash)
}

//tc：更新Method，发起调用gosdk绑定的私钥和初始化链上地址不一致
func (suite *Snake) Test017_MethodUpdateWithErrorKey() {
	//create doc. must start ipfs service
	//put data
	docAddr, err := suite.client.IPFSPutFromLocal("./testdata/method_cp.json")
	suite.Require().Nil(err)
	fmt.Println(string(docAddr.GetData()))
	//get data
	docData, err := suite.client.IPFSGet("/ipfs/" + string(docAddr.GetData()))
	suite.Require().Nil(err)
	fmt.Println(string(docData.GetData()))

	key := suite.pk
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	suite.client.SetPrivateKey(pk)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()), //admin
		rpcx.String("did:bitxhub:appchain001:."),
		rpcx.String(docAddr.String()),
		rpcx.Bytes(docData.GetData()),
		rpcx.Bytes([]byte{1, 2, 3}),
	}
	res, err := suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "Update", nil, args...)
	suite.client.SetPrivateKey(key)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	fmt.Println(string(res.Ret))
}

//tc：冻结Method，发起调用gosdk绑定的私钥和初始化链上地址不一致
func (suite *Snake) Test018_MethodFreezeWithErrorKey() {
	key := suite.pk
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	suite.client.SetPrivateKey(pk)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()), //admin
		rpcx.String("did:bitxhub:appchain001:."),
		rpcx.Bytes([]byte{1, 2, 3}),
	}
	res, err := suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "Freeze", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)

	//check status
	status, err := suite.getMethodStatus("did:bitxhub:appchain001:.")
	suite.client.SetPrivateKey(key)
	suite.Require().Nil(err)
	suite.Require().Equal("Normal", status)
}

//tc：非管理员冻结Method
func (suite *Snake) Test019_MethodFreezeWithNotAdmin() {
	key := suite.pk
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	suite.client.SetPrivateKey(pk)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + address.String()),
		rpcx.String("did:bitxhub:appchain001:."),
		rpcx.Bytes([]byte{1, 2, 3}),
	}
	res, err := suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "Freeze", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)

	//check status
	status, err := suite.getMethodStatus("did:bitxhub:appchain001:.")
	suite.client.SetPrivateKey(key)
	suite.Require().Nil(err)
	suite.Require().Equal("Normal", status)
}

//tc：正确冻结Method
func (suite *Snake) Test020_MethodFreeze() {
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()), //admin
		rpcx.String("did:bitxhub:appchain001:."),
		rpcx.Bytes([]byte{1, 2, 3}),
	}
	res, err := suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "Freeze", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	//check status
	status, err := suite.getMethodStatus("did:bitxhub:appchain001:.")
	suite.Require().Nil(err)
	suite.Require().Equal("Frozen", status)
}

//tc：重复冻结Method
func (suite *Snake) Test021_MethodFreezeRepeat() {
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()), //admin
		rpcx.String("did:bitxhub:appchain001:."),
		rpcx.Bytes([]byte{1, 2, 3}),
	}
	res, err := suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "Freeze", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)

	//check status
	status, err := suite.getMethodStatus("did:bitxhub:appchain001:.")
	suite.Require().Nil(err)
	suite.Require().Equal("Frozen", status)
}

//tc：解冻Method，发起调用gosdk绑定的私钥和初始化链上地址不一致
func (suite *Snake) Test022_MethodUnFreeWithErrorKey() {
	key := suite.pk
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	suite.client.SetPrivateKey(pk)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()), //admin
		rpcx.String("did:bitxhub:appchain001:."),
		rpcx.Bytes([]byte{1, 2, 3}),
	}
	res, err := suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "UnFreeze", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)

	//check status
	status, err := suite.getMethodStatus("did:bitxhub:appchain001:.")
	suite.client.SetPrivateKey(key)
	suite.Require().Nil(err)
	suite.Require().Equal("Frozen", status)
}

//tc：非管理员解冻Method
func (suite *Snake) Test023_MethodUnFreeWithNotAdmin() {
	key := suite.pk
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)
	suite.client.SetPrivateKey(pk)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + address.String()), //admin
		rpcx.String("did:bitxhub:appchain001:."),
		rpcx.Bytes([]byte{1, 2, 3}),
	}
	res, err := suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "UnFreeze", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)

	//check status
	status, err := suite.getMethodStatus("did:bitxhub:appchain001:.")
	suite.client.SetPrivateKey(key)
	suite.Require().Nil(err)
	suite.Require().Equal("Frozen", status)
}

//tc：正确解冻Method
func (suite *Snake) Test024_MethodUnFreeze() {
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()), //admin
		rpcx.String("did:bitxhub:appchain001:."),
		rpcx.Bytes([]byte{1, 2, 3}),
	}
	res, err := suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "UnFreeze", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	//check status
	status, err := suite.getMethodStatus("did:bitxhub:appchain001:.")
	suite.Require().Nil(err)
	suite.Require().Equal("Normal", status)
}

//tc：重复解冻Method
func (suite *Snake) Test025_MethodUnFreezeRepeat() {
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()), //admin
		rpcx.String("did:bitxhub:appchain001:."),
		rpcx.Bytes([]byte{1, 2, 3}),
	}
	res, err := suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "UnFreeze", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)

	//check status
	status, err := suite.getMethodStatus("did:bitxhub:appchain001:.")
	suite.Require().Nil(err)
	suite.Require().Equal("Normal", status)
}

//tc：删除Method，发起调用gosdk绑定的私钥和初始化链上地址不一致
func (suite *Snake) Test026_MethodDeleteMethodWithErrorKey() {
	key := suite.pk
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	suite.client.SetPrivateKey(pk)

	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()), //admin
		rpcx.String("did:bitxhub:appchain001:."),
		rpcx.Bytes([]byte{1, 2, 3}),
	}
	res, err := suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "Delete", nil, args...)
	suite.client.SetPrivateKey(key)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
}

//tc：正确删除Method
func (suite *Snake) Test027_MethodDeleteMethod() {
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()), //admin
		rpcx.String("did:bitxhub:appchain001:."),
		rpcx.Bytes([]byte{1, 2, 3}),
	}
	res, err := suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "Delete", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	//resolve
	args = []*pb.Arg{
		rpcx.String("did:bitxhub:appchain001:."),
	}
	res, err = suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "Resolve", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	mi := &Info{}
	err = Bytes2Struct(res.Ret, mi)
	suite.Require().Nil(err)
	fmt.Println("Resolve res.Ret: ", mi)
	fmt.Println(string(mi.DocHash))
}
//tc：重复删除Method
func (suite *Snake) Test028_MethodDeleteMethodRepeat() {
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()), //admin
		rpcx.String("did:bitxhub:appchain001:."),
		rpcx.Bytes([]byte{1, 2, 3}),
	}
	res, err := suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "Delete", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
}

//tc：超管增加管理员，发起调用gosdk绑定的私钥和初始化链上地址不一致
func (suite *Snake) Test029_MethodAddAdminWithErrorKey() {
	key := suite.pk
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	suite.client.SetPrivateKey(pk)
	admins, err := suite.getMethodAdmins()
	suite.Require().Nil(err)
	num1 := len(admins)
	fmt.Println(admins)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()),
		rpcx.String("did:bitxhub:appchain001:" + suite.from.String()),
	}
	res, err := suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "AddAdmin", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)

	admins, err = suite.getMethodAdmins()
	suite.client.SetPrivateKey(key)
	suite.Require().Nil(err)
	num2 := len(admins)
	fmt.Println(admins)
	suite.Require().Equal(num1, num2)
}

//tc：超管增加管理员
func (suite *Snake) Test030_MethodAddAdmin() {
	admins, err := suite.getMethodAdmins()
	suite.Require().Nil(err)
	num1 := len(admins)
	fmt.Println(admins)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()),
		rpcx.String("did:bitxhub:appchain001:" + suite.from.String()),
	}
	res, err := suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "AddAdmin", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	admins, err = suite.getMethodAdmins()
	suite.Require().Nil(err)
	num2 := len(admins)
	fmt.Println(admins)
	suite.Require().Equal(num1+1, num2)
}

//tc：超管重复增加管理员
func (suite *Snake) Test031_MethodAddAdminRepeat() {
	admins, err := suite.getMethodAdmins()
	suite.Require().Nil(err)
	num1 := len(admins)
	fmt.Println(admins)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()),
		rpcx.String("did:bitxhub:appchain001:" + suite.from.String()),
	}
	res, err := suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "AddAdmin", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)

	admins, err = suite.getMethodAdmins()
	suite.Require().Nil(err)
	num2 := len(admins)
	fmt.Println(admins)
	suite.Require().Equal(num1, num2)
}

//tc：非超管增加管理员
func (suite *Snake) Test032_MethodAddAdminWithNoAdmin() {
	key := suite.pk
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	suite.client.SetPrivateKey(pk)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)

	admins, err := suite.getMethodAdmins()
	suite.Require().Nil(err)
	num1 := len(admins)
	fmt.Println(admins)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + address.String()),
		rpcx.String("did:bitxhub:appchain001:" + suite.from.String()),
	}
	res, err := suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "AddAdmin", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)

	admins, err = suite.getMethodAdmins()
	suite.client.SetPrivateKey(key)
	suite.Require().Nil(err)
	num2 := len(admins)
	fmt.Println(admins)
	suite.Require().Equal(num1, num2)
}

//tc：超管增加自己为管理员
func (suite *Snake) Test033_MethodAddAdminWithSuperAdminSelf()  {
	admins, err := suite.getMethodAdmins()
	suite.Require().Nil(err)
	num1 := len(admins)
	fmt.Println(admins)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()),
		rpcx.String("did:bitxhub:root:" + suite.from.String()),
	}
	res, err := suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "AddAdmin", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)

	admins, err = suite.getMethodAdmins()
	suite.Require().Nil(err)
	num2 := len(admins)
	fmt.Println(admins)
	suite.Require().Equal(num1, num2)
}

//tc：超管删除管理员，发起调用gosdk绑定的私钥和初始化链上地址不一致
func (suite *Snake) Test034_MethodRemoveAdminWithErrorKey() {
	key := suite.pk
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	suite.client.SetPrivateKey(pk)
	admins, err := suite.getMethodAdmins()
	suite.Require().Nil(err)
	num1 := len(admins)
	fmt.Println(admins)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()),
		rpcx.String("did:bitxhub:appchain001:" + suite.from.String()),
	}
	res, err := suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "RemoveAdmin", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)

	admins, err = suite.getMethodAdmins()
	suite.client.SetPrivateKey(key)
	suite.Require().Nil(err)
	num2 := len(admins)
	fmt.Println(admins)
	suite.Require().Equal(num1, num2)

}

//tc：超管正确删除管理员
func (suite *Snake) Test035_MethodRemoveAdmin() {
	admins, err := suite.getMethodAdmins()
	suite.Require().Nil(err)
	num1 := len(admins)
	fmt.Println(admins)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()),
		rpcx.String("did:bitxhub:appchain001:" + suite.from.String()),
	}
	res, err := suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "RemoveAdmin", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)

	admins, err = suite.getMethodAdmins()
	suite.Require().Nil(err)
	num2 := len(admins)
	fmt.Println(admins)
	suite.Require().Equal(num1-1, num2)
}

//tc：超管重复删除管理员
func (suite *Snake) Test036_MethodRemoveAdminRepeat() {
	admins, err := suite.getMethodAdmins()
	suite.Require().Nil(err)
	num1 := len(admins)
	fmt.Println(admins)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()),
		rpcx.String("did:bitxhub:appchain001:" + suite.from.String()),
	}
	res, err := suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "RemoveAdmin", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)

	admins, err = suite.getMethodAdmins()
	suite.Require().Nil(err)
	num2 := len(admins)
	fmt.Println(admins)
	suite.Require().Equal(num1, num2)
}

//tc：非超管删除管理员
func (suite *Snake) Test037_MethodRemoveAdminWithNoAdmin() {
	key := suite.pk
	pk, err := asym.GenerateKeyPair(crypto.Secp256k1)
	suite.Require().Nil(err)
	suite.client.SetPrivateKey(pk)
	address, err := pk.PublicKey().Address()
	suite.Require().Nil(err)

	admins, err := suite.getMethodAdmins()
	suite.Require().Nil(err)
	num1 := len(admins)
	fmt.Println(admins)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + address.String()),
		rpcx.String("did:bitxhub:appchain001:" + suite.from.String()),
	}
	res, err := suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "RemoveAdmin", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)

	admins, err = suite.getMethodAdmins()
	suite.client.SetPrivateKey(key)
	suite.Require().Nil(err)
	num2 := len(admins)
	fmt.Println(admins)
	suite.Require().Equal(num1, num2)
}

//tc：超管删除自己为管理员
func (suite *Snake) Test038_MethodRemoveAdminWithSuperAdminSelf()  {
	admins, err := suite.getMethodAdmins()
	suite.Require().Nil(err)
	num1 := len(admins)
	fmt.Println(admins)
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()),
		rpcx.String("did:bitxhub:root:" + suite.from.String()),
	}
	res, err := suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "RemoveAdmin", nil, args...)
	suite.Require().Nil(err)
	fmt.Println(string(res.Ret))
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)

	admins, err = suite.getMethodAdmins()
	suite.Require().Nil(err)
	num2 := len(admins)
	fmt.Println(admins)
	suite.Require().Equal(num1, num2)
}

func (suite *Snake) Test039_MethodHasAdmin() {
	//1 is admin
	args := []*pb.Arg{
		rpcx.String("did:bitxhub:root:" + suite.from.String()),
	}
	res, err := suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "HasAdmin", nil, args...)
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_SUCCESS, res.Status)
	fmt.Println(string(res.Ret))
}

func (suite *Snake) getMethodAdmins() ([]string, error) {
	res, err := suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "GetAdmins", nil)
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

func (suite *Snake) getMethodStatus(method string) (string, error) {
	args := []*pb.Arg{
		rpcx.String(method),
	}
	res, err := suite.client.InvokeBVMContract(constant.MethodRegistryContractAddr.Address(), "Resolve", nil, args...)
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
