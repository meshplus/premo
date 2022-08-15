package bxh_tester

import (
	"encoding/json"
	"fmt"
	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-model/constant"
	"github.com/meshplus/bitxhub-model/pb"
	rpcx "github.com/meshplus/go-bitxhub-client"
	"github.com/meshplus/premo/internal/repo"
)

type ServDomainData struct {
	Addr        map[uint64]string `json:"addr"`
	ServiceName string            `json:"serviceName"`
	Des         string            `json:"des"`
	Dids        []string          `json:"dids"`
}

type Model19 struct {
	*Snake
}

func (suite *Model19) SetupTest() {
	suite.T().Parallel()
}

//tc：根据正确的域名设置服务域名数据，域名数据设置成功
func (suite *Model19) Test1901_SetServDomainDataWithGoodDomainIsSuccess() {
	pk, address, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, Year, constant.ServiceResolverContractAddr.String())
	suite.Require().Nil(err)

	err = suite.SetServDomainData(pk, domain+".hub", uint64(1), address.String(), "service", "desc", "dids")
	suite.Require().Nil(err)
}

//tc：根据不存在的域名设置服务域名数据，域名数据设置失败
func (suite *Model19) Test1902_SetServDomainDataWithNoExistDomainIsFail() {
	pk, address, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)

	err = suite.SetServDomainData(pk, domain+".hub", uint64(1), address.String(), "service", "desc", "dids")
	suite.Require().NotNil(err)
	suite.Require().Equal("call error: 2140000:The domain name does not belong to you", err.Error())
}

//tc：根据空的域名设置服务域名数据，域名数据设置失败
func (suite *Model19) Test1903_SetServDomainDataWithEmptyDomainIsFail() {
	pk, address, err := repo.KeyPriv()
	suite.Require().Nil(err)

	err = suite.SetServDomainData(pk, "", uint64(1), address.String(), "service", "desc", "dids")
	suite.Require().NotNil(err)
	suite.Require().Equal("call error: 2140000:The domain name does not belong to you", err.Error())
}

//tc：根据正确的链类型设置服务域名数据，域名数据设置成功
func (suite *Model19) Test1904_SetServDomainDataWithGoodTypeIsSuccess() {
	pk, address, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, Year, constant.ServiceResolverContractAddr.String())
	suite.Require().Nil(err)

	err = suite.SetServDomainData(pk, domain+".hub", uint64(1), address.String(), "service", "desc", "dids")
	suite.Require().Nil(err)
}

//tc：根据错误的类型设置服务域名数据，域名数据设置失败
func (suite *Model19) Test1905_SetServDomainDataWithErrorTypeIsFail() {

}

//tc：根据正确的地址设置服务域名数据，域名数据设置成功
func (suite *Model19) Test1906_SetServDomainDataWithGoodAddrIsSuccess() {
	pk, address, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, Year, constant.ServiceResolverContractAddr.String())
	suite.Require().Nil(err)

	err = suite.SetServDomainData(pk, domain+".hub", uint64(1), address.String(), "service", "desc", "dids")
	suite.Require().Nil(err)
}

//tc：根据空的地址设置服务域名数据，域名数据设置失败
func (suite *Model19) Test1907_SetServDomainDataWithEmptyAddrIsFail() {

}

//tc：根据正确的服务名称设置域名数据，域名数据设置成功
func (suite *Model19) Test1908_SetServDomainDataWithGoodServiceIsSuccess() {
	pk, address, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, Year, constant.ServiceResolverContractAddr.String())
	suite.Require().Nil(err)

	err = suite.SetServDomainData(pk, domain+".hub", uint64(1), address.String(), "service", "desc", "dids")
	suite.Require().Nil(err)
}

//tc：根据空的服务名称设置域名数据，域名数据设置失败
func (suite *Model19) Test1909_SetServDomainDataWithEmptyServiceIsFail() {

}

