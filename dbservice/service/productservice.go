package service

import (
	"seckill/common/entity"
	"seckill/common/interfaces"
	"seckill/dbservice/store"
	"sync"
)

type ProductServImpl struct {
	store store.ProductStore
}

func (p ProductServImpl) AddProduct(product *entity.Product) error {
	return p.store.Create(product)
}

func (p ProductServImpl) GetProduct(name string) (entity.Product, error) {
	return p.store.FindByName(name)
}

func (p ProductServImpl) GetProducts() ([]entity.Product, error) {
	return p.store.List()
}

func (p ProductServImpl) DeleteProduct(name string) error {
	return p.store.DeleteByName(name)
}

func (p ProductServImpl) SetStock(id uint, num int) error {
	return p.store.SetStock(id, num)
}

func (p ProductServImpl)GetStock(id uint) (int, error){
	return p.store.GetStock(id)
}

var productServOnce sync.Once

var productServ interfaces.ProductServ

func init() {
	productServOnce = sync.Once{}
}

func GetProductServ() interfaces.ProductServ {
	productServOnce.Do(func() {
		productServ = &ProductServImpl{
			store: store.Mysql.NewProductStore(),
		}
	})
	return productServ
}

var _ interfaces.ProductServ = (*ProductServImpl)(nil)
