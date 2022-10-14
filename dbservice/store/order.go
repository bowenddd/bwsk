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

	CreateByDbOLock(order *entity.Order) error

	CreateWithNoMeasure(order *entity.Order) error

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

func (o *OrderOp) CreateByDbOLock(order *entity.Order) error {
	succ := false
	retry := 0
	for !succ && retry < 5 {
		p := &entity.Product{}
		tx := o.DB().Begin()
		tx.Model(p).Where("id = ?", order.ProductId).Select("stock", "version").Scan(p)
		if p.Stock < order.Num {
			tx.Rollback()
			return fmt.Errorf("stock is not enough or product not exist")
		}
		// 这里用gorm封装的update导致乐观锁有问题，所以拿原生的sql试一下
		// exec := tx.Model(p).Where("id = ? and version = ?", order.ProductId, oldV).Updates(p)
		// 使用数据库乐观锁，会出现少买的情况，因为大量的请求都失效了，这里增加一个重试。
		exec := tx.Exec("update product set stock = stock - ?, version = version + 1 where id = ? and version = ?",
			order.Num, order.ProductId, p.Version)
		affected := exec.RowsAffected
		if affected != 1 {
			tx.Rollback()
			retry++
			continue
		}
		create := tx.Create(order)
		if create.Error != nil || exec.Error != nil {
			tx.Rollback()
			return fmt.Errorf("order transaction error")
		}
		tx.Commit()
		succ = true
	}
	if succ {
		return nil
	}
	return fmt.Errorf("抢购失败")
}

func (o *OrderOp) CreateByServLock(order *entity.Order) error {
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

func (o *OrderOp) CreateWithNoMeasure(order *entity.Order) error {
	tx := o.DB().Begin()
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

func (o *OrderOp) HandleOrderChan() {
	for ochan := range o.ch {
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
