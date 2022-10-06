package service

import (
	"seckill/dbservice/store"
	"seckill/entity"
	"sync"
)

type ProductServ interface {
	AddProduct(product *entity.Product) error

	GetProduct(name string) (entity.Product, error)

	GetProducts() ([]entity.Product, error)

	DeleteProduct(name string) error

	SetStock(id uint, num int) error
}

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

var productServOnce sync.Once

var productServ ProductServ

func init() {
	productServOnce = sync.Once{}
}

func GetProductServ() ProductServ {
	productServOnce.Do(func() {
		productServ = &ProductServImpl{
			store: store.Mysql.NewProductStore(),
		}
	})
	return productServ
}

var _ ProductServ = (*ProductServImpl)(nil)
