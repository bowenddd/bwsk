package store

import (
	"fmt"
	"gorm.io/gorm"
	"seckill/entity"
)

type OrderStore interface {
	Create(order *entity.Order) error

	DeleteById(id uint) error

	FindByOrderId(id uint) (entity.Order, error)

	FindByUserId(id uint) ([]entity.Order, error)

	FindByProductId(id uint) ([]entity.Order, error)

	List() ([]entity.Order, error)
}

type OrderOp struct {
	db *gorm.DB
}

func (o *OrderOp) Create(order *entity.Order) error {
	err := o.DB().Transaction(func(tx *gorm.DB) error {
		exec := tx.Exec("update product set stock = stock - ? where id = ? and stock >= ?",
			order.Num, order.ProductId, order.Num)
		affected := exec.RowsAffected
		if affected != 1 {
			return fmt.Errorf("stock is not enough or no product not exist")
		}
		create := tx.Create(order)
		if create.Error != nil || exec.Error != nil {
			fmt.Errorf("order transaction error")
		}
		return nil
	})
	return err
}

func (o *OrderOp) DeleteById(id uint) error {
	res := o.DB().Exec("delete from orders where id = ?", id)
	return res.Error
}

func (o *OrderOp) FindByOrderId(id uint) (entity.Order, error) {
	order := entity.Order{}
	res := o.DB().Where("id = ?", id).First(&order)
	return order, res.Error
}

func (o *OrderOp) FindByUserId(uid uint) ([]entity.Order, error) {
	var orders []entity.Order
	res := o.DB().Where("user_id = ?", uid).Find(&orders)
	return orders, res.Error
}

func (o *OrderOp) FindByProductId(pid uint) ([]entity.Order, error) {
	var orders []entity.Order
	res := o.DB().Where("product_id = ?", pid).Find(&orders)
	return orders, res.Error
}

func (o *OrderOp) List() ([]entity.Order, error) {
	var orders []entity.Order
	res := o.DB().Find(&orders)
	return orders, res.Error
}

var _ OrderStore = (*OrderOp)(nil)

func (o *OrderOp) DB() *gorm.DB {
	return o.db
}
