package bxh_tester

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"sync/atomic"
	"unsafe"

	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/meshplus/premo/internal/repo"
)

const Year = 365 * 24 * 60 * 60

type PriceLevel struct {
	Price1Letter uint64 `json:"price1Letter"`
	Price2Letter uint64 `json:"price2Letter"`
	Price3Letter uint64 `json:"price3Letter"`
	Price4Letter uint64 `json:"price4Letter"`
	Price5Letter uint64 `json:"price5Letter"`
}

type ServDomain struct {
	Name       string `json:"name"`
	Level      int    `json:"level"`
	Status     int    `json:"status"`
	ParentName string `json:"parent_name"`
}

type Model18 struct {
	*Snake
}

func (suite *Model18) SetupTest() {
	suite.T().Parallel()
}

//tc：根据正确的域名注册域名，域名注册成功
func (suite *Model18) Test1801_RegisterDomainWithGoodDomainIsSuccess() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, Year, constant.ServiceResolverContractAddr.String())
	suite.Require().Nil(err)
}

//tc：根据已经注册的域名注册域名，域名注册失败
func (suite *Model18) Test1802_RegisterDomainWithSameDomainIsFail() {
	pk1, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk1, domain, Year, constant.ServiceResolverContractAddr.String())
	suite.Require().Nil(err)

	pk2, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterDomain(pk2, domain, Year, constant.ServiceResolverContractAddr.String())
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "call error: 2140000:The domain id registered or in GRACEPERIOD")
}

//tc：根据非法的域名注册域名，域名注册失败
//tc：根据空的域名注册域名，域名注册失败
func (suite *Model18) Test1804_RegisterDomainWithEmptyDomainIsFail() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RegisterDomain(pk, "", Year, constant.ServiceResolverContractAddr.String())
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "call error: 2140000:The domain id can not be an empty string")
}

//tc：根据正确的有效期注册域名，域名注册成功
func (suite *Model18) Test1804_RegisterDomainWithGoodDurationIsSuccess() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, Year, constant.ServiceResolverContractAddr.String())
	suite.Require().Nil(err)
}

//tc：根据过大的有效期注册域名（账户金额不足），域名注册失败
func (suite *Model18) Test1805_RegisterDomainWithErrorDurationIsFail() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, 1000*Year, constant.ServiceResolverContractAddr.String())
	suite.Require().NotNil(err)
	suite.Require().Equal("call error: 2140000:Not enough Bitxhub Token provided", err.Error())
}

//tc：根据空的有效期注册域名，域名注册失败
func (suite *Model18) Test1806_RegisterDomainWithEmptyDurationIsFail() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ServiceRegistryContractAddr.Address(), "Register", nil, rpcx.String(domain), rpcx.String(constant.ServiceResolverContractAddr.String()))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	suite.Require().Equal("reflect: Call with too few input arguments", string(res.Ret))
}

//tc：根据正确的解析器合约地址注册域名，域名注册成功
func (suite *Model18) Test1807_RegisterDomainWithGoodResolverIsSuccess() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, Year, constant.ServiceResolverContractAddr.String())
	suite.Require().Nil(err)
}

//tc：根据不存在的解析器合约地址注册域名，域名注册失败
func (suite *Model18) Test1808_RegisterDomainWithNoExistDomainIsFail() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, Year, constant.AppchainMgrContractAddr.String())
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "call error: 2140000:The resolver is not in the list")
}

//tc：根据空的解析器合约地址注册域名，域名注册失败
func (suite *Model18) Test1809_RegisterDomainWithEmptyResolverIsFail() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, Year, "")
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "call error: 2140000:The resolver is not in the list")
}

//tc：根据正确的域名续费域名，域名续费成功
func (suite *Model18) Test1810_RenewDomainWithGoodDomainIsSuccess() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, Year, constant.ServiceResolverContractAddr.String())
	suite.Require().Nil(err)

	err = suite.RenewDomain(pk, Domain(domain), Year)
	suite.Require().Nil(err)
}

//tc：根据不存在的域名续费域名，域名续费失败
func (suite *Model18) Test1811_RenewDomainWithNoExistDomainIsFail() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RenewDomain(pk, Domain(domain), Year)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "call error: 2140000:The domain must register first")
}

//tc：根据空的域名续费域名，域名续费失败
func (suite *Model18) Test1812_RenewDomainWithEmptyDomainIsFail() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.RenewDomain(pk, "", Year)
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "call error: 2140000:The domain id can not be an empty string")
}

