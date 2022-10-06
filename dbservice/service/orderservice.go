package service

import (
	"seckill/dbservice/store"
	"seckill/entity"
	"sync"
)

type OrderServ interface {
	AddOrder(order *entity.Order) error

	GetOrderById(id uint) (entity.Order, error)

	GetOrdersByUID(uid uint) ([]entity.Order, error)

	GetOrdersByPID(pid uint) ([]entity.Order, error)

	DeleteOrder(id uint) error

	GetOrders() ([]entity.Order, error)
}

type OrderServImpl struct {
	store store.OrderStore
}

var orderServOnce sync.Once

var orderServ OrderServ

func init() {
	orderServOnce = sync.Once{}
}

func GetOrderServ() OrderServ {
	orderServOnce.Do(func() {
		orderServ = &OrderServImpl{
			store: store.Mysql.NewOrderStore(),
		}
	})
	return orderServ
}

func (o OrderServImpl) AddOrder(order *entity.Order) error {
	return o.store.Create(order)
}

func (o OrderServImpl) GetOrderById(id uint) (entity.Order, error) {
	return o.store.FindByOrderId(id)
}

func (o OrderServImpl) GetOrdersByUID(uid uint) ([]entity.Order, error) {
	return o.store.FindByUserId(uid)
}

func (o OrderServImpl) GetOrdersByPID(pid uint) ([]entity.Order, error) {
	return o.store.FindByProductId(pid)
}

func (o OrderServImpl) DeleteOrder(id uint) error {
	return o.store.DeleteById(id)
}

func (o OrderServImpl) GetOrders() (orders []entity.Order, err error) {
	return o.store.List()
}

var _ OrderServ = (*OrderServImpl)(nil)
