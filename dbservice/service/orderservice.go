package service

import (
	"fmt"
	"seckill/common/consts"
	"seckill/common/entity"
	"seckill/common/interfaces"
	"seckill/dbservice/store"
	"sync"
)

type OrderServImpl struct {
	store store.OrderStore
}

var orderServOnce sync.Once

var orderServ interfaces.OrderServ

func init() {
	orderServOnce = sync.Once{}
}

func GetOrderServ() interfaces.OrderServ {
	orderServOnce.Do(func() {
		orderServ = &OrderServImpl{
			store: store.Mysql.NewOrderStore(),
		}
	})
	return orderServ
}

func (o OrderServImpl) AddOrder(order *entity.Order, method string) error {
	switch method {
	case consts.DBPESSIMISTICLOCK:
		return o.store.CreateByDbPLock(order)
	case consts.SERVICELOCK:
		return o.store.CreateByServLock(order)
	case consts.SERVICECHANNEL:
		return o.store.CreateByServChan(order)
	default:
		return fmt.Errorf("method %s not supported", method)
	}
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