//tc：根据正确的有效期续费域名，域名续费成功
func (suite *Model18) Test1813_RenewDomainWithGoodDurationIsSuccess() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, Year, constant.ServiceResolverContractAddr.String())
	suite.Require().Nil(err)

	err = suite.RenewDomain(pk, Domain(domain), Year)
	suite.Require().Nil(err)
}

//tc：根据过大的有效期续费域名，域名续费失败
func (suite *Model18) Test1814_RenewDomainWithErrorDurationIsFail() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, Year, constant.ServiceResolverContractAddr.String())
	suite.Require().Nil(err)

	err = suite.RenewDomain(pk, Domain(domain), 1000*Year)
	suite.Require().NotNil(err)
	suite.Require().Equal("call error: 2140000:Not enough Bitxhub Token provided", err.Error())
}

//tc：根据空的有效期续费域名，域名续费失败
func (suite *Model18) Test1814_RenewDomainWithEmptyDurationIsFail() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, Year, constant.ServiceResolverContractAddr.String())
	suite.Require().Nil(err)

	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ServiceRegistryContractAddr.Address(), "Renew", nil, rpcx.String(domain))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	suite.Require().Equal("reflect: Call with too few input arguments", string(res.Ret))
}

//tc：根据正确的父域名分配二级域名，域名分配成功
func (suite *Model18) Test1815_AllocateSubDomainWithGoodDomainIsSuccess() {
	pk, address, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, Year, constant.ServiceResolverContractAddr.String())
	suite.Require().Nil(err)

	sonDomain := randomDomain(3)
	err = suite.AllocateSubDomain(pk, Domain(domain), sonDomain, address.String(), constant.ServiceResolverContractAddr.String(), "")
	suite.Require().Nil(err)
}

//tc：根据不存在的父域名分配二级域名，域名分配失败
func (suite *Model18) Test1816_AllocateSubDomainWithNoExistDomainIsFail() {
	pk, address, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	sonDomain := randomDomain(3)
	err = suite.AllocateSubDomain(pk, Domain(domain), sonDomain, address.String(), constant.ServiceResolverContractAddr.String(), "")
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "call error: 2140000:The domain must register first")
}

//tc：根据空的父域名分配二级域名，域名分配失败
func (suite *Model18) Test1817_AllocateSubDomainWithEmptyDomainIsFail() {
	pk, address, err := repo.KeyPriv()
	suite.Require().Nil(err)
	sonDomain := randomDomain(3)
	err = suite.AllocateSubDomain(pk, "", sonDomain, address.String(), constant.ServiceResolverContractAddr.String(), "")
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "call error: 2140000:The parentDomain name can not be an empty string")
}

//tc：根据正确的子域名分配二级域名，域名分配成功
func (suite *Model18) Test1818_AllocateSubDomainWithGoodSonDomainIsSuccess() {
	pk, address, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, Year, constant.ServiceResolverContractAddr.String())
	suite.Require().Nil(err)

	sonDomain := randomDomain(3)
	err = suite.AllocateSubDomain(pk, Domain(domain), sonDomain, address.String(), constant.ServiceResolverContractAddr.String(), "")
	suite.Require().Nil(err)
}

//tc：根据非法的子域名分配二级域名，域名分配失败
//tc：根据空的子域名分配二级域名，域名分配失败
func (suite *Model18) Test1819_AllocateSubDomainWithEmptySomDomainIsFail() {
	pk, address, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, Year, constant.ServiceResolverContractAddr.String())
	suite.Require().Nil(err)

	err = suite.AllocateSubDomain(pk, Domain(domain), "", address.String(), constant.ServiceResolverContractAddr.String(), "")
	suite.Require().NotNil(err)
	suite.Require().Equal("call error: 2140000:The sonDomain name can not be an empty string", err.Error())
}

//tc：根据正确的所有者分配二级域名，域名分配成功
func (suite *Model18) Test1820_AllocateSubDomainWithGoodOwnerIsSuccess() {
	pk, address, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, Year, constant.ServiceResolverContractAddr.String())
	suite.Require().Nil(err)

	sonDomain := randomDomain(3)
	err = suite.AllocateSubDomain(pk, Domain(domain), sonDomain, address.String(), constant.ServiceResolverContractAddr.String(), "")
	suite.Require().Nil(err)
}

