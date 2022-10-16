package service

import (
	cacherpc "seckill/cacheservice/rpc"
	"seckill/common/entity"
	"seckill/common/interfaces"
	dbrpc "seckill/dbservice/rpc"
	"strings"
	"sync"
)

type ProductServImpl struct {
	dbCli    *dbrpc.ProductRpcServCli
	cacheCli *cacherpc.CacheServCli
}

var productServOnce = new(sync.Once)

var producrServ *ProductServImpl

var _ interfaces.ProductServ = (*ProductServImpl)(nil)

func (p *ProductServImpl) AddProduct(product *entity.Product) error {
	return p.dbCli.AddProduct(product)

}

func (p *ProductServImpl) GetProduct(name string) (entity.Product, error) {
	return p.dbCli.GetProduct(name)
}

func (p *ProductServImpl) GetProducts() ([]entity.Product, error) {
	return p.dbCli.GetProducts()
}

func (p *ProductServImpl) DeleteProduct(name string) error {
	return p.dbCli.DeleteProduct(name)
}

func (p *ProductServImpl) SetStock(id uint, num int) error {
	err := p.dbCli.SetStock(id, num)
	if err != nil {
		return err
	}
	err = p.cacheCli.SetStock(id, num, 0)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProductServImpl) GetStock(id uint, method string) (int, error) {
	if strings.Contains(method, "CACHE") {
		return p.cacheCli.GetStock(id)
	}
	return p.dbCli.GetStock(id, method)
}

func GetProductService() interfaces.ProductServ {
	dbcli, err := dbrpc.GetDbServRpcCli()
	if err != nil {
		return (*ProductServImpl)(nil)
	}
	cachecli, err := cacherpc.NewCacheServClient()
	if err != nil {
		return (*ProductServImpl)(nil)
	}
	productServOnce.Do(func() {
		productServCli := dbcli.GetProductRpcServCli()
		producrServ = &ProductServImpl{
			dbCli:    &productServCli,
			cacheCli: cachecli,
		}
	})
	return producrServ
}
