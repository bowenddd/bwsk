package controller

import (
	"fmt"
	"net/http"
	"seckill/common/entity"
	"seckill/common/interfaces"
	"seckill/common/response"
	"seckill/clientservice/service"
	"sync"
	"strconv"
	"github.com/gin-gonic/gin"
)


type ProductController struct{
	serv interfaces.ProductServ
}

func (p *ProductController) Get(ctx *gin.Context){
	name := ctx.Param("name")
	product, err := p.serv.GetProduct(name)
	if err != nil{
		response.Error(ctx,http.StatusBadGateway,err.Error())
		return
	}
	response.Success(ctx,http.StatusOK,product,"ok")
}

func (p *ProductController)Create(ctx *gin.Context){
	product := &entity.Product{}
	err := ctx.BindJSON(product)
	if err != nil{
		response.Error(ctx,http.StatusBadRequest,err.Error())
		return
	}
	err = p.serv.AddProduct(product)
	if err != nil{
		response.Error(ctx, http.StatusBadGateway,err.Error())
		return
	}
	response.Success(ctx,http.StatusOK,"",fmt.Sprintf("create product %s success",product.Name))
}

func (p *ProductController)List(ctx *gin.Context){
	products, err := p.serv.GetProducts()
	if err != nil{
		response.Error(ctx,http.StatusBadGateway,err.Error())
		return
	}
	response.Success(ctx,http.StatusOK,products,"ok")
}

func (p *ProductController)Delete(ctx *gin.Context){
	name := ctx.Param("name")
	err := p.serv.DeleteProduct(name)
	if err != nil{
		response.Error(ctx, http.StatusBadGateway, err.Error())
		return
	}
	response.Success(ctx,http.StatusOK,"",fmt.Sprintf("delete product %s success",name))
}

func (p *ProductController)SetStock(ctx *gin.Context){
	id_str := ctx.Query("id")
	num_str := ctx.Query("num")
	if id_str == "" || num_str == ""{
		response.Error(ctx,http.StatusBadRequest,"id and num must not empty")
		return
	}
	num, err := strconv.Atoi(num_str)
	if err != nil{
		response.Error(ctx,http.StatusBadRequest,"num must be an integer")
		return
	}
	id, err := strconv.Atoi(id_str)
	if err != nil{
		response.Error(ctx,http.StatusBadRequest,"id must be an unsigned integer")
		return
	}
	err = p.serv.SetStock(uint(id), num)
	if err != nil{
		response.Error(ctx,http.StatusBadGateway,err.Error())
		return
	}
	response.Success(ctx,http.StatusOK,"",fmt.Sprintf("product %s 's stock is set to %s",id_str,num_str))
}

func (p *ProductController)GetStock(ctx *gin.Context){
	id_str := ctx.Param("id")
	id, err := strconv.Atoi(id_str)
	if err != nil{
		response.Error(ctx,http.StatusBadRequest,"id must be an unsigned integer")
		return
	}
	stock, err := p.serv.GetStock(uint(id))
	if err != nil{
		response.Error(ctx,http.StatusBadGateway,err.Error())
		return
	}
	response.Success(ctx,http.StatusOK,stock,"ok")
}


var productCtlOnce = new(sync.Once)

var productController *ProductController


func GetProductController() *ProductController{
	productCtlOnce.Do(func() {
		productController = & ProductController{
			serv: service.GetProductService(),
		}	
	})
	return productController
}