//tc：根据非法的所有者分配二级域名，域名分配失败
func (suite *Model18) Test1820_AllocateSubDomainWithErrorOwnerIsFail() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, Year, constant.ServiceResolverContractAddr.String())
	suite.Require().Nil(err)

	sonDomain := randomDomain(3)
	err = suite.AllocateSubDomain(pk, Domain(domain), sonDomain, "0x111", constant.ServiceResolverContractAddr.String(), "")
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "call error: 2140000:The address is not valid")
}

//tc：根据空的所有者分配二级域名，域名分配失败
func (suite *Model18) Test1821_AllocateSubDomainWithEmptyOwnerIsFail() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, Year, constant.ServiceResolverContractAddr.String())
	suite.Require().Nil(err)

	sonDomain := randomDomain(3)
	err = suite.AllocateSubDomain(pk, Domain(domain), sonDomain, "", constant.ServiceResolverContractAddr.String(), "")
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "call error: 2140000:The address is not valid")
}

//tc：根据正确的解析器合约地址分配二级域名，域名分配成功
func (suite *Model18) Test1822_AllocateSubDomainWithGoodResolverIsSuccess() {
	pk, address, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, Year, constant.ServiceResolverContractAddr.String())
	suite.Require().Nil(err)

	sonDomain := randomDomain(3)
	err = suite.AllocateSubDomain(pk, Domain(domain), sonDomain, address.String(), constant.ServiceResolverContractAddr.String(), "")
	suite.Require().Nil(err)
}

//tc：根据不存在的解析器合约地址分配二级域名，域名分配失败
func (suite *Model18) Test1823_AllocateSubDomainWithNoExistResolverIsFail() {
	pk, address, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, Year, constant.ServiceResolverContractAddr.String())
	suite.Require().Nil(err)

	sonDomain := randomDomain(3)
	err = suite.AllocateSubDomain(pk, Domain(domain), sonDomain, address.String(), constant.AppchainMgrContractAddr.String(), "")
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "call error: 2140000:The resolver is not in the list")
}

//tc：根据空的解析器合约地址分配二级域名，域名分配失败
func (suite *Model18) Test1824_AllocateSubDomainWithEmptyResolverIsFail() {
	pk, address, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, Year, constant.ServiceResolverContractAddr.String())
	suite.Require().Nil(err)

	sonDomain := randomDomain(3)
	err = suite.AllocateSubDomain(pk, Domain(domain), sonDomain, address.String(), "", "")
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "call error: 2140000:The resolver is not in the list")
}

//tc：非域名所有者分配二级域名，域名分配失败
func (suite *Model18) Test1825_AllocateSubDomainWithErrorCallerIsFail() {
	pk1, address, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk1, domain, Year, constant.ServiceResolverContractAddr.String())
	suite.Require().Nil(err)

	pk2, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	sonDomain := randomDomain(3)
	err = suite.AllocateSubDomain(pk2, Domain(domain), sonDomain, address.String(), constant.ServiceResolverContractAddr.String(), "")
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "call error: 2140000:The domain name does not belong to you")
}

//tc：根据正确的域名删除域名，域名删除成功
func (suite *Model18) Test1826_DeleteSecondDomainWithGoodDomainIsSuccess() {
	pk, address, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, Year, constant.ServiceResolverContractAddr.String())
	suite.Require().Nil(err)

	sonDomain := randomDomain(3)
	err = suite.AllocateSubDomain(pk, Domain(domain), sonDomain, address.String(), constant.ServiceResolverContractAddr.String(), "")
	suite.Require().Nil(err)

	err = suite.DeleteSecondDomain(pk, SonDomain(domain, sonDomain))
	suite.Require().Nil(err)
}

//tc：根据不存在的域名删除域名，域名删除失败
func (suite *Model18) Test1827_DeleteSecondDomainWithNoExistDomainIsFail() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, Year, constant.ServiceResolverContractAddr.String())
	suite.Require().Nil(err)

	sonDomain := randomDomain(3)

	err = suite.DeleteSecondDomain(pk, SonDomain(domain, sonDomain))
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "call error: 2140000:The domain must be Allocate first")
}

//tc：根据空的域名删除域名，域名删除失败
func (suite *Model18) Test1828_DeleteSecondDomainWithEmptyDomainIsFail() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, Year, constant.ServiceResolverContractAddr.String())
	suite.Require().Nil(err)

	err = suite.DeleteSecondDomain(pk, SonDomain(domain, ""))
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "call error: 2140000:The domain must be Allocate first")
}

