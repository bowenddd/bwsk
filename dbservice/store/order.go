package store

import (
	"fmt"
	"seckill/common/entity"
	"sync"

	"gorm.io/gorm"
)

type OrderStore interface {
	CreateByDbPLock(order *entity.Order) error

	CreateByServLock(order *entity.Order) error

	CreateByServChan(order *entity.Order) error

	DeleteById(id uint) error

	FindByOrderId(id uint) (entity.Order, error)

	FindByUserId(id uint) ([]entity.Order, error)

	FindByProductId(id uint) ([]entity.Order, error)

	List() ([]entity.Order, error)

	ClearOrders() error

	HandleOrderChan()
}

type OrderOp struct {
	db *gorm.DB
	mu *sync.Mutex
	ch chan orderChan
}

type orderChan struct {
	res   chan error
	order *entity.Order
}

func (o *OrderOp) CreateByDbPLock(order *entity.Order) error {
	// 执行事务，首先减库存，如果减库存失败，即rowsaffected!=1，则事务失败
	// 然后再执行新增订单表操作
	tx := o.DB().Begin()
	exec := tx.Exec("update product set stock = stock - ? where id = ? and stock >= ?",
		order.Num, order.ProductId, order.Num)
	affected := exec.RowsAffected
	if affected != 1 {
		tx.Rollback()
		return fmt.Errorf("stock is not enough or product not exist")
	}
	create := tx.Create(order)
	if create.Error != nil || exec.Error != nil {
		tx.Rollback()
		return fmt.Errorf("order transaction error")
	}
	tx.Commit()
	return nil
}

func (o *OrderOp) CreateByServLock(order *entity.Order) error {
	// 执行事务，首先减库存，如果减库存失败，即rowsaffected!=1，则事务失败
	// 然后再执行新增订单表操作
	o.mu.Lock()
	defer o.mu.Unlock()
	tx := o.DB().Begin()
	stock := 0
	o.DB().Model(&entity.Product{}).Where("id = ?", order.ProductId).Select("stock").Scan(&stock)
	if stock < order.Num {
		tx.Rollback()
		return fmt.Errorf("stock is not enough or product not exist")
	}
	exec := tx.Exec("update product set stock = stock - ? where id = ?",
		order.Num, order.ProductId)
	create := tx.Create(order)
	if create.Error != nil || exec.Error != nil {
		tx.Rollback()
		return fmt.Errorf("order transaction error")
	}
	tx.Commit()
	return nil
}

func (o *OrderOp) CreateByServChan(order *entity.Order) error {
	ochan := orderChan{
		res:   make(chan error),
		order: order,
	}
	o.ch <- ochan
	return <-ochan.res

}

func (o *OrderOp) HandleOrderChan() {
	for ochan := range o.ch {
		// 执行事务，首先减库存，如果减库存失败，即rowsaffected!=1，则事务失败
		// 然后再执行新增订单表操作
		tx := o.DB().Begin()
		stock := 0
		o.DB().Model(&entity.Product{}).Where("id = ?", ochan.order.ProductId).Select("stock").Scan(&stock)
		if stock < ochan.order.Num {
			tx.Rollback()
			ochan.res <- fmt.Errorf("stock is not enough or product not exist")
			continue
		}
		exec := tx.Exec("update product set stock = stock - ? where id = ?",
			ochan.order.Num, ochan.order.ProductId)
		create := tx.Create(ochan.order)
		if create.Error != nil || exec.Error != nil {
			tx.Rollback()
			ochan.res <- fmt.Errorf("order transaction error")
			continue
		}
		tx.Commit()
		ochan.res <- nil
	}
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

func (o *OrderOp) ClearOrders() error {
	res := o.DB().Exec("delete from orders")
	return res.Error
}

var _ OrderStore = (*OrderOp)(nil)

func (o *OrderOp) DB() *gorm.DB {
	return o.db
}
