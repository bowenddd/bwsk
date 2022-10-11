package service

import (
	"seckill/common/interfaces"
	"seckill/dbservice/rpc"
	"sync"
	"seckill/common/entity"
)

type ProductServImpl struct {
	cli *rpc.ProductRpcServCli
}

var productServOnce = new(sync.Once)

var producrServ *ProductServImpl

var _ interfaces.ProductServ = (*ProductServImpl)(nil)

func (p *ProductServImpl) AddProduct(product *entity.Product) error{
	return p.cli.AddProduct(product)

}

func (p *ProductServImpl)GetProduct(name string) (entity.Product, error){
	return p.cli.GetProduct(name)
}

func (p *ProductServImpl)GetProducts() ([]entity.Product, error){
	return p.cli.GetProducts()
}

func(p *ProductServImpl)DeleteProduct(name string) error{
	return p.cli.DeleteProduct(name)
}

func(p *ProductServImpl)SetStock(id uint, num int) error{
	return p.cli.SetStock(id, num)
}


func GetProductServ() interfaces.ProductServ {
	cli, err := GetDbServRpcCli()
	if err != nil {
		return (*ProductServImpl)(nil)
	}
	productServOnce.Do(func() {
		productServCli := cli.GetProductRpcServCli()
		producrServ = &ProductServImpl{
			cli: &productServCli,
		}
	})
	return producrServ
}