//tc：根据正确的价格更新注册价格，注册价格更新成功
func (suite *Model18) Test1829_SetPriceLevelWithGoodPriceIsSuccess() {
	err := suite.SetPriceLevel(1, 2, 3, 4, 5)
	suite.Require().Nil(err)
}

//tc：根据空的价格更新注册价格，注册价格更新失败
func (suite *Model18) Test1829_SetPriceLevelWithEmptyPriceIsFail() {
	pk, address, err := repo.Node1Priv()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ServiceRegistryContractAddr.Address(), "SetPriceLevel", &rpcx.TransactOpts{
		From:  address.String(),
		Nonce: atomic.AddUint64(&nonce1, 1),
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	suite.Require().Equal(string(res.Ret), "reflect: Call with too few input arguments")
}

//tc：非中继链管理员更新注册价格，注册价格更新失败
func (suite *Model18) Test1829_SetPriceLevelWithNoAdminIsFail() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ServiceRegistryContractAddr.Address(), "SetPriceLevel", nil, rpcx.Uint64(1), rpcx.Uint64(2), rpcx.Uint64(3), rpcx.Uint64(4), rpcx.Uint64(5))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	suite.Require().Contains(string(res.Ret), "call error: 2010000:you have no permission")
}

//tc：获取注册价格，价格获取成功
func (suite *Model18) Test1830_GetPriceLevelIsSuccess() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	price, err := suite.GetPriceLevel(pk)
	suite.Require().Nil(err)
	suite.Require().Equal(uint64(1), price.Price1Letter)
	suite.Require().Equal(uint64(2), price.Price2Letter)
	suite.Require().Equal(uint64(3), price.Price3Letter)
	suite.Require().Equal(uint64(4), price.Price4Letter)
	suite.Require().Equal(uint64(5), price.Price5Letter)
}

//tc：根据正确的价格设置币价，币价设置成功
func (suite *Model18) Test1831_SetTokenPriceWithGoodPriceIsSuccess() {
	err := suite.SetTokenPrice(1)
	suite.Require().Nil(err)
}

//tc：根据空的价格设置币价，币价设置失败
func (suite *Model18) Test1832_SetTokenPriceWithEmptyPriceIsFail() {
	pk, address, err := repo.Node1Priv()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ServiceRegistryContractAddr.Address(), "SetTokenPrice", &rpcx.TransactOpts{
		From:  address.String(),
		Nonce: atomic.AddUint64(&nonce1, 1),
	})
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	suite.Require().Equal(string(res.Ret), "reflect: Call with too few input arguments")
}

//tc：非中继链管理员设置币价，币价设置失败
func (suite *Model18) Test1833_SetTokenPriceWithNoAdminIsFail() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ServiceRegistryContractAddr.Address(), "SetTokenPrice", nil, rpcx.Uint64(1))
	suite.Require().Nil(err)
	suite.Require().Equal(pb.Receipt_FAILED, res.Status)
	suite.Require().Equal(string(res.Ret), "call error: 2010000:you have no permission")
}

//tc：获取币价，币价获取成功
func (suite *Model18) Test1833_GetTokenPriceIsSuccess() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	price, err := suite.GetTokenPrice(pk)
	suite.Require().Nil(err)
	suite.Require().Equal(uint64(1), price)
}

//tc：根据正确的域名获取一级域名过期时间，过期时间获取成功
func (suite *Model18) Test1834_GetDomainExpiresWithGoodDomainIsSuccess() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, Year, constant.ServiceResolverContractAddr.String())
	suite.Require().Nil(err)

	expires, err := suite.GetDomainExpires(pk, domain+".hub")
	suite.Require().Nil(err)
	suite.Require().Greater(expires, uint64(0))
}

//tc：根据不存在的域名获取一级域名过期时间，过期时间获取失败
func (suite *Model18) Test1835_GetDomainExpiresWithNoExistDomainIsFail() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	_, err = suite.GetDomainExpires(pk, domain+".hub")
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "call error: 2140000:The domain id is not registered")
}

//tc：根据空的域名获取一级域名过期时间，过期时间获取失败
func (suite *Model18) Test1836_GetDomainExpiresWithEmptyDomainIsFail() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	_, err = suite.GetDomainExpires(pk, "")
	suite.Require().NotNil(err)
	suite.Require().Contains(err.Error(), "call error: 2140000:The domain id can not be an empty string")
}

//tc：根据正确的域名获取域名是否被注册，获取成功
func (suite *Model18) Test1837_RecordExistsWithGoodDomainIsSuccess() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, Year, constant.ServiceResolverContractAddr.String())
	suite.Require().Nil(err)

	exists, err := suite.RecordExists(pk, domain+".hub")
	suite.Require().Nil(err)
	suite.Require().Equal(true, exists)
}