//tc：根据正确的描述设置域名数据，域名数据设置成功
func (suite *Model19) Test1910_SetServDomainDataWithGoodDescIsSuccess() {
	pk, address, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, Year, constant.ServiceResolverContractAddr.String())
	suite.Require().Nil(err)

	err = suite.SetServDomainData(pk, domain+".hub", uint64(1), address.String(), "service", "desc", "dids")
	suite.Require().Nil(err)
}

//tc：根据空的描述设置域名数据，域名数据设置成功
func (suite *Model19) Test1911_SetServDomainDataWithEmptyDescIsFail() {

}

//tc：根据正确的did设置域名数据，域名数据设置成功
func (suite *Model19) Test1912_SetServDomainDataWithGoodDidsIsSuccess() {
	pk, address, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, Year, constant.ServiceResolverContractAddr.String())
	suite.Require().Nil(err)

	err = suite.SetServDomainData(pk, domain+".hub", uint64(1), address.String(), "service", "desc", "dids")
	suite.Require().Nil(err)
}

//tc：根据空的did设置域名数据，域名数据设置成功
func (suite *Model19) Test1913_SetServDomainDataWithEmptyDidsIsFail() {

}

//tc：根据正确的域名获取域名数据，域名数据获取成功
func (suite *Model19) Test1914_GetServDomainDataWithGoodDomainIsSuccess() {
	pk, address, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, Year, constant.ServiceResolverContractAddr.String())
	suite.Require().Nil(err)

	err = suite.SetServDomainData(pk, domain+".hub", uint64(1), address.String(), domain, "desc", "dids")
	suite.Require().Nil(err)

	data, err := suite.GetServDomainData(pk, domain+".hub")
	suite.Require().Nil(err)
	suite.Require().Equal(domain, data.ServiceName)
}

//tc：根据不存在的域名获取域名数据，域名数据获取失败
func (suite *Model19) Test1915_GetServDomainDataWithNoExistDomainIsFail() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	_, err = suite.GetServDomainData(pk, domain+".hub")
	suite.Require().NotNil(err)
	suite.Require().Equal("call error: 2140000:The domain id must be registered", err.Error())
}

//tc：根据空的域名获取域名数据，域名数据获取失败
func (suite *Model19) Test1916_GetServDomainDataWithEmptyDomainIsFail() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	_, err = suite.GetServDomainData(pk, domain+".hub")
	suite.Require().NotNil(err)
	suite.Require().Equal("call error: 2140000:The domain id must be registered", err.Error())
}

//tc：根据正确的域名设置地址，地址设置成功
func (suite *Model19) Test1917_SetAddrWithGoodDomainIsSuccess() {
	pk, address, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, Year, constant.ServiceResolverContractAddr.String())
	suite.Require().Nil(err)

	err = suite.SetServDomainData(pk, domain+".hub", uint64(1), address.String(), domain, "desc", "dids")
	suite.Require().Nil(err)

	err = suite.SetAddr(pk, domain+".hub", address.String())
	suite.Require().Nil(err)
}

//tc：根据不存在的域名设置地址，地址设置失败
func (suite *Model19) Test1918_SetAddrWithNoExistDomainIsFail() {
	pk, address, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)

	err = suite.SetAddr(pk, domain+".hub", address.String())
	suite.Require().NotNil(err)
	suite.Require().Equal("call error: 2140000:The domain name does not belong to you", err.Error())
}

//tc：根据空的域名设置地址，地址设置失败
func (suite *Model19) Test1919_SetAddrWithEmptyDomainIsFail() {
	pk, address, err := repo.KeyPriv()
	suite.Require().Nil(err)

	err = suite.SetAddr(pk, "", address.String())
	suite.Require().NotNil(err)
	suite.Require().Equal("call error: 2140000:The domain name does not belong to you", err.Error())
}

