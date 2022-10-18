package service

import (
	cacherpc "seckill/cacheservice/rpc"
	"seckill/common/entity"
	"seckill/common/interfaces"
	dbrpc "seckill/dbservice/rpc"
	"seckill/registercenter/registerservice"
	"strings"
	"sync"
)

type OrderServImpl struct {
	registerCenter *registerservice.RegisterCenter
}

func (o *OrderServImpl) DbRpcCli() *dbrpc.OrderRpcServCli {
	sc, err := o.registerCenter.GetDbClient()
	if err != nil {
		return nil
	}
	ordersc := sc.GetOrderRpcServCli()
	return &ordersc
}

func (o *OrderServImpl) CacheRpcCli() *cacherpc.CacheServCli {
	csc, err := o.registerCenter.GetCacheClient()
	if err != nil {
		return nil
	}
	return csc
}

func (o *OrderServImpl) AddOrder(order *entity.Order, method string) error {
	if strings.Contains(method, "CACHE") {
		return o.CacheRpcCli().CreateOrder(order, method)
	}
	return o.DbRpcCli().AddOrder(order, method)
}

func (o *OrderServImpl) GetOrderById(id uint) (entity.Order, error) {
	return o.DbRpcCli().GetOrderById(id)
}

func (o *OrderServImpl) GetOrdersByUID(uid uint) ([]entity.Order, error) {
	return o.DbRpcCli().GetOrdersByUID(uid)
}

func (o *OrderServImpl) GetOrdersByPID(pid uint) ([]entity.Order, error) {
	return o.DbRpcCli().GetOrdersByPID(pid)
}

func (o *OrderServImpl) DeleteOrder(id uint) error {
	return o.DbRpcCli().DeleteOrder(id)
}

func (o *OrderServImpl) GetOrders() ([]entity.Order, error) {
	return o.DbRpcCli().GetOrders()
}

func (o *OrderServImpl) ClearOrders() error {
	return o.DbRpcCli().ClearOrders()
}

var _ interfaces.OrderServ = (*OrderServImpl)(nil)

var orderServOnce = new(sync.Once)

var orderServ *OrderServImpl

func GetOrderService() interfaces.OrderServ {
	orderServOnce.Do(func() {
		center := registerservice.GetRegisterCenter()
		go func() {
			ch := make(chan error, 0)
			registerservice.GetRegisterCenter().Discovery(ch)
			<-ch
		}()
		orderServ = &OrderServImpl{
			registerCenter: center,
		}
	})
	return orderServ
}
