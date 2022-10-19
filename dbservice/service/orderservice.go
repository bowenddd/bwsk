package service

import (
	"fmt"
	"seckill/common/consts"
	"seckill/common/entity"
	"seckill/common/interfaces"
	"seckill/dbservice/store"
)

type CreateOrderFunc func(order *entity.Order) error

type CreateOrderStrategy interface {
	CreateOrder(order *entity.Order) error
}

func (f CreateOrderFunc) CreateOrder(order *entity.Order) error {
	return f(order)
}

type OrderServImpl struct {
	store    store.OrderStore
	strategy CreateOrderStrategy
}

func GetOrderServ() interfaces.OrderServ {
	orderServ := &OrderServImpl{
		store: store.Mysql.NewOrderStore(),
	}
	return orderServ
}

func (o *OrderServImpl) SetStragegy(f CreateOrderFunc) {
	o.strategy = f
}

func (o *OrderServImpl) ExecCreateOrder(order *entity.Order) error {
	if o.strategy == nil {
		return fmt.Errorf("strategy not set")
	}
	return o.strategy.CreateOrder(order)
}

func (o OrderServImpl) AddOrder(order *entity.Order, method string) error {
	switch method {
	case consts.DBPESSIMISTICLOCK:
		o.SetStragegy(o.store.CreateByDbPLock)
	case consts.SERVICELOCK:
		o.SetStragegy(o.store.CreateByServLock)
	case consts.SERVICECHANNEL:
		o.SetStragegy(o.store.CreateByServChan)
	case consts.DBOPTIMISTICLOCK:
		o.SetStragegy(o.store.CreateByDbOLock)
	case consts.NOMEASURE:
		o.SetStragegy(o.store.CreateWithNoMeasure)
	default:
		return fmt.Errorf("method %s not supported", method)
	}
	return o.ExecCreateOrder(order)
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

func (o OrderServImpl) ClearOrders() error {
	return o.store.ClearOrders()
}

var _ interfaces.OrderServ = (*OrderServImpl)(nil)