//tc：根据正确的链类型设置地址，地址设置成功
func (suite *Model19) Test1920_SetAddrWithGoodTypeIsSuccess() {
	pk, address, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, Year, constant.ServiceResolverContractAddr.String())
	suite.Require().Nil(err)

	err = suite.SetServDomainData(pk, domain+".hub", uint64(1), address.String(), domain, "desc", "dids")
	suite.Require().Nil(err)

	err = suite.SetAddr(pk, domain+".hub", address.String())
	suite.Require().Nil(err)
}

//tc：根据错误的链类型设置地址，地址设置失败
func (suite *Model19) Test1921_SetAddrWithErrorTypeIsFail() {

}

//tc：根据正确的地址设置地址，地址设置成功
func (suite *Model19) Test1922_SetAddrWithGoodAddressIsSuccess() {
	pk, address, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, Year, constant.ServiceResolverContractAddr.String())
	suite.Require().Nil(err)

	err = suite.SetServDomainData(pk, domain+".hub", uint64(1), address.String(), domain, "desc", "dids")
	suite.Require().Nil(err)

	err = suite.SetAddr(pk, domain+".hub", address.String())
	suite.Require().Nil(err)
}

//tc：根据错误的地址设置地址，地址设置失败
func (suite *Model19) Test1923_SetAddrWithErrorAddressIsFail() {

}

//tc：根据正确的域名设置服务名，服务名设置成功
func (suite *Model19) Test1923_SetServiceNameWithGoodDomainIsSuccess() {
	pk, address, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, Year, constant.ServiceResolverContractAddr.String())
	suite.Require().Nil(err)

	err = suite.SetServDomainData(pk, domain+".hub", uint64(1), address.String(), domain, "desc", "dids")
	suite.Require().Nil(err)

	err = suite.SetServiceName(pk, domain+".hub", domain)
	suite.Require().Nil(err)
}

//tc：根据不存在的域名设置服务名，服务名设置失败
func (suite *Model19) Test1924_SetServiceNameWithNoExistDomainIsFail() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)

	err = suite.SetServiceName(pk, domain+".hub", domain)
	suite.Require().NotNil(err)
	suite.Require().Equal("call error: 2140000:The domain name does not belong to you", err.Error())
}

//tc：根据空的域名设置服务名，服务名设置失败
func (suite *Model19) Test1925_SetServiceNameWithEmptyIsFail() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)

	err = suite.SetServiceName(pk, "", "")
	suite.Require().NotNil(err)
	suite.Require().Equal("call error: 2140000:The domain name does not belong to you", err.Error())
}

//tc：根据正确的服务名设置服务名，服务名设置成功
func (suite *Model19) Test1926_SetServiceNameWithGoodServiceIsSuccess() {
	pk, address, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, Year, constant.ServiceResolverContractAddr.String())
	suite.Require().Nil(err)

	err = suite.SetServDomainData(pk, domain+".hub", uint64(1), address.String(), domain, "desc", "dids")
	suite.Require().Nil(err)

	err = suite.SetServiceName(pk, domain+".hub", domain)
	suite.Require().Nil(err)
}

//tc：根据空的服务名设置服务名，服务名设置失败
func (suite *Model19) Test1927_SetServiceNameWithEmptyServiceIsFail() {
	pk, address, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, Year, constant.ServiceResolverContractAddr.String())
	suite.Require().Nil(err)

	err = suite.SetServDomainData(pk, domain+".hub", uint64(1), address.String(), domain, "desc", "dids")
	suite.Require().Nil(err)

	err = suite.SetServiceName(pk, domain+".hub", "")
	suite.Require().NotNil(err)
	suite.Require().Equal("call error: 2140000:The serviceName can not be an empty string", err.Error())
}

//tc：根据正确的域名设置描述，描述设置成功
func (suite *Model19) Test1928_SetServiceDescWithGoodDomainIsSuccess() {
	pk, address, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, Year, constant.ServiceResolverContractAddr.String())
	suite.Require().Nil(err)

	err = suite.SetServDomainData(pk, domain+".hub", uint64(1), address.String(), domain, "desc", "dids")
	suite.Require().Nil(err)

	err = suite.SetServiceDesc(pk, domain+".hub", "desc")
	suite.Require().Nil(err)
}

