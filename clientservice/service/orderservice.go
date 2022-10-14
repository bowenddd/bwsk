package service

import (
	"seckill/common/entity"
	"seckill/common/interfaces"
	"seckill/dbservice/rpc"
	"sync"
)

type OrderServImpl struct {
	cli *rpc.OrderRpcServCli
}

func (o *OrderServImpl) AddOrder(order *entity.Order, method string) error {
	return o.cli.AddOrder(order,method)
}

func (o *OrderServImpl) GetOrderById(id uint) (entity.Order, error) {
	return o.cli.GetOrderById(id)
}

func (o *OrderServImpl) GetOrdersByUID(uid uint) ([]entity.Order, error) {
	return o.cli.GetOrdersByUID(uid)
}

func (o *OrderServImpl) GetOrdersByPID(pid uint) ([]entity.Order, error) {
	return o.cli.GetOrdersByPID(pid)
}

func (o *OrderServImpl) DeleteOrder(id uint) error {
	return o.cli.DeleteOrder(id)
}

func (o *OrderServImpl) GetOrders() ([]entity.Order, error) {
	return o.cli.GetOrders()
}

func (o *OrderServImpl) ClearOrders() error {
	return o.cli.ClearOrders()
}

var _ interfaces.OrderServ = (*OrderServImpl)(nil)

var orderServOnce = new(sync.Once)

var orderServ *OrderServImpl

func GetOrderService() interfaces.OrderServ {
	cli, err := rpc.GetDbServRpcCli()
	if err != nil {
		return (*OrderServImpl)(nil)
	}
	orderServOnce.Do(func() {
		orderServCli := cli.GetOrderRpcServCli()
		orderServ = &OrderServImpl{
			cli: &orderServCli,
		}
	})
	return orderServ
}
