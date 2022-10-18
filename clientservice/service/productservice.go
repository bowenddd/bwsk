package service

import (
	cacherpc "seckill/cacheservice/rpc"
	"seckill/common/entity"
	"seckill/common/interfaces"
	dbrpc "seckill/dbservice/rpc"
	"seckill/registercenter/registerservice"
	"strings"
	"sync"
)

type ProductServImpl struct {
	registerCenter *registerservice.RegisterCenter
}

var productServOnce = new(sync.Once)

var producrServ *ProductServImpl

var _ interfaces.ProductServ = (*ProductServImpl)(nil)

func (p *ProductServImpl) DbRpcCli() *dbrpc.ProductRpcServCli {
	sc, err := p.registerCenter.GetDbClient()
	if err != nil {
		return nil
	}
	psc := sc.GetProductRpcServCli()
	return &psc
}

func (p *ProductServImpl) CacheRpcCli() *cacherpc.CacheServCli {
	csc, err := p.registerCenter.GetCacheClient()
	if err != nil {
		return nil
	}
	return csc
}

func (p *ProductServImpl) AddProduct(product *entity.Product) error {
	return p.DbRpcCli().AddProduct(product)

}

func (p *ProductServImpl) GetProduct(name string) (entity.Product, error) {
	return p.DbRpcCli().GetProduct(name)
}

func (p *ProductServImpl) GetProducts() ([]entity.Product, error) {
	return p.DbRpcCli().GetProducts()
}

func (p *ProductServImpl) DeleteProduct(name string) error {
	return p.DbRpcCli().DeleteProduct(name)
}

func (p *ProductServImpl) SetStock(id uint, num int) error {
	err := p.DbRpcCli().SetStock(id, num)
	if err != nil {
		return err
	}
	err = p.CacheRpcCli().SetStock(id, num, 0)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProductServImpl) GetStock(id uint, method string) (int, error) {
	if strings.Contains(method, "CACHE") {
		return p.CacheRpcCli().GetStock(id)
	}
	return p.DbRpcCli().GetStock(id, method)
}

func GetProductService() interfaces.ProductServ {
	productServOnce.Do(func() {
		center := registerservice.GetRegisterCenter()
		go func() {
			ch := make(chan error, 0)
			registerservice.GetRegisterCenter().Discovery(ch)
			<-ch
		}()
		producrServ = &ProductServImpl{
			registerCenter: center,
		}
	})
	return producrServ
}