//tc：根据不存在的域名设置描述，描述设置失败
func (suite *Model19) Test1929_SetServiceDescWithNoExistDomainIsFail() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.SetServiceDesc(pk, domain+".hub", "desc")
	suite.Require().NotNil(err)
	suite.Require().Equal("call error: 2140000:The domain name does not belong to you", err.Error())
}

//tc：根据空的域名设置描述，描述设置失败
func (suite *Model19) Test1930_SetServiceDescWithEmptyDomainIsFail() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.SetServiceDesc(pk, "", "desc")
	suite.Require().NotNil(err)
	suite.Require().Equal("call error: 2140000:The domain name does not belong to you", err.Error())
}

//tc：根据正确的域名设置DID，DID设置成功
func (suite *Model19) Test1931_SetServiceDidsWithGoodDomainIsSuccess() {
	pk, address, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, Year, constant.ServiceResolverContractAddr.String())
	suite.Require().Nil(err)

	err = suite.SetServDomainData(pk, domain+".hub", uint64(1), address.String(), domain, "desc", "dids")
	suite.Require().Nil(err)

	err = suite.SetDids(pk, domain+".hub", "dids")
	suite.Require().Nil(err)
}

//tc：根据不存在的域名设置DID，DID设置失败
func (suite *Model19) Test1932_SetServiceDidsWithNoExistDomain() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)

	err = suite.SetDids(pk, domain+".hub", "dids")
	suite.Require().NotNil(err)
	suite.Require().Equal("call error: 2140000:The domain name does not belong to you", err.Error())
}

//tc：根据空的域名设置DID，DID设置失败
func (suite *Model19) Test1933_SetServiceDidsWithEmptyDomainIsFail() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)

	err = suite.SetDids(pk, "", "dids")
	suite.Require().NotNil(err)
	suite.Require().Equal("call error: 2140000:The domain name does not belong to you", err.Error())
}

//tc：根据正确的服务名设置反向映射，反向映射设置成功
func (suite *Model19) Test1934_SetReverseWithGoodDomainIsSuccess() {
	pk, address, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, Year, constant.ServiceResolverContractAddr.String())
	suite.Require().Nil(err)

	err = suite.SetServDomainData(pk, domain+".hub", uint64(1), address.String(), domain, "desc", "dids")
	suite.Require().Nil(err)

	err = suite.SetReverse(pk, domain+".hub", domain)
	suite.Require().Nil(err)
}

//tc：根据空的服务名设置反向映射，反向映射设置失败
func (suite *Model19) Test1935_SetReverseWithEmptyServiceIsFail() {

}

//tc：根据正确的域名设置反向映射，反向映射设置成功
func (suite *Model19) Test1936_SetReverseWithGoodDomainIsSuccess() {
	pk, address, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, Year, constant.ServiceResolverContractAddr.String())
	suite.Require().Nil(err)

	err = suite.SetServDomainData(pk, domain+".hub", uint64(1), address.String(), domain, "desc", "dids")
	suite.Require().Nil(err)

	err = suite.SetReverse(pk, domain+".hub", domain)
	suite.Require().Nil(err)
}

//tc：根据不存在的域名设置反向映射，反向映射设置失败
func (suite *Model19) Test1937_SetReverseWithNoExistDomainIsFail() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)

	err = suite.SetReverse(pk, domain+".hub", domain)
	suite.Require().NotNil(err)
	suite.Require().Equal("call error: 2140000:The domain name does not belong to you", err.Error())
}

//tc：根据空的域名设置反向映射，反向映射设置失败
func (suite *Model19) Test1938_SetReverseWithEmptyDomainIsFail() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)

	err = suite.SetReverse(pk, "", "")
	suite.Require().NotNil(err)
	suite.Require().Equal("call error: 2140000:The domain name does not belong to you", err.Error())
}