//tc：根据不存在的域名获取域名是否被注册，获取成功
func (suite *Model18) Test1838_RecordExistsWithNoExistDomainIsFail() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	exists, err := suite.RecordExists(pk, domain+".hub")
	suite.Require().Nil(err)
	suite.Require().Equal(false, exists)
}

//tc：根据空的域名获取域名是否被注册，获取失败
func (suite *Model18) Test1839_RecordExistsWithEmptyDomainIsFail() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	_, err = suite.RecordExists(pk, "")
	suite.Require().NotNil(err)
	suite.Require().Equal("call error: 2140000:The domain id can not be an empty string", err.Error())
}

//tc：根据正确的域名获取域名所有者，获取成功
func (suite *Model18) Test1840_OwnerWithGoodDomainIsSuccess() {
	pk, address, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, Year, constant.ServiceResolverContractAddr.String())
	suite.Require().Nil(err)

	owner, err := suite.Owner(pk, domain+".hub")
	suite.Require().Nil(err)
	suite.Require().Equal(address.String(), owner)
}

//tc：根据不存在的域名获取域名所有者，获取失败
func (suite *Model18) Test1841_OwnerWithNoExistDomainIsFail() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	_, err = suite.Owner(pk, domain+".hub")
	suite.Require().NotNil(err)
	suite.Require().Equal("call error: 2140000:The domain id must be registered", err.Error())
}

//tc：根据空的域名获取域名所有者，获取失败
func (suite *Model18) Test1842_OwnerWithEmptyDomainIsFail() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	_, err = suite.Owner(pk, "")
	suite.Require().NotNil(err)
	suite.Require().Equal("call error: 2140000:The domain id must be registered", err.Error())
}

//tc：根据正确的域名获取解析器合约地址，获取成功
func (suite *Model18) Test1843_ResolverWithGoodDomainIsSuccess() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, Year, constant.ServiceResolverContractAddr.String())
	suite.Require().Nil(err)

	resolver, err := suite.Resolver(pk, domain+".hub")
	suite.Require().Nil(err)
	suite.Require().Equal(constant.ServiceResolverContractAddr.String(), resolver)
}

//tc：根据不存在的域名获取解析器合约地址，获取失败
func (suite *Model18) Test1844_ResolverWithNoExistDomainIsFail() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	_, err = suite.Resolver(pk, domain+".hub")
	suite.Require().NotNil(err)
	suite.Require().Equal("call error: 2140000:The domain id must be registered", err.Error())
}

//tc：根据空的域名获取解析器合约地址，获取失败
func (suite *Model18) Test1845_ResolverWithEmptyDomainIsFail() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	_, err = suite.Resolver(pk, "")
	suite.Require().NotNil(err)
	suite.Require().Equal("call error: 2140000:The domain id must be registered", err.Error())
}

func (suite *Model18) Test1846_GetAllDomainsIsSuccess() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, Year, constant.ServiceResolverContractAddr.String())
	domains, err := suite.GetAllDomains(pk)
	suite.Require().Nil(err)
	suite.Require().Greater(len(domains), 1)
}

const letters = "abcdefghijklmnopqrstuvwxyz"
const (
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

func randomDomain(n int) string {
	b := make([]byte, n)
	for i, cache, remain := 0, rand.Int63(), letterIdxMax; i < n; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letters) {
			b[i] = letters[idx]
			i++
		}
		cache >>= letterIdxBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}

func (suite *Snake) RegisterDomain(pk crypto.PrivateKey, domain string, duration uint64, resolver string) error {
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ServiceRegistryContractAddr.Address(), "Register", nil, rpcx.String(domain), rpcx.Uint64(duration), rpcx.String(resolver))
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf(string(res.Ret))
	}
	return nil
}

func (suite *Model18) RenewDomain(pk crypto.PrivateKey, domain string, duration uint64) error {
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ServiceRegistryContractAddr.Address(), "Renew", nil, rpcx.String(domain), rpcx.Uint64(duration))
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf(string(res.Ret))
	}
	return nil
}

