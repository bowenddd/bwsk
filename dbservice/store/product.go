package store

import (
	"fmt"
	"math"
	"seckill/common/entity"

	"gorm.io/gorm"
)

type ProductStore interface {
	Create(product *entity.Product) error

	FindByName(name string) (entity.Product, error)

	DeleteByName(name string) error

	List() ([]entity.Product, error)

	DecreaseStock(name string, num int) error

	SetStock(id uint, num int) error

	GetStock(id uint) (int, error)
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

func (p ProductOp) GetStock(id uint) (int, error) {
	stock := math.MinInt32
	res := p.DB().Raw("select stock from product where id = ?", id).Scan(&stock)
	if stock == math.MinInt32 {
		return 0, fmt.Errorf("product not found")
	}
	return stock, res.Error
}

var _ ProductStore = (*ProductOp)(nil)
