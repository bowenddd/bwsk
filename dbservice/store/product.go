package store

import (
	"gorm.io/gorm"
	"seckill/common/entity"
)

type ProductStore interface {
	Create(product *entity.Product) error

	FindByName(name string) (entity.Product, error)

	DeleteByName(name string) error

	List() ([]entity.Product, error)

	DecreaseStock(name string, num int) error

	SetStock(id uint, num int) error
}

type ProductOp struct {
	db *gorm.DB
}

func (p ProductOp) SetStock(id uint, num int) error {
	res := p.DB().Exec("update product set stock = ? where id = ?", num, id)
	return res.Error
}

func (p ProductOp) DecreaseStock(name string, num int) error {
	res := p.DB().Exec("update product set stock = stock - ? where name = ? and stock >= ?",
		num, name, num)
	return res.Error
}

func (p ProductOp) DB() *gorm.DB {
	return p.db
}

func (p ProductOp) Create(product *entity.Product) error {
	res := p.DB().Create(product)
	return res.Error
}

func (p ProductOp) FindByName(name string) (entity.Product, error) {
	product := entity.Product{}
	res := p.DB().Where("name = ?", name).First(&product)
	return product, res.Error
}

func (p ProductOp) DeleteByName(name string) error {
	res := p.DB().Exec("delete from product where name = ?", name)
	return res.Error
}

func (p ProductOp) List() ([]entity.Product, error) {
	var products []entity.Product
	res := p.DB().Find(&products)
	return products, res.Error
}

var _ ProductStore = (*ProductOp)(nil)
