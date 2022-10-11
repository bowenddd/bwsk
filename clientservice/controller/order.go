package controller

import (
	"fmt"
	"net/http"
	"seckill/clientservice/service"
	"seckill/common/entity"
	"seckill/common/interfaces"
	"seckill/common/response"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	serv interfaces.OrderServ
}

func (o *OrderController) Create(ctx *gin.Context) {
	order := &entity.Order{}
	err := ctx.BindJSON(order)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	productServ := service.GetProductService()
	stock, err := productServ.GetStock(order.ProductId)
	if err != nil {
		response.Error(ctx, http.StatusBadGateway, err.Error())
		return
	}
	if order.Num > stock {
		response.Success(ctx, http.StatusOK, "", "库存不足")
		return
	}
	err = o.serv.AddOrder(order)
	if err != nil {
		response.Error(ctx, http.StatusBadGateway, err.Error())
		return
	}
	response.Success(ctx, http.StatusOK, "", fmt.Sprintf("create order(id%d) success", order.ID))
}

func (o *OrderController) GetById(ctx *gin.Context) {
	id_str := ctx.Param("id")
	id, err := strconv.Atoi(id_str)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "order id must be integer")
		return
	}
	order, err := o.serv.GetOrderById(uint(id))
	if err != nil {
		response.Error(ctx, http.StatusBadGateway, err.Error())
		return
	}
	response.Success(ctx, http.StatusOK, order, "ok")

}

func (o *OrderController) GetByUID(ctx *gin.Context) {
	uid_str := ctx.Param("uid")
	uid, err := strconv.Atoi(uid_str)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "user id must be integer")
		return
	}
	orders, err := o.serv.GetOrdersByUID(uint(uid))
	if err != nil {
		response.Error(ctx, http.StatusBadGateway, err.Error())
		return
	}
	response.Success(ctx, http.StatusOK, orders, "ok")
}

func (o *OrderController) GetByPID(ctx *gin.Context) {
	pid_str := ctx.Param("pid")
	pid, err := strconv.Atoi(pid_str)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "product id must be integer")
		return
	}
	orders, err := o.serv.GetOrdersByPID(uint(pid))
	if err != nil {
		response.Error(ctx, http.StatusBadGateway, err.Error())
		return
	}
	response.Success(ctx, http.StatusOK, orders, "ok")
}
func (o *OrderController) Delete(ctx *gin.Context) {
	id_str := ctx.Param("id")
	id, err := strconv.Atoi(id_str)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "order id must be integer")
		return
	}
	err = o.serv.DeleteOrder(uint(id))
	if err != nil {
		response.Error(ctx, http.StatusBadGateway, err.Error())
		return
	}
	response.Success(ctx, http.StatusOK, "", fmt.Sprintf("delete order(id:%s) success", id_str))
}
func (o *OrderController) List(ctx *gin.Context) {
	orders, err := o.serv.GetOrders()
	if err != nil {
		response.Error(ctx, http.StatusBadGateway, err.Error())
		return
	}
	response.Success(ctx, http.StatusOK, orders, "ok")
}

var orderCliOnce = new(sync.Once)

var orderController *OrderController

func GetOrderController() *OrderController {
	orderCliOnce.Do(func() {
		orderController = &OrderController{
			serv: service.GetOrderService(),
		}
	})
	return orderController
}
