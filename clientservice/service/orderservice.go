package service

import (
	cacherpc "seckill/cacheservice/rpc"
	"seckill/common/entity"
	"seckill/common/interfaces"
	dbrpc "seckill/dbservice/rpc"
	"strings"
	"sync"
)

type OrderServImpl struct {
	dbCli *dbrpc.OrderRpcServCli
	cacheCli *cacherpc.CacheServCli
}

func (o *OrderServImpl) AddOrder(order *entity.Order, method string) error {
	if strings.Contains(method, "CACHE") {
		return o.cacheCli.CreateOrder(order, method)
	}
	return o.dbCli.AddOrder(order,method)
}

func (o *OrderServImpl) GetOrderById(id uint) (entity.Order, error) {
	return o.dbCli.GetOrderById(id)
}

func (o *OrderServImpl) GetOrdersByUID(uid uint) ([]entity.Order, error) {
	return o.dbCli.GetOrdersByUID(uid)
}

func (o *OrderServImpl) GetOrdersByPID(pid uint) ([]entity.Order, error) {
	return o.dbCli.GetOrdersByPID(pid)
}

func (o *OrderServImpl) DeleteOrder(id uint) error {
	return o.dbCli.DeleteOrder(id)
}

func (o *OrderServImpl) GetOrders() ([]entity.Order, error) {
	return o.dbCli.GetOrders()
}

func (o *OrderServImpl) ClearOrders() error {
	return o.dbCli.ClearOrders()
}

var _ interfaces.OrderServ = (*OrderServImpl)(nil)

var orderServOnce = new(sync.Once)

var orderServ *OrderServImpl

func GetOrderService() interfaces.OrderServ {
	dbServCli, err := dbrpc.GetDbServRpcCli()
	if err != nil {
		return (*OrderServImpl)(nil)
	}
	cacheSevCli, err := cacherpc.NewCacheServClient()
	if err != nil {
		return (*OrderServImpl)(nil)
	}
	orderServOnce.Do(func() {
		orderServCli := dbServCli.GetOrderRpcServCli()
		orderServ = &OrderServImpl{
			dbCli: &orderServCli,
			cacheCli: cacheSevCli,
		}
	})
	return orderServ
}