//tc：根据正确的服务名获取域名，域名获取成功
func (suite *Model19) Test1939_GetReverseNameWithGoodServiceIsSuccess() {
	pk, address, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, Year, constant.ServiceResolverContractAddr.String())
	suite.Require().Nil(err)

	err = suite.SetServDomainData(pk, domain+".hub", uint64(1), address.String(), domain, "desc", "dids")
	suite.Require().Nil(err)

	err = suite.SetReverse(pk, domain+".hub", domain)
	suite.Require().Nil(err)

	name, err := suite.GetReverseName(pk, domain)
	suite.Require().Nil(err)
	suite.Require().Equal(domain+".hub", name)
}

//tc：根据不存在的服务名获取域名，域名获取失败
func (suite *Model19) Test1940_GetReverseNameWithNoExistServiceIsFail() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	name, err := suite.GetReverseName(pk, domain)
	suite.Require().Nil(err)
	suite.Require().Equal("", name)
}

//tc：根据空的服务名获取域名，域名获取失败
func (suite *Model19) Test1941_GetReverseNameWithEmptyDomainIsFail() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	name, err := suite.GetReverseName(pk, "")
	suite.Require().Nil(err)
	suite.Require().Equal("", name)
}

//tc：根据正确的域名获取服务名，服务名获取成功
func (suite *Model19) Test1942_GetServiceNameWithGoodDomainIsSuccess() {
	pk, address, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, Year, constant.ServiceResolverContractAddr.String())
	suite.Require().Nil(err)

	err = suite.SetServDomainData(pk, domain+".hub", uint64(1), address.String(), domain, "desc", "dids")
	suite.Require().Nil(err)

	name, err := suite.GetServiceName(pk, domain+".hub")
	suite.Require().Nil(err)
	suite.Require().Equal(domain, name)
}

//tc：根据不存在的域名获取服务名，服务名获取失败
func (suite *Model19) Test1943_GetServiceNameWithNoExistDomainIsFail() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	_, err = suite.GetServiceName(pk, domain+".hub")
	suite.Require().NotNil(err)
	suite.Require().Equal("call error: 2140000:The domain id must be registered", err.Error())
}

//tc：根据空的域名获取服务名，服务名获取失败
func (suite *Model19) Test1944_GetServiceNameWithEmptyDomainIsFail() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	_, err = suite.GetServiceName(pk, "")
	suite.Require().NotNil(err)
	suite.Require().Equal("call error: 2140000:The domain id must be registered", err.Error())
}

//tc：根据正确的域名获取删除域名数据，域名数据删除成功
func (suite *Model19) Test1945_DeleteServDomainDataWithGoodDomainIsSuccess() {
	pk, address, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	err = suite.RegisterDomain(pk, domain, Year, constant.ServiceResolverContractAddr.String())
	suite.Require().Nil(err)
	sonDomain := randomDomain(3)
	err = suite.AllocateSubDomain(pk, Domain(domain), sonDomain, address.String(), constant.ServiceResolverContractAddr.String())
	suite.Require().Nil(err)

	err = suite.SetServDomainData(pk, SonDomain(domain, sonDomain), uint64(1), address.String(), domain, "desc", "dids")
	suite.Require().Nil(err)
	name, err := suite.GetServiceName(pk, SonDomain(domain, sonDomain))
	suite.Require().Nil(err)
	suite.Require().Equal(domain, name)

	err = suite.DeleteServDomainData(pk, SonDomain(domain, sonDomain))
	suite.Require().Nil(err)
	name, err = suite.GetServiceName(pk, SonDomain(domain, sonDomain))
	suite.Require().Nil(err)
	suite.Require().Equal("", name)
}

//tc：根据不存在的域名获取删除域名数据，域名数据删除失败
func (suite *Model19) Test1946_DeleteServDomainDataWithNoExistDomainIsFail() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	domain := randomDomain(10)
	sonDomain := randomDomain(3)
	err = suite.DeleteServDomainData(pk, SonDomain(domain, sonDomain))
	suite.Require().NotNil(err)
	suite.Require().Equal("call error: 2140000:The domain name does not belong to you", err.Error())
}