func (suite *Snake) AllocateSubDomain(pk crypto.PrivateKey, parentName, sonName, owner, resolver, serviceName string) error {
	client := suite.NewClient(pk)
	args := []*pb.Arg{
		rpcx.String(parentName),
		rpcx.String(sonName),
		rpcx.String(owner),
		rpcx.String(resolver),
		rpcx.String(serviceName),
	}
	res, err := client.InvokeBVMContract(constant.ServiceRegistryContractAddr.Address(), "AllocateSubDomain", nil, args...)
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf(string(res.Ret))
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

func (suite *Model18) DeleteSecondDomain(pk crypto.PrivateKey, domain string) error {
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ServiceRegistryContractAddr.Address(), "DeleteSecondDomain", nil, rpcx.String(domain))
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf(string(res.Ret))
	}
	return nil
}

func (suite *Model18) SetPriceLevel(price1, price2, price3, price4, price5 uint64) error {
	pk, address, err := repo.Node1Priv()
	if err != nil {
		return err
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ServiceRegistryContractAddr.Address(), "SetPriceLevel", &rpcx.TransactOpts{
		From:  address.String(),
		Nonce: atomic.AddUint64(&nonce1, 1),
	}, rpcx.Uint64(price1), rpcx.Uint64(price2), rpcx.Uint64(price3), rpcx.Uint64(price4), rpcx.Uint64(price5))
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf(string(res.Ret))
	}
	// TODO: vote proposal
	return nil
}

func (suite *Model18) GetPriceLevel(pk crypto.PrivateKey) (*PriceLevel, error) {
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ServiceRegistryContractAddr.Address(), "GetPriceLevel", nil)
	if err != nil {
		return nil, err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return nil, fmt.Errorf(string(res.Ret))
	}
	var level PriceLevel
	err = json.Unmarshal(res.Ret, &level)
	if err != nil {
		return nil, err
	}
	return &level, nil
}

func (suite *Model18) SetTokenPrice(price uint64) error {
	pk, address, err := repo.Node1Priv()
	if err != nil {
		return err
	}
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ServiceRegistryContractAddr.Address(), "SetTokenPrice", &rpcx.TransactOpts{
		From:  address.String(),
		Nonce: atomic.AddUint64(&nonce1, 1),
	}, rpcx.Uint64(price))
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf(string(res.Ret))
	}
	// TODO: vote proposal
	return nil
}

func (suite *Model18) GetTokenPrice(pk crypto.PrivateKey) (uint64, error) {
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ServiceRegistryContractAddr.Address(), "GetTokenPrice", nil)
	if err != nil {
		return 0, err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return 0, fmt.Errorf(string(res.Ret))
	}
	return binary.BigEndian.Uint64(res.Ret), nil
}

func (suite *Model18) GetDomainExpires(pk crypto.PrivateKey, domain string) (uint64, error) {
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ServiceRegistryContractAddr.Address(), "GetDomainExpires", nil, rpcx.String(domain))
	if err != nil {
		return 0, err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return 0, fmt.Errorf(string(res.Ret))
	}
	return binary.BigEndian.Uint64(res.Ret), nil
}

func (suite *Model18) RecordExists(pk crypto.PrivateKey, domain string) (bool, error) {
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ServiceRegistryContractAddr.Address(), "RecordExists", nil, rpcx.String(domain))
	if err != nil {
		return false, err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return false, fmt.Errorf(string(res.Ret))
	}
	ok, err := strconv.ParseBool(string(res.Ret))
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (suite *Model18) Owner(pk crypto.PrivateKey, domain string) (string, error) {
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ServiceRegistryContractAddr.Address(), "Owner", nil, rpcx.String(domain))
	if err != nil {
		return "", err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return "", fmt.Errorf(string(res.Ret))
	}
	return string(res.Ret), nil
}

func (suite *Model18) Resolver(pk crypto.PrivateKey, domain string) (string, error) {
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ServiceRegistryContractAddr.Address(), "Resolver", nil, rpcx.String(domain))
	if err != nil {
		return "", err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return "", fmt.Errorf(string(res.Ret))
	}
	return string(res.Ret), nil
}

func (suite *Model18) GetAllDomains(pk crypto.PrivateKey) ([]ServDomain, error) {
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ServiceRegistryContractAddr.Address(), "GetAllDomains", nil)
	if err != nil {
		return nil, err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return nil, fmt.Errorf(string(res.Ret))
	}
	var servDomain []ServDomain
	err = json.Unmarshal(res.Ret, &servDomain)
	if err != nil {
		return nil, err
	}
	return servDomain, nil
}

func Domain(domain string) string {
	return domain + ".hub"
}

func SonDomain(domain, sonDomain string) string {
	return sonDomain + "." + domain + ".hub"
}