//tc：根据空的域名获取删除域名数据，域名数据删除失败
func (suite *Model19) Test1947_DeleteServDomainDataWithEmptyDomainIsFail() {
	pk, _, err := repo.KeyPriv()
	suite.Require().Nil(err)
	err = suite.DeleteServDomainData(pk, "")
	suite.Require().NotNil(err)
	suite.Require().Equal("call error: 2140000:The domain name does not belong to you", err.Error())
}

func (suite *Model19) SetServDomainData(pk crypto.PrivateKey, domain string, coinTyp uint64, addr, service, desc, dids string) error {
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ServiceResolverContractAddr.Address(), "SetServDomainData", nil,
		rpcx.String(domain), rpcx.Uint64(coinTyp), rpcx.String(addr), rpcx.String(service), rpcx.String(desc), rpcx.String(dids))
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf(string(res.Ret))
	}
	return nil
}

func (suite *Model19) GetServDomainData(pk crypto.PrivateKey, domain string) (*ServDomainData, error) {
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ServiceResolverContractAddr.Address(), "GetServDomainData", nil, rpcx.String(domain))
	if err != nil {
		return nil, err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return nil, fmt.Errorf(string(res.Ret))
	}
	var data ServDomainData
	err = json.Unmarshal(res.Ret, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (suite *Model19) SetAddr(pk crypto.PrivateKey, domain, address string) error {
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ServiceResolverContractAddr.Address(), "SetAddr", nil, rpcx.String(domain), rpcx.Uint64(1), rpcx.String(address))
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf(string(res.Ret))
	}
	return nil
}

func (suite *Model19) SetServiceName(pk crypto.PrivateKey, domain, service string) error {
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ServiceResolverContractAddr.Address(), "SetServiceName", nil, rpcx.String(domain), rpcx.String(service), rpcx.Bool(false))
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf(string(res.Ret))
	}
	return nil
}

func (suite *Model19) SetServiceDesc(pk crypto.PrivateKey, domain, desc string) error {
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ServiceResolverContractAddr.Address(), "SetServiceDes", nil, rpcx.String(domain), rpcx.String(desc))
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf(string(res.Ret))
	}
	return nil
}

func (suite *Model19) SetDids(pk crypto.PrivateKey, domain, dids string) error {
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ServiceResolverContractAddr.Address(), "SetDids", nil, rpcx.String(domain), rpcx.String(dids))
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf(string(res.Ret))
	}
	return nil
}

func (suite *Model19) SetReverse(pk crypto.PrivateKey, domain, service string) error {
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ServiceResolverContractAddr.Address(), "SetReverse", nil, rpcx.String(service), rpcx.String(domain))
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf(string(res.Ret))
	}
	return nil
}

func (suite *Model19) GetReverseName(pk crypto.PrivateKey, service string) (string, error) {
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ServiceResolverContractAddr.Address(), "GetReverseName", nil, rpcx.String(service))
	if err != nil {
		return "", err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return "", fmt.Errorf(string(res.Ret))
	}
	return string(res.Ret), nil
}

func (suite *Model19) GetServiceName(pk crypto.PrivateKey, name string) (string, error) {
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ServiceResolverContractAddr.Address(), "GetServiceName", nil, rpcx.String(name))
	if err != nil {
		return "", err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return "", fmt.Errorf(string(res.Ret))
	}
	return string(res.Ret), nil
}

func (suite *Model19) DeleteServDomainData(pk crypto.PrivateKey, domain string) error {
	client := suite.NewClient(pk)
	res, err := client.InvokeBVMContract(constant.ServiceResolverContractAddr.Address(), "DeleteServDomainData", nil, rpcx.String(domain))
	if err != nil {
		return err
	}
	if res.Status != pb.Receipt_SUCCESS {
		return fmt.Errorf(string(res.Ret))
	}
	return nil
}
